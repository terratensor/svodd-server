package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/qaparser"
)

func main() {
	cfg := config.MustLoad()

	wg := &sync.WaitGroup{}
	for _, parserCfg := range cfg.Parsers {

		wg.Add(1)
		parser := qaparser.NewParser(parserCfg, *cfg.Delay, *cfg.RandomDelay)
		go Run(parser, wg)
	}
	wg.Wait()
	log.Println("finished, all workers successfully stopped.")
}

func Run(p *qaparser.Parser, wg *sync.WaitGroup) {
	defer wg.Done()
	randomDelay := time.Duration(0)
		if p.RandomDelay != 0 {
			randomDelay = time.Duration(rand.Int63n(int64(p.RandomDelay)))
		}
		time.Sleep(p.Delay + randomDelay)

		log.Printf("started parser for given url: %v", p.Link)
		_, err := p.ParseURL(p.Link)
		if err != nil {
			log.Printf("failed to parse url %v, %v", p.Link, err)
		}
		log.Printf("fetched the contents of a given url %v", p.Link)
}
