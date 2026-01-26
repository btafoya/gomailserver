package carddav

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"

	"github.com/btafoya/gomailserver/internal/contact/domain"
)

// Simple mock implementations for testing

type MockAddressbookService struct {
	addressbooks []*domain.Addressbook
	contacts     []*domain.Contact
	shouldError  bool
}

func (m *MockAddressbookService) CreateAddressbook(userID int64, name, displayName, description string) (*domain.Addressbook, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	addressbook := &domain.Addressbook{
		ID:          1,
		UserID:      userID,
		Name:        name,
		DisplayName: displayName,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.addressbooks = append(m.addressbooks, addressbook)
	return addressbook, nil
}

func (m *MockAddressbookService) GetAddressbook(id int64) (*domain.Addressbook, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	for _, ab := range m.addressbooks {
		if ab.ID == id {
			return ab, nil
		}
	}
	return nil, fmt.Errorf("addressbook not found")
}

func (m *MockAddressbookService) GetUserAddressbooks(userID int64) ([]*domain.Addressbook, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	var userAddressbooks []*domain.Addressbook
	for _, ab := range m.addressbooks {
		if ab.UserID == userID {
			userAddressbooks = append(userAddressbooks, ab)
		}
	}
	return userAddressbooks, nil
}

func (m *MockAddressbookService) UpdateAddressbook(id int64, displayName, description *string) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for _, ab := range m.addressbooks {
		if ab.ID == id {
			if displayName != nil {
				ab.DisplayName = *displayName
			}
			if description != nil {
				ab.Description = *description
			}
			return nil
		}
	}
	return fmt.Errorf("addressbook not found")
}

func (m *MockAddressbookService) DeleteAddressbook(id int64) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for i, ab := range m.addressbooks {
		if ab.ID == id {
			m.addressbooks = append(m.addressbooks[:i], m.addressbooks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("addressbook not found")
}

func (m *MockAddressbookService) GenerateSyncToken(id int64) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("mock error")
	}
	return fmt.Sprintf("sync-token-%d", id), nil
}

type MockContactService struct {
	contacts    []*domain.Contact
	shouldError bool
}

