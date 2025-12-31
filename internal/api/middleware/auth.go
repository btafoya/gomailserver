package middleware

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// AuthMiddleware handles bearer token authentication for admin endpoints
type AuthMiddleware struct {
	adminToken string
	logger     *zap.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(adminToken string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		adminToken: adminToken,
		logger:     logger,
	}
}

// Wrap wraps an HTTP handler with authentication
func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If no admin token is configured, allow all requests (dev mode)
		if m.adminToken == "" {
			m.logger.Warn("admin API running without authentication - NOT RECOMMENDED FOR PRODUCTION")
			next.ServeHTTP(w, r)
			return
		}

		// Extract authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn("unauthorized API request - missing authorization header",
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("path", r.URL.Path),
			)
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Validate bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("unauthorized API request - invalid authorization format",
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("path", r.URL.Path),
			)
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		// Validate token
		token := parts[1]
		if token != m.adminToken {
			m.logger.Warn("unauthorized API request - invalid token",
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("path", r.URL.Path),
			)
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Token valid, proceed
		next.ServeHTTP(w, r)
	})
}
