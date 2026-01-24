package service

import (
	"context"
	"fmt"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"go.uber.org/zap"
)

const DefaultTemplateDomainName = "_default"

// DomainService handles domain business logic including default templates
type DomainService struct {
	repo   repository.DomainRepository
	logger *zap.Logger
}

// NewDomainService creates a new domain service
func NewDomainService(repos *repository.Repositories, logger *zap.Logger) *DomainService {
	return &DomainService{
		repo:   repos.Domain,
		logger: logger,
	}
}

// EnsureDefaultTemplate creates default domain template if it doesn't exist
func (s *DomainService) EnsureDefaultTemplate() error {
	// Check if default template exists
	_, err := s.repo.GetByName(DefaultTemplateDomainName)
	if err == nil {
		// Default template already exists
		return nil
	}

	// Create default template with recommended security settings
	defaultTemplate := &domain.Domain{
		Name:           DefaultTemplateDomainName,
		Status:         "active",
		MaxUsers:       0,          // unlimited
		MaxMailboxSize: 0,          // unlimited
		DefaultQuota:   1073741824, // 1GB

		// DKIM defaults
		DKIMSelector:       "default",
		DKIMPrivateKey:     "",
		DKIMPublicKey:      "",
		DKIMSigningEnabled: true,
		DKIMVerifyEnabled:  true,
		DKIMKeySize:        2048,
		DKIMKeyType:        "rsa",
		DKIMHeadersToSign:  `["From","To","Subject","Date","Message-ID","MIME-Version","Content-Type"]`,

		// SPF defaults
		SPFRecord:         "v=spf1 mx ~all",
		SPFEnabled:        true,
		SPFDNSServer:      "8.8.8.8:53",
		SPFDNSTimeout:     5,
		SPFMaxLookups:     10,
		SPFFailAction:     "reject",
		SPFSoftFailAction: "accept",

		// DMARC defaults
		DMARCPolicy:        "none",
		DMARCEnabled:       true,
		DMARCDNSServer:     "8.8.8.8:53",
		DMARCDNSTimeout:    5,
		DMARCReportEnabled: false,
		DMARCReportEmail:   "",

		// ClamAV defaults
		ClamAVEnabled:     true,
		ClamAVMaxScanSize: 52428800, // 50MB
		ClamAVVirusAction: "reject",
		ClamAVFailAction:  "accept",

		// SpamAssassin defaults
		SpamEnabled:         true,
		SpamRejectScore:     10.0,
		SpamQuarantineScore: 5.0,
		SpamLearningEnabled: true,

		// Greylisting defaults
		GreylistEnabled:         true,
		GreylistDelayMinutes:    5,
		GreylistExpiryDays:      30,
		GreylistCleanupInterval: 3600,
		GreylistWhitelistAfter:  3,

		// Rate limiting defaults
		RateLimitEnabled:         true,
		RateLimitSMTPPerIP:       `{"count":100,"window_minutes":60}`,
		RateLimitSMTPPerUser:     `{"count":500,"window_minutes":60}`,
		RateLimitSMTPPerDomain:   `{"count":1000,"window_minutes":60}`,
		RateLimitAuthPerIP:       `{"count":10,"window_minutes":15}`,
		RateLimitIMAPPerUser:     `{"count":1000,"window_minutes":60}`,
		RateLimitCleanupInterval: 300,

		// Authentication security defaults
		AuthTOTPEnforced:            false,
		AuthBruteForceEnabled:       true,
		AuthBruteForceThreshold:     5,
		AuthBruteForceWindowMinutes: 15,
		AuthBruteForceBlockMinutes:  60,
		AuthIPBlacklistEnabled:      true,
		AuthCleanupInterval:         0,
	}

	// Create default template
	if err := s.repo.Create(defaultTemplate); err != nil {
		return fmt.Errorf("failed to create default domain template: %w", err)
	}

	return nil
}

// List retrieves all domains
func (s *DomainService) List(ctx context.Context, offset, limit int) ([]*domain.Domain, error) {
	return s.repo.List(offset, limit)
}

// Create creates a new domain
func (s *DomainService) Create(ctx context.Context, domain *domain.Domain) error {
	return s.repo.Create(domain)
}

// GetByID retrieves a domain by ID
func (s *DomainService) GetByID(ctx context.Context, id int64) (*domain.Domain, error) {
	return s.repo.GetByID(id)
}

// Update updates an existing domain
func (s *DomainService) Update(ctx context.Context, domain *domain.Domain) error {
	return s.repo.Update(domain)
}

// Delete deletes a domain
func (s *DomainService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(id)
}

// GetDKIMConfig retrieves DKIM configuration for a domain
func (s *DomainService) GetDKIMConfig(domainName string) (*domain.DKIMConfig, error) {
	domain, err := s.repo.GetByName(domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return &domain.DKIMConfig{
		Domain:     domain.Name,
		Selector:   domain.DKIMSelector,
		PrivateKey: []byte(domain.DKIMPrivateKey),
		PublicKey:  domain.DKIMPublicKey,
	}, nil
}
