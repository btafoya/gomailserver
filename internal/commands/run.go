package commands

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/imap"
	"github.com/btafoya/gomailserver/internal/repository/sqlite"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/btafoya/gomailserver/internal/smtp"
	tlspkg "github.com/btafoya/gomailserver/internal/tls"
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
	defer func() {
		_ = logger.Sync() // Ignore sync errors on defer
	}()

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

	// Initialize TLS manager
	var tlsMgr *tlspkg.Manager
	var tlsCfg *tls.Config
	if cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
		tlsMgr, err = tlspkg.NewManager(&cfg.TLS, cfg.Server.Hostname, logger)
		if err != nil {
			return fmt.Errorf("failed to initialize TLS manager: %w", err)
		}
		tlsCfg = tlsMgr.GetTLSConfig()
		logger.Info("TLS initialized",
			zap.String("hostname", cfg.Server.Hostname),
		)

		// Check certificate expiry
		if err := tlsMgr.ValidateExpiry(30); err != nil {
			logger.Warn("TLS certificate validation warning", zap.Error(err))
		}
	} else {
		logger.Warn("TLS disabled - NOT RECOMMENDED FOR PRODUCTION")
	}

	// Create repositories
	userRepo := sqlite.NewUserRepository(db)
	mailboxRepo := sqlite.NewMailboxRepository(db)
	messageRepo := sqlite.NewMessageRepository(db)
	queueRepo := sqlite.NewQueueRepository(db)

	logger.Debug("repositories initialized")

	// Create services
	userSvc := service.NewUserService(userRepo, logger)
	mailboxSvc := service.NewMailboxService(mailboxRepo, logger)
	messageSvc := service.NewMessageService(messageRepo, "./data/mail", logger)
	queueSvc := service.NewQueueService(queueRepo, logger)

	logger.Debug("services initialized")

	// Create SMTP server
	smtpServer := smtp.NewServer(&cfg.SMTP, tlsCfg, userSvc, messageSvc, queueSvc, logger)

	// Create IMAP server
	imapServer := imap.NewServer(&cfg.IMAP, tlsCfg, userSvc, mailboxSvc, messageSvc, logger)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start SMTP server
	if err := smtpServer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start SMTP server: %w", err)
	}

	// Start IMAP server
	if err := imapServer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start IMAP server: %w", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("mail server ready",
		zap.Int("smtp_submission_port", cfg.SMTP.SubmissionPort),
		zap.Int("smtp_relay_port", cfg.SMTP.RelayPort),
		zap.Int("smtps_port", cfg.SMTP.SMTPSPort),
		zap.Int("imap_port", cfg.IMAP.Port),
		zap.Int("imaps_port", cfg.IMAP.IMAPSPort),
	)

	// Wait for shutdown signal
	sig := <-sigChan
	logger.Info("received shutdown signal",
		zap.String("signal", sig.String()),
	)

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	logger.Info("initiating graceful shutdown...")

	// Shutdown SMTP server
	if err := smtpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("SMTP server shutdown error", zap.Error(err))
	}

	// Shutdown IMAP server
	if err := imapServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("IMAP server shutdown error", zap.Error(err))
	}

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
