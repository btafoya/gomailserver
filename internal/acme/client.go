package acme

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"go.uber.org/zap"
)

// Client manages ACME certificate operations
type Client struct {
	logger     *zap.Logger
	email      string
	legoClient *lego.Client
	privateKey crypto.PrivateKey
}

// User implements registration.User for ACME
type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

// parseCertificateDates extracts NotBefore and NotAfter times from a PEM-encoded certificate
func parseCertificateDates(certPEM []byte) (notBefore, notAfter time.Time, err error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to decode PEM certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert.NotBefore, cert.NotAfter, nil
}

// NewClient creates a new ACME client
func NewClient(email string, production bool, logger *zap.Logger) (*Client, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	user := &User{
		Email: email,
		key:   privateKey,
	}

	config := lego.NewConfig(user)
	config.Certificate.KeyType = certcrypto.RSA2048

	if production {
		config.CADirURL = lego.LEDirectoryProduction
	} else {
		config.CADirURL = lego.LEDirectoryStaging
	}

	legoClient, err := lego.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create lego client: %w", err)
	}

	reg, err := legoClient.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	user.Registration = reg

	return &Client{
		logger:     logger,
		email:      email,
		legoClient: legoClient,
		privateKey: privateKey,
	}, nil
}

// ObtainCertificate requests a new certificate for the given domains
func (c *Client) ObtainCertificate(ctx context.Context, domains []string) (*certificate.Resource, error) {
	c.logger.Info("requesting certificate", zap.Strings("domains", domains))

	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}

	certificates, err := c.legoClient.Certificate.Obtain(request)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain certificate: %w", err)
	}

	// Parse certificate dates from the PEM-encoded certificate
	notBefore, notAfter, err := parseCertificateDates(certificates.Certificate)
	if err != nil {
		c.logger.Warn("failed to parse certificate dates", zap.Error(err))
		c.logger.Info("certificate obtained successfully", zap.String("domain", certificates.Domain))
	} else {
		c.logger.Info("certificate obtained successfully",
			zap.String("domain", certificates.Domain),
			zap.Time("not_before", notBefore),
			zap.Time("not_after", notAfter))
	}

	return certificates, nil
}

// RenewCertificate renews an existing certificate
func (c *Client) RenewCertificate(ctx context.Context, cert *certificate.Resource) (*certificate.Resource, error) {
	c.logger.Info("renewing certificate", zap.String("domain", cert.Domain))

	certificates, err := c.legoClient.Certificate.Renew(*cert, true, false, "")
	if err != nil {
		return nil, fmt.Errorf("failed to renew certificate: %w", err)
	}

	// Parse certificate dates from the PEM-encoded certificate
	_, notAfter, err := parseCertificateDates(certificates.Certificate)
	if err != nil {
		c.logger.Warn("failed to parse certificate dates", zap.Error(err))
		c.logger.Info("certificate renewed successfully", zap.String("domain", certificates.Domain))
	} else {
		c.logger.Info("certificate renewed successfully",
			zap.String("domain", certificates.Domain),
			zap.Time("not_after", notAfter))
	}

	return certificates, nil
}

// RevokeCertificate revokes a certificate
func (c *Client) RevokeCertificate(ctx context.Context, cert *certificate.Resource) error {
	c.logger.Info("revoking certificate", zap.String("domain", cert.Domain))

	err := c.legoClient.Certificate.Revoke(cert.Certificate)
	if err != nil {
		return fmt.Errorf("failed to revoke certificate: %w", err)
	}

	c.logger.Info("certificate revoked successfully", zap.String("domain", cert.Domain))
	return nil
}

// NeedsRenewal checks if a certificate needs renewal (within 30 days of expiry)
func NeedsRenewal(notAfter time.Time) bool {
	return time.Until(notAfter) < 30*24*time.Hour
}
