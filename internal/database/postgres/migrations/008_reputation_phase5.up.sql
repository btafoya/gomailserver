-- Migration V8: Add reputation management Phase 5 tables (DMARC, ARF, external metrics, predictions, alerts)

CREATE TABLE dmarc_reports (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL,
    report_id TEXT NOT NULL UNIQUE,
    begin_time BIGINT NOT NULL,
    end_time BIGINT NOT NULL,
    organization TEXT,
    total_messages INTEGER NOT NULL,
    spf_pass INTEGER NOT NULL,
    dkim_pass INTEGER NOT NULL,
    alignment_pass INTEGER NOT NULL,
    raw_xml TEXT,
    processed_at BIGINT NOT NULL,
    UNIQUE(report_id)
);

CREATE INDEX idx_dmarc_reports_domain ON dmarc_reports(domain);
CREATE INDEX idx_dmarc_reports_time ON dmarc_reports(begin_time);

CREATE TABLE dmarc_report_records (
    id SERIAL PRIMARY KEY,
    report_id INTEGER NOT NULL REFERENCES dmarc_reports(id) ON DELETE CASCADE,
    source_ip TEXT NOT NULL,
    count INTEGER NOT NULL,
    disposition TEXT,
    spf_result TEXT,
    dkim_result TEXT,
    spf_aligned BOOLEAN,
    dkim_aligned BOOLEAN,
    header_from TEXT,
    envelope_from TEXT
);

CREATE INDEX idx_dmarc_records_report ON dmarc_report_records(report_id);
CREATE INDEX idx_dmarc_records_ip ON dmarc_report_records(source_ip);

CREATE TABLE dmarc_auto_actions (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL,
    issue_type TEXT NOT NULL,
    description TEXT,
    action_taken TEXT,
    taken_at BIGINT NOT NULL,
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT
);

CREATE INDEX idx_dmarc_actions_domain ON dmarc_auto_actions(domain);
CREATE INDEX idx_dmarc_actions_time ON dmarc_auto_actions(taken_at);

CREATE TABLE postmaster_metrics (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL,
    fetched_at BIGINT NOT NULL,
    metric_date BIGINT NOT NULL,
    domain_reputation TEXT,
    spam_rate REAL,
    ip_reputation TEXT,
    authentication_rate REAL,
    user_spam_reports INTEGER,
    raw_response TEXT
);

CREATE INDEX idx_postmaster_domain ON postmaster_metrics(domain);
CREATE INDEX idx_postmaster_date ON postmaster_metrics(metric_date);
CREATE INDEX idx_postmaster_fetched ON postmaster_metrics(fetched_at);

CREATE TABLE snds_metrics (
    id SERIAL PRIMARY KEY,
    ip_address TEXT NOT NULL,
    fetched_at BIGINT NOT NULL,
    metric_date BIGINT NOT NULL,
    spam_trap_hits INTEGER DEFAULT 0,
    complaint_rate REAL,
    filter_level TEXT,
    message_count INTEGER,
    raw_response TEXT
);

CREATE INDEX idx_snds_ip ON snds_metrics(ip_address);
CREATE INDEX idx_snds_date ON snds_metrics(metric_date);
CREATE INDEX idx_snds_fetched ON snds_metrics(fetched_at);

CREATE TABLE provider_rate_limits (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL,
    provider TEXT NOT NULL,
    max_hourly_rate INTEGER NOT NULL,
    max_daily_rate INTEGER,
    current_hour_count INTEGER DEFAULT 0,
    current_day_count INTEGER DEFAULT 0,
    hour_reset_at BIGINT NOT NULL,
    day_reset_at BIGINT NOT NULL,
    circuit_breaker_active BOOLEAN DEFAULT FALSE,
    last_updated BIGINT NOT NULL,
    UNIQUE(domain, provider)
);

CREATE INDEX idx_provider_limits_domain ON provider_rate_limits(domain);
CREATE INDEX idx_provider_limits_provider ON provider_rate_limits(provider);

CREATE TABLE custom_warmup_schedules (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL,
    schedule_name TEXT NOT NULL,
    day INTEGER NOT NULL,
    max_volume INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(domain, day)
);

CREATE INDEX idx_custom_warmup_domain ON custom_warmup_schedules(domain);

CREATE TABLE arf_reports (
    id SERIAL PRIMARY KEY,
    received_at BIGINT NOT NULL,
    feedback_type TEXT,
    user_agent TEXT,
    version TEXT,
    original_rcpt_to TEXT,
    arrival_date BIGINT,
    reporting_mta TEXT,
    source_ip TEXT,
    authentication_results TEXT,
    message_id TEXT,
    subject TEXT,
    raw_report TEXT,
    processed BOOLEAN DEFAULT FALSE,
    suppressed_recipient TEXT
);

CREATE INDEX idx_arf_received ON arf_reports(received_at);
CREATE INDEX idx_arf_recipient ON arf_reports(original_rcpt_to);
CREATE INDEX idx_arf_processed ON arf_reports(processed);

CREATE TABLE reputation_predictions (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL,
    predicted_at BIGINT NOT NULL,
    prediction_horizon INTEGER NOT NULL,
    predicted_score INTEGER,
    predicted_complaint_rate REAL,
    predicted_bounce_rate REAL,
    confidence_level REAL,
    model_version TEXT,
    features_used TEXT
);

CREATE INDEX idx_predictions_domain ON reputation_predictions(domain);
CREATE INDEX idx_predictions_time ON reputation_predictions(predicted_at);

CREATE TABLE reputation_alerts (
    id SERIAL PRIMARY KEY,
    domain TEXT NOT NULL,
    alert_type TEXT NOT NULL,
    severity TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT,
    details TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    acknowledged BOOLEAN DEFAULT FALSE,
    acknowledged_at BIGINT,
    acknowledged_by TEXT,
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at BIGINT
);

CREATE INDEX idx_alerts_domain ON reputation_alerts(domain);
CREATE INDEX idx_alerts_type ON reputation_alerts(alert_type);
CREATE INDEX idx_alerts_created ON reputation_alerts(created_at);
CREATE INDEX idx_alerts_unacknowledged ON reputation_alerts(acknowledged) WHERE acknowledged = FALSE;
