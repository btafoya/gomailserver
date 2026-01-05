package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"github.com/btafoya/gomailserver/internal/security/ratelimit"
	"go.uber.org/zap"
)

var (
	// ErrCircuitBreakerActive indicates the domain is currently paused
	ErrCircuitBreakerActive = errors.New("circuit breaker active: sending paused for this domain")
	// ErrWarmUpLimitExceeded indicates the warm-up volume cap has been reached
	ErrWarmUpLimitExceeded = errors.New("warm-up limit exceeded: daily volume cap reached")
)

// AdaptiveLimiter extends the base rate limiter with reputation awareness
type AdaptiveLimiter struct {
	baseLimiter    *ratelimit.Limiter
	scoresRepo     repository.ScoresRepository
	warmUpRepo     repository.WarmUpRepository
	circuitRepo    repository.CircuitBreakerRepository
	logger         *zap.Logger
}

// NewAdaptiveLimiter creates a new adaptive rate limiter
func NewAdaptiveLimiter(
	baseLimiter *ratelimit.Limiter,
	scoresRepo repository.ScoresRepository,
	warmUpRepo repository.WarmUpRepository,
	circuitRepo repository.CircuitBreakerRepository,
	logger *zap.Logger,
) *AdaptiveLimiter {
	return &AdaptiveLimiter{
		baseLimiter: baseLimiter,
		scoresRepo:  scoresRepo,
		warmUpRepo:  warmUpRepo,
		circuitRepo: circuitRepo,
		logger:      logger,
	}
}

// GetLimit returns the effective rate limit for a domain based on reputation
// Returns the maximum messages per hour allowed for this domain
func (l *AdaptiveLimiter) GetLimit(ctx context.Context, domain string) (int, error) {
	// Check circuit breaker first - highest priority
	score, err := l.scoresRepo.GetReputationScore(ctx, domain)
	if err != nil {
		// No reputation score yet - use base limit
		l.logger.Debug("no reputation score found, using base limit",
			zap.String("domain", domain),
		)
		return ratelimit.DefaultLimits["smtp_per_domain"].Count, nil
	}

	// If circuit breaker is active, return 0 (no sending allowed)
	if score.CircuitBreakerActive {
		l.logger.Warn("circuit breaker active for domain",
			zap.String("domain", domain),
			zap.String("reason", score.CircuitBreakerReason),
		)
		return 0, ErrCircuitBreakerActive
	}

	// Check warm-up schedule - second priority
	if score.WarmUpActive {
		schedule, err := l.warmUpRepo.GetSchedule(ctx, domain)
		if err != nil {
			l.logger.Error("failed to get warm-up schedule",
				zap.String("domain", domain),
				zap.Error(err),
			)
			// Fall through to reputation-based limit
		} else {
			// Get the max volume for current day
			if score.WarmUpDay > 0 && score.WarmUpDay <= len(schedule) {
				maxVolume := schedule[score.WarmUpDay-1].MaxVolume
				l.logger.Info("warm-up limit applied",
					zap.String("domain", domain),
					zap.Int("day", score.WarmUpDay),
					zap.Int("max_volume", maxVolume),
				)
				return maxVolume, nil
			}
		}
	}

	// Apply reputation-based adjustment to base limit
	baseLimit := ratelimit.DefaultLimits["smtp_per_domain"].Count

	// Reputation score is 0-100, convert to multiplier 0.0-1.0
	// A score of 100 = full base limit
	// A score of 50 = 50% of base limit
	// A score of 0 = 0% of base limit (effectively paused)
	adjustment := float64(score.ReputationScore) / 100.0
	effectiveLimit := int(float64(baseLimit) * adjustment)

	// Ensure minimum limit of 10/hour for domains with very low reputation
	// This prevents complete blocking but still restricts poor senders
	if effectiveLimit < 10 && score.ReputationScore > 0 {
		effectiveLimit = 10
	}

	l.logger.Debug("reputation-adjusted limit calculated",
		zap.String("domain", domain),
		zap.Int("reputation_score", score.ReputationScore),
		zap.Int("base_limit", baseLimit),
		zap.Int("effective_limit", effectiveLimit),
	)

	return effectiveLimit, nil
}

// CheckDomain verifies if a domain can send based on reputation, circuit breaker, and warm-up
// This is the main entry point for SMTP backend to use
func (l *AdaptiveLimiter) CheckDomain(ctx context.Context, domain string) (bool, error) {
	limit, err := l.GetLimit(ctx, domain)
	if err != nil {
		// Circuit breaker or other critical error
		return false, err
	}

	// If limit is 0, domain cannot send
	if limit == 0 {
		return false, fmt.Errorf("sending limit is zero for domain %s", domain)
	}

	// Use base limiter to check actual usage against the effective limit
	// For now, delegate to base limiter's domain check
	// In future, we could implement custom tracking here
	return l.baseLimiter.Check("smtp_per_domain", domain)
}

// CheckWarmUpVolume verifies if a domain has exceeded its daily warm-up volume
// Returns current volume and max allowed
func (l *AdaptiveLimiter) CheckWarmUpVolume(ctx context.Context, domain string) (current int, max int, exceeded bool, err error) {
	score, err := l.scoresRepo.GetReputationScore(ctx, domain)
	if err != nil {
		return 0, 0, false, fmt.Errorf("failed to get reputation score: %w", err)
	}

	// Not in warm-up mode
	if !score.WarmUpActive {
		return 0, 0, false, nil
	}

	schedule, err := l.warmUpRepo.GetSchedule(ctx, domain)
	if err != nil {
		return 0, 0, false, fmt.Errorf("failed to get warm-up schedule: %w", err)
	}

	if score.WarmUpDay < 1 || score.WarmUpDay > len(schedule) {
		return 0, 0, false, fmt.Errorf("invalid warm-up day: %d", score.WarmUpDay)
	}

	daySchedule := schedule[score.WarmUpDay-1]
	exceeded = daySchedule.ActualVolume >= daySchedule.MaxVolume

	return daySchedule.ActualVolume, daySchedule.MaxVolume, exceeded, nil
}

// RecordSend increments the warm-up volume counter for a domain
// Should be called after successful message send during warm-up
func (l *AdaptiveLimiter) RecordSend(ctx context.Context, domain string) error {
	score, err := l.scoresRepo.GetReputationScore(ctx, domain)
	if err != nil || !score.WarmUpActive {
		// Not in warm-up, nothing to record
		return nil
	}

	return l.warmUpRepo.IncrementDayVolume(ctx, domain, score.WarmUpDay, 1)
}
