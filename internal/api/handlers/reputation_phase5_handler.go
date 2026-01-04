package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"github.com/btafoya/gomailserver/internal/reputation/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ReputationPhase5Handler handles Phase 5 reputation management endpoints
// (DMARC, ARF, external metrics, provider limits, custom warmup, predictions, alerts)
type ReputationPhase5Handler struct {
	// Repositories
	dmarcRepo         repository.DMARCReportsRepository
	arfRepo           repository.ARFReportsRepository
	postmasterRepo    repository.PostmasterMetricsRepository
	sndsRepo          repository.SNDSMetricsRepository
	providerLimitsRepo repository.ProviderRateLimitsRepository
	warmupRepo        repository.CustomWarmupSchedulesRepository
	predictionsRepo   repository.PredictionsRepository
	alertsRepo        repository.AlertsRepository

	// Services
	dmarcAnalyzer     *service.DMARCAnalyzerService
	arfParser         *service.ARFParserService
	gmailPostmaster   *service.GmailPostmasterService
	microsoftSNDS     *service.MicrosoftSNDSService
	providerLimits    *service.ProviderRateLimitsService
	customWarmup      *service.CustomWarmupService
	predictions       *service.PredictionsService
	alerts            *service.AlertsService

	logger *zap.Logger
}

// NewReputationPhase5Handler creates a new Phase 5 reputation handler
func NewReputationPhase5Handler(
	dmarcRepo repository.DMARCReportsRepository,
	arfRepo repository.ARFReportsRepository,
	postmasterRepo repository.PostmasterMetricsRepository,
	sndsRepo repository.SNDSMetricsRepository,
	providerLimitsRepo repository.ProviderRateLimitsRepository,
	warmupRepo repository.CustomWarmupSchedulesRepository,
	predictionsRepo repository.PredictionsRepository,
	alertsRepo repository.AlertsRepository,
	dmarcAnalyzer *service.DMARCAnalyzerService,
	arfParser *service.ARFParserService,
	gmailPostmaster *service.GmailPostmasterService,
	microsoftSNDS *service.MicrosoftSNDSService,
	providerLimits *service.ProviderRateLimitsService,
	customWarmup *service.CustomWarmupService,
	predictions *service.PredictionsService,
	alerts *service.AlertsService,
	logger *zap.Logger,
) *ReputationPhase5Handler {
	return &ReputationPhase5Handler{
		dmarcRepo:          dmarcRepo,
		arfRepo:            arfRepo,
		postmasterRepo:     postmasterRepo,
		sndsRepo:           sndsRepo,
		providerLimitsRepo: providerLimitsRepo,
		warmupRepo:         warmupRepo,
		predictionsRepo:    predictionsRepo,
		alertsRepo:         alertsRepo,
		dmarcAnalyzer:      dmarcAnalyzer,
		arfParser:          arfParser,
		gmailPostmaster:    gmailPostmaster,
		microsoftSNDS:      microsoftSNDS,
		providerLimits:     providerLimits,
		customWarmup:       customWarmup,
		predictions:        predictions,
		alerts:             alerts,
		logger:             logger,
	}
}

// ============================================================================
// DMARC Reports Endpoints
// ============================================================================

// DMARCReportResponse represents a DMARC report in API responses
type DMARCReportResponse struct {
	ID               int64                     `json:"id"`
	ReportID         string                    `json:"report_id"`
	Domain           string                    `json:"domain"`
	OrgName          string                    `json:"org_name"`
	EmailAddress     string                    `json:"email_address"`
	BeginTime        int64                     `json:"begin_time"`
	EndTime          int64                     `json:"end_time"`
	RecordCount      int                       `json:"record_count"`
	SPFAlignedCount  int                       `json:"spf_aligned_count"`
	DKIMAlignedCount int                       `json:"dkim_aligned_count"`
	ProcessedAt      int64                     `json:"processed_at"`
	Records          []*DMARCReportRecordResponse `json:"records,omitempty"`
}

// DMARCReportRecordResponse represents a DMARC report record in API responses
type DMARCReportRecordResponse struct {
	ID             int64  `json:"id"`
	SourceIP       string `json:"source_ip"`
	Count          int    `json:"count"`
	Disposition    string `json:"disposition"`
	DMARCResult    string `json:"dmarc_result"`
	SPFDomain      string `json:"spf_domain"`
	SPFResult      string `json:"spf_result"`
	DKIMDomain     string `json:"dkim_domain"`
	DKIMResult     string `json:"dkim_result"`
	HeaderFrom     string `json:"header_from"`
}

// DMARCStatsResponse represents DMARC statistics in API responses
type DMARCStatsResponse struct {
	Domain           string  `json:"domain"`
	TotalReports     int     `json:"total_reports"`
	TotalMessages    int     `json:"total_messages"`
	SPFAlignmentRate float64 `json:"spf_alignment_rate"`
	DKIMAlignmentRate float64 `json:"dkim_alignment_rate"`
	PassRate         float64 `json:"pass_rate"`
	FailRate         float64 `json:"fail_rate"`
	RecentTrend      string  `json:"recent_trend"`
}

