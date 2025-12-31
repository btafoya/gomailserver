package database

const migrationV3Up = `
-- Create API keys table for API authentication
CREATE TABLE IF NOT EXISTS api_keys (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	domain_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	key_hash TEXT NOT NULL UNIQUE,
	scopes TEXT NOT NULL DEFAULT '["read","write"]',
	last_used_at TIMESTAMP,
	last_used_ip TEXT,
	expires_at TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_domain_id ON api_keys(domain_id);

-- Create TLS certificates table for ACME/Let's Encrypt
CREATE TABLE IF NOT EXISTS tls_certificates (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	domain_name TEXT NOT NULL UNIQUE,
	certificate TEXT NOT NULL,
	private_key TEXT NOT NULL,
	issuer TEXT NOT NULL,
	not_before TIMESTAMP NOT NULL,
	not_after TIMESTAMP NOT NULL,
	status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active','expiring','expired','revoked')),
	acme_account_url TEXT,
	acme_order_url TEXT,
	auto_renew INTEGER NOT NULL DEFAULT 1,
	last_renewal_attempt TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tls_certificates_domain_name ON tls_certificates(domain_name);
CREATE INDEX idx_tls_certificates_status ON tls_certificates(status);
CREATE INDEX idx_tls_certificates_not_after ON tls_certificates(not_after);

-- Create setup wizard state table
CREATE TABLE IF NOT EXISTS setup_wizard_state (
	id INTEGER PRIMARY KEY CHECK(id = 1),
	current_step TEXT NOT NULL DEFAULT 'welcome' CHECK(current_step IN ('welcome','system','domain','admin','tls','complete')),
	completed_steps TEXT NOT NULL DEFAULT '[]',
	system_config TEXT,
	domain_config TEXT,
	admin_config TEXT,
	tls_config TEXT,
	started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	completed_at TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial wizard state
INSERT OR IGNORE INTO setup_wizard_state (id, current_step, completed_steps)
VALUES (1, 'welcome', '[]');
`

const migrationV3Down = `
-- Remove setup wizard state table
DROP TABLE IF EXISTS setup_wizard_state;

-- Remove TLS certificates table
DROP INDEX IF EXISTS idx_tls_certificates_not_after;
DROP INDEX IF EXISTS idx_tls_certificates_status;
DROP INDEX IF EXISTS idx_tls_certificates_domain_name;
DROP TABLE IF EXISTS tls_certificates;

-- Remove API keys table
DROP INDEX IF EXISTS idx_api_keys_domain_id;
DROP INDEX IF EXISTS idx_api_keys_user_id;
DROP INDEX IF EXISTS idx_api_keys_key_hash;
DROP TABLE IF EXISTS api_keys;
`
