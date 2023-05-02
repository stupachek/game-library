package handler

import (
	"errors"
	"fmt"
	"game-library/domens/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrFileExists = gin.H{"error": "file with the name already exists"}

type InputGame struct {
	Title          string                `form:"title" binding:"required"`
	Description    string                `form:"description"`
	File           *multipart.FileHeader `form:"file" binding:"required"`
	PublisherId    string                `form:"publisherId" binding:"required"`
	AgeRestriction int                   `form:"ageRestriction"`
	ReleaseYear    int                   `form:"releaseYear"`
	Genres         []string              `form:"genres"`
	Platforms      []string              `form:"platforms"`
}

func (g *GameHandler) GetGamesList(c *gin.Context) {
	games, err := g.GameService.GetGamesList()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
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
	inputGame := InputGame{}
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
	inputGame := InputGame{}
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
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	genres, err := g.fromStringToGenres(inputGame.Genres)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	plaforms, err := g.fromStringToPlatform(inputGame.Platforms)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	publisherId, err := uuid.Parse(inputGame.PublisherId)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("can't parse publisherId id"))
		return
	}
	game := models.NewGame(publisherId, inputGame.Title, inputGame.Description, dst, inputGame.AgeRestriction, inputGame.ReleaseYear)
	err = g.GameService.CreateGame(game, genres, plaforms)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game is successfully created", "data": gin.H{"link": dst}})
}
