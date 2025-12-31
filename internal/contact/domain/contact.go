package domain

import (
	"time"
)

// Contact represents a contact (vCard)
type Contact struct {
	ID            int64
	AddressbookID int64
	UID           string
	FN            string // Formatted name (required)
	Email         string
	Tel           string
	ETag          string
	VCardData     string // Full vCard data
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ContactRepository defines the interface for contact persistence
type ContactRepository interface {
	// Create creates a new contact
	Create(contact *Contact) error

	// GetByID retrieves a contact by ID
	GetByID(id int64) (*Contact, error)

	// GetByUID retrieves a contact by UID and addressbook ID
	GetByUID(addressbookID int64, uid string) (*Contact, error)

	// GetByAddressbook retrieves all contacts for an addressbook
	GetByAddressbook(addressbookID int64) ([]*Contact, error)

	// Search searches contacts by query string
	Search(addressbookID int64, query string) ([]*Contact, error)

	// Update updates an existing contact
	Update(contact *Contact) error

	// Delete deletes a contact
	Delete(id int64) error

	// UpdateETag updates the ETag for a contact
	UpdateETag(id int64, etag string) error
}

// ContactService defines business logic for contacts
type ContactService interface {
	// CreateContact creates a new contact from vCard data
	CreateContact(addressbookID int64, vcardData string) (*Contact, error)

	// GetContact retrieves a contact by ID
	GetContact(id int64) (*Contact, error)

	// GetAddressbookContacts retrieves all contacts for an addressbook
	GetAddressbookContacts(addressbookID int64) ([]*Contact, error)

	// SearchContacts searches contacts by query
	SearchContacts(addressbookID int64, query string) ([]*Contact, error)

	// UpdateContact updates a contact from vCard data
	UpdateContact(id int64, vcardData string) error

	// DeleteContact deletes a contact
	DeleteContact(id int64) error

	// GenerateETag generates a new ETag for a contact
	GenerateETag(contact *Contact) string
}
