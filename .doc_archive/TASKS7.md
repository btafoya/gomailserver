# TASKS7.md - Phase 7: Webmail Client (Weeks 20-25)

## Overview

Gmail-style webmail client with rich features including conversation threading, PGP encryption in browser, contact/calendar integration.

**Total Tasks**: 34
**Priority**: [FULL] - Post-MVP feature (Health endpoints + API versioning: MVP)
**Dependencies**: Phases 0-3, 4 (CalDAV/CardDAV), 5 (PGP)

---

## 7.0 API Versioning Strategy [MVP]

Establish API versioning strategy for backward compatibility and graceful evolution.

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| API-VER-001 | API versioning implementation | [ ] | - | MVP |

---

### API-VER-001: API Versioning Implementation

**File**: `internal/api/router.go`
```go
package api

import (
    "github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, handlers *Handlers) {
    // API v1 routes (current stable)
    v1 := e.Group("/api/v1")

    // Authentication (Phase 1)
    v1.POST("/auth/login", handlers.Auth.Login)
    v1.POST("/auth/logout", handlers.Auth.Logout)
    v1.POST("/auth/refresh", handlers.Auth.RefreshToken)
    v1.GET("/auth/me", handlers.Auth.GetCurrentUser, requireAuth)

    // Users (Phase 1)
    v1.GET("/users", handlers.User.List, requireAdmin)
    v1.POST("/users", handlers.User.Create, requireAdmin)
    v1.GET("/users/:id", handlers.User.Get, requireAuth)
    v1.PUT("/users/:id", handlers.User.Update, requireAuth)
    v1.DELETE("/users/:id", handlers.User.Delete, requireAdmin)

    // Domains (Phase 1)
    v1.GET("/domains", handlers.Domain.List, requireAdmin)
    v1.POST("/domains", handlers.Domain.Create, requireAdmin)
    v1.PUT("/domains/:id", handlers.Domain.Update, requireAdmin)
    v1.DELETE("/domains/:id", handlers.Domain.Delete, requireAdmin)

    // Queue (Phase 3)
    v1.GET("/queue", handlers.Queue.List, requireAdmin)
    v1.GET("/queue/:id", handlers.Queue.Get, requireAdmin)
    v1.POST("/queue/:id/retry", handlers.Queue.Retry, requireAdmin)
    v1.DELETE("/queue/:id", handlers.Queue.Delete, requireAdmin)

    // Webmail (Phase 7)
    v1.GET("/mailboxes", handlers.Mailbox.List, requireAuth)
    v1.GET("/mailboxes/:name/messages", handlers.Mailbox.ListMessages, requireAuth)
    v1.GET("/messages/:id", handlers.Message.Get, requireAuth)
    v1.POST("/messages", handlers.Message.Send, requireAuth)
    v1.PUT("/messages/:id", handlers.Message.Update, requireAuth)
    v1.DELETE("/messages/:id", handlers.Message.Delete, requireAuth)

    // Webhooks (Phase 8)
    v1.GET("/webhooks", handlers.Webhook.List, requireAuth)
    v1.POST("/webhooks", handlers.Webhook.Create, requireAuth)
    v1.PUT("/webhooks/:id", handlers.Webhook.Update, requireAuth)
    v1.DELETE("/webhooks/:id", handlers.Webhook.Delete, requireAuth)

    // Health endpoints (no versioning - stable contract)
    e.GET("/health/live", handlers.Health.Live)
    e.GET("/health/ready", handlers.Health.Ready)
}
```

**File**: `internal/api/middleware/versioning.go`
```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/labstack/echo/v4"
)

// APIVersionMiddleware ensures all API requests use a version prefix
func APIVersionMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            path := c.Request().URL.Path

            // Exempt health endpoints
            if strings.HasPrefix(path, "/health/") {
                return next(c)
            }

            // Require /api/v{N}/ prefix for all other API calls
            if strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/api/v") {
                return echo.NewHTTPError(http.StatusBadRequest,
                    "API version required. Use /api/v1/ prefix. See https://docs.gomailserver.com/api-versioning")
            }

            return next(c)
        }
    }
}

// DeprecationWarningMiddleware adds deprecation headers for sunset APIs
func DeprecationWarningMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            path := c.Request().URL.Path

            // Example: Mark v0 APIs as deprecated
            if strings.HasPrefix(path, "/api/v0/") {
                c.Response().Header().Set("Deprecation", "true")
                c.Response().Header().Set("Sunset", "2025-12-31T23:59:59Z")
                c.Response().Header().Set("Link", `</api/v1/>; rel="successor-version"`)
            }

            return next(c)
        }
    }
}
```

