package ratelimit

import (
	"time"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
)

type Limiter struct {
	repo repository.RateLimitRepository
}

type Limit struct {
	Count  int           // Max requests
	Window time.Duration // Time window
}

var DefaultLimits = map[string]Limit{
	"smtp_per_ip":     {Count: 100, Window: time.Hour},
	"smtp_per_user":   {Count: 500, Window: time.Hour},
	"smtp_per_domain": {Count: 1000, Window: time.Hour},
	"auth_per_ip":     {Count: 10, Window: 15 * time.Minute},
	"imap_per_user":   {Count: 1000, Window: time.Hour},
}

func NewLimiter(repo repository.RateLimitRepository) *Limiter {
	return &Limiter{repo: repo}
}

func (l *Limiter) Check(limitType, key string) (bool, error) {
	limit, ok := DefaultLimits[limitType]
	if !ok {
		// Unknown limit type, fail open
		return true, nil
	}

	entry, err := l.repo.Get(key, limitType)
	if err != nil || entry == nil {
		// No existing entry, create new one
		now := time.Now()
		entry = &domain.RateLimitEntry{
			Key:         key,
			Type:        limitType,
			Count:       1,
			WindowStart: now,
		}
		if err := l.repo.CreateOrUpdate(entry); err != nil {
			return true, err // Fail open
		}
		return true, nil
	}

	// Check if window has expired
	if time.Since(entry.WindowStart) >= limit.Window {
		// Reset window
		entry.Count = 1
		entry.WindowStart = time.Now()
		if err := l.repo.CreateOrUpdate(entry); err != nil {
			return true, err // Fail open
		}
		return true, nil
	}

	// Check if limit exceeded
	if entry.Count >= limit.Count {
		return false, nil // Rate limited
	}

	// Increment count
	entry.Count++
	if err := l.repo.CreateOrUpdate(entry); err != nil {
		return true, err // Fail open
	}

	return true, nil
}

func (l *Limiter) CheckIP(ip string) (bool, error) {
	return l.Check("smtp_per_ip", ip)
}

func (l *Limiter) CheckUser(userID string) (bool, error) {
	return l.Check("smtp_per_user", userID)
}

func (l *Limiter) CheckAuth(ip string) (bool, error) {
	return l.Check("auth_per_ip", ip)
}
