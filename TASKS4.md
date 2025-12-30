# Phase 4: CalDAV/CardDAV (Weeks 13-15)

**Status**: Not Started
**Priority**: Full Feature (Post-MVP)
**Estimated Duration**: 2-3 weeks
**Dependencies**: Phase 3 (Web Interfaces)

---

## Overview

Implement WebDAV-based protocols for calendar (CalDAV) and contact (CardDAV) synchronization, supporting major clients including Thunderbird, Apple, iOS, Android, and Outlook.

---

## 4.1 WebDAV Foundation [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| WD-001 | WebDAV base protocol (RFC 4918) | [ ] | API-001 |
| WD-002 | PROPFIND implementation | [ ] | WD-001 |
| WD-003 | PROPPATCH implementation | [ ] | WD-001 |
| WD-004 | MKCOL implementation | [ ] | WD-001 |
| WD-005 | DELETE implementation | [ ] | WD-001 |
| WD-006 | COPY/MOVE implementation | [ ] | WD-001 |

### WD-001: WebDAV Server

```go
// internal/webdav/server.go
package webdav

import (
    "net/http"
)

type Server struct {
    handler http.Handler
}

func NewServer(caldav *CalDAVHandler, carddav *CardDAVHandler) *Server {
    mux := http.NewServeMux()

    // CalDAV endpoints
    mux.Handle("/caldav/", caldav)

    // CardDAV endpoints
    mux.Handle("/carddav/", carddav)

    // Well-known redirects (RFC 6764)
    mux.HandleFunc("/.well-known/caldav", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/caldav/", http.StatusMovedPermanently)
    })
    mux.HandleFunc("/.well-known/carddav", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/carddav/", http.StatusMovedPermanently)
    })

    return &Server{handler: mux}
}
```

### WD-002: PROPFIND Handler

```go
// internal/webdav/propfind.go
package webdav

import (
    "encoding/xml"
)

type PropFind struct {
    XMLName  xml.Name  `xml:"DAV: propfind"`
    AllProp  *struct{} `xml:"allprop"`
    PropName *struct{} `xml:"propname"`
    Prop     *Prop     `xml:"prop"`
}

type Prop struct {
    ResourceType     *struct{} `xml:"DAV: resourcetype"`
    DisplayName      *struct{} `xml:"DAV: displayname"`
    GetContentType   *struct{} `xml:"DAV: getcontenttype"`
    GetETag          *struct{} `xml:"DAV: getetag"`
    GetLastModified  *struct{} `xml:"DAV: getlastmodified"`
    CurrentUserPrincipal *struct{} `xml:"DAV: current-user-principal"`
    // CalDAV specific
    CalendarData     *struct{} `xml:"urn:ietf:params:xml:ns:caldav calendar-data"`
    CalendarHomeSet  *struct{} `xml:"urn:ietf:params:xml:ns:caldav calendar-home-set"`
    // CardDAV specific
    AddressData      *struct{} `xml:"urn:ietf:params:xml:ns:carddav address-data"`
    AddressbookHomeSet *struct{} `xml:"urn:ietf:params:xml:ns:carddav addressbook-home-set"`
}

type MultiStatus struct {
    XMLName   xml.Name   `xml:"DAV: multistatus"`
    Responses []Response `xml:"response"`
}

type Response struct {
    Href     string   `xml:"href"`
    PropStat PropStat `xml:"propstat"`
}

type PropStat struct {
    Prop   interface{} `xml:"prop"`
    Status string      `xml:"status"`
}

func (h *BaseHandler) HandlePropFind(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var pf PropFind
    if err := xml.NewDecoder(r.Body).Decode(&pf); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    depth := r.Header.Get("Depth")
    if depth == "" {
        depth = "infinity"
    }

    // Build response based on path and depth
    responses := h.buildPropFindResponse(r.URL.Path, depth, &pf)

    // Write response
    w.Header().Set("Content-Type", "application/xml; charset=utf-8")
    w.WriteHeader(http.StatusMultiStatus)

    ms := MultiStatus{Responses: responses}
    xml.NewEncoder(w).Encode(ms)
}
```

---

## 4.2 CalDAV Server [FULL]

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

### CD-001: CalDAV Handler

