//go:generate swag init -o docs/api

// @title TwitterX API
// @version 1.0
// @description REST API for fetching Twitter user data via Nitter and FxTwitter
// @host localhost:8080
// @BasePath /api

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"twitterx-api/internal/logger"
	"twitterx-api/internal/service"
)

type TweetsResponse struct {
	Username string   `json:"username"`
	TweetIDs []string `json:"tweet_ids"`
}

// GetUserTweets godoc
// @Summary Get user tweet IDs
// @Description Fetch tweet IDs for a Twitter user from Nitter RSS feed
// @Tags tweets
// @Accept json
// @Produce json
// @Param username path string true "Twitter username"
// @Success 200 {object} TweetsResponse
// @Failure 500 {string} string "Failed to fetch tweets"
// @Router /users/{username}/tweets [get]
func makeGetUserTweetsHandler(nitterService *service.NitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		logger.Debug("Fetching tweets for user: %s", username)

		// Fetch tweet IDs from Nitter
		tweetIDs, err := nitterService.GetUserTweetIDs(username)
		if err != nil {
			logger.Error("Error fetching tweets for user %s: %v", username, err)
			http.Error(w, fmt.Sprintf("Failed to fetch tweets: %v", err), http.StatusInternalServerError)
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

// GetTweet godoc
// @Summary Get tweet details
// @Description Fetch detailed tweet information from FxTwitter API
// @Tags tweets
// @Accept json
// @Produce json
// @Param username path string true "Twitter username"
// @Param id path string true "Tweet ID"
// @Success 200 {object} models.FxTwitterResponse
// @Failure 500 {string} string "Failed to fetch tweet"
// @Router /users/{username}/tweets/{id} [get]
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
			http.Error(w, fmt.Sprintf("Failed to fetch tweet: %v", err), http.StatusInternalServerError)
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

// GetUser godoc
// @Summary Get user profile
// @Description Fetch Twitter user profile information from FxTwitter API
// @Tags users
// @Accept json
// @Produce json
// @Param username path string true "Twitter username"
// @Success 200 {object} models.FxTwitterUserResponse
// @Failure 500 {string} string "Failed to fetch user"
// @Router /users/{username} [get]
func makeGetUserHandler(fxTwitterService *service.FxTwitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		logger.Debug("Fetching user data for: %s", username)

		// Fetch user data from FxTwitter API
		userData, err := fxTwitterService.GetUserData(username)
		if err != nil {
			logger.Error("Error fetching user %s: %v", username, err)
			http.Error(w, fmt.Sprintf("Failed to fetch user: %v", err), http.StatusInternalServerError)
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

	// API endpoints (register specific routes first)
	router.HandleFunc("/api/users/{username}/tweets/{id}", makeGetTweetHandler(fxTwitterService)).Methods("GET")
	router.HandleFunc("/api/users/{username}/tweets", makeGetUserTweetsHandler(nitterService)).Methods("GET")
	router.HandleFunc("/api/users/{username}", makeGetUserHandler(fxTwitterService)).Methods("GET")

	// Root redirect to docs
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/docs/", http.StatusMovedPermanently)
	}).Methods("GET")

	// Swagger UI on /api/docs path (uses generated docs from go:generate)
	router.PathPrefix("/api/docs").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
	))

	port := ":8080"
	logger.Info("Server starting on http://localhost%s", port)
	logger.Info("Using Nitter instance: %s", nitterURL)
	logger.Info("Debug mode: %v", logger.IsDebugEnabled())
	logger.Info("Swagger UI available at: http://localhost%s/api/docs/", port)
	logger.Fatal(http.ListenAndServe(port, router).Error())
}
