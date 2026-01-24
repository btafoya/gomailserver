package phishing

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
)

const (
	// Default confidence threshold for phishing detection
	DefaultConfidenceThreshold = 0.7
	
	// High-risk indicators that warrant immediate flagging
	HighRiskScoreThreshold = 8.0
	
	// ML Model simulation confidence factors
	BrandConfidenceWeight = 0.3
	LinkAnalysisWeight = 0.25
	ContentAnalysisWeight = 0.25
	MetadataAnalysisWeight = 0.2
)

// PhishingDetectionService provides AI-powered phishing detection
type PhishingDetectionService struct {
	logger *zap.Logger
}

// NewPhishingDetectionService creates a new phishing detection service
func NewPhishingDetectionService(logger *zap.Logger) *PhishingDetectionService {
	return &PhishingDetectionService{
		logger: logger,
	}
}

// PhishingIndicator represents a potential phishing indicator
type PhishingIndicator struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`
	Severity    string  `json:"severity"` // "low", "medium", "high", "critical"
}

// PhishingResult represents the result of phishing analysis
type PhishingResult struct {
	IsPhishing           bool                `json:"is_phishing"`
	OverallConfidence   float64             `json:"overall_confidence"`
	PhishingScore      float64             `json:"phishing_score"`
	RiskLevel          string               `json:"risk_level"` // "low", "medium", "high", "critical"
	Indicators          []PhishingIndicator   `json:"indicators"`
	AnalysisTimestamp   time.Time            `json:"analysis_timestamp"`
	Recommendations    []string             `json:"recommendations"`
}

// AnalyzeMessage performs comprehensive phishing detection using AI techniques
func (s *PhishingDetectionService) AnalyzeMessage(ctx context.Context, message *domain.Message) (*PhishingResult, error) {
	s.logger.Info("Starting AI phishing analysis",
		zap.Int64("message_id", message.ID),
		zap.String("from", message.From),
		zap.String("subject", message.Subject),
	)

	indicators := []PhishingIndicator{}
	var totalScore float64
	var totalConfidence float64

	// 1. Sender Analysis
	senderIndicators, senderScore, senderConfidence := s.analyzeSender(message)
	indicators = append(indicators, senderIndicators...)
	totalScore += senderScore
	totalConfidence += senderConfidence

	// 2. Link Analysis (extract and analyze URLs)
	linkIndicators, linkScore, linkConfidence := s.analyzeLinks(message)
	indicators = append(indicators, linkIndicators...)
	totalScore += linkScore
	totalConfidence += linkConfidence

	// 3. Content Analysis
	contentIndicators, contentScore, contentConfidence := s.analyzeContent(message)
	indicators = append(indicators, contentIndicators...)
	totalScore += contentScore
	totalConfidence += contentConfidence

	// 4. Metadata Analysis
	metadataIndicators, metadataScore, metadataConfidence := s.analyzeMetadata(message)
	indicators = append(indicators, metadataIndicators...)
	totalScore += metadataScore
	totalConfidence += metadataConfidence

	// 5. Calculate overall results
	overallConfidence := totalConfidence / 4.0
	phishingScore := totalScore
	isPhishing := overallConfidence >= DefaultConfidenceThreshold || phishingScore >= HighRiskScoreThreshold

	riskLevel := s.calculateRiskLevel(phishingScore, overallConfidence)
	recommendations := s.generateRecommendations(indicators, riskLevel)

	result := &PhishingResult{
		IsPhishing:         isPhishing,
		OverallConfidence:   overallConfidence,
		PhishingScore:      phishingScore,
		RiskLevel:          riskLevel,
		Indicators:          indicators,
		AnalysisTimestamp:   time.Now(),
		Recommendations:    recommendations,
	}

	s.logger.Info("AI phishing analysis completed",
		zap.Int64("message_id", message.ID),
		zap.Bool("is_phishing", isPhishing),
		zap.Float64("confidence", overallConfidence),
		zap.String("risk_level", riskLevel),
	)

	return result, nil
}

// analyzeSender performs AI-based sender reputation analysis
func (s *PhishingDetectionService) analyzeSender(message *domain.Message) ([]PhishingIndicator, float64, float64) {
	var indicators []PhishingIndicator
	var score float64
	var confidence float64

	// Display name spoofing detection
	displayName := s.extractDisplayName(message.From)
	if displayName == "" {
		// Missing display name - suspicious
		indicators = append(indicators, PhishingIndicator{
			Type:        "display_name_missing",
			Description: "Sender display name is missing or empty",
			Confidence:  0.6,
			Severity:    "medium",
		})
		score += 2.0
		confidence += 0.6
	} else if s.isSuspiciousDisplayName(displayName) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "display_name_spoofing",
			Description: fmt.Sprintf("Potentially spoofed display name: %s", displayName),
			Confidence:  0.8,
			Severity:    "high",
		})
		score += 4.0
		confidence += 0.8 * BrandConfidenceWeight
	}

	// Domain mismatch detection
	senderDomain := s.extractDomain(message.From)
	if senderDomain != "" && s.isHighRiskDomain(senderDomain) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "high_risk_domain",
			Description: fmt.Sprintf("Message from high-risk domain: %s", senderDomain),
			Confidence:  0.9,
			Severity:    "high",
		})
		score += 5.0
		confidence += 0.9 * LinkAnalysisWeight
	}

	// Sender reputation simulation (would integrate with external services)
	if s.isUnusualSender(message.From) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "unusual_sender",
			Description: "Message from unusual or first-time sender",
			Confidence:  0.5,
			Severity:    "medium",
		})
		score += 2.5
		confidence += 0.5
	}

	return indicators, score, confidence
}

// analyzeLinks performs AI-based link analysis for phishing detection
func (s *PhishingDetectionService) analyzeLinks(message *domain.Message) ([]PhishingIndicator, float64, float64) {
	var indicators []PhishingIndicator
	var score float64
	var confidence float64

	// Extract URLs from message body
	urls := s.extractURLs(message.Body)
	urls = append(urls, s.extractURLsFromHeaders(message.Headers)...)

	for _, url := range urls {
		// URL reputation and structure analysis
		if s.isSuspiciousURL(url) {
			indicators = append(indicators, PhishingIndicator{
				Type:        "suspicious_url",
				Description: fmt.Sprintf("Suspicious URL detected: %s", s.sanitizeURL(url)),
				Confidence:  0.85,
				Severity:    "high",
			})
			score += 4.0
			confidence += 0.85
		} else if s.isURLMismatch(url) {
			indicators = append(indicators, PhishingIndicator{
				Type:        "url_display_mismatch",
				Description: fmt.Sprintf("URL display text doesn't match actual URL: %s", s.sanitizeURL(url)),
				Confidence:  0.7,
				Severity:    "medium",
			})
			score += 3.0
			confidence += 0.7
		} else if s.isShortenedURL(url) {
			indicators = append(indicators, PhishingIndicator{
				Type:        "url_shortener",
				Description: "URL uses known shortening service",
				Confidence:  0.4,
				Severity:    "low",
			})
			score += 1.5
			confidence += 0.4
		}
	}

	return indicators, score, confidence
}

// analyzeContent performs AI-based content analysis for phishing patterns
func (s *PhishingDetectionService) analyzeContent(message *domain.Message) ([]PhishingIndicator, float64, float64) {
	var indicators []PhishingIndicator
	var score float64
	var confidence float64

	subject := strings.ToLower(message.Subject)
	body := strings.ToLower(string(message.Content))

	// Urgency and pressure tactics
	if s.containsUrgencyKeywords(subject) || s.containsUrgencyKeywords(body) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "urgency_pressure",
			Description: "Message uses urgency or pressure tactics",
			Confidence:  0.7,
			Severity:    "medium",
		})
		score += 3.0
		confidence += 0.7
	}

	// Generic greeting patterns
	if s.containsGenericGreetings(body) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "generic_greeting",
			Description: "Uses generic greeting instead of personalized communication",
			Confidence:  0.5,
			Severity:    "low",
		})
		score += 1.5
		confidence += 0.5
	}

	// Grammar and spelling issues
	if s.hasPoorGrammar(subject, body) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "poor_grammar",
			Description: "Message contains grammar or spelling errors common in phishing",
			Confidence:  0.6,
			Severity:    "medium",
		})
		score += 2.5
		confidence += 0.6
	}

	// Request for sensitive information
	if s.containsSensitiveRequests(body) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "sensitive_request",
			Description: "Message requests sensitive information (passwords, financial details)",
			Confidence:  0.9,
			Severity:    "high",
		})
		score += 5.0
		confidence += 0.9 * ContentAnalysisWeight
	}

	// Brand impersonation detection
	if s.detectsBrandImpersonation(subject, body, message.From) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "brand_impersonation",
			Description: "Message impersonates well-known brand",
			Confidence:  0.8,
			Severity:    "high",
		})
		score += 4.5
		confidence += 0.8 * BrandConfidenceWeight
	}

	return indicators, score, confidence
}

// analyzeMetadata performs AI-based metadata analysis
func (s *PhishingDetectionService) analyzeMetadata(message *domain.Message) ([]PhishingIndicator, float64, float64) {
	var indicators []PhishingIndicator
	var score float64
	var confidence float64

	// Unusual message structure
	if message.MessageID == "" || message.ReceivedAt.IsZero() {
		indicators = append(indicators, PhishingIndicator{
			Type:        "missing_metadata",
			Description: "Message lacks essential metadata (Message-ID, proper dates)",
			Confidence:  0.4,
			Severity:    "low",
		})
		score += 1.0
		confidence += 0.4
	}

	// Suspicious headers analysis
	if s.hasSuspiciousHeaders(message.Headers) {
		indicators = append(indicators, PhishingIndicator{
			Type:        "suspicious_headers",
			Description: "Message contains suspicious email headers",
			Confidence:  0.7,
			Severity:    "medium",
		})
		score += 3.0
		confidence += 0.7
	}

	return indicators, score, confidence
}

// Helper methods for AI analysis

// isSuspiciousDisplayName checks if display name looks suspicious
func (s *PhishingDetectionService) isSuspiciousDisplayName(displayName string) bool {
	// Check for random strings, excessive capitalization, etc.
	return len(displayName) > 50 || // Display name too long
		regexp.MustCompile(`[A-Z]{3,}`).MatchString(displayName) || // Excessive capitalization
		strings.Contains(displayName, "support") || // Contains "support" (common in phishing)
		strings.Contains(displayName, "security") || // Contains "security" (common in phishing)
}

// extractDisplayName extracts display name from email address
func (s *PhishingDetectionService) extractDisplayName(from string) string {
	// Extract name from "Name <email@domain.com>" format
	re := regexp.MustCompile(`^([^<]+)`)
	matches := re.FindStringSubmatch(from)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractDomain extracts domain from email address
func (s *PhishingDetectionService) extractDomain(from string) string {
	re := regexp.MustCompile(`@([^>]+)`)
	matches := re.FindStringSubmatch(from)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// isHighRiskDomain checks if domain is known for phishing
func (s *PhishingDetectionService) isHighRiskDomain(domain string) bool {
	// List of known high-risk TLDs and patterns
	highRiskPatterns := []string{
		`\.tk$`,        // Tokelau
		`\.ml$`,        // Mali
		`\.ga$`,        // Gabon
		`tempmail-`,      // Temporary email services
		`10minutemail`, // Disposable email services
	}

	for _, pattern := range highRiskPatterns {
		if strings.Contains(strings.ToLower(domain), pattern) {
			return true
		}
	}

	// Check for numeric domains or suspicious patterns
	if regexp.MustCompile(`^\d+\.`).MatchString(domain) {
		return true
	}

	return false
}

// isUnusualSender checks if sender is unusual for this recipient
func (s *PhishingDetectionService) isUnusualSender(from string) bool {
	// This would integrate with sender frequency analysis
	// For now, use simple heuristics
	return strings.Contains(strings.ToLower(from), "noreply") && // No-reply with unusual content
		!strings.Contains(from, "@") // Invalid email format
}

// extractURLs finds all URLs in message content
func (s *PhishingDetectionService) extractURLs(content []byte) []string {
	// URL extraction regex
	urlRegex := regexp.MustCompile(`https?://[^\s<>"']`)
	matches := urlRegex.FindAllString(string(content))
	
	var urls []string
	seen := make(map[string]bool)
	for _, match := range matches {
		if !seen[match] {
			urls = append(urls, match)
			seen[match] = true
		}
	}
	
	return urls
}

