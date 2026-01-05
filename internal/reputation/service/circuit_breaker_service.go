package service

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// Circuit breaker thresholds
const (
	ComplaintRateThreshold = 0.1  // 0.1% complaint rate triggers circuit breaker
	BounceRateThreshold    = 10.0 // 10% bounce rate triggers circuit breaker
	BlockDetectionWindow   = 3    // 3 consecutive blocks from major providers
)

// Major email providers for block detection
var MajorProviders = []string{
	"gmail.com",
	"googlemail.com",
	"outlook.com",
	"hotmail.com",
	"live.com",
	"yahoo.com",
	"aol.com",
	"icloud.com",
}

// CircuitBreakerService manages automatic domain pause/resume based on reputation
type CircuitBreakerService struct {
	eventsRepo  repository.EventsRepository
	scoresRepo  repository.ScoresRepository
	circuitRepo repository.CircuitBreakerRepository
	telemetry   *TelemetryService
	logger      *zap.Logger
}

// NewCircuitBreakerService creates a new circuit breaker service
func NewCircuitBreakerService(
	eventsRepo repository.EventsRepository,
	scoresRepo repository.ScoresRepository,
	circuitRepo repository.CircuitBreakerRepository,
	telemetry *TelemetryService,
	logger *zap.Logger,
) *CircuitBreakerService {
	return &CircuitBreakerService{
		eventsRepo:  eventsRepo,
		scoresRepo:  scoresRepo,
		circuitRepo: circuitRepo,
		telemetry:   telemetry,
		logger:      logger,
	}
}

// CheckAndTrigger evaluates all circuit breaker thresholds and triggers if needed
// Should be called periodically (every 15 minutes)
func (s *CircuitBreakerService) CheckAndTrigger(ctx context.Context) error {
	s.logger.Debug("running circuit breaker threshold checks")

	// Get all domains with reputation scores
	scores, err := s.scoresRepo.ListAllScores(ctx)
	if err != nil {
		return fmt.Errorf("failed to list reputation scores: %w", err)
	}

	for _, score := range scores {
		// Skip domains already in circuit breaker state
		if score.CircuitBreakerActive {
			continue
		}

		// Check all threshold conditions
		triggered, triggerType, value, threshold := s.evaluateThresholds(score)
		if triggered {
			s.logger.Warn("circuit breaker threshold exceeded",
				zap.String("domain", score.Domain),
				zap.String("trigger_type", string(triggerType)),
				zap.Float64("value", value),
				zap.Float64("threshold", threshold),
			)

			if err := s.triggerCircuitBreaker(ctx, score.Domain, triggerType, value, threshold); err != nil {
				s.logger.Error("failed to trigger circuit breaker",
					zap.String("domain", score.Domain),
					zap.Error(err),
				)
				continue
			}
		}
	}

	return nil
}

// evaluateThresholds checks all circuit breaker conditions
// Returns: triggered, triggerType, actualValue, threshold
func (s *CircuitBreakerService) evaluateThresholds(score *domain.ReputationScore) (bool, domain.TriggerType, float64, float64) {
	// 1. Check complaint rate (most critical)
	if score.ComplaintRate > ComplaintRateThreshold {
		return true, domain.TriggerComplaint, score.ComplaintRate, ComplaintRateThreshold
	}

	// 2. Check bounce rate
	if score.BounceRate > BounceRateThreshold {
		return true, domain.TriggerBounce, score.BounceRate, BounceRateThreshold
	}

	// 3. Check for major provider blocks
	// This requires looking at recent events - implement separately
	blocked, blockRate := s.checkMajorProviderBlocks(score.Domain)
	if blocked {
		return true, domain.TriggerBlock, blockRate, float64(BlockDetectionWindow)
	}

	return false, "", 0, 0
}

