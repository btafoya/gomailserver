package smtp

import (
	"fmt"
	"io"
	"strings"

	"github.com/emersion/go-smtp"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/service"
)

// Backend implements SMTP backend interface
type Backend struct {
	userService    service.UserServiceInterface
	messageService service.MessageServiceInterface
	queueService   service.QueueServiceInterface
	logger         *zap.Logger
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

	user, err := s.backend.userService.Authenticate(username, password)
	if err != nil {
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

	s.authenticated = true
	s.username = username

	s.logger.Info("SMTP authentication successful",
		zap.String("username", username),
		zap.String("remote_addr", s.remoteAddr),
	)

	return nil
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
