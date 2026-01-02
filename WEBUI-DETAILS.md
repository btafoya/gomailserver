# gomailserver - Web UI and API Details

**Last Updated**: 2026-01-02
**Project**: gomailserver (github.com/btafoya/gomailserver)

---

## 🌐 Web UI Access Points

### Admin UI
- **URL**: `http://localhost:8980/admin`
- **HTTPS**: `https://your-domain.com:8980/admin` (with TLS enabled)
- **Framework**: Vue.js 3 + Vite
- **Type**: Embedded SPA (bundled in binary)
- **Features**:
  - Domain management
  - User account management
  - Alias configuration
  - DKIM key management
  - Queue monitoring
  - Statistics dashboard
  - Audit log viewer
  - System settings

**First-Time Setup**:
```bash
# Access setup wizard (no authentication required until admin created)
http://localhost:8980/api/v1/setup/status
```

**Login Credentials**:
- Created via setup wizard or CLI
- Authentication: JWT tokens
- Session management via cookies

---

### Webmail UI
- **URL**: `http://localhost:8980/webmail`
- **HTTPS**: `https://your-domain.com:8980/webmail`
- **Framework**: Nuxt 3 + Vue 3 + Tailwind CSS
- **Type**: Embedded SPA (bundled in binary)
- **Features**:
  - Email inbox and folder management
  - Rich text composer (TipTap editor)
  - Attachment handling (upload/download)
  - Draft auto-save
  - Dark mode support
  - Mobile responsive design
  - Keyboard shortcuts

**Login**:
- Uses email address and password
- Same credentials as IMAP/SMTP
- JWT-based session

---

### Self-Service Portal
- **URL**: `http://localhost:8980/portal` (planned)
- **Status**: ⚠️ Currently integrated with Admin UI
- **Features** (when implemented):
  - Password change
  - 2FA setup
  - Alias management
  - Quota usage display
  - Forwarding rules
  - Spam quarantine viewer

**Current Access**:
Users can access self-service features through the Admin UI with user-level permissions.

---

## 🔌 API Endpoints

### Base Configuration
- **Default Port**: `8980`
- **Environment Variable**: `API_PORT`
- **Configuration**: `api.port` in YAML config
- **Base Path**: `/api/v1`
- **Format**: JSON (application/json)
- **CORS**: Configurable origins

---

### Health & Status

#### Health Check
```http
GET /health
```
- **Authentication**: None
- **Response**: `{"status": "ok"}`
- **Use Case**: Load balancer health checks

---

### Authentication API

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "your-password"
}
```
- **Response**: JWT token
- **Rate Limited**: Yes (brute force protection)

#### Refresh Token
```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```
- **Response**: New JWT token

---

### Setup Wizard API

#### Get Setup Status
```http
GET /api/v1/setup/status
```
- **Authentication**: None (until admin created)
- **Returns**: Whether first-time setup is complete

#### Get Setup State
```http
GET /api/v1/setup/state
```
- **Authentication**: None
- **Returns**: Current setup progress

#### Create Admin User
```http
POST /api/v1/setup/admin
Content-Type: application/json

{
  "email": "admin@yourdomain.com",
  "password": "secure-password",
  "domain": "yourdomain.com"
}
```

#### Complete Setup
```http
POST /api/v1/setup/complete
```

---

### Domain Management API

#### List Domains
```http
GET /api/v1/domains
Authorization: Bearer <token>
```

#### Create Domain
```http
POST /api/v1/domains
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "example.com",
  "enabled": true
}
```

#### Get Domain
```http
GET /api/v1/domains/{id}
Authorization: Bearer <token>
```

#### Update Domain
```http
PUT /api/v1/domains/{id}
Authorization: Bearer <token>
```

#### Delete Domain
```http
DELETE /api/v1/domains/{id}
Authorization: Bearer <token>
```

#### Generate DKIM Keys
```http
POST /api/v1/domains/{id}/dkim
Authorization: Bearer <token>
Content-Type: application/json

{
  "selector": "default",
  "key_type": "rsa-2048"
}
```

---

### User Management API

#### List Users
```http
GET /api/v1/users
Authorization: Bearer <token>
```

#### Create User
```http
POST /api/v1/users
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure-password",
  "domain_id": 1,
  "quota_bytes": 1073741824
}
```

#### Get User
```http
GET /api/v1/users/{id}
Authorization: Bearer <token>
```

#### Update User
```http
PUT /api/v1/users/{id}
Authorization: Bearer <token>
```

#### Delete User
```http
DELETE /api/v1/users/{id}
Authorization: Bearer <token>
```

#### Reset Password
```http
POST /api/v1/users/{id}/password
Authorization: Bearer <token>
Content-Type: application/json

