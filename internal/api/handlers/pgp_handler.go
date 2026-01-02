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

// PGPHandler handles PGP key management endpoints
type PGPHandler struct {
	service *service.PGPService
	logger  *zap.Logger
}

// NewPGPHandler creates a new PGP handler
func NewPGPHandler(service *service.PGPService, logger *zap.Logger) *PGPHandler {
	return &PGPHandler{
		service: service,
		logger:  logger,
	}
}

// PGPKeyImportRequest represents a request to import a PGP key
type PGPKeyImportRequest struct {
	UserID    int64  `json:"user_id"`
	PublicKey string `json:"public_key"`
}

// PGPKeyResponse represents a PGP key in API responses
type PGPKeyResponse struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	KeyID       string `json:"key_id"`
	Fingerprint string `json:"fingerprint"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ExpiresAt   string `json:"expires_at,omitempty"`
	IsPrimary   bool   `json:"is_primary"`
}

// ImportKey imports a PGP public key for a user
func (h *PGPHandler) ImportKey(w http.ResponseWriter, r *http.Request) {
	var req PGPKeyImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PublicKey == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Public key is required")
		return
	}

	if req.UserID == 0 {
		middleware.RespondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	key, err := h.service.ImportKey(r.Context(), req.UserID, req.PublicKey)
	if err != nil {
		h.logger.Error("failed to import PGP key",
			zap.Error(err),
			zap.Int64("user_id", req.UserID),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to import PGP key")
		return
	}

	middleware.RespondJSON(w, http.StatusCreated, h.toResponse(key))
}

// ListKeys lists all PGP keys for a user
func (h *PGPHandler) ListKeys(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	keys, err := h.service.GetUserKeys(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to list PGP keys",
			zap.Error(err),
			zap.Int64("user_id", userID),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to list PGP keys")
		return
	}

	responses := make([]PGPKeyResponse, len(keys))
	for i, key := range keys {
		responses[i] = h.toResponse(key)
	}

	middleware.RespondJSON(w, http.StatusOK, responses)
}

// GetKey gets a specific PGP key
func (h *PGPHandler) GetKey(w http.ResponseWriter, r *http.Request) {
	keyIDStr := chi.URLParam(r, "id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	key, err := h.service.GetKey(r.Context(), keyID)
	if err != nil {
		h.logger.Error("failed to get PGP key",
			zap.Error(err),
			zap.Int64("key_id", keyID),
		)
		middleware.RespondError(w, http.StatusNotFound, "PGP key not found")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, h.toResponse(key))
}

// SetPrimary sets a key as the primary key for a user
func (h *PGPHandler) SetPrimary(w http.ResponseWriter, r *http.Request) {
	keyIDStr := chi.URLParam(r, "id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	if err := h.service.SetPrimaryKey(r.Context(), keyID); err != nil {
		h.logger.Error("failed to set primary PGP key",
			zap.Error(err),
			zap.Int64("key_id", keyID),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to set primary key")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]string{"message": "Primary key updated"})
}

// DeleteKey deletes a PGP key
func (h *PGPHandler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	keyIDStr := chi.URLParam(r, "id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	if err := h.service.DeleteKey(r.Context(), keyID); err != nil {
		h.logger.Error("failed to delete PGP key",
			zap.Error(err),
			zap.Int64("key_id", keyID),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to delete PGP key")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]string{"message": "PGP key deleted"})
}

// toResponse converts a domain PGP key to API response format
func (h *PGPHandler) toResponse(key *domain.PGPKey) PGPKeyResponse {
	resp := PGPKeyResponse{
		ID:          key.ID,
		UserID:      key.UserID,
		KeyID:       key.KeyID,
		Fingerprint: key.Fingerprint,
		CreatedAt:   key.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   key.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		IsPrimary:   key.IsPrimary,
	}

	if key.ExpiresAt != nil {
		resp.ExpiresAt = key.ExpiresAt.Format("2006-01-02T15:04:05Z")
	}

	return resp
}
