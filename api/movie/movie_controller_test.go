package movie_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hunterlemming/golang-microservice-example/api/model"
	"github.com/Hunterlemming/golang-microservice-example/api/movie"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Defining the mock MovieService
type mockServiceStruct struct {
	mock.Mock
}

func (s *mockServiceStruct) GetMovies() ([]model.Movie, error) {
	args := s.Called()
	return args.Get(0).([]model.Movie), args.Error(1)
}

func (s *mockServiceStruct) GetMovie(id int) (model.Movie, error) {
	args := s.Called(id)
	return args.Get(0).(model.Movie), args.Error(1)
}

func (s *mockServiceStruct) CreateMovie(m *model.Movie) error {
	args := s.Called(m)
	return args.Error(0)
}

func (s *mockServiceStruct) UpdateMovie(id int, m *model.Movie) error {
	args := s.Called(id, m)
	return args.Error(0)
}

func (s *mockServiceStruct) DeleteMovie(id int) error {
	args := s.Called(id)
	return args.Error(0)
}

// Setting up mockService and controller
var mockService = new(mockServiceStruct)
var controller = movie.NewMovieController(mockService)

func TestControllerGetMovies(t *testing.T) {
	// Arrange
	movies := []model.Movie{
		{ID: 1, Name: "test1"},
		{ID: 2, Name: "test2"},
	}
	mockService.On("GetMovies").Return(movies, nil).Once()

	// Act
	req, _ := http.NewRequest("GET", "/", nil)
	rr := execute("/", []string{"GET"}, req, controller.GetMovies)

	// Assert
	if !mockService.AssertCalled(t, "GetMovies") {
		t.Error("The service should be called")
	}
	var status = http.StatusOK
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
	assert.Equal(t, jsonString(movies), rr.Body.String(), "The returned json-array should contain all the records in the database")
}

