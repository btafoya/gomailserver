package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	ics "github.com/arran4/golang-ical"
	calendarService "github.com/btafoya/gomailserver/internal/calendar/service"
	"github.com/btafoya/gomailserver/internal/api/middleware"
	"go.uber.org/zap"
)

// WebmailCalendarHandler handles webmail calendar operations
type WebmailCalendarHandler struct {
	calendarService *calendarService.CalendarService
	eventService    *calendarService.EventService
	logger          *zap.Logger
}

// NewWebmailCalendarHandler creates a new webmail calendar handler
func NewWebmailCalendarHandler(
	calendarSvc *calendarService.CalendarService,
	eventSvc *calendarService.EventService,
	logger *zap.Logger,
) *WebmailCalendarHandler {
	return &WebmailCalendarHandler{
		calendarService: calendarSvc,
		eventService:    eventSvc,
		logger:          logger,
	}
}

// CalendarResponse represents a calendar for webmail
type CalendarResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
	Timezone    string `json:"timezone"`
}

// EventResponse represents an event for webmail
type EventResponse struct {
	ID          int64     `json:"id"`
	UID         string    `json:"uid"`
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	AllDay      bool      `json:"all_day"`
	Timezone    string    `json:"timezone,omitempty"`
}

// CreateEventRequest represents an event creation request
type CreateEventRequest struct {
	CalendarID  int64     `json:"calendar_id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	AllDay      bool      `json:"all_day"`
	Timezone    string    `json:"timezone,omitempty"`
	Attendees   []string  `json:"attendees,omitempty"`
}

// ProcessInvitationRequest represents a calendar invitation response
type ProcessInvitationRequest struct {
	CalendarID  int64  `json:"calendar_id"`
	ICalData    string `json:"ical_data"`
	ResponseStr string `json:"response"` // "ACCEPTED", "TENTATIVE", "DECLINED"
}

// ListCalendars lists user's calendars
// GET /api/v1/webmail/calendar/calendars
func (h *WebmailCalendarHandler) ListCalendars(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	calendars, err := h.calendarService.GetUserCalendars(userID)
	if err != nil {
		h.logger.Error("failed to get calendars", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get calendars")
		return
	}

	// Convert to response format
	var results []CalendarResponse
	for _, cal := range calendars {
		results = append(results, CalendarResponse{
			ID:          cal.ID,
			Name:        cal.Name,
			DisplayName: cal.DisplayName,
			Color:       cal.Color,
			Description: cal.Description,
			Timezone:    cal.Timezone,
		})
	}

	middleware.RespondJSON(w, http.StatusOK, results)
}

// GetUpcomingEvents gets upcoming events across all calendars
// GET /api/v1/webmail/calendar/upcoming?days=7
func (h *WebmailCalendarHandler) GetUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get days parameter (default 7)
	days := 7
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 90 {
			days = d
		}
	}

	// Get user's calendars
	calendars, err := h.calendarService.GetUserCalendars(userID)
	if err != nil {
		h.logger.Error("failed to get calendars", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get calendars")
		return
	}

	// Get events from all calendars
	now := time.Now()
	endTime := now.AddDate(0, 0, days)

	var allEvents []EventResponse
	for _, cal := range calendars {
		events, err := h.eventService.GetEventsInRange(cal.ID, now, endTime)
		if err != nil {
			h.logger.Error("failed to get events", zap.Error(err), zap.Int64("calendar_id", cal.ID))
			continue
		}

		for _, event := range events {
			allEvents = append(allEvents, EventResponse{
				ID:          event.ID,
				UID:         event.UID,
				Summary:     event.Summary,
				Description: event.Description,
				Location:    event.Location,
				StartTime:   event.StartTime,
				EndTime:     event.EndTime,
				AllDay:      event.AllDay,
				Timezone:    event.Timezone,
			})
		}
	}

	middleware.RespondJSON(w, http.StatusOK, allEvents)
}

// CreateEvent creates a new calendar event from webmail
// POST /api/v1/webmail/calendar/events
func (h *WebmailCalendarHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.CalendarID == 0 {
		middleware.RespondError(w, http.StatusBadRequest, "calendar_id is required")
		return
	}
	if req.Summary == "" {
		middleware.RespondError(w, http.StatusBadRequest, "summary is required")
		return
	}
	if req.StartTime.IsZero() {
		middleware.RespondError(w, http.StatusBadRequest, "start_time is required")
		return
	}
	if req.EndTime.IsZero() {
		middleware.RespondError(w, http.StatusBadRequest, "end_time is required")
		return
	}

	// Verify calendar belongs to user
	calendar, err := h.calendarService.GetCalendar(req.CalendarID)
	if err != nil {
		h.logger.Error("failed to get calendar", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get calendar")
		return
	}
	if calendar.UserID != userID {
		middleware.RespondError(w, http.StatusForbidden, "access denied")
		return
	}

	// Set default timezone if not provided
	timezone := req.Timezone
	if timezone == "" {
		timezone = calendar.Timezone
	}

	// Build iCalendar data using RFC 5545 compliant library (github.com/arran4/golang-ical)
	icalData := h.buildICalendarData(req.Summary, req.Description, req.Location, req.StartTime, req.EndTime, req.AllDay, timezone, req.Attendees)

	// Create event
	event, err := h.eventService.CreateEvent(req.CalendarID, icalData)
	if err != nil {
		h.logger.Error("failed to create event", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to create event")
		return
	}

	// Return created event
	middleware.RespondJSON(w, http.StatusCreated, EventResponse{
		ID:          event.ID,
		UID:         event.UID,
		Summary:     event.Summary,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		AllDay:      event.AllDay,
		Timezone:    event.Timezone,
	})
}

// ProcessInvitation processes a calendar invitation (accept/decline/tentative)
// POST /api/v1/webmail/calendar/invitations
func (h *WebmailCalendarHandler) ProcessInvitation(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request
	var req ProcessInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.CalendarID == 0 {
		middleware.RespondError(w, http.StatusBadRequest, "calendar_id is required")
		return
	}
	if req.ICalData == "" {
		middleware.RespondError(w, http.StatusBadRequest, "ical_data is required")
		return
	}
	if req.ResponseStr == "" {
		middleware.RespondError(w, http.StatusBadRequest, "response is required")
		return
	}

	// Validate response type
	validResponses := map[string]bool{
		"ACCEPTED":  true,
		"TENTATIVE": true,
		"DECLINED":  true,
	}
	if !validResponses[req.ResponseStr] {
		middleware.RespondError(w, http.StatusBadRequest, "response must be ACCEPTED, TENTATIVE, or DECLINED")
		return
	}

	// Verify calendar belongs to user
	calendar, err := h.calendarService.GetCalendar(req.CalendarID)
	if err != nil {
		h.logger.Error("failed to get calendar", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get calendar")
		return
	}
	if calendar.UserID != userID {
		middleware.RespondError(w, http.StatusForbidden, "access denied")
		return
	}

	// Update iCal data with user's response
	// TODO: Use proper iCalendar library to parse and modify
	// For now, create event from invitation if accepted
	if req.ResponseStr == "ACCEPTED" {
		// Create event from invitation
		event, err := h.eventService.CreateEvent(req.CalendarID, req.ICalData)
		if err != nil {
			h.logger.Error("failed to create event from invitation", zap.Error(err))
			middleware.RespondError(w, http.StatusInternalServerError, "failed to create event")
			return
		}

		middleware.RespondJSON(w, http.StatusCreated, EventResponse{
			ID:          event.ID,
			UID:         event.UID,
			Summary:     event.Summary,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			AllDay:      event.AllDay,
			Timezone:    event.Timezone,
		})
	} else {
		// For TENTATIVE or DECLINED, just acknowledge
		middleware.RespondJSON(w, http.StatusOK, map[string]string{
			"status":   "invitation processed",
			"response": req.ResponseStr,
		})
	}
}

// buildICalendarData builds basic iCalendar data
// TODO: Replace with proper iCalendar library (e.g., github.com/arran4/golang-ical)
func (h *WebmailCalendarHandler) buildICalendarData(summary, description, location string, start, end time.Time, allDay bool, timezone string, attendees []string) string {
	// Create a new calendar using the golang-ical library
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetProductId("-//gomailserver//NONSGML v1.0//EN")

	// Create a new event
	event := cal.AddEvent(time.Now().Format("20060102T150405") + "@gomailserver")

	// Set summary
	event.SetSummary(summary)

	// Set description if provided
	if description != "" {
		event.SetDescription(description)
	}

	// Set location if provided
	if location != "" {
		event.SetLocation(location)
	}

	// Set start and end times
	if allDay {
		event.SetAllDayStartAt(start)
		event.SetAllDayEndAt(end)
	} else {
		event.SetStartAt(start)
		event.SetEndAt(end)
	}

	// Set timestamps
	now := time.Now().UTC()
	event.SetDtStampTime(now)
	event.SetCreatedTime(now)
	event.SetModifiedAt(now)

	// Add attendees
	for _, attendee := range attendees {
		event.AddAttendee("mailto:" + attendee)
	}

	return cal.Serialize()
}
