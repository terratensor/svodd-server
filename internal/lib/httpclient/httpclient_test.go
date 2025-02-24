package httpclient

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew_UserAgentProvided(t *testing.T) {
	userAgent := "testUserAgent"
	client := New(&userAgent)

	if client.UserAgent != userAgent {
		t.Errorf("Expected user agent: %s, got: %s", userAgent, client.UserAgent)
	}

	if client.Client != http.DefaultClient {
		t.Error("Expected default http client")
	}
}

func TestNew_NoUserAgentProvided(t *testing.T) {
	client := New(nil)

	if client.UserAgent != defaultUserAgent {
		t.Errorf("Expected default user agent: %s, got: %s", defaultUserAgent, client.UserAgent)
	}

	if client.Client != http.DefaultClient {
		t.Error("Expected default http client")
	}
}

func TestNew_HttpClientInitializedCorrectly(t *testing.T) {
	userAgent := "testUserAgent"
	client := New(&userAgent)

	if client.UserAgent != userAgent {
		t.Errorf("Expected user agent: %s, got: %s", userAgent, client.UserAgent)
	}

	if client.Client != http.DefaultClient {
		t.Error("Expected default http client")
	}
}

func TestHttpClient_Get(t *testing.T) {
	// Test case: Successful request
	t.Run("Successful request", func(t *testing.T) {
		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Test response"))
		}))
		defer server.Close()

		// Create a new HttpClient with the mock server URL
		client := &HttpClient{
			Client:    &http.Client{},
			UserAgent: "Test User Agent",
		}

		// Parse the mock server URL
		u, err := url.Parse(server.URL)
		if err != nil {
			t.Fatalf("Failed to parse URL: %v", err)
		}

		// Make a GET request to the mock server
		response, err := client.Get(u)
		if err != nil {
			t.Fatalf("Failed to make GET request: %v", err)
		}

		// Check if the response body is correct
		expectedResponse := []byte("Test response")
		if !bytes.Equal(response, expectedResponse) {
			t.Errorf("Expected response: %s, but got: %s", expectedResponse, response)
		}
	})

	// Test case: Request fails
	t.Run("Request fails", func(t *testing.T) {
		// Create a new HttpClient with an invalid URL
		client := &HttpClient{
			Client:    &http.Client{},
			UserAgent: "Test User Agent",
		}

		// Parse an invalid URL
		u, err := url.Parse("http://invalid-url_com")
		if err != nil {
			t.Fatalf("Failed to parse URL: %v", err)
		}

		// Make a GET request to the invalid URL
		_, err = client.Get(u)
		if err == nil {
			t.Errorf("Expected an error, but got nil")
		}
	})

	// Test case: Response status code is not in the 200-299 range
	t.Run("Response status code is not in the 200-299 range", func(t *testing.T) {
		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		// Create a new HttpClient with the mock server URL
		client := &HttpClient{
			Client:    &http.Client{},
			UserAgent: "Test User Agent",
		}

		// Parse the mock server URL
		u, err := url.Parse(server.URL)
		if err != nil {
			t.Fatalf("Failed to parse URL: %v", err)
		}

		// Make a GET request to the mock server
		_, err = client.Get(u)
		if err == nil {
			t.Errorf("Expected an error, but got nil")
		}
	})
}
