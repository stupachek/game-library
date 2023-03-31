package models

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Email          string
	Username       string
	BadgeColor     string
	Role           string
	HashedPassword string
	CreatedAt      time.Time //should it be only in db?
	Comments       []Comment
	Likes          []CommentLike
	Ratings        []Rating
}

type RegisterModel struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3,max=20"`
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
