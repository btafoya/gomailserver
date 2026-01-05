package service

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// DMARCAnalyzerService analyzes DMARC reports and identifies issues

type DMARCAnalyzerService struct {
	reportsRepo repository.DMARCReportsRepository
	actionsRepo repository.DMARCActionsRepository
	alertsRepo  repository.AlertsRepository
	logger      *zap.Logger
}

func NewDMARCAnalyzerService(
	reportsRepo repository.DMARCReportsRepository,
	actionsRepo repository.DMARCActionsRepository,
	alertsRepo repository.AlertsRepository,
	logger *zap.Logger,
) *DMARCAnalyzerService {
	return &DMARCAnalyzerService{
		reportsRepo: reportsRepo,
		actionsRepo: actionsRepo,
		alertsRepo:  alertsRepo,
		logger:      logger,
	}
}

// AnalyzeReport performs alignment analysis on a DMARC report
func (s *DMARCAnalyzerService) AnalyzeReport(ctx context.Context, report *domain.DMARCReport) (*domain.AlignmentAnalysis, error) {
	analysis := &domain.AlignmentAnalysis{
		Domain:       report.Domain,
		ReportID:     report.ID,
		TotalMessages: report.TotalMessages,
		Issues:       make([]*domain.AlignmentIssue, 0),
		RecommendedActions: make([]string, 0),
	}

	if report.TotalMessages == 0 {
		return analysis, nil
	}

	// Calculate rates
	analysis.AlignmentPassRate = float64(report.AlignmentPass) / float64(report.TotalMessages)
	analysis.SPFPassRate = float64(report.SPFPass) / float64(report.TotalMessages)
	analysis.DKIMPassRate = float64(report.DKIMPass) / float64(report.TotalMessages)

	// Calculate alignment rates from records
	spfAlignedCount := 0
	dkimAlignedCount := 0

	for _, record := range report.Records {
		if record.SPFAligned {
			spfAlignedCount += record.Count
		}
		if record.DKIMAligned {
			dkimAlignedCount += record.Count
		}
	}

	analysis.SPFAlignmentRate = float64(spfAlignedCount) / float64(report.TotalMessages)
	analysis.DKIMAlignmentRate = float64(dkimAlignedCount) / float64(report.TotalMessages)

	// Identify issues
	s.identifyIssues(report, analysis)

	// Generate recommendations
	s.generateRecommendations(analysis)

	s.logger.Info("Completed DMARC analysis",
		zap.String("domain", report.Domain),
		zap.Float64("alignment_rate", analysis.AlignmentPassRate),
		zap.Int("issues_found", len(analysis.Issues)),
	)

	return analysis, nil
}

