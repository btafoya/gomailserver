package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const jwtExpiry = 24 * time.Hour

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	EmailKey    contextKey = "email"
	RoleKey     contextKey = "role"
	DomainIDKey contextKey = "domain_id"
)

// Claims represents JWT claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	DomainID *int64 `json:"domain_id,omitempty"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token for a user
func GenerateJWT(userID int64, email, role string, domainID *int64, secret string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		DomainID: domainID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates and parses a JWT token
func ValidateJWT(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Auth returns middleware that validates JWT tokens and API keys
func Auth(jwtSecret string, apiKeyRepo repository.APIKeyRepository, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("missing authorization header",
					zap.String("path", r.URL.Path),
					zap.String("remote_addr", r.RemoteAddr),
				)
				RespondError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 {
				logger.Warn("invalid authorization format",
					zap.String("path", r.URL.Path),
				)
				RespondError(w, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			scheme := parts[0]
			token := parts[1]

			switch scheme {
			case "Bearer":
				// JWT token authentication
				claims, err := ValidateJWT(token, jwtSecret)
				if err != nil {
					logger.Warn("invalid JWT token",
						zap.Error(err),
						zap.String("path", r.URL.Path),
					)
					RespondError(w, http.StatusUnauthorized, "Invalid token")
					return
				}

				// Add claims to request context
				ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
				ctx = context.WithValue(ctx, EmailKey, claims.Email)
				ctx = context.WithValue(ctx, RoleKey, claims.Role)
				if claims.DomainID != nil {
					ctx = context.WithValue(ctx, DomainIDKey, *claims.DomainID)
				}

				next.ServeHTTP(w, r.WithContext(ctx))

			case "ApiKey":
				// API key authentication
				apiKey, err := apiKeyRepo.GetByKeyHash(token)
				if err != nil {
					logger.Warn("invalid API key",
						zap.Error(err),
						zap.String("path", r.URL.Path),
					)
					RespondError(w, http.StatusUnauthorized, "Invalid API key")
					return
				}

				// Check if API key has expired
				if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
					logger.Warn("expired API key",
						zap.Int64("api_key_id", apiKey.ID),
						zap.String("path", r.URL.Path),
					)
					RespondError(w, http.StatusUnauthorized, "API key has expired")
					return
				}

				// Update last used timestamp and IP
				go func() {
					if err := apiKeyRepo.UpdateLastUsed(apiKey.ID, r.RemoteAddr); err != nil {
						logger.Error("failed to update API key last used",
							zap.Error(err),
							zap.Int64("api_key_id", apiKey.ID),
						)
					}
				}()

				// Add API key user info to request context
				ctx := context.WithValue(r.Context(), UserIDKey, apiKey.UserID)
				ctx = context.WithValue(ctx, DomainIDKey, apiKey.DomainID)
				ctx = context.WithValue(ctx, RoleKey, "api")

				next.ServeHTTP(w, r.WithContext(ctx))

			default:
				logger.Warn("unknown authorization scheme",
					zap.String("scheme", scheme),
					zap.String("path", r.URL.Path),
				)
				RespondError(w, http.StatusUnauthorized, "Unknown authorization scheme")
			}
		})
	}
}

// GetUserID retrieves user ID from request context
func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int64)
	return userID, ok
}

// GetEmail retrieves email from request context
func GetEmail(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(EmailKey).(string)
	return email, ok
}

// GetRole retrieves role from request context
func GetRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(RoleKey).(string)
	return role, ok
}

// GetDomainID retrieves domain ID from request context
func GetDomainID(r *http.Request) (int64, bool) {
	domainID, ok := r.Context().Value(DomainIDKey).(int64)
	return domainID, ok
}
