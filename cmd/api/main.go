package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/api/internal/api"
	"example.com/api/internal/setup"
)

func main() {
	initConfigErr := setup.InitializeConfig()
	if initConfigErr != nil {
		log.Fatalln(initConfigErr)
	}

	api.InitializeRoutes()

	cfg := setup.GetConfig()

	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: cfg.GinEngine,
	}
	errChan := make(chan error)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Printf("Server error: %v\n", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v\n", sig)
	}

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exiting")
}
