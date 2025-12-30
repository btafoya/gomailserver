# Phase 1: Core Mail Server (Weeks 1-4)

**Status**: Not Started
**Priority**: MVP - Required
**Estimated Duration**: 3-4 weeks
**Dependencies**: Phase 0 (Foundation)

---

## Overview

Implement the core mail server functionality including SMTP server (submission, relay, SMTPS), IMAP server with full mailbox operations, message storage with hybrid blob/file approach, and queue management.

---

## 1.1 SMTP Server [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| S-001 | Integrate go-smtp library | [ ] | F-002 | `emersion/go-smtp` |
| S-002 | Implement SMTP submission server (port 587) | [ ] | S-001 |
| S-003 | Implement SMTP relay server (port 25) | [ ] | S-001 |
| S-004 | Implement SMTPS server (port 465) | [ ] | S-001, T-001 |
| S-005 | STARTTLS support | [ ] | S-002, T-001 |
| S-006 | PLAIN authentication mechanism | [ ] | S-002, U-001 |
| S-007 | LOGIN authentication mechanism | [ ] | S-002, U-001 |
| S-008 | CRAM-MD5 authentication mechanism | [ ] | S-002, U-001 |
| S-009 | SIZE extension (RFC 1870) | [ ] | S-001 |
| S-010 | 8BITMIME support (RFC 6152) | [ ] | S-001 |
| S-011 | PIPELINING support (RFC 2920) | [ ] | S-001 |
| S-012 | CHUNKING support (RFC 3030) | [ ] | S-001 |

### S-001: go-smtp Integration

```go
// internal/smtp/server.go
package smtp

import (
    "github.com/emersion/go-smtp"
)

type Server struct {
    submission *smtp.Server  // Port 587
    relay      *smtp.Server  // Port 25
    smtps      *smtp.Server  // Port 465
    backend    *Backend
}

func NewServer(cfg *config.SMTP, backend *Backend) *Server {
    return &Server{
        backend: backend,
    }
}

func (s *Server) Start() error {
    // Start all listeners
}
```

**Acceptance Criteria**:
- [ ] Port 25 (relay/MX): Accept inbound mail from any server
- [ ] Port 587 (submission): Require STARTTLS + authentication
- [ ] Port 465 (SMTPS): Implicit TLS, require authentication
- [ ] Maximum 100 concurrent connections per port
- [ ] Idle timeout: 5 minutes per SMTP RFC
- [ ] Command timeout: 30 seconds per command
- [ ] Graceful shutdown: complete active sessions before exit
- [ ] TLS 1.2+ required for ports 465 and 587
- [ ] PIPELINING extension support
- [ ] 8BITMIME extension support

**Structured Logging (slog)**:
- [ ] **INFO**: SMTP server started (ports=[25,465,587], tls_enabled, max_connections=100, version)
- [ ] **INFO**: SMTP connection accepted (port, client_ip, connection_id, tls_active, protocol="SMTP|SMTPS")
- [ ] **DEBUG**: SMTP command received (connection_id, command="EHLO|MAIL|RCPT|DATA", session_id, client_ip)
- [ ] **INFO**: SMTP authentication successful (connection_id, username, auth_method="PLAIN|LOGIN", session_id, client_ip)
- [ ] **WARN**: SMTP authentication failed (connection_id, username, error_code="535", client_ip, session_id)
- [ ] **INFO**: Message accepted for delivery (connection_id, message_id, from, recipients_count, size_bytes, session_id)
- [ ] **ERROR**: SMTP connection limit reached (port, current_connections=100, rejected_ip, max_connections=100)
- [ ] **INFO**: SMTP server shutdown (graceful, active_sessions_completed, total_duration_ms)
- [ ] **TRACE**: STARTTLS negotiation (connection_id, cipher_suite, tls_version, client_ip)
- [ ] **Fields**: connection_id, session_id, client_ip, port, command, message_id, from, recipients_count, size_bytes, duration_ms

