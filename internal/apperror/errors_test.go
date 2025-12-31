package apperror

import (
	"errors"
	"net/http"
	"testing"
)

func TestHTTPStatusCodeMapping(t *testing.T) {
	if code := HTTPStatusCode(&ValidationError{}); code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", code)
	}
	if code := HTTPStatusCode(&NotFoundError{}); code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", code)
	}
	if code := HTTPStatusCode(&UpstreamError{}); code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", code)
	}
	if code := HTTPStatusCode(errors.New("other")); code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", code)
	}
}
