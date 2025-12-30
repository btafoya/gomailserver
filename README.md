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

### Running

```bash
# Start the server
./build/gomailserver run --config gomailserver.yaml

# Or with default config path
./build/gomailserver run
```

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
│   ├── commands/              # CLI commands
│   ├── config/                # Configuration management
│   ├── database/              # SQLite connection and migrations
│   ├── domain/                # Domain models
│   ├── repository/            # Data access layer
│   ├── service/               # Business logic
│   ├── smtp/                  # SMTP server
│   ├── imap/                  # IMAP server
│   ├── caldav/                # CalDAV server
│   ├── carddav/               # CardDAV server
│   ├── security/              # DKIM, SPF, DMARC, etc.
│   ├── api/                   # REST API
│   └── webmail/               # Webmail backend
├── pkg/
│   └── sieve/                 # Sieve interpreter
├── web/
│   ├── admin/                 # Admin UI (Vue.js)
│   ├── portal/                # User portal (Vue.js)
│   └── webmail/               # Webmail client (Vue.js)
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

### Phase 1: Core Mail Server (In Progress)
- [ ] SMTP server implementation
- [ ] IMAP server implementation
- [ ] Message storage and handling
- [ ] User authentication
- [ ] TLS support

### Phase 2: Security Foundation
- [ ] DKIM signing and verification
- [ ] SPF validation
- [ ] DMARC enforcement
- [ ] ClamAV integration
- [ ] SpamAssassin integration
- [ ] Greylisting

### Phase 3: Web Interfaces
- [ ] REST API
- [ ] Admin web UI
- [ ] User self-service portal
- [ ] Let's Encrypt integration
- [ ] Setup wizard

### Future Phases
- CalDAV/CardDAV servers
- Advanced security (DANE, MTA-STS, PGP)
- Sieve filtering
- Webmail client
- Webhooks

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
