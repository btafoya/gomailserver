package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/contact/domain"
	"github.com/emersion/go-vcard"
)

// ContactService implements domain.ContactService
type ContactService struct {
	contactRepo     domain.ContactRepository
	addressbookRepo domain.AddressbookRepository
}

// NewContactService creates a new contact service
func NewContactService(contactRepo domain.ContactRepository, addressbookRepo domain.AddressbookRepository) *ContactService {
	return &ContactService{
		contactRepo:     contactRepo,
		addressbookRepo: addressbookRepo,
	}
}

// CreateContact creates a new contact from vCard data
func (s *ContactService) CreateContact(addressbookID int64, vcardData string) (*domain.Contact, error) {
	// Verify addressbook exists
	addressbook, err := s.addressbookRepo.GetByID(addressbookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get addressbook: %w", err)
	}
	if addressbook == nil {
		return nil, fmt.Errorf("addressbook not found")
	}

	// Parse vCard data
	dec := vcard.NewDecoder(strings.NewReader(vcardData))
	card, err := dec.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to parse vCard data: %w", err)
	}

	// Extract contact properties
	contact := &domain.Contact{
		AddressbookID: addressbookID,
		VCardData:     vcardData,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Extract UID (required)
	if uid := card.Get(vcard.FieldUID); uid != nil {
		contact.UID = uid.Value
	} else {
		return nil, fmt.Errorf("UID field is required")
	}

	// Check if contact with same UID already exists
	existing, err := s.contactRepo.GetByUID(addressbookID, contact.UID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing contact: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("contact with UID %s already exists", contact.UID)
	}

	// Extract FN (formatted name - required)
	if fn := card.Get(vcard.FieldFormattedName); fn != nil {
		contact.FN = fn.Value
	} else {
		return nil, fmt.Errorf("FN field is required")
	}

	// Extract email (optional)
	if email := card.Get(vcard.FieldEmail); email != nil {
		contact.Email = email.Value
	}

	// Extract telephone (optional)
	if tel := card.Get(vcard.FieldTelephone); tel != nil {
		contact.Tel = tel.Value
	}

	// Generate ETag
	contact.ETag = s.GenerateETag(contact)

	// Create contact in repository
	if err := s.contactRepo.Create(contact); err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}

	return contact, nil
}

// GetContact retrieves a contact by ID
func (s *ContactService) GetContact(id int64) (*domain.Contact, error) {
	contact, err := s.contactRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}
	if contact == nil {
		return nil, fmt.Errorf("contact not found")
	}
	return contact, nil
}

// GetAddressbookContacts retrieves all contacts for an addressbook
func (s *ContactService) GetAddressbookContacts(addressbookID int64) ([]*domain.Contact, error) {
	contacts, err := s.contactRepo.GetByAddressbook(addressbookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get addressbook contacts: %w", err)
	}
	return contacts, nil
}

// SearchContacts searches contacts by query
func (s *ContactService) SearchContacts(addressbookID int64, query string) ([]*domain.Contact, error) {
	contacts, err := s.contactRepo.Search(addressbookID, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search contacts: %w", err)
	}
	return contacts, nil
}

// UpdateContact updates a contact from vCard data
func (s *ContactService) UpdateContact(id int64, vcardData string) error {
	// Get existing contact
	contact, err := s.contactRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get contact: %w", err)
	}
	if contact == nil {
		return fmt.Errorf("contact not found")
	}

	// Parse new vCard data
	dec := vcard.NewDecoder(strings.NewReader(vcardData))
	card, err := dec.Decode()
	if err != nil {
		return fmt.Errorf("failed to parse vCard data: %w", err)
	}

	// Update contact properties
	contact.VCardData = vcardData
	contact.UpdatedAt = time.Now()

	// Extract FN (formatted name - required)
	if fn := card.Get(vcard.FieldFormattedName); fn != nil {
		contact.FN = fn.Value
	}

	// Extract email (optional)
	if email := card.Get(vcard.FieldEmail); email != nil {
		contact.Email = email.Value
	} else {
		contact.Email = ""
	}

	// Extract telephone (optional)
	if tel := card.Get(vcard.FieldTelephone); tel != nil {
		contact.Tel = tel.Value
	} else {
		contact.Tel = ""
	}

	// Update ETag
	contact.ETag = s.GenerateETag(contact)

	// Update in repository
	if err := s.contactRepo.Update(contact); err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}

	return nil
}

// DeleteContact deletes a contact
func (s *ContactService) DeleteContact(id int64) error {
	if err := s.contactRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	return nil
}

// GenerateETag generates a new ETag for a contact
func (s *ContactService) GenerateETag(contact *domain.Contact) string {
	// Generate ETag based on contact content and update time
	data := fmt.Sprintf("%s-%s", contact.UID, contact.UpdatedAt.Format(time.RFC3339))
	hash := sha256.Sum256([]byte(data))
	return `"` + hex.EncodeToString(hash[:]) + `"`
}
