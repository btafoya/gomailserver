package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type sndsMetricsRepository struct {
	db *sql.DB
}

// NewSNDSMetricsRepository creates a new SQLite Microsoft SNDS metrics repository
func NewSNDSMetricsRepository(db *sql.DB) repository.SNDSMetricsRepository {
	return &sndsMetricsRepository{db: db}
}

// Create stores new Microsoft SNDS metrics
func (r *sndsMetricsRepository) Create(ctx context.Context, metrics *domain.SNDSMetrics) error {
	query := `
		INSERT INTO snds_metrics (
			ip_address, fetched_at, metric_date, spam_trap_hits,
			complaint_rate, filter_level, message_count, raw_response
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		metrics.IPAddress,
		metrics.FetchedAt,
		metrics.MetricDate,
		metrics.SpamTrapHits,
		metrics.ComplaintRate,
		metrics.FilterLevel,
		metrics.MessageCount,
		metrics.RawResponse,
	)
	if err != nil {
		return fmt.Errorf("failed to create SNDS metrics: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	metrics.ID = id

	return nil
}

// GetLatest returns the latest metrics for an IP
func (r *sndsMetricsRepository) GetLatest(ctx context.Context, ipAddress string) (*domain.SNDSMetrics, error) {
	query := `
		SELECT
			id, ip_address, fetched_at, metric_date, spam_trap_hits,
			complaint_rate, filter_level, message_count, raw_response
		FROM snds_metrics
		WHERE ip_address = ?
		ORDER BY metric_date DESC
		LIMIT 1
	`

	metrics := &domain.SNDSMetrics{}
	err := r.db.QueryRowContext(ctx, query, ipAddress).Scan(
		&metrics.ID,
		&metrics.IPAddress,
		&metrics.FetchedAt,
		&metrics.MetricDate,
		&metrics.SpamTrapHits,
		&metrics.ComplaintRate,
		&metrics.FilterLevel,
		&metrics.MessageCount,
		&metrics.RawResponse,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no SNDS metrics found for IP %s", ipAddress)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest SNDS metrics: %w", err)
	}

	return metrics, nil
}

// ListByIP returns metrics history for an IP
func (r *sndsMetricsRepository) ListByIP(ctx context.Context, ipAddress string, days int) ([]*domain.SNDSMetrics, error) {
	query := `
		SELECT
			id, ip_address, fetched_at, metric_date, spam_trap_hits,
			complaint_rate, filter_level, message_count, raw_response
		FROM snds_metrics
		WHERE ip_address = ? AND metric_date >= ?
		ORDER BY metric_date DESC
	`

	startTime := time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix()

	rows, err := r.db.QueryContext(ctx, query, ipAddress, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list SNDS metrics: %w", err)
	}
	defer rows.Close()

	var metricsList []*domain.SNDSMetrics
	for rows.Next() {
		metrics := &domain.SNDSMetrics{}
		err := rows.Scan(
			&metrics.ID,
			&metrics.IPAddress,
			&metrics.FetchedAt,
			&metrics.MetricDate,
			&metrics.SpamTrapHits,
			&metrics.ComplaintRate,
			&metrics.FilterLevel,
			&metrics.MessageCount,
			&metrics.RawResponse,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan SNDS metrics: %w", err)
		}
		metricsList = append(metricsList, metrics)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating SNDS metrics: %w", err)
	}

	return metricsList, nil
}

// GetFilterLevelTrend returns filter level trend
func (r *sndsMetricsRepository) GetFilterLevelTrend(ctx context.Context, ipAddress string, days int) ([]string, error) {
	query := `
		SELECT filter_level
		FROM snds_metrics
		WHERE ip_address = ? AND metric_date >= ?
		ORDER BY metric_date ASC
	`

	startTime := time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix()

	rows, err := r.db.QueryContext(ctx, query, ipAddress, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get filter level trend: %w", err)
	}
	defer rows.Close()

	var trend []string
	for rows.Next() {
		var filterLevel string
		err := rows.Scan(&filterLevel)
		if err != nil {
			return nil, fmt.Errorf("failed to scan filter level: %w", err)
		}
		trend = append(trend, filterLevel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating filter level trend: %w", err)
	}

	return trend, nil
}
