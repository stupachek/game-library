package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type QueryParams struct {
	Skip uint64
	Take uint64
}

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

type GameRespons struct {
	ID             string        `json:"id"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	ImageLink      string        `json:"imageLink"`
	AgeRestriction int           `json:"ageRestriction"`
	ReleaseYear    int           `json:"releaseYear"`
	Publisher      Publisher     `json:"publisher"`
	Genres         GenreList     `json:"genres"`
	Platform       PlatformsList `json:"platforms"`
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

type GenreList []Genre

func (gl *GenreList) Scan(src any) error {
	return pq.GenericArray{A: gl}.Scan(src)
}

func (g *Genre) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	}
	return json.Unmarshal(data, g)
}

type PlatformsList []Genre

func (pl *PlatformsList) Scan(src any) error {
	return pq.GenericArray{A: pl}.Scan(src)
}

func (p *Platform) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	}
	return json.Unmarshal(data, p)
}
