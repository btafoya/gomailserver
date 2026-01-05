#!/bin/bash

################################################################################
# gomailserver systemd Installer Script
################################################################################
#
# This script installs gomailserver as a systemd service with proper security
# hardening and user/group configuration.
#
# Usage:
#   sudo ./install-systemd.sh [options]
#
# Options:
#   --start           Start the service immediately after installation
#   --enable          Enable the service to start on boot (default)
#   --no-enable       Do not enable the service on boot
#   --user USER       Run as specified user (default: gomailserver)
#   --group GROUP     Run as specified group (default: gomailserver)
#   --binary PATH     Path to gomailserver binary (default: ./build/gomailserver)
#   --prefix PATH     Installation prefix (default: /usr/local)
#
# Prerequisites:
#   - Root privileges (sudo)
#   - systemd installed
#   - gomailserver binary built
#
################################################################################

set -e

# Color output for better readability
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored messages
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

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Default configuration
SERVICE_NAME="gomailserver"
SERVICE_USER="gomailserver"
SERVICE_GROUP="gomailserver"
BINARY_PATH="$PROJECT_ROOT/build/gomailserver"
INSTALL_PREFIX="/usr/local"
ENABLE_SERVICE=true
START_SERVICE=false

################################################################################
# Parse command line arguments
################################################################################
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --start)
                START_SERVICE=true
                shift
                ;;
            --enable)
                ENABLE_SERVICE=true
                shift
                ;;
            --no-enable)
                ENABLE_SERVICE=false
                shift
                ;;
            --user)
                SERVICE_USER="$2"
                shift 2
                ;;
            --group)
                SERVICE_GROUP="$2"
                shift 2
                ;;
            --binary)
                BINARY_PATH="$2"
                shift 2
                ;;
            --prefix)
                INSTALL_PREFIX="$2"
                shift 2
                ;;
            --help|-h)
                show_usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

################################################################################
# Show usage information
################################################################################
show_usage() {
    cat << EOF
Usage: sudo $0 [options]

Installs gomailserver as a systemd service with security hardening.

Options:
  --start           Start the service immediately after installation
  --enable          Enable the service to start on boot (default)
  --no-enable       Do not enable the service on boot
  --user USER       Run as specified user (default: gomailserver)
  --group GROUP     Run as specified group (default: gomailserver)
  --binary PATH     Path to binary (default: ./build/gomailserver)
  --prefix PATH     Installation prefix (default: /usr/local)
  --help, -h        Show this help message

Examples:
  # Install with defaults (enable but don't start)
  sudo $0

  # Install and start immediately
  sudo $0 --start

  # Install with custom user and start
  sudo $0 --user mailserver --start

  # Install but don't enable on boot
  sudo $0 --no-enable

EOF
}

################################################################################
# Check prerequisites
################################################################################
check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check if running as root
    if [ "$EUID" -ne 0 ]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi

    # Check if systemd is installed
    if ! command -v systemctl &> /dev/null; then
        log_error "systemd is not installed or systemctl not found"
        exit 1
    fi

    # Check if binary exists
    if [ ! -f "$BINARY_PATH" ]; then
        log_error "Binary not found at $BINARY_PATH"
        log_info "Please build the binary first with: make build"
        exit 1
    fi

    log_success "Prerequisites check passed"
}

################################################################################
# Create service user and group
################################################################################
create_user() {
    log_info "Creating service user and group..."

    # Create group if it doesn't exist
    if ! getent group "$SERVICE_GROUP" > /dev/null 2>&1; then
        groupadd --system "$SERVICE_GROUP"
        log_success "Created group: $SERVICE_GROUP"
    else
        log_info "Group already exists: $SERVICE_GROUP"
    fi

    # Create user if it doesn't exist
    if ! getent passwd "$SERVICE_USER" > /dev/null 2>&1; then
        useradd --system \
            --gid "$SERVICE_GROUP" \
            --home-dir /var/lib/gomailserver \
            --shell /usr/sbin/nologin \
            --comment "gomailserver daemon user" \
            "$SERVICE_USER"
        log_success "Created user: $SERVICE_USER"
    else
        log_info "User already exists: $SERVICE_USER"
    fi
}

################################################################################
# Create directories
################################################################################
create_directories() {
    log_info "Creating directories..."

    # Binary directory
    install -d "$INSTALL_PREFIX/bin"

    # Configuration directory
    install -d -m 750 -o root -g "$SERVICE_GROUP" /etc/gomailserver

    # Data/state directory
    install -d -m 750 -o "$SERVICE_USER" -g "$SERVICE_GROUP" /var/lib/gomailserver

    # Log directory
    install -d -m 750 -o "$SERVICE_USER" -g "$SERVICE_GROUP" /var/log/gomailserver

    log_success "Directories created"
}

################################################################################
# Install binary
################################################################################
install_binary() {
    log_info "Installing binary..."

    install -m 755 "$BINARY_PATH" "$INSTALL_PREFIX/bin/gomailserver"
    log_success "Binary installed to $INSTALL_PREFIX/bin/gomailserver"
}

