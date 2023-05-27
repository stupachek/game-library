package api

import (
	"crypto/ed25519"
	"database/sql"
	"game-library/domens/repository/game_repo"
	"game-library/domens/repository/genre_repo"
	"game-library/domens/repository/platform_repo"
	"game-library/domens/repository/publisher_repo"
	"game-library/domens/repository/user_repo"
	"game-library/domens/service/game"
	"game-library/domens/service/genre"
	"game-library/domens/service/jwt"
	"game-library/domens/service/platform"
	"game-library/domens/service/publisher"
	"game-library/domens/service/user"
	"game-library/handler"
	"game-library/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(DB *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(middleware.CORSMiddleware())
	//TODO: move init repo and servise to main
	userRepo := user_repo.NewPostgresUserRepo(DB)
	gameRepo := game_repo.NewPostgresGameRepo(DB)
	publisherRepo := publisher_repo.NewPostgresPublisherRepo(DB)
	platformRepo := platform_repo.NewPostgresPlatformRepo(DB)
	genreRepo := genre_repo.NewPostgresGenreRepo(DB)
	genreOnGameRepo := game_repo.NewPostgresGenresOnGamesRepo(DB)
	platformOnGameRepo := game_repo.NewPostgresPlatformsOnGamesRepo(DB)
	userService := user.NewUserService(userRepo)
	gameService := game.NewGameService(gameRepo, genreOnGameRepo, platformOnGameRepo)
	publisherService := publisher.NewPublisherService(publisherRepo)
	platformService := platform.NewPlatformService(platformRepo)
	genreService := genre.NewGenreService(genreRepo)
	var err error
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}
	userHandler := handler.NewUserHadler(userService)
	auth := r.Group("/auth")
	{
		auth.POST("/signup", userHandler.Register)
		auth.POST("/signin", userHandler.Login)
	}
	users := r.Group("/users")
	{
		users.GET("/:id", userHandler.GetUser)
		users.GET("", userHandler.GetUsers)
		users.Use(middleware.Auth())
		users.GET("/me", userHandler.GetUser)
		{
			users.Use(middleware.CheckIfAdmin(&userHandler))
			users.PATCH("/:id", userHandler.ChangerRole)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	gameHandler := handler.NewGameHandler(gameService, publisherService, genreService, platformService)
	games := r.Group("/games")
	games.Use(middleware.Auth())
	{
		games.GET("", gameHandler.GetGamesList)
		games.GET(":id", gameHandler.GetGame)
		games.Use(middleware.CheckIfManager(&userHandler))
		games.Use(middleware.DeleteFile(&gameHandler))
		games.POST("", gameHandler.CreateGame)
	}

	publisherHandler := handler.NewPublisherHandler(publisherService)
	publishers := r.Group("/publishers")
	{
		publishers.GET("", publisherHandler.GetPlatformsList)
		publishers.GET("/:id", publisherHandler.GetPublisher)
		{
			publishers.Use(middleware.Auth())
			publishers.Use(middleware.CheckIfManager(&userHandler))
			publishers.POST("", publisherHandler.CreatePublisher)
			publishers.PATCH("/:id", publisherHandler.UpdatePublisher)
			publishers.DELETE("/:id", publisherHandler.DeletePublisher)
		}
	}

	platformHandler := handler.NewPlatformHandler(platformService)
	platforms := r.Group("/platforms")
	{
		platforms.GET("", platformHandler.GetPlatformsList)
		{
			platforms.Use(middleware.Auth())
			platforms.Use(middleware.CheckIfManager(&userHandler))
			platforms.POST("", platformHandler.CreatePlatform)
		}
	}

	genreHandler := handler.NewGenreHandler(genreService)
	genres := r.Group("/genres")
	{
		genres.GET("", genreHandler.GetGenresList)
		{
			genres.Use(middleware.Auth())
			genres.Use(middleware.CheckIfManager(&userHandler))
			genres.POST("", genreHandler.CreateGenre)
		}
	}
	return r
}
