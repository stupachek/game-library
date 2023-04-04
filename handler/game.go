package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (g *GameHandler) GetGamesList(c *gin.Context) {
	games := g.GameService.GetGamesList()
	c.JSON(http.StatusOK, gin.H{"data": games})
}
