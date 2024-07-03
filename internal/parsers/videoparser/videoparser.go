package videoparser

import (
	"context"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/entities/answer"
)

type Parser struct {
	Link        *url.URL
	Delay       time.Duration
	RandomDelay time.Duration
}

type Entry struct {
	Url        string `json:"url"`
	Username   string `json:"username"`
	Text       string `json:"text"`
	AvatarFile string `json:"avatar_file"`
	Role       string `json:"role"`
	Datetime   string `json:"datetime"`
	DataID     string `json:"data_id,omitempty"`
	ParentID   string `json:"parent_id"`
	Type       string `json:"type"`
	Position   int    `json:"position"`
}

const TypeQuestion = "1"
const TypeLinkedQuestion = "2"
const TypeComment = "3"
const TypeAnswer = "4"

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
	np := &Parser{
		Link:        newLink,
		Delay:       delay,
		RandomDelay: randomDelay,
	}
	return np
}

func (p *Parser) Run(ch chan answer.Entry, wg *sync.WaitGroup) {

	log.Printf("🚩 run parser: delay: %v, random delay: %v, url: %v", p.Delay, p.RandomDelay, p.Link.String())
	
	defer wg.Done()
	// TODO: implement

	for {

		randomDelay := time.Duration(0)
		if p.RandomDelay != 0 {
			randomDelay = time.Duration(rand.Int63n(int64(p.RandomDelay)))
		}
		time.Sleep(p.Delay + randomDelay)

		log.Printf("started parser for given url: %v", p.Link)
		// entries := p.getEntries(fp)
		log.Printf("fetched the contents of a given url %v", p.Link)

		select {
		case <-context.Background().Done():
			break
		default:
		}

		// for _, entry := range entries {
		// 	ch <- entry
		// }
	}
}
