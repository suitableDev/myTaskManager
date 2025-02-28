package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/crypto/bcrypt"

	database "task-manager/server/database"
	helper "task-manager/server/helpers"
	model "task-manager/server/models"
)

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err, "HASH")
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword)) == nil
}

// Signup - create a new user
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user model.User

		if err := c.BindJSON(&user); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Validation error", validationErr.Error())
			return
		}

		userCollection := database.GetUserCollection()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Database error", err.Error())
			return
		}

		if count > 0 {
			helper.RespondWithError(c, http.StatusBadRequest, "Email already exists", "")
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now
		user.ID = primitive.NewObjectID()
		user.UserID = new(string)
		*user.UserID = user.ID.Hex()

		verificationToken := uuid.New().String()
		user.VerificationToken = &verificationToken
		user.Verified = false

		token, refreshToken, _, _, err := helper.GenerateAllTokens(*user.Email, *user.Username, *user.UserType, *user.UserID)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Failed to generate tokens", err.Error())
			return
		}

		user.Token = &token
		user.RefreshToken = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "User item was not created", insertErr.Error())
			return
		}

		helper.SendVerificationEmail(*user.Email, verificationToken)

		c.SetCookie("access_token", token, 3600, "/", "", true, true)
		c.SetCookie("refresh_token", refreshToken, 604800, "/", "", true, true)

		helper.RespondWithSuccess(c, http.StatusOK, "User created successfully, please verify your email", nil)
	}
}

// Login - authenticate user and issue tokens
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user model.User
		var foundUser model.User

		if err := c.BindJSON(&user); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		userCollection := database.GetUserCollection()

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			helper.RespondWithError(c, http.StatusUnauthorized, "Email or password incorrect", "")
			return
		}

		if !VerifyPassword(*user.Password, *foundUser.Password) {
			helper.RespondWithError(c, http.StatusUnauthorized, "Email or password incorrect", "")
			return
		}

		token, refreshToken, _, _, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.Username, *foundUser.UserType, *foundUser.UserID)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Failed to generate tokens", err.Error())
			return
		}

		err = helper.UpdateAllTokens(token, refreshToken, *foundUser.UserID)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Failed to update tokens", err.Error())
			return
		}

		c.SetCookie("access_token", token, 3600, "/", "", true, true)
		c.SetCookie("refresh_token", refreshToken, 604800, "/", "", true, true)

		helper.RespondWithSuccess(c, http.StatusOK, "Login successful", nil)
	}
}

// Logout - remove tokens and end session
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("access_token", "", -1, "/", "", true, true)
		c.SetCookie("refresh_token", "", -1, "/", "", true, true)

		helper.RespondWithSuccess(c, http.StatusOK, "Logout successful", nil)
	}
}

// RefreshAccessToken - refreshes access token
func RefreshAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		helper.RefreshToken(c)
	}
}
