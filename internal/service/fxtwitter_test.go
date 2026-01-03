package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"twitterx-api/internal/apperror"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestFxTwitterServiceGetTweetDataValidation(t *testing.T) {
	svc := &FxTwitterService{}
	_, err := svc.GetTweetData("", "123")
	var vErr *apperror.ValidationError
	if !errors.As(err, &vErr) {
		t.Fatalf("expected ValidationError for username, got %v", err)
	}

	_, err = svc.GetTweetData("user", "")
	if !errors.As(err, &vErr) {
		t.Fatalf("expected ValidationError for tweetID, got %v", err)
	}
}

func TestFxTwitterServiceGetTweetDataNotFound(t *testing.T) {
	svc := &FxTwitterService{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"code":404,"message":"NOT_FOUND"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})}}

	_, err := svc.GetTweetData("user", "123")
	var nfErr *apperror.NotFoundError
	if !errors.As(err, &nfErr) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
}

func TestFxTwitterServiceGetTweetDataUpstreamError(t *testing.T) {
	svc := &FxTwitterService{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"code":500,"message":"API_FAIL"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})}}

	_, err := svc.GetTweetData("user", "123")
	var upErr *apperror.UpstreamError
	if !errors.As(err, &upErr) {
		t.Fatalf("expected UpstreamError, got %v", err)
	}
}

func TestFxTwitterServiceGetTweetDataSuccess(t *testing.T) {
	svc := &FxTwitterService{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"code":200,"message":"OK","tweet":{"id":"123","text":"hi","author":{"id":"1","name":"A","screen_name":"a","avatar_url":"x","verified":false,"blue_badge":false},"created_at":"Mon Jan 02 15:04:05 -0700 2006","created_timestamp":1,"replies":0,"retweets":0,"likes":0,"possibly_scam":false,"possibly_sensitive":false,"lang":"en","source":"web"}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})}}

	resp, err := svc.GetTweetData("user", "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Tweet == nil || resp.Tweet.ID != "123" {
		t.Fatalf("unexpected tweet response: %#v", resp.Tweet)
	}
}

func TestFxTwitterServiceGetTweetDataWithReply(t *testing.T) {
	svc := &FxTwitterService{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"code":200,"message":"OK","tweet":{"id":"123","text":"hi","author":{"id":"1","name":"A","screen_name":"a","avatar_url":"x","verified":false,"blue_badge":false},"created_at":"Mon Jan 02 15:04:05 -0700 2006","created_timestamp":1,"replies":0,"retweets":0,"likes":0,"possibly_scam":false,"possibly_sensitive":false,"lang":"en","source":"web","replying_to":"otheruser","replying_to_status":"999"}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})}}

	resp, err := svc.GetTweetData("user", "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Tweet == nil || resp.Tweet.ID != "123" {
		t.Fatalf("unexpected tweet response: %#v", resp.Tweet)
	}
	if resp.Tweet.ReplyingTo == nil || *resp.Tweet.ReplyingTo != "otheruser" {
		t.Fatalf("unexpected replying_to: %v", resp.Tweet.ReplyingTo)
	}
	if resp.Tweet.ReplyingToStatus == nil || *resp.Tweet.ReplyingToStatus != "999" {
		t.Fatalf("unexpected replying_to_status: %v", resp.Tweet.ReplyingToStatus)
	}
}

func TestFxTwitterServiceGetUserDataValidation(t *testing.T) {
	svc := &FxTwitterService{}
	_, err := svc.GetUserData("")
	var vErr *apperror.ValidationError
	if !errors.As(err, &vErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestFxTwitterServiceGetUserDataNotFound(t *testing.T) {
	svc := &FxTwitterService{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"code":404,"message":"NOT_FOUND"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})}}

	_, err := svc.GetUserData("missing")
	var nfErr *apperror.NotFoundError
	if !errors.As(err, &nfErr) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
}

func TestFxTwitterServiceGetUserDataUpstreamError(t *testing.T) {
	svc := &FxTwitterService{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"code":500,"message":"API_FAIL"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})}}

	_, err := svc.GetUserData("user")
	var upErr *apperror.UpstreamError
	if !errors.As(err, &upErr) {
		t.Fatalf("expected UpstreamError, got %v", err)
	}
}

func TestFxTwitterServiceGetUserDataSuccess(t *testing.T) {
	svc := &FxTwitterService{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"code":200,"message":"OK","user":{"screen_name":"a","url":"u","id":"1","followers":1,"following":2,"likes":3,"media_count":4,"tweets":5,"name":"A","description":"d","location":"l","banner_url":"b","avatar_url":"a","joined":"j","protected":false}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})}}

	resp, err := svc.GetUserData("user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User == nil || resp.User.ID != "1" {
		t.Fatalf("unexpected user response: %#v", resp.User)
	}
}
