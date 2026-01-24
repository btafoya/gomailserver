# API Examples

## Setup Wizard

### Check Setup Status
```bash
curl http://localhost:8980/api/v1/setup/status
```

### Create First Admin User
```bash
curl -X POST http://localhost:8980/api/v1/setup/admin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "securepassword123",
    "name": "Administrator"
  }'
```

### Complete Setup
```bash
curl -X POST http://localhost:8980/api/v1/setup/complete \
  -H "Content-Type: application/json"
```

## Domain Management

### Create Domain with DKIM
```bash
curl -X POST http://localhost:8980/api/v1/domains \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.com",
    "dkim_enabled": true,
    "spf_enabled": true,
    "dmarc_enabled": true
  }'
```

### Update Domain Security Settings
```bash
curl -X PUT http://localhost:8980/api/v1/domains/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "dkim_enabled": true,
    "spf_enabled": true,
    "dmarc_enabled": true,
    "rate_limit_smtp_per_domain": "{\"count\":1000,\"window_minutes\":60}"
  }'
```

## User Management

### Create User with Quota
```bash
curl -X POST http://localhost:8980/api/v1/users \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "userpassword123",
    "name": "John Doe",
    "quota": 1073741824
  }'
```

### Update User Password
```bash
curl -X POST http://localhost:8980/api/v1/users/1/password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "new_password": "newsecurepassword123"
  }'
```

## Webmail Operations

### Send HTML Email
```bash
curl -X POST http://localhost:8980/api/v1/webmail/messages \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "recipient@example.com",
    "subject": "Welcome!",
    "body_html": "<h1>Welcome</h1><p>Thank you for joining us.</p>",
    "body_text": "Welcome\n\nThank you for joining us."
  }'
```

### Create and Send Draft
```bash
# Save draft first
curl -X POST http://localhost:8980/api/v1/webmail/drafts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["recipient@example.com"],
    "subject": "Draft Email",
    "body_text": "This is a draft message"
  }'

# Send draft (would modify draft then send)
```

### Search Messages
```bash
curl "http://localhost:8980/api/v1/webmail/search?q=important&limit=20" \
  -H "Authorization: Bearer <token>"
```

## Queue Management

### View Queue Status
```bash
curl http://localhost:8980/api/v1/queue \
  -H "Authorization: Bearer <token>"
```

### Retry Failed Message
```bash
curl -X POST http://localhost:8980/api/v1/queue/123/retry \
  -H "Authorization: Bearer <token>"
```

## Reputation Management

### Get Domain Reputation Score
```bash
curl http://localhost:8980/api/v1/reputation/domains/example.com/score \
  -H "Authorization: Bearer <token>"
```

### View DMARC Reports
```bash
curl http://localhost:8980/api/v1/reputation/dmarc-reports \
  -H "Authorization: Bearer <token>"
```

### Get External Metrics (Gmail/Microsoft)
```bash
curl http://localhost:8980/api/v1/reputation/external-metrics \
  -H "Authorization: Bearer <token>"
```

## Settings Management

### Update TLS Settings
```bash
curl -X PUT http://localhost:8980/api/v1/settings/tls \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "acme_enabled": true,
    "acme_email": "admin@example.com",
    "acme_provider": "cloudflare",
    "acme_api_token": "your-cloudflare-token"
  }'
```

## PostmarkApp Compatibility

### Send Email (PostmarkApp compatible)
```bash
curl -X POST http://localhost:8980/email \
  -H "X-Postmark-Server-Token: your-api-token" \
  -H "Content-Type: application/json" \
  -d '{
    "From": "sender@example.com",
    "To": "recipient@example.com",
    "Subject": "Hello from PostmarkApp API",
    "TextBody": "This works just like PostmarkApp!"
  }'
```

### Batch Send
```bash
curl -X POST http://localhost:8980/email/batch \
  -H "X-Postmark-Server-Token: your-api-token" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "From": "sender@example.com",
      "To": "user1@example.com",
      "Subject": "Batch Email 1",
      "TextBody": "First message"
    },
    {
      "From": "sender@example.com",
      "To": "user2@example.com",
      "Subject": "Batch Email 2",
      "TextBody": "Second message"
    }
  ]'
```

## Webhook Management

### Create Delivery Webhook
```bash
curl -X POST http://localhost:8980/api/v1/webhooks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Email Delivery Notifications",
    "url": "https://yourapp.com/webhooks/email",
    "secret": "webhook-secret-123",
    "event_types": [
      "email.delivered",
      "email.bounced",
      "email.spam_complaint"
    ]
  }'
```

### Test Webhook
```bash
curl -X POST http://localhost:8980/api/v1/webhooks/1/test \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "email.delivered",
    "test_data": {
      "message_id": "test-123",
      "recipient": "test@example.com"
    }
  }'
```

## Monitoring and Logs

### Get Server Logs
```bash
curl http://localhost:8980/api/v1/logs?level=error&limit=100 \
  -H "Authorization: Bearer <token>"
```

### Get System Statistics
```bash
curl http://localhost:8980/api/v1/stats \
  -H "Authorization: Bearer <token>"
```

## Error Examples

### Invalid Authentication
```bash
curl http://localhost:8980/api/v1/users
# Response: 401 Unauthorized
# {"error": "Missing authorization header"}
```

### Rate Limited
```bash
curl http://localhost:8980/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test","password":"test"}'
# Response: 429 Too Many Requests
# {"error": "Rate limit exceeded"}
```

### Validation Error
```bash
curl -X POST http://localhost:8980/api/v1/users \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"invalid-email","password":"short"}'
# Response: 400 Bad Request
# {"error": "Email format invalid and password too short"}
```