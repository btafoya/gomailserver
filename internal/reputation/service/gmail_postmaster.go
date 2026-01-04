package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
	"google.golang.org/api/gmailpostmastertools/v1"
	"google.golang.org/api/option"
)

// GmailPostmasterService integrates with Gmail Postmaster Tools API

type GmailPostmasterService struct {
	service         *gmailpostmastertools.Service
	metricsRepo     repository.PostmasterMetricsRepository
	alertsRepo      repository.AlertsRepository
	logger          *zap.Logger
	serviceAccountKey string
}

func NewGmailPostmasterService(
	serviceAccountKey string,
	metricsRepo repository.PostmasterMetricsRepository,
	alertsRepo repository.AlertsRepository,
	logger *zap.Logger,
) (*GmailPostmasterService, error) {
	ctx := context.Background()

	// Create Gmail Postmaster Tools service with service account
	service, err := gmailpostmastertools.NewService(ctx, option.WithCredentialsFile(serviceAccountKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail Postmaster service: %w", err)
	}

	return &GmailPostmasterService{
		service:           service,
		metricsRepo:       metricsRepo,
		alertsRepo:        alertsRepo,
		logger:            logger,
		serviceAccountKey: serviceAccountKey,
	}, nil
}

// FetchDomainReputation retrieves reputation metrics for a domain
func (s *GmailPostmasterService) FetchDomainReputation(ctx context.Context, domainName string) (*domain.PostmasterMetrics, error) {
	// Gmail Postmaster Tools uses domain resource names like "domains/example.com"
	domainResource := fmt.Sprintf("domains/%s", domainName)

	// Fetch traffic stats (contains domain reputation)
	trafficStats, err := s.service.Domains.TrafficStats.List(domainResource).Context(ctx).Do()
	if err != nil {
		s.logger.Error("Failed to fetch Gmail Postmaster metrics",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch traffic stats: %w", err)
	}

	if len(trafficStats.TrafficStats) == 0 {
		s.logger.Warn("No traffic stats available for domain", zap.String("domain", domainName))
		return nil, nil
	}

	// Get most recent stats
	stats := trafficStats.TrafficStats[0]

	// Parse metric date
	metricDate, err := time.Parse("2006-01-02", stats.Name)
	if err != nil {
		s.logger.Warn("Failed to parse metric date", zap.Error(err))
		metricDate = time.Now()
	}

	// Convert to domain model
	metrics := &domain.PostmasterMetrics{
		Domain:         domainName,
		FetchedAt:      time.Now().Unix(),
		MetricDate:     metricDate.Unix(),
		DomainReputation: stats.DomainReputation,
		IPReputation:    stats.IpReputation,
		UserSpamReports: int(stats.UserReportedSpamRatio * 100), // Convert to count
	}

	// Calculate spam rate (user reported spam ratio)
	if stats.UserReportedSpamRatio != nil {
		metrics.SpamRate = *stats.UserReportedSpamRatio
	}

	// Calculate authentication rate
	if stats.SpfSuccessRatio != nil && stats.DkimSuccessRatio != nil && stats.DmarcSuccessRatio != nil {
		authRate := (*stats.SpfSuccessRatio + *stats.DkimSuccessRatio + *stats.DmarcSuccessRatio) / 3.0
		metrics.AuthenticationRate = authRate
	}

	// Calculate encryption rate (TLS)
	if stats.InboundEncryptionRatio != nil {
		metrics.EncryptionRate = *stats.InboundEncryptionRatio
	}

	// Store raw response for debugging
	rawBytes, _ := json.Marshal(stats)
	metrics.RawResponse = string(rawBytes)

	s.logger.Info("Fetched Gmail Postmaster metrics",
		zap.String("domain", domainName),
		zap.String("domain_reputation", metrics.DomainReputation),
		zap.Float64("spam_rate", metrics.SpamRate),
	)

	return metrics, nil
}

// SyncDomain syncs metrics for a single domain
func (s *GmailPostmasterService) SyncDomain(ctx context.Context, domainName string) error {
	metrics, err := s.FetchDomainReputation(ctx, domainName)
	if err != nil {
		return err
	}

	if metrics == nil {
		return nil // No data available
	}

	// Store metrics
	if err := s.metricsRepo.Create(ctx, metrics); err != nil {
		s.logger.Error("Failed to store Postmaster metrics", zap.Error(err))
		return fmt.Errorf("failed to store metrics: %w", err)
	}

	// Check for reputation degradation and create alerts
	if err := s.checkReputationAlerts(ctx, metrics); err != nil {
		s.logger.Error("Failed to check reputation alerts", zap.Error(err))
		// Don't fail sync on alert error
	}

	return nil
}

// SyncAll syncs metrics for all configured domains
func (s *GmailPostmasterService) SyncAll(ctx context.Context, domains []string) error {
	s.logger.Info("Starting Gmail Postmaster sync for all domains", zap.Int("count", len(domains)))

	errors := make([]error, 0)
	for _, domainName := range domains {
		if err := s.SyncDomain(ctx, domainName); err != nil {
			s.logger.Error("Failed to sync domain",
				zap.String("domain", domainName),
				zap.Error(err),
			)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("sync completed with %d errors", len(errors))
	}

	s.logger.Info("Gmail Postmaster sync completed successfully")
	return nil
}

// checkReputationAlerts checks for reputation issues and creates alerts
func (s *GmailPostmasterService) checkReputationAlerts(ctx context.Context, metrics *domain.PostmasterMetrics) error {
	// Check for bad domain reputation
	if metrics.DomainReputation == "BAD" || metrics.DomainReputation == "LOW" {
		alert := &domain.ReputationAlert{
			Domain:    metrics.Domain,
			AlertType: domain.AlertTypeExternalFeedback,
			Severity:  getSeverityForReputation(metrics.DomainReputation),
			Title:     fmt.Sprintf("Gmail Domain Reputation Degraded: %s", metrics.DomainReputation),
			Message:   fmt.Sprintf("Gmail Postmaster Tools reports domain reputation as %s", metrics.DomainReputation),
			Details: map[string]interface{}{
				"domain_reputation":    metrics.DomainReputation,
				"spam_rate":           metrics.SpamRate,
				"authentication_rate": metrics.AuthenticationRate,
				"metric_date":         metrics.MetricDate,
			},
			CreatedAt: time.Now().Unix(),
		}

		if err := s.alertsRepo.Create(ctx, alert); err != nil {
			return err
		}
	}

	// Check for high spam rate
	if metrics.SpamRate > 0.001 { // 0.1%
		severity := domain.SeverityMedium
		if metrics.SpamRate > 0.003 {
			severity = domain.SeverityHigh
		}
		if metrics.SpamRate > 0.01 {
			severity = domain.SeverityCritical
		}

		alert := &domain.ReputationAlert{
			Domain:    metrics.Domain,
			AlertType: domain.AlertTypeExternalFeedback,
			Severity:  severity,
			Title:     fmt.Sprintf("High Spam Rate Detected: %.2f%%", metrics.SpamRate*100),
			Message:   "Gmail users are reporting messages as spam at an elevated rate",
			Details: map[string]interface{}{
				"spam_rate":    metrics.SpamRate,
				"metric_date":  metrics.MetricDate,
				"threshold":    0.001,
			},
			CreatedAt: time.Now().Unix(),
		}

		if err := s.alertsRepo.Create(ctx, alert); err != nil {
			return err
		}
	}

	// Check for low authentication rate
	if metrics.AuthenticationRate < 0.95 {
		alert := &domain.ReputationAlert{
			Domain:    metrics.Domain,
			AlertType: domain.AlertTypeExternalFeedback,
			Severity:  domain.SeverityMedium,
			Title:     fmt.Sprintf("Low Authentication Rate: %.1f%%", metrics.AuthenticationRate*100),
			Message:   "SPF/DKIM/DMARC authentication rate is below recommended 95%",
			Details: map[string]interface{}{
				"authentication_rate": metrics.AuthenticationRate,
				"metric_date":         metrics.MetricDate,
				"threshold":           0.95,
			},
			CreatedAt: time.Now().Unix(),
		}

		if err := s.alertsRepo.Create(ctx, alert); err != nil {
			return err
		}
	}

	return nil
}

// GetLatestMetrics returns the latest metrics for a domain
func (s *GmailPostmasterService) GetLatestMetrics(ctx context.Context, domainName string) (*domain.PostmasterMetrics, error) {
	return s.metricsRepo.GetLatest(ctx, domainName)
}

// GetMetricsHistory returns metrics history for a domain
func (s *GmailPostmasterService) GetMetricsHistory(ctx context.Context, domainName string, days int) ([]*domain.PostmasterMetrics, error) {
	return s.metricsRepo.ListByDomain(ctx, domainName, days)
}

// GetReputationTrend returns domain reputation trend over time
func (s *GmailPostmasterService) GetReputationTrend(ctx context.Context, domainName string, days int) ([]string, error) {
	return s.metricsRepo.GetReputationTrend(ctx, domainName, days)
}

// Helper function to map reputation to severity
func getSeverityForReputation(reputation string) domain.AlertSeverity {
	switch reputation {
	case "BAD":
		return domain.SeverityCritical
	case "LOW":
		return domain.SeverityHigh
	case "MEDIUM":
		return domain.SeverityMedium
	default:
		return domain.SeverityLow
	}
}
