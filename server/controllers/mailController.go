package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	database "task-manager/server/database"
	helper "task-manager/server/helpers"
	model "task-manager/server/models"
)

// VerifyEmail - Verifies user email
func VerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			helper.RespondWithError(c, http.StatusBadRequest, "Verification token is required", "")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userCollection := database.GetUserCollection()

		filter := bson.M{"verification_token": token}
		update := bson.M{"$set": bson.M{"verified": true, "verification_token": nil}}

		result, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Failed to verify email", err.Error())
			return
		}

		if result.ModifiedCount == 0 {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid verification token", "")
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "Email verified successfully", nil)
	}
}

// ForgotPassword - Sends password reset email
func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		type ForgotPasswordRequest struct {
			Email string `json:"email"`
		}

		var request ForgotPasswordRequest
		if err := c.BindJSON(&request); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userCollection := database.GetUserCollection()

		var user model.User
		err := userCollection.FindOne(ctx, bson.M{"email": request.Email}).Decode(&user)
		if err != nil {
			helper.RespondWithError(c, http.StatusNotFound, "User not found", "")
			return
		}

		resetToken := uuid.New().String()
		expiry := time.Now().Add(time.Hour)

		filter := bson.M{"email": request.Email}
		update := bson.M{"$set": bson.M{
			"reset_password_token":        &resetToken,
			"reset_password_token_expiry": expiry,
		}}

		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Failed to update user", err.Error())
			return
		}

		helper.SendPasswordResetEmail(*user.Email, resetToken)

		helper.RespondWithSuccess(c, http.StatusOK, "Password reset email sent", nil)
	}
}

// ResetPassword - Resets user password
func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		type ResetPasswordRequest struct {
			Token    string `json:"token"`
			Password string `json:"password"`
		}

		var request ResetPasswordRequest
		if err := c.BindJSON(&request); err != nil {
			helper.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userCollection := database.GetUserCollection()

		var user model.User
		err := userCollection.FindOne(ctx, bson.M{"reset_password_token": request.Token}).Decode(&user)
		if err != nil {
			helper.RespondWithError(c, http.StatusNotFound, "Invalid or expired token", "")
			return
		}

		currentTimeUTC := time.Now().UTC()
		expiryTimeUTC := user.ResetPasswordTokenExpiry.UTC()

		if expiryTimeUTC.Before(currentTimeUTC) {
			helper.RespondWithError(c, http.StatusBadRequest, "Token expired", "")
			return
		}

		if request.Password == "" {
			helper.RespondWithError(c, http.StatusBadRequest, "Password cannot be empty", "")
			return
		}

		hashedPassword := HashPassword(request.Password)

		filter := bson.M{"reset_password_token": request.Token}
		update := bson.M{"$set": bson.M{"password": hashedPassword, "reset_password_token": nil, "reset_password_token_expiry": time.Time{}}}

		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			helper.RespondWithError(c, http.StatusInternalServerError, "Failed to reset password", err.Error())
			return
		}

		helper.RespondWithSuccess(c, http.StatusOK, "Password reset successfully", nil)
	}
}
