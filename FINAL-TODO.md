# Final TODO List - gomailserver

**Generated:** 2026-01-24
**Last Updated:** 2026-01-24
**Total TODOs Found:** 54

## Summary by Category

| Category | Count | Status |
|----------|-------|--------|
| **🟢 Compilation Errors** | 12 | ✅ Fixed |
| API Implementation | 15 | ⚠️ Pending |
| SMTP/IMAP Protocol | 8 | ⚠️ Pending |
| WebDAV/CalDAV | 6 | ⚠️ Pending |
| Reputation System | 5 | ⚠️ Pending |
| Authentication | 2 | ⚠️ Pending |
| Queue/Delivery | 3 | ⚠️ Pending |
| Configuration | 3 | ⚠️ Pending |

---

## 🟢 FIXED: Compilation Errors

All 12 compilation errors have been resolved. The project now builds successfully.

### CE-1. CalendarRepository Missing GetAll Method
**File:** `internal/calendar/repository/sqlite/calendar.go`
**Status:** ❌ Build Error
**Description:** CalendarRepository does not implement `GetAll()` method required by interface
**Impact:** Cannot compile - CalendarService and EventService instantiation fails
**Interface:** `internal/calendar/domain/calendar.go:42-43`
```go
// Required method in interface:
GetAll() ([]*Calendar, error)
```
**Resolution:** Add GetAll method to CalendarRepository:
```go
func (r *CalendarRepository) GetAll() ([]*domain.Calendar, error) {
    query := `SELECT id, user_id, name, display_name, color, description, timezone, sync_token, created_at, updated_at FROM calendars ORDER BY created_at DESC`
    rows, err := r.db.Query(query)
    // ... scan and return
}
```

### CE-2. MessageServiceInterface Signature Mismatch
**File:** `internal/service/interfaces.go:26`
**Status:** ❌ Build Error
**Description:** Interface expects `GetByMailbox(mailboxID int64)` but implementation has `GetByMailbox(int64, int, int)` with pagination
**Impact:** Cannot compile - SMTP and IMAP backends fail interface check
**Current implementation:** `internal/service/message_service.go:208`
```go
// Implementation (has pagination):
func (s *MessageService) GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error)

// Interface (no pagination):
GetByMailbox(mailboxID int64) ([]*domain.Message, error)
```
**Resolution:** Either:
1. Update interface to include pagination: `GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error)`
2. Or add a wrapper method that calls the paginated version with defaults

### CE-3. CalDAV Handler Missing UserService Parameter
**File:** `internal/commands/run.go:316`
**Status:** ❌ Build Error
**Description:** `caldav.NewHandler` called with 3 arguments but requires 4 (includes UserService)
**Impact:** Cannot compile - WebDAV server initialization fails
**Handler signature:** `internal/webdav/caldav/handler.go:27`
```go
func NewHandler(logger *zap.Logger, calendarService domain.CalendarService, eventService domain.EventService, userService *service.UserService) *Handler
```
**Resolution:** Add `userSvc` parameter to caldav.NewHandler call:
```go
caldavHandler := caldav.NewHandler(logger, calendarSvc, eventSvc, userSvc)
```

### CE-4. DomainService Missing GetDefaultTemplate Method
**File:** `internal/api/handler/domain_handler.go:602`
**Status:** ❌ Build Error
**Description:** Handler calls `h.domainService.GetDefaultTemplate()` but method doesn't exist
**Impact:** Cannot compile - Domain handler API fails
**Resolution:** Add method to DomainService:
```go
func (s *DomainService) GetDefaultTemplate() (*domain.Domain, error) {
    return s.repo.GetByName(DefaultTemplateDomainName)
}
```

### CE-5. DomainService Missing UpdateDefaultTemplate Method
**File:** `internal/api/handler/domain_handler.go:621`
**Status:** ❌ Build Error
**Description:** Handler calls `h.domainService.UpdateDefaultTemplate()` but method doesn't exist
**Impact:** Cannot compile - Domain handler API fails
**Resolution:** Add method to DomainService:
```go
func (s *DomainService) UpdateDefaultTemplate(updates *domain.Domain) error {
    template, err := s.repo.GetByName(DefaultTemplateDomainName)
    if err != nil {
        return err
    }
    // Apply updates to template
    return s.repo.Update(template)
}
```

