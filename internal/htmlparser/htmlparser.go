package htmlparser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type HtmlParser struct {
	node *html.Node
}

// New parses the HTML from the io.Reader and returns a new HtmlParser instance.
//
// It takes an io.Reader as input.
// Returns a pointer to HtmlParser and an error.
func New(r io.Reader) (*HtmlParser, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %v", err)
	}
	return &HtmlParser{node: doc}, nil
}

// GetElementByID finds the HTML node with the specified ID.
//
// It takes a string ID as a parameter and returns an *html.Node.
func (p *HtmlParser) GetElementByID(id string) *html.Node {
	var f func(*html.Node) *html.Node
	f = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && hasID(n, id) {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if res := f(c); res != nil {
				return res
			}
		}
		return nil
	}
	return f(p.node)
}

// GetElementsByClassName retrieves all elements with a specific class name within the HTML node.
//
// Parameters:
// - class: the class name to search for.
// - n: the HTML node to search within.
// Returns a slice of HTML nodes that have the specified class name.
func GetElementsByClassName(class string, n *html.Node) []*html.Node {

	var f func(*html.Node) []*html.Node
	f = func(n *html.Node) []*html.Node {
		var res []*html.Node
		if n.Type == html.ElementNode && hasClass(n, class) {
			res = append(res, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			res = append(res, f(c)...)
		}
		return res
	}
	return f(n)
}

// GetFirstElementByClassName retrieves the first HTML element with a specific class name within the given HTML node.
//
// Parameters:
// - class: the class name to search for.
// - n: the HTML node to search within.
//
// Returns the first HTML element with the specified class name, or nil if no element is found.
func GetFirstElementByClassName(class string, n *html.Node) *html.Node {
	var f func(*html.Node) *html.Node
	f = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && hasClass(n, class) {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if res := f(c); res != nil {
				return res
			}
		}
		return nil
	}
	return f(n)
}


// GetAttribute retrieves the value of the specified attribute from the given HTML node.
//
// Parameters:
// - a: The name of the attribute to retrieve.
// - n: The HTML node to search within.
//
// Returns:
// - The value of the specified attribute, or an empty string if the attribute is not found.
func GetAttribute(a string, n *html.Node) string {
	attrMap := make(map[string]string)
	for _, attr := range n.Attr {
		attrMap[attr.Key] = attr.Val
	}
	if val, ok := attrMap[a]; ok {
		return val
	}
	return ""
}

// hasID checks if the given html.Node has the specified id attribute.
//
// n: the html.Node to check for the id attribute.
// id: the id string to search for.
// bool: returns true if the id attribute is found, false otherwise.
func hasID(n *html.Node, id string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "id" && attr.Val == id {
			return true
		}
	}
	return false
}

// hasClass checks if the given html.Node has the specified class attribute.
//
// n: the html.Node to check for the class attribute.
// class: the class string to search for.
// bool: returns true if the class attribute is found, false otherwise.
func hasClass(n *html.Node, class string) bool {
    classMap := make(map[string]bool)
    for _, attr := range n.Attr {
        if attr.Key == "class" {
            for _, c := range strings.Split(attr.Val, " ") {
                classMap[c] = true
            }
        }
    }
    _, ok := classMap[class]
    return ok
}