{
  "password": "new-password"
}
```

---

### Alias Management API

#### List Aliases
```http
GET /api/v1/aliases
Authorization: Bearer <token>
```

#### Create Alias
```http
POST /api/v1/aliases
Authorization: Bearer <token>
Content-Type: application/json

{
  "source": "alias@example.com",
  "destination": "user@example.com"
}
```

#### Get Alias
```http
GET /api/v1/aliases/{id}
Authorization: Bearer <token>
```

#### Delete Alias
```http
DELETE /api/v1/aliases/{id}
Authorization: Bearer <token>
```

---

### Statistics API

#### Dashboard Stats
```http
GET /api/v1/stats/dashboard
Authorization: Bearer <token>
```

#### Domain Stats
```http
GET /api/v1/stats/domains/{id}
Authorization: Bearer <token>
```

#### User Stats
```http
GET /api/v1/stats/users/{id}
Authorization: Bearer <token>
```

---

### Queue Management API

#### List Queue Items
```http
GET /api/v1/queue
Authorization: Bearer <token>
```

#### Get Queue Item
```http
GET /api/v1/queue/{id}
Authorization: Bearer <token>
```

#### Retry Queue Item
```http
POST /api/v1/queue/{id}/retry
Authorization: Bearer <token>
```

#### Delete Queue Item
```http
DELETE /api/v1/queue/{id}
Authorization: Bearer <token>
```

---

### Settings API

#### Get Server Settings
```http
GET /api/v1/settings/server
Authorization: Bearer <token>
```

#### Update Server Settings
```http
PUT /api/v1/settings/server
Authorization: Bearer <token>
```

#### Get Security Settings
```http
GET /api/v1/settings/security
Authorization: Bearer <token>
```

#### Update Security Settings
```http
PUT /api/v1/settings/security
Authorization: Bearer <token>
```

#### Get TLS Settings
```http
GET /api/v1/settings/tls
Authorization: Bearer <token>
```

#### Update TLS Settings
```http
PUT /api/v1/settings/tls
Authorization: Bearer <token>
```

---

### PGP Key Management API

#### Import PGP Key
```http
POST /api/v1/pgp/keys
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 1,
  "public_key": "-----BEGIN PGP PUBLIC KEY BLOCK-----..."
}
```

#### List User Keys
```http
GET /api/v1/pgp/users/{user_id}/keys
Authorization: Bearer <token>
```

#### Get Key
```http
GET /api/v1/pgp/keys/{id}
Authorization: Bearer <token>
```

#### Set Primary Key
```http
POST /api/v1/pgp/keys/{id}/primary
Authorization: Bearer <token>
```

#### Delete Key
```http
DELETE /api/v1/pgp/keys/{id}
Authorization: Bearer <token>
```

---

### Audit Log API

#### List Audit Logs
```http
GET /api/v1/audit/logs
Authorization: Bearer <token>
```

#### Get Audit Stats
```http
GET /api/v1/audit/stats
Authorization: Bearer <token>
```

---

### Webmail API

#### List Mailboxes
```http
GET /api/v1/webmail/mailboxes
Authorization: Bearer <token>
```

#### List Messages
```http
GET /api/v1/webmail/mailboxes/{id}/messages?offset=0&limit=50
Authorization: Bearer <token>
```

#### Get Message
```http
GET /api/v1/webmail/messages/{id}
Authorization: Bearer <token>
```

#### Send Message
```http
POST /api/v1/webmail/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "to": ["recipient@example.com"],
  "subject": "Hello",
  "body": "Message body",
  "html": true
}
```

#### Delete Message
```http
DELETE /api/v1/webmail/messages/{id}
Authorization: Bearer <token>
```

#### Move Message
```http
POST /api/v1/webmail/messages/{id}/move
Authorization: Bearer <token>
Content-Type: application/json

{
  "mailbox_id": 2
}
```

#### Update Message Flags
```http
POST /api/v1/webmail/messages/{id}/flags
Authorization: Bearer <token>
Content-Type: application/json

{
  "add": ["\\Seen", "\\Flagged"],
  "remove": []
}
```

#### Search Messages
```http
GET /api/v1/webmail/search?query=subject:test
Authorization: Bearer <token>
```

#### Download Attachment
```http
GET /api/v1/webmail/attachments/{id}
Authorization: Bearer <token>
```

#### Save Draft
```http
POST /api/v1/webmail/drafts
Authorization: Bearer <token>
Content-Type: application/json