**File**: `docs/API_VERSIONING.md`
```markdown
# API Versioning Policy

## Overview

gomailserver uses URL-based API versioning with the `/api/v{N}/` prefix pattern.

## Current Versions

- **v1**: Current stable API (launched 2025-01-01)
- Health endpoints (`/health/*`): Unversioned, stable contract

## Versioning Strategy

### Version Prefix

All API endpoints use the pattern: `/api/v{major}/endpoint`

Examples:
- `/api/v1/users`
- `/api/v1/mailboxes`
- `/api/v1/webhooks`

### Backward Compatibility

Within a major version (e.g., v1), we maintain strict backward compatibility:
- ✅ Adding new optional fields
- ✅ Adding new endpoints
- ✅ Expanding enum values
- ❌ Removing fields
- ❌ Changing field types
- ❌ Changing required fields
- ❌ Changing validation rules (stricter)

### Deprecation Timeline

1. **Announcement**: Deprecation announced with 6-month notice
2. **Warning Headers**: Deprecated endpoints return `Deprecation` and `Sunset` headers
3. **Sunset Date**: After sunset, endpoints return 410 Gone
4. **Documentation**: Migration guide published with deprecation announcement

### Migration Process

When migrating between versions:

1. Review the [migration guide](API_MIGRATION.md)
2. Test against new version in staging
3. Update client code to use new endpoints
4. Deploy client updates
5. Monitor for deprecation warnings
6. Complete migration before sunset date

## Version History

### v1 (2025-01-01 - Current)

Initial stable API release.

Endpoints:
- Authentication: `/api/v1/auth/*`
- Users: `/api/v1/users`
- Domains: `/api/v1/domains`
- Queue: `/api/v1/queue`
- Mailboxes: `/api/v1/mailboxes`
- Messages: `/api/v1/messages`
- Webhooks: `/api/v1/webhooks`
```

**Acceptance Criteria**:
- [ ] All API endpoints use `/api/v1/` prefix
- [ ] Middleware enforces version requirement
- [ ] Deprecation warning middleware implemented
- [ ] Health endpoints remain unversioned (stable contract)
- [ ] API versioning documentation published
- [ ] Migration guide template created

**Production Readiness**:
- [ ] Versioning policy: 6-month deprecation notice required
- [ ] Sunset headers: `Deprecation`, `Sunset`, `Link` headers on deprecated APIs
- [ ] Client compatibility: Backward-compatible changes only within major version
- [ ] Breaking changes: Require new major version (v2, v3, etc.)
- [ ] Documentation: API changelog maintained with version history

**Given/When/Then Scenarios**:
```
Given API v1 is current
When client calls /api/v1/users
Then request succeeds with 200 OK

Given client uses unversioned API path
When client calls /api/users (no version)
Then request fails with 400 Bad Request
And error message indicates "/api/v1/ prefix required"

Given API v0 is deprecated
When client calls /api/v0/users
Then request succeeds with 200 OK
And response includes Deprecation: true header
And response includes Sunset header with date
And response includes Link header to /api/v1/

Given API v0 sunset date has passed
When client calls /api/v0/users
Then request fails with 410 Gone
And error message indicates "API version sunset, use /api/v1/"
```

---

## 7.1 Webmail Backend

Backend API endpoints for webmail functionality.

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| WM-001 | Mailbox listing API | [ ] | API-002, I-002 | FULL |
| WM-002 | Message fetch API | [ ] | WM-001 | FULL |
| WM-003 | Message send API | [ ] | WM-001, S-002 | FULL |
| WM-004 | Message operations API (move, delete, flag) | [ ] | WM-002 | FULL |
| WM-005 | Attachment download API | [ ] | WM-002 | FULL |
| WM-006 | Attachment upload API | [ ] | WM-003 | FULL |
| WM-007 | Search API | [ ] | WM-001, I-014 | FULL |
| WM-008 | Labels/categories API | [ ] | WM-001 | FULL |

### Task Details

#### WM-001: Mailbox Listing API
**File**: `internal/webmail/mailbox_handler.go`
```go
package webmail

import (
    "github.com/labstack/echo/v4"
    "github.com/btafoya/gomailserver/internal/imap"
)

type MailboxHandler struct {
    imapBackend imap.Backend
}

type MailboxResponse struct {
    Name        string `json:"name"`
    Delimiter   string `json:"delimiter"`
    Attributes  []string `json:"attributes"`
    UnseenCount int    `json:"unseen_count"`
    TotalCount  int    `json:"total_count"`
    UIDValidity uint32 `json:"uid_validity"`
}

// GET /api/webmail/mailboxes
func (h *MailboxHandler) ListMailboxes(c echo.Context) error {
    userID := getUserFromContext(c)

    mailboxes, err := h.imapBackend.ListMailboxes(userID)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to list mailboxes")
    }

    response := make([]MailboxResponse, len(mailboxes))
    for i, mb := range mailboxes {
        status, _ := mb.Status([]imap.StatusItem{
            imap.StatusMessages,
            imap.StatusUnseen,
            imap.StatusUidValidity,
        })
        response[i] = MailboxResponse{
            Name:        mb.Name(),
            Attributes:  mb.Attributes(),
            UnseenCount: int(status.Unseen),
            TotalCount:  int(status.Messages),
            UIDValidity: status.UidValidity,
        }
    }

    return c.JSON(200, response)
}
```

**Acceptance Criteria**:
- [ ] Returns all mailboxes for authenticated user
- [ ] Includes message counts (total, unseen)
- [ ] Returns mailbox attributes (special-use flags)
- [ ] Supports nested mailbox hierarchies

#### WM-002: Message Fetch API
**File**: `internal/webmail/message_handler.go`
```go
type MessageListItem struct {
    UID         uint32    `json:"uid"`
    MessageID   string    `json:"message_id"`
    Subject     string    `json:"subject"`
    From        Address   `json:"from"`
    To          []Address `json:"to"`
    Date        time.Time `json:"date"`
    Size        uint32    `json:"size"`
    Flags       []string  `json:"flags"`
    HasAttach   bool      `json:"has_attachments"`
    ThreadID    string    `json:"thread_id,omitempty"`
    Preview     string    `json:"preview"`
    Labels      []string  `json:"labels,omitempty"`
}

type MessageDetail struct {
    MessageListItem
    Cc          []Address     `json:"cc,omitempty"`
    Bcc         []Address     `json:"bcc,omitempty"`
    ReplyTo     []Address     `json:"reply_to,omitempty"`
    InReplyTo   string        `json:"in_reply_to,omitempty"`
    References  []string      `json:"references,omitempty"`
    HTMLBody    string        `json:"html_body,omitempty"`
    TextBody    string        `json:"text_body"`
    Attachments []Attachment  `json:"attachments,omitempty"`
    Headers     []Header      `json:"headers,omitempty"`
}

// GET /api/webmail/mailboxes/:mailbox/messages
func (h *MessageHandler) ListMessages(c echo.Context) error {
    mailbox := c.Param("mailbox")
    page, _ := strconv.Atoi(c.QueryParam("page"))
    pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
    sort := c.QueryParam("sort") // date, from, subject
    order := c.QueryParam("order") // asc, desc

    if pageSize == 0 || pageSize > 100 {
        pageSize = 50
    }

    messages, total, err := h.fetchMessages(mailbox, page, pageSize, sort, order)
    if err != nil {
        return echo.NewHTTPError(500, err.Error())
    }

    return c.JSON(200, map[string]interface{}{
        "messages": messages,
        "total":    total,
        "page":     page,
        "page_size": pageSize,
    })
}

// GET /api/webmail/mailboxes/:mailbox/messages/:uid
func (h *MessageHandler) GetMessage(c echo.Context) error {
    mailbox := c.Param("mailbox")
    uid, _ := strconv.ParseUint(c.Param("uid"), 10, 32)

    message, err := h.fetchMessageDetail(mailbox, uint32(uid))
    if err != nil {
        return echo.NewHTTPError(404, "Message not found")
    }

    // Mark as seen
    if c.QueryParam("mark_read") != "false" {
        h.addFlag(mailbox, uint32(uid), "\\Seen")
    }

    return c.JSON(200, message)
}
```

**Acceptance Criteria**:
- [ ] Paginated message listing
- [ ] Sortable by date, from, subject
- [ ] Message preview (first 200 chars of body)
- [ ] Full message detail with HTML/text bodies
- [ ] Conversation threading via References header

#### WM-003: Message Send API
**File**: `internal/webmail/compose_handler.go`
```go
type ComposeRequest struct {
    To          []string          `json:"to" validate:"required,dive,email"`
    Cc          []string          `json:"cc,omitempty" validate:"omitempty,dive,email"`
    Bcc         []string          `json:"bcc,omitempty" validate:"omitempty,dive,email"`
    Subject     string            `json:"subject" validate:"required,max=998"`
    HTMLBody    string            `json:"html_body,omitempty"`
    TextBody    string            `json:"text_body"`
    ReplyTo     string            `json:"reply_to,omitempty"`
    InReplyTo   string            `json:"in_reply_to,omitempty"`
    References  []string          `json:"references,omitempty"`
    Attachments []AttachmentInput `json:"attachments,omitempty"`
    SaveDraft   bool              `json:"save_draft"`
    DraftUID    uint32            `json:"draft_uid,omitempty"` // for updating drafts
    Labels      []string          `json:"labels,omitempty"`
}

type AttachmentInput struct {
    ID          string `json:"id"` // From upload endpoint
    Filename    string `json:"filename"`
    ContentType string `json:"content_type"`
    Inline      bool   `json:"inline"`
    ContentID   string `json:"content_id,omitempty"` // For inline images
}

// POST /api/webmail/send
func (h *ComposeHandler) SendMessage(c echo.Context) error {
    var req ComposeRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    if err := c.Validate(req); err != nil {
        return echo.NewHTTPError(400, err.Error())
    }

    user := getUserFromContext(c)

    // Build message
    msg, err := h.buildMessage(user, &req)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to build message")
    }

    if req.SaveDraft {
        // Save to Drafts folder
        uid, err := h.saveDraft(user, msg, req.DraftUID)
        if err != nil {
            return echo.NewHTTPError(500, "Failed to save draft")
        }
        return c.JSON(200, map[string]interface{}{
            "status": "draft_saved",
            "uid":    uid,
        })
    }

    // Queue for sending
    messageID, err := h.smtpService.QueueMessage(user, msg)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to queue message")
    }

    // Save to Sent folder
    h.saveToSent(user, msg)

    // Delete draft if exists
    if req.DraftUID > 0 {
        h.deleteDraft(user, req.DraftUID)
    }

    return c.JSON(200, map[string]interface{}{
        "status":     "sent",
        "message_id": messageID,
    })
}
```

**Acceptance Criteria**:
- [ ] Compose with To, Cc, Bcc
- [ ] HTML and plain text bodies
- [ ] Attachment support
- [ ] Reply/forward with proper headers
- [ ] Save to Drafts
- [ ] Auto-save to Sent folder

#### WM-004: Message Operations API
```go
type MessageActionRequest struct {
    UIDs        []uint32 `json:"uids" validate:"required,min=1"`
    TargetBox   string   `json:"target_mailbox,omitempty"`
    AddFlags    []string `json:"add_flags,omitempty"`
    RemoveFlags []string `json:"remove_flags,omitempty"`
    AddLabels   []string `json:"add_labels,omitempty"`
    RemoveLabels []string `json:"remove_labels,omitempty"`
}

// POST /api/webmail/mailboxes/:mailbox/messages/move
func (h *MessageHandler) MoveMessages(c echo.Context) error {
    var req MessageActionRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    sourceBox := c.Param("mailbox")

    if req.TargetBox == "" {
        return echo.NewHTTPError(400, "Target mailbox required")
    }

    err := h.imapBackend.MoveMessages(sourceBox, req.TargetBox, req.UIDs)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to move messages")
    }

    return c.JSON(200, map[string]string{"status": "moved"})
}

// POST /api/webmail/mailboxes/:mailbox/messages/flag
func (h *MessageHandler) UpdateFlags(c echo.Context) error {
    var req MessageActionRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    mailbox := c.Param("mailbox")

    if len(req.AddFlags) > 0 {
        h.imapBackend.AddFlags(mailbox, req.UIDs, req.AddFlags)
    }
    if len(req.RemoveFlags) > 0 {
        h.imapBackend.RemoveFlags(mailbox, req.UIDs, req.RemoveFlags)
    }

    return c.JSON(200, map[string]string{"status": "updated"})
}

// DELETE /api/webmail/mailboxes/:mailbox/messages
func (h *MessageHandler) DeleteMessages(c echo.Context) error {
    var req MessageActionRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    mailbox := c.Param("mailbox")

    // Move to Trash or permanently delete if already in Trash
    if mailbox == "Trash" {
        err := h.imapBackend.ExpungeMessages(mailbox, req.UIDs)
        if err != nil {
            return echo.NewHTTPError(500, "Failed to delete messages")
        }
    } else {
        err := h.imapBackend.MoveMessages(mailbox, "Trash", req.UIDs)
        if err != nil {
            return echo.NewHTTPError(500, "Failed to move to trash")
        }
    }

    return c.JSON(200, map[string]string{"status": "deleted"})
}
```

**Acceptance Criteria**:
- [ ] Move messages between mailboxes
- [ ] Bulk flag operations (read, unread, flagged)
- [ ] Move to Trash / permanent delete
- [ ] Label management

#### WM-005: Attachment Download API
```go
// GET /api/webmail/mailboxes/:mailbox/messages/:uid/attachments/:part
func (h *MessageHandler) DownloadAttachment(c echo.Context) error {
    mailbox := c.Param("mailbox")
    uid, _ := strconv.ParseUint(c.Param("uid"), 10, 32)
    partID := c.Param("part")

    attachment, err := h.fetchAttachment(mailbox, uint32(uid), partID)
    if err != nil {
        return echo.NewHTTPError(404, "Attachment not found")
    }

    c.Response().Header().Set("Content-Disposition",
        fmt.Sprintf(`attachment; filename="%s"`, attachment.Filename))
    c.Response().Header().Set("Content-Type", attachment.ContentType)

    return c.Blob(200, attachment.ContentType, attachment.Data)
}
```

#### WM-006: Attachment Upload API
```go
type UploadedAttachment struct {
    ID          string `json:"id"`
    Filename    string `json:"filename"`
    Size        int64  `json:"size"`
    ContentType string `json:"content_type"`
}

// POST /api/webmail/attachments
func (h *ComposeHandler) UploadAttachment(c echo.Context) error {
    file, err := c.FormFile("file")
    if err != nil {
        return echo.NewHTTPError(400, "No file uploaded")
    }

    // Size limit: 25MB
    if file.Size > 25*1024*1024 {
        return echo.NewHTTPError(400, "File too large (max 25MB)")
    }

    src, err := file.Open()
    if err != nil {
        return echo.NewHTTPError(500, "Failed to read file")
    }
    defer src.Close()

    // Store temporarily with UUID
    attachmentID := uuid.NewString()
    contentType := file.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    err = h.storeTempAttachment(attachmentID, src, file.Size)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to store attachment")
    }

    return c.JSON(200, UploadedAttachment{
        ID:          attachmentID,
        Filename:    file.Filename,
        Size:        file.Size,
        ContentType: contentType,
    })
}
```

**Acceptance Criteria**:
- [ ] Upload with progress indication support
- [ ] 25MB size limit
- [ ] Temporary storage with cleanup
- [ ] Content-Type detection

#### WM-007: Search API
```go
type SearchRequest struct {
    Query      string   `json:"query"`
    Mailbox    string   `json:"mailbox,omitempty"` // Empty = all mailboxes
    From       string   `json:"from,omitempty"`
    To         string   `json:"to,omitempty"`
    Subject    string   `json:"subject,omitempty"`
    Body       string   `json:"body,omitempty"`
    DateFrom   string   `json:"date_from,omitempty"`
    DateTo     string   `json:"date_to,omitempty"`
    HasAttach  *bool    `json:"has_attachment,omitempty"`
    IsUnread   *bool    `json:"is_unread,omitempty"`
    IsFlagged  *bool    `json:"is_flagged,omitempty"`
    Labels     []string `json:"labels,omitempty"`
}

// POST /api/webmail/search
func (h *MessageHandler) Search(c echo.Context) error {
    var req SearchRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    // Build IMAP search criteria
    criteria := h.buildSearchCriteria(&req)

    results, err := h.imapBackend.Search(req.Mailbox, criteria)
    if err != nil {
        return echo.NewHTTPError(500, "Search failed")
    }

    return c.JSON(200, map[string]interface{}{
        "results": results,
        "total":   len(results),
    })
}
```

**Acceptance Criteria**:
- [ ] Full-text search
- [ ] Field-specific search (from, to, subject)
- [ ] Date range filtering
- [ ] Flag/label filtering
- [ ] Search across all mailboxes

#### WM-008: Labels/Categories API
```go
type Label struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Color     string `json:"color"`
    UserID    int64  `json:"user_id"`
}

// GET /api/webmail/labels
func (h *LabelHandler) ListLabels(c echo.Context) error {
    userID := getUserFromContext(c)

    labels, err := h.labelRepo.GetByUser(userID)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to fetch labels")
    }

    return c.JSON(200, labels)
}

// POST /api/webmail/labels
func (h *LabelHandler) CreateLabel(c echo.Context) error {
    var label Label
    if err := c.Bind(&label); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    userID := getUserFromContext(c)
    label.UserID = userID

    err := h.labelRepo.Create(&label)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to create label")
    }

    return c.JSON(201, label)
}
```

**Acceptance Criteria**:
- [ ] CRUD for custom labels
- [ ] Color assignment
- [ ] Per-user labels
- [ ] Apply/remove labels from messages

---

## 7.2 Webmail Frontend

Vue.js 3 + Nuxt webmail client.

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| WF-001 | Set up Vue.js 3 + Nuxt for webmail | [ ] | - | FULL |
| WF-002 | Authentication and session | [ ] | WF-001, USP-001 | FULL |
| WF-003 | Mailbox sidebar | [ ] | WF-002, WM-001 | FULL |
| WF-004 | Message list view | [ ] | WF-003 | FULL |
| WF-005 | Conversation/thread view | [ ] | WF-004 | FULL |
| WF-006 | Message detail view | [ ] | WF-004 | FULL |
| WF-007 | Rich text composer (TipTap) | [ ] | WF-002 | FULL |
| WF-008 | Plain text composer | [ ] | WF-007 | FULL |
| WF-009 | Attachment handling (drag-drop) | [ ] | WF-007, WM-006 | FULL |
| WF-010 | Inline images | [ ] | WF-009 | FULL |
| WF-011 | Gmail-like categories UI | [ ] | WF-003, WM-008 | FULL |
| WF-012 | Search interface | [ ] | WF-002, WM-007 | FULL |
| WF-013 | Keyboard shortcuts | [ ] | WF-002 | FULL |
| WF-014 | Dark mode | [ ] | WF-001 | FULL |
| WF-015 | Mobile responsive design | [ ] | WF-001 | FULL |
| WF-016 | PWA offline capability | [ ] | WF-001 | FULL |
| WF-017 | Auto-save drafts | [ ] | WF-007 | FULL |
| WF-018 | Message templates | [ ] | WF-007 | FULL |
| WF-019 | Spam reporting button | [ ] | WF-004, AS-005 | FULL |

### Task Details

#### WF-001: Vue.js 3 + Nuxt Setup
**Directory**: `web/webmail/`
```bash
# Project setup
pnpm create nuxt@latest webmail
cd webmail
pnpm add @tiptap/vue-3 @tiptap/starter-kit @tiptap/extension-image
pnpm add @vueuse/core pinia
pnpm add -D tailwindcss postcss autoprefixer
```

**File**: `web/webmail/nuxt.config.ts`
```typescript
export default defineNuxtConfig({
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@vueuse/nuxt',
  ],
  app: {
    head: {
      title: 'Webmail - gomailserver',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' }
      ]
    }
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE || '/api'
    }
  },
  pwa: {
    manifest: {
      name: 'gomailserver Webmail',
      short_name: 'Webmail',
      theme_color: '#1a73e8'
    }
  }
})
```

#### WF-003: Mailbox Sidebar
**File**: `web/webmail/components/MailboxSidebar.vue`
```vue
<template>
  <aside class="w-64 bg-gray-50 dark:bg-gray-900 h-screen overflow-y-auto">
    <div class="p-4">
      <button @click="openCompose"
        class="w-full bg-blue-500 text-white rounded-full py-3 px-6 flex items-center gap-2 hover:shadow-lg transition-shadow">
        <PencilIcon class="w-5 h-5" />
        <span>Compose</span>
      </button>
    </div>

    <nav class="px-2">
      <ul class="space-y-1">
        <li v-for="mailbox in mailboxes" :key="mailbox.name">
          <NuxtLink
            :to="`/mailbox/${encodeURIComponent(mailbox.name)}`"
            :class="[
              'flex items-center gap-3 px-3 py-2 rounded-full transition-colors',
              isActive(mailbox.name)
                ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300'
                : 'hover:bg-gray-200 dark:hover:bg-gray-800'
            ]">
            <component :is="getIcon(mailbox)" class="w-5 h-5" />
            <span class="flex-1">{{ getDisplayName(mailbox.name) }}</span>
            <span v-if="mailbox.unseen_count > 0"
              class="text-sm font-medium">
              {{ mailbox.unseen_count }}
            </span>
          </NuxtLink>

          <!-- Nested mailboxes -->
          <ul v-if="mailbox.children?.length" class="ml-6 mt-1 space-y-1">
            <li v-for="child in mailbox.children" :key="child.name">
              <!-- Recursive component -->
            </li>
          </ul>
        </li>
      </ul>
    </nav>

    <!-- Labels section -->
    <div class="px-4 mt-6">
      <h3 class="text-sm font-medium text-gray-500 mb-2">Labels</h3>
      <ul class="space-y-1">
        <li v-for="label in labels" :key="label.id">
          <NuxtLink
            :to="`/label/${label.id}`"
            class="flex items-center gap-2 px-3 py-1 rounded hover:bg-gray-200">
            <span class="w-3 h-3 rounded-full" :style="{ backgroundColor: label.color }"></span>
            <span>{{ label.name }}</span>
          </NuxtLink>
        </li>
      </ul>
      <button @click="createLabel" class="text-sm text-blue-500 mt-2">
        + Create new label
      </button>
    </div>

    <!-- Quota display -->
    <div class="px-4 mt-6 text-sm text-gray-500">
      <div class="flex justify-between mb-1">
        <span>Storage</span>
        <span>{{ formatBytes(quota.used) }} / {{ formatBytes(quota.limit) }}</span>
      </div>
      <div class="h-1.5 bg-gray-200 rounded-full">
        <div class="h-full bg-blue-500 rounded-full"
          :style="{ width: `${(quota.used / quota.limit) * 100}%` }"></div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { InboxIcon, PaperAirplaneIcon, DocumentIcon, TrashIcon,
         ExclamationIcon, ArchiveBoxIcon, StarIcon, PencilIcon } from '@heroicons/vue/24/outline'

const mailboxStore = useMailboxStore()
const { mailboxes, labels, quota } = storeToRefs(mailboxStore)

const route = useRoute()

const isActive = (name: string) => route.params.mailbox === name

const iconMap: Record<string, any> = {
  'INBOX': InboxIcon,
  'Sent': PaperAirplaneIcon,
  'Drafts': DocumentIcon,
  'Trash': TrashIcon,
  'Spam': ExclamationIcon,
  'Archive': ArchiveBoxIcon,
  'Starred': StarIcon,
}

const getIcon = (mailbox: Mailbox) => iconMap[mailbox.name] || FolderIcon
const getDisplayName = (name: string) => name === 'INBOX' ? 'Inbox' : name

const openCompose = () => {
  navigateTo('/compose')
}
</script>
```

#### WF-004: Message List View
**File**: `web/webmail/components/MessageList.vue`
```vue
<template>
  <div class="flex-1 flex flex-col">
    <!-- Toolbar -->
    <div class="flex items-center gap-2 p-2 border-b">
      <Checkbox v-model="selectAll" @change="toggleSelectAll" />
      <button @click="refresh" class="p-2 hover:bg-gray-100 rounded">
        <ArrowPathIcon class="w-5 h-5" />
      </button>
      <button v-if="selectedCount > 0" @click="deleteSelected" class="p-2 hover:bg-gray-100 rounded">
        <TrashIcon class="w-5 h-5" />
      </button>
      <button v-if="selectedCount > 0" @click="markAsRead" class="p-2 hover:bg-gray-100 rounded">
        <EnvelopeOpenIcon class="w-5 h-5" />
      </button>
      <button v-if="selectedCount > 0" @click="markAsUnread" class="p-2 hover:bg-gray-100 rounded">
        <EnvelopeIcon class="w-5 h-5" />
      </button>

      <div class="flex-1"></div>

      <span class="text-sm text-gray-500">
        {{ pagination.start }}-{{ pagination.end }} of {{ pagination.total }}
      </span>
      <button @click="prevPage" :disabled="!pagination.hasPrev"
        class="p-2 hover:bg-gray-100 rounded disabled:opacity-50">
        <ChevronLeftIcon class="w-5 h-5" />
      </button>
      <button @click="nextPage" :disabled="!pagination.hasNext"
        class="p-2 hover:bg-gray-100 rounded disabled:opacity-50">
        <ChevronRightIcon class="w-5 h-5" />
      </button>
    </div>

    <!-- Message list -->
    <div class="flex-1 overflow-y-auto">
      <div v-if="loading" class="flex items-center justify-center h-full">
        <Spinner />
      </div>

      <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center h-full text-gray-500">
        <InboxIcon class="w-16 h-16 mb-4" />
        <p>No messages</p>
      </div>

      <ul v-else>
        <li v-for="msg in messages" :key="msg.uid"
          @click="openMessage(msg)"
          :class="[
            'flex items-center gap-3 px-4 py-2 cursor-pointer border-b hover:bg-gray-50',
            msg.flags.includes('\\Seen') ? '' : 'font-semibold bg-blue-50',
            selected.has(msg.uid) ? 'bg-blue-100' : ''
          ]">
          <Checkbox :modelValue="selected.has(msg.uid)" @update:modelValue="toggleSelect(msg.uid)"
            @click.stop />

          <button @click.stop="toggleStar(msg)" class="p-1">
            <StarIcon :class="[
              'w-5 h-5',
              msg.flags.includes('\\Flagged') ? 'text-yellow-400 fill-yellow-400' : 'text-gray-400'
            ]" />
          </button>

          <div class="w-48 truncate">
            {{ msg.from.name || msg.from.address }}
          </div>

          <div class="flex-1 flex items-center gap-2 truncate">
            <span>{{ msg.subject }}</span>
            <span class="text-gray-500">- {{ msg.preview }}</span>
          </div>

          <div v-if="msg.has_attachments" class="text-gray-400">
            <PaperClipIcon class="w-4 h-4" />
          </div>

          <div class="text-sm text-gray-500 whitespace-nowrap">
            {{ formatDate(msg.date) }}
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  mailbox: string
}>()

const messageStore = useMessageStore()
const { messages, loading, pagination } = storeToRefs(messageStore)
const selected = ref(new Set<number>())
const selectAll = ref(false)

watch(() => props.mailbox, () => {
  messageStore.fetchMessages(props.mailbox)
  selected.value.clear()
}, { immediate: true })

const openMessage = (msg: Message) => {
  navigateTo(`/mailbox/${props.mailbox}/${msg.uid}`)
}

const toggleSelect = (uid: number) => {
  if (selected.value.has(uid)) {
    selected.value.delete(uid)
  } else {
    selected.value.add(uid)
  }
}

const deleteSelected = async () => {
  await messageStore.deleteMessages(props.mailbox, [...selected.value])
  selected.value.clear()
}

const formatDate = (date: string) => {
  const d = new Date(date)
  const today = new Date()
  if (d.toDateString() === today.toDateString()) {
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
  return d.toLocaleDateString([], { month: 'short', day: 'numeric' })
}
</script>
```

#### WF-007: Rich Text Composer (TipTap)
**File**: `web/webmail/components/Composer.vue`
```vue
<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] flex flex-col">
      <!-- Header -->
      <div class="flex items-center justify-between p-4 border-b">
        <h2 class="text-lg font-medium">{{ draft?.uid ? 'Edit Draft' : 'New Message' }}</h2>
        <button @click="close" class="p-2 hover:bg-gray-100 rounded">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>

      <!-- Recipients -->
      <div class="p-4 space-y-2 border-b">
        <div class="flex items-center gap-2">
          <label class="w-16 text-sm text-gray-500">To</label>
          <RecipientInput v-model="form.to" />
        </div>
        <div v-if="showCc" class="flex items-center gap-2">
          <label class="w-16 text-sm text-gray-500">Cc</label>
          <RecipientInput v-model="form.cc" />
        </div>
        <div v-if="showBcc" class="flex items-center gap-2">
          <label class="w-16 text-sm text-gray-500">Bcc</label>
          <RecipientInput v-model="form.bcc" />
        </div>
        <div class="flex items-center gap-2">
          <label class="w-16 text-sm text-gray-500">Subject</label>
          <input v-model="form.subject" type="text"
            class="flex-1 border-0 focus:ring-0 p-1" placeholder="Subject">
        </div>
        <div class="flex gap-2 text-sm text-blue-500">
          <button v-if="!showCc" @click="showCc = true">Cc</button>
          <button v-if="!showBcc" @click="showBcc = true">Bcc</button>
        </div>
      </div>

      <!-- Editor -->
      <div class="flex-1 overflow-y-auto p-4">
        <div class="border rounded-lg">
          <!-- Toolbar -->
          <div class="flex items-center gap-1 p-2 border-b bg-gray-50">
            <button @click="editor?.chain().focus().toggleBold().run()"
              :class="['p-1.5 rounded', editor?.isActive('bold') ? 'bg-gray-200' : 'hover:bg-gray-200']">
              <BoldIcon class="w-4 h-4" />
            </button>
            <button @click="editor?.chain().focus().toggleItalic().run()"
              :class="['p-1.5 rounded', editor?.isActive('italic') ? 'bg-gray-200' : 'hover:bg-gray-200']">
              <ItalicIcon class="w-4 h-4" />
            </button>
            <button @click="editor?.chain().focus().toggleUnderline().run()"
              :class="['p-1.5 rounded', editor?.isActive('underline') ? 'bg-gray-200' : 'hover:bg-gray-200']">
              <UnderlineIcon class="w-4 h-4" />
            </button>
            <div class="w-px h-6 bg-gray-300 mx-1"></div>
            <button @click="editor?.chain().focus().toggleBulletList().run()">
              <ListBulletIcon class="w-4 h-4" />
            </button>
            <button @click="editor?.chain().focus().toggleOrderedList().run()">
              <ListNumberedIcon class="w-4 h-4" />
            </button>
            <div class="w-px h-6 bg-gray-300 mx-1"></div>
            <button @click="insertLink">
              <LinkIcon class="w-4 h-4" />
            </button>
            <button @click="insertImage">
              <PhotoIcon class="w-4 h-4" />
            </button>
            <div class="flex-1"></div>
            <button @click="togglePlainText" class="text-sm text-gray-500">
              {{ isPlainText ? 'Rich text' : 'Plain text' }}
            </button>
          </div>

          <!-- Content -->
          <EditorContent v-if="!isPlainText" :editor="editor" class="prose max-w-none p-4 min-h-[200px]" />
          <textarea v-else v-model="form.textBody" class="w-full p-4 min-h-[200px] font-mono text-sm border-0 focus:ring-0"></textarea>
        </div>
      </div>

      <!-- Attachments -->
      <div v-if="attachments.length > 0" class="px-4 pb-2">
        <div class="flex flex-wrap gap-2">
          <div v-for="(att, i) in attachments" :key="att.id"
            class="flex items-center gap-2 bg-gray-100 rounded px-2 py-1 text-sm">
            <PaperClipIcon class="w-4 h-4" />
            <span>{{ att.filename }}</span>
            <span class="text-gray-500">({{ formatBytes(att.size) }})</span>
            <button @click="removeAttachment(i)" class="text-gray-400 hover:text-red-500">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-between p-4 border-t">
        <div class="flex items-center gap-2">
          <button @click="send" :disabled="sending"
            class="bg-blue-500 text-white px-6 py-2 rounded-full hover:bg-blue-600 disabled:opacity-50">
            {{ sending ? 'Sending...' : 'Send' }}
          </button>
          <button @click="saveDraft" class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded">
            Save Draft
          </button>
        </div>
        <div class="flex items-center gap-2">
          <label class="cursor-pointer p-2 hover:bg-gray-100 rounded">
            <input type="file" multiple @change="handleFileUpload" class="hidden" />
            <PaperClipIcon class="w-5 h-5" />
          </label>
          <button @click="togglePGP" v-if="hasPGPKey"
            :class="['p-2 rounded', pgpEnabled ? 'bg-green-100 text-green-600' : 'hover:bg-gray-100']">
            <LockClosedIcon class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import Link from '@tiptap/extension-link'
import Image from '@tiptap/extension-image'

const props = defineProps<{
  replyTo?: Message
  forward?: Message
  draft?: Message
}>()

const emit = defineEmits(['close', 'sent'])

const form = reactive({
  to: [] as string[],
  cc: [] as string[],
  bcc: [] as string[],
  subject: '',
  textBody: '',
})

const showCc = ref(false)
const showBcc = ref(false)
const isPlainText = ref(false)
const pgpEnabled = ref(false)
const sending = ref(false)
const attachments = ref<UploadedAttachment[]>([])

const editor = useEditor({
  content: '',
  extensions: [
    StarterKit,
    Underline,
    Link.configure({ openOnClick: false }),
    Image.configure({ inline: true }),
  ],
  autofocus: true,
})

// Auto-save draft every 30 seconds
const { pause, resume } = useIntervalFn(() => {
  saveDraft(true) // silent save
}, 30000)

onMounted(() => {
  if (props.replyTo) {
    setupReply()
  } else if (props.forward) {
    setupForward()
  } else if (props.draft) {
    loadDraft()
  }
})

const send = async () => {
  sending.value = true
  try {
    await $fetch('/api/webmail/send', {
      method: 'POST',
      body: {
        ...form,
        html_body: editor.value?.getHTML(),
        text_body: form.textBody || editor.value?.getText(),
        attachments: attachments.value,
        pgp_encrypt: pgpEnabled.value,
      }
    })
    emit('sent')
    close()
  } catch (err) {
    // Handle error
  } finally {
    sending.value = false
  }
}

const handleFileUpload = async (e: Event) => {
  const files = (e.target as HTMLInputElement).files
  if (!files) return

  for (const file of files) {
    const formData = new FormData()
    formData.append('file', file)

    const result = await $fetch<UploadedAttachment>('/api/webmail/attachments', {
      method: 'POST',
      body: formData,
    })
    attachments.value.push(result)
  }
}

// ... remaining methods
</script>
```

#### WF-013: Keyboard Shortcuts
**File**: `web/webmail/composables/useKeyboardShortcuts.ts`
```typescript
export const useKeyboardShortcuts = () => {
  const router = useRouter()
  const messageStore = useMessageStore()

  const shortcuts: Record<string, () => void> = {
    'c': () => navigateTo('/compose'),
    'r': () => messageStore.replyToSelected(),
    'a': () => messageStore.replyAllToSelected(),
    'f': () => messageStore.forwardSelected(),
    'e': () => messageStore.archiveSelected(),
    '#': () => messageStore.deleteSelected(),
    's': () => messageStore.toggleStarSelected(),
    'u': () => messageStore.markUnread(),
    '/': () => document.querySelector<HTMLInputElement>('.search-input')?.focus(),
    'j': () => messageStore.selectNext(),
    'k': () => messageStore.selectPrevious(),
    'o': () => messageStore.openSelected(),
    'Enter': () => messageStore.openSelected(),
    'Escape': () => messageStore.clearSelection(),
    'g i': () => navigateTo('/mailbox/INBOX'),
    'g s': () => navigateTo('/mailbox/Sent'),
    'g d': () => navigateTo('/mailbox/Drafts'),
    'g t': () => navigateTo('/mailbox/Trash'),
    '?': () => showShortcutsHelp(),
  }

  let keySequence = ''
  let keyTimeout: NodeJS.Timeout

  const handleKeydown = (e: KeyboardEvent) => {
    // Skip if in input/textarea
    if (['INPUT', 'TEXTAREA'].includes((e.target as HTMLElement).tagName)) {
      return
    }

    const key = e.key

    // Handle key sequences (g i, g s, etc.)
    clearTimeout(keyTimeout)
    keySequence += key + ' '
    keyTimeout = setTimeout(() => keySequence = '', 500)

    const action = shortcuts[keySequence.trim()] || shortcuts[key]
    if (action) {
      e.preventDefault()
      action()
      keySequence = ''
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
  })
}
```

**Acceptance Criteria**:
- [ ] Gmail-compatible shortcuts
- [ ] Key sequence support (g i, g s)
- [ ] Help modal (?)
- [ ] Disabled in text inputs

---

## 7.3 Contact Integration

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| CI-001 | Contact picker in composer | [ ] | WF-007, CAR-003 | FULL |
| CI-002 | Contact autocomplete | [ ] | CI-001 | FULL |
| CI-003 | Contact management view | [ ] | WF-002, CAR-003 | FULL |

#### CI-002: Contact Autocomplete
**File**: `web/webmail/components/RecipientInput.vue`
```vue
<template>
  <div class="flex-1 flex flex-wrap items-center gap-1 border rounded-lg p-1 focus-within:ring-2 focus-within:ring-blue-500">
    <div v-for="(addr, i) in modelValue" :key="i"
      class="flex items-center gap-1 bg-gray-100 rounded-full px-2 py-0.5 text-sm">
      <span>{{ formatAddress(addr) }}</span>
      <button @click="remove(i)" class="text-gray-400 hover:text-red-500">
        <XMarkIcon class="w-3 h-3" />
      </button>
    </div>

    <div class="relative flex-1 min-w-[120px]">
      <input
        ref="input"
        v-model="query"
        @input="search"
        @keydown="handleKeydown"
        @blur="handleBlur"
        type="text"
        class="w-full border-0 focus:ring-0 p-1 text-sm"
        placeholder="Type a name or email..."
      />

      <div v-if="suggestions.length > 0"
        class="absolute top-full left-0 w-64 bg-white shadow-lg rounded-lg border mt-1 z-50">
        <ul>
          <li v-for="(contact, i) in suggestions" :key="contact.email"
            @mousedown.prevent="selectSuggestion(contact)"
            :class="[
              'px-3 py-2 cursor-pointer',
              i === selectedIndex ? 'bg-blue-50' : 'hover:bg-gray-50'
            ]">
            <div class="font-medium">{{ contact.name }}</div>
            <div class="text-sm text-gray-500">{{ contact.email }}</div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: string[]
}>()

const emit = defineEmits(['update:modelValue'])

const query = ref('')
const suggestions = ref<Contact[]>([])
const selectedIndex = ref(0)
const input = ref<HTMLInputElement>()

const search = useDebounceFn(async () => {
  if (query.value.length < 2) {
    suggestions.value = []
    return
  }

  const results = await $fetch<Contact[]>('/api/contacts/search', {
    params: { q: query.value }
  })
  suggestions.value = results
  selectedIndex.value = 0
}, 200)

const handleKeydown = (e: KeyboardEvent) => {
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      selectedIndex.value = Math.min(selectedIndex.value + 1, suggestions.value.length - 1)
      break
    case 'ArrowUp':
      e.preventDefault()
      selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
      break
    case 'Enter':
    case 'Tab':
      if (suggestions.value.length > 0) {
        e.preventDefault()
        selectSuggestion(suggestions.value[selectedIndex.value])
      } else if (query.value.includes('@')) {
        e.preventDefault()
        addAddress(query.value)
      }
      break
    case 'Backspace':
      if (query.value === '' && props.modelValue.length > 0) {
        remove(props.modelValue.length - 1)
      }
      break
  }
}

const selectSuggestion = (contact: Contact) => {
  addAddress(contact.name ? `${contact.name} <${contact.email}>` : contact.email)
}

const addAddress = (addr: string) => {
  emit('update:modelValue', [...props.modelValue, addr])
  query.value = ''
  suggestions.value = []
}

const remove = (index: number) => {
  const updated = [...props.modelValue]
  updated.splice(index, 1)
  emit('update:modelValue', updated)
}
</script>
```

---

## 7.4 Calendar Integration

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| CLI-001 | Calendar widget/view | [ ] | WF-002, CD-003 | FULL |
| CLI-002 | Event creation from webmail | [ ] | CLI-001 | FULL |
| CLI-003 | Meeting invitation handling | [ ] | WF-004, CD-008 | FULL |

#### CLI-003: Meeting Invitation Handling
**File**: `web/webmail/components/MeetingInvitation.vue`
```vue
<template>
  <div class="bg-blue-50 rounded-lg p-4 my-4">
    <div class="flex items-start gap-4">
      <CalendarIcon class="w-8 h-8 text-blue-500 flex-shrink-0" />
      <div class="flex-1">
        <h3 class="font-medium">{{ event.summary }}</h3>
        <p class="text-sm text-gray-600 mt-1">
          {{ formatEventTime(event.start, event.end) }}
        </p>
        <p v-if="event.location" class="text-sm text-gray-600">
          <MapPinIcon class="w-4 h-4 inline" /> {{ event.location }}
        </p>
        <p class="text-sm text-gray-600 mt-2">
          Organizer: {{ event.organizer.name || event.organizer.email }}
        </p>

        <div class="flex gap-2 mt-4">
          <button @click="respond('ACCEPTED')"
            :class="['px-4 py-1.5 rounded-full text-sm',
              response === 'ACCEPTED' ? 'bg-green-500 text-white' : 'bg-white border hover:bg-gray-50']">
            Yes
          </button>
          <button @click="respond('TENTATIVE')"
            :class="['px-4 py-1.5 rounded-full text-sm',
              response === 'TENTATIVE' ? 'bg-yellow-500 text-white' : 'bg-white border hover:bg-gray-50']">
            Maybe
          </button>
          <button @click="respond('DECLINED')"
            :class="['px-4 py-1.5 rounded-full text-sm',
              response === 'DECLINED' ? 'bg-red-500 text-white' : 'bg-white border hover:bg-gray-50']">
            No
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  event: CalendarEvent
  messageUid: number
}>()

const response = ref(props.event.attendeeStatus)

const respond = async (status: 'ACCEPTED' | 'TENTATIVE' | 'DECLINED') => {
  await $fetch('/api/webmail/calendar/respond', {
    method: 'POST',
    body: {
      message_uid: props.messageUid,
      status,
    }
  })
  response.value = status
}
</script>
```

---

## 7.5 PGP in Webmail

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| WPG-001 | Integrate OpenPGP.js | [ ] | WF-001 | FULL |
| WPG-002 | Encrypt message composition | [ ] | WPG-001, WF-007 | FULL |
| WPG-003 | Decrypt message viewing | [ ] | WPG-001, WF-006 | FULL |
| WPG-004 | Sign messages | [ ] | WPG-001, WF-007 | FULL |
| WPG-005 | Verify signatures | [ ] | WPG-001, WF-006 | FULL |

#### WPG-001: OpenPGP.js Integration
**File**: `web/webmail/composables/usePGP.ts`
```typescript
import * as openpgp from 'openpgp'

export const usePGP = () => {
  const userStore = useUserStore()
  const privateKey = ref<openpgp.PrivateKey | null>(null)
  const publicKeys = ref<Map<string, openpgp.PublicKey>>(new Map())

  // Load user's private key (decrypted with passphrase)
  const loadPrivateKey = async (armoredKey: string, passphrase: string) => {
    try {
      privateKey.value = await openpgp.decryptKey({
        privateKey: await openpgp.readPrivateKey({ armoredKey }),
        passphrase,
      })
      return true
    } catch (err) {
      console.error('Failed to decrypt private key:', err)
      return false
    }
  }

  // Get public key for recipient
  const getPublicKey = async (email: string): Promise<openpgp.PublicKey | null> => {
    if (publicKeys.value.has(email)) {
      return publicKeys.value.get(email)!
    }

    try {
      const { armored_key } = await $fetch<{ armored_key: string }>(`/api/pgp/keys/${email}`)
      const key = await openpgp.readKey({ armoredKey: armored_key })
      publicKeys.value.set(email, key)
      return key
    } catch {
      return null
    }
  }

  // Encrypt message for recipients
  const encrypt = async (text: string, recipients: string[]): Promise<string | null> => {
    const keys: openpgp.PublicKey[] = []
    for (const email of recipients) {
      const key = await getPublicKey(email)
      if (!key) {
        throw new Error(`No PGP key found for ${email}`)
      }
      keys.push(key)
    }

    // Also encrypt to self
    if (privateKey.value) {
      keys.push(privateKey.value.toPublic())
    }

    const encrypted = await openpgp.encrypt({
      message: await openpgp.createMessage({ text }),
      encryptionKeys: keys,
      signingKeys: privateKey.value || undefined,
    })

    return encrypted as string
  }

  // Decrypt message
  const decrypt = async (armoredMessage: string): Promise<{ text: string; verified: boolean }> => {
    if (!privateKey.value) {
      throw new Error('Private key not loaded')
    }

    const message = await openpgp.readMessage({ armoredMessage })

    // Get signer's public key if signed
    let verificationKeys: openpgp.PublicKey[] = []
    const signingKeyIDs = message.getSigningKeyIDs()
    for (const keyID of signingKeyIDs) {
      // Try to find the public key
      // This is simplified - real implementation would search
    }

    const { data, signatures } = await openpgp.decrypt({
      message,
      decryptionKeys: privateKey.value,
      verificationKeys,
    })

    let verified = false
    if (signatures.length > 0) {
      try {
        await signatures[0].verified
        verified = true
      } catch {
        verified = false
      }
    }

    return { text: data as string, verified }
  }

  // Sign message
  const sign = async (text: string): Promise<string> => {
    if (!privateKey.value) {
      throw new Error('Private key not loaded')
    }

    const signed = await openpgp.sign({
      message: await openpgp.createCleartextMessage({ text }),
      signingKeys: privateKey.value,
    })

    return signed as string
  }

  // Verify signature
  const verify = async (armoredMessage: string, signerEmail: string): Promise<boolean> => {
    const publicKey = await getPublicKey(signerEmail)
    if (!publicKey) return false

    const message = await openpgp.readCleartextMessage({ cleartextMessage: armoredMessage })

    const { signatures } = await openpgp.verify({
      message,
      verificationKeys: publicKey,
    })

    try {
      await signatures[0].verified
      return true
    } catch {
      return false
    }
  }

  return {
    privateKey: readonly(privateKey),
    loadPrivateKey,
    getPublicKey,
    encrypt,
    decrypt,
    sign,
    verify,
  }
}
```

**Acceptance Criteria**:
- [ ] Key loading with passphrase
- [ ] Encrypt to multiple recipients
- [ ] Decrypt received messages
- [ ] Sign outgoing messages
- [ ] Verify signatures on incoming messages

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `@tiptap/vue-3` | Rich text editor |
| `@tiptap/starter-kit` | TipTap base extensions |
| `@tiptap/extension-image` | Image support |
| `@vueuse/core` | Vue composition utilities |
| `pinia` | State management |
| `openpgp` | Browser PGP encryption |
| `@nuxt/pwa` | PWA support |
| `tailwindcss` | Styling |

---

## SQL Schema Additions

```sql
-- Labels for Gmail-like categorization
CREATE TABLE IF NOT EXISTS message_labels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    color TEXT NOT NULL DEFAULT '#1a73e8',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Message to label mapping
CREATE TABLE IF NOT EXISTS message_label_assignments (
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    label_id INTEGER NOT NULL REFERENCES message_labels(id) ON DELETE CASCADE,
    PRIMARY KEY (message_id, label_id)
);

-- Draft attachments (temporary storage)
CREATE TABLE IF NOT EXISTS draft_attachments (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size INTEGER NOT NULL,
    data BLOB NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Message templates
CREATE TABLE IF NOT EXISTS message_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    subject TEXT,
    body_html TEXT,
    body_text TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_message_labels_user ON message_labels(user_id);
CREATE INDEX idx_draft_attachments_user ON draft_attachments(user_id);
CREATE INDEX idx_draft_attachments_created ON draft_attachments(created_at);
```

---

## 7.6 Health Monitoring [MVP]

Operational health check endpoints for Kubernetes liveness/readiness probes and monitoring.

| ID | Task | Status | Dependencies | Priority |
|----|------|--------|--------------|----------|
| WM-033 | Health check endpoints | [ ] | - | MVP |

---

### WM-033: Health Check Endpoints

**File**: `internal/api/health_handler.go`
```go
package api

import (
    "context"
    "database/sql"
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/btafoya/gomailserver/internal/smtp"
    "github.com/btafoya/gomailserver/internal/imap"
)

type HealthHandler struct {
    db          *sql.DB
    smtpServer  *smtp.Server
    imapServer  *imap.Server
}

type HealthResponse struct {
    Status    string            `json:"status"` // "healthy", "degraded", "unhealthy"
    Timestamp string            `json:"timestamp"`
    Services  map[string]string `json:"services,omitempty"`
    Version   string            `json:"version"`
}

// GET /health/live - Kubernetes liveness probe
// Returns 200 if process is alive (minimal check)
func (h *HealthHandler) Live(c echo.Context) error {
    return c.JSON(http.StatusOK, HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Version:   "1.0.0",
    })
}

// GET /health/ready - Kubernetes readiness probe
// Returns 200 only if all critical services are operational
func (h *HealthHandler) Ready(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    services := make(map[string]string)
    overallHealthy := true

    // Check database connectivity
    if err := h.checkDatabase(ctx); err != nil {
        services["database"] = "unhealthy: " + err.Error()
        overallHealthy = false
    } else {
        services["database"] = "healthy"
    }

    // Check SMTP server
    if err := h.checkSMTP(); err != nil {
        services["smtp"] = "unhealthy: " + err.Error()
        overallHealthy = false
    } else {
        services["smtp"] = "healthy"
    }

    // Check IMAP server
    if err := h.checkIMAP(); err != nil {
        services["imap"] = "unhealthy: " + err.Error()
        overallHealthy = false
    } else {
        services["imap"] = "healthy"
    }

    status := "healthy"
    statusCode := http.StatusOK
    if !overallHealthy {
        status = "unhealthy"
        statusCode = http.StatusServiceUnavailable
    }

    return c.JSON(statusCode, HealthResponse{
        Status:    status,
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Services:  services,
        Version:   "1.0.0",
    })
}

func (h *HealthHandler) checkDatabase(ctx context.Context) error {
    return h.db.PingContext(ctx)
}

func (h *HealthHandler) checkSMTP() error {
    if h.smtpServer == nil {
        return fmt.Errorf("SMTP server not initialized")
    }
    if !h.smtpServer.IsRunning() {
        return fmt.Errorf("SMTP server not running")
    }
    return nil
}

func (h *HealthHandler) checkIMAP() error {
    if h.imapServer == nil {
        return fmt.Errorf("IMAP server not initialized")
    }
    if !h.imapServer.IsRunning() {
        return fmt.Errorf("IMAP server not running")
    }
    return nil
}
```

**File**: `cmd/gomailserver/main.go` (register routes)
```go
// Health check endpoints
healthHandler := api.NewHealthHandler(db, smtpServer, imapServer)
e.GET("/health/live", healthHandler.Live)
e.GET("/health/ready", healthHandler.Ready)
```

**Acceptance Criteria**:
- [ ] `/health/live` endpoint returns 200 OK with minimal checks
- [ ] `/health/ready` endpoint checks database, SMTP, and IMAP health
- [ ] Ready endpoint returns 503 if any critical service is unhealthy
- [ ] Health checks complete in < 100ms (P99)
- [ ] JSON response includes service status breakdown
- [ ] Version information included in response
- [ ] Kubernetes liveness/readiness probe compatible

**Production Readiness**:
- [ ] Timeout: 5s for health checks (prevents hanging probes)
- [ ] Response time: P99 < 100ms
- [ ] Database: Ping with context timeout
- [ ] SMTP: Check server running status
- [ ] IMAP: Check server running status
- [ ] Status codes: 200 OK (healthy), 503 Service Unavailable (unhealthy)
- [ ] Format: Kubernetes-compatible health check responses

**Given/When/Then Scenarios**:
```
Given all services are running
When GET /health/ready is called
Then response status is 200 OK
And all services report "healthy"
And response time is < 100ms

Given database is unreachable
When GET /health/ready is called
Then response status is 503 Service Unavailable
And database service reports "unhealthy"
And other services still report their actual status

Given process is running but SMTP failed
When GET /health/live is called
Then response status is 200 OK (process is alive)
When GET /health/ready is called
Then response status is 503 Service Unavailable (not ready for traffic)
```

---

## Testing Checklist

- [ ] Message list loads with pagination
- [ ] Conversation threading works correctly
- [ ] Compose and send message
- [ ] Reply/Reply All/Forward
- [ ] Attachments upload and download
- [ ] Search returns accurate results
- [ ] Keyboard shortcuts work
- [ ] Dark mode displays correctly
- [ ] Mobile responsive on all screen sizes
- [ ] PWA installs and works offline (read cached)
- [ ] Contact autocomplete functions
- [ ] Calendar invitations display and respond
- [ ] PGP encryption/decryption works
- [ ] Auto-save drafts every 30 seconds
