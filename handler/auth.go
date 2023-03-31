package handler

import (
	"game-library/domens/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *UserHandler) Register(c *gin.Context) {
	var registerUser models.RegisterModel
	if err := c.ShouldBindJSON(&registerUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return

	}
	err := u.UserService.Register(registerUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sign up was successful"})
}
