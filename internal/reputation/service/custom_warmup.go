package service

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// CustomWarmupService manages custom warm-up schedules

type CustomWarmupService struct {
	warmupRepo repository.CustomWarmupRepository
	logger     *zap.Logger
}

func NewCustomWarmupService(
	warmupRepo repository.CustomWarmupRepository,
	logger *zap.Logger,
) *CustomWarmupService {
	return &CustomWarmupService{
		warmupRepo: warmupRepo,
		logger:     logger,
	}
}

// CreateSchedule creates a custom warm-up schedule for a domain
func (s *CustomWarmupService) CreateSchedule(ctx context.Context, domainName, scheduleName, createdBy string, schedule []*domain.CustomWarmupSchedule) error {
	// Validate schedule
	if len(schedule) == 0 {
		return fmt.Errorf("schedule must have at least one day")
	}

	// Set common fields
	now := time.Now().Unix()
	for i, day := range schedule {
		day.Domain = domainName
		day.ScheduleName = scheduleName
		day.Day = i + 1
		day.CreatedAt = now
		day.CreatedBy = createdBy
		day.IsActive = true
	}

	// Create schedule
	if err := s.warmupRepo.CreateSchedule(ctx, schedule); err != nil {
		s.logger.Error("Failed to create custom warm-up schedule",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Created custom warm-up schedule",
		zap.String("domain", domainName),
		zap.String("schedule", scheduleName),
		zap.Int("days", len(schedule)),
		zap.String("created_by", createdBy),
	)

	return nil
}

// GetSchedule retrieves the warm-up schedule for a domain
func (s *CustomWarmupService) GetSchedule(ctx context.Context, domainName string) ([]*domain.CustomWarmupSchedule, error) {
	return s.warmupRepo.GetSchedule(ctx, domainName)
}

// UpdateSchedule updates an existing schedule
func (s *CustomWarmupService) UpdateSchedule(ctx context.Context, schedule *domain.CustomWarmupSchedule) error {
	return s.warmupRepo.UpdateSchedule(ctx, schedule)
}

// DeleteSchedule deletes a custom schedule
func (s *CustomWarmupService) DeleteSchedule(ctx context.Context, domainName string) error {
	if err := s.warmupRepo.DeleteSchedule(ctx, domainName); err != nil {
		s.logger.Error("Failed to delete warm-up schedule",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Deleted warm-up schedule", zap.String("domain", domainName))
	return nil
}

// SetActive activates or deactivates a schedule
func (s *CustomWarmupService) SetActive(ctx context.Context, domainName string, active bool) error {
	if err := s.warmupRepo.SetActive(ctx, domainName, active); err != nil {
		return err
	}

	s.logger.Info("Updated warm-up schedule status",
		zap.String("domain", domainName),
		zap.Bool("active", active),
	)

	return nil
}

// GetVolumeForDay returns the max volume for a specific day
func (s *CustomWarmupService) GetVolumeForDay(ctx context.Context, domainName string, day int) (int, error) {
	schedule, err := s.GetSchedule(ctx, domainName)
	if err != nil {
		return 0, err
	}

	for _, daySchedule := range schedule {
		if daySchedule.Day == day {
			return daySchedule.MaxVolume, nil
		}
	}

	return 0, fmt.Errorf("day %d not found in schedule", day)
}

// CreateConservativeSchedule creates a conservative 30-day warm-up schedule
func (s *CustomWarmupService) CreateConservativeSchedule(ctx context.Context, domainName, createdBy string) error {
	schedule := make([]*domain.CustomWarmupSchedule, 30)

	// Very conservative ramp: 50 → 100K over 30 days
	volumes := []int{
		50, 100, 200, 400, 600, 800, 1000, 1500,    // Week 1
		2000, 3000, 4000, 5000, 6000, 7000, 8000,   // Week 2
		10000, 12000, 14000, 16000, 18000, 20000,   // Week 3
		25000, 30000, 35000, 40000, 50000,          // Week 4
		60000, 70000, 80000, 90000, 100000,         // Week 5
	}

	for i, volume := range volumes {
		schedule[i] = &domain.CustomWarmupSchedule{
			MaxVolume: volume,
		}
	}

	return s.CreateSchedule(ctx, domainName, "Conservative 30-day", createdBy, schedule)
}

// CreateAggressiveSchedule creates an aggressive 14-day warm-up schedule
func (s *CustomWarmupService) CreateAggressiveSchedule(ctx context.Context, domainName, createdBy string) error {
	schedule := make([]*domain.CustomWarmupSchedule, 14)

	// Aggressive ramp: 100 → 80K over 14 days
	volumes := []int{
		100, 200, 500, 1000, 2000, 5000, 10000, // Week 1
		15000, 20000, 30000, 40000, 50000, 65000, 80000, // Week 2
	}

	for i, volume := range volumes {
		schedule[i] = &domain.CustomWarmupSchedule{
			MaxVolume: volume,
		}
	}

	return s.CreateSchedule(ctx, domainName, "Aggressive 14-day", createdBy, schedule)
}

// CreateModerateSchedule creates a moderate 21-day warm-up schedule
func (s *CustomWarmupService) CreateModerateSchedule(ctx context.Context, domainName, createdBy string) error {
	schedule := make([]*domain.CustomWarmupSchedule, 21)

	// Moderate ramp: 100 → 60K over 21 days
	volumes := []int{
		100, 200, 500, 1000, 2000, 3000, 5000,      // Week 1
		7000, 10000, 12000, 15000, 18000, 21000, 25000, // Week 2
		30000, 35000, 40000, 45000, 50000, 55000, 60000, // Week 3
	}

	for i, volume := range volumes {
		schedule[i] = &domain.CustomWarmupSchedule{
			MaxVolume: volume,
		}
	}

	return s.CreateSchedule(ctx, domainName, "Moderate 21-day", createdBy, schedule)
}

// ListActiveSchedules returns all active custom schedules
func (s *CustomWarmupService) ListActiveSchedules(ctx context.Context) (map[string][]*domain.CustomWarmupSchedule, error) {
	return s.warmupRepo.ListActiveSchedules(ctx)
}