// DMARCActionResponse represents a DMARC auto-action in API responses
type DMARCActionResponse struct {
	ID          int64  `json:"id"`
	ReportID    int64  `json:"report_id"`
	ActionType  string `json:"action_type"`
	TargetIP    string `json:"target_ip"`
	Reason      string `json:"reason"`
	ActionTaken bool   `json:"action_taken"`
	TakenAt     *int64 `json:"taken_at,omitempty"`
}

// ListDMARCReports retrieves DMARC reports with optional filters
// GET /api/v1/reputation/dmarc/reports
func (h *ReputationPhase5Handler) ListDMARCReports(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	var reports []*domain.DMARCReport
	var err error

	if domain != "" {
		reports, err = h.dmarcRepo.GetReportsByDomain(r.Context(), domain, limit)
	} else {
		reports, err = h.dmarcRepo.GetRecentReports(r.Context(), limit)
	}

	if err != nil {
		h.logger.Error("Failed to list DMARC reports", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve DMARC reports")
		return
	}

	responses := make([]*DMARCReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = dmarcReportToResponse(report, false)
	}

	middleware.RespondSuccess(w, responses, "DMARC reports retrieved successfully")
}

// GetDMARCReport retrieves a specific DMARC report with all records
// GET /api/v1/reputation/dmarc/reports/:id
func (h *ReputationPhase5Handler) GetDMARCReport(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid report ID")
		return
	}

	report, err := h.dmarcRepo.GetReportByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get DMARC report", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "DMARC report not found")
		return
	}

	response := dmarcReportToResponse(report, true)
	middleware.RespondSuccess(w, response, "DMARC report retrieved successfully")
}

// GetDMARCStats retrieves DMARC statistics for a domain
// GET /api/v1/reputation/dmarc/stats/:domain
func (h *ReputationPhase5Handler) GetDMARCStats(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	stats, err := h.dmarcAnalyzer.GetDomainStats(r.Context(), domainName, days)
	if err != nil {
		h.logger.Error("Failed to get DMARC stats", zap.String("domain", domainName), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve DMARC statistics")
		return
	}

	middleware.RespondSuccess(w, stats, "DMARC statistics retrieved successfully")
}

// GetDMARCActions retrieves auto-actions taken based on DMARC reports
// GET /api/v1/reputation/dmarc/actions
func (h *ReputationPhase5Handler) GetDMARCActions(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	actions, err := h.dmarcRepo.GetRecentActions(r.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to list DMARC actions", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve DMARC actions")
		return
	}

	responses := make([]*DMARCActionResponse, len(actions))
	for i, action := range actions {
		responses[i] = dmarcActionToResponse(action)
	}

	middleware.RespondSuccess(w, responses, "DMARC actions retrieved successfully")
}

// ExportDMARCReport exports a DMARC report in specified format
// POST /api/v1/reputation/dmarc/reports/:id/export
func (h *ReputationPhase5Handler) ExportDMARCReport(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid report ID")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	report, err := h.dmarcRepo.GetReportByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get DMARC report for export", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "DMARC report not found")
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=dmarc-report-"+strconv.FormatInt(id, 10)+".json")
		json.NewEncoder(w).Encode(dmarcReportToResponse(report, true))
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=dmarc-report-"+strconv.FormatInt(id, 10)+".csv")
		// TODO: Implement CSV export
		middleware.RespondError(w, http.StatusNotImplemented, "CSV export not yet implemented")
	default:
		middleware.RespondError(w, http.StatusBadRequest, "Unsupported export format")
	}
}

// ============================================================================
// ARF Reports Endpoints
// ============================================================================

// ARFReportResponse represents an ARF report in API responses
type ARFReportResponse struct {
	ID               int64  `json:"id"`
	FeedbackType     string `json:"feedback_type"`
	UserAgent        string `json:"user_agent"`
	Version          string `json:"version"`
	SourceIP         string `json:"source_ip"`
	IncidentCount    int    `json:"incident_count"`
	OriginalMailFrom string `json:"original_mail_from"`
	OriginalRcptTo   string `json:"original_rcpt_to"`
	ReportedDomain   string `json:"reported_domain"`
	ReportedURI      string `json:"reported_uri"`
	AuthFailure      string `json:"auth_failure"`
	DeliveryResult   string `json:"delivery_result"`
	ReceivedDate     int64  `json:"received_date"`
	ProcessedAt      int64  `json:"processed_at"`
	ActionTaken      bool   `json:"action_taken"`
}

