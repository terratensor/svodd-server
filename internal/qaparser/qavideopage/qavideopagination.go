package qavideopage

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Pagination struct {
	Active *url.URL
	Next   *url.URL
	Prev   *url.URL
	First  *url.URL
	Last   *url.URL
}

var ErrNoNextLinkFound = fmt.Errorf("no next link found in the pagination block")

func (p *Pagination) Parse(doc *goquery.Document) error {
	el := doc.Find("#answer-list .pagination").First()
	p.ActiveLink(el)
	p.NextLink(el)
	p.PrevLink(el)
	p.FirstLink(el)
	p.LastLink(el)
	return nil
}

func (p *Pagination) ActiveLink(sel *goquery.Selection) {
	anchor := sel.Find(".active a").First()
	if href, exists := anchor.Attr("href"); exists {
		p.Active, _ = url.Parse(href)
	} else {
		p.Active = nil
	}
}

func (p *Pagination) NextLink(sel *goquery.Selection) {
	anchor := sel.Find(".next a").First()
	if href, exists := anchor.Attr("href"); exists {
		p.Next, _ = url.Parse(href)
	} else {
		p.Next = nil
	}
}

func (p *Pagination) PrevLink(sel *goquery.Selection) {
	anchor := sel.Find(".prev a").First()
	if href, exists := anchor.Attr("href"); exists {
		p.Prev, _ = url.Parse(href)
	} else {
		p.Prev = nil
	}
}

func (p *Pagination) FirstLink(sel *goquery.Selection) {
	anchor := sel.Find(".first a").First()
	if href, exists := anchor.Attr("href"); exists {
		p.First, _ = url.Parse(href)
	} else {
		p.First = nil
	}
}

func (p *Pagination) LastLink(sel *goquery.Selection) {
	anchor := sel.Find(".last a").First()
	if href, exists := anchor.Attr("href"); exists {
		p.Last, _ = url.Parse(href)
	} else {
		p.Last = nil
	}
}
