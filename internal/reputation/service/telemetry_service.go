package service

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// TelemetryService handles telemetry data collection and aggregation
type TelemetryService struct {
	eventsRepo repository.EventsRepository
	scoresRepo repository.ScoresRepository
	logger     *zap.Logger
}

// NewTelemetryService creates a new telemetry service
func NewTelemetryService(
	eventsRepo repository.EventsRepository,
	scoresRepo repository.ScoresRepository,
	logger *zap.Logger,
) *TelemetryService {
	return &TelemetryService{
		eventsRepo: eventsRepo,
		scoresRepo: scoresRepo,
		logger:     logger,
	}
}

// RecordDelivery records a successful email delivery
func (s *TelemetryService) RecordDelivery(ctx context.Context, domainName, recipientDomain, ip string) error {
	event := &domain.SendingEvent{
		Timestamp:       time.Now().Unix(),
		Domain:          domainName,
		RecipientDomain: recipientDomain,
		EventType:       domain.EventDelivery,
		IPAddress:       ip,
		Metadata:        make(map[string]interface{}),
	}

	if err := s.eventsRepo.RecordEvent(ctx, event); err != nil {
		s.logger.Error("Failed to record delivery event",
			zap.String("domain", domainName),
			zap.String("recipient_domain", recipientDomain),
			zap.Error(err),
		)
		return fmt.Errorf("failed to record delivery: %w", err)
	}

	s.logger.Debug("Recorded delivery event",
		zap.String("domain", domainName),
		zap.String("recipient_domain", recipientDomain),
	)

	return nil
}

// RecordBounce records an email bounce
func (s *TelemetryService) RecordBounce(
	ctx context.Context,
	domainName, recipientDomain, ip string,
	bounceType, statusCode, response string,
) error {
	event := &domain.SendingEvent{
		Timestamp:          time.Now().Unix(),
		Domain:             domainName,
		RecipientDomain:    recipientDomain,
		EventType:          domain.EventBounce,
		BounceType:         &bounceType,
		EnhancedStatusCode: &statusCode,
		SMTPResponse:       &response,
		IPAddress:          ip,
		Metadata:           make(map[string]interface{}),
	}

	if err := s.eventsRepo.RecordEvent(ctx, event); err != nil {
		s.logger.Error("Failed to record bounce event",
			zap.String("domain", domainName),
			zap.String("bounce_type", bounceType),
			zap.Error(err),
		)
		return fmt.Errorf("failed to record bounce: %w", err)
	}

	s.logger.Warn("Recorded bounce event",
		zap.String("domain", domainName),
		zap.String("bounce_type", bounceType),
		zap.String("status_code", statusCode),
	)

	return nil
}

