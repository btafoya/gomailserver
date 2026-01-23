# GoMailServer Project TODO List Summary

**Generated**: 2026-01-22  
**Total Tasks**: 57  
**Status**: Project in critical state - does not compile

## Quick Overview

### Critical Blockers (6 tasks - 1-2 weeks estimated)
These tasks must be completed before any other work can proceed.

| ID | Task | Est. Time | Priority |
|----|------|-----------|----------|
| CRITICAL-1 | Fix repository factory missing imports | 1 hour | 🔴 HIGH |
| CRITICAL-2 | Complete 8 PostgreSQL repository stubs | 3-5 days | 🔴 HIGH |
| CRITICAL-3 | Fix UserService interface (6 missing methods) | 1 day | 🔴 HIGH |
| CRITICAL-4 | Fix API handler compilation errors | 1 day | 🔴 HIGH |
| CRITICAL-5 | Fix integration test dependency | 1-2 hours | 🔴 HIGH |
| CRITICAL-6 | Verify project builds successfully | 2-4 hours | 🔴 HIGH |

### Phase Completion Status

| Phase | Description | Tasks | Status |
|-------|-------------|-------|--------|
| Phase 1 | Core Mail Server | 5 | 🟡 60% (needs verification) |
| Phase 2 | Security Foundation | 9 | 🟡 70% (needs testing) |
| Phase 3 | REST API & Admin | 6 | 🟢 90% (mostly complete) |
| Phase 4 | CalDAV/CardDAV | 6 | 🟢 90% (needs testing) |
| Phase 5 | PostmarkApp API | 5 | 🟡 60% (extensions incomplete) |
| Phase 6 | Sieve Filtering | 4 | 🔴 0% (not started) |
| Phase 7 | Webmail Client | 2 | 🟢 90% (nice-to-haves) |
| Phase 8 | Reputation Phase 5 | 4 | 🟡 85% (integration pending) |
| Phase 9 | Polish & Docs | 5 | 🟡 20% (not started) |
| Phase 10 | Testing | 6 | 🟡 30% (partial) |
| Infrastructure | CI/CD & Ops | 6 | 🔴 0% (not started) |

## Work Priority Order

### Week 1-2: Emergency Stabilization
1. CRITICAL-1 through CRITICAL-6 (restore build)
2. INFRA-1, INFRA-2, INFRA-3 (CI/CD foundation)

### Week 3-4: Core Protocol Verification
3. PHASE1-1 through PHASE1-5 (SMTP/IMAP verification)
4. PHASE2-1 through PHASE2-9 (security testing)
5. INFRA-4, INFRA-5 (race detector, health checks)

### Week 5-6: API & UI Verification
6. PHASE3-1 through PHASE3-6 (REST API & admin UI)
7. PHASE4-1 through PHASE4-6 (CalDAV/CardDAV testing)

### Week 7-8: Advanced Features
8. PHASE8-1 through PHASE8-4 (reputation phase 5 completion)
9. PHASE10-1 through PHASE10-3 (test suite completion)

### Week 9+: Polish & Completion
10. Remaining tasks as needed

## Agent Execution Guidelines

### When Working on Tasks:
1. **Always read the TODO list first** - use `todoread` to check current status
2. **Mark tasks as in_progress** when starting - use `todowrite`
3. **Complete one task at a time** - don't batch multiple tasks
4. **Mark tasks as completed** when done - update status
5. **Run tests after each task** - `make test` or `go test ./...`
6. **Run linter after each task** - `make lint` or `golangci-lint run`
7. **Update AGENTS.md** if you discover new patterns or conventions

### Task Dependencies:
- All PHASE tasks depend on CRITICAL-6 (successful build)
- All Phase tasks within a phase can be done in parallel
- INFRA-1, INFRA-2, INFRA-3 should be done after CRITICAL-6
- Testing tasks (Phase 10) depend on their respective feature phases

### File Locations Reference:
- PostgreSQL repos: `internal/repository/postgres/`
- SQLite repos: `internal/repository/sqlite/`
- Services: `internal/service/`
- API handlers: `internal/api/handlers/`
- Tests: `tests/`
- Unified UI: `unified/`
- Main command: `cmd/gomailserver/main.go`

## Success Metrics

### Immediate (Weeks 1-2):
- [ ] Project builds without errors
- [ ] All 6 CRITICAL tasks complete
- [ ] CI/CD pipeline established