func TestControllerGetMoviesInvalidMethodError(t *testing.T) {
	// Act
	req, _ := http.NewRequest("POST", "/", nil)
	rr := execute("/", []string{"POST"}, req, controller.GetMovies)

	// Assert
	var status = http.StatusMethodNotAllowed
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerGetMoviesServiceError(t *testing.T) {
	// Arrange
	mockService.On("GetMovies").Return([]model.Movie{}, errors.New("test-error-message")).Once()

	// Act
	req, _ := http.NewRequest("GET", "/", nil)
	rr := execute("/", []string{"GET"}, req, controller.GetMovies)

	// Assert
	status := http.StatusInternalServerError
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerGetMovie(t *testing.T) {
	// Arrange
	movie := model.Movie{
		ID: 1, Name: "test",
	}
	mockService.On("GetMovie", 1).Return(movie, nil).Once()

	// Act
	req, _ := http.NewRequest("GET", "/1", nil)
	rr := execute("/{id}", []string{"GET"}, req, controller.GetMovie)

	// Assert
	if !mockService.AssertCalled(t, "GetMovie", 1) {
		t.Error("The service should be called")
	}
	status := http.StatusOK
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
	assert.Equal(t, movie.ID, movieJson(rr.Body.Bytes()).ID, "The returned json should be correct")
}

func TestControllerGetMovieInvalidMethodError(t *testing.T) {
	req, _ := http.NewRequest("POST", "/1", nil)
	rr := execute("/{id}", []string{"POST"}, req, controller.GetMovie)

	var status = http.StatusMethodNotAllowed
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerGetMovieIdParsingError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/not-an-int", nil)
	rr := execute("/{id}", []string{"GET"}, req, controller.GetMovie)

	status := http.StatusBadRequest
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerGetMovieServiceError(t *testing.T) {
	mockService.On("GetMovie", mock.Anything).Return(model.Movie{}, errors.New("test-error-message")).Once()

	req, _ := http.NewRequest("GET", "/1", nil)
	rr := execute("/{id}", []string{"GET"}, req, controller.GetMovie)

	status := http.StatusInternalServerError
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerCreateMovie(t *testing.T) {
	movie := model.Movie{ID: 1, Name: "test"}
	movieBytes, _ := json.Marshal(movie)
	mockService.On("CreateMovie", &movie).Return(nil).Once()

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(movieBytes))
	rr := execute("/", []string{"POST"}, req, controller.CreateMovie)

	if !mockService.AssertCalled(t, "CreateMovie", &movie) {
		t.Error("The service should be called")
	}
	status := http.StatusCreated
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerCreateMovieInvalidMethodError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	rr := execute("/", []string{"GET"}, req, controller.CreateMovie)

	var status = http.StatusMethodNotAllowed
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerCreateMovieBodyParsingError(t *testing.T) {
	notMovieBytes, _ := json.Marshal(struct{ ID int }{ID: 1})

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(notMovieBytes))
	rr := execute("/", []string{"POST"}, req, controller.CreateMovie)

	status := http.StatusBadRequest
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerCreateMovieServiceError(t *testing.T) {
	movieBytes, _ := json.Marshal(model.Movie{ID: 1, Name: "test"})
	mockService.On("CreateMovie", mock.Anything).Return(errors.New("test-error-message")).Once()

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(movieBytes))
	rr := execute("/", []string{"POST"}, req, controller.CreateMovie)

	status := http.StatusInternalServerError
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerUpdateMovie(t *testing.T) {
	movie := model.Movie{Name: "test"}
	movieBytes, _ := json.Marshal(movie)
	mockService.On("UpdateMovie", 1, &movie).Return(nil).Once()

	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer(movieBytes))
	rr := execute("/{id}", []string{"PUT"}, req, controller.UpdateMovie)

	if !mockService.AssertCalled(t, "UpdateMovie", 1, &movie) {
		t.Error("The service should be called")
	}
	status := http.StatusOK
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerUpdateMovieInvalidMethodError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/1", nil)
	rr := execute("/{id}", []string{"GET"}, req, controller.UpdateMovie)

	var status = http.StatusMethodNotAllowed
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerUpdateMovieIdParsingError(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/not-an-int", nil)
	rr := execute("/{id}", []string{"PUT"}, req, controller.UpdateMovie)

	status := http.StatusBadRequest
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerUpdateMovieBodyParsingError(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer([]byte("not a json")))
	rr := execute("/{id}", []string{"PUT"}, req, controller.UpdateMovie)

	status := http.StatusBadRequest
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerUpdateMovieServiceError(t *testing.T) {
	movieBytes, _ := json.Marshal(model.Movie{Name: "test"})
	mockService.On("UpdateMovie", 1, mock.Anything).Return(errors.New("test-error-message")).Once()

	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer(movieBytes))
	rr := execute("/{id}", []string{"PUT"}, req, controller.UpdateMovie)

	status := http.StatusInternalServerError
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerDeleteMovie(t *testing.T) {
	mockService.On("DeleteMovie", 1).Return(nil).Once()

	req, _ := http.NewRequest("DELETE", "/1", nil)
	rr := execute("/{id}", []string{"DELETE"}, req, controller.DeleteMovie)

	if !mockService.AssertCalled(t, "DeleteMovie", 1) {
		t.Error("The service should be called")
	}
	status := http.StatusNoContent
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerDeleteMovieInvalidMethodError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/1", nil)
	rr := execute("/{id}", []string{"GET"}, req, controller.DeleteMovie)

	var status = http.StatusMethodNotAllowed
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerDeleteMovieIdParsingError(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/not-an-int", nil)
	rr := execute("/{id}", []string{"DELETE"}, req, controller.DeleteMovie)

	status := http.StatusBadRequest
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func TestControllerDeleteMovieServiceError(t *testing.T) {
	mockService.On("DeleteMovie", mock.Anything).Return(errors.New("test-error-message")).Once()

	req, _ := http.NewRequest("DELETE", "/1", nil)
	rr := execute("/{id}", []string{"DELETE"}, req, controller.DeleteMovie)

	status := http.StatusInternalServerError
	assert.Equal(t, status, rr.Code, fmt.Sprintf("Status code should be [%d]", status))
}

func execute(route string, methods []string, req *http.Request, handler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	r := mux.NewRouter()
	r.HandleFunc(route, handler).Methods(methods...)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func jsonString(obj interface{}) string {
	res, _ := json.Marshal(obj)
	return string(res)
}

func movieJson(obj []byte) *model.Movie {
	var movie model.Movie
	err := json.Unmarshal(obj, &movie)
	if err != nil {
		log.Fatal("the service returned an unknown object as a movie")
	}
	return &movie
}
