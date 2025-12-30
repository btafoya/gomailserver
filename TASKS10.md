# TASKS10.md - Phase 10: Testing (Weeks 30-31)

## Overview

Comprehensive testing suite covering unit tests, integration tests, external validation, performance benchmarks, and security audits.

**Total Tasks**: 20
**MVP Tasks**: 9 (Test Data + Unit + Integration tests)
**Priority**: Mixed - Unit/Integration MVP, others FULL (Chaos: OPTIONAL)
**Dependencies**: All previous phases

---

## 10.0 Test Data Management [MVP]

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| TDM-001 | Test fixtures and factories | [ ] | - | MVP |

---

### TDM-001: Test Fixtures and Factories

**File**: `internal/testutil/fixtures/user.go`
```go
package fixtures

import (
    "github.com/btafoya/gomailserver/internal/domain"
    "time"
)

// CreateTestUser creates a user for testing with sensible defaults
func CreateTestUser(overrides ...func(*domain.User)) *domain.User {
    user := &domain.User{
        Email:        "test@example.com",
        PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "password"
        DomainID:     1,
        Active:       true,
        QuotaBytes:   1073741824, // 1GB
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    for _, override := range overrides {
        override(user)
    }

    return user
}

// CreateTestMessage creates a message for testing
func CreateTestMessage(overrides ...func(*domain.Message)) *domain.Message {
    msg := &domain.Message{
        From:        "sender@example.com",
        To:          []string{"recipient@example.com"},
        Subject:     "Test Message",
        MessageID:   "<test-msg-id@example.com>",
        Size:        1024,
        ReceivedAt:  time.Now(),
    }

    for _, override := range overrides {
        override(msg)
    }

    return msg
}

// CreateTestMailbox creates a mailbox for testing
func CreateTestMailbox(overrides ...func(*domain.Mailbox)) *domain.Mailbox {
    mb := &domain.Mailbox{
        Name:      "INBOX",
        UserID:    1,
        UIDNext:   1,
        CreatedAt: time.Now(),
    }

    for _, override := range overrides {
        override(mb)
    }

    return mb
}
```

**File**: `testdata/fixtures.sql`
```sql
-- Test domain
INSERT INTO domains (id, name, active) VALUES (1, 'example.com', 1);

-- Test user
INSERT INTO users (id, email, password_hash, domain_id, active, quota_bytes, created_at, updated_at)
VALUES (1, 'test@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, 1, 1073741824, datetime('now'), datetime('now'));

-- Test mailbox
INSERT INTO mailboxes (id, name, user_id, uid_next, created_at)
VALUES (1, 'INBOX', 1, 1, datetime('now'));
```

**File**: `internal/testutil/db.go`
```go
package testutil

import (
    "database/sql"
    "testing"

    _ "github.com/mattn/go-sqlite3"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open test database: %v", err)
    }

    // Run migrations
    if err := runMigrations(db); err != nil {
        t.Fatalf("failed to run migrations: %v", err)
    }

    return db
}

// CleanDatabase truncates all tables for test isolation
func CleanDatabase(t *testing.T, db *sql.DB) {
    t.Helper()

    tables := []string{
        "messages", "mailboxes", "users", "domains",
        "aliases", "dkim_keys", "webhooks", "webhook_deliveries",
    }

    for _, table := range tables {
        _, err := db.Exec("DELETE FROM " + table)
        if err != nil {
            t.Fatalf("failed to clean table %s: %v", table, err)
        }
    }
}
```

**Acceptance Criteria**:
- [ ] Factory functions available for all domain entities (User, Message, Mailbox, Domain, Alias)
- [ ] All tests use factories instead of inline test data creation
- [ ] Integration tests can run in parallel with isolated database state
- [ ] Test data is realistic (valid email addresses, RFC-compliant headers, proper MIME structure)
- [ ] `CleanDatabase()` helper available for resetting state between tests
- [ ] SQL fixtures available for complex integration test scenarios

**Production Readiness**:
- Test isolation: Each test has clean database state
- Parallel execution: Tests can run concurrently without conflicts
- Realistic data: Matches production data patterns and constraints

---

## 10.1 Unit Tests [MVP]

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| UT-001 | Repository layer tests | [ ] | F-023 | MVP |
| UT-002 | Service layer tests | [ ] | All services | MVP |
| UT-003 | Security function tests | [ ] | Phase 2 | MVP |
| UT-004 | 80%+ code coverage | [ ] | UT-001-003 | MVP |

---

### UT-001: Repository Layer Tests
**File**: `internal/repository/user_repository_test.go`
```go
package repository

import (
    "database/sql"
    "testing"
    "time"

    "github.com/btafoya/gomailserver/internal/domain"
    _ "github.com/mattn/go-sqlite3"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
    suite.Suite
    db   *sql.DB
    repo *SQLiteUserRepository
}

func (s *UserRepositoryTestSuite) SetupTest() {
    // Create in-memory database
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(s.T(), err)

    // Run migrations
    err = runMigrations(db)
    require.NoError(s.T(), err)

    s.db = db
    s.repo = NewSQLiteUserRepository(db)
}

func (s *UserRepositoryTestSuite) TearDownTest() {
    s.db.Close()
}

func (s *UserRepositoryTestSuite) TestCreateUser() {
    user := &domain.User{
        Email:        "test@example.com",
        PasswordHash: "hashed_password",
        DomainID:     1,
        Active:       true,
        QuotaBytes:   1073741824, // 1GB
    }

    err := s.repo.Create(user)
    assert.NoError(s.T(), err)
    assert.NotZero(s.T(), user.ID)

    // Verify user was created
    found, err := s.repo.GetByID(user.ID)
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), user.Email, found.Email)
    assert.Equal(s.T(), user.Active, found.Active)
}

func (s *UserRepositoryTestSuite) TestGetByEmail() {
    // Create test user
    user := &domain.User{
        Email:        "lookup@example.com",
        PasswordHash: "hash",
        DomainID:     1,
        Active:       true,
    }
    s.repo.Create(user)

    // Test lookup
    found, err := s.repo.GetByEmail("lookup@example.com")
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), user.ID, found.ID)

    // Test not found
    _, err = s.repo.GetByEmail("notfound@example.com")
    assert.Equal(s.T(), domain.ErrNotFound, err)
}

func (s *UserRepositoryTestSuite) TestUpdateUser() {
    user := &domain.User{
        Email:        "update@example.com",
        PasswordHash: "old_hash",
        DomainID:     1,
        Active:       true,
    }
    s.repo.Create(user)

    // Update
    user.PasswordHash = "new_hash"
    user.Active = false
    err := s.repo.Update(user)
    assert.NoError(s.T(), err)

    // Verify
    found, _ := s.repo.GetByID(user.ID)
    assert.Equal(s.T(), "new_hash", found.PasswordHash)
    assert.False(s.T(), found.Active)
}

func (s *UserRepositoryTestSuite) TestDeleteUser() {
    user := &domain.User{
        Email:        "delete@example.com",
        PasswordHash: "hash",
        DomainID:     1,
        Active:       true,
    }
    s.repo.Create(user)

    err := s.repo.Delete(user.ID)
    assert.NoError(s.T(), err)

    // Verify deleted
    _, err = s.repo.GetByID(user.ID)
    assert.Equal(s.T(), domain.ErrNotFound, err)
}

func (s *UserRepositoryTestSuite) TestListByDomain() {
    // Create multiple users
    for i := 0; i < 5; i++ {
        s.repo.Create(&domain.User{
            Email:    fmt.Sprintf("user%d@domain1.com", i),
            DomainID: 1,
            Active:   true,
        })
    }
    for i := 0; i < 3; i++ {
        s.repo.Create(&domain.User{
            Email:    fmt.Sprintf("user%d@domain2.com", i),
            DomainID: 2,
            Active:   true,
        })
    }

    // List domain 1
    users, total, err := s.repo.ListByDomain(1, 10, 0)
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), 5, total)
    assert.Len(s.T(), users, 5)

    // List domain 2
    users, total, err = s.repo.ListByDomain(2, 10, 0)
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), 3, total)
}

func (s *UserRepositoryTestSuite) TestUniqueEmailConstraint() {
    user1 := &domain.User{
        Email:    "duplicate@example.com",
        DomainID: 1,
    }
    s.repo.Create(user1)

    user2 := &domain.User{
        Email:    "duplicate@example.com",
        DomainID: 1,
    }
    err := s.repo.Create(user2)
    assert.Error(s.T(), err)
    assert.Contains(s.T(), err.Error(), "UNIQUE constraint")
}

func TestUserRepositorySuite(t *testing.T) {
    suite.Run(t, new(UserRepositoryTestSuite))
}
```

