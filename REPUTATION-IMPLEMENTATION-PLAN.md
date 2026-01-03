# Reputation Management Implementation Plan

**Project**: gomailserver Mail Server
**Feature**: Automated Reputation Management System
**Created**: 2026-01-03
**Status**: Planning Phase
**Priority**: High (Critical for production deployment)

---

## 📋 Executive Summary

This plan implements comprehensive automated reputation management for gomailserver, transforming sender reputation from an external judgment into a managed engineering system with observability, adaptive policy, and automated remediation.

### Goals
- **Inbox Placement**: >90% delivery to inbox (not spam folder)
- **mail-tester.com Score**: 9+/10 rating
- **Complaint Rate**: <0.1% to prevent filtering
- **Issue Detection**: <15 minutes to identify reputation problems

### Approach
- **Iterative Development**: Flexible phased approach, building on existing foundation
- **Foundation First**: Telemetry pipeline enables all other features
- **High Automation**: Fully automatic warm-up, circuit breakers, and recovery
- **Comprehensive Monitoring**: Gmail Postmaster + Microsoft SNDS integration

---

## 🏗️ Architecture Overview

### Current Foundation (Already Implemented)
- ✅ DKIM signing/verification (RSA-2048/4096, Ed25519)
- ✅ SPF validation with IPv4/IPv6 support
- ✅ DMARC policy enforcement with alignment
- ✅ Rate limiting (per-IP, per-user, per-domain)
- ✅ Brute force protection
- ✅ Greylisting system
- ✅ Structured logging (zap JSON)
- ✅ Multiple domains support (shared IP)
- ✅ Clean architecture (repository/service pattern)

### New Components to Build

```
┌─────────────────────────────────────────────────────┐
│         Reputation Management System                 │
├─────────────────────────────────────────────────────┤
│  1. Telemetry Pipeline (separate SQLite DB)         │
│     ├─ Metrics Collection (deliveries, bounces...)  │
│     ├─ Event Aggregation (rolling windows)          │
│     └─ 90-day Retention Policy                      │
├─────────────────────────────────────────────────────┤
│  2. Deliverability Readiness Auditor                │
│     ├─ DNS Health (SPF/DKIM/DMARC/rDNS)            │
│     ├─ Reputation Scoring (0-100 per domain)        │
│     ├─ Circuit Breaker Status                       │
│     └─ Alert Management                             │
├─────────────────────────────────────────────────────┤
│  3. Adaptive Sending Policy Engine                  │
│     ├─ Enhanced Rate Limiter                        │
│     ├─ Circuit Breakers (3 triggers)                │
│     ├─ Auto Warm-up (14-30 day schedule)           │
│     └─ Auto-Resume with Backoff                    │
├─────────────────────────────────────────────────────┤
│  4. DMARC Report Processing                         │
│     ├─ RUA Report Parser (automatic)                │
│     ├─ Alignment Analysis                           │
│     ├─ Automated Corrective Actions                 │
│     └─ CSV/JSON Export                              │
├─────────────────────────────────────────────────────┤
│  5. External Feedback Integration                   │
│     ├─ Gmail Postmaster Tools API                   │
│     ├─ Microsoft SNDS API                           │
│     └─ ARF Complaint Processing                     │
├─────────────────────────────────────────────────────┤
│  6. Admin WebUI Extensions                          │
│     ├─ Dashboard: Deliverability Status             │
│     ├─ DMARC Reports Page                           │
│     └─ Operational Mailbox Inbox                    │
└─────────────────────────────────────────────────────┘
```

---

## 📊 Phased Implementation Plan

### Phase 1: Telemetry Foundation (Weeks 1-2)
**Objective**: Establish metrics collection and storage infrastructure

#### 1.1 Database Schema
**Location**: `internal/database/schema_reputation_v1.go`

