package qavideo

import (
	"net/url"
	"reflect"
	"testing"
)

func TestFirstLink(t *testing.T) {
	// Declare and initialize the variable p with a new Parser
	p := &Parser{}
	// Test case: HTML body with a valid link
	body := []byte(`
		<html>
			<body>
				<div id="answer-list">
					<div class="block">
						<a href="https://example.com">Link</a>
					</div>
				</div>
			</body>
		</html>
	`)
	expectedURL, _ := url.Parse("https://example.com")
	link, err := p.FirstLink(body)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(link, expectedURL) {
		t.Errorf("Expected %v, got %v", expectedURL, link)
	}

	// Test case: HTML body without a link
	body = []byte(`
		<html>
			<body>
				<div id="answer-list">
					<div class="block">
						No link here
					</div>
				</div>
			</body>
		</html>
	`)
	link, err = p.FirstLink(body)
	if err != ErrNoLinkFound {
		t.Errorf("Expected ErrNoLinkFound, got %v", err)
	}
	if link != nil {
		t.Errorf("Expected nil link, got %v", link)
	}

	// Test case: Invalid HTML body
	link, err = p.FirstLink([]byte("Not valid HTML"))
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if link != nil {
		t.Errorf("Expected nil link, got %v", link)
	}
}