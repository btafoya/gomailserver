package reputation

import (
	"context"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/service"
	"go.uber.org/zap"
)

// Scheduler handles periodic reputation management tasks
type Scheduler struct {
	telemetryService *service.TelemetryService
	logger           *zap.Logger
	stopChan         chan struct{}
}

// NewScheduler creates a new reputation scheduler
func NewScheduler(telemetryService *service.TelemetryService, logger *zap.Logger) *Scheduler {
	return &Scheduler{
		telemetryService: telemetryService,
		logger:           logger,
		stopChan:         make(chan struct{}),
	}
}

// Start begins the reputation scheduler goroutines
func (s *Scheduler) Start(ctx context.Context) error {
	s.logger.Info("starting reputation scheduler")

	// Start score calculation ticker (every 5 minutes)
	go s.runScoreCalculationLoop(ctx)

	// Start cleanup ticker (daily at 2 AM)
	go s.runCleanupLoop(ctx)

	return nil
}

// Stop gracefully stops the reputation scheduler
func (s *Scheduler) Stop() error {
	s.logger.Info("stopping reputation scheduler")
	close(s.stopChan)
	return nil
}

// runScoreCalculationLoop calculates reputation scores every 5 minutes
func (s *Scheduler) runScoreCalculationLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Run once immediately on startup
	s.calculateAllScores(ctx)

	for {
		select {
		case <-ticker.C:
			s.calculateAllScores(ctx)
		case <-s.stopChan:
			s.logger.Info("score calculation loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("score calculation loop context cancelled")
			return
		}
	}
}

// runCleanupLoop runs daily cleanup at 2 AM
func (s *Scheduler) runCleanupLoop(ctx context.Context) {
	// Calculate time until next 2 AM
	now := time.Now()
	next2AM := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if now.After(next2AM) {
		// If past 2 AM today, schedule for tomorrow
		next2AM = next2AM.Add(24 * time.Hour)
	}

	s.logger.Info("cleanup scheduled",
		zap.Time("next_run", next2AM),
	)

	// Wait until 2 AM
	timer := time.NewTimer(time.Until(next2AM))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Run cleanup
			s.runCleanup(ctx)

			// Schedule next cleanup (24 hours later)
			timer.Reset(24 * time.Hour)

		case <-s.stopChan:
			s.logger.Info("cleanup loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("cleanup loop context cancelled")
			return
		}
	}
}

// calculateAllScores calculates reputation scores for all domains
func (s *Scheduler) calculateAllScores(ctx context.Context) {
	s.logger.Debug("calculating reputation scores for all domains")

	start := time.Now()
	err := s.telemetryService.CalculateAllScores(ctx)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to calculate reputation scores",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Info("reputation scores calculated successfully",
		zap.Duration("duration", duration),
	)
}

// runCleanup runs the daily data cleanup task
func (s *Scheduler) runCleanup(ctx context.Context) {
	s.logger.Info("starting daily data cleanup")

	start := time.Now()
	err := s.telemetryService.CleanupOldData(ctx)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to cleanup old data",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Info("data cleanup completed successfully",
		zap.Duration("duration", duration),
	)
}
