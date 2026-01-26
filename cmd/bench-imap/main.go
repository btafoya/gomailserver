package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

var (
	host        = flag.String("host", "localhost", "IMAP server host")
	port        = flag.Int("port", 143, "IMAP server port")
	username    = flag.String("username", "test@example.com", "IMAP username")
	password    = flag.String("password", "password", "IMAP password")
	count       = flag.Int("count", 50, "Number of messages to process")
	concurrency = flag.Int("concurrency", 5, "Number of concurrent connections")
	mailbox     = flag.String("mailbox", "INBOX", "Mailbox to test")
	verbose     = flag.Bool("verbose", false, "Verbose output")
)

// Package-level counters for atomic operations across goroutines
var (
	successCount int64
	failCount    int64
	totalBytes   int64
	opCount      int64
)

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if *count <= 0 {
		logger.Fatal("Count must be greater than 0")
	}

	if *concurrency <= 0 {
		logger.Fatal("Concurrency must be greater than 0")
	}

	logger.Info("Starting IMAP throughput benchmark",
		zap.String("host", *host),
		zap.Int("port", *port),
		zap.String("mailbox", *mailbox),
		zap.Int("count", *count),
		zap.Int("concurrency", *concurrency),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *concurrency)

	for i := 0; i < *count; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			err := processMessage(ctx, logger, idx+1)
			if err != nil {
				atomic.AddInt64(&failCount, 1)
				if *verbose {
					logger.Error("Failed to process message",
						zap.Int("number", idx+1),
						zap.Error(err),
					)
				}
			} else {
				atomic.AddInt64(&successCount, 1)
				atomic.AddInt64(&opCount, 1)
			}

			progress := atomic.LoadInt64(&successCount)
			if progress%10 == 0 {
				logger.Info("Progress",
					zap.Int64("processed", progress),
					zap.Int64("operations", atomic.LoadInt64(&opCount)),
				)
			}
		}(i)
	}

	wg.Wait()

	duration := time.Since(startTime)

	throughput := float64(atomic.LoadInt64(&successCount)) / duration.Seconds()
	opsPerSec := float64(atomic.LoadInt64(&opCount)) / duration.Seconds()

	logger.Info("IMAP benchmark completed",
		zap.Duration("duration", duration),
		zap.Int64("total", int64(*count)),
		zap.Int64("success", atomic.LoadInt64(&successCount)),
		zap.Int64("failed", atomic.LoadInt64(&failCount)),
		zap.Float64("throughput_msg_s", throughput),
		zap.Float64("operations_s", opsPerSec),
		zap.Float64("avg_size_bytes", float64(atomic.LoadInt64(&totalBytes))/float64(atomic.LoadInt64(&successCount))),
	)

	fmt.Printf("\n=== IMAP Throughput Benchmark Results ===\n")
	fmt.Printf("Total Messages:      %d\n", *count)
	fmt.Printf("Success:              %d\n", atomic.LoadInt64(&successCount))
	fmt.Printf("Failed:               %d\n", atomic.LoadInt64(&failCount))
	fmt.Printf("Success Rate:         %.2f%%\n", float64(atomic.LoadInt64(&successCount))/float64(*count)*100)
	fmt.Printf("Duration:             %s\n", duration.Round(time.Millisecond))
	fmt.Printf("Throughput:            %.2f msg/s\n", throughput)
	fmt.Printf("Operations:            %.2f ops/s\n", opsPerSec)
	fmt.Printf("Average Message Size: %.2f bytes\n", float64(atomic.LoadInt64(&totalBytes))/float64(atomic.LoadInt64(&successCount)))
	fmt.Printf("Concurrency:           %d\n", *concurrency)
	fmt.Printf("Mailbox:              %s\n", *mailbox)
}

func processMessage(ctx context.Context, logger *zap.Logger, number int) error {
	imapAddr := net.JoinHostPort(*host, strconv.Itoa(*port))

	client, err := net.Dial("tcp", imapAddr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	reader := bufio.NewReader(client)

	// Read greeting
	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read greeting: %w", err)
	}

	// Login
	cmd := fmt.Sprintf("LOGIN %s %s", *username, *password)
	_, err = fmt.Fprintf(client, "%s\r\n", cmd)
	if err != nil {
		return fmt.Errorf("failed to send LOGIN: %w", err)
	}

	// Read login response
	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	// Select mailbox
	cmd = fmt.Sprintf("SELECT %s", *mailbox)
	_, err = fmt.Fprintf(client, "%s\r\n", cmd)
	if err != nil {
		return fmt.Errorf("failed to send SELECT: %w", err)
	}

	// Read select response
	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read select response: %w", err)
	}

	// Search for messages (simulate processing)
	cmd = fmt.Sprintf("SEARCH UNSEEN ALL")
	_, err = fmt.Fprintf(client, "%s\r\n", cmd)
	if err != nil {
		return fmt.Errorf("failed to send SEARCH: %w", err)
	}

	// Read search response
	searchResponse, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read search response: %w", err)
	}

	atomic.AddInt64(&totalBytes, int64(len(searchResponse)))

	// Fetch message (UID of first result)
	lines := strings.Split(searchResponse, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("no messages found")
	}

	uidLine := strings.TrimSpace(lines[1])
	if !strings.HasPrefix(uidLine, "* SEARCH") {
		return fmt.Errorf("unexpected search response")
	}

	// Fetch message body
	cmd = fmt.Sprintf("FETCH %s BODY", strings.TrimSpace(uidLine))
	_, err = fmt.Fprintf(client, "%s\r\n", cmd)
	if err != nil {
		return fmt.Errorf("failed to send FETCH: %w", err)
	}

	// Read fetch response
	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read fetch response: %w", err)
	}

	// Logout
	cmd = "LOGOUT"
	_, err = fmt.Fprintf(client, "%s\r\n", cmd)
	if err != nil {
		return fmt.Errorf("failed to send LOGOUT: %w", err)
	}

	return nil
}
