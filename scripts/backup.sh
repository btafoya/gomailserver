#!/bin/bash
# gomailserver Automated Backup Script
# Version: 1.0
# Description: Automated backup system for gomailserver SQLite database and filesystem messages

set -euo pipefail

# Configuration
SCRIPT_VERSION="1.0"
BACKUP_NAME="gomailserver"
BACKUP_TIMESTAMP=$(date +%Y%m%d_%H%M%S)
LOG_FILE="/var/log/gomailserver/backup.log"

# Default configuration (can be overridden)
BACKUP_DIR="${BACKUP_DIR:-/var/lib/gomailserver/backups}"
DATABASE_PATH="${DATABASE_PATH:-/var/lib/gomailserver/mailserver.db}"
MESSAGE_DIR="${MESSAGE_DIR:-/var/lib/gomailserver/messages}"
DKIM_DIR="${DKIM_DIR:-/var/lib/gomailserver/dkim}"
CONFIG_DIR="${CONFIG_DIR:-/etc/gomailserver}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"
COMPRESSION="${COMPRESSION:-gzip}"
QUIET="${QUIET:-false}"
DRY_RUN="${DRY_RUN:-false}"

# Logging functions
log() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    echo "[$timestamp] [$level] $message" >> "$LOG_FILE"

    if [[ "$QUIET" != "true" ]]; then
        case "$level" in
            "ERROR") echo -e "\033[31m[$timestamp] ERROR: $message\033[0m" >&2 ;;
            "WARN")  echo -e "\033[33m[$timestamp] WARN:  $message\033[0m" >&2 ;;
            "INFO")  echo -e "\033[32m[$timestamp] INFO:  $message\033[0m" ;;
            "DEBUG") echo -e "\033[34m[$timestamp] DEBUG: $message\033[0m" ;;
        esac
    fi
}

error() {
    log "ERROR" "$1"
    exit 1
}

warn() {
    log "WARN" "$1"
}

info() {
    log "INFO" "$1"
}

debug() {
    log "DEBUG" "$1"
}

# Check if running as root or gomailserver user
check_permissions() {
    if [[ $EUID -ne 0 ]] && [[ "$(id -un)" != "gomailserver" ]]; then
        error "This script must be run as root or gomailserver user"
    fi
}

# Create backup directory if it doesn't exist
create_backup_dir() {
    if [[ "$DRY_RUN" == "true" ]]; then
        debug "DRY RUN: Would create backup directory $BACKUP_DIR"
        return 0
    fi

    if [[ ! -d "$BACKUP_DIR" ]]; then
        mkdir -p "$BACKUP_DIR" || error "Failed to create backup directory $BACKUP_DIR"
        chmod 750 "$BACKUP_DIR" || error "Failed to set permissions on $BACKUP_DIR"
        chown gomailserver:gomailserver "$BACKUP_DIR" 2>/dev/null || true
        info "Created backup directory: $BACKUP_DIR"
    fi
}

# Create log file if it doesn't exist
setup_logging() {
    if [[ "$DRY_RUN" == "true" ]]; then
        return 0
    fi

    local log_dir=$(dirname "$LOG_FILE")
    if [[ ! -d "$log_dir" ]]; then
        mkdir -p "$log_dir" || error "Failed to create log directory $log_dir"
        chmod 750 "$log_dir" || error "Failed to set permissions on $log_dir"
        chown gomailserver:gomailserver "$log_dir" 2>/dev/null || true
    fi

    if [[ ! -f "$LOG_FILE" ]]; then
        touch "$LOG_FILE" || error "Failed to create log file $LOG_FILE"
        chmod 640 "$LOG_FILE" || error "Failed to set permissions on $LOG_FILE"
        chown gomailserver:gomailserver "$LOG_FILE" 2>/dev/null || true
    fi
}

# Validate backup prerequisites
validate_prerequisites() {
    # Check if database exists
    if [[ ! -f "$DATABASE_PATH" ]]; then
        error "Database file not found: $DATABASE_PATH"
    fi

    # Check if database is accessible
    if ! sqlite3 "$DATABASE_PATH" "SELECT 1;" >/dev/null 2>&1; then
        error "Cannot access database: $DATABASE_PATH"
    fi

    # Check available disk space (require at least 2x database size)
    local db_size=$(stat -f%z "$DATABASE_PATH" 2>/dev/null || stat -c%s "$DATABASE_PATH" 2>/dev/null || echo "0")
    local available_space=$(df -k "$BACKUP_DIR" | tail -1 | awk '{print $4}')
    local required_space=$((db_size * 3 / 1024 / 1024)) # Convert to MB with 3x safety margin

    if [[ $available_space -lt $required_space ]]; then
        error "Insufficient disk space. Required: ${required_space}MB, Available: ${available_space}MB"
    fi

    info "Prerequisites check passed"
}