**File**: `internal/repository/message_repository_test.go`
```go
package repository

import (
    "testing"
    "time"

    "github.com/btafoya/gomailserver/internal/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type MessageRepositoryTestSuite struct {
    suite.Suite
    db   *sql.DB
    repo *SQLiteMessageRepository
}

func (s *MessageRepositoryTestSuite) SetupTest() {
    db, _ := sql.Open("sqlite3", ":memory:")
    runMigrations(db)
    s.db = db
    s.repo = NewSQLiteMessageRepository(db)
}

func (s *MessageRepositoryTestSuite) TestStoreMessage() {
    msg := &domain.Message{
        UserID:     1,
        MailboxID:  1,
        MessageID:  "<test123@example.com>",
        Subject:    "Test Subject",
        From:       "sender@example.com",
        To:         []string{"recipient@example.com"},
        Date:       time.Now(),
        Size:       1024,
        Flags:      []string{"\\Seen"},
    }

    err := s.repo.Store(msg)
    assert.NoError(s.T(), err)
    assert.NotZero(s.T(), msg.ID)
    assert.NotZero(s.T(), msg.UID)
}

func (s *MessageRepositoryTestSuite) TestFetchByUID() {
    msg := &domain.Message{
        UserID:    1,
        MailboxID: 1,
        MessageID: "<fetch123@example.com>",
        Subject:   "Fetch Test",
    }
    s.repo.Store(msg)

    found, err := s.repo.GetByUID(1, msg.UID)
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), msg.MessageID, found.MessageID)
}

func (s *MessageRepositoryTestSuite) TestUpdateFlags() {
    msg := &domain.Message{
        UserID:    1,
        MailboxID: 1,
        MessageID: "<flags123@example.com>",
        Flags:     []string{},
    }
    s.repo.Store(msg)

    // Add flags
    err := s.repo.AddFlags(1, []uint32{msg.UID}, []string{"\\Seen", "\\Flagged"})
    assert.NoError(s.T(), err)

    found, _ := s.repo.GetByUID(1, msg.UID)
    assert.Contains(s.T(), found.Flags, "\\Seen")
    assert.Contains(s.T(), found.Flags, "\\Flagged")

    // Remove flag
    err = s.repo.RemoveFlags(1, []uint32{msg.UID}, []string{"\\Flagged"})
    assert.NoError(s.T(), err)

    found, _ = s.repo.GetByUID(1, msg.UID)
    assert.Contains(s.T(), found.Flags, "\\Seen")
    assert.NotContains(s.T(), found.Flags, "\\Flagged")
}

func (s *MessageRepositoryTestSuite) TestSearch() {
    // Create test messages
    messages := []*domain.Message{
        {UserID: 1, MailboxID: 1, Subject: "Important meeting", From: "boss@company.com"},
        {UserID: 1, MailboxID: 1, Subject: "Weekly report", From: "team@company.com"},
        {UserID: 1, MailboxID: 1, Subject: "Newsletter", From: "news@example.com"},
    }
    for _, msg := range messages {
        s.repo.Store(msg)
    }

    // Search by subject
    results, err := s.repo.Search(1, domain.SearchCriteria{Subject: "meeting"})
    assert.NoError(s.T(), err)
    assert.Len(s.T(), results, 1)

    // Search by from
    results, err = s.repo.Search(1, domain.SearchCriteria{From: "company.com"})
    assert.NoError(s.T(), err)
    assert.Len(s.T(), results, 2)
}

func TestMessageRepositorySuite(t *testing.T) {
    suite.Run(t, new(MessageRepositoryTestSuite))
}
```

**Acceptance Criteria**:
- [ ] All repository methods have tests
- [ ] In-memory SQLite for isolation
- [ ] Edge cases covered (not found, duplicates)
- [ ] Pagination tested
- [ ] Constraints verified

---