```go
// internal/caldav/handler.go
package caldav

import (
    "net/http"
)

type Handler struct {
    calendarService *service.CalendarService
    eventService    *service.EventService
    userService     *service.UserService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Authenticate
    user, err := h.authenticate(r)
    if err != nil {
        w.Header().Set("WWW-Authenticate", `Basic realm="CalDAV"`)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    switch r.Method {
    case "OPTIONS":
        h.handleOptions(w, r)
    case "PROPFIND":
        h.handlePropFind(w, r, user)
    case "PROPPATCH":
        h.handlePropPatch(w, r, user)
    case "MKCALENDAR":
        h.handleMkCalendar(w, r, user)
    case "REPORT":
        h.handleReport(w, r, user)
    case "PUT":
        h.handlePut(w, r, user)
    case "GET":
        h.handleGet(w, r, user)
    case "DELETE":
        h.handleDelete(w, r, user)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Allow", "OPTIONS, GET, HEAD, PUT, DELETE, PROPFIND, PROPPATCH, MKCALENDAR, REPORT")
    w.Header().Set("DAV", "1, 2, 3, calendar-access")
    w.WriteHeader(http.StatusOK)
}
```

### CD-003: Event Storage

```go
// internal/domain/event.go
package domain

type Calendar struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id"`
    Name        string    `json:"name"`
    DisplayName string    `json:"display_name"`
    Color       string    `json:"color"`
    Description string    `json:"description"`
    Timezone    string    `json:"timezone"`
    SyncToken   string    `json:"sync_token"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Event struct {
    ID          int64     `json:"id"`
    CalendarID  int64     `json:"calendar_id"`
    UID         string    `json:"uid"`
    Summary     string    `json:"summary"`
    Description string    `json:"description"`
    Location    string    `json:"location"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    AllDay      bool      `json:"all_day"`
    RRule       string    `json:"rrule"`       // Recurrence rule
    Attendees   string    `json:"attendees"`   // JSON array
    Organizer   string    `json:"organizer"`
    Status      string    `json:"status"`      // CONFIRMED, TENTATIVE, CANCELLED
    Sequence    int       `json:"sequence"`
    ETag        string    `json:"etag"`
    ICalData    string    `json:"ical_data"`   // Raw iCalendar data
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// SQL Schema
/*
CREATE TABLE calendars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    display_name TEXT,
    color TEXT DEFAULT '#3788d8',
    description TEXT,
    timezone TEXT DEFAULT 'UTC',
    sync_token TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    calendar_id INTEGER NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    uid TEXT NOT NULL,
    summary TEXT,
    description TEXT,
    location TEXT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    all_day INTEGER DEFAULT 0,
    rrule TEXT,
    attendees TEXT,
    organizer TEXT,
    status TEXT DEFAULT 'CONFIRMED',
    sequence INTEGER DEFAULT 0,
    etag TEXT NOT NULL,
    ical_data TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(calendar_id, uid)
);

CREATE INDEX idx_events_calendar ON events(calendar_id);
CREATE INDEX idx_events_start ON events(start_time);
CREATE INDEX idx_events_uid ON events(uid);
*/
```

### CD-006: Recurring Events

```go
// internal/caldav/recurrence.go
package caldav

import (
    "github.com/teambition/rrule-go"
)

type RecurrenceExpander struct{}

func (e *RecurrenceExpander) ExpandOccurrences(event *domain.Event, start, end time.Time) ([]time.Time, error) {
    if event.RRule == "" {
        return []time.Time{event.StartTime}, nil
    }

    rule, err := rrule.StrToRRule(event.RRule)
    if err != nil {
        return nil, err
    }

    // Set DTSTART
    rule.DTStart(event.StartTime)

    // Get occurrences in range
    occurrences := rule.Between(start, end, true)

    return occurrences, nil
}

func (e *RecurrenceExpander) ParseRRule(icalRRule string) (*rrule.RRule, error) {
    return rrule.StrToRRule(icalRRule)
}
```

### CD-009: Free/Busy

```go
// internal/caldav/freebusy.go
package caldav

type FreeBusyPeriod struct {
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
    Type   string    `json:"type"` // BUSY, FREE, BUSY-TENTATIVE, BUSY-UNAVAILABLE
}

