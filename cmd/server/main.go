package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celery-github/go-feature-flags/internal/api"
	"github.com/celery-github/go-feature-flags/internal/flags"
)

func main() {
	var (
		addr     = flag.String("addr", ":8080", "HTTP listen address")
		seedPath = flag.String("seed", "./configs/flags.json", "Path to seed flags JSON (optional)")
	)
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	store := flags.NewInMemoryStore()
	svc := flags.NewService(store)

	// Seed flags if file exists
	if *seedPath != "" {
		if err := svc.LoadFromFile(*seedPath); err != nil {
			// If file missing, don't fail hard; treat as optional
			if !errors.Is(err, os.ErrNotExist) {
				logger.Printf("seed load error: %v", err)
			} else {
				logger.Printf("seed file not found (ok): %s", *seedPath)
			}
		} else {
			logger.Printf("seeded flags from %s", *seedPath)
		}
	}

	router := api.NewRouter(svc, logger)

	srv := &http.Server{
		Addr:              *addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Printf("feature-flag service listening on %s", *addr)
		logger.Printf("try: curl http://localhost%s/healthz", *addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("server error: %v", err)
		}
	}()

	<-done
	logger.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Printf("shutdown error: %v", err)
	} else {
		logger.Println("shutdown complete")
	}

	fmt.Println("bye 👋")
}
