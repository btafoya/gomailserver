package phishing

import (
	"context"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/security/phishing"
)

// AI-powered phishing detection service provides comprehensive email security analysis

type PhishingDetectionServiceInterface interface {
	AnalyzeMessage(ctx context.Context, message *domain.Message) (*phishing.PhishingResult, error)
}
