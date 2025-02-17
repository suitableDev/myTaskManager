package helper

import (
	"log"
	"os"
	"task-manager/server/database"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email    string
	Username string
	Uid      string
	UserType string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.GetUserCollection()
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, userName string, userType string, uid string) (string, string, error) {
	claims := &SignedDetails{
		Email:    email,
		Username: userName,
		Uid:      uid,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(128)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	return token, refreshToken, nil
}
