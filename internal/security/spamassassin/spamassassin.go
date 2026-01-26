package spamassassin

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// SpamAssassinServiceInterface defines the interface for spam scanning
type SpamAssassinServiceInterface interface {
	Scan(message []byte) (*ScanResult, error)
}

// ScanResult represents the result of a SpamAssassin scan
type ScanResult struct {
	IsSpam    bool
	Score     float64
	Threshold float64
	Rules     []RuleResult
}

// RuleResult represents a single SpamAssassin rule result
type RuleResult struct {
	Rule        string
	Score       float64
	Description string
}

// SpamAssassin implements SpamAssassin scanning via spamd
type SpamAssassin struct {
	host string
	port int
}

// NewSpamAssassin creates a new SpamAssassin scanner
func NewSpamAssassin(host string, port int) *SpamAssassin {
	return &SpamAssassin{
		host: host,
		port: port,
	}
}

// Scan sends a message to SpamAssassin for analysis
func (s *SpamAssassin) Scan(message []byte) (*ScanResult, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(s.host, strconv.Itoa(s.port)))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to spamd: %w", err)
	}
	defer conn.Close()

	// Send PROCESS command with message
	_, err = fmt.Fprintf(conn, "PROCESS SPAMC/1.2\r\nContent-length: %d\r\n\r\n", len(message))
	if err != nil {
		return nil, fmt.Errorf("failed to send PROCESS command: %w", err)
	}

	// Send the message
	_, err = conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Read response
	scanner := bufio.NewScanner(conn)
	result := &ScanResult{
		Rules: make([]RuleResult, 0),
	}

	// Parse response
	parsingRules := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Spam: ") {
			// Parse spam status line: "Spam: True ; 5.5 / 5.0"
			parts := strings.Split(line, ";")
			if len(parts) >= 2 {
				statusPart := strings.TrimSpace(strings.TrimPrefix(parts[0], "Spam: "))
				result.IsSpam = statusPart == "True" || statusPart == "Yes"

				scorePart := strings.TrimSpace(parts[1])
				scoreParts := strings.Split(scorePart, "/")
				if len(scoreParts) >= 2 {
					if score, err := strconv.ParseFloat(strings.TrimSpace(scoreParts[0]), 64); err == nil {
						result.Score = score
					}
					if threshold, err := strconv.ParseFloat(strings.TrimSpace(scoreParts[1]), 64); err == nil {
						result.Threshold = threshold
					}
				}
			}
		} else if strings.TrimSpace(line) == "" && !parsingRules {
			// Empty line marks end of headers, start of rules
			parsingRules = true
		} else if parsingRules && strings.TrimSpace(line) != "" {
			// Parse rule line: "5.5 TEST_RULE Test rule description"
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				score, err := strconv.ParseFloat(parts[0], 64)
				if err == nil {
					rule := RuleResult{
						Score: score,
						Rule:  parts[1],
					}
					if len(parts) > 2 {
						rule.Description = strings.Join(parts[2:], " ")
					}
					result.Rules = append(result.Rules, rule)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return result, nil
}
