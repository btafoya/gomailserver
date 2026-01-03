package database

const migrationReputationV1Up = `
-- Reputation metrics database (separate SQLite: reputation.db)

-- Sending events table
CREATE TABLE sending_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,                 -- Unix timestamp
    domain TEXT NOT NULL,                       -- Sending domain
    recipient_domain TEXT NOT NULL,             -- Receiving domain
    event_type TEXT NOT NULL,                   -- delivery|bounce|defer|complaint
    bounce_type TEXT,                           -- hard|soft|null
    enhanced_status_code TEXT,                  -- e.g., "5.1.1"
    smtp_response TEXT,                         -- Full SMTP response
    ip_address TEXT NOT NULL,                   -- Sending IP
    metadata TEXT                               -- JSON: additional context
);

CREATE INDEX idx_sending_events_timestamp ON sending_events(timestamp);
CREATE INDEX idx_sending_events_domain ON sending_events(domain);
CREATE INDEX idx_sending_events_event_type ON sending_events(event_type);
CREATE INDEX idx_sending_events_recipient_domain ON sending_events(recipient_domain);

-- Domain reputation scores table
CREATE TABLE domain_reputation_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL UNIQUE,
    reputation_score INTEGER NOT NULL,          -- 0-100
    complaint_rate REAL NOT NULL,               -- Percentage
    bounce_rate REAL NOT NULL,                  -- Percentage
    delivery_rate REAL NOT NULL,                -- Percentage
    circuit_breaker_active BOOLEAN DEFAULT 0,
    circuit_breaker_reason TEXT,
    warm_up_active BOOLEAN DEFAULT 0,
    warm_up_day INTEGER DEFAULT 0,              -- Day in warm-up schedule
    last_updated INTEGER NOT NULL               -- Unix timestamp
);

CREATE INDEX idx_domain_reputation_domain ON domain_reputation_scores(domain);
CREATE INDEX idx_domain_reputation_score ON domain_reputation_scores(reputation_score);

-- Warm-up schedules table
CREATE TABLE warm_up_schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    day INTEGER NOT NULL,                       -- Day 1, 2, 3...
    max_volume INTEGER NOT NULL,                -- Max messages for this day
    actual_volume INTEGER DEFAULT 0,            -- Messages sent today
    created_at INTEGER NOT NULL
);

CREATE INDEX idx_warm_up_domain ON warm_up_schedules(domain);
CREATE INDEX idx_warm_up_day ON warm_up_schedules(domain, day);

-- Circuit breaker events table
CREATE TABLE circuit_breaker_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    trigger_type TEXT NOT NULL,                 -- complaint|bounce|block
    trigger_value REAL NOT NULL,                -- Rate or count
    threshold REAL NOT NULL,                    -- What triggered it
    paused_at INTEGER NOT NULL,                 -- Unix timestamp
    resumed_at INTEGER,                         -- Unix timestamp or null
    auto_resumed BOOLEAN DEFAULT 0,
    admin_notes TEXT
);

CREATE INDEX idx_circuit_breaker_domain ON circuit_breaker_events(domain);
CREATE INDEX idx_circuit_breaker_active ON circuit_breaker_events(domain, resumed_at);

-- Retention policy table
CREATE TABLE retention_policy (
    id INTEGER PRIMARY KEY,
    retention_days INTEGER NOT NULL DEFAULT 90,
    last_cleanup INTEGER NOT NULL              -- Unix timestamp
);

INSERT INTO retention_policy (id, retention_days, last_cleanup)
VALUES (1, 90, strftime('%s', 'now'));
`

const migrationReputationV1Down = `
DROP TABLE IF EXISTS retention_policy;
DROP TABLE IF EXISTS circuit_breaker_events;
DROP TABLE IF EXISTS warm_up_schedules;
DROP TABLE IF EXISTS domain_reputation_scores;
DROP TABLE IF EXISTS sending_events;
`
