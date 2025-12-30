package spf

import (
	"errors"
	"strings"

	"github.com/miekg/dns"
)

var ErrNoSPFRecord = errors.New("no SPF record found")

type Resolver struct {
	client *dns.Client
}

func NewResolver() *Resolver {
	return &Resolver{client: new(dns.Client)}
}

func (r *Resolver) LookupSPF(domain string) (string, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)

	resp, _, err := r.client.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return "", err
	}

	for _, ans := range resp.Answer {
		if txt, ok := ans.(*dns.TXT); ok {
			record := strings.Join(txt.Txt, "")
			if strings.HasPrefix(record, "v=spf1") {
				return record, nil
			}
		}
	}

	return "", ErrNoSPFRecord
}
