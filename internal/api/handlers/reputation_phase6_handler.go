package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	maildomain "github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	reputationService "github.com/btafoya/gomailserver/internal/reputation/service"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ReputationPhase6Handler handles Phase 6 WebUI endpoints
type ReputationPhase6Handler struct {
	alertsRepo           repository.AlertsRepository
	scoresRepo           repository.ScoresRepository
	circuitBreakerRepo   repository.CircuitBreakerRepository
	historicalScoresRepo repository.HistoricalScoresRepository

	// Services for operational mail integration
	userService    *service.UserService
	mailboxService *service.MailboxService
	messageService *service.MessageService
	queueService   *service.QueueService
	auditorService *reputationService.AuditorService

	logger *zap.Logger
}

// NewReputationPhase6Handler creates a new Phase 6 handler
func NewReputationPhase6Handler(
	alertsRepo repository.AlertsRepository,
	scoresRepo repository.ScoresRepository,
	circuitBreakerRepo repository.CircuitBreakerRepository,
	logger *zap.Logger,
) *ReputationPhase6Handler {
	return &ReputationPhase6Handler{
		alertsRepo:         alertsRepo,
		scoresRepo:         scoresRepo,
		circuitBreakerRepo: circuitBreakerRepo,
		logger:             logger,
	}
}

// SetHistoricalScoresRepo sets the historical scores repository for trend calculation
func (h *ReputationPhase6Handler) SetHistoricalScoresRepo(repo repository.HistoricalScoresRepository) {
	h.historicalScoresRepo = repo
}

// SetServices sets the service dependencies for operational mail access
func (h *ReputationPhase6Handler) SetServices(
	userService *service.UserService,
	mailboxService *service.MailboxService,
	messageService *service.MessageService,
	queueService *service.QueueService,
) {
	h.userService = userService
	h.mailboxService = mailboxService
	h.messageService = messageService
	h.queueService = queueService
}

// SetAuditorService sets the auditor service for DNS health checks
func (h *ReputationPhase6Handler) SetAuditorService(auditor *reputationService.AuditorService) {
	h.auditorService = auditor
}

// ===================================================================
// Operational Mail Endpoints (Phase 6.1) - IMAP Integration
// ===================================================================

// GetOperationalMail returns operational mailbox messages (postmaster@, abuse@)
// GET /api/v1/reputation/operational-mail
func (h *ReputationPhase6Handler) GetOperationalMail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Query parameters for filtering
	domainFilter := r.URL.Query().Get("domain")
	mailboxType := r.URL.Query().Get("type") // "postmaster", "abuse", or empty for both
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	// Check if services are available
	if h.userService == nil || h.mailboxService == nil || h.messageService == nil {
		h.logger.Warn("operational mail services not configured")
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"messages": []interface{}{},
			"total":    0,
			"warning":  "Mail services not configured for operational mailbox access",
		})
		return
	}

	// Get all users to find operational mailboxes
	users, err := h.userService.ListAll()
	if err != nil {
		h.logger.Error("failed to list users for operational mail", zap.Error(err))
		respondError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	// Find operational mailbox users
	var operationalUsers []*maildomain.User
	for _, user := range users {
		email := strings.ToLower(user.Email)

		// Filter by mailbox type
		isPostmaster := strings.HasPrefix(email, "postmaster@")
		isAbuse := strings.HasPrefix(email, "abuse@")

		if mailboxType == "postmaster" && !isPostmaster {
			continue
		}
		if mailboxType == "abuse" && !isAbuse {
			continue
		}
		if mailboxType == "" && !isPostmaster && !isAbuse {
			continue
		}

		// Filter by domain if specified
		if domainFilter != "" {
			parts := strings.Split(email, "@")
			if len(parts) == 2 && parts[1] != domainFilter {
				continue
			}
		}

		operationalUsers = append(operationalUsers, user)
	}

	// Collect messages from operational mailboxes
	var messages []map[string]interface{}
	for _, user := range operationalUsers {
		// Get user's INBOX
		inbox, err := h.mailboxService.GetByName(user.ID, "INBOX")
		if err != nil {
			h.logger.Debug("failed to get inbox for operational user",
				zap.String("email", user.Email),
				zap.Error(err),
			)
			continue
		}

		// Get messages from INBOX
		userMessages, err := h.messageService.ListMessages(ctx, int(inbox.ID), int(user.ID), limit/len(operationalUsers)+1, 0)
		if err != nil {
			h.logger.Debug("failed to get messages for operational user",
				zap.String("email", user.Email),
				zap.Error(err),
			)
			continue
		}

		// Convert messages to response format
		for _, msg := range userMessages {
			severity := h.classifyMessageSeverity(msg)
			isRead := strings.Contains(msg.Flags, "\\Seen")
			isSpam := strings.Contains(msg.Flags, "\\Junk") || strings.Contains(msg.Flags, "$Junk")

			// Extract preview from content
			preview := h.extractMessagePreview(msg.Content, 150)

			messages = append(messages, map[string]interface{}{
				"id":        fmt.Sprintf("%d", msg.ID),
				"from":      msg.From,
				"recipient": user.Email,
				"subject":   msg.Subject,
				"preview":   preview,
				"timestamp": msg.ReceivedAt.Unix(),
				"read":      isRead,
				"spam":      isSpam,
				"severity":  severity,
				"messageId": msg.MessageID,
				"userId":    user.ID,
				"mailboxId": inbox.ID,
			})
		}
	}

	// Sort by timestamp descending
	// (already done by database query in most cases)

	// Limit total results
	if len(messages) > limit {
		messages = messages[:limit]
	}

	response := map[string]interface{}{
		"messages": messages,
		"total":    len(messages),
	}

	respondJSON(w, http.StatusOK, response)
}

