package genre_repo

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

func (p *PostgresGenreRepo) GetGenre(name string) (models.Genre, error) {
	row := p.DB.QueryRow("SELECT id, name FROM genres WHERE name = $1", name)
	var genre models.Genre
	if err := row.Scan(&genre.ID, &genre.Name); err != nil {
		return models.Genre{}, err
	}
	return genre, nil
}

func (t *TestGenreRepo) GetGenresList() ([]models.Genre, error) {
	genres := make([]models.Genre, 0)
	for _, user := range t.Genres {
		genres = append(genres, *user)
	}
	return genres, nil

}

func (p *PostgresGenreRepo) GetGenresList() ([]models.Genre, error) {
	rows, err := p.DB.Query("SELECT id, name FROM genres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (t *TestGenreRepo) CreateGenre(genre models.Genre) error {
	t.Genres[genre.ID] = &genre
	return nil
}
func (p *PostgresGenreRepo) CreateGenre(genre models.Genre) error {
	_, err := p.DB.Exec("INSERT INTO genres(id, name) values($1, $2)", genre.ID, genre.Name)
	return err
}

func (t *TestGenreRepo) Setup() {
	t.Genres[uuid.UUID{111}] = &models.Genre{
		ID:   uuid.UUID{111},
		Name: "test",
	}
}
