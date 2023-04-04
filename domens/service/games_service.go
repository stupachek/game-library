package service

import (
	"game-library/domens/models"
	"game-library/domens/repository"
)

type GameService struct {
	GameRepo repository.IGameRepo
}

func NewGameService(repo repository.IGameRepo) GameService {
	return GameService{
		GameRepo: repo,
	}
}

func (g *GameService) GetGamesList() []models.Game {
	games := g.GameRepo.GetGamesList()
	return games
}
