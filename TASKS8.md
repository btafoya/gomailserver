# TASKS8.md - Phase 8: Webhooks (Week 26)

## Overview

HTTP webhook system for real-time event notifications to external services.

**Total Tasks**: 9
**Priority**: [FULL] - Post-MVP feature
**Dependencies**: Phase 0-3

---

## 8.1 Webhook System

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| WH-001 | Webhook table schema | [ ] | F-022 | FULL |
| WH-002 | Webhook registration API | [ ] | API-002, WH-001 | FULL |
| WH-003 | Email received event | [ ] | WH-002, S-003 | FULL |
| WH-004 | Email sent event | [ ] | WH-002, Q-002 | FULL |
| WH-005 | Delivery status events | [ ] | WH-002, Q-004 | FULL |
| WH-006 | Security events | [ ] | WH-002, AU-002 | FULL |
| WH-007 | Quota warning events | [ ] | WH-002, U-006 | FULL |
| WH-008 | Retry logic with exponential backoff | [ ] | WH-002 | FULL |
| WH-009 | Webhook testing UI | [ ] | AUI-001, WH-002 | FULL |

---

## Task Details

### WH-001: Webhook Table Schema
**File**: `internal/database/migrations/011_webhooks.sql`
```sql
-- Webhook configurations
CREATE TABLE IF NOT EXISTS webhooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    secret TEXT NOT NULL,
    events TEXT NOT NULL, -- JSON array of event types
    enabled INTEGER NOT NULL DEFAULT 1,
    headers TEXT, -- JSON object of custom headers
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER REFERENCES users(id),
    CHECK (domain_id IS NOT NULL OR user_id IS NOT NULL)
);

-- Webhook delivery attempts
CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    webhook_id INTEGER NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    event_id TEXT NOT NULL UNIQUE, -- UUID for idempotency
    payload TEXT NOT NULL, -- JSON payload
    attempt_count INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, success, failed, retrying
    response_code INTEGER,
    response_body TEXT,
    error_message TEXT,
    next_retry_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    delivered_at DATETIME,
    last_attempt_at DATETIME
);

-- Indexes
CREATE INDEX idx_webhooks_domain ON webhooks(domain_id);
CREATE INDEX idx_webhooks_user ON webhooks(user_id);
CREATE INDEX idx_webhooks_enabled ON webhooks(enabled);
CREATE INDEX idx_webhook_deliveries_webhook ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_next_retry ON webhook_deliveries(next_retry_at);
CREATE INDEX idx_webhook_deliveries_event ON webhook_deliveries(event_type, event_id);
```

**Acceptance Criteria**:
- [ ] Webhooks can be scoped to domain or user
- [ ] Support multiple event types per webhook
- [ ] Track delivery history and retries
- [ ] Unique event IDs for idempotency

---

