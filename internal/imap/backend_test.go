package imap

import (
	"errors"
	"net"
	"testing"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
)

// mockUserService for IMAP backend tests
type mockUserService struct {
	authenticateFunc func(string, string) (*domain.User, error)
}

func (m *mockUserService) Create(user *domain.User, password string) error {
	return nil
}

func (m *mockUserService) Authenticate(email, password string) (*domain.User, error) {
	if m.authenticateFunc != nil {
		return m.authenticateFunc(email, password)
	}
	return nil, errors.New("authentication failed")
}

func (m *mockUserService) GetByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserService) GetByID(id int64) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserService) Update(user *domain.User) error {
	return nil
}

func (m *mockUserService) UpdatePassword(userID int64, newPassword string) error {
	return nil
}

func (m *mockUserService) Delete(id int64) error {
	return nil
}

// mockMailboxService for IMAP backend tests
type mockMailboxService struct {
	listFunc       func(int64, bool) ([]*domain.Mailbox, error)
	getByNameFunc  func(int64, string) (*domain.Mailbox, error)
	createFunc     func(int64, string, string) error
	deleteFunc     func(int64) error
	renameFunc     func(int64, string) error
}

func (m *mockMailboxService) List(userID int64, subscribedOnly bool) ([]*domain.Mailbox, error) {
	if m.listFunc != nil {
		return m.listFunc(userID, subscribedOnly)
	}
	return []*domain.Mailbox{}, nil
}

func (m *mockMailboxService) GetByName(userID int64, name string) (*domain.Mailbox, error) {
	if m.getByNameFunc != nil {
		return m.getByNameFunc(userID, name)
	}
	return nil, errors.New("not found")
}

func (m *mockMailboxService) Create(userID int64, name, specialUse string) error {
	if m.createFunc != nil {
		return m.createFunc(userID, name, specialUse)
	}
	return nil
}

func (m *mockMailboxService) Delete(mailboxID int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(mailboxID)
	}
	return nil
}

func (m *mockMailboxService) Rename(mailboxID int64, newName string) error {
	if m.renameFunc != nil {
		return m.renameFunc(mailboxID, newName)
	}
	return nil
}

func (m *mockMailboxService) UpdateSubscription(id int64, subscribed bool) error {
	return nil
}

// mockMessageService for IMAP backend tests
type mockMessageService struct{}

func (m *mockMessageService) Store(userID, mailboxID, uid int64, messageData []byte) (*domain.Message, error) {
	return &domain.Message{ID: 1}, nil
}

func (m *mockMessageService) GetByID(id int64) (*domain.Message, error) {
	return nil, nil
}

func (m *mockMessageService) GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error) {
	return nil, nil
}

func (m *mockMessageService) Delete(id int64) error {
	return nil
}