// identifyIssues identifies specific alignment problems
func (s *DMARCAnalyzerService) identifyIssues(report *domain.DMARCReport, analysis *domain.AlignmentAnalysis) {
	// Track issues by source IP to avoid duplicates
	issuesBySeverity := make(map[string]*domain.AlignmentIssue)

	for _, record := range report.Records {
		// SPF failures
		if record.SPFResult != "pass" {
			key := fmt.Sprintf("spf_%s_%s", record.SourceIP, record.SPFResult)
			if issue, exists := issuesBySeverity[key]; exists {
				issue.Count += record.Count
			} else {
				severity := s.getSeverity(record.SPFResult, record.Count, report.TotalMessages)
				issuesBySeverity[key] = &domain.AlignmentIssue{
					IssueType:   domain.IssueTypeSPFFail,
					SourceIP:    record.SourceIP,
					Count:       record.Count,
					Description: fmt.Sprintf("SPF %s from %s", record.SPFResult, record.SourceIP),
					Severity:    severity,
				}
			}
		}

		// DKIM failures
		if record.DKIMResult != "pass" {
			key := fmt.Sprintf("dkim_%s_%s", record.SourceIP, record.DKIMResult)
			if issue, exists := issuesBySeverity[key]; exists {
				issue.Count += record.Count
			} else {
				severity := s.getSeverity(record.DKIMResult, record.Count, report.TotalMessages)
				issuesBySeverity[key] = &domain.AlignmentIssue{
					IssueType:   domain.IssueTypeDKIMFail,
					SourceIP:    record.SourceIP,
					Count:       record.Count,
					Description: fmt.Sprintf("DKIM %s from %s", record.DKIMResult, record.SourceIP),
					Severity:    severity,
				}
			}
		}

		// SPF misalignment
		if !record.SPFAligned && record.SPFResult == "pass" {
			key := fmt.Sprintf("spf_misalign_%s", record.SourceIP)
			if issue, exists := issuesBySeverity[key]; exists {
				issue.Count += record.Count
			} else {
				severity := s.getMisalignmentSeverity(record.Count, report.TotalMessages)
				issuesBySeverity[key] = &domain.AlignmentIssue{
					IssueType:   domain.IssueTypeSPFMisalign,
					SourceIP:    record.SourceIP,
					Count:       record.Count,
					Description: fmt.Sprintf("SPF passes but not aligned from %s (envelope: %s, header: %s)",
						record.SourceIP, record.EnvelopeFrom, record.HeaderFrom),
					Severity:    severity,
				}
			}
		}

		// DKIM misalignment
		if !record.DKIMAligned && record.DKIMResult == "pass" {
			key := fmt.Sprintf("dkim_misalign_%s", record.SourceIP)
			if issue, exists := issuesBySeverity[key]; exists {
				issue.Count += record.Count
			} else {
				severity := s.getMisalignmentSeverity(record.Count, report.TotalMessages)
				issuesBySeverity[key] = &domain.AlignmentIssue{
					IssueType:   domain.IssueTypeDKIMMisalign,
					SourceIP:    record.SourceIP,
					Count:       record.Count,
					Description: fmt.Sprintf("DKIM passes but not aligned from %s", record.SourceIP),
					Severity:    severity,
				}
			}
		}
	}

	// Convert map to slice
	for _, issue := range issuesBySeverity {
		analysis.Issues = append(analysis.Issues, issue)
	}
}

// getSeverity determines severity based on result type and volume
func (s *DMARCAnalyzerService) getSeverity(result string, count, total int) string {
	percentage := float64(count) / float64(total)

	// Permanent failures are always high severity
	if result == "permerror" || result == "fail" {
		if percentage > 0.1 {
			return "critical"
		}
		return "high"
	}

	// Temporary errors are medium severity
	if result == "temperror" || result == "softfail" {
		if percentage > 0.2 {
			return "high"
		}
		return "medium"
	}

	// Other issues are low severity
	if percentage > 0.3 {
		return "medium"
	}
	return "low"
}

// getMisalignmentSeverity determines severity for alignment issues
func (s *DMARCAnalyzerService) getMisalignmentSeverity(count, total int) string {
	percentage := float64(count) / float64(total)

	if percentage > 0.3 {
		return "high"
	} else if percentage > 0.1 {
		return "medium"
	}
	return "low"
}

