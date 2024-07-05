package workerpool

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/terratensor/svodd-server/internal/entities/answer"
	"github.com/terratensor/svodd-server/internal/qaparser"
	"github.com/terratensor/svodd-server/internal/splitter"
)

/**
Task содержит все необходимое для обработки задачи.
Мы передаем ей Data и функцию f, которая должна быть выполнена, с помощью функции process.
Функция f принимает Data в качестве параметра для обработки, а также храним возвращаемую ошибку
*/

type Task struct {
	Err               error
	Data              *qaparser.Entry
	f                 func(interface{}) error
	Splitter          splitter.Splitter
	ManticoreStorages *[]answer.Entries
	PsqlStorage       *answer.Entries
}

func NewTask(f func(interface{}) error, data qaparser.Entry, splitter *splitter.Splitter, storages *[]answer.Entries) *Task {
	return &Task{
		f:                 f,
		Data:              &data,
		Splitter:          *splitter,
		ManticoreStorages: storages,
	}
}

func process(workerID int, task *Task) {
	fmt.Printf("Worker %d processes task %v\n", workerID, task.Data.Url)

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Println(logger)
	// store := task.EntriesStorage

	task.Err = task.f(task.Data)
}
