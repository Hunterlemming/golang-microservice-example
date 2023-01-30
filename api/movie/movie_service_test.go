package movie_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/Hunterlemming/golang-microservice-example/api/model"
	"github.com/Hunterlemming/golang-microservice-example/api/movie"
	"github.com/Hunterlemming/golang-microservice-example/api/util"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const GetAllQuery = `^SELECT \* FROM [\p{L}\p{N}.]+$`

func TestServiceGetMovies(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	movies := []model.Movie{{ID: 1, Name: "test1"}, {ID: 2, Name: "test2"}}
	mock.ExpectQuery(GetAllQuery).WillReturnRows(newRows(&movies))

	res, err := service.GetMovies()

	assert.Equal(t, movies, res)
	assert.Equal(t, nil, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestServiceGetMoviesQueryError(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	queryError := errors.New("test-error-message")
	mock.ExpectQuery(GetAllQuery).WillReturnError(queryError)

	res, err := service.GetMovies()

	assert.Equal(t, []model.Movie{}, res)
	assert.Equal(t, queryError, err)
}

func TestServiceGetMoviesRowScanError(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"a", "s", "d"}).AddRow(1, 2, 3)
	mock.ExpectQuery(GetAllQuery).WillReturnRows(rows)

	res, err := service.GetMovies()

	assert.Equal(t, []model.Movie{}, res)
	assert.NotEqual(t, nil, err)
}

const GetOneQuery = `^SELECT \* FROM [\p{L}\p{N}.]+ WHERE [\p{L}\p{N}.]+ = \$1$`

func TestServiceGetMovie(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	movies := []model.Movie{{ID: 1, Name: "test1"}}
	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&movies))

	res, err := service.GetMovie(1)

	assert.Equal(t, movies[0], res)
	assert.Equal(t, nil, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestServiceGetMovieRowScanError(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"a", "s", "d"}).AddRow(1, 2, 3)
	mock.ExpectQuery(GetOneQuery).WillReturnRows(rows)

	res, err := service.GetMovie(1)

	assert.Equal(t, model.Movie{}, res)
	assert.NotEqual(t, nil, err)
}

const CreateMovieQuery = `^INSERT INTO [\p{L}\p{N}.]+ \([\p{L}\p{N},. ]+\) VALUES \([\p{N}$, ]+\)$`

func TestServiceCreateMovie(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&[]model.Movie{}))
	mock.ExpectExec(CreateMovieQuery).WithArgs(2, "test2").
		WillReturnResult(sqlmock.NewResult(2, 1))

	err := service.CreateMovie(&model.Movie{ID: 2, Name: "test2"})

	assert.Equal(t, nil, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestServiceCreateMovieRecordExistsError(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&[]model.Movie{{ID: 2, Name: "asd"}}))

	err := service.CreateMovie(&model.Movie{ID: 2, Name: "test2"})

	assert.IsType(t, &util.ExistingRecordError{}, err)
}

func TestServiceCreateMovieInsertError(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&[]model.Movie{}))
	insertError := errors.New("test-error-message")
	mock.ExpectExec(CreateMovieQuery).WithArgs(2, "test2").WillReturnError(insertError)

	err := service.CreateMovie(&model.Movie{ID: 2, Name: "test2"})

	assert.Equal(t, insertError, err)
}

const UpdateQuery = `^UPDATE [\p{L}\p{N}.]+ SET ([\p{L}\p{N}.]+ = \$\p{N}+[, ]+)+WHERE [\p{L}\p{N}.]+ = \$\p{N}+$`

func TestServiceUpdateMovie(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	movies := []model.Movie{{ID: 1, Name: "test1"}}
	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&movies))
	mock.ExpectExec(UpdateQuery).WithArgs("updated", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := service.UpdateMovie(1, &model.Movie{ID: 1, Name: "updated"})

	assert.Equal(t, nil, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestServiceUpdateMovieRecordDoesNotExistError(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&[]model.Movie{}))

	err := service.UpdateMovie(1, &model.Movie{ID: 1, Name: "updated"})

	assert.IsType(t, &util.NotExistingRecordError{}, err)
}

func TestServiceUpdateMovieUpdateError(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	movies := []model.Movie{{ID: 1, Name: "test1"}}
	mock.ExpectQuery(GetOneQuery).WillReturnRows(newRows(&movies))
	updateError := errors.New("test-error-message")
	mock.ExpectExec(UpdateQuery).WithArgs("updated", 1).WillReturnError(updateError)

	err := service.UpdateMovie(1, &model.Movie{ID: 1, Name: "updated"})

	assert.Equal(t, updateError, err)
}

const DeleteQuery = `^DELETE FROM [\p{L}\p{N}.]+ WHERE [\p{L}\p{N}.]+ = \$1$`

func TestServiceDeleteMovie(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	mock.ExpectExec(DeleteQuery).WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := service.DeleteMovie(1)

	assert.Equal(t, nil, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestServiceDeleteMovieDeleteError(t *testing.T) {
	service, mock, db := initNewService(t)
	defer db.Close()

	deleteError := errors.New("test-error-message")
	mock.ExpectExec(DeleteQuery).WithArgs(1).
		WillReturnError(deleteError)

	err := service.DeleteMovie(1)

	assert.Equal(t, deleteError, err)
}

func initNewService(t *testing.T) (movie.MovieService, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	service := movie.NewMovieService(db)
	return service, mock, db
}

func newRows(movies *[]model.Movie) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, m := range *movies {
		rows.AddRow(m.ID, m.Name)
	}
	return rows
}
