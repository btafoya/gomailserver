package service

import (
	"context"
	"fmt"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/security/dkim"
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

// CreateDomainFromTemplate creates a new domain from template with security defaults
func (s *DomainService) CreateDomainFromTemplate(domainName string) (*domain.Domain, error) {
	// Create default template with recommended security settings
	defaultTemplate := &domain.Domain{
		Name:           domainName,
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

	if err := s.repo.Create(defaultTemplate); err != nil {
		return nil, fmt.Errorf("failed to create domain from template: %w", err)
	}

	return defaultTemplate, nil
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
	dom, err := s.repo.GetByName(domainName)
	if err != nil {
		return nil, err
	}

	return &domain.DKIMConfig{
		Domain:     dom.Name,
		Selector:   dom.DKIMSelector,
		PrivateKey: []byte(dom.DKIMPrivateKey),
		PublicKey:  dom.DKIMPublicKey,
	}, nil
}

// DKIMKeyResult contains the generated DKIM keys and DNS record
type DKIMKeyResult struct {
	Selector   string `json:"selector"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	DNSRecord  string `json:"dns_record"`
	DNSName    string `json:"dns_name"`
}

// GenerateDKIMKeys generates new DKIM keys for a domain
func (s *DomainService) GenerateDKIMKeys(ctx context.Context, id int64, keyType string, keySize int) (*DKIMKeyResult, error) {
	// Get the domain
	dom, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("domain not found: %w", err)
	}

	var keyPair *dkim.KeyPair

	// Generate keys based on type
	switch keyType {
	case "ed25519":
		keyPair, err = dkim.GenerateEd25519KeyPair()
	case "rsa":
		fallthrough
	default:
		if keySize <= 0 {
			keySize = 2048
		}
		keyPair, err = dkim.GenerateRSAKeyPair(keySize)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate DKIM keys: %w", err)
	}

	// Update domain with new keys
	dom.DKIMSelector = keyPair.Selector
	dom.DKIMPrivateKey = keyPair.PrivateKey
	dom.DKIMPublicKey = keyPair.PublicKey
	dom.DKIMSigningEnabled = true

	if err := s.repo.Update(dom); err != nil {
		return nil, fmt.Errorf("failed to update domain with DKIM keys: %w", err)
	}

	s.logger.Info("DKIM keys generated",
		zap.Int64("domain_id", id),
		zap.String("domain", dom.Name),
		zap.String("selector", keyPair.Selector),
		zap.String("key_type", keyType),
	)

	return &DKIMKeyResult{
		Selector:   keyPair.Selector,
		PrivateKey: keyPair.PrivateKey,
		PublicKey:  keyPair.PublicKey,
		DNSRecord:  keyPair.DNSRecord(),
		DNSName:    fmt.Sprintf("%s._domainkey.%s", keyPair.Selector, dom.Name),
	}, nil
}

// GetAntivirusConfig retrieves antivirus configuration for a domain
func (s *DomainService) GetAntivirusConfig(domainName string) (*domain.AntivirusConfig, error) {
	dom, err := s.repo.GetByName(domainName)
	if err != nil {
		return nil, err
	}

	return &domain.AntivirusConfig{
		VirusAction: dom.ClamAVVirusAction,
	}, nil
}

// GetDefaultTemplate retrieves the default domain template
func (s *DomainService) GetDefaultTemplate() (*domain.Domain, error) {
	return s.repo.GetByName(DefaultTemplateDomainName)
}

// UpdateDefaultTemplate updates the default domain template
func (s *DomainService) UpdateDefaultTemplate(updates *domain.Domain) error {
	template, err := s.repo.GetByName(DefaultTemplateDomainName)
	if err != nil {
		return fmt.Errorf("default template not found: %w", err)
	}

	// Apply updates to template fields
	if updates.MaxUsers > 0 {
		template.MaxUsers = updates.MaxUsers
	}
	if updates.MaxMailboxSize > 0 {
		template.MaxMailboxSize = updates.MaxMailboxSize
	}
	if updates.DefaultQuota > 0 {
		template.DefaultQuota = updates.DefaultQuota
	}

	// DKIM settings
	if updates.DKIMSelector != "" {
		template.DKIMSelector = updates.DKIMSelector
	}
	template.DKIMSigningEnabled = updates.DKIMSigningEnabled
	template.DKIMVerifyEnabled = updates.DKIMVerifyEnabled
	if updates.DKIMKeySize > 0 {
		template.DKIMKeySize = updates.DKIMKeySize
	}
	if updates.DKIMKeyType != "" {
		template.DKIMKeyType = updates.DKIMKeyType
	}
	if updates.DKIMHeadersToSign != "" {
		template.DKIMHeadersToSign = updates.DKIMHeadersToSign
	}

	// SPF settings
	if updates.SPFRecord != "" {
		template.SPFRecord = updates.SPFRecord
	}
	template.SPFEnabled = updates.SPFEnabled
	if updates.SPFDNSServer != "" {
		template.SPFDNSServer = updates.SPFDNSServer
	}
	if updates.SPFDNSTimeout > 0 {
		template.SPFDNSTimeout = updates.SPFDNSTimeout
	}
	if updates.SPFMaxLookups > 0 {
		template.SPFMaxLookups = updates.SPFMaxLookups
	}
	if updates.SPFFailAction != "" {
		template.SPFFailAction = updates.SPFFailAction
	}
	if updates.SPFSoftFailAction != "" {
		template.SPFSoftFailAction = updates.SPFSoftFailAction
	}

	// DMARC settings
	if updates.DMARCPolicy != "" {
		template.DMARCPolicy = updates.DMARCPolicy
	}
	template.DMARCEnabled = updates.DMARCEnabled
	if updates.DMARCDNSServer != "" {
		template.DMARCDNSServer = updates.DMARCDNSServer
	}
	if updates.DMARCDNSTimeout > 0 {
		template.DMARCDNSTimeout = updates.DMARCDNSTimeout
	}
	template.DMARCReportEnabled = updates.DMARCReportEnabled
	if updates.DMARCReportEmail != "" {
		template.DMARCReportEmail = updates.DMARCReportEmail
	}

	// ClamAV settings
	template.ClamAVEnabled = updates.ClamAVEnabled
	if updates.ClamAVMaxScanSize > 0 {
		template.ClamAVMaxScanSize = updates.ClamAVMaxScanSize
	}
	if updates.ClamAVVirusAction != "" {
		template.ClamAVVirusAction = updates.ClamAVVirusAction
	}
	if updates.ClamAVFailAction != "" {
		template.ClamAVFailAction = updates.ClamAVFailAction
	}

	// Spam settings
	template.SpamEnabled = updates.SpamEnabled
	if updates.SpamRejectScore > 0 {
		template.SpamRejectScore = updates.SpamRejectScore
	}
	if updates.SpamQuarantineScore > 0 {
		template.SpamQuarantineScore = updates.SpamQuarantineScore
	}
	template.SpamLearningEnabled = updates.SpamLearningEnabled

	// Greylisting settings
	template.GreylistEnabled = updates.GreylistEnabled
	if updates.GreylistDelayMinutes > 0 {
		template.GreylistDelayMinutes = updates.GreylistDelayMinutes
	}
	if updates.GreylistExpiryDays > 0 {
		template.GreylistExpiryDays = updates.GreylistExpiryDays
	}
	if updates.GreylistCleanupInterval > 0 {
		template.GreylistCleanupInterval = updates.GreylistCleanupInterval
	}
	if updates.GreylistWhitelistAfter > 0 {
		template.GreylistWhitelistAfter = updates.GreylistWhitelistAfter
	}

	// Rate limiting settings
	template.RateLimitEnabled = updates.RateLimitEnabled
	if updates.RateLimitSMTPPerIP != "" {
		template.RateLimitSMTPPerIP = updates.RateLimitSMTPPerIP
	}
	if updates.RateLimitSMTPPerUser != "" {
		template.RateLimitSMTPPerUser = updates.RateLimitSMTPPerUser
	}
	if updates.RateLimitSMTPPerDomain != "" {
		template.RateLimitSMTPPerDomain = updates.RateLimitSMTPPerDomain
	}
	if updates.RateLimitAuthPerIP != "" {
		template.RateLimitAuthPerIP = updates.RateLimitAuthPerIP
	}
	if updates.RateLimitIMAPPerUser != "" {
		template.RateLimitIMAPPerUser = updates.RateLimitIMAPPerUser
	}
	if updates.RateLimitCleanupInterval > 0 {
		template.RateLimitCleanupInterval = updates.RateLimitCleanupInterval
	}

	// Auth settings
	template.AuthTOTPEnforced = updates.AuthTOTPEnforced
	template.AuthBruteForceEnabled = updates.AuthBruteForceEnabled
	if updates.AuthBruteForceThreshold > 0 {
		template.AuthBruteForceThreshold = updates.AuthBruteForceThreshold
	}
	if updates.AuthBruteForceWindowMinutes > 0 {
		template.AuthBruteForceWindowMinutes = updates.AuthBruteForceWindowMinutes
	}
	if updates.AuthBruteForceBlockMinutes > 0 {
		template.AuthBruteForceBlockMinutes = updates.AuthBruteForceBlockMinutes
	}
	template.AuthIPBlacklistEnabled = updates.AuthIPBlacklistEnabled

	return s.repo.Update(template)
}
