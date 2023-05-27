//go:build unit_test

package genre

import (
	"game-library/domens/models"
	"game-library/domens/repository/genre_repo"
	"testing"

	"github.com/google/uuid"
)

func TestWCreateSucces(t *testing.T) {
	repo := genre_repo.NewGenreRepo()
	service := NewGenreService(repo)
	t.Run("success create genre", func(t *testing.T) {
		genre, err := service.CreateGenre(models.Genre{
			Name: "test",
		})
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if genre.ID == (uuid.UUID{}) {
			t.Fatal("got empty uuid")
		}
	})
}

func TestGetSucces(t *testing.T) {
	testCases := []struct {
		description       string
		genre             string
		expectedError     error
		expectedGenreName string
	}{
		{
			description:       "success",
			genre:             "test",
			expectedError:     nil,
			expectedGenreName: "test",
		},
		{
			description:       "failed",
			genre:             "error",
			expectedError:     ErrUnknownGenre,
			expectedGenreName: "",
		},
	}
	repo := genre_repo.NewGenreRepo()
	repo.Setup()
	service := NewGenreService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			genre, err := service.GetGenre(tc.genre)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
			if genre.Name != tc.expectedGenreName {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedGenreName, genre.Name)
			}
		})

	}
}

func TestGetListGenres(t *testing.T) {
	repo := genre_repo.NewGenreRepo()
	repo.Setup()
	service := NewGenreService(repo)
	t.Run("success get genres", func(t *testing.T) {
		genres, err := service.GetGenresList()
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if len(genres) != 1 {
			t.Fatalf("expected %v, got %v", 1, len(genres))
		}
	})

}
