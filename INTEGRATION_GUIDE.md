# Security Integration Guide

This document outlines how to integrate Phase 2 security features (DKIM, SPF, DMARC, ClamAV, SpamAssassin, Greylisting, Rate Limiting, TOTP, Brute Force Protection) into the SMTP and IMAP services to read per-domain settings from SQLite.

## Current State

### Completed (Phase 2)
- ✅ All security modules implemented (`internal/security/`)
- ✅ SQLite schema with per-domain security configuration
- ✅ DomainService with default template system
- ✅ Admin API for managing domain security settings
- ✅ Security configuration stored in `domains` table

### Pending (Phase 3 Integration)
- ⏳ SMTP service integration with domain security settings
- ⏳ IMAP service integration with domain security settings
- ⏳ Security service initialization (ClamAV, SpamAssassin connections)
- ⏳ Message processing pipeline with security checks

---

## Architecture Overview

```
┌─────────────┐
│  SMTP/IMAP  │
│   Services  │
└──────┬──────┘
       │
       ├──> Extract domain from email address
       │
       ├──> Load domain security config from SQLite
       │    (via DomainRepository)
       │
       ├──> Apply security checks based on domain config:
       │    │
       │    ├──> Rate Limiting (SMTP connections, auth attempts)
       │    ├──> Brute Force Protection (login attempts)
       │    ├──> Greylisting (inbound SMTP)
       │    ├──> SPF Validation (inbound SMTP)
       │    ├──> DKIM Verification (inbound)
       │    ├──> DKIM Signing (outbound)
       │    ├──> DMARC Enforcement (inbound)
       │    ├──> ClamAV Scanning (all messages)
       │    └──> SpamAssassin Scoring (inbound messages)
       │
       └──> Deliver/Queue/Reject based on security results
```

---

## Integration Steps

### Step 1: SMTP Backend Enhancement

**File**: `internal/smtp/backend.go`

#### Add Domain Repository and Security Services

```go
type Backend struct {
	userService     service.UserServiceInterface
	messageService  service.MessageServiceInterface
	queueService    service.QueueServiceInterface
	domainRepo      repository.DomainRepository

	// Security services
	dkimSigner      *dkim.Signer
	dkimVerifier    *dkim.Verifier
	spfValidator    *spf.Validator
	dmarcEnforcer   *dmarc.Enforcer
	clamAVScanner   *antivirus.Scanner
	spamAssassin    *antispam.SpamAssassin
	greylist        *greylist.Greylist
	rateLimiter     *ratelimit.Limiter
	bruteForce      *bruteforce.Protection

	logger          *zap.Logger
}
```

#### Session-Level Domain Loading

```go
type Session struct {
	// ... existing fields ...
	domainConfig *domain.Domain  // Loaded per-domain security config
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	// Extract domain from sender address
	parts := strings.Split(from, "@")
	if len(parts) != 2 {
		return &smtp.SMTPError{Code: 553, Message: "Invalid sender address"}
	}
	senderDomain := parts[1]

	// Load domain configuration
	domainConfig, err := s.backend.domainRepo.GetByName(senderDomain)
	if err != nil {
		s.logger.Warn("domain not found", zap.String("domain", senderDomain))
		// Use default fallback or reject
		return &smtp.SMTPError{Code: 550, Message: "Domain not configured"}
	}

	s.domainConfig = domainConfig

	// Apply rate limiting based on domain config
	if domainConfig.RateLimitEnabled {
		allowed, err := s.backend.rateLimiter.CheckSMTPPerIP(
			s.remoteAddr,
			domainConfig.RateLimitSMTPPerIP,
		)
		if err != nil || !allowed {
			return &smtp.SMTPError{Code: 450, Message: "Rate limit exceeded"}
		}
	}

	// ... rest of Mail() logic
}
```

#### Authentication with Brute Force Protection

