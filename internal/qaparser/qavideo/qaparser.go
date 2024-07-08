package qavideo

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/lib/httpclient"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideopage"
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
	MaxPages    *int
	Client      *httpclient.HttpClient
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
	if cfg.UserAgent != nil {
		userAgent = *cfg.UserAgent
	}
	maxPages := 1
	if cfg.Pages == nil {
		cfg.Pages = &maxPages
	}
	client := httpclient.New(&userAgent)
	np := Parser{
		Link:        newLink,
		Delay:       delay,
		RandomDelay: randomDelay,
		UserAgent:   userAgent,
		Previous:    cfg.Previous,
		MaxPages:    cfg.Pages,
		Client:      client,
	}
	return &np
}

// Run is a method of the Parser struct that runs the parser.
//
// It takes a channel of URLs and a WaitGroup as parameters.
// It does not return anything.
func (p *Parser) Run(output chan *url.URL, wg *sync.WaitGroup) {
	log.Printf("Starting parser: delay: %v, random delay: %v, url: %v", p.Delay, p.RandomDelay, p.Link)

	defer wg.Done()

	for {
		delayWithRandomness := p.Delay + time.Duration(rand.Int63n(int64(p.RandomDelay)))
		time.Sleep(delayWithRandomness)

		log.Printf("Started parser for given URL: %v", p.Link)

		go func() {
			for page := range qavideopage.FetchAndParsePages(p.Client, *p.Link, *p.MaxPages) {
				for _, entry := range page.ListQALinks() {
					output <- entry
				}
			}
		}()

		log.Printf("Fetched the contents of a given URL: %v", p.Link)

		select {
		case <-context.Background().Done():
			return
		default:
		}
	}
}
