#!/bin/bash

################################################################################
# gomailserver Control Script
################################################################################
#
# This script manages the gomailserver daemon, providing start, stop, restart,
# and status operations. It supports both development and production modes.
#
# Usage:
#   ./gomailserver-control.sh start [--dev]     - Start the server
#   ./gomailserver-control.sh stop              - Stop the server
#   ./gomailserver-control.sh restart [--dev]   - Restart the server
#   ./gomailserver-control.sh status            - Check server status
#
# Modes:
#   --dev         Development mode (debug logging, auto-reload, local config)
#   (default)     Production mode (info logging, optimized, system config)
#
# Configuration:
#   Development:  Uses ./gomailserver.yaml for local testing
#   Production:   Uses /etc/gomailserver/gomailserver.yaml for deployment
#
################################################################################

set -e

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Binary location
BINARY_NAME="gomailserver"
BINARY_PATH="$PROJECT_ROOT/build/$BINARY_NAME"

# PID file location
PID_DIR="$PROJECT_ROOT/data"
PID_FILE="$PID_DIR/gomailserver.pid"

# Log file location
LOG_DIR="$PROJECT_ROOT/data"
LOG_FILE="$LOG_DIR/gomailserver.log"

# Configuration files
DEV_CONFIG="$PROJECT_ROOT/gomailserver.yaml"
PROD_CONFIG="/etc/gomailserver/gomailserver.yaml"

# Default mode is production
MODE="production"

################################################################################
# Color output for better readability
################################################################################
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored message
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

################################################################################
# Parse command line arguments
################################################################################
parse_args() {
    # Check for --dev flag
    if [[ " $* " =~ " --dev " ]] || [[ " $* " =~ " -d " ]]; then
        MODE="development"
        log_info "Mode: Development"
    else
        log_info "Mode: Production"
    fi
}

################################################################################
# Check if the server is running
################################################################################
# Returns: 0 if running, 1 if not running
is_running() {
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0  # Running
        else
            # PID file exists but process is not running
            log_warning "Stale PID file found, removing..."
            rm -f "$PID_FILE"
            return 1  # Not running
        fi
    fi
    return 1  # Not running
}

################################################################################
# Build the server binary
################################################################################
build_server() {
    log_info "Building gomailserver..."

    # Navigate to project root
    cd "$PROJECT_ROOT"

    # Build using make
    if make build; then
        log_success "Build completed successfully"
        return 0
    else
        log_error "Build failed"
        return 1
    fi
}

################################################################################
# Start the server
################################################################################
start_server() {
    # Check if already running
    if is_running; then
        log_warning "gomailserver is already running (PID: $(cat "$PID_FILE"))"
        return 1
    fi

    # Ensure binary exists, build if necessary
    if [ ! -f "$BINARY_PATH" ]; then
        log_warning "Binary not found, building..."
        if ! build_server; then
            log_error "Cannot start server, build failed"
            return 1
        fi
    fi

    # Create data directory if it doesn't exist
    mkdir -p "$PID_DIR"
    mkdir -p "$LOG_DIR"

    # Determine configuration file based on mode
    local config_file
    if [ "$MODE" = "development" ]; then
        config_file="$DEV_CONFIG"

        # Create development config if it doesn't exist
        if [ ! -f "$config_file" ]; then
            log_warning "Development config not found, copying from example..."
            if [ -f "$PROJECT_ROOT/gomailserver.example.yaml" ]; then
                cp "$PROJECT_ROOT/gomailserver.example.yaml" "$config_file"
                log_info "Created development config at $config_file"
            else
                log_error "Example config not found, cannot create development config"
                return 1
            fi
        fi
    else
        config_file="$PROD_CONFIG"

        # Ensure production config exists
        if [ ! -f "$config_file" ]; then
            log_error "Production config not found at $config_file"
            log_info "Please install configuration or use --dev flag for development mode"
            return 1
        fi
    fi

    log_info "Using configuration: $config_file"

    # Start the server in the background
    log_info "Starting gomailserver..."

    if [ "$MODE" = "development" ]; then
        # Development mode: more verbose logging, log to file
        "$BINARY_PATH" run --config "$config_file" > "$LOG_FILE" 2>&1 &
        local pid=$!
    else
        # Production mode: standard logging
        "$BINARY_PATH" run --config "$config_file" >> "$LOG_FILE" 2>&1 &
        local pid=$!
    fi

    # Save PID to file
    echo "$pid" > "$PID_FILE"

    # Wait a moment and check if the process is still running
    sleep 2
    if ps -p "$pid" > /dev/null 2>&1; then
        log_success "gomailserver started successfully (PID: $pid)"
        log_info "Logs: $LOG_FILE"
        return 0
    else
        log_error "gomailserver failed to start, check logs at $LOG_FILE"
        rm -f "$PID_FILE"
        return 1
    fi
}

