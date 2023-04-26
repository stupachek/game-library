package handler

import (
	"game-library/domens/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (p *PublisherHandler) GetPlatformsList(c *gin.Context) {
	publishers := p.PublisherService.GetPublishersList()
	c.JSON(http.StatusOK, gin.H{"data": publishers})
}

func (p *PublisherHandler) CreatePublisher(c *gin.Context) {
	var publisher models.PublisherModel
	if err := c.ShouldBindJSON(&publisher); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	createdPublisher := p.PublisherService.CreatePublisher(publisher)
	c.JSON(http.StatusOK, gin.H{"message": "Publisher is successfully created", "data": createdPublisher})
}

func (p *PublisherHandler) GetPublisher(c *gin.Context) {
	id := c.Param("id")

	publisher, err := p.PublisherService.GetPublisher(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": publisher})
}

func (p *PublisherHandler) UpdatePublisher(c *gin.Context) {
	var publisherModel models.PublisherModel
	if err := c.ShouldBindJSON(&publisherModel); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	idStr := c.Param("id")
	publisher, err := p.PublisherService.UpdatePublisher(idStr, publisherModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Publisher is successfully updated", "data": publisher})
}

func (p *PublisherHandler) DeletePublisher(c *gin.Context) {
	id := c.Param("id")
	err := p.PublisherService.DeletePublisher(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Publisher is successfully deleted"})
}
