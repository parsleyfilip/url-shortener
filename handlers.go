package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// URL represents the mapping between a shortened URL and its original URL
type URL struct {
	ID         string    `bson:"_id" json:"id"`
	OriginalURL string    `bson:"original_url" json:"original_url"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	Visits     int       `bson:"visits" json:"visits"`
}

// ShortenRequest is the expected JSON payload for shortening a URL
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse is the JSON response for a successful shortening
type ShortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// Shortens URLs and returns the shortened version
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal the request
	var req ShortenRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the URL
	if req.URL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	// Add http:// prefix if missing
	originalURL := req.URL
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	// Generate a short ID for the URL
	id := generateShortID(originalURL)

	// Check if the URL already exists in the database
	var existingURL URL
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingURL)
	
	// If it doesn't exist, create a new entry
	if err == mongo.ErrNoDocuments {
		newURL := URL{
			ID:         id,
			OriginalURL: originalURL,
			CreatedAt:  time.Now(),
			Visits:     0,
		}

		_, err = collection.InsertOne(ctx, newURL)
		if err != nil {
			log.Printf("Error inserting URL: %v", err)
			http.Error(w, "Error creating short URL", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get the host from the request or use a default
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}

	// Create the short URL
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	shortURL := fmt.Sprintf("%s://%s/%s", scheme, host, id)

	// Return the shortened URL
	response := ShortenResponse{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Redirects short URLs to their original destinations
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the short ID from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/")
	
	// Skip empty paths or favicon requests
	if id == "" || id == "favicon.ico" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Look up the URL in the database
	var urlData URL
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&urlData)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Increment visit count
	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"visits": 1}},
	)
	
	if err != nil {
		log.Printf("Error updating visit count: %v", err)
	}

	// Redirect to the original URL
	http.Redirect(w, r, urlData.OriginalURL, http.StatusFound)
}

// Generates a unique short ID for a URL
func generateShortID(url string) string {
	// Create a hash of the URL
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)
	
	// Encode the first 8 bytes as base64
	encoded := base64.URLEncoding.EncodeToString(hash[:8])
	
	// Remove padding and special characters
	return strings.Replace(strings.Replace(encoded, "+", "", -1), "/", "", -1)[:8]
} 