package main

import (
	"game-library/api"
	"log"
	"net/http"
)

func main() {
	app := api.SetupRouter()
	err := http.ListenAndServe(":8081", app)
	if err != nil {
		log.Fatal(err)
	}
}
