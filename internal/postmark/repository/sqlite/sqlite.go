package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/postmark/models"
	"github.com/btafoya/gomailserver/internal/postmark/repository"
	"golang.org/x/crypto/bcrypt"
)

// Repository implements PostmarkRepository using SQLite
type Repository struct {
	db *sql.DB
}

// New creates a new SQLite PostmarkRepository
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateServer creates a new PostmarkApp server
func (r *Repository) CreateServer(ctx context.Context, server *models.Server) error {
	// Hash the API token
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(server.APIToken), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash API token: %w", err)
	}

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO postmark_servers (name, api_token, account_id, message_stream, track_opens, track_links, active)
		VALUES (?, ?, ?, ?, ?, ?, 1)
	`, server.Name, string(hashedToken), server.ID, server.ServerLink, server.TrackOpens, server.TrackLinks)

	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	server.ID = int(id)
	return nil
}

// GetServer retrieves a server by ID
func (r *Repository) GetServer(ctx context.Context, id int) (*models.Server, error) {
	var server models.Server
	var trackOpens int
	var active int

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, account_id, message_stream, track_opens, track_links, active, created_at, updated_at
		FROM postmark_servers WHERE id = ?
	`, id).Scan(
		&server.ID,
		&server.Name,
		&server.ID, // account_id placeholder
		&server.ServerLink,
		&trackOpens,
		&server.TrackLinks,
		&active,
		&server.CreatedAt,
		&server.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	server.TrackOpens = trackOpens == 1
	server.SmtpApiActivated = active == 1

	return &server, nil
}

// GetServerByToken retrieves a server by API token
func (r *Repository) GetServerByToken(ctx context.Context, token string) (*models.Server, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, api_token, account_id, message_stream, track_opens, track_links, active, created_at, updated_at
		FROM postmark_servers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var server models.Server
		var hashedToken string
		var trackOpens int
		var active int

		err := rows.Scan(
			&server.ID,
			&server.Name,
			&hashedToken,
			&server.ID, // account_id placeholder
			&server.ServerLink,
			&trackOpens,
			&server.TrackLinks,
			&active,
			&server.CreatedAt,
			&server.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// Compare token with hashed version
		if bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token)) == nil {
			server.TrackOpens = trackOpens == 1
			server.SmtpApiActivated = active == 1
			return &server, nil
		}
	}

	return nil, sql.ErrNoRows
}

// UpdateServer updates a server's settings
func (r *Repository) UpdateServer(ctx context.Context, id int, req *models.UpdateServerRequest) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE postmark_servers
		SET name = COALESCE(?, name),
		    track_opens = COALESCE(?, track_opens),
		    track_links = COALESCE(?, track_links),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, req.Name, req.TrackOpens, req.TrackLinks, id)

	return err
}

