package database

const migrationV1Up = `
-- Domains table
CREATE TABLE domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'inactive', 'suspended')),
    max_users INTEGER DEFAULT 0,  -- 0 = unlimited
    max_mailbox_size INTEGER DEFAULT 0,  -- 0 = unlimited, in bytes
    default_quota INTEGER DEFAULT 1073741824,  -- 1GB default
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

-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    domain_id INTEGER NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    display_name TEXT,
    quota INTEGER NOT NULL DEFAULT 1073741824,  -- 1GB
    used_quota INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'disabled', 'suspended')),
    auth_method TEXT NOT NULL DEFAULT 'password' CHECK(auth_method IN ('password', 'totp')),
    totp_secret TEXT,
    totp_enabled INTEGER DEFAULT 0,
    forward_to TEXT,
    auto_reply_enabled INTEGER DEFAULT 0,
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

-- Aliases table
CREATE TABLE aliases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias_email TEXT NOT NULL UNIQUE,
    domain_id INTEGER NOT NULL,
    destination_emails TEXT NOT NULL,  -- JSON array or comma-separated
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

CREATE INDEX idx_aliases_email ON aliases(alias_email);
CREATE INDEX idx_aliases_domain_id ON aliases(domain_id);

-- Mailboxes/Folders table
CREATE TABLE mailboxes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    parent_id INTEGER,
    subscribed INTEGER DEFAULT 1,
    special_use TEXT,  -- NULL or 'Inbox', 'Sent', 'Drafts', 'Trash', 'Spam', 'Archive'
    uidvalidity INTEGER NOT NULL,
    uidnext INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES mailboxes(id) ON DELETE CASCADE,
    UNIQUE(user_id, name)
);

CREATE INDEX idx_mailboxes_user_id ON mailboxes(user_id);
CREATE INDEX idx_mailboxes_parent_id ON mailboxes(parent_id);

-- Messages table
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    mailbox_id INTEGER NOT NULL,
    uid INTEGER NOT NULL,
    size INTEGER NOT NULL,
    flags TEXT DEFAULT '',  -- Space-separated flags: \Seen \Deleted \Flagged \Answered \Draft
    categories TEXT DEFAULT '',  -- JSON array: ["Primary", "Social", etc.]
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
    headers TEXT,  -- JSON of all headers
    body_structure TEXT,  -- JSON of MIME structure
    storage_type TEXT NOT NULL DEFAULT 'blob' CHECK(storage_type IN ('blob', 'file')),
    content BLOB,  -- For small messages
    content_path TEXT,  -- For large messages
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

-- SMTP Queue table
CREATE TABLE smtp_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender TEXT NOT NULL,
    recipients TEXT NOT NULL,  -- JSON array
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

-- Failed login attempts
CREATE TABLE failed_logins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL,
    email TEXT,
    protocol TEXT NOT NULL,  -- SMTP, IMAP, API
    attempted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_failed_logins_ip ON failed_logins(ip_address);
CREATE INDEX idx_failed_logins_attempted_at ON failed_logins(attempted_at);

-- IP blacklist
CREATE TABLE ip_blacklist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL UNIQUE,
    reason TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ip_blacklist_ip ON ip_blacklist(ip_address);
CREATE INDEX idx_ip_blacklist_expires_at ON ip_blacklist(expires_at);

-- IP whitelist
CREATE TABLE ip_whitelist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL UNIQUE,
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ip_whitelist_ip ON ip_whitelist(ip_address);

-- Greylisting
CREATE TABLE greylist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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

-- Rate limiting
CREATE TABLE rate_limits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entity_type TEXT NOT NULL,  -- 'ip', 'user', 'domain'
    entity_value TEXT NOT NULL,
    action_type TEXT NOT NULL,  -- 'smtp_send', 'imap_connect', 'api_request'
    count INTEGER NOT NULL DEFAULT 1,
    window_start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(entity_type, entity_value, action_type, window_start)
);

CREATE INDEX idx_rate_limits_lookup ON rate_limits(entity_type, entity_value, action_type, window_start);

-- Sieve scripts
CREATE TABLE sieve_scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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

-- Spam quarantine
CREATE TABLE spam_quarantine (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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

-- Sessions
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    protocol TEXT NOT NULL,  -- SMTP, IMAP, WebDAV, API
    ip_address TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    state TEXT,  -- JSON of session state
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_last_activity ON sessions(last_activity);

-- Webhooks
CREATE TABLE webhooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id INTEGER,
    url TEXT NOT NULL,
    events TEXT NOT NULL,  -- JSON array of event types
    auth_token TEXT,
    active INTEGER DEFAULT 1,
    max_retries INTEGER DEFAULT 3,
    retry_delay INTEGER DEFAULT 60,  -- seconds
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

CREATE INDEX idx_webhooks_domain_id ON webhooks(domain_id);
CREATE INDEX idx_webhooks_active ON webhooks(active);

-- Audit log
CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    action TEXT NOT NULL,
    entity_type TEXT,
    entity_id INTEGER,
    ip_address TEXT,
    details TEXT,  -- JSON
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_action ON audit_log(action);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at);

-- Logs
CREATE TABLE logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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
`

const migrationV1Down = `
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS webhooks;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS spam_quarantine;
DROP TABLE IF EXISTS sieve_scripts;
DROP TABLE IF EXISTS rate_limits;
DROP TABLE IF EXISTS greylist;
DROP TABLE IF EXISTS ip_whitelist;
DROP TABLE IF EXISTS ip_blacklist;
DROP TABLE IF EXISTS failed_logins;
DROP TABLE IF EXISTS smtp_queue;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS mailboxes;
DROP TABLE IF EXISTS aliases;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS domains;
`
