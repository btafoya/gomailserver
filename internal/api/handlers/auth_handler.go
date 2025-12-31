package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/btafoya/gomailserver/internal/service"
	"go.uber.org/zap"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userService *service.UserService
	jwtSecret   string
	logger      *zap.Logger
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(userService *service.UserService, jwtSecret string, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
		logger:      logger,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TOTPCode string `json:"totp_code,omitempty"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserInfo represents user information in responses
type UserInfo struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Role        string `json:"role"`
	DomainID    int64  `json:"domain_id"`
	DomainName  string `json:"domain_name"`
	TOTPEnabled bool   `json:"totp_enabled"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Authenticate user
	user, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Authentication failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		middleware.RespondError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Check if TOTP is enabled
	if user.TOTPSecret != "" {
		if req.TOTPCode == "" {
			middleware.RespondError(w, http.StatusUnauthorized, "TOTP code required")
			return
		}
		// TODO: Validate TOTP code using security/totp package
		// For now, we'll skip TOTP validation
	}

	// Determine role (admin vs user)
	role := "user"
	// TODO: Check if user is admin (could be based on domain ownership or specific admin flag)
	// For now, all users are regular users

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Email, role, &user.DomainID, h.jwtSecret)
	if err != nil {
		h.logger.Error("Failed to generate JWT", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	// Generate refresh token (valid for 7 days)
	refreshToken, err := middleware.GenerateJWT(user.ID, user.Email, "refresh", &user.DomainID, h.jwtSecret)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Get domain information
	domain, err := h.userService.GetDomainByID(r.Context(), user.DomainID)
	domainName := ""
	if err == nil && domain != nil {
		domainName = domain.Name
	}

	// Prepare response
	response := LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: &UserInfo{
			ID:          user.ID,
			Email:       user.Email,
			FullName:    user.FullName,
			Role:        role,
			DomainID:    user.DomainID,
			DomainName:  domainName,
			TOTPEnabled: user.TOTPSecret != "",
		},
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	h.logger.Info("User logged in successfully",
		zap.String("email", user.Email),
		zap.Int64("user_id", user.ID),
	)

	middleware.RespondSuccess(w, response, "Login successful")
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate refresh token
	claims, err := middleware.ValidateJWT(req.RefreshToken, h.jwtSecret)
	if err != nil {
		middleware.RespondError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Verify it's a refresh token
	if claims.Role != "refresh" {
		middleware.RespondError(w, http.StatusUnauthorized, "Invalid token type")
		return
	}

	// Get user to verify they still exist and are active
	user, err := h.userService.GetByID(claims.UserID)
	if err != nil {
		middleware.RespondError(w, http.StatusUnauthorized, "User not found")
		return
	}

	if user.Status != "active" {
		middleware.RespondError(w, http.StatusUnauthorized, "User account is not active")
		return
	}

	// Generate new access token
	role := "user"
	// TODO: Check admin status
	newToken, err := middleware.GenerateJWT(user.ID, user.Email, role, &user.DomainID, h.jwtSecret)
	if err != nil {
		h.logger.Error("Failed to generate new JWT", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to generate new token")
		return
	}

	// Generate new refresh token
	newRefreshToken, err := middleware.GenerateJWT(user.ID, user.Email, "refresh", &user.DomainID, h.jwtSecret)
	if err != nil {
		h.logger.Error("Failed to generate new refresh token", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Get domain information
	domain, err := h.userService.GetDomainByID(r.Context(), user.DomainID)
	domainName := ""
	if err == nil && domain != nil {
		domainName = domain.Name
	}

	// Prepare response
	response := LoginResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		User: &UserInfo{
			ID:          user.ID,
			Email:       user.Email,
			FullName:    user.FullName,
			Role:        role,
			DomainID:    user.DomainID,
			DomainName:  domainName,
			TOTPEnabled: user.TOTPSecret != "",
		},
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	middleware.RespondSuccess(w, response, "Token refreshed successfully")
}
