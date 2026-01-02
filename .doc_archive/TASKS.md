# gomailserver - Detailed Task List

**Last Updated**: 2026-01-02
**Current Status**: Phase 5 (Advanced Security) & Phase 7 (Webmail) Complete - 232/303 tasks done (77%)
**Recent Achievement**: Phase 5 Advanced Security complete (DANE, MTA-STS, PGP/GPG, Audit Logging) with full shadcn-vue UI integration

Based on PR.md requirements. This is a greenfield project - building from scratch.

---

## Task Priority Legend
- **[MVP]** - Required for Minimum Viable Product (Phases 1-3)
- **[FULL]** - Full feature set
- **[OPT]** - Optional/Nice-to-have

## Status Legend
- `[ ]` Not started
- `[~]` In progress
- `[x]` Complete

---

## Phase 0: Foundation (Week 0)

### 0.1 Project Setup [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| F-001 | Initialize Go module (`github.com/btafoya/gomailserver`) | [x] | - |
| F-002 | Create package structure (clean architecture) | [x] | F-001 |
| F-003 | Set up golangci-lint configuration | [x] | F-001 |
| F-004 | Create Makefile for common tasks | [x] | F-001 |
| F-005 | Set up GitHub Actions CI/CD | [x] | F-003 |

### 0.2 Core Infrastructure [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| F-010 | Implement structured logging (JSON format) | [x] | F-002 |
| F-011 | Create configuration system (YAML + env vars) | [x] | F-002 |
| F-012 | Implement CLI framework (cobra) | [x] | F-002 |
| F-013 | Create graceful shutdown handler | [x] | F-010 |
| F-014 | Implement context-based cancellation | [x] | F-013 |

### 0.3 Database Foundation [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| F-020 | SQLite connection management with WAL mode | [x] | F-002 |
| F-021 | Database migration framework | [x] | F-020 |
| F-022 | Create schema version 1 (all tables) | [x] | F-021 |
| F-023 | Implement repository pattern interfaces | [x] | F-020 |
| F-024 | SQLite PRAGMA optimizations | [x] | F-020 |

**Package Structure:**
```
cmd/
  gomailserver/
    main.go
internal/
  config/         # Configuration management
  database/       # SQLite connection, migrations
  domain/         # Domain models
  repository/     # Data access layer
  service/        # Business logic
  smtp/           # SMTP server
  imap/           # IMAP server
  caldav/         # CalDAV server
  carddav/        # CardDAV server
  security/       # DKIM, SPF, DMARC, etc.
  api/            # REST API
  webmail/        # Webmail backend
pkg/
  sieve/          # Sieve interpreter (if custom)
web/
  admin/          # Admin UI assets
  portal/         # User portal assets
  webmail/        # Webmail client assets
```

---

## Phase 1: Core Mail Server (Weeks 1-4)

### 1.1 SMTP Server [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| S-001 | Integrate go-smtp library | [ ] | F-002 | emersion/go-smtp |
| S-002 | Implement SMTP submission server (port 587) | [ ] | S-001 |
| S-003 | Implement SMTP relay server (port 25) | [ ] | S-001 |
| S-004 | Implement SMTPS server (port 465) | [ ] | S-001, T-001 |
| S-005 | STARTTLS support | [ ] | S-002, T-001 |
| S-006 | PLAIN authentication mechanism | [ ] | S-002 |
| S-007 | LOGIN authentication mechanism | [ ] | S-002 |
| S-008 | CRAM-MD5 authentication mechanism | [ ] | S-002 |
| S-009 | SIZE extension (RFC 1870) | [ ] | S-001 |
| S-010 | 8BITMIME support (RFC 6152) | [ ] | S-001 |
| S-011 | PIPELINING support (RFC 2920) | [ ] | S-001 |
| S-012 | CHUNKING support (RFC 3030) | [ ] | S-001 |

### 1.2 SMTP Queue Management [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| Q-001 | Design queue table schema | [ ] | F-022 |
| Q-002 | Implement message queuing service | [ ] | Q-001 |
| Q-003 | Implement retry logic with exponential backoff | [ ] | Q-002 |
| Q-004 | Implement bounce handling | [ ] | Q-002 |
| Q-005 | DSN (Delivery Status Notifications) | [ ] | Q-002 |
| Q-006 | Queue cleanup and maintenance | [ ] | Q-002 |

### 1.3 Message Storage [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| M-001 | Integrate go-message for MIME parsing | [ ] | F-002 | emersion/go-message |
| M-002 | Implement hybrid storage (blob < 1MB, file >= 1MB) | [ ] | F-022 |
| M-003 | Message header parsing and indexing | [ ] | M-001 |
| M-004 | Attachment handling | [ ] | M-001 |
| M-005 | Message deduplication | [ ] | M-003 |
| M-006 | Thread ID generation for conversations | [ ] | M-003 |

