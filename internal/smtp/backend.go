package smtp

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	repService "github.com/btafoya/gomailserver/internal/reputation/service"
	"github.com/btafoya/gomailserver/internal/security/antispam"
	"github.com/btafoya/gomailserver/internal/security/antivirus"
	"github.com/btafoya/gomailserver/internal/security/bruteforce"
	"github.com/btafoya/gomailserver/internal/security/dkim"
	"github.com/btafoya/gomailserver/internal/security/dmarc"
	"github.com/btafoya/gomailserver/internal/security/greylist"
	"github.com/btafoya/gomailserver/internal/security/ratelimit"
	"github.com/btafoya/gomailserver/internal/security/spf"
	mailService "github.com/btafoya/gomailserver/internal/service"
)

// Backend implements SMTP backend interface
type Backend struct {
	userService      mailService.UserServiceInterface
	messageService   mailService.MessageServiceInterface
	queueService     mailService.QueueServiceInterface
	domainRepo       repository.DomainRepository
	telemetryService *repService.TelemetryService
	logger           *zap.Logger

	// Security services
	dkimSigner    *dkim.Signer
	dkimVerifier  *dkim.Verifier
	spfValidator  *spf.Validator
	dmarcEnforcer *dmarc.Enforcer
	greylister    *greylist.Greylister
	rateLimiter   *ratelimit.Limiter
	bruteForce    *bruteforce.Protection
	clamav        *antivirus.ClamAV
	spamAssassin  *antispam.SpamAssassin
}

// NewBackend creates a new SMTP backend with all dependencies
func NewBackend(
	userService mailService.UserServiceInterface,
	messageService mailService.MessageServiceInterface,
	queueService mailService.QueueServiceInterface,
	domainRepo repository.DomainRepository,
	telemetryService *repService.TelemetryService,
	dkimSigner *dkim.Signer,
	dkimVerifier *dkim.Verifier,
	spfValidator *spf.Validator,
	dmarcEnforcer *dmarc.Enforcer,
	greylister *greylist.Greylister,
	rateLimiter *ratelimit.Limiter,
	bruteForce *bruteforce.Protection,
	clamav *antivirus.ClamAV,
	spamAssassin *antispam.SpamAssassin,
	logger *zap.Logger,
) *Backend {
	return &Backend{
		userService:      userService,
		messageService:   messageService,
		queueService:     queueService,
		domainRepo:       domainRepo,
		telemetryService: telemetryService,
		logger:           logger,
		dkimSigner:       dkimSigner,
		dkimVerifier:     dkimVerifier,
		spfValidator:     spfValidator,
		dmarcEnforcer:    dmarcEnforcer,
		greylister:       greylister,
		rateLimiter:      rateLimiter,
		bruteForce:       bruteForce,
		clamav:           clamav,
		spamAssassin:     spamAssassin,
	}
}

// NewSession creates a new SMTP session
func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		conn:           c,
		backend:        b,
		logger:         b.logger,
		remoteAddr:     c.Conn().RemoteAddr().String(),
		authenticated:  false,
	}, nil
}

// Session represents an SMTP session
type Session struct {
	conn           *smtp.Conn
	backend        *Backend
	logger         *zap.Logger
	remoteAddr     string
	authenticated  bool
	username       string
	from           string
	to             []string
}