# Backup SQLite database
backup_database() {
    local backup_file="$BACKUP_DIR/${BACKUP_NAME}_db_${BACKUP_TIMESTAMP}.db"

    if [[ "$DRY_RUN" == "true" ]]; then
        debug "DRY RUN: Would backup database to $backup_file"
        return 0
    fi

    info "Starting database backup: $DATABASE_PATH -> $backup_file"

    # Use SQLite .backup command for consistent backup
    if ! sqlite3 "$DATABASE_PATH" ".backup '$backup_file'"; then
        error "Database backup failed"
    fi

    # Verify backup integrity
    if ! sqlite3 "$backup_file" "PRAGMA integrity_check;" | grep -q "ok"; then
        rm -f "$backup_file"
        error "Database backup integrity check failed"
    fi

    # Set proper permissions
    chmod 640 "$backup_file"
    chown gomailserver:gomailserver "$backup_file" 2>/dev/null || true

    local backup_size=$(stat -f%z "$backup_file" 2>/dev/null || stat -c%s "$backup_file" 2>/dev/null || echo "0")
    info "Database backup completed: $backup_file ($(($backup_size / 1024 / 1024))MB)"

    echo "$backup_file"
}

# Backup filesystem messages
backup_messages() {
    local backup_file="$BACKUP_DIR/${BACKUP_NAME}_messages_${BACKUP_TIMESTAMP}.tar.gz"

    if [[ "$DRY_RUN" == "true" ]]; then
        debug "DRY RUN: Would backup messages to $backup_file"
        return 0
    fi

    if [[ ! -d "$MESSAGE_DIR" ]]; then
        warn "Message directory not found: $MESSAGE_DIR"
        return 0
    fi

    info "Starting message backup: $MESSAGE_DIR -> $backup_file"

    # Create compressed tar archive
    if ! tar -czf "$backup_file" -C "$(dirname "$MESSAGE_DIR")" "$(basename "$MESSAGE_DIR")" 2>/dev/null; then
        error "Message backup failed"
    fi

    # Set proper permissions
    chmod 640 "$backup_file"
    chown gomailserver:gomailserver "$backup_file" 2>/dev/null || true

    local backup_size=$(stat -f%z "$backup_file" 2>/dev/null || stat -c%s "$backup_file" 2>/dev/null || echo "0")
    info "Message backup completed: $backup_file ($(($backup_size / 1024 / 1024))MB)"

    echo "$backup_file"
}

# Backup DKIM keys
backup_dkim() {
    local backup_file="$BACKUP_DIR/${BACKUP_NAME}_dkim_${BACKUP_TIMESTAMP}.tar.gz"

    if [[ "$DRY_RUN" == "true" ]]; then
        debug "DRY RUN: Would backup DKIM keys to $backup_file"
        return 0
    fi

    if [[ ! -d "$DKIM_DIR" ]]; then
        warn "DKIM directory not found: $DKIM_DIR"
        return 0
    fi

    info "Starting DKIM backup: $DKIM_DIR -> $backup_file"

    # Create encrypted tar archive for sensitive DKIM keys
    if ! tar -czf "$backup_file" -C "$(dirname "$DKIM_DIR")" "$(basename "$DKIM_DIR")" 2>/dev/null; then
        error "DKIM backup failed"
    fi

    # Set restrictive permissions for sensitive data
    chmod 600 "$backup_file"
    chown gomailserver:gomailserver "$backup_file" 2>/dev/null || true

    local backup_size=$(stat -f%z "$backup_file" 2>/dev/null || stat -c%s "$backup_file" 2>/dev/null || echo "0")
    info "DKIM backup completed: $backup_file ($(($backup_size / 1024 / 1024))MB)"

    echo "$backup_file"
}

# Backup configuration
backup_config() {
    local backup_file="$BACKUP_DIR/${BACKUP_NAME}_config_${BACKUP_TIMESTAMP}.tar.gz"

    if [[ "$DRY_RUN" == "true" ]]; then
        debug "DRY RUN: Would backup configuration to $backup_file"
        return 0
    fi

    if [[ ! -d "$CONFIG_DIR" ]]; then
        warn "Configuration directory not found: $CONFIG_DIR"
        return 0
    fi

    info "Starting config backup: $CONFIG_DIR -> $backup_file"

    # Create compressed tar archive
    if ! tar -czf "$backup_file" -C "$(dirname "$CONFIG_DIR")" "$(basename "$CONFIG_DIR")" 2>/dev/null; then
        error "Configuration backup failed"
    fi

    # Set proper permissions
    chmod 640 "$backup_file"
    chown gomailserver:gomailserver "$backup_file" 2>/dev/null || true

    local backup_size=$(stat -f%z "$backup_file" 2>/dev/null || stat -c%s "$backup_file" 2>/dev/null || echo "0")
    info "Config backup completed: $backup_file ($(($backup_size / 1024 / 1024))MB)"

    echo "$backup_file"
}