### UT-002: Service Layer Tests
**File**: `internal/service/user_service_test.go`
```go
package service

import (
    "testing"

    "github.com/btafoya/gomailserver/internal/domain"
    "github.com/btafoya/gomailserver/internal/repository/mock"
    "github.com/stretchr/testify/assert"
    "go.uber.org/mock/gomock"
)

func TestUserService_Create(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mock.NewMockUserRepository(ctrl)
    mockDomainRepo := mock.NewMockDomainRepository(ctrl)
    service := NewUserService(mockRepo, mockDomainRepo)

    // Test successful creation
    t.Run("success", func(t *testing.T) {
        mockDomainRepo.EXPECT().
            GetByID(int64(1)).
            Return(&domain.Domain{ID: 1, Name: "example.com"}, nil)

        mockRepo.EXPECT().
            GetByEmail("new@example.com").
            Return(nil, domain.ErrNotFound)

        mockRepo.EXPECT().
            Create(gomock.Any()).
            DoAndReturn(func(u *domain.User) error {
                u.ID = 1
                return nil
            })

        user, err := service.Create("new@example.com", "password123", 1)
        assert.NoError(t, err)
        assert.Equal(t, int64(1), user.ID)
        assert.NotEmpty(t, user.PasswordHash)
        assert.NotEqual(t, "password123", user.PasswordHash) // Password should be hashed
    })

    // Test duplicate email
    t.Run("duplicate email", func(t *testing.T) {
        mockRepo.EXPECT().
            GetByEmail("existing@example.com").
            Return(&domain.User{ID: 1}, nil)

        _, err := service.Create("existing@example.com", "password", 1)
        assert.Equal(t, domain.ErrDuplicateEmail, err)
    })

    // Test weak password
    t.Run("weak password", func(t *testing.T) {
        _, err := service.Create("test@example.com", "123", 1) // Too short
        assert.Equal(t, domain.ErrWeakPassword, err)
    })
}

func TestUserService_Authenticate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mock.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo, nil)

    // Create user with known password
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
    testUser := &domain.User{
        ID:           1,
        Email:        "auth@example.com",
        PasswordHash: string(hashedPassword),
        Active:       true,
    }

    t.Run("successful auth", func(t *testing.T) {
        mockRepo.EXPECT().
            GetByEmail("auth@example.com").
            Return(testUser, nil)

        user, err := service.Authenticate("auth@example.com", "correct_password")
        assert.NoError(t, err)
        assert.Equal(t, testUser.ID, user.ID)
    })

    t.Run("wrong password", func(t *testing.T) {
        mockRepo.EXPECT().
            GetByEmail("auth@example.com").
            Return(testUser, nil)

        _, err := service.Authenticate("auth@example.com", "wrong_password")
        assert.Equal(t, domain.ErrInvalidCredentials, err)
    })

    t.Run("inactive user", func(t *testing.T) {
        inactiveUser := &domain.User{
            ID:           2,
            Email:        "inactive@example.com",
            PasswordHash: string(hashedPassword),
            Active:       false,
        }
        mockRepo.EXPECT().
            GetByEmail("inactive@example.com").
            Return(inactiveUser, nil)

        _, err := service.Authenticate("inactive@example.com", "correct_password")
        assert.Equal(t, domain.ErrUserInactive, err)
    })
}

func TestUserService_ChangePassword(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mock.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo, nil)

    oldHash, _ := bcrypt.GenerateFromPassword([]byte("old_password"), bcrypt.DefaultCost)
    user := &domain.User{
        ID:           1,
        PasswordHash: string(oldHash),
    }

    t.Run("success", func(t *testing.T) {
        mockRepo.EXPECT().
            GetByID(int64(1)).
            Return(user, nil)

        mockRepo.EXPECT().
            Update(gomock.Any()).
            DoAndReturn(func(u *domain.User) error {
                // Verify password was changed
                assert.NotEqual(t, string(oldHash), u.PasswordHash)
                return nil
            })

        err := service.ChangePassword(1, "old_password", "new_secure_password")
        assert.NoError(t, err)
    })

    t.Run("wrong current password", func(t *testing.T) {
        mockRepo.EXPECT().
            GetByID(int64(1)).
            Return(user, nil)

        err := service.ChangePassword(1, "wrong_password", "new_password")
        assert.Equal(t, domain.ErrInvalidCredentials, err)
    })
}
```

**Acceptance Criteria**:
- [ ] Business logic tested with mocks
- [ ] Authentication flows tested
- [ ] Password hashing verified
- [ ] Error cases covered
- [ ] Input validation tested

---

### UT-003: Security Function Tests
**File**: `internal/security/dkim_test.go`
```go
package security

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDKIM_GenerateKeyPair(t *testing.T) {
    t.Run("RSA 2048", func(t *testing.T) {
        priv, pub, err := GenerateDKIMKeyPair(RSA2048)
        assert.NoError(t, err)
        assert.NotEmpty(t, priv)
        assert.NotEmpty(t, pub)

        // Verify it's valid PEM
        block, _ := pem.Decode([]byte(priv))
        assert.NotNil(t, block)
        assert.Equal(t, "RSA PRIVATE KEY", block.Type)

        // Verify key size
        key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
        assert.NoError(t, err)
        assert.Equal(t, 2048, key.Size()*8)
    })

    t.Run("Ed25519", func(t *testing.T) {
        priv, pub, err := GenerateDKIMKeyPair(Ed25519)
        assert.NoError(t, err)
        assert.NotEmpty(t, priv)
        assert.NotEmpty(t, pub)
    })
}

func TestDKIM_Sign(t *testing.T) {
    // Generate key for testing
    privateKey, _, _ := GenerateDKIMKeyPair(RSA2048)

    signer := NewDKIMSigner(privateKey, "example.com", "default")

    message := []byte(`From: sender@example.com
To: recipient@test.com
Subject: Test Message
Date: Mon, 1 Jan 2024 12:00:00 +0000
Message-ID: <test123@example.com>

This is a test message body.
`)

    signed, err := signer.Sign(message)
    assert.NoError(t, err)
    assert.Contains(t, string(signed), "DKIM-Signature:")
    assert.Contains(t, string(signed), "d=example.com")
    assert.Contains(t, string(signed), "s=default")
}

func TestDKIM_Verify(t *testing.T) {
    // Create signer
    privateKey, publicKey, _ := GenerateDKIMKeyPair(RSA2048)
    signer := NewDKIMSigner(privateKey, "example.com", "default")

    // Create verifier with mock DNS
    verifier := NewDKIMVerifier(&mockDNSResolver{
        records: map[string]string{
            "default._domainkey.example.com": publicKey,
        },
    })

    message := []byte(`From: sender@example.com
To: recipient@test.com
Subject: Test
Message-ID: <verify@example.com>

Body content.
`)

    signed, _ := signer.Sign(message)
    result, err := verifier.Verify(signed)

    assert.NoError(t, err)
    assert.True(t, result.Valid)
    assert.Equal(t, "example.com", result.Domain)
}
```

**File**: `internal/security/spf_test.go`
```go
package security

import (
    "net"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestSPF_Check(t *testing.T) {
    checker := NewSPFChecker(&mockDNSResolver{
        txtRecords: map[string][]string{
            "example.com": {"v=spf1 ip4:192.168.1.0/24 include:_spf.google.com -all"},
        },
    })

    tests := []struct {
        name     string
        ip       string
        domain   string
        expected SPFResult
    }{
        {
            name:     "Pass - IP in range",
            ip:       "192.168.1.100",
            domain:   "example.com",
            expected: SPFPass,
        },
        {
            name:     "Fail - IP not in range",
            ip:       "10.0.0.1",
            domain:   "example.com",
            expected: SPFFail,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := checker.Check(net.ParseIP(tt.ip), tt.domain, "sender@"+tt.domain)
            assert.Equal(t, tt.expected, result.Result)
        })
    }
}

func TestSPF_Parse(t *testing.T) {
    tests := []struct {
        record    string
        shouldErr bool
    }{
        {"v=spf1 ip4:192.168.1.0/24 -all", false},
        {"v=spf1 a mx include:example.com ~all", false},
        {"v=spf1 redirect=_spf.example.com", false},
        {"invalid record", true},
        {"v=spf2 something", true}, // Wrong version
    }

    for _, tt := range tests {
        t.Run(tt.record, func(t *testing.T) {
            _, err := ParseSPFRecord(tt.record)
            if tt.shouldErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**File**: `internal/security/password_test.go`
```go
package security

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestPasswordStrength(t *testing.T) {
    tests := []struct {
        password string
        minScore int
        valid    bool
    }{
        {"", 0, false},
        {"123", 0, false},
        {"password", 1, false},
        {"Password1", 2, true},
        {"P@ssw0rd!", 3, true},
        {"MyS3cur3P@ssw0rd!2024", 4, true},
    }

    for _, tt := range tests {
        t.Run(tt.password, func(t *testing.T) {
            score, feedback := CheckPasswordStrength(tt.password)
            if tt.valid {
                assert.GreaterOrEqual(t, score, tt.minScore)
            } else {
                assert.NotEmpty(t, feedback)
            }
        })
    }
}

