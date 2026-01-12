# Migration Plan: SQLite to PostgreSQL

## Executive Summary

This migration plan provides a comprehensive strategy for transitioning **gomailserver** from SQLite to PostgreSQL. The migration leverages the project's existing repository pattern for clean separation of concerns and minimal disruption to business logic.

**Estimated Effort**: 40-60 hours
**Risk Level**: Medium (mitigated with proper testing)
**Recommended Approach**: Dual-phase migration with feature flags

---

## Current Architecture Analysis

### Database Layer

**Location**: `internal/database/`

- **Current Driver**: `github.com/mattn/go-sqlite3` (v1.14.32)
- **Connection**: Single-file SQLite with WAL mode enabled
- **Migrations**: Custom migration system with 8 versions (V1-V8)
- **Configuration**: File path-based (`./data/mailserver.db`)

**Key Files**:
- `sqlite.go` - Database connection and PRAGMA settings
- `migrations.go` - Migration orchestration
- `schema_v1.go` through `schema_v6.go` - Migration definitions
- `migration_v8.go` - Latest reputation tables

### Repository Layer

**Location**: `internal/repository/sqlite/`

- **Pattern**: Interface-based repository with SQLite implementation
- **Interfaces**: 12 repository interfaces (User, Domain, Message, Mailbox, etc.)
- **Query Style**: Raw SQL with parameterized queries (`?` placeholders)
- **Transaction Handling**: Manual `BEGIN`/`COMMIT` with defer rollback

### Data Volume Characteristics

- **Hybrid Storage**: Messages < 1MB stored in DB, > 1MB on filesystem
- **Tables**: 20+ tables including mailboxes, messages, queues, reputation, webhooks
- **Concurrency**: Currently limited by SQLite's single-writer model

---

## Migration Strategy

### Recommended Approach: Dual-Database Support

**Rationale**: Gradual transition allows rollback and validation

1. **Phase 1** (Week 1): Add PostgreSQL support alongside SQLite
2. **Phase 2** (Week 2): Implement data migration and validation
3. **Phase 3** (Week 3): Production cutover with monitoring
4. **Phase 4** (Week 4): Remove SQLite code path

**Evidence**: Pocket-ID multi-database implementation shows successful pattern for runtime database switching.

---

## Phase 1: Infrastructure Setup

### 1.1 Choose PostgreSQL Driver

**Recommendation**: `jackc/pgx/v5` (native or stdlib adapter)

**Why**:
- De facto standard, actively maintained
- 44-76% faster than lib/pq for bulk operations
- Native connection pooling (`pgxpool`)
- Supports PostgreSQL-specific features (JSONB, UUID, arrays)

**Installation**:
```bash
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/stdlib  # For database/sql compatibility
```

**Comparison** (from research):

| Driver | Performance | Maintenance | Features | Recommendation |
|--------|------------|--------------|-----------|----------------|
| **lib/pq** | ⚠️ Slower | ⚠️ Deprecated | Basic | ❌ Not recommended |
| **pgx** | ✅ Fastest | ✅ Active | Advanced | ✅ **RECOMMENDED** |

### 1.2 Refactor Database Layer

**Create new package**: `internal/database/postgres/`

```go
// internal/database/postgres/postgres.go
package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"
    "go.uber.org/zap"
)

type DB struct {
    *sql.DB
    logger *zap.Logger
}

type Config struct {
    Host            string
    Port            int
    Database        string
    User            string
    Password        string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

func New(cfg Config, logger *zap.Logger) (*DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password, cfg.SSLMode,
    )

    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    wrapper := &DB{DB: db, logger: logger}
    logger.Info("PostgreSQL connection established",
        zap.String("host", cfg.Host),
        zap.Int("port", cfg.Port),
        zap.String("database", cfg.Database),
    )

    return wrapper, nil
}
```

**Update configuration**:

```yaml
# gomailserver.yaml
database:
  driver: postgres  # "sqlite3" or "postgres"

  # SQLite configuration
  path: ./data/mailserver.db
  wal_enabled: true

  # PostgreSQL configuration
  postgres:
    host: localhost
    port: 5432
    database: gomailserver
    user: gomailserver
    password: ${DB_PASSWORD}  # Environment variable
    ssl_mode: disable  # "disable", "require", "verify-full"
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 1h
```

**Create database factory**:

