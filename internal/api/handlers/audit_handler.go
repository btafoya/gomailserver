package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/service"
	"go.uber.org/zap"
)

// AuditHandler handles audit log viewing endpoints
type AuditHandler struct {
	service *service.AuditService
	logger  *zap.Logger
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(service *service.AuditService, logger *zap.Logger) *AuditHandler {
	return &AuditHandler{
		service: service,
		logger:  logger,
	}
}

// AuditLogResponse represents an audit log entry in API responses
type AuditLogResponse struct {
	ID           int64  `json:"id"`
	Timestamp    string `json:"timestamp"`
	UserID       *int64 `json:"user_id,omitempty"`
	Username     string `json:"username,omitempty"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id,omitempty"`
	Details      string `json:"details,omitempty"`
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	Severity     string `json:"severity"`
	Success      bool   `json:"success"`
}

// ListLogs lists audit logs with optional filtering
func (h *AuditHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse filters
	filter := service.AuditLogFilter{
		Limit:  20, // Default limit
		Offset: 0,
	}

	// Parse user_id filter
	if userIDStr := query.Get("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil {
			filter.UserID = &userID
		}
	}

	// Parse action filter
	if action := query.Get("action"); action != "" {
		filter.Action = action
	}

	// Parse resource_type filter
	if resourceType := query.Get("resource_type"); resourceType != "" {
		filter.ResourceType = resourceType
	}

	// Parse severity filter
	if severity := query.Get("severity"); severity != "" {
		filter.Severity = severity
	}

	// Parse start_time filter
	if startTimeStr := query.Get("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filter.StartTime = startTime
		}
	}

	// Parse end_time filter
	if endTimeStr := query.Get("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filter.EndTime = endTime
		}
	}

	// Parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	logs, err := h.service.GetLogs(r.Context(), filter)
	if err != nil {
		h.logger.Error("failed to list audit logs",
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve audit logs")
		return
	}

	responses := make([]AuditLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = h.toResponse(log)
	}

	middleware.RespondJSON(w, http.StatusOK, responses)
}

// GetStats returns audit log statistics
func (h *AuditHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse time range
	var startTime, endTime time.Time
	if startTimeStr := query.Get("start_time"); startTimeStr != "" {
		if st, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			startTime = st
		}
	}
	if endTimeStr := query.Get("end_time"); endTimeStr != "" {
		if et, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			endTime = et
		}
	}

	// Default to last 24 hours if not specified
	if startTime.IsZero() {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	if endTime.IsZero() {
		endTime = time.Now()
	}

	filter := service.AuditLogFilter{
		StartTime: startTime,
		EndTime:   endTime,
	}

	logs, err := h.service.GetLogs(r.Context(), filter)
	if err != nil {
		h.logger.Error("failed to get audit log stats",
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve statistics")
		return
	}

	// Calculate statistics
	stats := map[string]interface{}{
		"total":     len(logs),
		"period":    map[string]string{"start": startTime.Format(time.RFC3339), "end": endTime.Format(time.RFC3339)},
		"by_action": make(map[string]int),
		"by_severity": make(map[string]int),
		"by_resource": make(map[string]int),
		"success_rate": 0.0,
	}

	successCount := 0
	byAction := stats["by_action"].(map[string]int)
	bySeverity := stats["by_severity"].(map[string]int)
	byResource := stats["by_resource"].(map[string]int)

	for _, log := range logs {
		byAction[log.Action]++
		bySeverity[log.Severity]++
		byResource[log.ResourceType]++
		if log.Success {
			successCount++
		}
	}

	if len(logs) > 0 {
		stats["success_rate"] = float64(successCount) / float64(len(logs)) * 100
	}

	middleware.RespondJSON(w, http.StatusOK, stats)
}

// toResponse converts a domain audit log to API response format
func (h *AuditHandler) toResponse(log *domain.AuditLog) AuditLogResponse {
	return AuditLogResponse{
		ID:           log.ID,
		Timestamp:    log.Timestamp.Format(time.RFC3339),
		UserID:       log.UserID,
		Username:     log.Username,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		Details:      log.Details,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		Severity:     log.Severity,
		Success:      log.Success,
	}
}
