package sqlite

import (
	"database/sql"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
)

// GreylistRepository implements repository.GreylistRepository
type GreylistRepository struct {
	db *database.DB
}

// NewGreylistRepository creates a new greylist repository
func NewGreylistRepository(db *database.DB) *GreylistRepository {
	return &GreylistRepository{db: db}
}

// Create creates a new greylist triplet
func (r *GreylistRepository) Create(ip, sender, recipient string) (*domain.GreylistTriplet, error) {
	query := `
		INSERT INTO greylist (sender_ip, sender_email, recipient_email, first_seen, status)
		VALUES (?, ?, ?, ?, 'greylisted')
	`
	now := time.Now()
	result, err := r.db.Exec(query, ip, sender, recipient, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &domain.GreylistTriplet{
		ID:        id,
		IP:        ip,
		Sender:    sender,
		Recipient: recipient,
		FirstSeen: now,
		PassCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Get retrieves a greylist triplet
func (r *GreylistRepository) Get(ip, sender, recipient string) (*domain.GreylistTriplet, error) {
	query := `
		SELECT id, first_seen, passed_at
		FROM greylist
		WHERE sender_ip = ?
		AND sender_email = ?
		AND recipient_email = ?
		AND status != 'expired'
	`

	var triplet domain.GreylistTriplet
	var passedAt sql.NullTime

	err := r.db.QueryRow(query, ip, sender, recipient).Scan(
		&triplet.ID,
		&triplet.FirstSeen,
		&passedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	triplet.IP = ip
	triplet.Sender = sender
	triplet.Recipient = recipient

	// Count passes based on whether passed_at is set
	if passedAt.Valid {
		triplet.PassCount = 1
	} else {
		triplet.PassCount = 0
	}

	return &triplet, nil
}

// IncrementPass marks a greylist triplet as passed
func (r *GreylistRepository) IncrementPass(id int64) error {
	query := `
		UPDATE greylist
		SET passed_at = ?,
		    status = 'passed'
		WHERE id = ?
	`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// DeleteOlderThan deletes greylist entries older than the specified age
func (r *GreylistRepository) DeleteOlderThan(age time.Duration) error {
	query := `
		UPDATE greylist
		SET status = 'expired'
		WHERE first_seen < ?
		AND status != 'expired'
	`
	cutoff := time.Now().Add(-age)
	_, err := r.db.Exec(query, cutoff)
	if err != nil {
		return err
	}

	// Actually delete expired entries after marking them
	deleteQuery := `
		DELETE FROM greylist
		WHERE status = 'expired'
		AND first_seen < ?
	`
	oldCutoff := time.Now().Add(-age * 2) // Keep expired entries for double the age
	_, err = r.db.Exec(deleteQuery, oldCutoff)
	return err
}
