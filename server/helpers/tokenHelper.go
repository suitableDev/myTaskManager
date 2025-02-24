package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	database "task-manager/server/database"
	model "task-manager/server/models"
)

func HashKey() string {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("SECRET_KEY is not set in the .env file")
	}
	return secretKey
}

var hashKey = HashKey()

const (
	AccessTokenExpiry  = 24
	RefreshTokenExpiry = 128
)

// Generate access and refresh tokens
func GenerateAllTokens(email, userName, userType, uid string) (string, string, int64, int64, error) {
	accessExpiry := time.Now().Add(time.Hour * time.Duration(AccessTokenExpiry)).Unix()
	refreshExpiry := time.Now().Add(time.Hour * time.Duration(RefreshTokenExpiry)).Unix()

	claims := &model.SignedDetails{
		Email:    email,
		Username: userName,
		Uid:      uid,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpiry,
		},
	}

	refreshClaims := &model.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpiry,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(hashKey))
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return "", "", 0, 0, err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(hashKey))
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return "", "", 0, 0, err
	}

	return token, refreshToken, accessExpiry, refreshExpiry, nil
}

// Validate token (checks expiration and claims)
func ValidateToken(signedToken string) (claims *model.SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&model.SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(hashKey), nil
		},
	)

	if err != nil {
		return nil, "Error while parsing token: " + err.Error()
	}

	claims, ok := token.Claims.(*model.SignedDetails)
	if !ok {
		return nil, "Invalid token"
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, "Token is expired"
	}

	return claims, ""
}

// Update tokens in database
func UpdateAllTokens(signedToken, signedRefreshToken, UserID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updatedAt := time.Now()
	updateObj := bson.D{
		{Key: "token", Value: signedToken},
		{Key: "refresh_token", Value: signedRefreshToken},
		{Key: "updated_at", Value: updatedAt},
	}

	upsert := true
	filter := bson.M{"userid": UserID}
	opt := options.UpdateOptions{Upsert: &upsert}

	log.Printf("Updating tokens for user: %v", UserID)

	userCollection := database.GetUserCollection()
	if userCollection == nil {
		log.Printf("Error: userCollection is nil")
		return fmt.Errorf("userCollection is nil")
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	if err != nil {
		log.Printf("Failed to update tokens for user %s: %v", UserID, err)
		return err
	}

	log.Printf("Successfully updated tokens for user %s", UserID)
	return nil
}

// Refresh access token using refresh token
func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		RespondWithError(c, 401, "No refresh token provided", "")
		return
	}

	claims, msg := ValidateToken(refreshToken)
	if msg != "" {
		RespondWithError(c, 401, "Invalid refresh token", msg)
		return
	}

	newAccessToken, newRefreshToken, accessExpiry, refreshExpiry, err := GenerateAllTokens(
		claims.Email, claims.Username, claims.UserType, claims.Uid,
	)
	if err != nil {
		RespondWithError(c, 500, "Failed to generate new tokens", err.Error())
		return
	}

	c.SetCookie("access_token", newAccessToken, int(accessExpiry-time.Now().Unix()), "/", "", true, true)
	c.SetCookie("refresh_token", newRefreshToken, int(refreshExpiry-time.Now().Unix()), "/", "", true, true)

	RespondWithSuccess(c, 200, "Token refreshed successfully", gin.H{"access_token": newAccessToken})
}
