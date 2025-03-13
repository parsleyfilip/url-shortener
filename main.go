package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	collection *mongo.Collection
	ctx        context.Context
)

func main() {
	// Set up MongoDB connection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(getEnv("MONGODB_URI", "mongodb://localhost:27017")))
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	// Check the connection
	ctx = context.Background()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Could not ping MongoDB: ", err)
	}

	// Get a handle for the urls collection
	database := client.Database(getEnv("MONGODB_DATABASE", "url_shortener"))
	collection = database.Collection("urls")

	// Clean up MongoDB connection when the application exits
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Error disconnecting from MongoDB: ", err)
		}
	}()

	// Define server routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If path is /, serve the index page
		if r.URL.Path == "/" && r.Method == http.MethodGet {
			serveIndex(w, r)
			return
		}
		// Otherwise, try to redirect
		redirectHandler(w, r)
	})
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/health", healthCheckHandler)

	// Configure server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server listening on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error during server shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped")
}

// Get environment variable with fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
} 