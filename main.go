package main

import (
	"net/http"

	"github.com/Hunterlemming/golang-microservice-example/api"

	_ "github.com/gorilla/mux"
)

func main() {
	api := api.Start()
	defer api.DB.Close()
	http.ListenAndServe(":8080", api.Router)
}
