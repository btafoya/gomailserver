package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// DMARC XML Structure (RFC 7489)

type dmarcFeedback struct {
	XMLName         xml.Name        `xml:"feedback"`
	Version         string          `xml:"version"`
	ReportMetadata  reportMetadata  `xml:"report_metadata"`
	PolicyPublished policyPublished `xml:"policy_published"`
	Records         []dmarcRecord   `xml:"record"`
}

type reportMetadata struct {
	OrgName   string    `xml:"org_name"`
	Email     string    `xml:"email"`
	ExtraContact string `xml:"extra_contact_info"`
	ReportID  string    `xml:"report_id"`
	DateRange dateRange `xml:"date_range"`
}

type dateRange struct {
	Begin int64 `xml:"begin"`
	End   int64 `xml:"end"`
}

type policyPublished struct {
	Domain string `xml:"domain"`
	ADKIM  string `xml:"adkim"` // relaxed or strict
	ASPF   string `xml:"aspf"`  // relaxed or strict
	P      string `xml:"p"`     // none, quarantine, reject
	SP     string `xml:"sp"`    // subdomain policy
	Pct    int    `xml:"pct"`   // percentage
}

type dmarcRecord struct {
	Row        recordRow        `xml:"row"`
	Identifiers recordIdentifiers `xml:"identifiers"`
	AuthResults authResults      `xml:"auth_results"`
}

type recordRow struct {
	SourceIP   string       `xml:"source_ip"`
	Count      int          `xml:"count"`
	PolicyEval policyEval   `xml:"policy_evaluated"`
}

type policyEval struct {
	Disposition string `xml:"disposition"` // none, quarantine, reject
	DKIM        string `xml:"dkim"`        // pass, fail
	SPF         string `xml:"spf"`         // pass, fail
	Reason      []policyReason `xml:"reason"`
}

type policyReason struct {
	Type    string `xml:"type"`
	Comment string `xml:"comment"`
}

type recordIdentifiers struct {
	HeaderFrom   string `xml:"header_from"`
	EnvelopeFrom string `xml:"envelope_from"`
}

type authResults struct {
	DKIM []dkimAuthResult `xml:"dkim"`
	SPF  []spfAuthResult  `xml:"spf"`
}

type dkimAuthResult struct {
	Domain      string `xml:"domain"`
	Result      string `xml:"result"` // pass, fail, neutral, temperror, permerror
	HumanResult string `xml:"human_result"`
}

type spfAuthResult struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"` // pass, fail, neutral, softfail, temperror, permerror
	Scope  string `xml:"scope"`  // helo, mfrom
}

// DMARCParserService handles DMARC aggregate report parsing

type DMARCParserService struct {
	reportsRepo repository.DMARCReportsRepository
	actionsRepo repository.DMARCActionsRepository
	logger      *zap.Logger
}

func NewDMARCParserService(
	reportsRepo repository.DMARCReportsRepository,
	actionsRepo repository.DMARCActionsRepository,
	logger *zap.Logger,
) *DMARCParserService {
	return &DMARCParserService{
		reportsRepo: reportsRepo,
		actionsRepo: actionsRepo,
		logger:      logger,
	}
}

