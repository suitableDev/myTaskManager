package main

import (
	"fmt"
)

// Task represents data about a task
type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Tasks slice to seed task data
var tasks = []Task{
	{Title: "task 1"},
	{Title: "task 2", Completed: true},
}

func main() {
	fmt.Println(tasks)
}
