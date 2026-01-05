package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type providerRateLimitsRepository struct {
	db *sql.DB
}

// NewProviderRateLimitsRepository creates a new SQLite provider rate limits repository
func NewProviderRateLimitsRepository(db *sql.DB) repository.ProviderRateLimitsRepository {
	return &providerRateLimitsRepository{db: db}
}

// Get retrieves rate limit for a domain and provider
func (r *providerRateLimitsRepository) Get(ctx context.Context, domainName string, provider domain.MailProvider) (*domain.ProviderRateLimit, error) {
	query := `
		SELECT
			id, domain, provider, max_hourly_rate, max_daily_rate,
			current_hour_count, current_day_count, hour_reset_at, day_reset_at,
			circuit_breaker_active, last_updated
		FROM provider_rate_limits
		WHERE domain = ? AND provider = ?
	`

	limit := &domain.ProviderRateLimit{}
	err := r.db.QueryRowContext(ctx, query, domainName, string(provider)).Scan(
		&limit.ID,
		&limit.Domain,
		&limit.Provider,
		&limit.MaxHourlyRate,
		&limit.MaxDailyRate,
		&limit.CurrentHourCount,
		&limit.CurrentDayCount,
		&limit.HourResetAt,
		&limit.DayResetAt,
		&limit.CircuitBreakerActive,
		&limit.LastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rate limit not found for domain %s and provider %s", domainName, provider)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get provider rate limit: %w", err)
	}

	return limit, nil
}

// CreateOrUpdate creates or updates rate limit
func (r *providerRateLimitsRepository) CreateOrUpdate(ctx context.Context, limit *domain.ProviderRateLimit) error {
	query := `
		INSERT INTO provider_rate_limits (
			domain, provider, max_hourly_rate, max_daily_rate,
			current_hour_count, current_day_count, hour_reset_at, day_reset_at,
			circuit_breaker_active, last_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(domain, provider) DO UPDATE SET
			max_hourly_rate = excluded.max_hourly_rate,
			max_daily_rate = excluded.max_daily_rate,
			current_hour_count = excluded.current_hour_count,
			current_day_count = excluded.current_day_count,
			hour_reset_at = excluded.hour_reset_at,
			day_reset_at = excluded.day_reset_at,
			circuit_breaker_active = excluded.circuit_breaker_active,
			last_updated = excluded.last_updated
	`

	result, err := r.db.ExecContext(ctx, query,
		limit.Domain,
		string(limit.Provider),
		limit.MaxHourlyRate,
		limit.MaxDailyRate,
		limit.CurrentHourCount,
		limit.CurrentDayCount,
		limit.HourResetAt,
		limit.DayResetAt,
		limit.CircuitBreakerActive,
		limit.LastUpdated,
	)
	if err != nil {
		return fmt.Errorf("failed to create/update provider rate limit: %w", err)
	}

	if limit.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		limit.ID = id
	}

	return nil
}

// IncrementHourly increments hourly counter
func (r *providerRateLimitsRepository) IncrementHourly(ctx context.Context, domainName string, provider domain.MailProvider, count int) error {
	query := `
		UPDATE provider_rate_limits
		SET current_hour_count = current_hour_count + ?,
		    last_updated = ?
		WHERE domain = ? AND provider = ?
	`

	_, err := r.db.ExecContext(ctx, query, count, time.Now().Unix(), domainName, string(provider))
	if err != nil {
		return fmt.Errorf("failed to increment hourly counter: %w", err)
	}

	return nil
}

// IncrementDaily increments daily counter
func (r *providerRateLimitsRepository) IncrementDaily(ctx context.Context, domainName string, provider domain.MailProvider, count int) error {
	query := `
		UPDATE provider_rate_limits
		SET current_day_count = current_day_count + ?,
		    last_updated = ?
		WHERE domain = ? AND provider = ?
	`

	_, err := r.db.ExecContext(ctx, query, count, time.Now().Unix(), domainName, string(provider))
	if err != nil {
		return fmt.Errorf("failed to increment daily counter: %w", err)
	}

	return nil
}

// ResetHourly resets hourly counter
func (r *providerRateLimitsRepository) ResetHourly(ctx context.Context, domainName string, provider domain.MailProvider, newResetTime int64) error {
	query := `
		UPDATE provider_rate_limits
		SET current_hour_count = 0,
		    hour_reset_at = ?,
		    last_updated = ?
		WHERE domain = ? AND provider = ?
	`

	_, err := r.db.ExecContext(ctx, query, newResetTime, time.Now().Unix(), domainName, string(provider))
	if err != nil {
		return fmt.Errorf("failed to reset hourly counter: %w", err)
	}

	return nil
}

// ResetDaily resets daily counter
func (r *providerRateLimitsRepository) ResetDaily(ctx context.Context, domainName string, provider domain.MailProvider, newResetTime int64) error {
	query := `
		UPDATE provider_rate_limits
		SET current_day_count = 0,
		    day_reset_at = ?,
		    last_updated = ?
		WHERE domain = ? AND provider = ?
	`

	_, err := r.db.ExecContext(ctx, query, newResetTime, time.Now().Unix(), domainName, string(provider))
	if err != nil {
		return fmt.Errorf("failed to reset daily counter: %w", err)
	}

	return nil
}

