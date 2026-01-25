package caldav

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/btafoya/gomailserver/internal/webdav"
	"go.uber.org/zap"
)

// Handler handles CalDAV-specific requests
type Handler struct {
	webdavHandler   *webdav.Handler
	logger          *zap.Logger
	calendarService domain.CalendarService
	eventService    domain.EventService
	taskService     domain.TaskService
	userService     *service.UserService
}

// NewHandler creates a new CalDAV handler
func NewHandler(logger *zap.Logger, calendarService domain.CalendarService, eventService domain.EventService, userService *service.UserService) *Handler {
	return &Handler{
		webdavHandler:   webdav.NewHandler(logger, "/caldav"),
		logger:          logger,
		calendarService: calendarService,
		eventService:    eventService,
		taskService:     nil, // Set via SetTaskService
		userService:     userService,
	}
}

// SetTaskService sets the task service for VTODO support
func (h *Handler) SetTaskService(taskService domain.TaskService) {
	h.taskService = taskService
}

// Permission checking helpers
func (h *Handler) hasReadPermission(userID int64, calendarID int64) bool {
	return h.calendarService.HasReadAccess(userID, calendarID)
}

func (h *Handler) hasWritePermission(userID int64, calendarID int64) bool {
	return h.calendarService.HasWriteAccess(userID, calendarID)
}

func (h *Handler) hasAdminPermission(userID int64, calendarID int64) bool {
	return h.calendarService.HasAdminAccess(userID, calendarID)
}

// extractACLProperties extracts ACL-related properties from request
func (h *Handler) extractACLProperties(r *http.Request) map[string]string {
	properties := map[string]string{}

	// Check for ACL extension in request headers
	if aclHeader := r.Header.Get("Acl"); aclHeader != "" {
		properties[aclHeader] = aclHeader
	}

	return properties
}

// writeCalendarResponse writes CalDAV XML responses
func (h *Handler) writeCalendarResponse(w http.ResponseWriter, calendar *domain.Calendar) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)

	encoder := xml.NewEncoder(w)

	response := &webdav.MultiStatus{
		XMLName: xml.Name{Space: "DAV:", Local: "multistatus"},
		Responses: []webdav.Response{
			{
				Href: fmt.Sprintf("/caldav/calendars/%d/%s", calendar.UserID, calendar.Name),
				PropStats: []webdav.PropStat{
					{
						Prop:   webdav.PropValue{DisplayName: &calendar.Name},
						Status: "HTTP/1.1 200 OK",
					},
				},
			},
		},
	}

	if err := encoder.Encode(response); err != nil {
		h.logger.Error("failed to encode calendar response", zap.Error(err))
		return
	}
}

// MkCalendarRequest represents a MKCALENDAR request body (RFC 4791)
type MkCalendarRequest struct {
	XMLName xml.Name   `xml:"urn:ietf:params:xml:ns:caldav mkcalendar"`
	Set     *MkCalSet  `xml:"DAV: set"`
	Prop    *MkCalProp `xml:"urn:ietf:params:xml:ns:caldav prop"` // Alternative structure
}

type MkCalSet struct {
	Prop MkCalProp `xml:"DAV: prop"`
}

type MkCalProp struct {
	DisplayName               string `xml:"DAV: displayname"`
	CalendarDescription       string `xml:"urn:ietf:params:xml:ns:caldav calendar-description"`
	CalendarTimezone          string `xml:"urn:ietf:params:xml:ns:caldav calendar-timezone"`
	CalendarColor             string `xml:"http://apple.com/ns/ical/ calendar-color"`
	SupportedCalendarCompSet  *SupportedCalendarComponentSet
}

type SupportedCalendarComponentSet struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:caldav supported-calendar-component-set"`
	Comps   []Comp   `xml:"comp"`
}

type Comp struct {
	Name string `xml:"name,attr"`
}

// CalendarQueryRequest represents a calendar-query REPORT request
type CalendarQueryRequest struct {
	XMLName    xml.Name          `xml:"urn:ietf:params:xml:ns:caldav calendar-query"`
	Prop       *CalendarQueryProp `xml:"DAV: prop"`
	Filter     *CalendarFilter   `xml:"urn:ietf:params:xml:ns:caldav filter"`
}

