package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	database "task-manager/server/database"
	helper "task-manager/server/helpers"
	model "task-manager/server/models"
)

var userCollection *mongo.Collection = database.GetUserCollection()

// GetUsers - Responds with the list of all users as JSON
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Unauthorized", err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage := 10
		page := 1

		if rpp := c.Query("recordPerPage"); rpp != "" {
			rppInt, err := strconv.Atoi(rpp)
			if err == nil && rppInt > 0 {
				recordPerPage = rppInt
			}
		}

		if p := c.Query("page"); p != "" {
			pInt, err := strconv.Atoi(p)
			if err == nil && pInt > 0 {
				page = pInt
			}
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Database error", err.Error())
			return
		}
		defer result.Close(ctx)

		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Database error", err.Error())
			return
		}

		if len(allusers) == 0 {
			helper.RespondWithSuccess(c, http.StatusOK, "No users found", []interface{}{}) // Return empty array for consistency
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "Users retrieved successfully", allusers[0])
	}
}

// GetUser - Responds with a single user as JSON
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserID := c.Param("userid")

		if err := helper.MatchUserTypeToUid(c, UserID); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Unauthorized", err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user model.User
		err := userCollection.FindOne(ctx, bson.M{"userid": UserID}).Decode(&user)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "User not found", err.Error())
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "User retrieved successfully", user)
	}
}
