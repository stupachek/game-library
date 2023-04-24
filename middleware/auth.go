package middleware

import (
	"game-library/domens/models"
	"game-library/domens/service"
	"game-library/handler"
	"log"
	"net/http"
	"os"

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

func CheckIfAdmin(u *handler.UserHandler) gin.HandlerFunc {
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

func CheckIfManager(u *handler.UserHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := handler.GetId(ctx)
		user, err := u.UserService.GetUser(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		switch user.Role {
		case models.ADMIN:
			ctx.Next()
		case models.MANAGER:
			ctx.Next()
		default:
			ctx.AbortWithStatusJSON(http.StatusForbidden, ErrPermissionDenied)
			return
		}
	}
}

func DeleteFile(g *handler.GameHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if err := ctx.Errors.Last(); err != nil {
			dst, ok := ctx.Get("dst")
			if !ok {
				log.Fatal("can't find dst")
			}
			if err := os.Remove(dst.(string)); err != nil {
				log.Fatalf("error when deleting uploaded file: %s\n", err)
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
}