### Short-term (Weeks 3-6):
- [ ] Can send/receive email via SMTP
- [ ] Can retrieve mail via IMAP
- [ ] All security features tested
- [ ] Admin UI serves and functions

### Medium-term (Weeks 7-10):
- [ ] CalDAV/CardDAV clients connect
- [ ] Reputation management complete
- [ ] 80%+ test coverage
- [ ] mail-tester.com 10/10 score

### Long-term (Months 2-3):
- [ ] Sieve filtering implemented
- [ ] Installation scripts complete
- [ ] Docker builds working
- [ ] Full documentation available

## Notes for Agents

- **CRITICAL tasks are blockers** - do not attempt other tasks until these are done
- **Follow existing patterns** - look at SQLite repos for PostgreSQL implementation examples
- **Test incrementally** - verify each task works before moving on
- **Document discoveries** - update AGENTS.md if you find new conventions
- **Ask questions** - if a task is genuinely ambiguous, ask before implementing

## Related Documentation

- `TRUESTATUS.md` - Detailed project status from 2026-01-20
- `ROADMAP.md` - Original phase completion checklist
- `PR.md` - Full project requirements specification
- `AGENTS.md` - Project conventions and coding standards
- `README.md` - Project overview and usage

---

# Detailed Task Specifications

## CRITICAL Tasks (Must Complete First)

### CRITICAL-1: Fix Repository Factory Missing Imports

**File**: `internal/repository/factory.go:34-43`  
**Issue**: Missing imports for sqlite and postgres packages  
**Action Required**:
```go
// Add these imports:
import (
    "github.com/btafoya/gomailserver/internal/repository/sqlite"
    "github.com/btafoya/gomailserver/internal/repository/postgres"
)
```
**Verification**: `go build ./internal/repository/`

---

### CRITICAL-2: Complete 8 PostgreSQL Repository Stubs

**Files** (all in `internal/repository/postgres/`):
- `alias_repository.go` - Currently empty stub
- `mailbox_repository.go` - Currently empty stub
- `queue_repository.go` - Currently empty stub
- `greylist_repository.go` - Currently empty stub
- `ipblacklist_repository.go` - Currently empty stub
- `loginattempt_repository.go` - Currently empty stub
- `ratelimit_repository.go` - Currently empty stub
- `webhook_repository.go` - Currently empty stub

**Reference Implementation**: Use `internal/repository/sqlite/` as templates  
**Pattern to Follow**: Each repository should implement `Repository` interface with methods matching SQLite equivalents

**Key Methods to Implement** (per repository):
- `Create(ctx context.Context, entity interface{}) error`
- `Update(ctx context.Context, entity interface{}) error`
- `Delete(ctx context.Context, id int64) error`
- `GetByID(ctx context.Context, id int64) (interface{}, error)`
- `ListAll(ctx context.Context) ([]interface{}, error)`

**Verification**: `go build ./internal/repository/postgres/`

---

### CRITICAL-3: Fix UserService Interface

**File**: `internal/service/user_service.go`  
**Missing Methods**:
1. `ListAll() []*domain.User`
2. `GetByID(id int64) *domain.User`
3. `CreateWithPassword(user *domain.User, password string) error`
4. `Update(user *domain.User) error`
5. `Delete(id int64) error`
6. `UpdatePassword(id int64, newPassword string) error`

**Implementation Pattern**:
```go
func (s *UserService) ListAll() ([]*domain.User, error) {
    return s.repo.ListAll()
}

func (s *UserService) GetByID(id int64) (*domain.User, error) {
    return s.repo.GetByID(id)
}

func (s *UserService) CreateWithPassword(user *domain.User, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    return s.repo.Create(user)
}

func (s *UserService) Update(user *domain.User) error {
    return s.repo.Update(user)
}

func (s *UserService) Delete(id int64) error {
    return s.repo.Delete(id)
}

func (s *UserService) UpdatePassword(id int64, newPassword string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user, err := s.GetByID(id)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    return s.repo.Update(user)
}
```

**Verification**: `go build ./internal/service/`

---

### CRITICAL-4: Fix API Handler Compilation Errors

**Files**: All files in `internal/api/handlers/`  
**Issue**: References to missing UserService methods  
**Action Required**:
1. Read each handler file
2. Identify calls to missing UserService methods
3. Either implement the missing methods or update handler logic

