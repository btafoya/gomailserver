package antivirus

import (
	clamavPkg "github.com/btafoya/gomailserver/internal/clamav"
	"github.com/btafoya/gomailserver/internal/domain"
	"go.uber.org/zap"
)

// AntivirusServiceInterface defines methods needed by antivirus components
type AntivirusServiceInterface interface {
	GetAntivirusConfig(domainName string) (*domain.AntivirusConfig, error)
}

type Scanner struct {
	clamav        *ClamAV
	domainService AntivirusServiceInterface
	logger        *zap.Logger
}

func NewAntivirusScanner(clamav *ClamAV, domainService AntivirusServiceInterface, logger *zap.Logger) *Scanner {
	return &Scanner{
		clamav:        clamav,
		domainService: domainService,
		logger:        logger,
	}
}

type ScanAction string

const (
	ActionReject     ScanAction = "reject"
	ActionQuarantine ScanAction = "quarantine"
	ActionTag        ScanAction = "tag"
)

func (s *AntivirusScanner) ScanMessage(domainName string, message []byte) (*ScanResult, ScanAction, error) {
	// Get antivirus configuration for the domain
	config, err := s.domainService.GetAntivirusConfig(domainName)
	if err != nil {
		s.logger.Error("Failed to get antivirus config",
			zap.String("domain", domainName),
			zap.Error(err))
		return nil, ActionTag, err // Fail open with tag
	}

	result, err := s.clamav.Scan(message)
	if err != nil {
		s.logger.Error("ClamAV scan failed",
			zap.String("domain", domainName),
			zap.Error(err))
		return nil, ActionTag, err // Fail open with tag
	}

	if result.Clean {
		return result, "", nil
	}

	// Apply action based on domain policy
	return result, config.VirusAction, nil
}
