package service

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
	GetGameById(id uuid.UUID) (models.Game, error)
	GetGamesList() []models.Game
	//UpdateGame(id uuid.UUID, game models.Game) (models.Game, error)
	//Delete(id uuid.UUID)
}

type IGenreRepo interface {
	GetGenre(name string) (models.Genre, error)
	GetGenresList() []models.Genre
	CreateGenre(genre models.Genre) models.Genre
}

type IGenresOnGamesRepo interface {
	CreateGenresOnGames(genresOnGames models.GenresOnGames) error
}

type IPlatformsOnGamesRepo interface {
	CreatePlatformsOnGames(PlatformsOnGames models.PlatformsOnGames) error
}

type IPublisherRepo interface {
	GetPublishersList() []models.Publisher
	CreatePublisher(publisher models.Publisher) models.Publisher
	UpdatePublisher(id uuid.UUID, publisher models.PublisherModel) (models.Publisher, error)
	GetPublisherById(id uuid.UUID) (*models.Publisher, error)
	Delete(id uuid.UUID)
}

var ErrParseId = errors.New("can't parse id")
var ErrUnknownId = errors.New("unknown id")

func NewGameService(repo IGameRepo) GameService {
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
	game.ID = gameId
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
