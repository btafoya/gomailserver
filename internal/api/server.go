package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/api/handler"
	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/repository"
)

// Server represents the admin API HTTP server
type Server struct {
	config     *config.APIConfig
	httpServer *http.Server
	logger     *zap.Logger
	router     http.Handler
}

// NewServer creates a new admin API server
func NewServer(
	cfg *config.APIConfig,
	domainRepo repository.DomainRepository,
	logger *zap.Logger,
) *Server {
	// Create handlers
	domainHandler := handler.NewDomainHandler(domainRepo, logger)

	// Setup router with middleware
	mux := http.NewServeMux()

	// Health check endpoint (no auth required)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Domain management endpoints (require admin authentication)
	authMiddleware := middleware.NewAuthMiddleware(cfg.AdminToken, logger)

	// Domain CRUD
	mux.Handle("GET /api/domains", authMiddleware.Wrap(http.HandlerFunc(domainHandler.ListDomains)))
	mux.Handle("GET /api/domains/{name}", authMiddleware.Wrap(http.HandlerFunc(domainHandler.GetDomain)))
	mux.Handle("POST /api/domains", authMiddleware.Wrap(http.HandlerFunc(domainHandler.CreateDomain)))
	mux.Handle("PUT /api/domains/{name}", authMiddleware.Wrap(http.HandlerFunc(domainHandler.UpdateDomain)))
	mux.Handle("DELETE /api/domains/{name}", authMiddleware.Wrap(http.HandlerFunc(domainHandler.DeleteDomain)))

	// Domain security configuration
	mux.Handle("GET /api/domains/{name}/security", authMiddleware.Wrap(http.HandlerFunc(domainHandler.GetDomainSecurity)))
	mux.Handle("PUT /api/domains/{name}/security", authMiddleware.Wrap(http.HandlerFunc(domainHandler.UpdateDomainSecurity)))

	// Default template management
	mux.Handle("GET /api/domains/_default", authMiddleware.Wrap(http.HandlerFunc(domainHandler.GetDefaultTemplate)))
	mux.Handle("PUT /api/domains/_default", authMiddleware.Wrap(http.HandlerFunc(domainHandler.UpdateDefaultTemplate)))

	// Wrap with logging middleware
	router := middleware.NewLoggingMiddleware(logger).Wrap(mux)

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
