package repository

import (
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrDublicateTitle  error = errors.New("game with the title is already exist")
	ErrDublicateGameID error = errors.New("game with the ID is already exist")
)

type IGameRepo interface {
	CreateGame(game models.Game) error
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

func (t *TestGameRepo) CreateGame(game models.Game) error {
	if err := t.checkIfExist(game); err != nil {
		return err
	}
	t.Games[game.ID] = &game
	return nil
}

func (t *TestGameRepo) checkIfExist(game models.Game) error {
	for _, g := range t.Games {
		switch {
		case g.ID == game.ID:
			return ErrDublicateGameID
		case g.Title == game.Title:
			return ErrDublicateTitle
		}
	}
	return nil
}
