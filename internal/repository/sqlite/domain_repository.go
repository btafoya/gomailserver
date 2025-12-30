package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type domainRepository struct {
	db *database.DB
}

// NewDomainRepository creates a new SQLite domain repository
func NewDomainRepository(db *database.DB) repository.DomainRepository {
	return &domainRepository{db: db}
}

// Create inserts a new domain
func (r *domainRepository) Create(dom *domain.Domain) error {
	query := `
		INSERT INTO domains (
			name, status, max_users, max_mailbox_size, default_quota,
			catchall_email, backup_mx,
			dkim_selector, dkim_private_key, dkim_public_key,
			spf_record, dmarc_policy,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		dom.Name, dom.Status, dom.MaxUsers, dom.MaxMailboxSize, dom.DefaultQuota,
		dom.CatchallEmail, dom.BackupMX,
		dom.DKIMSelector, dom.DKIMPrivateKey, dom.DKIMPublicKey,
		dom.SPFRecord, dom.DMARCPolicy,
		time.Now(), time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create domain: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get domain ID: %w", err)
	}

	dom.ID = id
	dom.CreatedAt = time.Now()
	dom.UpdatedAt = time.Now()

	return nil
}

// GetByID retrieves a domain by ID
func (r *domainRepository) GetByID(id int64) (*domain.Domain, error) {
	query := `
		SELECT
			id, name, status, max_users, max_mailbox_size, default_quota,
			catchall_email, backup_mx,
			dkim_selector, dkim_private_key, dkim_public_key,
			spf_record, dmarc_policy,
			created_at, updated_at
		FROM domains
		WHERE id = ?
	`

	dom := &domain.Domain{}

	err := r.db.QueryRow(query, id).Scan(
		&dom.ID, &dom.Name, &dom.Status, &dom.MaxUsers, &dom.MaxMailboxSize, &dom.DefaultQuota,
		&dom.CatchallEmail, &dom.BackupMX,
		&dom.DKIMSelector, &dom.DKIMPrivateKey, &dom.DKIMPublicKey,
		&dom.SPFRecord, &dom.DMARCPolicy,
		&dom.CreatedAt, &dom.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("domain not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return dom, nil
}

// GetByName retrieves a domain by name
func (r *domainRepository) GetByName(name string) (*domain.Domain, error) {
	query := `
		SELECT
			id, name, status, max_users, max_mailbox_size, default_quota,
			catchall_email, backup_mx,
			dkim_selector, dkim_private_key, dkim_public_key,
			spf_record, dmarc_policy,
			created_at, updated_at
		FROM domains
		WHERE name = ?
	`

	dom := &domain.Domain{}

	err := r.db.QueryRow(query, name).Scan(
		&dom.ID, &dom.Name, &dom.Status, &dom.MaxUsers, &dom.MaxMailboxSize, &dom.DefaultQuota,
		&dom.CatchallEmail, &dom.BackupMX,
		&dom.DKIMSelector, &dom.DKIMPrivateKey, &dom.DKIMPublicKey,
		&dom.SPFRecord, &dom.DMARCPolicy,
		&dom.CreatedAt, &dom.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("domain not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return dom, nil
}

// Update updates a domain
func (r *domainRepository) Update(dom *domain.Domain) error {
	query := `
		UPDATE domains SET
			name = ?, status = ?, max_users = ?, max_mailbox_size = ?, default_quota = ?,
			catchall_email = ?, backup_mx = ?,
			dkim_selector = ?, dkim_private_key = ?, dkim_public_key = ?,
			spf_record = ?, dmarc_policy = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		dom.Name, dom.Status, dom.MaxUsers, dom.MaxMailboxSize, dom.DefaultQuota,
		dom.CatchallEmail, dom.BackupMX,
		dom.DKIMSelector, dom.DKIMPrivateKey, dom.DKIMPublicKey,
		dom.SPFRecord, dom.DMARCPolicy,
		time.Now(), dom.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update domain: %w", err)
	}

	dom.UpdatedAt = time.Now()
	return nil
}

// Delete deletes a domain
func (r *domainRepository) Delete(id int64) error {
	query := `DELETE FROM domains WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}
	return nil
}

// List lists domains with pagination
func (r *domainRepository) List(offset, limit int) ([]*domain.Domain, error) {
	query := `
		SELECT
			id, name, status, max_users, max_mailbox_size, default_quota,
			catchall_email, backup_mx,
			dkim_selector, dkim_private_key, dkim_public_key,
			spf_record, dmarc_policy,
			created_at, updated_at
		FROM domains
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}
	defer rows.Close()

	domains := make([]*domain.Domain, 0)
	for rows.Next() {
		dom := &domain.Domain{}

		err := rows.Scan(
			&dom.ID, &dom.Name, &dom.Status, &dom.MaxUsers, &dom.MaxMailboxSize, &dom.DefaultQuota,
			&dom.CatchallEmail, &dom.BackupMX,
			&dom.DKIMSelector, &dom.DKIMPrivateKey, &dom.DKIMPublicKey,
			&dom.SPFRecord, &dom.DMARCPolicy,
			&dom.CreatedAt, &dom.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan domain: %w", err)
		}

		domains = append(domains, dom)
	}

	return domains, rows.Err()
}
