package sieve

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
)

// Parser handles Sieve script parsing and execution
type Parser struct {
	logger *zap.Logger
}

// NewParser creates a new Sieve parser
func NewParser(logger *zap.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

// Filter represents a Sieve filter rule
type Filter struct {
	ID         string
	Conditions []Condition
	Actions    []Action
	Script     string
}

// Condition represents a Sieve condition
type Condition struct {
	Type     string
	Field    string
	Operator string
	Value    string
	Not      bool
}

// Action represents a Sieve action
type Action struct {
	Type  string
	Value string
}

// Parse parses a Sieve script into filter rules
func (p *Parser) Parse(script string) ([]Filter, error) {
	p.logger.Debug("parsing Sieve script", zap.String("script", script))

	var filters []Filter
	lines := strings.Split(script, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse if-elsif-else structure
		if strings.Contains(line, "if ") {
			filter, err := p.parseIfStatement(line)
			if err != nil {
				p.logger.Error("failed to parse if statement", zap.Error(err))
				continue
			}
			filters = append(filters, filter)
		}
	}

	p.logger.Debug("parsed Sieve filters", zap.Int("count", len(filters)))
	return filters, nil
}

// parseIfStatement parses an if statement
func (p *Parser) parseIfStatement(line string) (Filter, error) {
	// Simple pattern matching for common Sieve constructs

	// Match: if header :contains "X-Spam-Flag" "YES" {
	re := regexp.MustCompile(`if\s+(not\s+)?(\w+)\s+"([^"]+)"\s+"([^"]+)"`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return Filter{}, fmt.Errorf("invalid if statement: %s", line)
	}

	condition := Condition{
		Type: matches[2],
		Not:  matches[1] != "",
	}

	// Handle different condition types
	switch matches[2] {
	case "header":
		parts := strings.SplitN(matches[3], " ", 2)
		if len(parts) != 2 {
			return Filter{}, fmt.Errorf("invalid header condition: %s", matches[3])
		}
		condition.Field = parts[0]
		condition.Operator = parts[1]
		condition.Value = strings.Trim(matches[4], `"`)
	case "address":
		condition.Field = "from"
		condition.Operator = matches[3]
		condition.Value = strings.Trim(matches[4], `"`)
	case "size":
		condition.Field = "size"
		condition.Operator = matches[3]
		condition.Value = strings.Trim(matches[4], `"`)
	default:
		condition.Field = matches[2]
		condition.Operator = matches[3]
		condition.Value = strings.Trim(matches[4], `"`)
	}

		// Parse action (simple fileinto)
	action := Action{}
	if strings.Contains(line, "fileinto") {
		actionRe := regexp.MustCompile(`fileinto\s+"([^"]+)"`)
		actionMatches := actionRe.FindStringSubmatch(line)
		if actionMatches != nil {
			action.Type = "fileinto"
			action.Value = strings.Trim(actionMatches[1], `"`)
		}
	} else if strings.Contains(line, "discard") {
		action.Type = "discard"
	} else if strings.Contains(line, "redirect") {
		actionRe := regexp.MustCompile(`redirect\s+"([^"]+)"`)
		actionMatches := actionRe.FindStringSubmatch(line)
		if actionMatches != nil {
			action.Type = "redirect"
			action.Value = strings.Trim(actionMatches[1], `"`)
		}
	}

	return Filter{
		ID:         fmt.Sprintf("filter_%d", len(parsedFilters)),
		Conditions: []Condition{condition},
		Actions:    []Action{action},
		Script:     line,
	}, nil
}

// Evaluate applies filters to a message
func (p *Parser) Evaluate(filters []Filter, message *domain.Message) ([]Action, error) {
	p.logger.Debug("evaluating Sieve filters",
		zap.Int("filter_count", len(filters)),
		zap.String("message_id", message.MessageID),
	)

	for _, filter := range filters {
		// Evaluate all conditions
		conditionsMet := true
		for _, condition := range filter.Conditions {
			met, err := p.evaluateCondition(condition, message)
			if err != nil {
				p.logger.Error("failed to evaluate condition", zap.Error(err))
				continue
			}

			if condition.Not {
				met = !met
			} else {
				conditionsMet = conditionsMet && met
			}

			if !conditionsMet {
				break
			}
		}

		// If all conditions are met, execute actions
		if conditionsMet {
			p.logger.Debug("filter conditions met", zap.String("filter_id", filter.ID))
			return filter.Actions, nil
		}
	}

	return nil, nil
}

// evaluateCondition evaluates a single condition
func (p *Parser) evaluateCondition(condition Condition, message *domain.Message) (bool, error) {
	switch condition.Field {
	case "header":
		return p.evaluateHeaderCondition(condition, message)
	case "from", "address":
		return p.evaluateAddressCondition(condition, message, "from")
	case "subject":
		return p.evaluateTextCondition(condition, message, "subject")
	case "body":
		return p.evaluateTextCondition(condition, message, "body")
	case "size":
		return p.evaluateSizeCondition(condition, message)
	default:
		return false, fmt.Errorf("unsupported condition field: %s", condition.Field)
	}
}

// evaluateHeaderCondition evaluates header-based conditions
func (p *Parser) evaluateHeaderCondition(condition Condition, message *domain.Message) (bool, error) {
	var headerValue string
	switch strings.ToLower(condition.Field) {
	case "subject":
		headerValue = message.Subject
	case "from":
		headerValue = message.From
	case "to":
		if len(message.To) > 0 {
			headerValue = message.To
		} else {
			headerValue = ""
		}
	default:
		return false, fmt.Errorf("unsupported header field: %s", condition.Field)
	}

	return p.evaluateStringCondition(condition.Operator, headerValue, condition.Value)
}
	default:
		return false, fmt.Errorf("unsupported header field: %s", condition.Field)
	}

	return p.evaluateStringCondition(condition.Operator, headerValue, condition.Value)
}

