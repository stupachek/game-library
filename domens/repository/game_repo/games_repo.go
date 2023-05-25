package game_repo

import (
	"database/sql"
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrDublicateTitle  error = errors.New("game with the title is already exist")
	ErrDublicateGameID error = errors.New("game with the ID is already exist")
)

type TestGameRepo struct {
	Games map[uuid.UUID]*models.Game
	sync.Mutex
}

func NewGameRepo() *TestGameRepo {
	return &TestGameRepo{
		Games: make(map[uuid.UUID]*models.Game),
	}
}

type PostgresGameRepo struct {
	DB *sql.DB
}

func NewPostgresGameRepo(DB *sql.DB) *PostgresGameRepo {
	return &PostgresGameRepo{
		DB: DB,
	}
}

// TODO: publisherId REFERENCES
func (p *PostgresGameRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS games(
        id UUID PRIMARY KEY,
		publisherId UUID,
    	title VARCHAR NOT NULL UNIQUE,
        description VARCHAR ,
		imageLink VARCHAR,
		ageRestriction INT,
		releaseYear INT,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (t *TestGameRepo) GetGames() ([]models.Game, error) {
	games := make([]models.Game, 0)
	for _, game := range t.Games {
		games = append(games, *game)
	}
	return games, nil
}

func (p *PostgresGameRepo) GetGames() ([]models.Game, error) {
	rows, err := p.DB.Query("SELECT id, PublisherId, Title,Description, ImageLink, AgeRestriction,   ReleaseYear, UpdatedAt FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var games []models.Game
	for rows.Next() {
		var game models.Game
		if err := rows.Scan(&game.ID, &game.PublisherId, &game.Title, &game.Description, &game.ImageLink, &game.AgeRestriction, &game.ReleaseYear, &game.UpdatedAt); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

func (t *TestGameRepo) GetGameById(id uuid.UUID) (models.Game, error) {
	for _, game := range t.Games {
		if game.ID == id {
			return *game, nil
		}
	}
	return models.Game{}, nil
}

func (p *PostgresGameRepo) GetGameById(id uuid.UUID) (models.Game, error) {
	row := p.DB.QueryRow("SELECT id, PublisherId, Title,Description, ImageLink, AgeRestriction,   ReleaseYear, UpdatedAt  FROM games WHERE id = $1", id)

	var game models.Game
	if err := row.Scan(&game.ID, &game.PublisherId, &game.Title, &game.Description, &game.ImageLink, &game.AgeRestriction, &game.ReleaseYear, &game.UpdatedAt); err != nil {
		return models.Game{}, err
	}
	return game, nil
}

func (t *TestGameRepo) CreateGame(game models.Game) error {
	if err := t.checkIfExist(game); err != nil {
		return err
	}
	t.Games[game.ID] = &game
	return nil
}

func (p *PostgresGameRepo) CreateGame(game models.Game) error {
	_, err := p.DB.Exec("INSERT INTO games(id, publisherId, title, description, imageLink, ageRestriction, releaseYear) values($1, $2, $3,  $4, $5, $6, $7)", game.ID, game.PublisherId, game.Title, game.Description, game.ImageLink, game.AgeRestriction, game.ReleaseYear)
	return err
}

func (t *TestGameRepo) checkIfExist(game models.Game) error {
	for _, g := range t.Games {
		switch {
		case g.ID == game.ID:
			return ErrDublicateGameID
		case g.Title == game.Title:
			return ErrDublicateTitle
		}
	}
	return nil
}

func (t *TestGameRepo) Setup() {
	t.Games[uuid.UUID{111}] = &models.Game{
		ID:             [16]byte{111},
		PublisherId:    [16]byte{123},
		Title:          "test",
		Description:    "test",
		ImageLink:      "test",
		AgeRestriction: 12,
		ReleaseYear:    2012,
	}
}
