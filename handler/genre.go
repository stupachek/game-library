package handler

import (
	"game-library/domens/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (g *GenreHandler) CreateGenre(c *gin.Context) {
	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	createdGenre, err := g.GenreService.CreateGenre(genre)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Genre is successfully created", "data": createdGenre})
}

func (g *GenreHandler) GetGenresList(c *gin.Context) {
	genres := g.GenreService.GetGenresList()
	c.JSON(http.StatusOK, gin.H{"data": genres})
}
