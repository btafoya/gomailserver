package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

// GetMigrations returns all available migrations
func GetMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "Initial schema - create all tables",
			Up:          migrationV1Up,
			Down:        migrationV1Down,
		},
		{
			Version:     2,
			Description: "Add security configuration columns to domains table",
			Up:          migrationV2Up,
			Down:        migrationV2Down,
		},
		{
			Version:     3,
			Description: "Add API keys, TLS certificates, and setup wizard tables",
			Up:          migrationV3Up,
			Down:        migrationV3Down,
		},
		{
			Version:     4,
			Description: "Add role column to users table for admin/user distinction",
			Up:          migrationV4Up,
			Down:        migrationV4Down,
		},
	}
}

// Migrate runs all pending migrations
func (db *DB) Migrate() error {
	db.logger.Info("starting database migration")

	// Create migrations table if it doesn't exist
	if err := db.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current version
	currentVersion, err := db.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	db.logger.Info("current database version",
		zap.Int("version", currentVersion),
	)

	// Run pending migrations
	migrations := GetMigrations()
	for _, migration := range migrations {
		if migration.Version <= currentVersion {
			continue
		}

		db.logger.Info("applying migration",
			zap.Int("version", migration.Version),
			zap.String("description", migration.Description),
		)

		if err := db.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}

		db.logger.Info("migration applied successfully",
			zap.Int("version", migration.Version),
		)
	}

	db.logger.Info("database migration complete")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (db *DB) createMigrationsTable() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// getCurrentVersion gets the current schema version
func (db *DB) getCurrentVersion() (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

// applyMigration applies a single migration
func (db *DB) applyMigration(migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			db.logger.Debug("transaction rollback after commit", zap.Error(err))
		}
	}()

	// Execute migration statements one by one
	statements := splitSQL(migration.Up)
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.ExecContext(context.Background(), stmt); err != nil {
			return fmt.Errorf("failed to execute migration statement %d: %w\nStatement: %s", i+1, err, stmt)
		}
	}

	// Record migration
	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (version, description, applied_at) VALUES (?, ?, ?)",
		migration.Version,
		migration.Description,
		time.Now(),
	); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Rollback rolls back to a specific version
func (db *DB) Rollback(targetVersion int) error {
	db.logger.Info("rolling back database",
		zap.Int("target_version", targetVersion),
	)

	currentVersion, err := db.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if targetVersion >= currentVersion {
		return fmt.Errorf("target version %d is not less than current version %d", targetVersion, currentVersion)
	}

	migrations := GetMigrations()
	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]
		if migration.Version <= targetVersion {
			break
		}

		if migration.Version > currentVersion {
			continue
		}

		db.logger.Info("rolling back migration",
			zap.Int("version", migration.Version),
		)

		if err := db.rollbackMigration(migration); err != nil {
			return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
		}
	}

	db.logger.Info("database rollback complete")
	return nil
}

// rollbackMigration rolls back a single migration
func (db *DB) rollbackMigration(migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			db.logger.Debug("transaction rollback after commit", zap.Error(err))
		}
	}()

	// Execute rollback statements one by one
	statements := splitSQL(migration.Down)
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.ExecContext(context.Background(), stmt); err != nil {
			return fmt.Errorf("failed to execute rollback statement %d: %w", i+1, err)
		}
	}

	// Remove migration record
	if _, err := tx.Exec(
		"DELETE FROM schema_migrations WHERE version = ?",
		migration.Version,
	); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// splitSQL splits a SQL migration into individual statements
func splitSQL(sql string) []string {
	var statements []string
	var current strings.Builder

	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		current.WriteString(line)
		current.WriteString("\n")

		// Check if this line ends a statement
		if strings.HasSuffix(line, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	// Add any remaining SQL
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}
