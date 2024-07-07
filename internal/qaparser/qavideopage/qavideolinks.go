package qavideopage

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Links struct {
	Links []*url.URL
}

// ErrNoLinkFound is returned when no link is found in the QA block
var ErrNoLinkFound = fmt.Errorf("no link found in the QA block")

func (l *Links) Parse(doc *goquery.Document) error {
	linkEls := doc.Find("#answer-list .block a")
	if linkEls.Nodes == nil {
		return ErrNoLinkFound
	}
	for _, e := range linkEls.Nodes {
		l.Links = append(l.Links, &url.URL{Path: e.Attr[0].Val})
	}
	return nil
}