package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/service"
)

// Server represents the admin API HTTP server
type Server struct {
	config     *config.APIConfig
	httpServer *http.Server
	logger     *zap.Logger
	router     *Router
}

// NewServer creates a new admin API server
// This is a wrapper around NewRouter for backward compatibility
func NewServer(
	cfg *config.APIConfig,
	fullConfig *config.Config,
	configPath string,
	db *database.DB,
	domainRepo repository.DomainRepository,
	userRepo repository.UserRepository,
	aliasRepo repository.AliasRepository,
	mailboxRepo repository.MailboxRepository,
	messageRepo repository.MessageRepository,
	queueRepo repository.QueueRepository,
	apiKeyRepo repository.APIKeyRepository,
	rateLimitRepo repository.RateLimitRepository,
	logger *zap.Logger,
) *Server {
	// Create services
	domainService := service.NewDomainService(domainRepo)
	userService := service.NewUserService(userRepo, domainRepo, logger)
	aliasService := service.NewAliasService(aliasRepo)
	mailboxService := service.NewMailboxService(mailboxRepo, logger)
	messageService := service.NewMessageService(messageRepo, "./data/mail", logger)
	queueService := service.NewQueueService(queueRepo, logger)
	setupService := service.NewSetupService(db, userRepo, domainRepo, logger)
	settingsService := service.NewSettingsService(fullConfig, configPath, logger)
	pgpService := service.NewPGPService(db, logger)
	auditService := service.NewAuditService(db, logger)

	// Wire up cross-service dependencies for webmail
	messageService.SetQueueService(queueService)
	messageService.SetMailboxService(mailboxService)

	// Create router with all dependencies
	router := NewRouter(RouterConfig{
		Logger:          logger,
		DomainService:   domainService,
		UserService:     userService,
		AliasService:    aliasService,
		MailboxService:  mailboxService,
		MessageService:  messageService,
		QueueService:    queueService,
		SetupService:    setupService,
		SettingsService: settingsService,
		PGPService:      pgpService,
		AuditService:    auditService,
		APIKeyRepo:      apiKeyRepo,
		RateLimitRepo:   rateLimitRepo,
		DB:              db.DB,
		JWTSecret:       cfg.JWTSecret,
		CORSOrigins:     cfg.CORSOrigins,
	})

	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	return &Server{
		config:     cfg,
		httpServer: httpServer,
		logger:     logger,
		router:     router,
	}
}

// Start starts the admin API server
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("starting admin API server",
		zap.Int("port", s.config.Port),
	)

	// Start server in goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("admin API server error", zap.Error(err))
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the admin API server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down admin API server")
	return s.httpServer.Shutdown(ctx)
}
