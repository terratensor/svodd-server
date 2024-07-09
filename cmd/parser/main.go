package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/app"
	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/entities/answer"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideo"
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
		go parser.Run(ch, wg)
	}

	var allTask []*workerpool.Task
	sp := splitter.NewSplitter(cfg.Splitter.OptChunkSize, cfg.Splitter.MaxChunkSize)

	// Создаем срез клиетнов мантикоры по количеству индексов в конфиге
	var manticoreStorages []answer.Entries
	for _, index := range cfg.ManticoreIndex {
		manticoreStorages = append(manticoreStorages, *app.NewEntriesStorage(index.Name))
	}

	for page := range ch {
		log.Printf("page: %v", page)
		task := workerpool.NewTask(func(data interface{}) error {
			taskID := data.(*url.URL)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Task %v processed\n", taskID.String())
			return nil
		}, page, sp, &manticoreStorages)
		allTask = append(allTask, task)
	}

	pool := workerpool.NewPool(allTask, cfg.Workers)
	pool.Run()

	wg.Wait()
	log.Println("finished, all tasks successfully processed.")
}
