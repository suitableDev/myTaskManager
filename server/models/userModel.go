package model

import (
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        *string            `json:"email" validate:"email,required"`
	Username     *string            `json:"username" validate:"required,min=2,max=30"`
	Password     *string            `json:"password" validate:"required,min=6,max=30"`
	UserType     *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	UserID       *string            `json:"userid"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created"`
	UpdatedAt    time.Time          `json:"updated"`
}

type SignedDetails struct {
	Email    string
	Username string
	Uid      string
	UserType string
	jwt.StandardClaims
}
