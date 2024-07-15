package workerpool

import (
	"fmt"
	"net/url"
)

/**
Task содержит все необходимое для обработки задачи.
Мы передаем ей Data и функцию f, которая должна быть выполнена, с помощью функции process.
Функция f принимает Data в качестве параметра для обработки, а также храним возвращаемую ошибку
*/

type Task struct {
	Err  error
	Data *url.URL
	f    func(interface{}) error
}

func NewTask(f func(interface{}) error, data *url.URL) *Task {
	return &Task{
		f:    f,
		Data: data,
	}
}

func process(workerID int, task *Task) {

	fmt.Printf("Worker %d processes task %v\n", workerID, task.Data)
	task.Err = task.f(task.Data)
}
