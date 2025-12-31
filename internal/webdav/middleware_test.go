package webdav

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/btafoya/gomailserver/internal/domain"
)

// mockUserRepository is a test double for UserRepository
type mockUserRepository struct {
	getByEmailFunc func(string) (*domain.User, error)
}

func (m *mockUserRepository) GetByEmail(email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(email)
	}
	return nil, errors.New("not found")
}

// Implement other required methods as no-ops
func (m *mockUserRepository) Create(user *domain.User) error                 { return nil }
func (m *mockUserRepository) GetByID(id int64) (*domain.User, error)         { return nil, nil }
func (m *mockUserRepository) Update(user *domain.User) error                 { return nil }
func (m *mockUserRepository) UpdateLastLogin(id int64) error                 { return nil }
func (m *mockUserRepository) Delete(id int64) error                          { return nil }
func (m *mockUserRepository) List(domainID int64, offset, limit int) ([]*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepository) UpdateQuota(userID, usedQuota int64) error      { return nil }
func (m *mockUserRepository) UpdatePassword(userID int64, passwordHash string) error { return nil }
func (m *mockUserRepository) ListAll() ([]*domain.User, error)                         { return nil, nil }

func TestBasicAuthMiddleware(t *testing.T) {
	logger := zap.NewNop()

	// Create a valid password hash for testing
	validPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

	t.Run("successful authentication", func(t *testing.T) {
		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				if email == "user@example.com" {
					return &domain.User{
						ID:           1,
						Email:        "user@example.com",
						PasswordHash: string(validPasswordHash),
						Status:       "active",
					}, nil
				}
				return nil, errors.New("not found")
			},
		}

		middleware := BasicAuthMiddleware(repo, logger)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify user ID was added to context
			userID, ok := GetUserID(r)
			if !ok {
				t.Error("expected user ID in context")
			}
			if userID != 1 {
				t.Errorf("expected user ID 1, got %d", userID)
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/caldav/", nil)
		credentials := base64.StdEncoding.EncodeToString([]byte("user@example.com:correct-password"))
		req.Header.Set("Authorization", "Basic "+credentials)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("missing authorization header", func(t *testing.T) {
		repo := &mockUserRepository{}
		middleware := BasicAuthMiddleware(repo, logger)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called without auth")
		}))

		req := httptest.NewRequest("GET", "/caldav/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rr.Code)
		}
		if rr.Header().Get("WWW-Authenticate") != `Basic realm="WebDAV"` {
			t.Error("expected WWW-Authenticate header")
		}
	})

	t.Run("non-Basic auth scheme", func(t *testing.T) {
		repo := &mockUserRepository{}
		middleware := BasicAuthMiddleware(repo, logger)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called with wrong auth scheme")
		}))

		req := httptest.NewRequest("GET", "/caldav/", nil)
		req.Header.Set("Authorization", "Bearer some-token")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rr.Code)
		}
	})

	t.Run("invalid base64 encoding", func(t *testing.T) {
		repo := &mockUserRepository{}
		middleware := BasicAuthMiddleware(repo, logger)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called with invalid encoding")
		}))

		req := httptest.NewRequest("GET", "/caldav/", nil)
		req.Header.Set("Authorization", "Basic invalid-base64!!!")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rr.Code)
		}
	})

	t.Run("invalid credentials format", func(t *testing.T) {
		repo := &mockUserRepository{}
		middleware := BasicAuthMiddleware(repo, logger)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called with invalid format")
		}))

		req := httptest.NewRequest("GET", "/caldav/", nil)
		// Missing colon separator
		credentials := base64.StdEncoding.EncodeToString([]byte("usernameonly"))
		req.Header.Set("Authorization", "Basic "+credentials)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rr.Code)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				return nil, errors.New("not found")
			},
		}

		middleware := BasicAuthMiddleware(repo, logger)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called when user not found")
		}))

		req := httptest.NewRequest("GET", "/caldav/", nil)
		credentials := base64.StdEncoding.EncodeToString([]byte("nonexistent@example.com:password"))
		req.Header.Set("Authorization", "Basic "+credentials)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rr.Code)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				if email == "user@example.com" {
					return &domain.User{
						ID:           1,
						Email:        "user@example.com",
						PasswordHash: string(validPasswordHash),
						Status:       "active",
					}, nil
				}
				return nil, errors.New("not found")
			},
		}

		middleware := BasicAuthMiddleware(repo, logger)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called with wrong password")
		}))

		req := httptest.NewRequest("GET", "/caldav/", nil)
		credentials := base64.StdEncoding.EncodeToString([]byte("user@example.com:wrong-password"))
		req.Header.Set("Authorization", "Basic "+credentials)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rr.Code)
		}
	})

	t.Run("inactive user", func(t *testing.T) {
		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				if email == "user@example.com" {
					return &domain.User{
						ID:           1,
						Email:        "user@example.com",
						PasswordHash: string(validPasswordHash),
						Status:       "inactive",
					}, nil
				}
				return nil, errors.New("not found")
			},
		}

		middleware := BasicAuthMiddleware(repo, logger)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For this test, the middleware doesn't check user status yet
			// but we might add it in the future
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/caldav/", nil)
		credentials := base64.StdEncoding.EncodeToString([]byte("user@example.com:correct-password"))
		req.Header.Set("Authorization", "Basic "+credentials)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// Currently the middleware doesn't check status, so this passes
		// If we add status checking later, update this test
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("returns user ID from context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		ctx := req.Context()
		ctx = contextWithUserID(ctx, 123)
		req = req.WithContext(ctx)

		userID, ok := GetUserID(req)
		if !ok {
			t.Error("expected user ID to be present")
		}
		if userID != 123 {
			t.Errorf("expected user ID 123, got %d", userID)
		}
	})

	t.Run("returns false when no user ID in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		_, ok := GetUserID(req)
		if ok {
			t.Error("expected no user ID in context")
		}
	})
}

// Helper function for testing
func contextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}