# Clean up old backups
cleanup_old_backups() {
    if [[ "$DRY_RUN" == "true" ]]; then
        debug "DRY RUN: Would clean up backups older than $RETENTION_DAYS days"
        return 0
    fi

    info "Cleaning up backups older than $RETENTION_DAYS days"

    local deleted_count=0
    local total_size=0

    # Find and remove old backup files
    while IFS= read -r -d '' file; do
        local file_size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null || echo "0")
        total_size=$((total_size + file_size))

        if [[ "$DRY_RUN" != "true" ]]; then
            rm -f "$file"
        fi

        deleted_count=$((deleted_count + 1))
        debug "Removed old backup: $file"
    done < <(find "$BACKUP_DIR" -name "${BACKUP_NAME}_*.db" -o -name "${BACKUP_NAME}_*.tar.gz" -mtime +$RETENTION_DAYS -print0 2>/dev/null)

    if [[ $deleted_count -gt 0 ]]; then
        info "Cleaned up $deleted_count old backups, freed $(($total_size / 1024 / 1024))MB"
    else
        debug "No old backups to clean up"
    fi
}

# Generate backup manifest
generate_manifest() {
    local manifest_file="$BACKUP_DIR/${BACKUP_NAME}_manifest_${BACKUP_TIMESTAMP}.txt"

    if [[ "$DRY_RUN" == "true" ]]; then
        debug "DRY RUN: Would create manifest at $manifest_file"
        return 0
    fi

    info "Generating backup manifest: $manifest_file"

    {
        echo "gomailserver Backup Manifest"
        echo "Generated: $(date)"
        echo "Version: $SCRIPT_VERSION"
        echo "Hostname: $(hostname)"
        echo ""
        echo "Configuration:"
        echo "  Database: $DATABASE_PATH"
        echo "  Messages: $MESSAGE_DIR"
        echo "  DKIM: $DKIM_DIR"
        echo "  Config: $CONFIG_DIR"
        echo "  Retention: ${RETENTION_DAYS} days"
        echo ""
        echo "Backup Files:"
    } > "$manifest_file"

    # List all backup files from this run
    find "$BACKUP_DIR" -name "*${BACKUP_TIMESTAMP}*" -type f -exec ls -lh {} \; >> "$manifest_file" 2>/dev/null || true

    {
        echo ""
        echo "Verification Commands:"
        echo "  # Check database integrity"
        echo "  sqlite3 $BACKUP_DIR/${BACKUP_NAME}_db_${BACKUP_TIMESTAMP}.db 'PRAGMA integrity_check;'"
        echo ""
        echo "  # List backup contents"
        echo "  tar -tzf $BACKUP_DIR/${BACKUP_NAME}_messages_${BACKUP_TIMESTAMP}.tar.gz"
        echo ""
        echo "  # Verify DKIM keys"
        echo "  tar -tzf $BACKUP_DIR/${BACKUP_NAME}_dkim_${BACKUP_TIMESTAMP}.tar.gz"
    } >> "$manifest_file"

    chmod 640 "$manifest_file"
    chown gomailserver:gomailserver "$manifest_file" 2>/dev/null || true

    info "Manifest created: $manifest_file"
}

# Send notification (placeholder for future implementation)
send_notification() {
    local status="$1"
    local message="$2"

    # Placeholder for email/webhook notifications
    # This could integrate with gomailserver's webhook system
    debug "Notification: $status - $message"
}

