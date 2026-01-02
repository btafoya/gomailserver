# gomailserver Admin API Documentation

## Overview

The gomailserver Admin API provides RESTful endpoints for managing domain configuration and security settings. All security policies are stored in SQLite on a per-domain basis, enabling multi-tenant hosting with different security policies for each domain.

## Authentication

The API uses Bearer token authentication for all admin endpoints (except `/health`).

### Configuration

Set the admin token in your configuration file:

```yaml
api:
  port: 8080
  admin_token: "your-secret-token-here"
```

Or via environment variable:
```bash
export API_ADMIN_TOKEN="your-secret-token-here"
```

### Making Authenticated Requests

Include the bearer token in the `Authorization` header:

```bash
curl -H "Authorization: Bearer your-secret-token-here" \
     http://localhost:8080/api/domains
```

**Development Mode**: If no `admin_token` is configured, the API runs without authentication (NOT RECOMMENDED FOR PRODUCTION).

## Base URL

Default: `http://localhost:8080`

Configure via:
```yaml
api:
  port: 8080
```

Or environment variable: `API_PORT=8080`

## Endpoints

### Health Check

```
GET /health
```

Health check endpoint (no authentication required).

**Response:**
```json
{
  "status": "healthy"
}
```

---

### List All Domains

```
GET /api/domains
```

Returns all configured domains with their security settings.

**Response:**
```json
{
  "domains": [
    {
      "id": 1,
      "name": "example.com",
      "status": "active",
      "dkim_signing_enabled": true,
      "spf_enabled": true,
      ...
    }
  ],
  "count": 1
}
```

---

### Get Domain

```
GET /api/domains/{name}
```

Retrieves a specific domain by name.

**Parameters:**
- `name` (path): Domain name (e.g., "example.com")

**Response:**
```json
{
  "id": 1,
  "name": "example.com",
  "status": "active",
  "max_users": 0,
  "max_mailbox_size": 0,
  "default_quota": 1073741824,
  "dkim_signing_enabled": true,
  "dkim_verify_enabled": true,
  "dkim_key_size": 2048,
  "dkim_key_type": "rsa",
  ...
}
```

---

### Create Domain

```
POST /api/domains
```

Creates a new domain. You can either create from the default template or provide custom settings.

**Request Body (From Template):**
```json
{
  "name": "newdomain.com",
  "use_default_template": true
}
```

**Request Body (Custom Settings):**
```json
{
  "name": "newdomain.com",
  "use_default_template": false,
  "domain": {
    "status": "active",
    "max_users": 100,
    "dkim_signing_enabled": true,
    "spf_enabled": true,
    ...
  }
}
```

**Response:** `201 Created`
```json
{
  "id": 2,
  "name": "newdomain.com",
  "status": "active",
  ...
}
```

---

### Update Domain

```
PUT /api/domains/{name}
```

Updates an existing domain's settings.

**Parameters:**
- `name` (path): Domain name

**Request Body:**
```json
{
  "status": "active",
  "max_users": 200,
  "dkim_signing_enabled": false,
  ...
}
```

**Response:**
```json
{
  "id": 1,
  "name": "example.com",
  "status": "active",
  ...
}
```

**Note:** The domain name cannot be changed. `id`, `name`, and `created_at` are preserved.

---

### Delete Domain

```
DELETE /api/domains/{name}
```

Deletes a domain.

**Parameters:**
- `name` (path): Domain name

**Response:** `204 No Content`

**Note:** The default template domain (`_default`) cannot be deleted.

---

### Get Domain Security Configuration

```
GET /api/domains/{name}/security
```

Retrieves only the security configuration for a domain.

**Parameters:**
- `name` (path): Domain name

