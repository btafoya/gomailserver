package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
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

// GetEventsInRange retrieves events within a time range, expanding recurring events
func (s *EventService) GetEventsInRange(calendarID int64, start, end time.Time) ([]*domain.Event, error) {
	// Get all events for the calendar
	events, err := s.eventRepo.GetByCalendar(calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar events: %w", err)
	}

	// Expand recurring events and filter by time range
	var result []*domain.Event
	for _, event := range events {
		expanded, err := s.ExpandRecurrence(event, start, end)
		if err != nil {
			// Log error but continue with other events
			continue
		}
		result = append(result, expanded...)
	}

	return result, nil
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

// ExpandRecurrence expands recurring events within a time range using rrule-go
func (s *EventService) ExpandRecurrence(event *domain.Event, start, end time.Time) ([]*domain.Event, error) {
	// Non-recurring event: check if it overlaps with the time range
	if event.RRule == "" {
		if event.StartTime.Before(end) && event.EndTime.After(start) {
			return []*domain.Event{event}, nil
		}
		return []*domain.Event{}, nil
	}

	// Parse the RRULE
	rule, err := s.parseRRule(event.RRule, event.StartTime)
	if err != nil {
		// If we can't parse the RRULE, return the original event if it's in range
		if event.StartTime.Before(end) && event.EndTime.After(start) {
			return []*domain.Event{event}, nil
		}
		return []*domain.Event{}, nil
	}

	// Get occurrences in the time range
	occurrences := rule.Between(start.Add(-time.Hour*24), end.Add(time.Hour*24), true)

	// Limit to a reasonable number of occurrences
	const maxOccurrences = 100
	if len(occurrences) > maxOccurrences {
		occurrences = occurrences[:maxOccurrences]
	}

	// Calculate event duration
	duration := event.EndTime.Sub(event.StartTime)

	// Create event instances for each occurrence
	var expanded []*domain.Event
	for i, occStart := range occurrences {
		occEnd := occStart.Add(duration)

		// Check if this occurrence overlaps with the requested range
		if occStart.Before(end) && occEnd.After(start) {
			// Create a copy of the event with adjusted times
			instance := &domain.Event{
				ID:          event.ID, // Keep original ID for the master event
				CalendarID:  event.CalendarID,
				UID:         event.UID,
				Summary:     event.Summary,
				Description: event.Description,
				Location:    event.Location,
				StartTime:   occStart,
				EndTime:     occEnd,
				AllDay:      event.AllDay,
				Timezone:    event.Timezone,
				RRule:       event.RRule,
				Attendees:   event.Attendees,
				Organizer:   event.Organizer,
				Status:      event.Status,
				Sequence:    event.Sequence,
				ICalData:    s.generateInstanceICalData(event, occStart, occEnd, i),
				CreatedAt:   event.CreatedAt,
				UpdatedAt:   event.UpdatedAt,
			}
			// Generate unique ETag for this instance
			instance.ETag = s.generateInstanceETag(event, occStart)
			expanded = append(expanded, instance)
		}
	}

	return expanded, nil
}

// parseRRule parses an RRULE string into an rrule.RRule
func (s *EventService) parseRRule(rruleStr string, dtstart time.Time) (*rrule.RRule, error) {
	// Build the full RRULE string with DTSTART
	fullRule := fmt.Sprintf("DTSTART:%s\nRRULE:%s",
		dtstart.UTC().Format("20060102T150405Z"),
		rruleStr)

	rule, err := rrule.StrToRRule(fullRule)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RRULE: %w", err)
	}

	return rule, nil
}

// generateInstanceICalData generates iCalendar data for a specific occurrence
func (s *EventService) generateInstanceICalData(master *domain.Event, occStart, occEnd time.Time, instanceNum int) string {
	// Generate iCalendar data for this specific occurrence
	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//gomailserver//CalDAV Server//EN\r\n")
	sb.WriteString("BEGIN:VEVENT\r\n")
	sb.WriteString(fmt.Sprintf("UID:%s\r\n", master.UID))
	sb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", time.Now().UTC().Format("20060102T150405Z")))
	sb.WriteString(fmt.Sprintf("DTSTART:%s\r\n", occStart.UTC().Format("20060102T150405Z")))
	sb.WriteString(fmt.Sprintf("DTEND:%s\r\n", occEnd.UTC().Format("20060102T150405Z")))
	sb.WriteString(fmt.Sprintf("RECURRENCE-ID:%s\r\n", occStart.UTC().Format("20060102T150405Z")))

	if master.Summary != "" {
		sb.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapeICalText(master.Summary)))
	}
	if master.Description != "" {
		sb.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", escapeICalText(master.Description)))
	}
	if master.Location != "" {
		sb.WriteString(fmt.Sprintf("LOCATION:%s\r\n", escapeICalText(master.Location)))
	}
	if master.Status != "" {
		sb.WriteString(fmt.Sprintf("STATUS:%s\r\n", master.Status))
	}
	if master.Organizer != "" {
		sb.WriteString(fmt.Sprintf("ORGANIZER:%s\r\n", master.Organizer))
	}

	sb.WriteString("END:VEVENT\r\n")
	sb.WriteString("END:VCALENDAR\r\n")

	return sb.String()
}

// generateInstanceETag generates a unique ETag for a recurring event instance
func (s *EventService) generateInstanceETag(event *domain.Event, occStart time.Time) string {
	data := fmt.Sprintf("%s-%s-%s", event.UID, occStart.Format(time.RFC3339), event.UpdatedAt.Format(time.RFC3339))
	hash := sha256.Sum256([]byte(data))
	return `"` + hex.EncodeToString(hash[:16]) + `"`
}

// escapeICalText escapes special characters in iCalendar text values
func escapeICalText(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
