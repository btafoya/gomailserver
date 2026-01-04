package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type predictionsRepository struct {
	db *sql.DB
}

// NewPredictionsRepository creates a new SQLite predictions repository
func NewPredictionsRepository(db *sql.DB) repository.PredictionsRepository {
	return &predictionsRepository{db: db}
}

// Create stores a new prediction
func (r *predictionsRepository) Create(ctx context.Context, prediction *domain.ReputationPrediction) error {
	query := `
		INSERT INTO reputation_predictions (
			domain, predicted_at, prediction_horizon, predicted_score,
			predicted_complaint_rate, predicted_bounce_rate, confidence_level,
			model_version, features_used
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		prediction.Domain,
		prediction.PredictedAt,
		prediction.PredictionHorizon,
		prediction.PredictedScore,
		prediction.PredictedComplaintRate,
		prediction.PredictedBounceRate,
		prediction.ConfidenceLevel,
		prediction.ModelVersion,
		prediction.FeaturesUsedJSON(),
	)
	if err != nil {
		return fmt.Errorf("failed to create prediction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	prediction.ID = id

	return nil
}

// GetLatest returns the latest prediction for a domain
func (r *predictionsRepository) GetLatest(ctx context.Context, domainName string) (*domain.ReputationPrediction, error) {
	query := `
		SELECT
			id, domain, predicted_at, prediction_horizon, predicted_score,
			predicted_complaint_rate, predicted_bounce_rate, confidence_level,
			model_version, features_used
		FROM reputation_predictions
		WHERE domain = ?
		ORDER BY predicted_at DESC
		LIMIT 1
	`

	prediction := &domain.ReputationPrediction{}
	var featuresJSON string
	err := r.db.QueryRowContext(ctx, query, domainName).Scan(
		&prediction.ID,
		&prediction.Domain,
		&prediction.PredictedAt,
		&prediction.PredictionHorizon,
		&prediction.PredictedScore,
		&prediction.PredictedComplaintRate,
		&prediction.PredictedBounceRate,
		&prediction.ConfidenceLevel,
		&prediction.ModelVersion,
		&featuresJSON,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no prediction found for domain %s", domainName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest prediction: %w", err)
	}

	// Parse features JSON
	if err := prediction.ParseFeaturesJSON(featuresJSON); err != nil {
		return nil, fmt.Errorf("failed to parse features JSON: %w", err)
	}

	return prediction, nil
}

// ListByDomain returns prediction history for a domain
func (r *predictionsRepository) ListByDomain(ctx context.Context, domainName string, limit int) ([]*domain.ReputationPrediction, error) {
	query := `
		SELECT
			id, domain, predicted_at, prediction_horizon, predicted_score,
			predicted_complaint_rate, predicted_bounce_rate, confidence_level,
			model_version, features_used
		FROM reputation_predictions
		WHERE domain = ?
		ORDER BY predicted_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list predictions: %w", err)
	}
	defer rows.Close()

	var predictions []*domain.ReputationPrediction
	for rows.Next() {
		prediction := &domain.ReputationPrediction{}
		var featuresJSON string
		err := rows.Scan(
			&prediction.ID,
			&prediction.Domain,
			&prediction.PredictedAt,
			&prediction.PredictionHorizon,
			&prediction.PredictedScore,
			&prediction.PredictedComplaintRate,
			&prediction.PredictedBounceRate,
			&prediction.ConfidenceLevel,
			&prediction.ModelVersion,
			&featuresJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prediction: %w", err)
		}

		// Parse features JSON
		if err := prediction.ParseFeaturesJSON(featuresJSON); err != nil {
			return nil, fmt.Errorf("failed to parse features JSON: %w", err)
		}

		predictions = append(predictions, prediction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating predictions: %w", err)
	}

	return predictions, nil
}

// GetByHorizon returns predictions for a specific time horizon
func (r *predictionsRepository) GetByHorizon(ctx context.Context, domainName string, hours int) (*domain.ReputationPrediction, error) {
	query := `
		SELECT
			id, domain, predicted_at, prediction_horizon, predicted_score,
			predicted_complaint_rate, predicted_bounce_rate, confidence_level,
			model_version, features_used
		FROM reputation_predictions
		WHERE domain = ? AND prediction_horizon = ?
		ORDER BY predicted_at DESC
		LIMIT 1
	`

	prediction := &domain.ReputationPrediction{}
	var featuresJSON string
	err := r.db.QueryRowContext(ctx, query, domainName, hours).Scan(
		&prediction.ID,
		&prediction.Domain,
		&prediction.PredictedAt,
		&prediction.PredictionHorizon,
		&prediction.PredictedScore,
		&prediction.PredictedComplaintRate,
		&prediction.PredictedBounceRate,
		&prediction.ConfidenceLevel,
		&prediction.ModelVersion,
		&featuresJSON,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no prediction found for domain %s with horizon %d hours", domainName, hours)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction by horizon: %w", err)
	}

	// Parse features JSON
	if err := prediction.ParseFeaturesJSON(featuresJSON); err != nil {
		return nil, fmt.Errorf("failed to parse features JSON: %w", err)
	}

	return prediction, nil
}
