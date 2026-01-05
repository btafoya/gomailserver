package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	calendarService "github.com/btafoya/gomailserver/internal/calendar/service"
	"github.com/btafoya/gomailserver/internal/config"
	contactService "github.com/btafoya/gomailserver/internal/contact/service"
	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/repository"
	repRepository "github.com/btafoya/gomailserver/internal/reputation/repository"
	repService "github.com/btafoya/gomailserver/internal/reputation/service"
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
	webhookRepo repository.WebhookRepository,
	contactService *contactService.ContactService,
	addressbookService *contactService.AddressbookService,
	calendarService *calendarService.CalendarService,
	eventService *calendarService.EventService,
	telemetryService *repService.TelemetryService,
	auditorService *repService.AuditorService,
	reputationDB interface {
		GetEventRepo() repRepository.EventsRepository
		GetScoresRepo() repRepository.ScoresRepository
		GetCircuitBreakerRepo() repRepository.CircuitBreakerRepository
	},
	logger *zap.Logger,
) *Server {
	// Create services
	domainService := service.NewDomainService(domainRepo)
	userService := service.NewUserService(userRepo, domainRepo, logger)
	aliasService := service.NewAliasService(aliasRepo)
	mailboxService := service.NewMailboxService(mailboxRepo, logger)
	messageService := service.NewMessageService(messageRepo, "./data/mail", logger)
	queueService := service.NewQueueService(queueRepo, telemetryService, logger)
	setupService := service.NewSetupService(db, userRepo, domainRepo, logger)
	settingsService := service.NewSettingsService(fullConfig, configPath, logger)
	pgpService := service.NewPGPService(db, logger)
	auditService := service.NewAuditService(db, logger)
	webhookService := service.NewWebhookService(webhookRepo, logger)

	// Wire up cross-service dependencies for webmail
	messageService.SetQueueService(queueService)
	messageService.SetMailboxService(mailboxService)

	// Create router with all dependencies
	router := NewRouter(RouterConfig{
		Logger:             logger,
		DomainService:      domainService,
		UserService:        userService,
		AliasService:       aliasService,
		MailboxService:     mailboxService,
		MessageService:     messageService,
		QueueService:       queueService,
		SetupService:       setupService,
		SettingsService:    settingsService,
		PGPService:         pgpService,
		AuditService:       auditService,
		WebhookService:     webhookService,
		ContactService:     contactService,
		AddressbookService: addressbookService,
		CalendarService:    calendarService,
		EventService:       eventService,
		AuditorService:     auditorService,
		ScoresRepo:         reputationDB.GetScoresRepo(),
		EventsRepo:         reputationDB.GetEventRepo(),
		CircuitRepo:        reputationDB.GetCircuitBreakerRepo(),
		APIKeyRepo:         apiKeyRepo,
		RateLimitRepo:      rateLimitRepo,
		DB:                 db.DB,
		JWTSecret:          cfg.JWTSecret,
		CORSOrigins:        cfg.CORSOrigins,
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