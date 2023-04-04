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
	userService := service.NewUserService(userRepo)
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
			users.Use(middleware.CheckPermissions(&userHandler))
			users.PATCH("/:id", userHandler.ChangerRole)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}
	return r
}