**Given/When/Then Scenarios**:
```
Given SMTP server is configured for ports 25, 465, 587
When server starts
Then port 25 listener starts without TLS requirement
And port 465 listener starts with implicit TLS
And port 587 listener starts with STARTTLS advertised
And all listeners accept connections
And server logs "SMTP server started on ports 25, 465, 587"

Given client connects to port 587 (submission)
When EHLO command is sent
Then server responds with "250-STARTTLS"
And server responds with "250-AUTH PLAIN LOGIN"
And STARTTLS is required before authentication
And unauthenticated clients cannot send mail

Given authenticated user connects to port 587 with TLS
When message is submitted via DATA command
Then message is accepted
And queued for delivery
And server responds "250 2.0.0 Message accepted for delivery"

Given 100 concurrent connections exist on port 25
When 101st connection is attempted
Then connection is rejected
And server responds "421 4.3.2 Too many connections, try again later"
And connection is closed immediately

Given client connects to port 465 (SMTPS)
When connection is established
Then TLS handshake occurs immediately (implicit TLS)
And no STARTTLS command is needed
And connection is encrypted from start

Given SMTP server receives SIGTERM signal
When graceful shutdown is initiated
Then new connections are rejected
And existing sessions complete current transactions
And server waits up to 30 seconds for completion
And server exits cleanly with all messages queued
```

### S-002: SMTP Backend Implementation

```go
// internal/smtp/backend.go
package smtp

import (
    "github.com/emersion/go-smtp"
)

type Backend struct {
    userService    *service.UserService
    messageService *service.MessageService
    queueService   *service.QueueService
}

func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
    return &Session{
        conn:    c,
        backend: b,
    }, nil
}

type Session struct {
    conn    *smtp.Conn
    backend *Backend
    from    string
    to      []string
    user    *domain.User
}

func (s *Session) AuthPlain(username, password string) error {
    user, err := s.backend.userService.Authenticate(username, password)
    if err != nil {
        return &smtp.SMTPError{Code: 535, Message: "Authentication failed"}
    }
    s.user = user
    return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
    s.from = from
    return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
    s.to = append(s.to, to)
    return nil
}

func (s *Session) Data(r io.Reader) error {
    // Parse message and queue for delivery
}

func (s *Session) Reset() {
    s.from = ""
    s.to = nil
}

func (s *Session) Logout() error {
    return nil
}
```

---

## 1.2 SMTP Queue Management [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| Q-001 | Design queue table schema | [ ] | F-022 |
| Q-002 | Implement message queuing service | [ ] | Q-001 |
| Q-003 | Implement retry logic with exponential backoff | [ ] | Q-002 |
| Q-004 | Implement bounce handling | [ ] | Q-002 |
| Q-005 | DSN (Delivery Status Notifications) | [ ] | Q-002 |
| Q-006 | Queue cleanup and maintenance | [ ] | Q-002 |

### Q-002: Queue Service

```go
// internal/service/queue_service.go
package service

type QueueService struct {
    repo   repository.QueueRepository
    logger logger.Logger
}

type QueueItem struct {
    ID          int64
    Sender      string
    Recipients  []string
    Message     []byte
    RetryCount  int
    NextRetry   time.Time
    Status      string
    ErrorMsg    string
}

func (s *QueueService) Enqueue(sender string, recipients []string, message []byte) error {
    item := &QueueItem{
        Sender:     sender,
        Recipients: recipients,
        Message:    message,
        Status:     "pending",
        NextRetry:  time.Now(),
    }
    return s.repo.Create(item)
}

func (s *QueueService) ProcessQueue(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.processItems()
        }
    }
}

func (s *QueueService) processItems() {
    items, _ := s.repo.GetPending()
    for _, item := range items {
        if err := s.deliver(item); err != nil {
            s.handleFailure(item, err)
        } else {
            s.repo.MarkSent(item.ID)
        }
    }
}
```

### Q-003: Exponential Backoff