```sql
-- Reputation metrics database (separate SQLite: reputation.db)

CREATE TABLE sending_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,                 -- Unix timestamp
    domain TEXT NOT NULL,                       -- Sending domain
    recipient_domain TEXT NOT NULL,             -- Receiving domain
    event_type TEXT NOT NULL,                   -- delivery|bounce|defer|complaint
    bounce_type TEXT,                           -- hard|soft|null
    enhanced_status_code TEXT,                  -- e.g., "5.1.1"
    smtp_response TEXT,                         -- Full SMTP response
    ip_address TEXT NOT NULL,                   -- Sending IP
    metadata TEXT                               -- JSON: additional context
);

CREATE INDEX idx_sending_events_timestamp ON sending_events(timestamp);
CREATE INDEX idx_sending_events_domain ON sending_events(domain);
CREATE INDEX idx_sending_events_event_type ON sending_events(event_type);

CREATE TABLE domain_reputation_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL UNIQUE,
    reputation_score INTEGER NOT NULL,          -- 0-100
    complaint_rate REAL NOT NULL,               -- Percentage
    bounce_rate REAL NOT NULL,                  -- Percentage
    delivery_rate REAL NOT NULL,                -- Percentage
    circuit_breaker_active BOOLEAN DEFAULT 0,
    circuit_breaker_reason TEXT,
    warm_up_active BOOLEAN DEFAULT 0,
    warm_up_day INTEGER DEFAULT 0,              -- Day in warm-up schedule
    last_updated INTEGER NOT NULL               -- Unix timestamp
);

CREATE TABLE warm_up_schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    day INTEGER NOT NULL,                       -- Day 1, 2, 3...
    max_volume INTEGER NOT NULL,                -- Max messages for this day
    actual_volume INTEGER DEFAULT 0,            -- Messages sent today
    created_at INTEGER NOT NULL
);

CREATE INDEX idx_warm_up_domain ON warm_up_schedules(domain);

CREATE TABLE circuit_breaker_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    trigger_type TEXT NOT NULL,                 -- complaint|bounce|block
    trigger_value REAL NOT NULL,                -- Rate or count
    threshold REAL NOT NULL,                    -- What triggered it
    paused_at INTEGER NOT NULL,                 -- Unix timestamp
    resumed_at INTEGER,                         -- Unix timestamp or null
    auto_resumed BOOLEAN DEFAULT 0,
    admin_notes TEXT
);

CREATE INDEX idx_circuit_breaker_domain ON circuit_breaker_events(domain);

CREATE TABLE retention_policy (
    id INTEGER PRIMARY KEY,
    retention_days INTEGER NOT NULL DEFAULT 90,
    last_cleanup INTEGER NOT NULL              -- Unix timestamp
);

INSERT INTO retention_policy (id, retention_days, last_cleanup)
VALUES (1, 90, strftime('%s', 'now'));
```

#### 1.2 Repository Layer
**Location**: `internal/reputation/repository/sqlite/`

Files to create:
- `events_repository.go` - CRUD for sending_events
- `scores_repository.go` - Domain reputation scoring
- `warmup_repository.go` - Warm-up schedule management
- `circuit_breaker_repository.go` - Circuit breaker event tracking

**Interface**: `internal/reputation/repository/interfaces.go`
```go
package repository

type EventsRepository interface {
    RecordEvent(ctx context.Context, event *domain.SendingEvent) error
    GetEventsInWindow(ctx context.Context, domain string, startTime, endTime int64) ([]*domain.SendingEvent, error)
    GetEventCountsByType(ctx context.Context, domain string, startTime, endTime int64) (map[string]int64, error)
    CleanupOldEvents(ctx context.Context, olderThan int64) error
}

type ScoresRepository interface {
    GetReputationScore(ctx context.Context, domain string) (*domain.ReputationScore, error)
    UpdateReputationScore(ctx context.Context, score *domain.ReputationScore) error
    ListAllScores(ctx context.Context) ([]*domain.ReputationScore, error)
}

type WarmUpRepository interface {
    GetSchedule(ctx context.Context, domain string) ([]*domain.WarmUpDay, error)
    CreateSchedule(ctx context.Context, domain string, schedule []*domain.WarmUpDay) error
    UpdateDayVolume(ctx context.Context, domain string, day int, volume int) error
    DeleteSchedule(ctx context.Context, domain string) error
}

type CircuitBreakerRepository interface {
    RecordPause(ctx context.Context, event *domain.CircuitBreakerEvent) error
    RecordResume(ctx context.Context, domain string, autoResumed bool, notes string) error
    GetActiveBreakers(ctx context.Context) ([]*domain.CircuitBreakerEvent, error)
    GetBreakerHistory(ctx context.Context, domain string, limit int) ([]*domain.CircuitBreakerEvent, error)
}
```

#### 1.3 Domain Models
**Location**: `internal/reputation/domain/models.go`

