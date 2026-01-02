package service

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-message/mail"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
)

// Webmail-specific service methods for MessageService and MailboxService

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	From        string
	To          string
	Cc          string
	Bcc         string
	Subject     string
	BodyText    string
	BodyHTML    string
	Attachments []string
}

// DraftData represents draft message data
type DraftData struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	BodyHTML    string
	BodyText    string
	InReplyTo   string
	References  string
	Attachments []string
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// Draft represents a draft message
type Draft struct {
	ID        int
	UserID    int
	Data      *DraftData
	UpdatedAt string
}

// ListMailboxesByUserID lists all mailboxes for a user (alias for List method)
func (s *MailboxService) ListMailboxesByUserID(ctx context.Context, userID int) ([]*domain.Mailbox, error) {
	return s.List(int64(userID), false)
}

// ListMessages lists messages in a mailbox with pagination
func (s *MessageService) ListMessages(ctx context.Context, mailboxID, userID, limit, offset int) ([]*domain.Message, error) {
	// Get messages from mailbox
	messages, err := s.GetByMailbox(int64(mailboxID), offset, limit)
	if err != nil {
		return nil, err
	}

	// Validate user ownership for each message
	validMessages := make([]*domain.Message, 0, len(messages))
	for _, msg := range messages {
		if msg.UserID == int64(userID) {
			validMessages = append(validMessages, msg)
		} else {
			s.logger.Warn("message ownership mismatch",
				zap.Int64("message_id", msg.ID),
				zap.Int64("expected_user", int64(userID)),
				zap.Int64("actual_user", msg.UserID),
			)
		}
	}

	return validMessages, nil
}

// GetMessage retrieves a single message by ID with user ownership check
func (s *MessageService) GetMessage(ctx context.Context, messageID, userID int) (*domain.Message, error) {
	msg, err := s.GetByID(int64(messageID))
	if err != nil {
		return nil, err
	}

	// Verify user owns this message
	if msg.UserID != int64(userID) {
		return nil, fmt.Errorf("access denied: message does not belong to user")
	}

	return msg, nil
}

