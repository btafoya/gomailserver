-- Migration V5: Add PostmarkApp API tables for servers, messages, templates, webhooks, bounces, and events
CREATE TABLE postmark_servers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    api_token TEXT NOT NULL UNIQUE,
    account_id INTEGER,
    message_stream TEXT DEFAULT 'outbound',
    track_opens INTEGER DEFAULT 0,
    track_links TEXT DEFAULT 'None',
    active INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES users(id)
);

CREATE INDEX idx_postmark_servers_token ON postmark_servers(api_token);
CREATE INDEX idx_postmark_servers_account ON postmark_servers(account_id);

CREATE TABLE postmark_messages (
    id SERIAL PRIMARY KEY,
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

CREATE INDEX idx_postmark_messages_message_id ON postmark_messages(message_id);
CREATE INDEX idx_postmark_messages_server ON postmark_messages(server_id);
CREATE INDEX idx_postmark_messages_status ON postmark_messages(status);

CREATE TABLE postmark_templates (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    subject TEXT NOT NULL,
    html_body TEXT NOT NULL,
    text_body TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_templates_active ON postmark_templates(active);

CREATE TABLE postmark_bounces (
    id SERIAL PRIMARY KEY,
    message_id TEXT,
    server_id INTEGER REFERENCES postmark_servers(id),
    type TEXT NOT NULL CHECK(type IN ('HardBounce', 'SoftBounce', 'DsnError', 'SpamNotification')),
    description TEXT,
    details TEXT,
    email TEXT NOT NULL,
    bounce_id TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_bounces_message_id ON postmark_bounces(message_id);
CREATE INDEX idx_postmark_bounces_server ON postmark_bounces(server_id);
CREATE INDEX idx_postmark_bounces_email ON postmark_bounces(email);

CREATE TABLE postmark_events (
    id SERIAL PRIMARY KEY,
    message_id TEXT,
    server_id INTEGER REFERENCES postmark_servers(id),
    type TEXT NOT NULL CHECK(type IN ('Open', 'Click', 'Delivery', 'Bounce', 'SpamComplaint')),
    recipient_email TEXT,
    details TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_events_message_id ON postmark_events(message_id);
CREATE INDEX idx_postmark_events_server ON postmark_events(server_id);
CREATE INDEX idx_postmark_events_type ON postmark_events(type);
