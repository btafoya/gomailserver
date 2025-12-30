package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type userRepository struct {
	db *database.DB
}

// NewUserRepository creates a new SQLite user repository
func NewUserRepository(db *database.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user
func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (
			email, domain_id, password_hash, full_name, display_name,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		user.Email, user.DomainID, user.PasswordHash, user.FullName, user.DisplayName,
		user.Quota, user.UsedQuota, user.Status, user.AuthMethod, user.TOTPSecret, user.TOTPEnabled,
		user.ForwardTo, user.AutoReplyEnabled, user.AutoReplySubject, user.AutoReplyBody,
		user.SpamThreshold, user.Language, time.Now(), time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	user.ID = id
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	query := `
		SELECT
			id, email, domain_id, password_hash, full_name, display_name,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, last_login, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &domain.User{}
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.DomainID, &user.PasswordHash, &user.FullName, &user.DisplayName,
		&user.Quota, &user.UsedQuota, &user.Status, &user.AuthMethod, &user.TOTPSecret, &user.TOTPEnabled,
		&user.ForwardTo, &user.AutoReplyEnabled, &user.AutoReplySubject, &user.AutoReplyBody,
		&user.SpamThreshold, &user.Language, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT
			id, email, domain_id, password_hash, full_name, display_name,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, last_login, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &domain.User{}
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.DomainID, &user.PasswordHash, &user.FullName, &user.DisplayName,
		&user.Quota, &user.UsedQuota, &user.Status, &user.AuthMethod, &user.TOTPSecret, &user.TOTPEnabled,
		&user.ForwardTo, &user.AutoReplyEnabled, &user.AutoReplySubject, &user.AutoReplyBody,
		&user.SpamThreshold, &user.Language, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

// Update updates a user
func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users SET
			email = ?, domain_id = ?, password_hash = ?, full_name = ?, display_name = ?,
			quota = ?, used_quota = ?, status = ?, auth_method = ?, totp_secret = ?, totp_enabled = ?,
			forward_to = ?, auto_reply_enabled = ?, auto_reply_subject = ?, auto_reply_body = ?,
			spam_threshold = ?, language = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		user.Email, user.DomainID, user.PasswordHash, user.FullName, user.DisplayName,
		user.Quota, user.UsedQuota, user.Status, user.AuthMethod, user.TOTPSecret, user.TOTPEnabled,
		user.ForwardTo, user.AutoReplyEnabled, user.AutoReplySubject, user.AutoReplyBody,
		user.SpamThreshold, user.Language, time.Now(), user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.UpdatedAt = time.Now()
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (r *userRepository) UpdateLastLogin(id int64) error {
	query := `UPDATE users SET last_login = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// UpdatePassword updates a user's password hash
func (r *userRepository) UpdatePassword(userID int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, passwordHash, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// Delete deletes a user
func (r *userRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List lists users for a domain
func (r *userRepository) List(domainID int64, offset, limit int) ([]*domain.User, error) {
	query := `
		SELECT
			id, email, domain_id, password_hash, full_name, display_name,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, last_login, created_at, updated_at
		FROM users
		WHERE domain_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, domainID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		var lastLogin sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Email, &user.DomainID, &user.PasswordHash, &user.FullName, &user.DisplayName,
			&user.Quota, &user.UsedQuota, &user.Status, &user.AuthMethod, &user.TOTPSecret, &user.TOTPEnabled,
			&user.ForwardTo, &user.AutoReplyEnabled, &user.AutoReplySubject, &user.AutoReplyBody,
			&user.SpamThreshold, &user.Language, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}

		users = append(users, user)
	}

	return users, rows.Err()
}
