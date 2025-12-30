package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

const bcryptCost = 12

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user disabled")
	ErrUserNotFound       = errors.New("user not found")
)

// UserService handles user operations
type UserService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// Authenticate verifies user credentials
func (s *UserService) Authenticate(email, password string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		s.logger.Debug("authentication failed - user not found",
			zap.String("email", email),
		)
		// Return generic error to prevent user enumeration
		return nil, ErrInvalidCredentials
	}

	if user.Status != "active" {
		s.logger.Warn("authentication failed - user disabled",
			zap.String("email", email),
			zap.String("status", user.Status),
		)
		return nil, ErrUserDisabled
	}

	if !s.VerifyPassword(user.PasswordHash, password) {
		s.logger.Warn("authentication failed - invalid password",
			zap.String("email", email),
		)
		return nil, ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := s.repo.UpdateLastLogin(user.ID); err != nil {
		s.logger.Error("failed to update last login",
			zap.Error(err),
			zap.Int64("user_id", user.ID),
		)
	}

	s.logger.Info("authentication successful",
		zap.String("email", email),
		zap.Int64("user_id", user.ID),
	)

	return user, nil
}

// HashPassword hashes a password using bcrypt
func (s *UserService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func (s *UserService) VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Create creates a new user
func (s *UserService) Create(user *domain.User, password string) error {
	hash, err := s.HashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = hash

	if err := s.repo.Create(user); err != nil {
		s.logger.Error("failed to create user",
			zap.Error(err),
			zap.String("email", user.Email),
		)
		return err
	}

	s.logger.Info("user created",
		zap.String("email", user.Email),
		zap.Int64("user_id", user.ID),
	)

	return nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id int64) (*domain.User, error) {
	return s.repo.GetByID(id)
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	return s.repo.GetByEmail(email)
}

// Update updates a user
func (s *UserService) Update(user *domain.User) error {
	return s.repo.Update(user)
}

// Delete deletes a user
func (s *UserService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(userID int64, newPassword string) error {
	hash, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(userID, hash); err != nil {
		s.logger.Error("failed to update password",
			zap.Error(err),
			zap.Int64("user_id", userID),
		)
		return err
	}

	s.logger.Info("password updated",
		zap.Int64("user_id", userID),
	)

	return nil
}
