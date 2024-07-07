package qavideo

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Parser Вопрос — Ответ видео выпуски парсер
type Parser struct {
	Link        *url.URL
	FollowPages *int
}

func (p *Parser) Parse(r io.Reader) (*Feed, error) {
	var feed Feed // TODO: implement type

	node, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", p.Link.String(), err)
	}

	// Находим блок с id answer-list
	var processAnswerList func(*html.Node)
	processAnswerList = func(n *html.Node) {
		if n.Type == html.ElementNode && nodeHasRequiredID("answer-list", n) {
			parseNode(n)
		}
		// traverse the child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processAnswerList(c)
		}
	}
	// make a recursive call to your function
	processAnswerList(node)

	return &feed, nil
}

func parseNode(n *html.Node) {
	// var bufInnerHtml bytes.Buffer
	// w := io.Writer(&bufInnerHtml)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && nodeHasRequiredCssClass("block", n) {
			// process the Product details within each <li> element
			// err := html.Render(w, n)
			// if err != nil {
			// 	return
			// }
			log.Printf("parce block: %v \n", n.Data)
			parseBlock(n)
		}

		if n.Type == html.ElementNode && nodeHasRequiredCssClass("pagination", n) {
			log.Printf("parce pagination: %v \n", n.Data)
			parsePagination(n)
		}
		for cl := n.FirstChild; cl != nil; cl = cl.NextSibling {
			// log.Printf("cl.Type: %v, cl.Data: %v \n", cl.Type, cl.Data)

			f(cl)
		}
	}
	f(n)
}

func parseBlock(n *html.Node) {
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// process the Product details within each <li> element
			// err := html.Render(w, n)
			// if err != nil {
			// 	return
			// }
			log.Printf("href: %v \n", getRequiredDataAttr("href", n))
			links = append(links, getRequiredDataAttr("href", n))
		}
		for cl := n.FirstChild; cl != nil; cl = cl.NextSibling {
			// log.Printf("cl.Type: %v, cl.Data: %v \n", cl.Type, cl.Data)

			f(cl)
		}
	}
	f(n)
}

func parsePagination(n *html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && nodeHasRequiredCssClass("next", n) {
			if n.FirstChild.Type == html.ElementNode && n.FirstChild.Data == "a" {
				log.Printf("href next page: %v \n", getRequiredDataAttr("href", n.FirstChild))
			}
			// process the Product details within each <li> element
			// err := html.Render(w, n)
			// if err != nil {
			// 	return
			// }
		}
		for cl := n.FirstChild; cl != nil; cl = cl.NextSibling {
			// log.Printf("cl.Type: %v, cl.Data: %v \n", cl.Type, cl.Data)

			f(cl)
		}
	}
	f(n)
}

func getRequiredDataAttr(rda string, n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == rda {
			return attr.Val
		}
	}
	return ""
}

// Перебирает аттрибуты токена в цикле и возвращает bool
// если в html token найден переданный css class
func nodeHasRequiredCssClass(rcc string, n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Split(attr.Val, " ")
			for _, class := range classes {
				if class == rcc {
					return true
				}
			}
		}
	}
	return false
}

func nodeHasRequiredID(id string, n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == "id" {
			if attr.Val == id {
				log.Printf("n.Data %v, id: %v \n", n.Data, attr.Val)
				return true
			}
		}
	}
	return false
}
