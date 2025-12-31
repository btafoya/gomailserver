package service

import (
	"time"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"go.uber.org/zap"
)

type QuarantineService struct {
	repo   repository.QuarantineRepository
	logger *zap.Logger
}

func NewQuarantineService(repo repository.QuarantineRepository, logger *zap.Logger) *QuarantineService {
	return &QuarantineService{
		repo:   repo,
		logger: logger,
	}
}

func (s *QuarantineService) Quarantine(messageID, sender, recipient, subject, reason, messagePath string, score float64) error {
	message := &domain.QuarantineMessage{
		MessageID:   messageID,
		Sender:      sender,
		Recipient:   recipient,
		Subject:     subject,
		Reason:      reason,
		Score:       score,
		MessagePath: messagePath,
		Action:      "quarantined",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.repo.Create(message)
}

func (s *QuarantineService) Release(itemID int64) error {
	return s.repo.UpdateAction(itemID, "released")
}

func (s *QuarantineService) Delete(itemID int64) error {
	if err := s.repo.UpdateAction(itemID, "deleted"); err != nil {
		return err
	}
	return s.repo.Delete(itemID)
}

func (s *QuarantineService) List(offset, limit int) ([]*domain.QuarantineMessage, error) {
	return s.repo.List(offset, limit)
}

func (s *QuarantineService) CleanupOld(age time.Duration) error {
	return s.repo.DeleteOlderThan(age)
}