```go
func (s *QueueService) calculateNextRetry(retryCount int) time.Time {
    // Backoff: 5m, 15m, 30m, 1h, 2h, 4h, 8h, 16h, 24h
    delays := []time.Duration{
        5 * time.Minute,
        15 * time.Minute,
        30 * time.Minute,
        1 * time.Hour,
        2 * time.Hour,
        4 * time.Hour,
        8 * time.Hour,
        16 * time.Hour,
        24 * time.Hour,
    }
    if retryCount >= len(delays) {
        return time.Time{} // Give up
    }
    return time.Now().Add(delays[retryCount])
}
```

**Acceptance Criteria**:
- [ ] Exponential backoff schedule: 5m, 15m, 30m, 1h, 2h, 4h, 8h, 16h, 24h
- [ ] Maximum 9 retry attempts before permanent failure
- [ ] Retry count persisted across server restarts
- [ ] Failed deliveries marked with specific error codes

**Structured Logging (slog)**:
- [ ] **INFO**: Queue enqueue (message_id, sender, recipient_count, queue_size)
- [ ] **WARN**: Delivery failure with retry (message_id, retry_count, next_retry, error_msg, duration_ms)
- [ ] **ERROR**: Permanent delivery failure (message_id, sender, recipients, retry_count=9, error_msg)
- [ ] **DEBUG**: Retry calculation (message_id, retry_count, calculated_delay, next_retry_time)
- [ ] **INFO**: Successful delivery (message_id, sender, recipient, retry_count, total_duration_ms)
- [ ] **TRACE**: Queue processing (items_pending, items_processed, processing_duration_ms)
- [ ] **Fields**: message_id, sender, recipients[], retry_count, next_retry, status, error_msg, duration_ms, queue_size

**Given/When/Then Scenarios**:
```
Given a message fails delivery for the first time
When retry is scheduled
Then next attempt is in 5 minutes (retry count 0)

Given a message has failed 3 times (retry count 2)
When retry is scheduled
Then next attempt is in 30 minutes

Given a message has failed 9 times (retry count 8)
When retry is scheduled
Then next attempt is in 24 hours (maximum delay)

Given a message has failed 10 times (retry count 9)
When retry scheduling is attempted
Then message is marked as permanently failed
And delivery is abandoned
And bounce notification is generated

Given server restart occurs mid-retry
When server starts up
Then retry schedule resumes from persisted retry count
And no retries are skipped or duplicated
```

---

## 1.3 Message Storage [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| M-001 | Integrate go-message for MIME parsing | [ ] | F-002 | `emersion/go-message` |
| M-002 | Implement hybrid storage (blob < 1MB, file >= 1MB) | [ ] | F-022 |
| M-003 | Message header parsing and indexing | [ ] | M-001 |
| M-004 | Attachment handling | [ ] | M-001 |
| M-005 | Message deduplication | [ ] | M-003 |
| M-006 | Thread ID generation for conversations | [ ] | M-003 |

### M-002: Hybrid Storage

```go
// internal/storage/message_store.go
package storage

const BlobThreshold = 1024 * 1024 // 1MB

type MessageStore struct {
    db       *database.DB
    basePath string
}

func (s *MessageStore) Store(userID int64, mailboxID int64, message []byte) (*domain.Message, error) {
    msg := &domain.Message{
        UserID:    userID,
        MailboxID: mailboxID,
        Size:      len(message),
    }

    if len(message) < BlobThreshold {
        msg.StorageType = "blob"
        msg.Content = message
    } else {
        msg.StorageType = "file"
        path := s.generatePath(userID, msg.ID)
        if err := os.WriteFile(path, message, 0600); err != nil {
            return nil, err
        }
        msg.FilePath = path
    }

    return msg, s.repo.Create(msg)
}

func (s *MessageStore) Retrieve(msg *domain.Message) ([]byte, error) {
    if msg.StorageType == "blob" {
        return msg.Content, nil
    }
    return os.ReadFile(msg.FilePath)
}
```

### M-003: Header Parsing

