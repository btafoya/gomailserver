package acme

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/go-acme/lego/v4/certificate"
	"go.uber.org/zap"
)

// Service manages ACME certificate lifecycle
type Service struct {
	client *Client
	db     *database.DB
	logger *zap.Logger
}

// NewService creates a new ACME service
func NewService(email string, production bool, db *database.DB, logger *zap.Logger) (*Service, error) {
	client, err := NewClient(email, production, logger)
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
		db:     db,
		logger: logger,
	}, nil
}

// ObtainAndStoreCertificate requests and stores a new certificate
func (s *Service) ObtainAndStoreCertificate(ctx context.Context, domain string, altNames []string) error {
	domains := append([]string{domain}, altNames...)

	cert, err := s.client.ObtainCertificate(ctx, domains)
	if err != nil {
		return fmt.Errorf("failed to obtain certificate: %w", err)
	}

	return s.storeCertificate(domain, cert)
}

// RenewCertificate renews a certificate for a domain
func (s *Service) RenewCertificate(ctx context.Context, domain string) error {
	// Get existing certificate from database
	certData, err := s.getCertificateFromDB(domain)
	if err != nil {
		return fmt.Errorf("failed to get certificate: %w", err)
	}

	certResource := &certificate.Resource{
		Domain:            domain,
		CertURL:           certData.ACMEOrderURL,
		CertStableURL:     certData.ACMEAccountURL,
		PrivateKey:        []byte(certData.PrivateKey),
		Certificate:       []byte(certData.Certificate),
		IssuerCertificate: nil,
		CSR:               nil,
	}

	renewed, err := s.client.RenewCertificate(ctx, certResource)
	if err != nil {
		return fmt.Errorf("failed to renew certificate: %w", err)
	}

	return s.storeCertificate(domain, renewed)
}

// CheckAndRenewExpiring checks all certificates and renews those expiring soon
func (s *Service) CheckAndRenewExpiring(ctx context.Context) error {
	certs, err := s.getAllCertificates()
	if err != nil {
		return err
	}

	for _, cert := range certs {
		if cert.Status != "active" || !cert.AutoRenew {
			continue
		}

		if NeedsRenewal(cert.NotAfter) {
			s.logger.Info("renewing expiring certificate",
				zap.String("domain", cert.DomainName),
				zap.Time("expires", cert.NotAfter))

			if err := s.RenewCertificate(ctx, cert.DomainName); err != nil {
				s.logger.Error("failed to renew certificate",
					zap.String("domain", cert.DomainName),
					zap.Error(err))
				continue
			}
		}
	}

	return nil
}

func (s *Service) storeCertificate(domain string, cert *certificate.Resource) error {
	// Parse certificate to get validity dates
	block, _ := pem.Decode(cert.Certificate)
	if block == nil {
		return fmt.Errorf("failed to decode certificate PEM")
	}

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Store in database
	query := `
		INSERT INTO tls_certificates
		(domain_name, certificate, private_key, issuer, not_before, not_after,
		 status, acme_account_url, acme_order_url, auto_renew, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 'active', ?, ?, 1, datetime('now'), datetime('now'))
		ON CONFLICT(domain_name) DO UPDATE SET
		certificate = excluded.certificate,
		private_key = excluded.private_key,
		not_before = excluded.not_before,
		not_after = excluded.not_after,
		status = 'active',
		updated_at = datetime('now')
	`

	_, err = s.db.Exec(query,
		domain,
		string(cert.Certificate),
		string(cert.PrivateKey),
		x509Cert.Issuer.String(),
		x509Cert.NotBefore,
		x509Cert.NotAfter,
		cert.CertStableURL,
		cert.CertURL,
	)

	return err
}

type CertificateData struct {
	ID             int64
	DomainName     string
	Certificate    string
	PrivateKey     string
	Issuer         string
	NotBefore      time.Time
	NotAfter       time.Time
	Status         string
	ACMEAccountURL string
	ACMEOrderURL   string
	AutoRenew      bool
}

func (s *Service) getCertificateFromDB(domain string) (*CertificateData, error) {
	query := `
		SELECT id, domain_name, certificate, private_key, issuer,
		       not_before, not_after, status, acme_account_url, acme_order_url, auto_renew
		FROM tls_certificates
		WHERE domain_name = ?
	`

	var cert CertificateData
	err := s.db.QueryRow(query, domain).Scan(
		&cert.ID,
		&cert.DomainName,
		&cert.Certificate,
		&cert.PrivateKey,
		&cert.Issuer,
		&cert.NotBefore,
		&cert.NotAfter,
		&cert.Status,
		&cert.ACMEAccountURL,
		&cert.ACMEOrderURL,
		&cert.AutoRenew,
	)

	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (s *Service) getAllCertificates() ([]*CertificateData, error) {
	query := `
		SELECT id, domain_name, certificate, private_key, issuer,
		       not_before, not_after, status, acme_account_url, acme_order_url, auto_renew
		FROM tls_certificates
		ORDER BY domain_name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var certs []*CertificateData
	for rows.Next() {
		var cert CertificateData
		err := rows.Scan(
			&cert.ID,
			&cert.DomainName,
			&cert.Certificate,
			&cert.PrivateKey,
			&cert.Issuer,
			&cert.NotBefore,
			&cert.NotAfter,
			&cert.Status,
			&cert.ACMEAccountURL,
			&cert.ACMEOrderURL,
			&cert.AutoRenew,
		)
		if err != nil {
			return nil, err
		}
		certs = append(certs, &cert)
	}

	return certs, rows.Err()
}