### WH-002: Webhook Registration API
**File**: `internal/api/webhook_handler.go`
```go
package api

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "net/url"

    "github.com/labstack/echo/v4"
    "github.com/btafoya/gomailserver/internal/domain"
)

type WebhookHandler struct {
    webhookRepo domain.WebhookRepository
}

type WebhookRequest struct {
    Name     string   `json:"name" validate:"required,max=255"`
    URL      string   `json:"url" validate:"required,url,max=2048"`
    Events   []string `json:"events" validate:"required,min=1,dive,oneof=email.received email.sent email.bounced email.delivered security.failed_login security.ip_blocked quota.warning quota.exceeded"`
    Headers  map[string]string `json:"headers,omitempty"`
    DomainID *int64   `json:"domain_id,omitempty"`
    UserID   *int64   `json:"user_id,omitempty"`
    Enabled  bool     `json:"enabled"`
}

type WebhookResponse struct {
    ID        int64             `json:"id"`
    Name      string            `json:"name"`
    URL       string            `json:"url"`
    Secret    string            `json:"secret,omitempty"` // Only shown on create
    Events    []string          `json:"events"`
    Headers   map[string]string `json:"headers,omitempty"`
    DomainID  *int64            `json:"domain_id,omitempty"`
    UserID    *int64            `json:"user_id,omitempty"`
    Enabled   bool              `json:"enabled"`
    CreatedAt string            `json:"created_at"`
    UpdatedAt string            `json:"updated_at"`
}

// POST /api/admin/webhooks
func (h *WebhookHandler) Create(c echo.Context) error {
    var req WebhookRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request body")
    }

    if err := c.Validate(req); err != nil {
        return echo.NewHTTPError(400, err.Error())
    }

    // Validate URL
    u, err := url.Parse(req.URL)
    if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
        return echo.NewHTTPError(400, "Invalid webhook URL")
    }

    // Generate secret for HMAC signing
    secretBytes := make([]byte, 32)
    if _, err := rand.Read(secretBytes); err != nil {
        return echo.NewHTTPError(500, "Failed to generate secret")
    }
    secret := hex.EncodeToString(secretBytes)

    eventsJSON, _ := json.Marshal(req.Events)
    headersJSON, _ := json.Marshal(req.Headers)

    webhook := &domain.Webhook{
        Name:     req.Name,
        URL:      req.URL,
        Secret:   secret,
        Events:   string(eventsJSON),
        Headers:  string(headersJSON),
        DomainID: req.DomainID,
        UserID:   req.UserID,
        Enabled:  req.Enabled,
    }

    if err := h.webhookRepo.Create(webhook); err != nil {
        return echo.NewHTTPError(500, "Failed to create webhook")
    }

    return c.JSON(201, WebhookResponse{
        ID:        webhook.ID,
        Name:      webhook.Name,
        URL:       webhook.URL,
        Secret:    secret, // Only shown once on create
        Events:    req.Events,
        Headers:   req.Headers,
        DomainID:  webhook.DomainID,
        UserID:    webhook.UserID,
        Enabled:   webhook.Enabled,
        CreatedAt: webhook.CreatedAt.Format(time.RFC3339),
    })
}

// GET /api/admin/webhooks
func (h *WebhookHandler) List(c echo.Context) error {
    domainID, _ := strconv.ParseInt(c.QueryParam("domain_id"), 10, 64)
    userID, _ := strconv.ParseInt(c.QueryParam("user_id"), 10, 64)

    webhooks, err := h.webhookRepo.List(domainID, userID)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to list webhooks")
    }

    response := make([]WebhookResponse, len(webhooks))
    for i, wh := range webhooks {
        var events []string
        var headers map[string]string
        json.Unmarshal([]byte(wh.Events), &events)
        json.Unmarshal([]byte(wh.Headers), &headers)

        response[i] = WebhookResponse{
            ID:        wh.ID,
            Name:      wh.Name,
            URL:       wh.URL,
            Events:    events,
            Headers:   headers,
            DomainID:  wh.DomainID,
            UserID:    wh.UserID,
            Enabled:   wh.Enabled,
            CreatedAt: wh.CreatedAt.Format(time.RFC3339),
            UpdatedAt: wh.UpdatedAt.Format(time.RFC3339),
        }
    }

    return c.JSON(200, response)
}

// PUT /api/admin/webhooks/:id
func (h *WebhookHandler) Update(c echo.Context) error {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

    var req WebhookRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request body")
    }

    webhook, err := h.webhookRepo.GetByID(id)
    if err != nil {
        return echo.NewHTTPError(404, "Webhook not found")
    }

    eventsJSON, _ := json.Marshal(req.Events)
    headersJSON, _ := json.Marshal(req.Headers)

    webhook.Name = req.Name
    webhook.URL = req.URL
    webhook.Events = string(eventsJSON)
    webhook.Headers = string(headersJSON)
    webhook.Enabled = req.Enabled

    if err := h.webhookRepo.Update(webhook); err != nil {
        return echo.NewHTTPError(500, "Failed to update webhook")
    }

    return c.JSON(200, map[string]string{"status": "updated"})
}

// DELETE /api/admin/webhooks/:id
func (h *WebhookHandler) Delete(c echo.Context) error {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

    if err := h.webhookRepo.Delete(id); err != nil {
        return echo.NewHTTPError(500, "Failed to delete webhook")
    }

    return c.NoContent(204)
}

// POST /api/admin/webhooks/:id/regenerate-secret
func (h *WebhookHandler) RegenerateSecret(c echo.Context) error {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

    secretBytes := make([]byte, 32)
    if _, err := rand.Read(secretBytes); err != nil {
        return echo.NewHTTPError(500, "Failed to generate secret")
    }
    newSecret := hex.EncodeToString(secretBytes)

    if err := h.webhookRepo.UpdateSecret(id, newSecret); err != nil {
        return echo.NewHTTPError(500, "Failed to update secret")
    }

    return c.JSON(200, map[string]string{"secret": newSecret})
}
```

