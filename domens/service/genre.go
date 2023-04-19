package service

import (
	"game-library/domens/models"
	"game-library/domens/repository"
)

type GenreService struct {
	GenreRepo repository.IGenreRepo
}

func (p *GenreService) GetGenre(name string) (models.Genre, error) {
	return p.GenreRepo.GetGenre(name)
}