```go
// internal/service/message_service.go
package service

import (
    "github.com/emersion/go-message"
    "github.com/emersion/go-message/mail"
)

func (s *MessageService) ParseMessage(raw []byte) (*domain.Message, error) {
    r := bytes.NewReader(raw)
    m, err := mail.ReadMessage(r)
    if err != nil {
        return nil, err
    }

    header := m.Header
    msg := &domain.Message{
        MessageID:  header.Get("Message-ID"),
        Subject:    header.Get("Subject"),
        Sender:     header.Get("From"),
        Recipients: header.Get("To"),
        Date:       parseDate(header.Get("Date")),
        Headers:    headerToJSON(header),
    }

    // Generate thread ID from References/In-Reply-To
    msg.ThreadID = s.generateThreadID(header)

    return msg, nil
}
```

### M-006: Thread ID Generation

```go
func (s *MessageService) generateThreadID(header mail.Header) string {
    // Check In-Reply-To first
    if inReplyTo := header.Get("In-Reply-To"); inReplyTo != "" {
        // Find existing thread
        if thread := s.findThreadByMessageID(inReplyTo); thread != "" {
            return thread
        }
    }

    // Check References
    if refs := header.Get("References"); refs != "" {
        // Parse references, find existing thread
        for _, ref := range parseReferences(refs) {
            if thread := s.findThreadByMessageID(ref); thread != "" {
                return thread
            }
        }
    }

    // New thread - use Message-ID as thread ID
    return header.Get("Message-ID")
}
```

---

## 1.4 IMAP Server [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| I-001 | Integrate go-imap library | [ ] | F-002 | `emersion/go-imap` |
| I-002 | Implement IMAP backend interface | [ ] | I-001, M-002 |
| I-003 | Mailbox operations (CREATE, DELETE, RENAME) | [ ] | I-002 |
| I-004 | Message operations (FETCH, STORE, COPY) | [ ] | I-002 |
| I-005 | STARTTLS support | [ ] | I-001, T-001 |
| I-006 | PLAIN authentication | [ ] | I-001, U-001 |
| I-007 | LOGIN authentication | [ ] | I-001, U-001 |
| I-008 | IDLE support (RFC 2177) | [ ] | I-002 |
| I-009 | UIDPLUS extension (RFC 4315) | [ ] | I-002 |
| I-010 | QUOTA extension (RFC 2087) | [ ] | I-002 |
| I-011 | SORT extension (RFC 5256) | [ ] | I-002 |
| I-012 | NAMESPACE extension (RFC 2342) | [ ] | I-002 |
| I-013 | Special-use mailboxes (RFC 6154) | [ ] | I-002 |
| I-014 | Server-side search | [ ] | I-002 |

### I-001: go-imap Server Setup

```go
// internal/imap/server.go
package imap

import (
    "github.com/emersion/go-imap/server"
)

type Server struct {
    server  *server.Server
    backend *Backend
}

func NewServer(cfg *config.IMAP, backend *Backend) *Server {
    s := server.New(backend)
    s.Addr = fmt.Sprintf(":%d", cfg.Port)
    s.AllowInsecureAuth = false

    return &Server{
        server:  s,
        backend: backend,
    }
}
```

**Acceptance Criteria**:
- [ ] Port 143 (IMAP): Require STARTTLS before authentication
- [ ] Port 993 (IMAPS): Implicit TLS from connection start
- [ ] AllowInsecureAuth = false (no plaintext passwords without TLS)
- [ ] Maximum 200 concurrent IMAP connections
- [ ] Idle timeout: 30 minutes per IMAP RFC
- [ ] Command timeout: 60 seconds per command
- [ ] IDLE extension: 29-minute keepalive heartbeat
- [ ] TLS 1.2+ required for both ports
- [ ] LOGIN and PLAIN authentication methods
- [ ] Graceful shutdown: complete active sessions

