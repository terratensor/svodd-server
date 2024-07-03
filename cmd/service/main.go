package main

import (
	"log"
	"sync"

	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/entities/answer"
	"github.com/terratensor/svodd-server/internal/parsers/videoparser"
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

	log.Println("finished, all workers successfully stopped.")
}
