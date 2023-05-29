package genre_repo

import (
	"database/sql"
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var ErrUpdateFailed error = errors.New("update failed")

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

func (t *TestGenreRepo) GetGenre(id uuid.UUID) (models.Genre, error) {
	return *t.Genres[id], nil
}

// func (p *PostgresGenreRepo) GetGenre(name string) (models.Genre, error) {
// 	row := p.DB.QueryRow("SELECT id, name FROM genres WHERE name = $1", name)
// 	var genre models.Genre
// 	if err := row.Scan(&genre.ID, &genre.Name); err != nil {
// 		return models.Genre{}, err
// 	}
// 	return genre, nil
// }

func (p *PostgresGenreRepo) GetGenre(id uuid.UUID) (models.Genre, error) {
	row := p.DB.QueryRow("SELECT id, name FROM genres WHERE id = $1", id)
	var genre models.Genre
	if err := row.Scan(&genre.ID, &genre.Name); err != nil {
		return models.Genre{}, err
	}
	return genre, nil
}

func (p *PostgresGenreRepo) UpdateGenre(id uuid.UUID, genre models.Genre) error {
	res, err := p.DB.Exec("UPDATE genres SET name = $1 WHERE  id = $2", genre.Name, id)
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

func (t *TestGenreRepo) UpdateGenre(id uuid.UUID, genre models.Genre) error {
	t.Genres[id].Name = genre.Name
	return nil
}

func (p *PostgresGenreRepo) Delete(id uuid.UUID) error {
	_, err := p.DB.Exec("DELETE FROM genres WHERE id = $1", id)
	return err
}
func (t *TestGenreRepo) Delete(id uuid.UUID) error {
	delete(t.Genres, id)
	return nil
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

func (t *TestGenreRepo) GetGenreByName(name string) (models.Genre, error) {
	for _, genre := range t.Genres {
		if genre.Name == name {
			return *genre, nil
		}
	}
	return models.Genre{}, nil
}

func (p *PostgresGenreRepo) GetGenreByName(name string) (models.Genre, error) {
	row := p.DB.QueryRow("SELECT id, name FROM genres WHERE name = $1", name)
	var genre models.Genre
	if err := row.Scan(&genre.ID, &genre.Name); err != nil {
		return models.Genre{}, err
	}
	return genre, nil
}
