package domain

import "time"

// PGPKey represents a user's PGP/GPG public key for email encryption
type PGPKey struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	KeyID       string    `json:"key_id"`
	Fingerprint string    `json:"fingerprint"`
	PublicKey   string    `json:"public_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsPrimary   bool      `json:"is_primary"`
}
