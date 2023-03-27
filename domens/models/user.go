package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Email          string
	UserName       string
	BadgeColor     string
	Role           string
	HashedPassword string
	CreatedAt      time.Time //should it be only in db?
	Comments       []Comment
	Likes          []CommentLike
	Ratings        []Rating
}
