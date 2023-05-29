package handler

import (
	"errors"
	"fmt"
	"game-library/domens/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
	Genres         []string              `form:"genres" binding:"required"`
	Platforms      []string              `form:"platforms" binding:"required"`
}

type UpdateGame struct {
	Title          string                `form:"title" binding:"required"`
	Description    string                `form:"description"`
	File           *multipart.FileHeader `form:"file"`
	PublisherId    string                `form:"publisherId" binding:"required"`
	AgeRestriction int                   `form:"ageRestriction"`
	ReleaseYear    int                   `form:"releaseYear"`
	Genres         []string              `form:"genres" binding:"required"`
	Platforms      []string              `form:"platforms" binding:"required"`
}

func (g *GameHandler) GetGamesList(c *gin.Context) {
	query, err := query(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "unknown query params"})
		return
	}
	games, err := g.GameService.GetGamesList(query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't get list games"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": games, "meta": gin.H{"totalCount": len(games)}})
}

func (g *GameHandler) GetGame(c *gin.Context) {
	idStr := c.Param("id")
	game, err := g.GameService.GetGame(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't get game"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": game})
}

func (g *GameHandler) GetImage(c *gin.Context) {
	image := c.Param("impath")
	fileBytes, err := os.ReadFile("library/" + image)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "unknown image"})
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "image/png")
	_, err = c.Writer.Write(fileBytes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't write file"})
		return

	}

}

func (g *GameHandler) fromStringToGenres(stringGenres []string) ([]models.Genre, error) {
	genres := make([]models.Genre, len(stringGenres))
	for i, genre := range stringGenres {
		g, err := g.GenreService.GetGenreByName(genre)
		if err != nil {
			return []models.Genre{}, fmt.Errorf("unknown genre %s", genre)
		}
		genres[i] = g
	}
	return genres, nil
}

func (g *GameHandler) fromStringToPlatform(stringPlatform []string) ([]models.Platform, error) {
	plaforms := make([]models.Platform, len(stringPlatform))
	for i, platform := range stringPlatform {
		p, err := g.PlatformService.GetPlatformByName(platform)
		if err != nil {
			return []models.Platform{}, fmt.Errorf("unknown platform %s", platform)
		}
		plaforms[i] = p
	}
	return plaforms, nil
}

func (g *GameHandler) CreateGame(c *gin.Context) {
	inputGame := InputGame{}
	if err := c.ShouldBind(&inputGame); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't create game"})
		return
	}

	file := inputGame.File
	dst := fmt.Sprintf("library/%s", filepath.Base(file.Filename))
	if _, err := os.Stat(dst); err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrFileExists)
		return
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't save file"})
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
	dst = fmt.Sprintf("https://game-library-docker.onrender.com/image/library/%s", filepath.Base(file.Filename))
	game := models.NewGame(publisherId, inputGame.Title, inputGame.Description, dst, inputGame.AgeRestriction, inputGame.ReleaseYear)
	game, err = g.GameService.CreateGame(game, genres, plaforms)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game is successfully created", "data": gin.H{"gameId": game.ID, "link": dst}})
}

func query(ctx *gin.Context) (models.QueryParams, error) {
	skip, err := strconv.ParseUint(ctx.DefaultQuery("skip", "0"), 10, 32)
	if err != nil {
		return models.QueryParams{}, err
	}
	take, err := strconv.ParseUint(ctx.DefaultQuery("take", "99999"), 10, 32)
	if err != nil {
		return models.QueryParams{}, err
	}

	searchQuery := ctx.DefaultQuery("searchQuery", "")

	return models.QueryParams{
		Skip:        skip,
		Take:        take,
		SearchQuery: "%" + searchQuery + "%",
	}, nil
}

func (g *GameHandler) UpdateGame(c *gin.Context) {
	inputGame := UpdateGame{}
	if err := c.ShouldBind(&inputGame); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't parse input"})
		return
	}
	idStr := c.Param("id")
	gameR, err := g.GameService.GetGame(idStr)
	id, _ := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't get game"})
		return
	}
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
	if _, err := g.PublisherService.GetPublisher(publisherId.String()); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("unknown publisher"))
		return
	}
	file := inputGame.File
	dst := ""
	if file == nil {
		dst = gameR.ImageLink
	} else {
		dst = fmt.Sprintf("library/%s", filepath.Base(file.Filename))
		if _, err := os.Stat(dst); err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrFileExists)
			return
		}
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't save file"})
			return
		}
		c.Set("dst", dst)
		dst = fmt.Sprintf("http://localhost:8080/image/library/%s", filepath.Base(file.Filename))
	}
	game := models.NewGame(publisherId, inputGame.Title, inputGame.Description, dst, inputGame.AgeRestriction, inputGame.ReleaseYear)
	game.ID = id
	game, err = g.GameService.UpdateGame(game, genres, plaforms)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game is successfully updated", "data": gin.H{"game": game.ID, "link": dst}})
}

func (g *GameHandler) DeleteGame(c *gin.Context) {
	idStr := c.Param("id")
	err := g.GameService.DeleteGame(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't delete game"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game is successfully deleted"})
}
