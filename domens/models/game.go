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
}

func NewGame(publisherId uuid.UUID, title, description, imageLink string, ageRestriction, releaseYear int) Game {
	return Game{
		PublisherId:    publisherId,
		Title:          title,
		Description:    description,
		ImageLink:      imageLink,
		AgeRestriction: ageRestriction,
		ReleaseYear:    releaseYear,
	}
}

type Publisher struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Rating struct {
	ID     uuid.UUID
	UserId uuid.UUID
	GameId uuid.UUID
	Value  int
}

type Platform struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PlatformsOnGames struct {
	ID         uint
	GameId     uuid.UUID
	PlatformId uuid.UUID
}

type Genre struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GenresOnGames struct {
	ID      uint
	GameId  uuid.UUID
	GenreId uuid.UUID
}

type PublisherModel struct {
	Name string `json:"name" binding:"required"`
}
