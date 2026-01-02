# TASKS9.md - Phase 9: Polish & Documentation (Weeks 27-29)

## Overview

Installation tooling, Docker support, embedded assets, documentation, and backup system.

**Total Tasks**: 19
**MVP Tasks**: 4 (Embedded Assets + Database Migrations)
**Priority**: Mixed - Some MVP, mostly FULL
**Dependencies**: Phase 0-8

---

## 9.1 Installation

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| IN-001 | Debian/Ubuntu installation script | [ ] | All Phase 1-3 | FULL |
| IN-002 | Systemd service file | [ ] | IN-001 | FULL |
| IN-003 | Configuration validation tool | [ ] | F-011 | FULL |
| IN-004 | Pre-flight check utility | [ ] | IN-003 | FULL |

---

### IN-001: Debian/Ubuntu Installation Script
**File**: `scripts/install.sh`
```bash
#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Defaults
INSTALL_DIR="/opt/gomailserver"
CONFIG_DIR="/etc/gomailserver"
DATA_DIR="/var/lib/gomailserver"
LOG_DIR="/var/log/gomailserver"
USER="gomailserver"
GROUP="gomailserver"
VERSION="${VERSION:-latest}"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

check_os() {
    if [[ ! -f /etc/debian_version ]]; then
        log_error "This script is designed for Debian/Ubuntu systems"
        exit 1
    fi
    log_info "Detected OS: $(cat /etc/os-release | grep PRETTY_NAME | cut -d= -f2)"
}

install_dependencies() {
    log_info "Installing dependencies..."
    apt-get update
    apt-get install -y \
        curl \
        wget \
        ca-certificates \
        sqlite3 \
        clamav \
        clamav-daemon \
        spamassassin \
        spamd \
        openssl \
        certbot
}

create_user() {
    if id "$USER" &>/dev/null; then
        log_info "User $USER already exists"
    else
        log_info "Creating system user $USER..."
        groupadd --system "$GROUP"
        useradd --system --gid "$GROUP" --no-create-home --shell /usr/sbin/nologin "$USER"
    fi
}

create_directories() {
    log_info "Creating directories..."

    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$DATA_DIR"/{mail,db,certs,backups}
    mkdir -p "$LOG_DIR"

    chown -R "$USER:$GROUP" "$DATA_DIR"
    chown -R "$USER:$GROUP" "$LOG_DIR"
    chmod 750 "$DATA_DIR"
    chmod 750 "$LOG_DIR"
}

download_binary() {
    log_info "Downloading gomailserver $VERSION..."

    local ARCH=$(uname -m)
    case $ARCH in
        x86_64)  ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        *)       log_error "Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    if [[ "$VERSION" == "latest" ]]; then
        VERSION=$(curl -s https://api.github.com/repos/btafoya/gomailserver/releases/latest | grep tag_name | cut -d'"' -f4)
    fi

    local URL="https://github.com/btafoya/gomailserver/releases/download/${VERSION}/gomailserver-linux-${ARCH}"

    wget -O "$INSTALL_DIR/gomailserver" "$URL"
    chmod +x "$INSTALL_DIR/gomailserver"
    ln -sf "$INSTALL_DIR/gomailserver" /usr/local/bin/gomailserver
}

install_config() {
    if [[ -f "$CONFIG_DIR/gomailserver.yaml" ]]; then
        log_warn "Configuration file already exists, skipping..."
        return
    fi

    log_info "Creating default configuration..."

    cat > "$CONFIG_DIR/gomailserver.yaml" << 'EOF'
# gomailserver configuration
# See documentation for all options

server:
  hostname: mail.example.com

database:
  path: /var/lib/gomailserver/db/gomailserver.db

storage:
  mail_dir: /var/lib/gomailserver/mail
  threshold_mb: 1

smtp:
  submission:
    enabled: true
    port: 587
    require_tls: true
  relay:
    enabled: true
    port: 25
  smtps:
    enabled: true
    port: 465

imap:
  enabled: true
  port: 143
  imaps_port: 993

tls:
  cert_dir: /var/lib/gomailserver/certs
  acme:
    enabled: true
    email: admin@example.com
    provider: cloudflare

security:
  dkim:
    enabled: true
    selector: default
  spf:
    enabled: true
  dmarc:
    enabled: true
  clamav:
    enabled: true
    socket: /var/run/clamav/clamd.ctl
  spamassassin:
    enabled: true
    host: localhost
    port: 783

web:
  enabled: true
  port: 8080
  admin:
    enabled: true
  portal:
    enabled: true
  webmail:
    enabled: true

logging:
  level: info
  format: json
  output: /var/log/gomailserver/gomailserver.log
EOF

    chown "$USER:$GROUP" "$CONFIG_DIR/gomailserver.yaml"
    chmod 640 "$CONFIG_DIR/gomailserver.yaml"
}

install_systemd() {
    log_info "Installing systemd service..."

    cat > /etc/systemd/system/gomailserver.service << EOF
[Unit]
Description=gomailserver Mail Server
Documentation=https://github.com/btafoya/gomailserver
After=network.target clamav-daemon.service spamassassin.service
Wants=clamav-daemon.service spamassassin.service

[Service]
Type=simple
User=$USER
Group=$GROUP
ExecStart=$INSTALL_DIR/gomailserver run --config $CONFIG_DIR/gomailserver.yaml
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=$DATA_DIR $LOG_DIR
PrivateTmp=yes
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes

# Allow binding to privileged ports
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    log_info "Systemd service installed"
}

configure_services() {
    log_info "Configuring ClamAV..."
    systemctl enable clamav-daemon
    systemctl start clamav-daemon || log_warn "ClamAV failed to start, may need freshclam first"

    log_info "Configuring SpamAssassin..."
    systemctl enable spamassassin
    systemctl start spamassassin || log_warn "SpamAssassin failed to start"
}

print_post_install() {
    echo ""
    echo "==========================================="
    echo -e "${GREEN}gomailserver installation complete!${NC}"
    echo "==========================================="
    echo ""
    echo "Next steps:"
    echo "  1. Edit configuration: $CONFIG_DIR/gomailserver.yaml"
    echo "  2. Run setup wizard:   gomailserver setup"
    echo "  3. Start server:       systemctl start gomailserver"
    echo "  4. Enable on boot:     systemctl enable gomailserver"
    echo ""
    echo "Check status: systemctl status gomailserver"
    echo "View logs:    journalctl -u gomailserver -f"
    echo ""
    echo "Documentation: https://github.com/btafoya/gomailserver"
}

# Main
main() {
    echo ""
    echo "gomailserver Installer"
    echo "======================"
    echo ""

    check_root
    check_os
    install_dependencies
    create_user
    create_directories
    download_binary
    install_config
    install_systemd
    configure_services
    print_post_install
}

main "$@"
```

