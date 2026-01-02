package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type webhookRepository struct {
	db *database.DB
}

// NewWebhookRepository creates a new SQLite webhook repository
func NewWebhookRepository(db *database.DB) repository.WebhookRepository {
	return &webhookRepository{db: db}
}

// Create inserts a new webhook
func (r *webhookRepository) Create(ctx context.Context, webhook *domain.Webhook) error {
	query := `
		INSERT INTO webhooks (name, url, secret, event_types, active, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		webhook.Name,
		webhook.URL,
		webhook.Secret,
		webhook.EventTypes,
		webhook.Active,
		webhook.Description,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get webhook ID: %w", err)
	}

	webhook.ID = id
	webhook.CreatedAt = now
	webhook.UpdatedAt = now

	return nil
}

// GetByID retrieves a webhook by ID
func (r *webhookRepository) GetByID(ctx context.Context, id int64) (*domain.Webhook, error) {
	query := `
		SELECT id, name, url, secret, event_types, active, description, created_at, updated_at
		FROM webhooks
		WHERE id = ?
	`

	webhook := &domain.Webhook{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&webhook.ID,
		&webhook.Name,
		&webhook.URL,
		&webhook.Secret,
		&webhook.EventTypes,
		&webhook.Active,
		&webhook.Description,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	return webhook, nil
}

// List retrieves all webhooks
func (r *webhookRepository) List(ctx context.Context) ([]*domain.Webhook, error) {
	query := `
		SELECT id, name, url, secret, event_types, active, description, created_at, updated_at
		FROM webhooks
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []*domain.Webhook
	for rows.Next() {
		webhook := &domain.Webhook{}
		if err := rows.Scan(
			&webhook.ID,
			&webhook.Name,
			&webhook.URL,
			&webhook.Secret,
			&webhook.EventTypes,
			&webhook.Active,
			&webhook.Description,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan webhook: %w", err)
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// ListActive retrieves all active webhooks
func (r *webhookRepository) ListActive(ctx context.Context) ([]*domain.Webhook, error) {
	query := `
		SELECT id, name, url, secret, event_types, active, description, created_at, updated_at
		FROM webhooks
		WHERE active = 1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []*domain.Webhook
	for rows.Next() {
		webhook := &domain.Webhook{}
		if err := rows.Scan(
			&webhook.ID,
			&webhook.Name,
			&webhook.URL,
			&webhook.Secret,
			&webhook.EventTypes,
			&webhook.Active,
			&webhook.Description,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan webhook: %w", err)
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// Update updates a webhook
func (r *webhookRepository) Update(ctx context.Context, webhook *domain.Webhook) error {
	query := `
		UPDATE webhooks
		SET name = ?, url = ?, secret = ?, event_types = ?, active = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	webhook.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		webhook.Name,
		webhook.URL,
		webhook.Secret,
		webhook.EventTypes,
		webhook.Active,
		webhook.Description,
		webhook.UpdatedAt,
		webhook.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	return nil
}

// Delete deletes a webhook
func (r *webhookRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM webhooks WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}

// CreateDelivery creates a webhook delivery record
func (r *webhookRepository) CreateDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error {
	query := `
		INSERT INTO webhook_deliveries (
			webhook_id, event_type, payload, attempt_count, max_attempts, status,
			status_code, response_body, error_message, next_retry_at,
			first_attempted_at, last_attempted_at, completed_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		delivery.WebhookID,
		delivery.EventType,
		delivery.Payload,
		delivery.AttemptCount,
		delivery.MaxAttempts,
		delivery.Status,
		delivery.StatusCode,
		delivery.ResponseBody,
		delivery.ErrorMessage,
		delivery.NextRetryAt,
		delivery.FirstAttemptedAt,
		delivery.LastAttemptedAt,
		delivery.CompletedAt,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create delivery: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get delivery ID: %w", err)
	}

	delivery.ID = id
	delivery.CreatedAt = time.Now()

	return nil
}

// GetDeliveryByID retrieves a delivery by ID
func (r *webhookRepository) GetDeliveryByID(ctx context.Context, id int64) (*domain.WebhookDelivery, error) {
	query := `
		SELECT id, webhook_id, event_type, payload, attempt_count, max_attempts, status,
			status_code, response_body, error_message, next_retry_at,
			first_attempted_at, last_attempted_at, completed_at, created_at
		FROM webhook_deliveries
		WHERE id = ?
	`

	delivery := &domain.WebhookDelivery{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&delivery.ID,
		&delivery.WebhookID,
		&delivery.EventType,
		&delivery.Payload,
		&delivery.AttemptCount,
		&delivery.MaxAttempts,
		&delivery.Status,
		&delivery.StatusCode,
		&delivery.ResponseBody,
		&delivery.ErrorMessage,
		&delivery.NextRetryAt,
		&delivery.FirstAttemptedAt,
		&delivery.LastAttemptedAt,
		&delivery.CompletedAt,
		&delivery.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	return delivery, nil
}

// ListDeliveriesByWebhook retrieves deliveries for a webhook
func (r *webhookRepository) ListDeliveriesByWebhook(ctx context.Context, webhookID int64, limit int) ([]*domain.WebhookDelivery, error) {
	query := `
		SELECT id, webhook_id, event_type, payload, attempt_count, max_attempts, status,
			status_code, response_body, error_message, next_retry_at,
			first_attempted_at, last_attempted_at, completed_at, created_at
		FROM webhook_deliveries
		WHERE webhook_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, webhookID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list deliveries: %w", err)
	}
	defer rows.Close()

	return r.scanDeliveries(rows)
}

// ListPendingDeliveries retrieves pending deliveries
func (r *webhookRepository) ListPendingDeliveries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error) {
	query := `
		SELECT id, webhook_id, event_type, payload, attempt_count, max_attempts, status,
			status_code, response_body, error_message, next_retry_at,
			first_attempted_at, last_attempted_at, completed_at, created_at
		FROM webhook_deliveries
		WHERE status IN ('pending', 'retrying')
			AND (next_retry_at IS NULL OR next_retry_at <= ?)
			AND attempt_count < max_attempts
		ORDER BY created_at ASC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, time.Now(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending deliveries: %w", err)
	}
	defer rows.Close()

	return r.scanDeliveries(rows)
}

// UpdateDelivery updates a delivery record
func (r *webhookRepository) UpdateDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error {
	query := `
		UPDATE webhook_deliveries
		SET attempt_count = ?, status = ?, status_code = ?, response_body = ?,
			error_message = ?, next_retry_at = ?, first_attempted_at = ?,
			last_attempted_at = ?, completed_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		delivery.AttemptCount,
		delivery.Status,
		delivery.StatusCode,
		delivery.ResponseBody,
		delivery.ErrorMessage,
		delivery.NextRetryAt,
		delivery.FirstAttemptedAt,
		delivery.LastAttemptedAt,
		delivery.CompletedAt,
		delivery.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update delivery: %w", err)
	}

	return nil
}

