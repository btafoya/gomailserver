# ISSUE002: IMAP Authentication Configuration and Testing

**Status**: ✅ RESOLVED

## Problem

IMAP authentication was failing during testing, preventing successful login.

## Root Causes

### 1. Database Path Mismatch
The server configuration (`gomailserver.yaml`) pointed to `./data/mailserver.db`, but the test data (users and domains) was in `./mailserver.db`.

### 2. Brute Force Protection Blocking
Previous failed authentication attempts were recorded in the `failed_logins` table, causing brute force protection to block further authentication attempts from `127.0.0.1` (localhost).

## Solution

### 1. Corrected Database Path

Updated `gomailserver.yaml` to point to the correct database:

```yaml
database:
    path: ./mailserver.db  # Changed from ./data/mailserver.db
    wal_enabled: true
```

### 2. Cleared Brute Force Protection Records

Removed failed login records for localhost during testing:

```sql
DELETE FROM failed_logins WHERE ip_address = '127.0.0.1';
DELETE FROM ip_blacklist WHERE ip_address = '127.0.0.1';
```

## IMAP Authentication Features Verified

### AUTH Capability Advertisement ✅
- **Before STARTTLS**: `AUTH=PLAIN` advertised
- **After STARTTLS**: `AUTH=PLAIN` advertised
- **Compliance**: IMAP AUTHENTICATE extension working correctly

### Authentication Flow ✅
1. Client connects to port 2143
2. Server advertises CAPABILITY including AUTH=PLAIN
3. Client initiates STARTTLS
4. Server upgrades to TLS
5. Client re-checks CAPABILITY (still includes AUTH=PLAIN)
6. Client authenticates using LOGIN command
7. Server validates credentials via UserService
8. Server checks domain configuration
9. Server applies brute force protection
10. Server creates authenticated IMAP session

## Testing

### Verification Steps
1. ✅ IMAP server listens on port 2143 (STARTTLS)
2. ✅ IMAPS server listens on port 2993 (implicit TLS)
3. ✅ AUTH=PLAIN advertised before and after STARTTLS
4. ✅ Authentication succeeds with valid credentials
5. ✅ Session established successfully
6. ✅ Logout completes cleanly

### Test Results
```
✅ Connected successfully
✅ AUTH=PLAIN is advertised
✅ STARTTLS successful
✅ Login successful!
✅ Logout successful

Server log:
{"level":"info","msg":"IMAP authentication attempt","username":"test@localhost"}
{"level":"info","msg":"authentication successful","email":"test@localhost","user_id":4}
{"level":"info","msg":"IMAP authentication successful","username":"test@localhost","user_id":4}
{"level":"debug","msg":"IMAP session ended","username":"test@localhost","user_id":4}
```

## go-imap Architecture Notes

### Backend Interface
The IMAP backend implements `github.com/emersion/go-imap/backend.Backend` interface:

```go
type Backend interface {
    Login(connInfo *imap.ConnInfo, username, password string) (User, error)
}
```

### User Interface  
The authenticated user implements `backend.User` interface with methods:
- `Username() string`
- `ListMailboxes(subscribed bool) ([]Mailbox, error)`
- `GetMailbox(name string) (Mailbox, error)`
- `CreateMailbox(name string) error`
- `DeleteMailbox(name string) error`
- `RenameMailbox(oldName, newName string) error`
- `Logout() error`

### Security Features Integrated
- ✅ Brute force protection with configurable threshold and window
- ✅ Rate limiting per user/IP
- ✅ TOTP enforcement (when enabled)
- ✅ User status checking (active/disabled)
- ✅ Domain configuration loading
- ✅ Comprehensive logging

## Comparison with SMTP Authentication

Both SMTP and IMAP authentication share the same underlying `UserService.Authenticate()` method and security infrastructure:

| Feature | SMTP | IMAP |
|---------|------|------|
| **Authentication Method** | SASL PLAIN via AuthSession interface | LOGIN command via Backend.Login() |
| **User Service** | UserService.Authenticate() | UserService.Authenticate() |
| **Brute Force Protection** | ✅ Shared | ✅ Shared |
| **Rate Limiting** | ✅ Per method | ✅ Per user |
| **TOTP Support** | ✅ | ✅ |
| **Domain Config** | ✅ | ✅ |
| **TLS Support** | STARTTLS + implicit (SMTPS) | STARTTLS + implicit (IMAPS) |

## Files Modified

- `gomailserver.yaml` - Corrected database path
- Database: Cleared brute force protection records for testing

## References

- go-imap documentation: https://github.com/emersion/go-imap
- IMAP AUTHENTICATE extension: RFC 3501
- IMAP capabilities: RFC 3501 Section 6.1.1

## Impact

**Security**: No impact - existing security features preserved
**Compatibility**: Standard IMAP authentication working correctly
**Configuration**: Database path now consistent for all testing
**Breaking Changes**: None - configuration fix only
