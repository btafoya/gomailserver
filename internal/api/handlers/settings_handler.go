package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/service"
	"go.uber.org/zap"
)

// SettingsHandler handles settings-related API endpoints
type SettingsHandler struct {
	service *service.SettingsService
	logger  *zap.Logger
}

// NewSettingsHandler creates a new settings handler instance
func NewSettingsHandler(service *service.SettingsService, logger *zap.Logger) *SettingsHandler {
	return &SettingsHandler{
		service: service,
		logger:  logger,
	}
}

// GetServer retrieves current server configuration
func (h *SettingsHandler) GetServer(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.GetServerSettings(r.Context())
	if err != nil {
		h.logger.Error("Failed to get server settings", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve server settings")
		return
	}

	middleware.RespondSuccess(w, settings, "Server settings retrieved successfully")
}

// UpdateServer updates server configuration
func (h *SettingsHandler) UpdateServer(w http.ResponseWriter, r *http.Request) {
	var settings service.ServerSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if settings.Hostname == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Hostname is required")
		return
	}
	if settings.Domain == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Domain is required")
		return
	}

	// Validate port numbers
	if settings.SMTPSubmissionPort < 1 || settings.SMTPSubmissionPort > 65535 {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid SMTP submission port")
		return
	}
	if settings.SMTPRelayPort < 1 || settings.SMTPRelayPort > 65535 {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid SMTP relay port")
		return
	}
	if settings.SMTPSPort < 1 || settings.SMTPSPort > 65535 {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid SMTPS port")
		return
	}
	if settings.IMAPPort < 1 || settings.IMAPPort > 65535 {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid IMAP port")
		return
	}
	if settings.IMAPSPort < 1 || settings.IMAPSPort > 65535 {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid IMAPS port")
		return
	}
	if settings.APIPort < 1 || settings.APIPort > 65535 {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid API port")
		return
	}

	if err := h.service.UpdateServerSettings(r.Context(), &settings); err != nil {
		h.logger.Error("Failed to update server settings", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update server settings")
		return
	}

	h.logger.Info("Server settings updated",
		zap.String("hostname", settings.Hostname),
		zap.String("domain", settings.Domain),
	)

	// Return success with restart notice
	response := map[string]interface{}{
		"settings": settings,
		"message":  "Server settings updated successfully. Server restart required for changes to take effect.",
		"requires_restart": true,
	}

	middleware.RespondSuccess(w, response, "Server settings updated")
}

// GetSecurity retrieves current security configuration
func (h *SettingsHandler) GetSecurity(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.GetSecuritySettings(r.Context())
	if err != nil {
		h.logger.Error("Failed to get security settings", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve security settings")
		return
	}

	middleware.RespondSuccess(w, settings, "Security settings retrieved successfully")
}

// UpdateSecurity updates security configuration
func (h *SettingsHandler) UpdateSecurity(w http.ResponseWriter, r *http.Request) {
	var settings service.SecuritySettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate ClamAV settings if enabled
	if settings.ClamAVEnabled && settings.ClamAVSocketPath == "" {
		middleware.RespondError(w, http.StatusBadRequest, "ClamAV socket path is required when ClamAV is enabled")
		return
	}

	// Validate SpamAssassin settings if enabled
	if settings.SpamAssassinEnabled {
		if settings.SpamAssassinHost == "" {
			middleware.RespondError(w, http.StatusBadRequest, "SpamAssassin host is required when SpamAssassin is enabled")
			return
		}
		if settings.SpamAssassinPort < 1 || settings.SpamAssassinPort > 65535 {
			middleware.RespondError(w, http.StatusBadRequest, "Invalid SpamAssassin port")
			return
		}
	}

	if err := h.service.UpdateSecuritySettings(r.Context(), &settings); err != nil {
		h.logger.Error("Failed to update security settings", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update security settings")
		return
	}

	h.logger.Info("Security settings updated",
		zap.Bool("clamav_enabled", settings.ClamAVEnabled),
		zap.Bool("spamassassin_enabled", settings.SpamAssassinEnabled),
	)

	// Return success with restart notice
	response := map[string]interface{}{
		"settings": settings,
		"message":  "Security settings updated successfully. Server restart may be required for some changes to take effect.",
		"requires_restart": true,
	}

	middleware.RespondSuccess(w, response, "Security settings updated")
}

// GetTLS retrieves current TLS/certificate configuration
func (h *SettingsHandler) GetTLS(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.GetTLSSettings(r.Context())
	if err != nil {
		h.logger.Error("Failed to get TLS settings", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve TLS settings")
		return
	}

	middleware.RespondSuccess(w, settings, "TLS settings retrieved successfully")
}

// UpdateTLS updates TLS/certificate configuration
func (h *SettingsHandler) UpdateTLS(w http.ResponseWriter, r *http.Request) {
	var settings service.TLSSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate ACME settings if enabled
	if settings.ACMEEnabled {
		if settings.ACMEEmail == "" {
			middleware.RespondError(w, http.StatusBadRequest, "ACME email is required when ACME is enabled")
			return
		}
		if settings.ACMEProvider == "" {
			middleware.RespondError(w, http.StatusBadRequest, "ACME provider is required when ACME is enabled")
			return
		}
	}

	// Validate manual certificate paths if ACME is disabled
	if !settings.ACMEEnabled {
		if settings.CertFile == "" && settings.KeyFile == "" {
			middleware.RespondError(w, http.StatusBadRequest, "Certificate and key file paths are required when ACME is disabled")
			return
		}
	}

	if err := h.service.UpdateTLSSettings(r.Context(), &settings); err != nil {
		h.logger.Error("Failed to update TLS settings", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update TLS settings")
		return
	}

	h.logger.Info("TLS settings updated",
		zap.Bool("acme_enabled", settings.ACMEEnabled),
		zap.String("acme_provider", settings.ACMEProvider),
	)

	// Return success with restart notice
	response := map[string]interface{}{
		"settings": settings,
		"message":  "TLS settings updated successfully. Server restart required for changes to take effect.",
		"requires_restart": true,
	}

	middleware.RespondSuccess(w, response, "TLS settings updated")
}
