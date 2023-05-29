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

func (p *PlatformHandler) GetPlatform(c *gin.Context) {
	id := c.Param("id")

	platform, err := p.PlatformService.GetPlatform(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": platform})
}

func (p *PlatformHandler) UpdatePlatform(c *gin.Context) {
	var platformModel models.Platform
	if err := c.ShouldBindJSON(&platformModel); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idStr := c.Param("id")
	platform, err := p.PlatformService.UpdatePlatform(idStr, platformModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Platform is successfully updated", "data": platform})
}

func (p *PlatformHandler) DeletePlatform(c *gin.Context) {
	id := c.Param("id")
	err := p.PlatformService.DeletePlatform(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Platform is successfully deleted"})
}
