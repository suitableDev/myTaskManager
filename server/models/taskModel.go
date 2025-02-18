package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	User    string             `json:"User" validate:"required"`
	UserId  string             `json:"userId"`
	Title   string             `json:"title" validate:"required,min=1,max=140"`
	Status  bool               `json:"status"`
	Created time.Time          `json:"created" validate:"required"`
	Updated time.Time          `json:"updated"`
}
