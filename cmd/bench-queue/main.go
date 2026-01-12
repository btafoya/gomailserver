package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

var (
	host        = flag.String("host", "localhost", "gomailserver host")
	port        = flag.Int("port", 8980, "gomailserver port")
	apikey      = flag.String("apikey", "", "API key for authentication")
	count       = flag.Int("count", 100, "Number of emails to process")
	concurrency = flag.Int("concurrency", 5, "Number of concurrent workers")
	verbose     = flag.Bool("verbose", false, "Verbose output")
)

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if *count <= 0 {
		logger.Fatal("Count must be greater than 0")
	}

	if *apikey == "" {
		logger.Fatal("API key is required for authentication")
	}

	if *concurrency <= 0 {
		logger.Fatal("Concurrency must be greater than 0")
	}

	logger.Info("Starting Queue processing benchmark",
		zap.String("host", *host),
		zap.Int("port", *port),
		zap.Int("count", *count),
		zap.Int("concurrency", *concurrency),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now()

	var processedCount int64
	var processedBytes int64
	var errorCount int64

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *concurrency)

	for i := 0; i < *count; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			err := simulateEmailProcessing(ctx, logger, idx+1)
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				if *verbose {
					logger.Error("Failed to process email",
						zap.Int("number", idx+1),
						zap.Error(err),
					)
				}
			} else {
				atomic.AddInt64(&processedCount, 1)
				atomic.AddInt64(&processedBytes, 1024)
			}

			progress := atomic.LoadInt64(&processedCount)
			if progress%10 == 0 {
				logger.Info("Progress",
					zap.Int64("processed", progress),
					zap.Int64("errors", atomic.LoadInt64(&errorCount)),
					zap.Int64("total", progress+atomic.LoadInt64(&errorCount)),
				)
			}
		}(i)
	}

	wg.Wait()

	duration := time.Since(startTime)
	throughput := float64(atomic.LoadInt64(&processedCount)) / duration.Seconds()
	opsPerSec := float64(atomic.LoadInt64(&processedCount)) / duration.Seconds()

	logger.Info("Queue benchmark completed",
		zap.Duration("duration", duration),
		zap.Int64("total", int64(*count)),
		zap.Int64("processed", atomic.LoadInt64(&processedCount)),
		zap.Int64("errors", atomic.LoadInt64(&errorCount)),
		zap.Float64("throughput_msg_s", throughput),
		zap.Float64("throughput_bytes_s", float64(atomic.LoadInt64(&processedBytes))/duration.Seconds()/1024),
		zap.Float64("ops_s", opsPerSec),
	)

	fmt.Printf("\n=== Queue Processing Benchmark Results ===\n")
	fmt.Printf("Total Emails:         %d\n", *count)
	fmt.Printf("Processed:            %d\n", atomic.LoadInt64(&processedCount))
	fmt.Printf("Failed:                %d\n", atomic.LoadInt64(&errorCount))
	fmt.Printf("Success Rate:          %.2f%%\n", float64(atomic.LoadInt64(&processedCount))/float64(*count)*100)
	fmt.Printf("Duration:               %s\n", duration.Round(time.Millisecond))
	fmt.Printf("Throughput:              %.2f msg/s\n", throughput)
	fmt.Printf("Throughput:              %.2f KB/s\n", float64(atomic.LoadInt64(&processedBytes))/duration.Seconds()/1024)
	fmt.Printf("Average Message Size:     1KB (simulated)")
	fmt.Printf("Operations/s:           %.2f\n", opsPerSec)
	fmt.Printf("Concurrency:              %d\n", *concurrency)
}

func simulateEmailProcessing(ctx context.Context, logger *zap.Logger, number int) error {
	// This is a simulation for benchmarking purposes
	// In real usage, you would call actual gomailserver API endpoint

	// Simulate processing latency based on message number
	latency := time.Duration(5+number%11) * time.Millisecond

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(latency):
		return nil
	}
}
