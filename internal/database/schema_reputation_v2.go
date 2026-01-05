package database

// Schema v2: DMARC Reports and External Metrics
// This extends the reputation database with Phase 5 Advanced Automation features

const SchemaReputationV2 = `
-- DMARC Reports
CREATE TABLE IF NOT EXISTS dmarc_reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    report_id TEXT NOT NULL UNIQUE,          -- From DMARC report
    begin_time INTEGER NOT NULL,             -- Unix timestamp
    end_time INTEGER NOT NULL,               -- Unix timestamp
    organization TEXT,                        -- Reporter org (e.g., "google.com")
    total_messages INTEGER NOT NULL,
    spf_pass INTEGER NOT NULL,
    dkim_pass INTEGER NOT NULL,
    alignment_pass INTEGER NOT NULL,
    raw_xml TEXT,                            -- Full report for debugging
    processed_at INTEGER NOT NULL,           -- Unix timestamp
    UNIQUE(report_id)
);

CREATE INDEX IF NOT EXISTS idx_dmarc_reports_domain ON dmarc_reports(domain);
CREATE INDEX IF NOT EXISTS idx_dmarc_reports_time ON dmarc_reports(begin_time);

-- DMARC Report Records (individual source IPs and results)
CREATE TABLE IF NOT EXISTS dmarc_report_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL REFERENCES dmarc_reports(id) ON DELETE CASCADE,
    source_ip TEXT NOT NULL,
    count INTEGER NOT NULL,
    disposition TEXT,                        -- none|quarantine|reject
    spf_result TEXT,                         -- pass|fail|neutral|softfail|temperror|permerror
    dkim_result TEXT,                        -- pass|fail|neutral|temperror|permerror
    spf_aligned BOOLEAN,
    dkim_aligned BOOLEAN,
    header_from TEXT,                        -- RFC5322.From domain
    envelope_from TEXT                       -- SMTP MAIL FROM domain
);

CREATE INDEX IF NOT EXISTS idx_dmarc_records_report ON dmarc_report_records(report_id);
CREATE INDEX IF NOT EXISTS idx_dmarc_records_ip ON dmarc_report_records(source_ip);

-- DMARC Automated Actions Log
CREATE TABLE IF NOT EXISTS dmarc_auto_actions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    issue_type TEXT NOT NULL,               -- spf_misalign|dkim_misalign|dkim_fail|spf_fail
    description TEXT,
    action_taken TEXT,                      -- e.g., "Updated SPF record", "Rotated DKIM key"
    taken_at INTEGER NOT NULL,              -- Unix timestamp
    success BOOLEAN DEFAULT 1,
    error_message TEXT
);

CREATE INDEX IF NOT EXISTS idx_dmarc_actions_domain ON dmarc_auto_actions(domain);
CREATE INDEX IF NOT EXISTS idx_dmarc_actions_time ON dmarc_auto_actions(taken_at);

-- Gmail Postmaster Tools Metrics
CREATE TABLE IF NOT EXISTS postmaster_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    fetched_at INTEGER NOT NULL,            -- Unix timestamp
    metric_date INTEGER NOT NULL,           -- Date of the metric (Unix timestamp)
    domain_reputation TEXT,                  -- HIGH|MEDIUM|LOW|BAD
    spam_rate REAL,                          -- 0.0 to 1.0
    ip_reputation TEXT,                      -- HIGH|MEDIUM|LOW|BAD
    authentication_rate REAL,                -- 0.0 to 1.0 (SPF/DKIM/DMARC pass rate)
    encryption_rate REAL,                    -- 0.0 to 1.0 (TLS rate)
    user_spam_reports INTEGER,               -- Number of user spam reports
    raw_response TEXT                        -- JSON response for debugging
);

CREATE INDEX IF NOT EXISTS idx_postmaster_domain ON postmaster_metrics(domain);
CREATE INDEX IF NOT EXISTS idx_postmaster_date ON postmaster_metrics(metric_date);
CREATE INDEX IF NOT EXISTS idx_postmaster_fetched ON postmaster_metrics(fetched_at);

-- Microsoft SNDS Metrics
CREATE TABLE IF NOT EXISTS snds_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL,
    fetched_at INTEGER NOT NULL,            -- Unix timestamp
    metric_date INTEGER NOT NULL,           -- Date of the metric (Unix timestamp)
    spam_trap_hits INTEGER DEFAULT 0,
    complaint_rate REAL,                     -- 0.0 to 1.0
    filter_level TEXT,                       -- GREEN|YELLOW|RED
    message_count INTEGER,                   -- Total messages seen by Microsoft
    raw_response TEXT                        -- JSON response for debugging
);

CREATE INDEX IF NOT EXISTS idx_snds_ip ON snds_metrics(ip_address);
CREATE INDEX IF NOT EXISTS idx_snds_date ON snds_metrics(metric_date);
CREATE INDEX IF NOT EXISTS idx_snds_fetched ON snds_metrics(fetched_at);

-- Provider-Specific Rate Limits
CREATE TABLE IF NOT EXISTS provider_rate_limits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    provider TEXT NOT NULL,                 -- gmail|outlook|yahoo|generic
    max_hourly_rate INTEGER NOT NULL,       -- Messages per hour
    max_daily_rate INTEGER,                 -- Messages per day (optional)
    current_hour_count INTEGER DEFAULT 0,
    current_day_count INTEGER DEFAULT 0,
    hour_reset_at INTEGER NOT NULL,         -- Unix timestamp
    day_reset_at INTEGER NOT NULL,          -- Unix timestamp
    circuit_breaker_active BOOLEAN DEFAULT 0,
    last_updated INTEGER NOT NULL,
    UNIQUE(domain, provider)
);

CREATE INDEX IF NOT EXISTS idx_provider_limits_domain ON provider_rate_limits(domain);
CREATE INDEX IF NOT EXISTS idx_provider_limits_provider ON provider_rate_limits(provider);

-- Custom Warm-up Schedules
CREATE TABLE IF NOT EXISTS custom_warmup_schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    schedule_name TEXT NOT NULL,
    day INTEGER NOT NULL,                   -- Day 1, 2, 3, etc.
    max_volume INTEGER NOT NULL,            -- Max messages for this day
    created_at INTEGER NOT NULL,
    created_by TEXT,                        -- Admin user who created it
    is_active BOOLEAN DEFAULT 1,
    UNIQUE(domain, day)
);

CREATE INDEX IF NOT EXISTS idx_custom_warmup_domain ON custom_warmup_schedules(domain);

-- ARF Complaint Reports
CREATE TABLE IF NOT EXISTS arf_reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    received_at INTEGER NOT NULL,           -- Unix timestamp
    feedback_type TEXT,                     -- abuse|fraud|virus|other
    user_agent TEXT,                        -- Reporter (e.g., "Yahoo! Inc.")
    version TEXT,                           -- ARF version
    original_rcpt_to TEXT,                  -- Original recipient
    arrival_date INTEGER,                   -- When original message arrived
    reporting_mta TEXT,                     -- MTA reporting the complaint
    source_ip TEXT,                         -- IP that sent the complained message
    authentication_results TEXT,            -- SPF/DKIM/DMARC results
    message_id TEXT,                        -- Original Message-ID
    subject TEXT,                           -- Original subject
    raw_report TEXT,                        -- Full ARF message
    processed BOOLEAN DEFAULT 0,
    suppressed_recipient TEXT               -- Recipient that was suppressed
);

CREATE INDEX IF NOT EXISTS idx_arf_received ON arf_reports(received_at);
CREATE INDEX IF NOT EXISTS idx_arf_recipient ON arf_reports(original_rcpt_to);
CREATE INDEX IF NOT EXISTS idx_arf_processed ON arf_reports(processed);

-- Reputation Predictions (ML/Trend-based)
CREATE TABLE IF NOT EXISTS reputation_predictions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    predicted_at INTEGER NOT NULL,          -- Unix timestamp
    prediction_horizon INTEGER NOT NULL,    -- Hours ahead (e.g., 24, 48, 72)
    predicted_score INTEGER,                -- 0-100
    predicted_complaint_rate REAL,
    predicted_bounce_rate REAL,
    confidence_level REAL,                  -- 0.0 to 1.0
    model_version TEXT,                     -- e.g., "trend-v1" or "ml-v1"
    features_used TEXT                      -- JSON: features used for prediction
);

CREATE INDEX IF NOT EXISTS idx_predictions_domain ON reputation_predictions(domain);
CREATE INDEX IF NOT EXISTS idx_predictions_time ON reputation_predictions(predicted_at);

-- Alert System
CREATE TABLE IF NOT EXISTS reputation_alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    alert_type TEXT NOT NULL,               -- dns_failure|score_drop|circuit_breaker|external_feedback|dmarc_issue
    severity TEXT NOT NULL,                 -- low|medium|high|critical
    title TEXT NOT NULL,
    message TEXT,
    details TEXT,                           -- JSON: additional context
    created_at INTEGER NOT NULL,            -- Unix timestamp
    acknowledged BOOLEAN DEFAULT 0,
    acknowledged_at INTEGER,
    acknowledged_by TEXT,                   -- Admin user
    resolved BOOLEAN DEFAULT 0,
    resolved_at INTEGER
);

CREATE INDEX IF NOT EXISTS idx_alerts_domain ON reputation_alerts(domain);
CREATE INDEX IF NOT EXISTS idx_alerts_type ON reputation_alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_alerts_created ON reputation_alerts(created_at);
CREATE INDEX IF NOT EXISTS idx_alerts_unacknowledged ON reputation_alerts(acknowledged) WHERE acknowledged = 0;
`
