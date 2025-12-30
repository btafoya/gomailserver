package service

import (
	"time"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

// MailboxService handles mailbox operations
type MailboxService struct {
	repo   repository.MailboxRepository
	logger *zap.Logger
}

// NewMailboxService creates a new mailbox service
func NewMailboxService(repo repository.MailboxRepository, logger *zap.Logger) *MailboxService {
	return &MailboxService{
		repo:   repo,
		logger: logger,
	}
}

// List lists mailboxes for a user
func (s *MailboxService) List(userID int64, subscribed bool) ([]*domain.Mailbox, error) {
	s.logger.Debug("listing mailboxes",
		zap.Int64("user_id", userID),
		zap.Bool("subscribed", subscribed),
	)

	mailboxes, err := s.repo.GetByUser(userID)
	if err != nil {
		return nil, err
	}

	// Filter by subscription status if requested
	if subscribed {
		filtered := make([]*domain.Mailbox, 0)
		for _, mb := range mailboxes {
			if mb.Subscribed {
				filtered = append(filtered, mb)
			}
		}
		return filtered, nil
	}

	return mailboxes, nil
}

// GetByName retrieves a mailbox by name
func (s *MailboxService) GetByName(userID int64, name string) (*domain.Mailbox, error) {
	return s.repo.GetByName(userID, name)
}

// GetByID retrieves a mailbox by ID
func (s *MailboxService) GetByID(id int64) (*domain.Mailbox, error) {
	return s.repo.GetByID(id)
}

// Create creates a new mailbox
func (s *MailboxService) Create(userID int64, name, specialUse string) error {
	now := time.Now()
	mailbox := &domain.Mailbox{
		UserID:      userID,
		Name:        name,
		Subscribed:  true,
		SpecialUse:  specialUse,
		UIDValidity: now.Unix(),
		UIDNext:     1,
		CreatedAt:   now,
	}

	err := s.repo.Create(mailbox)
	if err != nil {
		s.logger.Error("failed to create mailbox",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.String("name", name),
		)
		return err
	}

	s.logger.Info("mailbox created",
		zap.Int64("mailbox_id", mailbox.ID),
		zap.Int64("user_id", userID),
		zap.String("name", name),
	)

	return nil
}

// Delete deletes a mailbox
func (s *MailboxService) Delete(id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("failed to delete mailbox",
			zap.Error(err),
			zap.Int64("mailbox_id", id),
		)
		return err
	}

	s.logger.Info("mailbox deleted",
		zap.Int64("mailbox_id", id),
	)

	return nil
}

// Rename renames a mailbox
func (s *MailboxService) Rename(id int64, newName string) error {
	mailbox, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	mailbox.Name = newName

	err = s.repo.Update(mailbox)
	if err != nil {
		s.logger.Error("failed to rename mailbox",
			zap.Error(err),
			zap.Int64("mailbox_id", id),
			zap.String("new_name", newName),
		)
		return err
	}

	s.logger.Info("mailbox renamed",
		zap.Int64("mailbox_id", id),
		zap.String("new_name", newName),
	)

	return nil
}

// UpdateSubscription updates mailbox subscription status
func (s *MailboxService) UpdateSubscription(id int64, subscribed bool) error {
	mailbox, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	mailbox.Subscribed = subscribed

	return s.repo.Update(mailbox)
}

// CreateDefaultMailboxes creates default mailboxes for a new user
func (s *MailboxService) CreateDefaultMailboxes(userID int64) error {
	defaults := []struct {
		Name       string
		SpecialUse string
	}{
		{"INBOX", ""},
		{"Drafts", "\\Drafts"},
		{"Sent", "\\Sent"},
		{"Trash", "\\Trash"},
		{"Spam", "\\Junk"},
		{"Archive", "\\Archive"},
	}

	for _, mb := range defaults {
		if err := s.Create(userID, mb.Name, mb.SpecialUse); err != nil {
			return err
		}
	}

	s.logger.Info("default mailboxes created",
		zap.Int64("user_id", userID),
		zap.Int("count", len(defaults)),
	)

	return nil
}
