package reputation

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"github.com/btafoya/gomailserver/internal/reputation/repository/sqlite"
	"github.com/btafoya/gomailserver/internal/reputation/service"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const migrationReputationV1Up = `
-- Reputation metrics database (separate SQLite: reputation.db)

-- Sending events table
CREATE TABLE sending_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    domain TEXT NOT NULL,
    recipient_domain TEXT NOT NULL,
    event_type TEXT NOT NULL,
    bounce_type TEXT,
    enhanced_status_code TEXT,
    smtp_response TEXT,
    ip_address TEXT NOT NULL,
    metadata TEXT
);

CREATE INDEX idx_sending_events_timestamp ON sending_events(timestamp);
CREATE INDEX idx_sending_events_domain ON sending_events(domain);
CREATE INDEX idx_sending_events_event_type ON sending_events(event_type);
CREATE INDEX idx_sending_events_recipient_domain ON sending_events(recipient_domain);

-- Domain reputation scores table
CREATE TABLE domain_reputation_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL UNIQUE,
    reputation_score INTEGER NOT NULL,
    complaint_rate REAL NOT NULL,
    bounce_rate REAL NOT NULL,
    delivery_rate REAL NOT NULL,
    circuit_breaker_active BOOLEAN DEFAULT 0,
    circuit_breaker_reason TEXT,
    warm_up_active BOOLEAN DEFAULT 0,
    warm_up_day INTEGER DEFAULT 0,
    last_updated INTEGER NOT NULL
);

CREATE INDEX idx_domain_reputation_domain ON domain_reputation_scores(domain);
CREATE INDEX idx_domain_reputation_score ON domain_reputation_scores(reputation_score);

-- Warm-up schedules table
CREATE TABLE warm_up_schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    day INTEGER NOT NULL,
    max_volume INTEGER NOT NULL,
    actual_volume INTEGER DEFAULT 0,
    created_at INTEGER NOT NULL
);

CREATE INDEX idx_warm_up_domain ON warm_up_schedules(domain);
CREATE INDEX idx_warm_up_day ON warm_up_schedules(domain, day);

-- Circuit breaker events table
CREATE TABLE circuit_breaker_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    trigger_type TEXT NOT NULL,
    trigger_value REAL NOT NULL,
    threshold REAL NOT NULL,
    paused_at INTEGER NOT NULL,
    resumed_at INTEGER,
    auto_resumed BOOLEAN DEFAULT 0,
    admin_notes TEXT
);

CREATE INDEX idx_circuit_breaker_domain ON circuit_breaker_events(domain);
CREATE INDEX idx_circuit_breaker_active ON circuit_breaker_events(domain, resumed_at);

-- Retention policy table
CREATE TABLE retention_policy (
    id INTEGER PRIMARY KEY,
    retention_days INTEGER NOT NULL DEFAULT 90,
    last_cleanup INTEGER NOT NULL
);

INSERT INTO retention_policy (id, retention_days, last_cleanup)
VALUES (1, 90, strftime('%s', 'now'));
`

// Database represents the reputation database connection and repositories
type Database struct {
	DB                       *sql.DB
	EventsRepo               repository.EventsRepository
	ScoresRepo               repository.ScoresRepository
	WarmUpRepo               repository.WarmUpRepository
	CircuitBreakerRepo       repository.CircuitBreakerRepository
	TelemetryService         *service.TelemetryService
	logger                   *zap.Logger
}

// Config holds reputation database configuration
type Config struct {
	Path string
}

// InitDatabase initializes the reputation database and creates all tables
func InitDatabase(cfg Config, logger *zap.Logger) (*Database, error) {
	// Ensure the database directory exists
	dbDir := filepath.Dir(cfg.Path)
	if dbDir != "." && dbDir != "/" {
		if err := os.MkdirAll(dbDir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create reputation database directory: %w", err)
		}
	}

	// Open database connection
	dsn := fmt.Sprintf("%s?_foreign_keys=1", cfg.Path)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open reputation database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping reputation database: %w", err)
	}

	// Apply PRAGMA settings
	if err := applyPragmas(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply PRAGMAs: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize repositories
	eventsRepo := sqlite.NewEventsRepository(db)
	scoresRepo := sqlite.NewScoresRepository(db)
	warmUpRepo := sqlite.NewWarmUpRepository(db)
	circuitBreakerRepo := sqlite.NewCircuitBreakerRepository(db)

	// Initialize services
	telemetryService := service.NewTelemetryService(eventsRepo, scoresRepo, logger)

	logger.Info("reputation database initialized",
		zap.String("path", cfg.Path),
	)

	return &Database{
		DB:                 db,
		EventsRepo:         eventsRepo,
		ScoresRepo:         scoresRepo,
		WarmUpRepo:         warmUpRepo,
		CircuitBreakerRepo: circuitBreakerRepo,
		TelemetryService:   telemetryService,
		logger:             logger,
	}, nil
}

// Close closes the reputation database connection
func (d *Database) Close() error {
	d.logger.Info("closing reputation database connection")
	return d.DB.Close()
}

// applyPragmas configures SQLite for optimal performance
func applyPragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = -32000", // 32MB cache
		"PRAGMA temp_store = MEMORY",
		"PRAGMA mmap_size = 134217728", // 128MB mmap
		"PRAGMA page_size = 4096",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %q: %w", pragma, err)
		}
	}

	return nil
}

// runMigrations applies database schema migrations
func runMigrations(db *sql.DB) error {
	ctx := context.Background()

	// Create migrations table if it doesn't exist
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at INTEGER NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Check if migration v1 has been applied
	var version int
	err = db.QueryRowContext(ctx, "SELECT version FROM schema_migrations WHERE version = 1").Scan(&version)
	if err == sql.ErrNoRows {
		// Apply migration v1
		if _, err := db.ExecContext(ctx, migrationReputationV1Up); err != nil {
			return fmt.Errorf("failed to apply migration v1: %w", err)
		}

		// Record migration
		_, err = db.ExecContext(ctx, "INSERT INTO schema_migrations (version, applied_at) VALUES (1, ?)",
			time.Now().Unix())
		if err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	return nil
}