// extractURLsFromHeaders extracts URLs from message headers
func (s *PhishingDetectionService) extractURLsFromHeaders(headers string) []string {
	if headers == "" {
		return nil
	}
	
	urlRegex := regexp.MustCompile(`https?://[^\s<>"']`)
	return urlRegex.FindAllString(headers)
}

// isSuspiciousURL performs URL analysis for phishing indicators
func (s *PhishingDetectionService) isSuspiciousURL(url string) bool {
	lowerURL := strings.ToLower(url)
	
	// IP address instead of domain
	ipRegex := regexp.MustCompile(`^https?://\d+\.\d+\.\d+\.\d+`)
	if ipRegex.MatchString(lowerURL) {
		return true
	}
	
	// Missing HTTPS
	if strings.HasPrefix(lowerURL, "http://") && !strings.Contains(lowerURL, "localhost") {
		return true
	}
	
	// Suspicious TLDs
	suspiciousTLDs := []string{".tk", ".ml", ".ga", ".cf"}
	for _, tld := range suspiciousTLDs {
		if strings.HasSuffix(lowerURL, tld) {
			return true
		}
	}
	
	// URL shorteners with suspicious patterns
	if strings.Contains(lowerURL, "bit.ly") && 
	   (strings.Contains(lowerURL, "login") || 
		strings.Contains(lowerURL, "verify") || 
		strings.Contains(lowerURL, "secure")) {
		return true
	}
	
	return false
}

