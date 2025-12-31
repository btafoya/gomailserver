package service

import (
	"context"
	"fmt"
	"os"

	"github.com/btafoya/gomailserver/internal/config"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// SettingsService handles configuration management operations
type SettingsService struct {
	config     *config.Config
	configPath string
	logger     *zap.Logger
}

// NewSettingsService creates a new settings service instance
func NewSettingsService(cfg *config.Config, configPath string, logger *zap.Logger) *SettingsService {
	return &SettingsService{
		config:     cfg,
		configPath: configPath,
		logger:     logger,
	}
}

// ServerSettings represents server configuration settings
type ServerSettings struct {
	Hostname           string `json:"hostname"`
	Domain             string `json:"domain"`
	SMTPSubmissionPort int    `json:"smtp_submission_port"`
	SMTPRelayPort      int    `json:"smtp_relay_port"`
	SMTPSPort          int    `json:"smtps_port"`
	IMAPPort           int    `json:"imap_port"`
	IMAPSPort          int    `json:"imaps_port"`
	APIPort            int    `json:"api_port"`
	MaxMessageSize     int64  `json:"max_message_size"`
}

// SecuritySettings represents security configuration settings
type SecuritySettings struct {
	JWTSecret              string `json:"jwt_secret,omitempty"`
	ClamAVEnabled          bool   `json:"clamav_enabled"`
	ClamAVSocketPath       string `json:"clamav_socket_path"`
	SpamAssassinEnabled    bool   `json:"spamassassin_enabled"`
	SpamAssassinHost       string `json:"spamassassin_host"`
	SpamAssassinPort       int    `json:"spamassassin_port"`
	RateLimitEnabled       bool   `json:"rate_limit_enabled"`
	RateLimitRequests      int    `json:"rate_limit_requests"`
	RateLimitWindow        int    `json:"rate_limit_window"`
}

// TLSSettings represents TLS/certificate configuration settings
type TLSSettings struct {
	ACMEEnabled   bool   `json:"acme_enabled"`
	ACMEEmail     string `json:"acme_email"`
	ACMEProvider  string `json:"acme_provider"`
	ACMEAPIToken  string `json:"acme_api_token,omitempty"`
	CertFile      string `json:"cert_file"`
	KeyFile       string `json:"key_file"`
}

// GetServerSettings retrieves current server configuration
func (s *SettingsService) GetServerSettings(ctx context.Context) (*ServerSettings, error) {
	settings := &ServerSettings{
		Hostname:           s.config.Server.Hostname,
		Domain:             s.config.Server.Domain,
		SMTPSubmissionPort: s.config.SMTP.SubmissionPort,
		SMTPRelayPort:      s.config.SMTP.RelayPort,
		SMTPSPort:          s.config.SMTP.SMTPSPort,
		IMAPPort:           s.config.IMAP.Port,
		IMAPSPort:          s.config.IMAP.IMAPSPort,
		APIPort:            s.config.API.Port,
		MaxMessageSize:     s.config.SMTP.MaxMessageSize,
	}

	return settings, nil
}

// UpdateServerSettings updates server configuration
func (s *SettingsService) UpdateServerSettings(ctx context.Context, settings *ServerSettings) error {
	// Update in-memory config
	s.config.Server.Hostname = settings.Hostname
	s.config.Server.Domain = settings.Domain
	s.config.SMTP.SubmissionPort = settings.SMTPSubmissionPort
	s.config.SMTP.RelayPort = settings.SMTPRelayPort
	s.config.SMTP.SMTPSPort = settings.SMTPSPort
	s.config.SMTP.MaxMessageSize = settings.MaxMessageSize
	s.config.IMAP.Port = settings.IMAPPort
	s.config.IMAP.IMAPSPort = settings.IMAPSPort
	s.config.API.Port = settings.APIPort

	// Write to YAML file
	if err := s.writeConfig(); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	s.logger.Info("Server settings updated",
		zap.String("hostname", settings.Hostname),
		zap.String("domain", settings.Domain),
	)

	return nil
}

// GetSecuritySettings retrieves current security configuration
func (s *SettingsService) GetSecuritySettings(ctx context.Context) (*SecuritySettings, error) {
	settings := &SecuritySettings{
		// JWT secret is sensitive, don't return it
		JWTSecret:              "",
		ClamAVEnabled:          s.config.Security.ClamAV.SocketPath != "",
		ClamAVSocketPath:       s.config.Security.ClamAV.SocketPath,
		SpamAssassinEnabled:    s.config.Security.SpamAssassin.Host != "",
		SpamAssassinHost:       s.config.Security.SpamAssassin.Host,
		SpamAssassinPort:       s.config.Security.SpamAssassin.Port,
		// Rate limiting settings would come from database or additional config
		RateLimitEnabled:       true,
		RateLimitRequests:      100,
		RateLimitWindow:        60,
	}

	return settings, nil
}

// UpdateSecuritySettings updates security configuration
func (s *SettingsService) UpdateSecuritySettings(ctx context.Context, settings *SecuritySettings) error {
	// Update JWT secret if provided
	if settings.JWTSecret != "" {
		s.config.API.JWTSecret = settings.JWTSecret
	}

	// Update ClamAV settings
	if settings.ClamAVEnabled {
		s.config.Security.ClamAV.SocketPath = settings.ClamAVSocketPath
	} else {
		s.config.Security.ClamAV.SocketPath = ""
	}

	// Update SpamAssassin settings
	if settings.SpamAssassinEnabled {
		s.config.Security.SpamAssassin.Host = settings.SpamAssassinHost
		s.config.Security.SpamAssassin.Port = settings.SpamAssassinPort
	} else {
		s.config.Security.SpamAssassin.Host = ""
		s.config.Security.SpamAssassin.Port = 0
	}

	// Write to YAML file
	if err := s.writeConfig(); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	s.logger.Info("Security settings updated",
		zap.Bool("clamav_enabled", settings.ClamAVEnabled),
		zap.Bool("spamassassin_enabled", settings.SpamAssassinEnabled),
	)

	return nil
}

// GetTLSSettings retrieves current TLS/certificate configuration
func (s *SettingsService) GetTLSSettings(ctx context.Context) (*TLSSettings, error) {
	settings := &TLSSettings{
		ACMEEnabled:   s.config.TLS.ACME.Enabled,
		ACMEEmail:     s.config.TLS.ACME.Email,
		ACMEProvider:  s.config.TLS.ACME.Provider,
		// API token is sensitive, don't return it
		ACMEAPIToken:  "",
		CertFile:      s.config.TLS.CertFile,
		KeyFile:       s.config.TLS.KeyFile,
	}

	return settings, nil
}

// UpdateTLSSettings updates TLS/certificate configuration
func (s *SettingsService) UpdateTLSSettings(ctx context.Context, settings *TLSSettings) error {
	// Update ACME settings
	s.config.TLS.ACME.Enabled = settings.ACMEEnabled
	s.config.TLS.ACME.Email = settings.ACMEEmail
	s.config.TLS.ACME.Provider = settings.ACMEProvider

	// Update API token if provided
	if settings.ACMEAPIToken != "" {
		s.config.TLS.ACME.APIToken = settings.ACMEAPIToken
	}

	// Update manual certificate paths (only if ACME is disabled)
	if !settings.ACMEEnabled {
		s.config.TLS.CertFile = settings.CertFile
		s.config.TLS.KeyFile = settings.KeyFile
	}

	// Write to YAML file
	if err := s.writeConfig(); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	s.logger.Info("TLS settings updated",
		zap.Bool("acme_enabled", settings.ACMEEnabled),
		zap.String("acme_provider", settings.ACMEProvider),
	)

	return nil
}

// writeConfig writes the current configuration to the YAML file
func (s *SettingsService) writeConfig() error {
	// Marshal config to YAML
	data, err := yaml.Marshal(s.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with proper permissions
	if err := os.WriteFile(s.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
