package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/security/dmarc"
	"github.com/btafoya/gomailserver/internal/security/spf"
	tlsManager "github.com/btafoya/gomailserver/internal/tls"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

// AuditorService performs deliverability readiness audits for domains
type AuditorService struct {
	spfResolver   *spf.Resolver
	dmarcResolver *dmarc.Resolver
	dnsClient     *dns.Client
	tlsManager    *tlsManager.Manager
	logger        *zap.Logger
	nameserver    string
}

// NewAuditorService creates a new auditor service
func NewAuditorService(tlsMgr *tlsManager.Manager, logger *zap.Logger) *AuditorService {
	return &AuditorService{
		spfResolver:   spf.NewResolver(),
		dmarcResolver: dmarc.NewResolver(),
		dnsClient:     new(dns.Client),
		tlsManager:    tlsMgr,
		logger:        logger,
		nameserver:    "8.8.8.8:53", // Google DNS
	}
}

// AuditDomain performs a comprehensive deliverability audit for a domain
func (s *AuditorService) AuditDomain(ctx context.Context, domainName string, sendingIP net.IP) (*domain.AuditResult, error) {
	s.logger.Info("Starting deliverability audit", zap.String("domain", domainName))

	result := &domain.AuditResult{
		Domain:    domainName,
		Timestamp: time.Now().Unix(),
		Issues:    make([]string, 0),
	}

	// Run all checks in parallel for efficiency
	spfChan := make(chan domain.CheckStatus, 1)
	dkimChan := make(chan domain.CheckStatus, 1)
	dmarcChan := make(chan domain.CheckStatus, 1)
	rdnsChan := make(chan domain.CheckStatus, 1)
	fcrDNSChan := make(chan domain.CheckStatus, 1)
	tlsChan := make(chan domain.CheckStatus, 1)
	mtastsChan := make(chan domain.CheckStatus, 1)
	postmasterChan := make(chan bool, 1)
	abuseChan := make(chan bool, 1)

	// Launch all checks concurrently
	go func() { spfChan <- s.checkSPF(ctx, domainName) }()
	go func() { dkimChan <- s.checkDKIM(ctx, domainName, "default") }()
	go func() { dmarcChan <- s.checkDMARC(ctx, domainName) }()
	go func() { rdnsChan <- s.checkRDNS(ctx, sendingIP) }()
	go func() { fcrDNSChan <- s.checkFCrDNS(ctx, domainName, sendingIP) }()
	go func() { tlsChan <- s.checkTLS(ctx, domainName) }()
	go func() { mtastsChan <- s.checkMTASTS(ctx, domainName) }()
	go func() { postmasterChan <- s.checkOperationalMailbox(ctx, domainName, "postmaster") }()
	go func() { abuseChan <- s.checkOperationalMailbox(ctx, domainName, "abuse") }()

	// Collect results
	result.SPFStatus = <-spfChan
	result.DKIMStatus = <-dkimChan
	result.DMARCStatus = <-dmarcChan
	result.RDNSStatus = <-rdnsChan
	result.FCrDNSStatus = <-fcrDNSChan
	result.TLSStatus = <-tlsChan
	result.MTASTSStatus = <-mtastsChan
	result.PostmasterOK = <-postmasterChan
	result.AbuseOK = <-abuseChan

	// Calculate overall score and collect issues
	result.OverallScore = s.calculateOverallScore(result)
	result.Issues = s.collectIssues(result)

	s.logger.Info("Audit completed",
		zap.String("domain", domainName),
		zap.Int("overall_score", result.OverallScore),
		zap.Int("issues_count", len(result.Issues)),
	)

	return result, nil
}

// checkSPF validates SPF record presence and syntax
func (s *AuditorService) checkSPF(ctx context.Context, domainName string) domain.CheckStatus {
	record, err := s.spfResolver.LookupSPF(domainName)
	if err != nil {
		if err == spf.ErrNoSPFRecord {
			return domain.CheckStatus{
				Passed:  false,
				Message: "No SPF record found",
				Details: map[string]interface{}{"error": err.Error()},
			}
		}
		return domain.CheckStatus{
			Passed:  false,
			Message: fmt.Sprintf("SPF lookup failed: %v", err),
			Details: map[string]interface{}{"error": err.Error()},
		}
	}

	// Validate SPF record syntax
	if !strings.HasPrefix(record, "v=spf1") {
		return domain.CheckStatus{
			Passed:  false,
			Message: "Invalid SPF record format",
			Details: map[string]interface{}{"record": record},
		}
	}

	return domain.CheckStatus{
		Passed:  true,
		Message: "SPF record valid",
		Details: map[string]interface{}{"record": record},
	}
}

