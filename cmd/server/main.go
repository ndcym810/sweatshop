// cmd/server/main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sweatshop/sweatshop/internal/app"
	"github.com/sweatshop/sweatshop/internal/shared/db"
	"github.com/sweatshop/sweatshop/pkg/logger"
)

func main() {
	port := flag.Int("port", 8000, "Server port")
	dataPath := flag.String("data", "./data", "Data directory path")
	flag.Parse()

	// Initialize database
	if err := db.Init(*dataPath); err != nil {
		logger.Error.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	logger.Info.Printf("Database initialized at %s", *dataPath)

	// Create application
	application := app.New()

	// Setup router
	e := application.SetupRouter()

	// Start server
	go func() {
		addr := fmt.Sprintf(":%d", *port)
		logger.Info.Printf("Server starting on %s", addr)
		if err := e.Start(addr); err != nil {
			logger.Error.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info.Println("Shutting down server...")
	if err := e.Close(); err != nil {
		logger.Error.Printf("Error during shutdown: %v", err)
	}
}
