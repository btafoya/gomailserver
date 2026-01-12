-- Migration V6: Add Advanced Security tables: audit logs, PGP keys, DANE cache, MTA-STS cache, TLS reports
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
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
    success BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp DESC);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_severity ON audit_logs(severity);

CREATE TABLE pgp_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    key_id TEXT NOT NULL,
    fingerprint TEXT NOT NULL,
    public_key TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, fingerprint)
);

CREATE INDEX idx_pgp_keys_user_id ON pgp_keys(user_id);
CREATE INDEX idx_pgp_keys_fingerprint ON pgp_keys(fingerprint);
CREATE INDEX idx_pgp_keys_key_id ON pgp_keys(key_id);

CREATE TABLE dane_tlsa_cache (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL,
    port INTEGER NOT NULL DEFAULT 25,
    usage INTEGER NOT NULL,
    selector INTEGER NOT NULL,
    matching_type INTEGER NOT NULL,
    certificate_association_data TEXT,
    tlsa_record TEXT,
    fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    UNIQUE(domain, port, usage, selector, matching_type)
);

CREATE INDEX idx_dane_tlsa_cache_domain ON dane_tlsa_cache(domain);
CREATE INDEX idx_dane_tlsa_cache_expires ON dane_tlsa_cache(expires_at);

CREATE TABLE mtasts_cache (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL UNIQUE,
    policy TEXT NOT NULL,
    mx TEXT NOT NULL,
    fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_mtasts_cache_expires ON mtasts_cache(expires_at);

CREATE TABLE tls_reports (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL,
    reporter TEXT NOT NULL,
    report_id TEXT NOT NULL,
    received_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    policy TEXT,
    tls_success INTEGER,
    tls_failure INTEGER,
    details TEXT,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

CREATE INDEX idx_tls_reports_domain ON tls_reports(domain_id);
CREATE INDEX idx_tls_reports_received ON tls_reports(received_at);
