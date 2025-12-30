package service

import (
	"errors"
	"time"

	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

var ErrUnauthorized = errors.New("unauthorized")

type QuarantineService struct {
	repo           repository.QuarantineRepository
	logger         *config.Logger
	messageService MessageService
}

func NewQuarantineService(repo repository.QuarantineRepository, logger *config.Logger, messageService MessageService) *QuarantineService {
	return &QuarantineService{
		repo:           repo,
		logger:         logger,
		messageService: messageService,
	}
}

func (s *QuarantineService) Quarantine(userID int64, messageID int64, score float64) error {
	item := &domain.QuarantineItem{
		MessageID:     messageID,
		UserID:        userID,
		SpamScore:     score,
		QuarantinedAt: time.Now(),
		AutoDeleteAt:  time.Now().AddDate(0, 0, 30), // 30 days
	}
	return s.repo.Create(item)
}

func (s *QuarantineService) Release(userID int64, itemID int64) error {
	item, err := s.repo.GetByID(itemID)
	if err != nil {
		return err
	}

	if item.UserID != userID {
		return ErrUnauthorized
	}

	// Move message back to inbox
	if err := s.messageService.MoveToInbox(item.MessageID); err != nil {
		return err
	}

	return s.repo.MarkReleased(itemID)
}

func (s *QuarantineService) ListForUser(userID int64) ([]*domain.QuarantineItem, error) {
	return s.repo.ListByUser(userID)
}

func (s *QuarantineService) CleanupExpired() error {
	return s.repo.DeleteExpired()
}
