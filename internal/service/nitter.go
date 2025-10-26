package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"twitter-api/internal/parser"
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
		return nil, fmt.Errorf("username cannot be empty")
	}

	// Construct RSS URL
	rssURL := fmt.Sprintf("%s/%s/rss", s.baseURL, username)

	// Make HTTP request
	resp, err := s.httpClient.Get(rssURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse RSS
	rss, err := parser.ParseRSS(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS: %w", err)
	}

	// Extract tweet IDs
	tweetIDs, err := parser.ExtractTweetIDs(rss)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tweet IDs: %w", err)
	}

	return tweetIDs, nil
}