**Response:**
```json
{
  "domain": "example.com",
  "dkim": {
    "signing_enabled": true,
    "verify_enabled": true,
    "key_size": 2048,
    "key_type": "rsa",
    "headers_to_sign": "[\"From\",\"To\",\"Subject\",\"Date\",\"Message-ID\"]",
    "selector": "default",
    "public_key": "..."
  },
  "spf": {
    "enabled": true,
    "dns_server": "8.8.8.8:53",
    "dns_timeout": 5,
    "max_lookups": 10,
    "fail_action": "reject",
    "softfail_action": "accept",
    "record": "v=spf1 mx -all"
  },
  "dmarc": {
    "enabled": true,
    "dns_server": "8.8.8.8:53",
    "dns_timeout": 5,
    "report_enabled": false,
    "report_email": "",
    "policy": "p=quarantine"
  },
  "clamav": {
    "enabled": true,
    "max_scan_size": 52428800,
    "virus_action": "reject",
    "fail_action": "accept"
  },
  "spam": {
    "enabled": true,
    "reject_score": 10.0,
    "quarantine_score": 5.0,
    "learning_enabled": true
  },
  "greylist": {
    "enabled": true,
    "delay_minutes": 5,
    "expiry_days": 30,
    "cleanup_interval": 3600,
    "whitelist_after": 3
  },
  "rate_limit": {
    "enabled": true,
    "smtp_per_ip": "{\"count\":100,\"window_minutes\":60}",
    "smtp_per_user": "{\"count\":500,\"window_minutes\":60}",
    "smtp_per_domain": "{\"count\":1000,\"window_minutes\":60}",
    "auth_per_ip": "{\"count\":10,\"window_minutes\":15}",
    "imap_per_user": "{\"count\":1000,\"window_minutes\":60}",
    "cleanup_interval": 300
  },
  "auth": {
    "totp_enforced": false,
    "brute_force_enabled": true,
    "brute_force_threshold": 5,
    "brute_force_window_minutes": 15,
    "brute_force_block_minutes": 60,
    "ip_blacklist_enabled": true,
    "cleanup_interval": 3600
  }
}
```

---

### Update Domain Security Configuration

```
PUT /api/domains/{name}/security
```

Updates security configuration for a domain. You can update specific sections without affecting others.

**Parameters:**
- `name` (path): Domain name

**Request Body (Partial Update Example):**
```json
{
  "dkim": {
    "signing_enabled": false
  },
  "spam": {
    "reject_score": 15.0,
    "quarantine_score": 7.0
  },
  "auth": {
    "totp_enforced": true
  }
}
```

**Response:**
```json
{
  "message": "security configuration updated successfully",
  "domain": "example.com"
}
```

**Note:** Only the fields you specify are updated. All other settings remain unchanged.

---

### Get Default Template

```
GET /api/domains/_default
```

Retrieves the default domain template used for creating new domains.

**Response:**
```json
{
  "id": 1,
  "name": "_default",
  "status": "active",
  "dkim_signing_enabled": true,
  "dkim_verify_enabled": true,
  "dkim_key_size": 2048,
  ...
}
```

---

### Update Default Template

```
PUT /api/domains/_default
```

Updates the default template. New domains created after this will inherit the updated settings.

**Request Body:**
```json
{
  "dkim_key_size": 4096,
  "spam_reject_score": 15.0,
  "auth_totp_enforced": true,
  ...
}
```

**Response:**
```json
{
  "message": "default template updated successfully"
}
```

**Note:** Updating the default template does not affect existing domains.

---

## Security Configuration Reference

### DKIM Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `signing_enabled` | boolean | true | Enable DKIM signing for outbound mail |
| `verify_enabled` | boolean | true | Enable DKIM verification for inbound mail |
| `key_size` | integer | 2048 | RSA key size (2048 or 4096) |
| `key_type` | string | "rsa" | Key type ("rsa" or "ed25519") |
| `headers_to_sign` | string | JSON array | Headers to include in DKIM signature |

### SPF Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | true | Enable SPF validation |
| `dns_server` | string | "8.8.8.8:53" | DNS server for SPF lookups |
| `dns_timeout` | integer | 5 | DNS lookup timeout (seconds) |
| `max_lookups` | integer | 10 | Maximum DNS lookups per SPF check |
| `fail_action` | string | "reject" | Action on SPF fail (reject/quarantine/accept/tag) |
| `softfail_action` | string | "accept" | Action on SPF softfail |

### DMARC Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | true | Enable DMARC policy enforcement |
| `dns_server` | string | "8.8.8.8:53" | DNS server for DMARC lookups |
| `dns_timeout` | integer | 5 | DNS lookup timeout (seconds) |
| `report_enabled` | boolean | false | Enable DMARC reporting |
| `report_email` | string | "" | Email address for DMARC reports |

### ClamAV Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | true | Enable virus scanning |
| `max_scan_size` | integer | 52428800 | Maximum message size to scan (bytes) |
| `virus_action` | string | "reject" | Action on virus detection (reject/quarantine/tag) |
| `fail_action` | string | "accept" | Action on scan failure |

### SpamAssassin Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | true | Enable spam scanning |
| `reject_score` | float | 10.0 | Score threshold for rejection |
| `quarantine_score` | float | 5.0 | Score threshold for quarantine |
| `learning_enabled` | boolean | true | Enable Bayes learning from user actions |

