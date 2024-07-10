package questionanswer

import (
	"bytes"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/terratensor/svodd-server/internal/lib/httpclient"
)

type Entry struct {
	Url      *url.URL
	Title    string
	Video    *url.URL
	Datetime *time.Time
	Content  []QuestionAnswer
	Comments []Comment
}

type QuestionAnswer struct {
	Text  string
	Chunk int
}

type Comment struct {
	Username   string
	Text       string
	AvatarFile *url.URL
	Role       string
	Datetime   *time.Time
	DataID     string
	Type       int
	Position   int
}

const TypeComment = 3

func NewEntry(url *url.URL) *Entry {
	return &Entry{
		Url: url,
	}
}

func (entry *Entry) FetchData(client *httpclient.HttpClient) error {

	responseBytes, err := client.Get(entry.Url)
	if err != nil {
		return err
	}
	err = entry.Parse(responseBytes)
	if err != nil {
		return err
	}

	return nil
}

func (e *Entry) Parse(resBytes []byte) error {

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resBytes))
	if err != nil {
		return err
	}
	el := doc.Find(".block").First()

	e.Title = el.Find("h1").First().Text()

	e.Video, err = url.Parse(el.Find(".embed-responsive iframe").AttrOr("src", ""))
	if err != nil {
		log.Printf("failed to parse video url: %v", err)
	}

	dateStr := el.Find(".datetime").First().Text()
	dateStr = strings.TrimSpace(dateStr)
	datetime, err := time.Parse("15:04 02.01.2006", dateStr)
	if err != nil {
		log.Printf("failed to parse datetime: %v", err)
	} else {
		e.Datetime = &datetime
	}

	els := doc.Find("#answer-content").First()
	els.Find("p").Each(func(i int, s *goquery.Selection) {
		e.Content = append(e.Content, QuestionAnswer{Text: strings.TrimSpace(s.Text())})
	})

	els = doc.Find(".comment-list").First()
	els.Find(".comment-item").Each(func(i int, s *goquery.Selection) {
		e.Comments = append(e.Comments, Comment{
			Username:   strings.TrimSpace(s.Find(".username").Text()),
			Text:       strings.TrimSpace(s.Find(".comment-text").Text()),
			AvatarFile: parseAvatarFile(s.Find(".ava-80").AttrOr("src", "")),
			Role:       strings.TrimSpace(s.Find(".role").Text()),
			Datetime:   parseDatetime(s.Find(".datetime").Text()),
			DataID:     s.Find(".comment-text").AttrOr("data-id", ""),
			Type:       TypeComment,
			Position:   i + 1,
		})
	})
	log.Printf("entry: %v", e)
	return nil
}

func parseDatetime(datetime string) *time.Time {
	t, err := time.Parse("15:04 02.01.2006", strings.TrimSpace(datetime))
	if err != nil {
		log.Printf("failed to parse datetime: %v", err)
	}
	return &t
}

func parseAvatarFile(avatarFile string) *url.URL {
	u, err := url.Parse(strings.TrimSpace(avatarFile))
	if err != nil {
		log.Printf("failed to parse avatar url: %v", err)
	}
	return u
}
