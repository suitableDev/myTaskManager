package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   string             `bson:"user_id" json:"user_id" validate:"required"`
	Username string             `bson:"username" json:"username" validate:"required"`
	Title    string             `bson:"title" json:"title" validate:"required,min=1,max=140"`
	Status   bool               `bson:"status" json:"status"`
	Created  time.Time          `bson:"created_at" json:"created_at"`
	Updated  time.Time          `bson:"updated_at" json:"updated_at"`
}
