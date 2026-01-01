package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/postmark/models"
	"github.com/btafoya/gomailserver/internal/postmark/repository"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EmailService handles PostmarkApp email operations
type EmailService struct {
	repo         repository.PostmarkRepository
	queueService *service.QueueService
	logger       *zap.Logger
}

// NewEmailService creates a new email service
func NewEmailService(repo repository.PostmarkRepository, queueService *service.QueueService, logger *zap.Logger) *EmailService {
	return &EmailService{
		repo:         repo,
		queueService: queueService,
		logger:       logger,
	}
}

// SendEmail sends a single email
func (s *EmailService) SendEmail(ctx context.Context, serverID int, req *models.EmailRequest) (*models.EmailResponse, error) {
	// Validate request
	if err := s.validateEmailRequest(req); err != nil {
		return nil, err
	}

	// Generate message ID
	messageID := uuid.New().String()
	submittedAt := time.Now()

	// Build MIME message
	message, err := s.buildMIMEMessage(req, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to build message: %w", err)
	}

	// Get recipient list
	recipients := s.parseRecipients(req.To, req.Cc, req.Bcc)

	// Enqueue message for delivery
	_, err = s.queueService.Enqueue(req.From, recipients, message)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue message: %w", err)
	}

	// Store message in PostmarkApp tracking
	metadataJSON, _ := json.Marshal(req.Metadata)
	err = s.repo.CreateMessage(ctx, &repository.Message{
		MessageID:     messageID,
		ServerID:      serverID,
		FromEmail:     req.From,
		ToEmail:       req.To,
		CcEmail:       req.Cc,
		BccEmail:      req.Bcc,
		Subject:       req.Subject,
		HtmlBody:      req.HtmlBody,
		TextBody:      req.TextBody,
		Tag:           req.Tag,
		Metadata:      string(metadataJSON),
		MessageStream: req.MessageStream,
		Status:        "pending",
		SubmittedAt:   submittedAt.Format(time.RFC3339),
	})

	if err != nil {
		s.logger.Error("failed to store message tracking", zap.Error(err))
	}

	return &models.EmailResponse{
		To:          req.To,
		SubmittedAt: submittedAt,
		MessageID:   messageID,
		ErrorCode:   models.ErrorCodeSuccess,
		Message:     models.MsgSuccess,
	}, nil
}

// SendBatchEmail sends multiple emails
func (s *EmailService) SendBatchEmail(ctx context.Context, serverID int, requests models.BatchEmailRequest) (models.BatchEmailResponse, error) {
	if len(requests) == 0 {
		return nil, &models.PostmarkError{
			ErrorCode: models.ErrorCodeInvalidEmail,
			Message:   "Batch request is empty",
		}
	}

	if len(requests) > 500 {
		return nil, &models.PostmarkError{
			ErrorCode: models.ErrorCodeBatchLimitExceeded,
			Message:   models.MsgBatchLimitExceeded,
		}
	}

	responses := make(models.BatchEmailResponse, 0, len(requests))

	for _, req := range requests {
		resp, err := s.SendEmail(ctx, serverID, &req)
		if err != nil {
			// Convert error to PostmarkError format
			if pmErr, ok := err.(*models.PostmarkError); ok {
				responses = append(responses, models.EmailResponse{
					To:        req.To,
					ErrorCode: pmErr.ErrorCode,
					Message:   pmErr.Message,
				})
			} else {
				responses = append(responses, models.EmailResponse{
					To:        req.To,
					ErrorCode: models.ErrorCodeInternalServerError,
					Message:   err.Error(),
				})
			}
		} else {
			responses = append(responses, *resp)
		}
	}

	return responses, nil
}

