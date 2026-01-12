# gomailserver - API Reference

**Version**: 1.0  
**Date**: January 11, 2026  
**gomailserver Version**: 0.10.0+  
**Base URL**: `http://your-server:8980/api/v1`

---

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Setup Wizard](#setup-wizard)
4. [Domains](#domains)
5. [Users](#users)
6. [Aliases](#aliases)
7. [Queue](#queue)
8. [Settings](#settings)
9. [Statistics](#statistics)
10. [PGP Keys](#pgp-keys)
11. [Audit Logs](#audit-logs)
12. [Webhooks](#webhooks)
13. [Reputation](#reputation)
14. [Webmail](#webmail)
15. [Error Codes](#error-codes)

---

## Overview

gomailserver provides a comprehensive REST API for system administration, email management, and monitoring.

### Base URL

```
Production: https://mail.example.com/api/v1
Development: http://localhost:8980/api/v1
```

### Authentication

All API endpoints (except `/auth/*` and `/setup/*`) require authentication via:
- **JWT Bearer Token**: `Authorization: Bearer <token>`
- **API Key**: `Authorization: ApiKey <key>`

### Rate Limiting

- **Auth endpoints**: 10 requests per 15 minutes per IP
- **API endpoints**: 100 requests per 60 minutes per IP
- **Per-user limits**: 500 requests per 60 minutes

### Common Response Format

**Success Response**:
```json
{
  "message": "Operation successful",
  "data": { ... }
}
```

**Error Response**:
```json
{
  "error": "Error message",
  "details": {
    "field": "field_name"
  }
}
```

**HTTP Status Codes**:
- `200 OK`: Request successful
- `201 Created`: Resource created
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Invalid or missing authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

---

## Authentication

### POST /auth/login

Authenticate user and receive JWT token.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "totp_code": "123456"
}
```

**Parameters**:
- `email` (string, required): User email address
- `password` (string, required): User password
- `totp_code` (string, optional): 6-digit TOTP code if 2FA enabled

**Response** (200 OK):
```json
{
  "token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "user",
    "domain_id": 1,
    "domain_name": "example.com",
    "totp_enabled": false
  },
  "expires_at": "2026-01-12T00:00:00Z"
}
```

### POST /auth/refresh

Refresh JWT token using refresh token.

**Request Body**:
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Parameters**:
- `refresh_token` (string, required): Refresh token from login

**Response** (200 OK):
```json
{
  "token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_at": "2026-01-13T00:00:00Z"
}
```

---

## Setup Wizard

### GET /setup/status

Check if setup wizard is complete.

**Response** (200 OK):
```json
{
  "completed": false,
  "current_step": "admin",
  "completed_steps": ["system", "domain"]
}
```

### GET /setup/state

Get current setup state.

**Response** (200 OK):
```json
{
  "current_step": "admin",
  "system_config": {
    "hostname": "mail.example.com",
    "domain": "example.com"
  },
  "domain_config": {
    "name": "example.com",
    "status": "active"
  },
  "admin_config": {
    "email": "admin@example.com",
    "created": true
  }
}
```

### POST /setup/admin

Create first admin user.

**Request Body**:
```json
{
  "email": "admin@example.com",
  "password": "SecurePass123!",
  "full_name": "System Administrator"
}
```

**Response** (201 Created):
```json
{
  "message": "Admin user created successfully",
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "full_name": "System Administrator",
    "role": "admin"
  }
}
```

### POST /setup/complete

Mark setup as complete.

**Request Body**:
```json
{}
```

**Response** (200 OK):
```json
{
  "message": "Setup completed successfully"
}
```

---

## Domains

### GET /domains

List all domains.

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 50)
- `search` (string, optional): Search by domain name

**Response** (200 OK):
```json
{
  "domains": [
    {
      "id": 1,
      "name": "example.com",
      "description": "Primary mail domain",
      "status": "active",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1,
    "pages": 1
  }
}
```

### POST /domains

Create new domain.

**Request Body**:
```json
{
  "name": "example.com",
  "description": "Primary mail domain",
  "dkim_enabled": true,
  "spf_enabled": true,
  "dmarc_enabled": true,
  "dane_enabled": false,
  "mtasts_enabled": false
}
```

**Response** (201 Created):
```json
{
  "message": "Domain created successfully",
  "domain": {
    "id": 1,
    "name": "example.com",
    "description": "Primary mail domain",
    "status": "active",
    "dkim_enabled": true,
    "spf_enabled": true,
    "dmarc_enabled": true
  }
}
```

### GET /domains/{id}

Get domain by ID.

**Response** (200 OK):
```json
{
  "id": 1,
  "name": "example.com",
  "description": "Primary mail domain",
  "status": "active",
  "dkim_enabled": true,
  "spf_enabled": true,
  "dmarc_enabled": true,
  "dkim_selector": "default",
  "dkim_key": "MIGfMA0GCs...",
  "spf_record": "v=spf1 ip4:192.168.1.1 -all",
  "dmarc_record": "v=DMARC1; p=none; rua=mailto:dmarc@example.com",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

### PUT /domains/{id}

Update domain.

**Request Body**:
```json
{
  "description": "Updated description",
  "dkim_enabled": false,
  "status": "disabled"
}
```

**Response** (200 OK):
```json
{
  "message": "Domain updated successfully",
  "domain": {
    "id": 1,
    "name": "example.com",
    "description": "Updated description"
  }
}
```

### DELETE /domains/{id}

Delete domain.

**Response** (200 OK):
```json
{
  "message": "Domain deleted successfully"
}
```

### POST /domains/{id}/dkim

Generate DKIM key for domain.

**Request Body**:
```json
{
  "selector": "mail1",
  "key_size": 2048
}
```

**Parameters**:
- `selector` (string, optional): DKIM selector (default: "default")
- `key_size` (integer, optional): Key size in bits - 2048 or 4096 (default: 2048)

**Response** (201 Created):
```json
{
  "message": "DKIM key generated successfully",
  "dkim": {
    "selector": "mail1",
    "public_key": "v=DKIM1; k=rsa; p=MIGfMA0GCs...",
    "dns_record": "mail1._domainkey.example.com.  IN TXT \"v=DKIM1; k=rsa; p=MIGfMA0GCs...\"",
    "created_at": "2026-01-01T00:00:00Z",
    "expires_at": "2026-07-01T00:00:00Z"
  }
}
```

---

## Users

### GET /users

List all users.

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 50)
- `domain_id` (integer, optional): Filter by domain
- `search` (string, optional): Search by email or name

**Response** (200 OK):
```json
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "full_name": "John Doe",
      "role": "user",
      "domain_id": 1,
      "domain_name": "example.com",
      "status": "active",
      "quota_mb": 1024,
      "used_mb": 256,
      "quota_percent": 25,
      "last_login": "2026-01-11T10:30:00Z",
      "created_at": "2026-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1,
    "pages": 1
  }
}
```

### POST /users

Create new user.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe",
  "domain_id": 1,
  "role": "user",
  "quota_mb": 1024
}
```

**Parameters**:
- `email` (string, required): User email address
- `password` (string, required): User password (min 8 characters)
- `full_name` (string, required): Full name
- `domain_id` (integer, required): Domain ID
- `role` (string, optional): Role - "user" or "admin" (default: "user")
- `quota_mb` (integer, optional): Quota in MB (0 = unlimited)

**Response** (201 Created):
```json
{
  "message": "User created successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "user",
    "domain_id": 1,
    "quota_mb": 1024
  }
}
```

### GET /users/{id}

Get user by ID.

**Response** (200 OK):
```json
{
  "id": 1,
  "email": "user@example.com",
  "full_name": "John Doe",
  "role": "user",
  "domain_id": 1,
  "domain_name": "example.com",
  "status": "active",
  "quota_mb": 1024,
  "used_mb": 256,
  "quota_percent": 25,
  "totp_enabled": false,
  "totp_secret": "",
  "totp_recovery_codes": [],
  "last_login": "2026-01-11T10:30:00Z",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-11T10:30:00Z"
}
```

### PUT /users/{id}

Update user.

**Request Body**:
```json
{
  "full_name": "Jane Doe",
  "quota_mb": 2048,
  "status": "active"
}
```

**Parameters**:
- `full_name` (string, optional): Full name
- `quota_mb` (integer, optional): Quota in MB
- `status` (string, optional): "active" or "disabled"
- `totp_enabled` (boolean, optional): Enable/disable 2FA

**Response** (200 OK):
```json
{
  "message": "User updated successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "Jane Doe",
    "quota_mb": 2048
  }
}
```

### DELETE /users/{id}

Delete user.

**Response** (200 OK):
```json
{
  "message": "User deleted successfully"
}
```

### POST /users/{id}/password

Reset user password.

**Request Body**:
```json
{
  "new_password": "NewSecurePass123!"
}
```

**Parameters**:
- `new_password` (string, required): New password (min 8 characters)

**Response** (200 OK):
```json
{
  "message": "Password reset successfully"
}
```

---

## Aliases

### GET /aliases

List all aliases.

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 50)
- `domain_id` (integer, optional): Filter by domain

**Response** (200 OK):
```json
{
  "aliases": [
    {
      "id": 1,
      "source": "support@example.com",
      "destination": "user@example.com",
      "domain_id": 1,
      "domain_name": "example.com",
      "active": true,
      "created_at": "2026-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1,
    "pages": 1
  }
}
```

### POST /aliases

Create new alias.

**Request Body**:
```json
{
  "source": "support@example.com",
  "destination": "user@example.com",
  "domain_id": 1
}
```

**Parameters**:
- `source` (string, required): Alias email address
- `destination` (string, required): Destination email address
- `domain_id` (integer, required): Domain ID

**Response** (201 Created):
```json
{
  "message": "Alias created successfully",
  "alias": {
    "id": 1,
    "source": "support@example.com",
    "destination": "user@example.com",
    "active": true
  }
}
```

### GET /aliases/{id}

Get alias by ID.

**Response** (200 OK):
```json
{
  "id": 1,
  "source": "support@example.com",
  "destination": "user@example.com",
  "domain_id": 1,
  "domain_name": "example.com",
  "active": true,
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

### DELETE /aliases/{id}

Delete alias.

**Response** (200 OK):
```json
{
  "message": "Alias deleted successfully"
}
```

---

## Queue

### GET /queue

List queued messages.

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 50)
- `status` (string, optional): Filter by status (pending, retrying, failed, sent)

**Response** (200 OK):
```json
{
  "messages": [
    {
      "id": 1,
      "message_id": "<message-id@example.com>",
      "sender": "user@example.com",
      "recipients": "recipient@example.com",
      "status": "pending",
      "retry_count": 0,
      "max_retries": 5,
      "next_retry": "2026-01-11T10:35:00Z",
      "created_at": "2026-01-11T10:30:00Z",
      "updated_at": "2026-01-11T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1,
    "pages": 1
  }
}
```

### GET /queue/{id}

Get queued message by ID.

**Response** (200 OK):
```json
{
  "id": 1,
  "message_id": "<message-id@example.com>",
  "sender": "user@example.com",
  "recipients": "recipient@example.com",
  "subject": "Test email",
  "status": "pending",
  "retry_count": 0,
  "max_retries": 5,
  "next_retry": "2026-01-11T10:35:00Z",
  "error_message": "",
  "created_at": "2026-01-11T10:30:00Z",
  "updated_at": "2026-01-11T10:30:00Z"
}
```

### POST /queue/{id}/retry

Retry queued message immediately.

**Response** (200 OK):
```json
{
  "message": "Message retry queued successfully"
}
```

### DELETE /queue/{id}

Delete queued message.

**Response** (200 OK):
```json
{
  "message": "Message deleted from queue successfully"
}
```

---

## Settings

### GET /settings/server

Get server settings.

**Response** (200 OK):
```json
{
  "hostname": "mail.example.com",
  "domain": "example.com",
  "listen_addr": "0.0.0.0:25",
  "smtp_enabled": true,
  "smtp_port": 25,
  "submission_port": 587,
  "smtps_port": 465,
  "imap_enabled": true,
  "imap_port": 143,
  "imaps_port": 993,
  "http_enabled": true,
  "http_port": 8980
}
```

### PUT /settings/server

Update server settings.

**Request Body**:
```json
{
  "hostname": "mail.example.com",
  "smtp_enabled": true,
  "imap_enabled": true,
  "http_port": 8980
}
```

**Response** (200 OK):
```json
{
  "message": "Server settings updated successfully. Server restart required for changes to take effect.",
  "server": {
    "hostname": "mail.example.com",
    "smtp_enabled": true
  }
}
```

### GET /settings/security

Get security settings.

**Response** (200 OK):
```json
{
  "dkim_enabled": true,
  "spf_enabled": true,
  "dmarc_enabled": true,
  "dane_enabled": false,
  "mtasts_enabled": false,
  "clamav_enabled": true,
  "clamav_socket": "/var/run/clamav/clamd.sock",
  "spamassassin_enabled": true,
  "spamassassin_host": "localhost",
  "spamassassin_port": 783,
  "totp_enabled": false,
  "jwt_secret": "",
  "api_admin_token": ""
}
```

### PUT /settings/security

Update security settings.

**Request Body**:
```json
{
  "dkim_enabled": true,
  "spf_enabled": true,
  "dmarc_enabled": true,
  "clamav_enabled": true,
  "clamav_socket": "/var/run/clamav/clamd.sock",
  "spamassassin_enabled": true,
  "spamassassin_host": "localhost",
  "spamassassin_port": 783
}
```

**Response** (200 OK):
```json
{
  "message": "Security settings updated successfully. Server restart may be required for some changes to take effect.",
  "security": {
    "dkim_enabled": true,
    "clamav_enabled": true
  }
}
```

### GET /settings/tls

Get TLS settings.

**Response** (200 OK):
```json
{
  "acme_enabled": true,
  "acme_email": "admin@example.com",
  "acme_provider": "cloudflare",
  "acme_api_token": "",
  "cert_file": "/etc/ssl/certs/gomailserver.crt",
  "key_file": "/etc/ssl/private/gomailserver.key"
}
```

### PUT /settings/tls

Update TLS settings.

**Request Body**:
```json
{
  "acme_enabled": true,
  "acme_email": "admin@example.com",
  "acme_provider": "cloudflare",
  "acme_api_token": "your_cloudflare_api_token",
  "cert_file": "/etc/ssl/certs/gomailserver.crt",
  "key_file": "/etc/ssl/private/gomailserver.key"
}
```

**Response** (200 OK):
```json
{
  "message": "TLS settings updated successfully. Server restart required for changes to take effect.",
  "tls": {
    "acme_enabled": true,
    "acme_provider": "cloudflare"
  }
}
```

---

## Statistics

### GET /stats/dashboard

Get dashboard statistics.

**Response** (200 OK):
```json
{
  "domains": {
    "total": 5,
    "active": 4,
    "disabled": 1
  },
  "users": {
    "total": 25,
    "active": 23,
    "disabled": 2
  },
  "emails": {
    "today": 150,
    "week": 850,
    "month": 3200
  },
  "queue": {
    "pending": 5,
    "retrying": 2,
    "failed": 1
  },
  "system": {
    "uptime": "99.9%",
    "cpu_usage": "45%",
    "memory_usage": "62%",
    "disk_usage": "38%"
  }
}
```

### GET /stats/domains/{id}

Get statistics for specific domain.

**Response** (200 OK):
```json
{
  "domain_id": 1,
  "domain_name": "example.com",
  "users": {
    "total": 5,
    "active": 5
  },
  "emails": {
    "sent_today": 50,
    "sent_week": 300,
    "received_today": 75,
    "received_week": 450
  },
  "quota": {
    "total_mb": 5120,
    "used_mb": 1280,
    "free_mb": 3840,
    "usage_percent": 25
  }
}
```

### GET /stats/users/{id}

Get statistics for specific user.

**Response** (200 OK):
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "emails": {
    "sent": 150,
    "received": 300,
    "drafts": 5,
    "spam": 2
  },
  "storage": {
    "quota_mb": 1024,
    "used_mb": 256,
    "free_mb": 768,
    "usage_percent": 25
  },
  "last_login": "2026-01-11T10:30:00Z"
}
```

---

## PGP Keys

### POST /pgp/keys

Import PGP public key.

**Request Body**:
```json
{
  "user_id": 1,
  "key_data": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n..."
}
```

**Parameters**:
- `user_id` (integer, required): User ID
- `key_data` (string, required): PGP public key in ASCII armor format

**Response** (201 Created):
```json
{
  "message": "PGP key imported successfully",
  "key": {
    "id": 1,
    "user_id": 1,
    "key_id": "ABCDEF123...",
    "email": "user@example.com",
    "created_at": "2026-01-11T10:30:00Z",
    "primary": false
  }
}
```

### GET /pgp/keys/users/{user_id}

List PGP keys for user.

**Response** (200 OK):
```json
{
  "keys": [
    {
      "id": 1,
      "user_id": 1,
      "key_id": "ABCDEF123...",
      "email": "user@example.com",
      "created_at": "2026-01-11T10:30:00Z",
      "primary": true
    }
  ]
}
```

### GET /pgp/keys/{id}

Get PGP key by ID.

**Response** (200 OK):
```json
{
  "id": 1,
  "user_id": 1,
  "key_id": "ABCDEF123...",
  "email": "user@example.com",
  "key_data": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n...",
  "created_at": "2026-01-11T10:30:00Z",
  "primary": true
}
```

### POST /pgp/keys/{id}/primary

Set PGP key as primary for user.

**Response** (200 OK):
```json
{
  "message": "PGP key set as primary successfully"
}
```

### DELETE /pgp/keys/{id}

Delete PGP key.

**Response** (200 OK):
```json
{
  "message": "PGP key deleted successfully"
}
```

---

## Audit Logs

### GET /audit/logs

List audit logs.

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 50)
- `user_id` (integer, optional): Filter by user
- `action` (string, optional): Filter by action type
- `start_date` (string, optional): Filter by start date (ISO 8601)
- `end_date` (string, optional): Filter by end date (ISO 8601)

**Response** (200 OK):
```json
{
  "logs": [
    {
      "id": 1,
      "user_id": 1,
      "user_email": "user@example.com",
      "action": "login",
      "details": "User logged in successfully",
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2026-01-11T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 100,
    "pages": 2
  }
}
```

### GET /audit/stats

Get audit statistics.

**Response** (200 OK):
```json
{
  "total_events": 1000,
  "unique_users": 25,
  "by_action": {
    "login": 500,
    "logout": 450,
    "password_change": 25,
    "user_create": 25
  },
  "by_user": [
    {
      "user_id": 1,
      "user_email": "user@example.com",
      "event_count": 75
    }
  ],
  "last_24_hours": {
    "total": 50,
    "unique_users": 10
  }
}
```

---

## Webhooks

### GET /webhooks

List webhooks.

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 50)
- `event_type` (string, optional): Filter by event type

**Response** (200 OK):
```json
{
  "webhooks": [
    {
      "id": 1,
      "url": "https://example.com/webhook",
      "events": ["email.sent", "email.failed"],
      "secret": "webhook_secret",
      "active": true,
      "created_at": "2026-01-01T00:00:00Z",
      "delivery_count": 150,
      "failure_count": 2,
      "last_delivery": "2026-01-11T10:25:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1,
    "pages": 1
  }
}
```

### POST /webhooks

Create webhook.

**Request Body**:
```json
{
  "url": "https://example.com/webhook",
  "events": ["email.sent", "email.failed"],
  "secret": "webhook_secret"
}
```

**Parameters**:
- `url` (string, required): Webhook URL (must be HTTPS)
- `events` (array, required): Event types to subscribe to
- `secret` (string, optional): HMAC secret for webhook verification

**Response** (201 Created):
```json
{
  "message": "Webhook created successfully",
  "webhook": {
    "id": 1,
    "url": "https://example.com/webhook",
    "events": ["email.sent", "email.failed"],
    "active": true
  }
}
```

### GET /webhooks/{id}

Get webhook by ID.

**Response** (200 OK):
```json
{
  "id": 1,
  "url": "https://example.com/webhook",
  "events": ["email.sent", "email.failed"],
  "secret": "webhook_secret",
  "active": true,
  "created_at": "2026-01-01T00:00:00Z",
  "delivery_count": 150,
  "failure_count": 2
}
```

### PUT /webhooks/{id}

Update webhook.

**Request Body**:
```json
{
  "url": "https://example.com/webhook-updated",
  "active": false
}
```

**Response** (200 OK):
```json
{
  "message": "Webhook updated successfully",
  "webhook": {
    "id": 1,
    "url": "https://example.com/webhook-updated",
    "active": false
  }
}
```

### DELETE /webhooks/{id}

Delete webhook.

**Response** (200 OK):
```json
{
  "message": "Webhook deleted successfully"
}
```

### POST /webhooks/{id}/test

Test webhook delivery.

**Request Body**:
```json
{
  "test_event": "email.sent"
}
```

**Response** (200 OK):
```json
{
  "message": "Webhook test sent successfully",
  "test_result": {
    "delivered": true,
    "response_code": 200,
    "response_body": "Received",
    "duration_ms": 150
  }
}
```

---

## Reputation

### POST /reputation/audit/{domain}

Run reputation audit for domain.

**Response** (200 OK):
```json
{
  "domain": "example.com",
  "audit": {
    "spf": {
      "status": "pass",
      "score": 100,
      "message": "SPF record valid"
    },
    "dkim": {
      "status": "pass",
      "score": 100,
      "message": "DKIM record valid"
    },
    "dmarc": {
      "status": "pass",
      "score": 100,
      "message": "DMARC record valid"
    },
    "overall_score": 95
  }
}
```

### GET /reputation/scores

List reputation scores for all domains.

**Response** (200 OK):
```json
{
  "scores": [
    {
      "domain": "example.com",
      "score": 85,
      "status": "good",
      "deliveries_90d": 5000,
      "bounces_90d": 250,
      "complaints_90d": 5,
      "deferrals_90d": 100,
      "circuit_breaker": {
        "active": false,
        "trigger": ""
      }
    }
  ]
}
```

### GET /reputation/scores/{domain}

Get reputation score for specific domain.

**Response** (200 OK):
```json
{
  "domain": "example.com",
  "score": 85,
  "status": "good",
  "metrics": {
    "deliveries_90d": 5000,
    "bounces_90d": 250,
    "complaints_90d": 5,
    "deferrals_90d": 100
  },
  "circuit_breaker": {
    "active": false,
    "trigger": ""
  },
  "provider_limits": {
    "gmail": {
      "hourly_limit": 100,
      "daily_limit": 500,
      "hourly_usage": 25,
      "daily_usage": 150
    }
  }
}
```

### GET /reputation/circuit-breakers

List circuit breakers.

**Response** (200 OK):
```json
{
  "circuit_breakers": [
    {
      "id": 1,
      "domain": "example.com",
      "provider": "gmail",
      "active": true,
      "trigger": "complaint_rate_high",
      "triggered_at": "2026-01-10T15:00:00Z",
      "auto_resume_at": "2026-01-11T15:00:00Z"
    }
  ]
}
```

### GET /reputation/circuit-breakers/{domain}/history

Get circuit breaker history for domain.

**Response** (200 OK):
```json
{
  "history": [
    {
      "id": 1,
      "domain": "example.com",
      "provider": "gmail",
      "trigger": "complaint_rate_high",
      "triggered_at": "2026-01-10T15:00:00Z",
      "resolved_at": "2026-01-11T15:00:00Z",
      "auto_resumed": true
    }
  ]
}
```

### GET /reputation/alerts

List reputation alerts.

**Query Parameters**:
- `status` (string, optional): Filter by status (pending, acknowledged, resolved)

**Response** (200 OK):
```json
{
  "alerts": [
    {
      "id": 1,
      "domain": "example.com",
      "severity": "warning",
      "message": "Reputation score dropped below 70",
      "status": "pending",
      "created_at": "2026-01-11T10:00:00Z"
    }
  ]
}
```

---

## Webmail

### GET /webmail/mailboxes

List user's mailboxes.

**Response** (200 OK):
```json
{
  "mailboxes": [
    {
      "id": 1,
      "name": "INBOX",
      "type": "inbox",
      "total": 100,
      "unread": 15,
      "size_bytes": 10485760
    },
    {
      "id": 2,
      "name": "Drafts",
      "type": "drafts",
      "total": 5,
      "unread": 0,
      "size_bytes": 5242880
    },
    {
      "id": 3,
      "name": "Sent",
      "type": "sent",
      "total": 150,
      "unread": 0,
      "size_bytes": 157286400
    }
  ]
}
```

### GET /webmail/mailboxes/{id}/messages

List messages in mailbox.

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 50)
- `folder` (string, optional): Filter by folder
- `search` (string, optional): Search query
- `sort` (string, optional): Sort field (date, subject, from)

**Response** (200 OK):
```json
{
  "messages": [
    {
      "id": 1,
      "message_id": "<message-id@example.com>",
      "subject": "Test email",
      "from": "sender@example.com",
      "to": ["user@example.com"],
      "cc": [],
      "bcc": [],
      "date": "2026-01-11T10:00:00Z",
      "flags": {
        "seen": false,
        "flagged": false,
        "answered": false,
        "draft": false
      },
      "has_attachments": true,
      "size_bytes": 10240
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 100,
    "pages": 2
  }
}
```

### GET /webmail/messages/{id}

Get message by ID.

**Response** (200 OK):
```json
{
  "id": 1,
  "message_id": "<message-id@example.com>",
  "subject": "Test email",
  "from": {
    "name": "Sender Name",
    "email": "sender@example.com"
  },
  "to": [
    {
      "name": "User Name",
      "email": "user@example.com"
    }
  ],
  "cc": [],
  "bcc": [],
  "date": "2026-01-11T10:00:00Z",
  "flags": {
    "seen": false,
    "flagged": false,
    "answered": false
    "draft": false
  },
  "body": {
    "plain": "Plain text content",
    "html": "<html>HTML content</html>"
  },
  "attachments": [
    {
      "id": 1,
      "filename": "document.pdf",
      "size_bytes": 102400,
      "content_type": "application/pdf"
    }
  ]
}
```

### POST /webmail/messages

Send new message.

**Request Body**:
```json
{
  "from": "user@example.com",
  "to": ["recipient@example.com"],
  "cc": ["cc@example.com"],
  "bcc": ["bcc@example.com"],
  "subject": "Test email",
  "body": "Email body content",
  "is_html": true,
  "attachments": [],
  "priority": "normal",
  "request_receipt": false
}
```

**Parameters**:
- `from` (string, required): Sender email address
- `to` (array, required): Recipients
- `cc` (array, optional): CC recipients
- `bcc` (array, optional): BCC recipients
- `subject` (string, required): Email subject
- `body` (string, required): Email body
- `is_html` (boolean, optional): HTML content (default: false)
- `attachments` (array, optional): Attachment IDs or paths
- `priority` (string, optional): "normal", "high", "low" (default: "normal")
- `request_receipt` (boolean, optional): Request read receipt (default: false)

**Response** (201 Created):
```json
{
  "message": "Message sent successfully",
  "message_id": "<message-id@example.com>",
  "queue_id": 1,
  "queued_at": "2026-01-11T10:30:00Z"
}
```

### DELETE /webmail/messages/{id}

Delete message.

**Response** (200 OK):
```json
{
  "message": "Message deleted successfully"
}
```

### POST /webmail/messages/{id}/move

Move message to mailbox.

**Request Body**:
```json
{
  "target_mailbox_id": 3
}
```

**Response** (200 OK):
```json
{
  "message": "Message moved successfully"
}
```

### POST /webmail/messages/{id}/flags

Update message flags.

**Request Body**:
```json
{
  "seen": true,
  "flagged": false,
  "answered": true
}
```

**Response** (200 OK):
```json
{
  "message": "Message flags updated successfully"
}
```

### GET /webmail/search

Search messages.

**Query Parameters**:
- `q` (string, required): Search query
- `mailbox_id` (integer, optional): Search specific mailbox
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 50)

**Response** (200 OK):
```json
{
  "results": [
    {
      "id": 1,
      "subject": "Test email",
      "from": "sender@example.com",
      "date": "2026-01-11T10:00:00Z",
      "snippet": "...",
      "score": 0.95
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 10,
    "pages": 1
  }
}
```

### GET /webmail/attachments/{id}

Download attachment.

**Response**: Binary file content

### POST /webmail/drafts

Save draft message.

**Request Body**:
```json
{
  "mailbox_id": 2,
  "subject": "Draft email",
  "body": "Draft content",
  "is_html": true
}
```

**Response** (201 Created):
```json
{
  "message": "Draft saved successfully",
  "draft_id": 1
}
```

### GET /webmail/drafts

List drafts.

**Response** (200 OK):
```json
{
  "drafts": [
    {
      "id": 1,
      "subject": "Draft email",
      "created_at": "2026-01-11T10:00:00Z",
      "updated_at": "2026-01-11T10:30:00Z"
    }
  ]
}
```

### GET /webmail/drafts/{id}

Get draft by ID.

**Response** (200 OK):
```json
{
  "id": 1,
  "subject": "Draft email",
  "body": "Draft content",
  "is_html": true,
  "from": "user@example.com",
  "to": [],
  "created_at": "2026-01-11T10:00:00Z",
  "updated_at": "2026-01-11T10:30:00Z"
}
```

### DELETE /webmail/drafts/{id}

Delete draft.

**Response** (200 OK):
```json
{
  "message": "Draft deleted successfully"
}
```

---

## Error Codes

### Common Error Responses

**400 Bad Request**:
```json
{
  "error": "Invalid request",
  "details": {
    "field": "email",
    "message": "Invalid email format"
  }
}
```

**401 Unauthorized**:
```json
{
  "error": "Invalid or missing authentication",
  "details": {
    "message": "JWT token expired"
  }
}
```

**403 Forbidden**:
```json
{
  "error": "Insufficient permissions",
  "details": {
    "message": "Admin access required"
  }
}
```

**404 Not Found**:
```json
{
  "error": "Resource not found",
  "details": {
    "resource": "user",
    "id": 999
  }
}
```

**409 Conflict**:
```json
{
  "error": "Resource conflict",
  "details": {
    "field": "email",
    "message": "Email address already exists"
  }
}
```

**422 Unprocessable Entity**:
```json
{
  "error": "Validation failed",
  "details": {
    "field": "password",
    "message": "Password must be at least 8 characters"
  }
}
```

**429 Too Many Requests**:
```json
{
  "error": "Rate limit exceeded",
  "details": {
    "limit": 100,
    "window": "60s",
    "retry_after": "2026-01-11T10:31:00Z"
  }
}
```

**500 Internal Server Error**:
```json
{
  "error": "Internal server error",
  "details": {
    "message": "An unexpected error occurred"
  }
}
```

---

## SDK Examples

### cURL

**Login**:
```bash
curl -X POST http://localhost:8980/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password"
  }'
```

**List Domains**:
```bash
curl -X GET http://localhost:8980/api/v1/domains \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Send Email**:
```bash
curl -X POST http://localhost:8980/api/v1/webmail/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["recipient@example.com"],
    "subject": "Test email",
    "body": "Hello World"
  }'
```

### JavaScript (Fetch API)

**Login**:
```javascript
const response = await fetch('http://localhost:8980/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password'
  })
});

const data = await response.json();
const token = data.token;
```

**List Users**:
```javascript
const response = await fetch('http://localhost:8980/api/v1/users', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const data = await response.json();
console.log(data.users);
```

### Python (requests)

**Login**:
```python
import requests

response = requests.post(
    'http://localhost:8980/api/v1/auth/login',
    json={
        'email': 'user@example.com',
        'password': 'password'
    }
)

data = response.json()
token = data['token']
```

**Create User**:
```python
response = requests.post(
    'http://localhost:8980/api/v1/users',
    headers={'Authorization': f'Bearer {token}'},
    json={
        'email': 'newuser@example.com',
        'password': 'SecurePass123!',
        'full_name': 'New User',
        'domain_id': 1
    }
)

data = response.json()
print(data)
```

---

**End of API Reference**

For the latest API documentation and updates, visit:
https://github.com/btafoya/gomailserver
