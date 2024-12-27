package main

import (
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

	if err := context.BindJSON(&newTask); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	context.IndentedJSON(http.StatusCreated, newTask)
}

func updateTask(context *gin.Context) {
	id := context.Param("id")
	var updatedTask Task

	if err := context.BindJSON(&updatedTask); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = updatedTask
			context.IndentedJSON(http.StatusOK, updatedTask)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// deleteTask deletes the task with an ID value matches the id parameter sent by the client
func deleteTask(context *gin.Context) {
	id := context.Param("id")

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			context.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted"})
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func main() {
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByID)
	router.POST("/tasks", postTask)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)
	router.Run("localhost:8080")
}
