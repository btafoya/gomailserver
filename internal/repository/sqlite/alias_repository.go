package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type aliasRepository struct {
	db *database.DB
}

// NewAliasRepository creates a new SQLite alias repository
func NewAliasRepository(db *database.DB) repository.AliasRepository {
	return &aliasRepository{db: db}
}

// Create inserts a new alias
func (r *aliasRepository) Create(alias *domain.Alias) error {
	query := `
		INSERT INTO aliases (
			alias_email, domain_id, destination_emails, status, created_at
		) VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		alias.AliasEmail, alias.DomainID, alias.DestinationEmails, alias.Status, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create alias: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get alias ID: %w", err)
	}

	alias.ID = id
	alias.CreatedAt = time.Now()

	return nil
}

// GetByID retrieves an alias by ID
func (r *aliasRepository) GetByID(id int64) (*domain.Alias, error) {
	query := `
		SELECT id, alias_email, domain_id, destination_emails, status, created_at
		FROM aliases
		WHERE id = ?
	`

	alias := &domain.Alias{}
	err := r.db.QueryRow(query, id).Scan(
		&alias.ID, &alias.AliasEmail, &alias.DomainID, &alias.DestinationEmails,
		&alias.Status, &alias.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("alias not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alias: %w", err)
	}

	return alias, nil
}

// GetByEmail retrieves an alias by email address
func (r *aliasRepository) GetByEmail(email string) (*domain.Alias, error) {
	query := `
		SELECT id, alias_email, domain_id, destination_emails, status, created_at
		FROM aliases
		WHERE alias_email = ?
	`

	alias := &domain.Alias{}
	err := r.db.QueryRow(query, email).Scan(
		&alias.ID, &alias.AliasEmail, &alias.DomainID, &alias.DestinationEmails,
		&alias.Status, &alias.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("alias not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alias: %w", err)
	}

	return alias, nil
}

// Update updates an alias
func (r *aliasRepository) Update(alias *domain.Alias) error {
	query := `
		UPDATE aliases SET
			alias_email = ?, domain_id = ?, destination_emails = ?, status = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		alias.AliasEmail, alias.DomainID, alias.DestinationEmails, alias.Status, alias.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update alias: %w", err)
	}

	return nil
}

// Delete deletes an alias
func (r *aliasRepository) Delete(id int64) error {
	query := `DELETE FROM aliases WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete alias: %w", err)
	}
	return nil
}

// ListAll retrieves all aliases
func (r *aliasRepository) ListAll() ([]*domain.Alias, error) {
	query := `
		SELECT id, alias_email, domain_id, destination_emails, status, created_at
		FROM aliases
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all aliases: %w", err)
	}
	defer rows.Close()

	aliases := make([]*domain.Alias, 0)
	for rows.Next() {
		alias := &domain.Alias{}
		err := rows.Scan(
			&alias.ID, &alias.AliasEmail, &alias.DomainID, &alias.DestinationEmails,
			&alias.Status, &alias.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alias: %w", err)
		}
		aliases = append(aliases, alias)
	}

	return aliases, rows.Err()
}

// ListByDomain retrieves all aliases for a domain
func (r *aliasRepository) ListByDomain(domainID int64) ([]*domain.Alias, error) {
	query := `
		SELECT id, alias_email, domain_id, destination_emails, status, created_at
		FROM aliases
		WHERE domain_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to list aliases by domain: %w", err)
	}
	defer rows.Close()

	aliases := make([]*domain.Alias, 0)
	for rows.Next() {
		alias := &domain.Alias{}
		err := rows.Scan(
			&alias.ID, &alias.AliasEmail, &alias.DomainID, &alias.DestinationEmails,
			&alias.Status, &alias.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alias: %w", err)
		}
		aliases = append(aliases, alias)
	}

	return aliases, rows.Err()
}