type CalendarQueryProp struct {
	GetETag      *struct{} `xml:"DAV: getetag"`
	CalendarData *struct{} `xml:"urn:ietf:params:xml:ns:caldav calendar-data"`
}

type CalendarFilter struct {
	CompFilter *CompFilter `xml:"urn:ietf:params:xml:ns:caldav comp-filter"`
}

type CompFilter struct {
	Name       string     `xml:"name,attr"`
	TimeRange  *TimeRange `xml:"urn:ietf:params:xml:ns:caldav time-range"`
	PropFilter []PropFilter `xml:"urn:ietf:params:xml:ns:caldav prop-filter"`
	CompFilter *CompFilter `xml:"urn:ietf:params:xml:ns:caldav comp-filter"`
}

type TimeRange struct {
	Start string `xml:"start,attr"`
	End   string `xml:"end,attr"`
}

type PropFilter struct {
	Name      string     `xml:"name,attr"`
	TextMatch *TextMatch `xml:"urn:ietf:params:xml:ns:caldav text-match"`
}

type TextMatch struct {
	Collation  string `xml:"collation,attr"`
	NegateCondition string `xml:"negate-condition,attr"`
	Value      string `xml:",chardata"`
}

// CalendarMultigetRequest represents a calendar-multiget REPORT request
type CalendarMultigetRequest struct {
	XMLName xml.Name              `xml:"urn:ietf:params:xml:ns:caldav calendar-multiget"`
	Prop    *CalendarQueryProp    `xml:"DAV: prop"`
	Hrefs   []string              `xml:"DAV: href"`
}

// FreeBusyQueryRequest represents a free-busy-query REPORT request
type FreeBusyQueryRequest struct {
	XMLName   xml.Name   `xml:"urn:ietf:params:xml:ns:caldav free-busy-query"`
	TimeRange *TimeRange `xml:"urn:ietf:params:xml:ns:caldav time-range"`
}

// ServeHTTP implements http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("CalDAV request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	switch r.Method {
	case "MKCALENDAR":
		h.handleMkCalendar(w, r)
	case "REPORT":
		h.handleReport(w, r)
	default:
		// Delegate to base WebDAV handler
		h.webdavHandler.ServeHTTP(w, r)
	}
}

