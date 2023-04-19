package handler

import (
	"fmt"
	"game-library/domens/models"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (g *GameHandler) GetGamesList(c *gin.Context) {
	games := g.GameService.GetGamesList()
	c.JSON(http.StatusOK, gin.H{"data": games})
}

func (g *GameHandler) CreateGame(c *gin.Context) {
	inputGame := models.InputGame{}
	if err := c.ShouldBind(&inputGame); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := inputGame.File
	dst := fmt.Sprintf("library/%s", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := g.PublisherService.GetPublisher(inputGame.PublisherId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	genres := make([]models.Genre, len(inputGame.Genres))
	for i, genre := range inputGame.Genres {
		g, err := g.GenreService.GetGenre(genre)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "genre": genre})
			return
		}
		genres[i] = g
	}
	plaforms := make([]models.Platform, len(inputGame.Platforms))
	for i, platform := range inputGame.Platforms {
		p, err := g.PlatformService.GetPlatform(platform)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "platform": platform})
			return
		}
		plaforms[i] = p
	}
	g.GameService.CreateGame(inputGame, dst, genres, plaforms)

	c.JSON(http.StatusOK, gin.H{"message": "Game is successfully created", "data": gin.H{"link": dst}})
}
