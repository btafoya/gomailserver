package postgres

import (
	"database/sql"
)

type webhookRepository struct {
	db *database.DB
}

func NewWebhookRepository(db *database.DB) repository.WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) Create(webhook *domain.Webhook) error {
	panic("postgres repository not implemented yet")
}

func (r *webhookRepository) GetByID(id int64) (*domain.Webhook, error) {
	panic("postgres repository not implemented yet")
}

func (r *webhookRepository) List(offset, limit int) ([]*domain.Webhook, error) {
	panic("postgres repository not implemented yet")
}

func (r *webhookRepository) Update(webhook *domain.Webhook) error {
	panic("postgres repository not implemented yet")
}

func (r *webhookRepository) Delete(id int64) error {
	panic("postgres repository not implemented yet")
}

func (r *webhookRepository) GetActive() ([]*domain.Webhook, error) {
	panic("postgres repository not implemented yet")
}
