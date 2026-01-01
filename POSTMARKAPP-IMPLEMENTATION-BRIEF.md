# PostmarkApp API Implementation - Technical Brief

**Status**: Planning
**Priority**: High
**Complexity**: Moderate
**Estimated Duration**: 3-4 weeks

---

## Executive Summary

Implement PostmarkApp-compatible REST API endpoints in gomailserver to enable drop-in replacement of PostmarkApp services. This allows existing applications using PostmarkApp client libraries to switch to gomailserver without code changes, providing full control over email infrastructure while maintaining API compatibility.

## Business Value

- **Migration Path**: Enables seamless migration from PostmarkApp to self-hosted gomailserver
- **Cost Reduction**: Eliminates PostmarkApp subscription fees for high-volume senders
- **Data Sovereignty**: Full control over email data and infrastructure
- **Feature Parity**: Maintains compatibility with existing PostmarkApp integrations
- **Client Library Support**: Works with official PostmarkApp libraries (Node.js, Python, Ruby, PHP, .NET, Go)

---

## Research Findings

### PostmarkApp API Architecture

**Base URL**: `https://api.postmarkapp.com`

**Authentication Methods**:
- `X-Postmark-Server-Token` - Server-level privileges
- `X-Postmark-Account-Token` - Account-level privileges
- Headers are case-insensitive
- Test mode uses `POSTMARK_API_TEST` token value

**Request/Response Format**:
- Protocol: HTTPS with TLS encryption
- Architecture: REST-based
- Content-Type: `application/json`
- Response codes: 200 (success), 401 (unauthorized), 422 (validation error), 429 (rate limit), 500 (server error)

**Error Structure**:
```json
{
  "ErrorCode": 123,
  "Message": "Error description"
}
```

Error codes range: 10-1302 (authentication, validation, resource, rate limit, service errors)

### Core API Endpoints (Priority Order)

#### 1. Email Sending (Critical - MVP)
- **POST /email** - Send single email
- **POST /email/batch** - Send up to 500 emails in one request
- **POST /email/withTemplate** - Send using template

**Request Structure** (POST /email):
```json
{
  "From": "sender@example.com",
  "To": "recipient@example.com",
  "Cc": "cc@example.com",
  "Bcc": "bcc@example.com",
  "Subject": "Email subject",
  "TextBody": "Plain text content",
  "HtmlBody": "<html><body>HTML content</body></html>",
  "ReplyTo": "reply@example.com",
  "Tag": "category-tag",
  "Metadata": {"key": "value"},
  "Headers": [{"Name": "X-Custom", "Value": "value"}],
  "TrackOpens": true,
  "TrackLinks": "HtmlAndText",
  "MessageStream": "outbound",
  "Attachments": [{
    "Name": "file.pdf",
    "Content": "base64-encoded-content",
    "ContentType": "application/pdf"
  }]
}
```

**Response Structure**:
```json
{
  "ErrorCode": 0,
  "Message": "OK",
  "MessageID": "b7bc2f4a-e38e-4336-af7d-e6c392c2f817",
  "SubmittedAt": "2010-11-26T12:01:05.1794748-05:00",
  "To": "recipient@example.com"
}
```

**Limits**:
- Max 50 recipients total (To + Cc + Bcc)
- Max message size: 10MB including attachments
- Subject max: 2000 characters
- Tag max: 1000 characters

#### 2. Templates (High Priority)
- **GET /templates** - List all templates
- **GET /templates/{templateIdOrAlias}** - Get template details
- **POST /templates** - Create template
- **PUT /templates/{templateId}** - Update template
- **DELETE /templates/{templateId}** - Delete template

**Template Structure**:
```json
{
  "TemplateId": 123,
  "Name": "Welcome Email",
  "Alias": "welcome",
  "Subject": "Welcome {{name}}!",
  "HtmlBody": "<html><body>Hello {{name}}</body></html>",
  "TextBody": "Hello {{name}}",
  "Active": true,
  "TemplateType": "Standard",
  "LayoutTemplate": null,
  "AssociatedServerId": 456
}
```

