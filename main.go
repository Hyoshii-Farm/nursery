package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hyoshii-Farm/nursery/config"
	"github.com/Hyoshii-Farm/nursery/feature"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup and start server
	app := fiber.New(fiber.Config{
		Prefork:        os.Getenv("ENVIRONMENT") == "production",
		ServerHeader:   "Fiber",
		ReadBufferSize: 16384,
	})
	app.Use(cors.New())

	// Apply Auth middleware to all /api/v2 routes
	api := app.Group("/api/v2")
	feature.RegisterAll(api, db)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3022"
	}

	// Start server in goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Server starting on port %s", port)

	// Handle graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
