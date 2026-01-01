package models

import "time"

// Attachment represents an email attachment
type Attachment struct {
	Name        string `json:"Name"`
	Content     string `json:"Content"` // Base64-encoded
	ContentType string `json:"ContentType"`
	ContentID   string `json:"ContentID,omitempty"`
}

// Header represents a custom email header
type Header struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

// EmailRequest represents a single email send request
type EmailRequest struct {
	From          string            `json:"From"`
	To            string            `json:"To"`
	Cc            string            `json:"Cc,omitempty"`
	Bcc           string            `json:"Bcc,omitempty"`
	Subject       string            `json:"Subject"`
	Tag           string            `json:"Tag,omitempty"`
	HtmlBody      string            `json:"HtmlBody,omitempty"`
	TextBody      string            `json:"TextBody,omitempty"`
	ReplyTo       string            `json:"ReplyTo,omitempty"`
	Headers       []Header          `json:"Headers,omitempty"`
	TrackOpens    bool              `json:"TrackOpens,omitempty"`
	TrackLinks    string            `json:"TrackLinks,omitempty"` // None, HtmlOnly, HtmlAndText, TextOnly
	Attachments   []Attachment      `json:"Attachments,omitempty"`
	Metadata      map[string]string `json:"Metadata,omitempty"`
	MessageStream string            `json:"MessageStream,omitempty"`
}

// EmailResponse represents a successful email send response
type EmailResponse struct {
	To          string    `json:"To"`
	SubmittedAt time.Time `json:"SubmittedAt"`
	MessageID   string    `json:"MessageID"`
	ErrorCode   int       `json:"ErrorCode"`
	Message     string    `json:"Message"`
}

// BatchEmailRequest represents a batch email send request
type BatchEmailRequest []EmailRequest

// BatchEmailResponse represents a batch email send response
type BatchEmailResponse []EmailResponse

// TemplateEmailRequest represents an email send with template
type TemplateEmailRequest struct {
	From              string            `json:"From"`
	To                string            `json:"To"`
	Cc                string            `json:"Cc,omitempty"`
	Bcc               string            `json:"Bcc,omitempty"`
	Tag               string            `json:"Tag,omitempty"`
	ReplyTo           string            `json:"ReplyTo,omitempty"`
	Headers           []Header          `json:"Headers,omitempty"`
	TrackOpens        bool              `json:"TrackOpens,omitempty"`
	TrackLinks        string            `json:"TrackLinks,omitempty"`
	TemplateID        int               `json:"TemplateId,omitempty"`
	TemplateAlias     string            `json:"TemplateAlias,omitempty"`
	TemplateModel     map[string]interface{} `json:"TemplateModel"`
	InlineCss         bool              `json:"InlineCss,omitempty"`
	Attachments       []Attachment      `json:"Attachments,omitempty"`
	Metadata          map[string]string `json:"Metadata,omitempty"`
	MessageStream     string            `json:"MessageStream,omitempty"`
}
