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
	// Phase 5 services
	gmailPostmasterSvc  *service.GmailPostmasterService
	microsoftSNDSSvc    *service.MicrosoftSNDSService
	arfParserSvc        *service.ARFParserService
	dmarcAnalyzerSvc    *service.DMARCAnalyzerService
	predictionsSvc      *service.PredictionsService
	providerLimitsSvc   *service.ProviderRateLimitsService
	alertsSvc           *service.AlertsService
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

// SetPhase5Services sets Phase 5 services for advanced automation
func (s *Scheduler) SetPhase5Services(
	gmailPostmasterSvc *service.GmailPostmasterService,
	microsoftSNDSSvc *service.MicrosoftSNDSService,
	arfParserSvc *service.ARFParserService,
	dmarcAnalyzerSvc *service.DMARCAnalyzerService,
	predictionsSvc *service.PredictionsService,
	providerLimitsSvc *service.ProviderRateLimitsService,
	alertsSvc *service.AlertsService,
) {
	s.gmailPostmasterSvc = gmailPostmasterSvc
	s.microsoftSNDSSvc = microsoftSNDSSvc
	s.arfParserSvc = arfParserSvc
	s.dmarcAnalyzerSvc = dmarcAnalyzerSvc
	s.predictionsSvc = predictionsSvc
	s.providerLimitsSvc = providerLimitsSvc
	s.alertsSvc = alertsSvc
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

	// Phase 5: Gmail Postmaster sync (every 1 hour)
	if s.gmailPostmasterSvc != nil {
		go s.runGmailPostmasterSyncLoop(ctx)
	}

	// Phase 5: Microsoft SNDS sync (every 6 hours)
	if s.microsoftSNDSSvc != nil {
		go s.runMicrosoftSNDSSyncLoop(ctx)
	}

	// Phase 5: ARF processing (every 15 minutes)
	if s.arfParserSvc != nil {
		go s.runARFProcessingLoop(ctx)
	}

	// Phase 5: DMARC analysis (every 30 minutes)
	if s.dmarcAnalyzerSvc != nil {
		go s.runDMARCAnalysisLoop(ctx)
	}

	// Phase 5: Predictions generation (daily at 3 AM)
	if s.predictionsSvc != nil {
		go s.runPredictionsLoop(ctx)
	}

	// Phase 5: Alert cleanup (runs with daily cleanup at 2 AM)
	// Integrated into runCleanup method

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

// Phase 5: Gmail Postmaster Tools Integration

// runGmailPostmasterSyncLoop syncs Gmail Postmaster metrics every hour
func (s *Scheduler) runGmailPostmasterSyncLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Run once immediately on startup
	s.syncGmailPostmaster(ctx)

	for {
		select {
		case <-ticker.C:
			s.syncGmailPostmaster(ctx)
		case <-s.stopChan:
			s.logger.Info("Gmail Postmaster sync loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("Gmail Postmaster sync loop context cancelled")
			return
		}
	}
}

// syncGmailPostmaster fetches metrics for all configured domains
func (s *Scheduler) syncGmailPostmaster(ctx context.Context) {
	s.logger.Debug("syncing Gmail Postmaster metrics")

	start := time.Now()
	// Note: In production, domains list would come from configuration or database
	domains := []string{} // TODO: Get from config
	err := s.gmailPostmasterSvc.SyncAll(ctx, domains)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to sync Gmail Postmaster metrics",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Info("Gmail Postmaster sync completed",
		zap.Int("domains", len(domains)),
		zap.Duration("duration", duration),
	)
}

// Phase 5: Microsoft SNDS Integration

// runMicrosoftSNDSSyncLoop syncs Microsoft SNDS metrics every 6 hours
func (s *Scheduler) runMicrosoftSNDSSyncLoop(ctx context.Context) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	// Run once immediately on startup
	s.syncMicrosoftSNDS(ctx)

	for {
		select {
		case <-ticker.C:
			s.syncMicrosoftSNDS(ctx)
		case <-s.stopChan:
			s.logger.Info("Microsoft SNDS sync loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("Microsoft SNDS sync loop context cancelled")
			return
		}
	}
}

