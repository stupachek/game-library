package publisher

import (
	"errors"
	"game-library/domens/models"

	"github.com/google/uuid"
)

var (
	ErrPublisherId = errors.New("can't parse publisher id")
)

type IPublisherRepo interface {
	GetPublishersList() ([]models.Publisher, error)
	CreatePublisher(pub models.Publisher) error
	GetPublisherById(id uuid.UUID) (*models.Publisher, error)
	UpdatePublisher(id uuid.UUID, publisher models.PublisherModel) error
	Delete(id uuid.UUID) error
}

type PublisherService struct {
	PublisherRepo IPublisherRepo
}

func NewPublisherService(repo IPublisherRepo) PublisherService {
	return PublisherService{
		PublisherRepo: repo,
	}
}

func (p *PublisherService) GetPublishersList() ([]models.Publisher, error) {
	return p.PublisherRepo.GetPublishersList()
}

func (p *PublisherService) CreatePublisher(publisherModel models.PublisherModel) (models.Publisher, error) {
	id, _ := uuid.NewRandom()
	publisher := models.Publisher{
		ID:   id,
		Name: publisherModel.Name,
	}
	err := p.PublisherRepo.CreatePublisher(publisher)
	return publisher, err
}

func (p *PublisherService) UpdatePublisher(idStr string, publisherModel models.PublisherModel) (models.Publisher, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Publisher{}, ErrPublisherId
	}
	err = p.PublisherRepo.UpdatePublisher(id, publisherModel)
	if err != nil {
		return models.Publisher{}, err
	}
	publisher := models.Publisher{
		ID:   id,
		Name: publisherModel.Name,
	}
	return publisher, nil
}

func (p *PublisherService) GetPublisher(idStr string) (models.Publisher, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Publisher{}, ErrPublisherId
	}
	publisher, err := p.PublisherRepo.GetPublisherById(id)
	if err != nil {
		return models.Publisher{}, err
	}
	return *publisher, err
}

func (p *PublisherService) DeletePublisher(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ErrPublisherId
	}
	if _, err := p.PublisherRepo.GetPublisherById(id); err != nil {
		return err
	}
	return p.PublisherRepo.Delete(id)
}