```go
package domain

type SendingEvent struct {
    ID                  int64
    Timestamp           int64
    Domain              string
    RecipientDomain     string
    EventType           EventType  // delivery, bounce, defer, complaint
    BounceType          *string    // hard, soft, nil
    EnhancedStatusCode  *string
    SMTPResponse        *string
    IPAddress           string
    Metadata            map[string]interface{}
}

type EventType string

const (
    EventDelivery  EventType = "delivery"
    EventBounce    EventType = "bounce"
    EventDefer     EventType = "defer"
    EventComplaint EventType = "complaint"
)

type ReputationScore struct {
    Domain               string
    ReputationScore      int     // 0-100
    ComplaintRate        float64 // Percentage
    BounceRate           float64 // Percentage
    DeliveryRate         float64 // Percentage
    CircuitBreakerActive bool
    CircuitBreakerReason string
    WarmUpActive         bool
    WarmUpDay            int
    LastUpdated          int64
}

type WarmUpDay struct {
    Domain        string
    Day           int
    MaxVolume     int
    ActualVolume  int
    CreatedAt     int64
}

type CircuitBreakerEvent struct {
    ID           int64
    Domain       string
    TriggerType  TriggerType  // complaint, bounce, block
    TriggerValue float64
    Threshold    float64
    PausedAt     int64
    ResumedAt    *int64
    AutoResumed  bool
    AdminNotes   string
}

type TriggerType string

const (
    TriggerComplaint TriggerType = "complaint"
    TriggerBounce    TriggerType = "bounce"
    TriggerBlock     TriggerType = "block"
)
```

#### 1.4 Telemetry Service
**Location**: `internal/reputation/service/telemetry_service.go`

**Responsibilities**:
- Collect events from SMTP delivery pipeline
- Aggregate metrics over time windows (1h, 24h, 7d)
- Calculate reputation scores
- Clean up old data per retention policy

```go
package service

type TelemetryService struct {
    eventsRepo  repository.EventsRepository
    scoresRepo  repository.ScoresRepository
    logger      *zap.Logger
}

// RecordDelivery - called from SMTP queue on successful delivery
func (s *TelemetryService) RecordDelivery(ctx context.Context, domain, recipientDomain, ip string) error

// RecordBounce - called from SMTP queue on bounce
func (s *TelemetryService) RecordBounce(ctx context.Context, domain, recipientDomain, ip string, bounceType string, statusCode, response string) error

// RecordComplaint - called from ARF processor
func (s *TelemetryService) RecordComplaint(ctx context.Context, domain, recipientDomain string) error

// CalculateReputationScore - periodic job to update scores
func (s *TelemetryService) CalculateReputationScore(ctx context.Context, domain string) (*domain.ReputationScore, error)

// CleanupOldData - periodic job for retention policy
func (s *TelemetryService) CleanupOldData(ctx context.Context) error
```

#### 1.5 Integration Points
**SMTP Queue Integration**: Modify `internal/smtp/backend.go` or queue service to call telemetry service on delivery events

**Cron Jobs**: Add periodic tasks for:
- Reputation score calculation (every 5 minutes)
- Data cleanup (daily)

---

### Phase 2: Deliverability Readiness Auditor (Weeks 3-4)
**Objective**: Real-time DNS/auth validation and dashboard display

#### 2.1 Auditor Service
**Location**: `internal/reputation/service/auditor_service.go`

**Checks to implement**:
1. **DNS Health**:
   - SPF record presence and syntax
   - DKIM selector DNS records
   - DMARC policy presence
   - rDNS (PTR) records
   - FCrDNS (forward-confirmed reverse DNS)

2. **TLS Health**:
   - Certificate validity and expiry
   - MTA-STS policy presence

3. **Operational Mailboxes**:
   - `postmaster@domain` deliverability
   - `abuse@domain` deliverability

**Example**:
```go
type AuditorService struct {
    dnsResolver   *net.Resolver
    tlsManager    *tls.Manager
    domainService *service.DomainService
}

type AuditResult struct {
    Domain          string
    Timestamp       int64
    SPFStatus       CheckStatus
    DKIMStatus      CheckStatus
    DMARCStatus     CheckStatus
    RDNSStatus      CheckStatus
    TLSStatus       CheckStatus
    PostmasterOK    bool
    AbuseOK         bool
    OverallScore    int  // 0-100
    Issues          []string
}

type CheckStatus struct {
    Passed  bool
    Message string
    Details map[string]interface{}
}
```

#### 2.2 Admin WebUI - Dashboard Integration
**Location**: `web-ui/src/components/reputation/`

**New Components**:
1. **DeliverabilityDashboard.vue**:
   - DNS health indicators (✅/❌ for each check)
   - Reputation score gauge (0-100)
   - Active circuit breakers alert box
   - Recent alerts timeline

2. **ReputationScoreCard.vue**:
   - Per-domain score display
   - 24h/7d trend graphs
   - Complaint/bounce rate meters

