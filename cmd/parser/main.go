package main

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/app"
	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/lib/httpclient"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideo"
	"github.com/terratensor/svodd-server/internal/qaparser/questionanswer"
	"github.com/terratensor/svodd-server/internal/workerpool"
)

func main() {
	cfg := config.MustLoad()

	app := app.NewApp(cfg)

	ch := make(chan *url.URL, cfg.EntryChanBuffer)

	wg := &sync.WaitGroup{}
	for _, parserCfg := range cfg.Parsers {

		wg.Add(1)
		parser := qavideo.NewParser(parserCfg, *cfg.Delay, *cfg.RandomDelay)
		go parser.Run(ch, wg)
	}

	var allTask []*workerpool.Task

	for page := range ch {
		log.Printf("page: %v", page)
		task := workerpool.NewTask(func(data interface{}) error {
			taskID := data.(*url.URL)
			time.Sleep(100 * time.Millisecond)
			log.Printf("Task %v processed\n", taskID.String())

			entry := questionanswer.NewEntry(taskID, cfg)

			client := httpclient.New(nil)
			err := entry.FetchData(client)
			if err != nil {
				return err
			}

			return app.Process(entry)

		}, page)

		allTask = append(allTask, task)
	}

	pool := workerpool.NewPool(allTask, cfg.Workers)
	pool.Run()

	wg.Wait()
	log.Println("finished, all tasks successfully processed.")
}