// isURLMismatch checks if displayed URL differs from actual URL
func (s *PhishingDetectionService) isURLMismatch(url string) bool {
	// This would compare displayed text vs actual href
	// Simplified implementation for demonstration
	return strings.Contains(url, "click HERE") || 
		   strings.Contains(url, "immediate action required")
}

// isShortenedURL checks if URL uses known shortening services
func (s *PhishingDetectionService) isShortenedURL(url string) bool {
	shorteners := []string{
		"bit.ly", "tinyurl.com", "goo.gl", "t.co",
		"ow.ly", "is.gd", "buff.ly", "adf.ly",
	}
	
	lowerURL := strings.ToLower(url)
	for _, shortener := range shorteners {
		if strings.Contains(lowerURL, shortener) {
			return true
		}
	}
	
	return false
}

// containsUrgencyKeywords checks for urgency tactics
func (s *PhishingDetectionService) containsUrgencyKeywords(text string) bool {
	keywords := []string{
		"urgent", "immediately", "asap", "hurry",
		"expire", "suspend", "terminate", "cancel",
		"verification required", "account locked", "security breach",
		"unusual activity", "limited time", "act now",
	}
	
	lowerText := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(lowerText, keyword) {
			return true
		}
	}
	
	return false
}

// containsGenericGreetings detects generic greeting patterns
func (s *PhishingDetectionService) containsGenericGreetings(text string) bool {
	greetings := []string{
		"dear user", "dear customer", "dear valued customer",
		"dear account holder", "hello user", "hello customer",
		"attention user", "important notice", "security alert",
	}
	
	lowerText := strings.ToLower(text)
	for _, greeting := range greetings {
		if strings.Contains(lowerText, greeting) {
			return true
		}
	}
	
	return false
}

