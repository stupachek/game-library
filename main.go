package main

import (
	"game-library/api"
	"game-library/domens/repository"
	"log"
	"net/http"
)

func main() {
	_ = repository.ConnectDataBase()
	app := api.SetupRouter()
	err := http.ListenAndServe(":8080", app)
	if err != nil {
		log.Fatal(err)
	}
}