// SendMessage sends a new message via the queue
func (s *MessageService) SendMessage(ctx context.Context, userID int, req *SendMessageRequest) (int, error) {
	// Build MIME message
	var buf strings.Builder

	// Write headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", req.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", req.To))
	if req.Cc != "" {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", req.Cc))
	}
	if req.Bcc != "" {
		buf.WriteString(fmt.Sprintf("Bcc: %s\r\n", req.Bcc))
	}
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", req.Subject))
	buf.WriteString("MIME-Version: 1.0\r\n")

	// Simple message format (text/plain or text/html)
	if req.BodyHTML != "" {
		buf.WriteString("Content-Type: text/html; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(req.BodyHTML)
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(req.BodyText)
	}

	messageData := []byte(buf.String())

	// Build recipient list for SMTP
	var recipients []string
	if req.To != "" {
		recipients = append(recipients, req.To)
	}
	if req.Cc != "" {
		recipients = append(recipients, strings.Split(req.Cc, ",")...)
	}
	if req.Bcc != "" {
		recipients = append(recipients, strings.Split(req.Bcc, ",")...)
	}

	// Queue message for SMTP delivery if QueueService is available
	if s.queueService != nil {
		_, err := s.queueService.Enqueue(req.From, recipients, messageData)
		if err != nil {
			return 0, fmt.Errorf("failed to queue message for delivery: %w", err)
		}
	}

	// Store copy in Sent folder if MailboxService is available
	var sentMessageID int64
	if s.mailboxService != nil {
		sentMailbox, err := s.mailboxService.GetByName(int64(userID), "Sent")
		if err != nil {
			s.logger.Warn("failed to get Sent mailbox, message queued but not stored in Sent",
				zap.Error(err),
				zap.Int("user_id", userID),
			)
		} else {
			// Store message in Sent folder
			msg, err := s.Store(int64(userID), sentMailbox.ID, 0, messageData)
			if err != nil {
				s.logger.Warn("failed to store message in Sent folder",
					zap.Error(err),
					zap.Int("user_id", userID),
				)
			} else {
				sentMessageID = msg.ID
			}
		}
	}

	return int(sentMessageID), nil
}

// DeleteMessage moves a message to trash or deletes permanently
func (s *MessageService) DeleteMessage(ctx context.Context, messageID, userID int) error {
	msg, err := s.GetByID(int64(messageID))
	if err != nil {
		return err
	}

	// Verify user owns this message
	if msg.UserID != int64(userID) {
		return fmt.Errorf("access denied: message does not belong to user")
	}

	// Try to move to Trash folder if MailboxService is available
	if s.mailboxService != nil {
		trashMailbox, err := s.mailboxService.GetByName(int64(userID), "Trash")
		if err != nil {
			// If Trash mailbox doesn't exist, fall back to hard delete
			s.logger.Warn("Trash mailbox not found, performing hard delete",
				zap.Error(err),
				zap.Int("user_id", userID),
			)
		} else {
			// Move message to Trash folder
			msg.MailboxID = trashMailbox.ID
			if err := s.repo.Update(msg); err != nil {
				return fmt.Errorf("failed to move message to Trash: %w", err)
			}
			return nil
		}
	}

	// Hard delete if no MailboxService or Trash folder doesn't exist
	return s.Delete(int64(messageID))
}

// MoveMessage moves a message to a different mailbox
func (s *MessageService) MoveMessage(ctx context.Context, messageID, targetMailboxID, userID int) error {
	msg, err := s.GetByID(int64(messageID))
	if err != nil {
		return err
	}

	// Verify user owns this message
	if msg.UserID != int64(userID) {
		return fmt.Errorf("access denied: message does not belong to user")
	}

	// Update mailbox ID
	msg.MailboxID = int64(targetMailboxID)

	// Update in repository
	return s.repo.Update(msg)
}

// UpdateFlags updates message flags (read, starred, etc)
func (s *MessageService) UpdateFlags(ctx context.Context, messageID, userID int, flags []string, action string) error {
	msg, err := s.GetByID(int64(messageID))
	if err != nil {
		return err
	}

	// Verify user owns this message
	if msg.UserID != int64(userID) {
		return fmt.Errorf("access denied: message does not belong to user")
	}

	// Parse current flags
	currentFlags := strings.Split(msg.Flags, " ")
	flagMap := make(map[string]bool)
	for _, f := range currentFlags {
		if f != "" {
			flagMap[f] = true
		}
	}

	// Add or remove flags
	if action == "add" {
		for _, f := range flags {
			flagMap[f] = true
		}
	} else if action == "remove" {
		for _, f := range flags {
			delete(flagMap, f)
		}
	}

	// Rebuild flags string
	newFlags := []string{}
	for f := range flagMap {
		newFlags = append(newFlags, f)
	}
	msg.Flags = strings.Join(newFlags, " ")

	// Update in repository
	return s.repo.Update(msg)
}

// SearchMessages searches messages for a user
func (s *MessageService) SearchMessages(ctx context.Context, userID int, query string) ([]*domain.Message, error) {
	if s.mailboxService == nil {
		return nil, fmt.Errorf("MailboxService not available for search")
	}

	// Get all mailboxes for the user
	mailboxes, err := s.mailboxService.List(int64(userID), false)
	if err != nil {
		return nil, fmt.Errorf("failed to get user mailboxes: %w", err)
	}

	// Search across all mailboxes
	query = strings.ToLower(query)
	var results []*domain.Message

	for _, mailbox := range mailboxes {
		// Get messages from this mailbox (limit to 100 per mailbox for performance)
		messages, err := s.repo.GetByMailbox(mailbox.ID, 0, 100)
		if err != nil {
			s.logger.Warn("failed to get messages from mailbox during search",
				zap.Error(err),
				zap.Int64("mailbox_id", mailbox.ID),
			)
			continue
		}

		// Filter messages that match the query
		for _, msg := range messages {
			if msg.UserID != int64(userID) {
				continue
			}

			// Search in subject, from, to, and body
			if strings.Contains(strings.ToLower(msg.Subject), query) ||
				strings.Contains(strings.ToLower(msg.From), query) ||
				strings.Contains(strings.ToLower(msg.To), query) {
				results = append(results, msg)
			}
		}
	}

	return results, nil
}

// GetAttachment retrieves an attachment from a message
func (s *MessageService) GetAttachment(ctx context.Context, attachmentID string, userID int) (*Attachment, error) {
	// Parse attachmentID format: "messageID:partIndex"
	parts := strings.Split(attachmentID, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid attachment ID format, expected messageID:partIndex")
	}

	messageID := parts[0]
	partIndex := parts[1]

	// Convert messageID to int64
	msgID, err := fmt.Sscanf(messageID, "%d", new(int64))
	if err != nil || msgID != 1 {
		return nil, fmt.Errorf("invalid message ID in attachment ID")
	}

	var msgIDInt int64
	fmt.Sscanf(messageID, "%d", &msgIDInt)

	// Get message and verify ownership
	msg, err := s.GetMessage(ctx, int(msgIDInt), userID)
	if err != nil {
		return nil, err
	}

	// Parse MIME message to extract attachment
	reader := strings.NewReader(string(msg.Content))
	mailReader, err := mail.CreateReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	// Iterate through parts to find the requested part
	currentIndex := 0
	for {
		part, err := mailReader.NextPart()
		if err != nil {
			break
		}

		// Check if this is an attachment
		contentDisposition := part.Header.Get("Content-Disposition")
		if !strings.HasPrefix(contentDisposition, "attachment") {
			continue
		}

		// Check if this is the requested part
		if fmt.Sprintf("%d", currentIndex) == partIndex {
			// Extract attachment data
			data, err := io.ReadAll(part.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read attachment: %w", err)
			}

			// Get filename from Content-Disposition header
			_, params, err := mime.ParseMediaType(contentDisposition)
			filename := "attachment"
			if err == nil {
				if name, ok := params["filename"]; ok {
					filename = name
				}
			}

			// Get content type
			contentType := part.Header.Get("Content-Type")
			if contentType == "" {
				contentType = "application/octet-stream"
			}

			return &Attachment{
				Filename:    filename,
				ContentType: contentType,
				Data:        data,
			}, nil
		}

		currentIndex++
	}

	return nil, fmt.Errorf("attachment not found")
}

// SaveDraft saves or updates a draft message
func (s *MessageService) SaveDraft(ctx context.Context, userID int, draftID *int, data *DraftData) (*Draft, error) {
	if s.mailboxService == nil {
		return nil, fmt.Errorf("MailboxService not available")
	}

	// Build draft MIME message
	var buf strings.Builder

	// Write headers
	if len(data.To) > 0 {
		buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(data.To, ", ")))
	}
	if len(data.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(data.Cc, ", ")))
	}
	if len(data.Bcc) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(data.Bcc, ", ")))
	}
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", data.Subject))
	if data.InReplyTo != "" {
		buf.WriteString(fmt.Sprintf("In-Reply-To: %s\r\n", data.InReplyTo))
	}
	if data.References != "" {
		buf.WriteString(fmt.Sprintf("References: %s\r\n", data.References))
	}
	buf.WriteString("MIME-Version: 1.0\r\n")

	// Write body
	if data.BodyHTML != "" {
		buf.WriteString("Content-Type: text/html; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(data.BodyHTML)
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(data.BodyText)
	}

	messageData := []byte(buf.String())

	// Get user's Drafts mailbox
	draftsMailbox, err := s.mailboxService.GetByName(int64(userID), "Drafts")
	if err != nil {
		return nil, fmt.Errorf("failed to get Drafts mailbox: %w", err)
	}

	var msg *domain.Message

	if draftID != nil && *draftID > 0 {
		// Update existing draft
		msg, err = s.GetByID(int64(*draftID))
		if err != nil {
			return nil, fmt.Errorf("failed to get existing draft: %w", err)
		}

		// Verify ownership
		if msg.UserID != int64(userID) {
			return nil, fmt.Errorf("access denied: draft does not belong to user")
		}

		// Delete old message file if it exists
		if msg.StorageType == "file" && msg.ContentPath != "" {
			_ = os.Remove(msg.ContentPath)
		}

		// Store new version
		newMsg, err := s.Store(int64(userID), draftsMailbox.ID, int64(msg.UID), messageData)
		if err != nil {
			return nil, fmt.Errorf("failed to update draft: %w", err)
		}

		// Delete old message record
		_ = s.repo.Delete(msg.ID)

		msg = newMsg
	} else {
		// Create new draft
		msg, err = s.Store(int64(userID), draftsMailbox.ID, 0, messageData)
		if err != nil {
			return nil, fmt.Errorf("failed to store draft: %w", err)
		}
	}

	// Return Draft object
	return &Draft{
		ID:        int(msg.ID),
		UserID:    userID,
		Data:      data,
		UpdatedAt: msg.ReceivedAt.Format(time.RFC3339),
	}, nil
}

