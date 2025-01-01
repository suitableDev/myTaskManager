package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
)

// Task represents data about a task
type Task struct {
	ID      ksuid.KSUID `json:"id" validate:"required"`
	Title   string      `json:"title" validate:"required,min=1,max=140"`
	Status  bool        `json:"status"`
	Created time.Time   `json:"created" validate:"required"`
}

var validate = validator.New()

// Tasks slice to seed task data - including dummy data
var tasks = []Task{
	{ID: ksuid.New(), Title: "task 1"},
	{ID: ksuid.New(), Title: "task 2", Status: true},
}

// Utility function for error response
func respondWithError(context *gin.Context, code int, message string, details string) {
	context.JSON(code, gin.H{
		"error":   message,
		"details": details,
	})
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
	newTask.Created = time.Now().UTC()
	newTask.Status = false

	if err := context.BindJSON(&newTask); err != nil {
		respondWithError(context, http.StatusBadRequest, "invalid JSON", err.Error())
		return
	}

	if newTask.Status {
		newTask.Status = false
	}

	if err := validate.Struct(newTask); err != nil {
		respondWithError(context, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	tasks = append(tasks, newTask)
	context.IndentedJSON(http.StatusCreated, newTask)
}

// updateTask updates the task with an ID value that matches the id parameter sent by the client
func updateTask(context *gin.Context) {
	id := context.Param("id")

	var updatedFields struct {
		Title  *string `json:"title" validate:"omitempty,min=1"`
		Status *bool   `json:"status"`
	}

	if err := context.BindJSON(&updatedFields); err != nil {
		respondWithError(context, http.StatusBadRequest, "invalid JSON", err.Error())
		return
	}

	if err := validate.Struct(updatedFields); err != nil {
		respondWithError(context, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	for i, task := range tasks {
		if task.ID.String() == id {

			if updatedFields.Title != nil {
				tasks[i].Title = *updatedFields.Title
			}

			if updatedFields.Status != nil {
				tasks[i].Status = *updatedFields.Status
			}

			context.IndentedJSON(http.StatusOK, tasks[i])
			return
		}
	}

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