**API Endpoints**:
```go
// internal/api/handlers/reputation_handler.go
GET  /api/v1/reputation/audit/:domain       // Get audit results
GET  /api/v1/reputation/scores              // List all domain scores
GET  /api/v1/reputation/circuit-breakers    // Active breakers
GET  /api/v1/reputation/alerts              // Recent alerts
```

---

### Phase 3: Adaptive Sending Policy Engine (Weeks 5-6)
**Objective**: Automatic throttling, circuit breakers, and warm-up

#### 3.1 Enhanced Rate Limiter
**Location**: `internal/reputation/service/adaptive_limiter.go`

Extends existing `internal/security/ratelimit/limiter.go` with:
- **Reputation-aware limits**: Reduce rate limits for domains with poor reputation scores
- **Warm-up enforcement**: Cap sending volume during warm-up period
- **Circuit breaker integration**: Zero rate limit when paused

```go
type AdaptiveLimiter struct {
    baseLimiter       *ratelimit.Limiter
    scoresRepo        repository.ScoresRepository
    warmUpRepo        repository.WarmUpRepository
    circuitBreakerRepo repository.CircuitBreakerRepository
}

// GetLimit returns the effective rate limit for a domain
func (l *AdaptiveLimiter) GetLimit(ctx context.Context, domain string) (int, error) {
    // Check circuit breaker first
    score, _ := l.scoresRepo.GetReputationScore(ctx, domain)
    if score.CircuitBreakerActive {
        return 0, ErrCircuitBreakerActive
    }

    // Check warm-up schedule
    if score.WarmUpActive {
        schedule, _ := l.warmUpRepo.GetSchedule(ctx, domain)
        return schedule[score.WarmUpDay].MaxVolume, nil
    }

    // Adjust base limit by reputation score
    baseLimit := l.baseLimiter.GetLimit(domain)
    adjustment := float64(score.ReputationScore) / 100.0
    return int(float64(baseLimit) * adjustment), nil
}
```

#### 3.2 Circuit Breaker Logic
**Location**: `internal/reputation/service/circuit_breaker_service.go`

**Triggers** (all enabled):
1. **Complaint Rate > 0.1%** (in last 24h window)
2. **Bounce Rate > 10%** (in last 24h window)
3. **Major Provider Blocks**: 3+ consecutive 4xx/5xx from Gmail/Outlook/Yahoo

**Behavior**:
- Pause sending for affected domain
- Record event with reason
- Auto-resume after exponential backoff (1h → 2h → 4h → 8h)

```go
type CircuitBreakerService struct {
    eventsRepo         repository.EventsRepository
    scoresRepo         repository.ScoresRepository
    circuitBreakerRepo repository.CircuitBreakerRepository
    limiter            *AdaptiveLimiter
    logger             *zap.Logger
}

// CheckAndTrigger evaluates thresholds and triggers circuit breaker if needed
func (s *CircuitBreakerService) CheckAndTrigger(ctx context.Context, domain string) error

// AutoResume attempts to resume sending after backoff period
func (s *CircuitBreakerService) AutoResume(ctx context.Context) error
```

#### 3.3 Automatic Warm-Up
**Location**: `internal/reputation/service/warmup_service.go`

**Schedule** (14-30 day ramp):
```
Day 1:    100 messages
Day 2:    200 messages
Day 3:    500 messages
Day 4:    1,000 messages
Day 5:    2,000 messages
Day 6:    5,000 messages
Day 7:    10,000 messages
Day 8-14: +10,000/day
Day 15+:  No limit (warm-up complete)
```

**Auto-detection**:
- New domain added → auto-start warm-up
- New sending IP configured → auto-start warm-up

```go
type WarmUpService struct {
    warmUpRepo repository.WarmUpRepository
    scoresRepo repository.ScoresRepository
}

// StartWarmUp creates a 14-30 day schedule for a new domain/IP
func (s *WarmUpService) StartWarmUp(ctx context.Context, domain string) error

// AdvanceDay moves to next day in schedule (called daily)
func (s *WarmUpService) AdvanceDay(ctx context.Context, domain string) error

// CompleteWarmUp marks warm-up as finished
func (s *WarmUpService) CompleteWarmUp(ctx context.Context, domain string) error
```

---

### Phase 4: DMARC Report Processing (Weeks 7-8)
**Objective**: Automatic RUA parsing, analysis, and corrective actions

#### 4.1 DMARC Report Parser
**Location**: `internal/reputation/service/dmarc_reports_service.go`