// evaluateAddressCondition evaluates address-based conditions
func (p *Parser) evaluateAddressCondition(condition Condition, message *domain.Message, addressType string) (bool, error) {
	var address string
	switch addressType {
	case "from":
		address = message.Sender
	default:
		return false, fmt.Errorf("unsupported address type: %s", addressType)
	}

	return p.evaluateStringCondition(condition.Operator, address, condition.Value)
}

// evaluateTextCondition evaluates text-based conditions
func (p *Parser) evaluateTextCondition(condition Condition, message *domain.Message, fieldType string) (bool, error) {
	var text string
	switch fieldType {
	case "subject":
		text = message.Subject
	case "body":
		text = message.Body
	default:
		return false, fmt.Errorf("unsupported text field: %s", fieldType)
	}

	return p.evaluateStringCondition(condition.Operator, text, condition.Value)
}

// evaluateSizeCondition evaluates size-based conditions
func (p *Parser) evaluateSizeCondition(condition Condition, message *domain.Message) (bool, error) {
	size := len(message.Body)

	var targetSize int
	_, err := fmt.Sscanf(condition.Value, "%d", &targetSize)
	if err != nil {
		return false, fmt.Errorf("invalid size value: %s", condition.Value)
	}

	switch condition.Operator {
	case "over":
		return size > targetSize, nil
	case "under":
		return size < targetSize, nil
	default:
		return false, fmt.Errorf("unsupported size operator: %s", condition.Operator)
	}
}

// evaluateStringCondition evaluates string-based conditions
func (p *Parser) evaluateStringCondition(operator, value, target string) (bool, error) {
	switch strings.ToLower(operator) {
	case "contains", ":contains":
		return strings.Contains(strings.ToLower(value), strings.ToLower(target)), nil
	case "matches", ":matches":
		matched, _ := regexp.MatchString(target, value)
		return matched, nil
	case "is", ":is":
		return strings.EqualFold(value, target), nil
	case "starts", ":starts":
		return strings.HasPrefix(strings.ToLower(value), strings.ToLower(target)), nil
	case "ends", ":ends":
		return strings.HasSuffix(strings.ToLower(value), strings.ToLower(target)), nil
	default:
		return false, fmt.Errorf("unsupported string operator: %s", operator)
	}
}