func TestBackend_Login(t *testing.T) {
	logger := zap.NewNop()

	t.Run("authenticates with valid credentials", func(t *testing.T) {
		userSvc := &mockUserService{
			authenticateFunc: func(email, password string) (*domain.User, error) {
				if email == "test@example.com" && password == "correctpassword" {
					return &domain.User{
						ID:     1,
						Email:  email,
						Status: "active",
					}, nil
				}
				return nil, errors.New("authentication failed")
			},
		}

		b := &Backend{
			userService:    userSvc,
			mailboxService: &mockMailboxService{},
			messageService: &mockMessageService{},
			logger:         logger,
		}

		// Create mock connection info
		connInfo := &imap.ConnInfo{
			RemoteAddr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345},
		}

		user, err := b.Login(connInfo, "test@example.com", "correctpassword")
		if err != nil {
			t.Fatalf("expected successful login, got error: %v", err)
		}

		if user == nil {
			t.Fatal("expected user to be returned")
		}

		// Verify user type
		imapUser, ok := user.(*User)
		if !ok {
			t.Fatal("expected user to be of type *User")
		}

		if imapUser.user.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", imapUser.user.Email)
		}
	})

	t.Run("rejects invalid credentials", func(t *testing.T) {
		userSvc := &mockUserService{
			authenticateFunc: func(email, password string) (*domain.User, error) {
				return nil, errors.New("authentication failed")
			},
		}

		b := &Backend{
			userService:    userSvc,
			mailboxService: &mockMailboxService{},
			messageService: &mockMessageService{},
			logger:         logger,
		}

		connInfo := &imap.ConnInfo{
			RemoteAddr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345},
		}

		user, err := b.Login(connInfo, "test@example.com", "wrongpassword")
		if err == nil {
			t.Error("expected login to fail with wrong password")
		}

		if user != nil {
			t.Error("expected no user to be returned on failed login")
		}

		// Verify error is ErrInvalidCredentials
		if err != backend.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("rejects disabled user", func(t *testing.T) {
		userSvc := &mockUserService{
			authenticateFunc: func(email, password string) (*domain.User, error) {
				return &domain.User{
					ID:     1,
					Email:  email,
					Status: "disabled",
				}, nil
			},
		}

		b := &Backend{
			userService:    userSvc,
			mailboxService: &mockMailboxService{},
			messageService: &mockMessageService{},
			logger:         logger,
		}

		connInfo := &imap.ConnInfo{
			RemoteAddr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345},
		}

		user, err := b.Login(connInfo, "test@example.com", "password")
		if err == nil {
			t.Error("expected login to fail for disabled user")
		}

		if user != nil {
			t.Error("expected no user to be returned for disabled user")
		}

		if err != backend.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}

func TestUser_Username(t *testing.T) {
	logger := zap.NewNop()

	domainUser := &domain.User{
		ID:    1,
		Email: "test@example.com",
	}

	user := &User{
		user:   domainUser,
		logger: logger,
	}

	if user.Username() != "test@example.com" {
		t.Errorf("expected username 'test@example.com', got '%s'", user.Username())
	}
}

func TestUser_ListMailboxes(t *testing.T) {
	logger := zap.NewNop()

	t.Run("lists all mailboxes", func(t *testing.T) {
		mailboxes := []*domain.Mailbox{
			{ID: 1, Name: "INBOX", SpecialUse: ""},
			{ID: 2, Name: "Sent", SpecialUse: "\\Sent"},
			{ID: 3, Name: "Drafts", SpecialUse: "\\Drafts"},
		}

		mailboxSvc := &mockMailboxService{
			listFunc: func(userID int64, subscribedOnly bool) ([]*domain.Mailbox, error) {
				if subscribedOnly {
					return []*domain.Mailbox{mailboxes[0]}, nil
				}
				return mailboxes, nil
			},
		}

		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: mailboxSvc,
			logger:         logger,
		}

		result, err := user.ListMailboxes(false)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result) != 3 {
			t.Errorf("expected 3 mailboxes, got %d", len(result))
		}
	})

	t.Run("lists only subscribed mailboxes", func(t *testing.T) {
		mailboxSvc := &mockMailboxService{
			listFunc: func(userID int64, subscribedOnly bool) ([]*domain.Mailbox, error) {
				if subscribedOnly {
					return []*domain.Mailbox{
						{ID: 1, Name: "INBOX", Subscribed: true},
					}, nil
				}
				return nil, nil
			},
		}

		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: mailboxSvc,
			logger:         logger,
		}

		result, err := user.ListMailboxes(true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result) != 1 {
			t.Errorf("expected 1 mailbox, got %d", len(result))
		}
	})
}

