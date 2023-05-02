package repository

import (
	"database/sql"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

type TestGenreRepo struct {
	Genres map[uuid.UUID]*models.Genre
	sync.Mutex
}

func NewGenreRepo() *TestGenreRepo {
	return &TestGenreRepo{
		Genres: make(map[uuid.UUID]*models.Genre),
	}
}

type PostgresGenreRepo struct {
	DB *sql.DB
}

func NewPostgresGenreRepo(DB *sql.DB) *PostgresGenreRepo {
	return &PostgresGenreRepo{
		DB: DB,
	}
}

func (p *PostgresGenreRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS genres(
        id UUID PRIMARY KEY,
    	name VARCHAR NOT NULL UNIQUE
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (t *TestGenreRepo) GetGenre(name string) (models.Genre, error) {
	for _, genre := range t.Genres {
		if genre.Name == name {
			return *genre, nil
		}
	}
	return models.Genre{}, nil
}

// TODO
func (p *PostgresGenreRepo) GetGenre(name string) (models.Genre, error) {
	return models.Genre{}, nil
}

func (t *TestGenreRepo) GetGenresList() []models.Genre {
	genres := make([]models.Genre, 0)
	for _, user := range t.Genres {
		genres = append(genres, *user)
	}
	return genres

}

// TODO
func (p *PostgresGenreRepo) GetGenresList() []models.Genre {
	return []models.Genre{}
}

func (t *TestGenreRepo) CreateGenre(genre models.Genre) error {
	t.Genres[genre.ID] = &genre
	return nil
}
func (p *PostgresGenreRepo) CreateGenre(genre models.Genre) error {
	_, err := p.DB.Exec("INSERT INTO genres(id, name) values($1, $2)", genre.ID, genre.Name)
	return err
}
