#!/bin/bash
#
# gomailserver Control Script
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

# Exit on error, undefined variables, and pipe failures
set -euo pipefail
# Trap to cleanup on script exit
trap cleanup EXIT

# ==============================================================================
# Configuration and Path Setup
# ==============================================================================

# Get script directory and project root (resolves symlinks)
get_script_dir() {
    local source="${BASH_SOURCE[0]}"
    while [ -L "$source" ]; do
        local dir
        dir="$(cd -P "$(dirname "$source")" && pwd)"
        source="$(readlink "$source")"
        [[ $source != /* ]] && source="$dir/$source"
    done
    cd -P "$(dirname "$source")" && pwd
}

readonly SCRIPT_DIR="$(get_script_dir)"
readonly PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Binary configuration
readonly BINARY_NAME="gomailserver"
readonly BINARY_PATH="$PROJECT_ROOT/build/$BINARY_NAME"

# Directory configuration
readonly PID_DIR="$PROJECT_ROOT/data"
readonly LOG_DIR="$PROJECT_ROOT/data"
readonly PID_FILE="$PID_DIR/gomailserver.pid"
readonly WEBUI_PID_FILE="$PID_DIR/webui.pid"
readonly LOG_FILE="$LOG_DIR/gomailserver.log"
readonly WEBUI_LOG_FILE="$LOG_DIR/webui.log"

# Configuration files
readonly DEV_CONFIG="$PROJECT_ROOT/gomailserver.yaml"
readonly PROD_CONFIG="/etc/gomailserver/gomailserver.yaml"

# WebUI configuration
readonly WEBUI_DIR="$PROJECT_ROOT/unified"

# Default mode
MODE="production"

################################################################################
# Color output for better readability
################################################################################
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ==============================================================================
# Logging Functions
# ==============================================================================

# Enhanced logging with consistent formatting and optional timestamps
log() {
    local level="$1"
    local message="$2"
    local color="$3"
    local timestamp

    # Add timestamp for warnings and errors
    if [[ "$level" =~ ^(WARNING|ERROR)$ ]]; then
        timestamp="$(date '+%Y-%m-%d %H:%M:%S') "
    else
        timestamp=""
    fi

    echo -e "${color}[$level]${NC} ${timestamp}$message" >&2
}

log_info() { log "INFO" "$1" "$BLUE"; }
log_success() { log "SUCCESS" "$1" "$GREEN"; }
log_warning() { log "WARNING" "$1" "$YELLOW"; }
log_error() { log "ERROR" "$1" "$RED"; }

# ==============================================================================
# Argument Parsing and Validation
# ==============================================================================

# Parse and validate command line arguments
parse_args() {
    local args=("$@")

    # Process flags
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --dev|-d)
                MODE="development"
                log_info "Mode: Development"
                shift
                ;;
            --help|-h)
                show_usage
                exit 0
                ;;
            -*)
                log_error "Unknown option: $1"
                echo ""
                show_usage
                exit 1
                ;;
            *)
                # Non-flag arguments handled elsewhere
                break
                ;;
        esac
    done

    # Set default mode if not specified
    if [[ "$MODE" != "development" ]]; then
        MODE="production"
        log_info "Mode: Production"
    fi
}

# ==============================================================================
# Process Management Functions
# ==============================================================================

# Check if the server is running
# Returns: 0 if running, 1 if not running
is_running() {
    local pid

    # Check if PID file exists
    if [[ ! -f "$PID_FILE" ]]; then
        return 1  # Not running
    fi

    # Read PID safely
    if ! pid=$(<"$PID_FILE"); then
        log_warning "Failed to read PID file: $PID_FILE"
        return 1
    fi

    # Validate PID format (should be numeric)
    if ! [[ "$pid" =~ ^[0-9]+$ ]]; then
        log_warning "Invalid PID format in $PID_FILE: $pid"
        rm -f "$PID_FILE"
        return 1
    fi

    # Check if process is actually running
    if ps -p "$pid" > /dev/null 2>&1; then
        return 0  # Running
    else
        # PID file exists but process is not running
        log_warning "Stale PID file found, removing..."
        rm -f "$PID_FILE"
        return 1  # Not running
    fi
}

################################################################################
# Check if the WebUI is running
################################################################################
# Returns: 0 if running, 1 if not running
is_webui_running() {
    if [ -f "$WEBUI_PID_FILE" ]; then
        local pid=$(cat "$WEBUI_PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0  # Running
        else
            # PID file exists but process is not running
            log_warning "Stale WebUI PID file found, removing..."
            rm -f "$WEBUI_PID_FILE"
            return 1  # Not running
        fi
    fi
    return 1  # Not running
}

################################################################################
# Build the server binary with validation
build_server() {
    local start_time
    start_time=$(date +%s)

    log_info "Building gomailserver..."

    # Validate we're in the correct directory
    if [[ ! -f "$PROJECT_ROOT/Makefile" ]]; then
        log_error "Makefile not found in $PROJECT_ROOT"
        return 1
    fi

    # Navigate to project root (with error checking)
    if ! cd "$PROJECT_ROOT"; then
        log_error "Failed to change to project root: $PROJECT_ROOT"
        return 1
    fi

    # Build using make with timing
    log_info "Running make build..."
    if make build; then
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        log_success "Build completed successfully in ${duration}s"

        # Verify binary was created
        if [[ ! -x "$BINARY_PATH" ]]; then
            log_error "Build succeeded but binary not found: $BINARY_PATH"
            return 1
        fi

        return 0
    else
        log_error "Build failed after ${duration:-?}s"
        return 1
    fi
}

################################################################################
# Extract webui.port from the configuration file
################################################################################
get_webui_port() {
    local config_file="$1"

    if [ -f "$config_file" ]; then
        # Try to extract port from yaml format (webui: port: XXXX)
        local port=$(grep -E '^\s*port:\s*[0-9]+' "$config_file" 2>/dev/null | head -1 | grep -oE '[0-9]+' || echo "")

        if [ -n "$port" ]; then
            echo "$port"
            return 0
        fi
    fi

    # Default to 8080 if not found
    echo "8080"
    return 1
}

################################################################################
# Start the WebUI development server
################################################################################
start_webui() {
    # Only start WebUI in development mode
    if [ "$MODE" != "development" ]; then
        return 0
    fi

    # Check if already running
    if is_webui_running; then
        log_warning "WebUI is already running (PID: $(cat "$WEBUI_PID_FILE"))"
        return 1
    fi

    # Check if WebUI directory exists
    if [ ! -d "$WEBUI_DIR" ]; then
        log_error "WebUI directory not found at $WEBUI_DIR"
        return 1
    fi

    # Check if pnpm is installed
    if ! command -v pnpm > /dev/null 2>&1; then
        log_error "pnpm is not installed. Please install pnpm to run the WebUI."
        return 1
    fi

    # Check if node_modules exists, install if needed
    if [ ! -d "$WEBUI_DIR/node_modules" ]; then
        log_info "Installing WebUI dependencies..."
        cd "$WEBUI_DIR"
        if ! pnpm install; then
            log_error "Failed to install WebUI dependencies"
            return 1
        fi
    fi

    log_info "Starting WebUI development server..."

    # Get the Go server port from config for API proxying
    local config_file="$DEV_CONFIG"
    local go_port
    go_port=$(get_webui_port "$config_file")

    # Export environment variables for Nuxt to use
    export NUXT_PUBLIC_API_BASE="http://localhost:${go_port}/api/v1"
    log_info "API base URL: $NUXT_PUBLIC_API_BASE"

    # Navigate to WebUI directory
    cd "$WEBUI_DIR"

    # Start the WebUI in the background with environment variables
    NUXT_PUBLIC_API_BASE="$NUXT_PUBLIC_API_BASE" pnpm dev > "$WEBUI_LOG_FILE" 2>&1 &
    local pid=$!

    # Save PID to file
    echo "$pid" > "$WEBUI_PID_FILE"
    
    # Extract actual WebUI port from logs
    WEBUI_PORT=$(grep -oPmP 'Local:.*http://localhost:\([0-9][0-9]*\)/' "$WEBUI_LOG_FILE" 2>/dev/null | head -1 | sed 's/.*http:\/\/localhost:\([0-9][0-9]*\)//')

    if [ -z "$WEBUI_PORT" ]; then
        WEBUI_PORT="5173"
    fi

    # Wait a moment and check if the process is still running
    sleep 2
    if ps -p "$pid" > /dev/null 2>&1; then
        log_success "WebUI started successfully (PID: $pid)"
        log_info "WebUI logs: $WEBUI_LOG_FILE"
        log_info "WebUI available at:"
        log_info "  - Local:   http://localhost:5173/admin/"
        log_info "  - Network: Check logs for all network URLs"
        return 0
    else
        log_error "WebUI failed to start, check logs at $WEBUI_LOG_FILE"
        rm -f "$WEBUI_PID_FILE"
        return 1
    fi
}

################################################################################
# Stop the WebUI development server
################################################################################
stop_webui() {
    # Check if running
    if ! is_webui_running; then
        return 0  # Not an error, just not running
    fi

    local pid=$(cat "$WEBUI_PID_FILE")
    log_info "Stopping WebUI (PID: $pid)..."

    # Send SIGTERM for graceful shutdown
    kill -TERM "$pid" 2>/dev/null

    # Wait for process to terminate (max 10 seconds)
    local timeout=10
    local count=0
    while ps -p "$pid" > /dev/null 2>&1; do
        if [ $count -ge $timeout ]; then
            log_warning "WebUI did not terminate gracefully, forcing shutdown..."
            kill -KILL "$pid" 2>/dev/null
            break
        fi
        sleep 1
        count=$((count + 1))
    done

    # Remove PID file
    rm -f "$WEBUI_PID_FILE"

    log_success "WebUI stopped"
    return 0
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
    if [[ ! -f "$config_file" ]]; then
        log_warning "Development config not found, copying from example..."

        local example_config="$PROJECT_ROOT/gomailserver.example.yaml"
        if [[ -f "$example_config" ]]; then
            if cp "$example_config" "$config_file"; then
                log_info "Created development config at $config_file"
            else
                log_error "Failed to copy example config to $config_file"
                return 1
            fi
        else
            log_error "Example config not found at $example_config"
            log_info "Please create a development config manually or run from project root"
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

        # Start WebUI in development mode
        if [ "$MODE" = "development" ]; then
            start_webui
        fi

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
    # Stop WebUI first if it's running
    stop_webui

    # Check if server is running
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
    local server_running=false
    local webui_running=false

    # Check server status
    if is_running; then
        local pid=$(cat "$PID_FILE")
        log_success "gomailserver is running (PID: $pid)"
        server_running=true

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
    else
        log_warning "gomailserver is not running"
    fi

    # Check WebUI status
    echo ""
    if is_webui_running; then
        local webui_pid=$(cat "$WEBUI_PID_FILE")
        log_success "WebUI is running (PID: $webui_pid)"
        webui_running=true
        log_info "WebUI URLs:"
        log_info "  - Local:   http://localhost:5173/admin/"
        log_info "  - Network: Check $WEBUI_LOG_FILE for all network URLs"
    else
        log_info "WebUI is not running (only runs in --dev mode)"
    fi

    # Return success if at least one service is running
    if [ "$server_running" = true ]; then
        return 0
    else
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
  --dev, -d         Run in development mode (debug logging, local config, WebUI)

Examples:
  $0 start          Start in production mode
  $0 start --dev    Start in development mode (includes WebUI at http://localhost:5173)
  $0 stop           Stop the server (and WebUI if running)
  $0 restart --dev  Restart in development mode
  $0 status         Check if server and WebUI are running

Configuration:
  Development:  $DEV_CONFIG
  Production:   $PROD_CONFIG

Logs:
  Server:  $LOG_FILE
  WebUI:   $WEBUI_LOG_FILE (dev mode only)

PID Files:
  Server:  $PID_FILE
  WebUI:   $WEBUI_PID_FILE (dev mode only)

Development Mode:
  In --dev mode, the script automatically starts the unified WebUI development
  server (Vite) at http://localhost:5173 alongside the gomailserver backend.
  The WebUI is automatically stopped when the server is stopped.

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

# ==============================================================================
# Utility Functions
# ==============================================================================

# Cleanup function called on script exit
cleanup() {
    # Placeholder for future cleanup operations
    return 0
}

# Validate script environment
validate_environment() {
    local errors=0

    # Check if we're in a reasonable directory
    if [[ ! -d "$PROJECT_ROOT" ]]; then
        log_error "Project root not found: $PROJECT_ROOT"
        ((errors++))
    fi

    # Check for required commands
    local required_cmds=("ps" "mkdir" "rm" "cat")
    for cmd in "${required_cmds[@]}"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            log_error "Required command not found: $cmd"
            ((errors++))
        fi
    done

    return $errors
}

# ==============================================================================
# Main Script Logic
# ==============================================================================

main() {
    # Validate environment first
    if ! validate_environment; then
        log_error "Environment validation failed. Please check your setup."
        exit 1
    fi

    # Check if no arguments provided
    if [[ $# -eq 0 ]]; then
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