################################################################################
# Install configuration
################################################################################
install_config() {
    log_info "Installing configuration..."

    local config_dest="/etc/gomailserver/gomailserver.yaml"

    # Check if config already exists
    if [ -f "$config_dest" ]; then
        log_warning "Configuration already exists at $config_dest"
        log_info "Skipping config installation to preserve existing settings"
    else
        # Copy example config
        if [ -f "$PROJECT_ROOT/gomailserver.example.yaml" ]; then
            install -m 640 -o root -g "$SERVICE_GROUP" \
                "$PROJECT_ROOT/gomailserver.example.yaml" \
                "$config_dest"
            log_success "Configuration template installed to $config_dest"
            log_warning "Please edit $config_dest with your settings before starting"
        else
            log_warning "Example config not found, skipping config installation"
        fi
    fi
}

################################################################################
# Install systemd service
################################################################################
install_service() {
    log_info "Installing systemd service..."

    local service_src="$SCRIPT_DIR/gomailserver.service"
    local service_dest="/etc/systemd/system/gomailserver.service"

    if [ ! -f "$service_src" ]; then
        log_error "Service file not found at $service_src"
        exit 1
    fi

    # Copy service file
    install -m 644 "$service_src" "$service_dest"

    # Update service file with custom user/group if specified
    if [ "$SERVICE_USER" != "gomailserver" ] || [ "$SERVICE_GROUP" != "gomailserver" ]; then
        sed -i "s/^User=.*/User=$SERVICE_USER/" "$service_dest"
        sed -i "s/^Group=.*/Group=$SERVICE_GROUP/" "$service_dest"
        log_info "Updated service file with user=$SERVICE_USER group=$SERVICE_GROUP"
    fi

    # Reload systemd daemon
    systemctl daemon-reload

    log_success "systemd service installed"
}

################################################################################
# Enable service
################################################################################
enable_service() {
    if [ "$ENABLE_SERVICE" = true ]; then
        log_info "Enabling service to start on boot..."
        systemctl enable gomailserver.service
        log_success "Service enabled"
    else
        log_info "Service not enabled (--no-enable specified)"
    fi
}

################################################################################
# Start service
################################################################################
start_service() {
    if [ "$START_SERVICE" = true ]; then
        log_info "Starting service..."

        # Check if config is properly set up
        if ! "$INSTALL_PREFIX/bin/gomailserver" run --config /etc/gomailserver/gomailserver.yaml --help &> /dev/null; then
            log_warning "Service not started - please configure /etc/gomailserver/gomailserver.yaml first"
            return
        fi

        systemctl start gomailserver.service
        sleep 2

        # Check if service started successfully
        if systemctl is-active --quiet gomailserver.service; then
            log_success "Service started successfully"
            systemctl status gomailserver.service --no-pager
        else
            log_error "Service failed to start"
            log_info "Check logs with: journalctl -u gomailserver.service -n 50"
            exit 1
        fi
    else
        log_info "Service not started (use --start to start immediately)"
    fi
}

################################################################################
# Print post-installation instructions
################################################################################
print_instructions() {
    cat << EOF

${GREEN}===================================================================${NC}
${GREEN}gomailserver systemd Installation Complete!${NC}
${GREEN}===================================================================${NC}

${BLUE}Configuration:${NC}
  Config file: /etc/gomailserver/gomailserver.yaml
  Binary:      $INSTALL_PREFIX/bin/gomailserver
  Data dir:    /var/lib/gomailserver
  Log dir:     /var/log/gomailserver

${BLUE}Service Management:${NC}
  Start:       sudo systemctl start gomailserver
  Stop:        sudo systemctl stop gomailserver
  Restart:     sudo systemctl restart gomailserver
  Status:      sudo systemctl status gomailserver
  Logs:        sudo journalctl -u gomailserver -f

${BLUE}Next Steps:${NC}
  1. Edit configuration: sudo nano /etc/gomailserver/gomailserver.yaml
  2. Create admin user:  sudo -u $SERVICE_USER $INSTALL_PREFIX/bin/gomailserver create-admin
  3. Start service:      sudo systemctl start gomailserver

${YELLOW}Important:${NC}
  - The service runs as user '$SERVICE_USER'
  - Mail ports (25, 587, 465, 143, 993) require privileged binding
  - Ensure firewall allows mail traffic
  - Configure TLS certificates before production use

${GREEN}===================================================================${NC}

EOF
}

################################################################################
# Main installation flow
################################################################################
main() {
    echo ""
    log_info "gomailserver systemd Installer"
    echo ""

    # Parse arguments
    parse_args "$@"

    # Run installation steps
    check_prerequisites
    create_user
    create_directories
    install_binary
    install_config
    install_service
    enable_service
    start_service

    # Print instructions
    print_instructions
}

# Run main function with all arguments
main "$@"
