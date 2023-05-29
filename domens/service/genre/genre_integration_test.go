//go:build integration_test

package genre

import (
	"crypto/ed25519"
	"game-library/domens/models"
	"game-library/domens/repository/database"
	"game-library/domens/repository/genre_repo"
	"game-library/domens/service/jwt"
	"testing"
)

func TestCreateGetGenres(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	repo := genre_repo.NewPostgresGenreRepo(DB)
	service := NewGenreService(repo)
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	t.Run("create genre, get genre, list genres", func(t *testing.T) {
		//create genre
		test1, err := service.CreateGenre(models.Genre{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create dublicate
		_, err = service.CreateGenre(models.Genre{
			Name: "test",
		})
		if err.Error() != "pq: duplicate key value violates unique constraint \"genres_name_key\"" {
			t.Fatalf("expected %v, got %v", "pq: duplicate key value violates unique constraint \"genres_name_key\"", err)
		}

		//create genre
		_, err = service.CreateGenre(models.Genre{
			Name: "test1",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create genre
		_, err = service.CreateGenre(models.Genre{
			Name: "test2",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//get list genres
		genres, err := service.GetGenresList()
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if len(genres) != 3 {
			t.Fatalf("expected %v, got %v", 3, len(genres))
		}

		//get genre
		genreTest, err := service.GetGenre(test1.ID.String())
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if genreTest != test1 {
			t.Fatalf("expected %v, got %v", test1, genreTest)
		}

	})
}
