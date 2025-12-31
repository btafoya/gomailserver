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

// QueueHandler handles queue management endpoints
type QueueHandler struct {
	service *service.QueueService
	logger  *zap.Logger
}

// NewQueueHandler creates a new queue handler
func NewQueueHandler(service *service.QueueService, logger *zap.Logger) *QueueHandler {
	return &QueueHandler{
		service: service,
		logger:  logger,
	}
}

// QueueItemResponse represents a queued message in API responses
type QueueItemResponse struct {
	ID           int64    `json:"id"`
	Sender       string   `json:"sender"`
	Recipients   []string `json:"recipients"`
	MessageID    string   `json:"message_id"`
	MessagePath  string   `json:"message_path"`
	Status       string   `json:"status"`
	RetryCount   int      `json:"retry_count"`
	MaxRetries   int      `json:"max_retries"`
	NextRetry    string   `json:"next_retry,omitempty"`
	ErrorMessage string   `json:"error_message,omitempty"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// List retrieves all queued messages
func (h *QueueHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	status := r.URL.Query().Get("status")

	// TODO: Add pagination support
	var items []*domain.QueueItem
	var err error

	if status != "" {
		// TODO: Filter by status
		items, err = h.service.GetPendingItems(r.Context())
	} else {
		// Get all queue items
		items, err = h.service.GetPendingItems(r.Context())
	}

	if err != nil {
		h.logger.Error("Failed to list queue items", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve queue items")
		return
	}

	// Convert to response format
	responses := make([]*QueueItemResponse, len(items))
	for i, item := range items {
		responses[i] = queueItemToResponse(item)
	}

	middleware.RespondSuccess(w, responses, "Queue items retrieved successfully")
}

// Get retrieves a specific queue item
func (h *QueueHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid queue item ID")
		return
	}

	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get queue item", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Queue item not found")
		return
	}

	middleware.RespondSuccess(w, queueItemToResponse(item), "Queue item retrieved successfully")
}

// Retry manually retries a failed queue item
func (h *QueueHandler) Retry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid queue item ID")
		return
	}

	// Get the queue item
	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, http.StatusNotFound, "Queue item not found")
		return
	}

	// Reset retry count and status
	err = h.service.RetryItem(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to retry queue item", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retry queue item")
		return
	}

	h.logger.Info("Queue item retry requested",
		zap.Int64("id", id),
		zap.String("sender", item.Sender),
	)

	middleware.RespondSuccess(w, nil, "Queue item scheduled for retry")
}

// Delete removes a queue item
func (h *QueueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid queue item ID")
		return
	}

	err = h.service.DeleteItem(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete queue item", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to delete queue item")
		return
	}

	h.logger.Info("Queue item deleted", zap.Int64("id", id))

	middleware.RespondNoContent(w)
}

// queueItemToResponse converts a queue item to API response format
func queueItemToResponse(item *domain.QueueItem) *QueueItemResponse {
	// Parse recipients from JSON
	var recipients []string
	if item.Recipients != "" {
		_ = json.Unmarshal([]byte(item.Recipients), &recipients)
	}

	response := &QueueItemResponse{
		ID:           item.ID,
		Sender:       item.Sender,
		Recipients:   recipients,
		MessageID:    item.MessageID,
		MessagePath:  item.MessagePath,
		Status:       item.Status,
		RetryCount:   item.RetryCount,
		MaxRetries:   item.MaxRetries,
		ErrorMessage: item.ErrorMessage,
		CreatedAt:    item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if item.NextRetry != nil {
		response.NextRetry = item.NextRetry.Format("2006-01-02T15:04:05Z07:00")
	}

	return response
}