// checkMajorProviderBlocks detects repeated blocks from major email providers
// Returns true if 3+ consecutive blocks detected in recent events
func (s *CircuitBreakerService) checkMajorProviderBlocks(domainName string) (bool, float64) {
	// Get recent events for this domain (last 24 hours)
	ctx := context.Background()
	now := time.Now().Unix()
	startTime := now - 86400 // 24 hours ago

	events, err := s.eventsRepo.GetEventsInWindow(ctx, domainName, startTime, now)
	if err != nil {
		s.logger.Error("failed to get events for block detection",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		return false, 0
	}

	// Track consecutive blocks per provider
	providerBlocks := make(map[string]int)
	consecutiveBlocks := 0

	for _, event := range events {
		// Only check bounce/defer events to major providers
		if event.EventType != domain.EventBounce && event.EventType != domain.EventDefer {
			continue
		}

		// Check if recipient is from a major provider
		isMajor := false
		for _, provider := range MajorProviders {
			if event.RecipientDomain == provider {
				isMajor = true
				break
			}
		}

		if !isMajor {
			continue
		}

		// Check SMTP response for block indicators (4xx/5xx codes)
		if event.SMTPResponse != nil {
			response := *event.SMTPResponse
			// Look for block-related codes (421, 450, 451, 550, 551, 554)
			if len(response) >= 3 {
				code := response[:3]
				if code == "421" || code == "450" || code == "451" ||
					code == "550" || code == "551" || code == "554" {
					providerBlocks[event.RecipientDomain]++
					consecutiveBlocks++

					// If we hit 3+ blocks from same provider or overall, trigger
					if providerBlocks[event.RecipientDomain] >= BlockDetectionWindow ||
						consecutiveBlocks >= BlockDetectionWindow {
						return true, float64(consecutiveBlocks)
					}
				}
			}
		}
	}

	return false, 0
}

// triggerCircuitBreaker activates the circuit breaker for a domain
func (s *CircuitBreakerService) triggerCircuitBreaker(
	ctx context.Context,
	domainName string,
	triggerType domain.TriggerType,
	value float64,
	threshold float64,
) error {
	// Record the circuit breaker event
	event := &domain.CircuitBreakerEvent{
		Domain:       domainName,
		TriggerType:  triggerType,
		TriggerValue: value,
		Threshold:    threshold,
		PausedAt:     time.Now().Unix(),
	}
	if err := s.circuitRepo.RecordPause(ctx, event); err != nil {
		return fmt.Errorf("failed to record circuit breaker pause: %w", err)
	}

	// Update reputation score to reflect circuit breaker state
	score, err := s.scoresRepo.GetReputationScore(ctx, domainName)
	if err != nil {
		return fmt.Errorf("failed to get reputation score: %w", err)
	}

	score.CircuitBreakerActive = true
	score.CircuitBreakerReason = fmt.Sprintf("%s threshold exceeded: %.2f > %.2f",
		triggerType, value, threshold)

	if err := s.scoresRepo.UpdateReputationScore(ctx, score); err != nil {
		return fmt.Errorf("failed to update reputation score: %w", err)
	}

	s.logger.Warn("circuit breaker triggered",
		zap.String("domain", domainName),
		zap.String("trigger_type", string(triggerType)),
		zap.String("reason", score.CircuitBreakerReason),
	)

	return nil
}

// AutoResume attempts to resume domains that have been paused
// Uses exponential backoff: 1h → 2h → 4h → 8h
// Should be called periodically (every hour)
func (s *CircuitBreakerService) AutoResume(ctx context.Context) error {
	s.logger.Debug("checking for domains eligible for auto-resume")

	// Get all active circuit breakers
	breakers, err := s.circuitRepo.GetActiveBreakers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active circuit breakers: %w", err)
	}

	for _, breaker := range breakers {
		// Calculate backoff duration
		history, err := s.circuitRepo.GetBreakerHistory(ctx, breaker.Domain, 10)
		if err != nil {
			s.logger.Error("failed to get circuit breaker history",
				zap.String("domain", breaker.Domain),
				zap.Error(err),
			)
			continue
		}

		// Count consecutive pauses (failures to resume)
		consecutivePauses := s.countConsecutivePauses(history)

		// Exponential backoff: 1h, 2h, 4h, 8h (max)
		backoffHours := 1
		for i := 0; i < consecutivePauses && backoffHours < 8; i++ {
			backoffHours *= 2
		}

		// Check if enough time has passed since pause
		pausedAt := breaker.PausedAt
		backoffSeconds := int64(backoffHours * 3600)
		resumeTime := pausedAt + backoffSeconds
		now := time.Now().Unix()

		if now >= resumeTime {
			// Attempt to resume
			if err := s.attemptResume(ctx, breaker.Domain); err != nil {
				s.logger.Error("failed to resume domain",
					zap.String("domain", breaker.Domain),
					zap.Error(err),
				)
				continue
			}

			s.logger.Info("domain auto-resumed from circuit breaker",
				zap.String("domain", breaker.Domain),
				zap.Int("backoff_hours", backoffHours),
			)
		} else {
			timeRemaining := resumeTime - now
			s.logger.Debug("domain not yet eligible for resume",
				zap.String("domain", breaker.Domain),
				zap.Int64("seconds_remaining", timeRemaining),
			)
		}
	}

	return nil
}

