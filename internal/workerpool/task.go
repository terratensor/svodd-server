package workerpool

import (
	"fmt"
	"net/url"

	"github.com/terratensor/svodd-server/internal/entities/answer"
)

/**
Task содержит все необходимое для обработки задачи.
Мы передаем ей Data и функцию f, которая должна быть выполнена, с помощью функции process.
Функция f принимает Data в качестве параметра для обработки, а также храним возвращаемую ошибку
*/

type Task struct {
	Err               error
	Data              *url.URL
	f                 func(interface{}) error
	ManticoreStorages *[]answer.Entries
	PsqlStorage       *answer.Entries
}

func NewTask(f func(interface{}) error, data *url.URL, storages *[]answer.Entries) *Task {
	return &Task{
		f:                 f,
		Data:              data,
		ManticoreStorages: storages,
	}
}

func process(workerID int, task *Task) {

	fmt.Printf("Worker %d processes task %v\n", workerID, task.Data)
	task.Err = task.f(task.Data)
}
