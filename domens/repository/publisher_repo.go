package repository

import (
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrPublisherNotFound = errors.New("publisher does not exist")
)

type TestPublisherRepo struct {
	Publishers map[uuid.UUID]*models.Publisher
	sync.Mutex
}

func NewPublisherRepo() *TestPublisherRepo {
	return &TestPublisherRepo{
		Publishers: make(map[uuid.UUID]*models.Publisher),
	}
}

func (t *TestPublisherRepo) GetPublishersList() []models.Publisher {
	publishers := make([]models.Publisher, 0)
	for _, user := range t.Publishers {
		publishers = append(publishers, *user)
	}
	return publishers

}

func (t *TestPublisherRepo) CreatePublisher(publisher models.Publisher) models.Publisher {
	t.Publishers[publisher.ID] = &publisher
	return publisher
}

func (t *TestPublisherRepo) GetPublisherById(id uuid.UUID) (*models.Publisher, error) {
	publisher, ok := t.Publishers[id]
	if !ok {
		return &models.Publisher{}, ErrPublisherNotFound
	}
	return publisher, nil
}

func (t *TestPublisherRepo) UpdatePublisher(id uuid.UUID, publisher models.PublisherModel) (models.Publisher, error) {
	publisherById, err := t.GetPublisherById(id)
	if err != nil {
		return *publisherById, err
	}
	publisherById.Name = publisher.Name
	return *publisherById, nil
}

func (t *TestPublisherRepo) Delete(id uuid.UUID) {
	delete(t.Publishers, id)
}
