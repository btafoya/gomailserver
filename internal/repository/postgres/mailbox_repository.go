package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type mailboxRepository struct {
	db *database.DB
}

// NewMailboxRepository creates a new PostgreSQL mailbox repository
func NewMailboxRepository(db *database.DB) repository.MailboxRepository {
	return &mailboxRepository{db: db}
}

// Create inserts a new mailbox
func (r *mailboxRepository) Create(mailbox *domain.Mailbox) error {
	query := `
		INSERT INTO mailboxes (
			user_id, name, parent_id, subscribed, special_use,
			uidvalidity, uidnext, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.db.QueryRow(query,
		mailbox.UserID, mailbox.Name, mailbox.ParentID, mailbox.Subscribed, mailbox.SpecialUse,
		mailbox.UIDValidity, mailbox.UIDNext, time.Now(),
	).Scan(&mailbox.ID)
	if err != nil {
		return fmt.Errorf("failed to create mailbox: %w", err)
	}

	mailbox.CreatedAt = time.Now()

	return nil
}

// GetByID retrieves a mailbox by ID
func (r *mailboxRepository) GetByID(id int64) (*domain.Mailbox, error) {
	query := `
		SELECT
			id, user_id, name, parent_id, subscribed, special_use,
			uidvalidity, uidnext, created_at
		FROM mailboxes
		WHERE id = $1
	`

	mailbox := &domain.Mailbox{}
	var parentID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&mailbox.ID, &mailbox.UserID, &mailbox.Name, &parentID, &mailbox.Subscribed, &mailbox.SpecialUse,
		&mailbox.UIDValidity, &mailbox.UIDNext, &mailbox.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("mailbox not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mailbox: %w", err)
	}

	if parentID.Valid {
		mailbox.ParentID = &parentID.Int64
	}

	return mailbox, nil
}

// GetByUser retrieves all mailboxes for a user
func (r *mailboxRepository) GetByUser(userID int64) ([]*domain.Mailbox, error) {
	query := `
		SELECT
			id, user_id, name, parent_id, subscribed, special_use,
			uidvalidity, uidnext, created_at
		FROM mailboxes
		WHERE user_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list mailboxes: %w", err)
	}
	defer rows.Close()

	mailboxes := make([]*domain.Mailbox, 0)
	for rows.Next() {
		mailbox := &domain.Mailbox{}
		var parentID sql.NullInt64

		err := rows.Scan(
			&mailbox.ID, &mailbox.UserID, &mailbox.Name, &parentID, &mailbox.Subscribed, &mailbox.SpecialUse,
			&mailbox.UIDValidity, &mailbox.UIDNext, &mailbox.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mailbox: %w", err)
		}

		if parentID.Valid {
			mailbox.ParentID = &parentID.Int64
		}

		mailboxes = append(mailboxes, mailbox)
	}

	return mailboxes, rows.Err()
}

// GetByName retrieves a mailbox by name
func (r *mailboxRepository) GetByName(userID int64, name string) (*domain.Mailbox, error) {
	query := `
		SELECT
			id, user_id, name, parent_id, subscribed, special_use,
			uidvalidity, uidnext, created_at
		FROM mailboxes
		WHERE user_id = $1 AND name = $2
	`

	mailbox := &domain.Mailbox{}
	var parentID sql.NullInt64

	err := r.db.QueryRow(query, userID, name).Scan(
		&mailbox.ID, &mailbox.UserID, &mailbox.Name, &parentID, &mailbox.Subscribed, &mailbox.SpecialUse,
		&mailbox.UIDValidity, &mailbox.UIDNext, &mailbox.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("mailbox not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mailbox: %w", err)
	}

	if parentID.Valid {
		mailbox.ParentID = &parentID.Int64
	}

	return mailbox, nil
}

// Update updates a mailbox
func (r *mailboxRepository) Update(mailbox *domain.Mailbox) error {
	query := `
		UPDATE mailboxes SET
			name = $1, parent_id = $2, subscribed = $3, special_use = $4,
			uidvalidity = $5, uidnext = $6
		WHERE id = $7
	`

	_, err := r.db.Exec(query,
		mailbox.Name, mailbox.ParentID, mailbox.Subscribed, mailbox.SpecialUse,
		mailbox.UIDValidity, mailbox.UIDNext, mailbox.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update mailbox: %w", err)
	}

	return nil
}

// Delete deletes a mailbox
func (r *mailboxRepository) Delete(id int64) error {
	query := `DELETE FROM mailboxes WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete mailbox: %w", err)
	}
	return nil
}
