package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
)

// CalendarService implements domain.CalendarService
type CalendarService struct {
	calendarRepo domain.CalendarRepository
	eventRepo    domain.EventRepository
}

// NewCalendarService creates a new calendar service
func NewCalendarService(calendarRepo domain.CalendarRepository, eventRepo domain.EventRepository) *CalendarService {
	return &CalendarService{
		calendarRepo: calendarRepo,
		eventRepo:    eventRepo,
	}
}

// CreateCalendar creates a new calendar for a user
func (s *CalendarService) CreateCalendar(userID int64, name, displayName, color, description, timezone string) (*domain.Calendar, error) {
	// Check if calendar with same name already exists
	existing, err := s.calendarRepo.GetByUserAndName(userID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing calendar: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("calendar with name %s already exists", name)
	}

	// Set defaults
	if timezone == "" {
		timezone = "UTC"
	}
	if displayName == "" {
		displayName = name
	}

	// Generate initial sync token
	syncToken, err := s.generateSyncToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate sync token: %w", err)
	}

	calendar := &domain.Calendar{
		UserID:      userID,
		Name:        name,
		DisplayName: displayName,
		Color:       color,
		Description: description,
		Timezone:    timezone,
		SyncToken:   syncToken,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.calendarRepo.Create(calendar); err != nil {
		return nil, fmt.Errorf("failed to create calendar: %w", err)
	}

	return calendar, nil
}

// GetCalendar retrieves a calendar by ID
func (s *CalendarService) GetCalendar(id int64) (*domain.Calendar, error) {
	calendar, err := s.calendarRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar: %w", err)
	}
	if calendar == nil {
		return nil, fmt.Errorf("calendar not found")
	}
	return calendar, nil
}

// GetUserCalendars retrieves all calendars for a user
func (s *CalendarService) GetUserCalendars(userID int64) ([]*domain.Calendar, error) {
	calendars, err := s.calendarRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user calendars: %w", err)
	}
	return calendars, nil
}

// UpdateCalendar updates calendar properties
func (s *CalendarService) UpdateCalendar(id int64, displayName, color, description, timezone *string) error {
	calendar, err := s.calendarRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get calendar: %w", err)
	}
	if calendar == nil {
		return fmt.Errorf("calendar not found")
	}

	// Update fields if provided
	if displayName != nil {
		calendar.DisplayName = *displayName
	}
	if color != nil {
		calendar.Color = *color
	}
	if description != nil {
		calendar.Description = *description
	}
	if timezone != nil {
		calendar.Timezone = *timezone
	}

	if err := s.calendarRepo.Update(calendar); err != nil {
		return fmt.Errorf("failed to update calendar: %w", err)
	}

	return nil
}

// DeleteCalendar deletes a calendar and all its events
func (s *CalendarService) DeleteCalendar(id int64) error {
	// Repository handles cascade delete of events
	if err := s.calendarRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete calendar: %w", err)
	}
	return nil
}

// GenerateSyncToken generates a new sync token for a calendar
func (s *CalendarService) GenerateSyncToken(id int64) (string, error) {
	token, err := s.generateSyncToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate sync token: %w", err)
	}

	if err := s.calendarRepo.UpdateSyncToken(id, token); err != nil {
		return "", fmt.Errorf("failed to update sync token: %w", err)
	}

	return token, nil
}

// generateSyncToken generates a random sync token
func (s *CalendarService) generateSyncToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
