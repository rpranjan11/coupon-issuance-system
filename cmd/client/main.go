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
)

func main() {
	// Parse command-line flags
	serverAddr := flag.String("server", "http://localhost:8080", "server address")
	command := flag.String("command", "issue", "command to run: create, get, or issue")
	campaignID := flag.String("campaign", "", "campaign ID for get and issue commands")
	campaignName := flag.String("name", "Test Campaign", "campaign name for create command")
	totalCoupons := flag.Int("total", 10, "total coupons for create command")
	startIn := flag.Duration("start-in", 0, "start time in duration from now for create command")
	flag.Parse()

	// Create HTTP client
	client := coupon.NewCouponServiceClient(
		http.DefaultClient,
		*serverAddr,
	)

	// Execute the requested command
	switch *command {
	case "create":