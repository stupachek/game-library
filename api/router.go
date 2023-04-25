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
	userRepo := repository.NewUserRepo()
	gameRepo := repository.NewGameRepo()
	publisherRepo := repository.NewPublisherRepo()
	platformRepo := repository.NewPlatformRepo()
	genreRepo := repository.NewGenreRepo()
	userService := service.NewUserService(userRepo)
	gameService := service.NewGameService(gameRepo)
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
		games.POST("", gameHandler.CreateGame)
		games.PATCH(":id")
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
