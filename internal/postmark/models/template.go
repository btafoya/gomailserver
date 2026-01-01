package models

import "time"

// Template represents a PostmarkApp template
type Template struct {
	TemplateID       int       `json:"TemplateId"`
	ServerID         int       `json:"AssociatedServerId"`
	Name             string    `json:"Name"`
	Alias            string    `json:"Alias,omitempty"`
	Subject          string    `json:"Subject,omitempty"`
	HtmlBody         string    `json:"HtmlBody,omitempty"`
	TextBody         string    `json:"TextBody,omitempty"`
	TemplateType     string    `json:"TemplateType"` // Standard, Layout
	LayoutTemplate   *int      `json:"LayoutTemplate,omitempty"`
	Active           bool      `json:"Active"`
	CreatedAt        time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt        time.Time `json:"UpdatedAt,omitempty"`
}

// TemplateListResponse represents a list of templates
type TemplateListResponse struct {
	TotalCount int        `json:"TotalCount"`
	Templates  []Template `json:"Templates"`
}

// CreateTemplateRequest represents a template creation request
type CreateTemplateRequest struct {
	Name           string `json:"Name"`
	Alias          string `json:"Alias,omitempty"`
	Subject        string `json:"Subject,omitempty"`
	HtmlBody       string `json:"HtmlBody,omitempty"`
	TextBody       string `json:"TextBody,omitempty"`
	TemplateType   string `json:"TemplateType,omitempty"`
	LayoutTemplate *int   `json:"LayoutTemplate,omitempty"`
}

// UpdateTemplateRequest represents a template update request
type UpdateTemplateRequest struct {
	Name           string `json:"Name,omitempty"`
	Alias          string `json:"Alias,omitempty"`
	Subject        string `json:"Subject,omitempty"`
	HtmlBody       string `json:"HtmlBody,omitempty"`
	TextBody       string `json:"TextBody,omitempty"`
	TemplateType   string `json:"TemplateType,omitempty"`
	LayoutTemplate *int   `json:"LayoutTemplate,omitempty"`
}
