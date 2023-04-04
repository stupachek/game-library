package middleware

import (
	"game-library/domens/models"
	"game-library/domens/service"
	"game-library/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrPermissionDenied = gin.H{"error": "permission denied"}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		id, err := service.ValidateToken(token)
		ctx.Set("UserID", id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, handler.ErrUnauthenticated)
			return
		}
		ctx.Next()
	}
}

func CheckPermissions(u *handler.UserHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := handler.GetId(ctx)
		user, err := u.UserService.GetUser(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.Role != models.ADMIN {
			ctx.AbortWithStatusJSON(http.StatusForbidden, ErrPermissionDenied)
			return
		}
		ctx.Next()
	}
}
