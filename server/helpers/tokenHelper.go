package helper

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	database "task-manager/server/database"
	model "task-manager/server/models"
)

var userCollection *mongo.Collection = database.GetUserCollection()
var SECRET_KEY string

const (
	AccessTokenExpiry  = 24
	RefreshTokenExpiry = 128
)

// GenerateAllTokens generates both access and refresh tokens
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

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("error generating access token: %v", err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("error generating refresh token: %v", err)
		return "", "", err
	}

	return token, refreshToken, nil
}

// ValidateToken validates the given token and returns the claims
func ValidateToken(signedToken string) (claims *model.SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&model.SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
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

// UpdateAllTokens updates the tokens and timestamp in the database
func UpdateAllTokens(signedToken, signedRefreshToken, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updatedAt := time.Now()
	updateObj := bson.D{
		{Key: "token", Value: signedToken},
		{Key: "refresh_token", Value: signedRefreshToken},
		{Key: "updated_at", Value: updatedAt},
	}

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{Upsert: &upsert}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	if err != nil {
		log.Printf("failed to update tokens for user %s: %v", userId, err)
		return err
	}
	return nil
}
