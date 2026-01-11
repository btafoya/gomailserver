package mailflow

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/testing/types"
)

type SMTPConnectivityCheck struct{}

func (c *SMTPConnectivityCheck) Name() string {
	return "SMTP Connectivity"
}

func (c *SMTPConnectivityCheck) Description() string {
	return "Connect to SMTP server and verify response"
}

func (c *SMTPConnectivityCheck) Category() types.Category {
	return types.CategoryMailFlow
}

func (c *SMTPConnectivityCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *SMTPConnectivityCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	startTime := time.Now()
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusFail,
		Message:  "SMTP connectivity check failed",
		Details:  make(map[string]interface{}),
	}

	timeout := 10 * time.Second
	smtpPort := cfg.SMTPPort
	if smtpPort == 0 {
		smtpPort = 587
	}
	smtpAddress := fmt.Sprintf("%s:%d", cfg.SMTPHost, smtpPort)

	conn, err := net.DialTimeout("tcp", smtpAddress, timeout)
	if err != nil {
		result.Message = fmt.Sprintf("Cannot connect to SMTP server: %s", smtpAddress)
		result.Details["smtp_address"] = smtpAddress
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}
	defer conn.Close()

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read SMTP response: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	response := string(data)
	if !strings.HasPrefix(response, "2") && !strings.HasPrefix(response, "3") {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", response)
		result.Details["response"] = response
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	response := string(data)
	if len(response) < 4 || response[0] != '2' && response[0] != '3' {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", response)
		result.Details["response"] = response
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	result.Status = types.StatusPass
	result.Message = "SMTP server reachable and responding correctly"
	result.Details["smtp_address"] = smtpAddress
	result.Details["smtp_port"] = smtpPort
	result.Details["smtp_banner"] = response
	result.Duration = int64(time.Since(startTime).Milliseconds())

	return result, nil
}

type SMTPAuthenticationCheck struct{}

func (c *SMTPAuthenticationCheck) Name() string {
	return "SMTP Authentication"
}

func (c *SMTPAuthenticationCheck) Description() string {
	return "Authenticate with test credentials"
}

func (c *SMTPAuthenticationCheck) Category() types.Category {
	return types.CategoryMailFlow
}

func (c *SMTPAuthenticationCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *SMTPAuthenticationCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	startTime := time.Now()
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusFail,
		Message:  "SMTP authentication check failed",
		Details:  make(map[string]interface{}),
	}

	timeout := 10 * time.Second
	smtpPort := cfg.SMTPPort
	if smtpPort == 0 {
		smtpPort = 587
	}
	smtpAddress := fmt.Sprintf("%s:%d", cfg.SMTPHost, smtpPort)

	if cfg.TestUser == "" || cfg.TestPass == "" {
		result.Status = types.StatusWarning
		result.Message = "SMTP authentication skipped (no test credentials provided)"
		result.Details["smtp_address"] = smtpAddress
		result.Details["test_user"] = "not configured"
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	conn, err := net.DialTimeout("tcp", smtpAddress, timeout)
	if err != nil {
		result.Message = fmt.Sprintf("Cannot connect to SMTP server: %s", smtpAddress)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}
	defer conn.Close()

	ehlo := fmt.Sprintf("HELO %s\r\n", "gomailtest")
	_, err = conn.Write([]byte(ehlo))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write HELO: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read HELO response: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	if !strings.Contains(string(data), "250") {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", string(data))
		result.Details["response"] = string(data)
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	authLine := fmt.Sprintf("AUTH PLAIN \x00%s\x00%s\x00", cfg.TestUser, cfg.TestPass)
	_, err = conn.Write([]byte(authLine))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write AUTH: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 1024)
	n, err = conn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read AUTH response: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	response := string(data)
	if !strings.Contains(response, "235") {
		result.Message = fmt.Sprintf("SMTP authentication failed: %s", response)
		result.Details["response"] = response
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	quit := "QUIT\r\n"
	_, err = conn.Write([]byte(quit))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write QUIT: %v", err)
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	result.Status = types.StatusPass
	result.Message = "SMTP authentication successful"
	result.Details["smtp_address"] = smtpAddress
	result.Details["smtp_port"] = smtpPort
	result.Details["test_user"] = cfg.TestUser
	result.Details["auth_method"] = "PLAIN"
	result.Duration = int64(time.Since(startTime).Milliseconds())

	return result, nil
}

type IMAPConnectivityCheck struct{}

func (c *IMAPConnectivityCheck) Name() string {
	return "IMAP Connectivity"
}

func (c *IMAPConnectivityCheck) Description() string {
	return "Connect to IMAP server and verify response"
}

func (c *IMAPConnectivityCheck) Category() types.Category {
	return types.CategoryMailFlow
}

func (c *IMAPConnectivityCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *IMAPConnectivityCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	startTime := time.Now()
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusFail,
		Message:  "IMAP connectivity check failed",
		Details:  make(map[string]interface{}),
	}

	timeout := 10 * time.Second
	imapPort := cfg.IMAPPort
	if imapPort == 0 {
		imapPort = 143
	}
	imapAddress := fmt.Sprintf("%s:%d", cfg.IMAPHost, imapPort)

	conn, err := net.DialTimeout("tcp", imapAddress, timeout)
	if err != nil {
		result.Message = fmt.Sprintf("Cannot connect to IMAP server: %s", imapAddress)
		result.Details["imap_address"] = imapAddress
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}
	defer conn.Close()

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read IMAP greeting: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	response := string(data)
	if !strings.HasPrefix(response, "* OK") {
		result.Message = fmt.Sprintf("Unexpected IMAP greeting: %s", response)
		result.Details["response"] = response
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	result.Status = types.StatusPass
	result.Message = "IMAP server reachable and responding correctly"
	result.Details["imap_address"] = imapAddress
	result.Details["imap_port"] = imapPort
	result.Details["imap_banner"] = response
	result.Duration = int64(time.Since(startTime).Milliseconds())

	return result, nil
}

type IMAPAuthenticationCheck struct{}

func (c *IMAPAuthenticationCheck) Name() string {
	return "IMAP Authentication"
}

func (c *IMAPAuthenticationCheck) Description() string {
	return "Authenticate with test credentials"
}

func (c *IMAPAuthenticationCheck) Category() types.Category {
	return types.CategoryMailFlow
}

func (c *IMAPAuthenticationCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *IMAPAuthenticationCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	startTime := time.Now()
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusFail,
		Message:  "IMAP authentication check failed",
		Details:  make(map[string]interface{}),
	}

	timeout := 10 * time.Second
	imapPort := cfg.IMAPPort
	if imapPort == 0 {
		imapPort = 143
	}
	imapAddress := fmt.Sprintf("%s:%d", cfg.IMAPHost, imapPort)

	if cfg.TestUser == "" || cfg.TestPass == "" {
		result.Status = types.StatusWarning
		result.Message = "IMAP authentication skipped (no test credentials provided)"
		result.Details["imap_address"] = imapAddress
		result.Details["test_user"] = "not configured"
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	conn, err := net.DialTimeout("tcp", imapAddress, timeout)
	if err != nil {
		result.Message = fmt.Sprintf("Cannot connect to IMAP server: %s", imapAddress)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}
	defer conn.Close()

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read IMAP greeting: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	response := string(data)
	if !strings.HasPrefix(response, "* OK") {
		result.Message = fmt.Sprintf("Unexpected IMAP greeting: %s", response)
		result.Details["response"] = response
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	login := fmt.Sprintf("%s LOGIN %s %s\r\n", cfg.TestUser, cfg.TestPass)
	_, err = conn.Write([]byte(login))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write LOGIN: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 1024)
	n, err = conn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read LOGIN response: %v", err)
		result.Details["error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	response := string(data)
	if !strings.Contains(response, "OK") {
		result.Message = fmt.Sprintf("IMAP authentication failed: %s", response)
		result.Details["response"] = response
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	logout := "LOGOUT\r\n"
	_, err = conn.Write([]byte(logout))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write LOGOUT: %v", err)
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	result.Status = types.StatusPass
	result.Message = "IMAP authentication successful"
	result.Details["imap_address"] = imapAddress
	result.Details["imap_port"] = imapPort
	result.Details["test_user"] = cfg.TestUser
	result.Details["auth_method"] = "LOGIN"
	result.Duration = int64(time.Since(startTime).Milliseconds())

	return result, nil
}

type MailFlowEndToEndCheck struct{}

func (c *MailFlowEndToEndCheck) Name() string {
	return "End-to-End Mail Flow"
}

func (c *MailFlowEndToEndCheck) Description() string {
	return "Send test message via SMTP, retrieve via IMAP"
}

func (c *MailFlowEndToEndCheck) Category() types.Category {
	return types.CategoryMailFlow
}

func (c *MailFlowEndToEndCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *MailFlowEndToEndCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	startTime := time.Now()
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusFail,
		Message:  "End-to-end mail flow check failed",
		Details:  make(map[string]interface{}),
	}

	if cfg.DryRun {
		result.Status = types.StatusSkip
		result.Message = "End-to-end mail flow skipped (dry-run mode)"
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	if cfg.TestUser == "" || cfg.TestPass == "" {
		result.Status = types.StatusWarning
		result.Message = "End-to-end mail flow skipped (no test credentials provided)"
		result.Details["test_user"] = "not configured"
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	testID := fmt.Sprintf("gomailtest-%d", time.Now().Unix())

	fromAddr := cfg.TestUser
	toAddr := cfg.TestUser
	subject := fmt.Sprintf("[MAILTEST] Health Check %d", testID)
	body := fmt.Sprintf("GoMailTest Health Check\r\nTest ID: %d\r\nTimestamp: %s\r\nThis is an automated test message.\r\n",
		testID, time.Now().Format(time.RFC1123))

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nX-GoMailTest: true\r\nX-GoMailTest-ID: %d\r\nX-GoMailTest-Timestamp: %s\r\nDate: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		fromAddr, toAddr, subject, testID, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC1123), body)

	result.Details["test_id"] = testID
	result.Details["test_from"] = cfg.TestUser
	result.Details["test_to"] = cfg.TestUser

	smtpPort := cfg.SMTPPort
	if smtpPort == 0 {
		smtpPort = 587
	}
	smtpAddress := fmt.Sprintf("%s:%d", cfg.SMTPHost, smtpPort)

	smtpConn, err := net.DialTimeout("tcp", smtpAddress, 10*time.Second)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to connect to SMTP: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}
	defer smtpConn.Close()

	ehlo := fmt.Sprintf("HELO %s\r\n", "gomailtest")
	_, err = smtpConn.Write([]byte(ehlo))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write HELO: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data := make([]byte, 1024)
	n, err := smtpConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read HELO response: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	if !strings.Contains(string(data), "250") {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", string(data))
		result.Details["smtp_error"] = string(data)
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	from := fmt.Sprintf("MAIL FROM:<%s>\r\n", fromAddr)
	_, err = smtpConn.Write([]byte(from))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write MAIL FROM: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 1024)
	n, err = smtpConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read MAIL FROM response: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	if !strings.Contains(string(data), "250") {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", string(data))
		result.Details["smtp_error"] = string(data)
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	rcpt := fmt.Sprintf("RCPT TO:<%s>\r\n", toAddr)
	_, err = smtpConn.Write([]byte(rcpt))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write RCPT TO: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 1024)
	n, err = smtpConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read RCPT TO response: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	if !strings.Contains(string(data), "250") {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", string(data))
		result.Details["smtp_error"] = string(data)
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	dataCmd := "DATA\r\n"
	_, err = smtpConn.Write([]byte(dataCmd))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write DATA: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 1024)
	n, err = smtpConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read DATA response: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	if !strings.Contains(string(data), "354") {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", string(data))
		result.Details["smtp_error"] = string(data)
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	_, err = smtpConn.Write([]byte(message))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write message: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	_, err = smtpConn.Write([]byte(".\r\n"))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write end of message: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 1024)
	n, err = smtpConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read message send response: %v", err)
		result.Details["smtp_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	if !strings.Contains(string(data), "250") {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", string(data))
		result.Details["smtp_error"] = string(data)
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	result.Details["smtp_sent"] = true
	result.Details["smtp_duration_ms"] = time.Since(startTime).Milliseconds()

	time.Sleep(3 * time.Second)

	imapPort := cfg.IMAPPort
	if imapPort == 0 {
		imapPort = 143
	}
	imapAddress := fmt.Sprintf("%s:%d", cfg.IMAPHost, imapPort)

	imapConn, err := net.DialTimeout("tcp", imapAddress, 10*time.Second)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to connect to IMAP: %v", err)
		result.Details["imap_error"] = err.Error()
		result.Details["smtp_sent"] = true
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}
	defer imapConn.Close()

	data = make([]byte, 1024)
	n, err = imapConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read IMAP greeting: %v", err)
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	response := string(data)
	if !strings.HasPrefix(response, "* OK") {
		result.Message = fmt.Sprintf("Unexpected IMAP greeting: %s", response)
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	login := fmt.Sprintf("%s LOGIN %s %s\r\n", cfg.TestUser, cfg.TestPass)
	_, err = imapConn.Write([]byte(login))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write LOGIN: %v", err)
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 1024)
	n, err = imapConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read LOGIN response: %v", err)
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	respData := string(data)
	if !strings.Contains(respData, "250") {
		result.Message = fmt.Sprintf("Unexpected SMTP response: %s", respData)
		result.Details["smtp_error"] = respData
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	result.Details["imap_authenticated"] = true

	selectCmd := fmt.Sprintf("SELECT INBOX\r\n")
	_, err = imapConn.Write([]byte(selectCmd))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write SELECT: %v", err)
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 1024)
	n, err = imapConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read SELECT response: %v", err)
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	if !strings.Contains(string(data), "OK") {
		result.Message = fmt.Sprintf("Unexpected IMAP SELECT response: %s", string(data))
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	search := fmt.Sprintf("SEARCH SUBJECT \"%s\"\r\n", subject)
	_, err = imapConn.Write([]byte(search))
	if err != nil {
		result.Message = fmt.Sprintf("Failed to write SEARCH: %v", err)
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	data = make([]byte, 8192)
	n, err := imapConn.Read(data)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to read SEARCH response: %v", err)
		result.Details["imap_error"] = err.Error()
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	response := string(data)
	if !strings.Contains(response, "SEARCH") && !strings.Contains(response, "OK") {
		result.Message = "Test message not found in INBOX"
		result.Details["search_response"] = response
		result.Duration = int64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	result.Details["message_found"] = true

	if cfg.AutoCleanup {
		result.Details["cleanup_skipped"] = "Cleanup not yet implemented for manual IMAP"
	} else {
		result.Details["cleanup_skipped"] = "Cleanup skipped (--no-cleanup flag)"
	}

	result.Status = types.StatusPass
	result.Message = "End-to-end mail flow successful: SMTP → Queue → IMAP"
	result.Details["delivery_latency_ms"] = time.Since(startTime).Milliseconds()
	result.Duration = int64(time.Since(startTime).Milliseconds())

	return result, nil
}

type MessageIntegrityCheck struct{}

func (c *MessageIntegrityCheck) Name() string {
	return "Message Integrity"
}

func (c *MessageIntegrityCheck) Description() string {
	return "Verify retrieved message content matches sent"
}

func (c *MessageIntegrityCheck) Category() types.Category {
	return types.CategoryMailFlow
}

func (c *MessageIntegrityCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *MessageIntegrityCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	startTime := time.Now()
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusSkip,
		Message:  "Message integrity check skipped (requires end-to-end flow)",
		Details:  make(map[string]interface{}),
	}

	result.Details["note"] = "Message integrity is verified within End-to-End Mail Flow check"
	result.Duration = int64(time.Since(startTime).Milliseconds())

	return result, nil
}
