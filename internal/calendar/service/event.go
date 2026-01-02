package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
	"github.com/emersion/go-ical"
)

// EventService implements domain.EventService
type EventService struct {
	eventRepo    domain.EventRepository
	calendarRepo domain.CalendarRepository
}

// NewEventService creates a new event service
func NewEventService(eventRepo domain.EventRepository, calendarRepo domain.CalendarRepository) *EventService {
	return &EventService{
		eventRepo:    eventRepo,
		calendarRepo: calendarRepo,
	}
}

// CreateEvent creates a new event from iCalendar data
func (s *EventService) CreateEvent(calendarID int64, icalData string) (*domain.Event, error) {
	// Verify calendar exists
	calendar, err := s.calendarRepo.GetByID(calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar: %w", err)
	}
	if calendar == nil {
		return nil, fmt.Errorf("calendar not found")
	}

	// Parse iCalendar data
	dec := ical.NewDecoder(strings.NewReader(icalData))
	cal, err := dec.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to parse iCalendar data: %w", err)
	}

	// Extract VEVENT component
	var vevent *ical.Component
	for _, comp := range cal.Children {
		if comp.Name == "VEVENT" {
			vevent = comp
			break
		}
	}
	if vevent == nil {
		return nil, fmt.Errorf("no VEVENT component found in iCalendar data")
	}

	// Extract event properties
	event := &domain.Event{
		CalendarID: calendarID,
		ICalData:   icalData,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Extract UID (required)
	if prop := vevent.Props.Get("UID"); prop != nil {
		event.UID = prop.Value
	} else {
		return nil, fmt.Errorf("UID property is required")
	}

	// Check if event with same UID already exists
	existing, err := s.eventRepo.GetByUID(calendarID, event.UID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing event: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("event with UID %s already exists", event.UID)
	}

	// Extract other properties
	if prop := vevent.Props.Get("SUMMARY"); prop != nil {
		event.Summary = prop.Value
	}
	if prop := vevent.Props.Get("DESCRIPTION"); prop != nil {
		event.Description = prop.Value
	}
	if prop := vevent.Props.Get("LOCATION"); prop != nil {
		event.Location = prop.Value
	}
	if prop := vevent.Props.Get("STATUS"); prop != nil {
		event.Status = prop.Value
	} else {
		event.Status = "CONFIRMED"
	}
	if prop := vevent.Props.Get("ORGANIZER"); prop != nil {
		event.Organizer = prop.Value
	}
	if prop := vevent.Props.Get("SEQUENCE"); prop != nil {
		fmt.Sscanf(prop.Value, "%d", &event.Sequence)
	}

	// Extract DTSTART (required)
	if prop := vevent.Props.Get("DTSTART"); prop != nil {
		dtstart, err := prop.DateTime(time.UTC)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DTSTART: %w", err)
		}
		event.StartTime = dtstart

		// Check if it's an all-day event
		if dateParam := prop.Params.Get("VALUE"); dateParam == "DATE" {
			event.AllDay = true
		}

		// Extract timezone
		if tzid := prop.Params.Get("TZID"); tzid != "" {
			event.Timezone = tzid
		}
	} else {
		return nil, fmt.Errorf("DTSTART property is required")
	}

	// Extract DTEND
	if prop := vevent.Props.Get("DTEND"); prop != nil {
		dtend, err := prop.DateTime(time.UTC)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DTEND: %w", err)
		}
		event.EndTime = dtend
	} else {
		// If no DTEND, use DTSTART + 1 hour (or next day for all-day events)
		if event.AllDay {
			event.EndTime = event.StartTime.Add(24 * time.Hour)
		} else {
			event.EndTime = event.StartTime.Add(1 * time.Hour)
		}
	}

	// Extract RRULE (recurrence rule)
	if prop := vevent.Props.Get("RRULE"); prop != nil {
		event.RRule = prop.Value
	}

	// Generate ETag
	event.ETag = s.GenerateETag(event)

	// Create event in repository
	if err := s.eventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

// GetEvent retrieves an event by ID
func (s *EventService) GetEvent(id int64) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}

