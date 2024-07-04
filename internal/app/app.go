package app

import (
	"log"
	"os"

	"github.com/terratensor/svodd-server/internal/entities/answer"
	"github.com/terratensor/svodd-server/internal/storage/manticore"
)

func NewEntriesStorage(index string) *answer.Entries {
	var storage answer.StorageInterface

	manticoreClient, err := manticore.New(index)
	if err != nil {
		log.Printf("failed to initialize manticore client for index %v, %v", index, err)
		os.Exit(1)
	}

	storage = manticoreClient

	return answer.NewAnswerStorage(storage)
}