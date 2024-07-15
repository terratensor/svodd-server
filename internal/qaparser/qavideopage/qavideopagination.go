package qavideopage

import (
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

// Parse parses the pagination elements from the provided goquery document.
//
// Parameters:
//   - doc: The goquery Document to parse.
//
// Return type: None.
func (p *Pagination) Parse(doc *goquery.Document) {
	el := doc.Find("#answer-list .pagination").First()
	p.Active = activeLink(el)
	p.Next = nextLink(el)
	p.Prev = prevLink(el)
	p.First = firstLink(el)
	p.Last = lastLink(el)
}

// activeLink returns the URL of the active link in the provided goquery Selection.
//
// Parameters:
//   - sel: The goquery Selection to search for active link.
//
// Return type: *url.URL.
func activeLink(sel *goquery.Selection) *url.URL {
	anchor := sel.Find(".active a").First()
	if href, exists := anchor.Attr("href"); exists {
		activeURL, err := url.Parse(href)
		if err != nil {
			return nil
		}
		return activeURL
	}
	return nil
}

// nextLink returns the URL of the next link from the provided goquery selection.
//
// Parameters:
//   - sel: The goquery Selection to extract the next link URL from.
//
// Return type: *url.URL.
func nextLink(sel *goquery.Selection) *url.URL {
	anchor := sel.Find(".next a").First()
	if href, exists := anchor.Attr("href"); exists {
		nextURL, err := url.Parse(href)
		if err != nil {
			return nil
		}
		return nextURL
	}
	return nil
}

// prevLink returns the URL of the previous link from the provided goquery Selection.
//
// Parameters:
//   - sel: The goquery Selection to extract the previous link URL from.
//
// Return type: *url.URL.
func prevLink(sel *goquery.Selection) *url.URL {
	anchor := sel.Find(".prev a").First()
	if href, exists := anchor.Attr("href"); exists {
		prevURL, err := url.Parse(href)
		if err != nil {
			return nil
		}
		return prevURL
	}
	return nil
}

// firstLink returns the URL of the first link from the provided goquery Selection.
//
// Parameters:
//   - sel: The goquery Selection to extract the first link URL from.
//
// Return type: *url.URL.
func firstLink(sel *goquery.Selection) *url.URL {
	anchor := sel.Find(".first a").First()
	if href, exists := anchor.Attr("href"); exists {
		firstURL, err := url.Parse(href)
		if err != nil {
			return nil
		}
		return firstURL
	}
	return nil
}

// lastLink returns the URL of the last link from the provided goquery Selection.
//
// Parameters:
//   - sel: The goquery Selection to extract the last link URL from.
//
// Return type: *url.URL.
func lastLink(sel *goquery.Selection) *url.URL {
	anchor := sel.Find(".last a").First()
	if href, exists := anchor.Attr("href"); exists {
		lastURL, err := url.Parse(href)
		if err != nil {
			return nil
		}
		return lastURL
	}
	return nil
}
