package handler

import (
	"fmt"
	"game-library/domens/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var ErrFileExists = gin.H{"error": "file with the name already exists"}

func (g *GameHandler) GetGamesList(c *gin.Context) {
	games := g.GameService.GetGamesList()
	c.JSON(http.StatusOK, gin.H{"data": games})
}

func (g *GameHandler) GetGame(c *gin.Context) {
	idStr := c.Param("id")
	game, err := g.GameService.GetGame(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": game})
}

func (g *GameHandler) UpdateGame(c *gin.Context) {
	inputGame := models.InputGame{}
	if err := c.ShouldBind(&inputGame); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	file := inputGame.File
	dst := fmt.Sprintf("library/%s", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}

func (g *GameHandler) fromStringToGenres(stringGenres []string) ([]models.Genre, error) {
	genres := make([]models.Genre, len(stringGenres))
	for i, genre := range stringGenres {
		g, err := g.GenreService.GetGenre(genre)
		if err != nil {
			return []models.Genre{}, err
		}
		genres[i] = g
	}
	return genres, nil
}

func (g *GameHandler) fromStringToPlatform(stringPlatform []string) ([]models.Platform, error) {
	plaforms := make([]models.Platform, len(stringPlatform))
	for i, platform := range stringPlatform {
		p, err := g.PlatformService.GetPlatform(platform)
		if err != nil {
			return []models.Platform{}, err
		}
		plaforms[i] = p
	}
	return plaforms, nil
}

func (g *GameHandler) CreateGame(c *gin.Context) {
	inputGame := models.InputGame{}
	if err := c.ShouldBind(&inputGame); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := inputGame.File
	dst := fmt.Sprintf("library/%s", filepath.Base(file.Filename))
	if _, err := os.Stat(dst); err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrFileExists)
		return
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Set("dst", dst)
	if _, err := g.PublisherService.GetPublisher(inputGame.PublisherId); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	genres, err := g.fromStringToGenres(inputGame.Genres)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	plaforms, err := g.fromStringToPlatform(inputGame.Platforms)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = g.GameService.CreateGame(inputGame, dst, genres, plaforms)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game is successfully created", "data": gin.H{"link": dst}})
}
