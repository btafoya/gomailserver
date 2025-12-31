package webdav

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/btafoya/gomailserver/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// BasicAuthMiddleware provides HTTP Basic Authentication for WebDAV endpoints
func BasicAuthMiddleware(userRepo repository.UserRepository, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Debug("WebDAV request without authorization",
					zap.String("path", r.URL.Path),
					zap.String("remote_addr", r.RemoteAddr),
				)
				requestAuth(w)
				return
			}

			// Parse Basic Auth header
			if !strings.HasPrefix(authHeader, "Basic ") {
				logger.Warn("WebDAV request with non-Basic auth",
					zap.String("path", r.URL.Path),
				)
				requestAuth(w)
				return
			}

			// Decode base64 credentials
			encoded := strings.TrimPrefix(authHeader, "Basic ")
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				logger.Warn("failed to decode Basic auth",
					zap.Error(err),
					zap.String("path", r.URL.Path),
				)
				requestAuth(w)
				return
			}

			// Split username:password
			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) != 2 {
				logger.Warn("invalid Basic auth format",
					zap.String("path", r.URL.Path),
				)
				requestAuth(w)
				return
			}

			username := parts[0]
			password := parts[1]

			// Authenticate user
			user, err := userRepo.GetByEmail(username)
			if err != nil {
				logger.Warn("WebDAV authentication failed - user not found",
					zap.String("username", username),
					zap.Error(err),
				)
				requestAuth(w)
				return
			}

			// Verify password
			if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
				logger.Warn("WebDAV authentication failed - invalid password",
					zap.String("username", username),
				)
				requestAuth(w)
				return
			}

			logger.Debug("WebDAV authentication successful",
				zap.Int64("user_id", user.ID),
				zap.String("email", user.Email),
			)

			// Add user ID to request context
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// requestAuth sends a 401 Unauthorized response with WWW-Authenticate header
func requestAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("401 Unauthorized\n"))
}

// GetUserID retrieves user ID from request context
func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int64)
	return userID, ok
}
