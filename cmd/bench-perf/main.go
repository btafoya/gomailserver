package main

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/emersion/go-imap/client"
)

type BenchmarkConfig struct {
	Host         string
	SMTPPort     int
	IMAPPort     int
	Username     string
	Password     string
	MessageCount int
	Concurrency  int
}

func main() {
	config := parseArgs()

	fmt.Printf("Starting gomailserver benchmark\n")
	fmt.Printf("Target: %s:%d (SMTP), %s:%d (IMAP)\n", config.Host, config.SMTPPort, config.Host, config.IMAPPort)
	fmt.Printf("Messages: %d, Concurrency: %d\n\n", config.MessageCount, config.Concurrency)

	// SMTP Benchmark
	fmt.Println("=== SMTP Benchmark ===")
	smtpResults := benchmarkSMTP(config)
	printSMTPResults(smtpResults)

	// IMAP Benchmark
	fmt.Println("\n=== IMAP Benchmark ===")
	imapResults := benchmarkIMAP(config)
	printIMAPResults(imapResults)

	fmt.Println("\nBenchmark complete")
}

func parseArgs() *BenchmarkConfig {
	config := &BenchmarkConfig{
		Host:         getEnv("BENCH_HOST", "localhost"),
		SMTPPort:     getEnvInt("BENCH_SMTP_PORT", 25),
		IMAPPort:     getEnvInt("BENCH_IMAP_PORT", 143),
		Username:     getEnv("BENCH_USER", "test@example.com"),
		Password:     getEnv("BENCH_PASS", "password"),
		MessageCount: getEnvInt("BENCH_MESSAGES", 100),
		Concurrency:  getEnvInt("BENCH_CONCURRENCY", 10),
	}
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

type SMTPResults struct {
	TotalTime    time.Duration
	MessagesSent int
	SuccessRate  float64
	AvgLatency   time.Duration
	Errors       []error
}

func benchmarkSMTP(config *BenchmarkConfig) *SMTPResults {
	results := &SMTPResults{
		Errors: make([]error, 0),
	}

	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex
	var sent int
	var latencies []time.Duration

	// Worker pool
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < config.MessageCount/config.Concurrency; j++ {
				msgStart := time.Now()

				err := sendSMTPMessage(config, workerID, j)
				latency := time.Since(msgStart)

				mu.Lock()
				latencies = append(latencies, latency)
				if err != nil {
					results.Errors = append(results.Errors, err)
				} else {
					sent++
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	results.TotalTime = time.Since(start)
	results.MessagesSent = sent
	results.SuccessRate = float64(sent) / float64(config.MessageCount) * 100

	// Calculate average latency
	if len(latencies) > 0 {
		var total time.Duration
		for _, lat := range latencies {
			total += lat
		}
		results.AvgLatency = total / time.Duration(len(latencies))
	}

	return results
}

func sendSMTPMessage(config *BenchmarkConfig, workerID, messageID int) error {
	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", config.Host, config.SMTPPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Send message
	from := config.Username
	to := []string{"benchmark@example.com"}

	subject := fmt.Sprintf("Benchmark Message %d-%d", workerID, messageID)
	body := fmt.Sprintf("This is benchmark message %d from worker %d at %s", messageID, workerID, time.Now().Format(time.RFC3339))

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to[0], subject, body)

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	if err := client.Rcpt(to[0]); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close DATA: %w", err)
	}

	client.Quit()
	return nil
}

func printSMTPResults(results *SMTPResults) {
	fmt.Printf("Total Time: %v\n", results.TotalTime)
	fmt.Printf("Messages Sent: %d/%d (%.1f%%)\n", results.MessagesSent, results.MessagesSent+len(results.Errors), results.SuccessRate)
	fmt.Printf("Average Latency: %v\n", results.AvgLatency)
	fmt.Printf("Throughput: %.1f msg/sec\n", float64(results.MessagesSent)/results.TotalTime.Seconds())

	if len(results.Errors) > 0 {
		fmt.Printf("Errors: %d\n", len(results.Errors))
		for i, err := range results.Errors {
			if i >= 5 { // Show only first 5 errors
				fmt.Printf("... and %d more errors\n", len(results.Errors)-5)
				break
			}
			fmt.Printf("  - %v\n", err)
		}
	}
}

type IMAPResults struct {
	TotalTime      time.Duration
	Connections    int
	AvgConnectTime time.Duration
	Errors         []error
}

func benchmarkIMAP(config *BenchmarkConfig) *IMAPResults {
	results := &IMAPResults{
		Errors: make([]error, 0),
	}

	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex
	var connectTimes []time.Duration

	// Concurrent IMAP connections
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			connectStart := time.Now()
			err := testIMAPConnection(config)
			connectTime := time.Since(connectStart)

			mu.Lock()
			connectTimes = append(connectTimes, connectTime)
			if err != nil {
				results.Errors = append(results.Errors, err)
			} else {
				results.Connections++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	results.TotalTime = time.Since(start)

	// Calculate average connect time
	if len(connectTimes) > 0 {
		var total time.Duration
		for _, t := range connectTimes {
			total += t
		}
		results.AvgConnectTime = total / time.Duration(len(connectTimes))
	}

	return results
}

func testIMAPConnection(config *BenchmarkConfig) error {
	addr := fmt.Sprintf("%s:%d", config.Host, config.IMAPPort)

	client, err := client.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Logout()

	if err := client.Login(config.Username, config.Password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Select INBOX
	_, err = client.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("failed to select INBOX: %w", err)
	}

	return nil
}

func printIMAPResults(results *IMAPResults) {
	fmt.Printf("Total Time: %v\n", results.TotalTime)
	fmt.Printf("Successful Connections: %d/%d\n", results.Connections, results.Connections+len(results.Errors))
	fmt.Printf("Average Connect Time: %v\n", results.AvgConnectTime)

	if len(results.Errors) > 0 {
		fmt.Printf("Errors: %d\n", len(results.Errors))
		for i, err := range results.Errors {
			if i >= 5 { // Show only first 5 errors
				fmt.Printf("... and %d more errors\n", len(results.Errors)-5)
				break
			}
			fmt.Printf("  - %v\n", err)
		}
	}
}