**Common Issues**:
- `handler.UserList()` - needs `ListAll()`
- `handler.UserGet()` - needs `GetByID()`
- `handler.UserCreate()` - needs `CreateWithPassword()`
- `handler.UserUpdate()` - needs `Update()`
- `handler.UserDelete()` - needs `Delete()`
- `handler.UserPassword()` - needs `UpdatePassword()`

**Verification**: `go build ./internal/api/handlers/`

---

### CRITICAL-5: Fix Integration Test Dependency

**File**: `tests/integration2/integration.go`  
**Issue**: References `github.com/emersion/go-smtp/v2` not in `go.mod`  
**Options**:
1. Add dependency: `go get github.com/emersion/go-smtp/v2`
2. Or update code to use v1 API: `github.com/emersion/go-smtp` (v1)

**Verification**: `go build ./tests/integration2/`

---

### CRITICAL-6: Verify Project Builds

**Command**: `go build ./...`  
**Expected Result**: All packages compile without errors  
**If Errors**: Fix remaining compilation errors before proceeding

---

## PHASE 1: Core Mail Server

### PHASE1-1: Implement SMTP Protocol

**Port**: 25 (submission), 587 (submission), 465 (implicit TLS)  
**RFC**: RFC 5321  
**Implementation**: `internal/smtp/`  
**Library**: `github.com/emersion/go-smtp/v2`  
**Key Features**:
- [ ] MAIL FROM command
- [ ] RCPT TO command
- [ ] DATA command
- [ ] AUTH PLAIN/LOGIN
- [ ] STARTTLS
- [ ] Queue management

**Verification**# Test with telnet or openssl
openssl s_client -connect localhost::
```bash
465
telnet localhost 25
```

---

### PHASE1-2: Implement IMAP4 Protocol

**Port**: 143 (IMAP), 993 (IMAPS)  
**RFC**: RFC 3501  
**Implementation**: `internal/imap/`  
**Library**: `github.com/emersion/go-imap/v2`  
**Key Features**:
- [ ] LOGIN/AUTHENTICATE
- [ ] LIST/LSUB
- [ ] SELECT/EXAMINE
- [ ] FETCH
- [ ] STORE
- [ ] SEARCH
- [ ] UID commands

**Verification**:
```bash
openssl s_client -connect localhost:993
```

---

### PHASE1-3: Implement Hybrid Storage

**Threshold**: >1MB messages stored on filesystem  
**Implementation**: `internal/repository/message_repository.go`  
**Pattern**:
```go
const BlobThreshold = 1024 * 1024 // 1MB

func (r *MessageRepository) SaveMessage(msg *domain.Message) error {
    if len(msg.Body) > BlobThreshold {
        path := filepath.Join(r.blobPath, msg.ID)
        return os.WriteFile(path, msg.Body, 0644)
    }
    return r.db.Save(msg)
}

func (r *MessageRepository) GetMessage(id string) (*domain.Message, error) {
    msg := &domain.Message{ID: id}
    path := filepath.Join(r.blobPath, id)
    if _, err := os.Stat(path); err == nil {
        data, err := os.ReadFile(path)
        if err != nil {
            return nil, err
        }
        msg.Body = data
        msg.IsBlob = true
        return msg, nil
    }
    return r.db.Get(id)
}
```

---

### PHASE1-4: Implement User Authentication

**Methods**:
- [ ] PLAIN authentication
- [ ] LOGIN authentication
- [ ] bcrypt password hashing
- [ ] TOTP 2FA support

**Implementation**: `internal/security/auth.go`  
**Password Hashing**:
```go
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

---

### PHASE1-5: Implement TLS Support

**Ports**: 465 (implicit), 587 (STARTTLS), 993 (implicit IMAPS)  
**Implementation**: `internal/security/tls.go`  
**Features**:
- [ ] Automatic certificate generation (development)
- [ ] Let's Encrypt integration (production)
- [ ] TLS 1.3 support
- [ ] Strong cipher suites

**Configuration**:
```yaml
tls:
  enabled: true
  port: 465
  cert_path: /etc/gomailserver/cert.pem
  key_path: /etc/gomailserver/key.pem
  lets_encrypt:
    enabled: true
    email: admin@example.com
```

---

## PHASE 2: Security Foundation

### PHASE2-1: Verify ClamAV Integration

**Implementation**: `internal/security/clamav/`  
**Features**:
- [ ] Scan incoming mail
- [ ] Scan outgoing mail
- [ ] Quarantine infected messages
- [ ] Automatic updates

**Configuration**:
```yaml
clamav:
  enabled: true
  host: localhost
  port: 3310
  timeout: 30s