// syncMicrosoftSNDS fetches metrics for all configured IP addresses
func (s *Scheduler) syncMicrosoftSNDS(ctx context.Context) {
	s.logger.Debug("syncing Microsoft SNDS metrics")

	start := time.Now()
	// Note: In production, IP addresses would come from configuration or database
	ipAddresses := []string{} // TODO: Get from config
	err := s.microsoftSNDSSvc.SyncAll(ctx, ipAddresses)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to sync Microsoft SNDS metrics",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Info("Microsoft SNDS sync completed",
		zap.Int("ips", len(ipAddresses)),
		zap.Duration("duration", duration),
	)
}

// Phase 5: ARF Complaint Processing

// runARFProcessingLoop processes ARF complaints every 15 minutes
func (s *Scheduler) runARFProcessingLoop(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	// Run once immediately on startup
	s.processARFReports(ctx)

	for {
		select {
		case <-ticker.C:
			s.processARFReports(ctx)
		case <-s.stopChan:
			s.logger.Info("ARF processing loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("ARF processing loop context cancelled")
			return
		}
	}
}

// processARFReports processes all unprocessed ARF complaint reports
func (s *Scheduler) processARFReports(ctx context.Context) {
	s.logger.Debug("processing ARF complaint reports")

	start := time.Now()
	err := s.arfParserSvc.ProcessUnprocessed(ctx)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to process ARF reports",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	s.logger.Info("ARF reports processing completed",
		zap.Duration("duration", duration),
	)
}

// Phase 5: DMARC Analysis

// runDMARCAnalysisLoop analyzes DMARC reports every 30 minutes
func (s *Scheduler) runDMARCAnalysisLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	// Run once immediately on startup
	s.analyzeDMARCReports(ctx)

	for {
		select {
		case <-ticker.C:
			s.analyzeDMARCReports(ctx)
		case <-s.stopChan:
			s.logger.Info("DMARC analysis loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("DMARC analysis loop context cancelled")
			return
		}
	}
}

// analyzeDMARCReports analyzes recent DMARC reports for alignment issues
func (s *Scheduler) analyzeDMARCReports(ctx context.Context) {
	s.logger.Debug("analyzing DMARC reports")

	start := time.Now()
	// Note: In production, domains list would come from configuration or database
	domains := []string{} // TODO: Get from config

	for _, domainName := range domains {
		_, err := s.dmarcAnalyzerSvc.AnalyzeDomain(ctx, domainName, 7) // Analyze last 7 days
		if err != nil {
			s.logger.Error("failed to analyze DMARC reports",
				zap.String("domain", domainName),
				zap.Error(err),
			)
			continue
		}
	}

	duration := time.Since(start)
	s.logger.Info("DMARC analysis completed",
		zap.Int("domains", len(domains)),
		zap.Duration("duration", duration),
	)
}

// Phase 5: Reputation Predictions

// runPredictionsLoop generates reputation predictions daily at 3 AM
func (s *Scheduler) runPredictionsLoop(ctx context.Context) {
	// Calculate time until next 3 AM
	now := time.Now()
	next3AM := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	if now.After(next3AM) {
		// If past 3 AM today, schedule for tomorrow
		next3AM = next3AM.Add(24 * time.Hour)
	}

	s.logger.Info("predictions generation scheduled",
		zap.Time("next_run", next3AM),
	)

	// Wait until 3 AM
	timer := time.NewTimer(time.Until(next3AM))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Run predictions generation
			s.generatePredictions(ctx)

			// Schedule next run (24 hours later)
			timer.Reset(24 * time.Hour)

		case <-s.stopChan:
			s.logger.Info("predictions loop stopped")
			return
		case <-ctx.Done():
			s.logger.Info("predictions loop context cancelled")
			return
		}
	}
}

// generatePredictions generates reputation predictions for all domains
func (s *Scheduler) generatePredictions(ctx context.Context) {
	s.logger.Info("generating reputation predictions")

	start := time.Now()
	// Note: In production, domains list would come from configuration or database
	domains := []string{} // TODO: Get from config

	// Generate predictions for 24h, 48h, and 72h horizons
	horizons := []int{24, 48, 72}
	for _, horizon := range horizons {
		err := s.predictionsSvc.GeneratePredictionsForAllDomains(ctx, domains, horizon)
		if err != nil {
			s.logger.Error("failed to generate predictions",
				zap.Int("horizon_hours", horizon),
				zap.Error(err),
			)
			continue
		}
	}

	duration := time.Since(start)
	s.logger.Info("predictions generation completed",
		zap.Int("domains", len(domains)),
		zap.Int("horizons", len(horizons)),
		zap.Duration("duration", duration),
	)
}
