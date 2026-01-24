package phishing

// AI-powered phishing detection service provides comprehensive email security analysis

type PhishingDetectionServiceInterface interface {
	AnalyzeMessage(ctx context.Context, message *domain.Message) (*PhishingResult, error)
}
