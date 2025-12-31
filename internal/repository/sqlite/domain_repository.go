package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type domainRepository struct {
	db *database.DB
}

// NewDomainRepository creates a new SQLite domain repository
func NewDomainRepository(db *database.DB) repository.DomainRepository {
	return &domainRepository{db: db}
}

// Create inserts a new domain
func (r *domainRepository) Create(dom *domain.Domain) error {
	query := `
		INSERT INTO domains (
			name, status, max_users, max_mailbox_size, default_quota,
			catchall_email, backup_mx,
			dkim_selector, dkim_private_key, dkim_public_key,
			dkim_signing_enabled, dkim_verify_enabled, dkim_key_size, dkim_key_type, dkim_headers_to_sign,
			spf_record, spf_enabled, spf_dns_server, spf_dns_timeout, spf_max_lookups, spf_fail_action, spf_softfail_action,
			dmarc_policy, dmarc_enabled, dmarc_dns_server, dmarc_dns_timeout, dmarc_report_enabled, dmarc_report_email,
			clamav_enabled, clamav_max_scan_size, clamav_virus_action, clamav_fail_action,
			spam_enabled, spam_reject_score, spam_quarantine_score, spam_learning_enabled,
			greylist_enabled, greylist_delay_minutes, greylist_expiry_days, greylist_cleanup_interval, greylist_whitelist_after,
			ratelimit_enabled, ratelimit_smtp_per_ip, ratelimit_smtp_per_user, ratelimit_smtp_per_domain, ratelimit_auth_per_ip, ratelimit_imap_per_user, ratelimit_cleanup_interval,
			auth_totp_enforced, auth_brute_force_enabled, auth_brute_force_threshold, auth_brute_force_window_minutes, auth_brute_force_block_minutes, auth_ip_blacklist_enabled, auth_cleanup_interval,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		dom.Name, dom.Status, dom.MaxUsers, dom.MaxMailboxSize, dom.DefaultQuota,
		dom.CatchallEmail, dom.BackupMX,
		dom.DKIMSelector, dom.DKIMPrivateKey, dom.DKIMPublicKey,
		dom.DKIMSigningEnabled, dom.DKIMVerifyEnabled, dom.DKIMKeySize, dom.DKIMKeyType, dom.DKIMHeadersToSign,
		dom.SPFRecord, dom.SPFEnabled, dom.SPFDNSServer, dom.SPFDNSTimeout, dom.SPFMaxLookups, dom.SPFFailAction, dom.SPFSoftFailAction,
		dom.DMARCPolicy, dom.DMARCEnabled, dom.DMARCDNSServer, dom.DMARCDNSTimeout, dom.DMARCReportEnabled, dom.DMARCReportEmail,
		dom.ClamAVEnabled, dom.ClamAVMaxScanSize, dom.ClamAVVirusAction, dom.ClamAVFailAction,
		dom.SpamEnabled, dom.SpamRejectScore, dom.SpamQuarantineScore, dom.SpamLearningEnabled,
		dom.GreylistEnabled, dom.GreylistDelayMinutes, dom.GreylistExpiryDays, dom.GreylistCleanupInterval, dom.GreylistWhitelistAfter,
		dom.RateLimitEnabled, dom.RateLimitSMTPPerIP, dom.RateLimitSMTPPerUser, dom.RateLimitSMTPPerDomain, dom.RateLimitAuthPerIP, dom.RateLimitIMAPPerUser, dom.RateLimitCleanupInterval,
		dom.AuthTOTPEnforced, dom.AuthBruteForceEnabled, dom.AuthBruteForceThreshold, dom.AuthBruteForceWindowMinutes, dom.AuthBruteForceBlockMinutes, dom.AuthIPBlacklistEnabled, dom.AuthCleanupInterval,
		time.Now(), time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create domain: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get domain ID: %w", err)
	}

	dom.ID = id
	dom.CreatedAt = time.Now()
	dom.UpdatedAt = time.Now()

	return nil
}

// GetByID retrieves a domain by ID
func (r *domainRepository) GetByID(id int64) (*domain.Domain, error) {
	query := `
		SELECT
			id, name, status, max_users, max_mailbox_size, default_quota,
			catchall_email, backup_mx,
			dkim_selector, dkim_private_key, dkim_public_key,
			dkim_signing_enabled, dkim_verify_enabled, dkim_key_size, dkim_key_type, dkim_headers_to_sign,
			spf_record, spf_enabled, spf_dns_server, spf_dns_timeout, spf_max_lookups, spf_fail_action, spf_softfail_action,
			dmarc_policy, dmarc_enabled, dmarc_dns_server, dmarc_dns_timeout, dmarc_report_enabled, dmarc_report_email,
			clamav_enabled, clamav_max_scan_size, clamav_virus_action, clamav_fail_action,
			spam_enabled, spam_reject_score, spam_quarantine_score, spam_learning_enabled,
			greylist_enabled, greylist_delay_minutes, greylist_expiry_days, greylist_cleanup_interval, greylist_whitelist_after,
			ratelimit_enabled, ratelimit_smtp_per_ip, ratelimit_smtp_per_user, ratelimit_smtp_per_domain, ratelimit_auth_per_ip, ratelimit_imap_per_user, ratelimit_cleanup_interval,
			auth_totp_enforced, auth_brute_force_enabled, auth_brute_force_threshold, auth_brute_force_window_minutes, auth_brute_force_block_minutes, auth_ip_blacklist_enabled, auth_cleanup_interval,
			created_at, updated_at
		FROM domains
		WHERE id = ?
	`

	dom := &domain.Domain{}

	err := r.db.QueryRow(query, id).Scan(
		&dom.ID, &dom.Name, &dom.Status, &dom.MaxUsers, &dom.MaxMailboxSize, &dom.DefaultQuota,
		&dom.CatchallEmail, &dom.BackupMX,
		&dom.DKIMSelector, &dom.DKIMPrivateKey, &dom.DKIMPublicKey,
		&dom.DKIMSigningEnabled, &dom.DKIMVerifyEnabled, &dom.DKIMKeySize, &dom.DKIMKeyType, &dom.DKIMHeadersToSign,
		&dom.SPFRecord, &dom.SPFEnabled, &dom.SPFDNSServer, &dom.SPFDNSTimeout, &dom.SPFMaxLookups, &dom.SPFFailAction, &dom.SPFSoftFailAction,
		&dom.DMARCPolicy, &dom.DMARCEnabled, &dom.DMARCDNSServer, &dom.DMARCDNSTimeout, &dom.DMARCReportEnabled, &dom.DMARCReportEmail,
		&dom.ClamAVEnabled, &dom.ClamAVMaxScanSize, &dom.ClamAVVirusAction, &dom.ClamAVFailAction,
		&dom.SpamEnabled, &dom.SpamRejectScore, &dom.SpamQuarantineScore, &dom.SpamLearningEnabled,
		&dom.GreylistEnabled, &dom.GreylistDelayMinutes, &dom.GreylistExpiryDays, &dom.GreylistCleanupInterval, &dom.GreylistWhitelistAfter,
		&dom.RateLimitEnabled, &dom.RateLimitSMTPPerIP, &dom.RateLimitSMTPPerUser, &dom.RateLimitSMTPPerDomain, &dom.RateLimitAuthPerIP, &dom.RateLimitIMAPPerUser, &dom.RateLimitCleanupInterval,
		&dom.AuthTOTPEnforced, &dom.AuthBruteForceEnabled, &dom.AuthBruteForceThreshold, &dom.AuthBruteForceWindowMinutes, &dom.AuthBruteForceBlockMinutes, &dom.AuthIPBlacklistEnabled, &dom.AuthCleanupInterval,
		&dom.CreatedAt, &dom.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("domain not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return dom, nil
}

// GetByName retrieves a domain by name
func (r *domainRepository) GetByName(name string) (*domain.Domain, error) {
	query := `
		SELECT
			id, name, status, max_users, max_mailbox_size, default_quota,
			catchall_email, backup_mx,
			dkim_selector, dkim_private_key, dkim_public_key,
			dkim_signing_enabled, dkim_verify_enabled, dkim_key_size, dkim_key_type, dkim_headers_to_sign,
			spf_record, spf_enabled, spf_dns_server, spf_dns_timeout, spf_max_lookups, spf_fail_action, spf_softfail_action,
			dmarc_policy, dmarc_enabled, dmarc_dns_server, dmarc_dns_timeout, dmarc_report_enabled, dmarc_report_email,
			clamav_enabled, clamav_max_scan_size, clamav_virus_action, clamav_fail_action,
			spam_enabled, spam_reject_score, spam_quarantine_score, spam_learning_enabled,
			greylist_enabled, greylist_delay_minutes, greylist_expiry_days, greylist_cleanup_interval, greylist_whitelist_after,
			ratelimit_enabled, ratelimit_smtp_per_ip, ratelimit_smtp_per_user, ratelimit_smtp_per_domain, ratelimit_auth_per_ip, ratelimit_imap_per_user, ratelimit_cleanup_interval,
			auth_totp_enforced, auth_brute_force_enabled, auth_brute_force_threshold, auth_brute_force_window_minutes, auth_brute_force_block_minutes, auth_ip_blacklist_enabled, auth_cleanup_interval,
			created_at, updated_at
		FROM domains
		WHERE name = ?
	`

	dom := &domain.Domain{}

	err := r.db.QueryRow(query, name).Scan(
		&dom.ID, &dom.Name, &dom.Status, &dom.MaxUsers, &dom.MaxMailboxSize, &dom.DefaultQuota,
		&dom.CatchallEmail, &dom.BackupMX,
		&dom.DKIMSelector, &dom.DKIMPrivateKey, &dom.DKIMPublicKey,
		&dom.DKIMSigningEnabled, &dom.DKIMVerifyEnabled, &dom.DKIMKeySize, &dom.DKIMKeyType, &dom.DKIMHeadersToSign,
		&dom.SPFRecord, &dom.SPFEnabled, &dom.SPFDNSServer, &dom.SPFDNSTimeout, &dom.SPFMaxLookups, &dom.SPFFailAction, &dom.SPFSoftFailAction,
		&dom.DMARCPolicy, &dom.DMARCEnabled, &dom.DMARCDNSServer, &dom.DMARCDNSTimeout, &dom.DMARCReportEnabled, &dom.DMARCReportEmail,
		&dom.ClamAVEnabled, &dom.ClamAVMaxScanSize, &dom.ClamAVVirusAction, &dom.ClamAVFailAction,
		&dom.SpamEnabled, &dom.SpamRejectScore, &dom.SpamQuarantineScore, &dom.SpamLearningEnabled,
		&dom.GreylistEnabled, &dom.GreylistDelayMinutes, &dom.GreylistExpiryDays, &dom.GreylistCleanupInterval, &dom.GreylistWhitelistAfter,
		&dom.RateLimitEnabled, &dom.RateLimitSMTPPerIP, &dom.RateLimitSMTPPerUser, &dom.RateLimitSMTPPerDomain, &dom.RateLimitAuthPerIP, &dom.RateLimitIMAPPerUser, &dom.RateLimitCleanupInterval,
		&dom.AuthTOTPEnforced, &dom.AuthBruteForceEnabled, &dom.AuthBruteForceThreshold, &dom.AuthBruteForceWindowMinutes, &dom.AuthBruteForceBlockMinutes, &dom.AuthIPBlacklistEnabled, &dom.AuthCleanupInterval,
		&dom.CreatedAt, &dom.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("domain not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return dom, nil
}

// Update updates a domain
func (r *domainRepository) Update(dom *domain.Domain) error {
	query := `
		UPDATE domains SET
			name = ?, status = ?, max_users = ?, max_mailbox_size = ?, default_quota = ?,
			catchall_email = ?, backup_mx = ?,
			dkim_selector = ?, dkim_private_key = ?, dkim_public_key = ?,
			dkim_signing_enabled = ?, dkim_verify_enabled = ?, dkim_key_size = ?, dkim_key_type = ?, dkim_headers_to_sign = ?,
			spf_record = ?, spf_enabled = ?, spf_dns_server = ?, spf_dns_timeout = ?, spf_max_lookups = ?, spf_fail_action = ?, spf_softfail_action = ?,
			dmarc_policy = ?, dmarc_enabled = ?, dmarc_dns_server = ?, dmarc_dns_timeout = ?, dmarc_report_enabled = ?, dmarc_report_email = ?,
			clamav_enabled = ?, clamav_max_scan_size = ?, clamav_virus_action = ?, clamav_fail_action = ?,
			spam_enabled = ?, spam_reject_score = ?, spam_quarantine_score = ?, spam_learning_enabled = ?,
			greylist_enabled = ?, greylist_delay_minutes = ?, greylist_expiry_days = ?, greylist_cleanup_interval = ?, greylist_whitelist_after = ?,
			ratelimit_enabled = ?, ratelimit_smtp_per_ip = ?, ratelimit_smtp_per_user = ?, ratelimit_smtp_per_domain = ?, ratelimit_auth_per_ip = ?, ratelimit_imap_per_user = ?, ratelimit_cleanup_interval = ?,
			auth_totp_enforced = ?, auth_brute_force_enabled = ?, auth_brute_force_threshold = ?, auth_brute_force_window_minutes = ?, auth_brute_force_block_minutes = ?, auth_ip_blacklist_enabled = ?, auth_cleanup_interval = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		dom.Name, dom.Status, dom.MaxUsers, dom.MaxMailboxSize, dom.DefaultQuota,
		dom.CatchallEmail, dom.BackupMX,
		dom.DKIMSelector, dom.DKIMPrivateKey, dom.DKIMPublicKey,
		dom.DKIMSigningEnabled, dom.DKIMVerifyEnabled, dom.DKIMKeySize, dom.DKIMKeyType, dom.DKIMHeadersToSign,
		dom.SPFRecord, dom.SPFEnabled, dom.SPFDNSServer, dom.SPFDNSTimeout, dom.SPFMaxLookups, dom.SPFFailAction, dom.SPFSoftFailAction,
		dom.DMARCPolicy, dom.DMARCEnabled, dom.DMARCDNSServer, dom.DMARCDNSTimeout, dom.DMARCReportEnabled, dom.DMARCReportEmail,
		dom.ClamAVEnabled, dom.ClamAVMaxScanSize, dom.ClamAVVirusAction, dom.ClamAVFailAction,
		dom.SpamEnabled, dom.SpamRejectScore, dom.SpamQuarantineScore, dom.SpamLearningEnabled,
		dom.GreylistEnabled, dom.GreylistDelayMinutes, dom.GreylistExpiryDays, dom.GreylistCleanupInterval, dom.GreylistWhitelistAfter,
		dom.RateLimitEnabled, dom.RateLimitSMTPPerIP, dom.RateLimitSMTPPerUser, dom.RateLimitSMTPPerDomain, dom.RateLimitAuthPerIP, dom.RateLimitIMAPPerUser, dom.RateLimitCleanupInterval,
		dom.AuthTOTPEnforced, dom.AuthBruteForceEnabled, dom.AuthBruteForceThreshold, dom.AuthBruteForceWindowMinutes, dom.AuthBruteForceBlockMinutes, dom.AuthIPBlacklistEnabled, dom.AuthCleanupInterval,
		time.Now(), dom.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update domain: %w", err)
	}

	dom.UpdatedAt = time.Now()
	return nil
}