func TestUser_GetMailbox(t *testing.T) {
	logger := zap.NewNop()

	t.Run("gets mailbox by name", func(t *testing.T) {
		expectedMailbox := &domain.Mailbox{
			ID:   1,
			Name: "INBOX",
		}

		mailboxSvc := &mockMailboxService{
			getByNameFunc: func(userID int64, name string) (*domain.Mailbox, error) {
				if name == "INBOX" {
					return expectedMailbox, nil
				}
				return nil, errors.New("not found")
			},
		}

		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: mailboxSvc,
			logger:         logger,
		}

		mailbox, err := user.GetMailbox("INBOX")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if mailbox == nil {
			t.Fatal("expected mailbox to be returned")
		}

		// Verify it's wrapped in our Mailbox type
		_, ok := mailbox.(*Mailbox)
		if !ok {
			t.Error("expected mailbox to be of type *Mailbox")
		}
	})

	t.Run("returns error for non-existent mailbox", func(t *testing.T) {
		mailboxSvc := &mockMailboxService{
			getByNameFunc: func(userID int64, name string) (*domain.Mailbox, error) {
				return nil, errors.New("not found")
			},
		}

		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: mailboxSvc,
			logger:         logger,
		}

		mailbox, err := user.GetMailbox("NonExistent")
		if err == nil {
			t.Error("expected error for non-existent mailbox")
		}

		if mailbox != nil {
			t.Error("expected no mailbox to be returned")
		}

		if err != backend.ErrNoSuchMailbox {
			t.Errorf("expected ErrNoSuchMailbox, got %v", err)
		}
	})
}

func TestUser_CreateMailbox(t *testing.T) {
	logger := zap.NewNop()

	t.Run("creates mailbox", func(t *testing.T) {
		created := false
		mailboxSvc := &mockMailboxService{
			createFunc: func(userID int64, name, specialUse string) error {
				created = true
				if name != "NewFolder" {
					t.Errorf("expected name 'NewFolder', got '%s'", name)
				}
				return nil
			},
		}

		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: mailboxSvc,
			logger:         logger,
		}

		err := user.CreateMailbox("NewFolder")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !created {
			t.Error("expected mailbox service to be called")
		}
	})
}

func TestUser_DeleteMailbox(t *testing.T) {
	logger := zap.NewNop()

	t.Run("deletes mailbox", func(t *testing.T) {
		mailboxSvc := &mockMailboxService{
			getByNameFunc: func(userID int64, name string) (*domain.Mailbox, error) {
				return &domain.Mailbox{ID: 2, Name: "OldFolder"}, nil
			},
			deleteFunc: func(mailboxID int64) error {
				if mailboxID != 2 {
					t.Errorf("expected mailbox ID 2, got %d", mailboxID)
				}
				return nil
			},
		}

		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: mailboxSvc,
			logger:         logger,
		}

		err := user.DeleteMailbox("OldFolder")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("prevents deletion of INBOX", func(t *testing.T) {
		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: &mockMailboxService{},
			logger:         logger,
		}

		err := user.DeleteMailbox("INBOX")
		if err == nil {
			t.Error("expected error when trying to delete INBOX")
		}
	})
}

func TestUser_RenameMailbox(t *testing.T) {
	logger := zap.NewNop()

	t.Run("renames mailbox", func(t *testing.T) {
		mailboxSvc := &mockMailboxService{
			getByNameFunc: func(userID int64, name string) (*domain.Mailbox, error) {
				return &domain.Mailbox{ID: 2, Name: "OldName"}, nil
			},
			renameFunc: func(mailboxID int64, newName string) error {
				if newName != "NewName" {
					t.Errorf("expected new name 'NewName', got '%s'", newName)
				}
				return nil
			},
		}

		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: mailboxSvc,
			logger:         logger,
		}

		err := user.RenameMailbox("OldName", "NewName")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("prevents renaming INBOX", func(t *testing.T) {
		user := &User{
			user:           &domain.User{ID: 1},
			mailboxService: &mockMailboxService{},
			logger:         logger,
		}

		err := user.RenameMailbox("INBOX", "NewName")
		if err == nil {
			t.Error("expected error when trying to rename INBOX")
		}
	})
}

func TestUser_Logout(t *testing.T) {
	logger := zap.NewNop()

	user := &User{
		user:   &domain.User{ID: 1, Email: "test@example.com"},
		logger: logger,
	}

	err := user.Logout()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
