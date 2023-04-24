package service

import (
	"errors"
	"game-library/domens/models"
	"game-library/domens/repository"

	"github.com/google/uuid"
)

type GameService struct {
	GameRepo repository.IGameRepo
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

func (g *GameService) CreateGame(inputGame models.InputGame, dst string, genres []models.Genre, plaforms []models.Platform) error {
	publisherId, err := uuid.Parse(inputGame.PublisherId)
	if err != nil {
		return errors.New("can't parse publisherId id")
	}
	gameId, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	genresOnGames := make([]models.GenresOnGames, len(genres))
	for i, genre := range genres {
		genresOnGames[i] = models.GenresOnGames{
			GameId:  gameId,
			GenreId: genre.ID,
		}
	}
	platformsOnGames := make([]models.PlatformsOnGames, len(plaforms))
	for i, plaform := range plaforms {
		platformsOnGames[i] = models.PlatformsOnGames{
			GameId:     gameId,
			PlatformId: plaform.ID,
		}
	}
	game := models.Game{
		ID:             gameId,
		PublisherId:    publisherId,
		Title:          inputGame.Title,
		Description:    inputGame.Description,
		ImageLink:      dst,
		AgeRestriction: inputGame.AgeRestriction,
		ReleaseYear:    inputGame.ReleaseYear,
		Platforms:      platformsOnGames,
		Genres:         genresOnGames,
	}
	return g.GameRepo.CreateGame(game)
}