### 1.4 IMAP Server [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| I-001 | Integrate go-imap library | [ ] | F-002 | emersion/go-imap |
| I-002 | Implement IMAP backend interface | [ ] | I-001, M-002 |
| I-003 | Mailbox operations (CREATE, DELETE, RENAME) | [ ] | I-002 |
| I-004 | Message operations (FETCH, STORE, COPY) | [ ] | I-002 |
| I-005 | STARTTLS support | [ ] | I-001, T-001 |
| I-006 | PLAIN authentication | [ ] | I-001 |
| I-007 | LOGIN authentication | [ ] | I-001 |
| I-008 | IDLE support (RFC 2177) | [ ] | I-002 |
| I-009 | UIDPLUS extension (RFC 4315) | [ ] | I-002 |
| I-010 | QUOTA extension (RFC 2087) | [ ] | I-002 |
| I-011 | SORT extension (RFC 5256) | [ ] | I-002 |
| I-012 | NAMESPACE extension (RFC 2342) | [ ] | I-002 |
| I-013 | Special-use mailboxes (RFC 6154) | [ ] | I-002 |
| I-014 | Server-side search | [ ] | I-002 |

### 1.5 User Management [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| U-001 | Users table schema and repository | [ ] | F-022 |
| U-002 | Bcrypt password hashing | [ ] | U-001 |
| U-003 | Domains table schema and repository | [ ] | F-022 |
| U-004 | Aliases table schema and repository | [ ] | F-022 |
| U-005 | Mailboxes/folders table and repository | [ ] | F-022 |
| U-006 | Quota tracking and enforcement | [ ] | U-001, M-002 |

### 1.6 TLS Support [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| T-001 | TLS configuration loader | [ ] | F-011 |
| T-002 | SNI support for multi-domain | [ ] | T-001 |
| T-003 | TLS 1.2+ enforcement | [ ] | T-001 |
| T-004 | Modern cipher suite configuration | [ ] | T-001 |

---

## Phase 2: Security Foundation (Weeks 5-7)

### 2.1 DKIM [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| DK-001 | RSA-2048/4096 key generation | [ ] | F-002 |
| DK-002 | Ed25519 key generation (RFC 8463) | [ ] | F-002 |
| DK-003 | DKIM signing for outbound mail | [ ] | DK-001, S-002 |
| DK-004 | DKIM verification for inbound mail | [ ] | DK-001, S-003 |
| DK-005 | Multiple selector support per domain | [ ] | DK-001 |
| DK-006 | Key rotation mechanism | [ ] | DK-005 |
| DK-007 | DKIM keys storage (per domain) | [ ] | F-022, U-003 |

### 2.2 SPF [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| SPF-001 | SPF record parsing | [ ] | F-002 | miekg/dns |
| SPF-002 | SPF validation for inbound mail | [ ] | SPF-001, S-003 |
| SPF-003 | SPF result headers | [ ] | SPF-002 |
| SPF-004 | IPv4/IPv6 support | [ ] | SPF-001 |
| SPF-005 | Configurable handling (none/softfail/fail) | [ ] | SPF-002 |

### 2.3 DMARC [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| DM-001 | DMARC record parsing | [ ] | SPF-001 |
| DM-002 | DMARC policy enforcement | [ ] | DM-001, DK-004, SPF-002 |
| DM-003 | Alignment checking (relaxed/strict) | [ ] | DM-002 |
| DM-004 | Aggregate report generation | [ ] | DM-002 |
| DM-005 | Forensic report generation | [ ] | DM-002 |
| DM-006 | Report sending scheduler | [ ] | DM-004, Q-002 |

### 2.4 Anti-Virus (ClamAV) [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| AV-001 | ClamAV socket connection | [ ] | F-002 |
| AV-002 | Message scanning integration | [ ] | AV-001, M-001 |
| AV-003 | Configurable actions (reject/quarantine/tag) | [ ] | AV-002 |
| AV-004 | Per-domain/user configuration | [ ] | AV-003, U-003 |
| AV-005 | Scan result logging | [ ] | AV-002, F-010 |

### 2.5 Anti-Spam (SpamAssassin) [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| AS-001 | Integrate spamc client | [ ] | F-002 | teamwork/spamc |
| AS-002 | Message scoring integration | [ ] | AS-001, M-001 |
| AS-003 | Per-user spam threshold | [ ] | AS-002, U-001 |
| AS-004 | Spam quarantine system | [ ] | AS-002, F-022 |
| AS-005 | Learn from user actions (spam/ham) | [ ] | AS-001 |
| AS-006 | Spam report generation | [ ] | AS-002 |

### 2.6 Greylisting [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| GL-001 | Greylisting table schema | [ ] | F-022 |
| GL-002 | Implement greylist check | [ ] | GL-001, S-003 |
| GL-003 | Auto-whitelist after pass | [ ] | GL-002 |
| GL-004 | Configurable delay and expiry | [ ] | GL-002 |

