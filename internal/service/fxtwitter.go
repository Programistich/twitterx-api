package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"twitter-api/internal/models"
)

const fxTwitterAPIBaseURL = "https://api.fxtwitter.com"

// FxTwitterService handles interactions with FxTwitter API
type FxTwitterService struct {
	httpClient *http.Client
}

// NewFxTwitterService creates a new FxTwitter service instance
func NewFxTwitterService() *FxTwitterService {
	return &FxTwitterService{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetTweetData fetches complete tweet data from FxTwitter API
func (s *FxTwitterService) GetTweetData(username, tweetID string) (*models.FxTwitterResponse, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if tweetID == "" {
		return nil, fmt.Errorf("tweet ID cannot be empty")
	}

	// Construct API URL
	// Format: https://api.fxtwitter.com/{username}/status/{id}
	apiURL := fmt.Sprintf("%s/%s/status/%s", fxTwitterAPIBaseURL, username, tweetID)

	// Make HTTP request
	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tweet data: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var fxResponse models.FxTwitterResponse
	if err := json.Unmarshal(body, &fxResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Check for API errors
	if fxResponse.Code != 200 {
		return nil, fmt.Errorf("FxTwitter API error: %s (code: %d)", fxResponse.Message, fxResponse.Code)
	}

	return &fxResponse, nil
}

// GetUserData fetches user profile data from FxTwitter API
func (s *FxTwitterService) GetUserData(username string) (*models.FxTwitterUserResponse, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	// Construct API URL
	// Format: https://api.fxtwitter.com/{username}
	apiURL := fmt.Sprintf("%s/%s", fxTwitterAPIBaseURL, username)

	// Make HTTP request
	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var fxUserResponse models.FxTwitterUserResponse
	if err := json.Unmarshal(body, &fxUserResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Check for API errors
	if fxUserResponse.Code != 200 {
		return nil, fmt.Errorf("FxTwitter API error: %s (code: %d)", fxUserResponse.Message, fxUserResponse.Code)
	}

	return &fxUserResponse, nil
}
