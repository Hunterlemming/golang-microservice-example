package movie_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hunterlemming/golang-microservice-example/api/model"
	"github.com/Hunterlemming/golang-microservice-example/api/movie"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestInitializeMoviesPipeline(t *testing.T) {
	r := mux.NewRouter()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	api := model.Api{Router: r, DB: db}
	movie.InitializeMoviesPipeline(&api)

	testIntegrationGetAll(t, mock, api.Router)
	testIntegrationGetOne(t, mock, api.Router)
	testIntegrationCreate(t, mock, api.Router)
	testIntegrationUpdate(t, mock, api.Router)
	testIntegrationDelete(t, mock, api.Router)
}

func testIntegrationGetAll(t *testing.T, mock sqlmock.Sqlmock, r *mux.Router) {
	getAllResult := []model.Movie{{ID: 1, Name: "t1"}, {ID: 2, Name: "t2"}}
	mock.ExpectQuery(GetAllQuery).WillReturnRows(newRows(&getAllResult))

	req, _ := http.NewRequest("GET", "/movies", nil)
	rr := executeWithRouter(r, req)

	assert.Equal(t, jsonString(getAllResult), rr.Body.String())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func testIntegrationGetOne(t *testing.T, mock sqlmock.Sqlmock, r *mux.Router) {
	getOneResult := model.Movie{ID: 1, Name: "t1"}
	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&[]model.Movie{getOneResult}))

	req, _ := http.NewRequest("GET", "/movies/1", nil)
	rr := executeWithRouter(r, req)

	assert.Equal(t, jsonString(getOneResult), rr.Body.String())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func testIntegrationCreate(t *testing.T, mock sqlmock.Sqlmock, r *mux.Router) {
	movie := model.Movie{ID: 1, Name: "test"}
	movieBytes, _ := json.Marshal(movie)
	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&[]model.Movie{}))
	mock.ExpectExec(CreateMovieQuery).WithArgs(movie.ID, movie.Name).
		WillReturnResult(sqlmock.NewResult(int64(movie.ID), 1))

	req, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(movieBytes))
	rr := executeWithRouter(r, req)

	status := http.StatusCreated
	assert.Equal(t, status, rr.Code)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func testIntegrationUpdate(t *testing.T, mock sqlmock.Sqlmock, r *mux.Router) {
	movies := []model.Movie{{Name: "test"}}
	updatedMovie := model.Movie{ID: 1, Name: "updated"}
	updatedMovieBytes, _ := json.Marshal(updatedMovie)
	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&movies))
	mock.ExpectExec(UpdateQuery).WithArgs(updatedMovie.Name, updatedMovie.ID).
		WillReturnResult(sqlmock.NewResult(int64(updatedMovie.ID), 1))

	req, _ := http.NewRequest("PUT", "/movies/1", bytes.NewBuffer(updatedMovieBytes))
	rr := executeWithRouter(r, req)

	status := http.StatusOK
	assert.Equal(t, status, rr.Code)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func testIntegrationDelete(t *testing.T, mock sqlmock.Sqlmock, r *mux.Router) {
	deletedIndex := 1
	mock.ExpectExec(DeleteQuery).WithArgs(deletedIndex).
		WillReturnResult(sqlmock.NewResult(int64(deletedIndex), 1))

	req, _ := http.NewRequest("DELETE", "/movies/1", nil)
	rr := executeWithRouter(r, req)

	status := http.StatusNoContent
	assert.Equal(t, status, rr.Code)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func executeWithRouter(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
