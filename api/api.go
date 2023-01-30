package api

import (
	"github.com/Hunterlemming/golang-microservice-example/api/model"
	"github.com/Hunterlemming/golang-microservice-example/api/movie"

	"github.com/gorilla/mux"
)

func Start() model.Api {
	api := model.Api{Router: mux.NewRouter(), DB: getDatabaseConnection()}
	movie.InitializeMoviesPipeline(&api)
	return api
}