// RecordComplaint records a spam complaint
func (s *TelemetryService) RecordComplaint(ctx context.Context, domainName, recipientDomain string) error {
	event := &domain.SendingEvent{
		Timestamp:       time.Now().Unix(),
		Domain:          domainName,
		RecipientDomain: recipientDomain,
		EventType:       domain.EventComplaint,
		IPAddress:       "", // Complaints may not have IP context
		Metadata:        make(map[string]interface{}),
	}

	if err := s.eventsRepo.RecordEvent(ctx, event); err != nil {
		s.logger.Error("Failed to record complaint event",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		return fmt.Errorf("failed to record complaint: %w", err)
	}

	s.logger.Warn("Recorded complaint event",
		zap.String("domain", domainName),
		zap.String("recipient_domain", recipientDomain),
	)

	return nil
}

// CalculateReputationScore calculates and updates the reputation score for a domain
func (s *TelemetryService) CalculateReputationScore(ctx context.Context, domainName string) (*domain.ReputationScore, error) {
	now := time.Now().Unix()
	window24h := now - (24 * 60 * 60) // 24 hours ago

	// Get event counts for the last 24 hours
	counts, err := s.eventsRepo.GetEventCountsByType(ctx, domainName, window24h, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get event counts: %w", err)
	}

	deliveries := counts[string(domain.EventDelivery)]
	bounces := counts[string(domain.EventBounce)]
	complaints := counts[string(domain.EventComplaint)]
	total := deliveries + bounces + complaints

	// Calculate rates
	var deliveryRate, bounceRate, complaintRate float64
	if total > 0 {
		deliveryRate = (float64(deliveries) / float64(total)) * 100
		bounceRate = (float64(bounces) / float64(total)) * 100
		complaintRate = (float64(complaints) / float64(total)) * 100
	}

	// Calculate reputation score (0-100)
	// Base score starts at 100 and is reduced by:
	// - High bounce rate (>5% significantly impacts)
	// - Any complaint rate (very sensitive)
	// - Low delivery rate
	reputationScore := 100

	// Bounce rate penalties
	if bounceRate > 10.0 {
		reputationScore -= 30
	} else if bounceRate > 5.0 {
		reputationScore -= 15
	} else if bounceRate > 2.0 {
		reputationScore -= 5
	}

	// Complaint rate penalties (very strict)
	if complaintRate > 0.1 {
		reputationScore -= 40 // Severe penalty
	} else if complaintRate > 0.05 {
		reputationScore -= 20
	} else if complaintRate > 0.01 {
		reputationScore -= 10
	}

	// Delivery rate bonuses
	if deliveryRate > 95.0 {
		reputationScore += 0 // No bonus, this is expected
	} else if deliveryRate > 90.0 {
		reputationScore -= 5
	} else if deliveryRate > 80.0 {
		reputationScore -= 15
	} else {
		reputationScore -= 25
	}

	// Ensure score is in valid range
	if reputationScore < 0 {
		reputationScore = 0
	}
	if reputationScore > 100 {
		reputationScore = 100
	}

	// Get existing score to preserve circuit breaker and warm-up state
	existing, err := s.scoresRepo.GetReputationScore(ctx, domainName)
	if err != nil {
		// If the score doesn't exist yet, that's OK - we'll create a new one
		// Only propagate unexpected errors
		s.logger.Debug("No existing score found, will create new one",
			zap.String("domain", domainName),
		)
		existing = nil
	}

	score := &domain.ReputationScore{
		Domain:          domainName,
		ReputationScore: reputationScore,
		ComplaintRate:   complaintRate,
		BounceRate:      bounceRate,
		DeliveryRate:    deliveryRate,
		LastUpdated:     now,
	}

	// Preserve circuit breaker and warm-up state if they exist
	if existing != nil {
		score.CircuitBreakerActive = existing.CircuitBreakerActive
		score.CircuitBreakerReason = existing.CircuitBreakerReason
		score.WarmUpActive = existing.WarmUpActive
		score.WarmUpDay = existing.WarmUpDay
	}

	if err := s.scoresRepo.UpdateReputationScore(ctx, score); err != nil {
		return nil, fmt.Errorf("failed to update reputation score: %w", err)
	}

	s.logger.Info("Calculated reputation score",
		zap.String("domain", domainName),
		zap.Int("score", reputationScore),
		zap.Float64("delivery_rate", deliveryRate),
		zap.Float64("bounce_rate", bounceRate),
		zap.Float64("complaint_rate", complaintRate),
	)

	return score, nil
}

// CleanupOldData removes telemetry data older than the retention period
func (s *TelemetryService) CleanupOldData(ctx context.Context) error {
	// Default to 90 days retention
	retentionDays := 90
	cutoff := time.Now().AddDate(0, 0, -retentionDays).Unix()

	if err := s.eventsRepo.CleanupOldEvents(ctx, cutoff); err != nil {
		s.logger.Error("Failed to cleanup old events", zap.Error(err))
		return fmt.Errorf("failed to cleanup old data: %w", err)
	}

	s.logger.Info("Cleaned up old telemetry data",
		zap.Int("retention_days", retentionDays),
		zap.Int64("cutoff_timestamp", cutoff),
	)

	return nil
}

// CalculateAllScores calculates reputation scores for all domains
func (s *TelemetryService) CalculateAllScores(ctx context.Context) error {
	// Get all domains that have scores
	scores, err := s.scoresRepo.ListAllScores(ctx)
	if err != nil {
		return fmt.Errorf("failed to list scores: %w", err)
	}

	for _, score := range scores {
		if _, err := s.CalculateReputationScore(ctx, score.Domain); err != nil {
			s.logger.Error("Failed to calculate score for domain",
				zap.String("domain", score.Domain),
				zap.Error(err),
			)
			// Continue with other domains even if one fails
		}
	}

	return nil
}
