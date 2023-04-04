package repository

import (
	"game-library/domens/models"
	"github.com/google/uuid"
	"sync"
)

type IGameRepo interface {
	//CreateGame(models.Game) error
	//GetGameById(id uuid.UUID) (*models.Game, error)
	GetGamesList() []models.Game
	//UpdateGame(id uuid.UUID, game models.Game) (models.Game, error)
	//Delete(id uuid.UUID)
}

type TestGameRepo struct {
	Games map[uuid.UUID]*models.Game
	sync.Mutex
}

func NewGameRepo() *TestGameRepo {
	return &TestGameRepo{
		Games: make(map[uuid.UUID]*models.Game),
	}
}

func (t *TestGameRepo) GetGamesList() []models.Game {
	games := make([]models.Game, 0)
	for _, user := range t.Games {
		games = append(games, *user)
	}
	return games

}
