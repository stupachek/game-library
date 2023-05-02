package api

import (
	"crypto/ed25519"
	"game-library/domens/repository"
	"game-library/domens/service"
	"game-library/handler"
	"game-library/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//TODO: move init repo and servise to main
	DB := repository.ConnectDataBase()
	userRepo := repository.NewPostgresUserRepo(DB)
	if err := userRepo.Migrate(); err != nil {
		log.Fatalf("%v", err)
	}
	gameRepo := repository.NewPostgresGameRepo(DB)
	if err := gameRepo.Migrate(); err != nil {
		log.Fatalf("%v", err)
	}
	publisherRepo := repository.NewPublisherRepo()
	platformRepo := repository.NewPostgresPlatformRepo(DB)
	if err := platformRepo.Migrate(); err != nil {
		log.Fatalf("%v", err)
	}
	genreRepo := repository.NewPostgresGenreRepo(DB)
	if err := genreRepo.Migrate(); err != nil {
		log.Fatalf("%v", err)
	}
	genreOnGameRepo := repository.NewPostgresGenresOnGamesRepo(DB)
	if err := genreOnGameRepo.Migrate(); err != nil {
		log.Fatalf("%v", err)
	}
	platformOnGameRepo := repository.NewPostgresPlatformsOnGamesRepo(DB)
	if err := platformOnGameRepo.Migrate(); err != nil {
		log.Fatalf("%v", err)
	}
	userService := service.NewUserService(userRepo)
	gameService := service.NewGameService(gameRepo, genreOnGameRepo, platformOnGameRepo)
	publisherService := service.NewPublisherService(publisherRepo)
	platformService := service.NewPlatformService(platformRepo)
	genreService := service.NewGenreService(genreRepo)
	err := userService.SetupAdmin()
	if err != nil {
		log.Fatal(err)
	}
	service.Public, service.Private, err = ed25519.GenerateKey(nil)
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
	users.Use(middleware.Auth())
	{
		users.GET("/me", userHandler.GetUser)
		users.GET("/:id", userHandler.GetUser)
		users.GET("", userHandler.GetUsers)
		{
			users.Use(middleware.CheckIfAdmin(&userHandler))
			users.PATCH("/:id", userHandler.ChangerRole)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	gameHandler := handler.NewGameHandler(gameService, publisherService)
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
			platforms.POST("", platformHandler.CreatePlatform)
		}
	}

	genreHandler := handler.NewGenreHandler(genreService)
	genres := r.Group("/genres")
	{
		genres.GET("", genreHandler.GetGenresList)
		{
			genres.Use(middleware.Auth())
			genres.POST("", genreHandler.CreateGenre)
		}
	}
	return r
}
