# gomailserver

[![CI](https://github.com/btafoya/gomailserver/actions/workflows/ci.yml/badge.svg)](https://github.com/btafoya/gomailserver/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/btafoya/gomailserver)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/btafoya/gomailserver)](https://goreportcard.com/report/github.com/btafoya/gomailserver)
[![codecov](https://codecov.io/gh/btafoya/gomailserver/branch/main/graph/badge.svg)](https://codecov.io/gh/btafoya/gomailserver)
[![License](https://img.shields.io/badge/license-TBD-blue.svg)](LICENSE)

A modern, composable, all-in-one mail server written in Go 1.23.5+. Implements SMTP, IMAP, CalDAV, CardDAV with comprehensive email security features including DKIM, SPF, DMARC, antivirus, and anti-spam capabilities.

## Features

### Core Protocols
- **SMTP**: Full RFC 5321 compliance with submission (587), relay (25), and SMTPS (465)
- **IMAP4**: RFC 3501 compliance with extensions (IDLE, UIDPLUS, QUOTA, SORT, THREAD)
- **CalDAV**: RFC 4791 calendar synchronization
- **CardDAV**: RFC 6352 contact synchronization

### Security
- **DKIM**: Outbound signing and inbound verification (RSA-2048/4096, Ed25519)
- **SPF**: Sender Policy Framework validation
- **DMARC**: Policy enforcement with aggregate/forensic reporting
- **DANE**: DNS-based Authentication of Named Entities
- **MTA-STS**: Strict Transport Security
- **Antivirus**: ClamAV integration
- **Anti-Spam**: SpamAssassin integration
- **Greylisting**: Enabled by default
- **2FA**: TOTP-based two-factor authentication
- **PGP/GPG**: End-to-end encryption support

### Storage
- **SQLite**: All data in single database file for easy backup
- **Hybrid Storage**: Small messages (< 1MB) in database, large messages on filesystem
- **Unlimited**: Domains, users, aliases with configurable quotas

### Web Interfaces
- **Admin UI**: Modern web interface for domain/user/alias management
- **User Portal**: Self-service portal for password changes, 2FA, filters
- **Webmail**: Gmail-like interface with categories, conversation view, PGP integration

### Advanced Features
- **Sieve Filtering**: Server-side mail filtering (RFC 5228)
- **Webhooks**: Event notifications for integrations
- **Auto-configuration**: Let's Encrypt ACME with Cloudflare DNS
- **Multi-domain**: Support for unlimited domains and subdomains

## Quick Start

### Requirements
- Go 1.23.5 or higher (build time only)
- ClamAV daemon (clamd)
- SpamAssassin daemon (spamd)
- Cloudflare account (for automatic TLS certificates)

### Installation

```bash
# Clone repository
git clone https://github.com/btafoya/gomailserver.git
cd gomailserver

# Build
make build

# Or install system-wide
make install
```

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
├── web/
│   ├── admin/                 # Admin UI (Vue.js, Vite, TailwindCSS)
│   ├── portal/                # User portal (future)
│   └── webmail/               # Webmail client (future)
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
- **Migrations**: Automatic schema management
- **Backup**: Simple file copy or built-in backup command

### Backup

```bash
# Manual backup
cp ./data/mailserver.db ./backups/mailserver-$(date +%Y%m%d).db

# Built-in backup (future)
./build/gomailserver backup
```

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

## Roadmap

### Phase 0: Foundation ✅
- [x] Go module initialization
- [x] Package structure
- [x] Development tooling
- [x] Logging and configuration
- [x] Database foundation

### Phase 1: Core Mail Server ✅
- [x] SMTP server implementation (ports 25, 587, 465)
- [x] IMAP server implementation (ports 143, 993)
- [x] Message storage and handling (hybrid blob/filesystem)
- [x] User authentication (bcrypt with PLAIN auth)
- [x] TLS support (STARTTLS + implicit TLS)
- [x] Service layer (UserService, MessageService, QueueService, MailboxService)
- [x] SQLite repositories (5 repositories with full CRUD)
- [x] Unit tests (58 tests, 55 passing, 3 skipped)

### Phase 2: Security Foundation ✅
- [x] DKIM signing and verification (RSA-2048/4096, Ed25519)
- [x] SPF validation with configurable DNS servers
- [x] DMARC enforcement with policy validation
- [x] ClamAV integration (virus scanning)
- [x] SpamAssassin integration (spam filtering)
- [x] Greylisting with configurable delays
- [x] Rate limiting (SMTP, IMAP, Authentication)
- [x] Brute force protection with IP blacklisting
- [x] Per-domain security configuration
- [x] Default security template system

### Phase 3: REST API & Admin Interface ✅
- [x] REST API with JWT authentication
- [x] API key authentication support
- [x] Domain management endpoints
- [x] User management endpoints
- [x] Alias management endpoints
- [x] Queue management endpoints
- [x] Statistics and monitoring endpoints
- [x] Admin web UI (Vue.js embedded)
- [x] Setup wizard (CLI + API endpoints)
- [x] Admin user creation (`create-admin` command)
- [x] Role-based access control (admin/user)
- [ ] User self-service portal
- [ ] Let's Encrypt ACME integration

### Phase 4: CalDAV/CardDAV ✅
- [x] CalDAV server (RFC 4791) with HTTP Basic Auth
- [x] CardDAV server (RFC 6352) with HTTP Basic Auth
- [x] Calendar service (create, update, delete, list)
- [x] Event service with timezone support
- [x] Addressbook service (create, update, delete, list)
- [x] Contact service with vCard 4.0 support
- [x] WebDAV server integration
- [x] Comprehensive test coverage

### Future Phases
- Advanced TLS (DANE, MTA-STS)
- PGP/GPG encryption support
- Sieve filtering (RFC 5228)
- Webmail client
- Webhooks and event notifications
- Integration tests (SMTP/IMAP end-to-end)
- Manual testing with real mail clients

## Repository Information

- **GitHub**: [btafoya/gomailserver](https://github.com/btafoya/gomailserver)
- **Issues**: [Bug Reports & Features](https://github.com/btafoya/gomailserver/issues)
- **Documentation**: [Project Wiki](https://github.com/btafoya/gomailserver/wiki)
- **CI/CD**: Automated testing via [GitHub Actions](https://github.com/btafoya/gomailserver/actions)
- **Code Quality**: Monitored via [Go Report Card](https://goreportcard.com/report/github.com/btafoya/gomailserver)
- **Coverage**: Tracked via [Codecov](https://codecov.io/gh/btafoya/gomailserver)

## License

[License TBD]

## Author

**btafoya** - [GitHub Profile](https://github.com/btafoya)

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