**Acceptance Criteria**:
- [ ] Root check
- [ ] OS detection (Debian/Ubuntu)
- [ ] Dependency installation
- [ ] User/group creation
- [ ] Directory structure
- [ ] Binary download
- [ ] Default config creation

---

### IN-002: Systemd Service File
**File**: `dist/gomailserver.service`
```ini
[Unit]
Description=gomailserver Mail Server
Documentation=https://github.com/btafoya/gomailserver
After=network.target clamav-daemon.service spamassassin.service
Wants=clamav-daemon.service spamassassin.service

[Service]
Type=simple
User=gomailserver
Group=gomailserver
ExecStart=/opt/gomailserver/gomailserver run --config /etc/gomailserver/gomailserver.yaml
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=gomailserver

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/var/lib/gomailserver /var/log/gomailserver
PrivateTmp=yes
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes
RestrictNamespaces=yes
RestrictRealtime=yes
RestrictSUIDSGID=yes

# Allow binding to privileged ports
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

---

### IN-003: Configuration Validation Tool
**File**: `internal/config/validator.go`
```go
package config

import (
    "fmt"
    "net"
    "os"
    "strings"

    "github.com/go-playground/validator/v10"
)

type ValidationResult struct {
    Valid    bool
    Errors   []string
    Warnings []string
}

type ConfigValidator struct {
    validate *validator.Validate
}

func NewConfigValidator() *ConfigValidator {
    v := validator.New()

    // Custom validations
    v.RegisterValidation("hostname", validateHostname)
    v.RegisterValidation("filepath", validateFilePath)
    v.RegisterValidation("dirpath", validateDirPath)

    return &ConfigValidator{validate: v}
}

func (cv *ConfigValidator) Validate(cfg *Config) *ValidationResult {
    result := &ValidationResult{Valid: true}

    // Struct validation
    if err := cv.validate.Struct(cfg); err != nil {
        for _, e := range err.(validator.ValidationErrors) {
            result.Errors = append(result.Errors,
                fmt.Sprintf("%s: %s validation failed", e.Field(), e.Tag()))
        }
        result.Valid = false
    }

    // Business logic validation
    cv.validatePorts(cfg, result)
    cv.validateTLS(cfg, result)
    cv.validateDatabase(cfg, result)
    cv.validateSecurity(cfg, result)

    return result
}

func (cv *ConfigValidator) validatePorts(cfg *Config, result *ValidationResult) {
    ports := make(map[int]string)

    if cfg.SMTP.Submission.Enabled {
        if _, exists := ports[cfg.SMTP.Submission.Port]; exists {
            result.Errors = append(result.Errors,
                fmt.Sprintf("Port %d already in use by %s", cfg.SMTP.Submission.Port, ports[cfg.SMTP.Submission.Port]))
            result.Valid = false
        }
        ports[cfg.SMTP.Submission.Port] = "SMTP Submission"
    }

    if cfg.SMTP.Relay.Enabled {
        if _, exists := ports[cfg.SMTP.Relay.Port]; exists {
            result.Errors = append(result.Errors,
                fmt.Sprintf("Port %d already in use by %s", cfg.SMTP.Relay.Port, ports[cfg.SMTP.Relay.Port]))
            result.Valid = false
        }
        ports[cfg.SMTP.Relay.Port] = "SMTP Relay"
    }

    // Check if ports are available on system
    for port, service := range ports {
        addr := fmt.Sprintf(":%d", port)
        ln, err := net.Listen("tcp", addr)
        if err != nil {
            result.Warnings = append(result.Warnings,
                fmt.Sprintf("Port %d (%s) may already be in use or require elevated permissions", port, service))
        } else {
            ln.Close()
        }
    }
}

func (cv *ConfigValidator) validateTLS(cfg *Config, result *ValidationResult) {
    if !cfg.TLS.ACME.Enabled {
        // Manual certs - check files exist
        if cfg.TLS.CertFile != "" {
            if _, err := os.Stat(cfg.TLS.CertFile); os.IsNotExist(err) {
                result.Errors = append(result.Errors,
                    fmt.Sprintf("TLS certificate file not found: %s", cfg.TLS.CertFile))
                result.Valid = false
            }
        }
        if cfg.TLS.KeyFile != "" {
            if _, err := os.Stat(cfg.TLS.KeyFile); os.IsNotExist(err) {
                result.Errors = append(result.Errors,
                    fmt.Sprintf("TLS key file not found: %s", cfg.TLS.KeyFile))
                result.Valid = false
            }
        }
    }

    // Check TLS is enabled for secure services
    if cfg.SMTP.Submission.Enabled && !cfg.SMTP.Submission.RequireTLS {
        result.Warnings = append(result.Warnings,
            "SMTP submission should require TLS for security")
    }
}

