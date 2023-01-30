package movie

import (
	"github.com/Hunterlemming/golang-microservice-example/api/model"

	"github.com/gorilla/mux"
)

func InitializeMoviesPipeline(api *model.Api) {
	s := NewMovieService(api.DB)
	c := NewMovieController(s)
	setRouting(api.Router, c)
}

func setRouting(main *mux.Router, c MovieController) {
	sr := main.PathPrefix("/movies").Subrouter()

	sr.HandleFunc("", c.GetMovies).
		Methods("GET")

	sr.HandleFunc("/{id}", c.GetMovie).
		Methods("GET")

	sr.HandleFunc("", c.CreateMovie).
		Methods("POST")

	sr.HandleFunc("/{id}", c.UpdateMovie).
		Methods("PUT")

	sr.HandleFunc("/{id}", c.DeleteMovie).
		Methods("DELETE")
}