```go
func (s *Session) AuthPlain(username, password string) error {
	// Extract domain from username
	parts := strings.Split(username, "@")
	if len(parts) != 2 {
		return &smtp.SMTPError{Code: 535, Message: "Authentication failed"}
	}
	userDomain := parts[1]

	// Load domain config
	domainConfig, err := s.backend.domainRepo.GetByName(userDomain)
	if err != nil {
		return &smtp.SMTPError{Code: 535, Message: "Authentication failed"}
	}

	// Check brute force protection
	if domainConfig.AuthBruteForceEnabled {
		blocked, err := s.backend.bruteForce.IsBlocked(
			s.remoteAddr,
			domainConfig.AuthBruteForceThreshold,
			domainConfig.AuthBruteForceWindowMinutes,
		)
		if err == nil && blocked {
			return &smtp.SMTPError{Code: 535, Message: "Too many failed attempts"}
		}
	}

	// Authenticate
	user, err := s.backend.userService.Authenticate(username, password)
	if err != nil {
		// Record failed attempt
		if domainConfig.AuthBruteForceEnabled {
			s.backend.bruteForce.RecordFailure(s.remoteAddr, username)
		}
		return &smtp.SMTPError{Code: 535, Message: "Authentication failed"}
	}

	// Check TOTP if enforced
	if domainConfig.AuthTOTPEnforced && user.TOTPSecret != "" {
		// TOTP validation logic here
		// (would require TOTP token in authentication flow)
	}

	// Success
	s.authenticated = true
	s.username = username
	s.domainConfig = domainConfig
	return nil
}
```

#### Message Data Processing with Security Checks

```go
func (s *Session) Data(r io.Reader) error {
	// Read message
	messageBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Apply greylisting for inbound relay
	if !s.authenticated && s.domainConfig.GreylistEnabled {
		shouldDelay, err := s.backend.greylist.Check(
			s.remoteAddr,
			s.from,
			s.to[0],
			s.domainConfig.GreylistDelayMinutes,
		)
		if err == nil && shouldDelay {
			return &smtp.SMTPError{
				Code: 451,
				Message: "Please try again later (greylisted)",
			}
		}
	}

	// SPF validation for inbound
	if !s.authenticated && s.domainConfig.SPFEnabled {
		result, err := s.backend.spfValidator.Validate(
			s.remoteAddr,
			s.from,
			s.domainConfig.SPFDNSServer,
			s.domainConfig.SPFDNSTimeout,
		)
		if err == nil {
			switch result {
			case "fail":
				if s.domainConfig.SPFFailAction == "reject" {
					return &smtp.SMTPError{Code: 550, Message: "SPF validation failed"}
				}
			case "softfail":
				if s.domainConfig.SPFSoftFailAction == "reject" {
					return &smtp.SMTPError{Code: 550, Message: "SPF soft fail"}
				}
			}
		}
	}

	// DKIM verification for inbound
	if !s.authenticated && s.domainConfig.DKIMVerifyEnabled {
		valid, err := s.backend.dkimVerifier.Verify(messageBytes)
		// Log result but don't reject solely on DKIM failure
		s.logger.Info("DKIM verification", zap.Bool("valid", valid), zap.Error(err))
	}

	// DMARC enforcement for inbound
	if !s.authenticated && s.domainConfig.DMARCEnabled {
		action, err := s.backend.dmarcEnforcer.Enforce(
			s.from,
			messageBytes,
			s.domainConfig.DMARCDNSServer,
		)
		if err == nil && action == "reject" {
			return &smtp.SMTPError{Code: 550, Message: "DMARC policy rejection"}
		}
	}

	// ClamAV scanning
	if s.domainConfig.ClamAVEnabled {
		infected, err := s.backend.clamAVScanner.Scan(messageBytes)
		if err == nil && infected {
			switch s.domainConfig.ClamAVVirusAction {
			case "reject":
				return &smtp.SMTPError{Code: 550, Message: "Virus detected"}
			case "quarantine":
				// Queue to quarantine instead of normal delivery
			case "tag":
				// Add X-Virus-Scanned header
			}
		}
	}

	// SpamAssassin scoring for inbound
	if !s.authenticated && s.domainConfig.SpamEnabled {
		score, err := s.backend.spamAssassin.Score(messageBytes)
		if err == nil {
			if score >= s.domainConfig.SpamRejectScore {
				return &smtp.SMTPError{Code: 550, Message: "Message rejected as spam"}
			} else if score >= s.domainConfig.SpamQuarantineScore {
				// Queue to quarantine
			}
		}
	}

	// DKIM signing for outbound
	if s.authenticated && s.domainConfig.DKIMSigningEnabled {
		signedMessage, err := s.backend.dkimSigner.Sign(
			messageBytes,
			s.domainConfig.DKIMSelector,
			s.domainConfig.DKIMPrivateKey,
			s.domainConfig.DKIMHeadersToSign,
		)
		if err == nil {
			messageBytes = signedMessage
		}
	}

	// Store/queue message
	// ... existing message storage logic ...

	return nil
}
```

---

### Step 2: IMAP Backend Enhancement

**File**: `internal/imap/backend.go`

#### Add Domain Repository

