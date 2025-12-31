package sqlite

import (
	"database/sql"
	"time"

	"github.com/btafoya/gomailserver/internal/contact/domain"
)

// AddressbookRepository implements domain.AddressbookRepository for SQLite
type AddressbookRepository struct {
	db *sql.DB
}

// NewAddressbookRepository creates a new SQLite addressbook repository
func NewAddressbookRepository(db *sql.DB) *AddressbookRepository {
	return &AddressbookRepository{db: db}
}

// Create creates a new addressbook
func (r *AddressbookRepository) Create(addressbook *domain.Addressbook) error {
	query := `
		INSERT INTO addressbooks (user_id, name, display_name, description, sync_token, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		addressbook.UserID,
		addressbook.Name,
		addressbook.DisplayName,
		addressbook.Description,
		addressbook.SyncToken,
		addressbook.CreatedAt,
		addressbook.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	addressbook.ID = id

	return nil
}

// GetByID retrieves an addressbook by ID
func (r *AddressbookRepository) GetByID(id int64) (*domain.Addressbook, error) {
	query := `
		SELECT id, user_id, name, display_name, description, sync_token, created_at, updated_at
		FROM addressbooks
		WHERE id = ?
	`
	addressbook := &domain.Addressbook{}
	err := r.db.QueryRow(query, id).Scan(
		&addressbook.ID,
		&addressbook.UserID,
		&addressbook.Name,
		&addressbook.DisplayName,
		&addressbook.Description,
		&addressbook.SyncToken,
		&addressbook.CreatedAt,
		&addressbook.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return addressbook, nil
}

// GetByUserID retrieves all addressbooks for a user
func (r *AddressbookRepository) GetByUserID(userID int64) ([]*domain.Addressbook, error) {
	query := `
		SELECT id, user_id, name, display_name, description, sync_token, created_at, updated_at
		FROM addressbooks
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addressbooks []*domain.Addressbook
	for rows.Next() {
		addressbook := &domain.Addressbook{}
		err := rows.Scan(
			&addressbook.ID,
			&addressbook.UserID,
			&addressbook.Name,
			&addressbook.DisplayName,
			&addressbook.Description,
			&addressbook.SyncToken,
			&addressbook.CreatedAt,
			&addressbook.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addressbooks = append(addressbooks, addressbook)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addressbooks, nil
}

// GetByUserAndName retrieves an addressbook by user ID and name
func (r *AddressbookRepository) GetByUserAndName(userID int64, name string) (*domain.Addressbook, error) {
	query := `
		SELECT id, user_id, name, display_name, description, sync_token, created_at, updated_at
		FROM addressbooks
		WHERE user_id = ? AND name = ?
	`
	addressbook := &domain.Addressbook{}
	err := r.db.QueryRow(query, userID, name).Scan(
		&addressbook.ID,
		&addressbook.UserID,
		&addressbook.Name,
		&addressbook.DisplayName,
		&addressbook.Description,
		&addressbook.SyncToken,
		&addressbook.CreatedAt,
		&addressbook.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return addressbook, nil
}

// Update updates an existing addressbook
func (r *AddressbookRepository) Update(addressbook *domain.Addressbook) error {
	query := `
		UPDATE addressbooks
		SET display_name = ?, description = ?, sync_token = ?, updated_at = ?
		WHERE id = ?
	`
	addressbook.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		addressbook.DisplayName,
		addressbook.Description,
		addressbook.SyncToken,
		addressbook.UpdatedAt,
		addressbook.ID,
	)
	return err
}

// Delete deletes an addressbook
func (r *AddressbookRepository) Delete(id int64) error {
	query := `DELETE FROM addressbooks WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// UpdateSyncToken updates the sync token for an addressbook
func (r *AddressbookRepository) UpdateSyncToken(id int64, token string) error {
	query := `UPDATE addressbooks SET sync_token = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, token, time.Now(), id)
	return err
}
