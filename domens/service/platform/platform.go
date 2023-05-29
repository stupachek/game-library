package platform

import (
	"errors"
	"game-library/domens/models"

	"github.com/google/uuid"
)

var (
	ErrPlatformId            = errors.New("can't parse platform id")
	ErrUnknownPlatform error = errors.New("unknown platform")
)

type IPlatformRepo interface {
	GetPlatform(id uuid.UUID) (models.Platform, error)
	GetPlatformByName(name string) (models.Platform, error)
	GetPlatformsList() ([]models.Platform, error)
	CreatePlatform(platform models.Platform) error
	UpdatePlatform(id uuid.UUID, platform models.Platform) error
	Delete(id uuid.UUID) error
}

type PlatformService struct {
	PlatformRepo IPlatformRepo
}

func NewPlatformService(repo IPlatformRepo) PlatformService {
	return PlatformService{
		PlatformRepo: repo,
	}
}

func (p *PlatformService) GetPlatformByName(name string) (models.Platform, error) {
	platform, err := p.PlatformRepo.GetPlatformByName(name)
	if err != nil {
		return models.Platform{}, err
	}
	if platform.ID == (uuid.UUID{}) {
		return models.Platform{}, ErrUnknownPlatform
	}
	return platform, nil
}

func (p *PlatformService) GetPlatform(idStr string) (models.Platform, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Platform{}, ErrPlatformId
	}
	platform, err := p.PlatformRepo.GetPlatform(id)
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

func (g *PlatformService) UpdatePlatform(idStr string, platformModel models.Platform) (models.Platform, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Platform{}, ErrPlatformId
	}
	err = g.PlatformRepo.UpdatePlatform(id, platformModel)
	if err != nil {
		return models.Platform{}, err
	}
	platform := models.Platform{
		ID:   id,
		Name: platformModel.Name,
	}
	return platform, nil
}

func (p *PlatformService) DeletePlatform(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ErrPlatformId
	}
	if _, err := p.PlatformRepo.GetPlatform(id); err != nil {
		return err
	}
	return p.PlatformRepo.Delete(id)
}