**Structured Logging (slog)**:
- [ ] **INFO**: IMAP server started (ports=[143,993], tls_enabled, max_connections=200, version)
- [ ] **INFO**: IMAP connection accepted (port, client_ip, connection_id, tls_active, protocol="IMAP|IMAPS")
- [ ] **DEBUG**: IMAP command received (connection_id, command="LOGIN|SELECT|FETCH|IDLE", session_id, client_ip, mailbox)
- [ ] **INFO**: IMAP authentication successful (connection_id, username, auth_method="PLAIN|LOGIN", session_id, client_ip)
- [ ] **WARN**: IMAP authentication failed (connection_id, username, error_code, client_ip, session_id)
- [ ] **INFO**: IMAP IDLE mode entered (connection_id, mailbox, session_id, heartbeat_interval="29m")
- [ ] **DEBUG**: IMAP IDLE notification (connection_id, event="new_message", message_count, session_id)
- [ ] **WARN**: IMAP idle timeout (connection_id, idle_duration="30m", session_id, client_ip)
- [ ] **ERROR**: IMAP connection limit reached (port, current_connections=200, rejected_ip, max_connections=200)
- [ ] **INFO**: IMAP server shutdown (graceful, active_sessions_completed, total_duration_ms)
- [ ] **TRACE**: STARTTLS completed (connection_id, cipher_suite, tls_version, client_ip)
- [ ] **Fields**: connection_id, session_id, client_ip, port, command, mailbox, message_count, username, duration_ms

**Given/When/Then Scenarios**:
```
Given IMAP server is configured for ports 143 and 993
When server starts
Then port 143 listener starts advertising STARTTLS capability
And port 993 listener starts with implicit TLS
And server logs "IMAP server started on ports 143, 993"
And authentication is disabled until TLS is active

Given client connects to port 143 without TLS
When CAPABILITY command is sent
Then server responds with "* CAPABILITY IMAP4rev1 STARTTLS"
And STARTTLS is advertised as available
And LOGIN command is rejected until STARTTLS completes

Given client connects to port 143 and issues STARTTLS
When TLS handshake completes
Then secure connection is established
And CAPABILITY now includes "AUTH=PLAIN" and "AUTH=LOGIN"
And LOGIN command is now accepted

Given authenticated user on port 993
When IDLE command is issued
Then server enters IDLE mode
And sends "* OK Still here" heartbeat every 29 minutes
And notifies client of new messages immediately
And client can exit IDLE with DONE command

Given client is authenticated and idle for 30 minutes
When idle timeout is reached
Then server sends "* BYE Idle timeout"
And connection is closed gracefully
And client must reconnect to continue

Given 200 concurrent IMAP connections exist
When 201st connection is attempted
Then connection is rejected
And server responds "* BYE Maximum connections reached"
And connection is closed immediately

Given IMAP server receives SIGTERM signal
When graceful shutdown is initiated
Then new connections are rejected
And existing sessions complete current commands
And server waits up to 60 seconds for completion
And all connections close cleanly
```

### I-002: IMAP Backend Interface

```go
// internal/imap/backend.go
package imap

import (
    "github.com/emersion/go-imap/backend"
)

type Backend struct {
    userService    *service.UserService
    mailboxService *service.MailboxService
    messageService *service.MessageService
}

func (b *Backend) Login(connInfo *backend.ConnInfo, username, password string) (backend.User, error) {
    user, err := b.userService.Authenticate(username, password)
    if err != nil {
        return nil, backend.ErrInvalidCredentials
    }
    return &User{
        user:           user,
        mailboxService: b.mailboxService,
        messageService: b.messageService,
    }, nil
}

type User struct {
    user           *domain.User
    mailboxService *service.MailboxService
    messageService *service.MessageService
}

func (u *User) Username() string {
    return u.user.Email
}

func (u *User) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
    mailboxes, err := u.mailboxService.List(u.user.ID, subscribed)
    if err != nil {
        return nil, err
    }
    result := make([]backend.Mailbox, len(mailboxes))
    for i, mb := range mailboxes {
        result[i] = &Mailbox{
            mailbox:        mb,
            messageService: u.messageService,
        }
    }
    return result, nil
}

func (u *User) GetMailbox(name string) (backend.Mailbox, error) {
    mb, err := u.mailboxService.Get(u.user.ID, name)
    if err != nil {
        return nil, backend.ErrNoSuchMailbox
    }
    return &Mailbox{
        mailbox:        mb,
        messageService: u.messageService,
    }, nil
}

func (u *User) CreateMailbox(name string) error {
    return u.mailboxService.Create(u.user.ID, name)
}

func (u *User) DeleteMailbox(name string) error {
    return u.mailboxService.Delete(u.user.ID, name)
}

func (u *User) RenameMailbox(oldName, newName string) error {
    return u.mailboxService.Rename(u.user.ID, oldName, newName)
}

func (u *User) Logout() error {
    return nil
}
```