#### 3. Webhooks (High Priority)
- **GET /webhooks** - List webhooks
- **GET /webhooks/{webhookId}** - Get webhook details
- **POST /webhooks** - Create webhook
- **PUT /webhooks/{webhookId}** - Update webhook
- **DELETE /webhooks/{webhookId}** - Delete webhook

**Webhook Structure**:
```json
{
  "ID": 123,
  "Url": "https://example.com/webhook",
  "MessageStream": "outbound",
  "HttpAuth": {
    "Username": "user",
    "Password": "pass"
  },
  "HttpHeaders": [{"Name": "X-Custom", "Value": "value"}],
  "Triggers": {
    "Open": {"Enabled": true, "PostFirstOpenOnly": false},
    "Click": {"Enabled": true},
    "Delivery": {"Enabled": true},
    "Bounce": {"Enabled": true, "IncludeContent": false},
    "SpamComplaint": {"Enabled": true, "IncludeContent": false},
    "SubscriptionChange": {"Enabled": false}
  }
}
```

**Webhook Payload Example** (Bounce):
```json
{
  "RecordType": "Bounce",
  "MessageID": "uuid",
  "Type": "HardBounce",
  "TypeCode": 1,
  "Email": "recipient@example.com",
  "BouncedAt": "2023-01-01T12:00:00Z",
  "Details": "SMTP 550 User unknown",
  "Subject": "Email subject",
  "Tag": "tag",
  "Metadata": {"key": "value"}
}
```

#### 4. Messages (Medium Priority)
- **GET /messages/outbound** - List sent messages
- **GET /messages/outbound/{messageId}/details** - Get message details
- **GET /messages/outbound/{messageId}/opens** - Get open tracking
- **GET /messages/outbound/{messageId}/clicks** - Get click tracking

#### 5. Bounces (Medium Priority)
- **GET /bounces** - List bounces
- **GET /bounces/{bounceId}** - Get bounce details
- **PUT /bounces/{bounceId}/activate** - Reactivate bounced email

#### 6. Server/Account Management (Low Priority)
- **GET /server** - Get server details
- **PATCH /server** - Update server settings
- **GET /servers** - List servers (account token)
- **POST /servers** - Create server (account token)

---

## Technical Architecture

### Component Structure

```
gomailserver/
├── internal/
│   ├── postmark/
│   │   ├── handlers/
│   │   │   ├── email_handler.go          # POST /email, /email/batch
│   │   │   ├── template_handler.go       # Templates CRUD
│   │   │   ├── webhook_handler.go        # Webhooks CRUD
│   │   │   ├── message_handler.go        # Message tracking/details
│   │   │   ├── bounce_handler.go         # Bounce management
│   │   │   └── server_handler.go         # Server management
│   │   ├── middleware/
│   │   │   ├── auth.go                   # X-Postmark-*-Token validation
│   │   │   ├── error_handler.go          # PostmarkApp error format
│   │   │   └── rate_limit.go             # API rate limiting
│   │   ├── models/
│   │   │   ├── email.go                  # Email request/response
│   │   │   ├── template.go               # Template structures
│   │   │   ├── webhook.go                # Webhook structures
│   │   │   └── errors.go                 # PostmarkApp error codes
│   │   ├── service/
│   │   │   ├── email_service.go          # Email sending logic
│   │   │   ├── template_service.go       # Template rendering
│   │   │   ├── webhook_service.go        # Webhook delivery
│   │   │   └── tracking_service.go       # Open/click tracking
│   │   ├── repository/
│   │   │   ├── postmark_repo.go          # PostmarkApp-specific DB ops
│   │   │   └── sqlite/
│   │   │       └── postmark_sqlite.go    # SQLite implementation
│   │   └── router.go                     # PostmarkApp API router
│   └── api/
│       └── router.go                     # Mount /postmark/* routes
└── web/
    └── admin/
        └── src/
            └── views/
                └── postmark/             # PostmarkApp management UI
                    ├── Servers.vue       # Server/token management
                    ├── Templates.vue     # Template editor
                    └── Webhooks.vue      # Webhook configuration
```

### Database Schema (Migration V5)

