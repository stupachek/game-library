package service

import (
	"game-library/domens/models"
	"game-library/domens/repository"

	"github.com/google/uuid"
)

type PlatformService struct {
	PlatformRepo repository.IPlatformRepo
}

func NewPlatformService(repo repository.IPlatformRepo) PlatformService {
	return PlatformService{
		PlatformRepo: repo,
	}
}

func (p *PlatformService) GetPlatform(name string) (models.Platform, error) {
	return p.PlatformRepo.GetPlatform(name)
}

func (p *PlatformService) GetPlatformsList() ([]models.Platform, error) {
	return p.PlatformRepo.GetPlatformsList()
}

func (p *PlatformService) CreatePlatform(platformModel models.Platform) (models.Platform, error) {
	id, _ := uuid.NewRandom()
	platform := models.Platform{
		ID:   id,
		Name: platformModel.Name,
	}
	err := p.PlatformRepo.CreatePlatform(platform)
	return platform, err
}
