package sqlite

import (
	"database/sql"
	"time"

	"github.com/btafoya/gomailserver/internal/contact/domain"
)

// ContactRepository implements domain.ContactRepository for SQLite
type ContactRepository struct {
	db *sql.DB
}

// NewContactRepository creates a new SQLite contact repository
func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// Create creates a new contact
func (r *ContactRepository) Create(contact *domain.Contact) error {
	query := `
		INSERT INTO contacts (addressbook_id, uid, fn, email, tel, etag, vcard_data, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		contact.AddressbookID,
		contact.UID,
		contact.FN,
		contact.Email,
		contact.Tel,
		contact.ETag,
		contact.VCardData,
		contact.CreatedAt,
		contact.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	contact.ID = id

	return nil
}

// GetByID retrieves a contact by ID
func (r *ContactRepository) GetByID(id int64) (*domain.Contact, error) {
	query := `
		SELECT id, addressbook_id, uid, fn, email, tel, etag, vcard_data, created_at, updated_at
		FROM contacts
		WHERE id = ?
	`
	contact := &domain.Contact{}
	err := r.db.QueryRow(query, id).Scan(
		&contact.ID,
		&contact.AddressbookID,
		&contact.UID,
		&contact.FN,
		&contact.Email,
		&contact.Tel,
		&contact.ETag,
		&contact.VCardData,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return contact, nil
}

// GetByUID retrieves a contact by UID and addressbook ID
func (r *ContactRepository) GetByUID(addressbookID int64, uid string) (*domain.Contact, error) {
	query := `
		SELECT id, addressbook_id, uid, fn, email, tel, etag, vcard_data, created_at, updated_at
		FROM contacts
		WHERE addressbook_id = ? AND uid = ?
	`
	contact := &domain.Contact{}
	err := r.db.QueryRow(query, addressbookID, uid).Scan(
		&contact.ID,
		&contact.AddressbookID,
		&contact.UID,
		&contact.FN,
		&contact.Email,
		&contact.Tel,
		&contact.ETag,
		&contact.VCardData,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return contact, nil
}

// GetByAddressbook retrieves all contacts for an addressbook
func (r *ContactRepository) GetByAddressbook(addressbookID int64) ([]*domain.Contact, error) {
	query := `
		SELECT id, addressbook_id, uid, fn, email, tel, etag, vcard_data, created_at, updated_at
		FROM contacts
		WHERE addressbook_id = ?
		ORDER BY fn ASC
	`
	rows, err := r.db.Query(query, addressbookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*domain.Contact
	for rows.Next() {
		contact := &domain.Contact{}
		err := rows.Scan(
			&contact.ID,
			&contact.AddressbookID,
			&contact.UID,
			&contact.FN,
			&contact.Email,
			&contact.Tel,
			&contact.ETag,
			&contact.VCardData,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

// Search searches contacts by query string
func (r *ContactRepository) Search(addressbookID int64, query string) ([]*domain.Contact, error) {
	searchQuery := `
		SELECT id, addressbook_id, uid, fn, email, tel, etag, vcard_data, created_at, updated_at
		FROM contacts
		WHERE addressbook_id = ? AND (
			fn LIKE ? OR
			email LIKE ? OR
			tel LIKE ?
		)
		ORDER BY fn ASC
	`
	searchPattern := "%" + query + "%"
	rows, err := r.db.Query(searchQuery, addressbookID, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*domain.Contact
	for rows.Next() {
		contact := &domain.Contact{}
		err := rows.Scan(
			&contact.ID,
			&contact.AddressbookID,
			&contact.UID,
			&contact.FN,
			&contact.Email,
			&contact.Tel,
			&contact.ETag,
			&contact.VCardData,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

// Update updates an existing contact
func (r *ContactRepository) Update(contact *domain.Contact) error {
	query := `
		UPDATE contacts
		SET fn = ?, email = ?, tel = ?, etag = ?, vcard_data = ?, updated_at = ?
		WHERE id = ?
	`
	contact.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		contact.FN,
		contact.Email,
		contact.Tel,
		contact.ETag,
		contact.VCardData,
		contact.UpdatedAt,
		contact.ID,
	)
	return err
}

// Delete deletes a contact
func (r *ContactRepository) Delete(id int64) error {
	query := `DELETE FROM contacts WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// UpdateETag updates the ETag for a contact
func (r *ContactRepository) UpdateETag(id int64, etag string) error {
	query := `UPDATE contacts SET etag = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, etag, time.Now(), id)
	return err
}
