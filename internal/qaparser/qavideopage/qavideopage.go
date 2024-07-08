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

func Parse(client *httpclient.HttpClient, initialLink url.URL, maxPages *int) chan *Page {
	pageCh := make(chan *Page, *maxPages)

	go func() {
		defer close(pageCh)

		for i := 0; i < *maxPages; i++ {
			resBytes, _ := client.Get(&initialLink)
			page, err := New(resBytes)
			if err != nil {
				continue
			}

			initialLink.RawQuery = page.Next().RawQuery
			pageCh <- page
		}
	}()

	return pageCh
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
