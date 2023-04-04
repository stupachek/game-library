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
	Platforms      []PlatformsOnGames
	Genres         []GenresOnGames
}

type Publisher struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Games []Game    `json:"games"`
}

type Rating struct {
	ID     uuid.UUID
	UserId uuid.UUID
	GameId uuid.UUID
	Value  int
}

type Platform struct {
	ID    uuid.UUID
	Name  string
	Games []PlatformsOnGames
}

type PlatformsOnGames struct {
	GameId     uuid.UUID
	PlatformId uuid.UUID
}

type Genre struct {
	ID    uuid.UUID
	Name  string
	Games []GenresOnGames
}

type GenresOnGames struct {
	GameId  uuid.UUID
	GenreId uuid.UUID
}

type PublisherModel struct {
	Name string `json:"name" binding:"required"`
}
