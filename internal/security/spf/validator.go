package spf

import (
	"net"
	"strings"
)

type Result string

const (
	ResultNone      Result = "none"
	ResultNeutral   Result = "neutral"
	ResultPass      Result = "pass"
	ResultFail      Result = "fail"
	ResultSoftFail  Result = "softfail"
	ResultTempError Result = "temperror"
	ResultPermError Result = "permerror"
)

type Validator struct {
	resolver *Resolver
}

func NewValidator(resolver *Resolver) *Validator {
	return &Validator{resolver: resolver}
}

func (v *Validator) Check(ip net.IP, domain, sender string) (Result, error) {
	record, err := v.resolver.LookupSPF(domain)
	if err != nil {
		if err == ErrNoSPFRecord {
			return ResultNone, nil
		}
		return ResultTempError, err
	}

	return v.evaluate(record, ip, domain, sender)
}

// Mechanism represents a single SPF mechanism (e.g., "a", "mx", "ip4:...")
type Mechanism struct {
	Qualifier Result
	Type      string
	Value     string
}

func (v *Validator) evaluate(record string, ip net.IP, domain, _ string) (Result, error) {
	mechanisms := parseSPF(record)

	for _, mech := range mechanisms {
		match, err := v.matchMechanism(mech, ip, domain)
		if err != nil {
			continue // Or handle error appropriately
		}
		if match {
			return mech.Qualifier, nil
		}
	}

	return ResultNeutral, nil
}

func parseSPF(record string) []Mechanism {
	var mechanisms []Mechanism
	parts := strings.Fields(record)

	for _, part := range parts {
		if strings.HasPrefix(part, "v=") {
			continue // Skip version directive
		}

		qualifier := ResultPass
		directive := part

		// Check for qualifier prefix
		switch part[0] {
		case '+':
			qualifier = ResultPass
			directive = part[1:]
		case '-':
			qualifier = ResultFail
			directive = part[1:]
		case '~':
			qualifier = ResultSoftFail
			directive = part[1:]
		case '?':
			qualifier = ResultNeutral
			directive = part[1:]
		}

		// Parse mechanism type and value
		colonIdx := strings.Index(directive, ":")
		slashIdx := strings.Index(directive, "/")

		var mechType, value string
		if colonIdx > 0 {
			mechType = directive[:colonIdx]
			if slashIdx > colonIdx {
				value = directive[colonIdx+1:]
			} else {
				value = directive[colonIdx+1:]
			}
		} else if slashIdx > 0 {
			mechType = directive[:slashIdx]
			value = directive[slashIdx+1:]
		} else {
			mechType = directive
			value = ""
		}

		mechanisms = append(mechanisms, Mechanism{
			Qualifier: qualifier,
			Type:      strings.ToLower(mechType),
			Value:     value,
		})
	}

	return mechanisms
}

func (v *Validator) matchMechanism(mech Mechanism, ip net.IP, domain string) (bool, error) {
	switch mech.Type {
	case "all":
		return true, nil

	case "ip4":
		return v.matchIP4(ip, mech.Value)

	case "ip6":
		return v.matchIP6(ip, mech.Value)

	case "a":
		target := domain
		if mech.Value != "" {
			target = mech.Value
		}
		return v.matchA(ip, target)

	case "mx":
		target := domain
		if mech.Value != "" {
			target = mech.Value
		}
		return v.matchMX(ip, target)

	case "ptr":
		target := domain
		if mech.Value != "" {
			target = mech.Value
		}
		return v.matchPTR(ip, target)

	case "exists":
		if mech.Value == "" {
			return false, nil
		}
		return v.matchExists(mech.Value)

	case "include":
		if mech.Value == "" {
			return false, nil
		}
		result, err := v.Check(ip, mech.Value, "")
		if err != nil {
			return false, err
		}
		return result == ResultPass, nil

	default:
		return false, nil
	}
}

func (v *Validator) matchIP4(ip net.IP, cidr string) (bool, error) {
	if ip.To4() == nil {
		return false, nil // Not an IPv4 address
	}

	if !strings.Contains(cidr, "/") {
		cidr = cidr + "/32"
	}

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	return ipNet.Contains(ip), nil
}

func (v *Validator) matchIP6(ip net.IP, cidr string) (bool, error) {
	if ip.To4() != nil {
		return false, nil // Not an IPv6 address
	}

	if !strings.Contains(cidr, "/") {
		cidr = cidr + "/128"
	}

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	return ipNet.Contains(ip), nil
}

func (v *Validator) matchA(ip net.IP, domain string) (bool, error) {
	ips, err := v.resolver.LookupA(domain)
	if err != nil {
		return false, err
	}

	for _, resolvedIP := range ips {
		if resolvedIP.Equal(ip) {
			return true, nil
		}
	}

	return false, nil
}

func (v *Validator) matchMX(ip net.IP, domain string) (bool, error) {
	mxRecords, err := v.resolver.LookupMX(domain)
	if err != nil {
		return false, err
	}

	for _, mx := range mxRecords {
		ips, err := v.resolver.LookupA(mx)
		if err != nil {
			continue
		}
		for _, resolvedIP := range ips {
			if resolvedIP.Equal(ip) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (v *Validator) matchPTR(ip net.IP, domain string) (bool, error) {
	names, err := v.resolver.LookupPTR(ip)
	if err != nil {
		return false, err
	}

	for _, name := range names {
		if strings.HasSuffix(strings.ToLower(name), strings.ToLower(domain)) {
			// Verify with forward lookup
			ips, err := v.resolver.LookupA(name)
			if err != nil {
				continue
			}
			for _, resolvedIP := range ips {
				if resolvedIP.Equal(ip) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (v *Validator) matchExists(domain string) (bool, error) {
	ips, err := v.resolver.LookupA(domain)
	if err != nil {
		return false, nil // No error, just no match
	}
	return len(ips) > 0, nil
}
