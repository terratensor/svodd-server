package qaparser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/qaparser/qaquestion"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideo"
)

// ErrFeedTypeNotDetected is returned when the detection system can not figure
// out the Feed format
var ErrFeedTypeNotDetected = errors.New("failed to detect feed type")

// HTTPError represents an HTTP error returned by a server.
type HTTPError struct {
	StatusCode int
	Status     string
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("http error: %s", err.Status)
}

type Parser struct {
	Link                 *url.URL
	Delay                time.Duration
	RandomDelay          time.Duration
	UserAgent            string
	Previous             bool
	FollowPages          *int
	Client               *http.Client
	QAVideoTranslator    Translator
	QAQuestionTranslator Translator
	qavideo              *qavideo.Parser
	qaquestion           *qaquestion.Parser
}

// NewParser creates a new Parser with the given URL, delay, and randomDelay.
func NewParser(cfg config.Parser, delay time.Duration, randomDelay time.Duration) *Parser {

	newLink, err := url.Parse(cfg.Url)
	if err != nil {
		log.Printf("ERROR: %v, %v", err, cfg.Url)
		return nil
	}
	if cfg.Delay != nil {
		delay = *cfg.Delay
	}
	if cfg.RandomDelay != nil {
		randomDelay = *cfg.RandomDelay
	}
	userAgent := "svodd/1.0"
	if cfg.RandomDelay != nil {
		userAgent = cfg.UserAgent
	}
	np := Parser{
		Link:        newLink,
		Delay:       delay,
		RandomDelay: randomDelay,
		UserAgent:   userAgent,
		Previous:    cfg.Previous,
		FollowPages: cfg.Pages,
		qavideo:     &qavideo.Parser{Link: newLink, FollowPages: cfg.Pages},
		qaquestion:  &qaquestion.Parser{},
	}
	return &np
}

func (p *Parser) Run(ch chan Entry, wg *sync.WaitGroup) {

	log.Printf("🚩 run parser: delay: %v, random delay: %v, url: %v", p.Delay, p.RandomDelay, p.Link.String())

	defer wg.Done()
	// TODO: implement
loop:
	for {

		randomDelay := time.Duration(0)
		if p.RandomDelay != 0 {
			randomDelay = time.Duration(rand.Int63n(int64(p.RandomDelay)))
		}
		time.Sleep(p.Delay + randomDelay)

		log.Printf("started parser for given url: %v", p.Link)
		entries, err := p.ParseURL(p.Link)
		if err != nil {
			log.Printf("failed to parse url %v, %v", p.Link, err)
			continue

		}
		log.Printf("fetched the contents of a given url %v", p.Link)

		select {
		case <-context.Background().Done():
			break loop
		default:
		}

		for _, entry := range *entries {
			ch <- entry
		}
	}
}

// func GetResBody(link *url.URL) ([]byte, error) {

// }

func (p *Parser) Parse(r io.Reader) (entries *[]Entry, err error) {

	feedType := DetectFeedType(p.Link)
	log.Printf("feed type: %v", feedType)

	switch feedType {
	case FeedTypeQA:
		// return p.parseQAFeed(r)
		// link, err := p.qavideo.ParseQAFirst(r)
		// if err != nil {
		// 	return nil, err
		// }
		// log.Printf("link: %v", link)
		// links, _ := p.qavideo.ParseQAList(r)
		// log.Printf("links: %v", links)
		return nil, nil
	case FeedTypeQAQuestion:
		return p.parseQAQuestionFeed(r)
	}

	return nil, ErrFeedTypeNotDetected
}

func (p *Parser) ParseURL(link *url.URL) (*[]Entry, error) {
	entries, err := p.ParseURLWithContext(link, context.Background())
	if err != nil {
		log.Printf("ERROR: %v, %v", err, link)
		return nil, err
	}
	return entries, nil

}

func (p *Parser) ParseURLWithContext(link *url.URL, ctx context.Context) (entries *[]Entry, err error) {
	client := p.httpClient()

	req, err := http.NewRequestWithContext(ctx, "GET", link.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", p.UserAgent)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	return p.Parse(resp.Body)
}

func (p *Parser) parseQAFeed(feed io.Reader) (*[]Entry, error) {
	qavideo, err := p.qavideo.Parse(feed)
	if err != nil {
		return nil, err
	}

	return p.qavideoTrans().Translate(qavideo)
}

func (p *Parser) parseQAQuestionFeed(feed io.Reader) (*[]Entry, error) {
	qaquestion, err := p.qaquestion.Parse(feed)
	if err != nil {
		return nil, err
	}

	return p.qaquestionTrans().Translate(qaquestion)
}

func (p *Parser) qavideoTrans() Translator {
	if p.QAVideoTranslator != nil {
		return p.QAVideoTranslator
	}
	p.QAVideoTranslator = &DefaultQAVideoTranslator{}
	return p.QAVideoTranslator
}

func (p *Parser) qaquestionTrans() Translator {
	if p.QAQuestionTranslator != nil {
		return p.QAQuestionTranslator
	}
	p.QAQuestionTranslator = &DefaultQAQuestionTranslator{}
	return p.QAQuestionTranslator
}

func (p *Parser) httpClient() *http.Client {
	if p.Client != nil {
		return p.Client
	}
	p.Client = &http.Client{}
	return p.Client
}
