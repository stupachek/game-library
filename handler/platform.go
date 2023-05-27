package handler

import (
	"game-library/domens/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *PlatformHandler) CreatePlatform(c *gin.Context) {
	var platform models.Platform
	if err := c.ShouldBindJSON(&platform); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't parse input"})
		return
	}
	createdPlatform, err := p.PlatformService.CreatePlatform(platform)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't create platform"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Platform is successfully created", "data": createdPlatform})
}

func (p *PlatformHandler) GetPlatformsList(c *gin.Context) {
	platforms, err := p.PlatformService.GetPlatformsList()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can't get platform list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": platforms})
}
