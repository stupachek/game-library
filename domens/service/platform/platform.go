package platform

import (
	"errors"
	"game-library/domens/models"

	"github.com/google/uuid"
)

var ErrUnknownPlatform error = errors.New("unknown platform")

type IPlatformRepo interface {
	GetPlatform(name string) (models.Platform, error)
	GetPlatformsList() ([]models.Platform, error)
	CreatePlatform(platform models.Platform) error
}

type PlatformService struct {
	PlatformRepo IPlatformRepo
}

func NewPlatformService(repo IPlatformRepo) PlatformService {
	return PlatformService{
		PlatformRepo: repo,
	}
}

func (p *PlatformService) GetPlatform(name string) (models.Platform, error) {
	platform, err := p.PlatformRepo.GetPlatform(name)
	if err != nil {
		return models.Platform{}, err
	}
	if platform.ID == (uuid.UUID{}) {
		return models.Platform{}, ErrUnknownPlatform
	}
	return platform, nil
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