{
  "to": ["recipient@example.com"],
  "subject": "Draft subject",
  "body": "Draft body"
}
```

#### List Drafts
```http
GET /api/v1/webmail/drafts
Authorization: Bearer <token>
```

#### Get Draft
```http
GET /api/v1/webmail/drafts/{id}
Authorization: Bearer <token>
```

#### Delete Draft
```http
DELETE /api/v1/webmail/drafts/{id}
Authorization: Bearer <token>
```

---

### Logs API

#### Get Logs
```http
GET /api/v1/logs?since=2026-01-01T00:00:00Z&level=error
Authorization: Bearer <token>
```

---

## 📧 PostmarkApp-Compatible API

### Base Configuration
- **Default Port**: `8980` (same as main API)
- **Base Path**: `/` (root-level endpoints)
- **Format**: JSON (PostmarkApp-compatible)
- **Authentication**: `X-Postmark-Server-Token` header

### Email Endpoints

#### Send Single Email
```http
POST /email
Content-Type: application/json
X-Postmark-Server-Token: <your-server-token>

{
  "From": "sender@example.com",
  "To": "recipient@example.com",
  "Subject": "Hello",
  "TextBody": "Plain text body",
  "HtmlBody": "<p>HTML body</p>",
  "Cc": "cc@example.com",
  "Bcc": "bcc@example.com",
  "ReplyTo": "reply@example.com",
  "Attachments": [
    {
      "Name": "document.pdf",
      "Content": "base64-encoded-content",
      "ContentType": "application/pdf"
    }
  ]
}
```

**Response**:
```json
{
  "To": "recipient@example.com",
  "SubmittedAt": "2026-01-02T12:00:00Z",
  "MessageID": "uuid-here",
  "ErrorCode": 0,
  "Message": "OK"
}
```

#### Send Batch Emails
```http
POST /email/batch
Content-Type: application/json
X-Postmark-Server-Token: <your-server-token>

[
  {
    "From": "sender@example.com",
    "To": "recipient1@example.com",
    "Subject": "Email 1",
    "TextBody": "Body 1"
  },
  {
    "From": "sender@example.com",
    "To": "recipient2@example.com",
    "Subject": "Email 2",
    "TextBody": "Body 2"
  }
]
```

**Response**: Array of individual send responses

#### List Templates (Placeholder)
```http
GET /templates
X-Postmark-Server-Token: <your-server-token>
```

**Response**:
```json
{
  "TotalCount": 0,
  "Templates": []
}
```

#### List Webhooks (Placeholder)
```http
GET /webhooks
X-Postmark-Server-Token: <your-server-token>
```

**Response**:
```json
{
  "Webhooks": []
}
```

#### Get Server Info
```http
GET /server
X-Postmark-Server-Token: <your-server-token>
```

**Response**:
```json
{
  "Name": "gomailserver",
  "ApiTokens": []
}
```

### Authentication

**Server Token**:
- Generated via database or admin UI
- Stored as bcrypt hash
- Sent in `X-Postmark-Server-Token` header

**Test Mode**:
- Use token `POSTMARK_API_TEST` for testing
- Validates requests without sending actual email

---

## 🔐 Authentication Methods

### JWT Authentication
- **Used By**: Admin UI, Webmail UI, REST API
- **Header**: `Authorization: Bearer <token>`
- **Obtain**: POST `/api/v1/auth/login`
- **Refresh**: POST `/api/v1/auth/refresh`
- **Expiration**: Configurable (default: 24 hours)

### API Key Authentication
- **Used By**: Programmatic API access
- **Header**: `X-API-Key: <key>`
- **Generate**: Via admin UI or database
- **Use Case**: Server-to-server integration

### PostmarkApp Token
- **Used By**: PostmarkApp API
- **Header**: `X-Postmark-Server-Token: <token>`
- **Generate**: Via database (postmark_servers table)
- **Format**: Bcrypt-hashed for security

### HTTP Basic Auth
- **Used By**: CalDAV, CardDAV (WebDAV services)
- **Port**: `8800` (default)
- **Format**: `Authorization: Basic base64(email:password)`
- **Clients**: Thunderbird, Apple Calendar/Contacts, DAVx5

---

## 🌐 CalDAV/CardDAV Endpoints

### Base Configuration
- **Default Port**: `8800`
- **Environment Variable**: `WEBDAV_PORT`
- **Configuration**: `webdav.port` in YAML config
- **Authentication**: HTTP Basic Auth

### CalDAV (Calendar)
```
https://your-domain.com:8800/.well-known/caldav
https://your-domain.com:8800/calendars/{user_email}/
https://your-domain.com:8800/calendars/{user_email}/{calendar_name}/
```

**Supported Clients**:
- Thunderbird (Lightning addon)
- Apple Calendar (macOS/iOS)
- Android (DAVx5)
- Microsoft Outlook

### CardDAV (Contacts)
```
https://your-domain.com:8800/.well-known/carddav
https://your-domain.com:8800/contacts/{user_email}/
```

**Supported Clients**:
- Apple Contacts (macOS/iOS)
- Android (DAVx5)
- Thunderbird (CardBook addon)

---

## ⚙️ Configuration

### YAML Configuration
```yaml
# gomailserver.yaml

