package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/btafoya/gomailserver/internal/admin"
	"github.com/btafoya/gomailserver/internal/api/handlers"
	"github.com/btafoya/gomailserver/internal/api/middleware"
	calendarService "github.com/btafoya/gomailserver/internal/calendar/service"
	contactService "github.com/btafoya/gomailserver/internal/contact/service"
	"github.com/btafoya/gomailserver/internal/postmark"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/btafoya/gomailserver/internal/webmail"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Router configuration for the REST API
type Router struct {
	*chi.Mux
	logger          *zap.Logger
	domainService   *service.DomainService
	userService     *service.UserService
	aliasService    *service.AliasService
	mailboxService  *service.MailboxService
	messageService  *service.MessageService
	queueService    *service.QueueService
	settingsService *service.SettingsService
	apiKeyRepo      repository.APIKeyRepository
	jwtSecret       string
}

// RouterConfig contains dependencies for the API router
type RouterConfig struct {
	Logger             *zap.Logger
	DomainService      *service.DomainService
	UserService        *service.UserService
	AliasService       *service.AliasService
	MailboxService     *service.MailboxService
	MessageService     *service.MessageService
	QueueService       *service.QueueService
	SetupService       *service.SetupService
	SettingsService    *service.SettingsService
	PGPService         *service.PGPService
	AuditService       *service.AuditService
	WebhookService     *service.WebhookService
	ContactService     *contactService.ContactService
	AddressbookService *contactService.AddressbookService
	CalendarService    *calendarService.CalendarService
	EventService       *calendarService.EventService
	APIKeyRepo         repository.APIKeyRepository
	RateLimitRepo      repository.RateLimitRepository
	DB                 *sql.DB
	JWTSecret          string
	CORSOrigins        []string
}