// ListByDomain returns all provider limits for a domain
func (r *providerRateLimitsRepository) ListByDomain(ctx context.Context, domainName string) ([]*domain.ProviderRateLimit, error) {
	query := `
		SELECT
			id, domain, provider, max_hourly_rate, max_daily_rate,
			current_hour_count, current_day_count, hour_reset_at, day_reset_at,
			circuit_breaker_active, last_updated
		FROM provider_rate_limits
		WHERE domain = ?
		ORDER BY provider ASC
	`

	rows, err := r.db.QueryContext(ctx, query, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to list provider rate limits: %w", err)
	}
	defer rows.Close()

	var limits []*domain.ProviderRateLimit
	for rows.Next() {
		limit := &domain.ProviderRateLimit{}
		err := rows.Scan(
			&limit.ID,
			&limit.Domain,
			&limit.Provider,
			&limit.MaxHourlyRate,
			&limit.MaxDailyRate,
			&limit.CurrentHourCount,
			&limit.CurrentDayCount,
			&limit.HourResetAt,
			&limit.DayResetAt,
			&limit.CircuitBreakerActive,
			&limit.LastUpdated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan provider rate limit: %w", err)
		}
		limits = append(limits, limit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating provider rate limits: %w", err)
	}

	return limits, nil
}

// SetCircuitBreaker activates/deactivates circuit breaker
func (r *providerRateLimitsRepository) SetCircuitBreaker(ctx context.Context, domainName string, provider domain.MailProvider, active bool) error {
	query := `
		UPDATE provider_rate_limits
		SET circuit_breaker_active = ?,
		    last_updated = ?
		WHERE domain = ? AND provider = ?
	`

	_, err := r.db.ExecContext(ctx, query, active, time.Now().Unix(), domainName, string(provider))
	if err != nil {
		return fmt.Errorf("failed to set circuit breaker: %w", err)
	}

	return nil
}

// GetLimitsByDomain returns all provider limits for a domain (alias for ListByDomain)
func (r *providerRateLimitsRepository) GetLimitsByDomain(ctx context.Context, domainName string) ([]*domain.ProviderRateLimit, error) {
	return r.ListByDomain(ctx, domainName)
}

// GetAllLimits returns all provider limits across all domains
func (r *providerRateLimitsRepository) GetAllLimits(ctx context.Context) ([]*domain.ProviderRateLimit, error) {
	query := `
		SELECT
			id, domain, provider, max_hourly_rate, max_daily_rate,
			current_hour_count, current_day_count, hour_reset_at, day_reset_at,
			circuit_breaker_active, last_updated
		FROM provider_rate_limits
		ORDER BY domain ASC, provider ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all provider rate limits: %w", err)
	}
	defer rows.Close()

	var limits []*domain.ProviderRateLimit
	for rows.Next() {
		limit := &domain.ProviderRateLimit{}
		err := rows.Scan(
			&limit.ID,
			&limit.Domain,
			&limit.Provider,
			&limit.MaxHourlyRate,
			&limit.MaxDailyRate,
			&limit.CurrentHourCount,
			&limit.CurrentDayCount,
			&limit.HourResetAt,
			&limit.DayResetAt,
			&limit.CircuitBreakerActive,
			&limit.LastUpdated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan provider rate limit: %w", err)
		}
		limits = append(limits, limit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating provider rate limits: %w", err)
	}

	return limits, nil
}

// GetLimitByID retrieves a rate limit by ID
func (r *providerRateLimitsRepository) GetLimitByID(ctx context.Context, id int64) (*domain.ProviderRateLimit, error) {
	query := `
		SELECT
			id, domain, provider, max_hourly_rate, max_daily_rate,
			current_hour_count, current_day_count, hour_reset_at, day_reset_at,
			circuit_breaker_active, last_updated
		FROM provider_rate_limits
		WHERE id = ?
	`

	limit := &domain.ProviderRateLimit{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&limit.ID,
		&limit.Domain,
		&limit.Provider,
		&limit.MaxHourlyRate,
		&limit.MaxDailyRate,
		&limit.CurrentHourCount,
		&limit.CurrentDayCount,
		&limit.HourResetAt,
		&limit.DayResetAt,
		&limit.CircuitBreakerActive,
		&limit.LastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rate limit not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get provider rate limit by id: %w", err)
	}

	return limit, nil
}

// UpdateLimit updates an existing rate limit
func (r *providerRateLimitsRepository) UpdateLimit(ctx context.Context, limit *domain.ProviderRateLimit) error {
	query := `
		UPDATE provider_rate_limits
		SET domain = ?,
		    provider = ?,
		    max_hourly_rate = ?,
		    max_daily_rate = ?,
		    current_hour_count = ?,
		    current_day_count = ?,
		    hour_reset_at = ?,
		    day_reset_at = ?,
		    circuit_breaker_active = ?,
		    last_updated = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		limit.Domain,
		string(limit.Provider),
		limit.MaxHourlyRate,
		limit.MaxDailyRate,
		limit.CurrentHourCount,
		limit.CurrentDayCount,
		limit.HourResetAt,
		limit.DayResetAt,
		limit.CircuitBreakerActive,
		time.Now().Unix(),
		limit.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update provider rate limit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rate limit found with id %d", limit.ID)
	}

	return nil
}
