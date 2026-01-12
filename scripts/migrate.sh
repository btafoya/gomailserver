#!/bin/bash
# Data Migration Tool for SQLite to PostgreSQL
# Part of DATABASE-MIGRATION-PLAN.md implementation

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SQLITE_DB_PATH="${SQLITE_DB_PATH:-./data/mailserver.db}"
PG_HOST="${PG_HOST:-localhost}"
PG_PORT="${PG_PORT:-5432}"
PG_DATABASE="${PG_DATABASE:-gomailserver}"
PG_USER="${PG_USER:-gomailserver}"
PG_PASSWORD="${PG_PASSWORD}"

function log_info() {
	echo "[INFO] $*"
}

function log_error() {
	echo "[ERROR] $*" >&2
}

function log_success() {
	echo "[SUCCESS] $*"
}

function show_help() {
	cat << 'EOF'
Data Migration Tool - SQLite to PostgreSQL

Usage: migrate.sh [OPTIONS]

Options:
  -h, --help              Show this help message
  -d, --dry-run           Dry run (show what would happen)
  -v, --verbose            Enable verbose output

Configuration:
  SQLite Database Path    -s, --sqlite-path     Path to SQLite DB (default: ./data/mailserver.db)
  PostgreSQL Host        -H, --pg-host        PostgreSQL host (default: localhost)
  PostgreSQL Port        -P, --pg-port        PostgreSQL port (default: 5432)
  PostgreSQL Database    -D, --pg-database    PostgreSQL database name (default: gomailserver)
  PostgreSQL User        -u, --pg-user        PostgreSQL user (default: gomailserver)
  PostgreSQL Password    -p, --pg-password    PostgreSQL password (default: from env)

Environment Variables:
  PG_HOST, PG_PORT, PG_DATABASE, PG_USER, PG_PASSWORD
  SQL_DB_PATH               Path to SQLite database

Examples:
  # Migrate using pgloader
  ./scripts/migrate.sh --pg-host localhost --pg-password mypass

  # Dry run to see what would happen
  ./scripts/migrate.sh --dry-run

  # Verbose output
  ./scripts/migrate.sh --verbose
EOF
}

function validate_prerequisites() {
	log_info "Validating prerequisites..."
	
	if ! command -v pgloader >/dev/null 2>&1; then
		log_error "pgloader not found. Install with: sudo apt install pgloader"
		return 1
	fi
	
	if [ ! -f "$SQLITE_DB_PATH" ]; then
		log_error "SQLite database not found: $SQLITE_DB_PATH"
		return 1
	fi
	
	log_success "Prerequisites validated"
	return 0
}

function backup_sqlite() {
	local backup_path="./backups/mailserver-pre-migration-$(date +%Y%m%d-%H%M%S).db"
	
	log_info "Creating SQLite backup: $backup_path"
	mkdir -p "$(dirname "$backup_path")"
	
	if ! cp "$SQLITE_DB_PATH" "$backup_path"; then
		log_error "Failed to create backup: $backup_path"
		return 1
	fi
	
	log_success "Backup created: $backup_path"
	return 0
}

function run_pgloader() {
	local pg_dsn="postgresql://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/$PG_DATABASE"
	
	log_info "Running pgloader migration..."
	log_info "  SQLite source: $SQLITE_DB_PATH"
	log_info "  PostgreSQL target: $pg_dsn"
	
	if [ "$DRY_RUN" = "true" ]; then
		echo "[DRY RUN] Would execute: pgloader sqlite://$SQLITE_DB_PATH postgresql://$pg_dsn"
		return 0
	fi
	
	if pgloader sqlite://"$SQLITE_DB_PATH" postgresql://"$pg_dsn"; then
		log_success "Migration completed successfully"
		return 0
	else
		log_error "Migration failed. Check pgloader logs above."
		return 1
	fi
}

function validate_migration() {
	log_info "Validating migration..."
	
	if [ "$DRY_RUN" = "true" ]; then
		echo "[DRY RUN] Would validate data integrity"
		return 0
	fi
	
	# Check row counts for key tables
	local tables=("users" "domains" "mailboxes" "messages" "aliases" "smtp_queue")
	
	for table in "${tables[@]}"; do
		local sqlite_count=$(sqlite3 "$SQLITE_DB_PATH" "SELECT COUNT(*) FROM $table" 2>/dev/null)
		local pg_count=$(psql -h "$PG_HOST" -p "$PG_PORT" -d "$PG_DATABASE" -U "$PG_USER" -t -c "SELECT COUNT(*) FROM $table" 2>/dev/null)
		
		if [ "$sqlite_count" = "$pg_count" ]; then
			log_success "  $table: $sqlite_count = $pg_count rows ✓"
		else
			log_error "  $table: SQLite=$sqlite_count vs PostgreSQL=$pg_count ✗ MISMATCH"
		fi
	done
	
	log_success "Migration validation complete"
	return 0
}

function main() {
	local dry_run=false
	local verbose=false
	
	while [[ $# -gt 0 ]]; do
		case "$1" in
			-h|--help)
				show_help
				exit 0
				;;
			-d|--dry-run)
				dry_run=true
				shift
				;;
			-v|--verbose)
				verbose=true
				shift
				;;
			-s|--sqlite-path)
				SQLITE_DB_PATH="$2"
				shift
				;;
			-H|--pg-host)
				PG_HOST="$2"
				shift
				;;
			-P|--pg-port)
				PG_PORT="$2"
				shift
				;;
			-D|--pg-database)
				PG_DATABASE="$2"
				shift
				;;
			-u|--pg-user)
				PG_USER="$2"
				shift
				;;
			-p|--pg-password)
				PG_PASSWORD="$2"
				shift
				;;
			*)
				log_error "Unknown option: $1"
				show_help
				exit 1
				;;
		esac
	done
	
	if ! validate_prerequisites; then
		exit 1
	fi
	
	if ! backup_sqlite; then
		exit 1
	fi
	
	if ! run_pgloader; then
		exit 1
	fi
	
	if ! validate_migration; then
		exit 1
	fi
	
	log_success "=== Migration complete ==="
	log_info "SQLite backup: ./backups/mailserver-pre-migration-$(date +%Y%m%d-%H%M%S).db"
	log_info "Next steps:"
	log_info "1. Test application with PostgreSQL: Set 'database.driver: postgres' in gomailserver.yaml"
	log_info "2. Verify data integrity: Compare counts, test critical operations"
	log_info "3. Monitor: Watch error logs, performance metrics"
	log_info "4. Rollback: If issues occur, restore from backup"
}

main "$@"