**Acceptance Criteria**:
- [ ] CRUD operations for webhooks
- [ ] Automatic HMAC secret generation
- [ ] Event type validation
- [ ] Secret regeneration endpoint

---

### WH-003: Email Received Event
**File**: `internal/webhook/events.go`
```go
package webhook

import (
    "time"
)

// Event types
const (
    EventEmailReceived    = "email.received"
    EventEmailSent        = "email.sent"
    EventEmailBounced     = "email.bounced"
    EventEmailDelivered   = "email.delivered"
    EventSecurityLogin    = "security.failed_login"
    EventSecurityIPBlocked = "security.ip_blocked"
    EventQuotaWarning     = "quota.warning"
    EventQuotaExceeded    = "quota.exceeded"
)

// Base webhook payload
type WebhookPayload struct {
    ID        string    `json:"id"`        // Unique event ID
    Event     string    `json:"event"`     // Event type
    Timestamp time.Time `json:"timestamp"`
    Data      any       `json:"data"`
}

// Email received event data
type EmailReceivedData struct {
    MessageID   string   `json:"message_id"`
    From        string   `json:"from"`
    To          []string `json:"to"`
    Cc          []string `json:"cc,omitempty"`
    Subject     string   `json:"subject"`
    Date        string   `json:"date"`
    Size        int64    `json:"size"`
    HasAttachments bool  `json:"has_attachments"`
    SpamScore   float64  `json:"spam_score,omitempty"`
    VirusStatus string   `json:"virus_status,omitempty"` // clean, infected, error
    Headers     map[string]string `json:"headers,omitempty"`
}

// Create email received event
func NewEmailReceivedEvent(msg *domain.Message) *WebhookPayload {
    return &WebhookPayload{
        ID:        uuid.NewString(),
        Event:     EventEmailReceived,
        Timestamp: time.Now().UTC(),
        Data: EmailReceivedData{
            MessageID:      msg.MessageID,
            From:           msg.From,
            To:             msg.To,
            Cc:             msg.Cc,
            Subject:        msg.Subject,
            Date:           msg.Date.Format(time.RFC3339),
            Size:           msg.Size,
            HasAttachments: msg.HasAttachments,
            SpamScore:      msg.SpamScore,
            VirusStatus:    msg.VirusStatus,
        },
    }
}
```

**File**: `internal/smtp/hooks.go`
```go
package smtp

import (
    "github.com/btafoya/gomailserver/internal/webhook"
)

// Called after message is successfully stored
func (s *Server) onMessageReceived(msg *domain.Message) {
    event := webhook.NewEmailReceivedEvent(msg)

    // Get webhooks for this recipient's domain
    for _, recipient := range msg.Recipients {
        domain := extractDomain(recipient)
        webhooks, _ := s.webhookService.GetWebhooksForDomain(domain, webhook.EventEmailReceived)

        for _, wh := range webhooks {
            s.webhookService.QueueDelivery(wh.ID, event)
        }
    }
}
```

**Acceptance Criteria**:
- [ ] Fires on successful message receipt
- [ ] Includes message metadata
- [ ] Supports domain and user scoped webhooks
- [ ] Spam/virus status included if available

---

### WH-004: Email Sent Event
```go
// Email sent event data
type EmailSentData struct {
    MessageID   string   `json:"message_id"`
    From        string   `json:"from"`
    To          []string `json:"to"`
    Subject     string   `json:"subject"`
    Size        int64    `json:"size"`
    QueueID     string   `json:"queue_id"`
}

func NewEmailSentEvent(msg *domain.QueuedMessage) *WebhookPayload {
    return &WebhookPayload{
        ID:        uuid.NewString(),
        Event:     EventEmailSent,
        Timestamp: time.Now().UTC(),
        Data: EmailSentData{
            MessageID: msg.MessageID,
            From:      msg.From,
            To:        msg.To,
            Subject:   msg.Subject,
            Size:      msg.Size,
            QueueID:   msg.QueueID,
        },
    }
}
```

---

