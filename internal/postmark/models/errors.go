package models

import (
	"encoding/json"
	"net/http"
)

// PostmarkError represents an error in PostmarkApp format
type PostmarkError struct {
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

// Error implements the error interface
func (e *PostmarkError) Error() string {
	return e.Message
}

// WriteError writes a PostmarkApp-formatted error response
func WriteError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")

	var status int
	switch code {
	case 0:
		status = http.StatusOK
	case 401:
		status = http.StatusUnauthorized
	case 300, 405, 406, 409:
		status = http.StatusUnprocessableEntity
	case 402, 410, 411:
		status = http.StatusBadRequest
	case 429:
		status = http.StatusTooManyRequests
	case 500:
		status = http.StatusInternalServerError
	case 503:
		status = http.StatusServiceUnavailable
	default:
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(PostmarkError{
		ErrorCode: code,
		Message:   message,
	})
}

// Common error codes
const (
	ErrorCodeSuccess                 = 0
	ErrorCodeInvalidEmail            = 300
	ErrorCodeSenderNotFound          = 400
	ErrorCodeUnauthorized            = 401
	ErrorCodeInactiveRecipient       = 402
	ErrorCodeInvalidJSON             = 405
	ErrorCodeInactiveServer          = 406
	ErrorCodeJSONRequired            = 409
	ErrorCodeBatchLimitExceeded      = 410
	ErrorCodeForbiddenAttachment     = 411
	ErrorCodeRateLimitExceeded       = 429
	ErrorCodeInternalServerError     = 500
	ErrorCodeServiceUnavailable      = 503
)

// Error messages
const (
	MsgSuccess                  = "OK"
	MsgInvalidEmail             = "Invalid email request"
	MsgSenderNotFound           = "Sender signature not found"
	MsgUnauthorized             = "Unauthorized: please provide valid API token"
	MsgInactiveRecipient        = "Inactive recipient"
	MsgInvalidJSON              = "Invalid JSON"
	MsgInactiveServer           = "Inactive server"
	MsgJSONRequired             = "JSON content type required"
	MsgBatchLimitExceeded       = "Too many messages in batch (max 500)"
	MsgForbiddenAttachment      = "Forbidden attachment type"
	MsgRateLimitExceeded        = "Rate limit exceeded"
	MsgInternalServerError      = "Internal server error"
	MsgServiceUnavailable       = "Service unavailable"
)
