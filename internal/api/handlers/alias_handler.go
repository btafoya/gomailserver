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

// AliasListResponse represents a paginated list of aliases
type AliasListResponse struct {
	Aliases    []*AliasResponse `json:"aliases"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
	TotalCount int64            `json:"total_count"`
}

// List retrieves aliases with pagination and optional filtering
func (h *AliasHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	domainIDStr := r.URL.Query().Get("domain_id")

	// Set defaults
	page := 1
	pageSize := 50

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	offset := (page - 1) * pageSize

	var aliases []*domain.Alias
	var totalCount int64
	var err error

	// Filter by domain if specified
	if domainIDStr != "" {
		domainID, parseErr := strconv.ParseInt(domainIDStr, 10, 64)
		if parseErr != nil {
			middleware.RespondError(w, http.StatusBadRequest, "Invalid domain_id")
			return
		}
		aliases, err = h.service.ListByDomain(r.Context(), domainID)
		if err != nil {
			h.logger.Error("Failed to list aliases by domain", zap.Error(err))
			middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve aliases")
			return
		}
		totalCount, err = h.service.CountByDomain(r.Context(), domainID)
	} else {
		aliases, err = h.service.List(r.Context(), offset, pageSize)
		if err != nil {
			h.logger.Error("Failed to list aliases", zap.Error(err))
			middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve aliases")
			return
		}
		totalCount, err = h.service.Count(r.Context())
	}

	if err != nil {
		h.logger.Error("Failed to count aliases", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to count aliases")
		return
	}

	// Convert to response format
	responses := make([]*AliasResponse, len(aliases))
	for i, a := range aliases {
		responses[i] = aliasToResponse(a)
	}

	// Calculate total pages
	totalPages := int(totalCount / int64(pageSize))
	if totalCount%int64(pageSize) > 0 {
		totalPages++
	}

	response := AliasListResponse{
		Aliases:    responses,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}

	middleware.RespondSuccess(w, response, "Aliases retrieved successfully")
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

	// Convert destinations to JSON string
	destinationsJSON, err := service.SetDestinations(req.Destinations)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid destinations format")
		return
	}

	// Convert request to alias model
	newAlias := &domain.Alias{
		AliasEmail:        req.Address,
		DestinationEmails: destinationsJSON,
		DomainID:          req.DomainID,
		Status:            req.Status,
	}

	// Set defaults
	if newAlias.Status == "" {
		newAlias.Status = "active"
	}

	// Create alias
	err = h.service.Create(r.Context(), newAlias)
	if err != nil {
		h.logger.Error("Failed to create alias",
			zap.String("address", req.Address),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to create alias")
		return
	}

	h.logger.Info("Alias created",
		zap.String("address", newAlias.AliasEmail),
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
	// Parse destinations from JSON
	destinations, err := service.GetDestinations(a.DestinationEmails)
	if err != nil {
		destinations = []string{} // Fallback to empty array on error
	}

	return &AliasResponse{
		ID:           a.ID,
		Address:      a.AliasEmail,
		Destinations: destinations,
		DomainID:     a.DomainID,
		Status:       a.Status,
		CreatedAt:    a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