// DeleteServer deletes a server
func (r *Repository) DeleteServer(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM postmark_servers WHERE id = ?`, id)
	return err
}

// ListServers lists all servers for an account
func (r *Repository) ListServers(ctx context.Context, accountID int) ([]*models.Server, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, account_id, message_stream, track_opens, track_links, active, created_at, updated_at
		FROM postmark_servers
		WHERE account_id = ?
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []*models.Server
	for rows.Next() {
		var server models.Server
		var trackOpens int
		var active int

		err := rows.Scan(
			&server.ID,
			&server.Name,
			&server.ID, // account_id placeholder
			&server.ServerLink,
			&trackOpens,
			&server.TrackLinks,
			&active,
			&server.CreatedAt,
			&server.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		server.TrackOpens = trackOpens == 1
		server.SmtpApiActivated = active == 1
		servers = append(servers, &server)
	}

	return servers, nil
}

// CreateMessage creates a new message record
func (r *Repository) CreateMessage(ctx context.Context, message *repository.Message) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO postmark_messages (
			message_id, server_id, from_email, to_email, cc_email, bcc_email,
			subject, html_body, text_body, tag, metadata, message_stream, status, submitted_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		message.MessageID, message.ServerID, message.FromEmail, message.ToEmail,
		message.CcEmail, message.BccEmail, message.Subject, message.HtmlBody,
		message.TextBody, message.Tag, message.Metadata, message.MessageStream,
		message.Status, message.SubmittedAt,
	)

	return err
}

// GetMessage retrieves a message by ID
func (r *Repository) GetMessage(ctx context.Context, messageID string) (*repository.Message, error) {
	var message repository.Message

	err := r.db.QueryRowContext(ctx, `
		SELECT id, message_id, server_id, from_email, to_email, cc_email, bcc_email,
		       subject, html_body, text_body, tag, metadata, message_stream, status,
		       submitted_at, sent_at, delivered_at
		FROM postmark_messages WHERE message_id = ?
	`, messageID).Scan(
		&message.ID, &message.MessageID, &message.ServerID, &message.FromEmail,
		&message.ToEmail, &message.CcEmail, &message.BccEmail, &message.Subject,
		&message.HtmlBody, &message.TextBody, &message.Tag, &message.Metadata,
		&message.MessageStream, &message.Status, &message.SubmittedAt,
		&message.SentAt, &message.DeliveredAt,
	)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

// UpdateMessageStatus updates a message's status
func (r *Repository) UpdateMessageStatus(ctx context.Context, messageID string, status string) error {
	now := time.Now().Format(time.RFC3339)

	_, err := r.db.ExecContext(ctx, `
		UPDATE postmark_messages
		SET status = ?,
		    sent_at = CASE WHEN ? = 'sent' THEN ? ELSE sent_at END,
		    delivered_at = CASE WHEN ? = 'delivered' THEN ? ELSE delivered_at END
		WHERE message_id = ?
	`, status, status, now, status, now, messageID)

	return err
}

// ListMessages lists messages for a server
func (r *Repository) ListMessages(ctx context.Context, serverID int, limit, offset int) ([]*repository.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, message_id, server_id, from_email, to_email, subject, status, submitted_at
		FROM postmark_messages
		WHERE server_id = ?
		ORDER BY submitted_at DESC
		LIMIT ? OFFSET ?
	`, serverID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*repository.Message
	for rows.Next() {
		var message repository.Message
		err := rows.Scan(
			&message.ID, &message.MessageID, &message.ServerID, &message.FromEmail,
			&message.ToEmail, &message.Subject, &message.Status, &message.SubmittedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// CreateTemplate creates a new template
func (r *Repository) CreateTemplate(ctx context.Context, template *models.Template) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO postmark_templates (
			server_id, name, alias, subject, html_body, text_body, template_type, layout_template, active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)
	`,
		template.ServerID, template.Name, template.Alias, template.Subject,
		template.HtmlBody, template.TextBody, template.TemplateType, template.LayoutTemplate,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	template.TemplateID = int(id)
	return nil
}

// GetTemplate retrieves a template by ID
func (r *Repository) GetTemplate(ctx context.Context, id int) (*models.Template, error) {
	var template models.Template
	var active int

	err := r.db.QueryRowContext(ctx, `
		SELECT id, server_id, name, alias, subject, html_body, text_body,
		       template_type, layout_template, active, created_at, updated_at
		FROM postmark_templates WHERE id = ?
	`, id).Scan(
		&template.TemplateID, &template.ServerID, &template.Name, &template.Alias,
		&template.Subject, &template.HtmlBody, &template.TextBody, &template.TemplateType,
		&template.LayoutTemplate, &active, &template.CreatedAt, &template.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	template.Active = active == 1
	return &template, nil
}

// GetTemplateByAlias retrieves a template by alias
func (r *Repository) GetTemplateByAlias(ctx context.Context, serverID int, alias string) (*models.Template, error) {
	var template models.Template
	var active int

	err := r.db.QueryRowContext(ctx, `
		SELECT id, server_id, name, alias, subject, html_body, text_body,
		       template_type, layout_template, active, created_at, updated_at
		FROM postmark_templates WHERE server_id = ? AND alias = ?
	`, serverID, alias).Scan(
		&template.TemplateID, &template.ServerID, &template.Name, &template.Alias,
		&template.Subject, &template.HtmlBody, &template.TextBody, &template.TemplateType,
		&template.LayoutTemplate, &active, &template.CreatedAt, &template.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	template.Active = active == 1
	return &template, nil
}

// UpdateTemplate updates a template
func (r *Repository) UpdateTemplate(ctx context.Context, id int, req *models.UpdateTemplateRequest) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE postmark_templates
		SET name = COALESCE(?, name),
		    alias = COALESCE(?, alias),
		    subject = COALESCE(?, subject),
		    html_body = COALESCE(?, html_body),
		    text_body = COALESCE(?, text_body),
		    template_type = COALESCE(?, template_type),
		    layout_template = COALESCE(?, layout_template),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		req.Name, req.Alias, req.Subject, req.HtmlBody, req.TextBody,
		req.TemplateType, req.LayoutTemplate, id,
	)

	return err
}

// DeleteTemplate deletes a template
func (r *Repository) DeleteTemplate(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM postmark_templates WHERE id = ?`, id)
	return err
}

// ListTemplates lists all templates for a server
func (r *Repository) ListTemplates(ctx context.Context, serverID int) ([]*models.Template, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, server_id, name, alias, subject, template_type, active, created_at, updated_at
		FROM postmark_templates
		WHERE server_id = ?
		ORDER BY name
	`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*models.Template
	for rows.Next() {
		var template models.Template
		var active int

		err := rows.Scan(
			&template.TemplateID, &template.ServerID, &template.Name, &template.Alias,
			&template.Subject, &template.TemplateType, &active, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		template.Active = active == 1
		templates = append(templates, &template)
	}

	return templates, nil
}

// CreateWebhook creates a new webhook
func (r *Repository) CreateWebhook(ctx context.Context, serverID int, webhook *models.Webhook) error {
	triggersJSON, err := json.Marshal(webhook.Triggers)
	if err != nil {
		return err
	}

	headersJSON, err := json.Marshal(webhook.HttpHeaders)
	if err != nil {
		return err
	}

	var authUsername, authPassword *string
	if webhook.HttpAuth != nil {
		authUsername = &webhook.HttpAuth.Username
		authPassword = &webhook.HttpAuth.Password
	}

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO postmark_webhooks (
			server_id, url, message_stream, http_auth_username, http_auth_password,
			http_headers, triggers, active
		) VALUES (?, ?, ?, ?, ?, ?, ?, 1)
	`,
		serverID, webhook.URL, webhook.MessageStream, authUsername, authPassword,
		string(headersJSON), string(triggersJSON),
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	webhook.ID = int(id)
	return nil
}

// GetWebhook retrieves a webhook by ID
func (r *Repository) GetWebhook(ctx context.Context, id int) (*models.Webhook, error) {
	var webhook models.Webhook
	var triggersJSON, headersJSON string
	var authUsername, authPassword *string
	var active int

	err := r.db.QueryRowContext(ctx, `
		SELECT id, url, message_stream, http_auth_username, http_auth_password,
		       http_headers, triggers, active, created_at, updated_at
		FROM postmark_webhooks WHERE id = ?
	`, id).Scan(
		&webhook.ID, &webhook.URL, &webhook.MessageStream, &authUsername, &authPassword,
		&headersJSON, &triggersJSON, &active, &webhook.CreatedAt, &webhook.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(triggersJSON), &webhook.Triggers); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(headersJSON), &webhook.HttpHeaders); err != nil {
		return nil, err
	}

	if authUsername != nil && authPassword != nil {
		webhook.HttpAuth = &models.HttpAuth{
			Username: *authUsername,
			Password: *authPassword,
		}
	}

	return &webhook, nil
}