// DeleteOldDeliveries deletes deliveries older than the specified number of days
func (r *webhookRepository) DeleteOldDeliveries(ctx context.Context, olderThan int) error {
	query := `
		DELETE FROM webhook_deliveries
		WHERE created_at < datetime('now', '-' || ? || ' days')
	`

	_, err := r.db.ExecContext(ctx, query, olderThan)
	if err != nil {
		return fmt.Errorf("failed to delete old deliveries: %w", err)
	}

	return nil
}

// scanDeliveries is a helper to scan delivery rows
func (r *webhookRepository) scanDeliveries(rows *sql.Rows) ([]*domain.WebhookDelivery, error) {
	var deliveries []*domain.WebhookDelivery
	for rows.Next() {
		delivery := &domain.WebhookDelivery{}
		if err := rows.Scan(
			&delivery.ID,
			&delivery.WebhookID,
			&delivery.EventType,
			&delivery.Payload,
			&delivery.AttemptCount,
			&delivery.MaxAttempts,
			&delivery.Status,
			&delivery.StatusCode,
			&delivery.ResponseBody,
			&delivery.ErrorMessage,
			&delivery.NextRetryAt,
			&delivery.FirstAttemptedAt,
			&delivery.LastAttemptedAt,
			&delivery.CompletedAt,
			&delivery.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan delivery: %w", err)
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}
