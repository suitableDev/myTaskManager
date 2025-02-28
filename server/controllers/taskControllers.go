package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	database "task-manager/server/database"
	helper "task-manager/server/helpers"
	model "task-manager/server/models"
)

// Helper function to handle context setup -- keepng it here for ease
func getContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// HealthCheck - Health endpoint, returns a simple status response
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		myData := map[string]string{"version": "1.0"}
		helper.RespondWithSuccess(c, http.StatusOK, "Health check", myData)
	}
}

// GetTasks - Retrieves all tasks
func GetTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, username, valid := helper.GetUserDetails(c)
		if !valid {
			helper.RespondWithError(c, http.StatusUnauthorized, "User not authorized", "UID or Username not found in context")
			return
		}

		ctx, cancel := getContextWithTimeout()
		defer cancel()

		taskCollection := database.GetTaskCollection()
		filter := bson.M{"user_id": userID}

		opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}})

		cursor, err := taskCollection.Find(ctx, filter, opts)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Error fetching tasks", err.Error())
			return
		}
		defer cursor.Close(ctx)

		var tasks []model.Task
		if err = cursor.All(ctx, &tasks); err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Error decoding tasks", err.Error())
			return
		}

		if len(tasks) == 0 {
			helper.RespondWithSuccess(c, http.StatusOK, "No tasks found for "+username, []model.Task{})
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "Tasks for "+username, tasks)
	}
}

// GetTaskByID - Retrieves a single task by its ID
func GetTaskByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, username, valid := helper.GetUserDetails(c)
		if !valid {
			helper.RespondWithError(c, http.StatusUnauthorized, "User not authorized", "UID or Username not found in context")
			return
		}

		taskId := c.Param("id")
		if taskId == "" {
			helper.RespondWithError(c, http.StatusBadRequest, "Task ID is required", "No ID provided in the request")
			return
		}

		objId, err := primitive.ObjectIDFromHex(taskId)
		if err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid task ID format", err.Error())
			return
		}

		ctx, cancel := getContextWithTimeout()
		defer cancel()

		taskCollection := database.GetTaskCollection()
		filter := bson.M{"_id": objId, "user_id": userID}

		var task model.Task
		err = taskCollection.FindOne(ctx, filter).Decode(&task)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				helper.RespondWithError(c, http.StatusNotFound, "Task not found", "No task found for the specified ID and user "+username)
				return
			}
			helper.RespondWithError(c, http.StatusInternalServerError, "Error fetching task", err.Error())
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "Task for "+username, task)
	}
}

// PostTask - Adds a task from JSON received in the request body
func PostTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTask model.Task
		if err := c.ShouldBindJSON(&newTask); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid JSON input", err.Error())
			return
		}

		userID, username, valid := helper.GetUserDetails(c)
		if !valid {
			helper.RespondWithError(c, http.StatusInternalServerError, "Invalid user details", "Failed to get username or UID")
			return
		}

		newTask.ID = primitive.NewObjectID()
		newTask.UserID = userID
		newTask.Username = username
		newTask.Created = time.Now().UTC()
		newTask.Updated = time.Time{}
		newTask.Status = false

		if err := validate.Struct(newTask); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Validation error", err.Error())
			return
		}

		collection := database.GetTaskCollection()
		if _, err := collection.InsertOne(c.Request.Context(), newTask); err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Error inserting task", err.Error())
			return
		}

		helper.RespondWithSuccess(c, http.StatusCreated, "Task created successfully", newTask)
	}
}

// UpdateTask - Updates the task with the specified ID
func UpdateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, username, valid := helper.GetUserDetails(c)
		if !valid {
			helper.RespondWithError(c, http.StatusUnauthorized, "User not authorized", "UID or Username not found in context")
			return
		}

		idStr := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid ID format", err.Error())
			return
		}

		var updatedFields struct {
			Title   *string    `json:"title" validate:"omitempty,min=1"`
			Status  *bool      `json:"status"`
			Updated *time.Time `json:"updated"`
		}

		if err := c.ShouldBindJSON(&updatedFields); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid JSON input", err.Error())
			return
		}

		if err := validate.Struct(updatedFields); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Validation error", err.Error())
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
		filter := bson.M{"_id": id, "user_id": userID}

		result, err := collection.UpdateOne(c.Request.Context(), filter, bson.M{"$set": update})
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Error updating task", err.Error())
			return
		}

		if result.MatchedCount == 0 {
			helper.RespondWithError(c, http.StatusNotFound, "Task not found", "No task found for the specified ID and user "+username)
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "Task updated successfully for "+username, update)
	}
}

// DeleteTask - Deletes the task with the specified ID
func DeleteTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, valid := helper.GetUserDetails(c)
		if !valid {
			helper.RespondWithError(c, http.StatusUnauthorized, "User not authorized", "UID not found in context")
			return
		}

		idStr := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid ID format", err.Error())
			return
		}

		collection := database.GetTaskCollection()
		filter := bson.M{"_id": id, "user_id": userID}
		result, err := collection.DeleteOne(c.Request.Context(), filter)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Error deleting task", err.Error())
			return
		}

		if result.DeletedCount == 0 {
			helper.RespondWithError(c, http.StatusNotFound, "Task not found", "No task found for the specified ID and user")
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "Task deleted successfully", nil)
	}
}

// DeleteAllTasks - Deletes all tasks
func DeleteAllTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, valid := helper.GetUserDetails(c)
		if !valid {
			helper.RespondWithError(c, http.StatusUnauthorized, "User not authorized", "UID not found in context")
			return
		}

		collection := database.GetTaskCollection()
		filter := bson.M{"user_id": userID}
		result, err := collection.DeleteMany(c.Request.Context(), filter)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Error deleting all tasks", err.Error())
			return
		}

		if result.DeletedCount == 0 {
			helper.RespondWithError(c, http.StatusNotFound, "No tasks found", "No tasks found for the user to delete")
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "All tasks deleted successfully", nil)
	}
}
