package repository

import (
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

//var ErrUnknownPlatform error = errors.New("unknown platform")

type IPlatformRepo interface {
	GetPlatform(name string) (models.Platform, error)
	GetPlatformsList() []models.Platform
	CreatePlatform(platform models.Platform) models.Platform
}

type TestPlatformRepo struct {
	Platforms map[uuid.UUID]*models.Platform
	sync.Mutex
}

func NewPlatformRepo() *TestPlatformRepo {
	return &TestPlatformRepo{
		Platforms: make(map[uuid.UUID]*models.Platform),
	}
}

func (t *TestPlatformRepo) GetPlatform(name string) (models.Platform, error) {
	for _, Platform := range t.Platforms {
		if Platform.Name == name {
			return *Platform, nil
		}
	}
	return models.Platform{}, ErrDublicateEmail
}

func (t *TestPlatformRepo) GetPlatformsList() []models.Platform {
	platforms := make([]models.Platform, 0)
	for _, user := range t.Platforms {
		platforms = append(platforms, *user)
	}
	return platforms

}

func (t *TestPlatformRepo) CreatePlatform(platform models.Platform) models.Platform {
	t.Platforms[platform.ID] = &platform
	return platform
}