// Delete deletes a domain
func (r *domainRepository) Delete(id int64) error {
	query := `DELETE FROM domains WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}
	return nil
}

// List lists domains with pagination
func (r *domainRepository) List(offset, limit int) ([]*domain.Domain, error) {
	query := `
		SELECT
			id, name, status, max_users, max_mailbox_size, default_quota,
			catchall_email, backup_mx,
			dkim_selector, dkim_private_key, dkim_public_key,
			dkim_signing_enabled, dkim_verify_enabled, dkim_key_size, dkim_key_type, dkim_headers_to_sign,
			spf_record, spf_enabled, spf_dns_server, spf_dns_timeout, spf_max_lookups, spf_fail_action, spf_softfail_action,
			dmarc_policy, dmarc_enabled, dmarc_dns_server, dmarc_dns_timeout, dmarc_report_enabled, dmarc_report_email,
			clamav_enabled, clamav_max_scan_size, clamav_virus_action, clamav_fail_action,
			spam_enabled, spam_reject_score, spam_quarantine_score, spam_learning_enabled,
			greylist_enabled, greylist_delay_minutes, greylist_expiry_days, greylist_cleanup_interval, greylist_whitelist_after,
			ratelimit_enabled, ratelimit_smtp_per_ip, ratelimit_smtp_per_user, ratelimit_smtp_per_domain, ratelimit_auth_per_ip, ratelimit_imap_per_user, ratelimit_cleanup_interval,
			auth_totp_enforced, auth_brute_force_enabled, auth_brute_force_threshold, auth_brute_force_window_minutes, auth_brute_force_block_minutes, auth_ip_blacklist_enabled, auth_cleanup_interval,
			created_at, updated_at
		FROM domains
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}
	defer rows.Close()

	domains := make([]*domain.Domain, 0)
	for rows.Next() {
		dom := &domain.Domain{}

		err := rows.Scan(
			&dom.ID, &dom.Name, &dom.Status, &dom.MaxUsers, &dom.MaxMailboxSize, &dom.DefaultQuota,
			&dom.CatchallEmail, &dom.BackupMX,
			&dom.DKIMSelector, &dom.DKIMPrivateKey, &dom.DKIMPublicKey,
			&dom.DKIMSigningEnabled, &dom.DKIMVerifyEnabled, &dom.DKIMKeySize, &dom.DKIMKeyType, &dom.DKIMHeadersToSign,
			&dom.SPFRecord, &dom.SPFEnabled, &dom.SPFDNSServer, &dom.SPFDNSTimeout, &dom.SPFMaxLookups, &dom.SPFFailAction, &dom.SPFSoftFailAction,
			&dom.DMARCPolicy, &dom.DMARCEnabled, &dom.DMARCDNSServer, &dom.DMARCDNSTimeout, &dom.DMARCReportEnabled, &dom.DMARCReportEmail,
			&dom.ClamAVEnabled, &dom.ClamAVMaxScanSize, &dom.ClamAVVirusAction, &dom.ClamAVFailAction,
			&dom.SpamEnabled, &dom.SpamRejectScore, &dom.SpamQuarantineScore, &dom.SpamLearningEnabled,
			&dom.GreylistEnabled, &dom.GreylistDelayMinutes, &dom.GreylistExpiryDays, &dom.GreylistCleanupInterval, &dom.GreylistWhitelistAfter,
			&dom.RateLimitEnabled, &dom.RateLimitSMTPPerIP, &dom.RateLimitSMTPPerUser, &dom.RateLimitSMTPPerDomain, &dom.RateLimitAuthPerIP, &dom.RateLimitIMAPPerUser, &dom.RateLimitCleanupInterval,
			&dom.AuthTOTPEnforced, &dom.AuthBruteForceEnabled, &dom.AuthBruteForceThreshold, &dom.AuthBruteForceWindowMinutes, &dom.AuthBruteForceBlockMinutes, &dom.AuthIPBlacklistEnabled, &dom.AuthCleanupInterval,
			&dom.CreatedAt, &dom.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan domain: %w", err)
		}

		domains = append(domains, dom)
	}

	return domains, rows.Err()
}
