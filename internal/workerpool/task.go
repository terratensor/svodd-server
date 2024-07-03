package workerpool

import (
	"log"
	"os"

	"github.com/terratensor/svodd-server/internal/entities/answer"
)

/**
Task содержит все необходимое для обработки задачи.
Мы передаем ей Data и функцию f, которая должна быть выполнена, с помощью функции process.
Функция f принимает Data в качестве параметра для обработки, а также храним возвращаемую ошибку
*/

type Task struct {
	Err error
	//Entries *feed.Entries
	Data           *answer.Entry
	f              func(interface{}) error
	// Splitter       splitter.Splitter
	EntriesStorage *answer.Entries
}

func NewTaskStorage() *answer.Entries {
	var storage answer.StorageInterface

	manticoreClient, err := manticore.New("feed")
	if err != nil {
		log.Printf("failed to initialize manticore client, %v", err)
		os.Exit(1)
	}

	storage = manticoreClient

	return answer.NewAnswerStorage(storage)
}