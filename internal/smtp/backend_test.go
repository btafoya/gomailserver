package smtp

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/emersion/go-smtp"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
)

// mockUserService for SMTP backend tests
type mockUserService struct {
	authenticateFunc func(string, string) (*domain.User, error)
}

func (m *mockUserService) Create(user *domain.User, password string) error {
	return nil
}

func (m *mockUserService) Authenticate(email, password string) (*domain.User, error) {
	if m.authenticateFunc != nil {
		return m.authenticateFunc(email, password)
	}
	return nil, errors.New("authentication failed")
}

func (m *mockUserService) GetByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserService) GetByID(id int64) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserService) Update(user *domain.User) error {
	return nil
}

func (m *mockUserService) UpdatePassword(userID int64, newPassword string) error {
	return nil
}

func (m *mockUserService) Delete(id int64) error {
	return nil
}

// mockDomainRepository for SMTP backend tests
type mockDomainRepository struct{}

func (m *mockDomainRepository) Create(domain *domain.Domain) error          { return nil }
func (m *mockDomainRepository) GetByID(id int64) (*domain.Domain, error)    { return nil, nil }
func (m *mockDomainRepository) GetByName(name string) (*domain.Domain, error) { return nil, nil }
func (m *mockDomainRepository) Update(domain *domain.Domain) error          { return nil }
func (m *mockDomainRepository) Delete(id int64) error                       { return nil }
func (m *mockDomainRepository) List(offset, limit int) ([]*domain.Domain, error) { return nil, nil }
func (m *mockDomainRepository) CreateTemplate(template *domain.Domain) error { return nil }
func (m *mockDomainRepository) GetDefaultTemplate() (*domain.Domain, error) { return nil, nil }

// mockMessageService for SMTP backend tests
type mockMessageService struct {
	storeFunc func(int64, int64, int64, []byte) (*domain.Message, error)
}

func (m *mockMessageService) Store(userID, mailboxID, uid int64, messageData []byte) (*domain.Message, error) {
	if m.storeFunc != nil {
		return m.storeFunc(userID, mailboxID, uid, messageData)
	}
	return &domain.Message{ID: 1}, nil
}

func (m *mockMessageService) GetByID(id int64) (*domain.Message, error) {
	return nil, nil
}

func (m *mockMessageService) GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error) {
	return nil, nil
}

func (m *mockMessageService) Delete(id int64) error {
	return nil
}

// mockQueueService for SMTP backend tests
type mockQueueService struct {
	enqueueFunc func(string, []string, []byte) (string, error)
}

func (m *mockQueueService) Enqueue(sender string, recipients []string, message []byte) (string, error) {
	if m.enqueueFunc != nil {
		return m.enqueueFunc(sender, recipients, message)
	}
	return "message-id", nil
}

func (m *mockQueueService) GetPending() ([]*domain.QueueItem, error) {
	return nil, nil
}

func (m *mockQueueService) MarkDelivered(id int64) error {
	return nil
}

func (m *mockQueueService) MarkFailed(id int64, errorMsg string) error {
	return nil
}

func (m *mockQueueService) IncrementRetry(id int64, currentRetryCount int, failedAt time.Time) error {
	return nil
}

func (m *mockQueueService) CalculateNextRetry(retryCount int, failedAt time.Time) time.Time {
	return time.Now()
}

func TestBackend_NewSession(t *testing.T) {
	// NewSession requires *smtp.Conn which we can't easily mock in unit tests
	// This test is skipped as it requires integration testing with actual SMTP connection
	t.Skip("Skipping NewSession test - requires actual smtp.Conn instance")
}