func (cv *ConfigValidator) validateDatabase(cfg *Config, result *ValidationResult) {
    // Check directory is writable
    dbDir := filepath.Dir(cfg.Database.Path)
    if _, err := os.Stat(dbDir); os.IsNotExist(err) {
        result.Errors = append(result.Errors,
            fmt.Sprintf("Database directory does not exist: %s", dbDir))
        result.Valid = false
    }
}

func (cv *ConfigValidator) validateSecurity(cfg *Config, result *ValidationResult) {
    if cfg.Security.ClamAV.Enabled {
        if _, err := os.Stat(cfg.Security.ClamAV.Socket); os.IsNotExist(err) {
            result.Warnings = append(result.Warnings,
                "ClamAV socket not found - ensure clamd is running")
        }
    }

    if !cfg.Security.DKIM.Enabled {
        result.Warnings = append(result.Warnings,
            "DKIM is disabled - outgoing mail may be flagged as spam")
    }

    if !cfg.Security.SPF.Enabled {
        result.Warnings = append(result.Warnings,
            "SPF checking is disabled")
    }
}

// CLI command
func ValidateConfigCmd(configPath string) {
    cfg, err := LoadConfig(configPath)
    if err != nil {
        fmt.Printf("❌ Failed to load configuration: %s\n", err)
        os.Exit(1)
    }

    validator := NewConfigValidator()
    result := validator.Validate(cfg)

    fmt.Println("\nConfiguration Validation Report")
    fmt.Println("================================\n")

    if len(result.Errors) > 0 {
        fmt.Println("❌ Errors:")
        for _, e := range result.Errors {
            fmt.Printf("   • %s\n", e)
        }
        fmt.Println()
    }

    if len(result.Warnings) > 0 {
        fmt.Println("⚠️  Warnings:")
        for _, w := range result.Warnings {
            fmt.Printf("   • %s\n", w)
        }
        fmt.Println()
    }

    if result.Valid {
        fmt.Println("✅ Configuration is valid")
        os.Exit(0)
    } else {
        fmt.Println("❌ Configuration has errors")
        os.Exit(1)
    }
}
```

**CLI Integration** (`cmd/gomailserver/validate.go`):
```go
var validateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate configuration file",
    Run: func(cmd *cobra.Command, args []string) {
        configPath, _ := cmd.Flags().GetString("config")
        config.ValidateConfigCmd(configPath)
    },
}

func init() {
    rootCmd.AddCommand(validateCmd)
}
```

**Acceptance Criteria**:
- [ ] Port conflict detection
- [ ] TLS configuration validation
- [ ] Path existence checks
- [ ] Security recommendations
- [ ] Clear error/warning output

---

### IN-004: Pre-flight Check Utility
**File**: `cmd/gomailserver/preflight.go`
```go
package main

import (
    "fmt"
    "net"
    "os"
    "os/exec"
    "syscall"

    "github.com/spf13/cobra"
)

var preflightCmd = &cobra.Command{
    Use:   "preflight",
    Short: "Run pre-flight checks before starting the server",
    Run:   runPreflight,
}

type PreflightCheck struct {
    Name    string
    Check   func() (bool, string)
    Fatal   bool
}

func runPreflight(cmd *cobra.Command, args []string) {
    checks := []PreflightCheck{
        {"System resources", checkResources, false},
        {"Required ports", checkPorts, true},
        {"DNS resolution", checkDNS, false},
        {"ClamAV daemon", checkClamAV, false},
        {"SpamAssassin", checkSpamAssassin, false},
        {"Disk space", checkDiskSpace, true},
        {"File permissions", checkPermissions, true},
        {"Network connectivity", checkNetwork, false},
    }

    fmt.Println("\n🔍 Running pre-flight checks...\n")

    allPassed := true
    for _, check := range checks {
        passed, message := check.Check()

        if passed {
            fmt.Printf("✅ %s: %s\n", check.Name, message)
        } else if check.Fatal {
            fmt.Printf("❌ %s: %s\n", check.Name, message)
            allPassed = false
        } else {
            fmt.Printf("⚠️  %s: %s\n", check.Name, message)
        }
    }

    fmt.Println()
    if allPassed {
        fmt.Println("✅ All pre-flight checks passed")
        os.Exit(0)
    } else {
        fmt.Println("❌ Some pre-flight checks failed")
        os.Exit(1)
    }
}

func checkResources() (bool, string) {
    var info syscall.Sysinfo_t
    if err := syscall.Sysinfo(&info); err != nil {
        return false, "Unable to get system info"
    }

    totalMB := info.Totalram * uint64(info.Unit) / 1024 / 1024
    if totalMB < 512 {
        return false, fmt.Sprintf("Low memory: %dMB (recommended: 512MB+)", totalMB)
    }

    return true, fmt.Sprintf("%dMB RAM available", totalMB)
}

func checkPorts() (bool, string) {
    ports := []int{25, 143, 465, 587, 993, 8080}
    blocked := []int{}

    for _, port := range ports {
        ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
        if err != nil {
            blocked = append(blocked, port)
        } else {
            ln.Close()
        }
    }

    if len(blocked) > 0 {
        return false, fmt.Sprintf("Ports unavailable: %v", blocked)
    }
    return true, "All required ports available"
}