// ARFStatsResponse represents ARF statistics in API responses
type ARFStatsResponse struct {
	TotalReports      int                `json:"total_reports"`
	ByFeedbackType    map[string]int     `json:"by_feedback_type"`
	ByDomain          map[string]int     `json:"by_domain"`
	ActionTakenCount  int                `json:"action_taken_count"`
	RecentTrend       string             `json:"recent_trend"`
}

// ListARFReports retrieves ARF complaint reports with optional filters
// GET /api/v1/reputation/arf/reports
func (h *ReputationPhase5Handler) ListARFReports(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	var reports []*domain.ARFReport
	var err error

	if domain != "" {
		reports, err = h.arfRepo.GetReportsByDomain(r.Context(), domain, limit)
	} else {
		reports, err = h.arfRepo.GetRecentReports(r.Context(), limit)
	}

	if err != nil {
		h.logger.Error("Failed to list ARF reports", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve ARF reports")
		return
	}

	responses := make([]*ARFReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = arfReportToResponse(report)
	}

	middleware.RespondSuccess(w, responses, "ARF reports retrieved successfully")
}

// GetARFStats retrieves ARF complaint statistics
// GET /api/v1/reputation/arf/stats
func (h *ReputationPhase5Handler) GetARFStats(w http.ResponseWriter, r *http.Request) {
	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	stats, err := h.arfParser.GetARFStats(r.Context(), days)
	if err != nil {
		h.logger.Error("Failed to get ARF stats", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve ARF statistics")
		return
	}

	middleware.RespondSuccess(w, stats, "ARF statistics retrieved successfully")
}

// ProcessARFReport manually processes a specific ARF report
// POST /api/v1/reputation/arf/reports/:id/process
func (h *ReputationPhase5Handler) ProcessARFReport(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid report ID")
		return
	}

	report, err := h.arfRepo.GetReportByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get ARF report", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "ARF report not found")
		return
	}

	if err := h.arfParser.ProcessReport(r.Context(), report); err != nil {
		h.logger.Error("Failed to process ARF report", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to process ARF report")
		return
	}

	middleware.RespondSuccess(w, nil, "ARF report processed successfully")
}

// ============================================================================
// External Metrics Endpoints
// ============================================================================

// PostmasterMetricsResponse represents Gmail Postmaster metrics in API responses
type PostmasterMetricsResponse struct {
	ID                int64   `json:"id"`
	Domain            string  `json:"domain"`
	Date              string  `json:"date"`
	SpamRate          float64 `json:"spam_rate"`
	IPReputation      string  `json:"ip_reputation"`
	DomainReputation  string  `json:"domain_reputation"`
	FeedbackLoopRate  float64 `json:"feedback_loop_rate"`
	AuthenticationRate float64 `json:"authentication_rate"`
	EncryptionRate    float64 `json:"encryption_rate"`
	DeliveryErrors    int     `json:"delivery_errors"`
	SyncedAt          int64   `json:"synced_at"`
}

// SNDSMetricsResponse represents Microsoft SNDS metrics in API responses
type SNDSMetricsResponse struct {
	ID             int64   `json:"id"`
	IPAddress      string  `json:"ip_address"`
	Date           string  `json:"date"`
	MessageCount   int     `json:"message_count"`
	FilterResult   string  `json:"filter_result"`
	ComplaintRate  float64 `json:"complaint_rate"`
	TrapHits       int     `json:"trap_hits"`
	SampleData     int     `json:"sample_data"`
	RCPT           int     `json:"rcpt"`
	SyncedAt       int64   `json:"synced_at"`
}

// MetricsTrendResponse represents trend data for external metrics
type MetricsTrendResponse struct {
	Metric     string              `json:"metric"`
	Domain     string              `json:"domain,omitempty"`
	IPAddress  string              `json:"ip_address,omitempty"`
	DataPoints []*TrendDataPoint   `json:"data_points"`
	Trend      string              `json:"trend"`
	Change     float64             `json:"change"`
}

// TrendDataPoint represents a single data point in a trend
type TrendDataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// GetPostmasterMetrics retrieves Gmail Postmaster metrics for a domain
// GET /api/v1/reputation/external/postmaster/:domain
func (h *ReputationPhase5Handler) GetPostmasterMetrics(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	metrics, err := h.postmasterRepo.GetMetricsByDomain(r.Context(), domainName, days)
	if err != nil {
		h.logger.Error("Failed to get Postmaster metrics", zap.String("domain", domainName), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve Postmaster metrics")
		return
	}

	responses := make([]*PostmasterMetricsResponse, len(metrics))
	for i, metric := range metrics {
		responses[i] = postmasterMetricsToResponse(metric)
	}

	middleware.RespondSuccess(w, responses, "Postmaster metrics retrieved successfully")
}

// GetSNDSMetrics retrieves Microsoft SNDS metrics for an IP address
// GET /api/v1/reputation/external/snds/:ip
func (h *ReputationPhase5Handler) GetSNDSMetrics(w http.ResponseWriter, r *http.Request) {
	ipAddress := chi.URLParam(r, "ip")
	if ipAddress == "" {
		middleware.RespondError(w, http.StatusBadRequest, "IP address parameter required")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	metrics, err := h.sndsRepo.GetMetricsByIP(r.Context(), ipAddress, days)
	if err != nil {
		h.logger.Error("Failed to get SNDS metrics", zap.String("ip", ipAddress), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve SNDS metrics")
		return
	}

	responses := make([]*SNDSMetricsResponse, len(metrics))
	for i, metric := range metrics {
		responses[i] = sndsMetricsToResponse(metric)
	}

	middleware.RespondSuccess(w, responses, "SNDS metrics retrieved successfully")
}

// GetExternalMetricsTrends retrieves trend analysis for external metrics
// GET /api/v1/reputation/external/trends
func (h *ReputationPhase5Handler) GetExternalMetricsTrends(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	ipAddress := r.URL.Query().Get("ip")

	if domain == "" && ipAddress == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Either domain or IP address parameter required")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	trends := make([]*MetricsTrendResponse, 0)

	if domain != "" {
		// Get Postmaster trends
		postmasterTrends, err := h.gmailPostmaster.GetTrends(r.Context(), domain, days)
		if err == nil {
			trends = append(trends, postmasterTrends...)
		}
	}

	if ipAddress != "" {
		// Get SNDS trends
		sndsTrends, err := h.microsoftSNDS.GetTrends(r.Context(), ipAddress, days)
		if err == nil {
			trends = append(trends, sndsTrends...)
		}
	}

	middleware.RespondSuccess(w, trends, "External metrics trends retrieved successfully")
}

// ============================================================================
// Provider Rate Limits Endpoints
// ============================================================================

// ProviderRateLimitResponse represents a provider rate limit in API responses
type ProviderRateLimitResponse struct {
	ID                  int64  `json:"id"`
	Domain              string `json:"domain"`
	Provider            string `json:"provider"`
	MessagesPerHour     int    `json:"messages_per_hour"`
	MessagesPerDay      int    `json:"messages_per_day"`
	ConnectionsPerHour  int    `json:"connections_per_hour"`
	MaxRecipientsPerMsg int    `json:"max_recipients_per_msg"`
	CurrentUsageHour    int    `json:"current_usage_hour"`
	CurrentUsageDay     int    `json:"current_usage_day"`
	LastResetHour       int64  `json:"last_reset_hour"`
	LastResetDay        int64  `json:"last_reset_day"`
	UpdatedAt           int64  `json:"updated_at"`
}

// ListProviderRateLimits retrieves all provider-specific rate limits
// GET /api/v1/reputation/provider-limits
func (h *ReputationPhase5Handler) ListProviderRateLimits(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")

	var limits []*domain.ProviderRateLimit
	var err error

	if domain != "" {
		limits, err = h.providerLimitsRepo.GetLimitsByDomain(r.Context(), domain)
	} else {
		limits, err = h.providerLimitsRepo.GetAllLimits(r.Context())
	}

	if err != nil {
		h.logger.Error("Failed to list provider rate limits", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve provider rate limits")
		return
	}

	responses := make([]*ProviderRateLimitResponse, len(limits))
	for i, limit := range limits {
		responses[i] = providerRateLimitToResponse(limit)
	}

	middleware.RespondSuccess(w, responses, "Provider rate limits retrieved successfully")
}

// UpdateProviderRateLimit updates a provider rate limit
// PUT /api/v1/reputation/provider-limits/:id
func (h *ReputationPhase5Handler) UpdateProviderRateLimit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid limit ID")
		return
	}

	var req struct {
		MessagesPerHour     *int `json:"messages_per_hour"`
		MessagesPerDay      *int `json:"messages_per_day"`
		ConnectionsPerHour  *int `json:"connections_per_hour"`
		MaxRecipientsPerMsg *int `json:"max_recipients_per_msg"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	limit, err := h.providerLimitsRepo.GetLimitByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get provider rate limit", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Provider rate limit not found")
		return
	}

	if req.MessagesPerHour != nil {
		limit.MessagesPerHour = *req.MessagesPerHour
	}
	if req.MessagesPerDay != nil {
		limit.MessagesPerDay = *req.MessagesPerDay
	}
	if req.ConnectionsPerHour != nil {
		limit.ConnectionsPerHour = *req.ConnectionsPerHour
	}
	if req.MaxRecipientsPerMsg != nil {
		limit.MaxRecipientsPerMsg = *req.MaxRecipientsPerMsg
	}

	if err := h.providerLimitsRepo.UpdateLimit(r.Context(), limit); err != nil {
		h.logger.Error("Failed to update provider rate limit", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update provider rate limit")
		return
	}

	response := providerRateLimitToResponse(limit)
	middleware.RespondSuccess(w, response, "Provider rate limit updated successfully")
}

// InitializeProviderLimits initializes rate limits for a domain
// POST /api/v1/reputation/provider-limits/init/:domain
func (h *ReputationPhase5Handler) InitializeProviderLimits(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	if err := h.providerLimits.InitializeLimits(r.Context(), domainName); err != nil {
		h.logger.Error("Failed to initialize provider limits", zap.String("domain", domainName), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to initialize provider limits")
		return
	}

	middleware.RespondSuccess(w, nil, "Provider limits initialized successfully")
}

// ResetProviderUsage resets usage counters for a provider limit
// POST /api/v1/reputation/provider-limits/:id/reset
func (h *ReputationPhase5Handler) ResetProviderUsage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid limit ID")
		return
	}

	if err := h.providerLimits.ResetUsage(r.Context(), id); err != nil {
		h.logger.Error("Failed to reset provider usage", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to reset provider usage")
		return
	}

	middleware.RespondSuccess(w, nil, "Provider usage reset successfully")
}

// ============================================================================
// Custom Warmup Endpoints
// ============================================================================

// CustomWarmupScheduleResponse represents a custom warmup schedule in API responses
type CustomWarmupScheduleResponse struct {
	ID          int64                  `json:"id"`
	Domain      string                 `json:"domain"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	TotalDays   int                    `json:"total_days"`
	CurrentDay  int                    `json:"current_day"`
	IsActive    bool                   `json:"is_active"`
	StartedAt   *int64                 `json:"started_at,omitempty"`
	CompletedAt *int64                 `json:"completed_at,omitempty"`
	CreatedAt   int64                  `json:"created_at"`
	DailyLimits []*WarmupDayLimit      `json:"daily_limits,omitempty"`
}

// WarmupDayLimit represents a daily limit in a warmup schedule
type WarmupDayLimit struct {
	Day          int `json:"day"`
	MessageLimit int `json:"message_limit"`
}

// GetCustomWarmupSchedule retrieves the active warmup schedule for a domain
// GET /api/v1/reputation/warmup/:domain
func (h *ReputationPhase5Handler) GetCustomWarmupSchedule(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	schedule, err := h.warmupRepo.GetActiveSchedule(r.Context(), domainName)
	if err != nil {
		h.logger.Error("Failed to get warmup schedule", zap.String("domain", domainName), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Warmup schedule not found")
		return
	}

	response := customWarmupToResponse(schedule, true)
	middleware.RespondSuccess(w, response, "Warmup schedule retrieved successfully")
}

// CreateCustomWarmupSchedule creates a new custom warmup schedule
// POST /api/v1/reputation/warmup
func (h *ReputationPhase5Handler) CreateCustomWarmupSchedule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain      string             `json:"domain"`
		Name        string             `json:"name"`
		Description string             `json:"description"`
		DailyLimits []*WarmupDayLimit  `json:"daily_limits"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Domain == "" || req.Name == "" || len(req.DailyLimits) == 0 {
		middleware.RespondError(w, http.StatusBadRequest, "Domain, name, and daily limits are required")
		return
	}

	schedule := &domain.CustomWarmupSchedule{
		Domain:      req.Domain,
		Name:        req.Name,
		Description: req.Description,
		TotalDays:   len(req.DailyLimits),
		CurrentDay:  0,
		IsActive:    false,
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.warmupRepo.CreateSchedule(r.Context(), schedule, req.DailyLimits); err != nil {
		h.logger.Error("Failed to create warmup schedule", zap.String("domain", req.Domain), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to create warmup schedule")
		return
	}

	response := customWarmupToResponse(schedule, true)
	middleware.RespondSuccess(w, response, "Warmup schedule created successfully")
}

// UpdateCustomWarmupSchedule updates an existing warmup schedule
// PUT /api/v1/reputation/warmup/:id
func (h *ReputationPhase5Handler) UpdateCustomWarmupSchedule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var req struct {
		Name        *string            `json:"name"`
		Description *string            `json:"description"`
		DailyLimits []*WarmupDayLimit  `json:"daily_limits"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	schedule, err := h.warmupRepo.GetScheduleByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get warmup schedule", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Warmup schedule not found")
		return
	}

	if req.Name != nil {
		schedule.Name = *req.Name
	}
	if req.Description != nil {
		schedule.Description = *req.Description
	}

	if err := h.warmupRepo.UpdateSchedule(r.Context(), schedule, req.DailyLimits); err != nil {
		h.logger.Error("Failed to update warmup schedule", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update warmup schedule")
		return
	}

	response := customWarmupToResponse(schedule, true)
	middleware.RespondSuccess(w, response, "Warmup schedule updated successfully")
}

// DeleteCustomWarmupSchedule deletes a warmup schedule
// DELETE /api/v1/reputation/warmup/:id
func (h *ReputationPhase5Handler) DeleteCustomWarmupSchedule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	if err := h.warmupRepo.DeleteSchedule(r.Context(), id); err != nil {
		h.logger.Error("Failed to delete warmup schedule", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to delete warmup schedule")
		return
	}

	middleware.RespondSuccess(w, nil, "Warmup schedule deleted successfully")
}

// GetWarmupTemplates retrieves available warmup schedule templates
// GET /api/v1/reputation/warmup/templates
func (h *ReputationPhase5Handler) GetWarmupTemplates(w http.ResponseWriter, r *http.Request) {
	templates := h.customWarmup.GetTemplates(r.Context())
	middleware.RespondSuccess(w, templates, "Warmup templates retrieved successfully")
}

// ============================================================================
// Predictions Endpoints
// ============================================================================

// PredictionResponse represents a reputation prediction in API responses
type PredictionResponse struct {
	ID                 int64              `json:"id"`
	Domain             string             `json:"domain"`
	PredictionDate     string             `json:"prediction_date"`
	Horizon            string             `json:"horizon"`
	PredictedScore     int                `json:"predicted_score"`
	Confidence         float64            `json:"confidence"`
	PredictedBounce    float64            `json:"predicted_bounce"`
	PredictedComplaint float64            `json:"predicted_complaint"`
	TrendDirection     string             `json:"trend_direction"`
	RiskLevel          string             `json:"risk_level"`
	RecommendedActions []string           `json:"recommended_actions"`
	FeatureImportance  map[string]float64 `json:"feature_importance,omitempty"`
	GeneratedAt        int64              `json:"generated_at"`
}

// GetLatestPredictions retrieves the latest predictions for all domains
// GET /api/v1/reputation/predictions/latest
func (h *ReputationPhase5Handler) GetLatestPredictions(w http.ResponseWriter, r *http.Request) {
	predictions, err := h.predictionsRepo.GetLatestPredictions(r.Context())
	if err != nil {
		h.logger.Error("Failed to get latest predictions", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve predictions")
		return
	}

	responses := make([]*PredictionResponse, len(predictions))
	for i, pred := range predictions {
		responses[i] = predictionToResponse(pred, true)
	}

	middleware.RespondSuccess(w, responses, "Predictions retrieved successfully")
}

// GetDomainPredictions retrieves predictions for a specific domain
// GET /api/v1/reputation/predictions/:domain
func (h *ReputationPhase5Handler) GetDomainPredictions(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	horizon := r.URL.Query().Get("horizon")
	if horizon == "" {
		horizon = "7d"
	}

	prediction, err := h.predictionsRepo.GetPredictionByDomain(r.Context(), domainName, horizon)
	if err != nil {
		h.logger.Error("Failed to get domain prediction", zap.String("domain", domainName), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Prediction not found")
		return
	}

	response := predictionToResponse(prediction, true)
	middleware.RespondSuccess(w, response, "Prediction retrieved successfully")
}

// GeneratePredictions triggers prediction generation for a domain
// POST /api/v1/reputation/predictions/generate/:domain
func (h *ReputationPhase5Handler) GeneratePredictions(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	if err := h.predictions.GeneratePredictions(r.Context(), domainName); err != nil {
		h.logger.Error("Failed to generate predictions", zap.String("domain", domainName), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to generate predictions")
		return
	}

	middleware.RespondSuccess(w, nil, "Predictions generated successfully")
}

// GetPredictionHistory retrieves historical predictions for a domain
// GET /api/v1/reputation/predictions/:domain/history
func (h *ReputationPhase5Handler) GetPredictionHistory(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	predictions, err := h.predictionsRepo.GetPredictionHistory(r.Context(), domainName, days)
	if err != nil {
		h.logger.Error("Failed to get prediction history", zap.String("domain", domainName), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve prediction history")
		return
	}

	responses := make([]*PredictionResponse, len(predictions))
	for i, pred := range predictions {
		responses[i] = predictionToResponse(pred, false)
	}

	middleware.RespondSuccess(w, responses, "Prediction history retrieved successfully")
}

// ============================================================================
// Alerts Endpoints (Phase 5 Extensions)
// ============================================================================

// Phase5AlertResponse represents a Phase 5 alert in API responses
type Phase5AlertResponse struct {
	ID             int64  `json:"id"`
	Domain         string `json:"domain"`
	AlertType      string `json:"alert_type"`
	Severity       string `json:"severity"`
	Title          string `json:"title"`
	Message        string `json:"message"`
	SourceType     string `json:"source_type"`
	SourceID       *int64 `json:"source_id,omitempty"`
	Acknowledged   bool   `json:"acknowledged"`
	AcknowledgedAt *int64 `json:"acknowledged_at,omitempty"`
	AcknowledgedBy string `json:"acknowledged_by,omitempty"`
	Resolved       bool   `json:"resolved"`
	ResolvedAt     *int64 `json:"resolved_at,omitempty"`
	CreatedAt      int64  `json:"created_at"`
}

// ListPhase5Alerts retrieves all Phase 5 alerts with optional filters
// GET /api/v1/reputation/alerts/phase5
func (h *ReputationPhase5Handler) ListPhase5Alerts(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	severity := r.URL.Query().Get("severity")
	unacknowledgedOnly := r.URL.Query().Get("unacknowledged") == "true"

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	var alerts []*domain.Alert
	var err error

	if unacknowledgedOnly {
		alerts, err = h.alertsRepo.GetUnacknowledgedAlerts(r.Context(), limit)
	} else if domain != "" {
		alerts, err = h.alertsRepo.GetAlertsByDomain(r.Context(), domain, limit)
	} else {
		alerts, err = h.alertsRepo.GetRecentAlerts(r.Context(), limit)
	}

	if err != nil {
		h.logger.Error("Failed to list alerts", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve alerts")
		return
	}

	// Filter by severity if specified
	if severity != "" {
		filtered := make([]*domain.Alert, 0)
		for _, alert := range alerts {
			if alert.Severity == severity {
				filtered = append(filtered, alert)
			}
		}
		alerts = filtered
	}

	responses := make([]*Phase5AlertResponse, len(alerts))
	for i, alert := range alerts {
		responses[i] = phase5AlertToResponse(alert)
	}

	middleware.RespondSuccess(w, responses, "Alerts retrieved successfully")
}

// AcknowledgeAlert marks an alert as acknowledged
// POST /api/v1/reputation/alerts/:id/acknowledge
func (h *ReputationPhase5Handler) AcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid alert ID")
		return
	}

	var req struct {
		AcknowledgedBy string `json:"acknowledged_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.alerts.AcknowledgeAlert(r.Context(), id, req.AcknowledgedBy); err != nil {
		h.logger.Error("Failed to acknowledge alert", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to acknowledge alert")
		return
	}

	middleware.RespondSuccess(w, nil, "Alert acknowledged successfully")
}

// ResolveAlert marks an alert as resolved
// POST /api/v1/reputation/alerts/:id/resolve
func (h *ReputationPhase5Handler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid alert ID")
		return
	}

	if err := h.alerts.ResolveAlert(r.Context(), id); err != nil {
		h.logger.Error("Failed to resolve alert", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to resolve alert")
		return
	}

	middleware.RespondSuccess(w, nil, "Alert resolved successfully")
}

// ============================================================================
// Helper Functions - Response Converters
// ============================================================================

func dmarcReportToResponse(report *domain.DMARCReport, includeRecords bool) *DMARCReportResponse {
	resp := &DMARCReportResponse{
		ID:               report.ID,
		ReportID:         report.ReportID,
		Domain:           report.Domain,
		OrgName:          report.OrgName,
		EmailAddress:     report.EmailAddress,
		BeginTime:        report.BeginTime,
		EndTime:          report.EndTime,
		RecordCount:      report.RecordCount,
		SPFAlignedCount:  report.SPFAlignedCount,
		DKIMAlignedCount: report.DKIMAlignedCount,
		ProcessedAt:      report.ProcessedAt,
	}

	if includeRecords && report.Records != nil {
		resp.Records = make([]*DMARCReportRecordResponse, len(report.Records))
		for i, record := range report.Records {
			resp.Records[i] = &DMARCReportRecordResponse{
				ID:             record.ID,
				SourceIP:       record.SourceIP,
				Count:          record.Count,
				Disposition:    record.Disposition,
				DMARCResult:    record.DMARCResult,
				SPFDomain:      record.SPFDomain,
				SPFResult:      record.SPFResult,
				DKIMDomain:     record.DKIMDomain,
				DKIMResult:     record.DKIMResult,
				HeaderFrom:     record.HeaderFrom,
			}
		}
	}

	return resp
}

func dmarcActionToResponse(action *domain.DMARCAutoAction) *DMARCActionResponse {
	return &DMARCActionResponse{
		ID:          action.ID,
		ReportID:    action.ReportID,
		ActionType:  action.ActionType,
		TargetIP:    action.TargetIP,
		Reason:      action.Reason,
		ActionTaken: action.ActionTaken,
		TakenAt:     action.TakenAt,
	}
}

func arfReportToResponse(report *domain.ARFReport) *ARFReportResponse {
	return &ARFReportResponse{
		ID:               report.ID,
		FeedbackType:     report.FeedbackType,
		UserAgent:        report.UserAgent,
		Version:          report.Version,
		SourceIP:         report.SourceIP,
		IncidentCount:    report.IncidentCount,
		OriginalMailFrom: report.OriginalMailFrom,
		OriginalRcptTo:   report.OriginalRcptTo,
		ReportedDomain:   report.ReportedDomain,
		ReportedURI:      report.ReportedURI,
		AuthFailure:      report.AuthFailure,
		DeliveryResult:   report.DeliveryResult,
		ReceivedDate:     report.ReceivedDate,
		ProcessedAt:      report.ProcessedAt,
		ActionTaken:      report.ActionTaken,
	}
}

func postmasterMetricsToResponse(metrics *domain.PostmasterMetrics) *PostmasterMetricsResponse {
	return &PostmasterMetricsResponse{
		ID:                 metrics.ID,
		Domain:             metrics.Domain,
		Date:               metrics.Date,
		SpamRate:           metrics.SpamRate,
		IPReputation:       metrics.IPReputation,
		DomainReputation:   metrics.DomainReputation,
		FeedbackLoopRate:   metrics.FeedbackLoopRate,
		AuthenticationRate: metrics.AuthenticationRate,
		EncryptionRate:     metrics.EncryptionRate,
		DeliveryErrors:     metrics.DeliveryErrors,
		SyncedAt:           metrics.SyncedAt,
	}
}

func sndsMetricsToResponse(metrics *domain.SNDSMetrics) *SNDSMetricsResponse {
	return &SNDSMetricsResponse{
		ID:            metrics.ID,
		IPAddress:     metrics.IPAddress,
		Date:          metrics.Date,
		MessageCount:  metrics.MessageCount,
		FilterResult:  metrics.FilterResult,
		ComplaintRate: metrics.ComplaintRate,
		TrapHits:      metrics.TrapHits,
		SampleData:    metrics.SampleData,
		RCPT:          metrics.RCPT,
		SyncedAt:      metrics.SyncedAt,
	}
}

func providerRateLimitToResponse(limit *domain.ProviderRateLimit) *ProviderRateLimitResponse {
	return &ProviderRateLimitResponse{
		ID:                  limit.ID,
		Domain:              limit.Domain,
		Provider:            limit.Provider,
		MessagesPerHour:     limit.MessagesPerHour,
		MessagesPerDay:      limit.MessagesPerDay,
		ConnectionsPerHour:  limit.ConnectionsPerHour,
		MaxRecipientsPerMsg: limit.MaxRecipientsPerMsg,
		CurrentUsageHour:    limit.CurrentUsageHour,
		CurrentUsageDay:     limit.CurrentUsageDay,
		LastResetHour:       limit.LastResetHour,
		LastResetDay:        limit.LastResetDay,
		UpdatedAt:           limit.UpdatedAt,
	}
}

func customWarmupToResponse(schedule *domain.CustomWarmupSchedule, includeLimits bool) *CustomWarmupScheduleResponse {
	resp := &CustomWarmupScheduleResponse{
		ID:          schedule.ID,
		Domain:      schedule.Domain,
		Name:        schedule.Name,
		Description: schedule.Description,
		TotalDays:   schedule.TotalDays,
		CurrentDay:  schedule.CurrentDay,
		IsActive:    schedule.IsActive,
		StartedAt:   schedule.StartedAt,
		CompletedAt: schedule.CompletedAt,
		CreatedAt:   schedule.CreatedAt,
	}

	if includeLimits && schedule.DailyLimits != nil {
		resp.DailyLimits = make([]*WarmupDayLimit, len(schedule.DailyLimits))
		for i, limit := range schedule.DailyLimits {
			resp.DailyLimits[i] = &WarmupDayLimit{
				Day:          limit.Day,
				MessageLimit: limit.MessageLimit,
			}
		}
	}

	return resp
}

func predictionToResponse(pred *domain.Prediction, includeFeatures bool) *PredictionResponse {
	resp := &PredictionResponse{
		ID:                 pred.ID,
		Domain:             pred.Domain,
		PredictionDate:     pred.PredictionDate,
		Horizon:            pred.Horizon,
		PredictedScore:     pred.PredictedScore,
		Confidence:         pred.Confidence,
		PredictedBounce:    pred.PredictedBounce,
		PredictedComplaint: pred.PredictedComplaint,
		TrendDirection:     pred.TrendDirection,
		RiskLevel:          pred.RiskLevel,
		RecommendedActions: pred.RecommendedActions,
		GeneratedAt:        pred.GeneratedAt,
	}

	if includeFeatures {
		resp.FeatureImportance = pred.FeatureImportance
	}

	return resp
}

func phase5AlertToResponse(alert *domain.Alert) *Phase5AlertResponse {
	return &Phase5AlertResponse{
		ID:             alert.ID,
		Domain:         alert.Domain,
		AlertType:      alert.AlertType,
		Severity:       alert.Severity,
		Title:          alert.Title,
		Message:        alert.Message,
		SourceType:     alert.SourceType,
		SourceID:       alert.SourceID,
		Acknowledged:   alert.Acknowledged,
		AcknowledgedAt: alert.AcknowledgedAt,
		AcknowledgedBy: alert.AcknowledgedBy,
		Resolved:       alert.Resolved,
		ResolvedAt:     alert.ResolvedAt,
		CreatedAt:      alert.CreatedAt,
	}
}
