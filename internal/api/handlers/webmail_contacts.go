package handlers

import (
	"net/http"
	"strconv"

	contactService "github.com/btafoya/gomailserver/internal/contact/service"
	"github.com/btafoya/gomailserver/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// WebmailContactsHandler handles webmail contact operations
type WebmailContactsHandler struct {
	contactService     *contactService.ContactService
	addressbookService *contactService.AddressbookService
	logger             *zap.Logger
}

// NewWebmailContactsHandler creates a new webmail contacts handler
func NewWebmailContactsHandler(
	contactSvc *contactService.ContactService,
	addressbookSvc *contactService.AddressbookService,
	logger *zap.Logger,
) *WebmailContactsHandler {
	return &WebmailContactsHandler{
		contactService:     contactSvc,
		addressbookService: addressbookSvc,
		logger:             logger,
	}
}

// ContactSearchResponse represents a contact search result
type ContactSearchResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Tel   string `json:"tel,omitempty"`
}

// SearchContacts handles contact autocomplete/search for composer
// GET /api/v1/webmail/contacts/search?q=john
func (h *WebmailContactsHandler) SearchContacts(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get search query
	query := r.URL.Query().Get("q")
	if query == "" {
		middleware.RespondJSON(w, http.StatusOK, []ContactSearchResponse{})
		return
	}

	// Get user's addressbooks
	addressbooks, err := h.addressbookService.GetUserAddressbooks(userID)
	if err != nil {
		h.logger.Error("failed to get addressbooks", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get addressbooks")
		return
	}

	// Search across all user's addressbooks
	var results []ContactSearchResponse
	for _, ab := range addressbooks {
		contacts, err := h.contactService.SearchContacts(ab.ID, query)
		if err != nil {
			h.logger.Error("failed to search contacts", zap.Error(err), zap.Int64("addressbook_id", ab.ID))
			continue
		}

		for _, contact := range contacts {
			results = append(results, ContactSearchResponse{
				ID:    contact.ID,
				Name:  contact.FN,
				Email: contact.Email,
				Tel:   contact.Tel,
			})
		}
	}

	middleware.RespondJSON(w, http.StatusOK, results)
}

// ListAddressbooks lists user's addressbooks
// GET /api/v1/webmail/contacts/addressbooks
func (h *WebmailContactsHandler) ListAddressbooks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	addressbooks, err := h.addressbookService.GetUserAddressbooks(userID)
	if err != nil {
		h.logger.Error("failed to get addressbooks", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get addressbooks")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, addressbooks)
}

// ListContacts lists contacts in an addressbook
// GET /api/v1/webmail/contacts/addressbooks/{id}/contacts
func (h *WebmailContactsHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		middleware.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get addressbook ID
	addressbookID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "invalid addressbook ID")
		return
	}

	// Verify addressbook belongs to user
	addressbook, err := h.addressbookService.GetAddressbook(addressbookID)
	if err != nil {
		h.logger.Error("failed to get addressbook", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get addressbook")
		return
	}
	if addressbook.UserID != userID {
		middleware.RespondError(w, http.StatusForbidden, "access denied")
		return
	}

	// Get contacts
	contacts, err := h.contactService.GetAddressbookContacts(addressbookID)
	if err != nil {
		h.logger.Error("failed to get contacts", zap.Error(err))
		middleware.RespondError(w, http.StatusInternalServerError, "failed to get contacts")
		return
	}

	// Convert to response format
	var results []ContactSearchResponse
	for _, contact := range contacts {
		results = append(results, ContactSearchResponse{
			ID:    contact.ID,
			Name:  contact.FN,
			Email: contact.Email,
			Tel:   contact.Tel,
		})
	}

	middleware.RespondJSON(w, http.StatusOK, results)
}
