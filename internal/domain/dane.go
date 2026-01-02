package domain

import "time"

// DANETLSARecord represents a cached DANE TLSA DNS record
type DANETLSARecord struct {
	ID               int64     `json:"id"`
	Domain           string    `json:"domain"`
	Port             int       `json:"port"`
	Usage            int       `json:"usage"`
	Selector         int       `json:"selector"`
	MatchingType     int       `json:"matching_type"`
	CertificateData  string    `json:"certificate_data"`
	FetchedAt        time.Time `json:"fetched_at"`
	TTL              int       `json:"ttl"`
	DNSSECVerified   bool      `json:"dnssec_verified"`
}

// DANE TLSA Usage types (RFC 6698)
const (
	TLSAUsageCAConstraint       = 0 // PKIX-TA
	TLSAUsageServiceConstraint  = 1 // PKIX-EE
	TLSAUsageTrustAnchor        = 2 // DANE-TA
	TLSAUsageDomainIssuedCert   = 3 // DANE-EE
)

// DANE TLSA Selector types
const (
	TLSASelectorFullCert        = 0
	TLSASelectorSubjectPublicKeyInfo = 1
)

// DANE TLSA Matching types
const (
	TLSAMatchingFull            = 0
	TLSAMatchingSHA256          = 1
	TLSAMatchingSHA512          = 2
)
