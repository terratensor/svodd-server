package htmlparser

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestGetElementByID(t *testing.T) {
	// Test case 1: Element with matching ID is found
	doc, _ := html.Parse(strings.NewReader(`<html><body><div id="test">Hello</div></body></html>`))
	parser := &HtmlParser{node: doc}
	element := parser.GetElementByID("test")
	if element == nil || element.Data != "div" || element.Attr[0].Val != "test" {
		t.Error("Failed to find element with matching ID")
	}

	// Test case 2: Element with non-matching ID is not found
	doc, _ = html.Parse(strings.NewReader(`<html><body><div id="test">Hello</div></body></html>`))
	parser = &HtmlParser{node: doc}
	element = parser.GetElementByID("non-existent")
	if element != nil {
		t.Error("Found element with non-matching ID")
	}

	// Test case 3: Element with matching ID is nested in other elements
	doc, _ = html.Parse(strings.NewReader(`<html><body><div><div id="test">Hello</div></div></body></html>`))
	parser = &HtmlParser{node: doc}
	element = parser.GetElementByID("test")
	if element == nil || element.Data != "div" || element.Attr[0].Val != "test" {
		t.Error("Failed to find nested element with matching ID")
	}
}

func TestGetElementsByClassName(t *testing.T) {
	// Test case 1: Empty class name
	root := &html.Node{
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: "class", Val: "foo"},
		},
	}
	result := GetElementsByClassName("", root)
	if len(result) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(result))
	}

	// Test case 2: Non-existent class name
	result = GetElementsByClassName("bar", root)
	if len(result) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(result))
	}

	// Test case 3: Exact match
	result = GetElementsByClassName("foo", root)
	if len(result) != 1 || result[0] != root {
		t.Errorf("Expected 1 element, got %d", len(result))
	}

	// Test case 4: Multiple elements
	root.FirstChild = &html.Node{
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: "class", Val: "foo"},
		},
	}
	result = GetElementsByClassName("foo", root)
	if len(result) != 2 || result[0] != root || result[1] != root.FirstChild {
		t.Errorf("Expected 2 elements, got %d", len(result))
	}

	// Test case 5: Deep nesting
	child := root.FirstChild
	child.FirstChild = &html.Node{
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: "class", Val: "foo"},
		},
	}
	result = GetElementsByClassName("foo", root)
	if len(result) != 3 || result[0] != root || result[1] != root.FirstChild || result[2] != child.FirstChild {
		t.Errorf("Expected 3 elements, got %d", len(result))
	}
}

func TestGetFirstElementByClassName(t *testing.T) {
	// Test case 1: Empty HTML node
	node := &html.Node{}
	result := GetFirstElementByClassName("class", node)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	// Test case 2: No matching class name
	node = &html.Node{
		Type: html.ElementNode,
		Attr: []html.Attribute{{Key: "class", Val: "other"}},
	}
	result = GetFirstElementByClassName("class", node)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	// Test case 3: Matching class name
	node = &html.Node{
		Type: html.ElementNode,
		Attr: []html.Attribute{{Key: "class", Val: "class"}},
	}
	result = GetFirstElementByClassName("class", node)
	if result == nil {
		t.Errorf("Expected non-nil, got nil")
	}

	// Test case 4: Nested HTML node with matching class name
	node = &html.Node{
		Type: html.DocumentNode,
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Attr: []html.Attribute{{Key: "class", Val: "class"}},
		},
	}
	result = GetFirstElementByClassName("class", node)
	if result == nil {
		t.Errorf("Expected non-nil, got nil")
	}
}
