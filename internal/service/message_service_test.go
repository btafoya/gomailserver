package service

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
)

// mockMessageRepository is a test double for MessageRepository
type mockMessageRepository struct {
	createFunc      func(*domain.Message) error
	getByIDFunc     func(int64) (*domain.Message, error)
	getByMailboxFunc func(int64, int, int) ([]*domain.Message, error)
	deleteFunc      func(int64) error
}

func (m *mockMessageRepository) Create(msg *domain.Message) error {
	if m.createFunc != nil {
		return m.createFunc(msg)
	}
	msg.ID = 1
	return nil
}

func (m *mockMessageRepository) GetByID(id int64) (*domain.Message, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *mockMessageRepository) GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error) {
	if m.getByMailboxFunc != nil {
		return m.getByMailboxFunc(mailboxID, offset, limit)
	}
	return []*domain.Message{}, nil
}

func (m *mockMessageRepository) Update(message *domain.Message) error {
	return nil
}

func (m *mockMessageRepository) Delete(id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func TestMessageService_Store_SmallMessage(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	repo := &mockMessageRepository{}
	svc := NewMessageService(repo, tempDir, logger)

	// Create a small test email (< 1MB)
	smallEmail := createTestEmail("test@example.com", "recipient@example.com", "Test Subject", "Small body")

	msg, err := svc.Store(1, 1, 100, []byte(smallEmail))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify blob storage
	if msg.StorageType != "blob" {
		t.Errorf("expected storage type 'blob', got '%s'", msg.StorageType)
	}

	if len(msg.Content) == 0 {
		t.Error("expected content to be stored in blob")
	}

	if msg.ContentPath != "" {
		t.Error("expected no content path for blob storage")
	}

	// Verify headers were parsed
	if msg.Subject != "Test Subject" {
		t.Errorf("expected subject 'Test Subject', got '%s'", msg.Subject)
	}

	if msg.From != "test@example.com" {
		t.Errorf("expected from 'test@example.com', got '%s'", msg.From)
	}

	if msg.To != "recipient@example.com" {
		t.Errorf("expected to 'recipient@example.com', got '%s'", msg.To)
	}
}

func TestMessageService_Store_LargeMessage(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	repo := &mockMessageRepository{}
	svc := NewMessageService(repo, tempDir, logger)

	// Create a large test email (> 1MB)
	largeBody := strings.Repeat("A", 1024*1024+1000) // > 1MB
	largeEmail := createTestEmail("sender@example.com", "recipient@example.com", "Large Email", largeBody)

	msg, err := svc.Store(1, 2, 101, []byte(largeEmail))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify file storage
	if msg.StorageType != "file" {
		t.Errorf("expected storage type 'file', got '%s'", msg.StorageType)
	}

	if len(msg.Content) != 0 {
		t.Error("expected no content in blob for file storage")
	}

	if msg.ContentPath == "" {
		t.Error("expected content path for file storage")
	}

	// Verify file exists
	if _, err := os.Stat(msg.ContentPath); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s", msg.ContentPath)
	}

	// Verify file content
	fileContent, err := os.ReadFile(msg.ContentPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if !bytes.Equal(fileContent, []byte(largeEmail)) {
		t.Error("file content doesn't match original email")
	}
}

func TestMessageService_Store_ThreadIDGeneration(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	repo := &mockMessageRepository{}
	svc := NewMessageService(repo, tempDir, logger)

	t.Run("generates thread ID from Message-ID", func(t *testing.T) {
		email := `From: sender@example.com
To: recipient@example.com
Subject: Test
Message-ID: <unique-id-123@example.com>

Body`

		msg, err := svc.Store(1, 1, 100, []byte(email))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if msg.ThreadID == "" {
			t.Error("expected thread ID to be generated")
		}

		if msg.MessageID != "<unique-id-123@example.com>" {
			t.Errorf("expected Message-ID to be parsed, got '%s'", msg.MessageID)
		}
	})

	t.Run("uses In-Reply-To for thread ID", func(t *testing.T) {
		email := `From: sender@example.com
To: recipient@example.com
Subject: Re: Test
Message-ID: <reply-id-456@example.com>
In-Reply-To: <original-id-789@example.com>

Reply body`

		msg, err := svc.Store(1, 1, 101, []byte(email))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if msg.ThreadID == "" {
			t.Error("expected thread ID to be generated from In-Reply-To")
		}

		if msg.InReplyTo != "<original-id-789@example.com>" {
			t.Errorf("expected In-Reply-To to be parsed, got '%s'", msg.InReplyTo)
		}

		// Thread ID should be based on In-Reply-To, not Message-ID
		// Both messages in same thread should have same thread ID (when properly normalized)
	})
}

func TestMessageService_GetByID(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	t.Run("loads blob message", func(t *testing.T) {
		content := []byte("test email content")
		repo := &mockMessageRepository{
			getByIDFunc: func(id int64) (*domain.Message, error) {
				return &domain.Message{
					ID:          1,
					StorageType: "blob",
					Content:     content,
				}, nil
			},
		}

		svc := NewMessageService(repo, tempDir, logger)
		msg, err := svc.GetByID(1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !bytes.Equal(msg.Content, content) {
			t.Error("content doesn't match expected")
		}
	})

	t.Run("loads file message", func(t *testing.T) {
		// Create a test file
		testContent := []byte("file-stored email content")
		testFile := filepath.Join(tempDir, "test.eml")
		err := os.WriteFile(testFile, testContent, 0640)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		repo := &mockMessageRepository{
			getByIDFunc: func(id int64) (*domain.Message, error) {
				return &domain.Message{
					ID:          1,
					StorageType: "file",
					ContentPath: testFile,
				}, nil
			},
		}

		svc := NewMessageService(repo, tempDir, logger)
		msg, err := svc.GetByID(1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !bytes.Equal(msg.Content, testContent) {
			t.Error("content doesn't match file content")
		}
	})
}

func TestMessageService_Delete(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	t.Run("deletes file when message uses file storage", func(t *testing.T) {
		// Create a test file
		testFile := filepath.Join(tempDir, "deleteme.eml")
		err := os.WriteFile(testFile, []byte("content"), 0640)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		repo := &mockMessageRepository{
			getByIDFunc: func(id int64) (*domain.Message, error) {
				return &domain.Message{
					ID:          1,
					StorageType: "file",
					ContentPath: testFile,
				}, nil
			},
		}

		svc := NewMessageService(repo, tempDir, logger)
		err = svc.Delete(1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify file was deleted
		if _, err := os.Stat(testFile); !os.IsNotExist(err) {
			t.Error("expected file to be deleted")
		}
	})

	t.Run("handles blob message without error", func(t *testing.T) {
		repo := &mockMessageRepository{
			getByIDFunc: func(id int64) (*domain.Message, error) {
				return &domain.Message{
					ID:          1,
					StorageType: "blob",
					Content:     []byte("content"),
				}, nil
			},
		}

		svc := NewMessageService(repo, tempDir, logger)
		err := svc.Delete(1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

// Helper function to create test email
func createTestEmail(from, to, subject, body string) string {
	return `From: ` + from + `
To: ` + to + `
Subject: ` + subject + `
Content-Type: text/plain

` + body
}
