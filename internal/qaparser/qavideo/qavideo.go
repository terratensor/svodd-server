package qavideo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"

	"golang.org/x/net/html"

	"github.com/terratensor/svodd-server/internal/entities/answer"
)

// Parser Вопрос — Ответ видео выпуски парсер
type Parser struct {
	Link *url.URL
}

func (p *Parser) Parse(r io.Reader) (*[]answer.Entry, error) {
	var entries []answer.Entry // TODO: implement type

	node, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", p.Link.String(), err)
	}

	var bufInnerHtml bytes.Buffer

	w := io.Writer(&bufInnerHtml)

	// find all <li> elements
	var processAllProduct func(*html.Node)
	processAllProduct = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h2" {
			// process the Product details within each <li> element
			err := html.Render(w, n)
			if err != nil {
				return
			}
			log.Printf("h2: %v \n", bufInnerHtml.String())

		}
		// traverse the child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processAllProduct(c)
		}
	}
	// make a recursive call to your function
	processAllProduct(node)

	return &entries, nil
}
