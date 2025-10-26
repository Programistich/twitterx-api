package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"twitter-api/internal/service"
)

type TweetsResponse struct {
	Username string   `json:"username"`
	TweetIDs []string `json:"tweet_ids"`
}

func makeGetUserTweetsHandler(nitterService *service.NitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		// Fetch tweet IDs from Nitter
		tweetIDs, err := nitterService.GetUserTweetIDs(username)
		if err != nil {
			log.Printf("Error fetching tweets for user %s: %v", username, err)
			http.Error(w, fmt.Sprintf("Failed to fetch tweets: %v", err), http.StatusInternalServerError)
			return
		}

		response := TweetsResponse{
			Username: username,
			TweetIDs: tweetIDs,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func makeGetTweetHandler(fxTwitterService *service.FxTwitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		tweetID := vars["id"]

		// Fetch tweet data from FxTwitter API
		tweetData, err := fxTwitterService.GetTweetData(username, tweetID)
		if err != nil {
			log.Printf("Error fetching tweet %s for user %s: %v", tweetID, username, err)
			http.Error(w, fmt.Sprintf("Failed to fetch tweet: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tweetData); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func makeGetUserHandler(fxTwitterService *service.FxTwitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		// Fetch user data from FxTwitter API
		userData, err := fxTwitterService.GetUserData(username)
		if err != nil {
			log.Printf("Error fetching user %s: %v", username, err)
			http.Error(w, fmt.Sprintf("Failed to fetch user: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get Nitter URL from environment
	nitterURL := os.Getenv("NITTER_URL")
	if nitterURL == "" {
		log.Fatal("NITTER_URL environment variable is required")
	}

	// Initialize Nitter service
	nitterService := service.NewNitterService(nitterURL)

	// Initialize FxTwitter service
	fxTwitterService := service.NewFxTwitterService()

	// Setup router
	router := mux.NewRouter()

	// API endpoints (register specific routes first)
	router.HandleFunc("/users/{username}/tweets/{id}", makeGetTweetHandler(fxTwitterService)).Methods("GET")
	router.HandleFunc("/users/{username}/tweets", makeGetUserTweetsHandler(nitterService)).Methods("GET")
	router.HandleFunc("/users/{username}", makeGetUserHandler(fxTwitterService)).Methods("GET")

	// OpenAPI spec endpoint
	router.HandleFunc("/openapi.yaml", serveOpenAPISpec).Methods("GET")

	// Swagger UI (register PathPrefix last to avoid conflicts)
	router.PathPrefix("/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/openapi.yaml"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
	))

	port := "127.0.0.1:8080"
	fmt.Printf("Server starting on http://%s\n", port)
	fmt.Printf("Using Nitter instance: %s\n", nitterURL)
	fmt.Printf("Swagger UI available at: http://%s/\n", port)
	fmt.Printf("OpenAPI spec available at: http://%s/openapi.yaml\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func serveOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	specPath := filepath.Join("docs", "openapi", "openapi.yaml")

	// Check if file exists
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		log.Printf("OpenAPI spec file not found at: %s", specPath)
		http.Error(w, "OpenAPI specification not found", http.StatusNotFound)
		return
	}

	// Set proper content type for YAML
	w.Header().Set("Content-Type", "application/x-yaml")
	http.ServeFile(w, r, specPath)
}