**Capabilities**:
- Parse DMARC aggregate reports (RUA) from XML
- Extract alignment statistics (SPF align, DKIM align)
- Identify sending sources and IPs
- Detect SPF/DKIM failures

**Integration**: Hook into IMAP to auto-detect reports sent to `dmarc-reports@domain`

```go
type DMARCReportsService struct {
    imapService *imap.Service
    parser      *dmarc.AggregateParser
    actionsRepo repository.DMARCActionsRepository
}

// ParseReport extracts data from RUA XML
func (s *DMARCReportsService) ParseReport(ctx context.Context, xmlData []byte) (*DMARCReport, error)

// AnalyzeAlignment checks SPF/DKIM alignment issues
func (s *DMARCReportsService) AnalyzeAlignment(ctx context.Context, report *DMARCReport) (*AlignmentAnalysis, error)

// TakeCorrectiveAction auto-fixes detected issues
func (s *DMARCReportsService) TakeCorrectiveAction(ctx context.Context, issue *AlignmentIssue) error
```

#### 4.2 Database Schema Extension
```sql
CREATE TABLE dmarc_reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    report_id TEXT NOT NULL UNIQUE,          -- From DMARC report
    begin_time INTEGER NOT NULL,
    end_time INTEGER NOT NULL,
    organization TEXT,                       -- Reporter org
    total_messages INTEGER NOT NULL,
    spf_pass INTEGER NOT NULL,
    dkim_pass INTEGER NOT NULL,
    alignment_pass INTEGER NOT NULL,
    raw_xml TEXT,                            -- Full report
    processed_at INTEGER NOT NULL
);

CREATE TABLE dmarc_report_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL REFERENCES dmarc_reports(id),
    source_ip TEXT NOT NULL,
    count INTEGER NOT NULL,
    disposition TEXT,                        -- none|quarantine|reject
    spf_result TEXT,                         -- pass|fail
    dkim_result TEXT,                        -- pass|fail
    spf_aligned BOOLEAN,
    dkim_aligned BOOLEAN
);

CREATE TABLE dmarc_auto_actions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    issue_type TEXT NOT NULL,               -- spf_misalign|dkim_misalign
    description TEXT,
    action_taken TEXT,                      -- e.g., "Updated SPF record"
    taken_at INTEGER NOT NULL,
    success BOOLEAN
);
```

#### 4.3 Admin WebUI - DMARC Reports Page
**Location**: `web-ui/src/views/reputation/DMARCReports.vue`

**Features**:
1. **Summary Stats** (top of page):
   - Total reports processed
   - Overall pass rate (SPF/DKIM/Alignment)
   - Trend graph (7d/30d)

2. **Reports Table** (middle):
   - Date range, organization, messages, pass/fail breakdown
   - Click to drill down into individual report

3. **Auto-Actions Log** (bottom):
   - Recent automated fixes with timestamps
   - Success/failure status

4. **Export** (toolbar button):
   - CSV/JSON export of report data

**API Endpoints**:
```go
GET  /api/v1/reputation/dmarc/reports       // List reports
GET  /api/v1/reputation/dmarc/reports/:id   // Report details
GET  /api/v1/reputation/dmarc/stats         // Summary statistics
GET  /api/v1/reputation/dmarc/actions       // Auto-actions log
POST /api/v1/reputation/dmarc/export        // Export data
```

---

### Phase 5: External Feedback Integration (Weeks 9-10)
**Objective**: Gmail Postmaster Tools + Microsoft SNDS API integration

#### 5.1 Gmail Postmaster Tools API
**Location**: `internal/reputation/service/gmail_postmaster_service.go`

**Metrics to collect**:
- Domain reputation (High/Medium/Low)
- IP reputation
- Spam rate
- Feedback loop data
- Authentication rate (SPF/DKIM/DMARC)

**Authentication**: OAuth 2.0 with service account

```go
type GmailPostmasterService struct {
    client      *postmaster.Service  // Google API client
    domainRepo  repository.DomainRepository
    metricsRepo repository.PostmasterMetricsRepository
}

// FetchDomainReputation polls Gmail Postmaster for domain data
func (s *GmailPostmasterService) FetchDomainReputation(ctx context.Context, domain string) (*PostmasterMetrics, error)

// SyncAll fetches data for all configured domains
func (s *GmailPostmasterService) SyncAll(ctx context.Context) error
```

#### 5.2 Microsoft SNDS API
**Location**: `internal/reputation/service/microsoft_snds_service.go`

**Metrics to collect**:
- IP reputation data
- Spam trap hits
- Complaint rates
- Filtering levels

**Authentication**: API key + IP-based access

