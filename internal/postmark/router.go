package postmark

import (
	"database/sql"
	"net/http"

	"github.com/btafoya/gomailserver/internal/postmark/handlers"
	"github.com/btafoya/gomailserver/internal/postmark/middleware"
	"github.com/btafoya/gomailserver/internal/postmark/repository/sqlite"
	"github.com/btafoya/gomailserver/internal/postmark/service"
	queueService "github.com/btafoya/gomailserver/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// NewRouter creates a new PostmarkApp API router
func NewRouter(db *sql.DB, queueSvc *queueService.QueueService, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	// Create repository
	repo := sqlite.New(db)

	// Create services
	emailSvc := service.NewEmailService(repo, queueSvc, logger)

	// Create handlers
	emailHandler := handlers.NewEmailHandler(emailSvc, logger)

	// Apply middleware
	r.Use(middleware.RequireJSONMiddleware)
	r.Use(middleware.AuthMiddleware(repo))

	// Email endpoints
	r.Post("/email", emailHandler.Send)
	r.Post("/email/batch", emailHandler.SendBatch)

	// Template endpoints (placeholder for future implementation)
	r.Get("/templates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"TotalCount": 0, "Templates": []}`))
	})

	// Webhook endpoints (placeholder for future implementation)
	r.Get("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Webhooks": []}`))
	})

	// Server endpoint (placeholder for future implementation)
	r.Get("/server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Name": "gomailserver", "ApiTokens": []}`))
	})

	return r
}
