package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"github.com/btafoya/gomailserver/internal/reputation/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ReputationHandler handles reputation management endpoints
type ReputationHandler struct {
	auditorService *service.AuditorService
	scoresRepo     repository.ScoresRepository
	eventsRepo     repository.EventsRepository
	circuitRepo    repository.CircuitBreakerRepository
	logger         *zap.Logger
}

// NewReputationHandler creates a new reputation handler
func NewReputationHandler(
	auditorService *service.AuditorService,
	scoresRepo repository.ScoresRepository,
	eventsRepo repository.EventsRepository,
	circuitRepo repository.CircuitBreakerRepository,
	logger *zap.Logger,
) *ReputationHandler {
	return &ReputationHandler{
		auditorService: auditorService,
		scoresRepo:     scoresRepo,
		eventsRepo:     eventsRepo,
		circuitRepo:    circuitRepo,
		logger:         logger,
	}
}

// AuditRequest represents a request to audit a domain
type AuditRequest struct {
	SendingIP string `json:"sending_ip,omitempty"`
}

// AuditResponse represents an audit result in API responses
type AuditResponse struct {
	Domain       string                 `json:"domain"`
	Timestamp    int64                  `json:"timestamp"`
	SPF          CheckStatusResponse    `json:"spf"`
	DKIM         CheckStatusResponse    `json:"dkim"`
	DMARC        CheckStatusResponse    `json:"dmarc"`
	RDNS         CheckStatusResponse    `json:"rdns"`
	FCrDNS       CheckStatusResponse    `json:"fcrdns"`
	TLS          CheckStatusResponse    `json:"tls"`
	MTASTS       CheckStatusResponse    `json:"mta_sts"`
	PostmasterOK bool                   `json:"postmaster_ok"`
	AbuseOK      bool                   `json:"abuse_ok"`
	OverallScore int                    `json:"overall_score"`
	Issues       []string               `json:"issues"`
}

// CheckStatusResponse represents a check status in API responses
type CheckStatusResponse struct {
	Passed  bool                   `json:"passed"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ScoreResponse represents a reputation score in API responses
type ScoreResponse struct {
	Domain               string  `json:"domain"`
	ReputationScore      int     `json:"reputation_score"`
	ComplaintRate        float64 `json:"complaint_rate"`
	BounceRate           float64 `json:"bounce_rate"`
	DeliveryRate         float64 `json:"delivery_rate"`
	CircuitBreakerActive bool    `json:"circuit_breaker_active"`
	CircuitBreakerReason string  `json:"circuit_breaker_reason,omitempty"`
	WarmUpActive         bool    `json:"warm_up_active"`
	WarmUpDay            int     `json:"warm_up_day,omitempty"`
	LastUpdated          int64   `json:"last_updated"`
}

// CircuitBreakerResponse represents a circuit breaker event in API responses
type CircuitBreakerResponse struct {
	ID           int64   `json:"id"`
	Domain       string  `json:"domain"`
	TriggerType  string  `json:"trigger_type"`
	TriggerValue float64 `json:"trigger_value"`
	Threshold    float64 `json:"threshold"`
	PausedAt     int64   `json:"paused_at"`
	ResumedAt    *int64  `json:"resumed_at,omitempty"`
	AutoResumed  bool    `json:"auto_resumed"`
	AdminNotes   string  `json:"admin_notes,omitempty"`
}

// AlertResponse represents a recent alert in API responses
type AlertResponse struct {
	Timestamp int64  `json:"timestamp"`
	Domain    string `json:"domain"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
}