// hasPoorGrammar detects grammar and spelling issues
func (s *PhishingDetectionService) hasPoorGrammar(subject, body string) bool {
	// Simplified grammar checking (would use ML model in production)
	text := subject + " " + body
	
	indicators := []string{
		// Common grammatical errors in phishing
		"we have notice", "your account will be", "click here immediate",
		"dear beneficiary", "congratulations winner", "you have been selected",
		// Poor sentence structure indicators
		 regexp.MustCompile(`\b[A-Z]{3,}\s+[A-Z]{3,}`).MatchString(text), // Random capitalization
		strings.Contains(text, "!!") || // Excessive punctuation
	}
	
	for _, indicator := range indicators {
		if strings.Contains(strings.ToLower(text), indicator) {
			return true
		}
	}
	
	return false
}

// containsSensitiveRequests detects requests for sensitive information
func (s *PhishingDetectionService) containsSensitiveRequests(text string) bool {
	sensitiveTerms := []string{
		"password", "credit card", "social security", "bank account",
		"routing number", "account number", "pin code",
		"verify identity", "confirm account", "update payment",
		"transaction details", "invoice download", "wire transfer",
		"personal information", "confidential documents",
	}
	
	lowerText := strings.ToLower(text)
	for _, term := range sensitiveTerms {
		if strings.Contains(lowerText, term) {
			return true
		}
	}
	
	return false
}

