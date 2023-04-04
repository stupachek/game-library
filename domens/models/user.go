package models

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

const (
	USER    = "user"
	ADMIN   = "admin"
	MANAGER = "manager"
)

type User struct {
	ID             uuid.UUID     `json:"id"`
	Email          string        `json:"email"`
	Username       string        `json:"usarname"`
	BadgeColor     string        `json:"badge_color"`
	Role           string        `json:"role"`
	HashedPassword string        `json:"-"`
	CreatedAt      time.Time     `json:"-"` //should it be only in db?
	Comments       []Comment     `json:"-"`
	Likes          []CommentLike `json:"-"`
	Ratings        []Rating      `json:"-"`
}

type RegisterModel struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3,max=20"`
}

type LoginModel struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3,max=20"`
}
type Role struct {
	Role string `json:"role"`
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Should be longer than " + fe.Param()
	case "max":
		return "Should be shorter than " + fe.Param()
	}
	return "Unknown error"
}
