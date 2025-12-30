package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// DB wraps the SQL database connection
type DB struct {
	*sql.DB
	logger *zap.Logger
}

// Config holds database configuration
type Config struct {
	Path       string
	WALEnabled bool
}

// New creates a new database connection
func New(cfg Config, logger *zap.Logger) (*DB, error) {
	// Ensure directory exists
	dbDir := filepath.Dir(cfg.Path)
	if dbDir != "." && dbDir != "/" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open database connection
	db, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	wrapper := &DB{
		DB:     db,
		logger: logger,
	}

	// Apply PRAGMA settings
	if err := wrapper.applyPragmas(cfg.WALEnabled); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply PRAGMAs: %w", err)
	}

	logger.Info("database connection established",
		zap.String("path", cfg.Path),
		zap.Bool("wal_enabled", cfg.WALEnabled),
	)

	return wrapper, nil
}

// applyPragmas configures SQLite for optimal performance
func (db *DB) applyPragmas(walEnabled bool) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = -64000", // 64MB cache
		"PRAGMA temp_store = MEMORY",
		"PRAGMA mmap_size = 268435456", // 256MB mmap
		"PRAGMA page_size = 4096",
	}

	if !walEnabled {
		pragmas[1] = "PRAGMA journal_mode = DELETE"
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %q: %w", pragma, err)
		}
	}

	db.logger.Debug("applied database PRAGMAs",
		zap.Bool("wal_enabled", walEnabled),
	)

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	db.logger.Info("closing database connection")
	return db.DB.Close()
}

// Vacuum performs database maintenance
func (db *DB) Vacuum() error {
	db.logger.Info("running VACUUM")
	_, err := db.Exec("VACUUM")
	return err
}

// Analyze updates database statistics
func (db *DB) Analyze() error {
	db.logger.Info("running ANALYZE")
	_, err := db.Exec("ANALYZE")
	return err
}
