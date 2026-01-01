package database

const migrationV5Up = `
-- PostmarkApp Servers (API token groups)
CREATE TABLE postmark_servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    api_token TEXT NOT NULL UNIQUE,
    account_id INTEGER REFERENCES users(id),
    message_stream TEXT DEFAULT 'outbound',
    track_opens INTEGER DEFAULT 0,
    track_links TEXT DEFAULT 'None',
    active INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_servers_token ON postmark_servers(api_token);
CREATE INDEX idx_postmark_servers_account ON postmark_servers(account_id);

-- PostmarkApp Message Tracking
CREATE TABLE postmark_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT NOT NULL UNIQUE,
    server_id INTEGER REFERENCES postmark_servers(id),
    from_email TEXT NOT NULL,
    to_email TEXT NOT NULL,
    cc_email TEXT,
    bcc_email TEXT,
    subject TEXT,
    html_body TEXT,
    text_body TEXT,
    tag TEXT,
    metadata TEXT,
    message_stream TEXT,
    status TEXT DEFAULT 'pending',
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_messages_id ON postmark_messages(message_id);
CREATE INDEX idx_postmark_messages_server ON postmark_messages(server_id);
CREATE INDEX idx_postmark_messages_status ON postmark_messages(status);
CREATE INDEX idx_postmark_messages_submitted ON postmark_messages(submitted_at);

-- PostmarkApp Templates
CREATE TABLE postmark_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER REFERENCES postmark_servers(id),
    name TEXT NOT NULL,
    alias TEXT,
    subject TEXT,
    html_body TEXT,
    text_body TEXT,
    template_type TEXT DEFAULT 'Standard',
    layout_template INTEGER REFERENCES postmark_templates(id),
    active INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(server_id, alias)
);

CREATE INDEX idx_postmark_templates_server ON postmark_templates(server_id);
CREATE INDEX idx_postmark_templates_alias ON postmark_templates(alias);

-- PostmarkApp Webhooks
CREATE TABLE postmark_webhooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER REFERENCES postmark_servers(id),
    url TEXT NOT NULL,
    message_stream TEXT DEFAULT 'outbound',
    http_auth_username TEXT,
    http_auth_password TEXT,
    http_headers TEXT,
    triggers TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_webhooks_server ON postmark_webhooks(server_id);

-- PostmarkApp Bounces
CREATE TABLE postmark_bounces (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT REFERENCES postmark_messages(message_id),
    type TEXT NOT NULL,
    type_code INTEGER,
    email TEXT NOT NULL,
    bounced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    details TEXT,
    inactive INTEGER DEFAULT 1,
    can_activate INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_bounces_message ON postmark_bounces(message_id);
CREATE INDEX idx_postmark_bounces_email ON postmark_bounces(email);

-- PostmarkApp Tracking Events
CREATE TABLE postmark_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT REFERENCES postmark_messages(message_id),
    event_type TEXT NOT NULL,
    recipient TEXT NOT NULL,
    user_agent TEXT,
    client_info TEXT,
    location TEXT,
    link_url TEXT,
    occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_events_message ON postmark_events(message_id);
CREATE INDEX idx_postmark_events_type ON postmark_events(event_type);
CREATE INDEX idx_postmark_events_occurred ON postmark_events(occurred_at);
`

const migrationV5Down = `
DROP TABLE IF EXISTS postmark_events;
DROP TABLE IF EXISTS postmark_bounces;
DROP TABLE IF EXISTS postmark_webhooks;
DROP TABLE IF EXISTS postmark_templates;
DROP TABLE IF EXISTS postmark_messages;
DROP TABLE IF EXISTS postmark_servers;
`
