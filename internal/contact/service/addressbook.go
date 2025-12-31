package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/contact/domain"
)

// AddressbookService implements domain.AddressbookService
type AddressbookService struct {
	addressbookRepo domain.AddressbookRepository
	contactRepo     domain.ContactRepository
}

// NewAddressbookService creates a new addressbook service
func NewAddressbookService(addressbookRepo domain.AddressbookRepository, contactRepo domain.ContactRepository) *AddressbookService {
	return &AddressbookService{
		addressbookRepo: addressbookRepo,
		contactRepo:     contactRepo,
	}
}

// CreateAddressbook creates a new addressbook for a user
func (s *AddressbookService) CreateAddressbook(userID int64, name, displayName, description string) (*domain.Addressbook, error) {
	// Check if addressbook with same name already exists
	existing, err := s.addressbookRepo.GetByUserAndName(userID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing addressbook: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("addressbook with name %s already exists", name)
	}

	// Set defaults
	if displayName == "" {
		displayName = name
	}

	// Generate initial sync token
	syncToken, err := s.generateSyncToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate sync token: %w", err)
	}

	addressbook := &domain.Addressbook{
		UserID:      userID,
		Name:        name,
		DisplayName: displayName,
		Description: description,
		SyncToken:   syncToken,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.addressbookRepo.Create(addressbook); err != nil {
		return nil, fmt.Errorf("failed to create addressbook: %w", err)
	}

	return addressbook, nil
}

// GetAddressbook retrieves an addressbook by ID
func (s *AddressbookService) GetAddressbook(id int64) (*domain.Addressbook, error) {
	addressbook, err := s.addressbookRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get addressbook: %w", err)
	}
	if addressbook == nil {
		return nil, fmt.Errorf("addressbook not found")
	}
	return addressbook, nil
}

// GetUserAddressbooks retrieves all addressbooks for a user
func (s *AddressbookService) GetUserAddressbooks(userID int64) ([]*domain.Addressbook, error) {
	addressbooks, err := s.addressbookRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user addressbooks: %w", err)
	}
	return addressbooks, nil
}

// UpdateAddressbook updates addressbook properties
func (s *AddressbookService) UpdateAddressbook(id int64, displayName, description *string) error {
	addressbook, err := s.addressbookRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get addressbook: %w", err)
	}
	if addressbook == nil {
		return fmt.Errorf("addressbook not found")
	}

	// Update fields if provided
	if displayName != nil {
		addressbook.DisplayName = *displayName
	}
	if description != nil {
		addressbook.Description = *description
	}

	if err := s.addressbookRepo.Update(addressbook); err != nil {
		return fmt.Errorf("failed to update addressbook: %w", err)
	}

	return nil
}

// DeleteAddressbook deletes an addressbook and all its contacts
func (s *AddressbookService) DeleteAddressbook(id int64) error {
	// Repository handles cascade delete of contacts
	if err := s.addressbookRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete addressbook: %w", err)
	}
	return nil
}

// GenerateSyncToken generates a new sync token for an addressbook
func (s *AddressbookService) GenerateSyncToken(id int64) (string, error) {
	token, err := s.generateSyncToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate sync token: %w", err)
	}

	if err := s.addressbookRepo.UpdateSyncToken(id, token); err != nil {
		return "", fmt.Errorf("failed to update sync token: %w", err)
	}

	return token, nil
}

// generateSyncToken generates a random sync token
func (s *AddressbookService) generateSyncToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
