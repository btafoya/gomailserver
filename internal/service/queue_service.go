package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

// QueueService handles SMTP queue management
type QueueService struct {
	repo   repository.QueueRepository
	logger *zap.Logger
}

// NewQueueService creates a new queue service
func NewQueueService(repo repository.QueueRepository, logger *zap.Logger) *QueueService {
	return &QueueService{
		repo:   repo,
		logger: logger,
	}
}

// Enqueue adds a message to the delivery queue
func (s *QueueService) Enqueue(from string, to []string, message []byte) (string, error) {
	messageID := generateMessageID()
	messagePath := "/var/spool/mail/queue/" + messageID + ".eml"

	// Store message to disk
	// TODO: Implement actual file storage

	// Create queue entry
	item := &domain.QueueItem{
		Sender:      from,
		Recipients:  encodeRecipients(to),
		MessagePath: messagePath,
		Status:      "pending",
		RetryCount:  0,
		MaxRetries:  9,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Enqueue(item); err != nil {
		return "", err
	}

	s.logger.Info("message queued",
		zap.String("message_id", messageID),
		zap.String("from", from),
		zap.Strings("to", to),
		zap.Int("size", len(message)),
	)

	return messageID, nil
}

// encodeRecipients converts recipient list to JSON
func encodeRecipients(recipients []string) string {
	// Simple implementation for now
	result := "["
	for i, r := range recipients {
		if i > 0 {
			result += ","
		}
		result += `"` + r + `"`
	}
	result += "]"
	return result
}

// generateMessageID generates a unique message ID
func generateMessageID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GetPending retrieves all pending queue items
func (s *QueueService) GetPending() ([]*domain.QueueItem, error) {
	return s.repo.GetPending()
}

// GetPendingItems retrieves all pending queue items (API handler method)
func (s *QueueService) GetPendingItems(ctx context.Context) ([]*domain.QueueItem, error) {
	return s.repo.GetPending()
}

// GetByID retrieves a specific queue item by ID
func (s *QueueService) GetByID(ctx context.Context, id int64) (*domain.QueueItem, error) {
	return s.repo.GetByID(id)
}

// RetryItem resets a queue item for retry
func (s *QueueService) RetryItem(ctx context.Context, id int64) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Reset to pending status with retry count reset
	if err := s.repo.UpdateStatus(id, "pending", ""); err != nil {
		return err
	}

	s.logger.Info("queue item reset for retry",
		zap.Int64("id", id),
		zap.String("sender", item.Sender),
	)

	return nil
}

// DeleteItem removes a queue item
func (s *QueueService) DeleteItem(ctx context.Context, id int64) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("failed to delete queue item",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return err
	}

	s.logger.Info("queue item deleted",
		zap.Int64("id", id),
	)

	return nil
}

// MarkDelivered marks a queue item as successfully delivered
func (s *QueueService) MarkDelivered(id int64) error {
	return s.repo.UpdateStatus(id, "delivered", "")
}

// MarkFailed marks a queue item as permanently failed
func (s *QueueService) MarkFailed(id int64, errorMsg string) error {
	return s.repo.UpdateStatus(id, "failed", errorMsg)
}

// IncrementRetry increments the retry count and schedules next retry
func (s *QueueService) IncrementRetry(id int64, currentRetryCount int, failedAt time.Time) error {
	nextRetry := s.CalculateNextRetry(currentRetryCount, failedAt)
	return s.repo.UpdateRetry(id, currentRetryCount+1, nextRetry)
}

// CalculateNextRetry calculates next retry time using exponential backoff
func (s *QueueService) CalculateNextRetry(retryCount int, failedAt time.Time) time.Time {
	delays := []time.Duration{
		5 * time.Minute,
		15 * time.Minute,
		30 * time.Minute,
		1 * time.Hour,
		2 * time.Hour,
		4 * time.Hour,
		8 * time.Hour,
		16 * time.Hour,
		24 * time.Hour,
	}

	if retryCount >= len(delays) {
		return time.Time{} // Give up
	}

	return failedAt.Add(delays[retryCount])
}

// ProcessQueue processes pending queue items
func (s *QueueService) ProcessQueue() error {
	// TODO: Implement queue processing with retry logic
	return nil
}
