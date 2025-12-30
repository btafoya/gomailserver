package dkim

import (
	"bytes"

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

	options := &dkim.SignOptions{
		Domain:   domainName,
		Selector: domainCfg.Selector,
		Signer:   bytes.NewReader(domainCfg.PrivateKey),
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
