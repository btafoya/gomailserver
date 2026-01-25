package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type historicalScoresRepository struct {
	db *sql.DB
}

// NewHistoricalScoresRepository creates a new SQLite historical scores repository
func NewHistoricalScoresRepository(db *sql.DB) repository.HistoricalScoresRepository {
	return &historicalScoresRepository{db: db}
}

// RecordScore stores a historical score snapshot
func (r *historicalScoresRepository) RecordScore(ctx context.Context, score *domain.HistoricalScore) error {
	query := `
		INSERT INTO historical_scores (domain, reputation_score, complaint_rate, bounce_rate, delivery_rate, recorded_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		score.Domain,
		score.ReputationScore,
		score.ComplaintRate,
		score.BounceRate,
		score.DeliveryRate,
		score.RecordedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to record historical score: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	score.ID = id

	return nil
}

// GetScoresInRange retrieves historical scores for a domain within a time range
func (r *historicalScoresRepository) GetScoresInRange(ctx context.Context, domainName string, startTime, endTime int64) ([]*domain.HistoricalScore, error) {
	query := `
		SELECT id, domain, reputation_score, complaint_rate, bounce_rate, delivery_rate, recorded_at
		FROM historical_scores
		WHERE domain = ? AND recorded_at BETWEEN ? AND ?
		ORDER BY recorded_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical scores: %w", err)
	}
	defer rows.Close()

	var scores []*domain.HistoricalScore
	for rows.Next() {
		score := &domain.HistoricalScore{}
		if err := rows.Scan(
			&score.ID,
			&score.Domain,
			&score.ReputationScore,
			&score.ComplaintRate,
			&score.BounceRate,
			&score.DeliveryRate,
			&score.RecordedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan historical score: %w", err)
		}
		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating historical scores: %w", err)
	}

	return scores, nil
}

// GetLatestScores retrieves the most recent N historical scores for a domain
func (r *historicalScoresRepository) GetLatestScores(ctx context.Context, domainName string, limit int) ([]*domain.HistoricalScore, error) {
	query := `
		SELECT id, domain, reputation_score, complaint_rate, bounce_rate, delivery_rate, recorded_at
		FROM historical_scores
		WHERE domain = ?
		ORDER BY recorded_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest scores: %w", err)
	}
	defer rows.Close()

	var scores []*domain.HistoricalScore
	for rows.Next() {
		score := &domain.HistoricalScore{}
		if err := rows.Scan(
			&score.ID,
			&score.Domain,
			&score.ReputationScore,
			&score.ComplaintRate,
			&score.BounceRate,
			&score.DeliveryRate,
			&score.RecordedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan historical score: %w", err)
		}
		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating historical scores: %w", err)
	}

	return scores, nil
}

// GetScoreAt retrieves the score closest to a specific timestamp
func (r *historicalScoresRepository) GetScoreAt(ctx context.Context, domainName string, timestamp int64) (*domain.HistoricalScore, error) {
	// Find the score closest to the given timestamp (either before or after)
	query := `
		SELECT id, domain, reputation_score, complaint_rate, bounce_rate, delivery_rate, recorded_at
		FROM historical_scores
		WHERE domain = ?
		ORDER BY ABS(recorded_at - ?)
		LIMIT 1
	`

	score := &domain.HistoricalScore{}
	err := r.db.QueryRowContext(ctx, query, domainName, timestamp).Scan(
		&score.ID,
		&score.Domain,
		&score.ReputationScore,
		&score.ComplaintRate,
		&score.BounceRate,
		&score.DeliveryRate,
		&score.RecordedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get score at timestamp: %w", err)
	}

	return score, nil
}

// GetDailyAverages retrieves daily average scores for trend analysis
func (r *historicalScoresRepository) GetDailyAverages(ctx context.Context, domainName string, days int) ([]*domain.HistoricalScore, error) {
	// Calculate the start time (days ago from now)
	startTime := time.Now().AddDate(0, 0, -days).Unix()

	query := `
		SELECT
			0 as id,
			domain,
			CAST(AVG(reputation_score) AS INTEGER) as reputation_score,
			AVG(complaint_rate) as complaint_rate,
			AVG(bounce_rate) as bounce_rate,
			AVG(delivery_rate) as delivery_rate,
			(recorded_at / 86400) * 86400 as day_timestamp
		FROM historical_scores
		WHERE domain = ? AND recorded_at >= ?
		GROUP BY domain, recorded_at / 86400
		ORDER BY day_timestamp ASC
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily averages: %w", err)
	}
	defer rows.Close()

	var scores []*domain.HistoricalScore
	for rows.Next() {
		score := &domain.HistoricalScore{}
		if err := rows.Scan(
			&score.ID,
			&score.Domain,
			&score.ReputationScore,
			&score.ComplaintRate,
			&score.BounceRate,
			&score.DeliveryRate,
			&score.RecordedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan daily average: %w", err)
		}
		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating daily averages: %w", err)
	}

	return scores, nil
}

// CleanupOldScores removes historical scores older than the specified timestamp
func (r *historicalScoresRepository) CleanupOldScores(ctx context.Context, olderThan int64) error {
	query := `DELETE FROM historical_scores WHERE recorded_at < ?`

	_, err := r.db.ExecContext(ctx, query, olderThan)
	if err != nil {
		return fmt.Errorf("failed to cleanup old historical scores: %w", err)
	}

	return nil
}
