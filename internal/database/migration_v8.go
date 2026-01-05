package database

// Migration v8: Reputation Management Phase 5 - Advanced Automation
// This migration adds tables for DMARC reports, ARF complaints, external metrics,
// provider rate limits, custom warm-up schedules, predictions, and alerts.

const migrationV8Up = SchemaReputationV2

const migrationV8Down = `
-- Drop all Phase 5 reputation tables in reverse order of dependencies
DROP TABLE IF EXISTS reputation_alerts;
DROP TABLE IF EXISTS reputation_predictions;
DROP TABLE IF EXISTS arf_reports;
DROP TABLE IF EXISTS custom_warmup_schedules;
DROP TABLE IF EXISTS provider_rate_limits;
DROP TABLE IF EXISTS snds_metrics;
DROP TABLE IF EXISTS postmaster_metrics;
DROP TABLE IF EXISTS dmarc_auto_actions;
DROP TABLE IF EXISTS dmarc_report_records;
DROP TABLE IF EXISTS dmarc_reports;
`
