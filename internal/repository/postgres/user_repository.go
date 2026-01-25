package postgres

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

func NewUserRepository(db *database.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (
			email, domain_id, password_hash, full_name, display_name, role,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id
	`

	err := r.db.QueryRow(query,
		user.Email, user.DomainID, user.PasswordHash, user.FullName, user.DisplayName, user.Role,
		user.Quota, user.UsedQuota, user.Status, user.AuthMethod, user.TOTPSecret, user.TOTPEnabled,
		user.ForwardTo, user.AutoReplyEnabled, user.AutoReplySubject, user.AutoReplyBody,
		user.SpamThreshold, user.Language, time.Now(), time.Now(),
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return nil
}

func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	query := `
		SELECT
			id, email, domain_id, password_hash, full_name, display_name, role,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, last_login, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.DomainID, &user.PasswordHash, &user.FullName, &user.DisplayName, &user.Role,
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

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT
			id, email, domain_id, password_hash, full_name, display_name, role,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, last_login, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.DomainID, &user.PasswordHash, &user.FullName, &user.DisplayName, &user.Role,
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

func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users SET
			email = $1, domain_id = $2, password_hash = $3, full_name = $4, display_name = $5, role = $6,
			quota = $7, used_quota = $8, status = $9, auth_method = $10, totp_secret = $11, totp_enabled = $12,
			forward_to = $13, auto_reply_enabled = $14, auto_reply_subject = $15, auto_reply_body = $16,
			spam_threshold = $17, language = $18, updated_at = $19
		WHERE id = $20
	`

	_, err := r.db.Exec(query,
		user.Email, user.DomainID, user.PasswordHash, user.FullName, user.DisplayName, user.Role,
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

func (r *userRepository) UpdateLastLogin(id int64) error {
	query := `UPDATE users SET last_login = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

func (r *userRepository) UpdatePassword(userID int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, passwordHash, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepository) List(domainID int64, offset, limit int) ([]*domain.User, error) {
	query := `
		SELECT
			id, email, domain_id, password_hash, full_name, display_name, role,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, last_login, created_at, updated_at
		FROM users
		WHERE domain_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
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
			&user.ID, &user.Email, &user.DomainID, &user.PasswordHash, &user.FullName, &user.DisplayName, &user.Role,
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

func (r *userRepository) ListAll() ([]*domain.User, error) {
	query := `
		SELECT
			id, email, domain_id, password_hash, full_name, display_name, role,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, last_login, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all users: %w", err)
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		var lastLogin sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Email, &user.DomainID, &user.PasswordHash, &user.FullName, &user.DisplayName, &user.Role,
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

func (r *userRepository) ListPaginated(offset, limit int) ([]*domain.User, error) {
	query := `
		SELECT
			id, email, domain_id, password_hash, full_name, display_name, role,
			quota, used_quota, status, auth_method, totp_secret, totp_enabled,
			forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
			spam_threshold, language, last_login, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users paginated: %w", err)
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		var lastLogin sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Email, &user.DomainID, &user.PasswordHash, &user.FullName, &user.DisplayName, &user.Role,
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

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

func (r *userRepository) CountByDomain(domainID int64) (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE domain_id = $1", domainID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by domain: %w", err)
	}
	return count, nil
}
