package postgres

import (
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/repository"
)

// IPBlacklistRepository implements repository.IPBlacklistRepository
type IPBlacklistRepository struct {
	db *database.DB
}

// NewIPBlacklistRepository creates a new PostgreSQL IP blacklist repository
func NewIPBlacklistRepository(db *database.DB) repository.IPBlacklistRepository {
	return &IPBlacklistRepository{db: db}
}

// Add adds an IP to the blacklist
func (r *IPBlacklistRepository) Add(ip, reason string, expiresAt *time.Time) error {
	query := `
		INSERT INTO ip_blacklist (ip_address, reason, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT(ip_address) DO UPDATE SET
			reason = excluded.reason,
			expires_at = excluded.expires_at
	`
	_, err := r.db.Exec(query, ip, reason, expiresAt)
	return err
}

// IsBlacklisted checks if an IP is currently blacklisted
func (r *IPBlacklistRepository) IsBlacklisted(ip string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM ip_blacklist
		WHERE ip_address = $1
		AND (expires_at IS NULL OR expires_at > $2)
	`
	var count int
	err := r.db.QueryRow(query, ip, time.Now()).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Remove removes an IP from the blacklist
func (r *IPBlacklistRepository) Remove(ip string) error {
	query := `DELETE FROM ip_blacklist WHERE ip_address = $1`
	_, err := r.db.Exec(query, ip)
	return err
}

// RemoveExpired removes all expired blacklist entries
func (r *IPBlacklistRepository) RemoveExpired() error {
	query := `
		DELETE FROM ip_blacklist
		WHERE expires_at IS NOT NULL
		AND expires_at <= $1
	`
	_, err := r.db.Exec(query, time.Now())
	return err
}