// GetCalendarEvents retrieves all events for a calendar
func (s *EventService) GetCalendarEvents(calendarID int64) ([]*domain.Event, error) {
	events, err := s.eventRepo.GetByCalendar(calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar events: %w", err)
	}
	return events, nil
}

// GetEventsInRange retrieves events within a time range
func (s *EventService) GetEventsInRange(calendarID int64, start, end time.Time) ([]*domain.Event, error) {
	events, err := s.eventRepo.GetByTimeRange(calendarID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get events in range: %w", err)
	}
	return events, nil
}

// UpdateEvent updates an event from iCalendar data
func (s *EventService) UpdateEvent(id int64, icalData string) error {
	// Get existing event
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}
	if event == nil {
		return fmt.Errorf("event not found")
	}

	// Parse new iCalendar data
	dec := ical.NewDecoder(strings.NewReader(icalData))
	cal, err := dec.Decode()
	if err != nil {
		return fmt.Errorf("failed to parse iCalendar data: %w", err)
	}

	// Extract VEVENT component
	var vevent *ical.Component
	for _, comp := range cal.Children {
		if comp.Name == "VEVENT" {
			vevent = comp
			break
		}
	}
	if vevent == nil {
		return fmt.Errorf("no VEVENT component found in iCalendar data")
	}

	// Update event properties
	event.ICalData = icalData
	event.UpdatedAt = time.Now()

	// Extract properties (similar to CreateEvent)
	if prop := vevent.Props.Get("SUMMARY"); prop != nil {
		event.Summary = prop.Value
	}
	if prop := vevent.Props.Get("DESCRIPTION"); prop != nil {
		event.Description = prop.Value
	}
	if prop := vevent.Props.Get("LOCATION"); prop != nil {
		event.Location = prop.Value
	}
	if prop := vevent.Props.Get("STATUS"); prop != nil {
		event.Status = prop.Value
	}
	if prop := vevent.Props.Get("SEQUENCE"); prop != nil {
		fmt.Sscanf(prop.Value, "%d", &event.Sequence)
	}

	if prop := vevent.Props.Get("DTSTART"); prop != nil {
		dtstart, err := prop.DateTime(time.UTC)
		if err != nil {
			return fmt.Errorf("failed to parse DTSTART: %w", err)
		}
		event.StartTime = dtstart
	}

	if prop := vevent.Props.Get("DTEND"); prop != nil {
		dtend, err := prop.DateTime(time.UTC)
		if err != nil {
			return fmt.Errorf("failed to parse DTEND: %w", err)
		}
		event.EndTime = dtend
	}

	if prop := vevent.Props.Get("RRULE"); prop != nil {
		event.RRule = prop.Value
	}

	// Update ETag
	event.ETag = s.GenerateETag(event)

	// Update in repository
	if err := s.eventRepo.Update(event); err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// DeleteEvent deletes an event
func (s *EventService) DeleteEvent(id int64) error {
	if err := s.eventRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

// GenerateETag generates a new ETag for an event
func (s *EventService) GenerateETag(event *domain.Event) string {
	// Generate ETag based on event content and update time
	data := fmt.Sprintf("%s-%s-%d", event.UID, event.UpdatedAt.Format(time.RFC3339), event.Sequence)
	hash := sha256.Sum256([]byte(data))
	return `"` + hex.EncodeToString(hash[:]) + `"`
}

// ExpandRecurrence expands recurring events within a time range
func (s *EventService) ExpandRecurrence(event *domain.Event, start, end time.Time) ([]*domain.Event, error) {
	// TODO: Implement recurrence expansion using rrule-go library
	// For now, just return the original event if it's in range
	if event.RRule == "" {
		// Non-recurring event
		if event.StartTime.Before(end) && event.EndTime.After(start) {
			return []*domain.Event{event}, nil
		}
		return []*domain.Event{}, nil
	}

	// Placeholder for recurrence expansion
	// This would use github.com/teambition/rrule-go to expand the RRULE
	return []*domain.Event{event}, nil
}
