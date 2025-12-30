package greylist

import (
	"time"

	"github.com/btafoya/gomailserver/internal/repository"
)

type Greylister struct {
	repo   repository.GreylistRepository
	delay  time.Duration // Minimum delay before accepting
	expiry time.Duration // How long to remember triplets
}

func NewGreylister(repo repository.GreylistRepository) *Greylister {
	return &Greylister{
		repo:   repo,
		delay:  5 * time.Minute,
		expiry: 30 * 24 * time.Hour, // 30 days
	}
}

type CheckResult struct {
	Action   string // "accept", "defer", "pass"
	Message  string
	WaitTime time.Duration
}

func (g *Greylister) Check(ip, sender, recipient string) (*CheckResult, error) {
	triplet, err := g.repo.Get(ip, sender, recipient)
	if err != nil {
		// First time seeing this triplet
		_, err := g.repo.Create(ip, sender, recipient)
		if err != nil {
			return nil, err
		}
		return &CheckResult{
			Action:   "defer",
			Message:  "Greylisting in effect, please retry later",
			WaitTime: g.delay,
		}, nil
	}

	// Check if enough time has passed
	if time.Since(triplet.FirstSeen) < g.delay {
		remaining := g.delay - time.Since(triplet.FirstSeen)
		return &CheckResult{
			Action:   "defer",
			Message:  "Greylisting in effect, please retry later",
			WaitTime: remaining,
		}, nil
	}

	// Passed greylisting
	err = g.repo.IncrementPass(triplet.ID)
	if err != nil {
		return nil, err
	}

	return &CheckResult{
		Action:  "accept",
		Message: "Greylisting passed",
	}, nil
}

func (g *Greylister) Cleanup() error {
	return g.repo.DeleteOlderThan(g.expiry)
}