### 2.7 Rate Limiting [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| RL-001 | Rate limit table schema | [ ] | F-022 |
| RL-002 | Per-IP rate limiting | [ ] | RL-001 |
| RL-003 | Per-user rate limiting | [ ] | RL-001 |
| RL-004 | Per-domain rate limiting | [ ] | RL-001 |
| RL-005 | Configurable limits | [ ] | RL-002 |

### 2.8 Authentication Security [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| AU-001 | TOTP 2FA implementation | [ ] | U-001 | pquerna/otp |
| AU-002 | Failed login tracking | [ ] | F-022 |
| AU-003 | IP blacklisting | [ ] | AU-002 |
| AU-004 | IP whitelisting | [ ] | F-022 |
| AU-005 | Brute force protection | [ ] | AU-002 |

---

## Phase 3: Web Interfaces (Weeks 8-12)

### 3.1 REST API Foundation [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| API-001 | Set up Echo web framework | [ ] | F-002 | labstack/echo |
| API-002 | JWT authentication middleware | [ ] | API-001 | golang-jwt/jwt |
| API-003 | API key authentication | [ ] | API-002 |
| API-004 | Request rate limiting middleware | [ ] | API-001 |
| API-005 | CORS configuration | [ ] | API-001 |
| API-006 | OpenAPI/Swagger documentation | [ ] | API-001 |
| API-007 | Request validation middleware | [ ] | API-001 |

### 3.2 Admin API Endpoints [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| AA-001 | Domain CRUD endpoints | [ ] | API-002, U-003 |
| AA-002 | User CRUD endpoints | [ ] | API-002, U-001 |
| AA-003 | Alias CRUD endpoints | [ ] | API-002, U-004 |
| AA-004 | Quota management endpoints | [ ] | API-002, U-006 |
| AA-005 | Statistics endpoints | [ ] | API-002 |
| AA-006 | Log retrieval endpoints | [ ] | API-002, F-010 |
| AA-007 | Queue management endpoints | [ ] | API-002, Q-002 |
| AA-008 | DKIM key management endpoints | [ ] | API-002, DK-007 |
| AA-009 | System health endpoints | [ ] | API-002 |
| AA-010 | Backup/restore endpoints | [ ] | API-002 |

### 3.3 Admin Web UI [MVP]
| ID | Task | Status | Dependencies | Framework |
|----|------|--------|--------------|-----------|
| AUI-001 | Set up Vue.js 3 project with Vite | [ ] | - | Vue 3 |
| AUI-002 | Admin authentication flow | [ ] | AUI-001, API-002 |
| AUI-003 | Domain management UI | [ ] | AUI-002, AA-001 |
| AUI-004 | User management UI | [ ] | AUI-002, AA-002 |
| AUI-005 | Alias management UI | [ ] | AUI-002, AA-003 |
| AUI-006 | Quota visualization | [ ] | AUI-002, AA-004 |
| AUI-007 | Real-time statistics dashboard | [ ] | AUI-002, AA-005 |
| AUI-008 | Log viewer with filtering | [ ] | AUI-002, AA-006 |
| AUI-009 | Queue management interface | [ ] | AUI-002, AA-007 |
| AUI-010 | DKIM/SPF/DMARC settings per domain | [ ] | AUI-002, AA-008 |
| AUI-011 | TLS certificate status | [ ] | AUI-002, LE-004 |
| AUI-012 | System health monitoring | [ ] | AUI-002, AA-009 |
| AUI-013 | Role-based access control | [ ] | AUI-002 |

### 3.4 User Self-Service Portal [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| USP-001 | User authentication API | [ ] | API-002 |
| USP-002 | Password change endpoint | [ ] | USP-001 |
| USP-003 | 2FA setup endpoint | [ ] | USP-001, AU-001 |
| USP-004 | User alias management | [ ] | USP-001 |
| USP-005 | Quota usage display | [ ] | USP-001, U-006 |
| USP-006 | Forwarding rules API | [ ] | USP-001 |
| USP-007 | Session management | [ ] | USP-001, F-022 |

### 3.5 User Portal UI [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| UP-001 | User portal Vue.js project | [ ] | AUI-001 |
| UP-002 | Password change UI | [ ] | UP-001, USP-002 |
| UP-003 | 2FA setup wizard | [ ] | UP-001, USP-003 |
| UP-004 | Alias management UI | [ ] | UP-001, USP-004 |
| UP-005 | Quota display widget | [ ] | UP-001, USP-005 |
| UP-006 | Forwarding rules editor | [ ] | UP-001, USP-006 |
| UP-007 | Spam quarantine viewer | [ ] | UP-001, AS-004 |

### 3.6 Let's Encrypt Integration [MVP]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| LE-001 | ACME client integration | [ ] | F-002 | go-acme/lego |
| LE-002 | Cloudflare DNS challenge | [ ] | LE-001 |
| LE-003 | Automatic certificate renewal | [ ] | LE-002 |
| LE-004 | Certificate storage and loading | [ ] | LE-001, T-001 |
| LE-005 | Per-domain certificate support | [ ] | LE-004 |