func (m *MockContactService) CreateContact(addressbookID int64, vcardData string) (*domain.Contact, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	contact := &domain.Contact{
		ID:            int64(len(m.contacts) + 1),
		AddressbookID: addressbookID,
		UID:           fmt.Sprintf("contact-%d", len(m.contacts)+1),
		FN:            "Mock Contact",
		VCardData:     vcardData,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	m.contacts = append(m.contacts, contact)
	return contact, nil
}

func (m *MockContactService) GetContact(id int64) (*domain.Contact, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	for _, contact := range m.contacts {
		if contact.ID == id {
			return contact, nil
		}
	}
	return nil, fmt.Errorf("contact not found")
}

func (m *MockContactService) GetContactByUID(addressbookID int64, uid string) (*domain.Contact, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	for _, contact := range m.contacts {
		if contact.AddressbookID == addressbookID && contact.UID == uid {
			return contact, nil
		}
	}
	return nil, fmt.Errorf("contact not found")
}

func (m *MockContactService) GetAddressbookContacts(addressbookID int64) ([]*domain.Contact, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	var addressbookContacts []*domain.Contact
	for _, contact := range m.contacts {
		if contact.AddressbookID == addressbookID {
			addressbookContacts = append(addressbookContacts, contact)
		}
	}
	return addressbookContacts, nil
}

func (m *MockContactService) SearchContacts(addressbookID int64, query string) ([]*domain.Contact, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	var results []*domain.Contact
	for _, contact := range m.contacts {
		if contact.AddressbookID == addressbookID && strings.Contains(strings.ToLower(contact.FN), strings.ToLower(query)) {
			results = append(results, contact)
		}
	}
	return results, nil
}

func (m *MockContactService) UpdateContact(id int64, vcardData string) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for _, contact := range m.contacts {
		if contact.ID == id {
			contact.VCardData = vcardData
			contact.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("contact not found")
}

func (m *MockContactService) DeleteContact(id int64) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for i, contact := range m.contacts {
		if contact.ID == id {
			m.contacts = append(m.contacts[:i], m.contacts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("contact not found")
}

func (m *MockContactService) UpdateETag(id int64, etag string) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	for _, contact := range m.contacts {
		if contact.ID == id {
			contact.ETag = etag
			return nil
		}
	}
	return fmt.Errorf("contact not found")
}

func (m *MockContactService) GenerateETag(contact *domain.Contact) string {
	return fmt.Sprintf("etag-%d", contact.ID)
}

// TestCardDAVHandler_ReportAddressbookQuery tests addressbook-query REPORT requests
func TestCardDAVHandler_ReportAddressbookQuery(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockAddressbookService := &MockAddressbookService{shouldError: false}
	mockContactService := &MockContactService{shouldError: false}

	// Setup test data
	testAddressbook := &domain.Addressbook{
		ID:          1,
		UserID:      123,
		Name:        "Test Addressbook",
		DisplayName: "Test Addressbook",
		Description: "Test Description",
	}
	mockAddressbookService.addressbooks = []*domain.Addressbook{testAddressbook}

	testContact := &domain.Contact{
		ID:            1,
		AddressbookID: 1,
		UID:           "test-contact-1",
		FN:            "Test Contact",
		ETag:          "test-etag-1",
		VCardData:     "BEGIN:VCARD\nVERSION:4.0\nFN:Test Contact\nEMAIL:test@example.com\nEND:VCARD",
	}
	mockContactService.contacts = []*domain.Contact{testContact}

	handler := NewHandler(logger, mockAddressbookService, mockContactService)

	// Create test request
	reportBody := `<C:addressbook-query xmlns:C="urn:ietf:params:xml:ns:carddav">
		<D:prop xmlns:D="DAV:">
			<D:getetag/>
			<C:address-data/>
		</D:prop>
	</C:addressbook-query>`

	req := httptest.NewRequest("REPORT", "/carddav/addressbooks/123/Test%20Addressbook", strings.NewReader(reportBody))
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMultiStatus {
		t.Errorf("Expected status %d, got %d", http.StatusMultiStatus, w.Code)
	}
	if !strings.Contains(w.Body.String(), "test-contact-1.vcf") {
		t.Errorf("Expected response to contain contact data")
	}
}

// TestCardDAVHandler_CreateContact tests PUT requests for creating contacts
func TestCardDAVHandler_CreateContact(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockAddressbookService := &MockAddressbookService{shouldError: false}
	mockContactService := &MockContactService{shouldError: false}

	// Setup test data
	testAddressbook := &domain.Addressbook{
		ID:     1,
		UserID: 123,
		Name:   "Test Addressbook",
	}
	mockAddressbookService.addressbooks = []*domain.Addressbook{testAddressbook}

	handler := NewHandler(logger, mockAddressbookService, mockContactService)

	// Create test request
	vcardData := `BEGIN:VCARD
VERSION:4.0
FN:John Doe
EMAIL:john@example.com
TEL:+1234567890
END:VCARD`

	req := httptest.NewRequest("PUT", "/carddav/addressbooks/123/new-contact.vcf", strings.NewReader(vcardData))
	req.Header.Set("Content-Type", "text/vcard")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

// TestCardDAVHandler_GetContact tests GET requests for individual contacts
func TestCardDAVHandler_GetContact(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockAddressbookService := &MockAddressbookService{shouldError: false}
	mockContactService := &MockContactService{shouldError: false}

	// Setup test data
	testContact := &domain.Contact{
		ID:            1,
		AddressbookID: 1,
		UID:           "test-contact-1",
		FN:            "Test Contact",
		ETag:          "test-etag-1",
		VCardData:     "BEGIN:VCARD\nVERSION:4.0\nFN:Test Contact\nEMAIL:test@example.com\nEND:VCARD",
	}
	mockContactService.contacts = []*domain.Contact{testContact}

	handler := NewHandler(logger, mockAddressbookService, mockContactService)

	req := httptest.NewRequest("GET", "/carddav/addressbooks/123/test-contact-1.vcf", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "BEGIN:VCARD") {
		t.Errorf("Expected response to contain vCard data")
	}
}

// TestCardDAVHandler_UpdateContact tests PROPPATCH requests for updating contacts
func TestCardDAVHandler_UpdateContact(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockAddressbookService := &MockAddressbookService{shouldError: false}
	mockContactService := &MockContactService{shouldError: false}

	// Setup test data
	testContact := &domain.Contact{
		ID:            1,
		AddressbookID: 1,
		UID:           "test-contact-1",
		FN:            "Test Contact",
		ETag:          "test-etag-1",
		VCardData:     "BEGIN:VCARD\nVERSION:4.0\nFN:Test Contact\nEMAIL:test@example.com\nEND:VCARD",
	}
	mockContactService.contacts = []*domain.Contact{testContact}

	handler := NewHandler(logger, mockAddressbookService, mockContactService)

	// Create test request
	patchData := `<D:propertyupdate xmlns:D="DAV:">
		<D:set>
		<D:prop>
			<VC:FN xmlns:VC="urn:ietf:params:xml:ns:vcard">John Updated</VC:FN>
		</D:prop>
		</D:set>
	</D:propertyupdate>`

	req := httptest.NewRequest("PROPPATCH", "/carddav/addressbooks/123/test-contact-1.vcf", strings.NewReader(patchData))
	req.Header.Set("Content-Type", "application/xml")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMultiStatus {
		t.Errorf("Expected status %d, got %d", http.StatusMultiStatus, w.Code)
	}
}

// TestCardDAVHandler_DeleteContact tests DELETE requests for contacts
func TestCardDAVHandler_DeleteContact(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockAddressbookService := &MockAddressbookService{shouldError: false}
	mockContactService := &MockContactService{shouldError: false}

	// Setup test data
	testContact := &domain.Contact{
		ID:            1,
		AddressbookID: 1,
		UID:           "test-contact-1",
		FN:            "Test Contact",
	}
	mockContactService.contacts = []*domain.Contact{testContact}

	handler := NewHandler(logger, mockAddressbookService, mockContactService)

	req := httptest.NewRequest("DELETE", "/carddav/addressbooks/123/test-contact-1.vcf", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// 204 No Content is the correct HTTP response for successful DELETE
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

// TestCardDAVHandler_InvalidPath tests handling of invalid paths
func TestCardDAVHandler_InvalidPath(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockAddressbookService := &MockAddressbookService{shouldError: false}
	mockContactService := &MockContactService{shouldError: false}

	handler := NewHandler(logger, mockAddressbookService, mockContactService)

	// Test invalid path (too short) - REPORT with empty body causes EOF error
	// which is logged and returns 400 Bad Request
	req := httptest.NewRequest("REPORT", "/carddav/addressbooks", strings.NewReader(""))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// The handler returns 400 for malformed REPORT requests
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCardDAVHandler_AddressbookNotFound tests handling of non-existent addressbooks
func TestCardDAVHandler_AddressbookNotFound(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockAddressbookService := &MockAddressbookService{shouldError: false}
	mockContactService := &MockContactService{shouldError: false}

	handler := NewHandler(logger, mockAddressbookService, mockContactService)

	// Create test request - no addressbooks set in mock
	// WebDAV REPORT returns 207 Multi-Status even for non-existent resources
	// (with an empty response or error status in the multistatus body)
	reportBody := `<C:addressbook-query xmlns:C="urn:ietf:params:xml:ns:carddav">
		<D:prop xmlns:D="DAV:">
			<D:getetag/>
			<C:address-data/>
		</D:prop>
	</C:addressbook-query>`

	req := httptest.NewRequest("REPORT", "/carddav/addressbooks/456/NonExistent", strings.NewReader(reportBody))
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// WebDAV returns 207 Multi-Status for REPORT requests, even with empty results
	if w.Code != http.StatusMultiStatus {
		t.Errorf("Expected status %d, got %d", http.StatusMultiStatus, w.Code)
	}
}
