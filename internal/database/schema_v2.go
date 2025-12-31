package database

const migrationV2Up = `
-- Add security configuration columns to domains table
-- All per-domain security policies and settings

-- DKIM configuration
ALTER TABLE domains ADD COLUMN dkim_signing_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN dkim_verify_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN dkim_key_size INTEGER DEFAULT 2048;
ALTER TABLE domains ADD COLUMN dkim_key_type TEXT DEFAULT 'rsa' CHECK(dkim_key_type IN ('rsa', 'ed25519'));
ALTER TABLE domains ADD COLUMN dkim_headers_to_sign TEXT DEFAULT '["From","To","Subject","Date","Message-ID","MIME-Version","Content-Type"]';

-- SPF configuration
ALTER TABLE domains ADD COLUMN spf_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN spf_dns_server TEXT DEFAULT '8.8.8.8:53';
ALTER TABLE domains ADD COLUMN spf_dns_timeout INTEGER DEFAULT 5;
ALTER TABLE domains ADD COLUMN spf_max_lookups INTEGER DEFAULT 10;
ALTER TABLE domains ADD COLUMN spf_fail_action TEXT DEFAULT 'reject' CHECK(spf_fail_action IN ('reject','quarantine','accept','tag'));
ALTER TABLE domains ADD COLUMN spf_softfail_action TEXT DEFAULT 'accept' CHECK(spf_softfail_action IN ('reject','quarantine','accept','tag'));

-- DMARC configuration
ALTER TABLE domains ADD COLUMN dmarc_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN dmarc_dns_server TEXT DEFAULT '8.8.8.8:53';
ALTER TABLE domains ADD COLUMN dmarc_dns_timeout INTEGER DEFAULT 5;
ALTER TABLE domains ADD COLUMN dmarc_report_enabled INTEGER DEFAULT 0;
ALTER TABLE domains ADD COLUMN dmarc_report_email TEXT;

-- ClamAV antivirus configuration
ALTER TABLE domains ADD COLUMN clamav_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN clamav_max_scan_size INTEGER DEFAULT 52428800;
ALTER TABLE domains ADD COLUMN clamav_virus_action TEXT DEFAULT 'reject' CHECK(clamav_virus_action IN ('reject','quarantine','tag'));
ALTER TABLE domains ADD COLUMN clamav_fail_action TEXT DEFAULT 'accept' CHECK(clamav_fail_action IN ('reject','quarantine','tag','accept'));

-- SpamAssassin configuration
ALTER TABLE domains ADD COLUMN spam_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN spam_reject_score REAL DEFAULT 10.0;
ALTER TABLE domains ADD COLUMN spam_quarantine_score REAL DEFAULT 5.0;
ALTER TABLE domains ADD COLUMN spam_learning_enabled INTEGER DEFAULT 1;

-- Greylisting configuration
ALTER TABLE domains ADD COLUMN greylist_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN greylist_delay_minutes INTEGER DEFAULT 5;
ALTER TABLE domains ADD COLUMN greylist_expiry_days INTEGER DEFAULT 30;
ALTER TABLE domains ADD COLUMN greylist_cleanup_interval INTEGER DEFAULT 3600;
ALTER TABLE domains ADD COLUMN greylist_whitelist_after INTEGER DEFAULT 3;

-- Rate limiting configuration (JSON objects for per-domain rules)
ALTER TABLE domains ADD COLUMN ratelimit_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN ratelimit_smtp_per_ip TEXT DEFAULT '{"count":100,"window_minutes":60}';
ALTER TABLE domains ADD COLUMN ratelimit_smtp_per_user TEXT DEFAULT '{"count":500,"window_minutes":60}';
ALTER TABLE domains ADD COLUMN ratelimit_smtp_per_domain TEXT DEFAULT '{"count":1000,"window_minutes":60}';
ALTER TABLE domains ADD COLUMN ratelimit_auth_per_ip TEXT DEFAULT '{"count":10,"window_minutes":15}';
ALTER TABLE domains ADD COLUMN ratelimit_imap_per_user TEXT DEFAULT '{"count":1000,"window_minutes":60}';
ALTER TABLE domains ADD COLUMN ratelimit_cleanup_interval INTEGER DEFAULT 300;

-- Authentication security configuration
ALTER TABLE domains ADD COLUMN auth_totp_enforced INTEGER DEFAULT 0;
ALTER TABLE domains ADD COLUMN auth_brute_force_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN auth_brute_force_threshold INTEGER DEFAULT 5;
ALTER TABLE domains ADD COLUMN auth_brute_force_window_minutes INTEGER DEFAULT 15;
ALTER TABLE domains ADD COLUMN auth_brute_force_block_minutes INTEGER DEFAULT 60;
ALTER TABLE domains ADD COLUMN auth_ip_blacklist_enabled INTEGER DEFAULT 1;
ALTER TABLE domains ADD COLUMN auth_cleanup_interval INTEGER DEFAULT 3600;
`

