package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/satyam-svg/resume-parser/config"
	"github.com/satyam-svg/resume-parser/internal/routes"
)

func main() {
	// Load .env if available
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	// Load config and initialize database
	config.LoadConfig()
	db := config.InitDB()

	// Register routes and pass DB
	handlerWithMiddleware := routes.RegisterRoutes(db)

	// Setup HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithMiddleware,
	}

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("✅ Server running at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	<-stop
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Graceful shutdown failed: %v", err)
	}

	log.Println("✅ Server shutdown complete")
}