// generateRecommendations creates actionable recommendations
func (s *DMARCAnalyzerService) generateRecommendations(analysis *domain.AlignmentAnalysis) {
	// Check overall alignment rate
	if analysis.AlignmentPassRate < 0.9 {
		analysis.RecommendedActions = append(analysis.RecommendedActions,
			"Overall alignment rate is below 90%. Review SPF and DKIM configuration.")
	}

	// SPF-specific recommendations
	if analysis.SPFPassRate < 0.95 {
		analysis.RecommendedActions = append(analysis.RecommendedActions,
			"SPF pass rate is low. Verify all sending IPs are included in SPF record.")
	}

	if analysis.SPFAlignmentRate < 0.9 && analysis.SPFPassRate > 0.9 {
		analysis.RecommendedActions = append(analysis.RecommendedActions,
			"SPF passes but alignment fails. Check envelope sender (Return-Path) matches header From domain.")
	}

	// DKIM-specific recommendations
	if analysis.DKIMPassRate < 0.95 {
		analysis.RecommendedActions = append(analysis.RecommendedActions,
			"DKIM pass rate is low. Verify DKIM signing is enabled and keys are valid.")
	}

	if analysis.DKIMAlignmentRate < 0.9 && analysis.DKIMPassRate > 0.9 {
		analysis.RecommendedActions = append(analysis.RecommendedActions,
			"DKIM passes but alignment fails. Ensure DKIM d= domain matches header From domain.")
	}

	// Issue-specific recommendations
	for _, issue := range analysis.Issues {
		switch issue.IssueType {
		case domain.IssueTypeSPFFail:
			if issue.Severity == "critical" || issue.Severity == "high" {
				analysis.RecommendedActions = append(analysis.RecommendedActions,
					fmt.Sprintf("High volume SPF failures from %s. Add this IP to SPF record if legitimate.", issue.SourceIP))
			}

		case domain.IssueTypeDKIMFail:
			if issue.Severity == "critical" || issue.Severity == "high" {
				analysis.RecommendedActions = append(analysis.RecommendedActions,
					fmt.Sprintf("High volume DKIM failures from %s. Verify DKIM signing configuration.", issue.SourceIP))
			}

		case domain.IssueTypeSPFMisalign:
			if issue.Severity == "high" {
				analysis.RecommendedActions = append(analysis.RecommendedActions,
					fmt.Sprintf("SPF misalignment from %s. Check envelope sender configuration.", issue.SourceIP))
			}

		case domain.IssueTypeDKIMMisalign:
			if issue.Severity == "high" {
				analysis.RecommendedActions = append(analysis.RecommendedActions,
					fmt.Sprintf("DKIM misalignment from %s. Verify DKIM d= parameter matches From domain.", issue.SourceIP))
			}
		}
	}

	// If no specific issues, provide general guidance
	if len(analysis.RecommendedActions) == 0 {
		analysis.RecommendedActions = append(analysis.RecommendedActions,
			"DMARC alignment is good. Continue monitoring for any changes.")
	}
}

// CreateAlert creates an alert for DMARC issues
func (s *DMARCAnalyzerService) CreateAlert(ctx context.Context, analysis *domain.AlignmentAnalysis) error {
	// Only create alert if there are high/critical severity issues
	hasCriticalIssues := false
	for _, issue := range analysis.Issues {
		if issue.Severity == "critical" || issue.Severity == "high" {
			hasCriticalIssues = true
			break
		}
	}

	if !hasCriticalIssues {
		return nil
	}

	// Determine alert severity
	severity := domain.SeverityMedium
	for _, issue := range analysis.Issues {
		if issue.Severity == "critical" {
			severity = domain.SeverityCritical
			break
		} else if issue.Severity == "high" {
			severity = domain.SeverityHigh
		}
	}

	alert := &domain.ReputationAlert{
		Domain:    analysis.Domain,
		AlertType: domain.AlertTypeDMARCIssue,
		Severity:  severity,
		Title:     fmt.Sprintf("DMARC Alignment Issues Detected for %s", analysis.Domain),
		Message:   fmt.Sprintf("Alignment pass rate: %.1f%% (%d issues found)",
			analysis.AlignmentPassRate*100, len(analysis.Issues)),
		Details: map[string]interface{}{
			"report_id":           analysis.ReportID,
			"total_messages":      analysis.TotalMessages,
			"alignment_pass_rate": analysis.AlignmentPassRate,
			"spf_pass_rate":       analysis.SPFPassRate,
			"dkim_pass_rate":      analysis.DKIMPassRate,
			"issues_count":        len(analysis.Issues),
			"recommendations":     analysis.RecommendedActions,
		},
		CreatedAt: time.Now().Unix(),
	}

	return s.alertsRepo.Create(ctx, alert)
}

// AnalyzeDomain analyzes all DMARC reports for a domain over a time period
func (s *DMARCAnalyzerService) AnalyzeDomain(ctx context.Context, domain string, days int) (*domain.AlignmentAnalysis, error) {
	return s.reportsRepo.GetDomainStats(ctx, domain, days)
}
