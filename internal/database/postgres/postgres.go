package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// DB wraps SQL database connection
type DB struct {
	*sql.DB
	logger *zap.Logger
}

// Config holds database configuration
type Config struct {
	Host            string
	Port            int
	Database        string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// New creates a new PostgreSQL database connection
func New(cfg Config, logger *zap.Logger) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password, cfg.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	wrapper := &DB{
		DB:     db,
		logger: logger,
	}

	logger.Info("PostgreSQL connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database),
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
	)

	return wrapper, nil
}

// Close closes database connection
func (db *DB) Close() error {
	db.logger.Info("closing PostgreSQL connection")
	return db.DB.Close()
}

// Vacuum performs database maintenance (PostgreSQL equivalent of SQLite VACUUM)
func (db *DB) Vacuum() error {
	db.logger.Info("running VACUUM ANALYZE")
	_, err := db.Exec("VACUUM ANALYZE")
	return err
}

// Analyze updates database statistics
func (db *DB) Analyze() error {
	db.logger.Info("running ANALYZE")
	_, err := db.Exec("ANALYZE")
	return err
}