```sql
-- PostmarkApp Servers (API token groups)
CREATE TABLE postmark_servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    api_token TEXT NOT NULL UNIQUE,  -- bcrypt hashed
    account_id INTEGER REFERENCES users(id),
    message_stream TEXT DEFAULT 'outbound',
    track_opens BOOLEAN DEFAULT false,
    track_links TEXT DEFAULT 'None',  -- None, HtmlOnly, HtmlAndText, TextOnly
    active BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_servers_token ON postmark_servers(api_token);
CREATE INDEX idx_postmark_servers_account ON postmark_servers(account_id);

-- PostmarkApp Message Tracking
CREATE TABLE postmark_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT NOT NULL UNIQUE,  -- UUID format
    server_id INTEGER REFERENCES postmark_servers(id),
    from_email TEXT NOT NULL,
    to_email TEXT NOT NULL,
    cc_email TEXT,
    bcc_email TEXT,
    subject TEXT,
    html_body TEXT,
    text_body TEXT,
    tag TEXT,
    metadata TEXT,  -- JSON
    message_stream TEXT,
    status TEXT DEFAULT 'pending',  -- pending, sent, bounced, delivered
    submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    sent_at DATETIME,
    delivered_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_messages_id ON postmark_messages(message_id);
CREATE INDEX idx_postmark_messages_server ON postmark_messages(server_id);
CREATE INDEX idx_postmark_messages_status ON postmark_messages(status);
CREATE INDEX idx_postmark_messages_submitted ON postmark_messages(submitted_at);

-- PostmarkApp Templates
CREATE TABLE postmark_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER REFERENCES postmark_servers(id),
    name TEXT NOT NULL,
    alias TEXT,
    subject TEXT,
    html_body TEXT,
    text_body TEXT,
    template_type TEXT DEFAULT 'Standard',  -- Standard, Layout
    layout_template INTEGER REFERENCES postmark_templates(id),
    active BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(server_id, alias)
);

CREATE INDEX idx_postmark_templates_server ON postmark_templates(server_id);
CREATE INDEX idx_postmark_templates_alias ON postmark_templates(alias);

-- PostmarkApp Webhooks
CREATE TABLE postmark_webhooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER REFERENCES postmark_servers(id),
    url TEXT NOT NULL,
    message_stream TEXT DEFAULT 'outbound',
    http_auth_username TEXT,
    http_auth_password TEXT,  -- encrypted
    http_headers TEXT,  -- JSON array
    triggers TEXT NOT NULL,  -- JSON object with enabled flags
    active BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_webhooks_server ON postmark_webhooks(server_id);

-- PostmarkApp Bounces
CREATE TABLE postmark_bounces (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT REFERENCES postmark_messages(message_id),
    type TEXT NOT NULL,  -- HardBounce, SoftBounce, Transient, etc.
    type_code INTEGER,
    email TEXT NOT NULL,
    bounced_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    details TEXT,
    inactive BOOLEAN DEFAULT true,
    can_activate BOOLEAN DEFAULT false,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_bounces_message ON postmark_bounces(message_id);
CREATE INDEX idx_postmark_bounces_email ON postmark_bounces(email);

-- PostmarkApp Tracking Events
CREATE TABLE postmark_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT REFERENCES postmark_messages(message_id),
    event_type TEXT NOT NULL,  -- Open, Click, Delivery, SpamComplaint
    recipient TEXT NOT NULL,
    user_agent TEXT,
    client_info TEXT,  -- JSON
    location TEXT,  -- JSON (geo)
    link_url TEXT,  -- for clicks
    occurred_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_postmark_events_message ON postmark_events(message_id);
CREATE INDEX idx_postmark_events_type ON postmark_events(event_type);
CREATE INDEX idx_postmark_events_occurred ON postmark_events(occurred_at);
```

### Integration Points

