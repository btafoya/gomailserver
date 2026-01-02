package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/openpgp"
)

// PGPService handles PGP/GPG key management for email encryption
type PGPService struct {
	db     *database.DB
	logger *zap.Logger
}

// NewPGPService creates a new PGP service
func NewPGPService(db *database.DB, logger *zap.Logger) *PGPService {
	return &PGPService{
		db:     db,
		logger: logger,
	}
}

// ImportKey imports a PGP public key for a user
func (s *PGPService) ImportKey(ctx context.Context, userID int64, publicKeyArmored string) (*domain.PGPKey, error) {
	// Parse the public key to extract metadata
	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicKeyArmored))
	if err != nil {
		return nil, fmt.Errorf("invalid PGP public key: %w", err)
	}

	if len(keyring) == 0 {
		return nil, fmt.Errorf("no keys found in provided data")
	}

	entity := keyring[0]
	fingerprint := fmt.Sprintf("%X", entity.PrimaryKey.Fingerprint)
	keyID := fmt.Sprintf("%X", entity.PrimaryKey.KeyId)

	// Extract expiration time from self-signature if available
	var expiresAt *time.Time
	for _, identity := range entity.Identities {
		if identity.SelfSignature != nil && identity.SelfSignature.KeyLifetimeSecs != nil && *identity.SelfSignature.KeyLifetimeSecs > 0 {
			expiry := entity.PrimaryKey.CreationTime.Add(time.Duration(*identity.SelfSignature.KeyLifetimeSecs) * time.Second)
			expiresAt = &expiry
			break
		}
	}

	// Check if key already exists
	existing, err := s.GetKeyByFingerprint(ctx, userID, fingerprint)
	if err == nil && existing != nil {
		// Update existing key
		return s.updateKey(ctx, existing.ID, publicKeyArmored, expiresAt)
	}

	// Insert new key
	query := `
		INSERT INTO pgp_keys (
			user_id, key_id, fingerprint, public_key, created_at, updated_at, expires_at, is_primary
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	// If this is the first key for the user, make it primary
	isPrimary, err := s.hasNoKeys(ctx, userID)
	if err != nil {
		isPrimary = false
	}

	now := time.Now()
	result, err := s.db.ExecContext(ctx, query,
		userID,
		keyID,
		fingerprint,
		publicKeyArmored,
		now,
		now,
		expiresAt,
		isPrimary,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert PGP key: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted key ID: %w", err)
	}

	s.logger.Info("PGP key imported",
		zap.Int64("user_id", userID),
		zap.String("fingerprint", fingerprint),
	)

	return &domain.PGPKey{
		ID:          id,
		UserID:      userID,
		KeyID:       keyID,
		Fingerprint: fingerprint,
		PublicKey:   publicKeyArmored,
		CreatedAt:   now,
		UpdatedAt:   now,
		ExpiresAt:   expiresAt,
		IsPrimary:   isPrimary,
	}, nil
}

// updateKey updates an existing PGP key
func (s *PGPService) updateKey(ctx context.Context, keyID int64, publicKey string, expiresAt *time.Time) (*domain.PGPKey, error) {
	query := `
		UPDATE pgp_keys
		SET public_key = ?, updated_at = ?, expires_at = ?
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		publicKey,
		time.Now(),
		expiresAt,
		keyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update PGP key: %w", err)
	}

	return s.GetKey(ctx, keyID)
}

// hasNoKeys checks if a user has any PGP keys
func (s *PGPService) hasNoKeys(ctx context.Context, userID int64) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM pgp_keys WHERE user_id = ?",
		userID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// GetKey retrieves a PGP key by ID
func (s *PGPService) GetKey(ctx context.Context, keyID int64) (*domain.PGPKey, error) {
	key := &domain.PGPKey{}
	var expiresAt sql.NullTime

	query := `
		SELECT id, user_id, key_id, fingerprint, public_key,
		       created_at, updated_at, expires_at, is_primary
		FROM pgp_keys
		WHERE id = ?
	`

	err := s.db.QueryRowContext(ctx, query, keyID).Scan(
		&key.ID,
		&key.UserID,
		&key.KeyID,
		&key.Fingerprint,
		&key.PublicKey,
		&key.CreatedAt,
		&key.UpdatedAt,
		&expiresAt,
		&key.IsPrimary,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("PGP key not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get PGP key: %w", err)
	}

	if expiresAt.Valid {
		key.ExpiresAt = &expiresAt.Time
	}

	return key, nil
}

