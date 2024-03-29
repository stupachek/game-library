package platform_repo

import (
	"database/sql"
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrUpdateFailed error = errors.New("update failed")
	ErrDeleteFailed error = errors.New("delete failed")
)

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

func (p *PostgresPlatformRepo) GetPlatform(id uuid.UUID) (models.Platform, error) {
	row := p.DB.QueryRow("SELECT id, name FROM platforms WHERE id = $1", id)
	var platform models.Platform
	if err := row.Scan(&platform.ID, &platform.Name); err != nil {
		return models.Platform{}, err
	}
	return platform, nil
}

func (t *TestPlatformRepo) GetPlatform(id uuid.UUID) (models.Platform, error) {
	return *t.Platforms[id], nil
}

func (t *TestPlatformRepo) GetPlatformsList() ([]models.Platform, error) {
	platforms := make([]models.Platform, 0)
	for _, user := range t.Platforms {
		platforms = append(platforms, *user)
	}
	return platforms, nil

}

func (p *PostgresPlatformRepo) GetPlatformsList() ([]models.Platform, error) {
	rows, err := p.DB.Query("SELECT id, name FROM platforms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var platforms []models.Platform
	for rows.Next() {
		var platform models.Platform
		if err := rows.Scan(&platform.ID, &platform.Name); err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}
	return platforms, nil

}

func (t *TestPlatformRepo) CreatePlatform(platform models.Platform) error {
	t.Platforms[platform.ID] = &platform
	return nil
}

func (p *PostgresPlatformRepo) CreatePlatform(platform models.Platform) error {
	_, err := p.DB.Exec("INSERT INTO platforms(id, name) values($1, $2)", platform.ID, platform.Name)
	return err
}

func (t *TestPlatformRepo) Setup() {
	t.Platforms[uuid.UUID{111}] = &models.Platform{
		ID:   uuid.UUID{111},
		Name: "test",
	}
}

func (p *PostgresPlatformRepo) UpdatePlatform(id uuid.UUID, platform models.Platform) error {
	res, err := p.DB.Exec("UPDATE platforms SET name = $1 WHERE  id = $2", platform.Name, id)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return ErrUpdateFailed
	}
	return nil
}

func (t *TestPlatformRepo) UpdatePlatform(id uuid.UUID, platform models.Platform) error {
	t.Platforms[id].Name = platform.Name
	return nil
}

func (p *PostgresPlatformRepo) Delete(id uuid.UUID) error {
	_, err := p.DB.Exec("DELETE FROM platforms WHERE id = $1", id)
	return err
}
func (t *TestPlatformRepo) Delete(id uuid.UUID) error {
	delete(t.Platforms, id)
	return nil
}

func (p *PostgresPlatformRepo) GetPlatformByName(name string) (models.Platform, error) {
	row := p.DB.QueryRow("SELECT id, name FROM platforms WHERE name = $1", name)
	var platform models.Platform
	if err := row.Scan(&platform.ID, &platform.Name); err != nil {
		return models.Platform{}, err
	}
	return platform, nil
}

func (t *TestPlatformRepo) GetPlatformByName(name string) (models.Platform, error) {
	for _, Platform := range t.Platforms {
		if Platform.Name == name {
			return *Platform, nil
		}
	}
	return models.Platform{}, nil
}
