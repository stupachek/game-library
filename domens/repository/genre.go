package repository

import (
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var ErrUnknownGenre error = errors.New("unknown genre")

type IGenreRepo interface {
	GetGenre(name string) (models.Genre, error)
}

type TestGenreRepo struct {
	Genres map[uuid.UUID]*models.Genre
	sync.Mutex
}

func (t *TestGenreRepo) GetGenre(name string) (models.Genre, error) {
	for _, genre := range t.Genres {
		if genre.Name == name {
			return *genre, nil
		}
	}
	return models.Genre{}, ErrUnknownGenre
}