### WH-005: Delivery Status Events
```go
// Delivery status event data
type DeliveryStatusData struct {
    MessageID   string `json:"message_id"`
    QueueID     string `json:"queue_id"`
    Recipient   string `json:"recipient"`
    Status      string `json:"status"` // delivered, bounced, deferred
    RemoteMTA   string `json:"remote_mta,omitempty"`
    ResponseCode int   `json:"response_code,omitempty"`
    ResponseText string `json:"response_text,omitempty"`
    BounceType  string `json:"bounce_type,omitempty"` // hard, soft
    RetryCount  int    `json:"retry_count,omitempty"`
}

func NewDeliveryStatusEvent(delivery *domain.DeliveryAttempt, status string) *WebhookPayload {
    eventType := EventEmailDelivered
    if status == "bounced" {
        eventType = EventEmailBounced
    }

    return &WebhookPayload{
        ID:        uuid.NewString(),
        Event:     eventType,
        Timestamp: time.Now().UTC(),
        Data: DeliveryStatusData{
            MessageID:    delivery.MessageID,
            QueueID:      delivery.QueueID,
            Recipient:    delivery.Recipient,
            Status:       status,
            RemoteMTA:    delivery.RemoteMTA,
            ResponseCode: delivery.ResponseCode,
            ResponseText: delivery.ResponseText,
            BounceType:   delivery.BounceType,
            RetryCount:   delivery.RetryCount,
        },
    }
}
```

**Acceptance Criteria**:
- [ ] Fires on successful delivery
- [ ] Fires on bounce (hard and soft)
- [ ] Includes remote MTA response
- [ ] Tracks retry attempts

---

### WH-006: Security Events
```go
// Security event data
type SecurityEventData struct {
    EventSubType string `json:"event_subtype"` // failed_login, ip_blocked, suspicious_activity
    Username     string `json:"username,omitempty"`
    Email        string `json:"email,omitempty"`
    IPAddress    string `json:"ip_address"`
    UserAgent    string `json:"user_agent,omitempty"`
    Service      string `json:"service"` // smtp, imap, webmail, api
    Reason       string `json:"reason,omitempty"`
    AttemptCount int    `json:"attempt_count,omitempty"`
    BlockDuration string `json:"block_duration,omitempty"`
}

func NewFailedLoginEvent(username, ip, userAgent, service, reason string, attemptCount int) *WebhookPayload {
    return &WebhookPayload{
        ID:        uuid.NewString(),
        Event:     EventSecurityLogin,
        Timestamp: time.Now().UTC(),
        Data: SecurityEventData{
            EventSubType: "failed_login",
            Username:     username,
            IPAddress:    ip,
            UserAgent:    userAgent,
            Service:      service,
            Reason:       reason,
            AttemptCount: attemptCount,
        },
    }
}

func NewIPBlockedEvent(ip, service, reason string, duration time.Duration) *WebhookPayload {
    return &WebhookPayload{
        ID:        uuid.NewString(),
        Event:     EventSecurityIPBlocked,
        Timestamp: time.Now().UTC(),
        Data: SecurityEventData{
            EventSubType:  "ip_blocked",
            IPAddress:     ip,
            Service:       service,
            Reason:        reason,
            BlockDuration: duration.String(),
        },
    }
}
```

**Acceptance Criteria**:
- [ ] Failed login attempts
- [ ] IP blocking events
- [ ] Includes service context (SMTP, IMAP, etc.)
- [ ] Attempt counts for brute force detection

---

### WH-007: Quota Warning Events
```go
// Quota event data
type QuotaEventData struct {
    UserID       int64  `json:"user_id"`
    Email        string `json:"email"`
    UsedBytes    int64  `json:"used_bytes"`
    LimitBytes   int64  `json:"limit_bytes"`
    UsagePercent int    `json:"usage_percent"`
    Threshold    int    `json:"threshold"` // 80, 90, 100
}

func NewQuotaWarningEvent(user *domain.User, usedBytes, limitBytes int64) *WebhookPayload {
    percent := int((float64(usedBytes) / float64(limitBytes)) * 100)

    eventType := EventQuotaWarning
    threshold := 80

    if percent >= 100 {
        eventType = EventQuotaExceeded
        threshold = 100
    } else if percent >= 90 {
        threshold = 90
    }

    return &WebhookPayload{
        ID:        uuid.NewString(),
        Event:     eventType,
        Timestamp: time.Now().UTC(),
        Data: QuotaEventData{
            UserID:       user.ID,
            Email:        user.Email,
            UsedBytes:    usedBytes,
            LimitBytes:   limitBytes,
            UsagePercent: percent,
            Threshold:    threshold,
        },
    }
}
```

**Acceptance Criteria**:
- [ ] Warning at 80% and 90%
- [ ] Exceeded at 100%
- [ ] Includes usage details
- [ ] Per-user tracking