```

---

### PHASE2-2: Verify SpamAssassin Integration

**Implementation**: `internal/security/spamassassin/`  
**Features**:
- [ ] Spam scoring
- [ ] Header tagging
- [ ] Threshold configuration
- [ ] Bayesian learning

**Configuration**:
```yaml
spamassassin:
  enabled: true
  host: localhost
  port: 783
  score_threshold: 5.0
  rewrite_subject: true
```

---

### PHASE2-3: Verify DKIM Signing/Verification

**RFC**: RFC 6376  
**Implementation**: `internal/security/dkim/`  
**Algorithms**: RSA-2048, RSA-4096, Ed25519  
**Features**:
- [ ] Sign outgoing mail
- [ ] Verify incoming mail
- [ ] DNS record generation
- [ ] Key rotation

**Configuration**:
```yaml
dkim:
  enabled: true
  selector: mail
  algorithm: Ed25519
  key_bits: 256
  expiry: 86400 # seconds
```

---

### PHASE2-4: Verify SPF Validation

**RFC**: RFC 7208  
**Implementation**: `internal/security/spf/`  
**Features**:
- [ ] DNS resolver integration
- [ ] Softfail/hardfail policies
- [ ] Sender ID validation
- [ ] Alignment checking

**Configuration**:
```yaml
spf:
  enabled: true
  fail_policy: softfail
  check_helo: true
  check_mailfrom: true
```

---

### PHASE2-5: Verify DMARC Enforcement

**RFC**: RFC 7489  
**Implementation**: `internal/security/dmarc/`  
**Features**:
- [ ] Policy validation
- [ ] Aggregate reporting
- [ ] Forensic reporting
- [ ] DNS record checking

**Configuration**:
```yaml
dmarc:
  enabled: true
  policy: reject
  pct: 100
  ruf:
    - mailto:dmarc-forensics@example.com
  rua:
    - mailto:dmarc-reports@example.com
```

---

### PHASE2-6: Verify Greylisting

**Implementation**: `internal/security/greylist/`  
**Features**:
- [ ] Configurable delays
- [ ] Whitelist management
- [ ] Triplet storage
- [ ] Automatic expiry

**Configuration**:
```yaml
greylist:
  enabled: true
  delay: 5m
  whitelisted:
    - 192.168.0.0/16
    - 10.0.0.0/8
  auto_whitelist:
    enabled: true
    threshold: 10
    expiry: 24h
```

---

### PHASE2-7: Verify Rate Limiting

**Scope**: SMTP, IMAP, Authentication  
**Implementation**: `internal/security/ratelimit/`  
**Features**:
- [ ] Sliding window algorithm
- [ ] Per-user limits
- [ ] Per-IP limits
- [ ] Exemptions

**Configuration**:
```yaml
ratelimit:
  smtp:
    max: 100
    window: 1m
  imap:
    max: 1000
    window: 1m
  auth:
    max: 5
    window: 1h
```

---

### PHASE2-8: Verify Brute Force Protection

**Implementation**: `internal/security/bruteforce/`  
**Features**:
- [ ] IP blacklisting
- [ ] Automatic unblock
- [ ] Login attempt tracking
- [ ] Alert generation

**Configuration**:
```yaml
bruteforce:
  enabled: true
  max_attempts: 5
  window: 1h
  block_duration: 24h
```

---

### PHASE2-9: Verify 2FA/TOTP Support

**Implementation**: `internal/security/totp/`  
**Features**:
- [ ] TOTP generation
- [ ] QR code display
- [ ] Backup codes
- [ ] Recovery options

**Configuration**:
```yaml
totp:
  enabled: true
  issuer: gomailserver
  digits: 6
  period: 30
  algorithm: SHA1
