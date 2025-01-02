package main

import (
	"time"

	"github.com/segmentio/ksuid"
)

// Task represents data about a task
type Task struct {
	ID      ksuid.KSUID `json:"id" validate:"required"`
	Title   string      `json:"title" validate:"required,min=1,max=140"`
	Status  bool        `json:"status"`
	Created time.Time   `json:"created" validate:"required"`
}