---

### WH-008: Retry Logic with Exponential Backoff
**File**: `internal/webhook/delivery.go`
```go
package webhook

import (
    "bytes"
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "net"
    "net/http"
    "sync"
    "time"

    "github.com/btafoya/gomailserver/internal/domain"
    "golang.org/x/sync/semaphore"
)

const (
    MaxRetries        = 5
    InitialDelay      = 30 * time.Second
    MaxDelay          = 4 * time.Hour
    RequestTimeout    = 30 * time.Second
    ConnectionTimeout = 5 * time.Second  // Production readiness: strict connection timeout
    MaxConcurrent     = 10               // Bulkhead: limit concurrent deliveries
    CircuitThreshold  = 5                // Open circuit after 5 consecutive failures
    CircuitTimeout    = 1 * time.Minute  // Circuit breaker reset timeout
)

type DeliveryService struct {
    webhookRepo      domain.WebhookRepository
    deliveryRepo     domain.WebhookDeliveryRepository
    httpClient       *http.Client
    logger           *zap.Logger
    concurrencySem   *semaphore.Weighted           // Bulkhead pattern
    circuitBreakers  map[int64]*CircuitBreaker     // Per-webhook circuit breaker
    circuitMu        sync.RWMutex
}

// CircuitBreaker tracks webhook health and prevents overwhelming failed endpoints
type CircuitBreaker struct {
    failures      int
    lastFailureAt time.Time
    state         string // "closed", "open", "half-open"
    mu            sync.RWMutex
}

func (cb *CircuitBreaker) RecordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    cb.failures = 0
    cb.state = "closed"
}

func (cb *CircuitBreaker) RecordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    cb.failures++
    cb.lastFailureAt = time.Now()
    if cb.failures >= CircuitThreshold {
        cb.state = "open"
    }
}

func (cb *CircuitBreaker) CanAttempt() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    if cb.state == "closed" {
        return true
    }

    if cb.state == "open" && time.Since(cb.lastFailureAt) > CircuitTimeout {
        cb.state = "half-open"
        return true
    }

    return cb.state == "half-open"
}

func NewDeliveryService(webhookRepo domain.WebhookRepository, deliveryRepo domain.WebhookDeliveryRepository, logger *zap.Logger) *DeliveryService {
    return &DeliveryService{
        webhookRepo:     webhookRepo,
        deliveryRepo:    deliveryRepo,
        concurrencySem:  semaphore.NewWeighted(MaxConcurrent),
        circuitBreakers: make(map[int64]*CircuitBreaker),
        httpClient: &http.Client{
            Timeout: RequestTimeout,
            Transport: &http.Transport{
                DialContext: (&net.Dialer{
                    Timeout:   ConnectionTimeout,
                    KeepAlive: 30 * time.Second,
                }).DialContext,
                MaxIdleConns:          100,
                MaxIdleConnsPerHost:   10,
                IdleConnTimeout:       90 * time.Second,
                TLSHandshakeTimeout:   10 * time.Second,
                ResponseHeaderTimeout: 10 * time.Second,
            },
        },
        logger: logger,
    }
}

func (s *DeliveryService) getCircuitBreaker(webhookID int64) *CircuitBreaker {
    s.circuitMu.RLock()
    cb, exists := s.circuitBreakers[webhookID]
    s.circuitMu.RUnlock()

    if exists {
        return cb
    }

    s.circuitMu.Lock()
    defer s.circuitMu.Unlock()

    // Double-check after acquiring write lock
    if cb, exists := s.circuitBreakers[webhookID]; exists {
        return cb
    }

    cb = &CircuitBreaker{state: "closed"}
    s.circuitBreakers[webhookID] = cb
    return cb
}

// Queue a delivery
func (s *DeliveryService) QueueDelivery(webhookID int64, payload *WebhookPayload) error {
    payloadJSON, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    delivery := &domain.WebhookDelivery{
        WebhookID: webhookID,
        EventType: payload.Event,
        EventID:   payload.ID,
        Payload:   string(payloadJSON),
        Status:    "pending",
    }

    return s.deliveryRepo.Create(delivery)
}

// Process pending deliveries
func (s *DeliveryService) ProcessPending(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.processBatch(ctx)
        }
    }
}

func (s *DeliveryService) processBatch(ctx context.Context) {
    deliveries, err := s.deliveryRepo.GetPendingDeliveries(100)
    if err != nil {
        s.logger.Error("Failed to get pending deliveries", zap.Error(err))
        return
    }

    for _, delivery := range deliveries {
        select {
        case <-ctx.Done():
            return
        default:
            // Bulkhead: acquire semaphore before delivery
            if err := s.concurrencySem.Acquire(ctx, 1); err != nil {
                continue
            }
            go func(d *domain.WebhookDelivery) {
                defer s.concurrencySem.Release(1)
                s.deliver(ctx, d)
            }(delivery)
        }
    }
}

func (s *DeliveryService) deliver(ctx context.Context, delivery *domain.WebhookDelivery) {
    webhook, err := s.webhookRepo.GetByID(delivery.WebhookID)
    if err != nil || !webhook.Enabled {
        delivery.Status = "failed"
        delivery.ErrorMessage = "Webhook not found or disabled"
        s.deliveryRepo.Update(delivery)
        return
    }

    // Circuit breaker: check if webhook is healthy
    cb := s.getCircuitBreaker(delivery.WebhookID)
    if !cb.CanAttempt() {
        delivery.Status = "retrying"
        delivery.ErrorMessage = "Circuit breaker open - too many consecutive failures"
        delivery.NextRetryAt = time.Now().Add(CircuitTimeout)
        s.deliveryRepo.Update(delivery)
        return
    }

    // Prepare request
    req, err := http.NewRequestWithContext(ctx, "POST", webhook.URL, bytes.NewReader([]byte(delivery.Payload)))
    if err != nil {
        s.handleFailure(delivery, 0, "", err.Error())
        cb.RecordFailure()
        return
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "gomailserver-webhook/1.0")
    req.Header.Set("X-Webhook-Event", delivery.EventType)
    req.Header.Set("X-Webhook-ID", fmt.Sprintf("%d", delivery.WebhookID))
    req.Header.Set("X-Webhook-Delivery-ID", delivery.EventID)
    req.Header.Set("X-Webhook-Timestamp", time.Now().UTC().Format(time.RFC3339))

    // HMAC signature
    signature := s.computeSignature([]byte(delivery.Payload), webhook.Secret)
    req.Header.Set("X-Webhook-Signature", "sha256="+signature)

    // Custom headers
    var customHeaders map[string]string
    if webhook.Headers != "" {
        json.Unmarshal([]byte(webhook.Headers), &customHeaders)
        for k, v := range customHeaders {
            req.Header.Set(k, v)
        }
    }

    // Send request
    resp, err := s.httpClient.Do(req)
    if err != nil {
        s.handleFailure(delivery, 0, "", err.Error())
        cb.RecordFailure()
        return
    }
    defer resp.Body.Close()

    // Read response
    body, _ := io.ReadAll(io.LimitReader(resp.Body, 10*1024)) // Max 10KB

    delivery.ResponseCode = resp.StatusCode
    delivery.ResponseBody = string(body)
    delivery.LastAttemptAt = time.Now()
    delivery.AttemptCount++

    // Check success (2xx status codes)
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        delivery.Status = "success"
        delivery.DeliveredAt = time.Now()
        cb.RecordSuccess()
        s.logger.Info("Webhook delivered successfully",
            zap.Int64("webhook_id", delivery.WebhookID),
            zap.String("event_id", delivery.EventID))
    } else {
        s.handleFailure(delivery, resp.StatusCode, string(body), "")
        cb.RecordFailure()
    }

    s.deliveryRepo.Update(delivery)
}

func (s *DeliveryService) handleFailure(delivery *domain.WebhookDelivery, statusCode int, responseBody, errorMsg string) {
    delivery.AttemptCount++
    delivery.ResponseCode = statusCode
    delivery.ResponseBody = responseBody
    if errorMsg != "" {
        delivery.ErrorMessage = errorMsg
    }
    delivery.LastAttemptAt = time.Now()

    if delivery.AttemptCount >= MaxRetries {
        delivery.Status = "failed"
        s.logger.Error("Webhook delivery failed after max retries",
            zap.Int64("webhook_id", delivery.WebhookID),
            zap.String("event_id", delivery.EventID),
            zap.Int("attempts", delivery.AttemptCount))
    } else {
        delivery.Status = "retrying"
        // Exponential backoff: 30s, 1m, 2m, 4m, 8m (capped at MaxDelay)
        delay := InitialDelay * time.Duration(1<<(delivery.AttemptCount-1))
        if delay > MaxDelay {
            delay = MaxDelay
        }
        delivery.NextRetryAt = time.Now().Add(delay)

        s.logger.Warn("Webhook delivery failed, will retry",
            zap.Int64("webhook_id", delivery.WebhookID),
            zap.String("event_id", delivery.EventID),
            zap.Int("attempt", delivery.AttemptCount),
            zap.Duration("next_retry", delay))
    }

    s.deliveryRepo.Update(delivery)
}

func (s *DeliveryService) computeSignature(payload []byte, secret string) string {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(payload)
    return hex.EncodeToString(mac.Sum(nil))
}

// Verify signature for incoming webhook verification requests
func VerifySignature(payload []byte, signature, secret string) bool {
    expected := "sha256=" + computeHMAC(payload, secret)
    return hmac.Equal([]byte(expected), []byte(signature))
}

func computeHMAC(payload []byte, secret string) string {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(payload)
    return hex.EncodeToString(mac.Sum(nil))
}
```

