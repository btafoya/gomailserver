package handlers

import (
	"net/http"
	"strconv"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// StatsHandler handles statistics and monitoring endpoints
type StatsHandler struct {
	domainService *service.DomainService
	userService   *service.UserService
	queueService  *service.QueueService
	logger        *zap.Logger
}

// NewStatsHandler creates a new statistics handler
func NewStatsHandler(
	domainService *service.DomainService,
	userService *service.UserService,
	queueService *service.QueueService,
	logger *zap.Logger,
) *StatsHandler {
	return &StatsHandler{
		domainService: domainService,
		userService:   userService,
		queueService:  queueService,
		logger:        logger,
	}
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalDomains       int64              `json:"total_domains"`
	ActiveDomains      int64              `json:"active_domains"`
	TotalUsers         int64              `json:"total_users"`
	ActiveUsers        int64              `json:"active_users"`
	QueuedMessages     int64              `json:"queued_messages"`
	FailedMessages     int64              `json:"failed_messages"`
	TotalStorageUsed   int64              `json:"total_storage_used"`
	TotalStorageQuota  int64              `json:"total_storage_quota"`
	MessagesToday      int64              `json:"messages_today"`
	MessagesThisWeek   int64              `json:"messages_this_week"`
	MessagesThisMonth  int64              `json:"messages_this_month"`
	RecentActivity     []ActivityItem     `json:"recent_activity"`
	TopDomainsByUsers  []DomainStat       `json:"top_domains_by_users"`
	TopDomainsByUsage  []DomainStat       `json:"top_domains_by_usage"`
	SystemHealth       SystemHealthStatus `json:"system_health"`
}

// ActivityItem represents a recent activity entry
type ActivityItem struct {
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
	Description string `json:"description"`
	UserEmail   string `json:"user_email,omitempty"`
	DomainName  string `json:"domain_name,omitempty"`
}

// DomainStat represents domain statistics
type DomainStat struct {
	DomainID     int64  `json:"domain_id"`
	DomainName   string `json:"domain_name"`
	UserCount    int64  `json:"user_count"`
	StorageUsed  int64  `json:"storage_used"`
	MessageCount int64  `json:"message_count"`
}

// SystemHealthStatus represents system health indicators
type SystemHealthStatus struct {
	Status           string  `json:"status"` // "healthy", "degraded", "critical"
	DatabaseStatus   string  `json:"database_status"`
	SMTPStatus       string  `json:"smtp_status"`
	IMAPStatus       string  `json:"imap_status"`
	QueueDepth       int64   `json:"queue_depth"`
	QueueHealthy     bool    `json:"queue_healthy"`
	DiskUsage        float64 `json:"disk_usage_percent"`
	MemoryUsage      float64 `json:"memory_usage_percent"`
	CPUUsage         float64 `json:"cpu_usage_percent"`
	UptimeSeconds    int64   `json:"uptime_seconds"`
}

// DomainStats represents statistics for a specific domain
type DomainStats struct {
	DomainID          int64  `json:"domain_id"`
	DomainName        string `json:"domain_name"`
	Status            string `json:"status"`
	TotalUsers        int64  `json:"total_users"`
	ActiveUsers       int64  `json:"active_users"`
	DisabledUsers     int64  `json:"disabled_users"`
	TotalStorageUsed  int64  `json:"total_storage_used"`
	TotalStorageQuota int64  `json:"total_storage_quota"`
	TotalAliases      int64  `json:"total_aliases"`
	MessagesToday     int64  `json:"messages_today"`
	MessagesThisWeek  int64  `json:"messages_this_week"`
	MessagesThisMonth int64  `json:"messages_this_month"`
	TopUsers          []UserStat `json:"top_users"`
}

// UserStat represents user statistics
type UserStat struct {
	UserID       int64  `json:"user_id"`
	Email        string `json:"email"`
	StorageUsed  int64  `json:"storage_used"`
	MessageCount int64  `json:"message_count"`
	QuotaPercent float64 `json:"quota_percent"`
}

// UserStats represents statistics for a specific user
type UserStats struct {
	UserID            int64  `json:"user_id"`
	Email             string `json:"email"`
	DomainID          int64  `json:"domain_id"`
	DomainName        string `json:"domain_name"`
	Status            string `json:"status"`
	StorageUsed       int64  `json:"storage_used"`
	StorageQuota      int64  `json:"storage_quota"`
	QuotaPercent      float64 `json:"quota_percent"`
	TotalMessages     int64  `json:"total_messages"`
	MessagesToday     int64  `json:"messages_today"`
	MessagesThisWeek  int64  `json:"messages_this_week"`
	MessagesThisMonth int64  `json:"messages_this_month"`
	LastLogin         string `json:"last_login,omitempty"`
	CreatedAt         string `json:"created_at"`
}

// Dashboard retrieves dashboard statistics
func (h *StatsHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual statistics gathering from database
	// For now, return placeholder data

	stats := DashboardStats{
		TotalDomains:      0,
		ActiveDomains:     0,
		TotalUsers:        0,
		ActiveUsers:       0,
		QueuedMessages:    0,
		FailedMessages:    0,
		TotalStorageUsed:  0,
		TotalStorageQuota: 0,
		MessagesToday:     0,
		MessagesThisWeek:  0,
		MessagesThisMonth: 0,
		RecentActivity:    []ActivityItem{},
		TopDomainsByUsers: []DomainStat{},
		TopDomainsByUsage: []DomainStat{},
		SystemHealth: SystemHealthStatus{
			Status:         "healthy",
			DatabaseStatus: "connected",
			SMTPStatus:     "running",
			IMAPStatus:     "running",
			QueueDepth:     0,
			QueueHealthy:   true,
			DiskUsage:      0.0,
			MemoryUsage:    0.0,
			CPUUsage:       0.0,
			UptimeSeconds:  0,
		},
	}

	h.logger.Info("Dashboard statistics requested")

	middleware.RespondSuccess(w, stats, "Dashboard statistics retrieved successfully")
}

// Domain retrieves statistics for a specific domain
func (h *StatsHandler) Domain(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	// Get domain
	domain, err := h.domainService.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get domain", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Domain not found")
		return
	}

	// TODO: Implement actual statistics gathering
	stats := DomainStats{
		DomainID:          domain.ID,
		DomainName:        domain.Name,
		Status:            domain.Status,
		TotalUsers:        0,
		ActiveUsers:       0,
		DisabledUsers:     0,
		TotalStorageUsed:  0,
		TotalStorageQuota: 0,
		TotalAliases:      0,
		MessagesToday:     0,
		MessagesThisWeek:  0,
		MessagesThisMonth: 0,
		TopUsers:          []UserStat{},
	}

	middleware.RespondSuccess(w, stats, "Domain statistics retrieved successfully")
}

// User retrieves statistics for a specific user
func (h *StatsHandler) User(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get user
	user, err := h.userService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Calculate quota percentage
	quotaPercent := 0.0
	if user.Quota > 0 {
		quotaPercent = (float64(user.UsedQuota) / float64(user.Quota)) * 100
	}

	// TODO: Implement actual statistics gathering
	stats := UserStats{
		UserID:            user.ID,
		Email:             user.Email,
		DomainID:          user.DomainID,
		DomainName:        "", // TODO: Get domain name
		Status:            user.Status,
		StorageUsed:       user.UsedQuota,
		StorageQuota:      user.Quota,
		QuotaPercent:      quotaPercent,
		TotalMessages:     0,
		MessagesToday:     0,
		MessagesThisWeek:  0,
		MessagesThisMonth: 0,
		CreatedAt:         user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if user.LastLogin != nil {
		stats.LastLogin = user.LastLogin.Format("2006-01-02T15:04:05Z07:00")
	}

	middleware.RespondSuccess(w, stats, "User statistics retrieved successfully")
}
