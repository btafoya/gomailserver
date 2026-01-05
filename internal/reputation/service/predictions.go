package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// PredictionsService generates reputation predictions using trend analysis

type PredictionsService struct {
	predictionsRepo repository.PredictionsRepository
	eventsRepo      repository.EventsRepository
	scoresRepo      repository.ScoresRepository
	logger          *zap.Logger
}

func NewPredictionsService(
	predictionsRepo repository.PredictionsRepository,
	eventsRepo repository.EventsRepository,
	scoresRepo repository.ScoresRepository,
	logger *zap.Logger,
) *PredictionsService {
	return &PredictionsService{
		predictionsRepo: predictionsRepo,
		eventsRepo:      eventsRepo,
		scoresRepo:      scoresRepo,
		logger:          logger,
	}
}

// GeneratePrediction generates a reputation prediction for a domain
func (s *PredictionsService) GeneratePrediction(ctx context.Context, domainName string, horizonHours int) (*domain.ReputationPrediction, error) {
	// Get historical scores (last 7 days)
	now := time.Now()
	startTime := now.Add(-7 * 24 * time.Hour).Unix()
	endTime := now.Unix()

	// Get event counts
	eventCounts, err := s.eventsRepo.GetEventCountsByType(ctx, domainName, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get event counts: %w", err)
	}

	totalEvents := int64(0)
	for _, count := range eventCounts {
		totalEvents += count
	}

	if totalEvents == 0 {
		s.logger.Warn("No events found for prediction", zap.String("domain", domainName))
		return nil, nil
	}

	// Calculate current rates
	complaints := eventCounts[string(domain.EventComplaint)]
	bounces := eventCounts[string(domain.EventBounce)]
	_ = eventCounts[string(domain.EventDelivery)] // delivered count not used in prediction

	currentComplaintRate := float64(complaints) / float64(totalEvents)
	currentBounceRate := float64(bounces) / float64(totalEvents)

	// Get current score
	currentScore, err := s.scoresRepo.GetReputationScore(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get current score: %w", err)
	}

	// Calculate trends (simple linear trend)
	trendScore := s.calculateScoreTrend(ctx, domainName)
	trendComplaint := s.calculateComplaintTrend(ctx, domainName)
	trendBounce := s.calculateBounceTrend(ctx, domainName)

	// Project future values
	predictionHours := float64(horizonHours)
	predictedScore := int(math.Max(0, math.Min(100, float64(currentScore.ReputationScore)+(trendScore*predictionHours/24))))
	predictedComplaintRate := math.Max(0, currentComplaintRate+(trendComplaint*predictionHours/24))
	predictedBounceRate := math.Max(0, currentBounceRate+(trendBounce*predictionHours/24))

	// Calculate confidence based on data volume
	confidence := s.calculateConfidence(totalEvents, 7)

	prediction := &domain.ReputationPrediction{
		Domain:                domainName,
		PredictedAt:           time.Now().Unix(),
		PredictionHorizon:     horizonHours,
		PredictedScore:        predictedScore,
		PredictedComplaintRate: predictedComplaintRate,
		PredictedBounceRate:   predictedBounceRate,
		ConfidenceLevel:       confidence,
		ModelVersion:          "trend-v1",
		FeaturesUsed: map[string]interface{}{
			"current_score":        currentScore.ReputationScore,
			"current_complaint_rate": currentComplaintRate,
			"current_bounce_rate":  currentBounceRate,
			"trend_score":          trendScore,
			"trend_complaint":      trendComplaint,
			"trend_bounce":         trendBounce,
			"total_events":         totalEvents,
			"days_analyzed":        7,
		},
	}

	// Store prediction
	if err := s.predictionsRepo.Create(ctx, prediction); err != nil {
		s.logger.Error("Failed to store prediction", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Generated reputation prediction",
		zap.String("domain", domainName),
		zap.Int("horizon_hours", horizonHours),
		zap.Int("predicted_score", predictedScore),
		zap.Float64("confidence", confidence),
	)

	return prediction, nil
}

// calculateScoreTrend calculates the score change trend
func (s *PredictionsService) calculateScoreTrend(ctx context.Context, domainName string) float64 {
	// Get current score
	_, err := s.scoresRepo.GetReputationScore(ctx, domainName)
	if err != nil {
		return 0
	}

	// Get score from 24 hours ago (simplified - would need historical scores table in real implementation)
	// For now, return 0 (no trend)
	// TODO: Implement historical score tracking for better trend analysis
	return 0
}

// calculateComplaintTrend calculates complaint rate trend
func (s *PredictionsService) calculateComplaintTrend(ctx context.Context, domainName string) float64 {
	// Compare complaint rates between last 24h and previous 24h
	now := time.Now()
	last24h := s.getComplaintRateForPeriod(ctx, domainName, now.Add(-24*time.Hour), now)
	prev24h := s.getComplaintRateForPeriod(ctx, domainName, now.Add(-48*time.Hour), now.Add(-24*time.Hour))

	return last24h - prev24h
}

// calculateBounceTrend calculates bounce rate trend
func (s *PredictionsService) calculateBounceTrend(ctx context.Context, domainName string) float64 {
	// Compare bounce rates between last 24h and previous 24h
	now := time.Now()
	last24h := s.getBounceRateForPeriod(ctx, domainName, now.Add(-24*time.Hour), now)
	prev24h := s.getBounceRateForPeriod(ctx, domainName, now.Add(-48*time.Hour), now.Add(-24*time.Hour))

	return last24h - prev24h
}

// getComplaintRateForPeriod gets complaint rate for a time period
func (s *PredictionsService) getComplaintRateForPeriod(ctx context.Context, domainName string, start, end time.Time) float64 {
	eventCounts, err := s.eventsRepo.GetEventCountsByType(ctx, domainName, start.Unix(), end.Unix())
	if err != nil {
		return 0
	}

	totalEvents := int64(0)
	for _, count := range eventCounts {
		totalEvents += count
	}

	if totalEvents == 0 {
		return 0
	}

	complaints := eventCounts[string(domain.EventComplaint)]
	return float64(complaints) / float64(totalEvents)
}

// getBounceRateForPeriod gets bounce rate for a time period
func (s *PredictionsService) getBounceRateForPeriod(ctx context.Context, domainName string, start, end time.Time) float64 {
	eventCounts, err := s.eventsRepo.GetEventCountsByType(ctx, domainName, start.Unix(), end.Unix())
	if err != nil {
		return 0
	}

	totalEvents := int64(0)
	for _, count := range eventCounts {
		totalEvents += count
	}

	if totalEvents == 0 {
		return 0
	}

	bounces := eventCounts[string(domain.EventBounce)]
	return float64(bounces) / float64(totalEvents)
}

// calculateConfidence calculates prediction confidence based on data volume
func (s *PredictionsService) calculateConfidence(totalEvents int64, days int) float64 {
	// Confidence increases with more data
	// Minimum 100 events/day for high confidence
	minEventsPerDay := int64(100)
	optimalEvents := minEventsPerDay * int64(days)

	if totalEvents >= optimalEvents {
		return 0.9
	} else if totalEvents >= optimalEvents/2 {
		return 0.7
	} else if totalEvents >= optimalEvents/4 {
		return 0.5
	}

	return 0.3
}

// GetLatestPrediction returns the latest prediction for a domain
func (s *PredictionsService) GetLatestPrediction(ctx context.Context, domainName string) (*domain.ReputationPrediction, error) {
	return s.predictionsRepo.GetLatest(ctx, domainName)
}

// GetPredictionHistory returns prediction history for a domain
func (s *PredictionsService) GetPredictionHistory(ctx context.Context, domainName string, limit int) ([]*domain.ReputationPrediction, error) {
	return s.predictionsRepo.ListByDomain(ctx, domainName, limit)
}

// GeneratePredictionsForAllDomains generates predictions for all active domains
func (s *PredictionsService) GeneratePredictionsForAllDomains(ctx context.Context, domains []string, horizonHours int) error {
	s.logger.Info("Generating predictions for all domains",
		zap.Int("count", len(domains)),
		zap.Int("horizon_hours", horizonHours),
	)

	errors := make([]error, 0)
	for _, domainName := range domains {
		if _, err := s.GeneratePrediction(ctx, domainName, horizonHours); err != nil {
			s.logger.Error("Failed to generate prediction",
				zap.String("domain", domainName),
				zap.Error(err),
			)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("prediction generation completed with %d errors", len(errors))
	}

	return nil
}

// GeneratePredictions generates predictions for a domain (alias for GeneratePrediction with default horizon)
func (s *PredictionsService) GeneratePredictions(ctx context.Context, domainName string) error {
	// Use default 24-hour prediction horizon
	_, err := s.GeneratePrediction(ctx, domainName, 24)
	return err
}