### 3.7 Setup Wizard [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| SW-001 | First-run detection | [ ] | F-020 |
| SW-002 | First domain configuration | [ ] | SW-001, U-003 |
| SW-003 | First admin user creation | [ ] | SW-002, U-001 |
| SW-004 | TLS certificate setup flow | [ ] | SW-002, LE-001 |
| SW-005 | DKIM key generation UI | [ ] | SW-002, DK-001 |
| SW-006 | DNS record suggestions | [ ] | SW-002 |
| SW-007 | Pre-flight checks (ports, services) | [ ] | SW-001 |

---

## Phase 4: CalDAV/CardDAV (Weeks 13-15)

### 4.1 WebDAV Foundation [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| WD-001 | WebDAV base protocol (RFC 4918) | [ ] | API-001 |
| WD-002 | PROPFIND implementation | [ ] | WD-001 |
| WD-003 | PROPPATCH implementation | [ ] | WD-001 |
| WD-004 | MKCOL implementation | [ ] | WD-001 |
| WD-005 | DELETE implementation | [ ] | WD-001 |
| WD-006 | COPY/MOVE implementation | [ ] | WD-001 |

### 4.2 CalDAV Server [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| CD-001 | CalDAV protocol (RFC 4791) | [ ] | WD-001 |
| CD-002 | Calendar collection management | [ ] | CD-001 |
| CD-003 | Event storage (RFC 5545 iCalendar) | [ ] | CD-002 |
| CD-004 | REPORT method implementation | [ ] | CD-001 |
| CD-005 | Calendar-query support | [ ] | CD-004 |
| CD-006 | Recurring events handling | [ ] | CD-003 |
| CD-007 | Event reminders | [ ] | CD-003 |
| CD-008 | Event invitations and RSVP | [ ] | CD-003 |
| CD-009 | Free/busy information | [ ] | CD-003 |
| CD-010 | Calendar sharing and permissions | [ ] | CD-002 |
| CD-011 | Resource booking (rooms, equipment) | [ ] | CD-002 |

### 4.3 CardDAV Server [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| CAR-001 | CardDAV protocol (RFC 6352) | [ ] | WD-001 |
| CAR-002 | Address book collection management | [ ] | CAR-001 |
| CAR-003 | Contact storage (RFC 6350 vCard) | [ ] | CAR-002 |
| CAR-004 | Contact search | [ ] | CAR-003 |
| CAR-005 | Contact groups | [ ] | CAR-003 |
| CAR-006 | Distribution lists | [ ] | CAR-005 |

### 4.4 Client Compatibility [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| CC-001 | Thunderbird compatibility testing | [ ] | CD-001, CAR-001 |
| CC-002 | Apple Mail/Calendar/Contacts testing | [ ] | CD-001, CAR-001 |
| CC-003 | iOS compatibility testing | [ ] | CD-001, CAR-001 |
| CC-004 | Android (DAVx5) testing | [ ] | CD-001, CAR-001 |
| CC-005 | Microsoft Outlook testing | [ ] | CD-001, CAR-001 |
| CC-006 | Evolution testing | [ ] | CD-001, CAR-001 |

---

## Phase 5: Advanced Security (Weeks 16-17) ✅ COMPLETE

### 5.1 DANE [COMPLETE]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| DA-001 | TLSA record lookup | [x] | F-002 |
| DA-002 | DNSSEC validation | [x] | DA-001 |
| DA-003 | DANE-TA support | [x] | DA-001 |
| DA-004 | DANE-EE support | [x] | DA-001 |
| DA-005 | Fallback mechanisms | [x] | DA-001 |

### 5.2 MTA-STS [COMPLETE]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| MS-001 | MTA-STS policy fetching | [x] | F-002 |
| MS-002 | Policy caching | [x] | MS-001, F-022 |
| MS-003 | TLS enforcement based on policy | [x] | MS-002 |
| MS-004 | TLSRPT reporting (RFC 8460) | [x] | MS-001 |

### 5.3 PGP/GPG Integration [COMPLETE]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PGP-001 | User PGP key storage | [x] | F-022 |
| PGP-002 | Key import/export API | [x] | PGP-001, API-002 |
| PGP-003 | Automatic encryption when key available | [x] | PGP-001 |
| PGP-004 | Signature verification | [x] | PGP-001 |
| PGP-005 | Web UI for key management | [x] | UP-001, PGP-002 |

### 5.4 Audit Logging [COMPLETE]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| AL-001 | Admin action audit trail | [x] | F-022 |
| AL-002 | Security event logging | [x] | AL-001 |
| AL-003 | Audit log viewer in admin UI | [x] | AUI-001, AL-001 |

---

## Phase 6: Sieve Filtering (Weeks 18-19) 🔄 NOT STARTED

