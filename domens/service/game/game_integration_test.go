//go:build integration_test

package game

import (
	"game-library/domens/models"
	"game-library/domens/repository/database"
	"game-library/domens/repository/game_repo"
	"game-library/domens/repository/genre_repo"
	"game-library/domens/repository/platform_repo"
	"game-library/domens/repository/publisher_repo"
	"game-library/domens/service/genre"
	"game-library/domens/service/platform"
	"game-library/domens/service/publisher"
	"testing"
)

func TestCreateGet(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	game := game_repo.NewPostgresGameRepo(DB)
	genreOnGame := game_repo.NewPostgresGenresOnGamesRepo(DB)
	platformOnGenre := game_repo.NewPostgresPlatformsOnGamesRepo(DB)
	gameService := NewGameService(game, genreOnGame, platformOnGenre)

	publisherRepo := publisher_repo.NewPostgresPublisherRepo(DB)
	publisherService := publisher.NewPublisherService(publisherRepo)

	genreRepo := genre_repo.NewPostgresGenreRepo(DB)
	genreService := genre.NewGenreService(genreRepo)

	platformRepo := platform_repo.NewPostgresPlatformRepo(DB)
	platformService := platform.NewPlatformService(platformRepo)
	t.Run("create, get", func(t *testing.T) {
		//create publisher
		publisherTest, err := publisherService.CreatePublisher(models.PublisherModel{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create platform1
		platformTest1, err := platformService.CreatePlatform(models.Platform{
			Name: "test1",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create platform2
		platformTest2, err := platformService.CreatePlatform(models.Platform{
			Name: "test2",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create genre 1
		genreTest1, err := genreService.CreateGenre(models.Genre{
			Name: "test1",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create genre 2
		genreTest2, err := genreService.CreateGenre(models.Genre{
			Name: "test2",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		//create genre 3
		genreTest3, err := genreService.CreateGenre(models.Genre{
			Name: "test3",
		})
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}

		platforms := make([]models.Platform, 0)
		platforms = append(platforms, platformTest1, platformTest2)

		genres := make([]models.Genre, 0)
		genres = append(genres, genreTest1, genreTest2, genreTest3)

		//create game
		game1, err := gameService.CreateGame(models.Game{
			PublisherId:    publisherTest.ID,
			Title:          "test1",
			Description:    "111",
			ImageLink:      "111",
			AgeRestriction: 111,
			ReleaseYear:    111,
		}, genres, platforms)
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}

		//create game
		_, err = gameService.CreateGame(models.Game{
			PublisherId:    publisherTest.ID,
			Title:          "test2",
			Description:    "111",
			ImageLink:      "111",
			AgeRestriction: 111,
			ReleaseYear:    111,
		}, genres, platforms)
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}

		//get game
		game1Get, err := gameService.GetGame(game1.ID.String())
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if game1Get.ID != game1.ID.String() {
			t.Fatalf(" expected %v, got %v", game1, game1Get)

		}

		//get games
		games, err := gameService.GetGamesList(models.QueryParams{
			Skip: 0,
			Take: 99999,
		})
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if len(games) != 2 {
			t.Fatalf(" expected %v, got %v", 2, len(games))

		}

	})
}
