package qavideopage

import (
	"bytes"
	"net/url"

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
// It takes in a pointer to an httpclient.HttpClient, a starting URL, and a
// maximum number of pages to fetch.
//
// It returns a channel of type *Page, which will be populated with parsed Page
// structs. The channel is buffered with the maximum number of pages.
func FetchAndParsePages(client *httpclient.HttpClient, startingURL url.URL, maxPages int) chan *Page {
	pageChannel := make(chan *Page, maxPages)

	go func() {
		defer close(pageChannel)

		currentURL := startingURL
		for i := 0; i < maxPages; i++ {
			resBytes, _ := client.Get(&currentURL)
			page, err := New(resBytes)
			if err != nil {
				continue
			}

			currentURL.RawQuery = page.Next().RawQuery
			pageChannel <- page
		}
	}()

	return pageChannel
}

func (p *Page) Active() *url.URL {
	return p.Pagination.Active
}

func (p *Page) Next() *url.URL {
	return p.Pagination.Next
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
