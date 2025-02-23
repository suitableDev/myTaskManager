package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

func GenerateAllTokens(email, userName, userType, uid string) (string, string, error) {
	claims := &model.SignedDetails{
		Email:    email,
		Username: userName,
		Uid:      uid,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(AccessTokenExpiry)).Unix(),
		},
	}

	refreshClaims := &model.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(RefreshTokenExpiry)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(hashKey))
	if err != nil {
		log.Printf("error generating access token: %v", err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(hashKey))
	if err != nil {
		log.Printf("error generating refresh token: %v", err)
		return "", "", err
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (claims *model.SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&model.SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(hashKey), nil
		},
	)

	if err != nil {
		return nil, "error while parsing token: " + err.Error()
	}

	claims, ok := token.Claims.(*model.SignedDetails)
	if !ok {
		return nil, "the token is invalid"
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, "token is expired"
	}
	return claims, ""
}

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
	log.Printf("filter: %v", filter)
	log.Printf("updateObj: %v", updateObj)

	userCollection := database.GetUserCollection()

	if userCollection == nil {
		log.Printf("Error: userCollection is nil")
		return fmt.Errorf("userCollection is nil")
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	if err != nil {
		log.Printf("failed to update tokens for user %s: %v", UserID, err)
		return err
	}

	log.Printf("Successfully updated tokens for user %s", UserID)
	return nil
}
