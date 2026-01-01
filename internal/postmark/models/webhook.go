package models

import "time"

// HttpAuth represents HTTP authentication for webhooks
type HttpAuth struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// WebhookTriggers represents which events trigger the webhook
type WebhookTriggers struct {
	Open               *WebhookOpenTrigger  `json:"Open,omitempty"`
	Click              *WebhookTrigger      `json:"Click,omitempty"`
	Delivery           *WebhookTrigger      `json:"Delivery,omitempty"`
	Bounce             *WebhookBounceTrigger `json:"Bounce,omitempty"`
	SpamComplaint      *WebhookBounceTrigger `json:"SpamComplaint,omitempty"`
	SubscriptionChange *WebhookTrigger       `json:"SubscriptionChange,omitempty"`
}

// WebhookTrigger represents a simple trigger
type WebhookTrigger struct {
	Enabled bool `json:"Enabled"`
}

// WebhookOpenTrigger represents an open event trigger
type WebhookOpenTrigger struct {
	Enabled            bool `json:"Enabled"`
	PostFirstOpenOnly  bool `json:"PostFirstOpenOnly"`
}

// WebhookBounceTrigger represents a bounce/spam event trigger
type WebhookBounceTrigger struct {
	Enabled        bool `json:"Enabled"`
	IncludeContent bool `json:"IncludeContent"`
}

// Webhook represents a PostmarkApp webhook
type Webhook struct {
	ID            int              `json:"ID,omitempty"`
	URL           string           `json:"Url"`
	MessageStream string           `json:"MessageStream"`
	HttpAuth      *HttpAuth        `json:"HttpAuth,omitempty"`
	HttpHeaders   []Header         `json:"HttpHeaders,omitempty"`
	Triggers      WebhookTriggers  `json:"Triggers"`
	CreatedAt     time.Time        `json:"CreatedAt,omitempty"`
	UpdatedAt     time.Time        `json:"UpdatedAt,omitempty"`
}

// WebhookListResponse represents a list of webhooks
type WebhookListResponse struct {
	Webhooks []Webhook `json:"Webhooks"`
}

// CreateWebhookRequest represents a webhook creation request
type CreateWebhookRequest struct {
	URL           string          `json:"Url"`
	MessageStream string          `json:"MessageStream,omitempty"`
	HttpAuth      *HttpAuth       `json:"HttpAuth,omitempty"`
	HttpHeaders   []Header        `json:"HttpHeaders,omitempty"`
	Triggers      WebhookTriggers `json:"Triggers"`
}

// UpdateWebhookRequest represents a webhook update request
type UpdateWebhookRequest struct {
	URL           string          `json:"Url,omitempty"`
	MessageStream string          `json:"MessageStream,omitempty"`
	HttpAuth      *HttpAuth       `json:"HttpAuth,omitempty"`
	HttpHeaders   []Header        `json:"HttpHeaders,omitempty"`
	Triggers      WebhookTriggers `json:"Triggers,omitempty"`
}

// WebhookEvent represents an event to be sent to webhooks
type WebhookEvent struct {
	RecordType  string            `json:"RecordType"` // Bounce, Delivery, Open, Click, SpamComplaint
	ServerID    int               `json:"ServerID,omitempty"`
	MessageID   string            `json:"MessageID"`
	Recipient   string            `json:"Recipient,omitempty"`
	Tag         string            `json:"Tag,omitempty"`
	Metadata    map[string]string `json:"Metadata,omitempty"`
	Subject     string            `json:"Subject,omitempty"`
	ReceivedAt  time.Time         `json:"ReceivedAt,omitempty"`

	// Bounce-specific fields
	Type        string    `json:"Type,omitempty"`
	TypeCode    int       `json:"TypeCode,omitempty"`
	Email       string    `json:"Email,omitempty"`
	BouncedAt   time.Time `json:"BouncedAt,omitempty"`
	Details     string    `json:"Details,omitempty"`

	// Open-specific fields
	FirstOpen   bool      `json:"FirstOpen,omitempty"`
	UserAgent   string    `json:"UserAgent,omitempty"`
	Platform    string    `json:"Platform,omitempty"`
	Client      string    `json:"Client,omitempty"`
	OS          string    `json:"OS,omitempty"`

	// Click-specific fields
	ClickLocation string  `json:"ClickLocation,omitempty"`
	OriginalLink  string  `json:"OriginalLink,omitempty"`
}