### I-013: Special-Use Mailboxes

```go
// Create default mailboxes for new users
func (s *MailboxService) CreateDefaultMailboxes(userID int64) error {
    defaults := []struct {
        Name       string
        SpecialUse string
    }{
        {"INBOX", ""},
        {"Drafts", "\\Drafts"},
        {"Sent", "\\Sent"},
        {"Trash", "\\Trash"},
        {"Spam", "\\Junk"},
        {"Archive", "\\Archive"},
    }

    for _, mb := range defaults {
        if err := s.CreateWithSpecialUse(userID, mb.Name, mb.SpecialUse); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 1.5 User Management [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| U-001 | Users table schema and repository | [ ] | F-022 |
| U-002 | Bcrypt password hashing | [ ] | U-001 |
| U-003 | Domains table schema and repository | [ ] | F-022 |
| U-004 | Aliases table schema and repository | [ ] | F-022 |
| U-005 | Mailboxes/folders table and repository | [ ] | F-022 |
| U-006 | Quota tracking and enforcement | [ ] | U-001, M-002 |

### U-001: User Repository

```go
// internal/repository/user_repository.go
package repository

type UserRepository interface {
    Create(user *domain.User) error
    GetByID(id int64) (*domain.User, error)
    GetByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int64) error
    List(domainID int64, offset, limit int) ([]*domain.User, error)
    UpdateQuota(id int64, usedQuota int64) error
}

type sqliteUserRepository struct {
    db *database.DB
}

