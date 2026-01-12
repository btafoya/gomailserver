#!/bin/bash
# Validation Scripts for PostgreSQL Migration Data Integrity
# Part of DATABASE-MIGRATION-PLAN.md implementation

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
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
PostgreSQL Data Validation Scripts

Usage: validate.sh [OPTIONS] [ACTION]

Options:
  -h, --help              Show this help message
  -v, --verbose            Enable verbose output

Configuration:
  PostgreSQL Host        -H, --pg-host        PostgreSQL host (default: localhost)
  PostgreSQL Port        -P, --pg-port        PostgreSQL port (default: 5432)
  PostgreSQL Database    -D, --pg-database    PostgreSQL database name (default: gomailserver)
  PostgreSQL User        -u, --pg-user        PostgreSQL user (default: gomailserver)
  PostgreSQL Password    -p, --pg-password    PostgreSQL password (default: from env)

Environment Variables:
  PG_HOST, PG_PORT, PG_DATABASE, PG_USER, PG_PASSWORD

Actions:
  row-counts               Compare row counts between SQLite and PostgreSQL
  foreign-keys              Validate foreign key constraints
  data-types                Check data type conversions
  indexes                   Verify indexes exist
  null-checks               Check for NULL/empty values

Examples:
  # Validate row counts
  ./scripts/validate.sh --pg-host localhost --pg-password pass row-counts
  
  # Validate foreign keys
  ./scripts/validate.sh foreign-keys
  
  # Run all validations
  ./scripts/validate.sh row-counts foreign-keys data-types indexes null-checks

  # Verbose mode
  ./scripts/validate.sh -v --pg-host localhost --pg-password pass row-counts
EOF
}

