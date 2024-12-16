package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Task represents data about a task
type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Tasks slice to seed task data
var tasks = []Task{
	{ID: "1", Title: "task 1"},
	{ID: "2", Title: "task 2", Completed: true},
}

// getTasks - responds with the list of all tasks as JSON.
func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, tasks)
}

// getTaskByID returns task with an ID value matches the id parameter sent by the client
func getTaskByID(context *gin.Context) {
	id := context.Param("id")

	// Loop over the list of tasks, looking for matching ID
	for _, a := range tasks {
		if a.ID == id {
			context.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})

}

// postTask adds an task from JSON received in the request body.
func postTask(context *gin.Context) {
	var newTask Task

	// Call BindJSON to bind the received JSON to newTask.
	if err := context.BindJSON(&newTask); err != nil {
		return
	}

	// Add the new task to the slice.
	tasks = append(tasks, newTask)
	context.IndentedJSON(http.StatusCreated, newTask)
}

func main() {
	fmt.Println(tasks)
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByID)
	router.POST("/tasks", postTask)
	router.Run("localhost:8080")
}
