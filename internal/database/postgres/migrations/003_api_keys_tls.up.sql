-- Migration V3: Add API keys, TLS certificates, and setup wizard tables
CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
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

CREATE TABLE tls_certificates (
    id SERIAL PRIMARY KEY,
    domain_name TEXT NOT NULL UNIQUE,
    certificate TEXT NOT NULL,
    private_key TEXT NOT NULL,
    issuer TEXT,
    not_before TIMESTAMP NOT NULL,
    not_after TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tls_certificates_domain_name ON tls_certificates(domain_name);
