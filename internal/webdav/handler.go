package webdav

import (
	"context"
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
	storage  Storage
}

// NewHandler creates a new WebDAV handler
func NewHandler(logger *zap.Logger, basePath string) *Handler {
	return &Handler{
		logger:   logger,
		basePath: basePath,
		storage:  nil, // Will use mock data if storage not set
	}
}

// NewHandlerWithStorage creates a new WebDAV handler with storage backend
func NewHandlerWithStorage(logger *zap.Logger, basePath string, storage Storage) *Handler {
	return &Handler{
		logger:   logger,
		basePath: basePath,
		storage:  storage,
	}
}

// SetStorage sets the storage backend
func (h *Handler) SetStorage(storage Storage) {
	h.storage = storage
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	if userID, ok := ctx.Value(UserIDKey).(int64); ok {
		return userID, true
	}
	return 0, false
}

// GetUsernameFromContext extracts username from context
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	if username, ok := ctx.Value(UsernameKey).(string); ok {
		return username, true
	}
	return "", false
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
	multistatus := h.buildMultiStatus(r.URL.Path, &propfind, depth, r.Context())

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

	// Use storage if available
	if h.storage != nil {
		if h.storage.Exists(r.URL.Path) {
			http.Error(w, "Resource already exists", http.StatusMethodNotAllowed)
			return
		}
		if err := h.storage.CreateCollection(r.URL.Path); err != nil {
			h.logger.Error("failed to create collection", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	// Use storage if available
	if h.storage != nil {
		if !h.storage.Exists(r.URL.Path) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if err := h.storage.DeleteResource(r.URL.Path); err != nil {
			h.logger.Error("failed to delete resource", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	overwrite := r.Header.Get("Overwrite") != "F"

	// Parse destination URL to get path
	destPath := h.parseDestinationPath(destination)

	// Use storage if available
	if h.storage != nil {
		if !h.storage.Exists(r.URL.Path) {
			http.Error(w, "Source not found", http.StatusNotFound)
			return
		}
		destExists := h.storage.Exists(destPath)
		if destExists && !overwrite {
			http.Error(w, "Destination exists", http.StatusPreconditionFailed)
			return
		}
		if err := h.storage.CopyResource(r.URL.Path, destPath, overwrite); err != nil {
			h.logger.Error("failed to copy resource", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if destExists {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		return
	}

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

	overwrite := r.Header.Get("Overwrite") != "F"

	// Parse destination URL to get path
	destPath := h.parseDestinationPath(destination)

	// Use storage if available
	if h.storage != nil {
		if !h.storage.Exists(r.URL.Path) {
			http.Error(w, "Source not found", http.StatusNotFound)
			return
		}
		destExists := h.storage.Exists(destPath)
		if destExists && !overwrite {
			http.Error(w, "Destination exists", http.StatusPreconditionFailed)
			return
		}
		if err := h.storage.MoveResource(r.URL.Path, destPath, overwrite); err != nil {
			h.logger.Error("failed to move resource", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if destExists {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		return
	}

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

	// Use storage if available
	if h.storage != nil {
		resourceInfo, err := h.storage.GetResourceInfo(r.URL.Path)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if resourceInfo.IsCollection {
			http.Error(w, "Method not allowed for collections", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", resourceInfo.ContentType)
		w.Header().Set("Content-Length", string(rune(resourceInfo.ContentLen)))
		w.Header().Set("ETag", resourceInfo.ETag)
		w.Header().Set("Last-Modified", FormatHTTPDate(resourceInfo.ModTime))

		if r.Method == "HEAD" {
			w.WriteHeader(http.StatusOK)
			return
		}

		content, err := h.storage.ReadResource(r.URL.Path)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer content.Close()

		w.WriteHeader(http.StatusOK)
		io.Copy(w, content)
		return
	}

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

	// Use storage if available
	if h.storage != nil {
		if err := h.storage.WriteResource(r.URL.Path, r.Body); err != nil {
			h.logger.Error("failed to write resource", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
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

// parseDestinationPath extracts the path from a Destination header URL
func (h *Handler) parseDestinationPath(destination string) string {
	// Handle full URLs like "http://host/path"
	if strings.Contains(destination, "://") {
		idx := strings.Index(destination, "://")
		rest := destination[idx+3:]
		if slashIdx := strings.Index(rest, "/"); slashIdx >= 0 {
			return rest[slashIdx:]
		}
	}
	return destination
}

// buildMultiStatus builds a multistatus response for PROPFIND
func (h *Handler) buildMultiStatus(urlPath string, propfind *PropFind, depth string, ctx context.Context) *MultiStatus {
	multistatus := &MultiStatus{
		Responses: []Response{},
	}

	// Clean path
	cleanPath := path.Clean(urlPath)

	// Add response for requested resource
	response := h.buildResponse(cleanPath, propfind, ctx)
	multistatus.Responses = append(multistatus.Responses, response)

	// Handle depth
	if depth == "1" || depth == "infinity" {
		h.addCollectionChildren(multistatus, cleanPath, propfind, ctx, depth == "infinity")
	}

	return multistatus
}

func (h *Handler) addCollectionChildren(multistatus *MultiStatus, parentPath string, propfind *PropFind, ctx context.Context, recursive bool) {
	// Use storage if available
	if h.storage != nil {
		children, err := h.storage.ListChildren(parentPath)
		if err != nil {
			h.logger.Warn("failed to list children", zap.String("path", parentPath), zap.Error(err))
			return
		}

		for _, child := range children {
			response := h.buildResponse(child.Path, propfind, ctx)
			multistatus.Responses = append(multistatus.Responses, response)

			if recursive && child.IsCollection {
				h.addCollectionChildren(multistatus, child.Path, propfind, ctx, true)
			}
		}
		return
	}

	// Fallback to mock data
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
					Prop:   h.buildPropValue(childPath, &PropFind{AllProp: &struct{}{}}, ctx),
					Status: "HTTP/1.1 200 OK",
				},
			},
		}
		multistatus.Responses = append(multistatus.Responses, response)
	}
}

// buildResponse builds a response for a single resource
func (h *Handler) buildResponse(urlPath string, propfind *PropFind, ctx context.Context) Response {
	response := Response{
		Href: urlPath,
		PropStats: []PropStat{
			{
				Prop:   h.buildPropValue(urlPath, propfind, ctx),
				Status: "HTTP/1.1 200 OK",
			},
		},
	}
	return response
}

// buildPropValue builds property values based on the requested properties
func (h *Handler) buildPropValue(urlPath string, propfind *PropFind, ctx context.Context) PropValue {
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

	// Get resource info from storage if available
	var resourceInfo *ResourceInfo
	if h.storage != nil {
		var err error
		resourceInfo, err = h.storage.GetResourceInfo(urlPath)
		if err != nil {
			h.logger.Warn("failed to get resource info", zap.String("path", urlPath), zap.Error(err))
		}
	}

	// Determine resource type based on path or storage info
	var resourceType string
	var isCollection bool
	if resourceInfo != nil {
		resourceType = resourceInfo.ResourceKind
		isCollection = resourceInfo.IsCollection
	} else {
		resourceType = h.getResourceType(urlPath)
		isCollection = resourceType == "collection" || resourceType == "calendar" || resourceType == "addressbook" || resourceType == "principal"
	}

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
		var contentType string
		if resourceInfo != nil {
			contentType = resourceInfo.ContentType
		} else {
			contentType = "application/octet-stream"
			if strings.HasSuffix(urlPath, ".ics") {
				contentType = "text/calendar; charset=utf-8"
			} else if strings.HasSuffix(urlPath, ".vcf") {
				contentType = "text/vcard; charset=utf-8"
			}
		}
		propValue.GetContentType = &contentType
	}

	if prop.GetETag != nil {
		var etag string
		if resourceInfo != nil {
			etag = resourceInfo.ETag
		} else {
			etag = h.generateETag(urlPath)
		}
		propValue.GetETag = &etag
	}

	if prop.GetLastModified != nil {
		var lastModified string
		if resourceInfo != nil {
			lastModified = FormatHTTPDate(resourceInfo.ModTime)
		} else {
			lastModified = FormatHTTPDate(time.Now())
		}
		propValue.GetLastModified = &lastModified
	}

	if prop.GetContentLength != nil && !isCollection {
		var length int64
		if resourceInfo != nil {
			length = resourceInfo.ContentLen
		}
		propValue.GetContentLength = &length
	}

	if prop.CreationDate != nil {
		var creationDate string
		if resourceInfo != nil {
			creationDate = FormatISO8601(resourceInfo.CreateTime)
		} else {
			creationDate = FormatISO8601(time.Now())
		}
		propValue.CreationDate = &creationDate
	}

	if prop.CurrentUserPrincipal != nil {
		// Get user principal from context
		var principalPath string
		if username, ok := GetUsernameFromContext(ctx); ok {
			principalPath = "/principals/users/" + username
		} else {
			principalPath = "/principals/users/admin"
		}
		propValue.CurrentUserPrincipal = &Href{Href: principalPath}
	}

	// CalDAV specific properties
	if prop.CalendarHomeSet != nil {
		var calendarHomePath string
		if username, ok := GetUsernameFromContext(ctx); ok {
			calendarHomePath = "/caldav/calendars/" + username
		} else {
			calendarHomePath = "/caldav/calendars/admin"
		}
		propValue.CalendarHomeSet = &Href{Href: calendarHomePath}
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
		var addressbookHomePath string
		if username, ok := GetUsernameFromContext(ctx); ok {
			addressbookHomePath = "/carddav/addressbooks/" + username
		} else {
			addressbookHomePath = "/carddav/addressbooks/admin"
		}
		propValue.AddressbookHomeSet = &Href{Href: addressbookHomePath}
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
