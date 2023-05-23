package genre

import (
	"fmt"
	"game-library/domens/models"

	"github.com/google/uuid"
)

type IGenreRepo interface {
	GetGenre(name string) (models.Genre, error)
	GetGenresList() ([]models.Genre, error)
	CreateGenre(genre models.Genre) error
}

type GenreService struct {
	GenreRepo IGenreRepo
}

func NewGenreService(repo IGenreRepo) GenreService {
	return GenreService{
		GenreRepo: repo,
	}
}

func (g *GenreService) GetGenre(name string) (models.Genre, error) {
	genre, err := g.GenreRepo.GetGenre(name)
	if err != nil {
		return models.Genre{}, err
	}
	if genre.ID == (uuid.UUID{}) {
		return models.Genre{}, fmt.Errorf("unknown genre: %s", name)
	}
	return genre, nil
}

func (g *GenreService) GetGenresList() ([]models.Genre, error) {
	return g.GenreRepo.GetGenresList()
}

func (g *GenreService) CreateGenre(genreModel models.Genre) (models.Genre, error) {
	id, _ := uuid.NewRandom()
	genre := models.Genre{
		ID:   id,
		Name: genreModel.Name,
	}
	err := g.GenreRepo.CreateGenre(genre)
	return genre, err
}
