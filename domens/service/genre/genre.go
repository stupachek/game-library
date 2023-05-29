package genre

import (
	"errors"
	"game-library/domens/models"

	"github.com/google/uuid"
)

var (
	ErrUnknownGenre error = errors.New("unknown genre")
	ErrGenreId            = errors.New("can't parse genre id")
)

type IGenreRepo interface {
	GetGenre(id uuid.UUID) (models.Genre, error)
	GetGenreByName(name string) (models.Genre, error)
	GetGenresList() ([]models.Genre, error)
	CreateGenre(genre models.Genre) error
	UpdateGenre(id uuid.UUID, genre models.Genre) error
	Delete(id uuid.UUID) error
}

type GenreService struct {
	GenreRepo IGenreRepo
}

func NewGenreService(repo IGenreRepo) GenreService {
	return GenreService{
		GenreRepo: repo,
	}
}

func (g *GenreService) GetGenreByName(name string) (models.Genre, error) {
	genre, err := g.GenreRepo.GetGenreByName(name)
	if err != nil {
		return models.Genre{}, err
	}
	if genre.ID == (uuid.UUID{}) {
		return models.Genre{}, ErrUnknownGenre
	}
	return genre, nil
}

func (g *GenreService) GetGenre(idStr string) (models.Genre, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Genre{}, ErrGenreId
	}
	genre, err := g.GenreRepo.GetGenre(id)
	if err != nil {
		return models.Genre{}, err
	}
	if genre.ID == (uuid.UUID{}) {
		return models.Genre{}, ErrUnknownGenre
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

func (g *GenreService) UpdateGenre(idStr string, genreModel models.Genre) (models.Genre, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Genre{}, ErrGenreId
	}
	err = g.GenreRepo.UpdateGenre(id, genreModel)
	if err != nil {
		return models.Genre{}, err
	}
	genre := models.Genre{
		ID:   id,
		Name: genreModel.Name,
	}
	return genre, nil
}

func (p *GenreService) DeleteGenre(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ErrGenreId
	}
	if _, err := p.GenreRepo.GetGenre(id); err != nil {
		return err
	}
	return p.GenreRepo.Delete(id)
}
