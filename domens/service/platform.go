package service

import (
	"game-library/domens/models"
	"game-library/domens/repository"
)

type PlatformService struct {
	PlatformRepo repository.IPlatformRepo
}

func (p *PlatformService) GetPlatform(name string) (models.Platform, error) {
	return p.PlatformRepo.GetPlatform(name)
}
