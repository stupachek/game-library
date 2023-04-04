package middleware

import (
	"game-library/domens/service"
	"game-library/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrUnauthenticated = gin.H{"error": "unauthenticated"}

func Auth(u *handler.UserHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		id, err := service.ValidateToken(token)
		ctx.Set("UserID", id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrUnauthenticated)
			return
		}
		ctx.Next()
	}
}
