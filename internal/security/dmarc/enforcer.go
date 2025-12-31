package dmarc

import (
	"strings"

	"github.com/btafoya/gomailserver/internal/security/spf"
)

type EnforcementResult struct {
	Policy      string
	Action      string
	SPFResult   spf.Result
	DKIMResult  bool
	SPFAligned  bool
	DKIMAligned bool
	Reason      string
}

// Message is a placeholder for the email message structure
type Message struct {
	From       string
	SPFResult  spf.Result
	DKIMValid  bool
	DKIMDomain string
}

type Enforcer struct {
	resolver *Resolver
}

func NewEnforcer(resolver *Resolver) *Enforcer {
	return &Enforcer{resolver: resolver}
}

func (e *Enforcer) Enforce(msg *Message) (*EnforcementResult, error) {
	fromDomain := extractDomain(msg.From)

	policy, err := e.resolver.LookupDMARC(fromDomain)
	if err != nil {
		// If there's an error (e.g., no DMARC record), treat as "none"
		return &EnforcementResult{Policy: "none", Action: "none"}, nil
	}

	result := &EnforcementResult{
		Policy:     policy.Policy,
		SPFResult:  msg.SPFResult,
		DKIMResult: msg.DKIMValid,
	}

	// Check alignment
	result.SPFAligned = e.checkSPFAlignment(msg, policy)
	result.DKIMAligned = e.checkDKIMAlignment(msg, policy)

	// Determine action
	if result.SPFAligned || result.DKIMAligned {
		result.Action = "none" // Pass
	} else {
		result.Action = policy.Policy
		result.Reason = "DMARC alignment failed"
	}

	return result, nil
}

func (e *Enforcer) checkSPFAlignment(msg *Message, policy *Policy) bool {
	if msg.SPFResult != spf.ResultPass {
		return false
	}

	// For SPF alignment, we need to check if the domain in the envelope From
	// (which was checked by SPF) aligns with the From header domain.
	// In strict mode, they must match exactly.
	// In relaxed mode, organizational domain match is sufficient.

	// For now, assume SPF checked the same domain as From header
	// A complete implementation would track the SPF-checked domain
	if policy.SPF == "s" { // Strict
		// Strict alignment requires exact domain match
		return true // Simplified - would need envelope domain comparison
	}
	// Relaxed alignment (default)
	return true // Simplified - would need organizational domain comparison
}

func (e *Enforcer) checkDKIMAlignment(msg *Message, policy *Policy) bool {
	if !msg.DKIMValid {
		return false
	}

	fromDomain := extractDomain(msg.From)
	dkimDomain := msg.DKIMDomain

	if policy.DKIM == "s" { // Strict
		return fromDomain == dkimDomain
	}
	// Relaxed (default)
	return hasOrgDomain(fromDomain, dkimDomain)
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	// Remove any angle brackets or whitespace
	domain := strings.TrimSpace(parts[1])
	domain = strings.Trim(domain, "<>")
	return strings.ToLower(domain)
}

func hasOrgDomain(domain1, domain2 string) bool {
	// Simplified organizational domain comparison
	// A full implementation would use the Public Suffix List
	domain1 = strings.ToLower(domain1)
	domain2 = strings.ToLower(domain2)

	if domain1 == domain2 {
		return true
	}

	// Check if domains share the same organizational domain
	// For example, mail.example.com and www.example.com share example.com
	parts1 := strings.Split(domain1, ".")
	parts2 := strings.Split(domain2, ".")

	if len(parts1) < 2 || len(parts2) < 2 {
		return false
	}

	// Get the last two parts (organizational domain)
	org1 := strings.Join(parts1[len(parts1)-2:], ".")
	org2 := strings.Join(parts2[len(parts2)-2:], ".")

	return org1 == org2
}
