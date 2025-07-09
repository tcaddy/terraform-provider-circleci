package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestClient_RateLimit(t *testing.T) {
	var requests int

	// Mock server that will return a 429 on the first call
	// and a 200 on the second call.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++

		if requests == 1 {
			w.Header().Set("x-ratelimit-reset", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New(ts.URL, "v2", "test-token")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/test", ts.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	statusCode, err := c.DoRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	duration := time.Since(start)

	if statusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, statusCode)
	}

	if requests != 2 {
		t.Errorf("expected 2 requests, got %d", requests)
	}

	if duration < 1*time.Second {
		t.Errorf("expected to wait at least 1 second, but waited %s", duration)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := New("https://circleci.com", "api/v2", "test-token")
	if c.baseURL.String() != "https://circleci.com/api/v2/" {
		t.Errorf("expected base URL to be https://circleci.com/api/v2/, got %s", c.baseURL.String())
	}

	u, _ := url.Parse("test")
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		t.Fatal(err)
	}

	if req.URL.String() != "https://circleci.com/api/v2/test" {
		t.Errorf("expected URL to be https://circleci.com/api/v2/test, got %s", req.URL.String())
	}

	if req.Header.Get("Circle-Token") != "test-token" {
		t.Error("missing Circle-Token header")
	}
}
