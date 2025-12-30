package antispam

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/teamwork/spamc"
)

type SpamAssassin struct {
	client *spamc.Client
}

func NewSpamAssassin(host string, port int) *SpamAssassin {
	return &SpamAssassin{
		client: spamc.New(fmt.Sprintf("%s:%d", host, port), 10*time.Second),
	}
}

func (s *SpamAssassin) Check(message []byte) (*SpamResult, error) {
	reply, err := s.client.Check(context.Background(), bytes.NewReader(message), nil)
	if err != nil {
		return nil, err
	}

	return &SpamResult{
		Score:     reply.Score,
		Threshold: reply.BaseScore,
		IsSpam:    reply.IsSpam,
		Rules:     parseRules(reply.Headers),
	}, nil
}

func (s *SpamAssassin) Learn(message []byte, isSpam bool) error {
	var cmd string
	if isSpam {
		cmd = "TELL -t 1 -C SPAM"
	} else {
		cmd = "TELL -t 1 -C HAM"
	}

	_, err := s.client.Tell(context.Background(), bytes.NewReader(message), cmd, nil)
	return err
}

type SpamResult struct {
	Score     float64
	Threshold float64
	IsSpam    bool
	Rules     []SpamRule
}

type SpamRule struct {
	Name        string
	Score       float64
	Description string
}

func parseRules(headers map[string][]string) []SpamRule {
	// Placeholder for parsing rules from headers
	return nil
}
