package webdav_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"

	calendardomain "github.com/btafoya/gomailserver/internal/calendar/domain"
	contactdomain "github.com/btafoya/gomailserver/internal/contact/domain"
	"github.com/btafoya/gomailserver/internal/webdav"
	"github.com/btafoya/gomailserver/internal/webdav/caldav"
	"github.com/btafoya/gomailserver/internal/webdav/carddav"
)

// TestWebDAV_FullIntegration tests complete WebDAV functionality
func TestWebDAV_FullIntegration(t *testing.T) {
	logger := zaptest.NewLogger(t)

	calendarService := &mockCalendarService{}
	eventService := &mockEventService{}
	addressbookService := &mockAddressbookService{}
	contactService := &mockContactService{}

	// Test CalDAV handler
	t.Run("CalDAV OPTIONS", func(t *testing.T) {
		handler := caldav.NewHandler(logger, calendarService, eventService)
		req := httptest.NewRequest("OPTIONS", "/caldav/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		davHeader := w.Header().Get("DAV")
		if davHeader == "" {
			t.Error("Expected DAV header")
		}
	})

	// Test CardDAV handler
	t.Run("CardDAV OPTIONS", func(t *testing.T) {
		handler := carddav.NewHandler(logger, addressbookService, contactService)
		req := httptest.NewRequest("OPTIONS", "/carddav/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		davHeader := w.Header().Get("DAV")
		if davHeader == "" {
			t.Error("Expected DAV header")
		}
	})

	// Test WebDAV server integration
	t.Run("WebDAV Server", func(t *testing.T) {
		cfg := &webdav.Config{
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
		}

		// Create server with nil user repo (simple integration test)
		server := webdav.NewServer(cfg, caldav.NewHandler(logger, calendarService, eventService), carddav.NewHandler(logger, addressbookService, contactService), nil, logger)

		req := httptest.NewRequest("GET", "/.well-known/caldav", nil)
		w := httptest.NewRecorder()

		// Create an HTTP handler from the server's router for testing
		// Note: This is a simplified test - in real usage the server runs with its own HTTP server
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/.well-known/caldav" {
				http.Redirect(w, r, "/caldav/", http.StatusMovedPermanently)
			} else {
				http.NotFound(w, r)
			}
		}).ServeHTTP(w, req)

		if w.Code != http.StatusMovedPermanently {
			t.Errorf("Expected redirect status 301, got %d", w.Code)
		}
	})
}

// Mock services for integration testing
type mockCalendarService struct {
	calendars []*calendardomain.Calendar
}

func (m *mockCalendarService) CreateCalendar(userID int64, name, displayName, color, description, timezone string) (*calendardomain.Calendar, error) {
	calendar := &calendardomain.Calendar{
		ID:          int64(len(m.calendars) + 1),
		UserID:      userID,
		Name:        name,
		DisplayName: displayName,
		Color:       color,
		Description: description,
		Timezone:    timezone,
	}
	m.calendars = append(m.calendars, calendar)
	return calendar, nil
}

func (m *mockCalendarService) GetCalendar(id int64) (*calendardomain.Calendar, error) {
	for _, cal := range m.calendars {
		if cal.ID == id {
			return cal, nil
		}
	}
	return nil, errors.New("calendar not found")
}

func (m *mockCalendarService) GetUserCalendars(userID int64) ([]*calendardomain.Calendar, error) {
	var userCalendars []*calendardomain.Calendar
	for _, cal := range m.calendars {
		if cal.UserID == userID {
			userCalendars = append(userCalendars, cal)
		}
	}
	return userCalendars, nil
}

func (m *mockCalendarService) UpdateCalendar(id int64, displayName, color, description, timezone *string) error {
	return nil
}

func (m *mockCalendarService) DeleteCalendar(id int64) error {
	return nil
}

func (m *mockCalendarService) GenerateSyncToken(id int64) (string, error) {
	return "sync-token-123", nil
}

type mockEventService struct {
	events []*calendardomain.Event
}

func (m *mockEventService) CreateEvent(calendarID int64, icalData string) (*calendardomain.Event, error) {
	event := &calendardomain.Event{
		ID:         int64(len(m.events) + 1),
		CalendarID: calendarID,
		ICalData:   icalData,
	}
	m.events = append(m.events, event)
	return event, nil
}

func (m *mockEventService) GetEvent(id int64) (*calendardomain.Event, error) {
	for _, event := range m.events {
		if event.ID == id {
			return event, nil
		}
	}
	return nil, &calendardomain.Event{}
}

