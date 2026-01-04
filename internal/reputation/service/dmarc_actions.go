package service

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// DMARCActionsService handles automated corrective actions for DMARC issues

type DMARCActionsService struct {
	actionsRepo repository.DMARCActionsRepository
	alertsRepo  repository.AlertsRepository
	logger      *zap.Logger
}

func NewDMARCActionsService(
	actionsRepo repository.DMARCActionsRepository,
	alertsRepo repository.AlertsRepository,
	logger *zap.Logger,
) *DMARCActionsService {
	return &DMARCActionsService{
		actionsRepo: actionsRepo,
		alertsRepo:  alertsRepo,
		logger:      logger,
	}
}

// TakeCorrectiveAction analyzes an issue and takes automated action
func (s *DMARCActionsService) TakeCorrectiveAction(ctx context.Context, issue *domain.AlignmentIssue, domainName string) error {
	var action *domain.DMARCAutoAction

	switch issue.IssueType {
	case domain.IssueTypeSPFMisalign:
		action = s.handleSPFMisalignment(ctx, issue, domainName)
	case domain.IssueTypeDKIMMisalign:
		action = s.handleDKIMMisalignment(ctx, issue, domainName)
	case domain.IssueTypeSPFFail:
		action = s.handleSPFFailure(ctx, issue, domainName)
	case domain.IssueTypeDKIMFail:
		action = s.handleDKIMFailure(ctx, issue, domainName)
	default:
		s.logger.Warn("Unknown issue type", zap.String("type", string(issue.IssueType)))
		return nil
	}

	if action != nil {
		if err := s.actionsRepo.RecordAction(ctx, action); err != nil {
			s.logger.Error("Failed to record action", zap.Error(err))
			return err
		}

		s.logger.Info("Recorded corrective action",
			zap.String("domain", domainName),
			zap.String("issue", string(action.IssueType)),
			zap.String("action", action.ActionTaken),
			zap.Bool("success", action.Success),
		)
	}

	return nil
}

func (s *DMARCActionsService) handleSPFMisalignment(ctx context.Context, issue *domain.AlignmentIssue, domainName string) *domain.DMARCAutoAction {
	// SPF misalignment means SPF passes but doesn't align with From domain
	// This requires manual intervention - log recommendation
	return &domain.DMARCAutoAction{
		Domain:      domainName,
		IssueType:   domain.IssueTypeSPFMisalign,
		Description: issue.Description,
		ActionTaken: "Logged recommendation: Review envelope sender (Return-Path) configuration to match From domain",
		TakenAt:     time.Now().Unix(),
		Success:     true,
	}
}

func (s *DMARCActionsService) handleDKIMMisalignment(ctx context.Context, issue *domain.AlignmentIssue, domainName string) *domain.DMARCAutoAction {
	// DKIM misalignment means DKIM passes but d= doesn't match From domain
	return &domain.DMARCAutoAction{
		Domain:      domainName,
		IssueType:   domain.IssueTypeDKIMMisalign,
		Description: issue.Description,
		ActionTaken: "Logged recommendation: Ensure DKIM d= parameter matches From domain in signing configuration",
		TakenAt:     time.Now().Unix(),
		Success:     true,
	}
}

func (s *DMARCActionsService) handleSPFFailure(ctx context.Context, issue *domain.AlignmentIssue, domainName string) *domain.DMARCAutoAction {
	// SPF failure means sending IP not authorized
	return &domain.DMARCAutoAction{
		Domain:      domainName,
		IssueType:   domain.IssueTypeSPFFail,
		Description: issue.Description,
		ActionTaken: fmt.Sprintf("Logged recommendation: Add IP %s to SPF record if legitimate sending source", issue.SourceIP),
		TakenAt:     time.Now().Unix(),
		Success:     true,
	}
}

func (s *DMARCActionsService) handleDKIMFailure(ctx context.Context, issue *domain.AlignmentIssue, domainName string) *domain.DMARCAutoAction {
	// DKIM failure could be due to key issues or misconfiguration
	return &domain.DMARCAutoAction{
		Domain:      domainName,
		IssueType:   domain.IssueTypeDKIMFail,
		Description: issue.Description,
		ActionTaken: "Logged recommendation: Verify DKIM keys are published and signing is enabled",
		TakenAt:     time.Now().Unix(),
		Success:     true,
	}
}

// ProcessAnalysis processes a full analysis and takes appropriate actions
func (s *DMARCActionsService) ProcessAnalysis(ctx context.Context, analysis *domain.AlignmentAnalysis) error {
	for _, issue := range analysis.Issues {
		// Only take action on high/critical severity issues
		if issue.Severity == "high" || issue.Severity == "critical" {
			if err := s.TakeCorrectiveAction(ctx, issue, analysis.Domain); err != nil {
				s.logger.Error("Failed to take corrective action",
					zap.Error(err),
					zap.String("issue", string(issue.IssueType)),
				)
				// Continue processing other issues
			}
		}
	}

	return nil
}

// ListActions returns recent actions for a domain
func (s *DMARCActionsService) ListActions(ctx context.Context, domain string, limit int) ([]*domain.DMARCAutoAction, error) {
	return s.actionsRepo.ListActions(ctx, domain, limit)
}

// ListAllActions returns all actions with pagination
func (s *DMARCActionsService) ListAllActions(ctx context.Context, limit, offset int) ([]*domain.DMARCAutoAction, error) {
	return s.actionsRepo.ListAllActions(ctx, limit, offset)
}