```go
type MicrosoftSNDSService struct {
    client     *http.Client
    apiKey     string
    metricsRepo repository.SNDSMetricsRepository
}

// FetchIPData polls SNDS for IP reputation
func (s *MicrosoftSNDSService) FetchIPData(ctx context.Context, ip string) (*SNDSMetrics, error)
```

#### 5.3 Cron Jobs for Sync
- Gmail Postmaster: Every 1 hour
- Microsoft SNDS: Every 6 hours

#### 5.4 Database Schema Extension
```sql
CREATE TABLE postmaster_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    fetched_at INTEGER NOT NULL,
    domain_reputation TEXT,                 -- HIGH|MEDIUM|LOW
    spam_rate REAL,
    ip_reputation TEXT,
    auth_rate REAL,                         -- SPF/DKIM/DMARC pass rate
    raw_response TEXT                       -- JSON
);

CREATE TABLE snds_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL,
    fetched_at INTEGER NOT NULL,
    spam_trap_hits INTEGER,
    complaint_rate REAL,
    filter_level TEXT,                      -- GREEN|YELLOW|RED
    raw_response TEXT                       -- JSON
);
```

---

### Phase 6: Admin WebUI Polish (Weeks 11-12)
**Objective**: Complete WebUI integration with all reputation features

#### 6.1 Operational Mailbox Inbox
**Location**: `web-ui/src/views/reputation/OperationalMail.vue`

**Features**:
- Dedicated inbox for `postmaster@*` and `abuse@*` mailboxes
- Filtered view showing only operational mail
- Quick actions: Mark as spam, forward to admin, delete
- Alert badges for unread operational messages

**Backend**:
- Filter IMAP mailboxes by operational address pattern
- Special handling for bounce messages and complaints

#### 6.2 Dashboard Enhancements
**Location**: `web-ui/src/views/admin/Dashboard.vue`

**Add Deliverability Section**:
```vue
<template>
  <v-row>
    <!-- Existing dashboard widgets -->

    <v-col cols="12" md="6">
      <DeliverabilityCard />
    </v-col>

    <v-col cols="12" md="6">
      <CircuitBreakersCard />
    </v-col>

    <v-col cols="12">
      <RecentAlertsTimeline />
    </v-col>
  </v-row>
</template>
```

**Components**:
1. **DeliverabilityCard**: Shows DNS health, reputation score
2. **CircuitBreakersCard**: Active pauses with resume actions
3. **RecentAlertsTimeline**: Last 10 reputation events

#### 6.3 Alert System
**Location**: `internal/reputation/service/alerts_service.go`

**Alert Types**:
- DNS validation failures
- Reputation score drops >20 points
- Circuit breaker triggered
- External feedback deterioration (Gmail/SNDS)
- DMARC alignment issues

**Delivery**:
- In-app notifications (WebUI badge)
- Email to admin (if configured)
- Webhook (if configured)

---

## 🔧 Technical Implementation Details

### Database Architecture

**Primary Database** (`gomailserver.db`): Existing mail data
**Reputation Database** (`reputation.db`): Separate database for:
- Sending events
- Reputation scores
- Warm-up schedules
- Circuit breaker events
- DMARC reports
- External metrics (Gmail Postmaster, SNDS)

**Rationale**: Isolation prevents reputation queries from impacting mail performance

### Service Layer Architecture

**New Services**:
1. **TelemetryService**: Event collection and aggregation
2. **AuditorService**: DNS/auth validation
3. **CircuitBreakerService**: Pause/resume logic
4. **WarmUpService**: Automatic volume ramping
5. **AdaptiveLimiter**: Reputation-aware rate limiting
6. **DMARCReportsService**: RUA parsing and analysis
7. **GmailPostmasterService**: External metrics sync
8. **MicrosoftSNDSService**: External metrics sync
9. **AlertsService**: Notification management

**Coordination**:
- All services use repository pattern
- Event-driven architecture where possible
- Cron jobs for periodic tasks

### Integration with Existing Code

**SMTP Backend** (`internal/smtp/backend.go`):
```go
// Add telemetry calls
func (b *Backend) Send(from string, to []string, r io.Reader) error {
    // ... existing send logic ...

    // Record delivery event
    b.telemetryService.RecordDelivery(ctx, domain, recipientDomain, ip)

    // Check rate limit (now reputation-aware)
    limit, err := b.adaptiveLimiter.GetLimit(ctx, domain)
    if err == ErrCircuitBreakerActive {
        return smtp.ErrPaused
    }

    // ... continue send ...
}
```

