package ratelimit

import (
	"time"

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
		// Or return an error
		return true, nil
	}

	count, err := l.repo.GetCount(limitType, key, limit.Window)
	if err != nil {
		return true, err // Fail open
	}

	if count >= limit.Count {
		return false, nil // Rate limited
	}

	err = l.repo.Increment(limitType, key)
	if err != nil {
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
