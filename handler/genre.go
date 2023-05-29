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

func (p *GenreHandler) GetGenre(c *gin.Context) {
	id := c.Param("id")

	genre, err := p.GenreService.GetGenre(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": genre})
}

func (p *GenreHandler) UpdateGenre(c *gin.Context) {
	var genreModel models.Genre
	if err := c.ShouldBindJSON(&genreModel); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idStr := c.Param("id")
	genre, err := p.GenreService.UpdateGenre(idStr, genreModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Genre is successfully updated", "data": genre})
}

func (p *GenreHandler) DeleteGenre(c *gin.Context) {
	id := c.Param("id")
	err := p.GenreService.DeleteGenre(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Genre is successfully deleted"})
}
