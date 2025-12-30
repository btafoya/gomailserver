package antivirus

import (
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/domain"
)

// This is a placeholder for the actual service
type placeholderDomainService struct{}

func (s *placeholderDomainService) GetAntivirusConfig(domainName string) (*domain.AntivirusConfig, error) {
	// In a real implementation, this would fetch the antivirus configuration
	// for the given domain from the database.
	return &domain.AntivirusConfig{
		VirusAction: string(ActionQuarantine),
	}, nil
}

type Scanner struct {
	clamav        *ClamAV
	domainService *placeholderDomainService
	logger        *zap.Logger
}

type ScanAction string

const (
	ActionReject     ScanAction = "reject"
	ActionQuarantine ScanAction = "quarantine"
	ActionTag        ScanAction = "tag"
)

func NewScanner(clamav *ClamAV, logger *config.Logger) *Scanner {
	return &Scanner{
		clamav:        clamav,
		domainService: &placeholderDomainService{},
		logger:        logger.L,
	}
}

func (s *Scanner) ScanMessage(domainName string, message []byte) (*ScanResult, ScanAction, error) {
	result, err := s.clamav.Scan(message)
	if err != nil {
		s.logger.Error("ClamAV scan failed", zap.Error(err))
		return nil, ActionTag, err // Fail open with tag
	}

	if result.Clean {
		return result, "", nil
	}

	// Get domain configuration
	cfg, _ := s.domainService.GetAntivirusConfig(domainName)
	action := ScanAction(cfg.VirusAction)
	if action == "" {
		action = ActionQuarantine // Default
	}

	s.logger.Warn("Virus detected",
		zap.String("virus", result.Virus),
		zap.String("domain", domainName),
		zap.String("action", string(action)))

	return result, action, nil
}
