package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"go.uber.org/zap"
)

// WebhookService handles webhook delivery and management
type WebhookService struct {
	Repo       repository.WebhookRepository
	logger     *zap.Logger
	httpClient *http.Client
}

// NewWebhookService creates a new webhook service
func NewWebhookService(repo repository.WebhookRepository, logger *zap.Logger) *WebhookService {
	return &WebhookService{
		Repo:   repo,
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TriggerEvent sends an event to all registered webhooks
func (s *WebhookService) TriggerEvent(ctx context.Context, event domain.WebhookEvent, data map[string]interface{}) error {
	webhooks, err := s.Repo.ListActive(ctx)
	if err != nil {
		s.logger.Error("failed to list active webhooks", zap.Error(err))
		return fmt.Errorf("failed to list active webhooks: %w", err)
	}

	if len(webhooks) == 0 {
		s.logger.Debug("no active webhooks found", zap.String("event", string(event)))
		return nil
	}

	// Create payload
	payload := domain.WebhookPayload{
		Event:     string(event),
		Timestamp: time.Now(),
		Data:      data,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create delivery records for matching webhooks
	for _, webhook := range webhooks {
		if !s.shouldTrigger(webhook, event) {
			continue
		}

		delivery := &domain.WebhookDelivery{
			WebhookID:    webhook.ID,
			EventType:    string(event),
			Payload:      string(payloadJSON),
			AttemptCount: 0,
			MaxAttempts:  10,
			Status:       domain.WebhookStatusPending,
		}

		if err := s.Repo.CreateDelivery(ctx, delivery); err != nil {
			s.logger.Error("failed to create delivery",
				zap.Error(err),
				zap.Int64("webhook_id", webhook.ID),
				zap.String("event", string(event)),
			)
			continue
		}

		s.logger.Info("webhook delivery queued",
			zap.Int64("delivery_id", delivery.ID),
			zap.Int64("webhook_id", webhook.ID),
			zap.String("event", string(event)),
		)
	}

	return nil
}

// ProcessPendingDeliveries processes all pending webhook deliveries
func (s *WebhookService) ProcessPendingDeliveries(ctx context.Context) error {
	deliveries, err := s.Repo.ListPendingDeliveries(ctx, 100)
	if err != nil {
		return fmt.Errorf("failed to list pending deliveries: %w", err)
	}

	if len(deliveries) == 0 {
		return nil
	}

	s.logger.Info("processing pending deliveries", zap.Int("count", len(deliveries)))

	for _, delivery := range deliveries {
		if err := s.deliverWebhook(ctx, delivery); err != nil {
			s.logger.Error("failed to deliver webhook",
				zap.Error(err),
				zap.Int64("delivery_id", delivery.ID),
			)
		}
	}

	return nil
}

// deliverWebhook attempts to deliver a single webhook
func (s *WebhookService) deliverWebhook(ctx context.Context, delivery *domain.WebhookDelivery) error {
	// Get webhook details
	webhook, err := s.Repo.GetByID(ctx, delivery.WebhookID)
	if err != nil {
		return fmt.Errorf("failed to get webhook: %w", err)
	}
	if webhook == nil || !webhook.Active {
		delivery.Status = domain.WebhookStatusFailed
		errMsg := "webhook not found or inactive"
		delivery.ErrorMessage = &errMsg
		_ = s.Repo.UpdateDelivery(ctx, delivery)
		return fmt.Errorf("webhook %d not found or inactive", delivery.WebhookID)
	}

	// Update attempt count
	delivery.AttemptCount++
	now := time.Now()
	if delivery.FirstAttemptedAt == nil {
		delivery.FirstAttemptedAt = &now
	}
	delivery.LastAttemptedAt = &now

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", webhook.URL, bytes.NewBufferString(delivery.Payload))
	if err != nil {
		return s.handleDeliveryError(ctx, delivery, fmt.Sprintf("failed to create request: %v", err))
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gomailserver-webhook/1.0")
	req.Header.Set("X-Webhook-Event", delivery.EventType)
	req.Header.Set("X-Webhook-Delivery-ID", fmt.Sprintf("%d", delivery.ID))
	req.Header.Set("X-Webhook-Timestamp", now.Format(time.RFC3339))

	// Add signature
	signature := s.generateSignature(webhook.Secret, delivery.Payload)
	req.Header.Set("X-Webhook-Signature", signature)

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return s.handleDeliveryError(ctx, delivery, fmt.Sprintf("request failed: %v", err))
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // 1MB limit
	if err != nil {
		bodyBytes = []byte("failed to read response")
	}
	responseBody := string(bodyBytes)
	delivery.ResponseBody = &responseBody
	delivery.StatusCode = &resp.StatusCode

	// Check status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Success
		delivery.Status = domain.WebhookStatusSuccess
		delivery.CompletedAt = &now
		s.logger.Info("webhook delivered successfully",
			zap.Int64("delivery_id", delivery.ID),
			zap.Int64("webhook_id", webhook.ID),
			zap.Int("status_code", resp.StatusCode),
			zap.Int("attempt", delivery.AttemptCount),
		)
	} else {
		// Failed
		errMsg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, responseBody)
		return s.handleDeliveryError(ctx, delivery, errMsg)
	}

	return s.Repo.UpdateDelivery(ctx, delivery)
}

// handleDeliveryError handles delivery failures with retry logic
func (s *WebhookService) handleDeliveryError(ctx context.Context, delivery *domain.WebhookDelivery, errorMessage string) error {
	delivery.ErrorMessage = &errorMessage

	if delivery.AttemptCount >= delivery.MaxAttempts {
		// Max attempts reached, mark as failed
		delivery.Status = domain.WebhookStatusFailed
		now := time.Now()
		delivery.CompletedAt = &now
		s.logger.Warn("webhook delivery failed permanently",
			zap.Int64("delivery_id", delivery.ID),
			zap.Int("attempts", delivery.AttemptCount),
			zap.String("error", errorMessage),
		)
	} else {
		// Schedule retry with exponential backoff
		delivery.Status = domain.WebhookStatusRetrying
		nextRetry := s.calculateNextRetry(delivery.AttemptCount)
		delivery.NextRetryAt = &nextRetry
		s.logger.Info("webhook delivery will retry",
			zap.Int64("delivery_id", delivery.ID),
			zap.Int("attempt", delivery.AttemptCount),
			zap.Time("next_retry", nextRetry),
			zap.String("error", errorMessage),
		)
	}

	return s.Repo.UpdateDelivery(ctx, delivery)
}

// calculateNextRetry calculates the next retry time using exponential backoff
func (s *WebhookService) calculateNextRetry(attemptCount int) time.Time {
	// Exponential backoff: 2^attempt * base delay
	// Attempts: 1=2s, 2=4s, 3=8s, 4=16s, 5=32s, 6=64s, 7=128s, 8=256s, 9=512s, 10=1024s
	baseDelay := 2 * time.Second
	maxDelay := 1 * time.Hour

	delay := time.Duration(math.Pow(2, float64(attemptCount))) * baseDelay
	if delay > maxDelay {
		delay = maxDelay
	}

	return time.Now().Add(delay)
}

// generateSignature generates an HMAC-SHA256 signature for the payload
func (s *WebhookService) generateSignature(secret, payload string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// shouldTrigger checks if a webhook should be triggered for an event
func (s *WebhookService) shouldTrigger(webhook *domain.Webhook, event domain.WebhookEvent) bool {
	// Parse event types (comma-separated)
	eventTypes := strings.Split(webhook.EventTypes, ",")
	for _, et := range eventTypes {
		et = strings.TrimSpace(et)
		if et == "*" || et == string(event) {
			return true
		}
		// Check for wildcard patterns (e.g., "email.*" matches "email.received")
		if strings.HasSuffix(et, ".*") {
			prefix := strings.TrimSuffix(et, ".*")
			if strings.HasPrefix(string(event), prefix+".") {
				return true
			}
		}
	}
	return false
}

// TestWebhook sends a test payload to a webhook
func (s *WebhookService) TestWebhook(ctx context.Context, webhookID int64) error {
	webhook, err := s.Repo.GetByID(ctx, webhookID)
	if err != nil {
		return fmt.Errorf("failed to get webhook: %w", err)
	}
	if webhook == nil {
		return fmt.Errorf("webhook not found")
	}

	// Create test delivery
	testData := map[string]interface{}{
		"test": true,
		"message": "This is a test webhook delivery",
	}

	delivery := &domain.WebhookDelivery{
		WebhookID:    webhook.ID,
		EventType:    "test.ping",
		Payload:      "",
		AttemptCount: 0,
		MaxAttempts:  1,
		Status:       domain.WebhookStatusPending,
	}

	payload := domain.WebhookPayload{
		Event:     "test.ping",
		Timestamp: time.Now(),
		Data:      testData,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal test payload: %w", err)
	}
	delivery.Payload = string(payloadJSON)

	if err := s.Repo.CreateDelivery(ctx, delivery); err != nil {
		return fmt.Errorf("failed to create test delivery: %w", err)
	}

	return s.deliverWebhook(ctx, delivery)
}

// CleanupOldDeliveries removes old delivery records
func (s *WebhookService) CleanupOldDeliveries(ctx context.Context, olderThanDays int) error {
	return s.Repo.DeleteOldDeliveries(ctx, olderThanDays)
}