### CE-6. Test File: Duplicate MockAddressbookService Declaration
**File:** `internal/webdav/carddav/simple_test.go:16`
**Status:** ❌ Build Error
**Description:** MockAddressbookService redeclared (also in handler_test.go:41)
**Impact:** Cannot run tests
**Resolution:** Remove duplicate declaration or rename to avoid collision

### CE-7. Test File: Missing Domain Import
**File:** `internal/webdav/carddav/simple_test.go:18-19`
**Status:** ❌ Build Error
**Description:** References `domain.Addressbook` and `domain.Contact` but missing import
**Impact:** Cannot run tests
**Resolution:** Add import: `"github.com/btafoya/gomailserver/internal/contacts/domain"`

### CE-8. Test File: Unused Import
**File:** `internal/webdav/carddav/simple_test.go:4`
**Status:** ❌ Build Error
**Description:** `encoding/xml` imported but not used
**Impact:** Cannot run tests
**Resolution:** Remove unused import or use it

### CE-9. SMTP Delivery Worker Syntax Error
**File:** `internal/smtp/delivery_worker.go:355-369`
**Status:** ❌ Build Error
**Description:** Extra closing brace at line 355, orphaned code after function end at line 367
**Impact:** Cannot compile - SMTP package broken
**Details:**
```go
// Line 355 has extra closing brace
}
}  // <-- Extra brace

// Lines 369+ have orphaned code outside any function
tlsConfig := &tls.Config{
```
**Resolution:** Remove extra closing braces and move orphaned code inside function or delete it

### CE-10. Postgres Domain Repository Syntax Error
**File:** `internal/repository/postgres/domain_repository.go:174`
**Status:** ❌ Build Error
**Description:** UPDATE query string is malformed (unclosed or missing continuation)
**Impact:** Cannot compile - Postgres repository package broken
**Resolution:** Fix the SQL query string syntax

### CE-11. Sieve Parser Duplicate Code
**File:** `pkg/sieve/parser.go:221-226`
**Status:** ❌ Build Error
**Description:** Duplicate default case and return statement (copy-paste error)
**Impact:** Cannot compile - Sieve parser package broken
**Details:**
```go
// Lines 215-220 (first occurrence):
default:
    return false, fmt.Errorf("unsupported header field: %s", condition.Field)
}
return p.evaluateStringCondition(...)
}

// Lines 221-226 (duplicate - should be deleted):
default:
    return false, fmt.Errorf("unsupported header field: %s", condition.Field)
}
return p.evaluateStringCondition(...)
}
```
**Resolution:** Delete the duplicate code block (lines 221-226)

### CE-12. Integration Test Missing Imports
**File:** `tests/integration2/integration.go`
**Status:** ❌ Build Error
**Description:** Multiple undefined references: domain, repository, context, time
**Impact:** Cannot run integration tests
**Resolution:** Add missing imports:
```go
import (
    "context"
    "time"
    "github.com/btafoya/gomailserver/internal/domain"
    "github.com/btafoya/gomailserver/internal/repository"
)
```

---

## API Implementation TODOs

### 1. Log Handler - Database Integration
**File:** `internal/api/handlers/log_handler.go:71`  
**Status:** ⚠️ Not Implemented  
**Description:** Implement actual log retrieval from database instead of returning empty logs  
**Impact:** Critical - Admin cannot view server logs  
**Resolution:** Need to implement database query for logs with filtering

### 2. Alias Handler - Pagination
**File:** `internal/api/handlers/alias_handler.go:50`  
**Status:** ⚠️ Not Implemented  
**Description:** Add pagination and filtering support for alias listing  
**Impact:** Performance issues with large alias lists  
**Resolution:** Implement pagination parameters and database queries

### 3. Reputation Handler - CSV Export
**File:** `internal/api/handlers/reputation_phase5_handler.go:281`  
**Status:** ❌ Not Implemented  
**Description:** CSV export for DMARC reports returns 501 error  
**Impact:** Users cannot export reports  
**Resolution:** Implement CSV generation for report data

### 4. Domain Handler - DKIM Key Generation
**File:** `internal/api/handlers/domain_handler.go:291`  
**Status:** ⚠️ Not Implemented  
**Description:** DKIM key generation returns placeholder response  
**Impact:** Cannot configure DKIM signing  
**Resolution:** Integrate with security/dkim package for actual key generation

