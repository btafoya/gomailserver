#!/bin/bash
# gomailserver systemd Management Script
# Provides monitoring, health checks, and service management utilities

set -euo pipefail

SCRIPT_VERSION="1.0"
SERVICE_NAME="gomailserver"
CONFIG_FILE="/etc/gomailserver/gomailserver.yaml"
DATA_DIR="/var/lib/gomailserver"
LOG_DIR="/var/log/gomailserver"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Logging functions
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

# Service status
status() {
    log_info "Checking $SERVICE_NAME status..."

    if ! systemctl is-active --quiet "$SERVICE_NAME"; then
        log_error "$SERVICE_NAME is not running"
        echo
        log_info "Recent logs:"
        journalctl -u "$SERVICE_NAME" -n 10 --no-pager
        return 1
    fi

    log_success "$SERVICE_NAME is running"
    systemctl status "$SERVICE_NAME" --no-pager -l
}

# Health check
health_check() {
    log_info "Performing health check..."

    # Check if service is running
    if ! systemctl is-active --quiet "$SERVICE_NAME"; then
        log_error "Service is not running"
        return 1
    fi

    # Check database
    if [[ ! -f "$DATA_DIR/mailserver.db" ]]; then
        log_error "Database file not found"
        return 1
    fi

    # Test database integrity
    if ! timeout 10 sqlite3 "$DATA_DIR/mailserver.db" "PRAGMA integrity_check;" | grep -q "ok"; then
        log_error "Database integrity check failed"
        return 1
    fi

    # Check configuration
    if [[ ! -f "$CONFIG_FILE" ]]; then
        log_warning "Configuration file not found"
    fi

    # Check web interface
    if command -v curl >/dev/null 2>&1; then
        if ! curl -f -s --max-time 5 http://localhost:8980/api/v1/health >/dev/null 2>&1; then
            log_warning "Web interface health check failed"
        else
            log_success "Web interface is responding"
        fi
    fi

    # Check resource usage
    local pid=$(systemctl show -p MainPID --value "$SERVICE_NAME")
    if [[ -n "$pid" && "$pid" != "0" ]]; then
        local mem_usage=$(ps -o pmem= -p "$pid" | tr -d ' ')
        local cpu_usage=$(ps -o pcpu= -p "$pid" | tr -d ' ')

        if [[ $(echo "$mem_usage > 80" | bc -l 2>/dev/null) -eq 1 ]]; then
            log_warning "High memory usage: ${mem_usage}%"
        else
            log_info "Memory usage: ${mem_usage}%"
        fi

        log_info "CPU usage: ${cpu_usage}%"
    fi

    # Check disk space
    local db_disk_usage=$(df "$DATA_DIR" | tail -1 | awk '{print $5}' | sed 's/%//')
    if [[ $db_disk_usage -gt 90 ]]; then
        log_error "Low disk space for database: ${db_disk_usage}% used"
        return 1
    elif [[ $db_disk_usage -gt 80 ]]; then
        log_warning "High disk usage for database: ${db_disk_usage}% used"
    fi

    log_success "Health check passed"
}

# Performance monitoring
performance() {
    log_info "Performance monitoring..."

    # Service resource usage
    echo "=== Service Resource Usage ==="
    systemctl status "$SERVICE_NAME" --no-pager | grep -E "(Memory|CPU|Tasks)"

    # Database statistics
    echo
    echo "=== Database Statistics ==="
    if [[ -f "$DATA_DIR/mailserver.db" ]]; then
        echo "Database size: $(du -h "$DATA_DIR/mailserver.db" | cut -f1)"
        echo "Table counts:"
        sqlite3 "$DATA_DIR/mailserver.db" << 'EOF' | column -t -s '|'
.mode list
.separator |
SELECT name, COUNT(*) as count FROM (
    SELECT 'users' as name FROM users
    UNION ALL SELECT 'domains' FROM domains
    UNION ALL SELECT 'mailboxes' FROM mailboxes
    UNION ALL SELECT 'messages' FROM messages
    UNION ALL SELECT 'queue' FROM queue
) GROUP BY name;
EOF
    fi

    # Queue status
    echo
    echo "=== Queue Status ==="
    if [[ -f "$DATA_DIR/mailserver.db" ]]; then
        sqlite3 "$DATA_DIR/mailserver.db" << 'EOF'
SELECT status, COUNT(*) as count
FROM queue
GROUP BY status;
EOF
    fi

    # Network connections
    echo
    echo "=== Network Connections ==="
    ss -tlnp | grep -E ":(25|587|465|143|993|8980) " || true
}

# Log analysis
logs() {
    local lines="${1:-50}"
    local follow="${2:-false}"

    if [[ "$follow" == "true" ]]; then
        log_info "Following logs (Ctrl+C to stop)..."
        journalctl -u "$SERVICE_NAME" -f
    else
        log_info "Showing last $lines log entries..."
        journalctl -u "$SERVICE_NAME" -n "$lines" --no-pager
    fi
}