func checkDNS() (bool, string) {
    _, err := net.LookupHost("gmail.com")
    if err != nil {
        return false, "DNS resolution failed"
    }
    return true, "DNS resolution working"
}

func checkClamAV() (bool, string) {
    socketPaths := []string{
        "/var/run/clamav/clamd.ctl",
        "/var/run/clamd.scan/clamd.sock",
        "/tmp/clamd.socket",
    }

    for _, path := range socketPaths {
        if _, err := os.Stat(path); err == nil {
            return true, fmt.Sprintf("Socket found at %s", path)
        }
    }
    return false, "ClamAV socket not found - install or start clamd"
}

func checkSpamAssassin() (bool, string) {
    conn, err := net.Dial("tcp", "localhost:783")
    if err != nil {
        return false, "SpamAssassin not responding on port 783"
    }
    conn.Close()
    return true, "SpamAssassin responding"
}

func checkDiskSpace() (bool, string) {
    var stat syscall.Statfs_t
    dataDir := "/var/lib/gomailserver"

    if err := syscall.Statfs(dataDir, &stat); err != nil {
        if os.IsNotExist(err) {
            return false, fmt.Sprintf("Data directory %s does not exist", dataDir)
        }
        return false, "Unable to check disk space"
    }

    freeMB := stat.Bavail * uint64(stat.Bsize) / 1024 / 1024
    if freeMB < 1024 {
        return false, fmt.Sprintf("Low disk space: %dMB (recommended: 1GB+)", freeMB)
    }

    return true, fmt.Sprintf("%dMB free space", freeMB)
}

func checkPermissions() (bool, string) {
    paths := []string{
        "/var/lib/gomailserver",
        "/var/log/gomailserver",
        "/etc/gomailserver",
    }

    for _, path := range paths {
        info, err := os.Stat(path)
        if os.IsNotExist(err) {
            continue // Will be created
        }
        if err != nil {
            return false, fmt.Sprintf("Cannot access %s", path)
        }
        if !info.IsDir() {
            return false, fmt.Sprintf("%s is not a directory", path)
        }
    }

    return true, "Permissions OK"
}

func checkNetwork() (bool, string) {
    // Try to connect to a known mail server
    conn, err := net.DialTimeout("tcp", "gmail-smtp-in.l.google.com:25", 5*time.Second)
    if err != nil {
        return false, "Cannot reach external mail servers (port 25 may be blocked)"
    }
    conn.Close()
    return true, "Can reach external mail servers"
}

func init() {
    rootCmd.AddCommand(preflightCmd)
}
```

**Acceptance Criteria**:
- [ ] Resource checks (RAM, CPU)
- [ ] Port availability
- [ ] DNS resolution
- [ ] ClamAV/SpamAssassin connectivity
- [ ] Disk space
- [ ] Network connectivity

---

## 9.2 Docker

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| DK-001 | Dockerfile (Alpine base) | [ ] | All Phase 1-3 | FULL |
| DK-002 | Docker Compose configuration | [ ] | DK-001 | FULL |
| DK-003 | Multi-architecture builds | [ ] | DK-001 | FULL |

---

### DK-001: Dockerfile
**File**: `Dockerfile`
```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with SQLite support
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o gomailserver ./cmd/gomailserver

# Build web assets
FROM node:20-alpine AS frontend

WORKDIR /app/web

COPY web/admin/package*.json ./admin/
COPY web/portal/package*.json ./portal/
COPY web/webmail/package*.json ./webmail/

RUN cd admin && npm ci && cd .. && \
    cd portal && npm ci && cd .. && \
    cd webmail && npm ci

COPY web/ ./

RUN cd admin && npm run build && cd .. && \
    cd portal && npm run build && cd .. && \
    cd webmail && npm run build

# Final stage
FROM alpine:3.19

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    sqlite \
    && addgroup -S gomailserver \
    && adduser -S -G gomailserver gomailserver

WORKDIR /app

COPY --from=builder /app/gomailserver /app/gomailserver
COPY --from=frontend /app/web/admin/dist /app/web/admin
COPY --from=frontend /app/web/portal/dist /app/web/portal
COPY --from=frontend /app/web/webmail/dist /app/web/webmail

# Create directories
RUN mkdir -p /data/db /data/mail /data/certs /data/backups \
    && chown -R gomailserver:gomailserver /data

# Default config
COPY dist/docker-config.yaml /etc/gomailserver/gomailserver.yaml

USER gomailserver

EXPOSE 25 143 465 587 993 8080

VOLUME ["/data"]

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/gomailserver"]
CMD ["run", "--config", "/etc/gomailserver/gomailserver.yaml"]
```

---

### DK-002: Docker Compose
**File**: `docker-compose.yml`
```yaml
version: '3.8'

services:
  gomailserver:
    build: .
    container_name: gomailserver
    hostname: mail.example.com
    ports:
      - "25:25"     # SMTP relay
      - "143:143"   # IMAP
      - "465:465"   # SMTPS
      - "587:587"   # SMTP submission
      - "993:993"   # IMAPS
      - "8080:8080" # Web UI
    volumes:
      - mail-data:/data
      - ./config/gomailserver.yaml:/etc/gomailserver/gomailserver.yaml:ro
    environment:
      - TZ=UTC
    depends_on:
      - clamav
      - spamassassin
    networks:
      - mail-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  clamav:
    image: clamav/clamav:latest
    container_name: clamav
    volumes:
      - clamav-data:/var/lib/clamav
    networks:
      - mail-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "clamdscan", "--ping", "1"]
      interval: 60s
      timeout: 10s
      retries: 3

  spamassassin:
    image: spamassassin/spamassassin:latest
    container_name: spamassassin
    volumes:
      - spamassassin-data:/var/lib/spamassassin
    networks:
      - mail-network
    restart: unless-stopped

