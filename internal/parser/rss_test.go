package parser

import "testing"

const sampleRSS = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test</title>
    <link>https://example.com</link>
    <description>Test</description>
    <item>
      <title>Tweet 1</title>
      <guid>2006027578998472912</guid>
    </item>
    <item>
      <title>Tweet 2</title>
      <guid>https://nitter.net/user/status/1982148508187500913#m</guid>
    </item>
    <item>
      <title>Tweet 3</title>
      <guid>not-a-tweet</guid>
    </item>
  </channel>
</rss>`

func TestParseRSSAndExtractTweetIDs(t *testing.T) {
	rss, err := ParseRSS([]byte(sampleRSS))
	if err != nil {
		t.Fatalf("ParseRSS error: %v", err)
	}

	ids, err := ExtractTweetIDs(rss)
	if err != nil {
		t.Fatalf("ExtractTweetIDs error: %v", err)
	}

	if len(ids) != 2 {
		t.Fatalf("expected 2 tweet IDs, got %d", len(ids))
	}

	if ids[0] != "2006027578998472912" || ids[1] != "1982148508187500913" {
		t.Fatalf("unexpected IDs: %#v", ids)
	}
}

func TestParseRSSInvalidXML(t *testing.T) {
	_, err := ParseRSS([]byte("<rss><channel>"))
	if err == nil {
		t.Fatal("expected ParseRSS error, got nil")
	}
}

func TestExtractTweetIDsNilRSS(t *testing.T) {
	_, err := ExtractTweetIDs(nil)
	if err == nil {
		t.Fatal("expected error for nil rss, got nil")
	}
}