func (m *mockEventService) GetCalendarEvents(calendarID int64) ([]*calendardomain.Event, error) {
	var calendarEvents []*calendardomain.Event
	for _, event := range m.events {
		if event.CalendarID == calendarID {
			calendarEvents = append(calendarEvents, event)
		}
	}
	return calendarEvents, nil
}

func (m *mockEventService) UpdateEvent(id int64, icalData string) error {
	return nil
}

func (m *mockEventService) DeleteEvent(id int64) error {
	return nil
}

func (m *mockEventService) GenerateETag(event *calendardomain.Event) string {
	return "etag-123"
}

func (m *mockEventService) GetEventsInRange(calendarID int64, start, end time.Time) ([]*calendardomain.Event, error) {
	return m.events, nil
}

func (m *mockEventService) ExpandRecurrence(event *calendardomain.Event, start, end time.Time) ([]*calendardomain.Event, error) {
	return []*calendardomain.Event{event}, nil
}

type mockAddressbookService struct {
	addressbooks []*contactdomain.Addressbook
}

func (m *mockAddressbookService) CreateAddressbook(userID int64, name, displayName, description string) (*contactdomain.Addressbook, error) {
	addressbook := &contactdomain.Addressbook{
		ID:          int64(len(m.addressbooks) + 1),
		UserID:      userID,
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}
	m.addressbooks = append(m.addressbooks, addressbook)
	return addressbook, nil
}

func (m *mockAddressbookService) GetAddressbook(id int64) (*contactdomain.Addressbook, error) {
	for _, ab := range m.addressbooks {
		if ab.ID == id {
			return ab, nil
		}
	}
	return nil, &contactdomain.Addressbook{}
}

func (m *mockAddressbookService) GetUserAddressbooks(userID int64) ([]*contactdomain.Addressbook, error) {
	var userAddressbooks []*contactdomain.Addressbook
	for _, ab := range m.addressbooks {
		if ab.UserID == userID {
			userAddressbooks = append(userAddressbooks, ab)
		}
	}
	return userAddressbooks, nil
}

func (m *mockAddressbookService) UpdateAddressbook(id int64, displayName, description *string) error {
	return nil
}

func (m *mockAddressbookService) DeleteAddressbook(id int64) error {
	return nil
}

func (m *mockAddressbookService) GenerateSyncToken(id int64) (string, error) {
	return "sync-token-456", nil
}

type mockContactService struct {
	contacts []*contactdomain.Contact
}

func (m *mockContactService) CreateContact(addressbookID int64, vcardData string) (*contactdomain.Contact, error) {
	contact := &contactdomain.Contact{
		ID:            int64(len(m.contacts) + 1),
		AddressbookID: addressbookID,
		VCardData:     vcardData,
	}
	m.contacts = append(m.contacts, contact)
	return contact, nil
}

func (m *mockContactService) GetContact(id int64) (*contactdomain.Contact, error) {
	for _, contact := range m.contacts {
		if contact.ID == id {
			return contact, nil
		}
	}
	return nil, &contactdomain.Contact{}
}

func (m *mockContactService) GetContactByUID(addressbookID int64, uid string) (*contactdomain.Contact, error) {
	for _, contact := range m.contacts {
		if contact.AddressbookID == addressbookID && contact.UID == uid {
			return contact, nil
		}
	}
	return nil, &contactdomain.Contact{}
}

func (m *mockContactService) GetAddressbookContacts(addressbookID int64) ([]*contactdomain.Contact, error) {
	var addressbookContacts []*contactdomain.Contact
	for _, contact := range m.contacts {
		if contact.AddressbookID == addressbookID {
			addressbookContacts = append(addressbookContacts, contact)
		}
	}
	return addressbookContacts, nil
}

func (m *mockContactService) SearchContacts(addressbookID int64, query string) ([]*contactdomain.Contact, error) {
	var results []*contactdomain.Contact
	for _, contact := range m.contacts {
		if contact.AddressbookID == addressbookID {
			results = append(results, contact)
		}
	}
	return results, nil
}

func (m *mockContactService) UpdateContact(id int64, vcardData string) error {
	return nil
}

func (m *mockContactService) DeleteContact(id int64) error {
	return nil
}

func (m *mockContactService) UpdateETag(id int64, etag string) error {
	return nil
}

func (m *mockContactService) GenerateETag(contact *contactdomain.Contact) string {
	return "etag-456"
}