// AuthPlain implements PLAIN authentication
func (s *Session) AuthPlain(username, password string) error {
	s.logger.Info("SMTP authentication attempt",
		zap.String("username", username),
		zap.String("remote_addr", s.remoteAddr),
		zap.String("method", "PLAIN"),
	)

	// Extract domain from username
	domain := extractDomain(username)
	if domain == "" {
		return &smtp.SMTPError{
			Code:         535,
			EnhancedCode: smtp.EnhancedCode{5, 7, 8},
			Message:      "Invalid username format",
		}
	}

	// Load domain configuration
	domainConfig, err := s.backend.domainRepo.GetByName(domain)
	if err != nil {
		s.logger.Error("failed to load domain config",
			zap.String("domain", domain),
			zap.Error(err),
		)
		// Continue with authentication even if domain config fails
		domainConfig = nil
	}

	// Extract IP address
	remoteIP := extractIP(s.remoteAddr)

	// Check brute force protection if enabled
	if domainConfig != nil && s.backend.bruteForce != nil && domainConfig.AuthBruteForceEnabled {
		blocked, err := s.backend.bruteForce.IsBlocked(remoteIP)
		if err != nil {
			s.logger.Error("brute force check failed", zap.Error(err))
		} else if blocked {
			s.logger.Warn("SMTP authentication blocked - brute force protection",
				zap.String("username", username),
				zap.String("remote_ip", remoteIP),
			)
			return &smtp.SMTPError{
				Code:         421,
				EnhancedCode: smtp.EnhancedCode{4, 7, 1},
				Message:      "Too many failed login attempts",
			}
		}
	}

	// Check auth rate limiting if enabled
	if domainConfig != nil && s.backend.rateLimiter != nil && domainConfig.RateLimitEnabled {
		allowed, err := s.backend.rateLimiter.CheckAuth(remoteIP)
		if err != nil {
			s.logger.Error("rate limit check failed", zap.Error(err))
		} else if !allowed {
			s.logger.Warn("SMTP authentication rate limited",
				zap.String("username", username),
				zap.String("remote_ip", remoteIP),
				zap.String("domain", domain),
			)
			return &smtp.SMTPError{
				Code:         421,
				EnhancedCode: smtp.EnhancedCode{4, 7, 1},
				Message:      "Rate limit exceeded",
			}
		}
	}

	user, err := s.backend.userService.Authenticate(username, password)
	if err != nil {
		// Record failed login attempt for brute force protection
		if domainConfig != nil && s.backend.bruteForce != nil && domainConfig.AuthBruteForceEnabled {
			if err := s.backend.bruteForce.RecordFailure(remoteIP, username); err != nil {
				s.logger.Error("failed to record login failure", zap.Error(err))
			}
		}

		s.logger.Warn("SMTP authentication failed",
			zap.String("username", username),
			zap.String("remote_addr", s.remoteAddr),
			zap.Error(err),
		)
		return &smtp.SMTPError{
			Code:         535,
			EnhancedCode: smtp.EnhancedCode{5, 7, 8},
			Message:      "Authentication failed",
		}
	}

	if user.Status != "active" {
		s.logger.Warn("SMTP authentication failed - user disabled",
			zap.String("username", username),
			zap.String("status", user.Status),
		)
		return &smtp.SMTPError{
			Code:         535,
			EnhancedCode: smtp.EnhancedCode{5, 7, 8},
			Message:      "Account disabled",
		}
	}

	// Record successful login for brute force protection
	if domainConfig != nil && s.backend.bruteForce != nil && domainConfig.AuthBruteForceEnabled {
		if err := s.backend.bruteForce.RecordSuccess(remoteIP, username); err != nil {
			s.logger.Error("failed to record successful login", zap.Error(err))
		}
	}

	s.authenticated = true
	s.username = username

	s.logger.Info("SMTP authentication successful",
		zap.String("username", username),
		zap.String("remote_addr", s.remoteAddr),
	)

	return nil
}

// AuthMechanisms returns the list of supported authentication mechanisms
// This method implements the AuthSession interface to enable AUTH advertisement
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth creates a SASL server for the specified mechanism
// This method implements the AuthSession interface to enable AUTH advertisement
func (s *Session) Auth(mech string) (sasl.Server, error) {
	if mech != sasl.Plain {
		return nil, &smtp.SMTPError{
			Code:         504,
			EnhancedCode: smtp.EnhancedCode{5, 7, 4},
			Message:      "Unsupported authentication mechanism",
		}
	}

	return sasl.NewPlainServer(func(identity, username, password string) error {
		// The identity parameter is typically empty for PLAIN auth
		// Use username if provided, otherwise fall back to identity
		authUser := username
		if authUser == "" {
			authUser = identity
		}

		return s.AuthPlain(authUser, password)
	}), nil
}

