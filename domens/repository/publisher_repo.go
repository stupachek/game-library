package repository

import (
	"database/sql"
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrPublisherNotFound = errors.New("publisher does not exist")
)

type IPublisherRepo interface {
	GetPublishersList() ([]models.Publisher, error)
	CreatePublisher(pub models.Publisher) error
	Delete(id uuid.UUID) error
}

type TestPublisherRepo struct {
	Publishers map[uuid.UUID]*models.Publisher
	sync.Mutex
}

func NewPublisherRepo() *TestPublisherRepo {
	return &TestPublisherRepo{
		Publishers: make(map[uuid.UUID]*models.Publisher),
	}
}

type PostgresPublisherRepo struct {
	DB *sql.DB
}

func NewPostgresPublisherRepo(DB *sql.DB) *PostgresPublisherRepo {
	return &PostgresPublisherRepo{
		DB: DB,
	}
}

func (p *PostgresPublisherRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS publishers(
        id UUID PRIMARY KEY,
    	name VARCHAR NOT NULL UNIQUE
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (t *TestPublisherRepo) GetPublishersList() ([]models.Publisher, error) {
	publishers := make([]models.Publisher, 0)
	for _, user := range t.Publishers {
		publishers = append(publishers, *user)
	}
	return publishers, nil
}

func (p *PostgresPublisherRepo) GetPublishersList() ([]models.Publisher, error) {
	rows, err := p.DB.Query("SELECT id, name FROM publishers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var publishers []models.Publisher
	for rows.Next() {
		var publisher models.Publisher
		if err := rows.Scan(&publisher.ID, &publisher.Name); err != nil {
			return nil, err
		}
		publishers = append(publishers, publisher)
	}
	return publishers, nil
}

func (t *TestPublisherRepo) CreatePublisher(publisher models.Publisher) error {
	t.Publishers[publisher.ID] = &publisher
	return nil
}

func (p *PostgresPublisherRepo) CreatePublisher(pub models.Publisher) error {
	_, err := p.DB.Exec("INSERT INTO publishers(id, name) values($1, $2)", pub.ID, pub.Name)
	return err
}

func (t *TestPublisherRepo) GetPublisherById(id uuid.UUID) (*models.Publisher, error) {
	publisher, ok := t.Publishers[id]
	if !ok {
		return &models.Publisher{}, ErrPublisherNotFound
	}
	return publisher, nil
}

func (p *PostgresPublisherRepo) GetPublisherById(id uuid.UUID) (*models.Publisher, error) {
	row := p.DB.QueryRow("SELECT id, name FROM publishers WHERE id = $1", id)
	var pub models.Publisher
	if err := row.Scan(&pub.ID, &pub.Name); err != nil {
		return &models.Publisher{}, err
	}
	return &pub, nil
}

func (t *TestPublisherRepo) UpdatePublisher(id uuid.UUID, publisher models.PublisherModel) error {
	publisherById, err := t.GetPublisherById(id)
	if err != nil {
		return err
	}
	publisherById.Name = publisher.Name
	return nil
}

func (p *PostgresPublisherRepo) UpdatePublisher(id uuid.UUID, publisher models.PublisherModel) error {
	res, err := p.DB.Exec("UPDATE publishers SET name = $1 WHERE  id = $2", publisher.Name, id)
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

func (t *TestPublisherRepo) Delete(id uuid.UUID) error {
	delete(t.Publishers, id)
	return nil
}

func (p *PostgresPublisherRepo) Delete(id uuid.UUID) error {
	res, err := p.DB.Exec("DELETE FROM publishers WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}
	return nil
}
