package handler

import (
	"game-library/domens/service"
)

type UserHandler struct {
	UserService service.UserService
}

type GameHandler struct {
	GameService service.GameService
}

type PublisherHandler struct {
	PublisherService service.PublisherService
}

func NewUserHadler(service service.UserService) UserHandler {
	return UserHandler{
		UserService: service,
	}
}

func NewGameHandler(service service.GameService) GameHandler {
	return GameHandler{
		GameService: service,
	}
}

func NewPublisherHandler(service service.PublisherService) PublisherHandler {
	return PublisherHandler{
		PublisherService: service,
	}
}
