package webdav

import (
	"encoding/xml"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Handler implements WebDAV HTTP methods
type Handler struct {
	logger   *zap.Logger
	basePath string
}

// NewHandler creates a new WebDAV handler
func NewHandler(logger *zap.Logger, basePath string) *Handler {
	return &Handler{
		logger:   logger,
		basePath: basePath,
	}
}

// ServeHTTP implements http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("WebDAV request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("depth", r.Header.Get("Depth")),
	)

	switch r.Method {
	case "PROPFIND":
		h.handlePropfind(w, r)
	case "PROPPATCH":
		h.handleProppatch(w, r)
	case "MKCOL":
		h.handleMkcol(w, r)
	case "DELETE":
		h.handleDelete(w, r)
	case "COPY":
		h.handleCopy(w, r)
	case "MOVE":
		h.handleMove(w, r)
	case "OPTIONS":
		h.handleOptions(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePropfind handles PROPFIND requests
func (h *Handler) handlePropfind(w http.ResponseWriter, r *http.Request) {
	depth := r.Header.Get("Depth")
	if depth == "" {
		depth = "0"
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read PROPFIND body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var propfind PropFind
	if len(body) > 0 {
		if err := xml.Unmarshal(body, &propfind); err != nil {
			h.logger.Error("failed to parse PROPFIND request", zap.Error(err))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	} else {
		// If no body, treat as allprop
		propfind.AllProp = &struct{}{}
	}

	// Build multistatus response
	multistatus := h.buildMultiStatus(r.URL.Path, &propfind, depth)

	// Marshal response
	xmlData, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		h.logger.Error("failed to marshal multistatus response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}

// buildMultiStatus builds a multistatus response for PROPFIND
func (h *Handler) buildMultiStatus(urlPath string, propfind *PropFind, depth string) *MultiStatus {
	multistatus := &MultiStatus{
		Responses: []Response{},
	}

	// Clean path
	cleanPath := path.Clean(urlPath)

	// Add response for requested resource
	response := h.buildResponse(cleanPath, propfind)
	multistatus.Responses = append(multistatus.Responses, response)

	// Handle depth
	if depth == "1" || depth == "infinity" {
		// For now, just return the root resource
		// TODO: Implement collection enumeration
	}

	return multistatus
}

// buildResponse builds a response for a single resource
func (h *Handler) buildResponse(urlPath string, propfind *PropFind) Response {
	response := Response{
		Href: urlPath,
		PropStats: []PropStat{
			{
				Prop:   h.buildPropValue(urlPath, propfind),
				Status: "HTTP/1.1 200 OK",
			},
		},
	}

	return response
}

// buildPropValue builds property values based on the requested properties
func (h *Handler) buildPropValue(urlPath string, propfind *PropFind) PropValue {
	propValue := PropValue{}

	// Determine what properties to return
	returnAll := propfind.AllProp != nil
	returnNames := propfind.PropName != nil

	if returnNames {
		// Just return property names, not values
		return propValue
	}

	prop := propfind.Prop
	if prop == nil && returnAll {
		// Create a prop with all properties requested
		prop = &Prop{
			ResourceType:         &struct{}{},
			DisplayName:          &struct{}{},
			GetContentType:       &struct{}{},
			GetETag:              &struct{}{},
			GetLastModified:      &struct{}{},
			GetContentLength:     &struct{}{},
			CreationDate:         &struct{}{},
			CurrentUserPrincipal: &struct{}{},
		}
	}

	if prop == nil {
		return propValue
	}

	// Determine resource type based on path
	resourceType := h.getResourceType(urlPath)
	isCollection := resourceType == "collection" || resourceType == "calendar" || resourceType == "addressbook" || resourceType == "principal"

	// Build property values
	if prop.ResourceType != nil {
		rt := &ResourceType{}
		switch resourceType {
		case "collection":
			rt.Collection = &struct{}{}
		case "calendar":
			rt.Collection = &struct{}{}
			rt.Calendar = &struct{}{}
		case "addressbook":
			rt.Collection = &struct{}{}
			rt.Addressbook = &struct{}{}
		case "principal":
			rt.Principal = &struct{}{}
		}
		propValue.ResourceType = rt
	}

	if prop.DisplayName != nil {
		displayName := h.getDisplayName(urlPath)
		propValue.DisplayName = &displayName
	}

	if prop.GetContentType != nil && !isCollection {
		contentType := "application/octet-stream"
		if strings.HasSuffix(urlPath, ".ics") {
			contentType = "text/calendar; charset=utf-8"
		} else if strings.HasSuffix(urlPath, ".vcf") {
			contentType = "text/vcard; charset=utf-8"
		}
		propValue.GetContentType = &contentType
	}

	if prop.GetETag != nil {
		etag := h.generateETag(urlPath)
		propValue.GetETag = &etag
	}

	if prop.GetLastModified != nil {
		// TODO: Get actual modification time from storage
		lastModified := FormatHTTPDate(h.getCurrentTime())
		propValue.GetLastModified = &lastModified
	}

	if prop.GetContentLength != nil && !isCollection {
		// TODO: Get actual content length from storage
		var length int64 = 0
		propValue.GetContentLength = &length
	}

	if prop.CreationDate != nil {
		// TODO: Get actual creation time from storage
		creationDate := FormatISO8601(h.getCurrentTime())
		propValue.CreationDate = &creationDate
	}

	if prop.CurrentUserPrincipal != nil {
		// TODO: Get actual user principal from authentication context
		propValue.CurrentUserPrincipal = &Href{Href: "/principals/users/admin"}
	}

	// CalDAV specific properties
	if prop.CalendarHomeSet != nil {
		// TODO: Get actual calendar home set from user context
		propValue.CalendarHomeSet = &Href{Href: "/caldav/calendars/admin"}
	}

	if prop.CalendarDescription != nil {
		description := h.getCalendarDescription(urlPath)
		if description != "" {
			propValue.CalendarDescription = &description
		}
	}

	if prop.CalendarColor != nil {
		color := h.getCalendarColor(urlPath)
		if color != "" {
			propValue.CalendarColor = &color
		}
	}

	if prop.CalendarOrder != nil {
		order := h.getCalendarOrder(urlPath)
		propValue.CalendarOrder = &order
	}

	if prop.SupportedCalendarComponentSet != nil {
		propValue.SupportedCalendarComponentSet = &SupportedCalendarComponentSet{
			Components: []CalendarComponent{
				{Name: "VEVENT"},
				{Name: "VTODO"},
				{Name: "VJOURNAL"},
			},
		}
	}

	// CardDAV specific properties
	if prop.AddressbookHomeSet != nil {
		// TODO: Get actual addressbook home set from user context
		propValue.AddressbookHomeSet = &Href{Href: "/carddav/addressbooks/admin"}
	}

	if prop.AddressbookDescription != nil {
		description := h.getAddressbookDescription(urlPath)
		if description != "" {
			propValue.AddressbookDescription = &description
		}
	}

	if prop.SupportedAddressData != nil {
		propValue.SupportedAddressData = &SupportedAddressData{
			AddressDataTypes: []AddressDataType{
				{ContentType: "text/vcard", Version: "3.0"},
				{ContentType: "text/vcard", Version: "4.0"},
			},
		}
	}

	return propValue
}

// getResourceType determines the resource type based on the path
func (h *Handler) getResourceType(urlPath string) string {
	// TODO: Implement actual resource type detection from storage
	if strings.HasPrefix(urlPath, "/principals/") {
		return "principal"
	}
	if strings.Contains(urlPath, "/calendars/") {
		if strings.HasSuffix(urlPath, ".ics") {
			return "event"
		}
		return "calendar"
	}
	if strings.Contains(urlPath, "/addressbooks/") {
		if strings.HasSuffix(urlPath, ".vcf") {
			return "contact"
		}
		return "addressbook"
	}
	return "collection"
}

// getDisplayName returns the display name for a resource
func (h *Handler) getDisplayName(urlPath string) string {
	// TODO: Get actual display name from storage
	return path.Base(urlPath)
}

// generateETag generates an ETag for a resource
func (h *Handler) generateETag(urlPath string) string {
	// TODO: Generate actual ETag based on resource content
	return `"` + urlPath + `"`
}

// getCalendarDescription returns the calendar description
func (h *Handler) getCalendarDescription(urlPath string) string {
	// TODO: Get from storage
	return ""
}

// getCalendarColor returns the calendar color
func (h *Handler) getCalendarColor(urlPath string) string {
	// TODO: Get from storage
	return ""
}

// getCalendarOrder returns the calendar order
func (h *Handler) getCalendarOrder(urlPath string) int {
	// TODO: Get from storage
	return 0
}

// getAddressbookDescription returns the addressbook description
func (h *Handler) getAddressbookDescription(urlPath string) string {
	// TODO: Get from storage
	return ""
}

// getCurrentTime returns current time (helper for testing)
func (h *Handler) getCurrentTime() time.Time {
	return time.Now()
}

// handleProppatch handles PROPPATCH requests
func (h *Handler) handleProppatch(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement PROPPATCH
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// handleMkcol handles MKCOL requests
func (h *Handler) handleMkcol(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement MKCOL
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// handleDelete handles DELETE requests
func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement DELETE
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// handleCopy handles COPY requests
func (h *Handler) handleCopy(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement COPY
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// handleMove handles MOVE requests
func (h *Handler) handleMove(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement MOVE
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// handleOptions handles OPTIONS requests
func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "OPTIONS, PROPFIND, PROPPATCH, MKCOL, DELETE, COPY, MOVE, GET, PUT")
	w.Header().Set("DAV", "1, 2, calendar-access, addressbook")
	w.WriteHeader(http.StatusOK)
}
