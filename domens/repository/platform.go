package repository

import (
	"database/sql"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

//var ErrUnknownPlatform error = errors.New("unknown platform")

type IPlatformRepo interface {
	GetPlatform(name string) (models.Platform, error)
	GetPlatformsList() []models.Platform
	CreatePlatform(platform models.Platform) error
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

type PostgresPlatformRepo struct {
	DB *sql.DB
}

func NewPostgresPlatformRepo(DB *sql.DB) *PostgresPlatformRepo {
	return &PostgresPlatformRepo{
		DB: DB,
	}
}

func (p *PostgresPlatformRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS platforms(
        id UUID PRIMARY KEY,
    	name VARCHAR NOT NULL UNIQUE
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (p *PostgresPlatformRepo) GetPlatform(name string) (models.Platform, error) {
	return models.Platform{}, nil
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

func (p *PostgresPlatformRepo) GetPlatformsList() []models.Platform {
	return []models.Platform{}

}

func (t *TestPlatformRepo) CreatePlatform(platform models.Platform) error {
	t.Platforms[platform.ID] = &platform
	return nil
}

func (p *PostgresPlatformRepo) CreatePlatform(platform models.Platform) error {
	_, err := p.DB.Exec("INSERT INTO platforms(id, name) values($1, $2)", platform.ID, platform.Name)
	return err
}
