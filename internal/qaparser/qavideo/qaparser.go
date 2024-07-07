package qavideo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/terratensor/svodd-server/internal/config"
)

type Links []*url.URL

type Page struct {
	Links Links
	Next  *url.URL
	Prev  *url.URL
}

// HTTPError represents an HTTP error returned by a server.
type HTTPError struct {
	StatusCode int
	Status     string
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("http error: %s", err.Status)
}

type Parser struct {
	Link        *url.URL
	Delay       time.Duration
	RandomDelay time.Duration
	UserAgent   string
	Previous    bool
	FollowPages *int
	Client      *http.Client
}

// NewParser creates a new Parser with the given URL, delay, and randomDelay.
func NewParser(cfg config.Parser, delay, randomDelay time.Duration) *Parser {

	newLink, err := url.Parse(cfg.Url)
	if err != nil {
		log.Printf("ERROR: %v, %v", err, cfg.Url)
		return nil
	}
	if cfg.Delay != nil {
		delay = *cfg.Delay
	}
	if cfg.RandomDelay != nil {
		randomDelay = *cfg.RandomDelay
	}
	userAgent := "svodd/1.0"
	if cfg.RandomDelay != nil {
		userAgent = cfg.UserAgent
	}
	np := Parser{
		Link:        newLink,
		Delay:       delay,
		RandomDelay: randomDelay,
		UserAgent:   userAgent,
		Previous:    cfg.Previous,
		FollowPages: cfg.Pages,
	}
	return &np
}

// Request sends an HTTP GET request to the specified URL using the provided Parser instance.
//
// Parameters:
// - link: A pointer to a url.URL struct representing the URL to send the request to.
//
// Returns:
// - []byte: The response body as a byte slice.
// - error: An error if the request fails or the response status code is not in the 200-299 range.
func (p *Parser) Request(link *url.URL) ([]byte, error) {
	client := p.httpClient()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, link.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", p.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	return io.ReadAll(resp.Body)
}

// FirstLink retrieves the first link in the QA block of the given HTML body.
//
// It takes a pointer to a Parser object, a pointer to a URL object, and a byte
// slice representing the HTML body.
//
// It returns a pointer to a URL object and an error.
func (p *Parser) FirstLink(body []byte) (*url.URL, error) {
	// Load the HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	// Find the first link in the QA block
	linkEl := doc.Find("#answer-list .block a").First()
	if linkEl.Nodes == nil {
		return nil, ErrNoLinkFound
	}
	href, _ := linkEl.Attr("href")

	return url.Parse(href)
}

// ErrNoLinkFound is returned when no link is found in the QA block
var ErrNoLinkFound = fmt.Errorf("no link found in the QA block")
var ErrNoNextPageFound = fmt.Errorf("no next page found")

func (p *Parser) httpClient() *http.Client {
	if p.Client != nil {
		return p.Client
	}
	p.Client = &http.Client{}
	return p.Client
}
