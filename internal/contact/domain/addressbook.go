package domain

import (
	"time"
)

// Addressbook represents an addressbook collection
type Addressbook struct {
	ID          int64
	UserID      int64
	Name        string
	DisplayName string
	Description string
	SyncToken   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AddressbookRepository defines the interface for addressbook persistence
type AddressbookRepository interface {
	// Create creates a new addressbook
	Create(addressbook *Addressbook) error

	// GetByID retrieves an addressbook by ID
	GetByID(id int64) (*Addressbook, error)

	// GetByUserID retrieves all addressbooks for a user
	GetByUserID(userID int64) ([]*Addressbook, error)

	// GetByUserAndName retrieves an addressbook by user ID and name
	GetByUserAndName(userID int64, name string) (*Addressbook, error)

	// Update updates an existing addressbook
	Update(addressbook *Addressbook) error

	// Delete deletes an addressbook
	Delete(id int64) error

	// UpdateSyncToken updates the sync token for an addressbook
	UpdateSyncToken(id int64, token string) error
}

// AddressbookService defines business logic for addressbooks
type AddressbookService interface {
	// CreateAddressbook creates a new addressbook for a user
	CreateAddressbook(userID int64, name, displayName, description string) (*Addressbook, error)

	// GetAddressbook retrieves an addressbook by ID
	GetAddressbook(id int64) (*Addressbook, error)

	// GetUserAddressbooks retrieves all addressbooks for a user
	GetUserAddressbooks(userID int64) ([]*Addressbook, error)

	// UpdateAddressbook updates addressbook properties
	UpdateAddressbook(id int64, displayName, description *string) error

	// DeleteAddressbook deletes an addressbook and all its contacts
	DeleteAddressbook(id int64) error

	// GenerateSyncToken generates a new sync token for an addressbook
	GenerateSyncToken(id int64) (string, error)
}
