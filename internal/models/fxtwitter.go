package models

import (
	"encoding/json"
	"time"
)

// TwitterTime is a custom type for parsing Twitter's time format
type TwitterTime struct {
	time.Time
}

// UnmarshalJSON implements custom unmarshaling for Twitter time format
func (tt *TwitterTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// Twitter API uses RFC1123 format: "Mon Jan 02 15:04:05 -0700 2006"
	t, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", s)
	if err != nil {
		// Try ISO 8601 format as fallback
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}

	tt.Time = t
	return nil
}

// MarshalJSON implements custom marshaling for Twitter time
func (tt TwitterTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(tt.Time)
}

// FxTwitterResponse represents the complete response from FxTwitter API
type FxTwitterResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Tweet   *Tweet `json:"tweet,omitempty"`
}

// Tweet represents a single tweet with all its data
type Tweet struct {
	URL               string       `json:"url"`
	ID                string       `json:"id"`
	Text              string       `json:"text"`
	Author            Author       `json:"author"`
	Replies           int64        `json:"replies"`
	Retweets          int64        `json:"retweets"`
	Likes             int64        `json:"likes"`
	Views             *int64       `json:"views,omitempty"`
	CreatedAt         TwitterTime  `json:"created_at"`
	CreatedTimestamp  int64        `json:"created_timestamp"`
	PossiblyScam      bool         `json:"possibly_scam"`
	PossiblySensitive bool         `json:"possibly_sensitive"`
	Lang              string       `json:"lang"`
	Source            string       `json:"source"`
	Media             *Media       `json:"media,omitempty"`
	Poll              *Poll        `json:"poll,omitempty"`
	Quote             *Tweet       `json:"quote,omitempty"`
	Translation       *Translation `json:"translation,omitempty"`
}

// Author represents the tweet author information
type Author struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ScreenName  string  `json:"screen_name"`
	AvatarURL   string  `json:"avatar_url"`
	AvatarColor *string `json:"avatar_color,omitempty"`
	BannerURL   *string `json:"banner_url,omitempty"`
	Description *string `json:"description,omitempty"`
	Location    *string `json:"location,omitempty"`
	URL         *string `json:"url,omitempty"`
	Followers   *int64  `json:"followers,omitempty"`
	Following   *int64  `json:"following,omitempty"`
	Joined      *string `json:"joined,omitempty"`
	Likes       *int64  `json:"likes,omitempty"`
	Tweets      *int64  `json:"tweets,omitempty"`
	Verified    bool    `json:"verified"`
	BlueBadge   bool    `json:"blue_badge"`
}

// Media represents media attachments in a tweet
type Media struct {
	All      []MediaItem    `json:"all,omitempty"`
	Photos   []Photo        `json:"photos,omitempty"`
	Videos   []Video        `json:"videos,omitempty"`
	Mosaic   *MosaicInfo    `json:"mosaic,omitempty"`
	External *ExternalMedia `json:"external,omitempty"`
}

// MediaItem represents a generic media item
type MediaItem struct {
	Type         string   `json:"type"`
	URL          string   `json:"url"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Format       string   `json:"format,omitempty"`
	ThumbnailURL string   `json:"thumbnail_url,omitempty"`
	Duration     *float64 `json:"duration,omitempty"`
}

// Photo represents a photo attachment
type Photo struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format,omitempty"`
}

// Video represents a video attachment
type Video struct {
	URL          string   `json:"url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Format       string   `json:"format,omitempty"`
	Duration     *float64 `json:"duration,omitempty"`
	Type         string   `json:"type,omitempty"`
}

// MosaicInfo represents mosaic layout information for multiple photos
type MosaicInfo struct {
	Type    string   `json:"type"`
	Formats []string `json:"formats,omitempty"`
}

// ExternalMedia represents external media embeds (e.g., YouTube)
type ExternalMedia struct {
	Type         string   `json:"type"`
	URL          string   `json:"url"`
	ThumbnailURL string   `json:"thumbnail_url,omitempty"`
	Width        int      `json:"width,omitempty"`
	Height       int      `json:"height,omitempty"`
	Duration     *float64 `json:"duration,omitempty"`
}

// Poll represents a poll in a tweet
type Poll struct {
	TotalVotes    int64        `json:"total_votes"`
	EndsAt        TwitterTime  `json:"ends_at"`
	TimeRemaining string       `json:"time_remaining"`
	Choices       []PollChoice `json:"choices"`
}

// PollChoice represents a single poll option
type PollChoice struct {
	Label      string `json:"label"`
	Count      int64  `json:"count"`
	Percentage int    `json:"percentage"`
}

// Translation represents a translated tweet
type Translation struct {
	Text           string `json:"text"`
	SourceLang     string `json:"source_lang"`
	TargetLang     string `json:"target_lang"`
	TranslationURL string `json:"translation_url,omitempty"`
}

// FxTwitterUserResponse represents the complete response from FxTwitter User API
type FxTwitterUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

// User represents Twitter user profile information
type User struct {
	ScreenName   string        `json:"screen_name"`
	URL          string        `json:"url"`
	ID           string        `json:"id"`
	Followers    int64         `json:"followers"`
	Following    int64         `json:"following"`
	Likes        int64         `json:"likes"`
	MediaCount   int64         `json:"media_count"`
	Tweets       int64         `json:"tweets"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Location     string        `json:"location"`
	BannerURL    string        `json:"banner_url"`
	AvatarURL    string        `json:"avatar_url"`
	Joined       string        `json:"joined"`
	Protected    bool          `json:"protected"`
	Website      *string       `json:"website"`
	Verification *Verification `json:"verification,omitempty"`
}

// Verification represents user verification information
type Verification struct {
	Verified bool   `json:"verified"`
	Type     string `json:"type"`
}