```go
// internal/database/database.go
package database

import (
    "go.uber.org/zap"
)

type Driver string

const (
    DriverSQLite  Driver = "sqlite3"
    DriverPostgres Driver = "postgres"
)

func New(driver Driver, cfg interface{}, logger *zap.Logger) (*DB, error) {
    switch driver {
    case DriverSQLite:
        return sqlite.New(cfg.(sqlite.Config), logger)
    case DriverPostgres:
        return postgres.New(cfg.(postgres.Config), logger)
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", driver)
    }
}
```

### 1.3 Create PostgreSQL Migrations

**Convert existing SQLite migrations to PostgreSQL**:

**Key Type Mappings**:

| SQLite | PostgreSQL | Notes |
|---------|-------------|--------|
| `INTEGER PRIMARY KEY AUTOINCREMENT` | `SERIAL PRIMARY KEY` | Auto-increment |
| `TEXT` | `TEXT` | Compatible |
| `BLOB` | `BYTEA` | Binary data |
| `DATETIME` | `TIMESTAMP` | Timezone-aware: `TIMESTAMPTZ` |
| `BOOLEAN` (stored as 0/1) | `BOOLEAN` | Native boolean |

**Example migration conversion**:

```sql
-- SQLite (existing)
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    totp_enabled INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- PostgreSQL (new)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    totp_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Create PostgreSQL migration files**:

```
internal/database/postgres/migrations/
├── 000001_initial_schema.up.sql
├── 000001_initial_schema.down.sql
├── 000002_security_columns.up.sql
├── 000002_security_columns.down.sql
...
└── 000008_reputation_tables.up.sql
```

**Use migration tool**:

**Recommendation**: `golang-migrate` for PostgreSQL

**Installation**:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go get github.com/golang-migrate/migrate/v4
go get github.com/jackc/pgx/v5  # pgx driver for migrations
```

**CLI usage**:
```bash
# Create migration
migrate create -ext sql -dir ./internal/database/postgres/migrations -seq add_reputation_tables

# Run migrations
migrate -database "postgres://user:pass@localhost/gomailserver" \
          -path ./internal/database/postgres/migrations up
```

---

## Phase 2: Repository Implementation

### 2.1 Create PostgreSQL Repository Package

**New package**: `internal/repository/postgres/`

