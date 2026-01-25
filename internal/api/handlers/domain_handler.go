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

// DomainHandler handles domain management endpoints
type DomainHandler struct {
	service *service.DomainService
	logger  *zap.Logger
}

// NewDomainHandler creates a new domain handler
func NewDomainHandler(service *service.DomainService, logger *zap.Logger) *DomainHandler {
	return &DomainHandler{
		service: service,
		logger:  logger,
	}
}

// DomainRequest represents a domain creation/update request
type DomainRequest struct {
	Name               string `json:"name"`
	Status             string `json:"status"`
	CatchallEmail      string `json:"catchall_email,omitempty"`
	MaxUsers           int    `json:"max_users,omitempty"`
	DefaultQuota       int64  `json:"default_quota,omitempty"`
	DKIMSelector       string `json:"dkim_selector,omitempty"`
	DKIMPrivateKey     string `json:"dkim_private_key,omitempty"`
	DKIMPublicKey      string `json:"dkim_public_key,omitempty"`
	SPFRecord          string `json:"spf_record,omitempty"`
	DMARCPolicy        string `json:"dmarc_policy,omitempty"`
	DMARCReportEmail   string `json:"dmarc_report_email,omitempty"`
	DKIMSigningEnabled bool   `json:"dkim_signing_enabled"`
	DKIMVerifyEnabled  bool   `json:"dkim_verify_enabled"`
}

