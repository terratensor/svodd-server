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
	FollowPages *int
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
	client := httpclient.New(&userAgent)
	np := Parser{
		Link:        newLink,
		Delay:       delay,
		RandomDelay: randomDelay,
		UserAgent:   userAgent,
		Previous:    cfg.Previous,
		FollowPages: cfg.Pages,
		Client:      client,
	}
	return &np
}

func (p *Parser) Run(ch chan *url.URL, wg *sync.WaitGroup) {
	log.Printf("🚩 run parser: delay: %v, random delay: %v, url: %v", p.Delay, p.RandomDelay, p.Link)

	defer wg.Done()
loop:
	for {

		randomDelay := time.Duration(0)
		if p.RandomDelay != 0 {
			randomDelay = time.Duration(rand.Int63n(int64(p.RandomDelay)))
		}
		time.Sleep(p.Delay + randomDelay)

		log.Printf("started parser for given url: %v", p.Link)
		// Передаем *p.link чтобы сделать копию и передать значение ссылки,
		// которое будет меняться только внутри функции
		chin := qavideopage.Parse(p.Client, *p.Link, p.FollowPages)

		go func() {
			// defer close(chout)
			for {
				select {
				case <-context.Background().Done():
					return
				case page, ok := <-chin:
					if !ok {
						return
					}
					for _, entry := range page.ListQALinks() {
						ch <- entry
					}
					// chout <- l
				}
			}
		}()

		log.Printf("fetched the contents of a given url %v", p.Link)

		select {
		case <-context.Background().Done():
			break loop
		default:
		}
	}
}