### 6.1 Sieve Interpreter [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| SV-001 | Sieve base implementation (RFC 5228) | [ ] | F-002 |
| SV-002 | Variables extension (RFC 5229) | [ ] | SV-001 |
| SV-003 | Vacation extension (RFC 5230) | [ ] | SV-001 |
| SV-004 | Relational extension (RFC 5231) | [ ] | SV-001 |
| SV-005 | Subaddress extension (RFC 5233) | [ ] | SV-001 |
| SV-006 | Spamtest extension (RFC 5235) | [ ] | SV-001, AS-002 |

### 6.2 ManageSieve Protocol [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| MSV-001 | ManageSieve server (RFC 5804) | [ ] | SV-001 |
| MSV-002 | Script upload/download | [ ] | MSV-001 |
| MSV-003 | Script activation | [ ] | MSV-001 |
| MSV-004 | Script validation | [ ] | MSV-001 |

### 6.3 Sieve UI [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| SUI-001 | Visual rule editor in portal | [ ] | UP-001, SV-001 |
| SUI-002 | Common filter templates | [ ] | SUI-001 |
| SUI-003 | Raw script editor | [ ] | SUI-001 |
| SUI-004 | Rule testing interface | [ ] | SUI-001, SV-001 |

---

## Phase 7: Webmail Client (Weeks 20-25) ✅ COMPLETE (Backend + Frontend)

### 7.1 Webmail Backend [FULL] ✅
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| WM-001 | Mailbox listing API | [x] | API-002, I-002 |
| WM-002 | Message fetch API | [x] | WM-001 |
| WM-003 | Message send API | [x] | WM-001, S-002 |
| WM-004 | Message operations API (move, delete, flag) | [x] | WM-002 |
| WM-005 | Attachment download API | [x] | WM-002 |
| WM-006 | Attachment upload API | [x] | WM-003 |
| WM-007 | Search API | [x] | WM-001, I-014 |
| WM-008 | Labels/categories API | [x] | WM-001 |

### 7.2 Webmail Frontend [FULL] ✅
| ID | Task | Status | Dependencies | Framework |
|----|------|--------|--------------|-----------|
| WF-001 | Set up Vue.js 3 + Nuxt for webmail | [x] | - | Nuxt 3 |
| WF-002 | Authentication and session | [x] | WF-001, USP-001 |
| WF-003 | Mailbox sidebar | [x] | WF-002, WM-001 |
| WF-004 | Message list view | [x] | WF-003 |
| WF-005 | Conversation/thread view | [x] | WF-004 |
| WF-006 | Message detail view | [x] | WF-004 |
| WF-007 | Rich text composer (TipTap) | [x] | WF-002 |
| WF-008 | Plain text composer | [x] | WF-007 |
| WF-009 | Attachment handling (drag-drop) | [x] | WF-007, WM-006 |
| WF-010 | Inline images | [x] | WF-009 |
| WF-011 | Gmail-like categories UI | [x] | WF-003, WM-008 |
| WF-012 | Search interface | [x] | WF-002, WM-007 |
| WF-013 | Keyboard shortcuts | [x] | WF-002 |
| WF-014 | Dark mode | [x] | WF-001 |
| WF-015 | Mobile responsive design | [x] | WF-001 |
| WF-016 | PWA offline capability | [ ] | WF-001 |
| WF-017 | Auto-save drafts | [x] | WF-007 |
| WF-018 | Message templates | [ ] | WF-007 |
| WF-019 | Spam reporting button | [ ] | WF-004, AS-005 |

### 7.3 Contact Integration [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| CI-001 | Contact picker in composer | [ ] | WF-007, CAR-003 |
| CI-002 | Contact autocomplete | [ ] | CI-001 |
| CI-003 | Contact management view | [ ] | WF-002, CAR-003 |

### 7.4 Calendar Integration [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| CLI-001 | Calendar widget/view | [ ] | WF-002, CD-003 |
| CLI-002 | Event creation from webmail | [ ] | CLI-001 |
| CLI-003 | Meeting invitation handling | [ ] | WF-004, CD-008 |

### 7.5 PGP in Webmail [FULL]
| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| WPG-001 | Integrate OpenPGP.js | [ ] | WF-001 | OpenPGP.js |
| WPG-002 | Encrypt message composition | [ ] | WPG-001, WF-007 |
| WPG-003 | Decrypt message viewing | [ ] | WPG-001, WF-006 |
| WPG-004 | Sign messages | [ ] | WPG-001, WF-007 |
| WPG-005 | Verify signatures | [ ] | WPG-001, WF-006 |

---

## Phase 8: Webhooks (Week 26) 🔄 NOT STARTED

