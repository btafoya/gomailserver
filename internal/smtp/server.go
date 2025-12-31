package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/emersion/go-smtp"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/config"
)

// Server manages SMTP server instances
type Server struct {
	submission *smtp.Server // Port 587
	relay      *smtp.Server // Port 25
	smtps      *smtp.Server // Port 465
	backend    *Backend
	cfg        *config.SMTPConfig
	tlsCfg     *tls.Config
	logger     *zap.Logger
	wg         sync.WaitGroup
	cancel     context.CancelFunc
}

// NewServer creates a new SMTP server manager
func NewServer(cfg *config.SMTPConfig, tlsCfg *tls.Config, backend *Backend, logger *zap.Logger) *Server {
	s := &Server{
		backend: backend,
		cfg:     cfg,
		tlsCfg:  tlsCfg,
		logger:  logger,
	}

	// Initialize SMTP servers
	s.submission = s.createSubmissionServer()
	s.relay = s.createRelayServer()
	s.smtps = s.createSMTPSServer()

	return s
}

// createSubmissionServer creates port 587 submission server
func (s *Server) createSubmissionServer() *smtp.Server {
	srv := smtp.NewServer(s.backend)
	srv.Addr = fmt.Sprintf(":%d", s.cfg.SubmissionPort)
	srv.Domain = s.cfg.Hostname
	srv.ReadTimeout = 30 * time.Second
	srv.WriteTimeout = 30 * time.Second
	srv.MaxMessageBytes = int64(s.cfg.MaxMessageSize)
	srv.MaxRecipients = 100
	srv.AllowInsecureAuth = false
	srv.EnableSMTPUTF8 = true
	srv.EnableREQUIRETLS = true

	// STARTTLS configuration
	if s.tlsCfg != nil {
		srv.TLSConfig = s.tlsCfg
	}

	return srv
}

// createRelayServer creates port 25 MX relay server
func (s *Server) createRelayServer() *smtp.Server {
	srv := smtp.NewServer(s.backend)
	srv.Addr = fmt.Sprintf(":%d", s.cfg.RelayPort)
	srv.Domain = s.cfg.Hostname
	srv.ReadTimeout = 30 * time.Second
	srv.WriteTimeout = 30 * time.Second
	srv.MaxMessageBytes = int64(s.cfg.MaxMessageSize)
	srv.MaxRecipients = 100
	srv.AllowInsecureAuth = true // Allow for receiving mail
	srv.EnableSMTPUTF8 = true

	// Optional TLS
	if s.tlsCfg != nil {
		srv.TLSConfig = s.tlsCfg
	}

	return srv
}

// createSMTPSServer creates port 465 SMTPS server (implicit TLS)
func (s *Server) createSMTPSServer() *smtp.Server {
	srv := smtp.NewServer(s.backend)
	srv.Addr = fmt.Sprintf(":%d", s.cfg.SMTPSPort)
	srv.Domain = s.cfg.Hostname
	srv.ReadTimeout = 30 * time.Second
	srv.WriteTimeout = 30 * time.Second
	srv.MaxMessageBytes = int64(s.cfg.MaxMessageSize)
	srv.MaxRecipients = 100
	srv.AllowInsecureAuth = false
	srv.EnableSMTPUTF8 = true

	return srv
}

// Start starts all SMTP servers
func (s *Server) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)

	// Start submission server (587)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.logger.Info("starting SMTP submission server",
			zap.Int("port", s.cfg.SubmissionPort),
			zap.String("hostname", s.cfg.Hostname),
		)
		if err := s.submission.ListenAndServe(); err != nil && ctx.Err() == nil {
			s.logger.Error("submission server error", zap.Error(err))
		}
	}()

	// Start relay server (25)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.logger.Info("starting SMTP relay server",
			zap.Int("port", s.cfg.RelayPort),
			zap.String("hostname", s.cfg.Hostname),
		)
		if err := s.relay.ListenAndServe(); err != nil && ctx.Err() == nil {
			s.logger.Error("relay server error", zap.Error(err))
		}
	}()

	// Start SMTPS server (465) with implicit TLS
	if s.tlsCfg != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.logger.Info("starting SMTPS server",
				zap.Int("port", s.cfg.SMTPSPort),
				zap.String("hostname", s.cfg.Hostname),
			)

			ln, err := net.Listen("tcp", s.smtps.Addr)
			if err != nil {
				s.logger.Error("failed to start SMTPS listener", zap.Error(err))
				return
			}

			tlsListener := tls.NewListener(ln, s.tlsCfg)
			if err := s.smtps.Serve(tlsListener); err != nil && ctx.Err() == nil {
				s.logger.Error("SMTPS server error", zap.Error(err))
			}
		}()
	}

	s.logger.Info("SMTP servers started",
		zap.Int("submission_port", s.cfg.SubmissionPort),
		zap.Int("relay_port", s.cfg.RelayPort),
		zap.Int("smtps_port", s.cfg.SMTPSPort),
	)

	return nil
}

// Shutdown performs graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down SMTP servers")

	if s.cancel != nil {
		s.cancel()
	}

	// Shutdown all servers
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var shutdownErr error
	shutdownDone := make(chan struct{})

	go func() {
		defer close(shutdownDone)

		if err := s.submission.Shutdown(shutdownCtx); err != nil {
			s.logger.Warn("submission server shutdown error", zap.Error(err))
			shutdownErr = err
		}

		if err := s.relay.Shutdown(shutdownCtx); err != nil {
			s.logger.Warn("relay server shutdown error", zap.Error(err))
			shutdownErr = err
		}

		if err := s.smtps.Shutdown(shutdownCtx); err != nil {
			s.logger.Warn("SMTPS server shutdown error", zap.Error(err))
			shutdownErr = err
		}
	}()

	select {
	case <-shutdownDone:
		s.logger.Info("SMTP servers shutdown complete")
	case <-shutdownCtx.Done():
		s.logger.Warn("SMTP servers shutdown timeout")
		return shutdownCtx.Err()
	}

	// Wait for all goroutines to finish
	s.wg.Wait()

	return shutdownErr
}
