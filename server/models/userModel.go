package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        *string            `json:"email" validate:"email, required"`
	Username     *string            `json:"username" validate:"required, min=2,max=30"`
	Password     *string            `json:"password" validate:"required,min=6,max=30"`
	UserType     *string            `json:"user_type" validate:"require, eq=ADMIN'eq=USER"`
	UserId       *string            `json:"user_id"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created"`
	UpdatedAt    time.Time          `json:"updated"`
}
