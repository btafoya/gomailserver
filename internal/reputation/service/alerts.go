package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// AlertsService manages reputation alerts

type AlertsService struct {
	alertsRepo repository.AlertsRepository
	logger     *zap.Logger
}

func NewAlertsService(
	alertsRepo repository.AlertsRepository,
	logger *zap.Logger,
) *AlertsService {
	return &AlertsService{
		alertsRepo: alertsRepo,
		logger:     logger,
	}
}

// CreateAlert creates a new reputation alert
func (s *AlertsService) CreateAlert(ctx context.Context, alert *domain.ReputationAlert) error {
	// Set creation time if not set
	if alert.CreatedAt == 0 {
		alert.CreatedAt = time.Now().Unix()
	}

	// Validate alert
	if err := s.validateAlert(alert); err != nil {
		return fmt.Errorf("invalid alert: %w", err)
	}

	// Store alert
	if err := s.alertsRepo.Create(ctx, alert); err != nil {
		s.logger.Error("Failed to create alert", zap.Error(err))
		return err
	}

	s.logger.Info("Created reputation alert",
		zap.String("domain", alert.Domain),
		zap.String("type", string(alert.AlertType)),
		zap.String("severity", string(alert.Severity)),
		zap.String("title", alert.Title),
	)

	return nil
}

// validateAlert validates alert fields
func (s *AlertsService) validateAlert(alert *domain.ReputationAlert) error {
	if alert.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if alert.AlertType == "" {
		return fmt.Errorf("alert type is required")
	}
	if alert.Severity == "" {
		return fmt.Errorf("severity is required")
	}
	if alert.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// GetUnacknowledged returns all unacknowledged alerts
func (s *AlertsService) GetUnacknowledged(ctx context.Context, limit int) ([]*domain.ReputationAlert, error) {
	return s.alertsRepo.ListUnacknowledged(ctx, limit)
}

// GetByDomain returns alerts for a specific domain
func (s *AlertsService) GetByDomain(ctx context.Context, domainName string, limit, offset int) ([]*domain.ReputationAlert, error) {
	return s.alertsRepo.ListByDomain(ctx, domainName, limit, offset)
}

// GetBySeverity returns alerts by severity level
func (s *AlertsService) GetBySeverity(ctx context.Context, severity domain.AlertSeverity, limit int) ([]*domain.ReputationAlert, error) {
	return s.alertsRepo.ListBySeverity(ctx, severity, limit)
}

// Acknowledge marks an alert as acknowledged
func (s *AlertsService) Acknowledge(ctx context.Context, id int64, acknowledgedBy string) error {
	if err := s.alertsRepo.Acknowledge(ctx, id, acknowledgedBy); err != nil {
		s.logger.Error("Failed to acknowledge alert", zap.Int64("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Acknowledged alert", zap.Int64("id", id), zap.String("by", acknowledgedBy))
	return nil
}

// Resolve marks an alert as resolved
func (s *AlertsService) Resolve(ctx context.Context, id int64) error {
	if err := s.alertsRepo.Resolve(ctx, id); err != nil {
		s.logger.Error("Failed to resolve alert", zap.Int64("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Resolved alert", zap.Int64("id", id))
	return nil
}

// GetUnacknowledgedCount returns count of unacknowledged alerts
func (s *AlertsService) GetUnacknowledgedCount(ctx context.Context) (int, error) {
	return s.alertsRepo.GetUnacknowledgedCount(ctx)
}

// GetUnacknowledgedCountByDomain returns count of unacknowledged alerts for a domain
func (s *AlertsService) GetUnacknowledgedCountByDomain(ctx context.Context, domainName string) (int, error) {
	return s.alertsRepo.GetUnacknowledgedCountByDomain(ctx, domainName)
}

// GetByID retrieves a specific alert
func (s *AlertsService) GetByID(ctx context.Context, id int64) (*domain.ReputationAlert, error) {
	return s.alertsRepo.GetByID(ctx, id)
}

// CreateDNSFailureAlert creates an alert for DNS validation failures
func (s *AlertsService) CreateDNSFailureAlert(ctx context.Context, domainName string, checkType string, details map[string]interface{}) error {
	alert := &domain.ReputationAlert{
		Domain:    domainName,
		AlertType: domain.AlertTypeDNSFailure,
		Severity:  domain.SeverityHigh,
		Title:     fmt.Sprintf("DNS Validation Failed: %s", checkType),
		Message:   fmt.Sprintf("DNS check for %s failed for domain %s", checkType, domainName),
		Details:   details,
		CreatedAt: time.Now().Unix(),
	}

	return s.CreateAlert(ctx, alert)
}

// CreateScoreDropAlert creates an alert for reputation score drops
func (s *AlertsService) CreateScoreDropAlert(ctx context.Context, domainName string, oldScore, newScore int, dropPercentage float64) error {
	severity := domain.SeverityMedium
	if dropPercentage > 30 {
		severity = domain.SeverityHigh
	}
	if dropPercentage > 50 {
		severity = domain.SeverityCritical
	}

	alert := &domain.ReputationAlert{
		Domain:    domainName,
		AlertType: domain.AlertTypeScoreDrop,
		Severity:  severity,
		Title:     fmt.Sprintf("Reputation Score Dropped by %.1f%%", dropPercentage),
		Message:   fmt.Sprintf("Domain reputation score dropped from %d to %d (%.1f%% decrease)", oldScore, newScore, dropPercentage),
		Details: map[string]interface{}{
			"old_score":       oldScore,
			"new_score":       newScore,
			"drop_percentage": dropPercentage,
		},
		CreatedAt: time.Now().Unix(),
	}

	return s.CreateAlert(ctx, alert)
}

// CreateCircuitBreakerAlert creates an alert for circuit breaker triggers
func (s *AlertsService) CreateCircuitBreakerAlert(ctx context.Context, domainName string, triggerType, reason string) error {
	alert := &domain.ReputationAlert{
		Domain:    domainName,
		AlertType: domain.AlertTypeCircuitBreaker,
		Severity:  domain.SeverityCritical,
		Title:     fmt.Sprintf("Circuit Breaker Triggered: %s", triggerType),
		Message:   fmt.Sprintf("Sending paused for domain %s due to %s", domainName, reason),
		Details: map[string]interface{}{
			"trigger_type": triggerType,
			"reason":       reason,
		},
		CreatedAt: time.Now().Unix(),
	}

	return s.CreateAlert(ctx, alert)
}

// ExportAlertsJSON exports alerts to JSON format
func (s *AlertsService) ExportAlertsJSON(ctx context.Context, alerts []*domain.ReputationAlert) ([]byte, error) {
	return json.MarshalIndent(alerts, "", "  ")
}

// AcknowledgeAlert marks an alert as acknowledged (alias for Acknowledge)
func (s *AlertsService) AcknowledgeAlert(ctx context.Context, id int64, acknowledgedBy string) error {
	return s.Acknowledge(ctx, id, acknowledgedBy)
}

// ResolveAlert marks an alert as resolved (alias for Resolve)
func (s *AlertsService) ResolveAlert(ctx context.Context, id int64) error {
	return s.Resolve(ctx, id)
}