**Acceptance Criteria**:
- [ ] Exponential backoff: 30s, 1m, 2m, 4m, 8m
- [ ] Maximum 5 retry attempts
- [ ] HMAC-SHA256 signature
- [ ] Delivery status tracking
- [ ] Response logging (truncated)

**Production Readiness** (Added from spec panel review):
- [ ] Connection timeout: 5s (prevents hanging on slow endpoints)
- [ ] Total request timeout: 30s (includes connection + response)
- [ ] Bulkhead pattern: max 10 concurrent webhook deliveries
- [ ] Circuit breaker: opens after 5 consecutive failures per webhook
- [ ] Circuit breaker recovery: 1 minute timeout before retry
- [ ] P99 delivery latency: < 5s for successful deliveries
- [ ] Success rate target: > 95% for healthy webhooks
- [ ] Concurrent delivery limit enforced via semaphore
- [ ] Per-webhook circuit breaker state tracking

---

### WH-009: Webhook Testing UI
**File**: `web/admin/components/WebhookTestDialog.vue`
```vue
<template>
  <Dialog :open="open" @close="$emit('close')">
    <DialogContent class="max-w-2xl">
      <DialogHeader>
        <DialogTitle>Test Webhook</DialogTitle>
        <DialogDescription>
          Send a test event to verify your webhook configuration.
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4 py-4">
        <!-- Event Type -->
        <div class="space-y-2">
          <Label>Event Type</Label>
          <Select v-model="eventType">
            <SelectTrigger>
              <SelectValue placeholder="Select event type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="email.received">Email Received</SelectItem>
              <SelectItem value="email.sent">Email Sent</SelectItem>
              <SelectItem value="email.bounced">Email Bounced</SelectItem>
              <SelectItem value="email.delivered">Email Delivered</SelectItem>
              <SelectItem value="security.failed_login">Failed Login</SelectItem>
              <SelectItem value="security.ip_blocked">IP Blocked</SelectItem>
              <SelectItem value="quota.warning">Quota Warning</SelectItem>
              <SelectItem value="quota.exceeded">Quota Exceeded</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <!-- Payload Preview -->
        <div class="space-y-2">
          <Label>Payload Preview</Label>
          <div class="bg-gray-100 dark:bg-gray-800 rounded-lg p-4">
            <pre class="text-sm overflow-x-auto"><code>{{ payloadPreview }}</code></pre>
          </div>
        </div>

        <!-- Response -->
        <div v-if="response" class="space-y-2">
          <Label>Response</Label>
          <div :class="[
            'rounded-lg p-4',
            response.success ? 'bg-green-50 dark:bg-green-900/20' : 'bg-red-50 dark:bg-red-900/20'
          ]">
            <div class="flex items-center gap-2 mb-2">
              <Badge :variant="response.success ? 'success' : 'destructive'">
                {{ response.status_code || 'Error' }}
              </Badge>
              <span class="text-sm">{{ response.duration }}ms</span>
            </div>
            <pre v-if="response.body" class="text-sm overflow-x-auto">{{ response.body }}</pre>
            <p v-if="response.error" class="text-sm text-red-600">{{ response.error }}</p>
          </div>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="$emit('close')">Cancel</Button>
        <Button @click="sendTest" :disabled="sending">
          <Loader2Icon v-if="sending" class="w-4 h-4 mr-2 animate-spin" />
          {{ sending ? 'Sending...' : 'Send Test' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
const props = defineProps<{
  open: boolean
  webhook: Webhook
}>()

const emit = defineEmits(['close'])

const eventType = ref('email.received')
const sending = ref(false)
const response = ref<TestResponse | null>(null)

const samplePayloads = {
  'email.received': {
    message_id: '<test-12345@example.com>',
    from: 'sender@example.com',
    to: ['recipient@yourdomain.com'],
    subject: 'Test Email',
    date: new Date().toISOString(),
    size: 1024,
    has_attachments: false,
    spam_score: 0.1,
    virus_status: 'clean',
  },
  'email.sent': {
    message_id: '<sent-12345@yourdomain.com>',
    from: 'user@yourdomain.com',
    to: ['recipient@example.com'],
    subject: 'Outgoing Test',
    size: 512,
    queue_id: 'Q12345',
  },
  'email.bounced': {
    message_id: '<bounced-12345@yourdomain.com>',
    queue_id: 'Q12345',
    recipient: 'invalid@example.com',
    status: 'bounced',
    remote_mta: 'mx.example.com',
    response_code: 550,
    response_text: '5.1.1 User unknown',
    bounce_type: 'hard',
  },
  // ... more sample payloads
}

const payloadPreview = computed(() => {
  const payload = {
    id: crypto.randomUUID(),
    event: eventType.value,
    timestamp: new Date().toISOString(),
    data: samplePayloads[eventType.value] || {},
  }
  return JSON.stringify(payload, null, 2)
})

const sendTest = async () => {
  sending.value = true
  response.value = null

  try {
    const result = await $fetch(`/api/admin/webhooks/${props.webhook.id}/test`, {
      method: 'POST',
      body: { event_type: eventType.value }
    })
    response.value = result
  } catch (err) {
    response.value = {
      success: false,
      error: err.message,
    }
  } finally {
    sending.value = false
  }
}
</script>
```

