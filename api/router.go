package api

import (
	"crypto/ed25519"
	"game-library/domens/repository"
	"game-library/domens/service"
	"game-library/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	userRepo := repository.NewUserRepo()
	userService := service.NewUserService(userRepo)
	var err error
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
	return r
}
