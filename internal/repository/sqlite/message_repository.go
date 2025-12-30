package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type messageRepository struct {
	db *database.DB
}

// NewMessageRepository creates a new SQLite message repository
func NewMessageRepository(db *database.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

// Create inserts a new message
func (r *messageRepository) Create(message *domain.Message) error {
	query := `
		INSERT INTO messages (
			user_id, mailbox_id, uid, size, flags, categories, thread_id,
			received_at, internal_date, subject, from_addr, to_addr, cc_addr, bcc_addr, reply_to,
			message_id, in_reply_to, refs, headers, body_structure,
			storage_type, content, content_path, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		message.UserID, message.MailboxID, message.UID, message.Size, message.Flags, message.Categories, message.ThreadID,
		message.ReceivedAt, message.InternalDate, message.Subject, message.From, message.To, message.CC, message.BCC, message.ReplyTo,
		message.MessageID, message.InReplyTo, message.Refs, message.Headers, message.BodyStructure,
		message.StorageType, message.Content, message.ContentPath, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get message ID: %w", err)
	}

	message.ID = id
	message.CreatedAt = time.Now()

	return nil
}

// GetByID retrieves a message by ID
func (r *messageRepository) GetByID(id int64) (*domain.Message, error) {
	query := `
		SELECT
			id, user_id, mailbox_id, uid, size, flags, categories, thread_id,
			received_at, internal_date, subject, from_addr, to_addr, cc_addr, bcc_addr, reply_to,
			message_id, in_reply_to, refs, headers, body_structure,
			storage_type, content, content_path, created_at
		FROM messages
		WHERE id = ?
	`

	message := &domain.Message{}

	err := r.db.QueryRow(query, id).Scan(
		&message.ID, &message.UserID, &message.MailboxID, &message.UID, &message.Size, &message.Flags, &message.Categories, &message.ThreadID,
		&message.ReceivedAt, &message.InternalDate, &message.Subject, &message.From, &message.To, &message.CC, &message.BCC, &message.ReplyTo,
		&message.MessageID, &message.InReplyTo, &message.Refs, &message.Headers, &message.BodyStructure,
		&message.StorageType, &message.Content, &message.ContentPath, &message.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("message not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}

// GetByMailbox retrieves messages for a mailbox with pagination
func (r *messageRepository) GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error) {
	query := `
		SELECT
			id, user_id, mailbox_id, uid, size, flags, categories, thread_id,
			received_at, internal_date, subject, from_addr, to_addr, cc_addr, bcc_addr, reply_to,
			message_id, in_reply_to, refs, headers, body_structure,
			storage_type, content, content_path, created_at
		FROM messages
		WHERE mailbox_id = ?
		ORDER BY received_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, mailboxID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	messages := make([]*domain.Message, 0)
	for rows.Next() {
		message := &domain.Message{}

		err := rows.Scan(
			&message.ID, &message.UserID, &message.MailboxID, &message.UID, &message.Size, &message.Flags, &message.Categories, &message.ThreadID,
			&message.ReceivedAt, &message.InternalDate, &message.Subject, &message.From, &message.To, &message.CC, &message.BCC, &message.ReplyTo,
			&message.MessageID, &message.InReplyTo, &message.Refs, &message.Headers, &message.BodyStructure,
			&message.StorageType, &message.Content, &message.ContentPath, &message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		messages = append(messages, message)
	}

	return messages, rows.Err()
}

// Update updates a message
func (r *messageRepository) Update(message *domain.Message) error {
	query := `
		UPDATE messages SET
			flags = ?, categories = ?, thread_id = ?,
			subject = ?, from_addr = ?, to_addr = ?, cc_addr = ?, bcc_addr = ?, reply_to = ?,
			message_id = ?, in_reply_to = ?, refs = ?, headers = ?, body_structure = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		message.Flags, message.Categories, message.ThreadID,
		message.Subject, message.From, message.To, message.CC, message.BCC, message.ReplyTo,
		message.MessageID, message.InReplyTo, message.Refs, message.Headers, message.BodyStructure,
		message.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}

// Delete deletes a message
func (r *messageRepository) Delete(id int64) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}
