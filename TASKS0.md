# Phase 0: Foundation (Week 0)

**Status**: Not Started
**Priority**: MVP - Required
**Estimated Duration**: 1 week
**Dependencies**: None (starting point)

---

## Overview

Establish the project foundation including Go module setup, package structure, core infrastructure (logging, configuration, CLI), and database foundation with SQLite.

---

## 0.1 Project Setup [MVP]

| ID | Task | Status | Dependencies | Notes |
|----|------|--------|--------------|-------|
| F-001 | Initialize Go module (`github.com/btafoya/gomailserver`) | [ ] | - | `go mod init` |
| F-002 | Create package structure (clean architecture) | [ ] | F-001 | See structure below |
| F-003 | Set up golangci-lint configuration | [ ] | F-001 | `.golangci.yml` |
| F-004 | Create Makefile for common tasks | [ ] | F-001 | build, test, lint, run |
| F-005 | Set up GitHub Actions CI/CD | [ ] | F-003 | `.github/workflows/` |

### Package Structure

```
cmd/
  gomailserver/
    main.go                 # Application entrypoint
internal/
  config/                   # Configuration management
    config.go
    loader.go
  database/                 # SQLite connection, migrations
    sqlite.go
    migrations/
  domain/                   # Domain models (entities)
    user.go
    domain.go
    message.go
    mailbox.go
  repository/               # Data access layer (interfaces + implementations)
    user_repository.go
    domain_repository.go
    message_repository.go
  service/                  # Business logic layer
    user_service.go
    domain_service.go
    mail_service.go
  smtp/                     # SMTP server implementation
  imap/                     # IMAP server implementation
  caldav/                   # CalDAV server
  carddav/                  # CardDAV server
  security/                 # DKIM, SPF, DMARC, etc.
    dkim/
    spf/
    dmarc/
  api/                      # REST API handlers
    handlers/
    middleware/
    routes.go
  webmail/                  # Webmail backend
pkg/
  sieve/                    # Sieve interpreter (if custom)
web/
  admin/                    # Admin UI assets (Vue.js)
  portal/                   # User portal assets
  webmail/                  # Webmail client assets
migrations/                 # SQL migration files
scripts/                    # Build and deployment scripts
docs/                       # Documentation
```

---

## 0.2 Core Infrastructure [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| F-010 | Implement structured logging (JSON format) | [ ] | F-002 | `go.uber.org/zap` |
| F-011 | Create configuration system (YAML + env vars) | [ ] | F-002 | `github.com/spf13/viper` |
| F-012 | Implement CLI framework | [ ] | F-002 | `github.com/spf13/cobra` |
| F-013 | Create graceful shutdown handler | [ ] | F-010 | `os/signal`, `context` |
| F-014 | Implement context-based cancellation | [ ] | F-013 | `context` |

### F-010: Structured Logging Details

```go
// internal/logger/logger.go
package logger

import "go.uber.org/zap"

type Logger interface {
    Debug(msg string, fields ...zap.Field)
    Info(msg string, fields ...zap.Field)
    Warn(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
    With(fields ...zap.Field) Logger
}
```

### F-011: Configuration Structure

```yaml
# gomailserver.yaml
server:
  hostname: mail.example.com

database:
  path: /var/lib/gomailserver/mailserver.db

smtp:
  submission_port: 587
  relay_port: 25
  smtps_port: 465

imap:
  port: 143
  imaps_port: 993

tls:
  mode: letsencrypt  # or "manual"
  cert_path: ""
  key_path: ""
  cloudflare_api_token: ""

clamav:
  socket: /var/run/clamav/clamd.ctl

spamassassin:
  host: localhost
  port: 783
```

### F-012: CLI Commands

```
gomailserver
├── run                 # Start the server
├── migrate             # Run database migrations
├── backup              # Create backup
├── restore             # Restore from backup
├── user                # User management
│   ├── create
│   ├── delete
│   └── list
├── domain              # Domain management
│   ├── create
│   ├── delete
│   └── list
└── version             # Show version
```

---

## 0.3 Database Foundation [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| F-020 | SQLite connection management with WAL mode | [ ] | F-002 | `github.com/mattn/go-sqlite3` |
| F-021 | Database migration framework | [ ] | F-020 | `github.com/golang-migrate/migrate` |
| F-022 | Create schema version 1 (all tables) | [ ] | F-021 | See schema below |
| F-023 | Implement repository pattern interfaces | [ ] | F-020 | - |
| F-024 | SQLite PRAGMA optimizations | [ ] | F-020 | See optimizations below |

### F-020: SQLite Connection

```go
// internal/database/sqlite.go
package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
    *sql.DB
}

func Open(path string) (*DB, error) {
    db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000")
    if err != nil {
        return nil, err
    }
    return &DB{db}, nil
}
```

### F-022: Schema V1 (Core Tables)