// ListDrafts lists all drafts for a user
func (s *MessageService) ListDrafts(ctx context.Context, userID int) ([]*Draft, error) {
	if s.mailboxService == nil {
		return nil, fmt.Errorf("MailboxService not available")
	}

	// Get user's Drafts mailbox
	draftsMailbox, err := s.mailboxService.GetByName(int64(userID), "Drafts")
	if err != nil {
		return nil, fmt.Errorf("failed to get Drafts mailbox: %w", err)
	}

	// Get all messages in Drafts folder
	messages, err := s.repo.GetByMailbox(draftsMailbox.ID, 0, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft messages: %w", err)
	}

	// Convert to Draft objects
	drafts := make([]*Draft, 0, len(messages))
	for _, msg := range messages {
		draftData := &DraftData{
			Subject:    msg.Subject,
			InReplyTo:  msg.InReplyTo,
			References: msg.Refs,
		}

		// Parse To, Cc, Bcc from headers
		if msg.To != "" {
			draftData.To = strings.Split(msg.To, ",")
		}
		if msg.CC != "" {
			draftData.Cc = strings.Split(msg.CC, ",")
		}
		if msg.BCC != "" {
			draftData.Bcc = strings.Split(msg.BCC, ",")
		}

		draft := &Draft{
			ID:        int(msg.ID),
			UserID:    userID,
			Data:      draftData,
			UpdatedAt: msg.ReceivedAt.Format(time.RFC3339),
		}

		drafts = append(drafts, draft)
	}

	return drafts, nil
}