// Mail is called when the client sends MAIL FROM
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	// For relay (port 25), no authentication required
	// For submission (port 587, 465), authentication required
	if !s.authenticated && s.conn.Server().Addr != fmt.Sprintf(":%d", 25) {
		return &smtp.SMTPError{
			Code:         530,
			EnhancedCode: smtp.EnhancedCode{5, 7, 0},
			Message:      "Authentication required",
		}
	}

	// Extract domain for rate limiting
	domain := extractDomain(from)
	if domain == "" {
		return &smtp.SMTPError{
			Code:         501,
			EnhancedCode: smtp.EnhancedCode{5, 1, 7},
			Message:      "Invalid sender address",
		}
	}

	// Load domain configuration
	domainConfig, err := s.backend.domainRepo.GetByName(domain)
	if err != nil {
		s.logger.Warn("failed to load domain config for MAIL FROM",
			zap.String("domain", domain),
			zap.Error(err),
		)
		// Continue even if domain config fails
		domainConfig = nil
	}

	// Check SMTP rate limiting if enabled
	if domainConfig != nil && s.backend.rateLimiter != nil && domainConfig.RateLimitEnabled {
		remoteIP := extractIP(s.remoteAddr)

		// Check per-IP rate limit
		allowed, err := s.backend.rateLimiter.CheckIP(remoteIP)
		if err != nil {
			s.logger.Error("rate limit check failed", zap.Error(err))
		} else if !allowed {
			s.logger.Warn("SMTP rate limited by IP",
				zap.String("from", from),
				zap.String("remote_ip", remoteIP),
				zap.String("domain", domain),
			)
			return &smtp.SMTPError{
				Code:         421,
				EnhancedCode: smtp.EnhancedCode{4, 7, 1},
				Message:      "Rate limit exceeded",
			}
		}

		// Check per-user rate limit if authenticated
		if s.authenticated {
			allowed, err := s.backend.rateLimiter.CheckUser(s.username)
			if err != nil {
				s.logger.Error("user rate limit check failed", zap.Error(err))
			} else if !allowed {
				s.logger.Warn("SMTP rate limited by user",
					zap.String("from", from),
					zap.String("username", s.username),
					zap.String("domain", domain),
				)
				return &smtp.SMTPError{
					Code:         421,
					EnhancedCode: smtp.EnhancedCode{4, 7, 1},
					Message:      "Rate limit exceeded",
				}
			}
		}
	}

	s.from = from
	s.logger.Debug("MAIL FROM",
		zap.String("from", from),
		zap.String("remote_addr", s.remoteAddr),
		zap.String("username", s.username),
	)

	return nil
}

// Rcpt is called when the client sends RCPT TO
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// TODO: Validate recipient exists
	// TODO: Check quota
	// TODO: Check greylisting
	// TODO: Check rate limiting

	s.to = append(s.to, to)
	s.logger.Debug("RCPT TO",
		zap.String("to", to),
		zap.String("from", s.from),
		zap.String("remote_addr", s.remoteAddr),
	)

	return nil
}