**Queue Service** (bounce handling):
```go
// Record bounce events
func (q *QueueService) ProcessBounce(msg *Message) error {
    // ... existing bounce logic ...

    q.telemetryService.RecordBounce(ctx, domain, recipientDomain, ip, bounceType, statusCode, response)

    // Trigger circuit breaker check
    q.circuitBreakerService.CheckAndTrigger(ctx, domain)
}
```

### Cron Jobs

**Required Periodic Tasks**:
```go
// internal/reputation/cron/scheduler.go

func SetupCronJobs(services *Services) {
    // Every 5 minutes: Calculate reputation scores
    cron.Schedule("*/5 * * * *", services.Telemetry.CalculateAllScores)

    // Every 15 minutes: Check circuit breaker thresholds
    cron.Schedule("*/15 * * * *", services.CircuitBreaker.CheckAndTrigger)

    // Every 1 hour: Auto-resume paused domains
    cron.Schedule("0 * * * *", services.CircuitBreaker.AutoResume)

    // Every 1 hour: Sync Gmail Postmaster data
    cron.Schedule("0 * * * *", services.GmailPostmaster.SyncAll)

    // Every 6 hours: Sync Microsoft SNDS data
    cron.Schedule("0 */6 * * *", services.MicrosoftSNDS.SyncAll)

    // Daily at 2 AM: Clean up old telemetry data
    cron.Schedule("0 2 * * *", services.Telemetry.CleanupOldData)

    // Daily at midnight: Advance warm-up day
    cron.Schedule("0 0 * * *", services.WarmUp.AdvanceDayForAll)
}
```

---

## 📈 Success Metrics & Validation

### KPIs to Track

1. **Inbox Placement Rate**: >90% (mail-tester.com + user feedback)
2. **mail-tester.com Score**: 9+/10 consistently
3. **Complaint Rate**: <0.1% (calculated from telemetry)
4. **Bounce Rate**: <5% (hard + soft bounces)
5. **Reputation Score**: >80/100 for all domains
6. **Issue Detection Time**: <15 minutes from problem start to alert

### Testing Plan

#### Integration Testing
1. **Telemetry Flow**:
   - Send test emails, verify events recorded
   - Trigger bounces, verify bounce events
   - Simulate complaints, verify complaint events

2. **Circuit Breaker**:
   - Force high complaint rate, verify pause
   - Wait for backoff, verify auto-resume
   - Manually resume, verify override works

3. **Warm-Up**:
   - Add new domain, verify schedule created
   - Send beyond limit, verify rejection
   - Advance days, verify volume increases

4. **DMARC Reports**:
   - Import sample RUA XML, verify parsing
   - Check alignment analysis, verify correctness
   - Trigger auto-action, verify SPF/DKIM update

5. **External APIs**:
   - Mock Gmail Postmaster response, verify parsing
   - Mock Microsoft SNDS response, verify parsing
   - Verify metrics stored correctly

#### External Validation
1. **mail-tester.com**: Test score with production config
2. **Gmail Deliverability**: Monitor Gmail Postmaster Tools for 30 days
3. **Microsoft Deliverability**: Monitor SNDS for 30 days
4. **Real User Feedback**: Collect inbox vs spam folder reports

---

## 🚀 Rollout Strategy

### Development Environment
1. Implement all phases incrementally
2. Use test domains and test IPs
3. Mock external API responses initially
4. Local SQLite databases for testing

### Staging Environment
1. Deploy with real DNS but test domains
2. Connect to real Gmail Postmaster (test domain)
3. Connect to real Microsoft SNDS (test IP)
4. Monitor for 1 week, validate metrics accuracy

### Production Rollout
1. **Soft Launch**: Enable telemetry only (read-only mode)
   - Collect data for 7 days
   - Verify no performance impact
   - Validate metric accuracy

2. **Phase 1 Enable**: Deliverability Auditor + Dashboard
   - Display-only, no enforcement
   - Gather admin feedback

3. **Phase 2 Enable**: Circuit Breakers (manual mode)
   - Alerts only, admin manually pauses
   - Validate thresholds are correct

4. **Phase 3 Enable**: Full Automation
   - Auto-pause with backoff
   - Auto-resume enabled
   - Warm-up automation for new domains

5. **Phase 4 Enable**: External Integrations
   - Gmail Postmaster sync
   - Microsoft SNDS sync
   - DMARC report processing

### Rollback Plan
- Separate database allows clean rollback
- Feature flags for each phase
- Disable cron jobs to stop automation
- Revert SMTP integration to remove telemetry calls

---

## 📚 Dependencies & Resources