api:
  port: 8980                    # Main API and UI port
  jwt_secret: "your-secret"     # JWT signing secret
  cors_origins:                 # Allowed CORS origins
    - "https://admin.example.com"

webdav:
  enabled: true
  port: 8800                    # CalDAV/CardDAV port

smtp:
  submission_port: 587          # SMTP submission
  relay_port: 25                # SMTP relay (MX)
  smtps_port: 465               # SMTP over TLS

imap:
  port: 143                     # IMAP
  imaps_port: 993               # IMAP over TLS
```

### Environment Variables
```bash
# API Configuration
export API_PORT=8980
export API_JWT_SECRET="your-jwt-secret"

# WebDAV Configuration
export WEBDAV_ENABLED=true
export WEBDAV_PORT=8800

# SMTP Configuration
export SMTP_SUBMISSION_PORT=587
export SMTP_RELAY_PORT=25
export SMTPS_PORT=465

# IMAP Configuration
export IMAP_PORT=143
export IMAPS_PORT=993
```

### Default Ports Summary
| Service | Protocol | Default Port | Config Key | Env Variable |
|---------|----------|--------------|------------|--------------|
| Admin UI | HTTP/HTTPS | 8980 | api.port | API_PORT |
| Webmail UI | HTTP/HTTPS | 8980 | api.port | API_PORT |
| REST API | HTTP/HTTPS | 8980 | api.port | API_PORT |
| PostmarkApp API | HTTP/HTTPS | 8980 | api.port | API_PORT |
| CalDAV/CardDAV | HTTP/HTTPS | 8800 | webdav.port | WEBDAV_PORT |
| SMTP Submission | SMTP | 587 | smtp.submission_port | SMTP_SUBMISSION_PORT |
| SMTP Relay | SMTP | 25 | smtp.relay_port | SMTP_RELAY_PORT |
| SMTPS | SMTP/TLS | 465 | smtp.smtps_port | SMTPS_PORT |
| IMAP | IMAP | 143 | imap.port | IMAP_PORT |
| IMAPS | IMAP/TLS | 993 | imap.imaps_port | IMAPS_PORT |

---

## 🚀 Quick Start Examples

### Access Admin UI
```bash
# First-time setup (creates admin user)
curl http://localhost:8980/api/v1/setup/status

# After setup, access UI in browser
open http://localhost:8980/admin
```

### Access Webmail
```bash
# Login with email credentials
open http://localhost:8980/webmail
```

### Send Email via PostmarkApp API
```bash
curl -X POST http://localhost:8980/email \
  -H "Content-Type: application/json" \
  -H "X-Postmark-Server-Token: your-token" \
  -d '{
    "From": "sender@example.com",
    "To": "recipient@example.com",review the source code and markdown files in .doc_archive/ and create a document called STATUS.md with the true project status
    "Subject": "Test Email",
    "TextBody": "This is a test email"
  }'
```

### Configure CalDAV in Apple Calendar
```
Server: https://your-domain.com:8800
Username: your-email@example.com
Password: your-password
```

### API Health Check
```bash
curl http://localhost:8980/health
# Response: {"status":"ok"}
```

---

## 📝 Notes

### Production Deployment
1. **Use HTTPS**: Configure TLS certificates (Let's Encrypt recommended)
2. **Reverse Proxy**: Consider nginx/Caddy for SSL termination
3. **Firewall**: Open required ports (8980, 8800, 25, 587, 465, 143, 993)
4. **DNS**: Configure MX, A, SPF, DKIM, DMARC records
5. **Security**: Use strong JWT secrets and API keys

### Development
- Default port 8980 allows running without root privileges
- All UIs are embedded in binary (no separate web server needed)
- Hot reload for development: webmail and admin UIs support proxy mode

### API Rate Limiting
- Authentication endpoints: Aggressive rate limiting (brute force protection)
- General API endpoints: Standard rate limiting
- PostmarkApp API: No rate limiting (relies on server token security)

---

**Last Updated**: 2026-01-02
**Documentation Version**: 1.0
**Project Status**: 77% Complete (232/303 tasks)

For full API documentation, see the OpenAPI/Swagger specification (planned).