// validateEmailRequest validates an email request
func (s *EmailService) validateEmailRequest(req *models.EmailRequest) error {
	if req.From == "" {
		return &models.PostmarkError{
			ErrorCode: models.ErrorCodeInvalidEmail,
			Message:   "The 'From' address is required",
		}
	}

	if req.To == "" {
		return &models.PostmarkError{
			ErrorCode: models.ErrorCodeInvalidEmail,
			Message:   "The 'To' address is required",
		}
	}

	// Count total recipients
	toCount := len(strings.Split(req.To, ","))
	ccCount := 0
	if req.Cc != "" {
		ccCount = len(strings.Split(req.Cc, ","))
	}
	bccCount := 0
	if req.Bcc != "" {
		bccCount = len(strings.Split(req.Bcc, ","))
	}

	if toCount+ccCount+bccCount > 50 {
		return &models.PostmarkError{
			ErrorCode: models.ErrorCodeInvalidEmail,
			Message:   "Maximum 50 recipients allowed (To + Cc + Bcc)",
		}
	}

	if req.Subject == "" && (req.HtmlBody == "" && req.TextBody == "") {
		return &models.PostmarkError{
			ErrorCode: models.ErrorCodeInvalidEmail,
			Message:   "Subject and Body are required",
		}
	}

	return nil
}

// parseRecipients extracts recipient addresses
func (s *EmailService) parseRecipients(to, cc, bcc string) []string {
	var recipients []string

	for _, addr := range strings.Split(to, ",") {
		if addr = strings.TrimSpace(addr); addr != "" {
			recipients = append(recipients, addr)
		}
	}

	for _, addr := range strings.Split(cc, ",") {
		if addr = strings.TrimSpace(addr); addr != "" {
			recipients = append(recipients, addr)
		}
	}

	for _, addr := range strings.Split(bcc, ",") {
		if addr = strings.TrimSpace(addr); addr != "" {
			recipients = append(recipients, addr)
		}
	}

	return recipients
}

// buildMIMEMessage builds a MIME message from the request
func (s *EmailService) buildMIMEMessage(req *models.EmailRequest, messageID string) ([]byte, error) {
	var buf bytes.Buffer

	// Write headers
	from, _ := mail.ParseAddress(req.From)
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from.String()))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", req.To))

	if req.Cc != "" {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", req.Cc))
	}

	if req.ReplyTo != "" {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", req.ReplyTo))
	}

	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", req.Subject))
	buf.WriteString(fmt.Sprintf("Message-ID: <%s@postmark>\r\n", messageID))
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString("MIME-Version: 1.0\r\n")

	// Add custom headers
	for _, h := range req.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", h.Name, h.Value))
	}

	// Add PostmarkApp metadata
	if req.Tag != "" {
		buf.WriteString(fmt.Sprintf("X-PM-Tag: %s\r\n", req.Tag))
	}

	if req.MessageStream != "" {
		buf.WriteString(fmt.Sprintf("X-PM-Message-Stream: %s\r\n", req.MessageStream))
	}

	// Build multipart message if needed
	hasAttachments := len(req.Attachments) > 0
	hasHTML := req.HtmlBody != ""
	hasText := req.TextBody != ""

	if hasAttachments || (hasHTML && hasText) {
		// Multipart message
		boundary := fmt.Sprintf("boundary_%s", messageID)

		if hasAttachments {
			buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary))
		} else {
			buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundary))
		}

		// Text part
		if hasText {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
			buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
			buf.WriteString(req.TextBody)
			buf.WriteString("\r\n\r\n")
		}

		// HTML part
		if hasHTML {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			buf.WriteString("Content-Type: text/html; charset=utf-8\r\n")
			buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
			buf.WriteString(req.HtmlBody)
			buf.WriteString("\r\n\r\n")
		}

		// Attachments
		for _, att := range req.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", att.ContentType, att.Name))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", att.Name))

			if att.ContentID != "" {
				buf.WriteString(fmt.Sprintf("Content-ID: <%s>\r\n", att.ContentID))
			}

			buf.WriteString("\r\n")
			buf.WriteString(att.Content)
			buf.WriteString("\r\n\r\n")
		}

		buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// Simple message
		if hasHTML {
			buf.WriteString("Content-Type: text/html; charset=utf-8\r\n")
			buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
			buf.WriteString(req.HtmlBody)
		} else {
			buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
			buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
			buf.WriteString(req.TextBody)
		}
	}

	return buf.Bytes(), nil
}
