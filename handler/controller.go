package handler

import (
	"game-library/domens/service/game"
	"game-library/domens/service/genre"
	"game-library/domens/service/platform"
	"game-library/domens/service/publisher"
	service "game-library/domens/service/user"
)

type UserHandler struct {
	UserService service.UserService
}

type GameHandler struct {
	GameService      game.GameService
	PublisherService publisher.PublisherService
	GenreService     genre.GenreService
	PlatformService  platform.PlatformService
}

type PublisherHandler struct {
	PublisherService publisher.PublisherService
}

type GenreHandler struct {
	GenreService genre.GenreService
}

type PlatformHandler struct {
	PlatformService platform.PlatformService
}

func NewUserHadler(service service.UserService) UserHandler {
	return UserHandler{
		UserService: service,
	}
}

func NewGameHandler(game game.GameService, publisher publisher.PublisherService, genre genre.GenreService, platform platform.PlatformService) GameHandler {
	return GameHandler{
		GameService:      game,
		PublisherService: publisher,
		GenreService:     genre,
		PlatformService:  platform,
	}
}

func NewPublisherHandler(service publisher.PublisherService) PublisherHandler {
	return PublisherHandler{
		PublisherService: service,
	}
}

func NewPlatformHandler(service platform.PlatformService) PlatformHandler {
	return PlatformHandler{
		PlatformService: service,
	}
}

func NewGenreHandler(service genre.GenreService) GenreHandler {
	return GenreHandler{
		GenreService: service,
	}
}
