package models

import (
	"mime/multipart"
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

type InputGame struct {
	Title          string                `form:"title" binding:"required"`
	Description    string                `form:"description"`
	File           *multipart.FileHeader `form:"file"`
	PublisherId    string                `form:"publisherId"`
	AgeRestriction int                   `form:"ageRestriction"`
	ReleaseYear    int                   `form:"releaseYear"`
	Genres         []string              `form:"genres"`
	Platforms      []string              `form:"platforms"`
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
	ID      uint
	GameId  uuid.UUID
	GenreId uuid.UUID
}

type PublisherModel struct {
	Name string `json:"name" binding:"required"`
}
