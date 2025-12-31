package antispam

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/teamwork/spamc"
)

type SpamAssassin struct {
	client *spamc.Client
}

func NewSpamAssassin(host string, port int) *SpamAssassin {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	return &SpamAssassin{
		client: spamc.New(fmt.Sprintf("%s:%d", host, port), dialer),
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
	}, nil
}

func (s *SpamAssassin) Learn(message []byte, isSpam bool) error {
	msgClass := "ham"
	if isSpam {
		msgClass = "spam"
	}

	header := spamc.Header{}.
		Set("Message-class", msgClass).
		Set("Set", "local")

	_, err := s.client.Tell(context.Background(), bytes.NewReader(message), header)
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
