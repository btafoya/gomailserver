package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/emersion/go-message/mail"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

const (
	// StorageThreshold is the size threshold for blob vs file storage (1MB)
	StorageThreshold = 1024 * 1024

	// MaxBlobSize is the maximum size for blob storage
	MaxBlobSize = 1024 * 1024
)

// MessageService handles message operations
type MessageService struct {
	repo        repository.MessageRepository
	logger      *zap.Logger
	storagePath string
}

// NewMessageService creates a new message service
func NewMessageService(repo repository.MessageRepository, storagePath string, logger *zap.Logger) *MessageService {
	return &MessageService{
		repo:        repo,
		logger:      logger,
		storagePath: storagePath,
	}
}

// Store stores a message with hybrid storage strategy
func (s *MessageService) Store(userID, mailboxID, uid int64, messageData []byte) (*domain.Message, error) {
	size := int64(len(messageData))

	// Parse MIME message
	reader := bytes.NewReader(messageData)
	mailReader, err := mail.CreateReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	// Extract headers
	headers := make(map[string]string)
	fields := mailReader.Header.Fields()
	for fields.Next() {
		headers[fields.Key()] = fields.Value()
	}

	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal headers: %w", err)
	}

	// Parse body structure
	bodyStructure, err := s.parseBodyStructure(mailReader)
	if err != nil {
		s.logger.Warn("failed to parse body structure", zap.Error(err))
		bodyStructure = ""
	}

	// Extract key fields
	subject := mailReader.Header.Get("Subject")
	from := mailReader.Header.Get("From")
	to := mailReader.Header.Get("To")
	cc := mailReader.Header.Get("Cc")
	bcc := mailReader.Header.Get("Bcc")
	replyTo := mailReader.Header.Get("Reply-To")
	messageID := mailReader.Header.Get("Message-ID")
	inReplyTo := mailReader.Header.Get("In-Reply-To")
	refs := mailReader.Header.Get("References")

	// Generate thread ID from message headers
	threadID := s.generateThreadID(messageID, inReplyTo)

	// Parse internal date (use current time for now)
	// TODO: Parse RFC 2822 date from Date header
	internalDate := time.Now()

	// Determine storage strategy
	var storageType string
	var content []byte
	var contentPath string

	if size < StorageThreshold {
		// Small message: store in blob
		storageType = "blob"
		content = messageData
	} else {
		// Large message: store as file
		storageType = "file"
		path, err := s.saveToFile(userID, mailboxID, uid, messageData)
		if err != nil {
			return nil, fmt.Errorf("failed to save message to file: %w", err)
		}
		contentPath = path
	}

	// Create message record
	msg := &domain.Message{
		UserID:        userID,
		MailboxID:     mailboxID,
		UID:           uint32(uid),
		Size:          size,
		Flags:         "",
		Categories:    "",
		ThreadID:      threadID,
		ReceivedAt:    time.Now(),
		InternalDate:  internalDate,
		Subject:       subject,
		From:          from,
		To:            to,
		CC:            cc,
		BCC:           bcc,
		ReplyTo:       replyTo,
		MessageID:     messageID,
		InReplyTo:     inReplyTo,
		Refs:          refs,
		Headers:       string(headersJSON),
		BodyStructure: bodyStructure,
		StorageType:   storageType,
		Content:       content,
		ContentPath:   contentPath,
	}

	// Store in database
	if err := s.repo.Create(msg); err != nil {
		// If database insert fails and we saved a file, clean it up
		if storageType == "file" && contentPath != "" {
			os.Remove(contentPath)
		}
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	s.logger.Info("message stored",
		zap.Int64("message_id", msg.ID),
		zap.Int64("user_id", userID),
		zap.Int64("mailbox_id", mailboxID),
		zap.Int64("size", size),
		zap.String("storage_type", storageType),
	)

	return msg, nil
}

// GetByID retrieves a message by ID and loads its content
func (s *MessageService) GetByID(id int64) (*domain.Message, error) {
	msg, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Load content from file if needed
	if msg.StorageType == "file" && msg.ContentPath != "" {
		content, err := os.ReadFile(msg.ContentPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read message file: %w", err)
		}
		msg.Content = content
	}

	return msg, nil
}

// GetByMailbox retrieves messages for a mailbox
func (s *MessageService) GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error) {
	return s.repo.GetByMailbox(mailboxID, offset, limit)
}

// Delete deletes a message and its file if it exists
func (s *MessageService) Delete(id int64) error {
	msg, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete file if it exists
	if msg.StorageType == "file" && msg.ContentPath != "" {
		if err := os.Remove(msg.ContentPath); err != nil && !os.IsNotExist(err) {
			s.logger.Warn("failed to delete message file",
				zap.Error(err),
				zap.String("path", msg.ContentPath),
			)
		}
	}

	return s.repo.Delete(id)
}

// saveToFile saves message content to a file
func (s *MessageService) saveToFile(userID, mailboxID, uid int64, content []byte) (string, error) {
	// Create directory structure: storagePath/userID/mailboxID/
	dir := filepath.Join(s.storagePath, fmt.Sprintf("%d", userID), fmt.Sprintf("%d", mailboxID))
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Generate filename: uid-hash.eml
	hash := sha256.Sum256(content)
	filename := fmt.Sprintf("%d-%s.eml", uid, hex.EncodeToString(hash[:8]))
	path := filepath.Join(dir, filename)

	// Write file
	if err := os.WriteFile(path, content, 0640); err != nil {
		return "", fmt.Errorf("failed to write message file: %w", err)
	}

	return path, nil
}

// parseBodyStructure parses the MIME body structure
func (s *MessageService) parseBodyStructure(mr *mail.Reader) (string, error) {
	var parts []string

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		contentType := part.Header.Get("Content-Type")
		parts = append(parts, contentType)
	}

	if len(parts) == 0 {
		return "text/plain", nil
	}

	// Create simplified body structure
	structure := map[string]interface{}{
		"parts": parts,
	}

	data, err := json.Marshal(structure)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// generateThreadID generates a thread ID from message headers
func (s *MessageService) generateThreadID(messageID, inReplyTo string) string {
	// Use In-Reply-To if available for threading
	if inReplyTo != "" {
		return s.normalizeMessageID(inReplyTo)
	}

	// Otherwise use Message-ID as thread root
	if messageID != "" {
		return s.normalizeMessageID(messageID)
	}

	// Generate random thread ID if no headers available
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return hex.EncodeToString(hash[:16])
}

// normalizeMessageID normalizes a Message-ID for use as thread ID
func (s *MessageService) normalizeMessageID(msgID string) string {
	// Remove angle brackets and whitespace
	msgID = strings.TrimSpace(msgID)
	msgID = strings.Trim(msgID, "<>")

	// Hash the message ID for consistent length
	hash := sha256.Sum256([]byte(msgID))
	return hex.EncodeToString(hash[:16])
}
