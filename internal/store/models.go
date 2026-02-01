package store

import (
	"time"

	"github.com/google/uuid"
)

type News struct {
	Id        uuid.UUID
	Author    string
	Title     string
	Summary   string
	Content   string
	CreatedAt time.Time
	Source    string
	Tags      []string
}
