// cmd/server/main.go
package main

import (
	"context"
	"fmt"

	// Remove unused log import
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Remove unused connect-go import
	"github.com/rs/zerolog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/rpranjan11/coupon-issuance-system/api/coupon/couponconnect"
	"github.com/rpranjan11/coupon-issuance-system/internal/repository/memory"
	"github.com/rpranjan11/coupon-issuance-system/internal/service"
	"github.com/rpranjan11/coupon-issuance-system/internal/service/rpc"
)

const (
	port = 8080
)

func main() {
	// Set up logger
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Create repositories
	campaignRepo := memory.NewCampaignRepository()
	couponRepo := memory.NewCouponRepository()

	// Create service
	campaignService := service.NewCampaignService(campaignRepo, couponRepo)

	// Create RPC server
	couponServer := rpc.NewCouponServiceServer(campaignService)

	// Set up Connect path
	// Change this line to use the correct function from couponconnect
	path, handler := couponconnect.NewCouponServiceHandler(couponServer)

	// Set up middleware for logging and error handling
	var loggingHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Msg("request started")

		handler.ServeHTTP(w, r)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Dur("elapsed", time.Since(startTime)).
			Msg("request completed")
	})

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle(path, loggingHandler)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})

	// Set up HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Start server in a goroutine
	go func() {
		log.Info().Int("port", port).Msg("starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited gracefully")
}
