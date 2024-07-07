package qavideopage

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	// Test case: empty body
	body := []byte{}
	_, err := New(body)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test case: valid body with no links or pagination
	body = []byte(`
		<html><body>
			<div id="answer-list">
				<div class="block"></div>
			</div>
		</body></html>
	`)
	_, err = New(body)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test case: valid body with links
	body = []byte(`
		<html><body>
			<div id="answer-list">
				<div class="block"><a href="/link1">Link 1</a></div>
				<div class="block"><a href="/link2">Link 2</a></div>
			</div>
		</body></html>
	`)
	page, err := New(body)
	if err != nil {
		t.Errorf("Expected no error1, got %v", err)
	}
	if len(page.Links.Links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(page.Links.Links))
	}
	expectedLinks := []*url.URL{
		{Path: "/link1"},
		{Path: "/link2"},
	}
	for i, link := range page.Links.Links {
		if link.Path != expectedLinks[i].Path {
			t.Errorf("Expected link %d to be %v, got %v", i, expectedLinks[i], link)
		}
	}

	// Test case: valid body with pagination
	body = []byte(`
		<html><body>
			<div id="answer-list">
				<div class="block"><a href="/link1">Link 1</a></div>
				<div class="block"><a href="/link2">Link 2</a></div>
				<div class="active"><a href="/page3">Active</a></div>
				<ul class="pagination">
					<li class="first"><a href="/page1">First</a></li>
					<li class="prev"><a href="/page2">Prev</a></li>
					<li class="active"><a href="/page3">Active</a></li>
					<li class="next"><a href="/page4">Next</a></li>
					<li class="last"><a href="/page5">Last</a></li>
				</ul>
			</div>
		</body></html>
	`)
	page, err = New(body)
	if err != nil {
		t.Errorf("Expected no error2, got %v", err)
	}
	expectedFirst := &url.URL{Path: "/page1"}
	expectedPrev := &url.URL{Path: "/page2"}
	expectedActive := &url.URL{Path: "/page3"}
	expectedNext := &url.URL{Path: "/page4"}
	expectedLast := &url.URL{Path: "/page5"}

	if !reflect.DeepEqual(page.Pagination.Active, expectedActive) {
		t.Errorf("Expected Active to be %v, but got %v", expectedActive, page.Pagination.Active)
	}
	if !reflect.DeepEqual(page.Pagination.Next, expectedNext) {
		t.Errorf("Expected Next to be %v, but got %v", expectedNext, page.Pagination.Next)
	}
	if !reflect.DeepEqual(page.Pagination.Prev, expectedPrev) {
		t.Errorf("Expected Prev to be %v, but got %v", expectedPrev, page.Pagination.Prev)
	}
	if !reflect.DeepEqual(page.Pagination.First, expectedFirst) {
		t.Errorf("Expected First to be %v, but got %v", expectedFirst, page.Pagination.First)
	}
	if !reflect.DeepEqual(page.Pagination.Last, expectedLast) {
		t.Errorf("Expected Last to be %v, but got %v", expectedLast, page.Pagination.Last)
	}
}
