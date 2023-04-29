package repository

import (
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

type TestGenreRepo struct {
	Genres map[uuid.UUID]*models.Genre
	sync.Mutex
}

func NewGenreRepo() *TestGenreRepo {
	return &TestGenreRepo{
		Genres: make(map[uuid.UUID]*models.Genre),
	}
}

func (t *TestGenreRepo) GetGenre(name string) (models.Genre, error) {
	for _, genre := range t.Genres {
		if genre.Name == name {
			return *genre, nil
		}
	}
	return models.Genre{}, nil
}

func (t *TestGenreRepo) GetGenresList() []models.Genre {
	genres := make([]models.Genre, 0)
	for _, user := range t.Genres {
		genres = append(genres, *user)
	}
	return genres

}

func (t *TestGenreRepo) CreateGenre(genre models.Genre) models.Genre {
	t.Genres[genre.ID] = &genre
	return genre
}