### Greylisting Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | true | Enable greylisting |
| `delay_minutes` | integer | 5 | Initial delay period (minutes) |
| `expiry_days` | integer | 30 | Triplet expiry period (days) |
| `cleanup_interval` | integer | 3600 | Cleanup interval (seconds) |
| `whitelist_after` | integer | 3 | Successful deliveries before whitelisting |

### Rate Limiting Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | true | Enable rate limiting |
| `smtp_per_ip` | string | JSON | SMTP rate limit per IP |
| `smtp_per_user` | string | JSON | SMTP rate limit per user |
| `smtp_per_domain` | string | JSON | SMTP rate limit per domain |
| `auth_per_ip` | string | JSON | Auth rate limit per IP |
| `imap_per_user` | string | JSON | IMAP rate limit per user |
| `cleanup_interval` | integer | 300 | Cleanup interval (seconds) |

**Rate Limit JSON Format:**
```json
{
  "count": 100,
  "window_minutes": 60
}
```

### Authentication Security Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `totp_enforced` | boolean | false | Require TOTP 2FA for all users |
| `brute_force_enabled` | boolean | true | Enable brute force protection |
| `brute_force_threshold` | integer | 5 | Failed attempts before blocking |
| `brute_force_window_minutes` | integer | 15 | Time window for counting failures |
| `brute_force_block_minutes` | integer | 60 | Block duration after threshold |
| `ip_blacklist_enabled` | boolean | true | Enable IP blacklisting |
| `cleanup_interval` | integer | 3600 | Cleanup interval (seconds) |

---

## Error Responses

All errors return JSON with the following format:

```json
{
  "error": "Error message description",
  "status": 400
}
```

### Common HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content (successful deletion) |
| 400 | Bad Request (invalid input) |
| 401 | Unauthorized (missing or invalid token) |
| 403 | Forbidden (operation not allowed) |
| 404 | Not Found (domain doesn't exist) |
| 409 | Conflict (domain already exists) |
| 500 | Internal Server Error |

---

## Examples

### Create a New Domain from Template

```bash
curl -X POST http://localhost:8080/api/domains \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "newcompany.com",
    "use_default_template": true
  }'
```

### Enable TOTP for a Specific Domain

```bash
curl -X PUT http://localhost:8080/api/domains/example.com/security \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json" \
  -d '{
    "auth": {
      "totp_enforced": true
    }
  }'
```

### Adjust Spam Scores

```bash
curl -X PUT http://localhost:8080/api/domains/example.com/security \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json" \
  -d '{
    "spam": {
      "reject_score": 15.0,
      "quarantine_score": 7.0
    }
  }'
```

### Disable DKIM Signing

```bash
curl -X PUT http://localhost:8080/api/domains/example.com/security \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json" \
  -d '{
    "dkim": {
      "signing_enabled": false
    }
  }'
```

### Update Default Template for New Domains

```bash
curl -X PUT http://localhost:8080/api/domains/_default \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json" \
  -d '{
    "dkim_key_size": 4096,
    "spam_reject_score": 12.0,
    "auth_brute_force_threshold": 3
  }'
```

---

## Best Practices

1. **Secure Your Admin Token**: Use a strong, randomly generated token for production
2. **Use HTTPS**: Always use HTTPS in production to protect authentication tokens
3. **Minimal Changes**: Update only the specific settings you need to change
4. **Test in Dev**: Test configuration changes in development before applying to production
5. **Backup Before Changes**: Backup your SQLite database before making bulk configuration changes
6. **Monitor Logs**: Check server logs after configuration changes to verify proper application
7. **Template Strategy**: Set sensible defaults in the `_default` template for consistent security policies

---

## Integration with Mail Services

Security settings are automatically applied by SMTP and IMAP services when they read domain configuration from the database. Changes take effect immediately without requiring server restart (hot-reload).

When processing mail:
1. SMTP/IMAP services lookup the domain from the database
2. Security settings are loaded per-domain
3. Appropriate security checks are applied based on domain configuration
4. Rate limits, greylisting, and brute force protection track state per-domain

---

## SQLite Direct Access

You can also manage domain configuration directly via SQL:

```sql
-- Enable TOTP for a domain
UPDATE domains SET auth_totp_enforced = 1 WHERE name = 'example.com';

-- Adjust spam scores
UPDATE domains
SET spam_reject_score = 15.0, spam_quarantine_score = 7.0
WHERE name = 'example.com';

-- Disable greylisting
UPDATE domains SET greylist_enabled = 0 WHERE name = 'example.com';
```

However, using the API is recommended as it provides:
- Input validation
- Proper error handling
- Audit logging
- Consistent interface