// checkDKIM validates DKIM selector DNS records
func (s *AuditorService) checkDKIM(ctx context.Context, domainName, selector string) domain.CheckStatus {
	// Query for DKIM public key at selector._domainkey.domain
	dkimDomain := fmt.Sprintf("%s._domainkey.%s", selector, domainName)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(dkimDomain), dns.TypeTXT)

	resp, _, err := s.dnsClient.Exchange(m, s.nameserver)
	if err != nil {
		return domain.CheckStatus{
			Passed:  false,
			Message: fmt.Sprintf("DKIM DNS lookup failed: %v", err),
			Details: map[string]interface{}{"error": err.Error(), "selector": selector},
		}
	}

	// Look for DKIM record
	for _, ans := range resp.Answer {
		if txt, ok := ans.(*dns.TXT); ok {
			record := strings.Join(txt.Txt, "")
			if strings.Contains(record, "v=DKIM1") {
				return domain.CheckStatus{
					Passed:  true,
					Message: "DKIM record found",
					Details: map[string]interface{}{"record": record, "selector": selector},
				}
			}
		}
	}

	return domain.CheckStatus{
		Passed:  false,
		Message: "No DKIM record found",
		Details: map[string]interface{}{"selector": selector, "query": dkimDomain},
	}
}

// checkDMARC validates DMARC policy presence
func (s *AuditorService) checkDMARC(ctx context.Context, domainName string) domain.CheckStatus {
	policy, err := s.dmarcResolver.LookupDMARC(domainName)
	if err != nil {
		if err == dmarc.ErrNoDMARCRecord {
			return domain.CheckStatus{
				Passed:  false,
				Message: "No DMARC record found",
				Details: map[string]interface{}{"error": err.Error()},
			}
		}
		return domain.CheckStatus{
			Passed:  false,
			Message: fmt.Sprintf("DMARC lookup failed: %v", err),
			Details: map[string]interface{}{"error": err.Error()},
		}
	}

	// Check for strict policy recommendation
	strictPolicy := policy.Policy == "quarantine" || policy.Policy == "reject"

	return domain.CheckStatus{
		Passed: true,
		Message: fmt.Sprintf("DMARC policy: %s (strict: %v)", policy.Policy, strictPolicy),
		Details: map[string]interface{}{
			"policy":    policy.Policy,
			"strict":    strictPolicy,
			"aggregate": policy.ReportAggregate,
			"forensic":  policy.ReportForensic,
		},
	}
}

// checkRDNS validates reverse DNS (PTR) records for the sending IP
func (s *AuditorService) checkRDNS(ctx context.Context, ip net.IP) domain.CheckStatus {
	if ip == nil {
		return domain.CheckStatus{
			Passed:  false,
			Message: "No IP address provided for rDNS check",
			Details: map[string]interface{}{},
		}
	}

	names, err := s.spfResolver.LookupPTR(ip)
	if err != nil {
		return domain.CheckStatus{
			Passed:  false,
			Message: fmt.Sprintf("No PTR record found: %v", err),
			Details: map[string]interface{}{"ip": ip.String(), "error": err.Error()},
		}
	}

	return domain.CheckStatus{
		Passed:  true,
		Message: fmt.Sprintf("PTR record found: %s", names[0]),
		Details: map[string]interface{}{"ip": ip.String(), "ptr_records": names},
	}
}

