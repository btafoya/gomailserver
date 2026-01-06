# gomailserver

[![CI](https://github.com/btafoya/gomailserver/actions/workflows/ci.yml/badge.svg)](https://github.com/btafoya/gomailserver/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/btafoya/gomailserver)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/btafoya/gomailserver)](https://goreportcard.com/report/github.com/btafoya/gomailserver)
[![codecov](https://codecov.io/gh/btafoya/gomailserver/branch/main/graph/badge.svg)](https://codecov.io/gh/btafoya/gomailserver)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](LICENSE.txt)

A modern, composable, all-in-one mail server written in Go 1.23.5+ designed to replace complex mail server stacks (Postfix, Dovecot, OpenDKIM, etc.) with a single daemon. **81% complete** (244/303 tasks) with core mail functionality operational and comprehensive reputation management with advanced automation complete.

Implements SMTP, IMAP, CalDAV, CardDAV with comprehensive email security features including DKIM, SPF, DMARC, DANE, MTA-STS, PGP/GPG, antivirus, and anti-spam capabilities. Features automated reputation management with external feedback integration (Gmail Postmaster Tools, Microsoft SNDS), DMARC report processing, and complete web interfaces with unified admin/portal (Vue.js) and dedicated webmail client (Nuxt.js) including contact/calendar integration.

## Features

### Core Protocols
- **SMTP**: Full RFC 5321 compliance with submission (587), relay (25), and SMTPS (465)
- **IMAP4**: RFC 3501 compliance with extensions (IDLE, UIDPLUS, QUOTA, SORT, THREAD)
- **CalDAV**: RFC 4791 calendar synchronization
- **CardDAV**: RFC 6352 contact synchronization

### Security & Reputation Management
- **DKIM**: Outbound signing and inbound verification (RSA-2048/4096, Ed25519)
- **SPF**: Sender Policy Framework validation
- **DMARC**: Policy enforcement with aggregate/forensic reporting and automated analysis
- **DANE**: DNS-based Authentication of Named Entities
- **MTA-STS**: Strict Transport Security
- **Antivirus**: ClamAV integration
- **Anti-Spam**: SpamAssassin integration
- **Greylisting**: Enabled by default
- **2FA**: TOTP-based two-factor authentication
- **PGP/GPG**: End-to-end encryption support
- **Reputation Telemetry**: Real-time metrics collection and scoring (0-100 scale)
- **External Feedback**: Gmail Postmaster Tools and Microsoft SNDS integration
- **Adaptive Sending**: Reputation-aware rate limiting with circuit breakers
- **Automatic Warm-up**: Progressive volume ramping for new domains/IPs
- **DMARC Processing**: Automated RUA report parsing and issue detection
- **ARF Complaints**: Automatic complaint handling and recipient suppression

### Storage
- **SQLite**: All data in single database file for easy backup
- **Hybrid Storage**: Small messages (< 1MB) in database, large messages on filesystem
- **Unlimited**: Domains, users, aliases with configurable quotas

### Web Interfaces
- **Unified Web Interface**: Single Nuxt.js application with three sections:
  - **Admin** (`/admin/*`): Domain/user/alias management and system administration
  - **Portal** (`/portal/*`): User self-service portal for account management and settings
  - **Webmail** (`/webmail/*`): Gmail-like email interface with rich text composer
- **Setup Wizard**: Guided first-run configuration for system, domain, and admin setup
- **Features**: Unified authentication, responsive design, dark mode support
  - Contact autocomplete and search (CardDAV integration)
  - Calendar widget with upcoming events (CalDAV integration)
  - Event creation and invitation handling
  - Rich text composer with TipTap
  - Dark mode and mobile responsive design

### Advanced Features
- **Sieve Filtering**: Server-side mail filtering (RFC 5228) (planned)
- **Webhooks**: Event notifications for integrations ✅ COMPLETE
  - 16 event types (email.*, security.*, dkim/spf/dmarc/user events)
  - HMAC-SHA256 signed payloads
  - Exponential backoff retry (up to 10 attempts)
  - Delivery tracking and monitoring
- **Auto-configuration**: Let's Encrypt ACME with Cloudflare DNS
- **Multi-domain**: Support for unlimited domains and subdomains
- **PostmarkApp API**: Compatible REST API for drop-in replacement of PostmarkApp services
  - Single email sending (POST /email)
  - Batch email sending (POST /email/batch, up to 500 messages)
  - X-Postmark-Server-Token authentication
  - Test mode support (POSTMARK_API_TEST)
  - Template system (planned)
  - Webhook delivery (planned)
- **Reputation Management**: Complete automated sender reputation system ✅ COMPLETE
  - **Telemetry Foundation**: Real-time metrics (deliveries, bounces, complaints, deferrals)
  - **Reputation Scoring**: 0-100 scale with 90-day rolling window
  - **Deliverability Auditor**: DNS/SPF/DKIM/DMARC/rDNS validation with scoring
  - **Adaptive Policy Engine**: Reputation-aware rate limiting with circuit breakers
  - **Circuit Breakers**: Automatic pause on high complaints (>0.1%), bounces (>10%), or provider blocks
  - **Auto-resume**: Exponential backoff retry (1h → 2h → 4h → 8h)
  - **Progressive Warm-up**: 14-day schedules (100 → 80,000 msgs/day) for new domains/IPs
  - **DMARC Processing**: Automated RUA report parsing, analysis, and actions
  - **ARF Complaints**: Automatic complaint handling and recipient suppression
  - **External Feedback**: Gmail Postmaster Tools and Microsoft SNDS integration
  - **Provider Rate Limits**: Gmail, Outlook, Yahoo-specific sending limits
  - **Custom Warmup**: Conservative/moderate/aggressive templates with progress tracking
  - **AI Predictions**: Trend-based reputation forecasting with confidence levels
  - **Alerts System**: Comprehensive alerts with acknowledgment/resolution workflow
  - **Dashboard UI**: Real-time visualization with 5 comprehensive views (DMARC, metrics, limits, warmup, predictions)

## Quick Start

### Requirements
- Go 1.23.5 or higher (build time only)
- ClamAV daemon (clamd)
- SpamAssassin daemon (spamd)
- Cloudflare account (for automatic TLS certificates)

### Installation

#### Option 1: APT Package (Recommended for Debian/Ubuntu)

Install via APT repository on supported distributions (Ubuntu: focal/jammy/noble, Debian: bullseye/bookworm):

```bash
# Add GPG key
curl -fsSL https://btafoya.github.io/gomailserver/repo/public.key | \
  sudo gpg --dearmor -o /usr/share/keyrings/gomailserver-archive-keyring.gpg

# Add repository (replace 'jammy' with your distribution codename)
echo "deb [signed-by=/usr/share/keyrings/gomailserver-archive-keyring.gpg] https://btafoya.github.io/gomailserver/repo jammy main" | \
  sudo tee /etc/apt/sources.list.d/gomailserver.list

# Install
sudo apt update
sudo apt install gomailserver
```

See [INSTALL-APT.md](INSTALL-APT.md) for detailed APT installation instructions.

#### Option 2: Build from Source

```bash
# Clone repository
git clone https://github.com/btafoya/gomailserver.git
cd gomailserver

# Build
make build

# Or install system-wide
make install
```

#### Option 3: systemd Installation (Manual Build for Production)

For production deployments with systemd:

```bash
# Build the binary
make build

# Install as systemd service (requires root)
sudo ./scripts/install-systemd.sh --start

# The installer will:
# - Create gomailserver user and group
# - Install binary to /usr/local/bin
# - Set up directories (/var/lib/gomailserver, /var/log/gomailserver)
# - Install configuration to /etc/gomailserver
# - Install and enable systemd service
# - Optionally start the service (--start flag)
```

**systemd Service Management:**
```bash
# Start the service
sudo systemctl start gomailserver

# Stop the service
sudo systemctl stop gomailserver

# Restart the service
sudo systemctl restart gomailserver

# Check status
sudo systemctl status gomailserver

# View logs
sudo journalctl -u gomailserver -f

# Enable on boot
sudo systemctl enable gomailserver
```

**Installer Options:**
- `--start` - Start service immediately after installation
- `--enable` - Enable service on boot (default)
- `--no-enable` - Don't enable on boot
- `--user USER` - Run as custom user (default: gomailserver)
- `--group GROUP` - Run as custom group (default: gomailserver)

### Configuration

Copy the example configuration:
```bash
cp gomailserver.example.yaml gomailserver.yaml
```

Edit `gomailserver.yaml` with your settings:
```yaml
server:
  hostname: mail.example.com
  domain: example.com

database:
  path: ./data/mailserver.db

tls:
  acme:
    enabled: true
    email: admin@example.com
    provider: cloudflare
    api_token: your_cloudflare_api_token
```

### First-Time Setup

Create the first admin user before starting the server:

```bash
# Create admin user interactively
./build/gomailserver create-admin

# Follow prompts for email, name, and password
```

### Running

```bash
# Start the server
./build/gomailserver run --config gomailserver.yaml

# Or with default config path
./build/gomailserver run
```

Access the admin UI at `http://localhost:8980/admin/` (or your configured API port).

### Using the Control Script

For easier daemon management, use the provided control script:

```bash
# Start in production mode (uses /etc/gomailserver/gomailserver.yaml)
./scripts/gomailserver-control.sh start

# Start in development mode (uses ./gomailserver.yaml)
./scripts/gomailserver-control.sh start --dev

# Check server status
./scripts/gomailserver-control.sh status

# Stop the server
./scripts/gomailserver-control.sh stop

# Restart the server
./scripts/gomailserver-control.sh restart

# Restart in development mode
./scripts/gomailserver-control.sh restart --dev
```

**Production Mode** (default):
- Uses system configuration at `/etc/gomailserver/gomailserver.yaml`
- Info-level logging for normal operations
- Suitable for deployment environments

**Development Mode** (`--dev` flag):
- Uses local configuration at `./gomailserver.yaml`
- Auto-creates config from example if missing
- Debug-level logging for troubleshooting
- Logs written to `./data/gomailserver.log`

The control script provides:
- PID-based process management (`./data/gomailserver.pid`)
- Graceful shutdown with SIGTERM
- Auto-build if binary is missing
- Status monitoring with port information
- Color-coded output for readability

## CLI Commands

```bash
# Display help
./build/gomailserver --help

# Create first admin user (interactive)
./build/gomailserver create-admin

# Start mail server
./build/gomailserver run [--config path/to/config.yaml]

# Show version information
./build/gomailserver version

# Generate shell completion
./build/gomailserver completion [bash|zsh|fish|powershell]
```

## API Endpoints

The REST API is available at `http://localhost:8980/api/v1/` (default port).

### Setup Wizard (No Authentication Required)
- `GET /api/v1/setup/status` - Check if setup is complete
- `GET /api/v1/setup/state` - Get current wizard state
- `POST /api/v1/setup/admin` - Create first admin user
- `POST /api/v1/setup/complete` - Mark setup as complete

### Authentication
- `POST /api/v1/auth/login` - Login with email/password
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Protected Endpoints (JWT or API Key Required)
- **Domains**: `/api/v1/domains` - CRUD operations for domains
- **Users**: `/api/v1/users` - CRUD operations for users
- **Aliases**: `/api/v1/aliases` - CRUD operations for aliases
- **Queue**: `/api/v1/queue` - View and manage mail queue
- **Statistics**: `/api/v1/stats` - Dashboard and domain/user stats
- **Logs**: `/api/v1/logs` - Server log retrieval
- **Webmail**: `/api/v1/webmail` - Email client endpoints
  - `/mailboxes` - List mailboxes
  - `/mailboxes/{id}/messages` - List messages
  - `/messages/{id}` - Get message details
  - `/messages` - Send email
  - `/messages/{id}/move` - Move message
  - `/messages/{id}/flags` - Update flags
  - `/drafts` - Draft management
  - `/contacts/search` - Contact autocomplete (CardDAV)
  - `/contacts/addressbooks` - List addressbooks
  - `/calendar/calendars` - List calendars (CalDAV)
  - `/calendar/upcoming` - Get upcoming events
  - `/calendar/events` - Create events
  - `/calendar/invitations` - Process meeting invitations

### PostmarkApp-Compatible Endpoints (X-Postmark-Server-Token Required)
- `POST /email` - Send single email
- `POST /email/batch` - Send up to 500 emails in batch
- `GET /templates` - Template listing (placeholder)
- `GET /webhooks` - Webhook listing (placeholder)
- `GET /server` - Server information (placeholder)

## Development

### Build Commands

```bash
# Build binary
make build

# Build static binary (for Docker)
make build-static

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Clean build artifacts
make clean

# Update dependencies
make deps
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Integration tests
cd tests/
./run.sh
```

### Docker

```bash
# Build image
make docker-build

# Run container
make docker-run

# Or manually
docker build -t gomailserver .
docker run -p 25:25 -p 143:143 -p 465:465 -p 587:587 -p 993:993 \
  -v gomailserver-data:/data \
  gomailserver:latest
```

## Project Structure

```
gomailserver/
├── cmd/
│   └── gomailserver/          # Main application entrypoint
├── internal/
│   ├── admin/                 # Admin UI handler (embedded Vue.js app)
│   ├── api/                   # REST API (handlers, middleware, router)
│   ├── caldav/                # CalDAV server (RFC 4791)
│   ├── carddav/               # CardDAV server (RFC 6352)
│   ├── commands/              # CLI commands (run, create-admin, version)
│   ├── config/                # Configuration management
│   ├── database/              # SQLite connection and migrations
│   ├── domain/                # Domain models
│   ├── imap/                  # IMAP server
│   ├── postmark/              # PostmarkApp-compatible API
│   │   ├── handlers/          # Email sending handlers
│   │   ├── middleware/        # Authentication middleware
│   │   ├── models/            # Request/response models
│   │   ├── repository/        # PostmarkApp data layer
│   │   └── service/           # Email business logic
│   ├── repository/            # Data access layer
│   │   ├── calendar/          # Calendar/Event repositories
│   │   ├── contact/           # Addressbook/Contact repositories
│   │   └── sqlite/            # SQLite implementations
│   ├── security/              # DKIM, SPF, DMARC, ClamAV, SpamAssassin
│   ├── service/               # Business logic (User, Domain, Queue, Setup, etc.)
│   ├── smtp/                  # SMTP server
│   └── webdav/                # WebDAV server (CalDAV/CardDAV integration)
├── pkg/
│   └── sieve/                 # Sieve interpreter (future)
├── unified/                  # Unified web interface (Nuxt.js 3 - admin/portal/webmail)
├── unified-go/                # Embedded frontend assets
├── tests/                     # Integration tests
├── Makefile                   # Build automation
├── Dockerfile                 # Docker container
└── README.md                  # This file
```

## Configuration Reference

See `gomailserver.example.yaml` for a complete configuration example with comments.

## Database

All configuration and metadata stored in SQLite:
- **File**: Single `mailserver.db` file
- **WAL Mode**: Enabled for better concurrency
- **Migrations**: Automatic schema management (currently V5)
  - V1: Core mail server tables
  - V2: Security and authentication enhancements
  - V3: CalDAV/CardDAV tables
  - V4: Settings and configuration
  - V5: PostmarkApp API tables (servers, messages, templates, webhooks, bounces, events)
- **Backup**: Simple file copy or built-in backup command

### Backup

```bash
# Manual backup
cp ./data/mailserver.db ./backups/mailserver-$(date +%Y%m%d).db

# Built-in backup (future)
./build/gomailserver backup
```

## Technical Stack

### Backend
- **Language**: Go 1.23.5+
- **Web Framework**: Chi Router v5
- **Database**: SQLite (hybrid message storage)
- **MIME**: emersion/go-message/mail
- **Authentication**: JWT with bcrypt

### Frontend
- **Unified Web Interface**: Nuxt.js 3 application serving admin, portal, and webmail under `/admin/*`, `/portal/*`, `/webmail/*`
  - Vue 3.5.26, Nuxt UI, Tailwind CSS 4.1.18, Pinia 3.0.4
  - Admin: Domain/user/alias management with modern dashboard
  - Portal: User self-service with account management and settings
  - Webmail: Gmail-like email interface with rich text composer
- **Features**: Dark mode, responsive design, modern UI components, unified authentication

### DevOps
- **Package Manager**: pnpm (frontend)
- **Build System**: Go modules, Nuxt build
- **Deployment**: Single 21MB binary
- **Assets**: Embedded with go:embed

## Contributing

Contributions are welcome! This is a greenfield project following the PR.md requirements. See CLAUDE.md for autonomous work guidelines.

### GitHub Resources

- **Issues**: [Report bugs or request features](https://github.com/btafoya/gomailserver/issues)
- **Pull Requests**: [Submit changes](https://github.com/btafoya/gomailserver/pulls)
- **Discussions**: [Ask questions or share ideas](https://github.com/btafoya/gomailserver/discussions)
- **CI/CD**: [GitHub Actions workflows](https://github.com/btafoya/gomailserver/actions)

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make changes following existing patterns
4. Add tests for new functionality
5. Run linter and tests (`make lint && make test`)
6. Commit with descriptive messages
7. Push to your fork (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Standards

- Go 1.23.5+ idiomatic code
- 80%+ test coverage required
- golangci-lint must pass (see `.golangci.yml`)
- Comprehensive error handling
- Context-based cancellation
- Structured logging with zap
- All CI checks must pass before merge

## Project Status

**Current Phase**: Webhooks Complete (Phase 8), Reputation Management Phase 5 In Progress
**Overall Progress**: 81% (244/303 tasks)
**Build Status**: ✅ Passing (21MB binary with embedded UI)
**Test Status**: ⚠️ Partial (ACME build failures, IMAP tests passing)
**Production Ready**: ❌ Not yet (requires testing and security audit)

### Reputation Management Status
- **Phase 1-4**: ✅ Complete (Telemetry, Auditor, Adaptive Sending, Dashboard)
- **Phase 5**: 🔧 85% Complete (All services + repositories implemented, integration pending)
- **Overall**: Operational with automated reputation scoring, circuit breakers, and warm-up
- **External APIs**: Ready for Gmail Postmaster Tools and Microsoft SNDS integration

### Known Issues
1. **ACME Service Build Failures** (Priority: High) - Let's Encrypt automatic certificate renewal may be broken
2. **Webmail Send Integration** (Priority: Medium) - MIME building complete, needs QueueService integration
3. **Draft Folder Integration** (Priority: Low) - Drafts saved to database but not visible in Drafts folder

### Next Steps (Priority Order)
**Critical** (Blocks Production):
1. Fix ACME Service - Resolve build failures for Let's Encrypt
2. Integration Testing - E2E tests for mail flow
3. Security Audit - Review authentication, TLS, SQL injection
4. Documentation - Admin guide, user guide, API reference

**High Priority** (MVP Features):
5. Queue Integration - Connect webmail SendMessage to SMTP queue
6. Draft Storage - Integrate draft management with mailbox system
7. Search Enhancement - Implement full-text search index
8. Performance Testing - Benchmark 100K emails/day throughput

## Documentation

### Project Status & Planning
- **PROJECT-STATUS-2026-01-04.md** - Comprehensive project status with complete phase breakdown and 81% completion tracking (244/303 tasks)
- **README.md** - Project overview and quick start (this file)

### Feature Documentation
- **REPUTATION-MANAGEMENT.md** - Complete reputation management strategy and architecture

### Development Guidelines
- **CLAUDE.md** - Development guidelines for autonomous work
- **PR.md** - Pull request guidelines and requirements
- **.doc_archive/** - Historical documentation and phase completion files (60+ archived documents)

## Repository Information

- **GitHub**: [btafoya/gomailserver](https://github.com/btafoya/gomailserver)
- **Issues**: [Bug Reports & Features](https://github.com/btafoya/gomailserver/issues)
- **Documentation**: [Project Wiki](https://github.com/btafoya/gomailserver/wiki)
- **CI/CD**: Automated testing via [GitHub Actions](https://github.com/btafoya/gomailserver/actions)
- **Code Quality**: Monitored via [Go Report Card](https://goreportcard.com/report/github.com/btafoya/gomailserver)
- **Coverage**: Tracked via [Codecov](https://codecov.io/gh/btafoya/gomailserver)

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

### What This Means

- ✅ **Free to Use**: You can use this software for any purpose, including commercial use
- ✅ **Free to Modify**: You can modify the source code to suit your needs
- ✅ **Free to Distribute**: You can distribute the original or modified versions
- ⚖️ **Network Use Requirement**: If you run a modified version on a server and let users interact with it over a network, you **must** provide them access to the modified source code
- 📝 **Share Alike**: Modifications must also be licensed under AGPL-3.0
- 🔓 **Source Code**: You must make source code available when you distribute the software

### Key AGPL-3.0 Provision

The AGPL-3.0 license includes a network copyleft provision (Section 13): if you modify this software and run it on a server where users can interact with it remotely (e.g., as a mail server service), you must offer those users access to the modified source code.

This ensures that improvements to mail server software benefit the entire community, even when the software is used to provide network services.

### Full License

See [LICENSE.txt](LICENSE.txt) for the complete GNU Affero General Public License v3.0 text.

For more information about AGPL-3.0, visit: https://www.gnu.org/licenses/agpl-3.0.html

## Author

**Brian Tafoya** - [GitHub Profile](https://github.com/btafoya)

## Acknowledgments

Built with excellent open-source libraries:
- [emersion/go-smtp](https://github.com/emersion/go-smtp) - SMTP server implementation
- [emersion/go-imap](https://github.com/emersion/go-imap) - IMAP server implementation
- [emersion/go-message](https://github.com/emersion/go-message) - MIME message parsing
- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [spf13/viper](https://github.com/spf13/viper) - Configuration management
- [uber-go/zap](https://github.com/uber-go/zap) - Structured logging
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite database driver

Special thanks to all contributors and the Go community!
