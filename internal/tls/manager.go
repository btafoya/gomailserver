package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/config"
)

// Manager handles TLS certificate operations
type Manager struct {
	cfg      *config.TLSConfig
	hostname string
	logger   *zap.Logger
	cert     *tls.Certificate
}

// NewManager creates a new TLS manager
func NewManager(cfg *config.TLSConfig, hostname string, logger *zap.Logger) (*Manager, error) {
	m := &Manager{
		cfg:      cfg,
		hostname: hostname,
		logger:   logger,
	}

	if err := m.loadOrGenerateCertificate(); err != nil {
		return nil, fmt.Errorf("failed to initialize TLS: %w", err)
	}

	return m, nil
}

// GetTLSConfig returns a TLS configuration
func (m *Manager) GetTLSConfig() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{*m.cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   false,
		Renegotiation:            tls.RenegotiateNever,
	}
}

// loadOrGenerateCertificate loads certificate from files or generates a self-signed one
func (m *Manager) loadOrGenerateCertificate() error {
	// Try loading from files first
	if m.cfg.CertFile != "" && m.cfg.KeyFile != "" {
		m.logger.Info("loading TLS certificate from files",
			zap.String("cert", m.cfg.CertFile),
			zap.String("key", m.cfg.KeyFile),
		)

		cert, err := tls.LoadX509KeyPair(m.cfg.CertFile, m.cfg.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load certificate: %w", err)
		}

		m.cert = &cert
		m.logger.Info("TLS certificate loaded successfully")
		return nil
	}

	// Generate self-signed certificate
	m.logger.Warn("generating self-signed certificate (NOT FOR PRODUCTION USE)")

	cert, err := m.generateSelfSignedCert()
	if err != nil {
		return fmt.Errorf("failed to generate self-signed certificate: %w", err)
	}

	m.cert = cert
	m.logger.Info("self-signed certificate generated",
		zap.String("common_name", m.hostname),
		zap.Int("valid_days", 365),
	)

	return nil
}

// generateSelfSignedCert generates a self-signed certificate for testing
func (m *Manager) generateSelfSignedCert() (*tls.Certificate, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year validity

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"gomailserver"},
			CommonName:   m.hostname,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{m.hostname},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode certificate and key to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Parse into tls.Certificate
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return &cert, nil
}

// SaveSelfSignedCert saves a self-signed certificate to files
func (m *Manager) SaveSelfSignedCert(certPath, keyPath string) error {
	if m.cert == nil {
		return fmt.Errorf("no certificate loaded")
	}

	// Get certificate and key
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: m.cert.Certificate[0],
	})

	privateKey, ok := m.cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("private key is not RSA")
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Write certificate file
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Write key file
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	m.logger.Info("self-signed certificate saved",
		zap.String("cert", certPath),
		zap.String("key", keyPath),
	)

	return nil
}

// Reload reloads the certificate from files
func (m *Manager) Reload() error {
	if m.cfg.CertFile == "" || m.cfg.KeyFile == "" {
		return fmt.Errorf("cannot reload: no certificate files configured")
	}

	m.logger.Info("reloading TLS certificate")

	cert, err := tls.LoadX509KeyPair(m.cfg.CertFile, m.cfg.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to reload certificate: %w", err)
	}

	m.cert = &cert
	m.logger.Info("TLS certificate reloaded successfully")

	return nil
}

// GetCertificate returns the loaded certificate
func (m *Manager) GetCertificate() *tls.Certificate {
	return m.cert
}

// ValidateExpiry checks if the certificate is expiring soon
func (m *Manager) ValidateExpiry(warningDays int) error {
	if m.cert == nil {
		return fmt.Errorf("no certificate loaded")
	}

	x509Cert, err := x509.ParseCertificate(m.cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	now := time.Now()
	daysUntilExpiry := int(x509Cert.NotAfter.Sub(now).Hours() / 24)

	if now.After(x509Cert.NotAfter) {
		return fmt.Errorf("certificate has expired")
	}

	if daysUntilExpiry <= warningDays {
		m.logger.Warn("certificate expiring soon",
			zap.Int("days_until_expiry", daysUntilExpiry),
			zap.Time("expires_at", x509Cert.NotAfter),
		)
	}

	return nil
}
