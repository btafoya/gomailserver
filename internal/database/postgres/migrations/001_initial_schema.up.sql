-- Migration V1: Initial schema - create all tables
CREATE TABLE domains (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'inactive', 'suspended')),
    max_users INTEGER DEFAULT 0,
    max_mailbox_size INTEGER DEFAULT 0,
    default_quota BIGINT DEFAULT 1073741824,
    catchall_email TEXT,
    backup_mx INTEGER DEFAULT 0,
    dkim_selector TEXT,
    dkim_private_key TEXT,
    dkim_public_key TEXT,
    spf_record TEXT,
    dmarc_policy TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_domains_name ON domains(name);
CREATE INDEX idx_domains_status ON domains(status);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    domain_id INTEGER NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    display_name TEXT,
    role TEXT DEFAULT 'user',
    quota BIGINT NOT NULL DEFAULT 1073741824,
    used_quota BIGINT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'disabled', 'suspended')),
    auth_method TEXT NOT NULL DEFAULT 'password' CHECK(auth_method IN ('password', 'totp')),
    totp_secret TEXT,
    totp_enabled BOOLEAN DEFAULT FALSE,
    forward_to TEXT,
    auto_reply_enabled BOOLEAN DEFAULT FALSE,
    auto_reply_subject TEXT,
    auto_reply_body TEXT,
    spam_threshold REAL DEFAULT 5.0,
    language TEXT DEFAULT 'en',
    last_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_domain_id ON users(domain_id);
CREATE INDEX idx_users_status ON users(status);

CREATE TABLE aliases (
    id SERIAL PRIMARY KEY,
    alias_email TEXT NOT NULL UNIQUE,
    domain_id INTEGER NOT NULL,
    destination_emails TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

CREATE INDEX idx_aliases_email ON aliases(alias_email);
CREATE INDEX idx_aliases_domain_id ON aliases(domain_id);

CREATE TABLE mailboxes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    parent_id INTEGER,
    subscribed INTEGER DEFAULT 1,
    special_use TEXT,
    uidvalidity INTEGER NOT NULL,
    uidnext INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES mailboxes(id) ON DELETE CASCADE,
    UNIQUE(user_id, name)
);

CREATE INDEX idx_mailboxes_user_id ON mailboxes(user_id);
CREATE INDEX idx_mailboxes_parent_id ON mailboxes(parent_id);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    mailbox_id INTEGER NOT NULL,
    uid INTEGER NOT NULL,
    size INTEGER NOT NULL,
    flags TEXT DEFAULT '',
    categories TEXT DEFAULT '',
    thread_id TEXT,
    received_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    internal_date TIMESTAMP NOT NULL,
    subject TEXT,
    from_addr TEXT,
    to_addr TEXT,
    cc_addr TEXT,
    bcc_addr TEXT,
    reply_to TEXT,
    message_id TEXT,
    in_reply_to TEXT,
    refs TEXT,
    headers TEXT,
    body_structure TEXT,
    storage_type TEXT NOT NULL DEFAULT 'blob' CHECK(storage_type IN ('blob', 'file')),
    content BYTEA,
    content_path TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (mailbox_id) REFERENCES mailboxes(id) ON DELETE CASCADE,
    UNIQUE(mailbox_id, uid)
);

CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_mailbox_id ON messages(mailbox_id);
CREATE INDEX idx_messages_uid ON messages(uid);
CREATE INDEX idx_messages_thread_id ON messages(thread_id);
CREATE INDEX idx_messages_subject ON messages(subject);
CREATE INDEX idx_messages_from_addr ON messages(from_addr);
CREATE INDEX idx_messages_received_at ON messages(received_at);

CREATE TABLE smtp_queue (
    id SERIAL PRIMARY KEY,
    sender TEXT NOT NULL,
    recipients TEXT NOT NULL,
    message_id TEXT,
    message_path TEXT NOT NULL,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 5,
    next_retry TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'processing', 'failed', 'delivered')),
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_smtp_queue_status ON smtp_queue(status);
CREATE INDEX idx_smtp_queue_next_retry ON smtp_queue(next_retry);