### 5. User Handler - Pagination
**File:** `internal/api/handlers/user_handler.go:73`  
**Status:** ⚠️ Not Implemented  
**Description:** User list lacks pagination support  
**Impact:** Performance issues with large user bases  
**Resolution:** Add pagination with query parameters

### 6. User Handler - Domain Name Lookup
**File:** `internal/api/handlers/user_handler.go:319`  
**Status:** ⚠️ Not Implemented  
**Description:** Missing domain name in user responses  
**Impact:** Incomplete user information  
**Resolution:** Add domain service integration for name lookup

### 7. Queue Handler - Pagination & Filtering
**File:** `internal/api/handlers/queue_handler.go:50,55`  
**Status:** ⚠️ Not Implemented  
**Description:** Queue items list lacks pagination and status filtering  
**Impact:** Poor performance and usability  
**Resolution:** Implement pagination parameters and status filters

### 8. Webmail Calendar - iCalendar Library
**File:** `internal/api/handlers/webmail_calendar.go:220,301,333`  
**Status:** ❌ Critical Missing Feature  
**Description:** Manual iCalendar data generation instead of proper library  
**Impact:** Non-compliant calendar functionality  
**Resolution:** Integrate proper iCalendar library (e.g., github.com/arran4/golang-ical)

### 9. Reputation Phase 6 - IMAP Integration
**File:** `internal/api/handlers/reputation_phase6_handler.go:24,49,90,106,122,151`  
**Status:** ❌ Major Gap  
**Description:** All operational mailbox actions are mocked  
**Impact:** Reputation system cannot process real operational mail  
**Resolution:** Complete IMAP service integration

### 10. Reputation Phase 6 - Trend Calculation
**File:** `internal/api/handlers/reputation_phase6_handler.go:211`  
**Status:** ⚠️ Simplified  
**Description:** Trend calculation based on current score only  
**Impact:** Inaccurate reputation trends  
**Resolution:** Implement historical score tracking

### 11. Reputation Phase 6 - DNS Health Checks
**File:** `internal/api/handlers/reputation_phase6_handler.go:223`  
**Status:** ❌ Mock Data  
**Description:** DNS health checks return mock data  
**Impact:** No real DNS monitoring  
**Resolution:** Implement actual DNS validation

### 12. Authentication - TOTP Validation
**File:** `internal/api/handlers/auth_handler.go:91`  
**Status:** ❌ Security Gap  
**Description:** TOTP validation is skipped  
**Impact:** 2FA not functional  
**Resolution:** Integrate security/totp package for validation

### 13. Authentication - Admin Status Check
**File:** `internal/api/handlers/auth_handler.go:183`  
**Status:** ⚠️ Missing  
**Description:** JWT generation uses hardcoded "user" role  
**Impact:** Admin privileges not granted  
**Resolution:** Check actual admin status from database

---

## SMTP/IMAP Protocol TODOs

### 14. SMTP Recipient Validation
**File:** `internal/smtp/backend.go:396-399`  
**Status:** ❌ Critical Missing  
**Description:** No recipient validation, quota checks, greylisting, or rate limiting  
**Impact:** Server accepts mail for non-existent recipients, no protection  
**Resolution:** Implement all recipient validation checks

### 15. SMTP Local Delivery
**File:** `internal/smtp/delivery_worker.go:185`  
**Status:** ❌ Not Implemented  
**Description:** Local delivery returns error instead of storing messages  
**Impact:** Internal mail delivery broken  
**Resolution:** Implement MessageService integration for local delivery

### 16. SMTP Message Reading
**File:** `internal/smtp/delivery_worker.go:542`  
**Status:** ❌ Not Implemented  
**Description:** Message reading from filesystem returns empty message  
**Impact:** Cannot process queued messages  
**Resolution:** Implement file-based message reading

### 17. SMTP Recipients Parsing
**File:** `internal/smtp/delivery_worker.go:550`  
**Status:** ❌ Not Implemented  
**Description:** JSON recipient parsing returns empty slice  
**Impact:** Cannot determine message recipients  
**Resolution:** Implement proper JSON parsing

### 18. SMTP Local Domain Check
**File:** `internal/smtp/delivery_worker.go:557`  
**Status:** ❌ Not Implemented  
**Description:** Always returns false for local domain detection  
**Impact:** Cannot distinguish local vs remote delivery  
**Resolution:** Check against configured local domains

