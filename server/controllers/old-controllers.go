package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"task-manager/server/database"
	"task-manager/server/helpers"
	"task-manager/server/models"
)

// var validate = validator.New()

// getTasks - Responds with the list of all tasks as JSON
func GetTasks(ctx *gin.Context) {
	collection := database.GetTaskCollection()
	cursor, err := collection.Find(ctx.Request.Context(), bson.D{})
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Error fetching tasks", err.Error())
		return
	}
	defer cursor.Close(ctx.Request.Context())

	var tasks []models.Task
	if err = cursor.All(ctx.Request.Context(), &tasks); err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Error decoding tasks", err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, tasks)
}

// getTaskByID - Returns the task with the specified ID
func GetTaskByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var task models.Task
	collection := database.GetTaskCollection()
	err = collection.FindOne(ctx.Request.Context(), bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			helpers.RespondWithError(ctx, http.StatusNotFound, "Task not found", "")
		} else {
			helpers.RespondWithError(ctx, http.StatusInternalServerError, "Error retrieving task", err.Error())
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, task)
}

// postTask - Adds a task from JSON received in the request body
func PostTask(ctx *gin.Context) {
	var newTask models.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid JSON input", err.Error())
		return
	}

	newTask.ID = primitive.NewObjectID()
	newTask.Created = time.Now().UTC()
	newTask.Updated = time.Time{}
	newTask.Status = false // Ensure new tasks are created with `false` status

	if err := validate.Struct(newTask); err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	collection := database.GetTaskCollection()
	_, err := collection.InsertOne(ctx.Request.Context(), newTask)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Error inserting task", err.Error())
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusCreated, "Task created successfully", newTask)
}

// updateTask - Updates the task with the specified ID
func UpdateTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var updatedFields struct {
		Title   *string    `json:"title" validate:"omitempty,min=1"`
		Status  *bool      `json:"status"`
		Updated *time.Time `json:"updated"`
	}

	if err := ctx.ShouldBindJSON(&updatedFields); err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid JSON input", err.Error())
		return
	}

	if err := validate.Struct(updatedFields); err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	update := bson.M{}
	if updatedFields.Title != nil {
		update["title"] = *updatedFields.Title
	}
	if updatedFields.Status != nil {
		update["status"] = *updatedFields.Status
	}
	update["updated"] = time.Now().UTC()

	collection := database.GetTaskCollection()
	result, err := collection.UpdateOne(ctx.Request.Context(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Error updating task", err.Error())
		return
	}

	if result.MatchedCount == 0 {
		helpers.RespondWithError(ctx, http.StatusNotFound, "Task not found", "")
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Task updated successfully", nil)
}

// deleteTask - Deletes the task with the specified ID
func DeleteTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	collection := database.GetTaskCollection()
	result, err := collection.DeleteOne(ctx.Request.Context(), bson.M{"_id": id})
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Error deleting task", err.Error())
		return
	}

	if result.DeletedCount == 0 {
		helpers.RespondWithError(ctx, http.StatusNotFound, "Task not found", "")
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Task deleted successfully", nil)
}

// DeleteAllTasks - Deletes all the tasks
func DeleteAllTasks(ctx *gin.Context) {
	collection := database.GetTaskCollection()
	result, err := collection.DeleteMany(ctx.Request.Context(), bson.D{{}})
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Error deleting all tasks", err.Error())
		return
	}

	if result.DeletedCount == 0 {
		helpers.RespondWithError(ctx, http.StatusNotFound, "No tasks found to delete", "")
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "All tasks deleted successfully", nil)
}
