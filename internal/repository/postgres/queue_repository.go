package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type queueRepository struct {
	db *database.DB
}

// NewQueueRepository creates a new PostgreSQL queue repository
func NewQueueRepository(db *database.DB) repository.QueueRepository {
	return &queueRepository{db: db}
}

// Enqueue inserts a new message into the queue
func (r *queueRepository) Enqueue(item *domain.QueueItem) error {
	query := `
		INSERT INTO smtp_queue (
			sender, recipients, message_id, message_path,
			retry_count, max_retries, next_retry, status,
			error_message, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	err := r.db.QueryRow(query,
		item.Sender, item.Recipients, item.MessageID, item.MessagePath,
		item.RetryCount, item.MaxRetries, item.NextRetry, item.Status,
		item.ErrorMessage, time.Now(), time.Now(),
	).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("failed to enqueue message: %w", err)
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	return nil
}

// GetPending retrieves all pending messages ready for delivery
func (r *queueRepository) GetPending() ([]*domain.QueueItem, error) {
	query := `
		SELECT
			id, sender, recipients, message_id, message_path,
			retry_count, max_retries, next_retry, status,
			error_message, created_at, updated_at
		FROM smtp_queue
		WHERE status = 'pending'
		  AND (next_retry IS NULL OR next_retry <= $1)
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get pending queue items: %w", err)
	}
	defer rows.Close()

	items := make([]*domain.QueueItem, 0)
	for rows.Next() {
		item := &domain.QueueItem{}
		var nextRetry sql.NullTime

		err := rows.Scan(
			&item.ID, &item.Sender, &item.Recipients, &item.MessageID, &item.MessagePath,
			&item.RetryCount, &item.MaxRetries, &nextRetry, &item.Status,
			&item.ErrorMessage, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan queue item: %w", err)
		}

		if nextRetry.Valid {
			item.NextRetry = &nextRetry.Time
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

// GetByID retrieves a queue item by ID
func (r *queueRepository) GetByID(id int64) (*domain.QueueItem, error) {
	query := `
		SELECT
			id, sender, recipients, message_id, message_path,
			retry_count, max_retries, next_retry, status,
			error_message, created_at, updated_at
		FROM smtp_queue
		WHERE id = $1
	`

	item := &domain.QueueItem{}
	var nextRetry sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Sender, &item.Recipients, &item.MessageID, &item.MessagePath,
		&item.RetryCount, &item.MaxRetries, &nextRetry, &item.Status,
		&item.ErrorMessage, &item.CreatedAt, &item.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("queue item not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get queue item: %w", err)
	}

	if nextRetry.Valid {
		item.NextRetry = &nextRetry.Time
	}

	return item, nil
}

// UpdateStatus updates the status and error message of a queue item
func (r *queueRepository) UpdateStatus(id int64, status string, errorMsg string) error {
	query := `
		UPDATE smtp_queue SET
			status = $1,
			error_message = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, status, errorMsg, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update queue item status: %w", err)
	}

	return nil
}

// UpdateRetry updates the retry count and next retry time
func (r *queueRepository) UpdateRetry(id int64, retryCount int, nextRetry time.Time) error {
	query := `
		UPDATE smtp_queue SET
			retry_count = $1,
			next_retry = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, retryCount, nextRetry, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update queue item retry: %w", err)
	}

	return nil
}

// Delete deletes a queue item
func (r *queueRepository) Delete(id int64) error {
	query := `DELETE FROM smtp_queue WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete queue item: %w", err)
	}
	return nil
}