func (h *Handler) GetFreeBusy(userID int64, start, end time.Time) ([]FreeBusyPeriod, error) {
    events, err := h.eventService.GetEventsInRange(userID, start, end)
    if err != nil {
        return nil, err
    }

    periods := []FreeBusyPeriod{}

    for _, event := range events {
        if event.Status == "CANCELLED" {
            continue
        }

        fbType := "BUSY"
        if event.Status == "TENTATIVE" {
            fbType = "BUSY-TENTATIVE"
        }

        // Handle recurring events
        occurrences, _ := h.recurrence.ExpandOccurrences(event, start, end)
        duration := event.EndTime.Sub(event.StartTime)

        for _, occ := range occurrences {
            periods = append(periods, FreeBusyPeriod{
                Start: occ,
                End:   occ.Add(duration),
                Type:  fbType,
            })
        }
    }

    return mergePeriods(periods), nil
}
```

---

## 4.3 CardDAV Server [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| CAR-001 | CardDAV protocol (RFC 6352) | [ ] | WD-001 |
| CAR-002 | Address book collection management | [ ] | CAR-001 |
| CAR-003 | Contact storage (RFC 6350 vCard) | [ ] | CAR-002 |
| CAR-004 | Contact search | [ ] | CAR-003 |
| CAR-005 | Contact groups | [ ] | CAR-003 |
| CAR-006 | Distribution lists | [ ] | CAR-005 |

### CAR-001: CardDAV Handler

```go
// internal/carddav/handler.go
package carddav

