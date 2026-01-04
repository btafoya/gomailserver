package service

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// Default warm-up schedule (14-30 day progression)
var DefaultWarmUpSchedule = []domain.WarmUpDay{
	{Day: 1, MaxVolume: 100},
	{Day: 2, MaxVolume: 200},
	{Day: 3, MaxVolume: 500},
	{Day: 4, MaxVolume: 1000},
	{Day: 5, MaxVolume: 2000},
	{Day: 6, MaxVolume: 5000},
	{Day: 7, MaxVolume: 10000},
	{Day: 8, MaxVolume: 20000},
	{Day: 9, MaxVolume: 30000},
	{Day: 10, MaxVolume: 40000},
	{Day: 11, MaxVolume: 50000},
	{Day: 12, MaxVolume: 60000},
	{Day: 13, MaxVolume: 70000},
	{Day: 14, MaxVolume: 80000},
	// After day 14, warm-up is considered complete
}

// Thresholds for detecting new domains/IPs that need warm-up
const (
	NewDomainAgeThreshold = 30 * 24 * time.Hour  // 30 days
	NewIPFirstSeenWindow  = 7 * 24 * time.Hour   // 7 days
	MinimumSendingHistory = 100                  // minimum messages before warm-up ends
)

// WarmUpService manages domain and IP warm-up schedules
type WarmUpService struct {
	eventsRepo  repository.EventsRepository
	scoresRepo  repository.ScoresRepository
	warmUpRepo  repository.WarmUpRepository
	telemetry   *TelemetryService
	logger      *zap.Logger
}

// NewWarmUpService creates a new warm-up service
func NewWarmUpService(
	eventsRepo repository.EventsRepository,
	scoresRepo repository.ScoresRepository,
	warmUpRepo repository.WarmUpRepository,
	telemetry *TelemetryService,
	logger *zap.Logger,
) *WarmUpService {
	return &WarmUpService{
		eventsRepo:  eventsRepo,
		scoresRepo:  scoresRepo,
		warmUpRepo:  warmUpRepo,
		telemetry:   telemetry,
		logger:      logger,
	}
}

// DetectNewDomains identifies domains that need warm-up
// Should be called periodically (daily) to catch new sending domains
func (s *WarmUpService) DetectNewDomains(ctx context.Context) error {
	s.logger.Debug("detecting new domains requiring warm-up")

	// Get all reputation scores
	scores, err := s.scoresRepo.ListAllScores(ctx)
	if err != nil {
		return fmt.Errorf("failed to list reputation scores: %w", err)
	}

	for _, score := range scores {
		// Skip domains already in warm-up or with enough history
		if score.WarmUpActive {
			continue
		}

		// Check if domain needs warm-up based on history
		needsWarmUp, reason := s.checkWarmUpNeeded(ctx, score)
		if needsWarmUp {
			s.logger.Info("starting warm-up for domain",
				zap.String("domain", score.Domain),
				zap.String("reason", reason),
			)

			if err := s.StartWarmUp(ctx, score.Domain, DefaultWarmUpSchedule); err != nil {
				s.logger.Error("failed to start warm-up",
					zap.String("domain", score.Domain),
					zap.Error(err),
				)
				continue
			}
		}
	}

	return nil
}

// checkWarmUpNeeded determines if a domain requires warm-up
func (s *WarmUpService) checkWarmUpNeeded(ctx context.Context, score *domain.ReputationScore) (bool, string) {
	// Get domain's sending history
	now := time.Now().Unix()
	startTime := now - int64(NewDomainAgeThreshold.Seconds())

	events, err := s.eventsRepo.GetEventsInWindow(ctx, score.Domain, startTime, now)
	if err != nil {
		s.logger.Error("failed to get domain events",
			zap.String("domain", score.Domain),
			zap.Error(err),
		)
		return false, ""
	}

	// If no sending history, definitely needs warm-up
	if len(events) == 0 {
		return true, "new domain with no sending history"
	}

	// Count total sends
	totalSends := 0
	for _, event := range events {
		if event.EventType == domain.EventSent || event.EventType == domain.EventDelivered {
			totalSends++
		}
	}

	// If very low volume, needs warm-up
	if totalSends < MinimumSendingHistory {
		return true, fmt.Sprintf("low sending volume: %d messages in last 30 days", totalSends)
	}

	// Check if using a new sending IP (7-day window)
	recentStartTime := now - int64(NewIPFirstSeenWindow.Seconds())
	recentEvents, err := s.eventsRepo.GetEventsInWindow(ctx, score.Domain, recentStartTime, now)
	if err != nil {
		return false, ""
	}

	// Track unique IPs
	ipFirstSeen := make(map[string]int64)
	for _, event := range recentEvents {
		if event.IPAddress != "" {
			if _, exists := ipFirstSeen[event.IPAddress]; !exists {
				ipFirstSeen[event.IPAddress] = event.Timestamp
			}
		}
	}

	// If using a very new IP (less than 7 days old), needs warm-up
	for ip, firstSeen := range ipFirstSeen {
		age := now - firstSeen
		if age < int64(NewIPFirstSeenWindow.Seconds()) {
			return true, fmt.Sprintf("new sending IP detected: %s", ip)
		}
	}

	return false, ""
}

