package database

import (
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/config"
	postgresdb "github.com/btafoya/gomailserver/internal/database/postgres"
	"go.uber.org/zap"
)

// Factory creates a database connection based on driver configuration
func Factory(cfg *config.Config, logger *zap.Logger) (*DB, error) {
	driver := Driver(cfg.Database.Driver)

	switch driver {
	case DriverPostgres:
		return newPostgresDB(cfg, logger)
	case DriverSQLite:
		return newSQLiteDB(cfg, logger)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}
}

func newPostgresDB(cfg *config.Config, logger *zap.Logger) (*DB, error) {
	pgCfg := postgresdb.Config{
		Host:            cfg.Database.Postgres.Host,
		Port:            cfg.Database.Postgres.Port,
		Database:        cfg.Database.Postgres.Database,
		User:            cfg.Database.Postgres.User,
		Password:        cfg.Database.Postgres.Password,
		SSLMode:         cfg.Database.Postgres.SSLMode,
		MaxOpenConns:    cfg.Database.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Database.Postgres.MaxIdleConns,
		ConnMaxLifetime: parseDuration(cfg.Database.Postgres.ConnMaxLifetime, time.Hour),
		ConnMaxIdleTime: parseDuration(cfg.Database.Postgres.ConnMaxIdleTime, 30*time.Minute),
	}

	pgDB, err := postgresdb.New(pgCfg, logger)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB:     pgDB.DB,
		logger: logger,
		driver: DriverPostgres,
	}, nil
}

func newSQLiteDB(cfg *config.Config, logger *zap.Logger) (*DB, error) {
	sqliteCfg := Config{
		Path:       cfg.Database.SQLite.Path,
		WALEnabled: cfg.Database.SQLite.WALEnabled,
	}

	sqliteDB, err := New(sqliteCfg, logger)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB:     sqliteDB.DB,
		logger: logger,
		driver: DriverSQLite,
	}, nil
}

func (db *DB) DriverType() Driver {
	return db.driver
}

func parseDuration(s string, defaultDuration time.Duration) time.Duration {
	if s == "" {
		return defaultDuration
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultDuration
	}

	return d
}