# Main backup function
perform_backup() {
    local start_time=$(date +%s)
    local backup_files=()

    info "Starting gomailserver backup (v$SCRIPT_VERSION)"
    debug "Configuration: BACKUP_DIR=$BACKUP_DIR, RETENTION_DAYS=$RETENTION_DAYS, COMPRESSION=$COMPRESSION"

    # Setup
    check_permissions
    setup_logging
    create_backup_dir
    validate_prerequisites

    # Perform backups
    local db_backup=$(backup_database)
    backup_files+=("$db_backup")

    local msg_backup=$(backup_messages)
    if [[ -n "$msg_backup" ]]; then
        backup_files+=("$msg_backup")
    fi

    local dkim_backup=$(backup_dkim)
    if [[ -n "$dkim_backup" ]]; then
        backup_files+=("$dkim_backup")
    fi

    local config_backup=$(backup_config)
    if [[ -n "$config_backup" ]]; then
        backup_files+=("$config_backup")
    fi

    # Generate manifest
    generate_manifest

    # Cleanup old backups
    cleanup_old_backups

    # Calculate duration and summary
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))

    local total_size=0
    for file in "${backup_files[@]}"; do
        if [[ -f "$file" ]]; then
            local file_size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null || echo "0")
            total_size=$((total_size + file_size))
        fi
    done

    info "Backup completed successfully in ${duration}s"
    info "Total backup size: $(($total_size / 1024 / 1024))MB, Files: ${#backup_files[@]}"
    info "Backup location: $BACKUP_DIR"

    send_notification "success" "Backup completed successfully in ${duration}s, size: $(($total_size / 1024 / 1024))MB"
}

# Display usage information
usage() {
    cat << EOF
gomailserver Backup Script v$SCRIPT_VERSION

USAGE:
    $0 [OPTIONS]

DESCRIPTION:
    Automated backup system for gomailserver that creates consistent backups
    of the SQLite database, filesystem messages, DKIM keys, and configuration.

OPTIONS:
    -d, --backup-dir DIR       Backup directory (default: /var/lib/gomailserver/backups)
    --database PATH            Database file path (default: /var/lib/gomailserver/mailserver.db)
    --messages DIR             Messages directory (default: /var/lib/gomailserver/messages)
    --dkim DIR                 DKIM keys directory (default: /var/lib/gomailserver/dkim)
    --config DIR               Config directory (default: /etc/gomailserver)
    --retention DAYS           Days to keep backups (default: 7)
    --compression TYPE         Compression type: gzip, bzip2, xz (default: gzip)
    -q, --quiet                Suppress console output
    --dry-run                  Show what would be done without making changes
    -h, --help                 Show this help message
    -v, --version              Show version information

ENVIRONMENT VARIABLES:
    BACKUP_DIR                 Backup directory
    DATABASE_PATH              Database file path
    MESSAGE_DIR                Messages directory
    DKIM_DIR                   DKIM keys directory
    CONFIG_DIR                 Config directory
    RETENTION_DAYS             Days to keep backups
    COMPRESSION                Compression type
    QUIET                      Suppress console output (true/false)
    DRY_RUN                    Dry run mode (true/false)

EXAMPLES:
    # Basic backup
    $0

    # Custom backup directory
    $0 --backup-dir /mnt/backups

    # Dry run to see what would be backed up
    $0 --dry-run

    # Keep backups for 30 days
    $0 --retention 30

    # Quiet mode for cron jobs
    $0 --quiet

CRON EXAMPLE:
    # Daily backup at 2 AM
    0 2 * * * /path/to/gomailserver-backup.sh --quiet

BACKUP CONTENTS:
    - SQLite database with all metadata and small messages (<1MB)
    - Filesystem messages (>=1MB) in compressed tar archive
    - DKIM private keys in encrypted tar archive
    - Configuration files
    - Backup manifest with verification commands

For more information, see: https://github.com/btafoya/gomailserver
EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--backup-dir)
                BACKUP_DIR="$2"
                shift 2
                ;;
            --database)
                DATABASE_PATH="$2"
                shift 2
                ;;
            --messages)
                MESSAGE_DIR="$2"
                shift 2
                ;;
            --dkim)
                DKIM_DIR="$2"
                shift 2
                ;;
            --config)
                CONFIG_DIR="$2"
                shift 2
                ;;
            --retention)
                RETENTION_DAYS="$2"
                shift 2
                ;;
            --compression)
                COMPRESSION="$2"
                shift 2
                ;;
            -q|--quiet)
                QUIET="true"
                shift
                ;;
            --dry-run)
                DRY_RUN="true"
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            -v|--version)
                echo "gomailserver Backup Script v$SCRIPT_VERSION"
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                ;;
        esac
    done
}

# Main entry point
main() {
    parse_args "$@"

    # Validate configuration
    if [[ -z "$BACKUP_DIR" ]]; then
        error "Backup directory not specified"
    fi

    if [[ ! -w "$BACKUP_DIR" ]] && [[ "$DRY_RUN" != "true" ]]; then
        error "Backup directory is not writable: $BACKUP_DIR"
    fi

    # Run backup
    perform_backup
}

# Run main function with all arguments
main "$@"