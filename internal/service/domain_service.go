package service

import (
	"context"
	"fmt"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

const DefaultTemplateDomainName = "_default"

// DomainService handles domain business logic including default templates
type DomainService struct {
	repo repository.DomainRepository
}

// NewDomainService creates a new domain service
func NewDomainService(repo repository.DomainRepository) *DomainService {
	return &DomainService{repo: repo}
}

// EnsureDefaultTemplate creates the default domain template if it doesn't exist
// This should be called during server initialization
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
		MaxUsers:       0, // unlimited
		MaxMailboxSize: 0, // unlimited
		DefaultQuota:   1073741824, // 1GB

		// DKIM defaults
		DKIMSigningEnabled: true,
		DKIMVerifyEnabled:  true,
		DKIMKeySize:        2048,
		DKIMKeyType:        "rsa",
		DKIMHeadersToSign:  `["From","To","Subject","Date","Message-ID","MIME-Version","Content-Type"]`,

		// SPF defaults
		SPFEnabled:        true,
		SPFDNSServer:      "8.8.8.8:53",
		SPFDNSTimeout:     5,
		SPFMaxLookups:     10,
		SPFFailAction:     "reject",
		SPFSoftFailAction: "accept",

		// DMARC defaults
		DMARCEnabled:       true,
		DMARCDNSServer:     "8.8.8.8:53",
		DMARCDNSTimeout:    5,
		DMARCReportEnabled: false,

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
		AuthCleanupInterval:         3600,
	}

	if err := s.repo.Create(defaultTemplate); err != nil {
		return fmt.Errorf("failed to create default template: %w", err)
	}

	return nil
}

// GetDefaultTemplate retrieves the default domain template
func (s *DomainService) GetDefaultTemplate() (*domain.Domain, error) {
	template, err := s.repo.GetByName(DefaultTemplateDomainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get default template: %w", err)
	}
	return template, nil
}

// CreateDomainFromTemplate creates a new domain using the default template for security settings
func (s *DomainService) CreateDomainFromTemplate(name string) (*domain.Domain, error) {
	// Get default template
	template, err := s.GetDefaultTemplate()
	if err != nil {
		return nil, fmt.Errorf("failed to get default template: %w", err)
	}

	// Create new domain with template's security settings
	newDomain := &domain.Domain{
		Name:           name,
		Status:         "active",
		MaxUsers:       template.MaxUsers,
		MaxMailboxSize: template.MaxMailboxSize,
		DefaultQuota:   template.DefaultQuota,

		// Copy all security settings from template
		DKIMSigningEnabled: template.DKIMSigningEnabled,
		DKIMVerifyEnabled:  template.DKIMVerifyEnabled,
		DKIMKeySize:        template.DKIMKeySize,
		DKIMKeyType:        template.DKIMKeyType,
		DKIMHeadersToSign:  template.DKIMHeadersToSign,

		SPFEnabled:        template.SPFEnabled,
		SPFDNSServer:      template.SPFDNSServer,
		SPFDNSTimeout:     template.SPFDNSTimeout,
		SPFMaxLookups:     template.SPFMaxLookups,
		SPFFailAction:     template.SPFFailAction,
		SPFSoftFailAction: template.SPFSoftFailAction,

		DMARCEnabled:       template.DMARCEnabled,
		DMARCDNSServer:     template.DMARCDNSServer,
		DMARCDNSTimeout:    template.DMARCDNSTimeout,
		DMARCReportEnabled: template.DMARCReportEnabled,

		ClamAVEnabled:     template.ClamAVEnabled,
		ClamAVMaxScanSize: template.ClamAVMaxScanSize,
		ClamAVVirusAction: template.ClamAVVirusAction,
		ClamAVFailAction:  template.ClamAVFailAction,

		SpamEnabled:         template.SpamEnabled,
		SpamRejectScore:     template.SpamRejectScore,
		SpamQuarantineScore: template.SpamQuarantineScore,
		SpamLearningEnabled: template.SpamLearningEnabled,

		GreylistEnabled:         template.GreylistEnabled,
		GreylistDelayMinutes:    template.GreylistDelayMinutes,
		GreylistExpiryDays:      template.GreylistExpiryDays,
		GreylistCleanupInterval: template.GreylistCleanupInterval,
		GreylistWhitelistAfter:  template.GreylistWhitelistAfter,

		RateLimitEnabled:         template.RateLimitEnabled,
		RateLimitSMTPPerIP:       template.RateLimitSMTPPerIP,
		RateLimitSMTPPerUser:     template.RateLimitSMTPPerUser,
		RateLimitSMTPPerDomain:   template.RateLimitSMTPPerDomain,
		RateLimitAuthPerIP:       template.RateLimitAuthPerIP,
		RateLimitIMAPPerUser:     template.RateLimitIMAPPerUser,
		RateLimitCleanupInterval: template.RateLimitCleanupInterval,

		AuthTOTPEnforced:            template.AuthTOTPEnforced,
		AuthBruteForceEnabled:       template.AuthBruteForceEnabled,
		AuthBruteForceThreshold:     template.AuthBruteForceThreshold,
		AuthBruteForceWindowMinutes: template.AuthBruteForceWindowMinutes,
		AuthBruteForceBlockMinutes:  template.AuthBruteForceBlockMinutes,
		AuthIPBlacklistEnabled:      template.AuthIPBlacklistEnabled,
		AuthCleanupInterval:         template.AuthCleanupInterval,
	}

	if err := s.repo.Create(newDomain); err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	return newDomain, nil
}

// UpdateDefaultTemplate updates the default template settings
// This allows administrators to change defaults for new domains
func (s *DomainService) UpdateDefaultTemplate(updates *domain.Domain) error {
	template, err := s.GetDefaultTemplate()
	if err != nil {
		return err
	}

	// Update template with provided values
	// Copy all security settings (keeping ID and timestamps)
	updates.ID = template.ID
	updates.Name = DefaultTemplateDomainName // Ensure name doesn't change
	updates.CreatedAt = template.CreatedAt

	if err := s.repo.Update(updates); err != nil {
		return fmt.Errorf("failed to update default template: %w", err)
	}

	return nil
}

// List retrieves all domains with pagination
func (s *DomainService) List(ctx context.Context) ([]*domain.Domain, error) {
	// For now, return all domains without pagination
	// TODO: Add pagination support with offset/limit
	domains, err := s.repo.List(0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}
	return domains, nil
}

// Create creates a new domain
func (s *DomainService) Create(ctx context.Context, newDomain *domain.Domain) error {
	if err := s.repo.Create(newDomain); err != nil {
		return fmt.Errorf("failed to create domain: %w", err)
	}
	return nil
}

// GetByID retrieves a domain by ID
func (s *DomainService) GetByID(ctx context.Context, id int64) (*domain.Domain, error) {
	domain, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}
	return domain, nil
}

// Update updates an existing domain
func (s *DomainService) Update(ctx context.Context, domain *domain.Domain) error {
	if err := s.repo.Update(domain); err != nil {
		return fmt.Errorf("failed to update domain: %w", err)
	}
	return nil
}

// Delete deletes a domain
func (s *DomainService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}
	return nil
}
