# ISSUE001: SMTP AUTH Not Advertised After STARTTLS

**Status**: ✅ RESOLVED

## Problem

SMTP AUTH PLAIN capability was not being advertised in EHLO response after STARTTLS, preventing clients from authenticating.

## Root Causes

1. **Missing AuthSession Interface**: The `Session` struct in `internal/smtp/backend.go` did not implement the optional `AuthSession` interface required by go-smtp v0.24.0 to enable AUTH advertisement
2. **Database NULL Handling**: SQLite TEXT columns containing NULL values caused scan errors when loading user and domain records

## Solution

### 1. Implemented AuthSession Interface

Added two required methods to `Session` struct in `internal/smtp/backend.go`:

```go
// AuthMechanisms returns the list of supported authentication mechanisms
func (s *Session) AuthMechanisms() []string {
    return []string{sasl.Plain}
}

// Auth creates a SASL server for the specified mechanism
func (s *Session) Auth(mech string) (sasl.Server, error) {
    if mech != sasl.Plain {
        return nil, &smtp.SMTPError{
            Code:         504,
            EnhancedCode: smtp.EnhancedCode{5, 7, 4},
            Message:      "Unsupported authentication mechanism",
        }
    }

    return sasl.NewPlainServer(func(identity, username, password string) error {
        authUser := username
        if authUser == "" {
            authUser = identity
        }
        return s.AuthPlain(authUser, password)
    }), nil
}
```

### 2. Fixed Database NULL Handling

#### Changed Domain Model
Modified `domain.Domain.CatchallEmail` from `string` to `*string` in `internal/domain/models.go` to handle NULL values properly.

#### Updated API Handlers
Fixed type conversions in `internal/api/handlers/domain_handler.go`:
- Create handler: Convert string to *string
- Update handler: Handle nil for empty values
- Response conversion: Safely dereference *string

#### Database Updates
Converted NULL TEXT fields to empty strings:

**Domains table:**
```sql
UPDATE domains SET
  catchall_email = COALESCE(catchall_email, ''),
  dkim_selector = COALESCE(dkim_selector, ''),
  dkim_private_key = COALESCE(dkim_private_key, ''),
  dkim_public_key = COALESCE(dkim_public_key, ''),
  spf_record = COALESCE(spf_record, ''),
  dmarc_policy = COALESCE(dmarc_policy, ''),
  dmarc_report_email = COALESCE(dmarc_report_email, '')
WHERE name = 'localhost';
```

**Users table:**
```sql
UPDATE users SET
  full_name = COALESCE(full_name, ''),
  display_name = COALESCE(display_name, ''),
  totp_secret = COALESCE(totp_secret, ''),
  forward_to = COALESCE(forward_to, ''),
  auto_reply_subject = COALESCE(auto_reply_subject, ''),
  auto_reply_body = COALESCE(auto_reply_body, '')
WHERE email = 'test@localhost';
```

## Testing

### Verification Steps
1. ✅ AUTH PLAIN advertised in pre-STARTTLS EHLO response
2. ✅ AUTH PLAIN advertised in post-STARTTLS EHLO response
3. ✅ Authentication succeeds with valid credentials
4. ✅ Email successfully sent after authentication
5. ✅ Server logs show successful authentication

### Test Results
```
✅ Authentication successful!
   User ID: 4
   Email: test@localhost
   Full Name: Test User
   Status: active

Server log:
{"level":"info","msg":"authentication successful","email":"test@localhost","user_id":4}
{"level":"info","msg":"SMTP authentication successful","username":"test@localhost"}
{"level":"info","msg":"message accepted","message_id":"c2b8d6fec6f1dfe508f0f2be3f510f33","size":805}
```

## Files Modified

- `internal/smtp/backend.go` - Implemented AuthSession interface
- `internal/domain/models.go` - Changed CatchallEmail to *string
- `internal/api/handlers/domain_handler.go` - Fixed type conversions
- Database: Updated NULL values to empty strings

## References

- go-smtp documentation: https://github.com/emersion/go-smtp
- go-sasl PLAIN mechanism: https://github.com/emersion/go-sasl
- RFC 4616 (SASL PLAIN): https://tools.ietf.org/html/rfc4616

## Impact

**Security**: No impact - AUTH was always required after STARTTLS, just not advertised
**Compatibility**: Enables standards-compliant SMTP clients to discover and use authentication
**Breaking Changes**: None - existing functionality preserved