// NewRouter creates a new API router with all routes configured
func NewRouter(config RouterConfig) *Router {
	r := &Router{
		Mux:             chi.NewRouter(),
		logger:          config.Logger,
		domainService:   config.DomainService,
		userService:     config.UserService,
		aliasService:    config.AliasService,
		mailboxService:  config.MailboxService,
		messageService:  config.MessageService,
		queueService:    config.QueueService,
		settingsService: config.SettingsService,
		apiKeyRepo:      config.APIKeyRepo,
		jwtSecret:       config.JWTSecret,
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
		// Authentication routes (no auth required, but rate limited)
		r.Group(func(r chi.Router) {
			// Aggressive rate limiting for auth endpoints to prevent brute force
			r.Use(middleware.RateLimit(config.RateLimitRepo, config.Logger))

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
		})

		// Setup wizard routes (no auth required - runs before admin user exists)
		r.Group(func(r chi.Router) {
			setupHandler := handlers.NewSetupHandler(config.SetupService, config.Logger)
			r.Route("/setup", func(r chi.Router) {
				r.Get("/status", setupHandler.GetStatus)
				r.Get("/state", setupHandler.GetState)
				r.Post("/admin", setupHandler.CreateAdmin)
				r.Post("/complete", setupHandler.CompleteSetup)
			})
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			// JWT or API Key authentication required
			r.Use(middleware.Auth(config.JWTSecret, config.APIKeyRepo, config.Logger))

			// Rate limiting for API calls
			r.Use(middleware.RateLimit(config.RateLimitRepo, config.Logger))

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
				config.AliasService,
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

			// Settings management (admin only)
			settingsHandler := handlers.NewSettingsHandler(config.SettingsService, config.Logger)
			r.Route("/settings", func(r chi.Router) {
				r.Get("/server", settingsHandler.GetServer)
				r.Put("/server", settingsHandler.UpdateServer)
				r.Get("/security", settingsHandler.GetSecurity)
				r.Put("/security", settingsHandler.UpdateSecurity)
				r.Get("/tls", settingsHandler.GetTLS)
				r.Put("/tls", settingsHandler.UpdateTLS)
			})

			// PGP key management
			if config.PGPService != nil {
				pgpHandler := handlers.NewPGPHandler(config.PGPService, config.Logger)
				r.Route("/pgp", func(r chi.Router) {
					r.Post("/keys", pgpHandler.ImportKey)
					r.Get("/users/{user_id}/keys", pgpHandler.ListKeys)
					r.Get("/keys/{id}", pgpHandler.GetKey)
					r.Post("/keys/{id}/primary", pgpHandler.SetPrimary)
					r.Delete("/keys/{id}", pgpHandler.DeleteKey)
				})
			}

			// Audit log viewer (admin only)
			if config.AuditService != nil {
				auditHandler := handlers.NewAuditHandler(config.AuditService, config.Logger)
				r.Route("/audit", func(r chi.Router) {
					r.Get("/logs", auditHandler.ListLogs)
					r.Get("/stats", auditHandler.GetStats)
				})
			}

			// Webhook management
			if config.WebhookService != nil {
				webhookHandler := handlers.NewWebhookHandler(config.WebhookService, config.Logger)
				r.Route("/webhooks", func(r chi.Router) {
					r.Get("/", webhookHandler.ListWebhooks)
					r.Post("/", webhookHandler.CreateWebhook)
					r.Get("/{id}", webhookHandler.GetWebhook)
					r.Put("/{id}", webhookHandler.UpdateWebhook)
					r.Delete("/{id}", webhookHandler.DeleteWebhook)
					r.Post("/{id}/test", webhookHandler.TestWebhook)
					r.Get("/{id}/deliveries", webhookHandler.ListDeliveries)
				})
				r.Get("/webhooks/deliveries/{id}", webhookHandler.GetDelivery)
			}

			// Webmail API
			webmailHandler := handlers.NewWebmailHandler(config.MailboxService, config.MessageService, config.Logger)
			r.Route("/webmail", func(r chi.Router) {
				r.Get("/mailboxes", webmailHandler.ListMailboxes)
				r.Get("/mailboxes/{id}/messages", webmailHandler.ListMessages)
				r.Get("/messages/{id}", webmailHandler.GetMessage)
				r.Post("/messages", webmailHandler.SendMessage)
				r.Delete("/messages/{id}", webmailHandler.DeleteMessage)
				r.Post("/messages/{id}/move", webmailHandler.MoveMessage)
				r.Post("/messages/{id}/flags", webmailHandler.UpdateFlags)
				r.Get("/search", webmailHandler.SearchMessages)
				r.Get("/attachments/{id}", webmailHandler.DownloadAttachment)
				r.Post("/drafts", webmailHandler.SaveDraft)
				r.Get("/drafts", webmailHandler.ListDrafts)
				r.Get("/drafts/{id}", webmailHandler.GetDraft)
				r.Delete("/drafts/{id}", webmailHandler.DeleteDraft)

				// Contact integration
				if config.ContactService != nil && config.AddressbookService != nil {
					contactHandler := handlers.NewWebmailContactsHandler(config.ContactService, config.AddressbookService, config.Logger)
					r.Get("/contacts/search", contactHandler.SearchContacts)
					r.Get("/contacts/addressbooks", contactHandler.ListAddressbooks)
					r.Get("/contacts/addressbooks/{id}/contacts", contactHandler.ListContacts)
				}

				// Calendar integration
				if config.CalendarService != nil && config.EventService != nil {
					calendarHandler := handlers.NewWebmailCalendarHandler(config.CalendarService, config.EventService, config.Logger)
					r.Get("/calendar/calendars", calendarHandler.ListCalendars)
					r.Get("/calendar/upcoming", calendarHandler.GetUpcomingEvents)
					r.Post("/calendar/events", calendarHandler.CreateEvent)
					r.Post("/calendar/invitations", calendarHandler.ProcessInvitation)
				}
			})
		})
	})

	// PostmarkApp API compatibility endpoints
	// Mount at root level for PostmarkApp client compatibility
	r.Mount("/", postmark.NewRouter(config.DB, config.QueueService, config.Logger))

	// Webmail UI - Serves at /webmail/* with embedded or proxied assets
	r.Mount("/webmail", webmail.Handler(config.Logger))

	// Admin UI - must be last to act as catch-all for SPA routing
	// Serves at /admin/* with embedded or proxied assets
	// Using unified handler for Phase 1 migration
	r.Mount("/admin", admin.UnifiedHandler(config.Logger))

	return r
}