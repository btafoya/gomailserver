package domain

import (
	"time"
)

// Webhook represents a webhook subscription
type Webhook struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Secret      string    `json:"secret"` // HMAC secret for signature validation
	EventTypes  string    `json:"event_types"` // Comma-separated event types
	Active      bool      `json:"active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID               int64     `json:"id"`
	WebhookID        int64     `json:"webhook_id"`
	EventType        string    `json:"event_type"`
	Payload          string    `json:"payload"` // JSON payload
	AttemptCount     int       `json:"attempt_count"`
	MaxAttempts      int       `json:"max_attempts"`
	Status           string    `json:"status"` // pending, success, failed, retrying
	StatusCode       *int      `json:"status_code,omitempty"`
	ResponseBody     *string   `json:"response_body,omitempty"`
	ErrorMessage     *string   `json:"error_message,omitempty"`
	NextRetryAt      *time.Time `json:"next_retry_at,omitempty"`
	FirstAttemptedAt *time.Time `json:"first_attempted_at,omitempty"`
	LastAttemptedAt  *time.Time `json:"last_attempted_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// WebhookEvent represents the event types that can trigger webhooks
type WebhookEvent string

const (
	// Email events
	WebhookEventEmailReceived      WebhookEvent = "email.received"
	WebhookEventEmailSent          WebhookEvent = "email.sent"
	WebhookEventEmailDelivered     WebhookEvent = "email.delivered"
	WebhookEventEmailBounced       WebhookEvent = "email.bounced"
	WebhookEventEmailFailed        WebhookEvent = "email.failed"
	WebhookEventEmailQueued        WebhookEvent = "email.queued"

	// Security events
	WebhookEventSecurityVirusDetected    WebhookEvent = "security.virus_detected"
	WebhookEventSecuritySpamDetected     WebhookEvent = "security.spam_detected"
	WebhookEventSecurityLoginFailed      WebhookEvent = "security.login_failed"
	WebhookEventSecurityLoginSuccess     WebhookEvent = "security.login_success"
	WebhookEventSecurityBruteForce       WebhookEvent = "security.brute_force"
	WebhookEventSecurityIPBlacklisted    WebhookEvent = "security.ip_blacklisted"

	// DKIM/SPF/DMARC events
	WebhookEventDKIMFailed         WebhookEvent = "dkim.failed"
	WebhookEventSPFFailed          WebhookEvent = "spf.failed"
	WebhookEventDMARCFailed        WebhookEvent = "dmarc.failed"

	// User events
	WebhookEventUserCreated        WebhookEvent = "user.created"
	WebhookEventUserDeleted        WebhookEvent = "user.deleted"
	WebhookEventUserQuotaExceeded  WebhookEvent = "user.quota_exceeded"
)

// WebhookDeliveryStatus represents the status of a webhook delivery
const (
	WebhookStatusPending  = "pending"
	WebhookStatusRetrying = "retrying"
	WebhookStatusSuccess  = "success"
	WebhookStatusFailed   = "failed"
)

// WebhookPayload is the structure sent to webhook endpoints
type WebhookPayload struct {
	Event     string                 `json:"event"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}
