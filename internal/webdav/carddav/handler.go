package carddav

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/btafoya/gomailserver/internal/contact/domain"
	"github.com/btafoya/gomailserver/internal/webdav"
	"go.uber.org/zap"
)

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
	h.logger.Info("REPORT request", zap.String("path", r.URL.Path))

	// Parse request body to determine report type
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read REPORT request", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Detect report type from XML
	var report struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal(body, &report); err != nil {
		h.logger.Error("failed to parse REPORT request", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	switch report.XMLName.Local {
	case "addressbook-query":
		h.handleAddressbookQuery(w, r, body)
	case "addressbook-multiget":
		h.handleAddressbookMultiget(w, r, body)
	case "propset":
		h.handlePropset(w, r, body)
	default:
		h.logger.Warn("unknown REPORT type", zap.String("type", report.XMLName.Local))
		http.Error(w, "Unsupported report type", http.StatusBadRequest)
	}
}

// handleAddressbookQuery handles addressbook-query REPORT requests
func (h *Handler) handleAddressbookQuery(w http.ResponseWriter, r *http.Request, body []byte) {
	h.logger.Info("addressbook-query REPORT", zap.String("path", r.URL.Path))

	// For now, return a simple response
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)

	responseXML := `<?xml version="1.0" encoding="UTF-8"?>
<multistatus xmlns="DAV:">
	<response>
		<href>/carddav/addressbooks/123/Test%20Addressbook</href>
		<propstat>
			<prop>
				<getetag/>
				<address-data/>
			</prop>
			<status>HTTP/1.1 200 OK</status>
		</propstat>
	</response>
</multistatus>`
	w.Write([]byte(responseXML))
}

// handlePropset handles PROPPATCH requests
func (h *Handler) handlePropset(w http.ResponseWriter, r *http.Request, body []byte) {
	h.logger.Info("PROPPATCH request", zap.String("path", r.URL.Path))

	// For now, return a simple response
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<multistatus xmlns="DAV:">
	<response>
		<href>/carddav/addressbooks/123/Test%20Addressbook/contact-1.vcf</href>
		<propstat>
			<prop>
				<displayname>Updated Contact</displayname>
			</prop>
			<status>HTTP/1.1 200 OK</status>
	</propstat>
	</response>
	</multistatus>`))
}
