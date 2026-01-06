# API HANDLERS KNOWLEDGE BASE

**Generated:** 2026-01-06
**Purpose:** REST API HTTP handlers implementing mail server management endpoints

## OVERVIEW
18 HTTP handlers implementing REST API endpoints for mail server management, authentication, reputation, webmail, and monitoring.

## STRUCTURE
Flat directory with 18 handler files organized by domain (no subdirectories).

## WHERE TO LOOK
| Task | Handler File | Notes |
|------|---------------|-------|
| Add/modify endpoint | Corresponding handler file | One handler per domain (e.g., user_handler.go, domain_handler.go) |
| Authentication flow | auth_handler.go | Login, JWT refresh, TOTP support |
| Domain management | domain_handler.go | CRUD + DKIM generation |
| User management | user_handler.go | CRUD + password reset + quota management |
| Alias management | alias_handler.go | CRUD for email aliases |
| Queue management | queue_handler.go | View/retry/delete queued messages |
| Statistics/monitoring | stats_handler.go | Dashboard, domain/user stats |
| Log retrieval | log_handler.go | Server log viewer |
| Settings management | settings_handler.go | Server/security/TLS settings |
| PGP key management | pgp_handler.go | Import/list/delete PGP keys |
| Audit logs | audit_handler.go | Audit log viewer and stats |
| Webhook management | webhook_handler.go | CRUD + test + delivery tracking |
| Reputation management | reputation_handler.go | Phase 1-4: audits, scores, circuit breakers, alerts |
| Reputation Phase 5 | reputation_phase5_handler.go | DMARC reports, ARF, external metrics, warmup, predictions |
| Reputation Phase 6 | reputation_phase6_handler.go | Operational mail, deliverability UI endpoints |
| Webmail core | webmail.go | Mailboxes, messages, drafts, search, attachments |
| Webmail contacts | webmail_contacts.go | CardDAV integration for contact autocomplete |
| Webmail calendar | webmail_calendar.go | CalDAV integration for calendar events |
| Setup wizard | setup_handler.go | First-run configuration (no auth required) |

## CONVENTIONS
- **Handler Pattern**: Struct with service + logger dependencies, NewXHandler factory function
- **Response Helpers**: Use middleware.RespondError, RespondJSON, RespondSuccess, RespondCreated, RespondPaginated
- **Auth Context**: Extract authenticated user via middleware.GetUserID(r)
- **Path Parameters**: chi.URLParam(r, "param_name") for route params
- **DTOs**: Handler-specific Request/Response structs with JSON tags (never return domain models)
- **Error Handling**: Log errors with zap.Error, return generic messages (never expose internals)
- **Method Signature**: `func (h *Handler) Method(w http.ResponseWriter, r *http.Request)`
- **Validation**: Validate all inputs before calling services, return 400 for invalid requests
- **Logging**: Structured logging with context (user_id, domain_id, request_id)

## ANTI-PATTERNS (THIS DIRECTORY)
- NEVER expose internal error messages in API responses (use generic messages)
- NEVER skip input validation - always check required fields and types
- NEVER return domain models directly - always use Response DTOs
- NEVER mix concerns - handlers orchestrate services, don't contain business logic
- NEVER forget to call middleware.RespondError on validation failures
- NEVER return sensitive data (passwords, tokens, keys) in responses
- NEVER implement business logic in handlers - delegate to services
- NEVER use panic() - always return errors via middleware.RespondError
- NEVER ignore pagination TODOs - implement when adding list endpoints
- NEVER log request bodies with sensitive data (passwords, tokens)