function check_row_counts() {
	local pg_dsn="postgresql://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/$PG_DATABASE"
	
	log_info "Checking row counts..."
	
	if [ "$DRY_RUN" = "true" ]; then
		echo "[DRY RUN] Would compare row counts"
		return 0
	fi
	
	local tables=("users" "domains" "mailboxes" "messages" "aliases" "smtp_queue"
	local all_valid=true
	
	for table in "${tables[@]}"; do
		local sqlite_count=$(sqlite3 "$SQLITE_DB_PATH" "SELECT COUNT(*) FROM $table" 2>/dev/null)
		local pg_count=$(psql -h "$PG_HOST" -p "$PG_PORT" -d "$PG_DATABASE" -U "$PG_USER" -t -c "SELECT COUNT(*) FROM $table" 2>/dev/null)
		
		if [ "$sqlite_count" = "$pg_count" ]; then
			log_success "  ✓ $table: $sqlite_count = $pg_count"
		else
			log_error "  ✗ $table: SQLite=$sqlite_count vs PostgreSQL=$pg_count"
			all_valid=false
		fi
	done
	
	if [ "$all_valid" = "true" ]; then
		log_success "All row counts match"
		return 0
	else
		return 1
}

function check_foreign_keys() {
	local pg_dsn="postgresql://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/$PG_DATABASE"
	
	log_info "Checking foreign key constraints..."
	
	# Check users.domain_id foreign key
	local orphan_users=$(psql -h "$PG_HOST" -p "$PG_PORT" -d "$PG_DATABASE" -U "$PG_USER" -t -c "
		SELECT COUNT(*) FROM users u WHERE u.domain_id IS NULL
	" 2>/dev/null)
	
	if [ "$orphan_users" = "0" ]; then
		log_success "  ✓ No orphaned users (users.domain_id)"
	else
		log_error "  ✗ Orphaned users found: $orphan_users"
	fi
	
	if [ "$DRY_RUN" = "true" ]; then
		echo "[DRY RUN] Would check more foreign keys"
	fi
	
	return 0
}

function check_data_types() {
	log_info "Checking data type conversions..."
	
	if [ "$DRY_RUN" = "true" ]; then
		echo "[DRY RUN] Would check data types"
	fi
	
	# Check for issues with boolean conversions (SQLite 0/1 → PostgreSQL BOOLEAN)
	local bool_issues=$(psql -h "$PG_HOST" -p "$PG_PORT" -d "$PG_DATABASE" -U "$PG_USER" -t -c "
		SELECT COUNT(*) FROM users WHERE totp_enabled NOT IN (TRUE, FALSE)
	" 2>/dev/null)
	
	if [ "$bool_issues" = "0" ]; then
		log_success "  ✓ Boolean values valid"
	else
		log_error "  ✗ Boolean conversion issues found: $bool_issues"
	fi
	
	return 0
}

function check_indexes() {
	local pg_dsn="postgresql://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/$PG_DATABASE"
	
	log_info "Checking indexes..."
	
	if [ "$DRY_RUN" = "true" ]; then
		echo "[DRY RUN] Would check indexes"
	fi
	
	# Check if critical indexes exist
	local expected_indexes=(
		"idx_users_email ON users(email)"
		"idx_users_domain_id ON users(domain_id)"
		"idx_messages_user_id ON messages(user_id)"
		"idx_messages_mailbox_id ON messages(mailbox_id)"
		"idx_messages_subject ON messages(subject)"
		"idx_smtp_queue_status ON smtp_queue(status)"
	)
	
	for index in "${expected_indexes[@]}"; do
		local exists=$(psql -h "$PG_HOST" -p "$PG_PORT" -d "$PG_DATABASE" -U "$PG_USER" -t -c "
			SELECT COUNT(*) FROM pg_indexes WHERE indexname = '$index'
		" 2>/dev/null)
		
		if [ "$exists" = "1" ]; then
			log_success "  ✓ $index"
		else
			log_warning "  ⚠ Missing index: $index"
		fi
	done
	
	return 0
}

function check_null_values() {
	local pg_dsn="postgresql://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/$PG_DATABASE"
	
	log_info "Checking for NULL values..."
	
	if [ "$DRY_RUN" = "true" ]; then
		echo "[DRY RUN] Would check NULL values"
	fi
	
	# Check for critical NULLs that shouldn't exist
	local issues=0
	
	# Users shouldn't have NULL email
	local null_emails=$(psql -h "$PG_HOST" -p "$PG_PORT" -d "$PG_DATABASE" -U "$PG_USER" -t -c "
		SELECT COUNT(*) FROM users WHERE email IS NULL
	" 2>/dev/null)
	if [ "$null_emails" -gt 0 ]; then
		log_error "  ✗ NULL emails in users table: $null_emails"
		issues=$((issues + 1))
	fi
	
	# Messages should have valid mailbox_id
	local null_mailbox=$(psql -h "$PG_HOST" -p "$PG_PORT" -d "$PG_DATABASE" -U "$PG_USER" -t -c "
		SELECT COUNT(*) FROM messages WHERE mailbox_id IS NULL
	" 2>/dev/null)
	if [ "$null_mailbox" -gt 0 ]; then
		log_error "  ✗ NULL mailbox_id in messages table: $null_mailbox"
		issues=$((issues + 1))
	fi
	
	if [ $issues -eq 0 ]; then
		log_success "No NULL value issues found"
	else
		return 1
	fi
	
	return 0
}

function main() {
	local action=""
	local verbose=false
	
	while [[ $# -gt 0 ]]; do
		case "$1" in
			-h|--help)
				show_help
				exit 0
				;;
			-v|--verbose)
				verbose=true
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
			row-counts)
				action="check_row_counts"
				shift
				;;
			foreign-keys)
				action="check_foreign_keys"
				shift
				;;
			data-types)
				action="check_data_types"
				shift
				;;
			indexes)
				action="check_indexes"
				shift
				;;
			null-checks)
				action="check_null_values"
				shift
				;;
			*)
				log_error "Unknown action: $1"
				show_help
				exit 1
				;;
		esac
	done
	
	if [ "$verbose" = "true" ]; then
		export DR_VERBOSE=true
	fi
	
	case "$action" in
		check_row_counts)
			check_foreign_keys
			check_data_types
			check_indexes
			check_null_values
			;;
		*)
			log_error "No action specified"
			show_help
			exit 1
			;;
	esac
	
	log_success "=== Validation complete ==="
}
