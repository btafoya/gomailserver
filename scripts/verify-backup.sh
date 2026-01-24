#!/bin/bash
# gomailserver Backup Verification Script
# This script verifies the integrity of gomailserver backups

set -euo pipefail

BACKUP_DIR="${1:-/var/lib/gomailserver/backups}"
LOG_FILE="/var/log/gomailserver/backup-verify.log"

# Logging function
log() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] [$level] $message" >> "$LOG_FILE"
    echo "[$timestamp] $level: $message"
}

info() { log "INFO" "$1"; }
error() { log "ERROR" "$1"; exit 1; }
warn() { log "WARN" "$1"; }

# Verify database backup
verify_database() {
    local db_file="$1"

    if [[ ! -f "$db_file" ]]; then
        error "Database backup file not found: $db_file"
    fi

    info "Verifying database backup: $db_file"

    # Check SQLite integrity
    if ! sqlite3 "$db_file" "PRAGMA integrity_check;" | grep -q "ok"; then
        error "Database integrity check failed for: $db_file"
    fi

    # Check if we can query basic tables
    local table_count=$(sqlite3 "$db_file" "SELECT COUNT(*) FROM sqlite_master WHERE type='table';" 2>/dev/null || echo "0")
    if [[ "$table_count" -lt 5 ]]; then
        error "Database appears incomplete, only $table_count tables found in: $db_file"
    fi

    info "Database verification passed: $db_file"
}

# Verify tar archive
verify_archive() {
    local archive_file="$1"
    local archive_type="$2"

    if [[ ! -f "$archive_file" ]]; then
        error "$archive_type backup file not found: $archive_file"
    fi

    info "Verifying $archive_type backup: $archive_file"

    # Test archive integrity
    if ! tar -tzf "$archive_file" >/dev/null 2>&1; then
        error "$archive_type archive is corrupted: $archive_file"
    fi

    # Check if archive is not empty
    local file_count=$(tar -tzf "$archive_file" | wc -l)
    if [[ "$file_count" -eq 0 ]]; then
        error "$archive_type archive is empty: $archive_file"
    fi

    info "$archive_type verification passed: $archive_file ($file_count files)"
}

# Find latest backup set
find_latest_backup() {
    local backup_prefix="$1"
    local latest_file=""

    # Find the most recent file matching the pattern
    latest_file=$(find "$BACKUP_DIR" -name "${backup_prefix}_*" -type f -printf '%T@ %p\n' 2>/dev/null | sort -n | tail -1 | cut -d' ' -f2-)

    if [[ -z "$latest_file" ]]; then
        error "No $backup_prefix backup files found in $BACKUP_DIR"
    fi

    echo "$latest_file"
}

# Main verification function
main() {
    info "Starting gomailserver backup verification"

    # Create log directory if needed
    local log_dir=$(dirname "$LOG_FILE")
    mkdir -p "$log_dir" 2>/dev/null || true

    # Find latest backup files
    local db_backup=$(find_latest_backup "*_db_*.db")
    local messages_backup=$(find_latest_backup "*_messages_*.tar.gz")
    local dkim_backup=$(find_latest_backup "*_dkim_*.tar.gz")
    local config_backup=$(find_latest_backup "*_config_*.tar.gz")

    # Verify each backup type
    verify_database "$db_backup"

    if [[ -f "$messages_backup" ]]; then
        verify_archive "$messages_backup" "messages"
    else
        warn "No messages backup found to verify"
    fi

    if [[ -f "$dkim_backup" ]]; then
        verify_archive "$dkim_backup" "DKIM"
    else
        warn "No DKIM backup found to verify"
    fi

    if [[ -f "$config_backup" ]]; then
        verify_archive "$config_backup" "config"
    fi

    # Check backup freshness (should be less than 48 hours old)
    local latest_timestamp=""
    for file in "$db_backup" "$messages_backup" "$dkim_backup" "$config_backup"; do
        if [[ -f "$file" ]]; then
            local file_timestamp=$(stat -c %Y "$file" 2>/dev/null || stat -f %m "$file" 2>/dev/null || echo "0")
            if [[ "$file_timestamp" -gt "$latest_timestamp" ]]; then
                latest_timestamp="$file_timestamp"
            fi
        fi
    done

    local current_time=$(date +%s)
    local age_hours=$(( (current_time - latest_timestamp) / 3600 ))

    if [[ $age_hours -gt 48 ]]; then
        warn "Latest backup is $age_hours hours old (should be less than 48 hours)"
    else
        info "Backup freshness OK: $age_hours hours old"
    fi

    # Check disk space for backups
    local available_mb=$(df -BM "$BACKUP_DIR" | tail -1 | awk '{print $4}' | sed 's/M//')
    if [[ $available_mb -lt 1024 ]]; then  # Less than 1GB
        warn "Low disk space for backups: ${available_mb}MB available"
    fi

    info "Backup verification completed successfully"
}

# Run verification
main "$@"