// detectsBrandImpersonation detects brand impersonation attempts
func (s *PhishingDetectionService) detectsBrandImpersonation(subject, body, from string) bool {
	text := strings.ToLower(subject + " " + body)
	lowerFrom := strings.ToLower(from)
	
	brands := []string{
		"microsoft", "apple", "google", "amazon", "paypal",
		"facebook", "twitter", "instagram", "linkedin", "bank",
		"chase", "wells fargo", "bank of america", "citibank",
	}
	
	// Check for brand mentions without proper authorization
	for _, brand := range brands {
		if strings.Contains(text, brand) &&
		   !strings.Contains(lowerFrom, brand) && // Sender not from brand domain
		   (strings.Contains(text, "verification") || 
			strings.Contains(text, "security") || 
			strings.Contains(text, "suspend") || 
			strings.Contains(text, "unusual activity")) {
			return true
		}
	}
	
	return false
}

// hasSuspiciousHeaders analyzes email headers for red flags
func (s *PhishingDetectionService) hasSuspiciousHeaders(headers string) bool {
	suspiciousPatterns := []string{
		// Header manipulation indicators
		"x-auto-generated-sender:", "x-spoofed-sender:",
		"list-unsubscribe:", "list-post:", "bulk-mail:",
		// Missing or unusual authentication headers
		"authentication-results:", "dkim-signature:", "received-spf:",
	}
	
	lowerHeaders := strings.ToLower(headers)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerHeaders, pattern) {
			return true
		}
	}
	
	return false
}

// calculateRiskLevel determines risk level based on score and confidence
func (s *PhishingDetectionService) calculateRiskLevel(score float64, confidence float64) string {
	if confidence >= 0.9 || score >= HighRiskScoreThreshold {
		return "critical"
	} else if confidence >= 0.7 || score >= 6.0 {
		return "high"
	} else if confidence >= 0.5 || score >= 3.0 {
		return "medium"
	}
	return "low"
}

// generateRecommendations creates security recommendations based on analysis
func (s *PhishingDetectionService) generateRecommendations(indicators []PhishingIndicator, riskLevel string) []string {
	var recommendations []string
	
	switch riskLevel {
	case "critical", "high":
		recommendations = append(recommendations,
			"Do not click any links in the message",
			"Verify sender identity through separate channel",
			"Contact IT security department immediately",
			"Delete the message without responding",
		)
	case "medium":
		recommendations = append(recommendations,
			"Be cautious with any links or requests",
			"Verify information through official website",
			"Check sender's email address format",
		)
	case "low":
		recommendations = append(recommendations,
			"Consider verifying sender if unexpected",
			"Be cautious with urgent requests",
		)
	}
	
	// Add specific recommendations based on indicators
	for _, indicator := range indicators {
		switch indicator.Type {
		case "suspicious_url":
			recommendations = append(recommendations, "Hover over links to preview actual URL before clicking")
		case "display_name_spoofing":
			recommendations = append(recommendations, "Check sender's email address for impersonation")
		case "brand_impersonation":
			recommendations = append(recommendations, "Verify brand communication through official channels")
		}
	}
	
	return recommendations
}

// sanitizeURL removes sensitive information from URLs for logging
func (s *PhishingDetectionService) sanitizeURL(url string) string {
	// Remove query parameters and truncate for logging
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}
	if len(url) > 100 {
		url = url[:100] + "..."
	}
	return url
}