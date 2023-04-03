package api

import (
	"game-library/domens/repository"
	"game-library/domens/service"
	"game-library/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	userRepo := repository.NewUserRepo()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHadler(userService)
	auth := r.Group("/auth")
	{
		auth.POST("/signup", userHandler.Register)
		auth.POST("/signin", userHandler.Login)
	}
	return r
}

// func (a *App) Run(port string) error {
// 	server := http.Server{Addr: port, Handler: a.Router}
// 	a.Server = server
// 	return a.Server.ListenAndServe()
// }

// func (a *App) Stop() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	if err := a.Server.Shutdown(ctx); err != nil {
// 		return err
// 	}
// 	return nil
// }