// StartWarmUp initiates warm-up schedule for a domain
func (s *WarmUpService) StartWarmUp(ctx context.Context, domainName string, schedule []domain.WarmUpDay) error {
	// Convert schedule to pointers
	schedulePointers := make([]*domain.WarmUpDay, len(schedule))
	for i := range schedule {
		schedulePointers[i] = &schedule[i]
	}

	// Create warm-up schedule in database
	if err := s.warmUpRepo.CreateSchedule(ctx, domainName, schedulePointers); err != nil {
		return fmt.Errorf("failed to create warm-up schedule: %w", err)
	}

	// Update reputation score to reflect warm-up state
	score, err := s.scoresRepo.GetReputationScore(ctx, domainName)
	if err != nil {
		return fmt.Errorf("failed to get reputation score: %w", err)
	}

	score.WarmUpActive = true
	score.WarmUpDay = 1 // Start on day 1

	if err := s.scoresRepo.UpdateReputationScore(ctx, score); err != nil {
		return fmt.Errorf("failed to update reputation score: %w", err)
	}

	s.logger.Info("warm-up started",
		zap.String("domain", domainName),
		zap.Int("schedule_days", len(schedule)),
	)

	return nil
}

// AdvanceWarmUp moves domains to the next day in their warm-up schedule
// Should be called daily (at midnight or similar)
func (s *WarmUpService) AdvanceWarmUp(ctx context.Context) error {
	s.logger.Debug("advancing warm-up schedules")

	// Get all domains currently in warm-up
	scores, err := s.scoresRepo.ListAllScores(ctx)
	if err != nil {
		return fmt.Errorf("failed to list reputation scores: %w", err)
	}

	for _, score := range scores {
		if !score.WarmUpActive {
			continue
		}

		// Get the schedule to check if we can advance
		schedule, err := s.warmUpRepo.GetSchedule(ctx, score.Domain)
		if err != nil {
			s.logger.Error("failed to get warm-up schedule",
				zap.String("domain", score.Domain),
				zap.Error(err),
			)
			continue
		}

		// Check if current day's volume was met
		if score.WarmUpDay > 0 && score.WarmUpDay <= len(schedule) {
			daySchedule := schedule[score.WarmUpDay-1]

			// Only advance if we reached at least 80% of the target volume
			targetVolume := float64(daySchedule.MaxVolume) * 0.8
			if float64(daySchedule.ActualVolume) < targetVolume {
				s.logger.Warn("warm-up day not completed, not advancing",
					zap.String("domain", score.Domain),
					zap.Int("day", score.WarmUpDay),
					zap.Int("actual_volume", daySchedule.ActualVolume),
					zap.Int("target_volume", daySchedule.MaxVolume),
				)
				continue
			}
		}

		// Advance to next day
		nextDay := score.WarmUpDay + 1

		// Check if warm-up is complete
		if nextDay > len(schedule) {
			s.logger.Info("warm-up completed",
				zap.String("domain", score.Domain),
				zap.Int("total_days", score.WarmUpDay),
			)

			if err := s.CompleteWarmUp(ctx, score.Domain); err != nil {
				s.logger.Error("failed to complete warm-up",
					zap.String("domain", score.Domain),
					zap.Error(err),
				)
			}
			continue
		}

		// Update to next day
		score.WarmUpDay = nextDay
		if err := s.scoresRepo.UpdateReputationScore(ctx, score); err != nil {
			s.logger.Error("failed to advance warm-up day",
				zap.String("domain", score.Domain),
				zap.Error(err),
			)
			continue
		}

		// Reset actual volume for new day
		if err := s.warmUpRepo.UpdateDayVolume(ctx, score.Domain, nextDay, 0); err != nil {
			s.logger.Error("failed to reset day volume",
				zap.String("domain", score.Domain),
				zap.Error(err),
			)
		}

		s.logger.Info("warm-up advanced",
			zap.String("domain", score.Domain),
			zap.Int("day", nextDay),
			zap.Int("max_volume", schedule[nextDay-1].MaxVolume),
		)
	}

	return nil
}

// CompleteWarmUp marks warm-up as finished for a domain
func (s *WarmUpService) CompleteWarmUp(ctx context.Context, domainName string) error {
	score, err := s.scoresRepo.GetReputationScore(ctx, domainName)
	if err != nil {
		return fmt.Errorf("failed to get reputation score: %w", err)
	}

	score.WarmUpActive = false
	score.WarmUpDay = 0

	if err := s.scoresRepo.UpdateReputationScore(ctx, score); err != nil {
		return fmt.Errorf("failed to update reputation score: %w", err)
	}

	s.logger.Info("warm-up completed successfully",
		zap.String("domain", domainName),
	)

	return nil
}

// GetWarmUpStatus returns the current warm-up status for a domain
func (s *WarmUpService) GetWarmUpStatus(ctx context.Context, domainName string) (*domain.WarmUpStatus, error) {
	score, err := s.scoresRepo.GetReputationScore(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get reputation score: %w", err)
	}

	if !score.WarmUpActive {
		return &domain.WarmUpStatus{
			Active:    false,
			Domain:    domainName,
			Completed: true,
		}, nil
	}

	schedule, err := s.warmUpRepo.GetSchedule(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get warm-up schedule: %w", err)
	}

	if score.WarmUpDay < 1 || score.WarmUpDay > len(schedule) {
		return nil, fmt.Errorf("invalid warm-up day: %d", score.WarmUpDay)
	}

	currentDay := schedule[score.WarmUpDay-1]

	return &domain.WarmUpStatus{
		Active:        true,
		Domain:        domainName,
		CurrentDay:    score.WarmUpDay,
		TotalDays:     len(schedule),
		MaxVolume:     currentDay.MaxVolume,
		ActualVolume:  currentDay.ActualVolume,
		VolumePercent: float64(currentDay.ActualVolume) / float64(currentDay.MaxVolume) * 100,
		Completed:     false,
	}, nil
}

// ManualComplete allows admin to manually complete warm-up
func (s *WarmUpService) ManualComplete(ctx context.Context, domainName string, adminNotes string) error {
	s.logger.Info("manually completing warm-up",
		zap.String("domain", domainName),
		zap.String("admin_notes", adminNotes),
	)

	return s.CompleteWarmUp(ctx, domainName)
}
