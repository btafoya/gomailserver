package middleware

import (
	"fmt"
	"net/http"

	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/security/ratelimit"
	"go.uber.org/zap"
)

// RateLimit returns middleware that enforces API rate limiting
func RateLimit(repo repository.RateLimitRepository, logger *zap.Logger) func(http.Handler) http.Handler {
	limiter := ratelimit.NewLimiter(repo)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP for rate limiting
			ip := getRealIP(r)

			// Check API rate limit per IP
			allowed, err := limiter.Check("api_per_ip", ip)
			if err != nil {
				logger.Error("rate limit check failed",
					zap.Error(err),
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
				)
				// Fail open - allow request on error
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				logger.Warn("rate limit exceeded",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
					zap.String("user_agent", r.UserAgent()),
				)
				RespondError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
				return
			}

			// If authenticated, also check per-user rate limit
			if userID, ok := GetUserID(r); ok {
				key := fmt.Sprintf("user:%d", userID)
				allowed, err = limiter.Check("api_per_user", key)
				if err != nil {
					logger.Error("user rate limit check failed",
						zap.Error(err),
						zap.Int64("user_id", userID),
						zap.String("path", r.URL.Path),
					)
					// Fail open - allow request on error
					next.ServeHTTP(w, r)
					return
				}

				if !allowed {
					logger.Warn("user rate limit exceeded",
						zap.Int64("user_id", userID),
						zap.String("ip", ip),
						zap.String("path", r.URL.Path),
					)
					RespondError(w, http.StatusTooManyRequests, "User rate limit exceeded. Please try again later.")
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getRealIP extracts the real client IP from the request
func getRealIP(r *http.Request) string {
	// Check X-Forwarded-For header first (if behind proxy)
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		return ip
	}

	// Check X-Real-IP header
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
