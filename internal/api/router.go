package api

import (
	"net/http"
	"time"

	"github.com/btafoya/gomailserver/internal/api/handlers"
	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Router configuration for the REST API
type Router struct {
	*chi.Mux
	logger         *zap.Logger
	domainService  *service.DomainService
	userService    *service.UserService
	aliasService   *service.AliasService
	mailboxService *service.MailboxService
	messageService *service.MessageService
	queueService   *service.QueueService
	apiKeyRepo     repository.APIKeyRepository
	jwtSecret      string
}

// RouterConfig contains dependencies for the API router
type RouterConfig struct {
	Logger         *zap.Logger
	DomainService  *service.DomainService
	UserService    *service.UserService
	AliasService   *service.AliasService
	MailboxService *service.MailboxService
	MessageService *service.MessageService
	QueueService   *service.QueueService
	APIKeyRepo     repository.APIKeyRepository
	JWTSecret      string
	CORSOrigins    []string
}

// NewRouter creates a new API router with all routes configured
func NewRouter(config RouterConfig) *Router {
	r := &Router{
		Mux:            chi.NewRouter(),
		logger:         config.Logger,
		domainService:  config.DomainService,
		userService:    config.UserService,
		aliasService:   config.AliasService,
		mailboxService: config.MailboxService,
		messageService: config.MessageService,
		queueService:   config.QueueService,
		apiKeyRepo:     config.APIKeyRepo,
		jwtSecret:      config.JWTSecret,
	}

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger(config.Logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		middleware.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Authentication routes (no auth required)
		r.Post("/auth/login", handlers.NewAuthHandler(
			config.UserService,
			config.JWTSecret,
			config.Logger,
		).Login)

		r.Post("/auth/refresh", handlers.NewAuthHandler(
			config.UserService,
			config.JWTSecret,
			config.Logger,
		).Refresh)

		// Protected routes
		r.Group(func(r chi.Router) {
			// JWT or API Key authentication required
			r.Use(middleware.Auth(config.JWTSecret, config.APIKeyRepo, config.Logger))

			// Domain management
			domainHandler := handlers.NewDomainHandler(config.DomainService, config.Logger)
			r.Route("/domains", func(r chi.Router) {
				r.Get("/", domainHandler.List)
				r.Post("/", domainHandler.Create)
				r.Get("/{id}", domainHandler.Get)
				r.Put("/{id}", domainHandler.Update)
				r.Delete("/{id}", domainHandler.Delete)
				r.Post("/{id}/dkim", domainHandler.GenerateDKIM)
			})

			// User management
			userHandler := handlers.NewUserHandler(config.UserService, config.Logger)
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.List)
				r.Post("/", userHandler.Create)
				r.Get("/{id}", userHandler.Get)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
				r.Post("/{id}/password", userHandler.ResetPassword)
			})

			// Alias management
			aliasHandler := handlers.NewAliasHandler(config.AliasService, config.Logger)
			r.Route("/aliases", func(r chi.Router) {
				r.Get("/", aliasHandler.List)
				r.Post("/", aliasHandler.Create)
				r.Get("/{id}", aliasHandler.Get)
				r.Delete("/{id}", aliasHandler.Delete)
			})

			// Statistics and monitoring
			statsHandler := handlers.NewStatsHandler(
				config.DomainService,
				config.UserService,
				config.QueueService,
				config.Logger,
			)
			r.Route("/stats", func(r chi.Router) {
				r.Get("/dashboard", statsHandler.Dashboard)
				r.Get("/domains/{id}", statsHandler.Domain)
				r.Get("/users/{id}", statsHandler.User)
			})

			// Queue management
			queueHandler := handlers.NewQueueHandler(config.QueueService, config.Logger)
			r.Route("/queue", func(r chi.Router) {
				r.Get("/", queueHandler.List)
				r.Get("/{id}", queueHandler.Get)
				r.Post("/{id}/retry", queueHandler.Retry)
				r.Delete("/{id}", queueHandler.Delete)
			})

			// Log retrieval
			logHandler := handlers.NewLogHandler(config.Logger)
			r.Get("/logs", logHandler.List)
		})
	})

	return r
}