### 19. SMTP Queue Integration
**File:** `internal/service/queue_service.go:271`  
**Status:** ❌ Not Implemented  
**Description:** No SMTP client integration for message delivery  
**Impact:** Queue processing cannot deliver messages  
**Resolution:** Implement SMTP client for outbound delivery

### 20. IMAP Message Fetching
**File:** `internal/imap/mailbox.go:128-130`  
**Status:** ❌ Not Implemented  
**Description:** IMAP FETCH command returns empty results  
**Impact:** Mail clients cannot retrieve messages  
**Resolution:** Implement database message retrieval with sequence filtering

### 21. IMAP Message Storing
**File:** `internal/imap/mailbox.go:264-267`  
**Status:** ❌ Not Implemented  
**Description:** IMAP APPEND doesn't store messages  
**Impact:** Cannot save messages via IMAP  
**Resolution:** Implement message storage with flag management

---

## WebDAV/CalDAV TODOs

### 22. WebDAV MultiStatus Response
**File:** `internal/webdav/handler.go:410`  
**Status:** ⚠️ Missing Implementation  
**Description:** MultiStatus response builder not implemented  
**Impact:** PROPFIND may not work correctly  
**Resolution:** Implement proper multistatus XML generation

### 23. WebDAV Property Values
**File:** `internal/webdav/handler.go:487,493,499,505,511,546,571`  
**Status:** ❌ All Mock Data  
**Description:** All WebDAV properties return hardcoded values  
**Impact:** No real file system integration  
**Resolution:** Connect to actual storage backend

### 24. CalDAV Calendar Creation
**File:** `internal/webdav/caldav/handler.go:130,150`  
**Status:** ⚠️ Simplified  
**Description:** Calendar creation lacks permission validation and XML parsing  
**Impact:** Limited CalDAV functionality  
**Resolution:** Add proper validation and XML parsing

### 25. CalDAV Free-Busy Query
**File:** `internal/webdav/caldav/handler.go:437`  
**Status:** ⚠️ Simplified  
**Description:** Free-busy query doesn't parse time range from request  
**Impact:** Incorrect free-busy data  
**Resolution:** Parse actual XML time range

### 26. Calendar Event Recurrence
**File:** `internal/calendar/service/event.go:287`  
**Status:** ❌ Not Implemented  
**Description:** Event recurrence expansion returns single event  
**Impact:** No recurring event support  
**Resolution:** Integrate rrule-go library

### 27. CalDAV VTODO Support
**File:** `internal/webdav/handler.go:538`  
**Status:** ✅ Planned  
**Description:** VTODO component declared but not implemented  
**Impact:** Task management not available  
**Resolution:** Implement task handling

---

## Reputation System TODOs

### 28. Scheduler Configuration
**File:** `internal/reputation/scheduler.go:464,512,605,670`  
**Status:** ❌ Hardcoded Empty Lists  
**Description:** All sync operations use empty domain/IP lists  
**Impact:** No reputation data collection  
**Resolution:** Connect to configuration/database

### 29. Historical Score Tracking
**File:** `internal/reputation/service/predictions.go:135`  
**Status:** ❌ Not Implemented  
**Description:** No historical scores for trend analysis  
**Impact:** Poor reputation predictions  
**Resolution:** Implement historical score storage

---

## Configuration TODOs

### 30. DANE DNS Resolver
**File:** `internal/service/dane_service.go:34`  
**Status:** ⚠️ Hardcoded  
**Description:** DNS resolver hardcoded to Cloudflare  
**Impact:** No configurability  
**Resolution:** Add to admin configuration

### 31. Message Date Parsing
**File:** `internal/service/message_service.go:119`  
**Status:** ⚠️ Simplified  
**Description:** Uses current time instead of parsing Date header  
**Impact:** Incorrect message timestamps  
**Resolution:** Implement RFC 2822 date parsing

---

## Priority Matrix

| Priority | Count | Items |
|----------|-------|-------|
| **🔴 Build Blocking** | 12 | CE-1 through CE-12 |
| **Critical** | 12 | #14, #15, #16, #17, #18, #19, #20, #21, #8, #9, #12, #28 |
| **High** | 15 | #1, #2, #4, #5, #6, #7, #10, #11, #13, #22, #23, #24, #25, #26, #29 |
| **Medium** | 10 | #3, #27, #30, #31, and others |
| **Low** | 5 | Documentation and minor improvements |

---

## Recommended Implementation Order

