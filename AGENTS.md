# PROJECT KNOWLEDGE BASE

**Generated:** 2026-01-06
**Commit:** $(git rev-parse HEAD 2>/dev/null | cut -c1-8 || echo "unknown")
**Branch:** $(git branch --show-current 2>/dev/null || echo "unknown")

## OVERVIEW
Go-based mail server replacing Postfix/Dovecot stack. Single daemon with SMTP/IMAP/Webmail/DKIM/SPF/DMARC/DANE/MTA-STS/PGP/antivirus/anti-spam. SQLite backend, unified Nuxt.js web interface (admin/portal/webmail), automated reputation management.

## GUIDELINES

### ❌ Do NOT Include:
- "Generated with Claude Code" in commit messages
- "Co-Authored-By: Claude Sonnet" in commits
- AI attribution in code comments
- References to Claude in documentation footer/header

### ✅ DO Include:
- Your name and email as the commit author
- Professional commit messages describing WHAT changed
- Standard documentation without AI tool references
- Human authorship for all contributions

## STRUCTURE
```
gomailserver/
├── cmd/gomailserver/          # CLI entry point
├── internal/                  # Core business logic
│   ├── api/                   # REST API (handlers, middleware)
│   ├── database/              # SQLite connection/migrations
│   ├── domain/                # Core models
│   ├── reputation/            # Automated sender reputation
│   ├── repository/            # Data access layer
│   ├── security/              # DKIM/SPF/DMARC/ClamAV/SpamAssassin
│   ├── service/               # Business logic layer
│   ├── smtp/                  # SMTP server
│   └── webdav/                # CalDAV/CardDAV
├── unified/                  # Unified web interface (Nuxt.js - admin/portal/webmail)
├── scripts/                  # Build/deployment scripts
└── tests/                    # Integration tests
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add new email feature | internal/service/ | Business logic + repository pattern |
| Add API endpoint | internal/api/handlers/ | REST handlers + middleware |
| Add security feature | internal/security/ | DKIM/SPF/DMARC implementation |
| Add UI feature | unified/ | Nuxt.js pages/components (unified app: admin/portal/webmail) |
| Database changes | internal/database/ + internal/repository/ | Migrations + queries |
| Reputation logic | internal/reputation/ | Automated scoring + rate limiting |
| Email protocols | internal/smtp/ or internal/imap/ | Protocol implementations |
| Build/deploy | scripts/ + Makefile | Control scripts + CI |

## CODE MAP
| Symbol | Type | Location | Refs | Role |
|--------|------|----------|------|------|
| UserService | Struct | internal/service/user_service.go | High | User management |
| QueueService | Struct | internal/service/queue_service.go | High | Email queuing |
| DomainService | Struct | internal/service/domain_service.go | High | Domain management |
| ReputationService | Struct | internal/reputation/service/ | High | Sender reputation |
| DKIMService | Struct | internal/security/dkim/ | Medium | DKIM signing/verification |
| SMTPService | Struct | internal/smtp/ | Medium | SMTP server |
| IMAPService | Struct | internal/imap/ | Medium | IMAP server |

## CONVENTIONS
- **Error handling**: Custom error types with context (ErrInvalidCredentials, etc.)
- **Logging**: Structured logging with zap, context-based
- **Configuration**: YAML-based with viper, environment overrides
- **Database**: SQLite with migrations, hybrid storage (<1MB in DB, >1MB filesystem)
- **Testing**: Unit tests alongside implementation files
- **Security**: bcrypt for passwords, JWT for sessions, TOTP for 2FA

## ANTI-PATTERNS (THIS PROJECT)
- NEVER use global variables - dependency injection only
- NEVER log sensitive data (passwords, tokens)
- NEVER skip input validation - all external data validated
- NEVER use panic() for error handling - return errors
- NEVER hardcode configuration - use config structs

## UNIQUE STYLES
- **Hybrid storage**: Small messages in DB, large messages on disk
- **Reputation-first**: All sending decisions reputation-aware
- **Webmail integrated**: CalDAV/CardDAV alongside IMAP
- **Postmark-compatible**: REST API matching PostmarkApp interface
- **Automated reputation**: Real-time scoring with external feedback

## COMMANDS
```bash
# Development
./scripts/gomailserver-control.sh start --dev    # Start with WebUI
make build                                      # Build binary
make test                                       # Run tests

# Production
./scripts/gomailserver-control.sh start          # Production mode
sudo systemctl start gomailserver               # systemd
docker run gomailserver                         # Container

# Development workflow
go test ./internal/service/...                  # Test services
golangci-lint run                               # Lint code
```

## NOTES
- **Security first**: All mail operations security-hardened by default
- **Reputation critical**: Sending blocked if reputation < threshold
- **Webmail critical**: Core server requires frontend for full functionality
- **Database critical**: SQLite file must be backed up regularly
- **Large files**: Messages >1MB stored on filesystem, not DB</content>
<parameter name="filePath">AGENTS.md