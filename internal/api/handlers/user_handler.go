package handlers

import (
	"context"
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
	Email            string  `json:"email"`
	Password         string  `json:"password,omitempty"`
	FullName         string  `json:"full_name"`
	DisplayName      string  `json:"display_name,omitempty"`
	DomainID         int64   `json:"domain_id"`
	Quota            int64   `json:"quota,omitempty"`
	Status           string  `json:"status,omitempty"`
	ForwardTo        string  `json:"forward_to,omitempty"`
	AutoReplyEnabled bool    `json:"auto_reply_enabled"`
	AutoReplySubject string  `json:"auto_reply_subject,omitempty"`
	AutoReplyBody    string  `json:"auto_reply_body,omitempty"`
	SpamThreshold    float64 `json:"spam_threshold,omitempty"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID               int64   `json:"id"`
	Email            string  `json:"email"`
	FullName         string  `json:"full_name"`
	DisplayName      string  `json:"display_name,omitempty"`
	DomainID         int64   `json:"domain_id"`
	DomainName       string  `json:"domain_name,omitempty"`
	Quota            int64   `json:"quota"`
	UsedQuota        int64   `json:"used_quota"`
	Status           string  `json:"status"`
	ForwardTo        string  `json:"forward_to,omitempty"`
	AutoReplyEnabled bool    `json:"auto_reply_enabled"`
	AutoReplySubject string  `json:"auto_reply_subject,omitempty"`
	AutoReplyBody    string  `json:"auto_reply_body,omitempty"`
	SpamThreshold    float64 `json:"spam_threshold"`
	TOTPEnabled      bool    `json:"totp_enabled"`
	CreatedAt        string  `json:"created_at"`
	LastLogin        string  `json:"last_login,omitempty"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	NewPassword string `json:"new_password"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
	TotalCount int64           `json:"total_count"`
}

// List retrieves users with pagination and optional domain filtering
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	var users []*domain.User
	var totalCount int64
	var err error

	// Filter by domain if specified
	if domainIDStr != "" {
		domainID, parseErr := strconv.ParseInt(domainIDStr, 10, 64)
		if parseErr != nil {
			middleware.RespondError(w, http.StatusBadRequest, "Invalid domain_id")
			return
		}
		users, err = h.service.ListPaginated(ctx, offset, pageSize)
		if err != nil {
			h.logger.Error("Failed to list users", zap.Error(err))
			middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve users")
			return
		}
		// Filter in memory for domain (or use dedicated method if we add it later)
		filteredUsers := make([]*domain.User, 0)
		for _, u := range users {
			if u.DomainID == domainID {
				filteredUsers = append(filteredUsers, u)
			}
		}
		users = filteredUsers
		totalCount, err = h.service.CountByDomain(ctx, domainID)
	} else {
		users, err = h.service.ListPaginated(ctx, offset, pageSize)
		if err != nil {
			h.logger.Error("Failed to list users", zap.Error(err))
			middleware.RespondError(w, http.StatusInternalServerError, "Failed to retrieve users")
			return
		}
		totalCount, err = h.service.Count(ctx)
	}

	if err != nil {
		h.logger.Error("Failed to count users", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to count users")
		return
	}

	// Convert to response format with domain name lookup
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = h.userToResponseWithDomain(ctx, u)
	}

	// Calculate total pages
	totalPages := int(totalCount / int64(pageSize))
	if totalCount%int64(pageSize) > 0 {
		totalPages++
	}

	response := UserListResponse{
		Users:      responses,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}

	middleware.RespondSuccess(w, response, "Users retrieved successfully")
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
		Email:            req.Email,
		FullName:         req.FullName,
		DisplayName:      req.DisplayName,
		DomainID:         req.DomainID,
		Quota:            req.Quota,
		Status:           req.Status,
		ForwardTo:        req.ForwardTo,
		AutoReplyEnabled: req.AutoReplyEnabled,
		AutoReplySubject: req.AutoReplySubject,
		AutoReplyBody:    req.AutoReplyBody,
		SpamThreshold:    req.SpamThreshold,
	}

	// Set defaults
	if newUser.Status == "" {
		newUser.Status = "active"
	}
	if newUser.SpamThreshold == 0 {
		newUser.SpamThreshold = 5.0 // Default spam threshold
	}

	// Create user (password will be hashed by service)
	err := h.service.CreateWithPassword(newUser, req.Password)
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

	middleware.RespondCreated(w, h.userToResponse(newUser), "User created successfully")
}

// Get retrieves a specific user
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	middleware.RespondSuccess(w, h.userToResponse(user), "User retrieved successfully")
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
	existingUser, err := h.service.GetByID(id)
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
	if req.ForwardTo != "" {
		existingUser.ForwardTo = req.ForwardTo
	}
	existingUser.AutoReplyEnabled = req.AutoReplyEnabled
	if req.AutoReplySubject != "" {
		existingUser.AutoReplySubject = req.AutoReplySubject
	}
	if req.AutoReplyBody != "" {
		existingUser.AutoReplyBody = req.AutoReplyBody
	}
	if req.SpamThreshold > 0 {
		existingUser.SpamThreshold = req.SpamThreshold
	}

	// Update user
	err = h.service.Update(existingUser)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	h.logger.Info("User updated",
		zap.Int64("id", id),
		zap.String("email", existingUser.Email),
	)

	middleware.RespondSuccess(w, h.userToResponse(existingUser), "User updated successfully")
}

// Delete deletes a user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.service.Delete(id)
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
	err = h.service.UpdatePassword(id, req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to reset password", zap.Int64("id", id), zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to reset password")
		return
	}

	h.logger.Info("Password reset", zap.Int64("id", id))

	middleware.RespondSuccess(w, nil, "Password reset successfully")
}

// userToResponse converts a user model to API response format
func (h *UserHandler) userToResponse(u *domain.User) *UserResponse {
	response := &UserResponse{
		ID:               u.ID,
		Email:            u.Email,
		FullName:         u.FullName,
		DisplayName:      u.DisplayName,
		DomainID:         u.DomainID,
		Quota:            u.Quota,
		UsedQuota:        u.UsedQuota,
		Status:           u.Status,
		ForwardTo:        u.ForwardTo,
		AutoReplyEnabled: u.AutoReplyEnabled,
		AutoReplySubject: u.AutoReplySubject,
		AutoReplyBody:    u.AutoReplyBody,
		SpamThreshold:    u.SpamThreshold,
		TOTPEnabled:      u.TOTPSecret != "",
		CreatedAt:        u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if u.LastLogin != nil {
		response.LastLogin = u.LastLogin.Format("2006-01-02T15:04:05Z07:00")
	}

	return response
}

// userToResponseWithDomain converts a user model to API response format with domain name lookup
func (h *UserHandler) userToResponseWithDomain(ctx context.Context, u *domain.User) *UserResponse {
	response := h.userToResponse(u)

	// Lookup domain name
	dom, err := h.service.GetDomainByID(ctx, u.DomainID)
	if err == nil && dom != nil {
		response.DomainName = dom.Name
	}

	return response
}