func TestSession_AuthPlain(t *testing.T) {
	logger := zap.NewNop()

	t.Run("authenticates with valid credentials", func(t *testing.T) {
		userSvc := &mockUserService{
			authenticateFunc: func(email, password string) (*domain.User, error) {
				if email == "test@example.com" && password == "correctpassword" {
					return &domain.User{
						ID:     1,
						Email:  email,
						Status: "active",
					}, nil
				}
				return nil, errors.New("authentication failed")
			},
		}

		backend := &Backend{
			userService:    userSvc,
			messageService: &mockMessageService{},
			queueService:   &mockQueueService{},
			domainRepo:     &mockDomainRepository{},
			logger:         logger,
		}

		session := &Session{backend: backend, logger: logger}

		err := session.AuthPlain("test@example.com", "correctpassword")
		if err != nil {
			t.Fatalf("expected authentication to succeed, got error: %v", err)
		}

		if !session.authenticated {
			t.Error("expected session to be authenticated")
		}

		if session.username != "test@example.com" {
			t.Errorf("expected username test@example.com, got %s", session.username)
		}
	})

	t.Run("rejects invalid credentials", func(t *testing.T) {
		userSvc := &mockUserService{
			authenticateFunc: func(email, password string) (*domain.User, error) {
				return nil, errors.New("authentication failed")
			},
		}

		backend := &Backend{
			userService:    userSvc,
			messageService: &mockMessageService{},
			queueService:   &mockQueueService{},
			domainRepo:     &mockDomainRepository{},
			logger:         logger,
		}

		session := &Session{backend: backend, logger: logger}

		err := session.AuthPlain("test@example.com", "wrongpassword")
		if err == nil {
			t.Error("expected authentication to fail with wrong password")
		}

		if session.authenticated {
			t.Error("expected session to not be authenticated")
		}
	})
}