#### 1. SMTP Service Integration
```go
// Map PostmarkApp email request to internal message
func (s *EmailService) SendEmail(req *models.EmailRequest) (*models.EmailResponse, error) {
    // Convert PostmarkApp format to internal format
    msg := &domain.Message{
        From:    req.From,
        To:      strings.Split(req.To, ","),
        Subject: req.Subject,
        HTML:    req.HtmlBody,
        Text:    req.TextBody,
        Headers: convertHeaders(req.Headers),
    }

    // Use existing SMTP queue service
    messageID, err := s.queueService.EnqueueMessage(msg)
    if err != nil {
        return nil, convertToPostmarkError(err)
    }

    // Store PostmarkApp tracking data
    err = s.repo.CreateMessage(&postmark.Message{
        MessageID:   messageID,
        ServerID:    req.ServerID,
        FromEmail:   req.From,
        ToEmail:     req.To,
        Subject:     req.Subject,
        SubmittedAt: time.Now(),
    })

    return &models.EmailResponse{
        ErrorCode:   0,
        Message:     "OK",
        MessageID:   messageID,
        SubmittedAt: time.Now(),
        To:          req.To,
    }, nil
}
```

#### 2. Template Rendering
```go
// Use Go html/template for variable substitution
func (s *TemplateService) RenderTemplate(templateID int, model map[string]interface{}) (string, string, error) {
    tmpl, err := s.repo.GetTemplate(templateID)
    if err != nil {
        return "", "", err
    }

    htmlTmpl, _ := template.New("html").Parse(tmpl.HtmlBody)
    textTmpl, _ := template.New("text").Parse(tmpl.TextBody)

    var htmlBuf, textBuf bytes.Buffer
    htmlTmpl.Execute(&htmlBuf, model)
    textTmpl.Execute(&textBuf, model)

    return htmlBuf.String(), textBuf.String(), nil
}
```

#### 3. Webhook Delivery
```go
// Leverage existing queue worker pattern
func (s *WebhookService) DeliverEvent(event *models.WebhookEvent) error {
    webhooks, err := s.repo.GetActiveWebhooks(event.ServerID, event.Type)

    for _, webhook := range webhooks {
        payload, _ := json.Marshal(event)

        req, _ := http.NewRequest("POST", webhook.Url, bytes.NewBuffer(payload))
        req.Header.Set("Content-Type", "application/json")

        if webhook.HttpAuth != nil {
            req.SetBasicAuth(webhook.HttpAuth.Username, webhook.HttpAuth.Password)
        }

        // Add retry logic using existing queue system
        s.queueService.EnqueueWebhook(req, 3) // 3 retries
    }

    return nil
}
```

#### 4. Admin UI Integration
Add PostmarkApp management to existing Admin UI:

**Router Update** (`web/admin/src/router/index.js`):
```javascript
{
  path: '/postmark',
  component: () => import('@/views/postmark/Layout.vue'),
  children: [
    { path: 'servers', component: () => import('@/views/postmark/Servers.vue') },
    { path: 'templates', component: () => import('@/views/postmark/Templates.vue') },
    { path: 'webhooks', component: () => import('@/views/postmark/Webhooks.vue') },
    { path: 'messages', component: () => import('@/views/postmark/Messages.vue') }
  ]
}
```

---

## Implementation Plan

### Milestone 1: Core Email API (Week 1-2)

**Goal**: Implement single email sending with PostmarkApp compatibility

**Tasks**:
1. Create `internal/postmark` package structure
2. Implement authentication middleware (`X-Postmark-Server-Token`)
3. Create database migration V5 (postmark_servers, postmark_messages tables)
4. Implement POST /email endpoint handler
5. Integrate with existing SMTP queue service
6. Add PostmarkApp error code mapping
7. Unit tests for email sending flow

**Deliverables**:
- Working POST /email endpoint
- Server token authentication
- Message tracking in database
- Integration with SMTP queue

### Milestone 2: Batch & Templates (Week 2-3)

**Goal**: Add batch sending and template support

**Tasks**:
1. Implement POST /email/batch endpoint (500 message limit)
2. Add template database tables (postmark_templates)
3. Implement template CRUD endpoints (GET, POST, PUT, DELETE /templates)
4. Build template rendering service with Go html/template
5. Implement POST /email/withTemplate endpoint
6. Add template validation and syntax checking
7. Unit tests for batch and template operations

