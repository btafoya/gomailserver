package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/service"
	"go.uber.org/zap"
)

// SetupHandler handles setup wizard endpoints
type SetupHandler struct {
	setupService *service.SetupService
	logger       *zap.Logger
}

// NewSetupHandler creates a new setup handler
func NewSetupHandler(setupService *service.SetupService, logger *zap.Logger) *SetupHandler {
	return &SetupHandler{
		setupService: setupService,
		logger:       logger,
	}
}

// StatusResponse represents the setup status response
type StatusResponse struct {
	SetupComplete bool   `json:"setup_complete"`
	CurrentStep   string `json:"current_step,omitempty"`
}

// GetStatus returns the current setup status
func (h *SetupHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	isComplete, err := h.setupService.IsSetupComplete(r.Context())
	if err != nil {
		h.logger.Error("Failed to check setup status", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to check setup status")
		return
	}

	var currentStep string
	if !isComplete {
		state, err := h.setupService.GetSetupState(r.Context())
		if err != nil {
			h.logger.Error("Failed to get setup state", zap.Error(err))
			// Don't fail, just don't include the current step
		} else {
			currentStep = state.CurrentStep
		}
	}

	response := StatusResponse{
		SetupComplete: isComplete,
		CurrentStep:   currentStep,
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// GetState returns the full setup wizard state
func (h *SetupHandler) GetState(w http.ResponseWriter, r *http.Request) {
	// Check if setup is already complete
	isComplete, err := h.setupService.IsSetupComplete(r.Context())
	if err != nil {
		h.logger.Error("Failed to check setup status", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to check setup status")
		return
	}

	if isComplete {
		middleware.RespondError(w, http.StatusForbidden, "Setup is already complete")
		return
	}

	state, err := h.setupService.GetSetupState(r.Context())
	if err != nil {
		h.logger.Error("Failed to get setup state", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get setup state")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, state)
}

// CreateAdmin creates the first admin user
func (h *SetupHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	// Check if setup is already complete
	isComplete, err := h.setupService.IsSetupComplete(r.Context())
	if err != nil {
		h.logger.Error("Failed to check setup status", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to check setup status")
		return
	}

	if isComplete {
		middleware.RespondError(w, http.StatusForbidden, "Setup is already complete")
		return
	}

	var req service.AdminUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if len(req.Password) < 8 {
		middleware.RespondError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	if err := h.setupService.CreateAdminUser(r.Context(), &req); err != nil {
		h.logger.Error("Failed to create admin user", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	middleware.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "Admin user created successfully",
	})
}

// CompleteSetup marks the setup wizard as complete
func (h *SetupHandler) CompleteSetup(w http.ResponseWriter, r *http.Request) {
	// Check if setup is already complete
	isComplete, err := h.setupService.IsSetupComplete(r.Context())
	if err != nil {
		h.logger.Error("Failed to check setup status", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to check setup status")
		return
	}

	if isComplete {
		middleware.RespondError(w, http.StatusForbidden, "Setup is already complete")
		return
	}

	if err := h.setupService.CompleteSetup(r.Context()); err != nil {
		h.logger.Error("Failed to complete setup", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to complete setup")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Setup completed successfully",
	})
}
