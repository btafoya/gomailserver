package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// WebhookHandler handles webhook-related API requests
type WebhookHandler struct {
	service *service.WebhookService
	logger  *zap.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(service *service.WebhookService, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{
		service: service,
		logger:  logger,
	}
}

// CreateWebhook handles POST /api/v1/webhooks
func (h *WebhookHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Secret      string `json:"secret"`
		EventTypes  string `json:"event_types"`
		Active      bool   `json:"active"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" || req.URL == "" || req.Secret == "" || req.EventTypes == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Missing required fields: name, url, secret, event_types")
		return
	}

	webhook := &domain.Webhook{
		Name:        req.Name,
		URL:         req.URL,
		Secret:      req.Secret,
		EventTypes:  req.EventTypes,
		Active:      req.Active,
		Description: req.Description,
	}

	if err := h.service.Repo.Create(r.Context(), webhook); err != nil {
		h.logger.Error("failed to create webhook", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to create webhook")
		return
	}

	middleware.RespondJSON(w, http.StatusCreated, webhook)
}

// GetWebhook handles GET /api/v1/webhooks/{id}
func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid webhook ID")
		return
	}

	webhook, err := h.service.Repo.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get webhook", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get webhook")
		return
	}

	if webhook == nil {
		middleware.RespondError(w, http.StatusNotFound, "Webhook not found")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, webhook)
}

// ListWebhooks handles GET /api/v1/webhooks
func (h *WebhookHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := h.service.Repo.List(r.Context())
	if err != nil {
		h.logger.Error("failed to list webhooks", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to list webhooks")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, webhooks)
}

// UpdateWebhook handles PUT /api/v1/webhooks/{id}
func (h *WebhookHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid webhook ID")
		return
	}

	var req struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Secret      string `json:"secret"`
		EventTypes  string `json:"event_types"`
		Active      bool   `json:"active"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing webhook
	webhook, err := h.service.Repo.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get webhook", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get webhook")
		return
	}

	if webhook == nil {
		middleware.RespondError(w, http.StatusNotFound, "Webhook not found")
		return
	}

	// Update fields
	webhook.Name = req.Name
	webhook.URL = req.URL
	webhook.Secret = req.Secret
	webhook.EventTypes = req.EventTypes
	webhook.Active = req.Active
	webhook.Description = req.Description

	if err := h.service.Repo.Update(r.Context(), webhook); err != nil{
		h.logger.Error("failed to update webhook", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update webhook")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, webhook)
}

// DeleteWebhook handles DELETE /api/v1/webhooks/{id}
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid webhook ID")
		return
	}

	if err := h.service.Repo.Delete(r.Context(), id); err != nil {
		h.logger.Error("failed to delete webhook", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to delete webhook")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]string{"message": "Webhook deleted successfully"})
}

// TestWebhook handles POST /api/v1/webhooks/{id}/test
func (h *WebhookHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid webhook ID")
		return
	}

	if err := h.service.TestWebhook(r.Context(), id); err != nil {
		h.logger.Error("failed to test webhook",
			zap.Error(err),
			zap.Int64("webhook_id", id),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to test webhook: "+err.Error())
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]string{"message": "Test webhook sent successfully"})
}

// ListDeliveries handles GET /api/v1/webhooks/{id}/deliveries
func (h *WebhookHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid webhook ID")
		return
	}

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	deliveries, err := h.service.Repo.ListDeliveriesByWebhook(r.Context(), id, limit)
	if err != nil {
		h.logger.Error("failed to list deliveries", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to list deliveries")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, deliveries)
}

// GetDelivery handles GET /api/v1/webhooks/deliveries/{id}
func (h *WebhookHandler) GetDelivery(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid delivery ID")
		return
	}

	delivery, err := h.service.Repo.GetDeliveryByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get delivery", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get delivery")
		return
	}

	if delivery == nil {
		middleware.RespondError(w, http.StatusNotFound, "Delivery not found")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, delivery)
}