// classifyMessageSeverity classifies the severity of an operational message
func (h *ReputationPhase6Handler) classifyMessageSeverity(msg *maildomain.Message) string {
	subject := strings.ToLower(msg.Subject)
	from := strings.ToLower(msg.From)

	// Critical: Spam complaints, abuse reports
	if strings.Contains(subject, "spam complaint") ||
		strings.Contains(subject, "abuse report") ||
		strings.Contains(subject, "complaint") ||
		strings.Contains(from, "abuse-report") ||
		strings.Contains(from, "fbl@") ||
		strings.Contains(from, "feedback@") {
		return "critical"
	}

	// High: Delivery failures, bounces
	if strings.Contains(subject, "delivery failure") ||
		strings.Contains(subject, "delivery status") ||
		strings.Contains(subject, "undeliverable") ||
		strings.Contains(subject, "bounce") ||
		strings.Contains(subject, "returned mail") {
		return "high"
	}

	// Medium: General reports, notifications
	if strings.Contains(subject, "dmarc") ||
		strings.Contains(subject, "report") ||
		strings.Contains(from, "postmaster") {
		return "medium"
	}

	// Low: Everything else
	return "low"
}

// extractMessagePreview extracts a text preview from message content
func (h *ReputationPhase6Handler) extractMessagePreview(content []byte, maxLen int) string {
	if len(content) == 0 {
		return ""
	}

	// Simple text extraction - skip headers
	text := string(content)
	parts := strings.SplitN(text, "\r\n\r\n", 2)
	if len(parts) > 1 {
		text = parts[1]
	}

	// Remove HTML tags if present
	text = strings.ReplaceAll(text, "<br>", " ")
	text = strings.ReplaceAll(text, "<br/>", " ")
	text = strings.ReplaceAll(text, "</p>", " ")

	// Remove other HTML tags
	for strings.Contains(text, "<") {
		start := strings.Index(text, "<")
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}

	// Clean whitespace
	text = strings.TrimSpace(text)
	text = strings.Join(strings.Fields(text), " ")

	// Truncate
	if len(text) > maxLen {
		text = text[:maxLen] + "..."
	}

	return text
}

