package spf

import (
	"errors"
	"net"
	"strings"

	"github.com/miekg/dns"
)

var ErrNoSPFRecord = errors.New("no SPF record found")

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

func (r *Resolver) LookupSPF(domain string) (string, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)

	resp, _, err := r.client.Exchange(m, r.nameserver)
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

func (r *Resolver) LookupA(domain string) ([]net.IP, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	resp, _, err := r.client.Exchange(m, r.nameserver)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, ans := range resp.Answer {
		if a, ok := ans.(*dns.A); ok {
			ips = append(ips, a.A)
		}
	}

	if len(ips) == 0 {
		return nil, errors.New("no A records found")
	}

	return ips, nil
}

func (r *Resolver) LookupMX(domain string) ([]string, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)

	resp, _, err := r.client.Exchange(m, r.nameserver)
	if err != nil {
		return nil, err
	}

	var mxRecords []string
	for _, ans := range resp.Answer {
		if mx, ok := ans.(*dns.MX); ok {
			mxRecords = append(mxRecords, mx.Mx)
		}
	}

	if len(mxRecords) == 0 {
		return nil, errors.New("no MX records found")
	}

	return mxRecords, nil
}

func (r *Resolver) LookupPTR(ip net.IP) ([]string, error) {
	addr, err := dns.ReverseAddr(ip.String())
	if err != nil {
		return nil, err
	}

	m := new(dns.Msg)
	m.SetQuestion(addr, dns.TypePTR)

	resp, _, err := r.client.Exchange(m, r.nameserver)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, ans := range resp.Answer {
		if ptr, ok := ans.(*dns.PTR); ok {
			names = append(names, ptr.Ptr)
		}
	}

	if len(names) == 0 {
		return nil, errors.New("no PTR records found")
	}

	return names, nil
}
