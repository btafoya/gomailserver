package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/btafoya/gomailserver/internal/postmark/models"
	"github.com/btafoya/gomailserver/internal/postmark/repository"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// ServerIDKey is the context key for the server ID
	ServerIDKey ContextKey = "server_id"
	// ServerKey is the context key for the full server object
	ServerKey ContextKey = "server"
)

// AuthMiddleware validates the PostmarkApp API token
func AuthMiddleware(repo repository.PostmarkRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from headers (case-insensitive)
			token := r.Header.Get("X-Postmark-Server-Token")
			if token == "" {
				token = r.Header.Get("X-Postmark-Account-Token")
			}

			if token == "" {
				models.WriteError(w, models.ErrorCodeUnauthorized, models.MsgUnauthorized)
				return
			}

			// Check for test mode
			if token == "POSTMARK_API_TEST" {
				// Test mode: allow but don't send actual emails
				ctx := context.WithValue(r.Context(), ServerIDKey, 0)
				ctx = context.WithValue(ctx, ServerKey, &models.Server{
					ID:   0,
					Name: "Test Server",
				})
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Validate token and get server
			server, err := repo.GetServerByToken(r.Context(), token)
			if err != nil {
				models.WriteError(w, models.ErrorCodeUnauthorized, models.MsgUnauthorized)
				return
			}

			// Check if server is active
			if server.ID > 0 {
				// Get full server details to check active status
				fullServer, err := repo.GetServer(r.Context(), server.ID)
				if err != nil || !fullServer.SmtpApiActivated {
					models.WriteError(w, models.ErrorCodeInactiveServer, models.MsgInactiveServer)
					return
				}
			}

			// Add server info to context
			ctx := context.WithValue(r.Context(), ServerIDKey, server.ID)
			ctx = context.WithValue(ctx, ServerKey, server)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetServerID retrieves the server ID from the request context
func GetServerID(r *http.Request) int {
	if id, ok := r.Context().Value(ServerIDKey).(int); ok {
		return id
	}
	return 0
}

// GetServer retrieves the server from the request context
func GetServer(r *http.Request) *models.Server {
	if server, ok := r.Context().Value(ServerKey).(*models.Server); ok {
		return server
	}
	return nil
}

// RequireJSONMiddleware ensures the Content-Type is application/json
func RequireJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				models.WriteError(w, models.ErrorCodeJSONRequired, models.MsgJSONRequired)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
