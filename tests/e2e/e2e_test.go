package tests

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/service"
)

// TestServiceIntegration tests that all services are properly wired together
func TestServiceIntegration(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Setup repositories
	userRepo := repository.NewSQLiteUserRepository(db)
	domainRepo := repository.NewSQLiteDomainRepository(db)
	mailboxRepo := repository.NewSQLiteMailboxRepository(db)
	messageRepo := repository.NewSQLiteMessageRepository(db)
	queueRepo := repository.NewSQLiteQueueRepository(db)

	// Setup services
	userSvc := service.NewUserService(&repository.Repositories{
		User:   userRepo,
		Domain: domainRepo,
	}, zap.NewNop())
	mailboxSvc := service.NewMailboxService(mailboxRepo, zap.NewNop())
	messageSvc := service.NewMessageService(messageRepo, "./testdata", zap.NewNop())
	queueSvc := service.NewQueueService(queueRepo, nil, zap.NewNop())

	// Wire services
	messageSvc.SetUserService(userSvc)
	messageSvc.SetMailboxService(mailboxSvc)
	messageSvc.SetQueueService(queueSvc)

	// Create test domain and user
	domain := &domain.Domain{
		Name:   "example.com",
		Status: "active",
	}
	if err := domainRepo.Create(domain); err != nil {
		t.Fatalf("Failed to create test domain: %v", err)
	}

	user := &domain.User{
		DomainID: domain.ID,
		Email:    "test@example.com",
		Name:     "Test User",
		Status:   "active",
	}
	if err := userSvc.Create(user, "password123"); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test queue service enqueues message
	t.Run("Queue_Enqueue", func(t *testing.T) {
		messageID, err := queueSvc.Enqueue("sender@example.com", []string{"test@example.com"}, []byte("test message"))
		if err != nil {
			t.Fatalf("Failed to enqueue message: %v", err)
		}
		if messageID == "" {
			t.Error("Expected non-empty message ID")
		}
	})

	// Test mailbox creation
	t.Run("Mailbox_Creation", func(t *testing.T) {
		err := mailboxSvc.Create(user.ID, "Drafts", "\\Drafts")
		if err != nil {
			t.Fatalf("Failed to create Drafts mailbox: %v", err)
		}

		mb, err := mailboxSvc.GetByName(user.ID, "Drafts")
		if err != nil {
			t.Fatalf("Failed to retrieve Drafts mailbox: %v", err)
		}
		if mb.Name != "Drafts" || mb.SpecialUse != "\\Drafts" {
			t.Errorf("Mailbox not created correctly: %+v", mb)
		}
	})

	// Test draft saving
	t.Run("Draft_Saving", func(t *testing.T) {
		draft, err := messageSvc.SaveDraft(context.Background(), int(user.ID), nil, &service.DraftData{
			To:       []string{"recipient@example.com"},
			Subject:  "Test Draft",
			BodyText: "This is a test draft",
		})
		if err != nil {
			t.Fatalf("Failed to save draft: %v", err)
		}
		if draft.ID == 0 {
			t.Error("Expected draft to have valid ID")
		}
	})

	// Test message sending (webmail integration)
	t.Run("Message_Sending", func(t *testing.T) {
		messageID, err := messageSvc.SendMessage(context.Background(), int(user.ID), &service.SendMessageRequest{
			From:     "test@example.com",
			To:       "recipient@example.com",
			Subject:  "Test Message",
			BodyText: "Hello World",
		})
		if err != nil {
			t.Fatalf("Failed to send message: %v", err)
		}
		if messageID == 0 {
			t.Error("Expected message to have valid ID")
		}
	})
}