**Deliverables**:
- Batch email sending (up to 500)
- Template CRUD operations
- Template-based email sending
- Template rendering engine

### Milestone 3: Webhooks & Tracking (Week 3-4)

**Goal**: Implement webhook system and event tracking

**Tasks**:
1. Add webhook database tables (postmark_webhooks, postmark_events)
2. Implement webhook CRUD endpoints
3. Build webhook delivery service with retry logic
4. Implement event tracking (opens, clicks, bounces, deliveries)
5. Add webhook payload generation for each event type
6. Integrate with bounce processing
7. Add webhook authentication (HTTP Basic Auth)
8. Unit tests for webhook delivery

**Deliverables**:
- Webhook CRUD operations
- Event tracking system
- Webhook delivery with retries
- Bounce webhook integration

### Milestone 4: Admin UI & Testing (Week 4)

**Goal**: Add management interface and comprehensive testing

**Tasks**:
1. Create Admin UI views for PostmarkApp management
2. Build Server/Token management interface
3. Implement template editor with preview
4. Add webhook configuration UI
5. Create message tracking/history view
6. Integration tests with PostmarkApp client libraries
7. Load testing for batch endpoints
8. Documentation and API reference

**Deliverables**:
- Complete Admin UI for PostmarkApp features
- Integration test suite
- Performance benchmarks
- API documentation

---

## PostmarkApp Error Codes

Map internal errors to PostmarkApp error codes:

| Code | Meaning | Internal Mapping |
|------|---------|------------------|
| 0 | Success | Success |
| 300 | Invalid email request | ValidationError |
| 400 | Sender signature not found | InvalidSender |
| 401 | Unauthorized | AuthenticationError |
| 402 | Inactive recipient | BounceInactive |
| 405 | Invalid JSON | JSONParseError |
| 406 | Inactive server | ServerInactive |
| 409 | JSON required | ContentTypeError |
| 410 | Too many batch messages | BatchLimitExceeded |
| 411 | Forbidden attachment type | AttachmentForbidden |
| 429 | Rate limit exceeded | RateLimitError |
| 500 | Internal server error | InternalError |
| 503 | Service unavailable | ServiceUnavailable |

**Error Response Format**:
```json
{
  "ErrorCode": 300,
  "Message": "The 'From' address is required"
}
```

---

## API Rate Limiting

Implement rate limiting per PostmarkApp specifications:

**Limits**:
- **Account-level**: 100,000 emails/hour
- **Server-level**: 50,000 emails/hour
- **API calls**: 1,000 requests/minute per server token
- **Batch endpoint**: 500 messages per request

**Implementation**:
```go
// Rate limit middleware
func RateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Every(window/time.Duration(limit)), limit)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                writePostmarkError(w, 429, "Rate limit exceeded")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Testing Strategy

### Unit Tests
- Handler logic for all endpoints
- Request/response serialization
- Error code mapping
- Template rendering
- Authentication middleware

### Integration Tests
- End-to-end email sending via POST /email
- Batch sending with 500 messages
- Template-based email sending
- Webhook delivery and retry
- PostmarkApp client library compatibility

### Load Tests
- Concurrent email sending (1000 req/sec)
- Batch endpoint performance
- Template rendering under load
- Webhook delivery queue

### Compatibility Tests
Test with official PostmarkApp client libraries:
- **Node.js**: `postmark` npm package
- **Python**: `postmarker` pip package
- **Ruby**: `postmark` gem
- **PHP**: `postmark/postmark-php` composer package
- **.NET**: `PostmarkDotNet` NuGet package
- **Go**: `postmark` go module

---

## Security Considerations

1. **API Token Storage**: Store tokens using bcrypt hashing (like API keys)
2. **Webhook Authentication**: Support HTTP Basic Auth for webhook URLs
3. **Rate Limiting**: Prevent abuse with per-token limits
4. **Input Validation**: Sanitize all email content, prevent injection
5. **TLS Required**: Enforce HTTPS for all PostmarkApp endpoints
6. **CORS**: Configure allowed origins for browser-based clients
7. **Content Security**: Scan attachments with ClamAV integration
8. **Token Rotation**: Support token regeneration in Admin UI

---

## Migration Guide (for PostmarkApp Users)

**Step 1**: Create PostmarkApp Server in gomailserver Admin UI
- Generate new server token
- Configure message stream settings

**Step 2**: Update application configuration
```javascript
// Before (PostmarkApp)
const client = new postmark.ServerClient("YOUR_POSTMARK_TOKEN");

