package game

import (
	"errors"
	"game-library/domens/models"

	"github.com/google/uuid"
)

type GameService struct {
	GameRepo             IGameRepo
	GenresOnGamesRepo    IGenresOnGamesRepo
	PlatformsOnGamesRepo IPlatformsOnGamesRepo
}

type IGameRepo interface {
	CreateGame(game models.Game) error
	GetGameById(id uuid.UUID) (models.GameRespons, error)
	GetGames() ([]models.GameRespons, error)
	//UpdateGame(id uuid.UUID, game models.Game) (models.Game, error)
	//Delete(id uuid.UUID)
}

type IGenresOnGamesRepo interface {
	CreateGenresOnGames(genresOnGames models.GenresOnGames) error
}

type IPlatformsOnGamesRepo interface {
	CreatePlatformsOnGames(PlatformsOnGames models.PlatformsOnGames) error
}

var ErrParseId = errors.New("can't parse id")
var ErrUnknownId = errors.New("unknown id")

func NewGameService(gameRepo IGameRepo, genresOnGamesRepo IGenresOnGamesRepo, platformsOnGamesRepo IPlatformsOnGamesRepo) GameService {
	return GameService{
		GameRepo:             gameRepo,
		GenresOnGamesRepo:    genresOnGamesRepo,
		PlatformsOnGamesRepo: platformsOnGamesRepo,
	}
}

func (g *GameService) GetGamesList() ([]models.GameRespons, error) {
	return g.GameRepo.GetGames()
}

func (g *GameService) GetGame(idStr string) (models.GameRespons, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.GameRespons{}, ErrParseId
	}
	game, err := g.GameRepo.GetGameById(id)
	if err != nil {
		return models.GameRespons{}, err
	}
	if game.ID == ("") {
		return models.GameRespons{}, ErrUnknownId
	}
	return game, nil
}

func (g *GameService) CreateGame(game models.Game, genres []models.Genre, plaforms []models.Platform) (models.Game, error) {
	gameId, err := uuid.NewRandom()
	if err != nil {
		return models.Game{}, err
	}
	game.ID = gameId
	err = g.GameRepo.CreateGame(game)
	if err != nil {
		return models.Game{}, err

	}
	for _, genre := range genres {
		err := g.GenresOnGamesRepo.CreateGenresOnGames(
			models.GenresOnGames{
				GameId:  gameId,
				GenreId: genre.ID,
			})
		if err != nil {
			return models.Game{}, err
		}
	}
	for _, plaform := range plaforms {
		err := g.PlatformsOnGamesRepo.CreatePlatformsOnGames(
			models.PlatformsOnGames{
				GameId:     gameId,
				PlatformId: plaform.ID,
			})
		if err != nil {
			return models.Game{}, err

		}
	}

	return game, nil
}