type Handler struct {
    addressbookService *service.AddressbookService
    contactService     *service.ContactService
    userService        *service.UserService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    user, err := h.authenticate(r)
    if err != nil {
        w.Header().Set("WWW-Authenticate", `Basic realm="CardDAV"`)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    switch r.Method {
    case "OPTIONS":
        h.handleOptions(w, r)
    case "PROPFIND":
        h.handlePropFind(w, r, user)
    case "PROPPATCH":
        h.handlePropPatch(w, r, user)
    case "MKCOL":
        h.handleMkCol(w, r, user)
    case "REPORT":
        h.handleReport(w, r, user)
    case "PUT":
        h.handlePut(w, r, user)
    case "GET":
        h.handleGet(w, r, user)
    case "DELETE":
        h.handleDelete(w, r, user)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Allow", "OPTIONS, GET, HEAD, PUT, DELETE, PROPFIND, PROPPATCH, MKCOL, REPORT")
    w.Header().Set("DAV", "1, 2, 3, addressbook")
    w.WriteHeader(http.StatusOK)
}
```

### CAR-003: Contact Storage

```go
// internal/domain/contact.go
package domain

type Addressbook struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id"`
    Name        string    `json:"name"`
    DisplayName string    `json:"display_name"`
    Description string    `json:"description"`
    SyncToken   string    `json:"sync_token"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Contact struct {
    ID            int64     `json:"id"`
    AddressbookID int64     `json:"addressbook_id"`
    UID           string    `json:"uid"`
    FullName      string    `json:"full_name"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    Nickname      string    `json:"nickname"`
    Email         string    `json:"email"`
    Emails        string    `json:"emails"`        // JSON array
    Phone         string    `json:"phone"`
    Phones        string    `json:"phones"`        // JSON array
    Organization  string    `json:"organization"`
    Title         string    `json:"title"`
    Birthday      string    `json:"birthday"`
    Photo         []byte    `json:"photo"`
    Notes         string    `json:"notes"`
    Categories    string    `json:"categories"`    // JSON array
    ETag          string    `json:"etag"`
    VCardData     string    `json:"vcard_data"`    // Raw vCard data
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

// SQL Schema
/*
CREATE TABLE addressbooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    display_name TEXT,
    description TEXT,
    sync_token TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

CREATE TABLE contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    addressbook_id INTEGER NOT NULL REFERENCES addressbooks(id) ON DELETE CASCADE,
    uid TEXT NOT NULL,
    full_name TEXT,
    first_name TEXT,
    last_name TEXT,
    nickname TEXT,
    email TEXT,
    emails TEXT,
    phone TEXT,
    phones TEXT,
    organization TEXT,
    title TEXT,
    birthday TEXT,
    photo BLOB,
    notes TEXT,
    categories TEXT,
    etag TEXT NOT NULL,
    vcard_data TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(addressbook_id, uid)
);

CREATE INDEX idx_contacts_addressbook ON contacts(addressbook_id);
CREATE INDEX idx_contacts_email ON contacts(email);
CREATE INDEX idx_contacts_name ON contacts(full_name);
*/
```

### CAR-004: Contact Search

```go
// internal/carddav/search.go
package carddav

type AddressbookQuery struct {
    Filter   *Filter
    Limit    int
    PropName bool
}

type Filter struct {
    Test  string       // anyof, allof
    Props []PropFilter
}

type PropFilter struct {
    Name      string
    TextMatch *TextMatch
}

type TextMatch struct {
    Collation string
    MatchType string // contains, starts-with, ends-with, equals
    Value     string
}

func (h *Handler) SearchContacts(addressbookID int64, query *AddressbookQuery) ([]*domain.Contact, error) {
    if query.Filter == nil {
        return h.contactService.ListByAddressbook(addressbookID, query.Limit)
    }

    // Build SQL query from filter
    conditions := []string{}
    args := []interface{}{addressbookID}

    for _, prop := range query.Filter.Props {
        if prop.TextMatch == nil {
            continue
        }

        field := mapPropToField(prop.Name)
        if field == "" {
            continue
        }

        switch prop.TextMatch.MatchType {
        case "contains":
            conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
            args = append(args, "%"+prop.TextMatch.Value+"%")
        case "starts-with":
            conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
            args = append(args, prop.TextMatch.Value+"%")
        case "equals":
            conditions = append(conditions, fmt.Sprintf("%s = ?", field))
            args = append(args, prop.TextMatch.Value)
        }
    }

    return h.contactService.SearchWithConditions(addressbookID, conditions, args, query.Filter.Test, query.Limit)
}
```

---

## 4.4 Client Compatibility [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| CC-001 | Thunderbird compatibility testing | [ ] | CD-001, CAR-001 |
| CC-002 | Apple Mail/Calendar/Contacts testing | [ ] | CD-001, CAR-001 |
| CC-003 | iOS compatibility testing | [ ] | CD-001, CAR-001 |
| CC-004 | Android (DAVx5) testing | [ ] | CD-001, CAR-001 |
| CC-005 | Microsoft Outlook testing | [ ] | CD-001, CAR-001 |
| CC-006 | Evolution testing | [ ] | CD-001, CAR-001 |

### Client Discovery Configuration

```go
// internal/webdav/discovery.go
package webdav

// Well-known service discovery (RFC 6764)
func (s *Server) SetupDiscovery() {
    // DNS SRV records (documented in setup wizard)
    // _caldavs._tcp.example.com. 443 mail.example.com
    // _carddavs._tcp.example.com. 443 mail.example.com

    // TXT records for path
    // _caldavs._tcp.example.com. path=/caldav
    // _carddavs._tcp.example.com. path=/carddav
}

// Client-specific quirks
type ClientQuirks struct {
    // Thunderbird needs specific header handling
    ThunderbirdCompat bool

    // Apple requires certain namespaces
    AppleCompat bool

    // Outlook has its own expectations
    OutlookCompat bool
}

func DetectClient(r *http.Request) string {
    ua := r.Header.Get("User-Agent")

    switch {
    case strings.Contains(ua, "Thunderbird"):
        return "thunderbird"
    case strings.Contains(ua, "DAVx5") || strings.Contains(ua, "DAVdroid"):
        return "davx5"
    case strings.Contains(ua, "iOS") || strings.Contains(ua, "macOS"):
        return "apple"
    case strings.Contains(ua, "Microsoft Outlook"):
        return "outlook"
    case strings.Contains(ua, "Evolution"):
        return "evolution"
    default:
        return "generic"
    }
}
```

---

## Acceptance Criteria

### WebDAV
- [ ] PROPFIND returns correct properties
- [ ] PROPPATCH updates properties
- [ ] MKCOL creates collections
- [ ] DELETE removes resources
- [ ] ETags handled correctly

### CalDAV
- [ ] Calendars can be created/deleted
- [ ] Events can be added/modified/deleted
- [ ] REPORT queries work correctly
- [ ] Recurring events expand properly
- [ ] Free/busy information accurate
- [ ] Sharing permissions work

### CardDAV
- [ ] Address books can be created/deleted
- [ ] Contacts can be added/modified/deleted
- [ ] vCard parsing works
- [ ] Search returns correct results
- [ ] Contact groups work

### Client Compatibility
- [ ] Thunderbird syncs calendars and contacts
- [ ] Apple devices sync properly
- [ ] iOS devices sync properly
- [ ] DAVx5 (Android) syncs properly
- [ ] Outlook can connect (if supported)
- [ ] Evolution syncs properly

---

## Go Dependencies for Phase 4

```go
// Additional go.mod entries
require (
    github.com/emersion/go-ical v0.0.0-20240127095438-fc1c9d8fb2b6
    github.com/emersion/go-vcard v0.0.0-20230815062825-8fda7d206ec9
    github.com/teambition/rrule-go v1.8.2
)
```

---

## Next Phase

After completing Phase 4, proceed to [TASKS5.md](TASKS5.md) - Advanced Security.
