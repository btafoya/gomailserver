package caldav

import (
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
	"github.com/btafoya/gomailserver/internal/webdav"
	"go.uber.org/zap"
)

// Handler handles CalDAV-specific requests
type Handler struct {
	webdavHandler   *webdav.Handler
	logger          *zap.Logger
	calendarService domain.CalendarService
	eventService    domain.EventService
}

// NewHandler creates a new CalDAV handler
func NewHandler(logger *zap.Logger, calendarService domain.CalendarService, eventService domain.EventService) *Handler {
	return &Handler{
		webdavHandler:   webdav.NewHandler(logger, "/caldav"),
		logger:          logger,
		calendarService: calendarService,
		eventService:    eventService,
	}
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

	// Extract calendar name from URL path
	// Path format: /caldav/calendars/{userID}/{calendarName}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid calendar path", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[2]
	calendarName := pathParts[3]

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// TODO: Validate user has permission to create calendar
	// For now, we'll allow any authenticated user to create calendars

	// Parse request body for calendar properties (optional)
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

		// Parse XML for calendar properties
		// This is a simplified implementation - full CalDAV would parse <C:mkcalendar> XML
		_ = body // TODO: Parse display-name, color, description from XML
	}

	// Create calendar
	calendar, err := h.calendarService.CreateCalendar(userID, calendarName, displayName, color, description, timezone)
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
	)

	// Return 201 Created
	w.Header().Set("Location", r.URL.Path)
	w.WriteHeader(http.StatusCreated)
}

// handleReport handles REPORT requests (RFC 3253, extended by CalDAV)
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
	case "calendar-query":
		h.handleCalendarQuery(w, r, body)
	case "calendar-multiget":
		h.handleCalendarMultiget(w, r, body)
	case "free-busy-query":
		h.handleFreeBusyQuery(w, r, body)
	default:
		h.logger.Warn("unknown REPORT type", zap.String("type", report.XMLName.Local))
		http.Error(w, "Unsupported report type", http.StatusBadRequest)
	}
}

// handleCalendarQuery handles calendar-query REPORT
func (h *Handler) handleCalendarQuery(w http.ResponseWriter, r *http.Request, body []byte) {
	h.logger.Info("calendar-query REPORT", zap.String("path", r.URL.Path))

	// Extract calendar ID from path
	// Path format: /caldav/calendars/{userID}/{calendarName}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid calendar path", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[2]
	calendarName := pathParts[3]

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get calendar by user and name
	calendars, err := h.calendarService.GetUserCalendars(userID)
	if err != nil {
		h.logger.Error("failed to get calendars", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var calendarID int64
	for _, cal := range calendars {
		if cal.Name == calendarName {
			calendarID = cal.ID
			break
		}
	}

	if calendarID == 0 {
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	// Parse request body for time-range filter
	// This is a simplified implementation - full CalDAV would parse complex filters
	// For now, we'll return all events in the calendar
	events, err := h.eventService.GetCalendarEvents(calendarID)
	if err != nil {
		h.logger.Error("failed to get events", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build multistatus response
	multistatus := &webdav.MultiStatus{
		Responses: []webdav.Response{},
	}

	for _, event := range events {
		eventPath := r.URL.Path + "/" + event.UID + ".ics"

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

// handleCalendarMultiget handles calendar-multiget REPORT
func (h *Handler) handleCalendarMultiget(w http.ResponseWriter, r *http.Request, body []byte) {
	h.logger.Info("calendar-multiget REPORT", zap.String("path", r.URL.Path))

	// Extract calendar ID from path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid calendar path", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[2]
	calendarName := pathParts[3]

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get calendar by user and name
	calendars, err := h.calendarService.GetUserCalendars(userID)
	if err != nil {
		h.logger.Error("failed to get calendars", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var calendarID int64
	for _, cal := range calendars {
		if cal.Name == calendarName {
			calendarID = cal.ID
			break
		}
	}

	if calendarID == 0 {
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	// Parse request body for hrefs
	// This is a simplified implementation - would parse <C:calendar-multiget> XML with <D:href> elements
	// For now, we'll return all events
	events, err := h.eventService.GetCalendarEvents(calendarID)
	if err != nil {
		h.logger.Error("failed to get events", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build multistatus response
	multistatus := &webdav.MultiStatus{
		Responses: []webdav.Response{},
	}

	for _, event := range events {
		eventPath := r.URL.Path + "/" + event.UID + ".ics"

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

// handleFreeBusyQuery handles free-busy-query REPORT
func (h *Handler) handleFreeBusyQuery(w http.ResponseWriter, r *http.Request, body []byte) {
	h.logger.Info("free-busy-query REPORT", zap.String("path", r.URL.Path))

	// Extract calendar ID from path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid calendar path", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[2]
	calendarName := pathParts[3]

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get calendar by user and name
	calendars, err := h.calendarService.GetUserCalendars(userID)
	if err != nil {
		h.logger.Error("failed to get calendars", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var calendarID int64
	for _, cal := range calendars {
		if cal.Name == calendarName {
			calendarID = cal.ID
			break
		}
	}

	if calendarID == 0 {
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	// Parse request body for time range
	// This is a simplified implementation - would parse <C:free-busy-query> XML with time-range
	// For now, we'll use a default 30-day range
	// TODO: Parse actual time range from request

	// Simple free/busy response (placeholder)
	freeBusyData := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//gomailserver//CalDAV Server//EN\nBEGIN:VFREEBUSY\nEND:VFREEBUSY\nEND:VCALENDAR"

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(freeBusyData))
}
