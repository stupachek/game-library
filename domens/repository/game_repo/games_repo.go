package game_repo

import (
	"database/sql"
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func (t *TestGameRepo) GetGames(params models.QueryParams) ([]models.GameRespons, error) {
	games := make([]models.GameRespons, 0)
	for _, game := range t.Games {
		games = append(games, models.GameRespons{
			ID:             game.ID.String(),
			Title:          game.Title,
			Description:    game.Description,
			ImageLink:      game.ImageLink,
			AgeRestriction: game.AgeRestriction,
			ReleaseYear:    game.ReleaseYear,
		})
	}
	return games, nil
}

func (p *PostgresGameRepo) GetGames(params models.QueryParams) ([]models.GameRespons, error) {
	rows, err := p.DB.Query(`
    SELECT games.id, games.title, games.description, games.imagelink, games.ageRestriction, games.releaseYear,
           publishers.id AS publishersId, publishers.name AS publishersName,
           ARRAY_AGG(DISTINCT jsonb_build_object('id', genres.id, 'name', genres.name)) AS genres,
           ARRAY_AGG(DISTINCT jsonb_build_object('id', platforms.id, 'name', platforms.name)) AS platforms
    FROM games
    JOIN publishers ON games.publisherId = publishers.id
    JOIN genresOnGames ON games.id = genresOnGames.gameId
    JOIN genres ON genresOnGames.genreId = genres.id
    JOIN platformsOnGames ON games.id = platformsOnGames.gameId
    JOIN platforms ON platformsOnGames.platformId = platforms.id
    WHERE games.title LIKE $1
    GROUP BY games.id, games.title, games.description, games.imagelink, games.ageRestriction, games.releaseYear,
             publishers.id, publishers.name
    LIMIT $2 OFFSET $3;
`, params.SearchQuery, params.Take, params.Skip)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var games []models.GameRespons
	for rows.Next() {
		var game models.GameRespons
		if err := rows.Scan(&game.ID, &game.Title, &game.Description, &game.ImageLink, &game.AgeRestriction, &game.ReleaseYear, &game.Publisher.ID, &game.Publisher.Name, pq.Array(&game.Genres), pq.Array(&game.Platform)); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

func (t *TestGameRepo) GetGameById(id uuid.UUID) (models.GameRespons, error) {
	for _, game := range t.Games {
		if game.ID == id {
			return models.GameRespons{
				ID:             game.ID.String(),
				Title:          game.Title,
				Description:    game.Description,
				ImageLink:      game.ImageLink,
				AgeRestriction: game.AgeRestriction,
				ReleaseYear:    game.ReleaseYear,
			}, nil
		}
	}
	return models.GameRespons{}, nil
}

func (p *PostgresGameRepo) GetGameById(id uuid.UUID) (models.GameRespons, error) {
	row := p.DB.QueryRow(`SELECT games.id, games.title, games.description, games.imagelink, games.ageRestriction, games.releaseYear,
       publishers.id AS publishersId, publishers.name AS publishersName,
       ARRAY_AGG(DISTINCT jsonb_build_object('id', genres.id, 'name', genres.name)) AS genres,
       ARRAY_AGG(DISTINCT jsonb_build_object('id', platforms.id, 'name', platforms.name)) AS platforms
FROM games 
JOIN publishers ON games.publisherId = publishers.id 
JOIN genresOnGames ON games.id = genresOnGames.gameId 
JOIN genres ON genresOnGames.genreId = genres.id 
JOIN platformsOnGames ON games.id = platformsOnGames.gameId 
JOIN platforms ON platformsOnGames.platformId = platforms.id 
WHERE games.id = $1 
GROUP BY games.id, games.title, games.description, games.imagelink, games.ageRestriction, games.releaseYear,
         publishers.id, publishers.name;`, id)

	var game models.GameRespons
	if err := row.Scan(&game.ID, &game.Title, &game.Description, &game.ImageLink, &game.AgeRestriction, &game.ReleaseYear, &game.Publisher.ID, &game.Publisher.Name, pq.Array(&game.Genres), pq.Array(&game.Platform)); err != nil {
		return models.GameRespons{}, err
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

func (t *TestGameRepo) UpdateGame(id uuid.UUID, game models.Game) error {
	if err := t.checkIfExist(game); err != nil {
		return err
	}
	t.Games[id] = &game
	return nil
}

func (p *PostgresGameRepo) CreateGame(game models.Game) error {
	_, err := p.DB.Exec("INSERT INTO games(id, publisherId, title, description, imageLink, ageRestriction, releaseYear) values($1, $2, $3,  $4, $5, $6, $7)", game.ID, game.PublisherId, game.Title, game.Description, game.ImageLink, game.AgeRestriction, game.ReleaseYear)
	return err
}
func (p *PostgresGameRepo) UpdateGame(id uuid.UUID, game models.Game) error {
	_, err := p.DB.Exec("UPDATE games SET publisherId = $1, title = $2, description = $3, imageLink = $4, ageRestriction = $5, releaseYear = $6 WHERE id = $7",
		game.PublisherId, game.Title, game.Description, game.ImageLink, game.AgeRestriction, game.ReleaseYear, id)

	return err
}

func (p *PostgresGameRepo) DeleteGame(id uuid.UUID) error {
	_, err := p.DB.Exec("DELETE FROM games WHERE id = $1", id)
	return err
}
func (t *TestGameRepo) DeleteGame(id uuid.UUID) error {
	delete(t.Games, id)
	return nil
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
