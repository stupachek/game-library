package handler

import (
	"game-library/domens/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrUnauthenticated = gin.H{"error": "unauthenticated"}

func (u *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		userID = GetId(c)
	}
	user, err := u.UserService.GetUser(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (u *UserHandler) GetUsers(c *gin.Context) {
	users := u.UserService.GetUsers()
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (u *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	err := u.UserService.DeleteUser(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User is successfully deleted"})
}

func (u *UserHandler) ChangerRole(c *gin.Context) {
	role := models.Role{}
	if err := c.ShouldBindJSON(&role); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	idStr := c.Param("id")
	user, err := u.UserService.ChangeRole(idStr, role.Role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User is successfully updated", "data": user})
}

func GetId(c *gin.Context) string {
	id, ok := c.Get("UserID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": ErrUnauthenticated})
	}
	userID, ok := id.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": ErrUnauthenticated})
	}
	return userID
}