// After (gomailserver)
const client = new postmark.ServerClient("YOUR_GOMAILSERVER_TOKEN");
client.setBaseUrl("https://your-gomailserver.com/postmark");
```

**Step 3**: Migrate templates
- Export templates from PostmarkApp
- Import via gomailserver Admin UI or API

**Step 4**: Configure webhooks
- Update webhook URLs to point to gomailserver
- Verify webhook delivery with test events

**Step 5**: Test and validate
- Send test emails via API
- Verify delivery and tracking
- Check webhook delivery

---

## Success Criteria

✅ PostmarkApp client libraries work without code changes
✅ All critical endpoints implemented (email, batch, templates, webhooks)
✅ Error codes match PostmarkApp specifications
✅ Admin UI provides complete management interface
✅ Integration tests pass with official client libraries
✅ Performance: 1000+ emails/sec via batch endpoint
✅ Documentation complete with migration guide
✅ Security audit passed (authentication, validation, rate limiting)

---

## Future Enhancements

**Phase 2 (Post-MVP)**:
- Advanced tracking (geolocation, device detection)
- Suppression lists management
- A/B testing for subject lines
- Scheduled email sending
- Email preview/testing tools
- Advanced analytics dashboard
- DKIM/SPF/DMARC reporting integration
- Inbound email processing API
- Message retention policies
- Data export tools

---

## Resources & References

### Official Documentation
- [PostmarkApp API Overview](https://postmarkapp.com/developer/api/overview)
- [Sending Email API](https://postmarkapp.com/developer/user-guide/send-email-with-api)
- [Templates API](https://postmarkapp.com/developer/api/templates-api)
- [Webhooks API](https://postmarkapp.com/developer/api/webhooks-api)
- [Webhooks Overview](https://postmarkapp.com/developer/webhooks/webhooks-overview)

### Client Libraries
- [Postmark.js (Node.js)](https://github.com/activecampaign/postmark.js)
- [PostmarkApp Templates](https://github.com/activecampaign/postmark-templates)

### Go Libraries
- [Chi Router Documentation](https://go-chi.github.io/chi/)
- [Go html/template](https://pkg.go.dev/html/template)

---

## Appendix: Complete Endpoint List

### Email Endpoints
- `POST /email` - Send single email
- `POST /email/batch` - Send batch emails
- `POST /email/withTemplate` - Send with template

### Template Endpoints
- `GET /templates` - List templates
- `GET /templates/{templateIdOrAlias}` - Get template
- `POST /templates` - Create template
- `PUT /templates/{templateId}` - Update template
- `DELETE /templates/{templateId}` - Delete template
- `POST /templates/validate` - Validate template

### Webhook Endpoints
- `GET /webhooks` - List webhooks
- `GET /webhooks/{webhookId}` - Get webhook
- `POST /webhooks` - Create webhook
- `PUT /webhooks/{webhookId}` - Update webhook
- `DELETE /webhooks/{webhookId}` - Delete webhook

### Message Endpoints
- `GET /messages/outbound` - List outbound messages
- `GET /messages/outbound/{messageId}/details` - Message details
- `GET /messages/outbound/{messageId}/opens` - Open events
- `GET /messages/outbound/{messageId}/clicks` - Click events
- `GET /messages/outbound/{messageId}/dump` - Raw message

### Bounce Endpoints
- `GET /bounces` - List bounces
- `GET /bounces/{bounceId}` - Get bounce
- `PUT /bounces/{bounceId}/activate` - Reactivate email
- `GET /deliverystats` - Delivery statistics

### Server Endpoints
- `GET /server` - Get server info
- `PATCH /server` - Update server

---

**End of Implementation Brief**