```sql
-- migrations/001_initial_schema.up.sql

-- Domains
CREATE TABLE domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    max_users INTEGER DEFAULT 0,
    max_mailbox_size INTEGER DEFAULT 0,
    default_quota INTEGER DEFAULT 1073741824,
    catchall_address TEXT,
    dkim_selector TEXT,
    dkim_private_key TEXT,
    spf_policy TEXT,
    dmarc_policy TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Users
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    display_name TEXT,
    quota INTEGER DEFAULT 1073741824,
    used_quota INTEGER DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'disabled', 'suspended')),
    totp_secret TEXT,
    totp_enabled INTEGER DEFAULT 0,
    language TEXT DEFAULT 'en',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME
);

-- Aliases
CREATE TABLE aliases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    address TEXT NOT NULL UNIQUE,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    destinations TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Mailboxes/Folders
CREATE TABLE mailboxes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    parent_id INTEGER REFERENCES mailboxes(id) ON DELETE CASCADE,
    subscribed INTEGER DEFAULT 1,
    special_use TEXT,
    uid_validity INTEGER NOT NULL,
    uid_next INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Messages
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mailbox_id INTEGER NOT NULL REFERENCES mailboxes(id) ON DELETE CASCADE,
    uid INTEGER NOT NULL,
    message_id TEXT,
    subject TEXT,
    sender TEXT,
    recipients TEXT,
    date DATETIME,
    size INTEGER NOT NULL,
    flags TEXT DEFAULT '',
    headers TEXT,
    storage_type TEXT NOT NULL CHECK (storage_type IN ('blob', 'file')),
    content BLOB,
    file_path TEXT,
    thread_id TEXT,
    labels TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(mailbox_id, uid)
);

-- SMTP Queue
CREATE TABLE smtp_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender TEXT NOT NULL,
    recipients TEXT NOT NULL,
    message_id INTEGER REFERENCES messages(id),
    raw_message BLOB,
    retry_count INTEGER DEFAULT 0,
    next_retry DATETIME,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'failed', 'sent')),
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Sessions
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    protocol TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_activity DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Failed Logins
CREATE TABLE failed_logins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL,
    username TEXT,
    attempted_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Blacklisted IPs
CREATE TABLE blacklisted_ips (
    ip_address TEXT PRIMARY KEY,
    reason TEXT,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Whitelisted IPs
CREATE TABLE whitelisted_ips (
    ip_address TEXT PRIMARY KEY,
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Greylisting
CREATE TABLE greylist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL,
    sender TEXT NOT NULL,
    recipient TEXT NOT NULL,
    first_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    pass_count INTEGER DEFAULT 0,
    UNIQUE(ip_address, sender, recipient)
);

-- Spam Quarantine
CREATE TABLE spam_quarantine (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id INTEGER REFERENCES messages(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    spam_score REAL,
    quarantined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    released INTEGER DEFAULT 0,
    auto_delete_at DATETIME
);

-- Sieve Scripts
CREATE TABLE sieve_scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    active INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Webhooks
CREATE TABLE webhooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    events TEXT NOT NULL,
    auth_token TEXT,
    active INTEGER DEFAULT 1,
    retry_count INTEGER DEFAULT 3,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Audit Log
CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    target_type TEXT,
    target_id INTEGER,
    details TEXT,
    ip_address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_users_domain ON users(domain_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_aliases_domain ON aliases(domain_id);
CREATE INDEX idx_mailboxes_user ON mailboxes(user_id);
CREATE INDEX idx_messages_mailbox ON messages(mailbox_id);
CREATE INDEX idx_messages_user ON messages(user_id);
CREATE INDEX idx_messages_date ON messages(date);
CREATE INDEX idx_messages_subject ON messages(subject);
CREATE INDEX idx_messages_sender ON messages(sender);
CREATE INDEX idx_smtp_queue_status ON smtp_queue(status);
CREATE INDEX idx_smtp_queue_next_retry ON smtp_queue(next_retry);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_failed_logins_ip ON failed_logins(ip_address);
CREATE INDEX idx_greylist_lookup ON greylist(ip_address, sender, recipient);
CREATE INDEX idx_audit_log_user ON audit_log(user_id);
CREATE INDEX idx_audit_log_action ON audit_log(action);
```

### F-024: SQLite PRAGMA Optimizations

```go
func (db *DB) Optimize() error {
    pragmas := []string{
        "PRAGMA journal_mode = WAL",
        "PRAGMA synchronous = NORMAL",
        "PRAGMA cache_size = -64000",  // 64MB cache
        "PRAGMA temp_store = MEMORY",
        "PRAGMA mmap_size = 268435456", // 256MB mmap
        "PRAGMA foreign_keys = ON",
    }
    for _, pragma := range pragmas {
        if _, err := db.Exec(pragma); err != nil {
            return err
        }
    }
    return nil
}
```

---

## Acceptance Criteria

- [ ] `go mod init` completed successfully
- [ ] All directories created per structure
- [ ] `golangci-lint run` passes with no errors
- [ ] `make build` produces working binary
- [ ] CLI `gomailserver version` works
- [ ] SQLite database created with all tables
- [ ] Migrations run successfully
- [ ] Configuration loads from YAML and env vars
- [ ] Structured JSON logging works
- [ ] Graceful shutdown handles SIGTERM/SIGINT

---

## Go Dependencies for Phase 0

```go
// go.mod additions
require (
    github.com/mattn/go-sqlite3 v1.14.22
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.2
    go.uber.org/zap v1.26.0
    github.com/golang-migrate/migrate/v4 v4.17.0
)
```

---

## Next Phase

After completing Phase 0, proceed to [TASKS1.md](TASKS1.md) - Core Mail Server.