// ParseReport parses a DMARC aggregate report from XML
func (s *DMARCParserService) ParseReport(ctx context.Context, xmlData []byte) (*domain.DMARCReport, error) {
	var feedback dmarcFeedback
	if err := xml.Unmarshal(xmlData, &feedback); err != nil {
		s.logger.Error("Failed to parse DMARC XML", zap.Error(err))
		return nil, fmt.Errorf("invalid DMARC XML: %w", err)
	}

	// Convert to domain model
	report := &domain.DMARCReport{
		Domain:        feedback.PolicyPublished.Domain,
		ReportID:      feedback.ReportMetadata.ReportID,
		BeginTime:     feedback.ReportMetadata.DateRange.Begin,
		EndTime:       feedback.ReportMetadata.DateRange.End,
		Organization:  feedback.ReportMetadata.OrgName,
		TotalMessages: 0,
		SPFPass:       0,
		DKIMPass:      0,
		AlignmentPass: 0,
		RawXML:        string(xmlData),
		ProcessedAt:   time.Now().Unix(),
		Records:       make([]*domain.DMARCReportRecord, 0, len(feedback.Records)),
	}

	// Process each record
	for _, rec := range feedback.Records {
		spfResult := "fail"
		dkimResult := "fail"

		// Get SPF result (use first SPF auth result)
		if len(rec.AuthResults.SPF) > 0 {
			spfResult = rec.AuthResults.SPF[0].Result
		}

		// Get DKIM result (use first DKIM auth result)
		if len(rec.AuthResults.DKIM) > 0 {
			dkimResult = rec.AuthResults.DKIM[0].Result
		}

		// Determine alignment
		spfAligned := rec.Row.PolicyEval.SPF == "pass"
		dkimAligned := rec.Row.PolicyEval.DKIM == "pass"

		record := &domain.DMARCReportRecord{
			SourceIP:     rec.Row.SourceIP,
			Count:        rec.Row.Count,
			Disposition:  rec.Row.PolicyEval.Disposition,
			SPFResult:    spfResult,
			DKIMResult:   dkimResult,
			SPFAligned:   spfAligned,
			DKIMAligned:  dkimAligned,
			HeaderFrom:   rec.Identifiers.HeaderFrom,
			EnvelopeFrom: rec.Identifiers.EnvelopeFrom,
		}

		report.Records = append(report.Records, record)

		// Aggregate statistics
		report.TotalMessages += rec.Row.Count

		if spfResult == "pass" {
			report.SPFPass += rec.Row.Count
		}

		if dkimResult == "pass" {
			report.DKIMPass += rec.Row.Count
		}

		if spfAligned || dkimAligned {
			report.AlignmentPass += rec.Row.Count
		}
	}

	s.logger.Info("Parsed DMARC report",
		zap.String("domain", report.Domain),
		zap.String("report_id", report.ReportID),
		zap.Int("total_messages", report.TotalMessages),
		zap.Int("records", len(report.Records)),
	)

	return report, nil
}

// StoreReport saves a parsed report to the database
func (s *DMARCParserService) StoreReport(ctx context.Context, report *domain.DMARCReport) error {
	// Check if report already exists
	existing, err := s.reportsRepo.GetByReportID(ctx, report.ReportID)
	if err == nil && existing != nil {
		s.logger.Info("DMARC report already processed, skipping",
			zap.String("report_id", report.ReportID),
		)
		return nil
	}

	// Create report
	if err := s.reportsRepo.Create(ctx, report); err != nil {
		s.logger.Error("Failed to store DMARC report", zap.Error(err))
		return fmt.Errorf("failed to store report: %w", err)
	}

	// Create records
	for _, record := range report.Records {
		record.ReportID = report.ID
		if err := s.reportsRepo.CreateRecord(ctx, record); err != nil {
			s.logger.Error("Failed to store DMARC record", zap.Error(err))
			// Continue processing other records
		}
	}

	s.logger.Info("Stored DMARC report",
		zap.String("domain", report.Domain),
		zap.String("report_id", report.ReportID),
		zap.Int64("id", report.ID),
	)

	return nil
}

// ParseAndStore combines parsing and storage
func (s *DMARCParserService) ParseAndStore(ctx context.Context, xmlData []byte) error {
	report, err := s.ParseReport(ctx, xmlData)
	if err != nil {
		return err
	}

	return s.StoreReport(ctx, report)
}

// ParseFromReader parses DMARC report from an io.Reader
func (s *DMARCParserService) ParseFromReader(ctx context.Context, r io.Reader) (*domain.DMARCReport, error) {
	xmlData, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML data: %w", err)
	}

	return s.ParseReport(ctx, xmlData)
}

// GetReportStats returns statistics for a domain's DMARC reports
func (s *DMARCParserService) GetReportStats(ctx context.Context, domain string, days int) (*domain.AlignmentAnalysis, error) {
	return s.reportsRepo.GetDomainStats(ctx, domain, days)
}

// ListReports returns DMARC reports for a domain with pagination
func (s *DMARCParserService) ListReports(ctx context.Context, domain string, limit, offset int) ([]*domain.DMARCReport, error) {
	return s.reportsRepo.ListByDomain(ctx, domain, limit, offset)
}

// GetReportByID retrieves a specific report with its records
func (s *DMARCParserService) GetReportByID(ctx context.Context, id int64) (*domain.DMARCReport, error) {
	report, err := s.reportsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load records
	records, err := s.reportsRepo.GetRecordsByReportID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to load report records", zap.Error(err))
		// Return report without records
	} else {
		report.Records = records
	}

	return report, nil
}