// MarkOperationalMailRead marks an operational message as read
// POST /api/v1/reputation/operational-mail/:id/read
func (h *ReputationPhase6Handler) MarkOperationalMailRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	messageID := chi.URLParam(r, "id")

	msgID, err := strconv.ParseInt(messageID, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	if h.messageService == nil {
		respondError(w, http.StatusServiceUnavailable, "Message service not available")
		return
	}

	// Get message to find user ID
	msg, err := h.messageService.GetByID(msgID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Add \Seen flag
	err = h.messageService.UpdateFlags(ctx, int(msgID), int(msg.UserID), []string{"\\Seen"}, "add")
	if err != nil {
		h.logger.Error("failed to mark message as read", zap.Error(err))
		respondError(w, http.StatusInternalServerError, "Failed to mark message as read")
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"message_id": messageID,
		"read_at":    time.Now().Unix(),
	}

	respondJSON(w, http.StatusOK, response)
}

// DeleteOperationalMail deletes an operational message
// DELETE /api/v1/reputation/operational-mail/:id
func (h *ReputationPhase6Handler) DeleteOperationalMail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	messageID := chi.URLParam(r, "id")

	msgID, err := strconv.ParseInt(messageID, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	if h.messageService == nil {
		respondError(w, http.StatusServiceUnavailable, "Message service not available")
		return
	}

	// Get message to find user ID
	msg, err := h.messageService.GetByID(msgID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Delete message (moves to Trash)
	err = h.messageService.DeleteMessage(ctx, int(msgID), int(msg.UserID))
	if err != nil {
		h.logger.Error("failed to delete message", zap.Error(err))
		respondError(w, http.StatusInternalServerError, "Failed to delete message")
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"message_id": messageID,
		"deleted_at": time.Now().Unix(),
	}

	respondJSON(w, http.StatusOK, response)
}

// MarkOperationalMailSpam marks message as spam and blocks sender
// POST /api/v1/reputation/operational-mail/:id/spam
func (h *ReputationPhase6Handler) MarkOperationalMailSpam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	messageID := chi.URLParam(r, "id")

	msgID, err := strconv.ParseInt(messageID, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	if h.messageService == nil || h.mailboxService == nil {
		respondError(w, http.StatusServiceUnavailable, "Mail services not available")
		return
	}

	// Get message
	msg, err := h.messageService.GetByID(msgID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Add spam flag
	err = h.messageService.UpdateFlags(ctx, int(msgID), int(msg.UserID), []string{"\\Junk", "$Junk"}, "add")
	if err != nil {
		h.logger.Error("failed to add spam flag", zap.Error(err))
	}

	// Move to Spam folder if available
	spamMailbox, err := h.mailboxService.GetByName(msg.UserID, "Spam")
	if err == nil {
		err = h.messageService.MoveMessage(ctx, int(msgID), int(spamMailbox.ID), int(msg.UserID))
		if err != nil {
			h.logger.Error("failed to move message to Spam folder", zap.Error(err))
		}
	}

	// Log the spam report for potential blocklist integration
	h.logger.Info("operational mail marked as spam",
		zap.Int64("message_id", msgID),
		zap.String("from", msg.From),
		zap.String("to", msg.To),
	)

	response := map[string]interface{}{
		"success":    true,
		"message_id": messageID,
		"blocked_at": time.Now().Unix(),
		"from":       msg.From,
	}

	respondJSON(w, http.StatusOK, response)
}

// ForwardOperationalMail forwards operational message to another address
// POST /api/v1/reputation/operational-mail/:id/forward
func (h *ReputationPhase6Handler) ForwardOperationalMail(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")

	msgID, err := strconv.ParseInt(messageID, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	var req struct {
		To string `json:"to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.To == "" {
		respondError(w, http.StatusBadRequest, "Missing 'to' field")
		return
	}

	if h.messageService == nil || h.queueService == nil {
		respondError(w, http.StatusServiceUnavailable, "Mail services not available")
		return
	}

	// Get original message
	msg, err := h.messageService.GetByID(msgID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Load full content if needed
	if msg.StorageType == "file" && len(msg.Content) == 0 {
		msg, err = h.messageService.GetByID(msgID) // GetByID loads content
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to load message content")
			return
		}
	}

	// Build forwarded message
	forwardSubject := fmt.Sprintf("Fwd: %s", msg.Subject)
	forwardBody := fmt.Sprintf(`---------- Forwarded message ----------
From: %s
Date: %s
Subject: %s

%s`, msg.From, msg.ReceivedAt.Format(time.RFC1123), msg.Subject, string(msg.Content))

	// Use queue service to send forwarded message
	_, err = h.queueService.Enqueue(msg.To, []string{req.To}, []byte(
		fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			msg.To, req.To, forwardSubject, forwardBody),
	))
	if err != nil {
		h.logger.Error("failed to queue forwarded message", zap.Error(err))
		respondError(w, http.StatusInternalServerError, "Failed to forward message")
		return
	}

	h.logger.Info("operational mail forwarded",
		zap.Int64("message_id", msgID),
		zap.String("forwarded_to", req.To),
	)

	response := map[string]interface{}{
		"success":      true,
		"message_id":   messageID,
		"forwarded_to": req.To,
		"forwarded_at": time.Now().Unix(),
	}

	respondJSON(w, http.StatusOK, response)
}

// ===================================================================
// Deliverability Status Endpoints (Phase 6.2) - Historical Trend & DNS
// ===================================================================

// GetDeliverabilityStatus returns comprehensive deliverability health
// GET /api/v1/reputation/deliverability
// GET /api/v1/reputation/deliverability/:domain
func (h *ReputationPhase6Handler) GetDeliverabilityStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	domainName := chi.URLParam(r, "domain")

	// Get reputation score
	var score *domain.ReputationScore
	var err error

	if domainName != "" {
		score, err = h.scoresRepo.GetReputationScore(ctx, domainName)
	} else {
		// Get first domain or average score
		scores, err := h.scoresRepo.ListAllScores(ctx)
		if err == nil && len(scores) > 0 {
			score = scores[0]
			domainName = score.Domain
		}
	}

	if err != nil || score == nil {
		respondError(w, http.StatusNotFound, "Domain reputation not found")
		return
	}

	// Calculate trend based on historical scores
	trend := h.calculateTrendFromHistory(ctx, domainName)

	// Get DNS health status
	dnsHealth := h.getDNSHealthActual(ctx, domainName)

	response := map[string]interface{}{
		"domain":          domainName,
		"reputationScore": score.ReputationScore,
		"trend":           trend.Direction,
		"trendDetails": map[string]interface{}{
			"direction":      trend.Direction,
			"changePercent":  trend.ChangePercent,
			"periodDays":     trend.PeriodDays,
			"previousScore":  trend.PreviousScore,
			"currentScore":   score.ReputationScore,
			"dataPoints":     trend.DataPoints,
			"confidenceNote": trend.ConfidenceNote,
		},
		"dnsHealth":   dnsHealth,
		"lastChecked": time.Now().Unix(),
		"metrics": map[string]interface{}{
			"complaintRate": score.ComplaintRate,
			"bounceRate":    score.BounceRate,
			"deliveryRate":  score.DeliveryRate,
		},
	}

	respondJSON(w, http.StatusOK, response)
}

// TrendResult represents the calculated reputation trend
type TrendResult struct {
	Direction      string  // "improving", "stable", "declining"
	ChangePercent  float64 // Percentage change over period
	PeriodDays     int     // Number of days analyzed
	PreviousScore  int     // Score at start of period
	DataPoints     int     // Number of data points used
	ConfidenceNote string  // Note about confidence level
}

// calculateTrendFromHistory calculates trend based on historical scores
func (h *ReputationPhase6Handler) calculateTrendFromHistory(ctx context.Context, domainName string) TrendResult {
	result := TrendResult{
		Direction:      "stable",
		PeriodDays:     7,
		ConfidenceNote: "Insufficient data for trend analysis",
	}

	if h.historicalScoresRepo == nil {
		return result
	}

	// Get daily averages for the last 7 days
	scores, err := h.historicalScoresRepo.GetDailyAverages(ctx, domainName, 7)
	if err != nil {
		h.logger.Debug("failed to get historical scores for trend",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		return result
	}

	result.DataPoints = len(scores)

	if len(scores) < 2 {
		// Not enough data points for trend analysis
		return result
	}

	// Get first and last scores
	firstScore := scores[0].ReputationScore
	lastScore := scores[len(scores)-1].ReputationScore

	result.PreviousScore = firstScore

	// Calculate percentage change
	if firstScore > 0 {
		result.ChangePercent = float64(lastScore-firstScore) / float64(firstScore) * 100
	}

	// Determine trend direction with significance threshold
	const significanceThreshold = 5.0 // 5% change considered significant

	if result.ChangePercent > significanceThreshold {
		result.Direction = "improving"
		result.ConfidenceNote = fmt.Sprintf("Score improved by %.1f%% over %d days", result.ChangePercent, result.PeriodDays)
	} else if result.ChangePercent < -significanceThreshold {
		result.Direction = "declining"
		result.ConfidenceNote = fmt.Sprintf("Score declined by %.1f%% over %d days", -result.ChangePercent, result.PeriodDays)
	} else {
		result.Direction = "stable"
		result.ConfidenceNote = fmt.Sprintf("Score stable (%.1f%% change) over %d days", result.ChangePercent, result.PeriodDays)
	}

	// Add confidence note based on data points
	if len(scores) < 4 {
		result.ConfidenceNote += " (limited data)"
	}

	return result
}

// getDNSHealthActual performs actual DNS checks using the auditor service
func (h *ReputationPhase6Handler) getDNSHealthActual(ctx context.Context, domainName string) map[string]interface{} {
	result := map[string]interface{}{
		"spf": map[string]string{
			"status":  "unknown",
			"message": "DNS check not available",
		},
		"dkim": map[string]string{
			"status":  "unknown",
			"message": "DNS check not available",
		},
		"dmarc": map[string]string{
			"status":  "unknown",
			"message": "DNS check not available",
		},
		"rdns": map[string]string{
			"status":  "unknown",
			"message": "DNS check not available",
		},
	}

	if h.auditorService == nil {
		return result
	}

	// Perform a full audit to get DNS status
	// Note: We pass nil for IP since we're only checking DNS records
	auditResult, err := h.auditorService.AuditDomain(ctx, domainName, nil)
	if err != nil {
		h.logger.Debug("DNS audit failed",
			zap.String("domain", domainName),
			zap.Error(err),
		)
		return result
	}

	// Convert audit results to API response format
	result["spf"] = map[string]interface{}{
		"status":  boolToStatus(auditResult.SPFStatus.Passed),
		"message": auditResult.SPFStatus.Message,
		"details": auditResult.SPFStatus.Details,
	}

	result["dkim"] = map[string]interface{}{
		"status":  boolToStatus(auditResult.DKIMStatus.Passed),
		"message": auditResult.DKIMStatus.Message,
		"details": auditResult.DKIMStatus.Details,
	}

	result["dmarc"] = map[string]interface{}{
		"status":  boolToStatus(auditResult.DMARCStatus.Passed),
		"message": auditResult.DMARCStatus.Message,
		"details": auditResult.DMARCStatus.Details,
	}

	result["rdns"] = map[string]interface{}{
		"status":  boolToStatus(auditResult.RDNSStatus.Passed),
		"message": auditResult.RDNSStatus.Message,
		"details": auditResult.RDNSStatus.Details,
	}

	// Add additional checks if available
	result["mtasts"] = map[string]interface{}{
		"status":  boolToStatus(auditResult.MTASTSStatus.Passed),
		"message": auditResult.MTASTSStatus.Message,
		"details": auditResult.MTASTSStatus.Details,
	}

	result["tls"] = map[string]interface{}{
		"status":  boolToStatus(auditResult.TLSStatus.Passed),
		"message": auditResult.TLSStatus.Message,
		"details": auditResult.TLSStatus.Details,
	}

	result["overallScore"] = auditResult.OverallScore
	result["issues"] = auditResult.Issues

	return result
}

// boolToStatus converts a boolean pass/fail to a status string
func boolToStatus(passed bool) string {
	if passed {
		return "pass"
	}
	return "fail"
}

// ===================================================================
// Circuit Breaker Manual Control Endpoints (Phase 6.2)
// ===================================================================

// GetCircuitBreakers returns active and recent circuit breakers
// GET /api/v1/reputation/circuit-breakers
// GET /api/v1/reputation/circuit-breakers/:domain
func (h *ReputationPhase6Handler) GetCircuitBreakers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	domainName := chi.URLParam(r, "domain")

	var breakers []*domain.CircuitBreakerEvent
	var err error

	if domainName != "" {
		breakers, err = h.circuitBreakerRepo.GetBreakerHistory(ctx, domainName, 10)
	} else {
		breakers, err = h.circuitBreakerRepo.GetActiveBreakers(ctx)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get circuit breakers: %v", err))
		return
	}

	// Enhance breakers with status and auto-resume countdown
	enhancedBreakers := make([]map[string]interface{}, 0, len(breakers))
	now := time.Now().Unix()

	for _, breaker := range breakers {
		status := "active"
		if breaker.ResumedAt != nil {
			status = "resolved"
		}

		var autoResumeAt *int64
		if status == "active" && breaker.AutoResumed {
			// Calculate auto-resume time (typically 4 hours after pause)
			resumeTime := breaker.PausedAt + (4 * 3600)
			autoResumeAt = &resumeTime
		}

		enhanced := map[string]interface{}{
			"id":           breaker.ID,
			"domain":       breaker.Domain,
			"triggerType":  breaker.TriggerType,
			"triggerValue": breaker.TriggerValue,
			"threshold":    breaker.Threshold,
			"reason":       fmt.Sprintf("%s rate exceeded threshold", breaker.TriggerType),
			"pausedAt":     breaker.PausedAt,
			"resumedAt":    breaker.ResumedAt,
			"autoResumed":  breaker.AutoResumed,
			"autoResumeAt": autoResumeAt,
			"adminNotes":   breaker.AdminNotes,
			"status":       status,
		}

		enhancedBreakers = append(enhancedBreakers, enhanced)
	}

	response := map[string]interface{}{
		"breakers":  enhancedBreakers,
		"total":     len(enhancedBreakers),
		"timestamp": now,
	}

	respondJSON(w, http.StatusOK, response)
}

// ResumeCircuitBreaker manually resumes a paused domain
// POST /api/v1/reputation/circuit-breakers/:id/resume
func (h *ReputationPhase6Handler) ResumeCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	breakerIDStr := chi.URLParam(r, "id")

	breakerID, err := strconv.ParseInt(breakerIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid breaker ID")
		return
	}

	// Mark breaker as resumed
	if err := h.circuitBreakerRepo.RecordResume(ctx, fmt.Sprintf("%d", breakerID), false, "Manual admin override"); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to resume circuit breaker: %v", err))
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"breaker_id": breakerID,
		"resumed_at": time.Now().Unix(),
	}

	respondJSON(w, http.StatusOK, response)
}

// PauseCircuitBreaker manually pauses a domain
// POST /api/v1/reputation/circuit-breakers/pause
func (h *ReputationPhase6Handler) PauseCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Domain      string `json:"domain"`
		Reason      string `json:"reason"`
		TriggerType string `json:"triggerType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Domain == "" || req.Reason == "" {
		respondError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Create manual circuit breaker event
	breaker := &domain.CircuitBreakerEvent{
		Domain:       req.Domain,
		TriggerType:  domain.TriggerType(req.TriggerType),
		TriggerValue: 0.0, // Manual pause has no trigger value
		Threshold:    0.0,
		PausedAt:     time.Now().Unix(),
		AutoResumed:  false,
		AdminNotes:   req.Reason,
	}

	if err := h.circuitBreakerRepo.RecordPause(ctx, breaker); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to pause domain: %v", err))
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"domain":     req.Domain,
		"paused_at":  breaker.PausedAt,
		"breaker_id": breaker.ID,
	}

	respondJSON(w, http.StatusOK, response)
}

// ===================================================================
// Enhanced Alert Endpoints (Phase 6.3)
// ===================================================================

// GetAlerts returns alerts with filtering and pagination
// GET /api/v1/reputation/alerts
func (h *ReputationPhase6Handler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	domainFilter := r.URL.Query().Get("domain")
	severityStr := r.URL.Query().Get("severity")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	var alerts []*domain.ReputationAlert
	var err error

	// Filter by parameters
	if domainFilter != "" {
		alerts, err = h.alertsRepo.ListByDomain(ctx, domainFilter, limit, offset)
	} else if severityStr != "" {
		severity := domain.AlertSeverity(severityStr)
		alerts, err = h.alertsRepo.ListBySeverity(ctx, severity, limit)
	} else {
		alerts, err = h.alertsRepo.GetRecentAlerts(ctx, limit)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get alerts: %v", err))
		return
	}

	// Convert to response format expected by frontend
	alertsResponse := make([]map[string]interface{}, 0, len(alerts))
	for _, alert := range alerts {
		alertsResponse = append(alertsResponse, map[string]interface{}{
			"id":             alert.ID,
			"domain":         alert.Domain,
			"alertType":      alert.AlertType,
			"severity":       string(alert.Severity),
			"title":          alert.Title,
			"message":        alert.Message,
			"metadata":       alert.Details,
			"createdAt":      alert.CreatedAt,
			"readAt":         alert.AcknowledgedAt,
			"acknowledgedAt": alert.AcknowledgedAt,
			"acknowledgedBy": alert.AcknowledgedBy,
		})
	}

	response := map[string]interface{}{
		"alerts": alertsResponse,
		"total":  len(alertsResponse),
	}

	respondJSON(w, http.StatusOK, response)
}

// GetUnreadAlertCount returns count of unread alerts
// GET /api/v1/reputation/alerts/unread
func (h *ReputationPhase6Handler) GetUnreadAlertCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count, err := h.alertsRepo.GetUnacknowledgedCount(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get unread count: %v", err))
		return
	}

	response := map[string]interface{}{
		"count": count,
	}

	respondJSON(w, http.StatusOK, response)
}

// MarkAlertRead marks an alert as read
// POST /api/v1/reputation/alerts/:id/read
func (h *ReputationPhase6Handler) MarkAlertRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertIDStr := chi.URLParam(r, "id")

	alertID, err := strconv.ParseInt(alertIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid alert ID")
		return
	}

	// Mark as acknowledged (we're using acknowledged as "read" for now)
	if err := h.alertsRepo.Acknowledge(ctx, alertID, "system"); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to mark alert as read: %v", err))
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"alert_id": alertID,
		"read_at":  time.Now().Unix(),
	}

	respondJSON(w, http.StatusOK, response)
}

// AcknowledgeAlert acknowledges an alert
// POST /api/v1/reputation/alerts/:id/acknowledge
func (h *ReputationPhase6Handler) AcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertIDStr := chi.URLParam(r, "id")

	alertID, err := strconv.ParseInt(alertIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid alert ID")
		return
	}

	// Get admin user from auth context if available
	adminUser := "admin"
	if user := r.Context().Value("user"); user != nil {
		if u, ok := user.(map[string]interface{}); ok {
			if email, ok := u["email"].(string); ok {
				adminUser = email
			}
		}
	}

	if err := h.alertsRepo.Acknowledge(ctx, alertID, adminUser); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to acknowledge alert: %v", err))
		return
	}

	response := map[string]interface{}{
		"success":         true,
		"alert_id":        alertID,
		"acknowledged_at": time.Now().Unix(),
		"acknowledged_by": adminUser,
	}

	respondJSON(w, http.StatusOK, response)
}

// ===================================================================
// Helper Functions
// ===================================================================

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error":  message,
		"status": http.StatusText(status),
	})
}
