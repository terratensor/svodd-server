package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideo"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideopage"
)

func main() {
	cfg := config.MustLoad()

	wg := &sync.WaitGroup{}
	for _, parserCfg := range cfg.Parsers {

		wg.Add(1)
		parser := qavideo.NewParser(parserCfg, *cfg.Delay, *cfg.RandomDelay)
		go Run(parser, wg)
	}
	wg.Wait()
	log.Println("finished, all workers successfully stopped.")
}

func Run(p *qavideo.Parser, wg *sync.WaitGroup) {
	log.Printf("*p.Link: %v", *p.Link)
	log.Printf("p.Link: %v", p.Link)
	log.Printf("&p.Link: %v", &p.Link)
	defer wg.Done()
	randomDelay := time.Duration(0)
	if p.RandomDelay != 0 {
		randomDelay = time.Duration(rand.Int63n(int64(p.RandomDelay)))
	}
	time.Sleep(p.Delay + randomDelay)

	log.Printf("started parser for given url: %v", p.Link)
	resBytes, _ := p.Client.Get(p.Link)
	// log.Printf("res: %v, err: %v", res, err)
	page, err := qavideopage.New(resBytes)
	if err != nil {
		log.Printf("failed to parse url %v, %v", p.Link, err)
	}
	log.Printf("link: %+v", page)
	log.Printf("active link: %+v", page.Active())
	log.Printf("last link: %+v", page.Last())
	log.Printf("first link: %+v", page.First())
	log.Printf("prev link: %+v", page.Prev())
	log.Printf("next link: %+v", page.Next())
	log.Printf("first qa link: %+v", page.FirstQALink())
	log.Printf("list of qa links: %+v", page.ListQALinks())
	if err != nil {
		log.Printf("failed to parse url %v, %v", p.Link, err)
	}
	log.Printf("fetched the contents of a given url %v", p.Link)
}
