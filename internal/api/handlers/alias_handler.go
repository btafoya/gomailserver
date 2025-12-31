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

// AliasHandler handles alias management endpoints
type AliasHandler struct {
	service *service.AliasService
	logger  *zap.Logger
}

// NewAliasHandler creates a new alias handler
func NewAliasHandler(service *service.AliasService, logger *zap.Logger) *AliasHandler {
	return &AliasHandler{
		service: service,
		logger:  logger,
	}
}

// AliasRequest represents an alias creation request
type AliasRequest struct {
	Address      string   `json:"address"`
	Destinations []string `json:"destinations"`
	DomainID     int64    `json:"domain_id"`
	Status       string   `json:"status,omitempty"`
}

// AliasResponse represents an alias in API responses
type AliasResponse struct {
	ID           int64    `json:"id"`
	Address      string   `json:"address"`
	Destinations []string `json:"destinations"`
	DomainID     int64    `json:"domain_id"`
	DomainName   string   `json:"domain_name,omitempty"`
	Status       string   `json:"status"`
	CreatedAt    string   `json:"created_at"`
}

// List retrieves all aliases
func (h *AliasHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Add pagination and filtering support
	aliases, err := h.service.ListAll(r.Context())
	if err != nil {
		h.logger.Error("Failed to list aliases", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve aliases")
		return
	}

	// Convert to response format
	responses := make([]*AliasResponse, len(aliases))
	for i, a := range aliases {
		responses[i] = aliasToResponse(a)
	}

	middleware.RespondSuccess(w, responses, "Aliases retrieved successfully")
}

// Create creates a new alias
func (h *AliasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req AliasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Address == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Alias address is required")
		return
	}
	if len(req.Destinations) == 0 {
		middleware.RespondError(w, http.StatusBadRequest, "At least one destination is required")
		return
	}
	if req.DomainID == 0 {
		middleware.RespondError(w, http.StatusBadRequest, "Domain ID is required")
		return
	}

	// Convert request to alias model
	newAlias := &domain.Alias{
		Address:      req.Address,
		Destinations: req.Destinations,
		DomainID:     req.DomainID,
		Status:       req.Status,
	}

	// Set defaults
	if newAlias.Status == "" {
		newAlias.Status = "active"
	}

	// Create alias
	err := h.service.Create(r.Context(), newAlias)
	if err != nil {
		h.logger.Error("Failed to create alias",
			zap.String("address", req.Address),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to create alias")
		return
	}

	h.logger.Info("Alias created",
		zap.String("address", newAlias.Address),
		zap.Int64("id", newAlias.ID),
	)

	middleware.RespondCreated(w, aliasToResponse(newAlias), "Alias created successfully")
}

// Get retrieves a specific alias
func (h *AliasHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid alias ID")
		return
	}

	alias, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get alias", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Alias not found")
		return
	}

	middleware.RespondSuccess(w, aliasToResponse(alias), "Alias retrieved successfully")
}

// Delete deletes an alias
func (h *AliasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid alias ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete alias", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to delete alias")
		return
	}

	h.logger.Info("Alias deleted", zap.Int64("id", id))

	middleware.RespondNoContent(w)
}

// aliasToResponse converts an alias model to API response format
func aliasToResponse(a *domain.Alias) *AliasResponse {
	return &AliasResponse{
		ID:           a.ID,
		Address:      a.Address,
		Destinations: a.Destinations,
		DomainID:     a.DomainID,
		Status:       a.Status,
		CreatedAt:    a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