volumes:
  mail-data:
  clamav-data:
  spamassassin-data:

networks:
  mail-network:
    driver: bridge
```

---

### DK-003: Multi-architecture Builds
**File**: `.github/workflows/docker.yml`
```yaml
name: Docker Build

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Login to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: |
            btafoya/gomailserver
            ghcr.io/btafoya/gomailserver
          tags: |
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

**Acceptance Criteria**:
- [ ] AMD64 and ARM64 builds
- [ ] Docker Hub publishing
- [ ] GitHub Container Registry publishing
- [ ] Version tagging

---

## 9.3 Embed Assets [MVP]

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| EA-001 | Embed admin UI in binary | [ ] | AUI-013 | MVP |
| EA-002 | Embed portal UI in binary | [ ] | UP-007 | MVP |
| EA-003 | Embed webmail in binary | [ ] | WF-019 | MVP |

---

### EA-001: Embed Admin UI
**File**: `internal/web/embed.go`
```go
package web

import (
    "embed"
    "io/fs"
    "net/http"

    "github.com/labstack/echo/v4"
)

//go:embed all:admin/dist
var adminFS embed.FS

//go:embed all:portal/dist
var portalFS embed.FS

//go:embed all:webmail/dist
var webmailFS embed.FS

// GetAdminFS returns the admin UI filesystem
func GetAdminFS() http.FileSystem {
    subFS, _ := fs.Sub(adminFS, "admin/dist")
    return http.FS(subFS)
}

// GetPortalFS returns the user portal filesystem
func GetPortalFS() http.FileSystem {
    subFS, _ := fs.Sub(portalFS, "portal/dist")
    return http.FS(subFS)
}

// GetWebmailFS returns the webmail filesystem
func GetWebmailFS() http.FileSystem {
    subFS, _ := fs.Sub(webmailFS, "webmail/dist")
    return http.FS(subFS)
}

// ServeEmbeddedUI serves embedded UI with SPA fallback
func ServeEmbeddedUI(e *echo.Echo, path string, filesystem http.FileSystem) {
    fileServer := http.FileServer(filesystem)

    e.GET(path+"*", echo.WrapHandler(http.StripPrefix(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Try to serve the file
        f, err := filesystem.Open(r.URL.Path)
        if err != nil {
            // Fallback to index.html for SPA routing
            r.URL.Path = "/"
        } else {
            f.Close()
        }
        fileServer.ServeHTTP(w, r)
    }))))
}
```

**Build Integration**:
```makefile
# Makefile
.PHONY: build-web build

build-web:
	cd web/admin && npm run build
	cd web/portal && npm run build
	cd web/webmail && npm run build
	cp -r web/admin/dist internal/web/admin/
	cp -r web/portal/dist internal/web/portal/
	cp -r web/webmail/dist internal/web/webmail/

build: build-web
	go build -o gomailserver ./cmd/gomailserver
```

**Acceptance Criteria**:
- [ ] Admin UI embedded in binary
- [ ] Portal UI embedded in binary
- [ ] Webmail embedded in binary
- [ ] SPA routing fallback
- [ ] Single binary deployment

---

## 9.4 Documentation

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| DOC-001 | README with quick start | [ ] | IN-001 | FULL |
| DOC-002 | Installation guide | [ ] | IN-001 | FULL |
| DOC-003 | Administration guide | [ ] | AUI-013 | FULL |
| DOC-004 | User guide | [ ] | UP-007 | FULL |
| DOC-005 | API documentation | [ ] | API-006 | FULL |
| DOC-006 | Architecture documentation | [ ] | F-002 | FULL |
| DOC-007 | Troubleshooting guide | [ ] | All | FULL |
| DOC-008 | DNS setup guide | [ ] | SW-006 | FULL |

---

### DOC-008: DNS Setup Guide
**File**: `docs/dns-setup.md`
```markdown
# DNS Setup Guide

## Required DNS Records

### MX Record
```
example.com.  IN  MX  10  mail.example.com.
```

### A/AAAA Records
```
mail.example.com.  IN  A     203.0.113.1
mail.example.com.  IN  AAAA  2001:db8::1
```

### SPF Record
```
example.com.  IN  TXT  "v=spf1 mx a:mail.example.com -all"
```

### DKIM Record
Get your DKIM public key from the admin UI or:
```bash
gomailserver dkim show --domain example.com
```

Add the record:
```
default._domainkey.example.com.  IN  TXT  "v=DKIM1; k=rsa; p=MIGfMA0..."
```

### DMARC Record
```
_dmarc.example.com.  IN  TXT  "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com; ruf=mailto:dmarc@example.com; pct=100"
```

### MTA-STS Record
```
_mta-sts.example.com.  IN  TXT  "v=STSv1; id=20240101T000000"
```

### TLSRPT Record
```
_smtp._tls.example.com.  IN  TXT  "v=TLSRPTv1; rua=mailto:tlsrpt@example.com"
```

## Verification

Use the built-in DNS checker:
```bash
gomailserver dns check --domain example.com
```
```

---

