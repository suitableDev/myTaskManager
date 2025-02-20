package controller

import (
	"context"
	"log"
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

// getUsers - Responds with the list of all user as JSON
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			log.Printf("CheckUserType error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		if c.Query("startIndex") != "" {
			startIndex, err = strconv.Atoi(c.Query("startIndex"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startIndex"})
				return
			}
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			log.Printf("MongoDB Aggregate Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
			return
		}

		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Printf("MongoDB Result Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while processing user data"})
			return
		}

		if len(allusers) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No users found"})
			return
		}

		c.JSON(http.StatusOK, allusers[0])
	}
}

// getUser - Responds with a single user as JSON
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserID := c.Param("userid")

		if err := helper.MatchUserTypeToUid(c, UserID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user model.User
		err := userCollection.FindOne(ctx, bson.M{"userid": UserID}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
