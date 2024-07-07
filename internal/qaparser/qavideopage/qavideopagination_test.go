package qavideopage

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestPagination_Parse(t *testing.T) {
	// Test case: valid document with pagination links
	t.Run("ValidDocument", func(t *testing.T) {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`
		<html><body>
			<div id="answer-list">
			<div class="block"></div>
			<div class="active"><a href="/page100">Active</a></div>
				<ul class="pagination">
					<li class="first"><a href="/page1">First</a></li>
					<li class="prev"><a href="/page2">Prev</a></li>
					<li class="active"><a href="/page3">3</a></li>
					<li class="next"><a href="/page4">Next</a></li>
					<li class="last"><a href="/page5">Last</a></li>
				</ul>
			</div>
		</body></html>
		`))

		p := &Pagination{}
		err := p.Parse(doc)
		if err != nil {
			t.Errorf("Parse() returned an error: %v", err)
		}

		expectedFirst := &url.URL{Path: "/page1"}
		expectedPrev := &url.URL{Path: "/page2"}
		expectedActive := &url.URL{Path: "/page3"}
		expectedNext := &url.URL{Path: "/page4"}
		expectedLast := &url.URL{Path: "/page5"}

		if !reflect.DeepEqual(p.Active, expectedActive) {
			t.Errorf("Expected Active to be %v, but got %v", expectedActive, p.Active)
		}
		if !reflect.DeepEqual(p.Next, expectedNext) {
			t.Errorf("Expected Next to be %v, but got %v", expectedNext, p.Next)
		}
		if !reflect.DeepEqual(p.Prev, expectedPrev) {
			t.Errorf("Expected Prev to be %v, but got %v", expectedPrev, p.Prev)
		}
		if !reflect.DeepEqual(p.First, expectedFirst) {
			t.Errorf("Expected First to be %v, but got %v", expectedFirst, p.First)
		}
		if !reflect.DeepEqual(p.Last, expectedLast) {
			t.Errorf("Expected Last to be %v, but got %v", expectedLast, p.Last)
		}
	})

	// Test case: document without pagination links
	t.Run("NoPaginationLinks", func(t *testing.T) {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`
			<div id="answer-list">
				<div class="pagination">
				</div>
			</div>
		`))

		p := &Pagination{}
		err := p.Parse(doc)
		if err != nil {
			t.Errorf("Parse() returned an error: %v", err)
		}

		// Check if all pagination links are nil
		if p.Active != nil {
			t.Errorf("Expected Active to be nil, but got %v", p.Active)
		}
		if p.Next != nil {
			t.Errorf("Expected Next to be nil, but got %v", p.Next)
		}
		if p.Prev != nil {
			t.Errorf("Expected Prev to be nil, but got %v", p.Prev)
		}
		if p.First != nil {
			t.Errorf("Expected First to be nil, but got %v", p.First)
		}
		if p.Last != nil {
			t.Errorf("Expected Last to be nil, but got %v", p.Last)
		}
	})
}
