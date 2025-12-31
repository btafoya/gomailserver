package dmarc

import (
	"errors"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var ErrNoDMARCRecord = errors.New("no DMARC record found")

type Policy struct {
	Version          string
	Policy           string // none, quarantine, reject
	SubdomainPolicy  string
	Percentage       int
	DKIM             string // r (relaxed) or s (strict)
	SPF              string // r (relaxed) or s (strict)
	ReportAggregate  string
	ReportForensic   string
	ReportInterval   int
	FailureReporting string
}

type Resolver struct {
	client     *dns.Client
	nameserver string
}

func NewResolver() *Resolver {
	return &Resolver{
		client:     new(dns.Client),
		nameserver: "8.8.8.8:53", // Google DNS
	}
}

func (r *Resolver) LookupDMARC(domain string) (*Policy, error) {
	dmarcDomain := "_dmarc." + domain
	record, err := r.lookupTXT(dmarcDomain)
	if err != nil {
		return nil, err
	}

	return parseDMARCRecord(record)
}

func (r *Resolver) lookupTXT(domain string) (string, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)

	resp, _, err := r.client.Exchange(m, r.nameserver)
	if err != nil {
		return "", err
	}

	for _, ans := range resp.Answer {
		if txt, ok := ans.(*dns.TXT); ok {
			record := strings.Join(txt.Txt, "")
			if strings.HasPrefix(record, "v=DMARC1") {
				return record, nil
			}
		}
	}

	return "", ErrNoDMARCRecord
}

func parseDMARCRecord(record string) (*Policy, error) {
	policy := &Policy{
		Version:    "DMARC1",
		Policy:     "none",
		DKIM:       "r", // Default relaxed
		SPF:        "r", // Default relaxed
		Percentage: 100,
	}

	parts := strings.Split(record, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "v":
			policy.Version = value
		case "p":
			policy.Policy = value
		case "sp":
			policy.SubdomainPolicy = value
		case "pct":
			if pct, err := strconv.Atoi(value); err == nil {
				policy.Percentage = pct
			}
		case "adkim":
			policy.DKIM = value
		case "aspf":
			policy.SPF = value
		case "rua":
			policy.ReportAggregate = value
		case "ruf":
			policy.ReportForensic = value
		case "ri":
			if interval, err := strconv.Atoi(value); err == nil {
				policy.ReportInterval = interval
			}
		case "fo":
			policy.FailureReporting = value
		}
	}

	if policy.Version != "DMARC1" {
		return nil, errors.New("invalid DMARC version")
	}

	return policy, nil
}
