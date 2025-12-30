package imap

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/emersion/go-imap/server"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/service"
)

// Server manages IMAP server instances
type Server struct {
	imap     *server.Server // Port 143 (STARTTLS)
	imaps    net.Listener   // Port 993 (implicit TLS)
	backend  *Backend
	cfg      *config.IMAPConfig
	tlsCfg   *tls.Config
	logger   *zap.Logger
	wg       sync.WaitGroup
	cancel   context.CancelFunc
}

// NewServer creates a new IMAP server manager
func NewServer(cfg *config.IMAPConfig, tlsCfg *tls.Config, userSvc *service.UserService, mailboxSvc *service.MailboxService, messageSvc *service.MessageService, logger *zap.Logger) *Server {
	backend := &Backend{
		userService:    userSvc,
		mailboxService: mailboxSvc,
		messageService: messageSvc,
		logger:         logger,
	}

	s := &Server{
		backend: backend,
		cfg:     cfg,
		tlsCfg:  tlsCfg,
		logger:  logger,
	}

	// Initialize IMAP server
	s.imap = s.createIMAPServer()

	return s
}

// createIMAPServer creates the IMAP server instance
func (s *Server) createIMAPServer() *server.Server {
	srv := server.New(s.backend)
	srv.Addr = fmt.Sprintf(":%d", s.cfg.Port)
	srv.AllowInsecureAuth = false // Require TLS for auth
	srv.AutoLogout = time.Duration(s.cfg.IdleTimeout) * time.Second

	// STARTTLS configuration
	if s.tlsCfg != nil {
		srv.TLSConfig = s.tlsCfg
	}

	return srv
}

// Start starts all IMAP servers
func (s *Server) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)

	// Start IMAP server (143) with STARTTLS
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.logger.Info("starting IMAP server",
			zap.Int("port", s.cfg.Port),
			zap.String("tls_mode", "STARTTLS"),
		)
		if err := s.imap.ListenAndServe(); err != nil && ctx.Err() == nil {
			s.logger.Error("IMAP server error", zap.Error(err))
		}
	}()

	// Start IMAPS server (993) with implicit TLS
	if s.tlsCfg != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.logger.Info("starting IMAPS server",
				zap.Int("port", s.cfg.IMAPSPort),
				zap.String("tls_mode", "implicit"),
			)

			ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.IMAPSPort))
			if err != nil {
				s.logger.Error("failed to start IMAPS listener", zap.Error(err))
				return
			}

			s.imaps = tls.NewListener(ln, s.tlsCfg)

			imapsServer := server.New(s.backend)
			imapsServer.AllowInsecureAuth = false
			imapsServer.AutoLogout = time.Duration(s.cfg.IdleTimeout) * time.Second

			if err := imapsServer.Serve(s.imaps); err != nil && ctx.Err() == nil {
				s.logger.Error("IMAPS server error", zap.Error(err))
			}
		}()
	}

	s.logger.Info("IMAP servers started",
		zap.Int("imap_port", s.cfg.Port),
		zap.Int("imaps_port", s.cfg.IMAPSPort),
		zap.Int("idle_timeout", s.cfg.IdleTimeout),
	)

	return nil
}

// Shutdown performs graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down IMAP servers")

	if s.cancel != nil {
		s.cancel()
	}

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var shutdownErr error
	shutdownDone := make(chan struct{})

	go func() {
		defer close(shutdownDone)

		// Close IMAP server
		if err := s.imap.Close(); err != nil {
			s.logger.Warn("IMAP server close error", zap.Error(err))
			shutdownErr = err
		}

		// Close IMAPS listener
		if s.imaps != nil {
			if err := s.imaps.Close(); err != nil {
				s.logger.Warn("IMAPS listener close error", zap.Error(err))
				shutdownErr = err
			}
		}
	}()

	select {
	case <-shutdownDone:
		s.logger.Info("IMAP servers shutdown complete")
	case <-shutdownCtx.Done():
		s.logger.Warn("IMAP servers shutdown timeout")
		return shutdownCtx.Err()
	}

	// Wait for all goroutines to finish
	s.wg.Wait()

	return shutdownErr
}