# Backup management
backup() {
    local backup_script="/opt/gomailserver/scripts/backup.sh"

    if [[ ! -x "$backup_script" ]]; then
        log_error "Backup script not found at $backup_script"
        log_info "Install gomailserver scripts or run backup manually"
        return 1
    fi

    log_info "Running automated backup..."
    "$backup_script"
}

# Configuration validation
validate_config() {
    log_info "Validating configuration..."

    if [[ ! -f "$CONFIG_FILE" ]]; then
        log_error "Configuration file not found: $CONFIG_FILE"
        return 1
    fi

    # Basic YAML syntax check
    if command -v python3 >/dev/null 2>&1; then
        if ! python3 -c "import yaml; yaml.safe_load(open('$CONFIG_FILE'))" 2>/dev/null; then
            log_error "Configuration file has invalid YAML syntax"
            return 1
        fi
    fi

    # Check required configuration sections
    local required_sections=("server" "database")
    for section in "${required_sections[@]}"; do
        if ! grep -q "^${section}:" "$CONFIG_FILE"; then
            log_error "Missing required configuration section: $section"
            return 1
        fi
    done

    log_success "Configuration validation passed"
}

# System information
system_info() {
    echo "=== System Information ==="
    echo "OS: $(lsb_release -d 2>/dev/null | cut -f2 || uname -s)"
    echo "Kernel: $(uname -r)"
    echo "Architecture: $(uname -m)"
    echo "Uptime: $(uptime -p)"

    echo
    echo "=== gomailserver Information ==="
    if command -v gomailserver >/dev/null 2>&1; then
        echo "Version: $(gomailserver version 2>/dev/null || echo "Unknown")"
    fi
    echo "Service Status: $(systemctl is-active "$SERVICE_NAME" 2>/dev/null || echo "Not installed")"
    echo "Config File: $CONFIG_FILE ($(test -f "$CONFIG_FILE" && echo "Present" || echo "Missing"))"
    echo "Database: $DATA_DIR/mailserver.db ($(test -f "$DATA_DIR/mailserver.db" && echo "Present" || echo "Missing"))"

    echo
    echo "=== Resource Usage ==="
    echo "Disk Usage (Data): $(df -h "$DATA_DIR" | tail -1 | awk '{print $5 " used (" $3 "/" $2 ")"}')"
    echo "Memory: $(free -h | grep "^Mem:" | awk '{print $3 "/" $2 " used"}')"
}

# Service restart with health verification
restart() {
    log_info "Restarting $SERVICE_NAME with health verification..."

    # Stop service
    systemctl stop "$SERVICE_NAME"

    # Wait for clean shutdown
    local count=0
    while systemctl is-active --quiet "$SERVICE_NAME" && [[ $count -lt 30 ]]; do
        sleep 1
        ((count++))
    done

    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_warning "Service did not stop cleanly, forcing stop..."
        systemctl kill "$SERVICE_NAME"
        sleep 5
    fi

    # Start service
    systemctl start "$SERVICE_NAME"

    # Wait for startup and verify health
    count=0
    while [[ $count -lt 30 ]]; do
        if systemctl is-active --quiet "$SERVICE_NAME"; then
            sleep 2  # Give it time to fully start
            if health_check >/dev/null 2>&1; then
                log_success "Service restarted successfully"
                return 0
            fi
        fi
        sleep 1
        ((count++))
    done

    log_error "Service failed to restart properly"
    logs 20
    return 1
}

# Show usage
usage() {
    cat << EOF
gomailserver systemd Management Script v$SCRIPT_VERSION

USAGE:
    $0 <command> [options]

COMMANDS:
    status          Show service status and recent logs
    health          Perform comprehensive health check
    performance     Show performance and resource usage statistics
    logs [lines]    Show recent log entries (default: 50)
    logs follow     Follow logs in real-time
    backup          Run automated backup
    validate        Validate configuration file
    info            Show system and service information
    restart         Restart service with health verification

EXAMPLES:
    $0 status                    # Show service status
    $0 health                    # Run health checks
    $0 performance               # Show performance stats
    $0 logs 100                  # Show last 100 log lines
    $0 logs follow               # Follow logs live
    $0 backup                    # Run backup
    $0 validate                  # Check configuration
    $0 info                      # Show system info
    $0 restart                   # Restart with verification

For systemd service management, use:
    sudo systemctl start/stop/restart/status gomailserver
    sudo journalctl -u gomailserver -f

EOF
}

# Main function
main() {
    local command="${1:-help}"

    case "$command" in
        status)
            check_root
            status
            ;;
        health)
            check_root
            health_check
            ;;
        performance)
            check_root
            performance
            ;;
        logs)
            check_root
            if [[ "${2:-}" == "follow" ]]; then
                logs "" true
            else
                logs "${2:-50}"
            fi
            ;;
        backup)
            check_root
            backup
            ;;
        validate)
            validate_config
            ;;
        info)
            system_info
            ;;
        restart)
            check_root
            restart
            ;;
        help|--help|-h)
            usage
            ;;
        *)
            log_error "Unknown command: $command"
            echo
            usage
            exit 1
            ;;
    esac
}

main "$@"