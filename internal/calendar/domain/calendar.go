package domain

import (
	"time"
)

// Calendar represents a calendar collection
type Calendar struct {
	ID          int64
	UserID      int64
	Name        string
	DisplayName string
	Color       string
	Description string
	Timezone    string
	SyncToken   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CalendarRepository defines the interface for calendar persistence
type CalendarRepository interface {
	// Create creates a new calendar
	Create(calendar *Calendar) error

	// GetByID retrieves a calendar by ID
	GetByID(id int64) (*Calendar, error)

	// GetByUserID retrieves all calendars for a user
	GetByUserID(userID int64) ([]*Calendar, error)

	// GetByUserAndName retrieves a calendar by user ID and name
	GetByUserAndName(userID int64, name string) (*Calendar, error)

	// Update updates an existing calendar
	Update(calendar *Calendar) error

	// Delete deletes a calendar
	Delete(id int64) error

	// UpdateSyncToken updates the sync token for a calendar
	UpdateSyncToken(id int64, token string) error
}

// CalendarService defines business logic for calendars
type CalendarService interface {
	// CreateCalendar creates a new calendar for a user
	CreateCalendar(userID int64, name, displayName, color, description, timezone string) (*Calendar, error)

	// GetCalendar retrieves a calendar by ID
	GetCalendar(id int64) (*Calendar, error)

	// GetUserCalendars retrieves all calendars for a user
	GetUserCalendars(userID int64) ([]*Calendar, error)

	// UpdateCalendar updates calendar properties
	UpdateCalendar(id int64, displayName, color, description, timezone *string) error

	// DeleteCalendar deletes a calendar and all its events
	DeleteCalendar(id int64) error

	// GenerateSyncToken generates a new sync token for a calendar
	GenerateSyncToken(id int64) (string, error)
}
