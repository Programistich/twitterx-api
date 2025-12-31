package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"twitterx-api/internal/apperror"
)

const nitterSampleRSS = `<?xml version="1.0" encoding="UTF-8"?>
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
  </channel>
</rss>`

func TestNitterServiceGetUserTweetIDsValidation(t *testing.T) {
	svc := &NitterService{}
	_, err := svc.GetUserTweetIDs("")
	var vErr *apperror.ValidationError
	if !errors.As(err, &vErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestNitterServiceGetUserTweetIDsNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	svc := &NitterService{baseURL: server.URL, httpClient: server.Client()}
	_, err := svc.GetUserTweetIDs("missing")
	var nfErr *apperror.NotFoundError
	if !errors.As(err, &nfErr) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
}

func TestNitterServiceGetUserTweetIDsUpstreamStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := &NitterService{baseURL: server.URL, httpClient: server.Client()}
	_, err := svc.GetUserTweetIDs("user")
	var upErr *apperror.UpstreamError
	if !errors.As(err, &upErr) {
		t.Fatalf("expected UpstreamError, got %v", err)
	}
}

func TestNitterServiceGetUserTweetIDsParseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte("<rss><channel>"))
	}))
	defer server.Close()

	svc := &NitterService{baseURL: server.URL, httpClient: server.Client()}
	_, err := svc.GetUserTweetIDs("user")
	var upErr *apperror.UpstreamError
	if !errors.As(err, &upErr) {
		t.Fatalf("expected UpstreamError, got %v", err)
	}
}

func TestNitterServiceGetUserTweetIDsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(nitterSampleRSS))
	}))
	defer server.Close()

	svc := &NitterService{baseURL: server.URL, httpClient: server.Client()}
	ids, err := svc.GetUserTweetIDs("user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 tweet IDs, got %d", len(ids))
	}
	if ids[0] != "2006027578998472912" || ids[1] != "1982148508187500913" {
		t.Fatalf("unexpected IDs: %#v", ids)
	}
}
