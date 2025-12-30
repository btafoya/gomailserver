package bruteforce

import (
	"time"

	"github.com/btafoya/gomailserver/internal/repository"
)

type Protection struct {
	repo      repository.FailedLoginRepository
	blacklist repository.BlacklistRepository
	threshold int           // Max failures before block
	window    time.Duration // Time window for failures
	blockTime time.Duration // How long to block
}

func NewProtection(repo repository.FailedLoginRepository, bl repository.BlacklistRepository) *Protection {
	return &Protection{
		repo:      repo,
		blacklist: bl,
		threshold: 5,
		window:    15 * time.Minute,
		blockTime: 1 * time.Hour,
	}
}

func (p *Protection) RecordFailure(ip, username string) error {
	err := p.repo.Create(ip, username)
	if err != nil {
		return err
	}

	// Check if threshold exceeded
	count, err := p.repo.CountByIP(ip, p.window)
	if err != nil {
		return err
	}

	if count >= p.threshold {
		err = p.blacklist.Add(ip, "Brute force protection", time.Now().Add(p.blockTime))
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Protection) IsBlocked(ip string) (bool, error) {
	return p.blacklist.Exists(ip)
}

func (p *Protection) ClearOnSuccess(ip string) error {
	return p.repo.DeleteByIP(ip)
}
