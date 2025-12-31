package bruteforce

import (
	"time"

	"github.com/btafoya/gomailserver/internal/repository"
)

type Protection struct {
	loginAttemptRepo repository.LoginAttemptRepository
	blacklistRepo    repository.IPBlacklistRepository
	threshold        int           // Max failures before block
	window           time.Duration // Time window for failures
	blockTime        time.Duration // How long to block
}

func NewProtection(loginRepo repository.LoginAttemptRepository, blacklistRepo repository.IPBlacklistRepository) *Protection {
	return &Protection{
		loginAttemptRepo: loginRepo,
		blacklistRepo:    blacklistRepo,
		threshold:        5,
		window:           15 * time.Minute,
		blockTime:        1 * time.Hour,
	}
}

func (p *Protection) RecordFailure(ip, email string) error {
	err := p.loginAttemptRepo.Record(ip, email, false)
	if err != nil {
		return err
	}

	// Check if threshold exceeded
	count, err := p.loginAttemptRepo.GetRecentFailures(ip, p.window)
	if err != nil {
		return err
	}

	if count >= p.threshold {
		expiresAt := time.Now().Add(p.blockTime)
		err = p.blacklistRepo.Add(ip, "Brute force protection", &expiresAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Protection) RecordSuccess(ip, email string) error {
	return p.loginAttemptRepo.Record(ip, email, true)
}

func (p *Protection) IsBlocked(ip string) (bool, error) {
	return p.blacklistRepo.IsBlacklisted(ip)
}

func (p *Protection) Cleanup() error {
	// Clean up old login attempts and expired blacklist entries
	if err := p.loginAttemptRepo.Cleanup(7 * 24 * time.Hour); err != nil {
		return err
	}
	return p.blacklistRepo.RemoveExpired()
}
