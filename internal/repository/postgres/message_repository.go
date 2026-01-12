package postgres

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

func NewMessageRepository(db *database.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *domain.Message) error {
	query := `
		INSERT INTO messages (
			user_id, mailbox_id, uid, size, flags, categories, thread_id,
			received_at, internal_date, subject, from_addr, to_addr, cc_addr, bcc_addr, reply_to,
			message_id, in_reply_to, refs, headers, body_structure,
			storage_type, content, content_path, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id
	`

	err := r.db.QueryRow(query,
		message.UserID, message.MailboxID, message.UID, message.Size, message.Flags, message.Categories, message.ThreadID,
		message.ReceivedAt, message.InternalDate, message.Subject, message.From, message.To, message.CC, message.BCC, message.ReplyTo,
		message.MessageID, message.InReplyTo, message.Refs, message.Headers, message.BodyStructure,
		message.StorageType, message.Content, message.ContentPath, time.Now(),
	).Scan(&message.ID)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	message.CreatedAt = time.Now()

	return nil
}

func (r *messageRepository) GetByID(id int64) (*domain.Message, error) {
	query := `
		SELECT
			id, user_id, mailbox_id, uid, size, flags, categories, thread_id,
			received_at, internal_date, subject, from_addr, to_addr, cc_addr, bcc_addr, reply_to,
			message_id, in_reply_to, refs, headers, body_structure,
			storage_type, content, content_path, created_at
		FROM messages
		WHERE id = $1
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

func (r *messageRepository) GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error) {
	query := `
		SELECT
			id, user_id, mailbox_id, uid, size, flags, categories, thread_id,
			received_at, internal_date, subject, from_addr, to_addr, cc_addr, bcc_addr, reply_to,
			message_id, in_reply_to, refs, headers, body_structure,
			storage_type, content, content_path, created_at
		FROM messages
		WHERE mailbox_id = $1
		ORDER BY received_at DESC
		LIMIT $2 OFFSET $3
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

func (r *messageRepository) Update(message *domain.Message) error {
	query := `
		UPDATE messages SET
			flags = ?, categories = ?, thread_id = ?,
			subject = ?, from_addr = ?, to_addr = ?, cc_addr = ?, bcc_addr = ?, reply_to = ?,
			message_id = ?, in_reply_to = ?, refs = ?, headers = ?, body_structure = ?
		WHERE id = $1
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

func (r *messageRepository) Delete(id int64) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}
