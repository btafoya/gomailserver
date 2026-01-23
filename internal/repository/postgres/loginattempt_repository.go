package postgres

import (
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/repository"
)

// LoginAttemptRepository implements repository.LoginAttemptRepository
type LoginAttemptRepository struct {
	db *database.DB
}

// NewLoginAttemptRepository creates a new PostgreSQL login attempt repository
func NewLoginAttemptRepository(db *database.DB) repository.LoginAttemptRepository {
	return &LoginAttemptRepository{db: db}
}

// Record records a login attempt (success or failure)
func (r *LoginAttemptRepository) Record(ip, email string, success bool) error {
	// For failed logins, we record in failed_logins table
	// Successful logins are just counted for clearing failed attempts
	if success {
		// Don't record successful logins, just used to potentially clear failures
		return nil
	}

	query := `
		INSERT INTO failed_logins (ip_address, email, protocol, attempted_at)
		VALUES ($1, $2, 'SMTP/IMAP', $3)
	`
	_, err := r.db.Exec(query, ip, email, time.Now())
	return err
}

// GetRecentFailures counts failed login attempts from an IP within the duration
func (r *LoginAttemptRepository) GetRecentFailures(ip string, duration time.Duration) (int, error) {
	query := `
		SELECT COUNT(*) FROM failed_logins
		WHERE ip_address = $1
		AND attempted_at > $2
	`
	cutoff := time.Now().Add(-duration)
	var count int
	err := r.db.QueryRow(query, ip, cutoff).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetRecentUserFailures counts failed login attempts for an email within the duration
func (r *LoginAttemptRepository) GetRecentUserFailures(email string, duration time.Duration) (int, error) {
	query := `
		SELECT COUNT(*) FROM failed_logins
		WHERE email = $1
		AND attempted_at > $2
	`
	cutoff := time.Now().Add(-duration)
	var count int
	err := r.db.QueryRow(query, email, cutoff).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Cleanup removes old login attempt records
func (r *LoginAttemptRepository) Cleanup(age time.Duration) error {
	query := `
		DELETE FROM failed_logins
		WHERE attempted_at < $1
	`
	cutoff := time.Now().Add(-age)
	_, err := r.db.Exec(query, cutoff)
	return err
}
