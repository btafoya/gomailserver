package caldav

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
	"github.com/btafoya/gomailserver/internal/webdav"
)

// Simple mock implementations for testing

type MockCalendarService struct {
	calendars   []*domain.Calendar
	events      []*domain.Event
	shouldError bool
}

func (m *MockCalendarService) CreateCalendar(userID int64, name, displayName, color, description, timezone string) (*domain.Calendar, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	calendar := &domain.Calendar{
		ID:          1,
		UserID:      userID,
		Name:        name,
		DisplayName: displayName,
		Color:       color,
		Description: description,
		Timezone:    timezone,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.calendars = append(m.calendars, calendar)
	return calendar, nil
}

func (m *MockCalendarService) GetCalendar(id int64) (*domain.Calendar, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	for _, cal := range m.calendars {
		if cal.ID == id {
			return cal, nil
		}
	}
	return nil, fmt.Errorf("calendar not found")
}

func (m *MockCalendarService) GetUserCalendars(userID int64) ([]*domain.Calendar, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	var userCalendars []*domain.Calendar
	for _, cal := range m.calendars {
		if cal.UserID == userID {
			userCalendars = append(userCalendars, cal)
		}
	}
	return userCalendars, nil
}

func (m *MockCalendarService) UpdateCalendar(id int64, displayName, color, description, timezone *string) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for _, cal := range m.calendars {
		if cal.ID == id {
			if displayName != nil {
				cal.DisplayName = *displayName
			}
			if color != nil {
				cal.Color = *color
			}
			if description != nil {
				cal.Description = *description
			}
			if timezone != nil {
				cal.Timezone = *timezone
			}
			return nil
		}
	}
	return fmt.Errorf("calendar not found")
}

func (m *MockCalendarService) DeleteCalendar(id int64) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for i, cal := range m.calendars {
		if cal.ID == id {
			m.calendars = append(m.calendars[:i], m.calendars[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("calendar not found")
}

func (m *MockCalendarService) GenerateSyncToken(id int64) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("mock error")
	}
	return fmt.Sprintf("sync-token-%d", id), nil
}

func (m *MockCalendarService) GetSharedCalendars(userID int64) ([]*domain.Calendar, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return []*domain.Calendar{}, nil
}

func (m *MockCalendarService) GetAll() ([]*domain.Calendar, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.calendars, nil
}

func (m *MockCalendarService) ShareCalendar(calendarID int64, readUsers, writeUsers, adminUsers []int64, readAll bool) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockCalendarService) UnshareCalendar(calendarID int64, userID int64) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockCalendarService) HasReadAccess(userID, calendarID int64) bool {
	return true
}

func (m *MockCalendarService) HasWriteAccess(userID, calendarID int64) bool {
	return true
}

func (m *MockCalendarService) HasAdminAccess(userID, calendarID int64) bool {
	return true
}

// addUserContext adds user authentication context to a request
func addUserContext(req *http.Request, userID int64) *http.Request {
	ctx := context.WithValue(req.Context(), webdav.UserIDKey, userID)
	ctx = context.WithValue(ctx, webdav.UsernameKey, "testuser")
	ctx = context.WithValue(ctx, webdav.UserEmailKey, "testuser@example.com")
	return req.WithContext(ctx)
}

type MockEventService struct {
	events      []*domain.Event
	shouldError bool
}

