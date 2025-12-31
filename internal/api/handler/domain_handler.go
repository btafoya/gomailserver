package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/service"
)

// DomainHandler handles domain management API requests
type DomainHandler struct {
	domainRepo    repository.DomainRepository
	domainService *service.DomainService
	logger        *zap.Logger
}

// NewDomainHandler creates a new domain handler
func NewDomainHandler(domainRepo repository.DomainRepository, logger *zap.Logger) *DomainHandler {
	return &DomainHandler{
		domainRepo:    domainRepo,
		domainService: service.NewDomainService(domainRepo),
		logger:        logger,
	}
}

// ListDomains returns all domains
// GET /api/domains
func (h *DomainHandler) ListDomains(w http.ResponseWriter, r *http.Request) {
	// Get all domains (offset=0, limit=10000)
	domains, err := h.domainRepo.List(0, 10000)
	if err != nil {
		h.logger.Error("failed to list domains", zap.Error(err))
		h.writeError(w, "failed to list domains", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, map[string]interface{}{
		"domains": domains,
		"count":   len(domains),
	})
}

// GetDomain returns a specific domain by name
// GET /api/domains/{name}
func (h *DomainHandler) GetDomain(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		h.writeError(w, "domain name is required", http.StatusBadRequest)
		return
	}

	dom, err := h.domainRepo.GetByName(name)
	if err != nil {
		h.logger.Error("failed to get domain",
			zap.String("domain", name),
			zap.Error(err),
		)
		h.writeError(w, "domain not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, dom)
}

// CreateDomain creates a new domain
// POST /api/domains
func (h *DomainHandler) CreateDomain(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name              string `json:"name"`
		UseDefaultTemplate bool   `json:"use_default_template"`
		Domain            *domain.Domain `json:"domain,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		h.writeError(w, "domain name is required", http.StatusBadRequest)
		return
	}

	// Check if domain already exists
	existing, _ := h.domainRepo.GetByName(req.Name)
	if existing != nil {
		h.writeError(w, "domain already exists", http.StatusConflict)
		return
	}

	var newDomain *domain.Domain
	var err error

	if req.UseDefaultTemplate {
		// Create from default template
		newDomain, err = h.domainService.CreateDomainFromTemplate(req.Name)
		if err != nil {
			h.logger.Error("failed to create domain from template",
				zap.String("domain", req.Name),
				zap.Error(err),
			)
			h.writeError(w, "failed to create domain", http.StatusInternalServerError)
			return
		}
	} else if req.Domain != nil {
		// Create with custom settings
		req.Domain.Name = req.Name
		req.Domain.Status = "active"
		if err := h.domainRepo.Create(req.Domain); err != nil {
			h.logger.Error("failed to create domain",
				zap.String("domain", req.Name),
				zap.Error(err),
			)
			h.writeError(w, "failed to create domain", http.StatusInternalServerError)
			return
		}
		newDomain = req.Domain
	} else {
		h.writeError(w, "either use_default_template or domain configuration is required", http.StatusBadRequest)
		return
	}

	h.logger.Info("domain created",
		zap.String("domain", newDomain.Name),
	)

	w.WriteHeader(http.StatusCreated)
	h.writeJSON(w, newDomain)
}

// UpdateDomain updates an existing domain
// PUT /api/domains/{name}
func (h *DomainHandler) UpdateDomain(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		h.writeError(w, "domain name is required", http.StatusBadRequest)
		return
	}

	// Get existing domain
	existing, err := h.domainRepo.GetByName(name)
	if err != nil {
		h.writeError(w, "domain not found", http.StatusNotFound)
		return
	}

	// Decode update request
	var updates domain.Domain
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Preserve immutable fields
	updates.ID = existing.ID
	updates.Name = existing.Name
	updates.CreatedAt = existing.CreatedAt

	// Update domain
	if err := h.domainRepo.Update(&updates); err != nil {
		h.logger.Error("failed to update domain",
			zap.String("domain", name),
			zap.Error(err),
		)
		h.writeError(w, "failed to update domain", http.StatusInternalServerError)
		return
	}

	h.logger.Info("domain updated",
		zap.String("domain", name),
	)

	h.writeJSON(w, &updates)
}

// DeleteDomain deletes a domain
// DELETE /api/domains/{name}
func (h *DomainHandler) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		h.writeError(w, "domain name is required", http.StatusBadRequest)
		return
	}

	// Prevent deletion of default template
	if name == service.DefaultTemplateDomainName {
		h.writeError(w, "cannot delete default template domain", http.StatusForbidden)
		return
	}

	// Get domain to verify it exists
	dom, err := h.domainRepo.GetByName(name)
	if err != nil {
		h.writeError(w, "domain not found", http.StatusNotFound)
		return
	}

	// Delete domain
	if err := h.domainRepo.Delete(dom.ID); err != nil {
		h.logger.Error("failed to delete domain",
			zap.String("domain", name),
			zap.Error(err),
		)
		h.writeError(w, "failed to delete domain", http.StatusInternalServerError)
		return
	}

	h.logger.Info("domain deleted",
		zap.String("domain", name),
	)

	w.WriteHeader(http.StatusNoContent)
}

// GetDomainSecurity returns security configuration for a domain
// GET /api/domains/{name}/security
func (h *DomainHandler) GetDomainSecurity(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		h.writeError(w, "domain name is required", http.StatusBadRequest)
		return
	}

	dom, err := h.domainRepo.GetByName(name)
	if err != nil {
		h.writeError(w, "domain not found", http.StatusNotFound)
		return
	}

	// Return only security-related fields
	security := map[string]interface{}{
		"domain": name,
		"dkim": map[string]interface{}{
			"signing_enabled":   dom.DKIMSigningEnabled,
			"verify_enabled":    dom.DKIMVerifyEnabled,
			"key_size":          dom.DKIMKeySize,
			"key_type":          dom.DKIMKeyType,
			"headers_to_sign":   dom.DKIMHeadersToSign,
			"selector":          dom.DKIMSelector,
			"public_key":        dom.DKIMPublicKey,
		},
		"spf": map[string]interface{}{
			"enabled":           dom.SPFEnabled,
			"dns_server":        dom.SPFDNSServer,
			"dns_timeout":       dom.SPFDNSTimeout,
			"max_lookups":       dom.SPFMaxLookups,
			"fail_action":       dom.SPFFailAction,
			"softfail_action":   dom.SPFSoftFailAction,
			"record":            dom.SPFRecord,
		},
		"dmarc": map[string]interface{}{
			"enabled":           dom.DMARCEnabled,
			"dns_server":        dom.DMARCDNSServer,
			"dns_timeout":       dom.DMARCDNSTimeout,
			"report_enabled":    dom.DMARCReportEnabled,
			"report_email":      dom.DMARCReportEmail,
			"policy":            dom.DMARCPolicy,
		},
		"clamav": map[string]interface{}{
			"enabled":           dom.ClamAVEnabled,
			"max_scan_size":     dom.ClamAVMaxScanSize,
			"virus_action":      dom.ClamAVVirusAction,
			"fail_action":       dom.ClamAVFailAction,
		},
		"spam": map[string]interface{}{
			"enabled":           dom.SpamEnabled,
			"reject_score":      dom.SpamRejectScore,
			"quarantine_score":  dom.SpamQuarantineScore,
			"learning_enabled":  dom.SpamLearningEnabled,
		},
		"greylist": map[string]interface{}{
			"enabled":           dom.GreylistEnabled,
			"delay_minutes":     dom.GreylistDelayMinutes,
			"expiry_days":       dom.GreylistExpiryDays,
			"cleanup_interval":  dom.GreylistCleanupInterval,
			"whitelist_after":   dom.GreylistWhitelistAfter,
		},
		"rate_limit": map[string]interface{}{
			"enabled":           dom.RateLimitEnabled,
			"smtp_per_ip":       dom.RateLimitSMTPPerIP,
			"smtp_per_user":     dom.RateLimitSMTPPerUser,
			"smtp_per_domain":   dom.RateLimitSMTPPerDomain,
			"auth_per_ip":       dom.RateLimitAuthPerIP,
			"imap_per_user":     dom.RateLimitIMAPPerUser,
			"cleanup_interval":  dom.RateLimitCleanupInterval,
		},
		"auth": map[string]interface{}{
			"totp_enforced":              dom.AuthTOTPEnforced,
			"brute_force_enabled":        dom.AuthBruteForceEnabled,
			"brute_force_threshold":      dom.AuthBruteForceThreshold,
			"brute_force_window_minutes": dom.AuthBruteForceWindowMinutes,
			"brute_force_block_minutes":  dom.AuthBruteForceBlockMinutes,
			"ip_blacklist_enabled":       dom.AuthIPBlacklistEnabled,
			"cleanup_interval":           dom.AuthCleanupInterval,
		},
	}

	h.writeJSON(w, security)
}

// UpdateDomainSecurity updates security configuration for a domain
// PUT /api/domains/{name}/security
func (h *DomainHandler) UpdateDomainSecurity(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		h.writeError(w, "domain name is required", http.StatusBadRequest)
		return
	}

	// Get existing domain
	existing, err := h.domainRepo.GetByName(name)
	if err != nil {
		h.writeError(w, "domain not found", http.StatusNotFound)
		return
	}

	// Decode security updates
	var securityUpdates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&securityUpdates); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Apply updates to domain (preserving non-security fields)
	updated := *existing

	// Update DKIM settings
	if dkim, ok := securityUpdates["dkim"].(map[string]interface{}); ok {
		if v, exists := dkim["signing_enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.DKIMSigningEnabled = b
			}
		}
		if v, exists := dkim["verify_enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.DKIMVerifyEnabled = b
			}
		}
		if v, exists := dkim["key_size"]; exists {
			if f, ok := v.(float64); ok {
				updated.DKIMKeySize = int(f)
			}
		}
		if v, exists := dkim["key_type"]; exists {
			if s, ok := v.(string); ok {
				updated.DKIMKeyType = s
			}
		}
		if v, exists := dkim["headers_to_sign"]; exists {
			if s, ok := v.(string); ok {
				updated.DKIMHeadersToSign = s
			}
		}
	}

	// Update SPF settings
	if spf, ok := securityUpdates["spf"].(map[string]interface{}); ok {
		if v, exists := spf["enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.SPFEnabled = b
			}
		}
		if v, exists := spf["dns_server"]; exists {
			if s, ok := v.(string); ok {
				updated.SPFDNSServer = s
			}
		}
		if v, exists := spf["dns_timeout"]; exists {
			if f, ok := v.(float64); ok {
				updated.SPFDNSTimeout = int(f)
			}
		}
		if v, exists := spf["max_lookups"]; exists {
			if f, ok := v.(float64); ok {
				updated.SPFMaxLookups = int(f)
			}
		}
		if v, exists := spf["fail_action"]; exists {
			if s, ok := v.(string); ok {
				updated.SPFFailAction = s
			}
		}
		if v, exists := spf["softfail_action"]; exists {
			if s, ok := v.(string); ok {
				updated.SPFSoftFailAction = s
			}
		}
	}

	// Update DMARC settings
	if dmarc, ok := securityUpdates["dmarc"].(map[string]interface{}); ok {
		if v, exists := dmarc["enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.DMARCEnabled = b
			}
		}
		if v, exists := dmarc["dns_server"]; exists {
			if s, ok := v.(string); ok {
				updated.DMARCDNSServer = s
			}
		}
		if v, exists := dmarc["dns_timeout"]; exists {
			if f, ok := v.(float64); ok {
				updated.DMARCDNSTimeout = int(f)
			}
		}
		if v, exists := dmarc["report_enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.DMARCReportEnabled = b
			}
		}
		if v, exists := dmarc["report_email"]; exists {
			if s, ok := v.(string); ok {
				updated.DMARCReportEmail = s
			}
		}
	}

	// Update ClamAV settings
	if clamav, ok := securityUpdates["clamav"].(map[string]interface{}); ok {
		if v, exists := clamav["enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.ClamAVEnabled = b
			}
		}
		if v, exists := clamav["max_scan_size"]; exists {
			if f, ok := v.(float64); ok {
				updated.ClamAVMaxScanSize = int64(f)
			}
		}
		if v, exists := clamav["virus_action"]; exists {
			if s, ok := v.(string); ok {
				updated.ClamAVVirusAction = s
			}
		}
		if v, exists := clamav["fail_action"]; exists {
			if s, ok := v.(string); ok {
				updated.ClamAVFailAction = s
			}
		}
	}

	// Update Spam settings
	if spam, ok := securityUpdates["spam"].(map[string]interface{}); ok {
		if v, exists := spam["enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.SpamEnabled = b
			}
		}
		if v, exists := spam["reject_score"]; exists {
			if f, ok := v.(float64); ok {
				updated.SpamRejectScore = f
			}
		}
		if v, exists := spam["quarantine_score"]; exists {
			if f, ok := v.(float64); ok {
				updated.SpamQuarantineScore = f
			}
		}
		if v, exists := spam["learning_enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.SpamLearningEnabled = b
			}
		}
	}

	// Update Greylist settings
	if greylist, ok := securityUpdates["greylist"].(map[string]interface{}); ok {
		if v, exists := greylist["enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.GreylistEnabled = b
			}
		}
		if v, exists := greylist["delay_minutes"]; exists {
			if f, ok := v.(float64); ok {
				updated.GreylistDelayMinutes = int(f)
			}
		}
		if v, exists := greylist["expiry_days"]; exists {
			if f, ok := v.(float64); ok {
				updated.GreylistExpiryDays = int(f)
			}
		}
		if v, exists := greylist["cleanup_interval"]; exists {
			if f, ok := v.(float64); ok {
				updated.GreylistCleanupInterval = int(f)
			}
		}
		if v, exists := greylist["whitelist_after"]; exists {
			if f, ok := v.(float64); ok {
				updated.GreylistWhitelistAfter = int(f)
			}
		}
	}

	// Update Rate Limit settings
	if rateLimit, ok := securityUpdates["rate_limit"].(map[string]interface{}); ok {
		if v, exists := rateLimit["enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.RateLimitEnabled = b
			}
		}
		if v, exists := rateLimit["smtp_per_ip"]; exists {
			if s, ok := v.(string); ok {
				updated.RateLimitSMTPPerIP = s
			}
		}
		if v, exists := rateLimit["smtp_per_user"]; exists {
			if s, ok := v.(string); ok {
				updated.RateLimitSMTPPerUser = s
			}
		}
		if v, exists := rateLimit["smtp_per_domain"]; exists {
			if s, ok := v.(string); ok {
				updated.RateLimitSMTPPerDomain = s
			}
		}
		if v, exists := rateLimit["auth_per_ip"]; exists {
			if s, ok := v.(string); ok {
				updated.RateLimitAuthPerIP = s
			}
		}
		if v, exists := rateLimit["imap_per_user"]; exists {
			if s, ok := v.(string); ok {
				updated.RateLimitIMAPPerUser = s
			}
		}
		if v, exists := rateLimit["cleanup_interval"]; exists {
			if f, ok := v.(float64); ok {
				updated.RateLimitCleanupInterval = int(f)
			}
		}
	}

	// Update Auth settings
	if auth, ok := securityUpdates["auth"].(map[string]interface{}); ok {
		if v, exists := auth["totp_enforced"]; exists {
			if b, ok := v.(bool); ok {
				updated.AuthTOTPEnforced = b
			}
		}
		if v, exists := auth["brute_force_enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.AuthBruteForceEnabled = b
			}
		}
		if v, exists := auth["brute_force_threshold"]; exists {
			if f, ok := v.(float64); ok {
				updated.AuthBruteForceThreshold = int(f)
			}
		}
		if v, exists := auth["brute_force_window_minutes"]; exists {
			if f, ok := v.(float64); ok {
				updated.AuthBruteForceWindowMinutes = int(f)
			}
		}
		if v, exists := auth["brute_force_block_minutes"]; exists {
			if f, ok := v.(float64); ok {
				updated.AuthBruteForceBlockMinutes = int(f)
			}
		}
		if v, exists := auth["ip_blacklist_enabled"]; exists {
			if b, ok := v.(bool); ok {
				updated.AuthIPBlacklistEnabled = b
			}
		}
		if v, exists := auth["cleanup_interval"]; exists {
			if f, ok := v.(float64); ok {
				updated.AuthCleanupInterval = int(f)
			}
		}
	}

	// Update domain in database
	if err := h.domainRepo.Update(&updated); err != nil {
		h.logger.Error("failed to update domain security",
			zap.String("domain", name),
			zap.Error(err),
		)
		h.writeError(w, "failed to update security configuration", http.StatusInternalServerError)
		return
	}

	h.logger.Info("domain security updated",
		zap.String("domain", name),
	)

	h.writeJSON(w, map[string]interface{}{
		"message": "security configuration updated successfully",
		"domain":  name,
	})
}

// GetDefaultTemplate returns the default domain template
// GET /api/domains/_default
func (h *DomainHandler) GetDefaultTemplate(w http.ResponseWriter, r *http.Request) {
	template, err := h.domainService.GetDefaultTemplate()
	if err != nil {
		h.logger.Error("failed to get default template", zap.Error(err))
		h.writeError(w, "default template not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, template)
}

// UpdateDefaultTemplate updates the default domain template
// PUT /api/domains/_default
func (h *DomainHandler) UpdateDefaultTemplate(w http.ResponseWriter, r *http.Request) {
	var updates domain.Domain
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.domainService.UpdateDefaultTemplate(&updates); err != nil {
		h.logger.Error("failed to update default template", zap.Error(err))
		h.writeError(w, "failed to update default template", http.StatusInternalServerError)
		return
	}

	h.logger.Info("default template updated")

	h.writeJSON(w, map[string]interface{}{
		"message": "default template updated successfully",
	})
}

// Helper methods

func (h *DomainHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode JSON response", zap.Error(err))
	}
}

func (h *DomainHandler) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  message,
		"status": status,
	})
}
