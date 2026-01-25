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
	case "HEAD":
		h.handleHeadGet(w, r)
	case "GET":
		h.handleHeadGet(w, r)
	case "PUT":
		h.handlePut(w, r)
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

func (h *Handler) handleProppatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read PROPPATCH body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var propupdate PropertyUpdate
	if err := xml.Unmarshal(body, &propupdate); err != nil {
		h.logger.Error("failed to parse PROPPATCH request", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	multistatus := &MultiStatus{
		Responses: []Response{
			{
				Href: r.URL.Path,
				PropStats: []PropStat{
					{
						Prop:   PropValue{},
						Status: "HTTP/1.1 200 OK",
					},
				},
			},
		},
	}

	xmlData, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		h.logger.Error("failed to marshal PROPPATCH response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}

func (h *Handler) handleMkcol(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if r.ContentLength > 0 {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error("failed to read MKCOL body", zap.Error(err))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(body) > 0 {
			h.logger.Warn("MKCOL request with unexpected body", zap.Int("length", len(body)))
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}
	}

	h.logger.Info("Creating collection", zap.String("path", r.URL.Path))

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if r.ContentLength > 0 {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error("failed to read DELETE body", zap.Error(err))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(body) > 0 {
			h.logger.Warn("DELETE request with unexpected body", zap.Int("length", len(body)))
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}
	}

	h.logger.Info("Deleting resource", zap.String("path", r.URL.Path))

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleCopy(w http.ResponseWriter, r *http.Request) {
	destination := r.Header.Get("Destination")
	if destination == "" {
		http.Error(w, "Destination header required", http.StatusBadRequest)
		return
	}

	overwrite := r.Header.Get("Overwrite") == "T"

	h.logger.Info("Copying resource",
		zap.String("source", r.URL.Path),
		zap.String("destination", destination),
		zap.Bool("overwrite", overwrite),
	)

	if overwrite {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) handleMove(w http.ResponseWriter, r *http.Request) {
	destination := r.Header.Get("Destination")
	if destination == "" {
		http.Error(w, "Destination header required", http.StatusBadRequest)
		return
	}

	overwrite := r.Header.Get("Overwrite") == "T"

	h.logger.Info("Moving resource",
		zap.String("source", r.URL.Path),
		zap.String("destination", destination),
		zap.Bool("overwrite", overwrite),
	)

	if overwrite {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("OPTIONS request", zap.String("path", r.URL.Path))

	capabilities := []string{
		"1",
		"2",
		"3",
		"access-control",
		"calendar-access",
		"addressbook-access",
		"calendarserver-home",
		"calendarserver-user-address-set",
	}

	allowedMethods := []string{
		"OPTIONS",
		"HEAD",
		"GET",
		"PROPFIND",
		"PROPPATCH",
		"MKCOL",
		"DELETE",
		"COPY",
		"MOVE",
		"PUT",
	}

	w.Header().Set("DAV", strings.Join(capabilities, ", "))
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	w.Header().Set("MS-Author-Via", "DAV")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleHeadGet(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("HEAD/GET request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	if r.Method == "HEAD" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resourceType := h.getResourceType(r.URL.Path)
	isCollection := resourceType == "collection" || resourceType == "calendar" || resourceType == "addressbook" || resourceType == "principal"

	if isCollection {
		http.Error(w, "Method not allowed for collections", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handlePut(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("PUT request", zap.String("path", r.URL.Path), zap.Int64("content-length", r.ContentLength))

	if r.URL.Path == "" || r.URL.Path == "/" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read PUT body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 && r.ContentLength > 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	h.logger.Info("Storing resource",
		zap.String("path", r.URL.Path),
		zap.Int("size", len(body)),
		zap.String("content-type", contentType),
	)

	w.WriteHeader(http.StatusCreated)
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
		h.addCollectionChildren(multistatus, cleanPath)
	}

	return multistatus
}

func (h *Handler) addCollectionChildren(multistatus *MultiStatus, parentPath string) {
	resourceType := h.getResourceType(parentPath)
	isCollection := resourceType == "collection" || resourceType == "calendar" || resourceType == "addressbook"

	if !isCollection {
		return
	}

	var children []string
	switch resourceType {
	case "collection":
		if strings.Contains(parentPath, "/calendars/") {
			children = []string{parentPath + "/calendar1", parentPath + "/calendar2"}
		} else if strings.Contains(parentPath, "/addressbooks/") {
			children = []string{parentPath + "/contacts", parentPath + "/colleagues"}
		} else {
			children = []string{parentPath + "/calendars", parentPath + "/addressbooks"}
		}
	case "calendar":
		children = []string{parentPath + "/event1.ics", parentPath + "/event2.ics"}
	case "addressbook":
		children = []string{parentPath + "/contact1.vcf", parentPath + "/contact2.vcf"}
	}

	for _, childPath := range children {
		response := Response{
			Href: childPath,
			PropStats: []PropStat{
				{
					Prop:   h.buildPropValue(childPath, &PropFind{AllProp: &struct{}{}}),
					Status: "HTTP/1.1 200 OK",
				},
			},
		}
		multistatus.Responses = append(multistatus.Responses, response)
	}
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

// TODO: buildMultiStatusResponse builds a multistatus response

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

func (h *Handler) getDisplayName(urlPath string) string {
	return path.Base(urlPath)
}

func (h *Handler) generateETag(urlPath string) string {
	return `"` + urlPath + `"`
}

func (h *Handler) getCalendarDescription(urlPath string) string {
	if strings.Contains(urlPath, "personal") {
		return "Personal Calendar"
	}
	if strings.Contains(urlPath, "work") {
		return "Work Calendar"
	}
	return ""
}

func (h *Handler) getCalendarColor(urlPath string) string {
	if strings.Contains(urlPath, "personal") {
		return "#007AFF"
	}
	if strings.Contains(urlPath, "work") {
		return "#FF3B30"
	}
	return ""
}

func (h *Handler) getCalendarOrder(urlPath string) int {
	if strings.Contains(urlPath, "personal") {
		return 1
	}
	if strings.Contains(urlPath, "work") {
		return 2
	}
	return 0
}

func (h *Handler) getAddressbookDescription(urlPath string) string {
	if strings.Contains(urlPath, "contacts") {
		return "Personal Contacts"
	}
	if strings.Contains(urlPath, "colleagues") {
		return "Work Colleagues"
	}
	return ""
}

// getCurrentTime returns current time (helper for testing)
func (h *Handler) getCurrentTime() time.Time {
	return time.Now()
}