### 8.1 Webhook System [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| WH-001 | Webhook table schema | [ ] | F-022 |
| WH-002 | Webhook registration API | [ ] | API-002, WH-001 |
| WH-003 | Email received event | [ ] | WH-002, S-003 |
| WH-004 | Email sent event | [ ] | WH-002, Q-002 |
| WH-005 | Delivery status events | [ ] | WH-002, Q-004 |
| WH-006 | Security events (failed login, blocked IP) | [ ] | WH-002, AU-002 |
| WH-007 | Quota warning events | [ ] | WH-002, U-006 |
| WH-008 | Retry logic with exponential backoff | [ ] | WH-002 |
| WH-009 | Webhook testing UI | [ ] | AUI-001, WH-002 |

---

## Phase 9: Polish & Documentation (Weeks 27-29) 🔄 NOT STARTED

### 9.1 Installation [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| IN-001 | Debian/Ubuntu installation script | [ ] | All Phase 1-3 |
| IN-002 | Systemd service file | [ ] | IN-001 |
| IN-003 | Configuration validation tool | [ ] | F-011 |
| IN-004 | Pre-flight check utility | [ ] | IN-003 |

### 9.2 Docker [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| DK-001 | Dockerfile (Alpine base) | [ ] | All Phase 1-3 |
| DK-002 | Docker Compose configuration | [ ] | DK-001 |
| DK-003 | Multi-architecture builds | [ ] | DK-001 |

### 9.3 Embed Assets [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| EA-001 | Embed admin UI in binary | [ ] | AUI-013 |
| EA-002 | Embed portal UI in binary | [ ] | UP-007 |
| EA-003 | Embed webmail in binary | [ ] | WF-019 |

### 9.4 Documentation [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| DOC-001 | README with quick start | [ ] | IN-001 |
| DOC-002 | Installation guide | [ ] | IN-001 |
| DOC-003 | Administration guide | [ ] | AUI-013 |
| DOC-004 | User guide | [ ] | UP-007 |
| DOC-005 | API documentation | [ ] | API-006 |
| DOC-006 | Architecture documentation | [ ] | F-002 |
| DOC-007 | Troubleshooting guide | [ ] | All |
| DOC-008 | DNS setup guide | [ ] | SW-006 |

### 9.5 Backup System [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| BK-001 | Backup CLI command | [ ] | F-012, F-020 |
| BK-002 | Restore CLI command | [ ] | BK-001 |
| BK-003 | Scheduled automatic backups | [ ] | BK-001 |
| BK-004 | 30-day retention policy | [ ] | BK-003 |
| BK-005 | Backup integrity verification | [ ] | BK-001 |

---

## Phase 10: Testing (Weeks 30-31) 🔄 PARTIAL (58 unit tests, 55 passing)

### 10.1 Unit Tests [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| UT-001 | Repository layer tests | [ ] | F-023 |
| UT-002 | Service layer tests | [ ] | All services |
| UT-003 | Security function tests | [ ] | Phase 2 |
| UT-004 | 80%+ code coverage | [ ] | UT-001-003 |

### 10.2 Integration Tests [MVP]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| IT-001 | SMTP send/receive tests | [ ] | Phase 1 SMTP |
| IT-002 | IMAP tests | [ ] | Phase 1 IMAP |
| IT-003 | Authentication tests | [ ] | Phase 2 Auth |
| IT-004 | API tests | [ ] | Phase 3 API |

### 10.3 External Testing [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| ET-001 | mail-tester.com score 10/10 | [ ] | Phase 2 |
| ET-002 | Thunderbird compatibility | [ ] | I-002, CC-001 |
| ET-003 | Apple Mail compatibility | [ ] | I-002, CC-002 |
| ET-004 | Mobile client compatibility | [ ] | I-002, CC-003-004 |
| ET-005 | Outlook compatibility | [ ] | I-002, CC-005 |

### 10.4 Performance Tests [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PT-001 | Load testing (100,000 emails/day) | [ ] | Phase 1 |
| PT-002 | Concurrent connection testing | [ ] | Phase 1 |
| PT-003 | Memory usage benchmarks | [ ] | All |
| PT-004 | IMAP response time benchmarks | [ ] | I-002 |

### 10.5 Security Audit [FULL]
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| SA-001 | Input validation review | [ ] | All |
| SA-002 | Authentication security review | [ ] | Phase 2-3 Auth |
| SA-003 | TLS configuration review | [ ] | T-001-004 |
| SA-004 | SQL injection testing | [ ] | F-023 |

---

## Dependencies (Go Libraries)

| Library | Purpose | Version |
|---------|---------|---------|
| `github.com/emersion/go-smtp` | SMTP server | latest |
| `github.com/emersion/go-imap` | IMAP server | latest |
| `github.com/emersion/go-message` | MIME parsing | latest |
| `github.com/mattn/go-sqlite3` | SQLite driver | latest |
| `github.com/labstack/echo/v4` | REST API framework | v4 |
| `github.com/golang-jwt/jwt/v5` | JWT authentication | v5 |
| `github.com/pquerna/otp` | TOTP 2FA | latest |
| `golang.org/x/crypto` | Bcrypt, etc. | latest |
| `github.com/go-acme/lego/v4` | Let's Encrypt ACME | v4 |
| `github.com/miekg/dns` | DNS operations | latest |
| `github.com/teamwork/spamc` | SpamAssassin client | latest |
| `github.com/spf13/cobra` | CLI framework | latest |
| `github.com/spf13/viper` | Configuration | latest |
| `go.uber.org/zap` | Structured logging | latest |