### Phase 0: Fix Compilation Errors (BLOCKING)
**Must complete before any other work - project doesn't build**

**Syntax Errors (fix first - these break parsing):**
1. **CE-9: delivery_worker.go** - Remove extra braces, fix orphaned code
2. **CE-10: domain_repository.go** - Fix SQL query string syntax
3. **CE-11: parser.go** - Delete duplicate code block

**Interface/Method Errors:**
4. **CE-2: MessageServiceInterface** - Update interface signature or add wrapper method
5. **CE-1: CalendarRepository.GetAll()** - Add missing method
6. **CE-3: CalDAV NewHandler** - Add UserService parameter to call
7. **CE-4 & CE-5: DomainService methods** - Add GetDefaultTemplate and UpdateDefaultTemplate

**Test File Fixes:**
8. **CE-6, CE-7, CE-8: simple_test.go** - Remove duplicates, add imports
9. **CE-12: integration.go** - Add missing imports

### Phase 1: Core Mail Functionality (Critical)
1. SMTP recipient validation (#14)
2. IMAP message fetching/storing (#20, #21)
3. SMTP local delivery (#15)
4. SMTP queue integration (#19)

### Phase 2: Security & Authentication
5. TOTP validation (#12)
6. Admin status checks (#13)
7. DKIM key generation (#4)

### Phase 3: Reputation System
8. Scheduler configuration (#28)
9. IMAP integration for reputation (#9)
10. Historical score tracking (#29)

### Phase 4: Web Interface & APIs
11. Log retrieval (#1)
12. Pagination across APIs (#2, #5, #7)
13. iCalendar library integration (#8)

### Phase 5: Advanced Features
14. WebDAV properties (#23)
15. Calendar recurrence (#26)
16. Export functionality (#3)

---

## Risk Assessment

| Risk Area | Risk Level | Mitigation |
|-----------|------------|------------|
| **Build Status** | **🔴 Critical** | Fix all 12 compilation errors immediately |
| Mail Delivery | **High** | Prioritize SMTP/IMAP core functionality |
| Security | **High** | Implement TOTP and admin checks early |
| Reputation | **Medium** | System affects deliverability but mail still works |
| WebDAV/CalDAV | **Low** | Nice-to-have feature, doesn't break core mail |
| Performance | **Medium** | Pagination issues become problematic with scale |
| Testing | **Medium** | Test file issues prevent test suite from running |

---

## Notes

### Current State
- **⛔ PROJECT DOES NOT BUILD** - 12 compilation errors must be fixed first
- **Estimated fix time for compilation errors:** 2-3 hours
- **Total Implementation Effort:** ~3-4 months for full completion
- **Minimum Viable Product:** Phase 0 + Phase 1 items
- **Production Ready:** Phases 0-3

### Key Observations
- **All TODOs are trackable** via file and line numbers
- **Many TODOs have placeholder implementations** that return mock data
- **Critical path:** Fix syntax errors → Fix interfaces → SMTP/IMAP core → Security → APIs
- **Syntax errors** in 3 files indicate incomplete editing or merge conflicts
- **Interface mismatches** indicate rapid development without integration testing
- **Test files have quality issues** (duplicate declarations, missing imports)
- **Copy-paste errors** found in sieve parser (duplicate code blocks)

### Quick Wins
1. CE-11: Delete duplicate code in parser.go (~2 min)
2. CE-9: Fix delivery_worker.go braces (~10 min)
3. CE-8: Remove unused import (~1 min)
4. CE-3: Fix CalDAV handler call (~5 min)
5. CE-2: Fix MessageServiceInterface signature (~10 min)
6. CE-1: Add CalendarRepository.GetAll() (~15 min)
7. CE-4 & CE-5: Add DomainService template methods (~20 min)

### Dependencies Between Items
```
CE-9 (delivery_worker.go)     ──► SMTP Package compiles
CE-10 (domain_repository.go)  ──► Postgres Package compiles
CE-11 (parser.go)             ──► Sieve Package compiles

CE-2 (MessageServiceInterface) ──► SMTP Backend
                              └──► IMAP Backend

CE-1 (CalendarRepository)     ──► CalendarService
                              └──► EventService

CE-3 (CalDAV Handler)         ──► WebDAV Server

CE-4/CE-5 (DomainService)     ──► Domain Handler API

CE-12 (integration.go)        ──► Integration Tests
```

**Last Updated:** 2026-01-24