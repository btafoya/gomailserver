package middleware

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse represents an API success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	TotalItems int         `json:"total_items"`
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	})
}

// RespondSuccess sends a success response
func RespondSuccess(w http.ResponseWriter, data interface{}, message string) {
	RespondJSON(w, http.StatusOK, SuccessResponse{
		Data:    data,
		Message: message,
	})
}

// RespondPaginated sends a paginated response
func RespondPaginated(w http.ResponseWriter, data interface{}, page, pageSize, totalItems int) {
	totalPages := (totalItems + pageSize - 1) / pageSize
	RespondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
	})
}

// RespondCreated sends a 201 Created response
func RespondCreated(w http.ResponseWriter, data interface{}, message string) {
	RespondJSON(w, http.StatusCreated, SuccessResponse{
		Data:    data,
		Message: message,
	})
}

// RespondNoContent sends a 204 No Content response
func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