```

---

## PHASE 3: REST API & Admin

### PHASE3-1: Verify REST API Handlers

**Implementation**: `internal/api/handlers/`  
**Authentication**:
- [ ] JWT tokens
- [ ] API keys
- [ ] Session management

**Endpoints**:
- `/api/v1/users`
- `/api/v1/domains`
- `/api/v1/messages`
- `/api/v1/settings`

---

### PHASE3-2: Verify Admin UI

**Path**: `unified/admin/`  
**Features** (17 Phase 6 API endpoints):
- [ ] Domain management
- [ ] User management
- [ ] Alias management
- [ ] Mailbox management
- [ ] Queue monitoring
- [ ] Blacklist management
- [ ] Rate limit configuration
- [ ] Security settings
- [ ] DKIM key management
- [ ] SPF/DMARC configuration
- [ ] Report viewing
- [ ] System status
- [ ] Log viewing
- [ ] Backup/restore
- [ ] Certificate management
- [ ] Plugin management
- [ ] API key management

---

### PHASE3-3: Verify User Portal

**Path**: `unified/portal/`  
**Features**:
- [ ] Domain listing
- [ ] User profile management
- [ ] Password change
- [ ] 2FA setup
- [ ] Alias management
- [ ] Forwarding rules
- [ ] Auto-reply setup
- [ ] Mail filter rules
- [ ] Storage usage
- [ ] Export data

---

### PHASE3-4: Verify Setup Wizard

**Path**: `unified/setup/`  
**Steps**:
1. [ ] Welcome & requirements check
2. [ ] Database configuration
3. [ ] Admin account creation
4. [ ] Domain configuration
5. [ ] Security settings
6. [ ] TLS certificate setup

---

### PHASE3-5: Verify RBAC

**Roles**: Admin, User, Read-only  
**Implementation**: `internal/api/middleware/auth.go`  
**Permissions**:
```go
const (
    RoleAdmin = "admin"
    RoleUser = "user"
    RoleReadOnly = "readonly"
)

var permissions = map[string][]string{
    RoleAdmin: {"*"},
    RoleUser: {
        "users:read", "users:write:own",
        "domains:read",
        "messages:read", "messages:write",
        "aliases:read", "aliases:write:own",
    },
    RoleReadOnly: {
        "users:read",
        "domains:read",
        "messages:read",
    },
}
```

---

### PHASE3-6: Verify Let's Encrypt ACME

**Implementation**: `internal/security/acme/`  
**Features**:
- [ ] Certificate issuance
- [ ] Automatic renewal
- [ ] HTTP-01 challenge
- [ ] TLS-ALPN-01 challenge

**Configuration**:
```yaml
acme:
  enabled: true
  email: admin@example.com
  directory: https://acme-v02.api.letsencrypt.org/directory
  http_port: 80
  tls_port: 443
