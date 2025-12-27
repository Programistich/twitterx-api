package service

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"twitterx-api/internal/apperror"
	"twitterx-api/internal/logger"
	"twitterx-api/internal/models"
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
		return nil, &apperror.ValidationError{Field: "username", Message: "cannot be empty"}
	}
	if tweetID == "" {
		return nil, &apperror.ValidationError{Field: "tweetID", Message: "cannot be empty"}
	}

	// Construct API URL
	// Format: https://api.fxtwitter.com/{username}/status/{id}
	apiURL := fxTwitterAPIBaseURL + "/" + username + "/status/" + tweetID
	logger.Debug("FxTwitter: fetching tweet from %s", apiURL)

	// Make HTTP request
	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		logger.Error("FxTwitter: failed to fetch tweet data: %v", err)
		return nil, &apperror.UpstreamError{Service: "FxTwitter", Message: "failed to fetch tweet data", Err: err}
	}
	defer resp.Body.Close()

	logger.Debug("FxTwitter: received response with status %d", resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("FxTwitter: failed to read response body: %v", err)
		return nil, &apperror.UpstreamError{Service: "FxTwitter", Message: "failed to read response body", Err: err}
	}

	logger.Debug("FxTwitter: received %d bytes", len(body))

	// Parse JSON response
	var fxResponse models.FxTwitterResponse
	if err := json.Unmarshal(body, &fxResponse); err != nil {
		logger.Error("FxTwitter: failed to parse JSON response: %v", err)
		return nil, &apperror.UpstreamError{Service: "FxTwitter", Message: "failed to parse JSON response", Err: err}
	}

	// Check for API errors (404 = NOT_FOUND, 401 = PRIVATE_TWEET, 500 = API_FAIL)
	if fxResponse.Code == 404 {
		return nil, &apperror.NotFoundError{Resource: "tweet", ID: tweetID}
	}
	if fxResponse.Code != 200 {
		logger.Error("FxTwitter: API error: %s (code: %d)", fxResponse.Message, fxResponse.Code)
		return nil, &apperror.UpstreamError{Service: "FxTwitter", StatusCode: fxResponse.Code, Message: fxResponse.Message}
	}

	logger.Debug("FxTwitter: successfully fetched tweet %s", tweetID)
	return &fxResponse, nil
}

// GetUserData fetches user profile data from FxTwitter API
func (s *FxTwitterService) GetUserData(username string) (*models.FxTwitterUserResponse, error) {
	if username == "" {
		return nil, &apperror.ValidationError{Field: "username", Message: "cannot be empty"}
	}

	// Construct API URL
	// Format: https://api.fxtwitter.com/{username}
	apiURL := fxTwitterAPIBaseURL + "/" + username
	logger.Debug("FxTwitter: fetching user from %s", apiURL)

	// Make HTTP request
	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		logger.Error("FxTwitter: failed to fetch user data: %v", err)
		return nil, &apperror.UpstreamError{Service: "FxTwitter", Message: "failed to fetch user data", Err: err}
	}
	defer resp.Body.Close()

	logger.Debug("FxTwitter: received response with status %d", resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("FxTwitter: failed to read response body: %v", err)
		return nil, &apperror.UpstreamError{Service: "FxTwitter", Message: "failed to read response body", Err: err}
	}

	logger.Debug("FxTwitter: received %d bytes", len(body))

	// Parse JSON response
	var fxUserResponse models.FxTwitterUserResponse
	if err := json.Unmarshal(body, &fxUserResponse); err != nil {
		logger.Error("FxTwitter: failed to parse JSON response: %v", err)
		return nil, &apperror.UpstreamError{Service: "FxTwitter", Message: "failed to parse JSON response", Err: err}
	}

	// Check for API errors (404 = NOT_FOUND, 500 = API_FAIL)
	if fxUserResponse.Code == 404 {
		return nil, &apperror.NotFoundError{Resource: "user", ID: username}
	}
	if fxUserResponse.Code != 200 {
		logger.Error("FxTwitter: API error: %s (code: %d)", fxUserResponse.Message, fxUserResponse.Code)
		return nil, &apperror.UpstreamError{Service: "FxTwitter", StatusCode: fxUserResponse.Code, Message: fxUserResponse.Message}
	}

	logger.Debug("FxTwitter: successfully fetched user %s", username)
	return &fxUserResponse, nil
}