func TestSession_Mail(t *testing.T) {
	logger := zap.NewNop()
	backend := &Backend{
		userService:    &mockUserService{},
		messageService: &mockMessageService{},
		queueService:   &mockQueueService{},
		domainRepo:     &mockDomainRepository{},
		logger:         logger,
	}

	t.Run("accepts valid sender", func(t *testing.T) {
		session := &Session{backend: backend, logger: logger, authenticated: true}

		err := session.Mail("sender@example.com", &smtp.MailOptions{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if session.from != "sender@example.com" {
			t.Errorf("expected from 'sender@example.com', got '%s'", session.from)
		}
	})

	t.Run("rejects mail without authentication for submission", func(t *testing.T) {
		// Skip: This test requires smtp.Conn to check server port
		// Can't mock smtp.Conn easily in unit tests
		t.Skip("Skipping port-based authentication check - requires smtp.Conn")
	})
}

func TestSession_Rcpt(t *testing.T) {
	logger := zap.NewNop()
	backend := &Backend{
		userService:    &mockUserService{},
		messageService: &mockMessageService{},
		queueService:   &mockQueueService{},
		domainRepo:     &mockDomainRepository{},
		logger:         logger,
	}

	t.Run("accepts valid recipient", func(t *testing.T) {
		session := &Session{
			backend:       backend,
			logger:        logger,
			authenticated: true,
			from:          "sender@example.com",
		}

		err := session.Rcpt("recipient@example.com", &smtp.RcptOptions{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(session.to) != 1 {
			t.Errorf("expected 1 recipient, got %d", len(session.to))
		}

		if session.to[0] != "recipient@example.com" {
			t.Errorf("expected recipient 'recipient@example.com', got '%s'", session.to[0])
		}
	})

	t.Run("accepts multiple recipients", func(t *testing.T) {
		session := &Session{
			backend:       backend,
			logger:        logger,
			authenticated: true,
			from:          "sender@example.com",
		}

		recipients := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
		for _, rcpt := range recipients {
			err := session.Rcpt(rcpt, &smtp.RcptOptions{})
			if err != nil {
				t.Fatalf("expected no error for recipient %s, got %v", rcpt, err)
			}
		}

		if len(session.to) != len(recipients) {
			t.Errorf("expected %d recipients, got %d", len(recipients), len(session.to))
		}
	})
}

func TestSession_Data(t *testing.T) {
	logger := zap.NewNop()

	t.Run("accepts and processes message", func(t *testing.T) {
		queueCalled := false
		queueSvc := &mockQueueService{
			enqueueFunc: func(sender string, recipients []string, message []byte) (string, error) {
				queueCalled = true
				if sender != "sender@example.com" {
					t.Errorf("expected sender 'sender@example.com', got '%s'", sender)
				}
				if len(recipients) != 2 {
					t.Errorf("expected 2 recipients, got %d", len(recipients))
				}
				return "test-message-id", nil
			},
		}

		backend := &Backend{
			userService:    &mockUserService{},
			messageService: &mockMessageService{},
			queueService:   queueSvc,
			domainRepo:     &mockDomainRepository{},
			logger:         logger,
		}

		session := &Session{
			backend:       backend,
			logger:        logger,
			authenticated: true,
			from:          "sender@example.com",
			to:            []string{"user1@example.com", "user2@example.com"},
		}

		emailData := `From: sender@example.com
To: user1@example.com, user2@example.com
Subject: Test Email

This is a test message.`

		reader := strings.NewReader(emailData)
		err := session.Data(reader)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !queueCalled {
			t.Error("expected queue service to be called")
		}
	})

	t.Run("reads full message body", func(t *testing.T) {
		var capturedData []byte
		queueSvc := &mockQueueService{
			enqueueFunc: func(sender string, recipients []string, message []byte) (string, error) {
				capturedData = message
				return "test-message-id", nil
			},
		}

		backend := &Backend{
			userService:    &mockUserService{},
			messageService: &mockMessageService{},
			queueService:   queueSvc,
			domainRepo:     &mockDomainRepository{},
			logger:         logger,
		}

		session := &Session{
			backend:       backend,
			logger:        logger,
			authenticated: true,
			from:          "sender@example.com",
			to:            []string{"recipient@example.com"},
		}

		emailData := "From: sender@example.com\r\nSubject: Test\r\n\r\nBody content here"
		reader := strings.NewReader(emailData)
		err := session.Data(reader)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !bytes.Contains(capturedData, []byte("Body content here")) {
			t.Error("message body not captured correctly")
		}
	})

	t.Run("handles empty message", func(t *testing.T) {
		backend := &Backend{
			userService:    &mockUserService{},
			messageService: &mockMessageService{},
			queueService:   &mockQueueService{},
			domainRepo:     &mockDomainRepository{},
			logger:         logger,
		}

		session := &Session{
			backend:       backend,
			logger:        logger,
			authenticated: true,
			from:          "sender@example.com",
			to:            []string{"recipient@example.com"},
		}

		reader := strings.NewReader("")
		err := session.Data(reader)

		// Should handle empty message gracefully
		if err != nil && err != io.EOF {
			t.Errorf("unexpected error for empty message: %v", err)
		}
	})
}

func TestSession_Reset(t *testing.T) {
	logger := zap.NewNop()
	backend := &Backend{
		userService:    &mockUserService{},
		messageService: &mockMessageService{},
		queueService:   &mockQueueService{},
		domainRepo:     &mockDomainRepository{},
		logger:         logger,
	}

	session := &Session{
		backend:       backend,
		logger:        logger,
		authenticated: true,
		from:          "sender@example.com",
		to:            []string{"recipient@example.com"},
	}

	session.Reset()

	if session.from != "" {
		t.Error("expected from to be reset")
	}

	if len(session.to) != 0 {
		t.Error("expected recipients to be cleared")
	}

	// Authentication should persist across RSET
	if !session.authenticated {
		t.Error("expected authentication to persist")
	}
}

func TestSession_Logout(t *testing.T) {
	logger := zap.NewNop()
	backend := &Backend{
		userService:    &mockUserService{},
		messageService: &mockMessageService{},
		queueService:   &mockQueueService{},
		domainRepo:     &mockDomainRepository{},
		logger:         logger,
	}

	session := &Session{
		backend:       backend,
		logger:        logger,
		authenticated: true,
		username:      "test@example.com",
	}

	err := session.Logout()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Session should still be usable after logout (for next command)
}
