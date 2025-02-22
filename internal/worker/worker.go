package worker

import "fmt"

type Worker struct {
	Name         string
	TasksStorage *TasksStorage
}

func (w *Worker) CollectStats() {
	fmt.Println("will collect stats")
}

func (w *Worker) RunTask() {
	fmt.Println("will start or stop task")
}

func (w *Worker) StartTask() {
	fmt.Println("will start task")
}

func (w *Worker) StopTask() {
	fmt.Println("will stop task")
}
