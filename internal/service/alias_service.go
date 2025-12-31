package service

import (
	"context"
	"encoding/json"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

// AliasService provides business logic for alias management
type AliasService struct {
	repo repository.AliasRepository
}

// NewAliasService creates a new alias service
func NewAliasService(repo repository.AliasRepository) *AliasService {
	return &AliasService{
		repo: repo,
	}
}

// Create creates a new alias
func (s *AliasService) Create(ctx context.Context, alias *domain.Alias) error {
	return s.repo.Create(ctx, alias)
}

// GetByID retrieves an alias by ID
func (s *AliasService) GetByID(ctx context.Context, id int64) (*domain.Alias, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByEmail retrieves an alias by email address
func (s *AliasService) GetByEmail(ctx context.Context, email string) (*domain.Alias, error) {
	return s.repo.GetByEmail(ctx, email)
}

// ListAll retrieves all aliases
func (s *AliasService) ListAll(ctx context.Context) ([]*domain.Alias, error) {
	return s.repo.ListAll(ctx)
}

// ListByDomain retrieves all aliases for a domain
func (s *AliasService) ListByDomain(ctx context.Context, domainID int64) ([]*domain.Alias, error) {
	return s.repo.ListByDomain(ctx, domainID)
}

// Update updates an alias
func (s *AliasService) Update(ctx context.Context, alias *domain.Alias) error {
	return s.repo.Update(ctx, alias)
}

// Delete deletes an alias
func (s *AliasService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// Helper methods for working with JSON destinations

// SetDestinations converts a slice of email addresses to JSON format for storage
func SetDestinations(destinations []string) (string, error) {
	if len(destinations) == 0 {
		return "[]", nil
	}
	bytes, err := json.Marshal(destinations)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// GetDestinations parses JSON destinations into a slice of email addresses
func GetDestinations(destinationsJSON string) ([]string, error) {
	if destinationsJSON == "" {
		return []string{}, nil
	}
	var destinations []string
	err := json.Unmarshal([]byte(destinationsJSON), &destinations)
	if err != nil {
		return nil, err
	}
	return destinations, nil
}