### External Libraries (Go)
```go
// DMARC report parsing
"github.com/emersion/go-message/mail"  // Already used

// Gmail Postmaster API
"google.golang.org/api/gmailpostmastertools/v1"

// Microsoft Graph API (if using newer API)
"github.com/microsoftgraph/msgraph-sdk-go"

// Cron scheduling
"github.com/robfig/cron/v3"
```

### Admin WebUI (Vue.js)
```json
{
  "dependencies": {
    "chart.js": "^4.4.0",           // Graphs and charts
    "vue-chartjs": "^5.3.0",        // Vue wrapper for Chart.js
    "date-fns": "^2.30.0"           // Date formatting
  }
}
```

### Configuration

**gomailserver.conf** additions:
```yaml
reputation:
  enabled: true
  database_path: "/var/lib/gomailserver/reputation.db"
  retention_days: 90

  circuit_breaker:
    enabled: true
    complaint_threshold: 0.1        # 0.1% complaint rate
    bounce_threshold: 10.0          # 10% bounce rate
    block_threshold: 3              # 3 consecutive blocks
    auto_resume: true
    backoff_schedule: [1h, 2h, 4h, 8h]

  warm_up:
    enabled: true
    auto_detect: true
    schedule: [100, 200, 500, 1000, 2000, 5000, 10000]  # Day 1-7

  external_apis:
    gmail_postmaster:
      enabled: true
      service_account_key: "/etc/gomailserver/gmail-postmaster-sa.json"
      sync_interval: "1h"

    microsoft_snds:
      enabled: true
      api_key: "${MICROSOFT_SNDS_API_KEY}"
      sync_interval: "6h"

  alerts:
    email_enabled: true
    email_to: "admin@example.com"
    webhook_enabled: false
    webhook_url: ""
```

---

## 🎯 Next Steps

### Immediate Actions
1. ✅ **Review this plan** with stakeholders
2. ✅ **Finalize priorities** for phasing (already decided: Telemetry first)
3. ⏳ **Set up development environment** with test domains
4. ⏳ **Create GitHub issues** for each phase (optional)

### Phase 1 Kickoff Checklist
- [ ] Create `reputation.db` schema migration
- [ ] Implement repository layer interfaces
- [ ] Build domain models
- [ ] Create TelemetryService with basic event recording
- [ ] Add integration points to SMTP backend
- [ ] Write unit tests for repository layer
- [ ] Test event recording end-to-end

### Long-Term Roadmap (Beyond MVP)
- **AI/ML Integration**: Predict reputation issues before they happen
- **Multi-IP Support**: Per-domain dedicated IPs (deferred per user preference)
- **Advanced Analytics**: Cohort analysis, A/B testing for sending strategies
- **Competitive Intelligence**: Benchmark against industry averages

---

## 📝 Open Questions & Decisions Needed

### Resolved
- ✅ Storage backend: Separate SQLite database
- ✅ DMARC priority: High - automatic processing
- ✅ Retention: 90 days
- ✅ Rate limiting approach: Enhance existing rate limiter
- ✅ Circuit breaker triggers: All three (complaint, bounce, blocks)
- ✅ Warm-up: Fully automatic
- ✅ Resume policy: Auto-resume after backoff
- ✅ Gmail/SNDS: Both high priority (Phase 1-2)
- ✅ Multi-IP: No - keep shared IP simple

### Still Open
- ⏳ Gmail Postmaster service account setup - **needs credentials**
- ⏳ Microsoft SNDS API key - **needs registration**
- ⏳ Exact warm-up schedule tuning - **validate 14-30 day range**
- ⏳ Circuit breaker thresholds fine-tuning - **may need adjustment after initial data**

---

## 📞 Support & Resources

### Documentation References
- [REPUTATION-MANAGEMENT.md](/home/btafoya/projects/gomailserver/REPUTATION-MANAGEMENT.md) - Original requirements
- [PROJECT-STATUS.md](/home/btafoya/projects/gomailserver/PROJECT-STATUS.md) - Current project state
- [Gmail Postmaster Tools](https://postmaster.google.com/)
- [Microsoft SNDS](https://sendersupport.olc.protection.outlook.com/snds/)
- [DMARC.org](https://dmarc.org/)

### Context7 References Used
- Postfix rate limiting patterns (`/vdukhovni/postfix`)
- Email authentication standards
- SMTP best practices

---

**Document Status**: ✅ Complete - Ready for Implementation
**Next Review Date**: After Phase 1 completion
**Owner**: btafoya
**Last Updated**: 2026-01-03
