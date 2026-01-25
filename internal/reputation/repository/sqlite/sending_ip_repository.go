package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type sendingIPRepository struct {
	db *sql.DB
}

// NewSendingIPRepository creates a new SQLite sending IP repository
func NewSendingIPRepository(db *sql.DB) repository.SendingIPRepository {
	return &sendingIPRepository{db: db}
}

// Create stores a new sending IP configuration
func (r *sendingIPRepository) Create(ctx context.Context, ipConfig *domain.SendingIPConfig) error {
	query := `
		INSERT INTO sending_ip_configs (ip_address, domain, description, active, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		ipConfig.IPAddress,
		ipConfig.Domain,
		ipConfig.Description,
		ipConfig.Active,
		ipConfig.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create sending IP config: %w", err)
	}

	return nil
}

// GetByIP retrieves configuration for an IP address
func (r *sendingIPRepository) GetByIP(ctx context.Context, ipAddress string) (*domain.SendingIPConfig, error) {
	query := `
		SELECT ip_address, domain, description, active, created_at
		FROM sending_ip_configs
		WHERE ip_address = ?
	`

	config := &domain.SendingIPConfig{}
	err := r.db.QueryRowContext(ctx, query, ipAddress).Scan(
		&config.IPAddress,
		&config.Domain,
		&config.Description,
		&config.Active,
		&config.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get sending IP config: %w", err)
	}

	return config, nil
}

// ListAll retrieves all configured sending IPs
func (r *sendingIPRepository) ListAll(ctx context.Context) ([]*domain.SendingIPConfig, error) {
	query := `
		SELECT ip_address, domain, description, active, created_at
		FROM sending_ip_configs
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sending IP configs: %w", err)
	}
	defer rows.Close()

	var configs []*domain.SendingIPConfig
	for rows.Next() {
		config := &domain.SendingIPConfig{}
		if err := rows.Scan(
			&config.IPAddress,
			&config.Domain,
			&config.Description,
			&config.Active,
			&config.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sending IP config: %w", err)
		}
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sending IP configs: %w", err)
	}

	return configs, nil
}

// ListActive retrieves all active sending IPs
func (r *sendingIPRepository) ListActive(ctx context.Context) ([]*domain.SendingIPConfig, error) {
	query := `
		SELECT ip_address, domain, description, active, created_at
		FROM sending_ip_configs
		WHERE active = 1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active sending IP configs: %w", err)
	}
	defer rows.Close()

	var configs []*domain.SendingIPConfig
	for rows.Next() {
		config := &domain.SendingIPConfig{}
		if err := rows.Scan(
			&config.IPAddress,
			&config.Domain,
			&config.Description,
			&config.Active,
			&config.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sending IP config: %w", err)
		}
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating active sending IP configs: %w", err)
	}

	return configs, nil
}

// ListByDomain retrieves sending IPs associated with a domain
func (r *sendingIPRepository) ListByDomain(ctx context.Context, domainName string) ([]*domain.SendingIPConfig, error) {
	query := `
		SELECT ip_address, domain, description, active, created_at
		FROM sending_ip_configs
		WHERE domain = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to list sending IP configs by domain: %w", err)
	}
	defer rows.Close()

	var configs []*domain.SendingIPConfig
	for rows.Next() {
		config := &domain.SendingIPConfig{}
		if err := rows.Scan(
			&config.IPAddress,
			&config.Domain,
			&config.Description,
			&config.Active,
			&config.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sending IP config: %w", err)
		}
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sending IP configs: %w", err)
	}

	return configs, nil
}

// Update updates a sending IP configuration
func (r *sendingIPRepository) Update(ctx context.Context, ipConfig *domain.SendingIPConfig) error {
	query := `
		UPDATE sending_ip_configs
		SET domain = ?, description = ?, active = ?
		WHERE ip_address = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		ipConfig.Domain,
		ipConfig.Description,
		ipConfig.Active,
		ipConfig.IPAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to update sending IP config: %w", err)
	}

	return nil
}

// Delete removes a sending IP configuration
func (r *sendingIPRepository) Delete(ctx context.Context, ipAddress string) error {
	query := `DELETE FROM sending_ip_configs WHERE ip_address = ?`

	_, err := r.db.ExecContext(ctx, query, ipAddress)
	if err != nil {
		return fmt.Errorf("failed to delete sending IP config: %w", err)
	}

	return nil
}