```go
type Backend struct {
	userService    service.UserServiceInterface
	mailboxService service.MailboxServiceInterface
	messageService service.MessageServiceInterface
	domainRepo     repository.DomainRepository
	bruteForce     *bruteforce.Protection
	rateLimiter    *ratelimit.Limiter
	logger         *zap.Logger
}
```

#### Login with Security Checks

```go
func (b *Backend) Login(connInfo *imap.ConnInfo, username, password string) (imap.User, error) {
	// Extract domain
	parts := strings.Split(username, "@")
	if len(parts) != 2 {
		return nil, errors.New("invalid username")
	}
	userDomain := parts[1]

	// Load domain config
	domainConfig, err := b.domainRepo.GetByName(userDomain)
	if err != nil {
		return nil, errors.New("authentication failed")
	}

	// Check brute force protection
	if domainConfig.AuthBruteForceEnabled {
		remoteAddr := connInfo.RemoteAddr.String()
		blocked, _ := b.bruteForce.IsBlocked(
			remoteAddr,
			domainConfig.AuthBruteForceThreshold,
			domainConfig.AuthBruteForceWindowMinutes,
		)
		if blocked {
			return nil, errors.New("too many failed attempts")
		}
	}

	// Check IMAP rate limit
	if domainConfig.RateLimitEnabled {
		allowed, _ := b.rateLimiter.CheckIMAPPerUser(
			username,
			domainConfig.RateLimitIMAPPerUser,
		)
		if !allowed {
			return nil, errors.New("rate limit exceeded")
		}
	}

	// Authenticate
	user, err := b.userService.Authenticate(username, password)
	if err != nil {
		if domainConfig.AuthBruteForceEnabled {
			b.bruteForce.RecordFailure(connInfo.RemoteAddr.String(), username)
		}
		return nil, err
	}

	// Check TOTP if enforced
	if domainConfig.AuthTOTPEnforced && user.TOTPSecret != "" {
		// TOTP validation
		// (requires TOTP support in IMAP authentication flow)
	}

	return &User{
		user:           user,
		domainConfig:   domainConfig,
		mailboxService: b.mailboxService,
		messageService: b.messageService,
		logger:         b.logger,
	}, nil
}
```

---

### Step 3: Server Initialization with Security Services

**File**: `internal/commands/run.go`

#### Initialize Security Services

```go
func run(cmd *cobra.Command, args []string) error {
	// ... existing initialization ...

	// Initialize security repositories
	greylistRepo := sqlite.NewGreylistRepository(db)
	rateLimitRepo := sqlite.NewRateLimitRepository(db)
	loginAttemptRepo := sqlite.NewLoginAttemptRepository(db)
	ipBlacklistRepo := sqlite.NewIPBlacklistRepository(db)
	quarantineRepo := sqlite.NewQuarantineRepository(db)

	// Initialize security services
	greylistSvc := greylist.NewGreylist(greylistRepo)
	rateLimiter := ratelimit.NewLimiter(rateLimitRepo)
	bruteForce := bruteforce.NewProtection(loginAttemptRepo, ipBlacklistRepo)

	// Initialize DKIM services
	dkimSigner := dkim.NewSigner()
	dkimVerifier := dkim.NewVerifier()

	// Initialize SPF/DMARC
	spfValidator := spf.NewValidator()
	dmarcEnforcer := dmarc.NewEnforcer(spfValidator, dkimVerifier)

	// Initialize ClamAV connection
	clamAVScanner, err := antivirus.NewScanner(
		cfg.Security.ClamAV.SocketPath,
		cfg.Security.ClamAV.Timeout,
		logger,
	)
	if err != nil {
		logger.Warn("ClamAV connection failed - virus scanning disabled", zap.Error(err))
	}

	// Initialize SpamAssassin connection
	spamAssassin, err := antispam.NewSpamAssassin(
		cfg.Security.SpamAssassin.Host,
		cfg.Security.SpamAssassin.Port,
		cfg.Security.SpamAssassin.Timeout,
		logger,
	)
	if err != nil {
		logger.Warn("SpamAssassin connection failed - spam filtering disabled", zap.Error(err))
	}

	// Create SMTP server with security services
	smtpServer := smtp.NewServerWithSecurity(
		&cfg.SMTP,
		tlsCfg,
		userSvc,
		messageSvc,
		queueSvc,
		domainRepo,
		dkimSigner,
		dkimVerifier,
		spfValidator,
		dmarcEnforcer,
		clamAVScanner,
		spamAssassin,
		greylistSvc,
		rateLimiter,
		bruteForce,
		logger,
	)

	// Create IMAP server with security services
	imapServer := imap.NewServerWithSecurity(
		&cfg.IMAP,
		tlsCfg,
		userSvc,
		mailboxSvc,
		messageSvc,
		domainRepo,
		rateLimiter,
		bruteForce,
		logger,
	)

	// ... rest of startup logic ...
}
```

