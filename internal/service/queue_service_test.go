package service

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
)

// mockQueueRepository is a test double for QueueRepository
type mockQueueRepository struct {
	enqueueFunc       func(*domain.QueueItem) error
	getPendingFunc    func() ([]*domain.QueueItem, error)
	getByIDFunc       func(int64) (*domain.QueueItem, error)
	updateStatusFunc  func(int64, string, string) error
	updateRetryFunc   func(int64, int, time.Time) error
	deleteFunc        func(int64) error
}

func (m *mockQueueRepository) Enqueue(item *domain.QueueItem) error {
	if m.enqueueFunc != nil {
		return m.enqueueFunc(item)
	}
	item.ID = 1
	return nil
}

func (m *mockQueueRepository) GetPending() ([]*domain.QueueItem, error) {
	if m.getPendingFunc != nil {
		return m.getPendingFunc()
	}
	return []*domain.QueueItem{}, nil
}

func (m *mockQueueRepository) GetByID(id int64) (*domain.QueueItem, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, errors.New("not found")
}

func (m *mockQueueRepository) UpdateStatus(id int64, status, errorMsg string) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(id, status, errorMsg)
	}
	return nil
}

func (m *mockQueueRepository) UpdateRetry(id int64, retryCount int, nextRetry time.Time) error {
	if m.updateRetryFunc != nil {
		return m.updateRetryFunc(id, retryCount, nextRetry)
	}
	return nil
}

func (m *mockQueueRepository) Delete(id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func TestQueueService_Enqueue(t *testing.T) {
	logger := zap.NewNop()

	t.Run("enqueues message successfully", func(t *testing.T) {
		repo := &mockQueueRepository{}
		svc := NewQueueService(repo, logger)

		_, err := svc.Enqueue("sender@example.com", []string{"recipient@example.com"}, []byte("test message data"))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("enqueues with multiple recipients", func(t *testing.T) {
		var capturedItem *domain.QueueItem
		repo := &mockQueueRepository{
			enqueueFunc: func(item *domain.QueueItem) error {
				capturedItem = item
				return nil
			},
		}
		svc := NewQueueService(repo, logger)

		recipients := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
		_, err := svc.Enqueue("sender@example.com", recipients, []byte("test message data"))

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if capturedItem == nil {
			t.Fatal("expected item to be captured")
		}

		if capturedItem.Sender != "sender@example.com" {
			t.Errorf("expected sender 'sender@example.com', got '%s'", capturedItem.Sender)
		}

		if capturedItem.MessagePath == "" {
			t.Error("expected message path to be set")
		}

		// Recipients should be JSON array
		if capturedItem.Recipients == "" {
			t.Error("expected recipients to be set")
		}

		if capturedItem.Status != "pending" {
			t.Errorf("expected status 'pending', got '%s'", capturedItem.Status)
		}

		if capturedItem.RetryCount != 0 {
			t.Errorf("expected retry count 0, got %d", capturedItem.RetryCount)
		}

		if capturedItem.MaxRetries != 9 {
			t.Errorf("expected max retries 9, got %d", capturedItem.MaxRetries)
		}
	})

	t.Run("returns error if repository fails", func(t *testing.T) {
		repo := &mockQueueRepository{
			enqueueFunc: func(item *domain.QueueItem) error {
				return errors.New("database error")
			},
		}
		svc := NewQueueService(repo, logger)

		_, err := svc.Enqueue("sender@example.com", []string{"recipient@example.com"}, []byte("test data"))

		if err == nil {
			t.Error("expected error from repository failure")
		}
	})
}

func TestQueueService_CalculateNextRetry(t *testing.T) {
	logger := zap.NewNop()
	repo := &mockQueueRepository{}
	svc := NewQueueService(repo, logger)

	// Test exponential backoff schedule
	tests := []struct {
		retryCount      int
		expectedMinutes int
	}{
		{0, 5},    // First retry: 5 minutes
		{1, 15},   // Second retry: 15 minutes
		{2, 30},   // Third retry: 30 minutes
		{3, 60},   // Fourth retry: 1 hour
		{4, 120},  // Fifth retry: 2 hours
		{5, 240},  // Sixth retry: 4 hours
		{6, 480},  // Seventh retry: 8 hours
		{7, 960},  // Eighth retry: 16 hours
		{8, 1440}, // Ninth retry: 24 hours
	}

	now := time.Now()
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			nextRetry := svc.CalculateNextRetry(tt.retryCount, now)
			expectedTime := now.Add(time.Duration(tt.expectedMinutes) * time.Minute)

			// Allow 1 second tolerance for timing differences
			diff := nextRetry.Sub(expectedTime)
			if diff < -time.Second || diff > time.Second {
				t.Errorf("retry count %d: expected %v, got %v (diff: %v)",
					tt.retryCount, expectedTime, nextRetry, diff)
			}
		})
	}
}

