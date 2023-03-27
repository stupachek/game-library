package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID             uuid.UUID
	PublisherId    uuid.UUID
	Title          string
	Description    string
	ImageLink      string
	AgeRestriction int
	ReleaseYear    int
	UpdatedAt      time.Time
	Ratings        []Rating
	Comments       []Comment
}

type Publisher struct {
	ID    uuid.UUID
	Name  string
	Games []Game
}

type Rating struct {
	ID     uuid.UUID
	UserId uuid.UUID
	GameId uuid.UUID
	Value  int
}
