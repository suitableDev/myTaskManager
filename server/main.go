package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
)

// Task represents data about a task
type Task struct {
	ID     ksuid.KSUID `json:"id"`
	Title  string      `json:"title" validate:"required"`
	Status bool        `json:"status"`
}

var validate = validator.New()

// Tasks slice to seed task data - including dummy data
var tasks = []Task{
	{ID: ksuid.New(), Title: "task 1"},
	{ID: ksuid.New(), Title: "task 2", Status: true},
}

// getTasks - responds with the list of all tasks as JSON.
func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, tasks)
}

// getTaskByID returns task with an ID value matches the id parameter sent by the client
func getTaskByID(context *gin.Context) {
	id := context.Param("id")

	for _, a := range tasks {
		if a.ID.String() == id {
			context.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// postTask adds a task from JSON received in the request body.
func postTask(context *gin.Context) {
	var newTask Task
	newTask.ID = ksuid.New()

	if err := context.BindJSON(&newTask); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(newTask); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "validation error", "details": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	context.IndentedJSON(http.StatusCreated, newTask)
}

// updateTask updates the task with an ID value that matches the id parameter sent by the client
func updateTask(context *gin.Context) {
	id := context.Param("id")

	var updatedFields struct {
		Title  *string `json:"title"`  // Pointer to check if field is present
		Status *bool   `json:"status"` // Pointer to check if field is present
	}

	// Parse the JSON payload
	if err := context.BindJSON(&updatedFields); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the task by ID
	for i, task := range tasks {
		if task.ID.String() == id {
			// Update the title if present and non-empty
			if updatedFields.Title != nil {
				if *updatedFields.Title == "" {
					context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
					return
				}
				tasks[i].Title = *updatedFields.Title
			}

			// Update the status if present
			if updatedFields.Status != nil {
				tasks[i].Status = *updatedFields.Status
			}

			// Respond with the updated task
			context.IndentedJSON(http.StatusOK, tasks[i])
			return
		}
	}

	// Task not found
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// deleteTask deletes the task with an ID value matches the id parameter sent by the client
func deleteTask(context *gin.Context) {
	id := context.Param("id")

	for i, task := range tasks {
		if task.ID.String() == id {
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