// AuditDomain performs a deliverability audit for a domain
// GET /api/v1/reputation/audit/:domain
func (h *ReputationHandler) AuditDomain(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	// Parse optional sending IP from query or body
	var req AuditRequest
	if r.Method == http.MethodPost {
		json.NewDecoder(r.Body).Decode(&req)
	} else {
		req.SendingIP = r.URL.Query().Get("sending_ip")
	}

	var sendingIP net.IP
	if req.SendingIP != "" {
		sendingIP = net.ParseIP(req.SendingIP)
		if sendingIP == nil {
			middleware.RespondError(w, http.StatusBadRequest, "Invalid IP address format")
			return
		}
	}

	// Perform audit
	result, err := h.auditorService.AuditDomain(r.Context(), domainName, sendingIP)
	if err != nil {
		h.logger.Error("Failed to audit domain",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to perform audit")
		return
	}

	// Convert to response format
	response := auditToResponse(result)
	middleware.RespondSuccess(w, response, "Audit completed successfully")
}

// ListScores retrieves all domain reputation scores
// GET /api/v1/reputation/scores
func (h *ReputationHandler) ListScores(w http.ResponseWriter, r *http.Request) {
	scores, err := h.scoresRepo.ListAllScores(r.Context())
	if err != nil {
		h.logger.Error("Failed to list reputation scores", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve reputation scores")
		return
	}

	// Convert to response format
	responses := make([]*ScoreResponse, len(scores))
	for i, score := range scores {
		responses[i] = scoreToResponse(score)
	}

	middleware.RespondSuccess(w, responses, "Reputation scores retrieved successfully")
}

// GetScore retrieves reputation score for a specific domain
// GET /api/v1/reputation/scores/:domain
func (h *ReputationHandler) GetScore(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	score, err := h.scoresRepo.GetReputationScore(r.Context(), domainName)
	if err != nil {
		h.logger.Error("Failed to get reputation score",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusNotFound, "Reputation score not found")
		return
	}

	response := scoreToResponse(score)
	middleware.RespondSuccess(w, response, "Reputation score retrieved successfully")
}

// ListCircuitBreakers retrieves all active circuit breakers
// GET /api/v1/reputation/circuit-breakers
func (h *ReputationHandler) ListCircuitBreakers(w http.ResponseWriter, r *http.Request) {
	breakers, err := h.circuitRepo.GetActiveBreakers(r.Context())
	if err != nil {
		h.logger.Error("Failed to list circuit breakers", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve circuit breakers")
		return
	}

	// Convert to response format
	responses := make([]*CircuitBreakerResponse, len(breakers))
	for i, breaker := range breakers {
		responses[i] = circuitBreakerToResponse(breaker)
	}

	middleware.RespondSuccess(w, responses, "Circuit breakers retrieved successfully")
}

// GetCircuitBreakerHistory retrieves circuit breaker history for a domain
// GET /api/v1/reputation/circuit-breakers/:domain/history
func (h *ReputationHandler) GetCircuitBreakerHistory(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domain")
	if domainName == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain parameter required")
		return
	}

	// Get limit from query parameter (default: 10, max: 100)
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := parseInt(limitStr, 1, 100); err == nil {
			limit = l
		}
	}

	history, err := h.circuitRepo.GetBreakerHistory(r.Context(), domainName, limit)
	if err != nil {
		h.logger.Error("Failed to get circuit breaker history",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve circuit breaker history")
		return
	}

	// Convert to response format
	responses := make([]*CircuitBreakerResponse, len(history))
	for i, breaker := range history {
		responses[i] = circuitBreakerToResponse(breaker)
	}

	middleware.RespondSuccess(w, responses, "Circuit breaker history retrieved successfully")
}

// ListAlerts retrieves recent reputation alerts
// GET /api/v1/reputation/alerts
func (h *ReputationHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	// Get all domains with recent events
	scores, err := h.scoresRepo.ListAllScores(r.Context())
	if err != nil {
		h.logger.Error("Failed to list scores for alerts", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve alerts")
		return
	}

	// Generate alerts based on reputation scores and circuit breakers
	alerts := make([]*AlertResponse, 0)

	for _, score := range scores {
		// Alert for active circuit breakers
		if score.CircuitBreakerActive {
			alerts = append(alerts, &AlertResponse{
				Timestamp: score.LastUpdated,
				Domain:    score.Domain,
				Type:      "circuit_breaker",
				Message:   score.CircuitBreakerReason,
				Severity:  "critical",
			})
		}

		// Alert for poor reputation
		if score.ReputationScore < 50 {
			alerts = append(alerts, &AlertResponse{
				Timestamp: score.LastUpdated,
				Domain:    score.Domain,
				Type:      "low_reputation",
				Message:   "Reputation score below 50",
				Severity:  "high",
			})
		}

		// Alert for high complaint rate
		if score.ComplaintRate > 0.1 {
			alerts = append(alerts, &AlertResponse{
				Timestamp: score.LastUpdated,
				Domain:    score.Domain,
				Type:      "high_complaint_rate",
				Message:   "Complaint rate exceeds 0.1%",
				Severity:  "critical",
			})
		}

		// Alert for high bounce rate
		if score.BounceRate > 10.0 {
			alerts = append(alerts, &AlertResponse{
				Timestamp: score.LastUpdated,
				Domain:    score.Domain,
				Type:      "high_bounce_rate",
				Message:   "Bounce rate exceeds 10%",
				Severity:  "high",
			})
		}
	}

	middleware.RespondSuccess(w, alerts, "Alerts retrieved successfully")
}

// Helper functions

func auditToResponse(audit *domain.AuditResult) *AuditResponse {
	return &AuditResponse{
		Domain:       audit.Domain,
		Timestamp:    audit.Timestamp,
		SPF:          checkStatusToResponse(audit.SPFStatus),
		DKIM:         checkStatusToResponse(audit.DKIMStatus),
		DMARC:        checkStatusToResponse(audit.DMARCStatus),
		RDNS:         checkStatusToResponse(audit.RDNSStatus),
		FCrDNS:       checkStatusToResponse(audit.FCrDNSStatus),
		TLS:          checkStatusToResponse(audit.TLSStatus),
		MTASTS:       checkStatusToResponse(audit.MTASTSStatus),
		PostmasterOK: audit.PostmasterOK,
		AbuseOK:      audit.AbuseOK,
		OverallScore: audit.OverallScore,
		Issues:       audit.Issues,
	}
}

func checkStatusToResponse(status domain.CheckStatus) CheckStatusResponse {
	return CheckStatusResponse{
		Passed:  status.Passed,
		Message: status.Message,
		Details: status.Details,
	}
}

func scoreToResponse(score *domain.ReputationScore) *ScoreResponse {
	return &ScoreResponse{
		Domain:               score.Domain,
		ReputationScore:      score.ReputationScore,
		ComplaintRate:        score.ComplaintRate,
		BounceRate:           score.BounceRate,
		DeliveryRate:         score.DeliveryRate,
		CircuitBreakerActive: score.CircuitBreakerActive,
		CircuitBreakerReason: score.CircuitBreakerReason,
		WarmUpActive:         score.WarmUpActive,
		WarmUpDay:            score.WarmUpDay,
		LastUpdated:          score.LastUpdated,
	}
}

func circuitBreakerToResponse(breaker *domain.CircuitBreakerEvent) *CircuitBreakerResponse {
	return &CircuitBreakerResponse{
		ID:           breaker.ID,
		Domain:       breaker.Domain,
		TriggerType:  string(breaker.TriggerType),
		TriggerValue: breaker.TriggerValue,
		Threshold:    breaker.Threshold,
		PausedAt:     breaker.PausedAt,
		ResumedAt:    breaker.ResumedAt,
		AutoResumed:  breaker.AutoResumed,
		AdminNotes:   breaker.AdminNotes,
	}
}

// Helper to parse int with bounds checking
func parseInt(s string, min, max int) (int, error) {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err != nil {
		return 0, err
	}
	if val < min {
		val = min
	}
	if val > max {
		val = max
	}
	return val, nil
}

// Helper to get current Unix timestamp
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
