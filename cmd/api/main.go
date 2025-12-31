package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"twitterx-api/internal/apperror"
	"twitterx-api/internal/logger"
	"twitterx-api/internal/service"
)

type TweetsResponse struct {
	Username string   `json:"username"`
	TweetIDs []string `json:"tweet_ids"`
}

func makeGetUserTweetsHandler(nitterService *service.NitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		logger.Debug("Fetching tweets for user: %s", username)

		// Fetch tweet IDs from Nitter
		tweetIDs, err := nitterService.GetUserTweetIDs(username)
		if err != nil {
			logger.Error("Error fetching tweets for user %s: %v", username, err)
			http.Error(w, err.Error(), apperror.HTTPStatusCode(err))
			return
		}

		logger.Debug("Found %d tweets for user: %s", len(tweetIDs), username)

		response := TweetsResponse{
			Username: username,
			TweetIDs: tweetIDs,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Error encoding response: %v", err)
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

		logger.Debug("Fetching tweet %s for user: %s", tweetID, username)

		// Fetch tweet data from FxTwitter API
		tweetData, err := fxTwitterService.GetTweetData(username, tweetID)
		if err != nil {
			logger.Error("Error fetching tweet %s for user %s: %v", tweetID, username, err)
			http.Error(w, err.Error(), apperror.HTTPStatusCode(err))
			return
		}

		logger.Debug("Successfully fetched tweet %s", tweetID)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tweetData); err != nil {
			logger.Error("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func makeGetUserHandler(fxTwitterService *service.FxTwitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		logger.Debug("Fetching user data for: %s", username)

		// Fetch user data from FxTwitter API
		userData, err := fxTwitterService.GetUserData(username)
		if err != nil {
			logger.Error("Error fetching user %s: %v", username, err)
			http.Error(w, err.Error(), apperror.HTTPStatusCode(err))
			return
		}

		logger.Debug("Successfully fetched user data for: %s", username)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			logger.Error("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("public", "index.html"))
}

func serveProfile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("public", "profile.html"))
}

func main() {
	// Get Nitter URL from environment
	nitterURL := os.Getenv("NITTER_URL")
	if nitterURL == "" {
		logger.Fatal("NITTER_URL environment variable is required")
	}

	// Initialize Nitter service
	nitterService := service.NewNitterService(nitterURL)

	// Initialize FxTwitter service
	fxTwitterService := service.NewFxTwitterService()

	// Setup router
	router := mux.NewRouter()

	// API endpoints
	router.HandleFunc("/api/users/{username}/tweets/{id}", makeGetTweetHandler(fxTwitterService)).Methods("GET")
	router.HandleFunc("/api/users/{username}/tweets", makeGetUserTweetsHandler(nitterService)).Methods("GET")
	router.HandleFunc("/api/users/{username}", makeGetUserHandler(fxTwitterService)).Methods("GET")

	// Static files
	staticFileServer := http.FileServer(http.Dir(filepath.Join("public", "static")))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	// UI routes
	router.HandleFunc("/", serveIndex).Methods("GET")
	router.HandleFunc("/{username}", serveProfile).Methods("GET")

	port := ":8080"
	logger.Info("Server starting on http://127.0.0.1%s", port)
	logger.Info("Using Nitter instance: %s", nitterURL)
	logger.Info("Debug mode: %v", logger.IsDebugEnabled())
	if err := http.ListenAndServe(port, router); err != nil {
		logger.Fatal("Server error: %v", err)
	}
}