// checkFCrDNS validates forward-confirmed reverse DNS
func (s *AuditorService) checkFCrDNS(ctx context.Context, domainName string, ip net.IP) domain.CheckStatus {
	if ip == nil {
		return domain.CheckStatus{
			Passed:  false,
			Message: "No IP address provided for FCrDNS check",
			Details: map[string]interface{}{},
		}
	}

	// Step 1: Get PTR records for IP
	ptrNames, err := s.spfResolver.LookupPTR(ip)
	if err != nil {
		return domain.CheckStatus{
			Passed:  false,
			Message: "FCrDNS failed: no PTR record",
			Details: map[string]interface{}{"ip": ip.String(), "error": err.Error()},
		}
	}

	// Step 2: For each PTR record, look up A/AAAA records
	for _, ptrName := range ptrNames {
		// Try A records
		ips, err := s.spfResolver.LookupA(ptrName)
		if err == nil {
			// Check if original IP is in the A records
			for _, foundIP := range ips {
				if foundIP.Equal(ip) {
					return domain.CheckStatus{
						Passed:  true,
						Message: "FCrDNS validated successfully",
						Details: map[string]interface{}{
							"ip":         ip.String(),
							"ptr_record": ptrName,
							"matched":    true,
						},
					}
				}
			}
		}
	}

	return domain.CheckStatus{
		Passed:  false,
		Message: "FCrDNS failed: forward lookup does not match",
		Details: map[string]interface{}{
			"ip":          ip.String(),
			"ptr_records": ptrNames,
			"matched":     false,
		},
	}
}

// checkTLS validates TLS certificate validity and expiry
func (s *AuditorService) checkTLS(ctx context.Context, domainName string) domain.CheckStatus {
	// Get MX records first
	mxRecords, err := s.spfResolver.LookupMX(domainName)
	if err != nil {
		return domain.CheckStatus{
			Passed:  false,
			Message: fmt.Sprintf("No MX records found: %v", err),
			Details: map[string]interface{}{"error": err.Error()},
		}
	}

	// Test TLS connection to first MX
	mxHost := strings.TrimSuffix(mxRecords[0], ".")
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:465", mxHost), &tls.Config{
		ServerName: mxHost,
	})
	if err != nil {
		return domain.CheckStatus{
			Passed:  false,
			Message: fmt.Sprintf("TLS connection failed: %v", err),
			Details: map[string]interface{}{"mx": mxHost, "error": err.Error()},
		}
	}
	defer conn.Close()

	// Validate certificate
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return domain.CheckStatus{
			Passed:  false,
			Message: "No TLS certificate presented",
			Details: map[string]interface{}{"mx": mxHost},
		}
	}

	cert := state.PeerCertificates[0]
	now := time.Now()

	// Check expiry
	if now.After(cert.NotAfter) {
		return domain.CheckStatus{
			Passed:  false,
			Message: "TLS certificate expired",
			Details: map[string]interface{}{
				"mx":         mxHost,
				"expired_at": cert.NotAfter,
			},
		}
	}

	daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)
	if daysUntilExpiry <= 30 {
		return domain.CheckStatus{
			Passed:  true,
			Message: fmt.Sprintf("TLS certificate valid but expiring soon (%d days)", daysUntilExpiry),
			Details: map[string]interface{}{
				"mx":               mxHost,
				"days_until_expiry": daysUntilExpiry,
				"expires_at":       cert.NotAfter,
			},
		}
	}

	return domain.CheckStatus{
		Passed:  true,
		Message: fmt.Sprintf("TLS certificate valid (%d days until expiry)", daysUntilExpiry),
		Details: map[string]interface{}{
			"mx":                mxHost,
			"days_until_expiry": daysUntilExpiry,
			"expires_at":        cert.NotAfter,
		},
	}
}

// checkMTASTS validates MTA-STS policy presence
func (s *AuditorService) checkMTASTS(ctx context.Context, domainName string) domain.CheckStatus {
	// Query for MTA-STS DNS record: _mta-sts.domain
	mtastsDomain := fmt.Sprintf("_mta-sts.%s", domainName)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(mtastsDomain), dns.TypeTXT)

	resp, _, err := s.dnsClient.Exchange(m, s.nameserver)
	if err != nil {
		return domain.CheckStatus{
			Passed:  false,
			Message: fmt.Sprintf("MTA-STS DNS lookup failed: %v", err),
			Details: map[string]interface{}{"error": err.Error()},
		}
	}

	// Look for MTA-STS record
	for _, ans := range resp.Answer {
		if txt, ok := ans.(*dns.TXT); ok {
			record := strings.Join(txt.Txt, "")
			if strings.Contains(record, "v=STSv1") {
				return domain.CheckStatus{
					Passed:  true,
					Message: "MTA-STS record found",
					Details: map[string]interface{}{"record": record},
				}
			}
		}
	}

	return domain.CheckStatus{
		Passed:  false,
		Message: "No MTA-STS record found",
		Details: map[string]interface{}{"query": mtastsDomain},
	}
}

