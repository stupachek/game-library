//go:build integration_test

package publisher

import (
	"crypto/ed25519"
	"game-library/domens/models"
	"game-library/domens/repository/database"
	"game-library/domens/repository/publisher_repo"
	"game-library/domens/service/jwt"
	"testing"

	"github.com/google/uuid"
)

func TestCreateGetpublishers(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	repo := publisher_repo.NewPostgresPublisherRepo(DB)
	service := NewPublisherService(repo)
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	t.Run("create publisher, get,  get list", func(t *testing.T) {
		//create publisher
		test, err := service.CreatePublisher(models.PublisherModel{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create dublicate
		_, err = service.CreatePublisher(models.PublisherModel{
			Name: "test",
		})
		if err.Error() != "pq: duplicate key value violates unique constraint \"publishers_name_key\"" {
			t.Fatalf("expected %v, got %v", "pq: duplicate key value violates unique constraint \"publishers_name_key\"", err)
		}

		//create publisher
		_, err = service.CreatePublisher(models.PublisherModel{
			Name: "test1",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create publisher
		_, err = service.CreatePublisher(models.PublisherModel{
			Name: "test2",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//get list publishers
		publishers, err := service.GetPublishersList()
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if len(publishers) != 3 {
			t.Fatalf("expected %v, got %v", 3, len(publishers))
		}

		//get publisher
		publisherTest, err := service.GetPublisher(test.ID.String())
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if publisherTest != test {
			t.Fatalf("expected %v, got %v", test, publisherTest)
		}

		//get unknown publisher
		publisherUnknown, err := service.GetPublisher(uuid.NewString())
		if err.Error() != "sql: no rows in result set" {
			t.Fatalf("expected %v, got %v", "sql: no rows in result set", err)
		}
		if publisherUnknown != (models.Publisher{}) {
			t.Fatalf("expected empty publisher")
		}

		//get publisher with wrong uuid
		publisherUnknown, err = service.GetPublisher("error")
		if err != ErrPublisherId {
			t.Fatalf("expected %v, got %v", ErrPublisherId, err)
		}
		if publisherUnknown != (models.Publisher{}) {
			t.Fatalf("expected empty publisher")
		}

	})
}

func TestCreateUpdateGet(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	repo := publisher_repo.NewPostgresPublisherRepo(DB)
	service := NewPublisherService(repo)
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	t.Run("create publisher, get,  get list", func(t *testing.T) {
		//create publisher
		test, err := service.CreatePublisher(models.PublisherModel{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//update publisher
		_, err = service.UpdatePublisher(test.ID.String(), models.PublisherModel{
			Name: "test_update",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//update unknown publisher
		publisherUnknown, err := service.UpdatePublisher(uuid.NewString(), models.PublisherModel{
			Name: "error",
		})
		if err != publisher_repo.ErrUpdateFailed {
			t.Fatalf("expected %v, got %v", publisher_repo.ErrUpdateFailed, err)
		}
		if publisherUnknown != (models.Publisher{}) {
			t.Fatalf("expected empty publisher")
		}

		//get publisher with wrong uuid
		publisherUnknown, err = service.UpdatePublisher("error", models.PublisherModel{
			Name: "error",
		})
		if err != ErrPublisherId {
			t.Fatalf("expected %v, got %v", ErrPublisherId, err)
		}
		if publisherUnknown != (models.Publisher{}) {
			t.Fatalf("expected empty publisher")
		}

		test, err = service.GetPublisher(test.ID.String())
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if test.Name != "test_update" {
			t.Fatalf("expected %v, got %v", "test_update", test.Name)
		}

	})
}

func TestCreateDeleteGetList(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	repo := publisher_repo.NewPostgresPublisherRepo(DB)
	service := NewPublisherService(repo)
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	t.Run("create publisher, get,  get list", func(t *testing.T) {
		//create publisher
		test1, err := service.CreatePublisher(models.PublisherModel{
			Name: "test1",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		//create publisher
		_, err = service.CreatePublisher(models.PublisherModel{
			Name: "test2",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//get list publishers
		publishers, err := service.GetPublishersList()
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if len(publishers) != 2 {
			t.Fatalf("expected %v, got %v", 2, len(publishers))
		}

		//delete publisher
		err = service.DeletePublisher(test1.ID.String())
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//delete unknown publisher
		err = service.DeletePublisher(uuid.NewString())
		if err.Error() != "sql: no rows in result set" {
			t.Fatalf("expected %v, got %v", "sql: no rows in result set", err)
		}

		//delete publisher with wrong uuid
		err = service.DeletePublisher("error")
		if err != ErrPublisherId {
			t.Fatalf("expected %v, got %v", ErrPublisherId, err)
		}

		//get list publishers
		publishers, err = service.GetPublishersList()
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if len(publishers) != 1 {
			t.Fatalf("expected %v, got %v", 1, len(publishers))
		}

	})
}
