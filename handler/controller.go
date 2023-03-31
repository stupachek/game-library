package handler

import (
	"game-library/domens/service"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHadler(service service.UserService) UserHandler {
	return UserHandler{
		UserService: service,
	}
}
