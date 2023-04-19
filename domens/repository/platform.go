package repository

import (
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var ErrUnknownPlatform error = errors.New("unknown platform")

type IPlatformRepo interface {
	GetPlatform(name string) (models.Platform, error)
}

type TestPlatformRepo struct {
	Platforms map[uuid.UUID]*models.Platform
	sync.Mutex
}

func (t *TestPlatformRepo) GetPlatform(name string) (models.Platform, error) {
	for _, Platform := range t.Platforms {
		if Platform.Name == name {
			return *Platform, nil
		}
	}
	return models.Platform{}, ErrDublicateEmail
}
