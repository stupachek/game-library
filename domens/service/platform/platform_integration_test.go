//go:build integration_test

package platform

import (
	"crypto/ed25519"
	"game-library/domens/models"
	"game-library/domens/repository/database"
	"game-library/domens/repository/platform_repo"
	"game-library/domens/service/jwt"
	"testing"
)

func TestCreateGetPlatforms(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	repo := platform_repo.NewPostgresPlatformRepo(DB)
	service := NewPlatformService(repo)
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	t.Run("create platform, get platform, list platforms", func(t *testing.T) {
		//create platform
		test1, err := service.CreatePlatform(models.Platform{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create dublicate
		_, err = service.CreatePlatform(models.Platform{
			Name: "test",
		})
		if err.Error() != "pq: duplicate key value violates unique constraint \"platforms_name_key\"" {
			t.Fatalf("expected %v, got %v", "pq: duplicate key value violates unique constraint \"platforms_name_key\"", err)
		}

		//create platform
		_, err = service.CreatePlatform(models.Platform{
			Name: "test1",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create platform
		_, err = service.CreatePlatform(models.Platform{
			Name: "test2",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//get list platforms
		platforms, err := service.GetPlatformsList()
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if len(platforms) != 3 {
			t.Fatalf("expected %v, got %v", 3, len(platforms))
		}

		//get platform
		platformTest, err := service.GetPlatform("test")
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if platformTest != test1 {
			t.Fatalf("expected %v, got %v", test1, platformTest)
		}

	})
}
