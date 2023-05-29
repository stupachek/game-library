package game_repo

import (
	"database/sql"
	"errors"
	"game-library/domens/models"
)

var ErrDublicateIDModel error = errors.New("model with the ID already exist")

type TestGenresOnGamesRepo struct {
	GenresOnGames map[uint]*models.GenresOnGames
	counter       uint
}

func NewTestGenresOnGamesRepo() *TestGenresOnGamesRepo {
	return &TestGenresOnGamesRepo{
		GenresOnGames: map[uint]*models.GenresOnGames{},
	}
}

type PostgresGenresOnGamesRepo struct {
	DB *sql.DB
}

func NewPostgresGenresOnGamesRepo(DB *sql.DB) *PostgresGenresOnGamesRepo {
	return &PostgresGenresOnGamesRepo{
		DB: DB,
	}
}

func (p *PostgresGenresOnGamesRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS genresOnGames(
        id SERIAL PRIMARY KEY ,
		gameId UUID REFERENCES games(id),
		genreId UUID REFERENCES genres(id)
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (t *TestGenresOnGamesRepo) CreateGenresOnGames(genresOnGames models.GenresOnGames) error {
	if _, ok := t.GenresOnGames[t.counter]; ok {
		return ErrDublicateIDModel
	}
	t.GenresOnGames[genresOnGames.ID] = &genresOnGames
	t.counter++
	return nil
}

func (p *PostgresGenresOnGamesRepo) CreateGenresOnGames(genresOnGames models.GenresOnGames) error {
	_, err := p.DB.Exec("INSERT INTO genresOnGames(gameId, genreId) values($1, $2)", genresOnGames.GameId, genresOnGames.GenreId)
	return err
}

type TestPlatformsOnGamesRepo struct {
	PlatformsOnGames map[uint]*models.PlatformsOnGames
	counter          uint
}

func NewTestPlatformsOnGamesRepo() *TestPlatformsOnGamesRepo {
	return &TestPlatformsOnGamesRepo{
		PlatformsOnGames: map[uint]*models.PlatformsOnGames{},
	}
}

type PostgresPlatformsOnGamesRepo struct {
	DB *sql.DB
}

func NewPostgresPlatformsOnGamesRepo(DB *sql.DB) *PostgresPlatformsOnGamesRepo {
	return &PostgresPlatformsOnGamesRepo{
		DB: DB,
	}
}

func (p *PostgresPlatformsOnGamesRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS platformsOnGames(
        id SERIAL PRIMARY KEY ,
		gameId UUID REFERENCES games(id),
		platformId UUID REFERENCES platforms(id)
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (t *TestPlatformsOnGamesRepo) CreatePlatformsOnGames(platformsOnGames models.PlatformsOnGames) error {
	if _, ok := t.PlatformsOnGames[t.counter]; ok {
		return ErrDublicateIDModel
	}
	t.PlatformsOnGames[platformsOnGames.ID] = &platformsOnGames
	t.counter++
	return nil
}

func (p *PostgresPlatformsOnGamesRepo) CreatePlatformsOnGames(platformsOnGames models.PlatformsOnGames) error {
	_, err := p.DB.Exec("INSERT INTO platformsOnGames(gameId, platformId) values($1, $2)", platformsOnGames.GameId, platformsOnGames.PlatformId)
	return err
}
