package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpClient struct {
	Client    *http.Client
	UserAgent string
}

// HTTPError represents an HTTP error returned by a server.
type HTTPError struct {
	StatusCode int
	Status     string
}

var defaultUserAgent = "svodd/1.0"

// Error returns a string representation of the HTTP error.
//
// No parameters.
// Returns a string.
func (err HTTPError) Error() string {
	return fmt.Sprintf("http error: %s", err.Status)
}

// New initializes a new HttpClient with the provided user agent.
//
// userAgent: A pointer to a string containing the user agent.
// Returns a pointer to a HttpClient.
func New(userAgent *string) *HttpClient {
	var u string
	if userAgent != nil {
		u = *userAgent
	} else {
		u = defaultUserAgent
	}
	return &HttpClient{
		Client:    http.DefaultClient,
		UserAgent: u,
	}
}

// Get retrieves data from the specified URL using an HTTP GET request.
//
// Parameters:
// - link: A pointer to a url.URL struct representing the URL to send the request to.
// Returns:
// - []byte: The response body as a byte slice.
// - error: An error if the request fails or the response status code is not in the 200-299 range.
func (c *HttpClient) Get(link *url.URL) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, link.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	return io.ReadAll(resp.Body)
}