```

---

## PHASE 4: CalDAV/CardDAV

### PHASE4-1: Verify CalDAV Server

**RFC**: RFC 4791  
**Implementation**: `internal/webdav/caldav/`  
**Features**:
- [ ] Calendar collections
- [ ] Calendar objects (iCalendar)
- [ ] Principal collections
- [ ] Scheduling

**Authentication**: HTTP Basic + OAuth2

---

### PHASE4-2: Verify CardDAV Server

**RFC**: RFC 6352  
**Implementation**: `internal/webdav/carddav/`  
**Features**:
- [ ] Address book collections
- [ ] vCard 4.0 support
- [ ] vCard 3.0 support
- [ ] Group addressing

---

### PHASE4-3: Verify WebDAV Integration

**Implementation**: `internal/webdav/`  
**Features**:
- [ ] PROPFIND
- [ ] PROPPATCH
- [ ] MKCOL
- [ ] DELETE
- [ ] COPY/MOVE
- [ ] LOCK/UNLOCK

---

### PHASE4-4: Test Client Compatibility

**Test Clients**:
- [ ] Thunderbird
- [ ] Apple Mail (macOS)
- [ ] Apple Calendar (iOS)
- [ ] Apple Contacts (iOS)
- [ ] Android (Google)
- [ ] Outlook
- [ ] Evolution
- [ ] iOS native apps

---

### PHASE4-5: Verify Calendar Service

**Implementation**: `internal/service/calendar_service.go`  
**Features**:
- [ ] Create events
- [ ] Update events
- [ ] Delete events
- [ ] List events
- [ ] Recurring events
- [ ] Timezones
- [ ] Reminders

---

### PHASE4-6: Verify Addressbook Service

**Implementation**: `internal/service/addressbook_service.go`  
**Features**:
- [ ] Create contacts
- [ ] Update contacts
- [ ] Delete contacts
- [ ] List contacts
- [ ] Contact groups
- [ ] Photo support
- [ ] vCard export

---

## PHASE 5: PostmarkApp API

### PHASE5-1: Complete Template-Based Email

**Endpoint**: `POST /email/withTemplate`  
**Implementation**: `internal/api/handlers/postmark.go`  
**Features**:
- [ ] Template ID support
- [ ] Template model
- [ ] Inline CSS
- [ ] Track opens

---

### PHASE5-2: Complete Template CRUD

**Endpoints**:
- `GET /templates` - List templates
- `GET /templates/{id}` - Get template
- `POST /templates` - Create template
- `PUT /templates/{id}` - Update template
- `DELETE /templates/{id}` - Delete template

---

### PHASE5-3: Implement Webhook Delivery

**Endpoints**:
- `POST /webhooks/delivery` - Delivery notifications
- `POST /webhooks/bounce` - Bounce notifications
- `POST /webhooks/spam` - Spam complaints

**Features**:
- [ ] Retry logic
- [ ] Signature verification
- [ ] Processing queue

---

### PHASE5-4: Implement Open/Click Tracking

**Features**:
- [ ] Pixel tracking
- [ ] Link rewriting
- [ ] Click counting
- [ ] Geolocation

---

### PHASE5-5: Implement Bounce Processing

**Features**:
- [ ] Bounce detection
- [ ] Bounce categorization
- [ ] Suppression list
- [ ] Reactivation

---

## PHASE 6: Sieve Filtering (NOT STARTED)

### PHASE6-1: Implement Sieve Interpreter

**RFC**: RFC 5228  
**Implementation**: `internal/sieve/interpreter.go`  
**Support**:
- [ ] require statement
- [ ] if/elsif/else
- [ ] keep/discard
- [ ] fileinto
- [ ] redirect
- [ ] reject

---

### PHASE6-2: Implement Sieve Extensions

**Extensions**:
- [ ] variables (RFC 5229)
- [ ] vacation (RFC 5230)
- [ ] relational (RFC 5231)
- [ ] subaddress (RFC 5233)
- [ ] spamtest (RFC 5235)

---

### PHASE6-3: Implement ManageSieve

**RFC**: RFC 5804  
**Port**: 4190  
**Implementation**: `internal/sieve/manage.go`  
**Commands**:
- [ ] AUTHENTICATE
- [ ] GETSCHELLS
- [ ] PUTSCRIPT
- [ ] LISTSCRIPTS
- [ ] SETACTIVE
- [ ] DELETESCRIPT

---

### PHASE6-4: Visual Rule Editor

**UI**: `unified/portal/filters/`  
**Features**:
- [ ] Drag-and-drop rules
- [ ] Condition builder
- [ ] Action selector
- [ ] Script preview

---

## PHASE 7: Webmail Client

### PHASE7-1: Offline Capability

**Implementation**: Service Worker + IndexedDB  
**Features**:
- [ ] Offline email reading
- [ ] Offline compose (queued)
- [ ] Background sync
- [ ] Push notifications

---

### PHASE7-2: Message Templates

**Features**:
- [ ] Save as template
- [ ] Insert template
- [ ] Template variables
- [ ] Drag-and-drop editor

---

## PHASE 8: Reputation Phase 5

### PHASE8-1: Database Migration Scripts

**Files**: `internal/database/migrations/`  
**Tables**:
- [ ] reputation_scores
- [ ] sending_history
- [ ] provider_metrics
- [ ] warmup_schedules
- [ ] alerts

---

### PHASE8-2: API Endpoints

**Endpoints**:
- `GET /api/v1/reputation` - Current reputation
- `GET /api/v1/reputation/history` - Historical data
- `GET /api/v1/reputation/providers` - Provider metrics
- `POST /api/v1/reputation/warmup` - Warmup schedule
- `GET /api/v1/reputation/alerts` - Active alerts

---

### PHASE8-3: Cron Job Integration

**Tasks**:
- [ ] Reputation score updates
- [ ] Provider metrics collection
- [ ] Alert generation
- [ ] Warmup progression
- [ ] Report generation

---

### PHASE8-4: WebUI Components

**Components**:
- [ ] DMARC report viewer
- [ ] External metrics dashboard
- [ ] Provider limits configuration
- [ ] Warm-up scheduling UI
- [ ] Prediction graphs
- [ ] Alert configuration

---

## PHASE 9: Polish & Docs

### PHASE9-1: Installation Scripts

**Scripts**: `scripts/install/`  
**Distros**:
- [ ] Debian 12
- [ ] Ubuntu 22.04+
- [ ] RHEL 9+
- [ ] Rocky Linux 9

**Features**:
- [ ] Dependency installation
- [ ] Configuration generation
- [ ] Service setup
- [ ] Firewall rules

---

### PHASE9-2: Docker Configuration

**Files**: `Dockerfile`, `docker-compose.yml`  
**Features**:
- [ ] Multi-arch builds (amd64, arm64)
- [ ] Health checks
- [ ] Volume management
- [ ] Nginx reverse proxy
- [ ] TLS termination

---

### PHASE9-3: Documentation

**Docs**: `docs/`  
**Sections**:
- [ ] Administrator guide
- [ ] User guide
- [ ] API reference
- [ ] Architecture overview
- [ ] Security hardening
- [ ] Troubleshooting

---

### PHASE9-4: Backup/Restore System

**Implementation**: `internal/service/backup.go`  
**Features**:
- [ ] Full database backup
- [ ] Incremental backups
- [ ] Blob backup
- [ ] Encrypted backups
- [ ] Scheduled backups
- [ ] Point-in-time recovery

---

### PHASE9-5: 30-Day Retention Policy

**Implementation**: `internal/service/retention.go`  
**Features**:
- [ ] Message retention
- [ ] Log retention
- [ ] Audit log retention
- [ ] Automatic cleanup

---

## PHASE 10: Testing

### PHASE10-1: Fix ACME Service Build

**File**: `internal/security/acme/service.go`  
**Issue**: Build failures  
**Fix**: Resolve compilation errors

---

### PHASE10-2: 80%+ Unit Test Coverage

**Target**: `make test` with 80%+ coverage  
**Command**: `go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out`

---

### PHASE10-3: Complete Integration Tests

**Tests**:
- [ ] SMTP sending
- [ ] IMAP retrieval
- [ ] Authentication flows
- [ ] API endpoints
- [ ] Security features

---

### PHASE10-4: External Testing

**Target**: mail-tester.com 10/10 score  
**Requirements**:
- [ ] Valid SPF
- [ ] Valid DKIM
- [ ] Valid DMARC
- [ ] TLS encryption
- [ ] No blacklists
- [ ] Proper headers

---

### PHASE10-5: Performance Benchmarks

**Target**: 100K emails/day  
**Tools**:
- [ ] wrk/hey for API
- [ ] smtpbench for SMTP
- [ ] Custom benchmarks

---

### PHASE10-6: Security Audit

**Tasks**:
- [ ] Penetration testing
- [ ] Code review
- [ ] Dependency audit
- [ ] Configuration review

---

## INFRASTRUCTURE

### INFRA-1: CI/CD Pipeline

**File**: `.github/workflows/ci.yml`  
**Jobs**:
- [ ] Build
- [ ] Unit tests
- [ ] Integration tests
- [ ] Lint
- [ ] Security scan
- [ ] Docker build
- [ ] Deploy (manual)

---

### INFRA-2: Pre-commit Hooks

**File**: `.pre-commit-config.yaml`  
**Hooks**:
- [ ] go fmt
- [ ] go vet
- [ ] go build
- [ ] Unit tests

---

### INFRA-3: golangci-lint Integration

**File**: `.golangci.yml`  
**Linters**:
- [ ] errcheck
- [ ] gofmt
- [ ] goimports
- [ ] gosimple
- [ ] govet
- [ ] ineffassign
- [ ] staticcheck
- [ ] typecheck
- [ ] unused

---

### INFRA-4: Race Detector

**Command**: `go test -race ./...`  
**Integration**: Add to CI pipeline

---

### INFRA-5: Health Check Endpoint

**Endpoint**: `GET /health`  
**Implementation**: `internal/api/handlers/health.go`  
**Checks**:
- [ ] Database connectivity
- [ ] SMTP server status
- [ ] IMAP server status
- [ ] Disk space
- [ ] Memory usage

---

### INFRA-6: Metrics/Observability

**Implementation**: `internal/observability/`  
**Features**:
- [ ] Prometheus metrics
- [ ] OpenTelemetry tracing
- [ ] Structured logging (zap)
- [ ] Dashboard (Grafana)

---

## Notes

This document provides comprehensive task specifications for the GoMailServer project. Agents should follow the priority order, complete tasks sequentially, and verify each task before moving to the next.

For questions or clarifications, refer to the related documentation files listed at the top of this document.
