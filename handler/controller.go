package handler

import (
	"game-library/domens/service"
)

type UserHandler struct {
	UserService service.UserService
}

type GameHandler struct {
	GameService      service.GameService
	PublisherService service.PublisherService
	GenreService     service.GenreService
	PlatformService  service.PlatformService
}

type PublisherHandler struct {
	PublisherService service.PublisherService
}

type GenreHandler struct {
	GenreService service.GenreService
}

type PlatformHandler struct {
	PlatformService service.PlatformService
}

func NewUserHadler(service service.UserService) UserHandler {
	return UserHandler{
		UserService: service,
	}
}

func NewGameHandler(serviceG service.GameService, serviceP service.PublisherService) GameHandler {
	return GameHandler{
		GameService:      serviceG,
		PublisherService: serviceP,
	}
}

func NewPublisherHandler(service service.PublisherService) PublisherHandler {
	return PublisherHandler{
		PublisherService: service,
	}
}

func NewPlatformHandler(service service.PlatformService) PlatformHandler {
	return PlatformHandler{
		PlatformService: service,
	}
}

func NewGenreHandler(service service.GenreService) GenreHandler {
	return GenreHandler{
		GenreService: service,
	}
}