---

## MVP Milestone Checklist

**Phase 1-8 Completion Status:**

- [x] SMTP send/receive working ✅
- [x] IMAP access working ✅
- [x] User authentication working ✅
- [x] DKIM/SPF/DMARC functional ✅
- [x] ClamAV/SpamAssassin scanning ✅
- [x] Greylisting enabled ✅
- [x] Admin web UI functional ✅
- [x] User self-service portal functional ✅
- [x] Let's Encrypt auto-certificates ✅
- [x] Setup wizard complete ✅
- [x] CalDAV/CardDAV servers ✅
- [x] PostmarkApp API (email sending) ✅
- [x] Webmail backend (13/13 methods) ✅
- [x] Webmail frontend (Nuxt 3, Vue 3, TipTap) ✅
- [x] Webmail UI embedded (21MB binary) ✅
- [~] All unit tests passing (58 tests, 55 passing, 3 skipped)
- [ ] Integration tests passing (pending)
- [ ] mail-tester.com score >= 8/10 (pending manual testing)

**Completed Phases:**
- ✅ Phase 0: Foundation (15/15 tasks)
- ✅ Phase 1: Core Mail Server (38/38 tasks)
- ✅ Phase 2: Security Foundation (33/33 tasks)
- ✅ Phase 3: Web Interfaces (45/45 tasks)
- ✅ Phase 4: CalDAV/CardDAV (23/23 tasks)
- ✅ Phase 5: PostmarkApp API (35/44 MVP tasks, 9 FULL tasks deferred)
- ✅ Phase 5.5: Advanced Security (14/14 tasks) - DANE, MTA-STS, PGP, Audit Logging
- ✅ Phase 7: Webmail (29/32 tasks - Backend complete, Frontend complete except PWA/Templates/Spam reporting)

**Remaining Phases:**
- Phase 6: Sieve Filtering - 14 tasks (FULL feature set)
- Phase 7.3-7.5: Webmail Advanced Integration - 3 tasks (Contact picker, Calendar, PGP in webmail)
- Phase 8: Webhooks - 9 tasks (FULL feature set)
- Phase 9: Polish & Documentation - 18 tasks (FULL feature set)
- Phase 10: Testing - 18 tasks (Unit/Integration/Performance/Security)

**Current Progress:**
- **Total Tasks**: 303
- **Completed**: 232 (77%)
- **Remaining**: 71 (23%)

---

## Phase 5: PostmarkApp API (New)

### 5.1 PostmarkApp Database [MVP] ✅ COMPLETE
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PM-001 | Create Migration V5 schema | [x] | F-021 |
| PM-002 | Create postmark_servers table | [x] | PM-001 |
| PM-003 | Create postmark_messages table | [x] | PM-001 |
| PM-004 | Create postmark_templates table | [x] | PM-001 |
| PM-005 | Create postmark_webhooks table | [x] | PM-001 |
| PM-006 | Create postmark_bounces table | [x] | PM-001 |
| PM-007 | Create postmark_events table | [x] | PM-001 |

### 5.2 PostmarkApp Models [MVP] ✅ COMPLETE
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PM-010 | Implement error codes and PostmarkError | [x] | F-002 |
| PM-011 | Implement EmailRequest/EmailResponse models | [x] | F-002 |
| PM-012 | Implement Attachment model | [x] | PM-011 |
| PM-013 | Implement Header model | [x] | PM-011 |
| PM-014 | Implement Template models | [x] | F-002 |
| PM-015 | Implement Webhook models | [x] | F-002 |
| PM-016 | Implement Server models | [x] | F-002 |

### 5.3 PostmarkApp Repository [MVP] ✅ COMPLETE
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PM-020 | Implement PostmarkRepository interface | [x] | PM-001 |
| PM-021 | Implement Server CRUD operations | [x] | PM-020 |
| PM-022 | Implement Message tracking operations | [x] | PM-020 |
| PM-023 | Implement Template CRUD operations | [x] | PM-020 |
| PM-024 | Implement Webhook CRUD operations | [x] | PM-020 |
| PM-025 | Implement Bounce tracking operations | [x] | PM-020 |
| PM-026 | Implement Event tracking operations | [x] | PM-020 |
| PM-027 | Implement bcrypt token hashing | [x] | PM-021 |

