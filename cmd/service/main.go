package main

import (
	"log"
	"sync"

	"github.com/terratensor/svodd-server/internal/app"
	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/entities/answer"
	"github.com/terratensor/svodd-server/internal/parsers/videoparser"
	"github.com/terratensor/svodd-server/internal/splitter"
	"github.com/terratensor/svodd-server/internal/workerpool"
)

func main() {
	cfg := config.MustLoad()

	ch := make(chan answer.Entry, cfg.EntryChanBuffer)

	wg := &sync.WaitGroup{}
	for _, parserCfg := range cfg.Parsers {

		wg.Add(1)
		parser := videoparser.NewParser(parserCfg, *cfg.Delay, *cfg.RandomDelay)
		go parser.Run(ch, wg)
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
