package carddav

import (
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/btafoya/gomailserver/internal/contact/domain"
	"github.com/btafoya/gomailserver/internal/webdav"
	"go.uber.org/zap"
)

// Handler handles CardDAV-specific requests
type Handler struct {
	webdavHandler      *webdav.Handler
	logger             *zap.Logger
	addressbookService domain.AddressbookService
	contactService     domain.ContactService
}

// NewHandler creates a new CardDAV handler
func NewHandler(logger *zap.Logger, addressbookService domain.AddressbookService, contactService domain.ContactService) *Handler {
	return &Handler{
		webdavHandler:      webdav.NewHandler(logger, "/carddav"),
		logger:             logger,
		addressbookService: addressbookService,
		contactService:     contactService,
	}
}

// ServeHTTP implements http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("CardDAV request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	switch r.Method {
	case "REPORT":
		h.handleReport(w, r)
	default:
		// Delegate to base WebDAV handler
		h.webdavHandler.ServeHTTP(w, r)
	}
}

// handleReport handles REPORT requests (RFC 3253, extended by CardDAV)
func (h *Handler) handleReport(w http.ResponseWriter, r *http.Request) {
	// Parse request body to determine report type
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read REPORT body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Detect report type from XML
	var report struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal(body, &report); err != nil {
		h.logger.Error("failed to parse REPORT request", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	switch report.XMLName.Local {
	case "addressbook-query":
		h.handleAddressbookQuery(w, r, body)
	case "addressbook-multiget":
		h.handleAddressbookMultiget(w, r, body)
	default:
		h.logger.Warn("unknown REPORT type", zap.String("type", report.XMLName.Local))
		http.Error(w, "Unsupported report type", http.StatusBadRequest)
	}
}

// handleAddressbookQuery handles addressbook-query REPORT
func (h *Handler) handleAddressbookQuery(w http.ResponseWriter, r *http.Request, body []byte) {
	h.logger.Info("addressbook-query REPORT", zap.String("path", r.URL.Path))

	// Extract addressbook ID from path
	// Path format: /carddav/addressbooks/{userID}/{addressbookName}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid addressbook path", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[2]
	addressbookName := pathParts[3]

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get addressbook by user and name
	addressbooks, err := h.addressbookService.GetUserAddressbooks(userID)
	if err != nil {
		h.logger.Error("failed to get addressbooks", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var addressbookID int64
	for _, ab := range addressbooks {
		if ab.Name == addressbookName {
			addressbookID = ab.ID
			break
		}
	}

	if addressbookID == 0 {
		http.Error(w, "Addressbook not found", http.StatusNotFound)
		return
	}

	// Parse request body for filter
	// This is a simplified implementation - full CardDAV would parse complex filters
	// For now, we'll return all contacts in the addressbook
	contacts, err := h.contactService.GetAddressbookContacts(addressbookID)
	if err != nil {
		h.logger.Error("failed to get contacts", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build multistatus response
	multistatus := &webdav.MultiStatus{
		Responses: []webdav.Response{},
	}

	for _, contact := range contacts {
		contactPath := r.URL.Path + "/" + contact.UID + ".vcf"

		propValue := webdav.PropValue{
			GetETag:     &contact.ETag,
			AddressData: &contact.VCardData,
		}

		response := webdav.Response{
			Href: contactPath,
			PropStats: []webdav.PropStat{
				{
					Prop:   propValue,
					Status: "HTTP/1.1 200 OK",
				},
			},
		}
		multistatus.Responses = append(multistatus.Responses, response)
	}

	// Marshal and send response
	xmlData, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		h.logger.Error("failed to marshal multistatus", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}

// handleAddressbookMultiget handles addressbook-multiget REPORT
func (h *Handler) handleAddressbookMultiget(w http.ResponseWriter, r *http.Request, body []byte) {
	h.logger.Info("addressbook-multiget REPORT", zap.String("path", r.URL.Path))

	// Extract addressbook ID from path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid addressbook path", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[2]
	addressbookName := pathParts[3]

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get addressbook by user and name
	addressbooks, err := h.addressbookService.GetUserAddressbooks(userID)
	if err != nil {
		h.logger.Error("failed to get addressbooks", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var addressbookID int64
	for _, ab := range addressbooks {
		if ab.Name == addressbookName {
			addressbookID = ab.ID
			break
		}
	}

	if addressbookID == 0 {
		http.Error(w, "Addressbook not found", http.StatusNotFound)
		return
	}

	// Parse request body for hrefs
	// This is a simplified implementation - would parse <C:addressbook-multiget> XML with <D:href> elements
	// For now, we'll return all contacts
	contacts, err := h.contactService.GetAddressbookContacts(addressbookID)
	if err != nil {
		h.logger.Error("failed to get contacts", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build multistatus response
	multistatus := &webdav.MultiStatus{
		Responses: []webdav.Response{},
	}

	for _, contact := range contacts {
		contactPath := r.URL.Path + "/" + contact.UID + ".vcf"

		propValue := webdav.PropValue{
			GetETag:     &contact.ETag,
			AddressData: &contact.VCardData,
		}

		response := webdav.Response{
			Href: contactPath,
			PropStats: []webdav.PropStat{
				{
					Prop:   propValue,
					Status: "HTTP/1.1 200 OK",
				},
			},
		}
		multistatus.Responses = append(multistatus.Responses, response)
	}

	// Marshal and send response
	xmlData, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		h.logger.Error("failed to marshal multistatus", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}
