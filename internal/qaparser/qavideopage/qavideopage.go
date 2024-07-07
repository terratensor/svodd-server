package qavideopage

import (
	"bytes"
	"net/url"

	"github.com/PuerkitoBio/goquery"
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

func (p * Page) Active() *url.URL {
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