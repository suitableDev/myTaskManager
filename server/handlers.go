package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

// getTasks - responds with the list of all tasks as JSON.
func getTasks(ctx *gin.Context) {
	cursor, err := mongoClient.Database("task_manager").Collection("tasks").Find(context.TODO(), bson.D{{}})
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "error fetching tasks", err.Error())
		return
	}

	var tasks []bson.M
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "error decoding tasks", err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, tasks)
}

// getTaskByID returns task with an ID value matches the id parameter sent by the client
func getTaskByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	result := mongoClient.Database("task_manager").Collection("tasks").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}})

	var task bson.M
	if err := result.Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(ctx, http.StatusNotFound, "Task not found", err.Error())
		} else {
			respondWithError(ctx, http.StatusInternalServerError, "Error retrieving task", err.Error())
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, task)
}

// postTask adds a task from JSON received in the request body.
func postTask(ctx *gin.Context) {
	var newTask Task
	newTask.ID = primitive.NewObjectID()
	newTask.Created = time.Now().UTC()
	newTask.Status = false

	if err := ctx.BindJSON(&newTask); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	if newTask.Status {
		respondWithError(ctx, http.StatusBadRequest, "Task cannot be created with true status", "")
		return
	}

	if err := validate.Struct(newTask); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	_, err := mongoClient.Database("task_manager").Collection("tasks").InsertOne(context.TODO(), newTask)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Error inserting task", err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"data":    newTask,
	})
}

// updateTask updates the task with an ID value that matches the id parameter sent by the client
func updateTask(ctx *gin.Context) {
	// id := context.Param("id")

	// var updatedFields struct {
	// 	Title  *string `json:"title" validate:"omitempty,min=1"`
	// 	Status *bool   `json:"status"`
	// }

	// if err := context.BindJSON(&updatedFields); err != nil {
	// 	respondWithError(context, http.StatusBadRequest, "invalid JSON", err.Error())
	// 	return
	// }

	// if err := validate.Struct(updatedFields); err != nil {
	// 	respondWithError(context, http.StatusBadRequest, "validation error", err.Error())
	// 	return
	// }

	// for i, task := range tasks {
	// 	if task.ID.String() == id {

	// 		if updatedFields.Title != nil {
	// 			tasks[i].Title = *updatedFields.Title
	// 		}

	// 		if updatedFields.Status != nil {
	// 			tasks[i].Status = *updatedFields.Status
	// 		}

	// 		context.IndentedJSON(http.StatusOK, tasks[i])
	// 		return
	// 	}
	// }

	// context.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// deleteTask deletes the task with an ID value matches the id parameter sent by the client
func deleteTask(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	result, err := mongoClient.Database("task_manager").Collection("tasks").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Error deleting task", err.Error())
		return
	}

	if result.DeletedCount == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
	}
}
