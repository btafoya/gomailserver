package models

import "time"

// Server represents a PostmarkApp server (API token group)
type Server struct {
	ID            int       `json:"ID,omitempty"`
	Name          string    `json:"Name"`
	APIToken      string    `json:"ApiToken,omitempty"` // Not returned in responses for security
	Color         string    `json:"Color,omitempty"`
	SmtpApiActivated bool   `json:"SmtpApiActivated,omitempty"`
	RawEmailEnabled  bool   `json:"RawEmailEnabled,omitempty"`
	ServerLink    string    `json:"ServerLink,omitempty"`
	InboundAddress string   `json:"InboundAddress,omitempty"`
	InboundHookUrl string   `json:"InboundHookUrl,omitempty"`
	BounceHookUrl  string   `json:"BounceHookUrl,omitempty"`
	OpenHookUrl    string   `json:"OpenHookUrl,omitempty"`
	PostFirstOpenOnly bool  `json:"PostFirstOpenOnly,omitempty"`
	TrackOpens     bool     `json:"TrackOpens"`
	TrackLinks     string   `json:"TrackLinks"` // None, HtmlOnly, HtmlAndText, TextOnly
	IncludeBounceContentInHook bool `json:"IncludeBounceContentInHook,omitempty"`
	EnableSmtpApiErrorHooks   bool  `json:"EnableSmtpApiErrorHooks,omitempty"`
	DeliveryHookUrl string   `json:"DeliveryHookUrl,omitempty"`
	CreatedAt      time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt      time.Time `json:"UpdatedAt,omitempty"`
}

// UpdateServerRequest represents a server update request
type UpdateServerRequest struct {
	Name                  string `json:"Name,omitempty"`
	Color                 string `json:"Color,omitempty"`
	RawEmailEnabled       bool   `json:"RawEmailEnabled,omitempty"`
	SmtpApiActivated      bool   `json:"SmtpApiActivated,omitempty"`
	InboundHookUrl        string `json:"InboundHookUrl,omitempty"`
	BounceHookUrl         string `json:"BounceHookUrl,omitempty"`
	OpenHookUrl           string `json:"OpenHookUrl,omitempty"`
	PostFirstOpenOnly     bool   `json:"PostFirstOpenOnly,omitempty"`
	TrackOpens            bool   `json:"TrackOpens,omitempty"`
	TrackLinks            string `json:"TrackLinks,omitempty"`
	IncludeBounceContentInHook bool `json:"IncludeBounceContentInHook,omitempty"`
	EnableSmtpApiErrorHooks    bool `json:"EnableSmtpApiErrorHooks,omitempty"`
	DeliveryHookUrl       string `json:"DeliveryHookUrl,omitempty"`
}
