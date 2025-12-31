package webdav

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/btafoya/gomailserver/internal/repository"
	"go.uber.org/zap"
)

// Server represents the WebDAV HTTP server
type Server struct {
	config     *Config
	httpServer *http.Server
	logger     *zap.Logger
	caldav     http.Handler
	carddav    http.Handler
}

// Config contains WebDAV server configuration
type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

// NewServer creates a new WebDAV server with CalDAV and CardDAV handlers
func NewServer(cfg *Config, caldavHandler, carddavHandler http.Handler, userRepo repository.UserRepository, logger *zap.Logger) *Server {
	mux := http.NewServeMux()

	// Create authentication middleware
	authMiddleware := BasicAuthMiddleware(userRepo, logger)

	// CalDAV endpoints with authentication
	mux.Handle("/caldav/", authMiddleware(caldavHandler))

	// CardDAV endpoints with authentication
	mux.Handle("/carddav/", authMiddleware(carddavHandler))

	// Well-known redirects (RFC 6764) - no authentication required
	mux.HandleFunc("/.well-known/caldav", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/caldav/", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/.well-known/carddav", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/carddav/", http.StatusMovedPermanently)
	})

	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &Server{
		config:     cfg,
		httpServer: httpServer,
		logger:     logger,
		caldav:     caldavHandler,
		carddav:    carddavHandler,
	}
}

// Start starts the WebDAV server
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("starting WebDAV server",
		zap.Int("port", s.config.Port),
	)

	// Start server in goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("WebDAV server error", zap.Error(err))
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the WebDAV server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down WebDAV server")
	return s.httpServer.Shutdown(ctx)
}