---

### Step 4: Background Cleanup Tasks

Add cleanup goroutines for expired data:

```go
// Start cleanup tasks
go func() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Cleanup expired greylist entries
			greylistSvc.Cleanup()

			// Cleanup expired rate limit entries
			rateLimiter.Cleanup()

			// Cleanup old login attempts
			bruteForce.Cleanup()

			// Remove expired IP blacklist entries
			ipBlacklistRepo.RemoveExpired()
		}
	}
}()
```

---

## Testing Strategy

### Unit Tests

For each integration point:

1. **Rate Limiting Tests**
   ```go
   func TestSMTPRateLimiting(t *testing.T) {
       // Test per-IP rate limiting
       // Test per-user rate limiting
       // Test per-domain rate limiting
   }
   ```

2. **Authentication Security Tests**
   ```go
   func TestBruteForceProtection(t *testing.T) {
       // Test failed login tracking
       // Test automatic IP blocking
       // Test block expiration
   }
   ```

3. **Mail Security Tests**
   ```go
   func TestMailSecurityPipeline(t *testing.T) {
       // Test SPF validation
       // Test DKIM signing/verification
       // Test DMARC enforcement
       // Test greylisting
       // Test virus scanning
       // Test spam filtering
   }
   ```

### Integration Tests

Test full mail flow with security checks:

```bash
# Test inbound mail with all security checks
./tests/integration/test_inbound_security.sh

# Test outbound mail with DKIM signing
./tests/integration/test_outbound_dkim.sh

# Test rate limiting enforcement
./tests/integration/test_rate_limits.sh

# Test brute force protection
./tests/integration/test_brute_force.sh
```

---

## Migration Path

### For Existing Installations

1. **Backup**: Backup SQLite database before migration
2. **Migrate Schema**: Run migration v2 (automatic on server start)
3. **Configure Domains**: Use Admin API to configure domain security settings
4. **External Services**: Ensure ClamAV and SpamAssassin are running
5. **Test**: Verify security checks are working with test messages
6. **Monitor**: Check logs for security events

### Configuration Migration

If migrating from YAML config to per-domain SQLite:

```bash
# Use Admin API to replicate YAML settings per-domain
curl -X PUT http://localhost:8080/api/domains/example.com/security \
  -H "Authorization: Bearer token" \
  -H "Content-Type: application/json" \
  -d '{
    "dkim": { "signing_enabled": true },
    "spf": { "enabled": true },
    ...
  }'
```

---

## Performance Considerations

1. **Domain Config Caching**: Cache loaded domain configurations for a short period (e.g., 60 seconds) to reduce database queries
2. **Security Service Pooling**: Use connection pooling for ClamAV and SpamAssassin
3. **Async Processing**: Move heavy security checks (virus scanning, spam filtering) to async queue processing where possible
4. **Database Indexing**: Ensure indexes exist on frequently queried columns (domain name, rate limit keys, greylist triplets)
5. **Cleanup Scheduling**: Run cleanup tasks during off-peak hours

---

## Monitoring and Observability

Add metrics for:
- Security check durations
- Rate limit hits
- Greylisting deferrals
- SPF/DKIM/DMARC failures
- Virus detections
- Spam scores distribution
- Authentication failures and blocks

Structured logging examples:

```go
logger.Info("security_check",
	zap.String("type", "spf"),
	zap.String("domain", domain),
	zap.String("result", "pass"),
	zap.Duration("duration", duration),
)

logger.Warn("rate_limit_exceeded",
	zap.String("type", "smtp_per_ip"),
	zap.String("ip", remoteAddr),
	zap.String("domain", domain),
)

logger.Error("security_service_failure",
	zap.String("service", "clamav"),
	zap.Error(err),
)
```

---

## Next Steps

1. Implement SMTP backend integration (Step 1)
2. Implement IMAP backend integration (Step 2)
3. Update server initialization (Step 3)
4. Add background cleanup tasks (Step 4)
5. Write comprehensive tests
6. Update deployment documentation
7. Create admin dashboard for security monitoring

---

## See Also

- [API.md](API.md) - Admin API documentation
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment guide with SQLite-first configuration
- [ISSUE003.md](ISSUE003.md) - Phase 2 security implementation details
