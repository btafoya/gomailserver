package service

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

// DANEServiceConfig holds DNS configuration for the DANE service
type DANEServiceConfig struct {
	Resolver        string
	FallbackServers []string
	Timeout         int
	UseTCP          bool
}

// DefaultDANEServiceConfig returns the default DNS configuration
func DefaultDANEServiceConfig() *DANEServiceConfig {
	return &DANEServiceConfig{
		Resolver:        "1.1.1.1:53",
		FallbackServers: []string{"8.8.8.8:53", "9.9.9.9:53"},
		Timeout:         5,
		UseTCP:          false,
	}
}

// DANEService handles DANE TLSA record lookups and caching
type DANEService struct {
	db              *database.DB
	logger          *zap.Logger
	dnsClient       *dns.Client
	resolver        string
	fallbackServers []string
}

// NewDANEService creates a new DANE service with default configuration
func NewDANEService(db *database.DB, logger *zap.Logger) *DANEService {
	return NewDANEServiceWithConfig(db, logger, nil)
}

// NewDANEServiceWithConfig creates a new DANE service with custom configuration
func NewDANEServiceWithConfig(db *database.DB, logger *zap.Logger, cfg *DANEServiceConfig) *DANEService {
	if cfg == nil {
		cfg = DefaultDANEServiceConfig()
	}

	// Normalize resolver address
	resolver := normalizeResolverAddress(cfg.Resolver)

	// Normalize fallback servers
	fallbacks := make([]string, 0, len(cfg.FallbackServers))
	for _, server := range cfg.FallbackServers {
		fallbacks = append(fallbacks, normalizeResolverAddress(server))
	}

	client := &dns.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}
	if cfg.UseTCP {
		client.Net = "tcp"
	}

	logger.Info("DANE service initialized",
		zap.String("resolver", resolver),
		zap.Strings("fallback_servers", fallbacks),
		zap.Int("timeout_seconds", cfg.Timeout),
		zap.Bool("use_tcp", cfg.UseTCP),
	)

	return &DANEService{
		db:              db,
		logger:          logger,
		dnsClient:       client,
		resolver:        resolver,
		fallbackServers: fallbacks,
	}
}

// normalizeResolverAddress ensures the resolver has a port
func normalizeResolverAddress(resolver string) string {
	host, port, err := net.SplitHostPort(resolver)
	if err != nil {
		// No port specified, add default DNS port
		return net.JoinHostPort(resolver, "53")
	}
	return net.JoinHostPort(host, port)
}

// LookupTLSA performs a DANE TLSA DNS lookup with caching
func (s *DANEService) LookupTLSA(ctx context.Context, domainName string, port int) ([]*domain.DANETLSARecord, error) {
	// Check cache first
	cached, err := s.getCachedRecords(ctx, domainName, port)
	if err == nil && len(cached) > 0 {
		// Verify cache hasn't expired
		if !s.isCacheExpired(cached[0]) {
			s.logger.Debug("DANE TLSA cache hit",
				zap.String("domain", domainName),
				zap.Int("port", port),
			)
			return cached, nil
		}
	}

	// Perform DNS lookup
	records, err := s.fetchTLSARecords(ctx, domainName, port)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLSA records: %w", err)
	}

	// Cache the results
	for _, record := range records {
		if err := s.cacheRecord(ctx, record); err != nil {
			s.logger.Error("failed to cache TLSA record",
				zap.Error(err),
				zap.String("domain", domainName),
			)
		}
	}

	return records, nil
}

// fetchTLSARecords performs the actual DNS query for TLSA records
func (s *DANEService) fetchTLSARecords(ctx context.Context, domainName string, port int) ([]*domain.DANETLSARecord, error) {
	// Construct TLSA query name: _port._tcp.domain
	queryName := fmt.Sprintf("_%d._tcp.%s.", port, domainName)

	msg := &dns.Msg{}
	msg.SetQuestion(queryName, dns.TypeTLSA)
	msg.SetEdns0(4096, true) // Request DNSSEC

	// Try primary resolver, then fallbacks
	resolvers := append([]string{s.resolver}, s.fallbackServers...)
	var lastErr error
	var resp *dns.Msg

	for _, resolver := range resolvers {
		var err error
		resp, _, err = s.dnsClient.Exchange(msg, resolver)
		if err == nil {
			break
		}
		lastErr = err
		s.logger.Debug("DNS query failed, trying next resolver",
			zap.String("resolver", resolver),
			zap.Error(err),
		)
	}

	if resp == nil {
		return nil, fmt.Errorf("DNS query failed on all resolvers: %w", lastErr)
	}

	if resp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS query returned error: %s", dns.RcodeToString[resp.Rcode])
	}

	// Check DNSSEC authentication
	dnssecVerified := resp.AuthenticatedData

	var records []*domain.DANETLSARecord
	now := time.Now()

	for _, answer := range resp.Answer {
		tlsa, ok := answer.(*dns.TLSA)
		if !ok {
			continue
		}

		certData := hex.EncodeToString([]byte(tlsa.Certificate))

		record := &domain.DANETLSARecord{
			Domain:          domainName,
			Port:            port,
			Usage:           int(tlsa.Usage),
			Selector:        int(tlsa.Selector),
			MatchingType:    int(tlsa.MatchingType),
			CertificateData: certData,
			FetchedAt:       now,
			TTL:             int(tlsa.Hdr.Ttl),
			DNSSECVerified:  dnssecVerified,
		}

		records = append(records, record)
	}

	s.logger.Info("DANE TLSA records fetched",
		zap.String("domain", domainName),
		zap.Int("port", port),
		zap.Int("count", len(records)),
		zap.Bool("dnssec_verified", dnssecVerified),
	)

	return records, nil
}

