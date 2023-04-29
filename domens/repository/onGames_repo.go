package repository

import (
	"errors"
	"game-library/domens/models"
)

var ErrDublicateIDModel error = errors.New("model with the ID already exist")

type TestGenresOnGamesRepo struct {
	GenresOnGames map[uint]*models.GenresOnGames
}

func (t *TestGenresOnGamesRepo) CreateGenresOnGames(genresOnGames models.GenresOnGames) error {
	if _, ok := t.GenresOnGames[genresOnGames.ID]; ok {
		return ErrDublicateIDModel
	}
	t.GenresOnGames[genresOnGames.ID] = &genresOnGames
	return nil
}

type TestPlatformsOnGamesRepo struct {
	PlatformsOnGames map[uint]*models.PlatformsOnGames
}

func (t *TestPlatformsOnGamesRepo) CreatePlatformsOnGames(platformsOnGames models.PlatformsOnGames) error {
	if _, ok := t.PlatformsOnGames[platformsOnGames.ID]; ok {
		return ErrDublicateIDModel
	}
	t.PlatformsOnGames[platformsOnGames.ID] = &platformsOnGames
	return nil
}
