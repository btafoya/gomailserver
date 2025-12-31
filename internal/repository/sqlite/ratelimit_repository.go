package sqlite

import (
	"database/sql"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
)

// RateLimitRepository implements repository.RateLimitRepository
type RateLimitRepository struct {
	db *database.DB
}

// NewRateLimitRepository creates a new rate limit repository
func NewRateLimitRepository(db *database.DB) *RateLimitRepository {
	return &RateLimitRepository{db: db}
}

// Get retrieves a rate limit entry
func (r *RateLimitRepository) Get(key string, limitType string) (*domain.RateLimitEntry, error) {
	// Determine entity type from the key (IP vs email/user)
	entityType := determineEntityType(key)

	query := `
		SELECT id, count, window_start
		FROM rate_limits
		WHERE entity_type = ?
		AND entity_value = ?
		AND action_type = ?
		AND window_start > ?
		ORDER BY window_start DESC
		LIMIT 1
	`

	// Only count entries from the last hour (adjust based on need)
	cutoff := time.Now().Add(-1 * time.Hour)

	var entry domain.RateLimitEntry
	var windowStart time.Time
	err := r.db.QueryRow(query, entityType, key, limitType, cutoff).Scan(
		&entry.ID,
		&entry.Count,
		&windowStart,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entry.Key = key
	entry.Type = limitType
	entry.WindowStart = windowStart

	return &entry, nil
}

// CreateOrUpdate creates or updates a rate limit entry
func (r *RateLimitRepository) CreateOrUpdate(entry *domain.RateLimitEntry) error {
	entityType := determineEntityType(entry.Key)

	query := `
		INSERT INTO rate_limits (entity_type, entity_value, action_type, count, window_start)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(entity_type, entity_value, action_type, window_start)
		DO UPDATE SET count = excluded.count
	`

	_, err := r.db.Exec(query,
		entityType,
		entry.Key,
		entry.Type,
		entry.Count,
		entry.WindowStart,
	)
	return err
}

// Cleanup removes old rate limit entries
func (r *RateLimitRepository) Cleanup(windowDuration time.Duration) error {
	query := `
		DELETE FROM rate_limits
		WHERE window_start < ?
	`
	cutoff := time.Now().Add(-windowDuration)
	_, err := r.db.Exec(query, cutoff)
	return err
}

// determineEntityType determines if the key is an IP or user identifier
func determineEntityType(key string) string {
	// Simple heuristic: if it contains '@' it's a user, otherwise IP
	for _, c := range key {
		if c == '@' {
			return "user"
		}
	}
	return "ip"
}
