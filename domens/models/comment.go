package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	GameId    uuid.UUID
	ParentId  uuid.UUID
	Value     string
	CreatedAt time.Time
}

type CommentLike struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	CommentId uuid.UUID
}