CREATE TABLE failed_logins (
    id SERIAL PRIMARY KEY,
    ip_address TEXT NOT NULL,
    email TEXT,
    protocol TEXT NOT NULL,
    attempted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_failed_logins_ip ON failed_logins(ip_address);
CREATE INDEX idx_failed_logins_attempted_at ON failed_logins(attempted_at);

CREATE TABLE ip_blacklist (
    id SERIAL PRIMARY KEY,
    ip_address TEXT NOT NULL UNIQUE,
    reason TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ip_blacklist_ip ON ip_blacklist(ip_address);
CREATE INDEX idx_ip_blacklist_expires_at ON ip_blacklist(expires_at);

CREATE TABLE ip_whitelist (
    id SERIAL PRIMARY KEY,
    ip_address TEXT NOT NULL UNIQUE,
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ip_whitelist_ip ON ip_whitelist(ip_address);

CREATE TABLE greylist (
    id SERIAL PRIMARY KEY,
    sender_ip TEXT NOT NULL,
    sender_email TEXT NOT NULL,
    recipient_email TEXT NOT NULL,
    first_seen TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    passed_at TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'greylisted' CHECK(status IN ('greylisted', 'passed', 'expired')),
    UNIQUE(sender_ip, sender_email, recipient_email)
);

CREATE INDEX idx_greylist_lookup ON greylist(sender_ip, sender_email, recipient_email);
CREATE INDEX idx_greylist_status ON greylist(status);

CREATE TABLE rate_limits (
    id SERIAL PRIMARY KEY,
    entity_type TEXT NOT NULL,
    entity_value TEXT NOT NULL,
    action_type TEXT NOT NULL,
    count INTEGER NOT NULL DEFAULT 1,
    window_start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(entity_type, entity_value, action_type, window_start)
);

CREATE INDEX idx_rate_limits_lookup ON rate_limits(entity_type, entity_value, action_type, window_start);

CREATE TABLE sieve_scripts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    active INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, name)
);

CREATE INDEX idx_sieve_scripts_user_id ON sieve_scripts(user_id);
CREATE INDEX idx_sieve_scripts_active ON sieve_scripts(active);

CREATE TABLE spam_quarantine (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    message_id INTEGER NOT NULL,
    spam_score REAL NOT NULL,
    quarantined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    released_at TIMESTAMP,
    auto_delete_at TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'quarantined' CHECK(status IN ('quarantined', 'released', 'deleted')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
);

CREATE INDEX idx_spam_quarantine_user_id ON spam_quarantine(user_id);
CREATE INDEX idx_spam_quarantine_status ON spam_quarantine(status);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    protocol TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    state TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_last_activity ON sessions(last_activity);

CREATE TABLE webhooks (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER,
    url TEXT NOT NULL,
    events TEXT NOT NULL,
    auth_token TEXT,
    active INTEGER DEFAULT 1,
    max_retries INTEGER DEFAULT 3,
    retry_delay INTEGER DEFAULT 60,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

CREATE INDEX idx_webhooks_domain_id ON webhooks(domain_id);
CREATE INDEX idx_webhooks_active ON webhooks(active);

CREATE TABLE audit_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    action TEXT NOT NULL,
    entity_type TEXT,
    entity_id INTEGER,
    ip_address TEXT,
    details TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_action ON audit_log(action);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at);

CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    level TEXT NOT NULL,
    service TEXT NOT NULL,
    user_email TEXT,
    ip_address TEXT,
    action TEXT,
    result TEXT,
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_logs_level ON logs(level);
CREATE INDEX idx_logs_service ON logs(service);
CREATE INDEX idx_logs_created_at ON logs(created_at);
