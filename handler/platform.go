package handler

import (
	"game-library/domens/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (p *PlatformHandler) CreatePlatform(c *gin.Context) {
	var platform models.Platform
	if err := c.ShouldBindJSON(&platform); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	createdPlatform := p.PlatformService.CreatePlatform(platform)
	c.JSON(http.StatusOK, gin.H{"message": "Platform is successfully created", "data": createdPlatform})
}

func (p *PlatformHandler) GetPlatformsList(c *gin.Context) {
	platforms := p.PlatformService.GetPlatformsList()
	c.JSON(http.StatusOK, gin.H{"data": platforms})
}
