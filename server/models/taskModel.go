package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents data about a task
type Task struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	User    string             `json:"User" validate:"required"`
	Title   string             `json:"title" validate:"required,min=1,max=140"`
	Status  bool               `json:"status"`
	Created time.Time          `json:"created" validate:"required"`
	Updated time.Time          `json:"updated"`
}
