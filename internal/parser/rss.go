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
// GUID format: http://127.0.0.1:8049/elonmusk/status/1982148508187500913#m
func ExtractTweetIDs(rss *RSS) ([]string, error) {
	if rss == nil {
		return nil, fmt.Errorf("rss is nil")
	}

	var tweetIDs []string
	// Regex to extract tweet ID from GUID
	// Pattern: /status/(\d+)
	re := regexp.MustCompile(`/status/(\d+)`)

	for _, item := range rss.Channel.Items {
		matches := re.FindStringSubmatch(item.GUID)
		if len(matches) >= 2 {
			tweetIDs = append(tweetIDs, matches[1])
		}
	}

	return tweetIDs, nil
}
