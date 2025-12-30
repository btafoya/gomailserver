package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/database"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the mail server",
	Long:  "Start the gomailserver mail server with the specified configuration",
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	logger, err := config.NewLogger(cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

	logger.Info("starting gomailserver",
		zap.String("version", version),
		zap.String("hostname", cfg.Server.Hostname),
		zap.String("domain", cfg.Server.Domain),
	)

	// Initialize database
	db, err := initDatabase(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start services (placeholder for now)
	logger.Info("mail server ready",
		zap.Int("smtp_submission_port", cfg.SMTP.SubmissionPort),
		zap.Int("smtp_relay_port", cfg.SMTP.RelayPort),
		zap.Int("imap_port", cfg.IMAP.Port),
	)

	// Wait for shutdown signal
	sig := <-sigChan
	logger.Info("received shutdown signal",
		zap.String("signal", sig.String()),
	)

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	logger.Info("initiating graceful shutdown...")

	// Shutdown services (placeholder for now)
	<-shutdownCtx.Done()

	logger.Info("shutdown complete")
	return nil
}

func initDatabase(cfg *config.Config, logger *zap.Logger) (*database.DB, error) {
	dbConfig := database.Config{
		Path:       cfg.Database.Path,
		WALEnabled: cfg.Database.WALEnabled,
	}

	db, err := database.New(dbConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return db, nil
}
