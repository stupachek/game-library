package service

import (
	"fmt"
	"game-library/domens/models"
	"game-library/domens/repository"

	"github.com/google/uuid"
)

type GenreService struct {
	GenreRepo repository.IGenreRepo
}

func (p *GenreService) GetGenre(name string) (models.Genre, error) {
	genre, err := p.GenreRepo.GetGenre(name)
	if err != nil {
		return models.Genre{}, err
	}
	if genre.ID == (uuid.UUID{}) {
		return models.Genre{}, fmt.Errorf("unknown genre: %s", name)
	}
	return genre, nil
}
