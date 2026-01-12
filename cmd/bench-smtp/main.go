package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/smtp"
	"net/textproto"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

var (
	host        = flag.String("host", "localhost", "SMTP server host")
	port        = flag.Int("port", 25, "SMTP server port")
	sender      = flag.String("sender", "test@example.com", "Sender email")
	recipient   = flag.String("recipient", "recipient@example.com", "Recipient email")
	count       = flag.Int("count", 100, "Number of emails to send")
	concurrency = flag.Int("concurrency", 10, "Number of concurrent connections")
	subject     = flag.String("subject", "SMTP Benchmark Test", "Email subject")
	body        = flag.String("body", "This is a test email for SMTP benchmarking.", "Email body")
	verbose     = flag.Bool("verbose", false, "Verbose output")
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

	logger.Info("Starting SMTP throughput benchmark",
		zap.String("host", *host),
		zap.Int("port", *port),
		zap.String("sender", *sender),
		zap.String("recipient", *recipient),
		zap.Int("count", *count),
		zap.Int("concurrency", *concurrency),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now()

	var successCount int64
	var failCount int64
	var totalBytes int64

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *concurrency)

	for i := 0; i < *count; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			err := sendEmail(ctx, logger, idx+1)
			if err != nil {
				atomic.AddInt64(&failCount, 1)
				if *verbose {
					logger.Error("Failed to send email",
						zap.Int("number", idx+1),
						zap.Error(err),
					)
				}
			} else {
				atomic.AddInt64(&successCount, 1)
				atomic.AddInt64(&totalSize, int64(len(*body)))
			}

			progress := atomic.LoadInt64(&successCount)
			if progress%10 == 0 {
				logger.Info("Progress",
					zap.Int64("sent", progress),
					zap.Int64("failed", atomic.LoadInt64(&failCount)),
					zap.Int64("total", progress+atomic.LoadInt64(&failCount)),
				)
			}
		}(i)
	}

	wg.Wait()

	duration := time.Since(startTime)

	throughput := float64(atomic.LoadInt64(&successCount)) / duration.Seconds()
	avgSize := float64(atomic.LoadInt64(&totalSize)) / float64(atomic.LoadInt64(&successCount))

	logger.Info("SMTP benchmark completed",
		zap.Duration("duration", duration),
		zap.Int64("total", int64(*count)),
		zap.Int64("success", atomic.LoadInt64(&successCount)),
		zap.Int64("failed", atomic.LoadInt64(&failCount)),
		zap.Float64("throughput_msg_s", throughput),
		zap.Float64("throughput_bytes_s", float64(atomic.LoadInt64(&totalBytes))/duration.Seconds()),
		zap.Float64("avg_size_bytes", avgSize),
	)

	fmt.Printf("\n=== SMTP Throughput Benchmark Results ===\n")
	fmt.Printf("Total Messages:      %d\n", *count)
	fmt.Printf("Success:              %d\n", atomic.LoadInt64(&successCount))
	fmt.Printf("Failed:               %d\n", atomic.LoadInt64(&failCount))
	fmt.Printf("Success Rate:         %.2f%%\n", float64(atomic.LoadInt64(&successCount))/float64(*count)*100)
	fmt.Printf("Duration:             %s\n", duration.Round(time.Millisecond))
	fmt.Printf("Throughput:            %.2f msg/s\n", throughput)
	fmt.Printf("Throughput:            %.2f KB/s\n", float64(atomic.LoadInt64(&totalBytes))/duration.Seconds()/1024)
	fmt.Printf("Average Message Size: %.2f bytes\n", avgSize)
	fmt.Printf("Concurrency:           %d\n", *concurrency)
}

func sendEmail(ctx context.Context, logger *zap.Logger, number int) error {
	smtpAddr := fmt.Sprintf("%s:%d", *host, *port)

	client, err := smtp.Dial(smtpAddr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	headers := make(map[string]string)
	headers["From"] = *sender
	headers["To"] = *recipient
	headers["Subject"] = *subject
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	headers["Message-ID"] = fmt.Sprintf("<%d@benchmark.example.com>", number)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMessage-ID: %s\r\n\r\n%s",
		*sender, *recipient, *subject, headers["Date"], headers["Message-ID"], *body)

	if err := client.Mail(*sender); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(*recipient); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := wc.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}
