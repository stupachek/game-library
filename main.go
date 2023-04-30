package main

import (
	"game-library/api"
	"log"
	"net/http"
)

func main() {
	app := api.SetupRouter()
	err := http.ListenAndServe(":8080", app)
	if err != nil {
		log.Fatal(err)
	}
}