const migrationV2Down = `
-- Remove security configuration columns from domains table
ALTER TABLE domains DROP COLUMN auth_cleanup_interval;
ALTER TABLE domains DROP COLUMN auth_ip_blacklist_enabled;
ALTER TABLE domains DROP COLUMN auth_brute_force_block_minutes;
ALTER TABLE domains DROP COLUMN auth_brute_force_window_minutes;
ALTER TABLE domains DROP COLUMN auth_brute_force_threshold;
ALTER TABLE domains DROP COLUMN auth_brute_force_enabled;
ALTER TABLE domains DROP COLUMN auth_totp_enforced;
ALTER TABLE domains DROP COLUMN ratelimit_cleanup_interval;
ALTER TABLE domains DROP COLUMN ratelimit_imap_per_user;
ALTER TABLE domains DROP COLUMN ratelimit_auth_per_ip;
ALTER TABLE domains DROP COLUMN ratelimit_smtp_per_domain;
ALTER TABLE domains DROP COLUMN ratelimit_smtp_per_user;
ALTER TABLE domains DROP COLUMN ratelimit_smtp_per_ip;
ALTER TABLE domains DROP COLUMN ratelimit_enabled;
ALTER TABLE domains DROP COLUMN greylist_whitelist_after;
ALTER TABLE domains DROP COLUMN greylist_cleanup_interval;
ALTER TABLE domains DROP COLUMN greylist_expiry_days;
ALTER TABLE domains DROP COLUMN greylist_delay_minutes;
ALTER TABLE domains DROP COLUMN greylist_enabled;
ALTER TABLE domains DROP COLUMN spam_learning_enabled;
ALTER TABLE domains DROP COLUMN spam_quarantine_score;
ALTER TABLE domains DROP COLUMN spam_reject_score;
ALTER TABLE domains DROP COLUMN spam_enabled;
ALTER TABLE domains DROP COLUMN clamav_fail_action;
ALTER TABLE domains DROP COLUMN clamav_virus_action;
ALTER TABLE domains DROP COLUMN clamav_max_scan_size;
ALTER TABLE domains DROP COLUMN clamav_enabled;
ALTER TABLE domains DROP COLUMN dmarc_report_email;
ALTER TABLE domains DROP COLUMN dmarc_report_enabled;
ALTER TABLE domains DROP COLUMN dmarc_dns_timeout;
ALTER TABLE domains DROP COLUMN dmarc_dns_server;
ALTER TABLE domains DROP COLUMN dmarc_enabled;
ALTER TABLE domains DROP COLUMN spf_softfail_action;
ALTER TABLE domains DROP COLUMN spf_fail_action;
ALTER TABLE domains DROP COLUMN spf_max_lookups;
ALTER TABLE domains DROP COLUMN spf_dns_timeout;
ALTER TABLE domains DROP COLUMN spf_dns_server;
ALTER TABLE domains DROP COLUMN spf_enabled;
ALTER TABLE domains DROP COLUMN dkim_headers_to_sign;
ALTER TABLE domains DROP COLUMN dkim_key_type;
ALTER TABLE domains DROP COLUMN dkim_key_size;
ALTER TABLE domains DROP COLUMN dkim_verify_enabled;
ALTER TABLE domains DROP COLUMN dkim_signing_enabled;
`
