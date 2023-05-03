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
	GetGames() ([]models.Game, error)
	//UpdateGame(id uuid.UUID, game models.Game) (models.Game, error)
	//Delete(id uuid.UUID)
}

type IGenreRepo interface {
	GetGenre(name string) (models.Genre, error)
	GetGenresList() []models.Genre
	CreateGenre(genre models.Genre) error
}

type IGenresOnGamesRepo interface {
	CreateGenresOnGames(genresOnGames models.GenresOnGames) error
}

type IPlatformsOnGamesRepo interface {
	CreatePlatformsOnGames(PlatformsOnGames models.PlatformsOnGames) error
}

type IPublisherRepo interface {
	GetPublishersList() ([]models.Publisher, error)
	CreatePublisher(pub models.Publisher) error
	GetPublisherById(id uuid.UUID) (*models.Publisher, error)
	UpdatePublisher(id uuid.UUID, publisher models.PublisherModel) error
	Delete(id uuid.UUID) error
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

func (g *GameService) GetGamesList() ([]models.Game, error) {
	return g.GameRepo.GetGames()
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
