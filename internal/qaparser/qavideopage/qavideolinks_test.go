package qavideopage

import (
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestLinksParse(t *testing.T) {
	// Test case: No links found
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	links := &Links{}
	err = links.Parse(doc)
	if err != ErrNoLinkFound {
		t.Errorf("Expected ErrNoLinkFound, got %v", err)
	}
	if len(links.Links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links.Links))
	}

	// Test case: Multiple links found
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(`<html><body><div id="answer-list"><div class="block"><a href="/link1">Link 1</a></div><div class="block"><a href="/link2">Link 2</a></div></div></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	links = &Links{}
	err = links.Parse(doc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(links.Links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(links.Links))
	}
	expectedLinks := []*url.URL{
		{Path: "/link1"},
		{Path: "/link2"},
	}
	for i, link := range links.Links {
		if link.Path != expectedLinks[i].Path {
			t.Errorf("Expected link %d to be %v, got %v", i, expectedLinks[i], link)
		}
	}
}
