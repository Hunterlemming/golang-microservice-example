package movie

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Hunterlemming/golang-microservice-example/api/model"

	"github.com/gorilla/mux"
)

type controller struct {
	service MovieService
}

type MovieController interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetMovie(w http.ResponseWriter, r *http.Request)
	CreateMovie(w http.ResponseWriter, r *http.Request)
	UpdateMovie(w http.ResponseWriter, r *http.Request)
	DeleteMovie(w http.ResponseWriter, r *http.Request)
}

func NewMovieController(s MovieService) MovieController {
	return &controller{
		service: s,
	}
}

func (c *controller) GetMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleMethodNotAllowed(w, fmt.Sprintf("%s method to GetMovies", r.Method))
		return
	}

	movies, err := c.service.GetMovies()
	if err != nil {
		handleServerError(w, err, "Service unreachable")
		return
	}

	res, _ := json.Marshal(movies)
	fmt.Fprintf(w, "%s", res)
}

func (c *controller) GetMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleMethodNotAllowed(w, fmt.Sprintf("%s method to GetMovie", r.Method))
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 0)
	if err != nil {
		handleBadRequest(w, "Invalid ID", err.Error())
		return
	}

	movie, err := c.service.GetMovie(int(id))
	if err != nil {
		handleServerError(w, err, "Service unreachable")
		return
	}

	res, _ := json.Marshal(movie)
	fmt.Fprintf(w, "%s", res)
}

func (c *controller) CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleMethodNotAllowed(w, fmt.Sprintf("%s method to CreateMovie", r.Method))
		return
	}

	// Extracting Movie object from request-body
	m, err := parseValidMovie(r)
	if err != nil {
		handleBadRequest(w, "Invalid request body", err.Error())
		return
	}

	// Create the Movie object in the database
	if err := c.service.CreateMovie(m); err != nil {
		handleServerError(w, err, "Service unreachable")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "success")
}

func (c *controller) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		handleMethodNotAllowed(w, fmt.Sprintf("%s method to UpdateMovie", r.Method))
		return
	}

	// Converting the ID to an integer
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 0)
	if err != nil {
		handleBadRequest(w, "Invalid ID", err.Error())
		return
	}

	// Extracting Movie object from request-body
	m, err := parseValidMovie(r)
	if err != nil {
		handleBadRequest(w, "Invalid request body", err.Error())
		return
	}

	// Updating Movie object in the database
	if err := c.service.UpdateMovie(int(id), m); err != nil {
		handleServerError(w, err, "Service unreachable")
		return
	}

	fmt.Fprintln(w, "success")
}

func (c *controller) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		handleMethodNotAllowed(w, fmt.Sprintf("%s method to DeleteMovie", r.Method))
		return
	}

	// Converting the ID to an integer
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 0)
	if err != nil {
		handleBadRequest(w, "Invalid ID", err.Error())
		return
	}

	// Deleting Movie object from the database
	if err := c.service.DeleteMovie(int(id)); err != nil {
		handleServerError(w, err, "Service unreachable")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "success")
}

func parseValidMovie(r *http.Request) (*model.Movie, error) {
	var m model.Movie

	// Return if the request-body cannot be decoded into a Movie object
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, err
	}

	// Return if the requested Movie object is invalid
	if err := m.Validate(); err != nil {
		return nil, errors.New("invalid movie-object")
	}

	return &m, nil
}

func handleMethodNotAllowed(w http.ResponseWriter, logMessage string) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	log.Println("[405 - Method Not Allowed] ", logMessage)
}

func handleServerError(w http.ResponseWriter, err error, message string) {
	http.Error(w, message, http.StatusInternalServerError)
	log.Println("[500 - Internal Server Error] ", err.Error())
}

func handleBadRequest(w http.ResponseWriter, responseMessage, logMessage string) {
	http.Error(w, responseMessage, http.StatusBadRequest)
	log.Println("[400 - Bad Request] ", logMessage)
}
