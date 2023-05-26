//go:build integration_test

package game

//	TODO:after publisher tests
// func TestCreateGet(t *testing.T) {
// 	DB := database.ConnectDataBase()
// 	err := database.ClearData(DB)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	game := game_repo.NewPostgresGameRepo(DB)
// 	genre := game_repo.NewPostgresGenresOnGamesRepo(DB)
// 	platform := game_repo.NewPostgresPlatformsOnGamesRepo(DB)
// 	service := NewGameService(game, genre, platform)
// 	t.Run("create, get", func(t *testing.T) {
// 		err := service.CreateGame(models.Game{
// 			PublisherId:    [16]byte{111},
// 			Title:          "111",
// 			Description:    "111",
// 			ImageLink:      "111",
// 			AgeRestriction: 111,
// 			ReleaseYear:    111,
// 		}, make([]models.Genre, 1), make([]models.Platform, 1))
// 		if err != nil {
// 			t.Fatalf(" expected %v, got %v", nil, err)
// 		}
// 	})
// }