func (r *sqliteUserRepository) Create(user *domain.User) error {
    query := `
        INSERT INTO users (email, domain_id, password_hash, full_name, display_name, quota)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    result, err := r.db.Exec(query, user.Email, user.DomainID, user.PasswordHash,
        user.FullName, user.DisplayName, user.Quota)
    if err != nil {
        return err
    }
    id, _ := result.LastInsertId()
    user.ID = id
    return nil
}
```

### U-002: Password Hashing

```go
// internal/service/user_service.go
package service

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12

func (s *UserService) HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func (s *UserService) VerifyPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *UserService) Authenticate(email, password string) (*domain.User, error) {
    user, err := s.repo.GetByEmail(email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    if user.Status != "active" {
        return nil, ErrUserDisabled
    }

    if !s.VerifyPassword(user.PasswordHash, password) {
        return nil, ErrInvalidCredentials
    }

    // Update last login
    s.repo.UpdateLastLogin(user.ID)

    return user, nil
}
```

**Acceptance Criteria**:
- [ ] bcrypt with cost factor 12 for password hashing
- [ ] Generic error messages to prevent user enumeration
- [ ] Last login timestamp updated on successful authentication
- [ ] Disabled users rejected even with correct credentials

**Structured Logging (slog)**:
- [ ] **INFO**: Successful authentication (user_id, email, ip_address, auth_method="password", session_id, duration_ms)
- [ ] **WARN**: Failed authentication - wrong password (email, ip_address, auth_method, error="invalid_credentials", session_id)
- [ ] **WARN**: Failed authentication - disabled user (user_id, email, ip_address, status="inactive", session_id)
- [ ] **ERROR**: Failed authentication - user not found (email, ip_address, error="user_not_found", session_id)
- [ ] **DEBUG**: Password hash generation (user_id, bcrypt_cost=12, duration_ms)
- [ ] **TRACE**: Authentication attempt start (email, ip_address, auth_method, request_id, session_id)
- [ ] **Fields**: user_id, email, ip_address, auth_method, session_id, request_id, trace_id, event_type, status, error_code, duration_ms

**Given/When/Then Scenarios**:
```
Given a new user account is created with password "SecurePass123"
When password is hashed
Then hash is generated using bcrypt cost 12
And original password cannot be recovered from hash
And hash verification succeeds with correct password

Given user "user@example.com" exists with correct password
When authentication is attempted with correct credentials
Then authentication succeeds
And user object is returned
And last_login timestamp is updated

Given user "user@example.com" exists
When authentication is attempted with incorrect password
Then authentication fails with "Invalid credentials" error
And error message does not reveal that email exists
And last_login timestamp is NOT updated

Given user "disabled@example.com" has status "inactive"
When authentication is attempted with correct credentials
Then authentication fails with "User disabled" error
And user is not authenticated

Given user "nonexistent@example.com" does not exist
When authentication is attempted
Then authentication fails with "Invalid credentials" error
And error message is identical to wrong password case (prevents enumeration)
```

---

### U-006: Quota Enforcement

```go
// internal/service/quota_service.go
package service

type QuotaService struct {
    userRepo    repository.UserRepository
    messageRepo repository.MessageRepository
}

func (s *QuotaService) CheckQuota(userID int64, additionalSize int64) error {
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return err
    }

    if user.UsedQuota+additionalSize > user.Quota {
        return ErrQuotaExceeded
    }
    return nil
}

func (s *QuotaService) UpdateUsage(userID int64) error {
    total, err := s.messageRepo.GetTotalSize(userID)
    if err != nil {
        return err
    }
    return s.userRepo.UpdateQuota(userID, total)
}

func (s *QuotaService) GetUsage(userID int64) (*QuotaInfo, error) {
    user, _ := s.userRepo.GetByID(userID)
    return &QuotaInfo{
        Used:       user.UsedQuota,
        Limit:      user.Quota,
        Percentage: float64(user.UsedQuota) / float64(user.Quota) * 100,
    }, nil
}
```

**Acceptance Criteria**:
- [ ] Default quota: 1GB (1073741824 bytes) per user
- [ ] Quota check before accepting messages (SMTP RCPT TO stage)
- [ ] Quota updated atomically after message storage
- [ ] Warning notifications at 80%, 90%, 95% thresholds
- [ ] Over-quota rejection returns SMTP 552 "Mailbox full"
- [ ] Quota calculation includes message size + metadata overhead
- [ ] Admin API to view/modify user quotas

**Structured Logging (slog)**:
- [ ] **INFO**: Quota check passed (user_id, email, used_bytes, limit_bytes, percentage, message_size, remaining_bytes)
- [ ] **WARN**: Quota threshold exceeded (user_id, email, threshold="80%|90%|95%", used_bytes, limit_bytes, percentage)
- [ ] **ERROR**: Quota exceeded - message rejected (user_id, email, used_bytes, limit_bytes, message_size, overflow_bytes)
- [ ] **INFO**: Quota updated after storage (user_id, email, old_quota, new_quota, message_id, message_size, duration_ms)
- [ ] **DEBUG**: Quota recalculation (user_id, email, message_count, total_size, metadata_overhead, final_quota)
- [ ] **FATAL**: Quota update failed - data integrity issue (user_id, email, error_msg, transaction_id)
- [ ] **Fields**: user_id, email, used_bytes, limit_bytes, percentage, message_size, message_id, threshold, duration_ms

**Given/When/Then Scenarios**:
```
Given user has quota 1GB and currently uses 500MB
When 200MB message arrives
Then quota check passes
And message is accepted
And used quota updates to 700MB

Given user has quota 1GB and currently uses 1000MB
When 100MB message arrives
Then quota check fails with "Mailbox full"
And SMTP returns 552 error code
And message is rejected at RCPT TO stage
And sender receives bounce notification

Given user has quota 1GB and currently uses 800MB
When quota usage is calculated
Then warning level is "80%" (first threshold)
And notification email is queued
And user is alerted via webmail interface

Given user has quota 1GB and currently uses 970MB
When quota usage is calculated
Then warning level is "97%" (>95% critical threshold)
And urgent notification is sent
And admin dashboard shows critical quota status

Given user has 1000 messages totaling 950MB
When quota recalculation is triggered
Then total size is computed from message table
And metadata overhead (headers, indexes) is added (~5%)
And final used quota is updated to ~997MB
```

---

## 1.6 TLS Support [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| T-001 | TLS configuration loader | [ ] | F-011 |
| T-002 | SNI support for multi-domain | [ ] | T-001 |
| T-003 | TLS 1.2+ enforcement | [ ] | T-001 |
| T-004 | Modern cipher suite configuration | [ ] | T-001 |

### T-001: TLS Configuration

```go
// internal/tls/config.go
package tls

import (
    "crypto/tls"
)

type Manager struct {
    certs map[string]*tls.Certificate
    cfg   *config.TLS
}

func NewManager(cfg *config.TLS) (*Manager, error) {
    m := &Manager{
        certs: make(map[string]*tls.Certificate),
        cfg:   cfg,
    }

    if err := m.loadCertificates(); err != nil {
        return nil, err
    }

    return m, nil
}

func (m *Manager) GetConfig() *tls.Config {
    return &tls.Config{
        GetCertificate:     m.getCertificate,
        MinVersion:         tls.VersionTLS12,
        CipherSuites:       modernCipherSuites(),
        PreferServerCipherSuites: true,
    }
}

func (m *Manager) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
    // SNI lookup
    if cert, ok := m.certs[hello.ServerName]; ok {
        return cert, nil
    }
    // Return default cert
    return m.certs["default"], nil
}

func modernCipherSuites() []uint16 {
    return []uint16{
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
    }
}
```

---

## Acceptance Criteria

### SMTP
- [ ] Can receive email via SMTP relay (port 25)
- [ ] Can send email via SMTP submission (port 587) with auth
- [ ] SMTPS (port 465) works with TLS
- [ ] STARTTLS upgrade works
- [ ] All auth mechanisms work (PLAIN, LOGIN, CRAM-MD5)
- [ ] Extensions advertised correctly (SIZE, 8BITMIME, PIPELINING)

### IMAP
- [ ] Can authenticate and list mailboxes
- [ ] FETCH returns correct message data
- [ ] STORE updates flags correctly
- [ ] COPY moves messages between mailboxes
- [ ] CREATE/DELETE/RENAME mailboxes work
- [ ] IDLE notifies on new messages
- [ ] Search returns correct results
- [ ] Special-use mailboxes created for new users

### Queue
- [ ] Messages queued for delivery
- [ ] Retry logic with exponential backoff works
- [ ] Bounces generated for permanent failures
- [ ] Queue can be inspected and managed

### Storage
- [ ] Small messages stored as blobs
- [ ] Large messages stored as files
- [ ] Headers parsed and indexed
- [ ] Thread IDs generated correctly

---

## Go Dependencies for Phase 1

```go
// Additional go.mod entries
require (
    github.com/emersion/go-smtp v0.21.0
    github.com/emersion/go-imap v1.2.1
    github.com/emersion/go-message v0.18.0
    github.com/emersion/go-sasl v0.0.0-20231106173351-e73c9f7bad43
    golang.org/x/crypto v0.19.0
)
```

---

## Testing Commands

```bash
# Test SMTP submission
swaks --to user@example.com --from sender@example.com \
      --server localhost:587 --auth PLAIN \
      --auth-user test@example.com --auth-password secret

# Test IMAP
openssl s_client -connect localhost:993
# Then: A LOGIN user@example.com password
# Then: A SELECT INBOX
# Then: A FETCH 1:* FLAGS

# Test STARTTLS
openssl s_client -starttls smtp -connect localhost:587
```

---

## Next Phase

After completing Phase 1, proceed to [TASKS2.md](TASKS2.md) - Security Foundation.
