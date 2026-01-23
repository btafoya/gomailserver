package postgres

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

func NewDomainRepository(db *database.DB) repository.DomainRepository {
	return &domainRepository{db: db}
}

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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40)
		RETURNING id
	`

	err := r.db.QueryRow(query,
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
	).Scan(&dom.ID)

	if err != nil {
		return fmt.Errorf("failed to create domain: %w", err)
	}

	dom.CreatedAt = time.Now()
	dom.UpdatedAt = time.Now()

	return nil
}

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
		WHERE id = $1
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
		WHERE name = $1
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

func (r *domainRepository) Update(dom *domain.Domain) error {
	query := `
		UPDATE domains SET
			name = $1, status = $2, max_users = $3, max_mailbox_size = $4, default_quota = $5,
			catchall_email = $6, backup_mx = $7,
			dkim_selector = $8, dkim_private_key = $9, dkim_public_key = $10,
			dkim_signing_enabled = $11, dkim_verify_enabled = $12, dkim_key_size = $13, dkim_key_type = $14, dkim_headers_to_sign = $15,
			spf_record = $16, spf_enabled = $17, spf_dns_server = $18, spf_dns_timeout = $19, spf_max_lookups = $20, spf_fail_action = $21, spf_softfail_action = $22,
			dmarc_policy = $23, dmarc_enabled = $24, dmarc_dns_server = $25, dmarc_dns_timeout = $26, dmarc_report_enabled = $27, dmarc_report_email = $28,
			clamav_enabled = $29, clamav_max_scan_size = $30, clamav_virus_action = $31, clamav_fail_action = $32,
			spam_enabled = $33, spam_reject_score = $34, spam_quarantine_score = $35, spam_learning_enabled = $36,
			greylist_enabled = $37, greylist_delay_minutes = $38, greylist_expiry_days = $39, greylist_cleanup_interval = $40, greylist_whitelist_after = $41,
			ratelimit_enabled = $42, ratelimit_smtp_per_ip = $43, ratelimit_smtp_per_user = $44, ratelimit_smtp_per_domain = $45, ratelimit_auth_per_ip = $46, ratelimit_imap_per_user = $47, ratelimit_cleanup_interval = $48,
			auth_totp_enforced = $49, auth_brute_force_enabled = $50, auth_brute_force_threshold = $51, auth_brute_force_window_minutes = $52, auth_brute_force_block_minutes = $53, auth_ip_blacklist_enabled = $54, auth_cleanup_interval = $55,
			updated_at = $56
		WHERE id = $57
	`

// GetDKIMConfig retrieves DKIM configuration for a domain
func (r *domainRepository) GetDKIMConfig(domainName string) (*domain.DKIMConfig, error) {
	dom, err := r.GetByName(domainName)
	if err != nil {
		return nil, fmt.Errorf("domain not found: %w", err)
	}

	return &domain.DKIMConfig{
		Domain:     dom.Name,
		Selector:   dom.DKIMSelector,
		PrivateKey: []byte(dom.DKIMPrivateKey),
		PublicKey:  dom.DKIMPublicKey,
	}, nil
}

// Update updates an existing domain
func (r *domainRepository) Update(dom *domain.Domain) error {
	query := `
		UPDATE domains SET
			name = $1, status = $2, max_users = $3, max_mailbox_size = $4, default_quota = $5,
			catchall_email = $6, backup_mx = $7,
			dkim_selector = $8, dkim_private_key = $9, dkim_public_key = $10,
			dkim_signing_enabled = $11, dkim_verify_enabled = $12, dkim_key_size = $13, dkim_key_type = $14, dkim_headers_to_sign = $15,
			spf_record = $16, spf_enabled = $17, spf_dns_server = $18, spf_dns_timeout = $19, spf_max_lookups = $20, spf_fail_action = $21, spf_softfail_action = $22,
			dmarc_policy = $23, dmarc_enabled = $24, dmarc_dns_server = $25, dmarc_dns_timeout = $26, dmarc_report_enabled = $27, dmarc_report_email = $28,
			clamav_enabled = $29, clamav_max_scan_size = $30, clamav_virus_action = $31, clamav_fail_action = $32,
			spam_enabled = $33, spam_reject_score = $34, spam_quarantine_score = $35, spam_learning_enabled = $36,
			greylist_enabled = $37, greylist_delay_minutes = $38, greylist_expiry_days = $39, greylist_cleanup_interval = $40, greylist_whitelist_after = $41,
			rate_limit_enabled = $42, rate_limit_smtp_per_ip = $43, rate_limit_smtp_per_user = $44, rate_limit_smtp_per_domain = $45, rate_limit_auth_per_ip = $46, rate_limit_imap_per_user = $47, rate_limit_cleanup_interval = $48,
			auth_totp_enforced = $49, auth_brute_force_enabled = $50, auth_brute_force_threshold = $51, auth_brute_force_window_minutes = $52, auth_brute_force_block_minutes = $53, auth_ip_blacklist_enabled = $54, auth_cleanup_interval = $55,
			created_at = $56, updated_at = $57
		WHERE id = $1
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	dom.UpdatedAt = time.Now()
	return nil
}

func (r *domainRepository) Delete(id int64) error {
	query := `DELETE FROM domains WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}
	return nil
}

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
		LIMIT $1 OFFSET $2
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
