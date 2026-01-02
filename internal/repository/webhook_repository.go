package repository

import (
	"context"
	"github.com/btafoya/gomailserver/internal/domain"
)

// WebhookRepository defines the interface for webhook data access
type WebhookRepository interface {
	// Webhook management
	Create(ctx context.Context, webhook *domain.Webhook) error
	GetByID(ctx context.Context, id int64) (*domain.Webhook, error)
	List(ctx context.Context) ([]*domain.Webhook, error)
	ListActive(ctx context.Context) ([]*domain.Webhook, error)
	Update(ctx context.Context, webhook *domain.Webhook) error
	Delete(ctx context.Context, id int64) error

	// Delivery management
	CreateDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error
	GetDeliveryByID(ctx context.Context, id int64) (*domain.WebhookDelivery, error)
	ListDeliveriesByWebhook(ctx context.Context, webhookID int64, limit int) ([]*domain.WebhookDelivery, error)
	ListPendingDeliveries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error)
	UpdateDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error
	DeleteOldDeliveries(ctx context.Context, olderThan int) error
}
