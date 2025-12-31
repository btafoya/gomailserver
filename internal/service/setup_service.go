package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// SetupService handles setup wizard operations
type SetupService struct {
	db         *database.DB
	userRepo   repository.UserRepository
	domainRepo repository.DomainRepository
	logger     *zap.Logger
}

// SetupState represents the setup wizard state
type SetupState struct {
	CurrentStep    string   `json:"current_step"`
	CompletedSteps []string `json:"completed_steps"`
	SystemConfig   *string  `json:"system_config,omitempty"`
	DomainConfig   *string  `json:"domain_config,omitempty"`
	AdminConfig    *string  `json:"admin_config,omitempty"`
	TLSConfig      *string  `json:"tls_config,omitempty"`
}

// AdminUserRequest represents a request to create an admin user
type AdminUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// NewSetupService creates a new setup service
func NewSetupService(
	db *database.DB,
	userRepo repository.UserRepository,
	domainRepo repository.DomainRepository,
	logger *zap.Logger,
) *SetupService {
	return &SetupService{
		db:         db,
		userRepo:   userRepo,
		domainRepo: domainRepo,
		logger:     logger,
	}
}

// IsSetupComplete checks if the setup wizard has been completed
func (s *SetupService) IsSetupComplete(ctx context.Context) (bool, error) {
	var currentStep string
	err := s.db.QueryRowContext(ctx, "SELECT current_step FROM setup_wizard_state WHERE id = 1").Scan(&currentStep)
	if err == sql.ErrNoRows {
		// Setup wizard state doesn't exist, setup not complete
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check setup status: %w", err)
	}

	return currentStep == "complete", nil
}

// GetSetupState retrieves the current setup wizard state
func (s *SetupService) GetSetupState(ctx context.Context) (*SetupState, error) {
	state := &SetupState{}

	query := `
		SELECT current_step, completed_steps, system_config, domain_config, admin_config, tls_config
		FROM setup_wizard_state
		WHERE id = 1
	`

	var completedStepsJSON string
	err := s.db.QueryRowContext(ctx, query).Scan(
		&state.CurrentStep,
		&completedStepsJSON,
		&state.SystemConfig,
		&state.DomainConfig,
		&state.AdminConfig,
		&state.TLSConfig,
	)
	if err == sql.ErrNoRows {
		// Return default state if no record exists
		return &SetupState{
			CurrentStep:    "welcome",
			CompletedSteps: []string{},
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get setup state: %w", err)
	}

	// Parse completed steps JSON
	if err := json.Unmarshal([]byte(completedStepsJSON), &state.CompletedSteps); err != nil {
		return nil, fmt.Errorf("failed to parse completed steps: %w", err)
	}

	return state, nil
}

// CreateAdminUser creates the first admin user and their domain
func (s *SetupService) CreateAdminUser(ctx context.Context, req *AdminUserRequest) error {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return fmt.Errorf("email and password are required")
	}

	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Extract domain from email
	parts := splitEmail(req.Email)
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}
	domainName := parts[1]

	// Check if domain exists
	existingDomain, err := s.domainRepo.GetByName(domainName)
	if err != nil && !isDomainNotFoundError(err) {
		return fmt.Errorf("failed to check domain: %w", err)
	}

	var domainID int64
	if existingDomain != nil {
		domainID = existingDomain.ID
	} else {
		// Create domain from template to ensure all security settings are initialized
		domainSvc := NewDomainService(s.domainRepo)
		newDomain, err := domainSvc.CreateDomainFromTemplate(domainName)
		if err != nil {
			return fmt.Errorf("failed to create domain: %w", err)
		}
		domainID = newDomain.ID
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && !isUserNotFoundError(err) {
		return fmt.Errorf("failed to check user: %w", err)
	}

	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	user := &domain.User{
		Email:        req.Email,
		DomainID:     domainID,
		PasswordHash: string(passwordHash),
		FullName:     req.FullName,
		DisplayName:  req.FullName,
		Role:         "admin",
		Quota:        1073741824, // 1GB default
		UsedQuota:    0,
		Status:       "active",
		AuthMethod:   "password",
		TOTPEnabled:  false,
		Language:     "en",
	}

	if err := s.userRepo.Create(user); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	// Update setup state to mark admin step as complete
	if err := s.updateSetupStep(ctx, "admin", "tls"); err != nil {
		s.logger.Warn("failed to update setup state", zap.Error(err))
	}

	s.logger.Info("admin user created successfully",
		zap.String("email", req.Email),
		zap.String("domain", domainName),
	)

	return nil
}

// CompleteSetup marks the setup wizard as complete
func (s *SetupService) CompleteSetup(ctx context.Context) error {
	query := `
		UPDATE setup_wizard_state
		SET current_step = 'complete',
		    completed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to complete setup: %w", err)
	}

	s.logger.Info("setup wizard completed")
	return nil
}

// updateSetupStep updates the current step and adds to completed steps
func (s *SetupService) updateSetupStep(ctx context.Context, completedStep, nextStep string) error {
	state, err := s.GetSetupState(ctx)
	if err != nil {
		return err
	}

	// Add to completed steps if not already there
	found := false
	for _, step := range state.CompletedSteps {
		if step == completedStep {
			found = true
			break
		}
	}

	if !found {
		state.CompletedSteps = append(state.CompletedSteps, completedStep)
	}

	completedStepsJSON, err := json.Marshal(state.CompletedSteps)
	if err != nil {
		return fmt.Errorf("failed to marshal completed steps: %w", err)
	}

	query := `
		UPDATE setup_wizard_state
		SET current_step = ?,
		    completed_steps = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`

	if _, err := s.db.ExecContext(ctx, query, nextStep, string(completedStepsJSON)); err != nil {
		return fmt.Errorf("failed to update setup step: %w", err)
	}

	return nil
}

// Helper functions
func splitEmail(email string) []string {
	result := make([]string, 0, 2)
	atIndex := -1
	for i, ch := range email {
		if ch == '@' {
			atIndex = i
			break
		}
	}

	if atIndex == -1 {
		return result
	}

	result = append(result, email[:atIndex])
	result = append(result, email[atIndex+1:])
	return result
}

func isDomainNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return containsString(err.Error(), "not found")
}

func isUserNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return containsString(err.Error(), "not found")
}

func containsString(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