// DomainResponse represents a domain in API responses
type DomainResponse struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	Status             string `json:"status"`
	CatchallEmail      string `json:"catchall_email,omitempty"`
	MaxUsers           int    `json:"max_users"`
	DefaultQuota       int64  `json:"default_quota"`
	DKIMSelector       string `json:"dkim_selector,omitempty"`
	DKIMPublicKey      string `json:"dkim_public_key,omitempty"`
	SPFRecord          string `json:"spf_record,omitempty"`
	DMARCPolicy        string `json:"dmarc_policy,omitempty"`
	DMARCReportEmail   string `json:"dmarc_report_email,omitempty"`
	DKIMSigningEnabled bool   `json:"dkim_signing_enabled"`
	DKIMVerifyEnabled  bool   `json:"dkim_verify_enabled"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// List retrieves all domains
func (h *DomainHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 10000 {
		limit = 100
	}

	offset := (page - 1) * limit

	domains, err := h.service.List(r.Context(), offset, limit)
	if err != nil {
		h.logger.Error("Failed to list domains", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve domains")
		return
	}

	// Convert to response format
	responses := make([]*DomainResponse, len(domains))
	for i, d := range domains {
		responses[i] = domainToResponse(d)
	}

	middleware.RespondSuccess(w, responses, "Domains retrieved successfully")
}

// Create creates a new domain
func (h *DomainHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req DomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Name == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain name is required")
		return
	}

	// Convert request to domain model
	newDomain := &domain.Domain{
		Name:               req.Name,
		Status:             req.Status,
		MaxUsers:           req.MaxUsers,
		DefaultQuota:       req.DefaultQuota,
		DKIMSelector:       req.DKIMSelector,
		DKIMPrivateKey:     req.DKIMPrivateKey,
		DKIMPublicKey:      req.DKIMPublicKey,
		SPFRecord:          req.SPFRecord,
		DMARCPolicy:        req.DMARCPolicy,
		DMARCReportEmail:   req.DMARCReportEmail,
		DKIMSigningEnabled: req.DKIMSigningEnabled,
		DKIMVerifyEnabled:  req.DKIMVerifyEnabled,
	}

	// Convert CatchallEmail string to *string
	if req.CatchallEmail != "" {
		newDomain.CatchallEmail = &req.CatchallEmail
	}

	// Set defaults
	if newDomain.Status == "" {
		newDomain.Status = "active"
	}

	// Create domain
	err := h.service.Create(r.Context(), newDomain)
	if err != nil {
		h.logger.Error("Failed to create domain",
			zap.String("domain", req.Name),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to create domain")
		return
	}

	h.logger.Info("Domain created",
		zap.String("domain", newDomain.Name),
		zap.Int64("id", newDomain.ID),
	)

	middleware.RespondCreated(w, domainToResponse(newDomain), "Domain created successfully")
}

// Get retrieves a specific domain
func (h *DomainHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	domain, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get domain", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "Domain not found")
		return
	}

	middleware.RespondSuccess(w, domainToResponse(domain), "Domain retrieved successfully")
}

// Update updates a domain
func (h *DomainHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	var req DomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing domain
	existingDomain, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, http.StatusNotFound, "Domain not found")
		return
	}

	// Update fields
	if req.Name != "" {
		existingDomain.Name = req.Name
	}
	if req.Status != "" {
		existingDomain.Status = req.Status
	}

	// Update CatchallEmail (convert string to *string)
	if req.CatchallEmail != "" {
		existingDomain.CatchallEmail = &req.CatchallEmail
	} else {
		existingDomain.CatchallEmail = nil
	}

	if req.MaxUsers > 0 {
		existingDomain.MaxUsers = req.MaxUsers
	}
	if req.DefaultQuota > 0 {
		existingDomain.DefaultQuota = req.DefaultQuota
	}
	if req.DKIMSelector != "" {
		existingDomain.DKIMSelector = req.DKIMSelector
	}
	if req.DKIMPrivateKey != "" {
		existingDomain.DKIMPrivateKey = req.DKIMPrivateKey
	}
	if req.DKIMPublicKey != "" {
		existingDomain.DKIMPublicKey = req.DKIMPublicKey
	}
	if req.SPFRecord != "" {
		existingDomain.SPFRecord = req.SPFRecord
	}
	if req.DMARCPolicy != "" {
		existingDomain.DMARCPolicy = req.DMARCPolicy
	}
	if req.DMARCReportEmail != "" {
		existingDomain.DMARCReportEmail = req.DMARCReportEmail
	}
	existingDomain.DKIMSigningEnabled = req.DKIMSigningEnabled
	existingDomain.DKIMVerifyEnabled = req.DKIMVerifyEnabled

	// Update domain
	err = h.service.Update(r.Context(), existingDomain)
	if err != nil {
		h.logger.Error("Failed to update domain", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update domain")
		return
	}

	h.logger.Info("Domain updated",
		zap.Int64("id", id),
		zap.String("domain", existingDomain.Name),
	)

	middleware.RespondSuccess(w, domainToResponse(existingDomain), "Domain updated successfully")
}

// Delete deletes a domain
func (h *DomainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete domain", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to delete domain")
		return
	}

	h.logger.Info("Domain deleted", zap.Int64("id", id))

	middleware.RespondNoContent(w)
}

// DKIMGenerateRequest represents a DKIM key generation request
type DKIMGenerateRequest struct {
	KeyType string `json:"key_type,omitempty"` // "rsa" or "ed25519", defaults to "rsa"
	KeySize int    `json:"key_size,omitempty"` // For RSA: 1024, 2048, 4096; defaults to 2048
}

// DKIMGenerateResponse represents a DKIM key generation response
type DKIMGenerateResponse struct {
	Selector   string `json:"selector"`
	PublicKey  string `json:"public_key"`
	DNSRecord  string `json:"dns_record"`
	DNSName    string `json:"dns_name"`
	DomainName string `json:"domain_name"`
}

// GenerateDKIM generates DKIM keys for a domain
func (h *DomainHandler) GenerateDKIM(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	// Parse request body for key type and size options
	var req DKIMGenerateRequest
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	}

	// Default values
	if req.KeyType == "" {
		req.KeyType = "rsa"
	}
	if req.KeySize <= 0 {
		req.KeySize = 2048
	}

	// Validate key type
	if req.KeyType != "rsa" && req.KeyType != "ed25519" {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid key type. Use 'rsa' or 'ed25519'")
		return
	}

	// Validate key size for RSA
	if req.KeyType == "rsa" && req.KeySize < 1024 {
		middleware.RespondError(w, http.StatusBadRequest, "RSA key size must be at least 1024 bits")
		return
	}

	// Get domain to validate it exists and get the name
	dom, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, http.StatusNotFound, "Domain not found")
		return
	}

	// Generate DKIM keys
	result, err := h.service.GenerateDKIMKeys(r.Context(), id, req.KeyType, req.KeySize)
	if err != nil {
		h.logger.Error("Failed to generate DKIM keys",
			zap.Int64("id", id),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to generate DKIM keys")
		return
	}

	h.logger.Info("DKIM keys generated",
		zap.Int64("id", id),
		zap.String("domain", dom.Name),
		zap.String("selector", result.Selector),
		zap.String("key_type", req.KeyType),
	)

	response := DKIMGenerateResponse{
		Selector:   result.Selector,
		PublicKey:  result.PublicKey,
		DNSRecord:  result.DNSRecord,
		DNSName:    result.DNSName,
		DomainName: dom.Name,
	}

	middleware.RespondSuccess(w, response, "DKIM keys generated successfully")
}

// domainToResponse converts a domain model to API response format
func domainToResponse(d *domain.Domain) *DomainResponse {
	// Convert *string to string for CatchallEmail
	catchallEmail := ""
	if d.CatchallEmail != nil {
		catchallEmail = *d.CatchallEmail
	}

	return &DomainResponse{
		ID:                 d.ID,
		Name:               d.Name,
		Status:             d.Status,
		CatchallEmail:      catchallEmail,
		MaxUsers:           d.MaxUsers,
		DefaultQuota:       d.DefaultQuota,
		DKIMSelector:       d.DKIMSelector,
		DKIMPublicKey:      d.DKIMPublicKey, // Don't expose private key
		SPFRecord:          d.SPFRecord,
		DMARCPolicy:        d.DMARCPolicy,
		DMARCReportEmail:   d.DMARCReportEmail,
		DKIMSigningEnabled: d.DKIMSigningEnabled,
		DKIMVerifyEnabled:  d.DKIMVerifyEnabled,
		CreatedAt:          d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
