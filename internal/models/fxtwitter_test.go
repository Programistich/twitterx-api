package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTwitterTimeUnmarshalRFC1123(t *testing.T) {
	var tt TwitterTime
	input := `"Mon Jan 02 15:04:05 -0700 2006"`
	if err := json.Unmarshal([]byte(input), &tt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	expected := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.FixedZone("-0700", -7*60*60))
	if !tt.Time.Equal(expected) {
		t.Fatalf("unexpected time: %v", tt.Time)
	}
}

func TestTwitterTimeUnmarshalRFC3339(t *testing.T) {
	var tt TwitterTime
	input := `"2006-01-02T15:04:05Z"`
	if err := json.Unmarshal([]byte(input), &tt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	expected := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	if !tt.Time.Equal(expected) {
		t.Fatalf("unexpected time: %v", tt.Time)
	}
}

func TestTwitterTimeMarshalJSON(t *testing.T) {
	tt := TwitterTime{Time: time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)}
	b, err := json.Marshal(tt)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(b) != `"2006-01-02T15:04:05Z"` {
		t.Fatalf("unexpected JSON: %s", string(b))
	}
}

func TestFxTwitterMosaicFormatsUnmarshal(t *testing.T) {
	input := []byte(`{
		"code": 200,
		"message": "OK",
		"tweet": {
			"id": "1",
			"url": "https://x.com/example/status/1",
			"text": "hello",
			"author": {
				"id": "1",
				"name": "Example",
				"screen_name": "example",
				"avatar_url": "https://example.com/avatar.jpg",
				"verified": false,
				"blue_badge": false
			},
			"replies": 0,
			"retweets": 0,
			"likes": 0,
			"created_at": "Mon Jan 02 15:04:05 -0700 2006",
			"created_timestamp": 1,
			"possibly_scam": false,
			"possibly_sensitive": false,
			"lang": "en",
			"source": "web",
			"media": {
				"mosaic": {
					"type": "mosaic_photo",
					"formats": {
						"jpeg": "https://mosaic.example/jpeg",
						"webp": "https://mosaic.example/webp"
					}
				}
			}
		}
	}`)

	var resp FxTwitterResponse
	if err := json.Unmarshal(input, &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if resp.Tweet == nil || resp.Tweet.Media == nil || resp.Tweet.Media.Mosaic == nil {
		t.Fatalf("missing mosaic data after unmarshal")
	}

	formats := resp.Tweet.Media.Mosaic.Formats
	if formats["jpeg"] != "https://mosaic.example/jpeg" {
		t.Fatalf("unexpected jpeg format: %s", formats["jpeg"])
	}
	if formats["webp"] != "https://mosaic.example/webp" {
		t.Fatalf("unexpected webp format: %s", formats["webp"])
	}
}
