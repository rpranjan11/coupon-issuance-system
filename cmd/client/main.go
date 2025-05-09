// cmd/client/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	coupon "github.com/rpranjan11/coupon-issuance-system/api/coupon"
	"github.com/rpranjan11/coupon-issuance-system/api/coupon/couponconnect"
)

func main() {
	// Parse command-line flags
	serverAddr := flag.String("server", "http://localhost:8080", "server address")
	command := flag.String("command", "issue", "command to run: create, get, issue, or delete")
	campaignID := flag.String("campaign", "", "campaign ID for get, issue, and delete commands")
	campaignName := flag.String("name", "Test Campaign", "campaign name for create and delete commands")
	totalCoupons := flag.Int("total", 10, "total coupons for create command")
	startIn := flag.Duration("start-in", 0, "start time in duration from now for create command")
	flag.Parse()

	// Create HTTP client
	client := couponconnect.NewCouponServiceClient(
		http.DefaultClient,
		*serverAddr,
	)

	// Execute the requested command
	switch *command {
	case "create":
		// Calculate start time
		startTime := time.Now().Add(*startIn)

		// Create request
		req := connect.NewRequest(&coupon.CreateCampaignRequest{
			Name:         *campaignName,
			TotalCoupons: int32(*totalCoupons),
			StartTime:    timestamppb.New(startTime),
		})

		// Call API
		resp, err := client.CreateCampaign(context.Background(), req)
		if err != nil {
			log.Fatalf("Error creating campaign: %v", err)
		}

		// Print response
		fmt.Printf("Campaign created successfully!\n")
		fmt.Printf("ID: %s\n", resp.Msg.Campaign.Id)
		fmt.Printf("Name: %s\n", resp.Msg.Campaign.Name)
		fmt.Printf("Total Coupons: %d\n", resp.Msg.Campaign.TotalCoupons)
		fmt.Printf("Start Time: %s\n", resp.Msg.Campaign.StartTime.AsTime().Format(time.RFC3339))

	case "get":
		// Validate campaign ID
		if *campaignID == "" {
			log.Fatal("Campaign ID is required for get command")
		}

		// Create request
		req := connect.NewRequest(&coupon.GetCampaignRequest{
			CampaignId: *campaignID,
		})

		// Call API
		resp, err := client.GetCampaign(context.Background(), req)
		if err != nil {
			log.Fatalf("Error getting campaign: %v", err)
		}

		// Print campaign details
		fmt.Printf("Campaign Details:\n")
		fmt.Printf("ID: %s\n", resp.Msg.Campaign.Id)
		fmt.Printf("Name: %s\n", resp.Msg.Campaign.Name)
		fmt.Printf("Total Coupons: %d\n", resp.Msg.Campaign.TotalCoupons)
		fmt.Printf("Issued Coupons: %d\n", resp.Msg.Campaign.IssuedCoupons)
		fmt.Printf("Start Time: %s\n", resp.Msg.Campaign.StartTime.AsTime().Format(time.RFC3339))

		// Print coupons
		fmt.Printf("\nIssued Coupons (%d):\n", len(resp.Msg.Coupons))
		for i, c := range resp.Msg.Coupons {
			fmt.Printf("%d. Code: %s, Issued At: %s\n",
				i+1, c.Code, c.IssuedAt.AsTime().Format(time.RFC3339))
		}

	case "issue":
		// Validate campaign ID
		if *campaignID == "" {
			log.Fatal("Campaign ID is required for issue command")
		}

		// Create request
		req := connect.NewRequest(&coupon.IssueCouponRequest{
			CampaignId: *campaignID,
		})

		// Call API
		resp, err := client.IssueCoupon(context.Background(), req)
		if err != nil {
			log.Fatalf("Error issuing coupon: %v", err)
		}

		// Print result
		if resp.Msg.Success {
			fmt.Printf("Coupon issued successfully!\n")
			fmt.Printf("Code: %s\n", resp.Msg.Coupon.Code)
			fmt.Printf("Campaign ID: %s\n", resp.Msg.Coupon.CampaignId)
			fmt.Printf("Issued At: %s\n", resp.Msg.Coupon.IssuedAt.AsTime().Format(time.RFC3339))
		} else {
			fmt.Printf("Failed to issue coupon: %s\n", resp.Msg.Error)
		}

	case "delete":
		// Create request
		req := connect.NewRequest(&coupon.DeleteCampaignRequest{
			CampaignId:   *campaignID,
			CampaignName: *campaignName,
		})

		// Call API
		resp, err := client.DeleteCampaign(context.Background(), req)
		if err != nil {
			log.Fatalf("Error deleting campaign: %v", err)
		}

		// Print result
		if resp.Msg.Success {
			fmt.Printf("Success: %s\n", resp.Msg.Message)
		} else {
			fmt.Printf("Failed: %s\n", resp.Msg.Message)
		}

	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: create, get, issue, delete")
		os.Exit(1)
	}
}
