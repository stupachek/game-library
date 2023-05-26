package main

import (
	"game-library/api"
	"game-library/domens/repository/database"
	"log"
	"net/http"
)

func main() {
	DB := database.ConnectDataBase()
	app := api.SetupRouter(DB)
	err := http.ListenAndServe(":8080", app)
	if err != nil {
		log.Fatal(err)
	}
}
