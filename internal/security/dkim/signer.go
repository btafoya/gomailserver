package dkim

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/emersion/go-msgauth/dkim"

	"github.com/btafoya/gomailserver/internal/domain"
)

// This is a placeholder for the actual service
type placeholderDomainService struct{}

func (s *placeholderDomainService) GetDKIMConfig(domainName string) (*domain.DKIMConfig, error) {
	// In a real implementation, this would fetch the DKIM configuration
	// for the given domain from the database.
	// For now, we'll return a dummy config.
	kp, err := GenerateRSAKeyPair(2048)
	if err != nil {
		return nil, err
	}

	return &domain.DKIMConfig{
		Domain:     domainName,
		Selector:   kp.Selector,
		PrivateKey: []byte(kp.PrivateKey),
	}, nil
}

type Signer struct {
	domainService *placeholderDomainService
}

func NewSigner() *Signer {
	return &Signer{domainService: &placeholderDomainService{}}
}

func (s *Signer) Sign(domainName string, message []byte) ([]byte, error) {
	domainCfg, err := s.domainService.GetDKIMConfig(domainName)
	if err != nil {
		return message, nil // Return unsigned if no DKIM config
	}

	// Parse the private key from PEM format
	privateKey, err := parsePrivateKey(domainCfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	options := &dkim.SignOptions{
		Domain:   domainName,
		Selector: domainCfg.Selector,
		Signer:   privateKey,
		HeaderKeys: []string{
			"From", "To", "Subject", "Date", "Message-ID",
			"MIME-Version", "Content-Type",
		},
	}

	r := bytes.NewReader(message)
	var buf bytes.Buffer
	if err := dkim.Sign(&buf, r, options); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parsePrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try parsing as PKCS1 (RSA PRIVATE KEY)
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Try parsing as PKCS8 (PRIVATE KEY)
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("not an RSA private key")
	}

	return nil, fmt.Errorf("unable to parse private key")
}
