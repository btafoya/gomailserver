package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type dmarcActionsRepository struct {
	db *sql.DB
}

// NewDMARCActionsRepository creates a new SQLite DMARC actions repository
func NewDMARCActionsRepository(db *sql.DB) repository.DMARCActionsRepository {
	return &dmarcActionsRepository{db: db}
}

// RecordAction logs an automated action taken
func (r *dmarcActionsRepository) RecordAction(ctx context.Context, action *domain.DMARCAutoAction) error {
	query := `
		INSERT INTO dmarc_auto_actions (
			domain, issue_type, description, action_taken,
			taken_at, success, error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		action.Domain,
		action.IssueType,
		action.Description,
		action.ActionTaken,
		action.TakenAt,
		action.Success,
		action.ErrorMessage,
	)
	if err != nil {
		return fmt.Errorf("failed to record DMARC action: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	action.ID = id

	return nil
}

// ListActions returns recent automated actions
func (r *dmarcActionsRepository) ListActions(ctx context.Context, domainName string, limit int) ([]*domain.DMARCAutoAction, error) {
	query := `
		SELECT
			id, domain, issue_type, description, action_taken,
			taken_at, success, error_message
		FROM dmarc_auto_actions
		WHERE domain = ?
		ORDER BY taken_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list DMARC actions: %w", err)
	}
	defer rows.Close()

	var actions []*domain.DMARCAutoAction
	for rows.Next() {
		action := &domain.DMARCAutoAction{}
		err := rows.Scan(
			&action.ID,
			&action.Domain,
			&action.IssueType,
			&action.Description,
			&action.ActionTaken,
			&action.TakenAt,
			&action.Success,
			&action.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan DMARC action: %w", err)
		}
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating DMARC actions: %w", err)
	}

	return actions, nil
}

// ListAllActions returns all actions with pagination
func (r *dmarcActionsRepository) ListAllActions(ctx context.Context, limit, offset int) ([]*domain.DMARCAutoAction, error) {
	query := `
		SELECT
			id, domain, issue_type, description, action_taken,
			taken_at, success, error_message
		FROM dmarc_auto_actions
		ORDER BY taken_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list all DMARC actions: %w", err)
	}
	defer rows.Close()

	var actions []*domain.DMARCAutoAction
	for rows.Next() {
		action := &domain.DMARCAutoAction{}
		err := rows.Scan(
			&action.ID,
			&action.Domain,
			&action.IssueType,
			&action.Description,
			&action.ActionTaken,
			&action.TakenAt,
			&action.Success,
			&action.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan DMARC action: %w", err)
		}
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating DMARC actions: %w", err)
	}

	return actions, nil
}