// Data is called when the client sends DATA
func (s *Session) Data(r io.Reader) error {
	s.logger.Info("receiving message",
		zap.String("from", s.from),
		zap.Strings("to", s.to),
		zap.String("remote_addr", s.remoteAddr),
	)

	// Read entire message
	data, err := io.ReadAll(r)
	if err != nil {
		s.logger.Error("failed to read message data", zap.Error(err))
		return &smtp.SMTPError{
			Code:         451,
			EnhancedCode: smtp.EnhancedCode{4, 0, 0},
			Message:      "Failed to read message",
		}
	}

	// Determine if this is inbound relay or authenticated submission
	isInboundRelay := !s.authenticated

	// Get domain configuration for first recipient
	var domainConfig *domain.Domain
	if len(s.to) > 0 {
		recipientDomain := extractDomain(s.to[0])
		if recipientDomain != "" {
			domainConfig, err = s.backend.domainRepo.GetByName(recipientDomain)
			if err != nil {
				s.logger.Warn("failed to load domain config for DATA",
					zap.String("domain", recipientDomain),
					zap.Error(err),
				)
				// Continue even if domain config fails
				domainConfig = nil
			}
		}
	}

	remoteIP := extractIP(s.remoteAddr)

	// For inbound relay, apply security checks
	if isInboundRelay && domainConfig != nil {
		// 1. Greylisting
		if s.backend.greylister != nil && domainConfig.GreylistEnabled {
			result, err := s.backend.greylister.Check(remoteIP, s.from, s.to[0])
			if err != nil {
				s.logger.Error("greylist check failed", zap.Error(err))
			} else if result.Action == "defer" {
				s.logger.Info("message greylisted - temporary rejection",
					zap.String("from", s.from),
					zap.String("to", s.to[0]),
					zap.String("remote_ip", remoteIP),
					zap.Duration("wait_time", result.WaitTime),
				)
				return &smtp.SMTPError{
					Code:         451,
					EnhancedCode: smtp.EnhancedCode{4, 7, 1},
					Message:      "Greylisted - please try again later",
				}
			}
		}

		// 2. SPF Validation
		if s.backend.spfValidator != nil && domainConfig.SPFEnabled {
			senderDomain := extractDomain(s.from)
			ipAddr := net.ParseIP(remoteIP)
			if ipAddr != nil {
				spfResult, err := s.backend.spfValidator.Check(ipAddr, senderDomain, s.from)
				if err != nil {
					s.logger.Error("SPF validation failed", zap.Error(err))
				} else {
					s.logger.Info("SPF validation result",
						zap.String("result", string(spfResult)),
						zap.String("from", s.from),
						zap.String("remote_ip", remoteIP),
					)

					// Apply SPF policy
					if spfResult == "fail" && domainConfig.SPFFailAction == "reject" {
						return &smtp.SMTPError{
							Code:         550,
							EnhancedCode: smtp.EnhancedCode{5, 7, 1},
							Message:      "SPF validation failed",
						}
					}
				}
			}
		}

		// 3. DKIM Verification
		if s.backend.dkimVerifier != nil && domainConfig.DKIMVerifyEnabled {
			verifications, err := s.backend.dkimVerifier.Verify(data)
			if err != nil {
				s.logger.Warn("DKIM verification failed",
					zap.Error(err),
					zap.String("from", s.from),
				)
			} else {
				// Check if any signature is valid
				for _, v := range verifications {
					if v.Valid {
						s.logger.Info("DKIM signature verified",
							zap.String("domain", v.Domain),
							zap.String("selector", v.Selector),
							zap.String("from", s.from),
						)
						break
					}
				}
			}
		}

		// 4. DMARC Enforcement (requires SPF and DKIM results)
		// Note: This is simplified - full DMARC would require SPF and DKIM results
		// For now, we skip DMARC enforcement as it needs proper message parsing

		// 5. Virus Scanning (ClamAV)
		if s.backend.clamav != nil && domainConfig.ClamAVEnabled {
			scanResult, err := s.backend.clamav.Scan(data)
			if err != nil {
				s.logger.Error("virus scan failed", zap.Error(err))
				// Apply fail action
				if domainConfig.ClamAVFailAction == "reject" {
					return &smtp.SMTPError{
						Code:         451,
						EnhancedCode: smtp.EnhancedCode{4, 7, 0},
						Message:      "Unable to scan message",
					}
				}
			} else if !scanResult.Clean {
				s.logger.Warn("virus detected in message",
					zap.String("virus", scanResult.Virus),
					zap.String("from", s.from),
					zap.Strings("to", s.to),
				)
				// Apply virus action
				if domainConfig.ClamAVVirusAction == "reject" {
					return &smtp.SMTPError{
						Code:         550,
						EnhancedCode: smtp.EnhancedCode{5, 7, 1},
						Message:      fmt.Sprintf("Virus detected: %s", scanResult.Virus),
					}
				}
			}
		}

		// 6. Spam Filtering (SpamAssassin)
		if s.backend.spamAssassin != nil && domainConfig.SpamEnabled {
			spamResult, err := s.backend.spamAssassin.Check(data)
			if err != nil {
				s.logger.Error("spam check failed", zap.Error(err))
			} else {
				s.logger.Info("spam check result",
					zap.Float64("score", spamResult.Score),
					zap.Bool("is_spam", spamResult.IsSpam),
					zap.String("from", s.from),
				)

				// Apply spam policy
				if spamResult.Score >= domainConfig.SpamRejectScore {
					s.logger.Info("message rejected as spam",
						zap.Float64("score", spamResult.Score),
						zap.Float64("threshold", domainConfig.SpamRejectScore),
					)
					return &smtp.SMTPError{
						Code:         550,
						EnhancedCode: smtp.EnhancedCode{5, 7, 1},
						Message:      "Message rejected as spam",
					}
				}
			}
		}
	}

	// For outbound authenticated mail, apply DKIM signing
	if !isInboundRelay && s.backend.dkimSigner != nil && domainConfig != nil && domainConfig.DKIMSigningEnabled {
		senderDomain := extractDomain(s.from)
		signedData, err := s.backend.dkimSigner.Sign(senderDomain, data)
		if err != nil {
			s.logger.Error("DKIM signing failed",
				zap.Error(err),
				zap.String("from", s.from),
			)
			// Continue without signing on error
		} else {
			data = signedData
			s.logger.Debug("DKIM signature added",
				zap.String("from", s.from),
			)
		}
	}

	// Queue message for delivery
	messageID, err := s.backend.queueService.Enqueue(s.from, s.to, data)
	if err != nil {
		s.logger.Error("failed to queue message",
			zap.Error(err),
			zap.String("from", s.from),
			zap.Strings("to", s.to),
		)
		return &smtp.SMTPError{
			Code:         451,
			EnhancedCode: smtp.EnhancedCode{4, 3, 0},
			Message:      "Failed to queue message",
		}
	}

	s.logger.Info("message accepted",
		zap.String("message_id", messageID),
		zap.String("from", s.from),
		zap.Strings("to", s.to),
		zap.Int("size", len(data)),
	)

	// Record telemetry for outbound authenticated mail
	if !isInboundRelay && s.backend.telemetryService != nil {
		senderDomain := extractDomain(s.from)
		for _, recipient := range s.to {
			recipientDomain := extractDomain(recipient)
			if recipientDomain != "" {
				// Record as queued for delivery (will be updated when actually delivered)
				// Note: This is optimistic - actual delivery telemetry should be recorded
				// by the delivery worker when processing the queue
				ctx := context.Background()
				if err := s.backend.telemetryService.RecordDelivery(ctx, senderDomain, recipientDomain, remoteIP); err != nil {
					s.logger.Warn("failed to record telemetry",
						zap.Error(err),
						zap.String("domain", senderDomain),
					)
				}
			}
		}
	}

	return nil
}

// Reset is called when the client sends RSET
func (s *Session) Reset() {
	s.from = ""
	s.to = nil
}

// Logout is called when the session ends
func (s *Session) Logout() error {
	s.logger.Debug("SMTP session ended",
		zap.String("username", s.username),
		zap.String("remote_addr", s.remoteAddr),
	)
	return nil
}

// extractDomain extracts domain from email address
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// extractIP extracts IP address from remote address string
func extractIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}