## 9.5 Backup System

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| BK-001 | Backup CLI command | [ ] | F-012, F-020 | FULL |
| BK-002 | Restore CLI command | [ ] | BK-001 | FULL |
| BK-003 | Scheduled automatic backups | [ ] | BK-001 | FULL |
| BK-004 | 30-day retention policy | [ ] | BK-003 | FULL |
| BK-005 | Backup integrity verification | [ ] | BK-001 | FULL |

---

### BK-001: Backup CLI Command
**File**: `cmd/gomailserver/backup.go`
```go
package main

import (
    "archive/tar"
    "compress/gzip"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"

    "github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
    Use:   "backup",
    Short: "Create a backup of all data",
    Run:   runBackup,
}

func init() {
    rootCmd.AddCommand(backupCmd)
    backupCmd.Flags().StringP("output", "o", "", "Output file path (default: backups/gomailserver-{timestamp}.tar.gz)")
    backupCmd.Flags().Bool("include-mail", true, "Include mail storage in backup")
}

func runBackup(cmd *cobra.Command, args []string) {
    output, _ := cmd.Flags().GetString("output")
    includeMail, _ := cmd.Flags().GetBool("include-mail")

    cfg := loadConfig(cmd)

    if output == "" {
        timestamp := time.Now().Format("20060102-150405")
        output = filepath.Join(cfg.Data.BackupDir, fmt.Sprintf("gomailserver-%s.tar.gz", timestamp))
    }

    // Ensure backup directory exists
    os.MkdirAll(filepath.Dir(output), 0750)

    fmt.Printf("Creating backup: %s\n", output)

    // Create tar.gz file
    file, err := os.Create(output)
    if err != nil {
        fmt.Printf("Error creating backup file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    gzw := gzip.NewWriter(file)
    defer gzw.Close()

    tw := tar.NewWriter(gzw)
    defer tw.Close()

    // Backup database
    fmt.Println("  📦 Backing up database...")
    if err := backupDatabase(tw, cfg.Database.Path); err != nil {
        fmt.Printf("Error backing up database: %v\n", err)
        os.Exit(1)
    }

    // Backup configuration
    fmt.Println("  📦 Backing up configuration...")
    if err := addFileToTar(tw, "/etc/gomailserver/gomailserver.yaml", "config/gomailserver.yaml"); err != nil {
        fmt.Printf("Warning: Could not backup config: %v\n", err)
    }

    // Backup certificates
    fmt.Println("  📦 Backing up certificates...")
    if err := addDirToTar(tw, cfg.TLS.CertDir, "certs"); err != nil {
        fmt.Printf("Warning: Could not backup certs: %v\n", err)
    }

    // Backup mail storage
    if includeMail {
        fmt.Println("  📦 Backing up mail storage (this may take a while)...")
        if err := addDirToTar(tw, cfg.Storage.MailDir, "mail"); err != nil {
            fmt.Printf("Error backing up mail: %v\n", err)
            os.Exit(1)
        }
    }

    tw.Close()
    gzw.Close()
    file.Close()

    // Calculate checksum
    checksum, _ := calculateChecksum(output)
    checksumFile := output + ".sha256"
    os.WriteFile(checksumFile, []byte(checksum+"  "+filepath.Base(output)+"\n"), 0644)

    // Get file size
    info, _ := os.Stat(output)
    sizeMB := float64(info.Size()) / 1024 / 1024

    fmt.Printf("\n✅ Backup complete!\n")
    fmt.Printf("   File: %s\n", output)
    fmt.Printf("   Size: %.2f MB\n", sizeMB)
    fmt.Printf("   Checksum: %s\n", checksumFile)
}

func backupDatabase(tw *tar.Writer, dbPath string) error {
    // For SQLite, we need to handle WAL mode properly
    // Use .backup command or VACUUM INTO for consistent backup

    backupPath := dbPath + ".backup"
    defer os.Remove(backupPath)

    // Use SQLite backup API
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return err
    }
    defer db.Close()

    _, err = db.Exec(fmt.Sprintf("VACUUM INTO '%s'", backupPath))
    if err != nil {
        return err
    }

    return addFileToTar(tw, backupPath, "db/gomailserver.db")
}

func addFileToTar(tw *tar.Writer, filePath, tarPath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    info, err := file.Stat()
    if err != nil {
        return err
    }

    header, err := tar.FileInfoHeader(info, "")
    if err != nil {
        return err
    }
    header.Name = tarPath

    if err := tw.WriteHeader(header); err != nil {
        return err
    }

    _, err = io.Copy(tw, file)
    return err
}

func addDirToTar(tw *tar.Writer, dirPath, tarPrefix string) error {
    return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        relPath, _ := filepath.Rel(dirPath, path)
        tarPath := filepath.Join(tarPrefix, relPath)

        if info.IsDir() {
            return nil
        }

        return addFileToTar(tw, path, tarPath)
    })
}

func calculateChecksum(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(hash.Sum(nil)), nil
}
```

---