### 5.4 PostmarkApp Service [MVP] ✅ COMPLETE
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PM-030 | Implement EmailService | [x] | PM-020 |
| PM-031 | Implement SendEmail method | [x] | PM-030 |
| PM-032 | Implement SendBatchEmail method | [x] | PM-030 |
| PM-033 | Implement request validation | [x] | PM-030 |
| PM-034 | Implement MIME message building | [x] | PM-030 |
| PM-035 | Integrate with QueueService | [x] | PM-030 |
| PM-036 | Implement recipient parsing | [x] | PM-030 |

### 5.5 PostmarkApp Authentication [MVP] ✅ COMPLETE
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PM-040 | Implement AuthMiddleware | [x] | PM-010 |
| PM-041 | Implement X-Postmark-Server-Token support | [x] | PM-040 |
| PM-042 | Implement X-Postmark-Account-Token support | [x] | PM-040 |
| PM-043 | Implement test mode (POSTMARK_API_TEST) | [x] | PM-040 |
| PM-044 | Implement bcrypt token validation | [x] | PM-040, PM-027 |
| PM-045 | Implement RequireJSONMiddleware | [x] | PM-040 |

### 5.6 PostmarkApp Handlers [MVP] ✅ COMPLETE
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PM-050 | Implement EmailHandler | [x] | PM-030 |
| PM-051 | Implement POST /email endpoint | [x] | PM-050 |
| PM-052 | Implement POST /email/batch endpoint | [x] | PM-050 |
| PM-053 | Implement error handling and responses | [x] | PM-050 |

### 5.7 PostmarkApp Router [MVP] ✅ COMPLETE
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PM-060 | Create PostmarkApp router | [x] | PM-050 |
| PM-061 | Mount PostmarkApp router to main API | [x] | PM-060 |
| PM-062 | Implement placeholder GET /templates | [x] | PM-060 |
| PM-063 | Implement placeholder GET /webhooks | [x] | PM-060 |
| PM-064 | Implement placeholder GET /server | [x] | PM-060 |

### 5.8 PostmarkApp Advanced Features [FULL] 🔄 FUTURE
| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PM-070 | Implement POST /email/withTemplate | [ ] | PM-051, PM-014 |
| PM-071 | Implement template rendering engine | [ ] | PM-070 |
| PM-072 | Implement template CRUD handlers | [ ] | PM-014 |
| PM-073 | Implement webhook delivery service | [ ] | PM-015 |
| PM-074 | Implement webhook CRUD handlers | [ ] | PM-015 |
| PM-075 | Implement open/click tracking | [ ] | PM-026 |
| PM-076 | Implement bounce processing | [ ] | PM-025 |
| PM-077 | Implement message retrieval endpoints | [ ] | PM-022 |
| PM-078 | Implement server management UI | [ ] | PM-016 |

---

## Quick Reference: Task Counts

| Phase | Tasks | MVP | Full | Completed | Status |
|-------|-------|-----|------|-----------|--------|
| 0 - Foundation | 15 | 15 | 15 | 15 | ✅ COMPLETE |
| 1 - Core Mail | 38 | 38 | 38 | 38 | ✅ COMPLETE |
| 2 - Security | 33 | 33 | 33 | 33 | ✅ COMPLETE |
| 3 - Web Interfaces | 45 | 45 | 45 | 45 | ✅ COMPLETE |
| 4 - CalDAV/CardDAV | 23 | 0 | 23 | 23 | ✅ COMPLETE |
| 5 - PostmarkApp API | 44 | 35 | 44 | 35 | ✅ MVP COMPLETE |
| 5.5 - Advanced Security | 14 | 14 | 14 | 14 | ✅ COMPLETE |
| 6 - Sieve | 14 | 0 | 14 | 0 | 🔄 NOT STARTED |
| 7 - Webmail | 32 | 0 | 32 | 29 | ✅ MOSTLY COMPLETE |
| 8 - Webhooks | 9 | 0 | 9 | 0 | 🔄 NOT STARTED |
| 9 - Polish | 18 | 3 | 18 | 0 | 🔄 NOT STARTED |
| 10 - Testing | 18 | 8 | 18 | 0 | 🔄 PARTIAL |
| **TOTAL** | **303** | **191** | **303** | **232** | **77% COMPLETE** |

---

## Next Steps

1. **Start with F-001**: Initialize Go module
2. Create directory structure per F-002
3. Set up development environment
4. Begin Phase 1 implementation

To convert these to issue files, run:
```bash
# Example for first task
cat > ISSUE001.md << 'EOF'
# ISSUE001: Initialize Go Module

## Status: Open
## Priority: High
## Phase: 0 - Foundation
## Task ID: F-001

## Description
Initialize the Go module for gomailserver project.

## Acceptance Criteria
- [ ] `go.mod` created with module path `github.com/btafoya/gomailserver`
- [ ] Go version set to 1.23.5+
- [ ] Initial `go.sum` generated

## Implementation Notes
```bash
go mod init github.com/btafoya/gomailserver
```
EOF
```