```go
// internal/repository/postgres/user_repository.go
package postgres

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/btafoya/gomailserver/internal/database"
    "github.com/btafoya/gomailserver/internal/domain"
    "github.com/btafoya/gomailserver/internal/repository"
)

type userRepository struct {
    db *database.DB
}

func NewUserRepository(db *database.DB) repository.UserRepository {
    return &userRepository{db: db}
}

// Create inserts a new user
func (r *userRepository) Create(user *domain.User) error {
    query := `
        INSERT INTO users (
            email, domain_id, password_hash, full_name, display_name, role,
            quota, used_quota, status, auth_method, totp_secret, totp_enabled,
            forward_to, auto_reply_enabled, auto_reply_subject, auto_reply_body,
            spam_threshold, language, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
        RETURNING id
    `

    err := r.db.QueryRow(query,
        user.Email, user.DomainID, user.PasswordHash, user.FullName, user.DisplayName, user.Role,
        user.Quota, user.UsedQuota, user.Status, user.AuthMethod, user.TOTPSecret, user.TOTPEnabled,
        user.ForwardTo, user.AutoReplyEnabled, user.AutoReplySubject, user.AutoReplyBody,
        user.SpamThreshold, user.Language, time.Now(), time.Now(),
    ).Scan(&user.ID)

    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    return nil
}
```

**Key Changes from SQLite**:

1. **Placeholder syntax**: `?` → `$1`, `$2`, etc.
2. **LastInsertId** → Use `RETURNING id` clause
3. **Boolean handling**: Native `BOOLEAN` instead of `INTEGER`

### 2.2 Create Repository Factory

```go
// internal/repository/factory.go
package repository

import (
    "database/sql"

    "github.com/btafoya/gomailserver/internal/database"
    postgresRepo "github.com/btafoya/gomailserver/internal/repository/postgres"
    sqliteRepo "github.com/btafoya/gomailserver/internal/repository/sqlite"
)

type Repositories struct {
    User      UserRepository
    Domain    DomainRepository
    Message   MessageRepository
    Mailbox   MailboxRepository
    Alias     AliasRepository
    Queue     QueueRepository
    // ... other repositories
}

func NewRepositories(db *database.DB, driver database.Driver) *Repositories {
    switch driver {
    case database.DriverPostgres:
        return &Repositories{
            User:    postgresRepo.NewUserRepository(db),
            Domain:  postgresRepo.NewDomainRepository(db),
            Message: postgresRepo.NewMessageRepository(db),
            Mailbox: postgresRepo.NewMailboxRepository(db),
            // ... other repositories
        }
    case database.DriverSQLite:
        return &Repositories{
            User:    sqliteRepo.NewUserRepository(db),
            Domain:  sqliteRepo.NewDomainRepository(db),
            Message: sqliteRepo.NewMessageRepository(db),
            Mailbox: sqliteRepo.NewMailboxRepository(db),
            // ... other repositories
        }
    default:
        panic("unsupported database driver")
    }
}
```

---

## Phase 3: Data Migration

### 3.1 Data Export from SQLite

**Use `sqlite3` CLI**:

```bash
# Export to SQL
sqlite3 ./data/mailserver.db .dump > sqlite_dump.sql

# Export specific tables as CSV
sqlite3 ./data/mailserver.db <<EOF
.mode csv
.headers on
.output users.csv
SELECT * FROM users;
.output domains.csv
SELECT * FROM domains;
.quit
EOF
```

### 3.2 Data Import to PostgreSQL

**Option A: Using pgloader (RECOMMENDED)**

**Why**: Automatic type conversion, progress tracking, error handling

**Installation**:
```bash
# Ubuntu/Debian
sudo apt install pgloader

# Or build from source
git clone https://github.com/dimitri/pgloader.git
cd pgloader
make install
```

**Migration command**:
```bash
pgloader sqlite://./data/mailserver.db \
          postgresql://gomailserver:password@localhost/gomailserver
```

**Progress output**:
```
2025-01-12 01:45:00 INFO  Starting pgloader version 3.6.7
2025-01-12 01:45:00 INFO  Fetching SQLite database schema
2025-01-12 01:45:01 INFO  Creating schema in PostgreSQL
2025-01-12 01:45:05 INFO  Loading table 'users'
2025-01-12 01:45:06 INFO  Loaded 152 rows from 'users' (0.987s)
2025-01-12 01:45:07 INFO  Loading table 'domains'
2025-01-12 01:45:07 INFO  Loaded 8 rows from 'domains' (0.234s)
...
2025-01-12 01:45:45 INFO  Migration complete
2025-01-12 01:45:45 INFO  Total rows migrated: 3,421
```

**Option B: Using PostgreSQL COPY**

```sql
-- After exporting to CSV
COPY users FROM '/path/to/users.csv' DELIMITER ',' CSV HEADER;
COPY domains FROM '/path/to/domains.csv' DELIMITER ',' CSV HEADER;
```

### 3.3 Validation

**Row count comparison**:

```sql
-- SQLite
SELECT COUNT(*) FROM users;

-- PostgreSQL
SELECT COUNT(*) FROM users;
```

**Data integrity checks**:

```sql
-- Check for NULL constraints violated
SELECT id FROM users WHERE email IS NULL;

-- Check foreign key constraints
SELECT u.id FROM users u
LEFT JOIN domains d ON u.domain_id = d.id
WHERE d.id IS NULL;
```

**Application-level validation**:

```bash
# Run test suite against PostgreSQL
make test DB_DRIVER=postgres

# Integration tests
./tests/integration-test.sh --db=postgres
```

---

## Phase 4: Production Migration

### 4.1 Prerequisites Checklist

- [ ] PostgreSQL server provisioned (production/staging)
- [ ] Backup of SQLite database created
- [ ] Migration scripts tested in staging
- [ ] Performance benchmarks completed
- [ ] Monitoring dashboards configured
- [ ] Rollback plan documented

### 4.2 Migration Timeline

| Time | Action |
|-------|--------|
| **T-24 hours** | Announce maintenance window to users |
| **T-1 hour** | Stop SMTP/IMAP servers, allow final queue processing |
| **T-10 minutes** | Create final SQLite backup |
| **T-5 minutes** | Run data migration with pgloader |
| **T-0 minutes** | Switch to PostgreSQL, update configuration |
| **T+10 minutes** | Start SMTP/IMAP servers |
| **T+30 minutes** | Monitor metrics, validate functionality |
| **T+1 hour** | Declare migration successful |
| **T+24 hours** | Archive SQLite database, remove SQLite code path |

### 4.3 Cutover Steps

```bash
# 1. Stop services
systemctl stop gomailserver

# 2. Create final backup
cp ./data/mailserver.db ./backups/mailserver-pre-migration-$(date +%Y%m%d-%H%M%S).db

# 3. Run data migration
pgloader sqlite://./data/mailserver.db \
          postgresql://gomailserver:${DB_PASSWORD}@localhost:5432/gomailserver

# 4. Validate migration
psql -U gomailserver -d gomailserver -c "SELECT COUNT(*) FROM users;"
psql -U gomailserver -d gomailserver -c "SELECT COUNT(*) FROM domains;"

# 5. Update configuration
sed -i 's/driver: sqlite3/driver: postgres/' /etc/gomailserver/gomailserver.yaml

# 6. Start services
systemctl start gomailserver

# 7. Monitor logs
journalctl -u gomailserver -f
```

### 4.4 Post-Migration Validation

```bash
# 1. Test SMTP sending
echo "Test email" | mail -s "Migration Test" admin@example.com

# 2. Test IMAP connectivity
telnet localhost 143
# IMAP commands...

# 3. Check queue processing
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8980/api/v1/queue

# 4. Verify web UI access
curl http://localhost:8980/admin/

# 5. Monitor error logs
journalctl -u gomailserver --since "5 minutes ago" | grep -i error
```

---

## PostgreSQL-Specific Optimizations

### 5.1 Connection Pooling

**Production configuration**:

```go
// Based on formula: (CPU cores * 2) + 1
db.SetMaxOpenConns(25)   // 4 cores * 2 + 1 = 9, rounded up to 25
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(1 * time.Hour)
db.SetConnMaxIdleTime(30 * time.Minute)
```

**Evidence**: OneUptime connection pooling guide recommends this formula.

### 5.2 Index Optimization

**Add indexes for common queries**:

```sql
-- Message lookup by user and mailbox
CREATE INDEX idx_messages_user_mailbox ON messages(user_id, mailbox_id);

-- Message search by date
CREATE INDEX idx_messages_date ON messages(internal_date DESC);

-- Queue pending lookup
CREATE INDEX idx_queue_pending ON smtp_queue(status, next_retry);

-- Greylist lookup
CREATE INDEX idx_greylist_triplet ON greylist(sender_ip, sender_email, recipient_email);
```

**Analyze query performance**:

```sql
EXPLAIN ANALYZE
SELECT * FROM messages
WHERE user_id = 1 AND mailbox_id = 5
ORDER BY internal_date DESC
LIMIT 50;
```

### 5.3 Leverage PostgreSQL Features

**JSONB for flexible data**:

```sql
-- Replace TEXT JSON columns with JSONB
ALTER TABLE messages
ALTER COLUMN categories TYPE JSONB USING categories::jsonb;

-- Add GIN index for JSONB queries
CREATE INDEX idx_messages_categories_gin ON messages USING GIN (categories);

-- Query JSONB
SELECT * FROM messages
WHERE categories @> '["Primary"]'::jsonb;
```

**Partial indexes for active data**:

```sql
-- Index only active users
CREATE INDEX idx_active_users ON users(created_at)
WHERE status = 'active';

-- Index only pending queue items
CREATE INDEX idx_queue_retry ON smtp_queue(next_retry)
WHERE status = 'pending';
```

---

## Risk Mitigation

### 6.1 Common Issues

| Issue | Mitigation |
|-------|------------|
| **Data loss during migration** | Multiple backups, validation, rollback plan |
| **Performance degradation** | Benchmark in staging, tune connection pool |
| **Application crashes** | Feature flags, gradual rollout, monitoring |
| **Query incompatibility** | Comprehensive testing, error logs review |
| **Connection pool exhaustion** | Monitor metrics, auto-tune pool size |

### 6.2 Rollback Plan

```bash
# If migration fails, rollback to SQLite

# 1. Stop services
systemctl stop gomailserver

# 2. Restore SQLite database
cp ./backups/mailserver-pre-migration-YYYYMMDD-HHMMSS.db ./data/mailserver.db

# 3. Revert configuration
sed -i 's/driver: postgres/driver: sqlite3/' /etc/gomailserver/gomailserver.yaml

# 4. Start services
systemctl start gomailserver

# 5. Verify functionality
./tests/integration-test.sh --db=sqlite3
```

### 6.3 Monitoring

**Key metrics to monitor**:

```sql
-- Connection pool usage
SELECT state, COUNT(*) FROM pg_stat_activity GROUP BY state;

-- Slow queries
SELECT query, mean_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Lock contention
SELECT relation::regclass, mode, pid
FROM pg_locks
WHERE pid != pg_backend_pid();
```

**Application metrics**:
- SMTP/IMAP connection success rate
- Queue processing time
- Web UI response time
- Error rate by service

---

## Timeline Summary

### Week 1: Infrastructure Setup
- [ ] Install PostgreSQL dependencies
- [ ] Create `internal/database/postgres/` package
- [ ] Convert migrations to PostgreSQL
- [ ] Create PostgreSQL repository package structure
- [ ] Set up staging PostgreSQL instance

### Week 2: Implementation & Testing
- [ ] Implement all PostgreSQL repositories (12 total)
- [ ] Create repository factory for driver switching
- [ ] Write tests for PostgreSQL code path
- [ ] Set up pgloader migration pipeline
- [ ] Test data migration in staging

### Week 3: Staging Validation
- [ ] Load test staging environment with production-like data
- [ ] Run full integration test suite
- [ ] Benchmark performance (SQLite vs PostgreSQL)
- [ ] Fix performance bottlenecks
- [ ] Document cutover procedure

### Week 4: Production Migration
- [ ] Announce maintenance window
- [ ] Create production PostgreSQL instance
- [ ] Perform production migration
- [ ] Monitor and validate
- [ ] Remove SQLite code path
- [ ] Archive SQLite database

---

## Recommendations

### Driver Choice

**Use `pgx/v5`** for:
- Performance (44-76% faster than lib/pq)
- Active maintenance
- Native PostgreSQL features support

**Use stdlib adapter** for:
- Compatibility with existing `database/sql` code
- Minimal repository code changes
- Easy switching between drivers

### Migration Tool

**Use `golang-migrate`** for:
- Most widely used tool
- CLI + library support
- PostgreSQL driver support

### Data Migration

**Use `pgloader`** for:
- Automatic type conversion
- Progress tracking
- Error handling and recovery

### Storage Strategy

**Keep hybrid storage**:
- Large messages (>1MB) remain on filesystem
- PostgreSQL stores metadata and small message content
- No changes to message retrieval logic needed

---

## Post-Migration Tasks

### 1. Remove SQLite Code Path

```bash
# Remove SQLite-specific files
rm -rf internal/database/sqlite.go
rm -rf internal/repository/sqlite/

# Remove SQLite dependencies
go mod tidy

# Update documentation
sed -i 's/SQLite/PostgreSQL/g' README.md
```

### 2. Update Backups

**PostgreSQL backup strategy**:

```bash
# Daily backups
pg_dump -U gomailserver -F c gomailserver > backup-$(date +%Y%m%d).dump

# Automated backup cron
0 2 * * * pg_dump -U gomailserver -F c gomailserver > /backups/daily/gomailserver-$(date +\%Y\%m\%d).dump
```

### 3. Performance Tuning

```sql
-- Update PostgreSQL configuration
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;

-- Reload configuration
SELECT pg_reload_conf();
```

### 4. Update Documentation

- Update README.md with PostgreSQL setup instructions
- Document PostgreSQL-specific features (JSONB, partial indexes)
- Add backup/restore procedures
- Update installation guides (APT, Docker, systemd)

---

## Conclusion

This migration plan provides a comprehensive, phased approach to transitioning gomailserver from SQLite to PostgreSQL. The dual-database strategy allows for gradual migration with minimal risk, while the repository pattern ensures clean separation of concerns.

**Key Success Factors**:
1. Thorough testing in staging before production cutover
2. Comprehensive backup and rollback procedures
3. Performance monitoring and optimization
4. Documentation updates for new database architecture

**Expected Benefits**:
- Improved concurrency (multi-writer support)
- Better performance at scale
- Advanced PostgreSQL features (JSONB, partial indexes)
- Easier horizontal scaling (read replicas)

---

## Appendix: Quick Reference

### Data Type Mapping Quick Reference

| SQLite | PostgreSQL | Example |
|---------|-------------|----------|
| `INTEGER PRIMARY KEY AUTOINCREMENT` | `SERIAL PRIMARY KEY` | Auto-incrementing IDs |
| `INTEGER` | `BIGINT` | Large numbers |
| `TEXT` | `TEXT` | Strings |
| `BLOB` | `BYTEA` | Binary data |
| `DATETIME` | `TIMESTAMP` | Timestamps |
| `BOOLEAN` (0/1) | `BOOLEAN` | True/false |

### SQL Placeholder Mapping

| SQLite | PostgreSQL |
|---------|-------------|
| `?` | `$1`, `$2`, `$3`... |
| `INSERT INTO t VALUES (?, ?)` | `INSERT INTO t VALUES ($1, $2)` |

### Migration Tools Summary

| Tool | Best For | Install |
|-------|-----------|---------|
| **pgloader** | Data migration | `apt install pgloader` |
| **golang-migrate** | Schema migrations | `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest` |
| **pg_dump** | PostgreSQL backups | Included with PostgreSQL |

---

**Document Version**: 1.0
**Created**: January 12, 2026
**Status**: Planning Phase
**Next Review**: After Phase 1 completion
