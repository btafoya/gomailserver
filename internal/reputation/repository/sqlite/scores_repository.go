package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type scoresRepository struct {
	db *sql.DB
}

// NewScoresRepository creates a new SQLite scores repository
func NewScoresRepository(db *sql.DB) repository.ScoresRepository {
	return &scoresRepository{db: db}
}

// GetReputationScore retrieves the reputation score for a domain
func (r *scoresRepository) GetReputationScore(ctx context.Context, domainName string) (*domain.ReputationScore, error) {
	query := `
		SELECT
			domain, reputation_score, complaint_rate, bounce_rate, delivery_rate,
			circuit_breaker_active, circuit_breaker_reason, warm_up_active, warm_up_day, last_updated
		FROM domain_reputation_scores
		WHERE domain = ?
	`

	score := &domain.ReputationScore{}
	err := r.db.QueryRowContext(ctx, query, domainName).Scan(
		&score.Domain,
		&score.ReputationScore,
		&score.ComplaintRate,
		&score.BounceRate,
		&score.DeliveryRate,
		&score.CircuitBreakerActive,
		&score.CircuitBreakerReason,
		&score.WarmUpActive,
		&score.WarmUpDay,
		&score.LastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("reputation score not found for domain %s", domainName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reputation score: %w", err)
	}

	return score, nil
}

// UpdateReputationScore updates or creates a reputation score for a domain
func (r *scoresRepository) UpdateReputationScore(ctx context.Context, score *domain.ReputationScore) error {
	query := `
		INSERT INTO domain_reputation_scores (
			domain, reputation_score, complaint_rate, bounce_rate, delivery_rate,
			circuit_breaker_active, circuit_breaker_reason, warm_up_active, warm_up_day, last_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(domain) DO UPDATE SET
			reputation_score = excluded.reputation_score,
			complaint_rate = excluded.complaint_rate,
			bounce_rate = excluded.bounce_rate,
			delivery_rate = excluded.delivery_rate,
			circuit_breaker_active = excluded.circuit_breaker_active,
			circuit_breaker_reason = excluded.circuit_breaker_reason,
			warm_up_active = excluded.warm_up_active,
			warm_up_day = excluded.warm_up_day,
			last_updated = excluded.last_updated
	`

	_, err := r.db.ExecContext(ctx, query,
		score.Domain,
		score.ReputationScore,
		score.ComplaintRate,
		score.BounceRate,
		score.DeliveryRate,
		score.CircuitBreakerActive,
		score.CircuitBreakerReason,
		score.WarmUpActive,
		score.WarmUpDay,
		score.LastUpdated,
	)
	if err != nil {
		return fmt.Errorf("failed to update reputation score: %w", err)
	}

	return nil
}

// ListAllScores retrieves all domain reputation scores
func (r *scoresRepository) ListAllScores(ctx context.Context) ([]*domain.ReputationScore, error) {
	query := `
		SELECT
			domain, reputation_score, complaint_rate, bounce_rate, delivery_rate,
			circuit_breaker_active, circuit_breaker_reason, warm_up_active, warm_up_day, last_updated
		FROM domain_reputation_scores
		ORDER BY domain ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list scores: %w", err)
	}
	defer rows.Close()

	var scores []*domain.ReputationScore
	for rows.Next() {
		score := &domain.ReputationScore{}
		err := rows.Scan(
			&score.Domain,
			&score.ReputationScore,
			&score.ComplaintRate,
			&score.BounceRate,
			&score.DeliveryRate,
			&score.CircuitBreakerActive,
			&score.CircuitBreakerReason,
			&score.WarmUpActive,
			&score.WarmUpDay,
			&score.LastUpdated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan score: %w", err)
		}
		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scores: %w", err)
	}

	return scores, nil
}
