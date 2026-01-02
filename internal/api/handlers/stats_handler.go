package handlers

import (
	"net/http"
	"strconv"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// StatsHandler handles statistics and monitoring endpoints
type StatsHandler struct {
	domainService *service.DomainService
	userService   *service.UserService
	queueService  *service.QueueService
	aliasService  *service.AliasService
	logger        *zap.Logger
}

// NewStatsHandler creates a new statistics handler
func NewStatsHandler(
	domainService *service.DomainService,
	userService *service.UserService,
	queueService *service.QueueService,
	aliasService *service.AliasService,
	logger *zap.Logger,
) *StatsHandler {
	return &StatsHandler{
		domainService: domainService,
		userService:   userService,
		queueService:  queueService,
		aliasService:  aliasService,
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
	ctx := r.Context()

	// Get all domains
	domains, err := h.domainService.List(ctx)
	if err != nil {
		h.logger.Error("Failed to list domains for dashboard", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve dashboard statistics")
		return
	}

	totalDomains := int64(len(domains))
	activeDomains := int64(0)
	for _, d := range domains {
		if d.Status == "active" {
			activeDomains++
		}
	}

	// Get all users
	users, err := h.userService.ListAll(ctx)
	if err != nil {
		h.logger.Error("Failed to list users for dashboard", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve dashboard statistics")
		return
	}

	totalUsers := int64(len(users))
	activeUsers := int64(0)
	var totalStorageUsed, totalStorageQuota int64
	for _, u := range users {
		if u.Status == "active" {
			activeUsers++
		}
		totalStorageUsed += u.UsedQuota
		totalStorageQuota += u.Quota
	}

	// Get queue statistics
	queueItems, err := h.queueService.GetPendingItems(ctx)
	if err != nil {
		h.logger.Warn("Failed to get queue items for dashboard", zap.Error(err))
		queueItems = []*domain.QueueItem{}
	}

	queuedMessages := int64(0)
	failedMessages := int64(0)
	for _, item := range queueItems {
		if item.Status == "pending" || item.Status == "retry" {
			queuedMessages++
		} else if item.Status == "failed" {
			failedMessages++
		}
	}

	stats := DashboardStats{
		TotalDomains:      totalDomains,
		ActiveDomains:     activeDomains,
		TotalUsers:        totalUsers,
		ActiveUsers:       activeUsers,
		QueuedMessages:    queuedMessages,
		FailedMessages:    failedMessages,
		TotalStorageUsed:  totalStorageUsed,
		TotalStorageQuota: totalStorageQuota,
		MessagesToday:     0, // Would require message repository with time-based queries
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
			QueueDepth:     queuedMessages,
			QueueHealthy:   queuedMessages < 100,
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
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	// Get domain
	domainRecord, err := h.domainService.GetByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get domain", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Domain not found")
		return
	}

	// Get all users and filter by domain
	allUsers, err := h.userService.ListAll(ctx)
	if err != nil {
		h.logger.Error("Failed to list users for domain stats", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve domain statistics")
		return
	}

	var totalUsers, activeUsers, disabledUsers int64
	var totalStorageUsed, totalStorageQuota int64
	var topUsers []UserStat

	for _, u := range allUsers {
		if u.DomainID != id {
			continue
		}

		totalUsers++
		if u.Status == "active" {
			activeUsers++
		} else {
			disabledUsers++
		}
		totalStorageUsed += u.UsedQuota
		totalStorageQuota += u.Quota

		// Build top users list
		quotaPercent := 0.0
		if u.Quota > 0 {
			quotaPercent = (float64(u.UsedQuota) / float64(u.Quota)) * 100
		}
		topUsers = append(topUsers, UserStat{
			UserID:       u.ID,
			Email:        u.Email,
			StorageUsed:  u.UsedQuota,
			MessageCount: 0, // Would require message repository with user queries
			QuotaPercent: quotaPercent,
		})
	}

	// Sort top users by storage used (descending) and limit to top 10
	if len(topUsers) > 1 {
		for i := 0; i < len(topUsers)-1; i++ {
			for j := i + 1; j < len(topUsers); j++ {
				if topUsers[j].StorageUsed > topUsers[i].StorageUsed {
					topUsers[i], topUsers[j] = topUsers[j], topUsers[i]
				}
			}
		}
	}
	if len(topUsers) > 10 {
		topUsers = topUsers[:10]
	}

	// Get aliases for this domain
	aliases, err := h.aliasService.ListByDomain(ctx, id)
	if err != nil {
		h.logger.Warn("Failed to get aliases for domain stats", zap.Error(err))
		aliases = []*domain.Alias{}
	}

	stats := DomainStats{
		DomainID:          domainRecord.ID,
		DomainName:        domainRecord.Name,
		Status:            domainRecord.Status,
		TotalUsers:        totalUsers,
		ActiveUsers:       activeUsers,
		DisabledUsers:     disabledUsers,
		TotalStorageUsed:  totalStorageUsed,
		TotalStorageQuota: totalStorageQuota,
		TotalAliases:      int64(len(aliases)),
		MessagesToday:     0, // Would require message repository with time-based queries
		MessagesThisWeek:  0,
		MessagesThisMonth: 0,
		TopUsers:          topUsers,
	}

	middleware.RespondSuccess(w, stats, "Domain statistics retrieved successfully")
}

// User retrieves statistics for a specific user
func (h *StatsHandler) User(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	// Get domain name
	domainName := ""
	domain, err := h.domainService.GetByID(ctx, user.DomainID)
	if err != nil {
		h.logger.Warn("Failed to get domain for user stats", zap.Int64("domain_id", user.DomainID), zap.Error(err))
	} else {
		domainName = domain.Name
	}

	// Calculate quota percentage
	quotaPercent := 0.0
	if user.Quota > 0 {
		quotaPercent = (float64(user.UsedQuota) / float64(user.Quota)) * 100
	}

	stats := UserStats{
		UserID:            user.ID,
		Email:             user.Email,
		DomainID:          user.DomainID,
		DomainName:        domainName,
		Status:            user.Status,
		StorageUsed:       user.UsedQuota,
		StorageQuota:      user.Quota,
		QuotaPercent:      quotaPercent,
		TotalMessages:     0, // Would require message repository with user queries
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
