package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/service"
	"go.uber.org/zap"
)

// LogHandler handles log retrieval endpoints
type LogHandler struct {
	auditService *service.AuditService
	logger       *zap.Logger
}

// NewLogHandler creates a new log handler
func NewLogHandler(auditService *service.AuditService, logger *zap.Logger) *LogHandler {
	return &LogHandler{
		auditService: auditService,
		logger:       logger,
	}
}

// LogEntry represents a log entry in API responses
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	UserEmail string `json:"user_email,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
	Action    string `json:"action,omitempty"`
	Result    string `json:"result,omitempty"`
	Message   string `json:"message"`
}

// LogsResponse represents a paginated logs response
type LogsResponse struct {
	Logs       []LogEntry `json:"logs"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
	TotalCount int        `json:"total_count"`
}

// List retrieves logs with optional filtering
func (h *LogHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	level := r.URL.Query().Get("level")
	resourceType := r.URL.Query().Get("service")
	action := r.URL.Query().Get("action")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Set defaults
	page := 1
	pageSize := 50

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Build filter
	filter := service.AuditLogFilter{
		Severity:     level,
		ResourceType: resourceType,
		Action:       action,
		Limit:        pageSize,
		Offset:       (page - 1) * pageSize,
	}

	// Parse date filters
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartTime = t
		} else if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			filter.StartTime = t
		}
	}

	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndTime = t.Add(24 * time.Hour) // Include the end date
		} else if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			filter.EndTime = t
		}
	}

	h.logger.Debug("Logs requested",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("level", level),
		zap.String("resource_type", resourceType),
		zap.String("action", action),
	)

	// Retrieve logs from database
	logs, err := h.auditService.GetLogs(r.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to retrieve logs", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve logs")
		return
	}

	// Convert to response format
	entries := make([]LogEntry, 0, len(logs))
	for _, log := range logs {
		entry := LogEntry{
			Timestamp: log.Timestamp.Format(time.RFC3339),
			Level:     log.Severity,
			Service:   log.ResourceType,
			IPAddress: log.IPAddress,
			Action:    log.Action,
			Message:   log.Details,
		}
		if log.Username != "" {
			entry.UserEmail = log.Username
		}
		if log.Success {
			entry.Result = "success"
		} else {
			entry.Result = "failure"
		}
		entries = append(entries, entry)
	}

	// Calculate total (for proper pagination, we'd need a count query)
	totalCount := len(entries)
	if totalCount == pageSize {
		// There might be more, but we don't know exactly how many
		totalCount = page * pageSize
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	response := LogsResponse{
		Logs:       entries,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}

	middleware.RespondSuccess(w, response, "Logs retrieved successfully")
}