### BK-002: Restore CLI Command
```go
var restoreCmd = &cobra.Command{
    Use:   "restore [backup-file]",
    Short: "Restore from a backup",
    Args:  cobra.ExactArgs(1),
    Run:   runRestore,
}

func init() {
    rootCmd.AddCommand(restoreCmd)
    restoreCmd.Flags().Bool("verify-only", false, "Only verify the backup, don't restore")
    restoreCmd.Flags().Bool("force", false, "Overwrite existing data without confirmation")
}

func runRestore(cmd *cobra.Command, args []string) {
    backupPath := args[0]
    verifyOnly, _ := cmd.Flags().GetBool("verify-only")
    force, _ := cmd.Flags().GetBool("force")

    // Verify checksum
    fmt.Println("🔍 Verifying backup integrity...")
    if err := verifyBackup(backupPath); err != nil {
        fmt.Printf("❌ Backup verification failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("✅ Backup integrity verified")

    if verifyOnly {
        os.Exit(0)
    }

    if !force {
        fmt.Print("\n⚠️  This will overwrite existing data. Continue? [y/N] ")
        var response string
        fmt.Scanln(&response)
        if response != "y" && response != "Y" {
            fmt.Println("Restore cancelled")
            os.Exit(0)
        }
    }

    cfg := loadConfig(cmd)

    fmt.Println("\n📦 Restoring backup...")

    file, err := os.Open(backupPath)
    if err != nil {
        fmt.Printf("Error opening backup: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    gzr, err := gzip.NewReader(file)
    if err != nil {
        fmt.Printf("Error reading backup: %v\n", err)
        os.Exit(1)
    }
    defer gzr.Close()

    tr := tar.NewReader(gzr)

    for {
        header, err := tr.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("Error reading backup: %v\n", err)
            os.Exit(1)
        }

        targetPath := mapBackupPath(header.Name, cfg)
        if targetPath == "" {
            continue
        }

        fmt.Printf("  📄 Restoring: %s\n", header.Name)

        os.MkdirAll(filepath.Dir(targetPath), 0750)

        outFile, err := os.Create(targetPath)
        if err != nil {
            fmt.Printf("Error creating file: %v\n", err)
            continue
        }

        io.Copy(outFile, tr)
        outFile.Close()
        os.Chmod(targetPath, os.FileMode(header.Mode))
    }

    fmt.Println("\n✅ Restore complete!")
    fmt.Println("   Please restart gomailserver to apply changes")
}
```

**Acceptance Criteria**:
- [ ] Full system backup
- [ ] Full system restore
- [ ] Checksum verification
- [ ] SQLite WAL mode handling
- [ ] 30-day retention
- [ ] Scheduled backups

---

## 9.6 Database Migrations [MVP]

Structured database schema evolution using golang-migrate for safe, reversible migrations.

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| MIG-001 | Database migration system | [ ] | - | MVP |

---

### MIG-001: Database Migration System

**File**: `internal/database/migrations/README.md`
```markdown
# Database Migrations

gomailserver uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema evolution.

## Migration Files

Migrations are stored in `internal/database/migrations/` with sequential numbering:

```
001_initial_schema.up.sql
001_initial_schema.down.sql
002_add_webhooks.up.sql
002_add_webhooks.down.sql
003_add_calendar.up.sql
003_add_calendar.down.sql
```

## Creating Migrations

```bash
# Using migrate CLI
migrate create -ext sql -dir internal/database/migrations -seq add_feature_name

# Manual creation
touch internal/database/migrations/004_add_feature.up.sql
touch internal/database/migrations/004_add_feature.down.sql
```

## Running Migrations

```bash
# Migrate up to latest
gomailserver migrate up

# Migrate to specific version
gomailserver migrate up 3

# Rollback one version
gomailserver migrate down 1

# Show current version
gomailserver migrate version
```

## Migration Best Practices

1. **Always provide down migrations** for rollback capability
2. **Test migrations on production-like data** before deployment
3. **Make migrations reversible** whenever possible
4. **Avoid data loss** - migrations should preserve data
5. **Keep migrations atomic** - one logical change per migration
6. **Never edit applied migrations** - create a new migration instead
```

**File**: `internal/database/migrate.go`
```go
package database

import (
    "database/sql"
    "embed"
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/sqlite3"
    "github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migrator struct {
    db      *sql.DB
    migrate *migrate.Migrate
}

func NewMigrator(db *sql.DB) (*Migrator, error) {
    // Create migration source from embedded FS
    source, err := iofs.New(migrationsFS, "migrations")
    if err != nil {
        return nil, fmt.Errorf("failed to create migration source: %w", err)
    }

    // Create database driver
    driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to create database driver: %w", err)
    }

    // Create migrator
    m, err := migrate.NewWithInstance("iofs", source, "sqlite3", driver)
    if err != nil {
        return nil, fmt.Errorf("failed to create migrator: %w", err)
    }

    return &Migrator{
        db:      db,
        migrate: m,
    }, nil
}

// Up migrates to the latest version
func (m *Migrator) Up() error {
    if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration up failed: %w", err)
    }
    return nil
}

// Down rolls back one migration
func (m *Migrator) Down() error {
    if err := m.migrate.Down(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration down failed: %w", err)
    }
    return nil
}

// Steps migrates up/down by n steps
func (m *Migrator) Steps(n int) error {
    if err := m.migrate.Steps(n); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration steps failed: %w", err)
    }
    return nil
}

// Version returns the current migration version
func (m *Migrator) Version() (uint, bool, error) {
    version, dirty, err := m.migrate.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return 0, false, fmt.Errorf("failed to get version: %w", err)
    }
    return version, dirty, nil
}

// Force sets the migration version without running migrations
// Use only for recovery from failed migrations
func (m *Migrator) Force(version int) error {
    if err := m.migrate.Force(version); err != nil {
        return fmt.Errorf("force version failed: %w", err)
    }
    return nil
}

// Close closes the migrator
func (m *Migrator) Close() error {
    sourceErr, dbErr := m.migrate.Close()
    if sourceErr != nil {
        return sourceErr
    }
    return dbErr
}
```

