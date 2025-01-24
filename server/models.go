package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents data about a task
type Task struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `json:"title" validate:"required,min=1,max=140"`
	Status  bool               `json:"status"`
	Created time.Time          `json:"created" validate:"required"`
}