// UpdateWebhook updates a webhook
func (r *Repository) UpdateWebhook(ctx context.Context, id int, req *models.UpdateWebhookRequest) error {
	triggersJSON, err := json.Marshal(req.Triggers)
	if err != nil {
		return err
	}

	headersJSON, err := json.Marshal(req.HttpHeaders)
	if err != nil {
		return err
	}

	var authUsername, authPassword *string
	if req.HttpAuth != nil {
		authUsername = &req.HttpAuth.Username
		authPassword = &req.HttpAuth.Password
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE postmark_webhooks
		SET url = COALESCE(?, url),
		    message_stream = COALESCE(?, message_stream),
		    http_auth_username = COALESCE(?, http_auth_username),
		    http_auth_password = COALESCE(?, http_auth_password),
		    http_headers = COALESCE(?, http_headers),
		    triggers = COALESCE(?, triggers),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		req.URL, req.MessageStream, authUsername, authPassword,
		string(headersJSON), string(triggersJSON), id,
	)

	return err
}

// DeleteWebhook deletes a webhook
func (r *Repository) DeleteWebhook(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM postmark_webhooks WHERE id = ?`, id)
	return err
}

// ListWebhooks lists all webhooks for a server
func (r *Repository) ListWebhooks(ctx context.Context, serverID int) ([]*models.Webhook, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, url, message_stream, triggers, active, created_at, updated_at
		FROM postmark_webhooks
		WHERE server_id = ?
		ORDER BY created_at DESC
	`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []*models.Webhook
	for rows.Next() {
		var webhook models.Webhook
		var triggersJSON string
		var active int

		err := rows.Scan(
			&webhook.ID, &webhook.URL, &webhook.MessageStream, &triggersJSON,
			&active, &webhook.CreatedAt, &webhook.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(triggersJSON), &webhook.Triggers); err != nil {
			return nil, err
		}

		webhooks = append(webhooks, &webhook)
	}

	return webhooks, nil
}