// VerifyTLSConnection verifies a TLS connection against DANE TLSA records
func (s *DANEService) VerifyTLSConnection(ctx context.Context, domainName string, port int, tlsState *tls.ConnectionState) (bool, error) {
	records, err := s.LookupTLSA(ctx, domainName, port)
	if err != nil {
		return false, fmt.Errorf("failed to lookup TLSA records: %w", err)
	}

	if len(records) == 0 {
		// No DANE records, cannot verify
		return false, nil
	}

	// Get the peer certificate chain
	if len(tlsState.PeerCertificates) == 0 {
		return false, fmt.Errorf("no peer certificates in TLS connection")
	}

	// Verify against each TLSA record
	for _, record := range records {
		if s.verifyRecord(record, tlsState.PeerCertificates) {
			s.logger.Info("DANE verification successful",
				zap.String("domain", domainName),
				zap.Int("usage", record.Usage),
			)
			return true, nil
		}
	}

	return false, fmt.Errorf("no TLSA record matched the certificate chain")
}

// verifyRecord verifies a certificate chain against a TLSA record
func (s *DANEService) verifyRecord(record *domain.DANETLSARecord, certs []*x509.Certificate) bool {
	var cert *x509.Certificate

	// Select certificate based on usage
	switch record.Usage {
	case domain.TLSAUsageCAConstraint, domain.TLSAUsageTrustAnchor:
		// Use the root/CA certificate (last in chain)
		if len(certs) > 0 {
			cert = certs[len(certs)-1]
		}
	case domain.TLSAUsageServiceConstraint, domain.TLSAUsageDomainIssuedCert:
		// Use the leaf certificate (first in chain)
		cert = certs[0]
	default:
		return false
	}

	if cert == nil {
		return false
	}

	// Extract data based on selector
	var data []byte
	switch record.Selector {
	case domain.TLSASelectorFullCert:
		data = cert.Raw
	case domain.TLSASelectorSubjectPublicKeyInfo:
		data = cert.RawSubjectPublicKeyInfo
	default:
		return false
	}

	// Hash data based on matching type
	var hash string
	switch record.MatchingType {
	case domain.TLSAMatchingFull:
		hash = hex.EncodeToString(data)
	case domain.TLSAMatchingSHA256:
		h := sha256.Sum256(data)
		hash = hex.EncodeToString(h[:])
	case domain.TLSAMatchingSHA512:
		h := sha512.Sum512(data)
		hash = hex.EncodeToString(h[:])
	default:
		return false
	}

	return hash == record.CertificateData
}

// getCachedRecords retrieves cached TLSA records
func (s *DANEService) getCachedRecords(ctx context.Context, domainName string, port int) ([]*domain.DANETLSARecord, error) {
	query := `
		SELECT id, domainName, port, usage, selector, matching_type,
		       certificate_data, fetched_at, ttl, dnssec_verified
		FROM dane_tlsa_cache
		WHERE domain = ? AND port = ?
	`

	rows, err := s.db.QueryContext(ctx, query, domainName, port)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.DANETLSARecord
	for rows.Next() {
		record := &domain.DANETLSARecord{}
		err := rows.Scan(
			&record.ID,
			&record.Domain,
			&record.Port,
			&record.Usage,
			&record.Selector,
			&record.MatchingType,
			&record.CertificateData,
			&record.FetchedAt,
			&record.TTL,
			&record.DNSSECVerified,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

// cacheRecord stores a TLSA record in the cache
func (s *DANEService) cacheRecord(ctx context.Context, record *domain.DANETLSARecord) error {
	query := `
		INSERT OR REPLACE INTO dane_tlsa_cache (
			domainName, port, usage, selector, matching_type,
			certificate_data, fetched_at, ttl, dnssec_verified
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		record.Domain,
		record.Port,
		record.Usage,
		record.Selector,
		record.MatchingType,
		record.CertificateData,
		record.FetchedAt,
		record.TTL,
		record.DNSSECVerified,
	)

	return err
}

// isCacheExpired checks if a cached record has expired
func (s *DANEService) isCacheExpired(record *domain.DANETLSARecord) bool {
	expiresAt := record.FetchedAt.Add(time.Duration(record.TTL) * time.Second)
	return time.Now().After(expiresAt)
}

// ClearCache removes expired TLSA records from the cache
func (s *DANEService) ClearCache(ctx context.Context) (int64, error) {
	// Delete records where fetched_at + ttl < now
	query := `
		DELETE FROM dane_tlsa_cache
		WHERE datetime(fetched_at, '+' || ttl || ' seconds') < datetime('now')
	`

	result, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to clear DANE cache: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if deleted > 0 {
		s.logger.Info("cleared expired DANE TLSA records",
			zap.Int64("deleted", deleted),
		)
	}

	return deleted, nil
}

// SetResolver sets the DNS resolver address
func (s *DANEService) SetResolver(resolver string) {
	host, port, err := net.SplitHostPort(resolver)
	if err != nil {
		// If no port specified, add default DNS port
		s.resolver = net.JoinHostPort(resolver, "53")
	} else {
		s.resolver = net.JoinHostPort(host, port)
	}
	s.logger.Info("DANE DNS resolver updated", zap.String("resolver", s.resolver))
}