// countConsecutivePauses counts how many times domain was re-paused after resume attempts
func (s *CircuitBreakerService) countConsecutivePauses(history []*domain.CircuitBreakerEvent) int {
	count := 0
	for i := len(history) - 1; i >= 0; i-- {
		// If we find a resume, stop counting
		if history[i].ResumedAt != nil {
			break
		}
		count++
	}
	return count
}

// attemptResume tries to resume a paused domain
func (s *CircuitBreakerService) attemptResume(ctx context.Context, domainName string) error {
	// Recalculate reputation score to see if conditions have improved
	score, err := s.telemetry.CalculateReputationScore(ctx, domainName)
	if err != nil {
		return fmt.Errorf("failed to recalculate reputation: %w", err)
	}

	// Check if conditions are still bad
	triggered, triggerType, value, threshold := s.evaluateThresholds(score)
	if triggered {
		// Conditions haven't improved, don't resume
		s.logger.Warn("resume attempt failed: conditions still bad",
			zap.String("domain", domainName),
			zap.String("trigger_type", string(triggerType)),
			zap.Float64("value", value),
			zap.Float64("threshold", threshold),
		)
		// Record new pause event
		return s.triggerCircuitBreaker(ctx, domainName, triggerType, value, threshold)
	}

	// Conditions have improved, resume sending
	if err := s.circuitRepo.RecordResume(ctx, domainName, true, "auto-resumed: conditions improved"); err != nil {
		return fmt.Errorf("failed to record resume: %w", err)
	}

	// Update reputation score
	score.CircuitBreakerActive = false
	score.CircuitBreakerReason = ""
	if err := s.scoresRepo.UpdateReputationScore(ctx, score); err != nil {
		return fmt.Errorf("failed to update reputation score: %w", err)
	}

	return nil
}

// ManualResume allows admin to manually resume a paused domain
// Bypasses automatic checks
func (s *CircuitBreakerService) ManualResume(ctx context.Context, domainName string, adminNotes string) error {
	// Record resume with admin notes
	if err := s.circuitRepo.RecordResume(ctx, domainName, false, adminNotes); err != nil {
		return fmt.Errorf("failed to record manual resume: %w", err)
	}

	// Update reputation score
	score, err := s.scoresRepo.GetReputationScore(ctx, domainName)
	if err != nil {
		return fmt.Errorf("failed to get reputation score: %w", err)
	}

	score.CircuitBreakerActive = false
	score.CircuitBreakerReason = ""
	if err := s.scoresRepo.UpdateReputationScore(ctx, score); err != nil {
		return fmt.Errorf("failed to update reputation score: %w", err)
	}

	s.logger.Info("domain manually resumed",
		zap.String("domain", domainName),
		zap.String("admin_notes", adminNotes),
	)

	return nil
}
