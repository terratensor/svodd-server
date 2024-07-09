package qavideopage

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/terratensor/svodd-server/internal/lib/httpclient"
)

type Page struct {
	Links      Links
	Pagination Pagination
}

func New(body []byte) (*Page, error) {
	// Load the HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	links := Links{}
	err = links.Parse(doc)
	if err != nil {
		return nil, err
	}

	pagination := Pagination{}
	pagination.Parse(doc)

	return &Page{Links: links, Pagination: pagination}, nil
}

// FetchAndParsePages fetches and parses pages from a given URL using the provided HTTP client.
//
// It takes in a pointer to an httpclient.HttpClient, a starting URL, and a maximum number of pages to fetch.
// It returns a channel of type *Page, which will be populated with parsed Page structs.
func FetchAndParsePages(client *httpclient.HttpClient, startURL url.URL, maxPages int) <-chan *Page {
	pageChan := make(chan *Page, maxPages)
	if maxPages <= 0 {
		maxPages = 1
	}

	go func() {
		defer close(pageChan)

		currentURL := startURL
		for i := 0; i < maxPages; i++ {
			resBytes, _ := client.Get(&currentURL)
			page, _ := New(resBytes)
			pageChan <- page

			nextURL, err := page.Next()
			if err != nil {
				break
			}

			currentURL.RawQuery = nextURL.RawQuery
			if page.Last() == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return pageChan
}

// FetchAndParseAll fetches and parses all pages from a given URL using the provided HTTP client.
//
// It takes in a pointer to an httpclient.HttpClient and a starting URL. It returns a channel of type *Page,
// which will be populated with parsed Page structs.
func FetchAndParseAll(client *httpclient.HttpClient, startURL url.URL) <-chan *Page {
	pageChan := make(chan *Page, 100)

	go func() {
		defer close(pageChan)

		currentURL := startURL
		for {
			resBytes, _ := client.Get(&currentURL)
			page, _ := New(resBytes)
			pageChan <- page

			nextURL, err := page.Next()
			if err != nil {
				break
			}

			currentURL.RawQuery = nextURL.RawQuery
			if page.Last() == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return pageChan
}

func (p *Page) Active() *url.URL {
	return p.Pagination.Active
}

func (p *Page) Next() (*url.URL, error) {
	if condition := p.Pagination.Next == nil; condition {
		return nil, fmt.Errorf("no next link found")
	}
	return p.Pagination.Next, nil
}

func (p *Page) Prev() *url.URL {
	return p.Pagination.Prev
}

func (p *Page) First() *url.URL {
	return p.Pagination.First
}

func (p *Page) Last() *url.URL {
	return p.Pagination.Last
}

func (p *Page) FirstQALink() *url.URL {
	return p.Links.Links[0]
}

func (p *Page) ListQALinks() []*url.URL {
	return p.Links.Links
}
