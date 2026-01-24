# gomailserver API Documentation

## Overview

The gomailserver REST API provides comprehensive access to mail server management, user administration, webmail functionality, and reputation monitoring.

## Authentication

The API supports two authentication methods:

### JWT Bearer Token
```bash
curl -H "Authorization: Bearer <jwt_token>" \
     https://mail.example.com/api/v1/users
```

### API Key
```bash
curl -H "X-API-Key: <api_key>" \
     https://mail.example.com/api/v1/users
```

## Core Endpoints

### Authentication

#### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "yourpassword"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "name": "Administrator",
    "role": "admin"
  }
}
```

#### Refresh Token
```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Domain Management

#### List Domains
```bash
GET /api/v1/domains
Authorization: Bearer <token>
```

#### Create Domain
```bash
POST /api/v1/domains
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "example.com",
  "dkim_enabled": true,
  "spf_enabled": true,
  "dmarc_enabled": true
}
```

### User Management

#### List Users
```bash
GET /api/v1/users?page=1&limit=50
Authorization: Bearer <token>
```

#### Create User
```bash
POST /api/v1/users
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe"
}
```

### Webmail

#### List Mailboxes
```bash
GET /api/v1/webmail/mailboxes
Authorization: Bearer <token>
```

Response:
```json
{
  "mailboxes": [
    {
      "id": 1,
      "name": "INBOX",
      "special_use": "",
      "subscribed": true
    },
    {
      "id": 2,
      "name": "Drafts",
      "special_use": "\\Drafts",
      "subscribed": true
    }
  ]
}
```

#### Send Message
```bash
POST /api/v1/webmail/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "to": "recipient@example.com",
  "cc": "cc@example.com",
  "subject": "Hello World",
  "body_text": "Plain text message",
  "body_html": "<p>HTML message</p>",
  "attachments": ["/tmp/file.pdf"]
}
```

Response:
```json
{
  "message_id": "20231201120000.123@example.com",
  "status": "queued"
}
```

#### List Messages
```bash
GET /api/v1/webmail/mailboxes/1/messages?page=1&limit=50
Authorization: Bearer <token>
```

### Statistics

#### Get Server Stats
```bash
GET /api/v1/stats
Authorization: Bearer <token>
```

Response:
```json
{
  "total_users": 150,
  "total_domains": 5,
  "total_messages": 12500,
  "queue_size": 25,
  "uptime": "2d 4h 30m"
}
```

### Reputation Management

#### Get Domain Reputation
```bash
GET /api/v1/reputation/domains/example.com/score
Authorization: Bearer <token>
```

Response:
```json
{
  "domain": "example.com",
  "score": 85.5,
  "last_updated": "2023-12-01T12:00:00Z"
}
```

## Error Handling

All API errors follow this format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

## Rate Limiting

The API implements rate limiting:
- Authentication endpoints: 10 requests per IP per 15 minutes
- General API: Configurable per user/domain
- SMTP/IMAP: Separate rate limits

Rate limited responses include:
```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Reset: 1640995200
```

## PostmarkApp Compatibility

The API is compatible with PostmarkApp's REST API for easy migration:

```bash
POST /email
X-Postmark-Server-Token: your-api-token
Content-Type: application/json

{
  "From": "sender@example.com",
  "To": "recipient@example.com",
  "Subject": "Test",
  "TextBody": "Hello World"
}
```

## Webhooks

Configure webhooks for real-time notifications:

```bash
POST /api/v1/webhooks
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Email Delivery Webhook",
  "url": "https://yourapp.com/webhook",
  "secret": "webhook-secret",
  "event_types": ["email.delivered", "email.bounced"]
}
```

## Pagination

List endpoints support pagination:

```bash
GET /api/v1/users?page=2&limit=25
```

Response includes pagination metadata:
```json
{
  "users": [...],
  "pagination": {
    "page": 2,
    "limit": 25,
    "total": 150
  }
}
```

## File Uploads

Message attachments are uploaded separately:

```bash
POST /api/v1/webmail/attachments
Authorization: Bearer <token>
Content-Type: multipart/form-data

# File upload form data
```

Response returns attachment ID for use in message sending.