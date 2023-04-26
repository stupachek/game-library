package handler

import (
	"game-library/domens/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (g *GenreHandler) CreateGenre(c *gin.Context) {
	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	createdGenre := g.GenreService.CreateGenre(genre)
	c.JSON(http.StatusOK, gin.H{"message": "Genre is successfully created", "data": createdGenre})
}

func (g *GenreHandler) GetGenresList(c *gin.Context) {
	genres := g.GenreService.GetGenresList()
	c.JSON(http.StatusOK, gin.H{"data": genres})
}
