package service

import (
	"io"
	"net/http"
	"time"

	"twitterx-api/internal/apperror"
	"twitterx-api/internal/logger"
	"twitterx-api/internal/parser"
)

// NitterService handles interactions with Nitter API
type NitterService struct {
	baseURL    string
	httpClient *http.Client
}

// NewNitterService creates a new Nitter service instance
func NewNitterService(baseURL string) *NitterService {
	return &NitterService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUserTweetIDs fetches tweet IDs for a given username from Nitter RSS feed
func (s *NitterService) GetUserTweetIDs(username string) ([]string, error) {
	if username == "" {
		return nil, &apperror.ValidationError{Field: "username", Message: "cannot be empty"}
	}

	// Construct RSS URL
	rssURL := s.baseURL + "/" + username + "/rss"
	logger.Debug("Nitter: fetching RSS from %s", rssURL)

	// Make HTTP request
	resp, err := s.httpClient.Get(rssURL)
	if err != nil {
		logger.Error("Nitter: failed to fetch RSS feed: %v", err)
		return nil, &apperror.UpstreamError{Service: "Nitter", Message: "failed to fetch RSS feed", Err: err}
	}
	defer resp.Body.Close()

	logger.Debug("Nitter: received response with status %d", resp.StatusCode)

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		return nil, &apperror.NotFoundError{Resource: "user", ID: username}
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("Nitter: unexpected status code: %d", resp.StatusCode)
		return nil, &apperror.UpstreamError{Service: "Nitter", StatusCode: resp.StatusCode, Message: "unexpected status code"}
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Nitter: failed to read response body: %v", err)
		return nil, &apperror.UpstreamError{Service: "Nitter", Message: "failed to read response body", Err: err}
	}

	logger.Debug("Nitter: received %d bytes", len(body))

	// Parse RSS
	rss, err := parser.ParseRSS(body)
	if err != nil {
		logger.Error("Nitter: failed to parse RSS: %v", err)
		return nil, &apperror.UpstreamError{Service: "Nitter", Message: "failed to parse RSS", Err: err}
	}

	// Extract tweet IDs
	tweetIDs, err := parser.ExtractTweetIDs(rss)
	if err != nil {
		logger.Error("Nitter: failed to extract tweet IDs: %v", err)
		return nil, &apperror.UpstreamError{Service: "Nitter", Message: "failed to extract tweet IDs", Err: err}
	}

	logger.Debug("Nitter: extracted %d tweet IDs", len(tweetIDs))
	return tweetIDs, nil
}
