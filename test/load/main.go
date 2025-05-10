// test/load/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rpranjan11/coupon-issuance-system/api/coupon"
	"github.com/rpranjan11/coupon-issuance-system/api/coupon/couponconnect"
)

func main() {
	// Command-line flags
	serverAddr := flag.String("server", "http://localhost:8080", "Server address")
	campaignID := flag.String("campaign", "", "Campaign ID to test (required)")
	concurrency := flag.Int("concurrency", 100, "Number of concurrent clients")
	requestRate := flag.Int("rate", 500, "Target requests per second")
	duration := flag.Duration("duration", 10*time.Second, "Test duration")
	flag.Parse()

	if *campaignID == "" {
		log.Fatal("Campaign ID is required. Use -campaign flag")
	}

	// Create client
	client := couponconnect.NewCouponServiceClient(
		http.DefaultClient,
		*serverAddr,
	)

	fmt.Printf("Starting load test for campaign %s\n", *campaignID)
	fmt.Printf("Concurrency: %d, Target rate: %d req/sec, Duration: %s\n",
		*concurrency, *requestRate, *duration)

	var (
		wg           sync.WaitGroup
		successCount int64
		failCount    int64
		rateLimiter  = time.NewTicker(time.Second / time.Duration(*requestRate) * time.Duration(*concurrency))
		testStart    = time.Now()
		testEnd      = testStart.Add(*duration)
		uniqueCodes  = make(map[string]bool)
		codesMutex   sync.Mutex
	)

	// Start concurrent workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for time.Now().Before(testEnd) {
				// Wait for rate limiter
				<-rateLimiter.C

				// Make the request
				req := connect.NewRequest(&coupon.IssueCouponRequest{
					CampaignId: *campaignID,
				})
				resp, err := client.IssueCoupon(context.Background(), req)

				if err != nil {
					atomic.AddInt64(&failCount, 1)
					continue
				}

				// Check response
				if resp.Msg.Success {
					atomic.AddInt64(&successCount, 1)

					// Track unique codes
					code := resp.Msg.Coupon.Code
					codesMutex.Lock()
					uniqueCodes[code] = true
					codesMutex.Unlock()
				} else {
					atomic.AddInt64(&failCount, 1)
				}
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	rateLimiter.Stop()

	// Calculate metrics
	elapsed := time.Since(testStart)
	totalRequests := successCount + failCount
	rps := float64(totalRequests) / elapsed.Seconds()

	// Print results
	fmt.Println("\nLoad Test Results:")
	fmt.Printf("Duration: %.2f seconds\n", elapsed.Seconds())
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Successful Requests: %d\n", successCount)
	fmt.Printf("Failed Requests: %d\n", failCount)
	fmt.Printf("Requests per second: %.2f\n", rps)

	// Check for duplicate codes (data consistency)
	fmt.Printf("Unique coupon codes issued: %d\n", len(uniqueCodes))

	// Validate consistency
	if int(successCount) != len(uniqueCodes) {
		fmt.Printf("WARNING: Successful requests (%d) doesn't match unique codes count (%d)\n",
			successCount, len(uniqueCodes))
		fmt.Println("This indicates potential duplicate coupon issuance - a concurrency issue!")
	} else {
		fmt.Println("Data consistency check PASSED: All issued coupons had unique codes")
	}

	// Check coupon limit
	// Get the campaign to check if the count matches
	campaignReq := connect.NewRequest(&coupon.GetCampaignRequest{
		CampaignId: *campaignID,
	})
	campaignResp, err := client.GetCampaign(context.Background(), campaignReq)
	if err == nil {
		issuedCount := campaignResp.Msg.Campaign.IssuedCoupons
		fmt.Printf("Campaign issued coupon count: %d\n", issuedCount)

		if int(issuedCount) != len(uniqueCodes) {
			fmt.Printf("WARNING: Campaign issued count (%d) doesn't match unique codes (%d)\n",
				issuedCount, len(uniqueCodes))
			fmt.Println("This indicates a potential race condition in the counter!")
		} else {
			fmt.Println("Counter consistency check PASSED: Campaign count matches issued coupons")
		}
	}
}
