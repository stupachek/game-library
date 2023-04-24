package service

import (
	"errors"
	"game-library/domens/models"
	"game-library/domens/repository"

	"github.com/google/uuid"
)

type GameService struct {
	GameRepo             repository.IGameRepo
	GenresOnGamesRepo    repository.IGenresOnGamesRepo
	PlatformsOnGamesRepo repository.IPlatformsOnGamesRepo
}

var ErrParseId = errors.New("can't parse id")
var ErrUnknownId = errors.New("unknown id")

func NewGameService(repo repository.IGameRepo) GameService {
	return GameService{
		GameRepo: repo,
	}
}

func (g *GameService) GetGamesList() []models.Game {
	games := g.GameRepo.GetGamesList()
	return games
}

func (g *GameService) GetGame(idStr string) (models.Game, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Game{}, ErrParseId
	}
	game, err := g.GameRepo.GetGameById(id)
	if err != nil {
		return models.Game{}, err
	}
	if game.ID == (uuid.UUID{}) {
		return models.Game{}, ErrUnknownId
	}
	return game, nil
}

func (g *GameService) CreateGame(game models.Game, genres []models.Genre, plaforms []models.Platform) error {
	gameId, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	err = g.GameRepo.CreateGame(game)
	if err != nil {
		return err
	}
	for _, genre := range genres {
		err := g.GenresOnGamesRepo.CreateGenresOnGames(
			models.GenresOnGames{
				GameId:  gameId,
				GenreId: genre.ID,
			})
		if err != nil {
			return err
		}
	}
	for _, plaform := range plaforms {
		err := g.PlatformsOnGamesRepo.CreatePlatformsOnGames(
			models.PlatformsOnGames{
				GameId:     gameId,
				PlatformId: plaform.ID,
			})
		if err != nil {
			return err
		}
	}

	return nil
}
