package movie

import (
	"database/sql"
	"fmt"

	"github.com/Hunterlemming/golang-microservice-example/api/model"
	"github.com/Hunterlemming/golang-microservice-example/api/util"
)

type service struct {
	db *sql.DB
}

type MovieService interface {
	GetMovies() ([]model.Movie, error)
	GetMovie(id int) (model.Movie, error)
	CreateMovie(m *model.Movie) error
	UpdateMovie(id int, m *model.Movie) error
	DeleteMovie(id int) error
}

func NewMovieService(db *sql.DB) MovieService {
	return &service{
		db: db,
	}
}

func (s *service) GetMovies() ([]model.Movie, error) {
	const q = "SELECT * FROM movies"
	qr, err := s.db.Query(q)
	if err != nil {
		return []model.Movie{}, err
	}
	defer qr.Close()

	result := make([]model.Movie, 0)
	for qr.Next() {
		m := model.Movie{}
		err = qr.Scan(&m.ID, &m.Name)
		if err != nil {
			return []model.Movie{}, err
		}
		result = append(result, m)
	}

	return result, nil
}

func (s *service) GetMovie(id int) (model.Movie, error) {
	const q = "SELECT * FROM movies WHERE id = $1"
	qr := s.db.QueryRow(q, id)

	result := model.Movie{}
	err := qr.Scan(&result.ID, &result.Name)
	if err != nil {
		return model.Movie{}, err
	}

	return result, nil
}

func (s *service) CreateMovie(m *model.Movie) error {
	// Returning if a record by the parameter id already exists
	if _, err := s.GetMovie(m.ID); err == nil {
		return &util.ExistingRecordError{Identification: fmt.Sprintf("ID: %v", m.ID)}
	}

	// Inserting the new record
	const q = "INSERT INTO movies (id, name) VALUES ($1, $2)"
	if _, err := s.db.Exec(q, m.ID, m.Name); err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateMovie(id int, m *model.Movie) error {
	// Returning if the record to update was not found in the database
	if _, err := s.GetMovie(id); err != nil {
		return &util.NotExistingRecordError{Identification: fmt.Sprintf("ID: %v", m.ID)}
	}

	// Updating existing record
	const q = "UPDATE movies SET name = $1 WHERE id = $2"
	if _, err := s.db.Exec(q, m.Name, id); err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteMovie(id int) error {
	const q = "DELETE FROM movies WHERE id = $1"
	if _, err := s.db.Exec(q, id); err != nil {
		return err
	}
	return nil
}
