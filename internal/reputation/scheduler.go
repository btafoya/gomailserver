package reputation

import (
	"context"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/service"
	"go.uber.org/zap"
)

// Scheduler handles periodic reputation management tasks
type Scheduler struct {
	telemetryService    *service.TelemetryService
	circuitBreakerSvc   *service.CircuitBreakerService
	warmUpSvc           *service.WarmUpService
	logger              *zap.Logger
	stopChan            chan struct{}
}

// NewScheduler creates a new reputation scheduler
func NewScheduler(
	telemetryService *service.TelemetryService,
	circuitBreakerSvc *service.CircuitBreakerService,
	warmUpSvc *service.WarmUpService,
	logger *zap.Logger,
) *Scheduler {
	return &Scheduler{
		telemetryService:  telemetryService,
		circuitBreakerSvc: circuitBreakerSvc,
		warmUpSvc:         warmUpSvc,
		logger:            logger,
		stopChan:          make(chan struct{}),
	}
}

// Start begins the reputation scheduler goroutines
func (s *Scheduler) Start(ctx context.Context) error {
	s.logger.Info("starting reputation scheduler")

	// Phase 1: Score calculation (every 5 minutes)
	go s.runScoreCalculationLoop(ctx)

	// Phase 1: Cleanup (daily at 2 AM)
	go s.runCleanupLoop(ctx)

	// Phase 3: Circuit breaker check (every 15 minutes)
	if s.circuitBreakerSvc != nil {
		go s.runCircuitBreakerCheckLoop(ctx)
	}

	// Phase 3: Auto-resume check (every hour)
	if s.circuitBreakerSvc != nil {
		go s.runAutoResumeLoop(ctx)
	}

	// Phase 3: Warm-up advancement (daily at midnight)
	if s.warmUpSvc != nil {
		go s.runWarmUpAdvancementLoop(ctx)
	}

	// Phase 3: New domain detection (daily at 1 AM)
	if s.warmUpSvc != nil {
		go s.runNewDomainDetectionLoop(ctx)
	}

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

// runCircuitBreakerCheckLoop checks circuit breaker thresholds every 15 minutes
func (s *Scheduler) runCircuitBreakerCheckLoop(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	// Run once immediately on startup
	s.checkCircuitBreakers(ctx)

	for {
		select {
		case <-ticker.C:
			s.checkCircuitBreakers(ctx)
		case <-s.stopChan:
			s.logger.Info("circuit breaker check loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("circuit breaker check loop context cancelled")
			return
		}
	}
}

// checkCircuitBreakers evaluates all domains for circuit breaker triggers
func (s *Scheduler) checkCircuitBreakers(ctx context.Context) {
	s.logger.Debug("checking circuit breaker thresholds")

	start := time.Now()
	err := s.circuitBreakerSvc.CheckAndTrigger(ctx)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to check circuit breakers",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Debug("circuit breaker check completed",
		zap.Duration("duration", duration),
	)
}

// runAutoResumeLoop attempts to resume paused domains every hour
func (s *Scheduler) runAutoResumeLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Run once immediately on startup
	s.attemptAutoResume(ctx)

	for {
		select {
		case <-ticker.C:
			s.attemptAutoResume(ctx)
		case <-s.stopChan:
			s.logger.Info("auto-resume loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("auto-resume loop context cancelled")
			return
		}
	}
}

// attemptAutoResume tries to resume domains with exponential backoff
func (s *Scheduler) attemptAutoResume(ctx context.Context) {
	s.logger.Debug("attempting to auto-resume paused domains")

	start := time.Now()
	err := s.circuitBreakerSvc.AutoResume(ctx)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to auto-resume domains",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Debug("auto-resume check completed",
		zap.Duration("duration", duration),
	)
}

// runWarmUpAdvancementLoop advances warm-up schedules daily at midnight
func (s *Scheduler) runWarmUpAdvancementLoop(ctx context.Context) {
	// Calculate time until next midnight
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if now.After(nextMidnight) {
		// If past midnight today, schedule for tomorrow
		nextMidnight = nextMidnight.Add(24 * time.Hour)
	}

	s.logger.Info("warm-up advancement scheduled",
		zap.Time("next_run", nextMidnight),
	)

	// Wait until midnight
	timer := time.NewTimer(time.Until(nextMidnight))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Run warm-up advancement
			s.advanceWarmUp(ctx)

			// Schedule next run (24 hours later)
			timer.Reset(24 * time.Hour)

		case <-s.stopChan:
			s.logger.Info("warm-up advancement loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("warm-up advancement loop context cancelled")
			return
		}
	}
}

// advanceWarmUp moves domains to the next day in their warm-up schedule
func (s *Scheduler) advanceWarmUp(ctx context.Context) {
	s.logger.Info("advancing warm-up schedules")

	start := time.Now()
	err := s.warmUpSvc.AdvanceWarmUp(ctx)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to advance warm-up",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Info("warm-up advancement completed successfully",
		zap.Duration("duration", duration),
	)
}

// runNewDomainDetectionLoop detects new domains requiring warm-up daily at 1 AM
func (s *Scheduler) runNewDomainDetectionLoop(ctx context.Context) {
	// Calculate time until next 1 AM
	now := time.Now()
	next1AM := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, now.Location())
	if now.After(next1AM) {
		// If past 1 AM today, schedule for tomorrow
		next1AM = next1AM.Add(24 * time.Hour)
	}

	s.logger.Info("new domain detection scheduled",
		zap.Time("next_run", next1AM),
	)

	// Wait until 1 AM
	timer := time.NewTimer(time.Until(next1AM))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Run new domain detection
			s.detectNewDomains(ctx)

			// Schedule next run (24 hours later)
			timer.Reset(24 * time.Hour)

		case <-s.stopChan:
			s.logger.Info("new domain detection loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("new domain detection loop context cancelled")
			return
		}
	}
}

// detectNewDomains identifies domains that need warm-up
func (s *Scheduler) detectNewDomains(ctx context.Context) {
	s.logger.Info("detecting new domains requiring warm-up")

	start := time.Now()
	err := s.warmUpSvc.DetectNewDomains(ctx)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to detect new domains",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Info("new domain detection completed successfully",
		zap.Duration("duration", duration),
	)
}
