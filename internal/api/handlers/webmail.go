package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// WebmailHandler handles webmail-related HTTP requests
type WebmailHandler struct {
	mailboxService *service.MailboxService
	messageService *service.MessageService
	logger         *zap.Logger
}

// NewWebmailHandler creates a new webmail handler
func NewWebmailHandler(
	mailboxService *service.MailboxService,
	messageService *service.MessageService,
	logger *zap.Logger,
) *WebmailHandler {
	return &WebmailHandler{
		mailboxService: mailboxService,
		messageService: messageService,
		logger:         logger,
	}
}

// ListMailboxes handles GET /api/v1/webmail/mailboxes
func (h *WebmailHandler) ListMailboxes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	mailboxes, err := h.mailboxService.ListMailboxesByUserID(ctx, int(userID))
	if err != nil {
		h.logger.Error("failed to list mailboxes", zap.Error(err), zap.Int64("user_id", userID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to list mailboxes")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"mailboxes": mailboxes,
	})
}

// ListMessages handles GET /api/v1/webmail/mailboxes/:id/messages
func (h *WebmailHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	mailboxID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid mailbox ID")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	messages, err := h.messageService.ListMessages(ctx, int(mailboxID), int(userID), limit, offset)
	if err != nil {
		h.logger.Error("failed to list messages", zap.Error(err), zap.Int64("mailbox_id", mailboxID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to list messages")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"messages": messages,
		"page":     page,
		"limit":    limit,
	})
}

// GetMessage handles GET /api/v1/webmail/messages/:id
func (h *WebmailHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	messageID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid message ID")
		return
	}

	message, err := h.messageService.GetMessage(ctx, int(messageID), int(userID))
	if err != nil {
		h.logger.Error("failed to get message", zap.Error(err), zap.Int64("message_id", messageID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get message")
		return
	}

	if message == nil {
		middleware.RespondError(w, http.StatusNotFound, "message not found")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, message)
}

// SendMessage handles POST /api/v1/webmail/messages
func (h *WebmailHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req struct {
		To          string   `json:"to"`
		Cc          string   `json:"cc"`
		Bcc         string   `json:"bcc"`
		Subject     string   `json:"subject"`
		BodyText    string   `json:"body_text"`
		BodyHTML    string   `json:"body_html"`
		Attachments []string `json:"attachments"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validation
	if req.To == "" {
		middleware.RespondError(w, http.StatusBadRequest, "recipient required")
		return
	}

	if req.Subject == "" {
		middleware.RespondError(w, http.StatusBadRequest, "subject required")
		return
	}

	messageID, err := h.messageService.SendMessage(ctx, int(userID), &service.SendMessageRequest{
		To:          req.To,
		Cc:          req.Cc,
		Bcc:         req.Bcc,
		Subject:     req.Subject,
		BodyText:    req.BodyText,
		BodyHTML:    req.BodyHTML,
		Attachments: req.Attachments,
	})

	if err != nil {
		h.logger.Error("failed to send message", zap.Error(err), zap.Int64("user_id", userID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to send message")
		return
	}

	middleware.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message_id": messageID,
		"status":     "queued",
	})
}

// DeleteMessage handles DELETE /api/v1/webmail/messages/:id
func (h *WebmailHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	messageID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid message ID")
		return
	}

	err = h.messageService.DeleteMessage(ctx, int(messageID), int(userID))
	if err != nil {
		h.logger.Error("failed to delete message", zap.Error(err), zap.Int64("message_id", messageID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to delete message")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "message deleted successfully",
	})
}

// MoveMessage handles POST /api/v1/webmail/messages/:id/move
func (h *WebmailHandler) MoveMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	messageID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid message ID")
		return
	}

	var req struct {
		MailboxID int `json:"mailbox_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.messageService.MoveMessage(ctx, int(messageID), req.MailboxID, int(userID))
	if err != nil {
		h.logger.Error("failed to move message", zap.Error(err), zap.Int64("message_id", messageID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to move message")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "message moved successfully",
	})
}

// UpdateFlags handles POST /api/v1/webmail/messages/:id/flags
func (h *WebmailHandler) UpdateFlags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	messageID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid message ID")
		return
	}

	var req struct {
		Flags  []string `json:"flags"`
		Action string   `json:"action"` // "add" or "remove"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.messageService.UpdateFlags(ctx, int(messageID), int(userID), req.Flags, req.Action)
	if err != nil {
		h.logger.Error("failed to update flags", zap.Error(err), zap.Int64("message_id", messageID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to update flags")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "flags updated successfully",
	})
}

// SearchMessages handles GET /api/v1/webmail/search
func (h *WebmailHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		middleware.RespondError(w, http.StatusBadRequest, "search query required")
		return
	}

	messages, err := h.messageService.SearchMessages(ctx, int(userID), query)
	if err != nil {
		h.logger.Error("failed to search messages", zap.Error(err), zap.String("query", query))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to search messages")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"messages": messages,
		"query":    query,
	})
}

// DownloadAttachment handles GET /api/v1/webmail/attachments/:id
func (h *WebmailHandler) DownloadAttachment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	attachmentID := chi.URLParam(r, "id")
	if attachmentID == "" {
		middleware.RespondError(w, http.StatusBadRequest, "invalid attachment ID")
		return
	}

	attachment, err := h.messageService.GetAttachment(ctx, attachmentID, int(userID))
	if err != nil {
		h.logger.Error("failed to get attachment", zap.Error(err), zap.String("attachment_id", attachmentID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get attachment")
		return
	}

	if attachment == nil {
		middleware.RespondError(w, http.StatusNotFound, "attachment not found")
		return
	}

	w.Header().Set("Content-Type", attachment.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+attachment.Filename+"\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(attachment.Data)))
	w.Write(attachment.Data)
}

// SaveDraft handles POST /api/v1/webmail/drafts
func (h *WebmailHandler) SaveDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req struct {
		DraftID     *int     `json:"draft_id,omitempty"`
		To          []string `json:"to"`
		Cc          []string `json:"cc,omitempty"`
		Bcc         []string `json:"bcc,omitempty"`
		Subject     string   `json:"subject"`
		BodyHTML    string   `json:"body_html,omitempty"`
		BodyText    string   `json:"body_text,omitempty"`
		InReplyTo   string   `json:"in_reply_to,omitempty"`
		References  string   `json:"references,omitempty"`
		Attachments []string `json:"attachments,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	draft, err := h.messageService.SaveDraft(ctx, int(userID), req.DraftID, &service.DraftData{
		To:          req.To,
		Cc:          req.Cc,
		Bcc:         req.Bcc,
		Subject:     req.Subject,
		BodyHTML:    req.BodyHTML,
		BodyText:    req.BodyText,
		InReplyTo:   req.InReplyTo,
		References:  req.References,
		Attachments: req.Attachments,
	})

	if err != nil {
		h.logger.Error("failed to save draft", zap.Error(err), zap.Int64("user_id", userID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to save draft")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, draft)
}

// ListDrafts handles GET /api/v1/webmail/drafts
func (h *WebmailHandler) ListDrafts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	drafts, err := h.messageService.ListDrafts(ctx, int(userID))
	if err != nil {
		h.logger.Error("failed to list drafts", zap.Error(err), zap.Int64("user_id", userID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to list drafts")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"drafts": drafts,
	})
}

// GetDraft handles GET /api/v1/webmail/drafts/:id
func (h *WebmailHandler) GetDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	draftID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid draft ID")
		return
	}

	draft, err := h.messageService.GetDraft(ctx, int(draftID), int(userID))
	if err != nil {
		h.logger.Error("failed to get draft", zap.Error(err), zap.Int64("draft_id", draftID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get draft")
		return
	}

	if draft == nil {
		middleware.RespondError(w, http.StatusNotFound, "draft not found")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, draft)
}

// DeleteDraft handles DELETE /api/v1/webmail/drafts/:id
func (h *WebmailHandler) DeleteDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	draftID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid draft ID")
		return
	}

	if err := h.messageService.DeleteDraft(ctx, int(draftID), int(userID)); err != nil {
		h.logger.Error("failed to delete draft", zap.Error(err), zap.Int64("draft_id", draftID))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to delete draft")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
	})
}