// GetDraft retrieves a specific draft
func (s *MessageService) GetDraft(ctx context.Context, draftID, userID int) (*Draft, error) {
	// Get draft message and verify ownership
	msg, err := s.GetMessage(ctx, draftID, userID)
	if err != nil {
		return nil, err
	}

	// Verify message is in Drafts folder if MailboxService is available
	if s.mailboxService != nil {
		draftsMailbox, err := s.mailboxService.GetByName(int64(userID), "Drafts")
		if err == nil && msg.MailboxID != draftsMailbox.ID {
			return nil, fmt.Errorf("message is not in Drafts folder")
		}
	}

	// Load full message content if needed
	var content []byte
	if msg.StorageType == "file" && msg.ContentPath != "" {
		content, err = os.ReadFile(msg.ContentPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read message content: %w", err)
		}
	} else {
		content = msg.Content
	}

	// Parse MIME message to extract draft data
	draftData := &DraftData{
		Subject:    msg.Subject,
		InReplyTo:  msg.InReplyTo,
		References: msg.Refs,
	}

	// Parse To, Cc, Bcc from message headers
	if msg.To != "" {
		draftData.To = strings.Split(msg.To, ",")
	}
	if msg.CC != "" {
		draftData.Cc = strings.Split(msg.CC, ",")
	}
	if msg.BCC != "" {
		draftData.Bcc = strings.Split(msg.BCC, ",")
	}

	// Parse MIME to extract body
	if len(content) > 0 {
		reader := strings.NewReader(string(content))
		mailReader, err := mail.CreateReader(reader)
		if err == nil {
			// Read body parts
			for {
				part, err := mailReader.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					break
				}

				contentType := part.Header.Get("Content-Type")
				body, err := io.ReadAll(part.Body)
				if err != nil {
					continue
				}

				if strings.Contains(contentType, "text/html") {
					draftData.BodyHTML = string(body)
				} else if strings.Contains(contentType, "text/plain") {
					draftData.BodyText = string(body)
				}
			}
		}
	}

	draft := &Draft{
		ID:        int(msg.ID),
		UserID:    int(msg.UserID),
		Data:      draftData,
		UpdatedAt: msg.ReceivedAt.Format(time.RFC3339),
	}

	return draft, nil
}

// DeleteDraft deletes a draft
func (s *MessageService) DeleteDraft(ctx context.Context, draftID, userID int) error {
	// Verify it's a draft before deleting
	msg, err := s.GetMessage(ctx, draftID, userID)
	if err != nil {
		return err
	}

	if !strings.Contains(msg.Flags, "\\Draft") {
		return fmt.Errorf("message is not a draft")
	}

	return s.DeleteMessage(ctx, draftID, userID)
}
