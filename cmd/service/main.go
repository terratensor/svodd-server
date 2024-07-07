package main

import (
	"context"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/app"
	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/entities/answer"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideo"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideopage"
	"github.com/terratensor/svodd-server/internal/splitter"
	"github.com/terratensor/svodd-server/internal/workerpool"
)

func main() {
	cfg := config.MustLoad()

	ch := make(chan *url.URL, cfg.EntryChanBuffer)

	wg := &sync.WaitGroup{}
	for _, parserCfg := range cfg.Parsers {

		wg.Add(1)
		parser := qavideo.NewParser(parserCfg, *cfg.Delay, *cfg.RandomDelay)
		go Run(parser, ch, wg)
	}

	var allTask []*workerpool.Task
	pool := workerpool.NewPool(allTask, cfg.Workers)
	sp := splitter.NewSplitter(cfg.Splitter.OptChunkSize, cfg.Splitter.MaxChunkSize)

	// Создаем срез клиетнов мантикоры по количеству индексов в конфиге
	var manticoreStorages []answer.Entries
	for _, index := range cfg.ManticoreIndex {
		manticoreStorages = append(manticoreStorages, *app.NewEntriesStorage(index.Name))
	}

	go func() {
		for {
			task := workerpool.NewTask(func(data interface{}) error {
				if cfg.Env != "prod" {
					return nil
				}
				// e := data.(answer.Entry)
				return nil
			}, <-ch, sp, &manticoreStorages)
			pool.AddTask(task)
		}
	}()

	pool.RunBackground()

	wg.Wait()
	log.Println("finished, all workers successfully stopped.")
}

func Run(p *qavideo.Parser, ch chan *url.URL, wg *sync.WaitGroup) {
	log.Printf("🚩 run parser: delay: %v, random delay: %v, url: %v", p.Delay, p.RandomDelay, p.Link)

	defer wg.Done()
loop:
	for {

		randomDelay := time.Duration(0)
		if p.RandomDelay != 0 {
			randomDelay = time.Duration(rand.Int63n(int64(p.RandomDelay)))
		}
		time.Sleep(p.Delay + randomDelay)

		log.Printf("started parser for given url: %v", p.Link)
		chin := ProcessPages(p, p.Link, ch)

		go func() {
			// defer close(chout)
			for {
				select {
				case <-context.Background().Done():
					return
				case page, ok := <-chin:
					if !ok {
						return
					}
					for _, entry := range page.ListQALinks() {
						ch <- entry
					}
					// chout <- l
				}
			}
		}()

		log.Printf("fetched the contents of a given url %v", p.Link)

		select {
		case <-context.Background().Done():
			break loop
		default:
		}
	}
}

func ProcessPages(p *qavideo.Parser, link *url.URL, ch chan *url.URL) chan *qavideopage.Page {

	// var page qavideopage.Page
	pch := make(chan *qavideopage.Page, 5)
	go func() {
		defer close(pch)
		log.Println(*p.FollowPages, "follow pages")
		for i := 0; i < *p.FollowPages; i++ {
			resBytes, _ := p.Request(link)
			page, err := qavideopage.New(resBytes)
			if err != nil {
				log.Printf("failed to parse url %v, %v", link, err)
			}
			link.RawQuery = page.Next().RawQuery
			log.Printf("next link: %v", link)
			pch <- page
		}
	}()

	return pch
}
