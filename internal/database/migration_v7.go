package database

const migrationV7Up = `
-- Create webhooks table
CREATE TABLE IF NOT EXISTS webhooks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	url TEXT NOT NULL,
	secret TEXT NOT NULL,
	event_types TEXT NOT NULL,
	active BOOLEAN NOT NULL DEFAULT 1,
	description TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_webhooks_active ON webhooks(active);
CREATE INDEX idx_webhooks_created_at ON webhooks(created_at);

-- Create webhook_deliveries table
CREATE TABLE IF NOT EXISTS webhook_deliveries (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	webhook_id INTEGER NOT NULL,
	event_type TEXT NOT NULL,
	payload TEXT NOT NULL,
	attempt_count INTEGER NOT NULL DEFAULT 0,
	max_attempts INTEGER NOT NULL DEFAULT 10,
	status TEXT NOT NULL DEFAULT 'pending',
	status_code INTEGER,
	response_body TEXT,
	error_message TEXT,
	next_retry_at TIMESTAMP,
	first_attempted_at TIMESTAMP,
	last_attempted_at TIMESTAMP,
	completed_at TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (webhook_id) REFERENCES webhooks(id) ON DELETE CASCADE
);

CREATE INDEX idx_webhook_deliveries_webhook_id ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_next_retry ON webhook_deliveries(next_retry_at);
CREATE INDEX idx_webhook_deliveries_event_type ON webhook_deliveries(event_type);
CREATE INDEX idx_webhook_deliveries_created_at ON webhook_deliveries(created_at);
`

const migrationV7Down = `
DROP TABLE IF EXISTS webhook_deliveries;
DROP TABLE IF NOT EXISTS webhooks;
`