// handleMkCalendar handles MKCALENDAR requests (RFC 4791)
func (h *Handler) handleMkCalendar(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("MKCALENDAR request", zap.String("path", r.URL.Path))

	// Get authenticated user from context
	userID, ok := webdav.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract calendar name from URL path
	// Path format: /caldav/calendars/{userID}/{calendarName}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid calendar path", http.StatusBadRequest)
		return
	}

	targetUserIDStr := pathParts[2]
	calendarName := pathParts[3]

	// Parse target user ID
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Permission check: user can only create calendars for themselves unless admin
	if targetUserID != userID {
		// Check if user is admin
		user, err := h.userService.GetByID(userID)
		if err != nil || user == nil || user.Role != "admin" {
			h.logger.Warn("MKCALENDAR permission denied",
				zap.Int64("user_id", userID),
				zap.Int64("target_user_id", targetUserID),
			)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	// Parse request body for calendar properties
	displayName := calendarName
	color := ""
	description := ""
	timezone := "UTC"

	if r.ContentLength > 0 {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error("failed to read MKCALENDAR body", zap.Error(err))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(body) > 0 {
			var mkCalRequest MkCalendarRequest
			if err := xml.Unmarshal(body, &mkCalRequest); err != nil {
				h.logger.Warn("failed to parse MKCALENDAR XML, using defaults",
					zap.Error(err),
					zap.String("body", string(body)),
				)
			} else {
				// Extract properties from request
				var prop *MkCalProp
				if mkCalRequest.Set != nil {
					prop = &mkCalRequest.Set.Prop
				} else if mkCalRequest.Prop != nil {
					prop = mkCalRequest.Prop
				}

				if prop != nil {
					if prop.DisplayName != "" {
						displayName = prop.DisplayName
					}
					if prop.CalendarDescription != "" {
						description = prop.CalendarDescription
					}
					if prop.CalendarColor != "" {
						color = prop.CalendarColor
					}
					if prop.CalendarTimezone != "" {
						timezone = prop.CalendarTimezone
					}
				}

				h.logger.Debug("parsed MKCALENDAR properties",
					zap.String("displayName", displayName),
					zap.String("description", description),
					zap.String("color", color),
					zap.String("timezone", timezone),
				)
			}
		}
	}

	// Create calendar
	calendar, err := h.calendarService.CreateCalendar(targetUserID, calendarName, displayName, color, description, timezone)
	if err != nil {
		h.logger.Error("failed to create calendar", zap.Error(err))
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("calendar created",
		zap.Int64("id", calendar.ID),
		zap.String("name", calendar.Name),
		zap.Int64("user_id", targetUserID),
	)

	// Return 201 Created
	w.Header().Set("Location", r.URL.Path)
	w.WriteHeader(http.StatusCreated)
}

// handleReport handles REPORT requests (RFC 3253, extended by CalDAV)
func (h *Handler) handleReport(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := webdav.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
	case "calendar-query":
		h.handleCalendarQuery(w, r, body, userID)
	case "calendar-multiget":
		h.handleCalendarMultiget(w, r, body, userID)
	case "free-busy-query":
		h.handleFreeBusyQuery(w, r, body, userID)
	default:
		h.logger.Warn("unknown REPORT type", zap.String("type", report.XMLName.Local))
		http.Error(w, "Unsupported report type", http.StatusBadRequest)
	}
}

// handleCalendarQuery handles calendar-query REPORT for both VEVENT and VTODO
func (h *Handler) handleCalendarQuery(w http.ResponseWriter, r *http.Request, body []byte, userID int64) {
	h.logger.Info("calendar-query REPORT", zap.String("path", r.URL.Path))

	// Use the combined handler that supports both VEVENT and VTODO
	h.handleCalendarQueryWithTasks(w, r, body, userID)
}

// handleCalendarMultiget handles calendar-multiget REPORT
func (h *Handler) handleCalendarMultiget(w http.ResponseWriter, r *http.Request, body []byte, userID int64) {
	h.logger.Info("calendar-multiget REPORT", zap.String("path", r.URL.Path))

	// Parse the multiget request
	var multigetRequest CalendarMultigetRequest
	if err := xml.Unmarshal(body, &multigetRequest); err != nil {
		h.logger.Warn("failed to parse calendar-multiget", zap.Error(err))
	}

	// Extract calendar from path
	calendar, err := h.getCalendarFromPath(r.URL.Path, userID)
	if err != nil {
		h.logger.Error("failed to get calendar", zap.Error(err))
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	// Check read permission
	if !h.hasReadPermission(userID, calendar.ID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var events []*domain.Event
	if len(multigetRequest.Hrefs) > 0 {
		// Fetch specific events by UID from href
		for _, href := range multigetRequest.Hrefs {
			uid := extractUIDFromHref(href)
			if uid != "" {
				allEvents, err := h.eventService.GetCalendarEvents(calendar.ID)
				if err != nil {
					continue
				}
				for _, e := range allEvents {
					if e.UID == uid {
						events = append(events, e)
						break
					}
				}
			}
		}
	} else {
		// Return all events
		events, err = h.eventService.GetCalendarEvents(calendar.ID)
		if err != nil {
			h.logger.Error("failed to get events", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Build multistatus response
	h.writeEventsResponse(w, r.URL.Path, events)
}

// handleFreeBusyQuery handles free-busy-query REPORT
func (h *Handler) handleFreeBusyQuery(w http.ResponseWriter, r *http.Request, body []byte, userID int64) {
	h.logger.Info("free-busy-query REPORT", zap.String("path", r.URL.Path))

	// Parse the free-busy-query request
	var fbQuery FreeBusyQueryRequest
	if err := xml.Unmarshal(body, &fbQuery); err != nil {
		h.logger.Warn("failed to parse free-busy-query", zap.Error(err))
	}

	// Extract calendar from path
	calendar, err := h.getCalendarFromPath(r.URL.Path, userID)
	if err != nil {
		h.logger.Error("failed to get calendar", zap.Error(err))
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	// Check read permission
	if !h.hasReadPermission(userID, calendar.ID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse time range from request
	var start, end time.Time
	if fbQuery.TimeRange != nil {
		start, end = parseTimeRange(fbQuery.TimeRange)
	} else {
		// Default to next 30 days
		start = time.Now()
		end = start.Add(30 * 24 * time.Hour)
	}

	// Get events in range
	events, err := h.eventService.GetEventsInRange(calendar.ID, start, end)
	if err != nil {
		h.logger.Error("failed to get events for free-busy", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build VFREEBUSY response
	freeBusyData := buildFreeBusyResponse(start, end, events)

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(freeBusyData))
}

// Helper functions

func (h *Handler) getCalendarFromPath(urlPath string, userID int64) (*domain.Calendar, error) {
	pathParts := strings.Split(strings.Trim(urlPath, "/"), "/")
	if len(pathParts) < 4 {
		return nil, fmt.Errorf("invalid calendar path")
	}

	targetUserIDStr := pathParts[2]
	calendarName := pathParts[3]

	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get calendars for user
	calendars, err := h.calendarService.GetUserCalendars(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendars: %w", err)
	}

	for _, cal := range calendars {
		if cal.Name == calendarName {
			return cal, nil
		}
	}

	return nil, fmt.Errorf("calendar not found: %s", calendarName)
}

func (h *Handler) writeEventsResponse(w http.ResponseWriter, basePath string, events []*domain.Event) {
	multistatus := &webdav.MultiStatus{
		Responses: []webdav.Response{},
	}

	for _, event := range events {
		eventPath := basePath
		if !strings.HasSuffix(eventPath, "/") {
			eventPath += "/"
		}
		eventPath += event.UID + ".ics"

		propValue := webdav.PropValue{
			GetETag:      &event.ETag,
			CalendarData: &event.ICalData,
		}

		response := webdav.Response{
			Href: eventPath,
			PropStats: []webdav.PropStat{
				{
					Prop:   propValue,
					Status: "HTTP/1.1 200 OK",
				},
			},
		}
		multistatus.Responses = append(multistatus.Responses, response)
	}

	xmlData, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}

func parseTimeRange(tr *TimeRange) (start, end time.Time) {
	// Parse iCal datetime format: 20060102T150405Z
	const iCalFormat = "20060102T150405Z"
	const iCalDateFormat = "20060102"

	start = time.Now()
	end = start.Add(30 * 24 * time.Hour)

	if tr.Start != "" {
		if parsed, err := time.Parse(iCalFormat, tr.Start); err == nil {
			start = parsed
		} else if parsed, err := time.Parse(iCalDateFormat, tr.Start); err == nil {
			start = parsed
		}
	}

	if tr.End != "" {
		if parsed, err := time.Parse(iCalFormat, tr.End); err == nil {
			end = parsed
		} else if parsed, err := time.Parse(iCalDateFormat, tr.End); err == nil {
			end = parsed
		}
	}

	return start, end
}

func extractUIDFromHref(href string) string {
	// Extract UID from href like /caldav/calendars/1/work/event-uid.ics
	if strings.HasSuffix(href, ".ics") {
		base := href[strings.LastIndex(href, "/")+1:]
		return strings.TrimSuffix(base, ".ics")
	}
	return ""
}

func buildFreeBusyResponse(start, end time.Time, events []*domain.Event) string {
	var fb strings.Builder
	fb.WriteString("BEGIN:VCALENDAR\r\n")
	fb.WriteString("VERSION:2.0\r\n")
	fb.WriteString("PRODID:-//gomailserver//CalDAV Server//EN\r\n")
	fb.WriteString("METHOD:REPLY\r\n")
	fb.WriteString("BEGIN:VFREEBUSY\r\n")
	fb.WriteString(fmt.Sprintf("DTSTART:%s\r\n", start.UTC().Format("20060102T150405Z")))
	fb.WriteString(fmt.Sprintf("DTEND:%s\r\n", end.UTC().Format("20060102T150405Z")))
	fb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", time.Now().UTC().Format("20060102T150405Z")))

	// Add FREEBUSY entries for each event
	for _, event := range events {
		if event.Status == "CANCELLED" {
			continue
		}
		fbType := "BUSY"
		if event.Status == "TENTATIVE" {
			fbType = "BUSY-TENTATIVE"
		}
		fb.WriteString(fmt.Sprintf("FREEBUSY;FBTYPE=%s:%s/%s\r\n",
			fbType,
			event.StartTime.UTC().Format("20060102T150405Z"),
			event.EndTime.UTC().Format("20060102T150405Z"),
		))
	}

	fb.WriteString("END:VFREEBUSY\r\n")
	fb.WriteString("END:VCALENDAR\r\n")

	return fb.String()
}

// handleCalendarQueryWithTasks handles calendar-query for both VEVENT and VTODO
func (h *Handler) handleCalendarQueryWithTasks(w http.ResponseWriter, r *http.Request, body []byte, userID int64) {
	// Parse the calendar-query request
	var queryRequest CalendarQueryRequest
	if err := xml.Unmarshal(body, &queryRequest); err != nil {
		h.logger.Warn("failed to parse calendar-query", zap.Error(err))
	}

	// Extract calendar from path
	calendar, err := h.getCalendarFromPath(r.URL.Path, userID)
	if err != nil {
		h.logger.Error("failed to get calendar", zap.Error(err))
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	// Check read permission
	if !h.hasReadPermission(userID, calendar.ID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Determine which component type is being queried
	compType := "VEVENT" // Default to VEVENT
	if queryRequest.Filter != nil && queryRequest.Filter.CompFilter != nil {
		if queryRequest.Filter.CompFilter.CompFilter != nil {
			compType = queryRequest.Filter.CompFilter.CompFilter.Name
		}
	}

	multistatus := &webdav.MultiStatus{
		Responses: []webdav.Response{},
	}

	switch compType {
	case "VTODO":
		// Handle VTODO query
		if h.taskService == nil {
			h.logger.Warn("task service not configured, returning empty response")
			break
		}

		tasks, err := h.taskService.GetCalendarTasks(calendar.ID)
		if err != nil {
			h.logger.Error("failed to get tasks", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		for _, task := range tasks {
			taskPath := r.URL.Path
			if !strings.HasSuffix(taskPath, "/") {
				taskPath += "/"
			}
			taskPath += task.UID + ".ics"

			propValue := webdav.PropValue{
				GetETag:      &task.ETag,
				CalendarData: &task.ICalData,
			}

			response := webdav.Response{
				Href: taskPath,
				PropStats: []webdav.PropStat{
					{
						Prop:   propValue,
						Status: "HTTP/1.1 200 OK",
					},
				},
			}
			multistatus.Responses = append(multistatus.Responses, response)
		}

	default:
		// Handle VEVENT query (existing implementation)
		var events []*domain.Event
		if queryRequest.Filter != nil && queryRequest.Filter.CompFilter != nil {
			var timeRange *TimeRange
			if queryRequest.Filter.CompFilter.CompFilter != nil {
				timeRange = queryRequest.Filter.CompFilter.CompFilter.TimeRange
			} else {
				timeRange = queryRequest.Filter.CompFilter.TimeRange
			}

			if timeRange != nil {
				start, end := parseTimeRange(timeRange)
				events, err = h.eventService.GetEventsInRange(calendar.ID, start, end)
			} else {
				events, err = h.eventService.GetCalendarEvents(calendar.ID)
			}
		} else {
			events, err = h.eventService.GetCalendarEvents(calendar.ID)
		}

		if err != nil {
			h.logger.Error("failed to get events", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		for _, event := range events {
			eventPath := r.URL.Path
			if !strings.HasSuffix(eventPath, "/") {
				eventPath += "/"
			}
			eventPath += event.UID + ".ics"

			propValue := webdav.PropValue{
				GetETag:      &event.ETag,
				CalendarData: &event.ICalData,
			}

			response := webdav.Response{
				Href: eventPath,
				PropStats: []webdav.PropStat{
					{
						Prop:   propValue,
						Status: "HTTP/1.1 200 OK",
					},
				},
			}
			multistatus.Responses = append(multistatus.Responses, response)
		}
	}

	xmlData, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}

// writeTasksResponse writes tasks as a CalDAV multistatus response
func (h *Handler) writeTasksResponse(w http.ResponseWriter, basePath string, tasks []*domain.Task) {
	multistatus := &webdav.MultiStatus{
		Responses: []webdav.Response{},
	}

	for _, task := range tasks {
		taskPath := basePath
		if !strings.HasSuffix(taskPath, "/") {
			taskPath += "/"
		}
		taskPath += task.UID + ".ics"

		propValue := webdav.PropValue{
			GetETag:      &task.ETag,
			CalendarData: &task.ICalData,
		}

		response := webdav.Response{
			Href: taskPath,
			PropStats: []webdav.PropStat{
				{
					Prop:   propValue,
					Status: "HTTP/1.1 200 OK",
				},
			},
		}
		multistatus.Responses = append(multistatus.Responses, response)
	}

	xmlData, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}
