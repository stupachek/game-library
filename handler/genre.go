package handler

import (
	"game-library/domens/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (g *GenreHandler) CreateGenre(c *gin.Context) {
	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't parse input"})
		return
	}
	createdGenre, err := g.GenreService.CreateGenre(genre)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't create genre"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Genre is successfully created", "data": createdGenre})
}

func (g *GenreHandler) GetGenresList(c *gin.Context) {
	genres, err := g.GenreService.GetGenresList()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't get list genres"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": genres})
}
