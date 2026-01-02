package database

// migrationV6Up adds tables for Advanced Security features:
// - Audit logging for admin actions and security events
// - PGP/GPG key storage for email encryption
// - DANE TLSA record cache
// - MTA-STS policy cache
const migrationV6Up = `
-- Audit Log Table
-- Tracks administrative actions and security events
CREATE TABLE IF NOT EXISTS audit_logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id INTEGER,
	username TEXT,
	action TEXT NOT NULL,
	resource_type TEXT NOT NULL,
	resource_id TEXT,
	details TEXT,
	ip_address TEXT,
	user_agent TEXT,
	severity TEXT NOT NULL DEFAULT 'info',
	success BOOLEAN NOT NULL DEFAULT 1,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_severity ON audit_logs(severity);

-- PGP Keys Table
-- Stores user PGP/GPG public keys for encryption
CREATE TABLE IF NOT EXISTS pgp_keys (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	key_id TEXT NOT NULL,
	fingerprint TEXT NOT NULL,
	public_key TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TIMESTAMP,
	is_primary BOOLEAN NOT NULL DEFAULT 0,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	UNIQUE(user_id, fingerprint)
);

CREATE INDEX IF NOT EXISTS idx_pgp_keys_user_id ON pgp_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_pgp_keys_fingerprint ON pgp_keys(fingerprint);
CREATE INDEX IF NOT EXISTS idx_pgp_keys_key_id ON pgp_keys(key_id);

-- DANE TLSA Cache Table
-- Caches DANE TLSA records for SMTP TLS verification
CREATE TABLE IF NOT EXISTS dane_tlsa_cache (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	domain TEXT NOT NULL,
	port INTEGER NOT NULL DEFAULT 25,
	usage INTEGER NOT NULL,
	selector INTEGER NOT NULL,
	matching_type INTEGER NOT NULL,
	certificate_data TEXT NOT NULL,
	fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ttl INTEGER NOT NULL DEFAULT 3600,
	dnssec_verified BOOLEAN NOT NULL DEFAULT 0,
	UNIQUE(domain, port, usage, selector, matching_type)
);

CREATE INDEX IF NOT EXISTS idx_dane_tlsa_domain ON dane_tlsa_cache(domain);
CREATE INDEX IF NOT EXISTS idx_dane_tlsa_fetched ON dane_tlsa_cache(fetched_at);

-- MTA-STS Policy Cache Table
-- Caches MTA-STS policies for enforcing TLS
CREATE TABLE IF NOT EXISTS mtasts_policy_cache (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	domain TEXT NOT NULL UNIQUE,
	version TEXT NOT NULL,
	mode TEXT NOT NULL,
	max_age INTEGER NOT NULL,
	mx_patterns TEXT NOT NULL,
	fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TIMESTAMP NOT NULL,
	policy_text TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mtasts_domain ON mtasts_policy_cache(domain);
CREATE INDEX IF NOT EXISTS idx_mtasts_expires ON mtasts_policy_cache(expires_at);

-- TLS Reporting Table
-- Stores TLS reporting data for TLSRPT (RFC 8460)
CREATE TABLE IF NOT EXISTS tls_reports (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	report_id TEXT NOT NULL,
	domain TEXT NOT NULL,
	date_range_start TIMESTAMP NOT NULL,
	date_range_end TIMESTAMP NOT NULL,
	contact_info TEXT,
	report_json TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	sent_at TIMESTAMP,
	UNIQUE(report_id, domain)
);

CREATE INDEX IF NOT EXISTS idx_tls_reports_domain ON tls_reports(domain);
CREATE INDEX IF NOT EXISTS idx_tls_reports_date ON tls_reports(date_range_start, date_range_end);
CREATE INDEX IF NOT EXISTS idx_tls_reports_sent ON tls_reports(sent_at);
`

// migrationV6Down removes Advanced Security tables
const migrationV6Down = `
DROP TABLE IF EXISTS tls_reports;
DROP TABLE IF EXISTS mtasts_policy_cache;
DROP TABLE IF EXISTS dane_tlsa_cache;
DROP TABLE IF EXISTS pgp_keys;
DROP TABLE IF EXISTS audit_logs;
`
