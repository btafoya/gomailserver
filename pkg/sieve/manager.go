package sieve

import (
	"context"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/service"
	"go.uber.org/zap"
)

// Manager handles Sieve filtering operations
type Manager struct {
	parser         *Parser
	messageService service.MessageServiceInterface
	logger         *zap.Logger
}

// NewManager creates a new Sieve manager
func NewManager(logger *zap.Logger, messageService service.MessageServiceInterface) *Manager {
	return &Manager{
		parser:         NewParser(logger),
		messageService: messageService,
		logger:         logger,
	}
}

// ApplyFilters applies Sieve filters to an incoming message
func (m *Manager) ApplyFilters(ctx context.Context, userID int64, message *domain.Message) (*domain.Message, error) {
	m.logger.Debug("applying Sieve filters",
		zap.Int64("user_id", userID),
		zap.String("message_id", message.MessageID),
	)

	return message, nil
}
