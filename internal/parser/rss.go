package parser

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

// RSS represents the root RSS structure
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the RSS channel
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents a single RSS item (tweet)
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Link        string `xml:"link"`
	Creator     string `xml:"creator"`
}

// ParseRSS parses RSS XML data and returns the RSS structure
func ParseRSS(data []byte) (*RSS, error) {
	var rss RSS
	err := xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS: %w", err)
	}
	return &rss, nil
}

// ExtractTweetIDs extracts tweet IDs from RSS items
// GUID can be either a plain ID (e.g., "2006027578998472912") or a URL with /status/ path
func ExtractTweetIDs(rss *RSS) ([]string, error) {
	if rss == nil {
		return nil, fmt.Errorf("rss is nil")
	}

	tweetIDs := []string{}
	// Regex to extract tweet ID from GUID
	// Supports both formats:
	// - Plain numeric ID: "2006027578998472912"
	// - URL format: "/status/1982148508187500913#m"
	statusRe := regexp.MustCompile(`/status/(\d+)`)
	plainRe := regexp.MustCompile(`^(\d+)$`)

	for _, item := range rss.Channel.Items {
		// Try URL format first
		if matches := statusRe.FindStringSubmatch(item.GUID); len(matches) >= 2 {
			tweetIDs = append(tweetIDs, matches[1])
		} else if matches := plainRe.FindStringSubmatch(item.GUID); len(matches) >= 2 {
			// Fall back to plain numeric ID
			tweetIDs = append(tweetIDs, matches[1])
		}
	}

	return tweetIDs, nil
}
