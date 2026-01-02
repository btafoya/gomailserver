package domain

import (
	"time"
)

// Event represents a calendar event
type Event struct {
	ID          int64
	CalendarID  int64
	UID         string
	Summary     string
	Description string
	Location    string
	StartTime   time.Time
	EndTime     time.Time
	AllDay      bool
	Timezone    string
	RRule       string // Recurrence rule (RFC 5545)
	Attendees   string // JSON array of attendees
	Organizer   string
	Status      string // CONFIRMED, TENTATIVE, CANCELLED
	Sequence    int
	ETag        string
	ICalData    string // Full iCalendar data
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// EventRepository defines the interface for event persistence
type EventRepository interface {
	// Create creates a new event
	Create(event *Event) error

	// GetByID retrieves an event by ID
	GetByID(id int64) (*Event, error)

	// GetByUID retrieves an event by UID and calendar ID
	GetByUID(calendarID int64, uid string) (*Event, error)

	// GetByCalendar retrieves all events for a calendar
	GetByCalendar(calendarID int64) ([]*Event, error)

	// GetByTimeRange retrieves events within a time range
	GetByTimeRange(calendarID int64, start, end time.Time) ([]*Event, error)

	// Update updates an existing event
	Update(event *Event) error

	// Delete deletes an event
	Delete(id int64) error

	// UpdateETag updates the ETag for an event
	UpdateETag(id int64, etag string) error
}

// EventService defines business logic for events
type EventService interface {
	// CreateEvent creates a new event from iCalendar data
	CreateEvent(calendarID int64, icalData string) (*Event, error)

	// GetEvent retrieves an event by ID
	GetEvent(id int64) (*Event, error)

	// GetCalendarEvents retrieves all events for a calendar
	GetCalendarEvents(calendarID int64) ([]*Event, error)

	// GetEventsInRange retrieves events within a time range
	GetEventsInRange(calendarID int64, start, end time.Time) ([]*Event, error)

	// UpdateEvent updates an event from iCalendar data
	UpdateEvent(id int64, icalData string) error

	// DeleteEvent deletes an event
	DeleteEvent(id int64) error

	// GenerateETag generates a new ETag for an event
	GenerateETag(event *Event) string

	// ExpandRecurrence expands recurring events within a time range
	ExpandRecurrence(event *Event, start, end time.Time) ([]*Event, error)
}

// Attendee represents an event attendee
type Attendee struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`   // CHAIR, REQ-PARTICIPANT, OPT-PARTICIPANT
	Status string `json:"status"` // NEEDS-ACTION, ACCEPTED, DECLINED, TENTATIVE
}
