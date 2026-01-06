# PROJECT KNOWLEDGE BASE

**Generated:** 2026-01-06
**Commit:** $(git rev-parse HEAD 2>/dev/null | cut -c1-8 || echo "unknown")
**Branch:** $(git branch --show-current 2>/dev/null || echo "unknown")

## OVERVIEW
Core business logic layer implementing all mail server operations through 19 specialized services with repository pattern and dependency injection.

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| User authentication | user_service.go | bcrypt hashing, session management |
| Message storage | message_service.go | Hybrid strategy: <1MB in DB, >1MB filesystem |
| Email queuing | queue_service.go | Exponential backoff (5m → 24h) |
| Security features | dane_service.go, mtasts_service.go, pgp_service.go | DKIM/SPF/DMARC, DANE, MTA-STS, PGP |
| Domain management | domain_service.go | CRUD, validation, DNS checks |
| Mailbox operations | mailbox_service.go | IMAP folder management |
| Setup wizard | setup_service.go | First-run configuration |
| Webhooks | webhook_service.go | Event delivery (16 types, HMAC-SHA256) |
| Webmail methods | webmail_methods.go | Rich UI integration (648 lines) |
| Audit logging | audit_service.go | All operations tracking |
| Settings | settings_service.go | System configuration |

## CONVENTIONS
- **Dependency Injection**: All services use `New*Service(repo, logger)` constructors
- **Interface Definition**: Service interfaces in interfaces.go (UserServiceInterface, etc.)
- **Circular Dependencies**: Optional setters (`SetQueueService`, `SetMailboxService`)
- **Context Methods**: API handlers use `ctx context.Context` parameter pattern
- **Custom Errors**: Service-specific errors (ErrInvalidCredentials, ErrUserNotFound)
- **Logging**: Structured zap logging with operation context (user_id, email, etc.)

## ANTI-PATTERNS (THIS PROJECT)
- NEVER call repository directly from handlers - always use service layer
- NEVER skip validation in service methods - all inputs validated before repo calls
- NEVER return repository errors directly - wrap with service context
- NEVER log sensitive data (passwords, tokens, full message bodies)
- NEVER use global state - all dependencies injected via constructor
- NEVER mix storage strategies - MessageService enforces hybrid threshold at 1MB
- NEVER implement business logic in repositories - service layer only