**API Endpoint**:
```go
// POST /api/admin/webhooks/:id/test
func (h *WebhookHandler) Test(c echo.Context) error {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

    var req struct {
        EventType string `json:"event_type"`
    }
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    webhook, err := h.webhookRepo.GetByID(id)
    if err != nil {
        return echo.NewHTTPError(404, "Webhook not found")
    }

    // Create test payload
    payload := createTestPayload(req.EventType)

    // Deliver synchronously for testing
    start := time.Now()
    result, err := h.deliveryService.DeliverSync(webhook, payload)
    duration := time.Since(start).Milliseconds()

    if err != nil {
        return c.JSON(200, map[string]interface{}{
            "success":  false,
            "error":    err.Error(),
            "duration": duration,
        })
    }

    return c.JSON(200, map[string]interface{}{
        "success":     result.StatusCode >= 200 && result.StatusCode < 300,
        "status_code": result.StatusCode,
        "body":        result.Body,
        "duration":    duration,
    })
}
```

**Acceptance Criteria**:
- [ ] Event type selection
- [ ] Payload preview
- [ ] Synchronous delivery for immediate feedback
- [ ] Response display with timing
- [ ] Error handling and display

---

## Webhook Integration Points

| Event Source | Hook Location | Event Type |
|--------------|---------------|------------|
| SMTP Server | `onMessageReceived()` | `email.received` |
| Queue Processor | `onMessageQueued()` | `email.sent` |
| Delivery Service | `onDeliverySuccess()` | `email.delivered` |
| Delivery Service | `onDeliveryBounce()` | `email.bounced` |
| Auth Handler | `onLoginFailed()` | `security.failed_login` |
| Rate Limiter | `onIPBlocked()` | `security.ip_blocked` |
| Quota Service | `onQuotaThreshold()` | `quota.warning`, `quota.exceeded` |

---

## Security Considerations

1. **HMAC Signatures**: All webhooks signed with SHA-256
2. **Secret Rotation**: API to regenerate secrets
3. **URL Validation**: Only HTTP/HTTPS allowed
4. **Response Limits**: Body truncated to 10KB
5. **Timeout**: 30 second connection timeout
6. **No Redirects**: Don't follow redirects by default
7. **Rate Limiting**: Limit delivery attempts per second

---

## Testing Checklist

- [ ] Webhook CRUD operations
- [ ] HMAC signature validation
- [ ] Event firing for all event types
- [ ] Retry logic with exponential backoff
- [ ] Max retry limit (5 attempts)
- [ ] Test webhook UI functionality
- [ ] Delivery history tracking
- [ ] Secret regeneration
