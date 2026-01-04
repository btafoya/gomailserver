package service

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// ProviderRateLimitsService manages provider-specific rate limiting

type ProviderRateLimitsService struct {
	limitsRepo repository.ProviderRateLimitsRepository
	alertsRepo repository.AlertsRepository
	logger     *zap.Logger
}

func NewProviderRateLimitsService(
	limitsRepo repository.ProviderRateLimitsRepository,
	alertsRepo repository.AlertsRepository,
	logger *zap.Logger,
) *ProviderRateLimitsService {
	return &ProviderRateLimitsService{
		limitsRepo: limitsRepo,
		alertsRepo: alertsRepo,
		logger:     logger,
	}
}

// GetLimit retrieves rate limit for a domain and provider
func (s *ProviderRateLimitsService) GetLimit(ctx context.Context, domainName string, provider domain.MailProvider) (*domain.ProviderRateLimit, error) {
	limit, err := s.limitsRepo.Get(ctx, domainName, provider)
	if err != nil {
		return nil, err
	}

	// Check if resets are needed
	now := time.Now()
	if limit.ShouldResetHour(now) {
		if err := s.limitsRepo.ResetHourly(ctx, domainName, provider, now.Add(1*time.Hour).Unix()); err != nil {
			s.logger.Error("Failed to reset hourly counter", zap.Error(err))
		}
	}

	if limit.ShouldResetDay(now) {
		if err := s.limitsRepo.ResetDaily(ctx, domainName, provider, now.Add(24*time.Hour).Unix()); err != nil {
			s.logger.Error("Failed to reset daily counter", zap.Error(err))
		}
	}

	return limit, nil
}

// CheckLimit checks if sending is allowed under current rate limits
func (s *ProviderRateLimitsService) CheckLimit(ctx context.Context, domainName string, provider domain.MailProvider) (bool, error) {
	limit, err := s.GetLimit(ctx, domainName, provider)
	if err != nil {
		return false, err
	}

	// Check circuit breaker
	if limit.CircuitBreakerActive {
		return false, fmt.Errorf("circuit breaker active for %s to %s", domainName, provider)
	}

	// Check hourly limit
	if limit.IsAtHourlyLimit() {
		return false, nil
	}

	// Check daily limit
	if limit.IsAtDailyLimit() {
		return false, nil
	}

	return true, nil
}

// IncrementCount increments message count for a provider
func (s *ProviderRateLimitsService) IncrementCount(ctx context.Context, domainName string, provider domain.MailProvider, count int) error {
	// Increment hourly
	if err := s.limitsRepo.IncrementHourly(ctx, domainName, provider, count); err != nil {
		return err
	}

	// Increment daily
	if err := s.limitsRepo.IncrementDaily(ctx, domainName, provider, count); err != nil {
		return err
	}

	return nil
}

// SetCircuitBreaker activates or deactivates circuit breaker for a provider
func (s *ProviderRateLimitsService) SetCircuitBreaker(ctx context.Context, domainName string, provider domain.MailProvider, active bool) error {
	if err := s.limitsRepo.SetCircuitBreaker(ctx, domainName, provider, active); err != nil {
		return err
	}

	if active {
		s.logger.Warn("Circuit breaker activated",
			zap.String("domain", domainName),
			zap.String("provider", string(provider)),
		)
	} else {
		s.logger.Info("Circuit breaker deactivated",
			zap.String("domain", domainName),
			zap.String("provider", string(provider)),
		)
	}

	return nil
}

// CreateOrUpdateLimit creates or updates rate limit configuration
func (s *ProviderRateLimitsService) CreateOrUpdateLimit(ctx context.Context, limit *domain.ProviderRateLimit) error {
	// Set reset times if not set
	now := time.Now()
	if limit.HourResetAt == 0 {
		limit.HourResetAt = now.Add(1 * time.Hour).Unix()
	}
	if limit.DayResetAt == 0 {
		limit.DayResetAt = now.Add(24 * time.Hour).Unix()
	}
	if limit.LastUpdated == 0 {
		limit.LastUpdated = now.Unix()
	}

	return s.limitsRepo.CreateOrUpdate(ctx, limit)
}

// GetDomainLimits returns all provider limits for a domain
func (s *ProviderRateLimitsService) GetDomainLimits(ctx context.Context, domainName string) ([]*domain.ProviderRateLimit, error) {
	return s.limitsRepo.ListByDomain(ctx, domainName)
}

// InitializeDefaultLimits initializes default rate limits for a domain
func (s *ProviderRateLimitsService) InitializeDefaultLimits(ctx context.Context, domainName string) error {
	now := time.Now()

	// Gmail: Conservative limits
	gmailLimit := &domain.ProviderRateLimit{
		Domain:          domainName,
		Provider:        domain.ProviderGmail,
		MaxHourlyRate:   500,  // 500 emails/hour
		MaxDailyRate:    10000, // 10K emails/day
		HourResetAt:     now.Add(1 * time.Hour).Unix(),
		DayResetAt:      now.Add(24 * time.Hour).Unix(),
		LastUpdated:     now.Unix(),
	}

	// Outlook: Moderate limits
	outlookLimit := &domain.ProviderRateLimit{
		Domain:          domainName,
		Provider:        domain.ProviderOutlook,
		MaxHourlyRate:   300,  // 300 emails/hour
		MaxDailyRate:    7500,  // 7.5K emails/day
		HourResetAt:     now.Add(1 * time.Hour).Unix(),
		DayResetAt:      now.Add(24 * time.Hour).Unix(),
		LastUpdated:     now.Unix(),
	}

	// Yahoo: Conservative limits
	yahooLimit := &domain.ProviderRateLimit{
		Domain:          domainName,
		Provider:        domain.ProviderYahoo,
		MaxHourlyRate:   200,  // 200 emails/hour
		MaxDailyRate:    5000,  // 5K emails/day
		HourResetAt:     now.Add(1 * time.Hour).Unix(),
		DayResetAt:      now.Add(24 * time.Hour).Unix(),
		LastUpdated:     now.Unix(),
	}

	// Generic: Higher limits
	genericLimit := &domain.ProviderRateLimit{
		Domain:          domainName,
		Provider:        domain.ProviderGeneric,
		MaxHourlyRate:   1000, // 1K emails/hour
		MaxDailyRate:    20000, // 20K emails/day
		HourResetAt:     now.Add(1 * time.Hour).Unix(),
		DayResetAt:      now.Add(24 * time.Hour).Unix(),
		LastUpdated:     now.Unix(),
	}

	// Create all limits
	limits := []*domain.ProviderRateLimit{gmailLimit, outlookLimit, yahooLimit, genericLimit}
	for _, limit := range limits {
		if err := s.CreateOrUpdateLimit(ctx, limit); err != nil {
			s.logger.Error("Failed to initialize provider limit",
				zap.String("domain", domainName),
				zap.String("provider", string(limit.Provider)),
				zap.Error(err),
			)
			return err
		}
	}

	s.logger.Info("Initialized default provider rate limits", zap.String("domain", domainName))
	return nil
}
