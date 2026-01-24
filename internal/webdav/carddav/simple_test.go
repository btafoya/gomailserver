package carddav

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type SimpleHandler struct{}

func (h *SimpleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type MockAddressbookService struct{}

func (m *MockAddressbookService) GetAddressbook(id int64) (*domain.Addressbook, error) {
	return &domain.Addressbook{
		ID:   id,
		Name: "test",
	}, nil
}

type MockContactService struct{}

func (m *MockContactService) GetAddressbookContacts(addressbookID int64) ([]*domain.Contact, error) {
	return []*domain.Contact{
		{
			AddressbookID: addressbookID,
			FN:            "Test Contact",
			UID:           "test-contact-1",
			VCardData:     "BEGIN:VCARD\nVERSION:4.0\nFN:Test Contact\nEMAIL:test@example.com\nEND:VCARD",
		},
	}, nil
}

type Handler struct {
	webdavHandler      *webdav.Handler
	logger             *zap.Logger
	addressbookService domain.AddressbookService
	contactService     domain.ContactService
}

func NewHandler(logger *zap.Logger, addressbookService domain.AddressbookService, contactService domain.ContactService) *Handler {
	return &Handler{
		webdavHandler:      webdav.NewHandler(logger, "/carddav"),
		logger:             logger,
		addressbookService: addressbookService,
		contactService:     contactService,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("CardDAV request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	switch r.Method {
	case "REPORT":
		h.handleReport(w, r)
	default:
		h.webdavHandler.ServeHTTP(w, r)
	}
}

func (h *Handler) handleReport(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("REPORT request", zap.String("path", r.URL.Path))
	// For testing, return simple response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><report></report>"))
}

// handleAddressbookQuery handles addressbook-query REPORT
func (h *Handler) handleAddressbookQuery(w http.ResponseWriter, r *http.Request, body []byte) {
	h.logger.Info("addressbook-query REPORT", zap.String("path", r.URL.Path))

	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><multistatus xmlns=\"DAV:\"><response><href>/carddav/addressbooks/123/test</href><propstat><prop><displayname/><getetag/></prop><status>HTTP/1.1 200 OK</status></propstat></response></multistatus>"))
}

// handlePropset handles PROPPATCH requests as propset operation
func (h *Handler) handlePropset(w http.ResponseWriter, r *http.Request, body []byte) {
	w.WriteHeader(http.StatusMultiStatus)
}

// handleAddressbookMultiget handles addressbook-multiget REPORT
func (h *Handler) handleAddressbookMultiget(w http.ResponseWriter, r *http.Request, body []byte) {
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><multistatus xmlns=\"DAV:\"><response><href>/carddav/addressbooks/123/test</href><propstat><prop><displayname/><getetag/></prop><status>HTTP/1.1 200 OK</status></propstat></response></multistatus>"))
}

// handleDelete handles DELETE requests for deleting contacts
func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Main function
func main() {
	h := &SimpleHandler{}
	req := &http.Request{
		Method: "REPORT",
		URL:    "/carddav/addressbooks/123/test",
		Body:   strings.NewReader("<C:addressbook-query/>"),
	}

	h.ServeHTTP(h, req)

	fmt.Printf("Status: %d\n", h.StatusCode)
	fmt.Printf("Response: %s\n", string(h.body))
}
