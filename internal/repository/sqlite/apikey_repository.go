package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type apiKeyRepository struct {
	db *database.DB
}

// NewAPIKeyRepository creates a new SQLite API key repository
func NewAPIKeyRepository(db *database.DB) repository.APIKeyRepository {
	return &apiKeyRepository{db: db}
}

// Create inserts a new API key
func (r *apiKeyRepository) Create(apiKey *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (
			user_id, domain_id, name, key_hash, scopes,
			expires_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		apiKey.UserID, apiKey.DomainID, apiKey.Name, apiKey.KeyHash, apiKey.Scopes,
		apiKey.ExpiresAt, time.Now(), time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get API key ID: %w", err)
	}

	apiKey.ID = id
	apiKey.CreatedAt = time.Now()
	apiKey.UpdatedAt = time.Now()

	return nil
}

// GetByKeyHash retrieves an API key by its hash
func (r *apiKeyRepository) GetByKeyHash(keyHash string) (*domain.APIKey, error) {
	query := `
		SELECT
			id, user_id, domain_id, name, key_hash, scopes,
			last_used_at, last_used_ip, expires_at, created_at, updated_at
		FROM api_keys
		WHERE key_hash = ?
	`

	apiKey := &domain.APIKey{}
	var lastUsedAt, expiresAt sql.NullTime
	var lastUsedIP sql.NullString

	err := r.db.QueryRow(query, keyHash).Scan(
		&apiKey.ID, &apiKey.UserID, &apiKey.DomainID, &apiKey.Name, &apiKey.KeyHash, &apiKey.Scopes,
		&lastUsedAt, &lastUsedIP, &expiresAt, &apiKey.CreatedAt, &apiKey.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if lastUsedAt.Valid {
		apiKey.LastUsedAt = &lastUsedAt.Time
	}
	if lastUsedIP.Valid {
		apiKey.LastUsedIP = lastUsedIP.String
	}
	if expiresAt.Valid {
		apiKey.ExpiresAt = &expiresAt.Time
	}

	return apiKey, nil
}

// GetByID retrieves an API key by ID
func (r *apiKeyRepository) GetByID(id int64) (*domain.APIKey, error) {
	query := `
		SELECT
			id, user_id, domain_id, name, key_hash, scopes,
			last_used_at, last_used_ip, expires_at, created_at, updated_at
		FROM api_keys
		WHERE id = ?
	`

	apiKey := &domain.APIKey{}
	var lastUsedAt, expiresAt sql.NullTime
	var lastUsedIP sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&apiKey.ID, &apiKey.UserID, &apiKey.DomainID, &apiKey.Name, &apiKey.KeyHash, &apiKey.Scopes,
		&lastUsedAt, &lastUsedIP, &expiresAt, &apiKey.CreatedAt, &apiKey.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if lastUsedAt.Valid {
		apiKey.LastUsedAt = &lastUsedAt.Time
	}
	if lastUsedIP.Valid {
		apiKey.LastUsedIP = lastUsedIP.String
	}
	if expiresAt.Valid {
		apiKey.ExpiresAt = &expiresAt.Time
	}

	return apiKey, nil
}

// ListByUser retrieves all API keys for a user
func (r *apiKeyRepository) ListByUser(userID int64) ([]*domain.APIKey, error) {
	query := `
		SELECT
			id, user_id, domain_id, name, key_hash, scopes,
			last_used_at, last_used_ip, expires_at, created_at, updated_at
		FROM api_keys
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []*domain.APIKey
	for rows.Next() {
		apiKey := &domain.APIKey{}
		var lastUsedAt, expiresAt sql.NullTime
		var lastUsedIP sql.NullString

		err := rows.Scan(
			&apiKey.ID, &apiKey.UserID, &apiKey.DomainID, &apiKey.Name, &apiKey.KeyHash, &apiKey.Scopes,
			&lastUsedAt, &lastUsedIP, &expiresAt, &apiKey.CreatedAt, &apiKey.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}

		if lastUsedAt.Valid {
			apiKey.LastUsedAt = &lastUsedAt.Time
		}
		if lastUsedIP.Valid {
			apiKey.LastUsedIP = lastUsedIP.String
		}
		if expiresAt.Valid {
			apiKey.ExpiresAt = &expiresAt.Time
		}

		apiKeys = append(apiKeys, apiKey)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating API keys: %w", err)
	}

	return apiKeys, nil
}

// UpdateLastUsed updates the last used timestamp and IP for an API key
func (r *apiKeyRepository) UpdateLastUsed(id int64, ip string) error {
	query := `
		UPDATE api_keys
		SET last_used_at = ?, last_used_ip = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, time.Now(), ip, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update API key last used: %w", err)
	}

	return nil
}

// Delete removes an API key
func (r *apiKeyRepository) Delete(id int64) error {
	query := `DELETE FROM api_keys WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}
