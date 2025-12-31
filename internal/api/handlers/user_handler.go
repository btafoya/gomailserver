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

// UserHandler handles user management endpoints
type UserHandler struct {
	service *service.UserService
	logger  *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// UserRequest represents a user creation/update request
type UserRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password,omitempty"`
	FullName        string `json:"full_name"`
	DisplayName     string `json:"display_name,omitempty"`
	DomainID        int64  `json:"domain_id"`
	Quota           int64  `json:"quota,omitempty"`
	Status          string `json:"status,omitempty"`
	ForwardingRules string `json:"forwarding_rules,omitempty"`
	AutoReply       string `json:"auto_reply,omitempty"`
	SpamThreshold   int    `json:"spam_threshold,omitempty"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID              int64  `json:"id"`
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	DisplayName     string `json:"display_name,omitempty"`
	DomainID        int64  `json:"domain_id"`
	DomainName      string `json:"domain_name,omitempty"`
	Quota           int64  `json:"quota"`
	CurrentUsage    int64  `json:"current_usage"`
	Status          string `json:"status"`
	ForwardingRules string `json:"forwarding_rules,omitempty"`
	AutoReply       string `json:"auto_reply,omitempty"`
	SpamThreshold   int    `json:"spam_threshold"`
	TOTPEnabled     bool   `json:"totp_enabled"`
	CreatedAt       string `json:"created_at"`
	LastLogin       string `json:"last_login,omitempty"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	NewPassword string `json:"new_password"`
}

// List retrieves all users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Add pagination support with query parameters
	// For now, retrieve all users
	users, err := h.service.ListAll(r.Context())
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	// Convert to response format
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = h.userToResponse(r.Context(), u)
	}

	middleware.RespondSuccess(w, responses, "Users retrieved successfully")
}

// Create creates a new user
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if req.Password == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Password is required")
		return
	}
	if req.DomainID == 0 {
		middleware.RespondError(w, http.StatusBadRequest, "Domain ID is required")
		return
	}

	// Convert request to user model
	newUser := &domain.User{
		Email:           req.Email,
		FullName:        req.FullName,
		DisplayName:     req.DisplayName,
		DomainID:        req.DomainID,
		Quota:           req.Quota,
		Status:          req.Status,
		ForwardingRules: req.ForwardingRules,
		AutoReply:       req.AutoReply,
		SpamThreshold:   req.SpamThreshold,
	}

	// Set defaults
	if newUser.Status == "" {
		newUser.Status = "active"
	}
	if newUser.SpamThreshold == 0 {
		newUser.SpamThreshold = 5 // Default spam threshold
	}

	// Create user (password will be hashed by service)
	err := h.service.CreateWithPassword(r.Context(), newUser, req.Password)
	if err != nil {
		h.logger.Error("Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	h.logger.Info("User created",
		zap.String("email", newUser.Email),
		zap.Int64("id", newUser.ID),
	)

	middleware.RespondCreated(w, h.userToResponse(r.Context(), newUser), "User created successfully")
}

// Get retrieves a specific user
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	middleware.RespondSuccess(w, h.userToResponse(r.Context(), user), "User retrieved successfully")
}

// Update updates a user
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing user
	existingUser, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Update fields (don't allow email change)
	if req.FullName != "" {
		existingUser.FullName = req.FullName
	}
	if req.DisplayName != "" {
		existingUser.DisplayName = req.DisplayName
	}
	if req.Quota > 0 {
		existingUser.Quota = req.Quota
	}
	if req.Status != "" {
		existingUser.Status = req.Status
	}
	if req.ForwardingRules != "" {
		existingUser.ForwardingRules = req.ForwardingRules
	}
	if req.AutoReply != "" {
		existingUser.AutoReply = req.AutoReply
	}
	if req.SpamThreshold > 0 {
		existingUser.SpamThreshold = req.SpamThreshold
	}

	// Update user
	err = h.service.Update(r.Context(), existingUser)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	h.logger.Info("User updated",
		zap.Int64("id", id),
		zap.String("email", existingUser.Email),
	)

	middleware.RespondSuccess(w, h.userToResponse(r.Context(), existingUser), "User updated successfully")
}

// Delete deletes a user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete user", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	h.logger.Info("User deleted", zap.Int64("id", id))

	middleware.RespondNoContent(w)
}

// ResetPassword resets a user's password
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.NewPassword == "" {
		middleware.RespondError(w, http.StatusBadRequest, "New password is required")
		return
	}

	// Update password
	err = h.service.UpdatePassword(r.Context(), id, req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to reset password", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to reset password")
		return
	}

	h.logger.Info("Password reset", zap.Int64("id", id))

	middleware.RespondSuccess(w, nil, "Password reset successfully")
}

// userToResponse converts a user model to API response format
func (h *UserHandler) userToResponse(ctx any, u *domain.User) *UserResponse {
	response := &UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FullName:        u.FullName,
		DisplayName:     u.DisplayName,
		DomainID:        u.DomainID,
		Quota:           u.Quota,
		CurrentUsage:    u.CurrentUsage,
		Status:          u.Status,
		ForwardingRules: u.ForwardingRules,
		AutoReply:       u.AutoReply,
		SpamThreshold:   u.SpamThreshold,
		TOTPEnabled:     u.TOTPSecret != "",
		CreatedAt:       u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if u.LastLogin != nil {
		response.LastLogin = u.LastLogin.Format("2006-01-02T15:04:05Z07:00")
	}

	// Get domain name if possible
	// Note: This would require access to domain service, skipping for now
	// TODO: Add domain name lookup

	return response
}