// checkOperationalMailbox validates operational mailbox deliverability
func (s *AuditorService) checkOperationalMailbox(ctx context.Context, domainName, mailbox string) bool {
	email := fmt.Sprintf("%s@%s", mailbox, domainName)

	// Get MX records
	mxRecords, err := s.spfResolver.LookupMX(domainName)
	if err != nil {
		s.logger.Debug("MX lookup failed for operational mailbox check",
			zap.String("mailbox", email),
			zap.Error(err),
		)
		return false
	}

	// Try to connect to first MX and verify address
	mxHost := strings.TrimSuffix(mxRecords[0], ".")

	// Connect to SMTP server
	client, err := smtp.Dial(fmt.Sprintf("%s:25", mxHost))
	if err != nil {
		s.logger.Debug("SMTP connection failed for operational mailbox check",
			zap.String("mailbox", email),
			zap.String("mx", mxHost),
			zap.Error(err),
		)
		return false
	}
	defer client.Close()

	// Try MAIL FROM (use noreply@domain)
	if err := client.Mail(fmt.Sprintf("noreply@%s", domainName)); err != nil {
		s.logger.Debug("MAIL FROM failed for operational mailbox check",
			zap.String("mailbox", email),
			zap.Error(err),
		)
		return false
	}

	// Try RCPT TO to verify mailbox exists
	if err := client.Rcpt(email); err != nil {
		s.logger.Debug("RCPT TO failed for operational mailbox check",
			zap.String("mailbox", email),
			zap.Error(err),
		)
		return false
	}

	// If we get here, the mailbox appears to exist
	s.logger.Debug("Operational mailbox validated",
		zap.String("mailbox", email),
	)
	return true
}

// calculateOverallScore calculates the overall audit score (0-100)
func (s *AuditorService) calculateOverallScore(result *domain.AuditResult) int {
	score := 100

	// Critical checks (major deductions)
	if !result.SPFStatus.Passed {
		score -= 20 // SPF is critical
	}
	if !result.DKIMStatus.Passed {
		score -= 20 // DKIM is critical
	}
	if !result.DMARCStatus.Passed {
		score -= 15 // DMARC is very important
	}

	// Important checks (moderate deductions)
	if !result.RDNSStatus.Passed {
		score -= 10 // rDNS is important
	}
	if !result.FCrDNSStatus.Passed {
		score -= 10 // FCrDNS is important
	}
	if !result.TLSStatus.Passed {
		score -= 10 // TLS is important
	}

	// Nice-to-have checks (small deductions)
	if !result.MTASTSStatus.Passed {
		score -= 5 // MTA-STS is nice but not critical
	}
	if !result.PostmasterOK {
		score -= 5 // Operational mailboxes are recommended
	}
	if !result.AbuseOK {
		score -= 5 // Operational mailboxes are recommended
	}

	// Ensure score is in valid range
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// collectIssues collects all issues found during the audit
func (s *AuditorService) collectIssues(result *domain.AuditResult) []string {
	issues := make([]string, 0)

	if !result.SPFStatus.Passed {
		issues = append(issues, fmt.Sprintf("SPF: %s", result.SPFStatus.Message))
	}
	if !result.DKIMStatus.Passed {
		issues = append(issues, fmt.Sprintf("DKIM: %s", result.DKIMStatus.Message))
	}
	if !result.DMARCStatus.Passed {
		issues = append(issues, fmt.Sprintf("DMARC: %s", result.DMARCStatus.Message))
	}
	if !result.RDNSStatus.Passed {
		issues = append(issues, fmt.Sprintf("rDNS: %s", result.RDNSStatus.Message))
	}
	if !result.FCrDNSStatus.Passed {
		issues = append(issues, fmt.Sprintf("FCrDNS: %s", result.FCrDNSStatus.Message))
	}
	if !result.TLSStatus.Passed {
		issues = append(issues, fmt.Sprintf("TLS: %s", result.TLSStatus.Message))
	}
	if !result.MTASTSStatus.Passed {
		issues = append(issues, fmt.Sprintf("MTA-STS: %s", result.MTASTSStatus.Message))
	}
	if !result.PostmasterOK {
		issues = append(issues, "postmaster@ mailbox not deliverable")
	}
	if !result.AbuseOK {
		issues = append(issues, "abuse@ mailbox not deliverable")
	}

	return issues
}
