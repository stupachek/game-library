//go:build unit_test

package publisher

import (
	"game-library/domens/models"
	"game-library/domens/repository/publisher_repo"
	"testing"

	"github.com/google/uuid"
)

func TestCreatePablisherSucces(t *testing.T) {
	repo := publisher_repo.NewPublisherRepo()
	service := NewPublisherService(repo)
	t.Run("succes create publisher", func(t *testing.T) {
		publisher, err := service.CreatePublisher(models.PublisherModel{
			Name: "test",
		})
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if publisher.ID == (uuid.UUID{}) {
			t.Fatal("got empty uuid")
		}
	})
}

func TestGetPublisher(t *testing.T) {
	testCases := []struct {
		description           string
		publisherId           string
		expectedError         error
		expectedPublisherName string
	}{
		{
			description:           "success",
			publisherId:           uuid.UUID{111}.String(),
			expectedError:         nil,
			expectedPublisherName: "test",
		},
		{
			description:           "unknown publisher",
			publisherId:           uuid.UUID{000}.String(),
			expectedError:         publisher_repo.ErrPublisherNotFound,
			expectedPublisherName: "",
		},
		{
			description:           "can't parse",
			publisherId:           "error",
			expectedError:         ErrPublisherId,
			expectedPublisherName: "",
		},
	}
	repo := publisher_repo.NewPublisherRepo()
	repo.Setup()
	service := NewPublisherService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			publisher, err := service.GetPublisher(tc.publisherId)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
			if publisher.Name != tc.expectedPublisherName {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedPublisherName, publisher.Name)
			}
		})

	}
}

func TestGetListPublishers(t *testing.T) {
	repo := publisher_repo.NewPublisherRepo()
	repo.Setup()
	service := NewPublisherService(repo)
	t.Run("success get publishers", func(t *testing.T) {
		publishers, err := service.GetPublishersList()
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if len(publishers) != 1 {
			t.Fatalf("expected %v, got %v", 1, len(publishers))
		}
	})

}

func TestDeletePublisher(t *testing.T) {
	testCases := []struct {
		description   string
		idStr         string
		expectedError error
	}{
		{
			description:   "success",
			idStr:         uuid.UUID{111}.String(),
			expectedError: nil,
		},
		{
			description:   "can't parse",
			idStr:         "error",
			expectedError: ErrPublisherId,
		},
	}
	repo := publisher_repo.NewPublisherRepo()
	repo.Setup()
	servive := NewPublisherService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := servive.DeletePublisher(tc.idStr)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
		})

	}

}

func TestUpdatePublisher(t *testing.T) {
	testCases := []struct {
		description   string
		idStr         string
		publisher     models.PublisherModel
		expectedError error
	}{
		{
			description: "success",
			idStr:       uuid.UUID{111}.String(),
			publisher: models.PublisherModel{
				Name: "test_update",
			},
			expectedError: nil,
		},
		{
			description:   "can't parse",
			idStr:         "error",
			expectedError: ErrPublisherId,
		},
		{
			description:   "unknown publisher",
			idStr:         uuid.UUID{000}.String(),
			expectedError: publisher_repo.ErrPublisherNotFound,
		},
	}
	repo := publisher_repo.NewPublisherRepo()
	repo.Setup()
	servive := NewPublisherService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			publisher, err := servive.UpdatePublisher(tc.idStr, tc.publisher)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
			if publisher.Name != tc.publisher.Name {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.publisher, publisher)
			}
		})

	}

}