// GetActiveWebhooks retrieves active webhooks for an event type
func (r *Repository) GetActiveWebhooks(ctx context.Context, serverID int, eventType string) ([]*models.Webhook, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, url, message_stream, http_auth_username, http_auth_password,
		       http_headers, triggers
		FROM postmark_webhooks
		WHERE server_id = ? AND active = 1
	`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []*models.Webhook
	for rows.Next() {
		var webhook models.Webhook
		var triggersJSON, headersJSON string
		var authUsername, authPassword *string

		err := rows.Scan(
			&webhook.ID, &webhook.URL, &webhook.MessageStream, &authUsername,
			&authPassword, &headersJSON, &triggersJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(triggersJSON), &webhook.Triggers); err != nil {
			continue
		}

		// Check if this webhook is enabled for the event type
		enabled := false
		switch eventType {
		case "Open":
			enabled = webhook.Triggers.Open != nil && webhook.Triggers.Open.Enabled
		case "Click":
			enabled = webhook.Triggers.Click != nil && webhook.Triggers.Click.Enabled
		case "Delivery":
			enabled = webhook.Triggers.Delivery != nil && webhook.Triggers.Delivery.Enabled
		case "Bounce":
			enabled = webhook.Triggers.Bounce != nil && webhook.Triggers.Bounce.Enabled
		case "SpamComplaint":
			enabled = webhook.Triggers.SpamComplaint != nil && webhook.Triggers.SpamComplaint.Enabled
		}

		if !enabled {
			continue
		}

		if err := json.Unmarshal([]byte(headersJSON), &webhook.HttpHeaders); err != nil {
			continue
		}

		if authUsername != nil && authPassword != nil {
			webhook.HttpAuth = &models.HttpAuth{
				Username: *authUsername,
				Password: *authPassword,
			}
		}

		webhooks = append(webhooks, &webhook)
	}

	return webhooks, nil
}

// CreateBounce creates a new bounce record
func (r *Repository) CreateBounce(ctx context.Context, bounce *repository.Bounce) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO postmark_bounces (
			message_id, type, type_code, email, bounced_at, details, inactive, can_activate
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		bounce.MessageID, bounce.Type, bounce.TypeCode, bounce.Email,
		bounce.BouncedAt, bounce.Details, bounce.Inactive, bounce.CanActivate,
	)

	return err
}

// GetBounce retrieves a bounce by ID
func (r *Repository) GetBounce(ctx context.Context, id int) (*repository.Bounce, error) {
	var bounce repository.Bounce

	err := r.db.QueryRowContext(ctx, `
		SELECT id, message_id, type, type_code, email, bounced_at, details, inactive, can_activate
		FROM postmark_bounces WHERE id = ?
	`, id).Scan(
		&bounce.ID, &bounce.MessageID, &bounce.Type, &bounce.TypeCode,
		&bounce.Email, &bounce.BouncedAt, &bounce.Details, &bounce.Inactive, &bounce.CanActivate,
	)

	if err != nil {
		return nil, err
	}

	return &bounce, nil
}

// ListBounces lists bounces for a server
func (r *Repository) ListBounces(ctx context.Context, serverID int, limit, offset int) ([]*repository.Bounce, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT b.id, b.message_id, b.type, b.type_code, b.email, b.bounced_at, b.inactive
		FROM postmark_bounces b
		JOIN postmark_messages m ON b.message_id = m.message_id
		WHERE m.server_id = ?
		ORDER BY b.bounced_at DESC
		LIMIT ? OFFSET ?
	`, serverID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bounces []*repository.Bounce
	for rows.Next() {
		var bounce repository.Bounce
		err := rows.Scan(
			&bounce.ID, &bounce.MessageID, &bounce.Type, &bounce.TypeCode,
			&bounce.Email, &bounce.BouncedAt, &bounce.Inactive,
		)
		if err != nil {
			return nil, err
		}
		bounces = append(bounces, &bounce)
	}

	return bounces, nil
}

// CreateEvent creates a new tracking event
func (r *Repository) CreateEvent(ctx context.Context, event *repository.Event) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO postmark_events (
			message_id, event_type, recipient, user_agent, client_info, location, link_url, occurred_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		event.MessageID, event.EventType, event.Recipient, event.UserAgent,
		event.ClientInfo, event.Location, event.LinkURL, event.OccurredAt,
	)

	return err
}

// ListEvents lists events for a message
func (r *Repository) ListEvents(ctx context.Context, messageID string, eventType string) ([]*repository.Event, error) {
	query := `
		SELECT id, message_id, event_type, recipient, user_agent, occurred_at
		FROM postmark_events
		WHERE message_id = ?
	`
	args := []interface{}{messageID}

	if eventType != "" {
		query += " AND event_type = ?"
		args = append(args, eventType)
	}

	query += " ORDER BY occurred_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*repository.Event
	for rows.Next() {
		var event repository.Event
		err := rows.Scan(
			&event.ID, &event.MessageID, &event.EventType, &event.Recipient,
			&event.UserAgent, &event.OccurredAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}
