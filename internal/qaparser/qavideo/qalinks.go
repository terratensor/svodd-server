package qavideo

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/terratensor/svodd-server/internal/htmlparser"
	"golang.org/x/net/html"
)

type Links []string

func ParseQACurrent(link string) (*url.URL, error) {
	var u *url.URL
// Request the HTML page.
  res, err := http.Get(link)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Find the QA block items
  doc.Find("#answer-list .block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
	return u, nil
}

// ParseQAFirst parses the first link from the QA list
func (p *Parser) ParseQAFirst(r io.Reader) (*url.URL, error) {
	node, err := htmlparser.New(r)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", p.Link.String(), err)
	}

	n := node.GetElementByID("answer-list")
	if n == nil {
		return nil, fmt.Errorf("answer-list not found")
	}

	firstEl := htmlparser.GetFirstElementByClassName("block", n)
	if firstEl == nil {
		return nil, fmt.Errorf("no blocks found")
	}

	link := p.Link
	link.Path = parseQABlock(firstEl)[0]
	return link, nil
}

func (p *Parser) ParseQAList(r io.Reader) (Links, error) {
	node, err := htmlparser.New(r)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", p.Link.String(), err)
	}

	n := node.GetElementByID("answer-list")
	if n == nil {
		return nil, fmt.Errorf("answer-list not found")
	}

	els := htmlparser.GetElementsByClassName("block", n)
	if len(els) == 0 {
		return nil, fmt.Errorf("no blocks found")
	}

	var links Links
	for _, e := range els {
		links = append(links, parseQABlock(e)...)
	}

	return links, nil
}

func parseQABlock(n *html.Node) []string {
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			links = append(links, htmlparser.GetAttribute("href", n))
		}
		for cl := n.FirstChild; cl != nil; cl = cl.NextSibling {
			f(cl)
		}
	}
	f(n)
	return links
}
