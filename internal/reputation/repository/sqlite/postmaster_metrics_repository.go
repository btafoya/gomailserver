package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type postmasterMetricsRepository struct {
	db *sql.DB
}

// NewPostmasterMetricsRepository creates a new SQLite Gmail Postmaster metrics repository
func NewPostmasterMetricsRepository(db *sql.DB) repository.PostmasterMetricsRepository {
	return &postmasterMetricsRepository{db: db}
}

// Create stores new Gmail Postmaster metrics
func (r *postmasterMetricsRepository) Create(ctx context.Context, metrics *domain.PostmasterMetrics) error {
	query := `
		INSERT INTO postmaster_metrics (
			domain, fetched_at, metric_date, domain_reputation, spam_rate,
			ip_reputation, authentication_rate, encryption_rate,
			user_spam_reports, raw_response
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		metrics.Domain,
		metrics.FetchedAt,
		metrics.MetricDate,
		metrics.DomainReputation,
		metrics.SpamRate,
		metrics.IPReputation,
		metrics.AuthenticationRate,
		metrics.EncryptionRate,
		metrics.UserSpamReports,
		metrics.RawResponse,
	)
	if err != nil {
		return fmt.Errorf("failed to create Postmaster metrics: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	metrics.ID = id

	return nil
}

// GetLatest returns the latest metrics for a domain
func (r *postmasterMetricsRepository) GetLatest(ctx context.Context, domainName string) (*domain.PostmasterMetrics, error) {
	query := `
		SELECT
			id, domain, fetched_at, metric_date, domain_reputation, spam_rate,
			ip_reputation, authentication_rate, encryption_rate,
			user_spam_reports, raw_response
		FROM postmaster_metrics
		WHERE domain = ?
		ORDER BY metric_date DESC
		LIMIT 1
	`

	metrics := &domain.PostmasterMetrics{}
	err := r.db.QueryRowContext(ctx, query, domainName).Scan(
		&metrics.ID,
		&metrics.Domain,
		&metrics.FetchedAt,
		&metrics.MetricDate,
		&metrics.DomainReputation,
		&metrics.SpamRate,
		&metrics.IPReputation,
		&metrics.AuthenticationRate,
		&metrics.EncryptionRate,
		&metrics.UserSpamReports,
		&metrics.RawResponse,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no Postmaster metrics found for domain %s", domainName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest Postmaster metrics: %w", err)
	}

	return metrics, nil
}

// ListByDomain returns metrics history for a domain
func (r *postmasterMetricsRepository) ListByDomain(ctx context.Context, domainName string, days int) ([]*domain.PostmasterMetrics, error) {
	query := `
		SELECT
			id, domain, fetched_at, metric_date, domain_reputation, spam_rate,
			ip_reputation, authentication_rate, encryption_rate,
			user_spam_reports, raw_response
		FROM postmaster_metrics
		WHERE domain = ? AND metric_date >= ?
		ORDER BY metric_date DESC
	`

	startTime := time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix()

	rows, err := r.db.QueryContext(ctx, query, domainName, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list Postmaster metrics: %w", err)
	}
	defer rows.Close()

	var metricsList []*domain.PostmasterMetrics
	for rows.Next() {
		metrics := &domain.PostmasterMetrics{}
		err := rows.Scan(
			&metrics.ID,
			&metrics.Domain,
			&metrics.FetchedAt,
			&metrics.MetricDate,
			&metrics.DomainReputation,
			&metrics.SpamRate,
			&metrics.IPReputation,
			&metrics.AuthenticationRate,
			&metrics.EncryptionRate,
			&metrics.UserSpamReports,
			&metrics.RawResponse,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan Postmaster metrics: %w", err)
		}
		metricsList = append(metricsList, metrics)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating Postmaster metrics: %w", err)
	}

	return metricsList, nil
}

// GetReputationTrend returns domain reputation trend
func (r *postmasterMetricsRepository) GetReputationTrend(ctx context.Context, domainName string, days int) ([]string, error) {
	query := `
		SELECT domain_reputation
		FROM postmaster_metrics
		WHERE domain = ? AND metric_date >= ?
		ORDER BY metric_date ASC
	`

	startTime := time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix()

	rows, err := r.db.QueryContext(ctx, query, domainName, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get reputation trend: %w", err)
	}
	defer rows.Close()

	var trend []string
	for rows.Next() {
		var reputation string
		err := rows.Scan(&reputation)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reputation: %w", err)
		}
		trend = append(trend, reputation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reputation trend: %w", err)
	}

	return trend, nil
}
