package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository/sqlite"
	"github.com/btafoya/gomailserver/internal/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var createAdminCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "Create the first admin user",
	Long:  "Create the first administrative user for the mail server. This should be run before starting the server for the first time.",
	RunE:  createAdmin,
}

func init() {
	rootCmd.AddCommand(createAdminCmd)
}

func createAdmin(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	logger, err := config.NewLogger(cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	// Initialize database
	dbConfig := database.Config{
		Path:       cfg.Database.Path,
		WALEnabled: cfg.Database.WALEnabled,
	}

	db, err := database.New(dbConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	defer db.Close()

	// Run migrations to ensure database is up to date
	if err := db.Migrate(); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	// Create repositories
	userRepo := sqlite.NewUserRepository(db)
	domainRepo := sqlite.NewDomainRepository(db)

	// Create services
	userSvc := service.NewUserService(userRepo, domainRepo, logger)
	domainSvc := service.NewDomainService(domainRepo)

	// Initialize default domain template
	if err := domainSvc.EnsureDefaultTemplate(); err != nil {
		return fmt.Errorf("failed to initialize default domain template: %w", err)
	}

	// Prompt for admin details
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Create Admin User")
	fmt.Println("================")
	fmt.Println()

	// Get email
	fmt.Print("Email address: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}
	email = strings.TrimSpace(email)

	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Extract domain from email
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}
	domainName := parts[1]

	// Get full name
	fmt.Print("Full name: ")
	fullName, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read full name: %w", err)
	}
	fullName = strings.TrimSpace(fullName)

	// Get password (hidden)
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()
	password := string(passwordBytes)

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Confirm password
	fmt.Print("Confirm password: ")
	confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password confirmation: %w", err)
	}
	fmt.Println()
	confirm := string(confirmBytes)

	if password != confirm {
		return fmt.Errorf("passwords do not match")
	}

	// Check if domain exists
	existingDomain, err := domainRepo.GetByName(domainName)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to check domain: %w", err)
	}

	var domainID int64
	if existingDomain != nil {
		// Domain exists
		domainID = existingDomain.ID
		logger.Info("using existing domain", zap.String("domain", domainName))
	} else {
		// Create domain from template to ensure all security settings are initialized
		newDomain, err := domainSvc.CreateDomainFromTemplate(domainName)
		if err != nil {
			return fmt.Errorf("failed to create domain: %w", err)
		}
		domainID = newDomain.ID
		logger.Info("created new domain", zap.String("domain", domainName))
	}

	// Check if user already exists
	existingUser, err := userRepo.GetByEmail(email)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to check user: %w", err)
	}

	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	user := &domain.User{
		Email:        email,
		DomainID:     domainID,
		PasswordHash: string(passwordHash),
		FullName:     fullName,
		DisplayName:  fullName,
		Role:         "admin",
		Quota:        1073741824, // 1GB default
		UsedQuota:    0,
		Status:       "active",
		AuthMethod:   "password",
		TOTPEnabled:  false,
		Language:     "en",
	}

	if err := userSvc.Create(user, password); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	fmt.Println()
	fmt.Println("✓ Admin user created successfully!")
	fmt.Printf("  Email: %s\n", email)
	fmt.Printf("  Name: %s\n", fullName)
	fmt.Printf("  Role: admin\n")
	fmt.Printf("  Domain: %s\n", domainName)
	fmt.Println()
	fmt.Println("You can now start the server and log in with these credentials.")

	return nil
}
