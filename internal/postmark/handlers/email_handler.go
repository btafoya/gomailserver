package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/btafoya/gomailserver/internal/postmark/middleware"
	"github.com/btafoya/gomailserver/internal/postmark/models"
	"github.com/btafoya/gomailserver/internal/postmark/service"
	"go.uber.org/zap"
)

// EmailHandler handles email sending endpoints
type EmailHandler struct {
	emailService *service.EmailService
	logger       *zap.Logger
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(emailService *service.EmailService, logger *zap.Logger) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
		logger:       logger,
	}
}

// Send handles POST /email
func (h *EmailHandler) Send(w http.ResponseWriter, r *http.Request) {
	var req models.EmailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, models.ErrorCodeInvalidJSON, models.MsgInvalidJSON)
		return
	}

	serverID := middleware.GetServerID(r)

	resp, err := h.emailService.SendEmail(r.Context(), serverID, &req)
	if err != nil {
		if pmErr, ok := err.(*models.PostmarkError); ok {
			models.WriteError(w, pmErr.ErrorCode, pmErr.Message)
		} else {
			models.WriteError(w, models.ErrorCodeInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SendBatch handles POST /email/batch
func (h *EmailHandler) SendBatch(w http.ResponseWriter, r *http.Request) {
	var requests models.BatchEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		models.WriteError(w, models.ErrorCodeInvalidJSON, models.MsgInvalidJSON)
		return
	}

	serverID := middleware.GetServerID(r)

	responses, err := h.emailService.SendBatchEmail(r.Context(), serverID, requests)
	if err != nil {
		if pmErr, ok := err.(*models.PostmarkError); ok {
			models.WriteError(w, pmErr.ErrorCode, pmErr.Message)
		} else {
			models.WriteError(w, models.ErrorCodeInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
