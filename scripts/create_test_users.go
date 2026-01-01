package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// TestUser represents a test user to create
type TestUser struct {
	Email    string
	Password string
	FullName string
	Domain   string
}

var defaultUsers = []TestUser{
	{
		Email:    "test@localhost",
		Password: "testpass123",
		FullName: "Test User",
		Domain:   "localhost",
	},
	{
		Email:    "alice@localhost",
		Password: "alice123",
		FullName: "Alice Smith",
		Domain:   "localhost",
	},
	{
		Email:    "bob@localhost",
		Password: "bob123",
		FullName: "Bob Jones",
		Domain:   "localhost",
	},
}

func main() {
	dbPath := flag.String("db", "./mailserver.db", "Path to SQLite database")
	flag.Parse()

	// Open database
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Ensure default domain exists
	if err := ensureDomain(ctx, db, "localhost"); err != nil {
		log.Fatalf("Failed to ensure domain: %v", err)
	}

	// Create test users
	for _, user := range defaultUsers {
		if err := createUser(ctx, db, user); err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
			continue
		}
		fmt.Printf("✅ Created user: %s (password: %s)\n", user.Email, user.Password)
	}

	// Create default mailboxes for each user
	for _, user := range defaultUsers {
		if err := createDefaultMailboxes(ctx, db, user.Email); err != nil {
			log.Printf("Failed to create mailboxes for %s: %v", user.Email, err)
			continue
		}
		fmt.Printf("✅ Created default mailboxes for: %s\n", user.Email)
	}

	fmt.Println("\n🎉 Test users created successfully!")
	fmt.Println("\nYou can now test SMTP authentication with:")
	for _, user := range defaultUsers {
		fmt.Printf("  Email: %s, Password: %s\n", user.Email, user.Password)
	}
}

func ensureDomain(ctx context.Context, db *sql.DB, domain string) error {
	// Check if domain exists
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM domains WHERE name = ?)", domain).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check domain: %w", err)
	}

	if exists {
		fmt.Printf("Domain %s already exists\n", domain)
		return nil
	}

	// Create domain
	now := time.Now()
	_, err = db.ExecContext(ctx, `
		INSERT INTO domains (
			name, status, max_users, max_mailbox_size, default_quota,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, domain, "active", 100, int64(10240*1024*1024), int64(1024*1024*1024), now, now)

	if err != nil {
		return fmt.Errorf("create domain: %w", err)
	}

	fmt.Printf("✅ Created domain: %s\n", domain)
	return nil
}

func createUser(ctx context.Context, db *sql.DB, user TestUser) error {
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	// Get domain ID
	var domainID int64
	err = db.QueryRowContext(ctx, "SELECT id FROM domains WHERE name = ?", user.Domain).Scan(&domainID)
	if err != nil {
		return fmt.Errorf("get domain id: %w", err)
	}

	// Check if user exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check user: %w", err)
	}

	if exists {
		// Update password
		_, err = db.ExecContext(ctx, `
			UPDATE users
			SET password_hash = ?, full_name = ?, updated_at = ?
			WHERE email = ?
		`, string(hash), user.FullName, time.Now(), user.Email)
		if err != nil {
			return fmt.Errorf("update user: %w", err)
		}
		fmt.Printf("Updated existing user: %s\n", user.Email)
		return nil
	}

	// Create user
	now := time.Now()
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (
			email, domain_id, password_hash, full_name, display_name,
			quota, used_quota, status, created_at, updated_at, last_login
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, user.Email, domainID, string(hash), user.FullName, user.FullName,
		int64(1024*1024*1024), int64(0), "active", now, now, nil)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func createDefaultMailboxes(ctx context.Context, db *sql.DB, email string) error {
	mailboxes := []struct {
		Name        string
		SpecialUse  string
		Subscribed  bool
	}{
		{"INBOX", "", true},
		{"Drafts", "\\Drafts", true},
		{"Sent", "\\Sent", true},
		{"Trash", "\\Trash", true},
		{"Junk", "\\Junk", true},
		{"Archive", "\\Archive", false},
	}

	for _, mb := range mailboxes {
		// Check if mailbox exists
		var exists bool
		err := db.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM mailboxes WHERE user_email = ? AND name = ?)",
			email, mb.Name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check mailbox %s: %w", mb.Name, err)
		}

		if exists {
			continue
		}

		// Create mailbox
		now := time.Now()
		uidvalidity := uint32(now.Unix())
		_, err = db.ExecContext(ctx, `
			INSERT INTO mailboxes (
				user_email, name, parent_id, subscribed, special_use,
				uidvalidity, uidnext, message_count, recent_count,
				unseen_count, created_at, updated_at
			) VALUES (?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, email, mb.Name, mb.Subscribed, mb.SpecialUse, uidvalidity, 1, 0, 0, 0, now, now)

		if err != nil {
			return fmt.Errorf("create mailbox %s: %w", mb.Name, err)
		}
	}

	return nil
}