################################################################################
# Stop the server
################################################################################
stop_server() {
    # Check if running
    if ! is_running; then
        log_warning "gomailserver is not running"
        return 1
    fi

    local pid=$(cat "$PID_FILE")
    log_info "Stopping gomailserver (PID: $pid)..."

    # Send SIGTERM for graceful shutdown
    kill -TERM "$pid" 2>/dev/null

    # Wait for process to terminate (max 30 seconds)
    local timeout=30
    local count=0
    while ps -p "$pid" > /dev/null 2>&1; do
        if [ $count -ge $timeout ]; then
            log_warning "Process did not terminate gracefully, forcing shutdown..."
            kill -KILL "$pid" 2>/dev/null
            break
        fi
        sleep 1
        count=$((count + 1))
    done

    # Remove PID file
    rm -f "$PID_FILE"

    log_success "gomailserver stopped"
    return 0
}

################################################################################
# Restart the server
################################################################################
restart_server() {
    log_info "Restarting gomailserver..."

    # Stop if running
    if is_running; then
        stop_server
    fi

    # Small delay to ensure clean shutdown
    sleep 1

    # Start server
    start_server
}

################################################################################
# Show server status
################################################################################
show_status() {
    if is_running; then
        local pid=$(cat "$PID_FILE")
        log_success "gomailserver is running (PID: $pid)"

        # Show process information
        echo ""
        ps -p "$pid" -o pid,ppid,user,%cpu,%mem,etime,cmd

        # Show listening ports if netstat/ss is available
        echo ""
        if command -v ss > /dev/null 2>&1; then
            log_info "Listening ports:"
            ss -tlnp | grep "$pid" || log_warning "No listening ports found"
        elif command -v netstat > /dev/null 2>&1; then
            log_info "Listening ports:"
            netstat -tlnp 2>/dev/null | grep "$pid" || log_warning "No listening ports found"
        fi

        return 0
    else
        log_warning "gomailserver is not running"
        return 1
    fi
}

################################################################################
# Show usage information
################################################################################
show_usage() {
    cat << EOF
Usage: $0 <command> [options]

Commands:
  start [--dev]     Start the gomailserver daemon
  stop              Stop the gomailserver daemon
  restart [--dev]   Restart the gomailserver daemon
  status            Show server status
  build             Build the server binary
  help              Show this help message

Options:
  --dev, -d         Run in development mode (debug logging, local config)

Examples:
  $0 start          Start in production mode
  $0 start --dev    Start in development mode
  $0 stop           Stop the server
  $0 restart --dev  Restart in development mode
  $0 status         Check if server is running

Configuration:
  Development:  $DEV_CONFIG
  Production:   $PROD_CONFIG

Logs:
  $LOG_FILE

PID File:
  $PID_FILE

EOF
}

################################################################################
# Main script logic
################################################################################
main() {
    # Check if no arguments provided
    if [ $# -eq 0 ]; then
        show_usage
        exit 1
    fi

    # Parse command
    local command="$1"
    shift

    # Parse remaining arguments for flags
    parse_args "$@"

    # Execute command
    case "$command" in
        start)
            start_server
            exit $?
            ;;
        stop)
            stop_server
            exit $?
            ;;
        restart)
            restart_server
            exit $?
            ;;
        status)
            show_status
            exit $?
            ;;
        build)
            build_server
            exit $?
            ;;
        help|--help|-h)
            show_usage
            exit 0
            ;;
        *)
            log_error "Unknown command: $command"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