// GetKeyByFingerprint retrieves a PGP key by fingerprint
func (s *PGPService) GetKeyByFingerprint(ctx context.Context, userID int64, fingerprint string) (*domain.PGPKey, error) {
	key := &domain.PGPKey{}
	var expiresAt sql.NullTime

	query := `
		SELECT id, user_id, key_id, fingerprint, public_key,
		       created_at, updated_at, expires_at, is_primary
		FROM pgp_keys
		WHERE user_id = ? AND fingerprint = ?
	`

	err := s.db.QueryRowContext(ctx, query, userID, fingerprint).Scan(
		&key.ID,
		&key.UserID,
		&key.KeyID,
		&key.Fingerprint,
		&key.PublicKey,
		&key.CreatedAt,
		&key.UpdatedAt,
		&expiresAt,
		&key.IsPrimary,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get PGP key: %w", err)
	}

	if expiresAt.Valid {
		key.ExpiresAt = &expiresAt.Time
	}

	return key, nil
}

// GetUserKeys retrieves all PGP keys for a user
func (s *PGPService) GetUserKeys(ctx context.Context, userID int64) ([]*domain.PGPKey, error) {
	query := `
		SELECT id, user_id, key_id, fingerprint, public_key,
		       created_at, updated_at, expires_at, is_primary
		FROM pgp_keys
		WHERE user_id = ?
		ORDER BY is_primary DESC, created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query PGP keys: %w", err)
	}
	defer rows.Close()

	var keys []*domain.PGPKey
	for rows.Next() {
		key := &domain.PGPKey{}
		var expiresAt sql.NullTime

		err := rows.Scan(
			&key.ID,
			&key.UserID,
			&key.KeyID,
			&key.Fingerprint,
			&key.PublicKey,
			&key.CreatedAt,
			&key.UpdatedAt,
			&expiresAt,
			&key.IsPrimary,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PGP key: %w", err)
		}

		if expiresAt.Valid {
			key.ExpiresAt = &expiresAt.Time
		}

		keys = append(keys, key)
	}

	return keys, rows.Err()
}

// GetPrimaryKey retrieves the primary PGP key for a user
func (s *PGPService) GetPrimaryKey(ctx context.Context, userID int64) (*domain.PGPKey, error) {
	key := &domain.PGPKey{}
	var expiresAt sql.NullTime

	query := `
		SELECT id, user_id, key_id, fingerprint, public_key,
		       created_at, updated_at, expires_at, is_primary
		FROM pgp_keys
		WHERE user_id = ? AND is_primary = 1
	`

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&key.ID,
		&key.UserID,
		&key.KeyID,
		&key.Fingerprint,
		&key.PublicKey,
		&key.CreatedAt,
		&key.UpdatedAt,
		&expiresAt,
		&key.IsPrimary,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get primary PGP key: %w", err)
	}

	if expiresAt.Valid {
		key.ExpiresAt = &expiresAt.Time
	}

	return key, nil
}

// SetPrimaryKey sets a key as the primary key for a user
func (s *PGPService) SetPrimaryKey(ctx context.Context, keyID int64) error {
	// Get the key to find the user_id
	key, err := s.GetKey(ctx, keyID)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Unset all primary keys for this user
	_, err = tx.ExecContext(ctx,
		"UPDATE pgp_keys SET is_primary = 0 WHERE user_id = ?",
		key.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to unset primary keys: %w", err)
	}

	// Set the new primary key
	_, err = tx.ExecContext(ctx,
		"UPDATE pgp_keys SET is_primary = 1 WHERE id = ?",
		keyID,
	)
	if err != nil {
		return fmt.Errorf("failed to set primary key: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("PGP primary key updated",
		zap.Int64("user_id", key.UserID),
		zap.Int64("key_id", keyID),
	)

	return nil
}

// DeleteKey deletes a PGP key
func (s *PGPService) DeleteKey(ctx context.Context, keyID int64) error {
	result, err := s.db.ExecContext(ctx,
		"DELETE FROM pgp_keys WHERE id = ?",
		keyID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete PGP key: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("PGP key not found")
	}

	s.logger.Info("PGP key deleted", zap.Int64("key_id", keyID))

	return nil
}