func TestQueueService_GetPending(t *testing.T) {
	logger := zap.NewNop()

	t.Run("returns pending items", func(t *testing.T) {
		expectedItems := []*domain.QueueItem{
			{ID: 1, Status: "pending", Sender: "sender1@example.com"},
			{ID: 2, Status: "pending", Sender: "sender2@example.com"},
		}

		repo := &mockQueueRepository{
			getPendingFunc: func() ([]*domain.QueueItem, error) {
				return expectedItems, nil
			},
		}

		svc := NewQueueService(repo, logger)
		items, err := svc.GetPending()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(items) != 2 {
			t.Errorf("expected 2 items, got %d", len(items))
		}
	})

	t.Run("handles repository error", func(t *testing.T) {
		repo := &mockQueueRepository{
			getPendingFunc: func() ([]*domain.QueueItem, error) {
				return nil, errors.New("database error")
			},
		}

		svc := NewQueueService(repo, logger)
		items, err := svc.GetPending()

		if err == nil {
			t.Error("expected error from repository")
		}

		if items != nil {
			t.Error("expected nil items on error")
		}
	})
}

func TestQueueService_MarkDelivered(t *testing.T) {
	logger := zap.NewNop()

	t.Run("marks item as delivered", func(t *testing.T) {
		var capturedStatus string
		repo := &mockQueueRepository{
			updateStatusFunc: func(id int64, status, errorMsg string) error {
				capturedStatus = status
				return nil
			},
		}

		svc := NewQueueService(repo, logger)
		err := svc.MarkDelivered(1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if capturedStatus != "delivered" {
			t.Errorf("expected status 'delivered', got '%s'", capturedStatus)
		}
	})
}

func TestQueueService_MarkFailed(t *testing.T) {
	logger := zap.NewNop()

	t.Run("marks item as failed with error message", func(t *testing.T) {
		var capturedStatus string
		var capturedError string

		repo := &mockQueueRepository{
			updateStatusFunc: func(id int64, status, errorMsg string) error {
				capturedStatus = status
				capturedError = errorMsg
				return nil
			},
		}

		svc := NewQueueService(repo, logger)
		err := svc.MarkFailed(1, "SMTP connection refused")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if capturedStatus != "failed" {
			t.Errorf("expected status 'failed', got '%s'", capturedStatus)
		}

		if capturedError != "SMTP connection refused" {
			t.Errorf("expected error message to be set, got '%s'", capturedError)
		}
	})
}

func TestQueueService_IncrementRetry(t *testing.T) {
	logger := zap.NewNop()

	t.Run("increments retry count and sets next retry time", func(t *testing.T) {
		var capturedRetryCount int
		var capturedNextRetry time.Time

		repo := &mockQueueRepository{
			updateRetryFunc: func(id int64, retryCount int, nextRetry time.Time) error {
				capturedRetryCount = retryCount
				capturedNextRetry = nextRetry
				return nil
			},
		}

		svc := NewQueueService(repo, logger)
		now := time.Now()
		err := svc.IncrementRetry(1, 2, now)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if capturedRetryCount != 3 {
			t.Errorf("expected retry count 3, got %d", capturedRetryCount)
		}

		// Verify next retry is in the future (30 minutes for retry count 2)
		expectedRetry := now.Add(30 * time.Minute)
		diff := capturedNextRetry.Sub(expectedRetry)
		if diff < -time.Second || diff > time.Second {
			t.Errorf("expected next retry around %v, got %v", expectedRetry, capturedNextRetry)
		}
	})
}