func TestPasswordHashing(t *testing.T) {
    password := "test_password_123"

    hash, err := HashPassword(password)
    assert.NoError(t, err)
    assert.NotEqual(t, password, hash)

    // Verify correct password
    assert.True(t, VerifyPassword(password, hash))

    // Verify wrong password
    assert.False(t, VerifyPassword("wrong_password", hash))
}
```

**Acceptance Criteria**:
- [ ] DKIM signing/verification tested
- [ ] SPF parsing and checking tested
- [ ] DMARC validation tested
- [ ] Password strength validation tested
- [ ] Rate limiting tested

---

### UT-004: Code Coverage Target
**File**: `.github/workflows/test.yml`
```yaml
name: Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y gcc libsqlite3-dev

      - name: Run tests with coverage
        run: |
          go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $COVERAGE%"
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "Coverage is below 80%!"
            exit 1
          fi

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: ./coverage.out
          fail_ci_if_error: true
```

**Coverage Exclusions** (`.coverignore`):
```
# Generated code
internal/api/openapi.gen.go
internal/repository/mock/*.go

# Main entry points
cmd/gomailserver/main.go

# Vendor
vendor/
```

**Acceptance Criteria**:
- [ ] 80%+ overall coverage
- [ ] Critical paths 90%+ coverage
- [ ] Coverage tracked in CI
- [ ] Coverage reports generated

---

## 10.2 Integration Tests [MVP]

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| IT-001 | SMTP send/receive tests | [ ] | Phase 1 SMTP | MVP |
| IT-002 | IMAP tests | [ ] | Phase 1 IMAP | MVP |
| IT-003 | Authentication tests | [ ] | Phase 2 Auth | MVP |
| IT-004 | API tests | [ ] | Phase 3 API | MVP |

---

### IT-001: SMTP Integration Tests
**File**: `tests/integration/smtp_test.go`
```go
package integration

import (
    "net/smtp"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSMTPSubmission(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Start test server
    server := startTestServer(t)
    defer server.Stop()

    t.Run("authenticated submission", func(t *testing.T) {
        auth := smtp.PlainAuth("", "testuser@example.com", "testpass", "localhost")

        err := smtp.SendMail(
            "localhost:10587",
            auth,
            "testuser@example.com",
            []string{"recipient@example.com"},
            []byte("Subject: Test\r\n\r\nTest body"),
        )
        assert.NoError(t, err)

        // Verify message was queued
        time.Sleep(100 * time.Millisecond)
        queue := server.GetQueuedMessages()
        assert.Len(t, queue, 1)
    })

    t.Run("unauthenticated submission rejected", func(t *testing.T) {
        err := smtp.SendMail(
            "localhost:10587",
            nil,
            "anonymous@external.com",
            []string{"recipient@example.com"},
            []byte("Subject: Spam\r\n\r\nSpam body"),
        )
        assert.Error(t, err)
    })

    t.Run("relay - valid domain", func(t *testing.T) {
        // Connect to relay port
        conn, err := smtp.Dial("localhost:10025")
        require.NoError(t, err)
        defer conn.Close()

        err = conn.Mail("sender@external.com")
        assert.NoError(t, err)

        err = conn.Rcpt("localuser@example.com")
        assert.NoError(t, err)

        wc, err := conn.Data()
        require.NoError(t, err)

        _, err = wc.Write([]byte("Subject: Incoming\r\n\r\nBody"))
        assert.NoError(t, err)
        wc.Close()

        // Verify delivered
        time.Sleep(100 * time.Millisecond)
        messages := server.GetUserMessages("localuser@example.com")
        assert.Len(t, messages, 1)
    })

    t.Run("relay - unknown domain rejected", func(t *testing.T) {
        conn, _ := smtp.Dial("localhost:10025")
        defer conn.Close()

        conn.Mail("sender@external.com")
        err := conn.Rcpt("user@unknown.com")
        assert.Error(t, err) // Should be rejected
    })
}

func TestSMTPTLS(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    server := startTestServer(t)
    defer server.Stop()

    t.Run("STARTTLS", func(t *testing.T) {
        conn, err := smtp.Dial("localhost:10587")
        require.NoError(t, err)
        defer conn.Close()

        // Should support STARTTLS
        err = conn.StartTLS(&tls.Config{
            InsecureSkipVerify: true,
        })
        assert.NoError(t, err)
    })

    t.Run("SMTPS", func(t *testing.T) {
        conn, err := tls.Dial("tcp", "localhost:10465", &tls.Config{
            InsecureSkipVerify: true,
        })
        require.NoError(t, err)
        conn.Close()
    })
}

func TestSMTPSizeLimit(t *testing.T) {
    server := startTestServer(t)
    defer server.Stop()

    // Server has 10MB limit
    auth := smtp.PlainAuth("", "testuser@example.com", "testpass", "localhost")

    // Create message larger than limit
    largeBody := make([]byte, 11*1024*1024) // 11MB
    err := smtp.SendMail(
        "localhost:10587",
        auth,
        "testuser@example.com",
        []string{"recipient@example.com"},
        append([]byte("Subject: Large\r\n\r\n"), largeBody...),
    )
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "552") // Message too big
}
```

---

### IT-002: IMAP Integration Tests
**File**: `tests/integration/imap_test.go`
```go
package integration

import (
    "testing"

    "github.com/emersion/go-imap/client"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestIMAPBasicOperations(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    server := startTestServer(t)
    defer server.Stop()

    // Seed test data
    seedTestUser(server, "imapuser@example.com", "testpass")
    seedTestMessages(server, "imapuser@example.com", 10)

    t.Run("login", func(t *testing.T) {
        c, err := client.Dial("localhost:10143")
        require.NoError(t, err)
        defer c.Logout()

        err = c.Login("imapuser@example.com", "testpass")
        assert.NoError(t, err)
    })

    t.Run("list mailboxes", func(t *testing.T) {
        c := connectAndLogin(t, "imapuser@example.com", "testpass")
        defer c.Logout()

        mailboxes := make(chan *imap.MailboxInfo, 10)
        done := make(chan error, 1)

        go func() {
            done <- c.List("", "*", mailboxes)
        }()

        var names []string
        for m := range mailboxes {
            names = append(names, m.Name)
        }

        assert.NoError(t, <-done)
        assert.Contains(t, names, "INBOX")
        assert.Contains(t, names, "Sent")
        assert.Contains(t, names, "Drafts")
        assert.Contains(t, names, "Trash")
    })

    t.Run("select and fetch", func(t *testing.T) {
        c := connectAndLogin(t, "imapuser@example.com", "testpass")
        defer c.Logout()

        // Select INBOX
        mbox, err := c.Select("INBOX", false)
        require.NoError(t, err)
        assert.Equal(t, uint32(10), mbox.Messages)

        // Fetch first message
        seqset := new(imap.SeqSet)
        seqset.AddNum(1)

        messages := make(chan *imap.Message, 1)
        done := make(chan error, 1)

        go func() {
            done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
        }()

        msg := <-messages
        assert.NotNil(t, msg)
        assert.NotEmpty(t, msg.Envelope.Subject)

        assert.NoError(t, <-done)
    })

    t.Run("IDLE", func(t *testing.T) {
        c := connectAndLogin(t, "imapuser@example.com", "testpass")
        defer c.Logout()

        c.Select("INBOX", false)

        // Start IDLE
        idleClient := idle.NewClient(c)
        updates := make(chan client.Update, 1)
        c.Updates = updates

        done := make(chan error, 1)
        stop := make(chan struct{})

        go func() {
            done <- idleClient.IdleWithFallback(stop, 0)
        }()

        // Deliver new message
        go func() {
            time.Sleep(100 * time.Millisecond)
            deliverMessage(server, "imapuser@example.com")
        }()

        // Should receive update
        select {
        case update := <-updates:
            _, ok := update.(*client.MailboxUpdate)
            assert.True(t, ok)
        case <-time.After(5 * time.Second):
            t.Fatal("Timeout waiting for IDLE update")
        }

        close(stop)
        <-done
    })
}

func TestIMAPSearch(t *testing.T) {
    server := startTestServer(t)
    defer server.Stop()

    seedTestUser(server, "search@example.com", "testpass")
    seedSearchTestMessages(server, "search@example.com")

    c := connectAndLogin(t, "search@example.com", "testpass")
    defer c.Logout()
    c.Select("INBOX", false)

    t.Run("search by subject", func(t *testing.T) {
        criteria := imap.NewSearchCriteria()
        criteria.Header.Set("Subject", "Important")

        uids, err := c.Search(criteria)
        assert.NoError(t, err)
        assert.NotEmpty(t, uids)
    })

    t.Run("search by date", func(t *testing.T) {
        criteria := imap.NewSearchCriteria()
        criteria.Since = time.Now().AddDate(0, 0, -7)

        uids, err := c.Search(criteria)
        assert.NoError(t, err)
        assert.NotEmpty(t, uids)
    })

    t.Run("search unseen", func(t *testing.T) {
        criteria := imap.NewSearchCriteria()
        criteria.WithoutFlags = []string{imap.SeenFlag}

        uids, err := c.Search(criteria)
        assert.NoError(t, err)
        assert.NotEmpty(t, uids)
    })
}
```

---

### IT-003: Authentication Integration Tests
**File**: `tests/integration/auth_test.go`
```go
package integration

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestAuthenticationMechanisms(t *testing.T) {
    server := startTestServer(t)
    defer server.Stop()

    seedTestUser(server, "authtest@example.com", "correctpass")

    t.Run("PLAIN auth", func(t *testing.T) {
        c, _ := smtp.Dial("localhost:10587")
        defer c.Close()

        auth := smtp.PlainAuth("", "authtest@example.com", "correctpass", "localhost")
        err := c.Auth(auth)
        assert.NoError(t, err)
    })

    t.Run("wrong password", func(t *testing.T) {
        c, _ := smtp.Dial("localhost:10587")
        defer c.Close()

        auth := smtp.PlainAuth("", "authtest@example.com", "wrongpass", "localhost")
        err := c.Auth(auth)
        assert.Error(t, err)
    })

    t.Run("brute force protection", func(t *testing.T) {
        // Attempt multiple failed logins
        for i := 0; i < 5; i++ {
            c, _ := smtp.Dial("localhost:10587")
            auth := smtp.PlainAuth("", "authtest@example.com", "wrongpass", "localhost")
            c.Auth(auth)
            c.Close()
        }

        // Should be rate limited
        c, _ := smtp.Dial("localhost:10587")
        defer c.Close()
        auth := smtp.PlainAuth("", "authtest@example.com", "correctpass", "localhost")
        err := c.Auth(auth)
        assert.Error(t, err) // Should be blocked
    })
}

func TestTOTP2FA(t *testing.T) {
    server := startTestServer(t)
    defer server.Stop()

    // Create user with 2FA enabled
    user, secret := seedUserWith2FA(server, "2fa@example.com", "password")

    t.Run("valid TOTP", func(t *testing.T) {
        totp := generateTOTP(secret, time.Now())
        token, err := server.AuthenticateWith2FA("2fa@example.com", "password", totp)
        assert.NoError(t, err)
        assert.NotEmpty(t, token)
    })

    t.Run("invalid TOTP", func(t *testing.T) {
        _, err := server.AuthenticateWith2FA("2fa@example.com", "password", "000000")
        assert.Error(t, err)
    })

    t.Run("expired TOTP", func(t *testing.T) {
        // Generate TOTP from 2 minutes ago
        totp := generateTOTP(secret, time.Now().Add(-2*time.Minute))
        _, err := server.AuthenticateWith2FA("2fa@example.com", "password", totp)
        assert.Error(t, err)
    })
}
```

---

### IT-004: API Integration Tests
**File**: `tests/integration/api_test.go`
```go
package integration

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAPIAuth(t *testing.T) {
    server := startTestAPIServer(t)

    t.Run("login returns token", func(t *testing.T) {
        resp := server.POST("/api/auth/login", `{
            "email": "admin@example.com",
            "password": "adminpass"
        }`)

        assert.Equal(t, http.StatusOK, resp.Code)

        var result map[string]interface{}
        json.Unmarshal(resp.Body.Bytes(), &result)

        assert.NotEmpty(t, result["token"])
        assert.NotEmpty(t, result["expires_at"])
    })

    t.Run("protected route requires token", func(t *testing.T) {
        resp := server.GET("/api/admin/users")
        assert.Equal(t, http.StatusUnauthorized, resp.Code)
    })

    t.Run("protected route with valid token", func(t *testing.T) {
        token := getAuthToken(server, "admin@example.com", "adminpass")

        resp := server.GETWithAuth("/api/admin/users", token)
        assert.Equal(t, http.StatusOK, resp.Code)
    })
}

func TestAPIDomainCRUD(t *testing.T) {
    server := startTestAPIServer(t)
    token := getAdminToken(server)

    var domainID int64

    t.Run("create domain", func(t *testing.T) {
        resp := server.POSTWithAuth("/api/admin/domains", `{
            "name": "newdomain.com",
            "active": true
        }`, token)

        assert.Equal(t, http.StatusCreated, resp.Code)

        var result map[string]interface{}
        json.Unmarshal(resp.Body.Bytes(), &result)
        domainID = int64(result["id"].(float64))
        assert.NotZero(t, domainID)
    })

    t.Run("list domains", func(t *testing.T) {
        resp := server.GETWithAuth("/api/admin/domains", token)
        assert.Equal(t, http.StatusOK, resp.Code)

        var domains []map[string]interface{}
        json.Unmarshal(resp.Body.Bytes(), &domains)
        assert.NotEmpty(t, domains)
    })

    t.Run("get domain", func(t *testing.T) {
        resp := server.GETWithAuth(fmt.Sprintf("/api/admin/domains/%d", domainID), token)
        assert.Equal(t, http.StatusOK, resp.Code)

        var domain map[string]interface{}
        json.Unmarshal(resp.Body.Bytes(), &domain)
        assert.Equal(t, "newdomain.com", domain["name"])
    })

    t.Run("update domain", func(t *testing.T) {
        resp := server.PUTWithAuth(fmt.Sprintf("/api/admin/domains/%d", domainID), `{
            "active": false
        }`, token)

        assert.Equal(t, http.StatusOK, resp.Code)
    })

    t.Run("delete domain", func(t *testing.T) {
        resp := server.DELETEWithAuth(fmt.Sprintf("/api/admin/domains/%d", domainID), token)
        assert.Equal(t, http.StatusNoContent, resp.Code)

        // Verify deleted
        resp = server.GETWithAuth(fmt.Sprintf("/api/admin/domains/%d", domainID), token)
        assert.Equal(t, http.StatusNotFound, resp.Code)
    })
}

func TestAPIRateLimiting(t *testing.T) {
    server := startTestAPIServer(t)

    // Make many requests quickly
    for i := 0; i < 100; i++ {
        server.POST("/api/auth/login", `{"email":"test@test.com","password":"wrong"}`)
    }

    // Should be rate limited
    resp := server.POST("/api/auth/login", `{"email":"test@test.com","password":"wrong"}`)
    assert.Equal(t, http.StatusTooManyRequests, resp.Code)
}
```

**Acceptance Criteria**:
- [ ] SMTP submission tested
- [ ] SMTP relay tested
- [ ] IMAP operations tested
- [ ] Authentication flows tested
- [ ] API endpoints tested
- [ ] Rate limiting tested

---

## 10.3 External Testing [FULL]

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| ET-001 | mail-tester.com score 10/10 | [ ] | Phase 2 | FULL |
| ET-002 | Thunderbird compatibility | [ ] | I-002, CC-001 | FULL |
| ET-003 | Apple Mail compatibility | [ ] | I-002, CC-002 | FULL |
| ET-004 | Mobile client compatibility | [ ] | I-002, CC-003-004 | FULL |
| ET-005 | Outlook compatibility | [ ] | I-002, CC-005 | FULL |

---

### ET-001: Mail Tester Score
**File**: `tests/external/mail_tester.md`
```markdown
# Mail Tester Validation

## Test Procedure

1. Start gomailserver with production configuration
2. Send test email to mail-tester.com address
3. Check results at provided URL

## Required Score: 10/10

## Checklist

### Authentication (3 points)
- [ ] SPF record valid and passes
- [ ] DKIM signature valid
- [ ] DMARC policy published and passes

### Blacklist Check (2 points)
- [ ] IP not on any blacklists
- [ ] Domain not on any blacklists

### Message Content (3 points)
- [ ] No SpamAssassin triggers
- [ ] Valid HTML formatting
- [ ] Proper MIME structure

### Technical (2 points)
- [ ] Reverse DNS matches
- [ ] Valid HELO/EHLO
- [ ] Proper TLS configuration

## Automated Testing Script
```bash
#!/bin/bash
# Send test email to mail-tester.com
MAIL_TESTER_ADDRESS="test-xxxxxxxx@mail-tester.com"

gomailserver sendtest \
    --to "$MAIL_TESTER_ADDRESS" \
    --subject "Mail Server Test" \
    --body "Testing mail configuration"

echo "Check results at: https://www.mail-tester.com/test-xxxxxxxx"
```
```

---

## 10.4 Performance Tests [FULL]

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| PT-001 | Load testing (100,000 emails/day) | [ ] | Phase 1 | FULL |
| PT-002 | Concurrent connection testing | [ ] | Phase 1 | FULL |
| PT-003 | Memory usage benchmarks | [ ] | All | FULL |
| PT-004 | IMAP response time benchmarks | [ ] | I-002 | FULL |

---

### PT-001: Load Testing
**File**: `tests/performance/load_test.go`
```go
package performance

import (
    "sync"
    "sync/atomic"
    "testing"
    "time"
)

func TestSMTPThroughput(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping performance test")
    }

    server := startTestServer(t)
    defer server.Stop()

    const (
        numWorkers    = 10
        messagesPerWorker = 1000
        targetThroughput = 100000 / 24 / 60 // ~70 messages/minute
    )

    var (
        sent     int64
        failed   int64
        wg       sync.WaitGroup
    )

    start := time.Now()

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()

            for j := 0; j < messagesPerWorker; j++ {
                err := sendTestMessage(server, workerID, j)
                if err != nil {
                    atomic.AddInt64(&failed, 1)
                } else {
                    atomic.AddInt64(&sent, 1)
                }
            }
        }(i)
    }

    wg.Wait()
    duration := time.Since(start)

    totalMessages := numWorkers * messagesPerWorker
    throughput := float64(sent) / duration.Seconds() * 60 // per minute

    t.Logf("Performance Results:")
    t.Logf("  Total messages: %d", totalMessages)
    t.Logf("  Sent: %d", sent)
    t.Logf("  Failed: %d", failed)
    t.Logf("  Duration: %v", duration)
    t.Logf("  Throughput: %.2f msg/min", throughput)

    if failed > int64(totalMessages/100) {
        t.Errorf("Too many failures: %d/%d", failed, totalMessages)
    }

    if throughput < float64(targetThroughput) {
        t.Errorf("Throughput below target: %.2f < %d msg/min", throughput, targetThroughput)
    }
}

func BenchmarkMessageDelivery(b *testing.B) {
    server := startBenchServer(b)
    defer server.Stop()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sendBenchMessage(server, i)
    }
}
```

---

### PT-002: Concurrent Connection Testing
```go
func TestConcurrentConnections(t *testing.T) {
    server := startTestServer(t)
    defer server.Stop()

    const maxConnections = 500

    var (
        connected int64
        errors    int64
        wg        sync.WaitGroup
    )

    for i := 0; i < maxConnections; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            conn, err := net.DialTimeout("tcp", "localhost:10143", 5*time.Second)
            if err != nil {
                atomic.AddInt64(&errors, 1)
                return
            }
            atomic.AddInt64(&connected, 1)

            // Hold connection
            time.Sleep(5 * time.Second)
            conn.Close()
        }()

        // Stagger connections
        time.Sleep(10 * time.Millisecond)
    }

    wg.Wait()

    t.Logf("Connected: %d/%d", connected, maxConnections)
    t.Logf("Errors: %d", errors)

    if connected < int64(maxConnections*0.95) {
        t.Errorf("Too few connections succeeded: %d", connected)
    }
}
```

---

## 10.5 Security Audit [FULL]

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| SA-001 | Input validation review | [ ] | All | FULL |
| SA-002 | Authentication security review | [ ] | Phase 2-3 Auth | FULL |
| SA-003 | TLS configuration review | [ ] | T-001-004 | FULL |
| SA-004 | SQL injection testing | [ ] | F-023 | FULL |

---

### SA-001: Input Validation Tests
**File**: `tests/security/input_validation_test.go`
```go
package security

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestEmailAddressValidation(t *testing.T) {
    tests := []struct {
        email string
        valid bool
    }{
        {"user@example.com", true},
        {"user+tag@example.com", true},
        {"user.name@sub.example.com", true},
        {"", false},
        {"invalid", false},
        {"@example.com", false},
        {"user@", false},
        {"user@.com", false},
        {"user@example..com", false},
        {"user@example.com\nBcc: attacker@evil.com", false}, // Header injection
        {"user\x00@example.com", false}, // Null byte
        {strings.Repeat("a", 300) + "@example.com", false}, // Too long
    }

    for _, tt := range tests {
        t.Run(tt.email, func(t *testing.T) {
            err := validateEmail(tt.email)
            if tt.valid {
                assert.NoError(t, err)
            } else {
                assert.Error(t, err)
            }
        })
    }
}

func TestHeaderInjection(t *testing.T) {
    tests := []struct {
        name  string
        input string
        safe  bool
    }{
        {"normal subject", "Hello World", true},
        {"CR injection", "Hello\rBcc: attacker@evil.com", false},
        {"LF injection", "Hello\nBcc: attacker@evil.com", false},
        {"CRLF injection", "Hello\r\nBcc: attacker@evil.com", false},
        {"null byte", "Hello\x00World", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sanitized := sanitizeHeader(tt.input)
            if tt.safe {
                assert.Equal(t, tt.input, sanitized)
            } else {
                assert.NotEqual(t, tt.input, sanitized)
                assert.NotContains(t, sanitized, "\r")
                assert.NotContains(t, sanitized, "\n")
                assert.NotContains(t, sanitized, "\x00")
            }
        })
    }
}

func TestPathTraversal(t *testing.T) {
    tests := []struct {
        path string
        safe bool
    }{
        {"inbox/message1", true},
        {"../etc/passwd", false},
        {"..\\windows\\system32", false},
        {"/absolute/path", false},
        {"valid/../attempt", false},
        {"normal.folder/file.eml", true},
    }

    for _, tt := range tests {
        t.Run(tt.path, func(t *testing.T) {
            safe := isPathSafe(tt.path)
            assert.Equal(t, tt.safe, safe)
        })
    }
}
```

---

### SA-004: SQL Injection Tests
**File**: `tests/security/sql_injection_test.go`
```go
package security

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestSQLInjectionPrevention(t *testing.T) {
    server := startTestServer(t)
    defer server.Stop()

    // Seed normal user
    seedTestUser(server, "normal@example.com", "password")

    injectionPayloads := []string{
        "' OR '1'='1",
        "'; DROP TABLE users; --",
        "' UNION SELECT * FROM users --",
        "admin'--",
        "1; SELECT * FROM users",
        "' OR 1=1#",
        "') OR ('1'='1",
        "admin'; UPDATE users SET password='hacked' WHERE '1'='1",
    }

    t.Run("login injection", func(t *testing.T) {
        for _, payload := range injectionPayloads {
            t.Run(payload, func(t *testing.T) {
                _, err := server.Authenticate(payload, "password")
                assert.Error(t, err)

                _, err = server.Authenticate("normal@example.com", payload)
                assert.Error(t, err)
            })
        }
    })

    t.Run("search injection", func(t *testing.T) {
        for _, payload := range injectionPayloads {
            t.Run(payload, func(t *testing.T) {
                // Should not cause error or return unexpected data
                results, err := server.SearchMessages("normal@example.com", payload)
                assert.NoError(t, err)
                // Results should be empty (nothing matches injection payload)
                assert.Empty(t, results)
            })
        }
    })

    t.Run("api parameter injection", func(t *testing.T) {
        token := getAuthToken(server, "normal@example.com", "password")

        for _, payload := range injectionPayloads {
            resp := server.GETWithAuth("/api/admin/users?search="+url.QueryEscape(payload), token)
            // Should return 200 or 403, not 500 (database error)
            assert.NotEqual(t, 500, resp.Code)
        }
    })
}
```

---

## Test Infrastructure

### Test Server Helper
**File**: `tests/testutil/server.go`
```go
package testutil

import (
    "context"
    "testing"

    "github.com/btafoya/gomailserver/internal/server"
)

type TestServer struct {
    *server.Server
    t      *testing.T
    ctx    context.Context
    cancel context.CancelFunc
}

func StartTestServer(t *testing.T) *TestServer {
    ctx, cancel := context.WithCancel(context.Background())

    cfg := &server.Config{
        Database: server.DatabaseConfig{
            Path: ":memory:",
        },
        SMTP: server.SMTPConfig{
            Submission: server.SubmissionConfig{
                Port: 10587,
            },
            Relay: server.RelayConfig{
                Port: 10025,
            },
        },
        IMAP: server.IMAPConfig{
            Port: 10143,
        },
    }

    srv, err := server.New(cfg)
    if err != nil {
        t.Fatalf("Failed to create test server: %v", err)
    }

    go srv.Run(ctx)

    // Wait for server to be ready
    waitForServer(t, "localhost:10587")
    waitForServer(t, "localhost:10143")

    return &TestServer{
        Server: srv,
        t:      t,
        ctx:    ctx,
        cancel: cancel,
    }
}

func (ts *TestServer) Stop() {
    ts.cancel()
}
```

---

## Coverage Requirements

| Package | Min Coverage |
|---------|--------------|
| `internal/repository` | 85% |
| `internal/service` | 80% |
| `internal/security` | 90% |
| `internal/smtp` | 75% |
| `internal/imap` | 75% |
| `internal/api` | 80% |
| **Overall** | **80%** |

---

## CI/CD Test Configuration

**File**: `.github/workflows/test.yml`
```yaml
name: Tests

on:
  push:
    branches: [main]
  pull_request:

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test -v -race -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v4

  integration-tests:
    runs-on: ubuntu-latest
    services:
      clamav:
        image: clamav/clamav:latest
      spamassassin:
        image: spamassassin/spamassassin:latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test -v -tags=integration ./tests/integration/...

  performance-tests:
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test -v -tags=performance -bench=. ./tests/performance/...

  security-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test -v -tags=security ./tests/security/...
      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: ./...
```

---

## 10.6 Chaos Engineering Tests [OPTIONAL]

Chaos engineering tests to validate system resilience under failure conditions.

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| CHAOS-001 | Chaos testing suite | [ ] | All tests | OPTIONAL |

---

### CHAOS-001: Chaos Testing Suite

**File**: `tests/chaos/process_kill_test.go`
```go
package chaos

import (
    "context"
    "os"
    "os/exec"
    "syscall"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestProcessKillRecovery(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping chaos test in short mode")
    }

    // Start gomailserver process
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    cmd := exec.CommandContext(ctx, "./build/gomailserver", "run", "--config", "testdata/test.conf")
    require.NoError(t, cmd.Start())

    // Wait for startup
    time.Sleep(5 * time.Second)

    // Verify process is running
    assert.True(t, isProcessRunning(cmd.Process.Pid))

    // SIGKILL (abrupt termination)
    t.Log("Sending SIGKILL to process...")
    cmd.Process.Signal(syscall.SIGKILL)
    cmd.Wait()

    // Restart process
    cmd = exec.CommandContext(ctx, "./build/gomailserver", "run", "--config", "testdata/test.conf")
    require.NoError(t, cmd.Start())
    defer cmd.Process.Kill()

    // Wait for recovery
    time.Sleep(5 * time.Second)

    // Verify process recovered successfully
    assert.True(t, isProcessRunning(cmd.Process.Pid))

    // Verify database integrity
    // Verify no data corruption
    // Verify services operational
}

func TestGracefulShutdown(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping chaos test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    cmd := exec.CommandContext(ctx, "./build/gomailserver", "run", "--config", "testdata/test.conf")
    require.NoError(t, cmd.Start())

    time.Sleep(5 * time.Second)

    // SIGTERM (graceful shutdown)
    t.Log("Sending SIGTERM to process...")
    cmd.Process.Signal(syscall.SIGTERM)

    // Wait for graceful shutdown (max 30s)
    done := make(chan error)
    go func() {
        done <- cmd.Wait()
    }()

    select {
    case <-time.After(30 * time.Second):
        t.Fatal("Process did not shut down gracefully within 30s")
    case err := <-done:
        assert.NoError(t, err)
    }

    // Verify clean shutdown
    // - No zombie processes
    // - Database connections closed
    // - SMTP/IMAP listeners stopped
    // - Queue flushed or persisted
}

func isProcessRunning(pid int) bool {
    process, err := os.FindProcess(pid)
    if err != nil {
        return false
    }
    err = process.Signal(syscall.Signal(0))
    return err == nil
}
```

**File**: `tests/chaos/network_partition_test.go`
```go
package chaos

import (
    "net"
    "os/exec"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestNetworkPartition(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping chaos test in short mode")
    }

    // Simulate network partition using iptables or network namespace
    t.Log("Simulating network partition...")

    // Block outbound SMTP (port 25)
    cmd := exec.Command("iptables", "-A", "OUTPUT", "-p", "tcp", "--dport", "25", "-j", "DROP")
    cmd.Run()
    defer exec.Command("iptables", "-D", "OUTPUT", "-p", "tcp", "--dport", "25", "-j", "DROP").Run()

    // Attempt to send email (should queue for retry)
    // Verify message is queued, not lost
    // Verify exponential backoff starts

    // Restore network
    exec.Command("iptables", "-D", "OUTPUT", "-p", "tcp", "--dport", "25", "-j", "DROP").Run()

    // Wait for queue retry
    time.Sleep(35 * time.Second) // First retry at 30s

    // Verify message delivered after network recovery
}

func TestDatabaseConnectionLoss(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping chaos test in short mode")
    }

    // Simulate database connection loss
    t.Log("Simulating database connection loss...")

    // Use iptables to block SQLite file access (if remote)
    // Or use file permissions to simulate I/O error

    // Verify system degrades gracefully:
    // - Read operations continue from cache
    // - Write operations are queued
    // - Health check reports unhealthy
    // - System doesn't crash

    // Restore database access
    // Verify system recovers:
    // - Queued writes are flushed
    // - Health check reports healthy
    // - Normal operation resumes
}
```

**File**: `tests/chaos/disk_full_test.go`
```go
package chaos

import (
    "io"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDiskFullScenario(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping chaos test in short mode")
    }

    // Create small filesystem for testing
    // Use loopback device with limited space

    t.Log("Simulating disk full scenario...")

    // Fill disk to 100% (leave only 1MB free)
    // Attempt operations:
    // - Receive email (should reject with 4xx code)
    // - Queue message (should fail gracefully)
    // - Write logs (should handle gracefully)

    // Verify:
    // - No data corruption
    // - Clear error messages
    // - System remains operational (doesn't crash)
    // - Recovery when space available
}

func TestLowMemoryPressure(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping chaos test in short mode")
    }

    t.Log("Simulating memory pressure...")

    // Use cgroups to limit memory
    // Set memory limit to 512MB

    // Create high load:
    // - 100 concurrent SMTP connections
    // - 100 concurrent IMAP connections
    // - Large message processing

    // Verify:
    // - System stays within memory limits
    // - No OOM kills
    // - Graceful connection rejection when at capacity
    // - Memory cleanup (GC) works correctly
}
```

**File**: `tests/chaos/README.md`
```markdown
# Chaos Engineering Tests

## Purpose

Validate system resilience and graceful degradation under failure conditions.

## Test Scenarios

### Process Resilience
- **Kill Testing**: Abrupt process termination (SIGKILL) and recovery
- **Graceful Shutdown**: Clean shutdown on SIGTERM within 30s
- **Restart Testing**: Data integrity after restart

### Network Failures
- **Network Partition**: Outbound network blocked, queue and retry behavior
- **DNS Failure**: DNS resolution fails, fallback to IP
- **Connection Timeout**: Slow/hanging connections, timeout enforcement

### Resource Exhaustion
- **Disk Full**: 100% disk usage, reject new messages gracefully
- **Memory Pressure**: Limited memory (512MB), stay within bounds
- **File Descriptor Exhaustion**: Max file descriptors, connection limiting

### Database Failures
- **Connection Loss**: Database unavailable, graceful degradation
- **Corruption Detection**: Detect and report corrupted database
- **Lock Contention**: High concurrent writes, deadlock prevention

## Running Chaos Tests

```bash
# Run all chaos tests (requires root for iptables)
sudo go test -v ./tests/chaos/...

# Run specific test
sudo go test -v ./tests/chaos/ -run TestProcessKillRecovery

# Skip chaos tests (default for short mode)
go test -short ./...
```

## Requirements

- Root access (for iptables, cgroups)
- Docker (for container-based isolation)
- 2GB+ available memory
- 1GB+ available disk space

## Safety

**⚠️ WARNING**: These tests modify system state (iptables, cgroups).
- Run in isolated environment (VM, container)
- Do NOT run on production systems
- Tests automatically clean up, but may require manual cleanup on failure
```

**Acceptance Criteria**:
- [ ] Process kill/restart recovery tests
- [ ] Graceful shutdown within 30s (SIGTERM)
- [ ] Network partition simulation (queue and retry)
- [ ] Disk full scenario (reject new messages gracefully)
- [ ] Memory pressure testing (stay within 512MB limit)
- [ ] Database connection loss (graceful degradation)
- [ ] All chaos tests automated and reproducible

**Production Readiness**:
- [ ] Recovery time: < 5s after process restart
- [ ] Graceful shutdown: < 30s for clean exit
- [ ] Queue persistence: Zero message loss on abrupt termination
- [ ] Resource limits: Enforce connection limits before exhaustion
- [ ] Error handling: Clear error messages for resource exhaustion
- [ ] Health monitoring: Health checks accurately reflect degraded state

**Given/When/Then Scenarios**:
```
Given gomailserver is running with active connections
When process receives SIGKILL (abrupt termination)
Then process terminates immediately
When process restarts
Then database integrity is verified
And all queued messages are preserved
And services resume within 5 seconds

Given gomailserver is running
When process receives SIGTERM (graceful shutdown)
Then active connections complete (up to 30s grace period)
And queue is flushed to disk
And database connections are closed
And process exits cleanly within 30 seconds

Given gomailserver is queueing messages
When network partition blocks outbound SMTP
Then messages remain queued (not lost)
And exponential backoff retry begins (30s, 1m, 2m...)
When network partition is resolved
Then queued messages are delivered successfully

Given disk usage is at 95%
When receiving new message
Then message is accepted
When disk usage reaches 100%
Then new messages are rejected with 452 (insufficient storage)
And existing messages remain intact
And services continue operating (no crash)

Given memory limit is 512MB
When under high load (100 SMTP + 100 IMAP connections)
Then memory usage stays < 512MB
And new connections are rejected gracefully when at capacity
And no OOM kills occur
```

---

## Testing Checklist Summary

### MVP (Must Pass)
- [ ] 80%+ unit test coverage
- [ ] All repository tests pass
- [ ] All service tests pass
- [ ] Security function tests pass
- [ ] SMTP integration tests pass
- [ ] IMAP integration tests pass
- [ ] Authentication tests pass
- [ ] API tests pass

### FULL (Should Pass)
- [ ] mail-tester.com 10/10 score
- [ ] Client compatibility (Thunderbird, Apple Mail, Outlook)
- [ ] 100K emails/day throughput
- [ ] 500+ concurrent connections
- [ ] Memory usage within limits
- [ ] Security audit clean
- [ ] SQL injection tests pass
- [ ] Input validation tests pass
