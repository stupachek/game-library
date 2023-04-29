package service

import (
	"errors"
	"game-library/domens/models"

	"github.com/google/uuid"
)

var (
// ErrPublisher = errors.New("error on publisher create")
)

type PublisherService struct {
	PublisherRepo IPublisherRepo
}

func NewPublisherService(repo IPublisherRepo) PublisherService {
	return PublisherService{
		PublisherRepo: repo,
	}
}

func (p *PublisherService) GetPublishersList() []models.Publisher {
	publishers := p.PublisherRepo.GetPublishersList()
	return publishers
}

func (p *PublisherService) CreatePublisher(platformModel models.PublisherModel) models.Publisher {
	id, _ := uuid.NewRandom()
	publisher := models.Publisher{
		ID:   id,
		Name: platformModel.Name,
	}
	createdPublisher := p.PublisherRepo.CreatePublisher(publisher)
	return createdPublisher
}

func (p *PublisherService) UpdatePublisher(idStr string, publisherModel models.PublisherModel) (models.Publisher, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Publisher{}, errors.New("can't parse publisher id")
	}
	publisher, err := p.PublisherRepo.UpdatePublisher(id, publisherModel)
	if err != nil {
		return models.Publisher{}, err
	}
	return publisher, nil
}

func (p *PublisherService) GetPublisher(idStr string) (models.Publisher, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Publisher{}, errors.New("can't parse publisher id")
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
		return errors.New("can't parse publisher id")
	}
	if _, err := p.PublisherRepo.GetPublisherById(id); err != nil {
		return err
	}
	p.PublisherRepo.Delete(id)
	return nil
}
