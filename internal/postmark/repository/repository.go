package repository

import (
	"context"

	"github.com/btafoya/gomailserver/internal/postmark/models"
)

// PostmarkRepository defines the interface for PostmarkApp data operations
type PostmarkRepository interface {
	// Server operations
	CreateServer(ctx context.Context, server *models.Server) error
	GetServer(ctx context.Context, id int) (*models.Server, error)
	GetServerByToken(ctx context.Context, token string) (*models.Server, error)
	UpdateServer(ctx context.Context, id int, req *models.UpdateServerRequest) error
	DeleteServer(ctx context.Context, id int) error
	ListServers(ctx context.Context, accountID int) ([]*models.Server, error)

	// Message operations
	CreateMessage(ctx context.Context, message *Message) error
	GetMessage(ctx context.Context, messageID string) (*Message, error)
	UpdateMessageStatus(ctx context.Context, messageID string, status string) error
	ListMessages(ctx context.Context, serverID int, limit, offset int) ([]*Message, error)

	// Template operations
	CreateTemplate(ctx context.Context, template *models.Template) error
	GetTemplate(ctx context.Context, id int) (*models.Template, error)
	GetTemplateByAlias(ctx context.Context, serverID int, alias string) (*models.Template, error)
	UpdateTemplate(ctx context.Context, id int, req *models.UpdateTemplateRequest) error
	DeleteTemplate(ctx context.Context, id int) error
	ListTemplates(ctx context.Context, serverID int) ([]*models.Template, error)

	// Webhook operations
	CreateWebhook(ctx context.Context, serverID int, webhook *models.Webhook) error
	GetWebhook(ctx context.Context, id int) (*models.Webhook, error)
	UpdateWebhook(ctx context.Context, id int, req *models.UpdateWebhookRequest) error
	DeleteWebhook(ctx context.Context, id int) error
	ListWebhooks(ctx context.Context, serverID int) ([]*models.Webhook, error)
	GetActiveWebhooks(ctx context.Context, serverID int, eventType string) ([]*models.Webhook, error)

	// Bounce operations
	CreateBounce(ctx context.Context, bounce *Bounce) error
	GetBounce(ctx context.Context, id int) (*Bounce, error)
	ListBounces(ctx context.Context, serverID int, limit, offset int) ([]*Bounce, error)

	// Event operations
	CreateEvent(ctx context.Context, event *Event) error
	ListEvents(ctx context.Context, messageID string, eventType string) ([]*Event, error)
}

// Message represents a PostmarkApp message record
type Message struct {
	ID            int
	MessageID     string
	ServerID      int
	FromEmail     string
	ToEmail       string
	CcEmail       string
	BccEmail      string
	Subject       string
	HtmlBody      string
	TextBody      string
	Tag           string
	Metadata      string // JSON
	MessageStream string
	Status        string
	SubmittedAt   string
	SentAt        *string
	DeliveredAt   *string
}

// Bounce represents a bounce record
type Bounce struct {
	ID          int
	MessageID   string
	Type        string
	TypeCode    int
	Email       string
	BouncedAt   string
	Details     string
	Inactive    bool
	CanActivate bool
}

// Event represents a tracking event
type Event struct {
	ID         int
	MessageID  string
	EventType  string
	Recipient  string
	UserAgent  string
	ClientInfo string // JSON
	Location   string // JSON
	LinkURL    string
	OccurredAt string
}