**File**: `cmd/gomailserver/migrate.go`
```go
package main

import (
    "fmt"
    "os"
    "strconv"

    "github.com/btafoya/gomailserver/internal/database"
    "github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Database migration commands",
}

var migrateUpCmd = &cobra.Command{
    Use:   "up [steps]",
    Short: "Migrate database up (default: latest)",
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        db, err := database.Open(config.Database.Path)
        if err != nil {
            return fmt.Errorf("failed to open database: %w", err)
        }
        defer db.Close()

        migrator, err := database.NewMigrator(db)
        if err != nil {
            return err
        }
        defer migrator.Close()

        if len(args) > 0 {
            steps, err := strconv.Atoi(args[0])
            if err != nil {
                return fmt.Errorf("invalid steps: %w", err)
            }
            return migrator.Steps(steps)
        }

        return migrator.Up()
    },
}

var migrateDownCmd = &cobra.Command{
    Use:   "down [steps]",
    Short: "Migrate database down (default: 1 step)",
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        db, err := database.Open(config.Database.Path)
        if err != nil {
            return fmt.Errorf("failed to open database: %w", err)
        }
        defer db.Close()

        migrator, err := database.NewMigrator(db)
        if err != nil {
            return err
        }
        defer migrator.Close()

        steps := -1
        if len(args) > 0 {
            n, err := strconv.Atoi(args[0])
            if err != nil {
                return fmt.Errorf("invalid steps: %w", err)
            }
            steps = -n
        }

        return migrator.Steps(steps)
    },
}

var migrateVersionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show current migration version",
    RunE: func(cmd *cobra.Command, args []string) error {
        db, err := database.Open(config.Database.Path)
        if err != nil {
            return fmt.Errorf("failed to open database: %w", err)
        }
        defer db.Close()

        migrator, err := database.NewMigrator(db)
        if err != nil {
            return err
        }
        defer migrator.Close()

        version, dirty, err := migrator.Version()
        if err != nil {
            return err
        }

        if dirty {
            fmt.Printf("Version: %d (dirty - migration failed)\n", version)
            fmt.Println("Run 'gomailserver migrate force <version>' to recover")
        } else {
            fmt.Printf("Version: %d\n", version)
        }

        return nil
    },
}

var migrateForceCmd = &cobra.Command{
    Use:   "force <version>",
    Short: "Force migration version (recovery only)",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        version, err := strconv.Atoi(args[0])
        if err != nil {
            return fmt.Errorf("invalid version: %w", err)
        }

        db, err := database.Open(config.Database.Path)
        if err != nil {
            return fmt.Errorf("failed to open database: %w", err)
        }
        defer db.Close()

        migrator, err := database.NewMigrator(db)
        if err != nil {
            return err
        }
        defer migrator.Close()

        return migrator.Force(version)
    },
}

func init() {
    migrateCmd.AddCommand(migrateUpCmd)
    migrateCmd.AddCommand(migrateDownCmd)
    migrateCmd.AddCommand(migrateVersionCmd)
    migrateCmd.AddCommand(migrateForceCmd)
    rootCmd.AddCommand(migrateCmd)
}
```

**File**: `internal/database/migrations/001_initial_schema.up.sql`
```sql
-- See TASKS0.md for complete schema
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    domain_id INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    active INTEGER NOT NULL DEFAULT 1,
    quota_bytes INTEGER NOT NULL DEFAULT 1073741824,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ... (rest of schema)
```

**File**: `internal/database/migrations/001_initial_schema.down.sql`
```sql
-- Rollback initial schema
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS mailboxes;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS domains;
-- ... (all tables in reverse dependency order)
```

**Acceptance Criteria**:
- [ ] golang-migrate integrated for schema evolution
- [ ] Migrations embedded in binary using `//go:embed`
- [ ] CLI commands: `migrate up`, `migrate down`, `migrate version`, `migrate force`
- [ ] Up and down migrations for all schema changes
- [ ] Migration versioning with sequential numbering
- [ ] Dirty migration detection and recovery with `force` command
- [ ] Migration documentation with best practices

**Production Readiness**:
- [ ] Migration timeout: 30s per migration (prevents hanging deployments)
- [ ] Atomic migrations: Each migration is a single transaction
- [ ] Rollback capability: All migrations have down scripts
- [ ] Version tracking: Current version stored in `schema_migrations` table
- [ ] Dirty state detection: Failed migrations marked as dirty
- [ ] Recovery procedure: `migrate force <version>` for manual recovery
- [ ] Testing: All migrations tested on production-like data before deployment

**Given/When/Then Scenarios**:
```
Given database is at version 0 (empty)
When 'gomailserver migrate up' is run
Then database is migrated to latest version
And all tables are created successfully
And version table shows current version

Given database is at version 5
When 'gomailserver migrate down 1' is run
Then database is rolled back to version 4
And changes from migration 5 are reversed
And data is preserved where possible

Given migration 3 failed mid-execution
When 'gomailserver migrate version' is run
Then version shows as "3 (dirty)"
When admin runs 'gomailserver migrate force 2'
Then version is set to 2
And admin can investigate failed migration
And admin can fix issues before re-running

Given new deployment includes migration 6
When application starts with auto-migrate enabled
Then migration 6 runs automatically
And application starts successfully
And no manual intervention required
```

---

## Testing Checklist

- [ ] Installation script works on Debian/Ubuntu
- [ ] Systemd service starts and stops correctly
- [ ] Config validation catches errors
- [ ] Pre-flight checks run successfully
- [ ] Docker image builds and runs
- [ ] Multi-arch Docker builds work
- [ ] Embedded assets serve correctly
- [ ] Backup creates valid archive
- [ ] Restore works from backup
- [ ] Documentation is accurate and complete
