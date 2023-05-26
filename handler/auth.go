package handler

import (
	"game-library/domens/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *UserHandler) Register(c *gin.Context) {
	var registerUser models.RegisterModel
	if err := c.ShouldBindJSON(&registerUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	err := u.UserService.Register(registerUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sign up was successful"})
}

func (u *UserHandler) Login(c *gin.Context) {
	var loginUser models.LoginModel
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
		return
	}
	token, err := u.UserService.Login(loginUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sign up was successful", "token": token})
}
