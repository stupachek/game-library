//go:build unit_test

package game

import (
	"game-library/domens/models"
	"game-library/domens/repository/game_repo"
	"testing"

	"github.com/google/uuid"
)

func TestCreateGame(t *testing.T) {
	testCases := []struct {
		description   string
		game          models.Game
		genres        []models.Genre
		plaforms      []models.Platform
		expectedError error
	}{
		{
			description: "succes",
			game: models.Game{
				PublisherId:    [16]byte{1},
				Title:          "test",
				Description:    "test",
				ImageLink:      "library/test",
				AgeRestriction: 12,
				ReleaseYear:    2012,
			},
			genres:        make([]models.Genre, 0),
			plaforms:      make([]models.Platform, 0),
			expectedError: nil,
		},
		{
			description: "dublicate",
			game: models.Game{
				PublisherId:    [16]byte{1},
				Title:          "test",
				Description:    "test",
				ImageLink:      "library/test",
				AgeRestriction: 12,
				ReleaseYear:    2012,
			},
			genres:        make([]models.Genre, 0),
			plaforms:      make([]models.Platform, 0),
			expectedError: game_repo.ErrDublicateTitle,
		},
	}
	testCases[0].genres = append(testCases[0].genres, models.Genre{
		ID:   [16]byte{12},
		Name: "smt_g",
	})
	testCases[0].plaforms = append(testCases[0].plaforms, models.Platform{
		ID:   [16]byte{111},
		Name: "smt_p1",
	}, models.Platform{
		ID:   [16]byte{112},
		Name: "smt_p2",
	})
	repo := game_repo.NewGameRepo()
	repo_g := game_repo.NewTestGenresOnGamesRepo()
	repo_p := game_repo.NewTestPlatformsOnGamesRepo()
	servive := NewGameService(repo, repo_g, repo_p)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := servive.CreateGame(tc.game, tc.genres, tc.plaforms)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
		})

	}

}

func TestGetGame(t *testing.T) {
	testCases := []struct {
		description   string
		idStr         string
		expectedError error
	}{
		{
			description:   "succes",
			idStr:         uuid.UUID{111}.String(),
			expectedError: nil,
		},
		{
			description:   "can't parse id",
			idStr:         "error",
			expectedError: ErrParseId,
		},
		{
			description:   "unknown id",
			idStr:         uuid.UUID{9}.String(),
			expectedError: ErrUnknownId,
		},
	}

	repo := game_repo.NewGameRepo()
	repo_g := game_repo.NewTestGenresOnGamesRepo()
	repo_p := game_repo.NewTestPlatformsOnGamesRepo()
	servive := NewGameService(repo, repo_g, repo_p)
	repo.Setup()
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := servive.GetGame(tc.idStr)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
		})

	}

}
