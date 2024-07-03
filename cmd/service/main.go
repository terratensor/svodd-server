package main

import (
	"log"

	"github.com/terratensor/svodd-server/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log.Println(cfg)
	log.Println("finished, all workers successfully stopped.")
}