func (m *MockEventService) CreateEvent(calendarID int64, icalData string) (*domain.Event, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	event := &domain.Event{
		ID:         int64(len(m.events) + 1),
		CalendarID: calendarID,
		UID:        fmt.Sprintf("event-%d", len(m.events)+1),
		Summary:    "Mock Event",
		ICalData:   icalData,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	m.events = append(m.events, event)
	return event, nil
}

func (m *MockEventService) GetEvent(id int64) (*domain.Event, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	for _, event := range m.events {
		if event.ID == id {
			return event, nil
		}
	}
	return nil, fmt.Errorf("event not found")
}

func (m *MockEventService) GetCalendarEvents(calendarID int64) ([]*domain.Event, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	var calendarEvents []*domain.Event
	for _, event := range m.events {
		if event.CalendarID == calendarID {
			calendarEvents = append(calendarEvents, event)
		}
	}
	return calendarEvents, nil
}

func (m *MockEventService) GetEventsInRange(calendarID int64, start, end time.Time) ([]*domain.Event, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	var events []*domain.Event
	for _, event := range m.events {
		if event.CalendarID == calendarID && event.StartTime.After(start) && event.EndTime.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (m *MockEventService) UpdateEvent(id int64, icalData string) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for _, event := range m.events {
		if event.ID == id {
			event.ICalData = icalData
			event.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("event not found")
}

func (m *MockEventService) DeleteEvent(id int64) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for i, event := range m.events {
		if event.ID == id {
			m.events = append(m.events[:i], m.events[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("event not found")
}

func (m *MockEventService) GenerateETag(event *domain.Event) string {
	return fmt.Sprintf("etag-%d", event.ID)
}

func (m *MockEventService) ExpandRecurrence(event *domain.Event, start, end time.Time) ([]*domain.Event, error) {
	return []*domain.Event{event}, nil
}

// TestCalDAVHandler_MkCalendar tests MKCALENDAR requests
func TestCalDAVHandler_MkCalendar(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockCalendarService := &MockCalendarService{shouldError: false}
	mockEventService := &MockEventService{shouldError: false}

	handler := NewHandler(logger, mockCalendarService, mockEventService, nil)

	// Create test request with user context
	req := httptest.NewRequest("MKCALENDAR", "/caldav/calendars/123/Test%20Calendar", strings.NewReader(""))
	req = addUserContext(req, 123)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	// Check if location contains expected path components
	location := w.Header().Get("Location")
	if !strings.Contains(location, "123") || !strings.Contains(location, "Test") {
		t.Errorf("Expected location to contain calendar path components, got %s", location)
	}
}

// TestCalDAVHandler_ReportCalendarQuery tests calendar-query REPORT requests
func TestCalDAVHandler_ReportCalendarQuery(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockCalendarService := &MockCalendarService{shouldError: false}
	mockEventService := &MockEventService{shouldError: false}

	// Setup test data
	testCalendar := &domain.Calendar{
		ID:     1,
		UserID: 123,
		Name:   "Test Calendar",
	}
	mockCalendarService.calendars = []*domain.Calendar{testCalendar}

	testEvent := &domain.Event{
		ID:         1,
		CalendarID: 1,
		UID:        "test-event-1",
		Summary:    "Test Event 1",
		ETag:       "test-etag-1",
		ICalData:   "BEGIN:VCALENDAR\nVERSION:2.0\nBEGIN:VEVENT\nUID:test-event-1\nSUMMARY:Test Event 1\nEND:VEVENT\nEND:VCALENDAR",
	}
	mockEventService.events = []*domain.Event{testEvent}

	handler := NewHandler(logger, mockCalendarService, mockEventService, nil)

	// Create test request
	reportBody := `<C:calendar-query xmlns:C="urn:ietf:params:xml:ns:caldav">
		<D:prop xmlns:D="DAV:">
			<D:getetag/>
			<C:calendar-data/>
		</D:prop>
		<C:filter>
			<C:comp-filter name="VEVENT"/>
		</C:filter>
	</C:calendar-query>`

	req := httptest.NewRequest("REPORT", "/caldav/calendars/123/Test%20Calendar", strings.NewReader(reportBody))
	req = addUserContext(req, 123)
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMultiStatus {
		t.Errorf("Expected status %d, got %d", http.StatusMultiStatus, w.Code)
	}
	if !strings.Contains(w.Body.String(), "test-event-1.ics") {
		t.Errorf("Expected response to contain event data")
	}
}

// TestCalDAVHandler_ReportFreeBusyQuery tests free-busy-query REPORT requests
func TestCalDAVHandler_ReportFreeBusyQuery(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockCalendarService := &MockCalendarService{shouldError: false}
	mockEventService := &MockEventService{shouldError: false}

	// Setup test data
	testCalendar := &domain.Calendar{
		ID:     1,
		UserID: 123,
		Name:   "Test Calendar",
	}
	mockCalendarService.calendars = []*domain.Calendar{testCalendar}

	handler := NewHandler(logger, mockCalendarService, mockEventService, nil)

	// Create test request
	reportBody := `<C:free-busy-query xmlns:C="urn:ietf:params:xml:ns:caldav">
		<D:prop xmlns:D="DAV:">
			<C:calendar-data/>
		</D:prop>
		<C:time-range start="20230101T000000Z" end="20230131T235959Z"/>
	</C:free-busy-query>`

	req := httptest.NewRequest("REPORT", "/caldav/calendars/123/Test%20Calendar", strings.NewReader(reportBody))
	req = addUserContext(req, 123)
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "BEGIN:VCALENDAR") {
		t.Errorf("Expected response to contain calendar data")
	}
	if !strings.Contains(w.Body.String(), "BEGIN:VFREEBUSY") {
		t.Errorf("Expected response to contain free/busy data")
	}
}

// TestCalDAVHandler_InvalidPath tests handling of invalid paths
func TestCalDAVHandler_InvalidPath(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockCalendarService := &MockCalendarService{shouldError: false}
	mockEventService := &MockEventService{shouldError: false}

	handler := NewHandler(logger, mockCalendarService, mockEventService, nil)

	// Test invalid path (too short)
	req := httptest.NewRequest("MKCALENDAR", "/caldav/calendars/123", strings.NewReader(""))
	req = addUserContext(req, 123)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Invalid calendar path") {
		t.Errorf("Expected response to contain error message")
	}
}

// TestCalDAVHandler_CalendarNotFound tests handling of non-existent calendars
func TestCalDAVHandler_CalendarNotFound(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockCalendarService := &MockCalendarService{shouldError: false}
	mockEventService := &MockEventService{shouldError: false}

	handler := NewHandler(logger, mockCalendarService, mockEventService, nil)

	// Create test request - no calendars set in mock
	reportBody := `<C:calendar-query xmlns:C="urn:ietf:params:xml:ns:caldav">
		<D:prop xmlns:D="DAV:">
			<D:getetag/>
			<C:calendar-data/>
		</D:prop>
	</C:calendar-query>`

	req := httptest.NewRequest("REPORT", "/caldav/calendars/123/NonExistent", strings.NewReader(reportBody))
	req = addUserContext(req, 123)
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Calendar not found") {
		t.Errorf("Expected response to contain error message")
	}
}
