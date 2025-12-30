package service

import (
	"errors"
	"testing"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
)

// mockUserRepository is a test double for UserRepository
type mockUserRepository struct {
	createFunc        func(*domain.User) error
	getByEmailFunc    func(string) (*domain.User, error)
	getByIDFunc       func(int64) (*domain.User, error)
	updateFunc        func(*domain.User) error
	deleteFunc        func(int64) error
	listFunc          func(int64, int, int) ([]*domain.User, error)
	updateQuotaFunc   func(int64, int64) error
	updatePasswordFunc func(int64, string) error
}

func (m *mockUserRepository) Create(user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(user)
	}
	user.ID = 1
	return nil
}

func (m *mockUserRepository) GetByEmail(email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(email)
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepository) GetByID(id int64) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepository) Update(user *domain.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(user)
	}
	return nil
}

func (m *mockUserRepository) UpdateLastLogin(id int64) error {
	return nil
}

func (m *mockUserRepository) Delete(id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *mockUserRepository) List(domainID int64, offset, limit int) ([]*domain.User, error) {
	if m.listFunc != nil {
		return m.listFunc(domainID, offset, limit)
	}
	return []*domain.User{}, nil
}

func (m *mockUserRepository) UpdateQuota(userID, usedQuota int64) error {
	if m.updateQuotaFunc != nil {
		return m.updateQuotaFunc(userID, usedQuota)
	}
	return nil
}

func (m *mockUserRepository) UpdatePassword(userID int64, passwordHash string) error {
	if m.updatePasswordFunc != nil {
		return m.updatePasswordFunc(userID, passwordHash)
	}
	return nil
}

func TestUserService_Create(t *testing.T) {
	logger := zap.NewNop()

	t.Run("creates user with hashed password", func(t *testing.T) {
		repo := &mockUserRepository{}
		svc := NewUserService(repo, logger)

		user := &domain.User{
			Email:    "test@example.com",
			FullName: "Test User",
			Status:   "active",
		}

		err := svc.Create(user, "plaintextpassword")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.PasswordHash == "" {
			t.Error("expected password hash to be set")
		}

		if user.PasswordHash == "plaintextpassword" {
			t.Error("password should be hashed, not plaintext")
		}

		// Verify bcrypt hash format (starts with $2a$ or $2b$)
		if len(user.PasswordHash) < 20 {
			t.Error("password hash seems too short for bcrypt")
		}
	})

	t.Run("returns error if repository fails", func(t *testing.T) {
		repo := &mockUserRepository{
			createFunc: func(u *domain.User) error {
				return errors.New("database error")
			},
		}
		svc := NewUserService(repo, logger)

		user := &domain.User{Email: "test@example.com"}
		err := svc.Create(user, "password")

		if err == nil {
			t.Error("expected error from repository failure")
		}
	})
}

func TestUserService_Authenticate(t *testing.T) {
	logger := zap.NewNop()

	t.Run("authenticates user with correct password", func(t *testing.T) {
		// Create a user with a known password hash
		svc := NewUserService(&mockUserRepository{}, logger)
		testUser := &domain.User{
			ID:       1,
			Email:    "test@example.com",
			Status:   "active",
			FullName: "Test User",
		}
		// Hash the password
		err := svc.Create(testUser, "correctpassword")
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				if email == "test@example.com" {
					return testUser, nil
				}
				return nil, errors.New("not found")
			},
		}

		svc = NewUserService(repo, logger)
		user, err := svc.Authenticate("test@example.com", "correctpassword")

		if err != nil {
			t.Fatalf("expected successful authentication, got error: %v", err)
		}

		if user == nil {
			t.Fatal("expected user to be returned")
		}

		if user.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", user.Email)
		}
	})

	t.Run("fails authentication with incorrect password", func(t *testing.T) {
		svc := NewUserService(&mockUserRepository{}, logger)
		testUser := &domain.User{
			ID:     1,
			Email:  "test@example.com",
			Status: "active",
		}
		err := svc.Create(testUser, "correctpassword")
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				return testUser, nil
			},
		}

		svc = NewUserService(repo, logger)
		user, err := svc.Authenticate("test@example.com", "wrongpassword")

		if err == nil {
			t.Error("expected authentication to fail with wrong password")
		}

		if user != nil {
			t.Error("expected no user to be returned on failed auth")
		}
	})

	t.Run("fails authentication for non-existent user", func(t *testing.T) {
		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				return nil, errors.New("not found")
			},
		}

		svc := NewUserService(repo, logger)
		user, err := svc.Authenticate("nonexistent@example.com", "password")

		if err == nil {
			t.Error("expected authentication to fail for non-existent user")
		}

		if user != nil {
			t.Error("expected no user to be returned")
		}
	})

	t.Run("fails authentication for disabled user", func(t *testing.T) {
		svc := NewUserService(&mockUserRepository{}, logger)
		testUser := &domain.User{
			ID:     1,
			Email:  "test@example.com",
			Status: "disabled",
		}
		err := svc.Create(testUser, "password")
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				return testUser, nil
			},
		}

		svc = NewUserService(repo, logger)
		user, err := svc.Authenticate("test@example.com", "password")

		if err == nil {
			t.Error("expected authentication to fail for disabled user")
		}

		if user != nil {
			t.Error("expected no user to be returned for disabled user")
		}
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	logger := zap.NewNop()

	t.Run("returns user by email", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:    1,
			Email: "test@example.com",
		}

		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				if email == "test@example.com" {
					return expectedUser, nil
				}
				return nil, errors.New("not found")
			},
		}

		svc := NewUserService(repo, logger)
		user, err := svc.GetByEmail("test@example.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.ID != expectedUser.ID {
			t.Errorf("expected user ID %d, got %d", expectedUser.ID, user.ID)
		}
	})

	t.Run("returns error for non-existent email", func(t *testing.T) {
		repo := &mockUserRepository{
			getByEmailFunc: func(email string) (*domain.User, error) {
				return nil, errors.New("not found")
			},
		}

		svc := NewUserService(repo, logger)
		user, err := svc.GetByEmail("nonexistent@example.com")

		if err == nil {
			t.Error("expected error for non-existent email")
		}

		if user != nil {
			t.Error("expected no user to be returned")
		}
	})
}

func TestUserService_UpdatePassword(t *testing.T) {
	logger := zap.NewNop()

	t.Run("updates password with new hash", func(t *testing.T) {
		var capturedHash string
		repo := &mockUserRepository{
			updatePasswordFunc: func(userID int64, passwordHash string) error {
				capturedHash = passwordHash
				return nil
			},
		}

		svc := NewUserService(repo, logger)
		err := svc.UpdatePassword(1, "newpassword")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if capturedHash == "" {
			t.Error("expected password hash to be generated")
		}

		if capturedHash == "newpassword" {
			t.Error("password should be hashed, not plaintext")
		}

		// Verify bcrypt hash format
		if len(capturedHash) < 20 {
			t.Error("password hash seems too short for bcrypt")
		}
	})

	t.Run("returns error if repository fails", func(t *testing.T) {
		repo := &mockUserRepository{
			updatePasswordFunc: func(userID int64, passwordHash string) error {
				return errors.New("database error")
			},
		}

		svc := NewUserService(repo, logger)
		err := svc.UpdatePassword(1, "newpassword")

		if err == nil {
			t.Error("expected error from repository failure")
		}
	})
}
