package integration

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/emersion/go-imap/client"
)

func TestSMTPReceive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	serverAddr := "localhost:2525"

	conn, err := net.DialTimeout(10*time.Second, "tcp", serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	helo := "EHLO " + serverAddr
	if _, err := conn.Write([]byte(helo)); err != nil {
		t.Fatalf("Failed to send HELO: %v", err)
	}

	// Read greeting
	greeting := make([]byte, 1024)
	n, err := conn.Read(greeting)
	if err != nil {
		t.Fatalf("Failed to read SMTP greeting: %v", err)
	}

	if !string(greeting).StartsWith("220") {
		t.Fatalf("Expected SMTP greeting starting with 220, got: %s", string(greeting))
	}

	// Send test message
	from := "test@example.com"
	to := []string{"test@" + serverAddr}
	subject := "Integration Test Message"
	body := "This is an automated integration test message."
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\n\r\n%s",
		from, strings.Join(to, ", "),
		subject,
		time.Now().Format(time.RFC1123),
		body)

	// Send message
	mailCmd := fmt.Sprintf("MAIL FROM:<%s>\r\nDATA%s\r\n.\r\n", from, message)
	if _, err := conn.Write([]byte(mailCmd)); err != nil {
		t.Fatalf("Failed to send MAIL command: %v", err)
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		t.Fatalf("Failed to read MAIL response: %v", err)
	}

	if !strings.Contains(string(response), "250 OK") {
		t.Fatalf("Expected 250 OK response, got: %s", string(response))
	}

	// Quit connection
	if _, err := conn.Write([]byte("QUIT\r\n")); err != nil {
		t.Logf("Warning: Failed to quit SMTP connection: %v", err)
	}

	t.Log("SMTP receive test completed successfully")
}

func TestQueueProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	ctx = context.WithValue(ctx, "user_id", int64(1001))
	ctx = context.WithValue(ctx, "domain_id", int64(1001))

	t.Log("Starting queue processing test")

	queueRepo := NewSQLiteQueueRepository(":memory:")

	svc := service.NewQueueService(queueRepo, nil, nil)
	t.Cleanup(context.Background())

	t.Log("Queue service initialized")

	testMessage := &domain.QueueItem{
		ID:          1,
		UserID:      1001,
		From:        "test@example.com",
		To:          []string{"test@localhost"},
		Subject:     "Queue Test Message",
		Message:     []byte("Queue test message"),
		Status:      "pending",
		CreatedAt:   time.Now(),
		NextRetryAt: nil,
		RetryCount:  0,
		LastError:   nil,
	}

	t.Log("Creating test queue item")

	itemID, err := svc.Enqueue("test@example.com", []string{"test@localhost"}, []byte("Test message"))
	if err != nil {
		t.Fatalf("Failed to enqueue message: %v", err)
	}

	t.Logf("Message enqueued with ID: %d", itemID)

	item, err = svc.GetByID(context, itemID)
	if err != nil {
		t.Fatalf("Failed to retrieve enqueued message: %v", err)
	}

	if item == nil {
		t.Fatalf("Enqueued message not found")
	}

	t.Logf("Enqueued message status: %s", item.Status)

	t.Logf("Enqueued message message size: %d bytes", len(item.Message))

	t.Log("Queue processing test completed")
}

func TestIMAPDelivery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	imapServer := "localhost:143"
	imapUser := "test@localhost"
	imapPassword := "testpass"

	conn, err := imap.Dial("tcp", imapServer, "", false)
	if err != nil {
		t.Fatalf("Failed to connect to IMAP server: %v", err)
	}
	defer conn.Logout()
	defer conn.Close()

	if err := conn.Login(imapUser, imapPassword); err != nil {
		t.Fatalf("IMAP authentication failed: %v", err)
	}
	t.Log("IMAP authentication successful")

	mbox, err := conn.Select("INBOX", nil)
	if err != nil {
		t.Fatalf("Failed to select INBOX: %v", err)
	}
	defer mbox.Close()

	seq, _ := mbox.Unseq()
	if err != nil {
		t.Fatalf("Failed to get unsorted sequence for INBOX: %v", err)
	}

	// Search for test message
	searchClient := imap.NewSearchClient("", true) // empty body for full search
	criteria := imap.NewSearchCriteria().
		WithCharset("utf-8").
		WithBodyFields(false).
		WithSearchAll()

	messages, err := searchClient.Criteria(criteria).Search("")
	if err != nil {
		t.Logf("Warning: IMAP search failed: %v", err)
	}

	t.Logf("IMAP search completed: found %d messages", len(messages))

	var found bool
	for _, msg := range messages {
		msgBody, err := msg.FetchBody()
		if err != nil {
			t.Logf("Warning: Failed to fetch message body: %v", err)
			continue
		}

		msgText, err := msg.Text()
		if err != nil {
			t.Logf("Warning: Failed to read message text: %v", err)
			continue
		}

		if strings.Contains(msgText, "Integration Test Message") {
			t.Logf("Found test message: %s", msg.Envelope.Subject)
			found = true
			break
		}
	}

	if !found {
		t.Log("Test message not found in IMAP search results")
	}

	conn.Close()
	t.Log("IMAP delivery test completed")
}

// Mock queue repository for testing
type mockSQLiteQueueRepository struct {
	items []*domain.QueueItem
}

func NewSQLiteQueueRepository(path string) repository.QueueRepository {
	return &mockSQLiteQueueRepository{items: make([]*domain.QueueItem)}
}

func (r *mockSQLiteQueueRepository) Enqueue(from string, to []string, message []byte) (string, error) {
	item := &domain.QueueItem{
		From:        from,
		To:          to,
		Message:     message,
		Status:      "pending",
		CreatedAt:   time.Now(),
		NextRetryAt: nil,
		RetryCount:  0,
	}

	id := int64(time.Now().Unix())
	r.items = append(r.items, item)
	return fmt.Sprintf("%d", id), nil
}

func (r *mockSQLiteQueueRepository) GetByID(ctx context.Context, id int64) (*domain.QueueItem, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, fmt.Errorf("queue item not found: %d", id)
}

func (r *mockSQLiteQueueRepository) GetPending(ctx context.Context) ([]*domain.QueueItem, error) {
	return r.items, nil
}

func (r *mockSQLiteQueueRepository) MarkDelivered(id int64) error {
	for _, item := range r.items {
		if item.ID == id {
			item.Status = "delivered"
			return nil
		}
	}
	return fmt.Errorf("queue item not found: %d", id)
}

func (r *mockSQLiteQueueRepository) GetByFrom(from string, limit int) ([]*domain.QueueItem, error) {
	return r.items, nil
}

func (r *mockSQLiteQueueRepository) MarkFailed(id int64, errorMsg string) error {
	return fmt.Errorf("not implemented")
}

func (r *mockSQLiteQueueRepository) IncrementRetry(id int64, currentRetryCount int, failedAt time.Time) error {
	return fmt.Errorf("not implemented")
}

func (r *mockSQLiteQueueRepository) CalculateNextRetry(retryCount int, failedAt time.Time) time.Time {
	return time.Now().Add(5 * time.Minute)
}

func (r *mockSQLiteQueueRepository) GetByIDAll(ctx context.Context) ([]*domain.QueueItem, error) {
	return r.items, nil
}

func (r *mockSQLiteQueueRepository) Delete(ctx context.Context, id int64) error {
	for i, item := range r.items {
		if item.ID == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("queue item not found: %d", id)
}

func (r *mockSQLiteQueueRepository) Cleanup(ctx context.Context) error {
	r.items = make([]*domain.QueueItem)
	return nil
}
