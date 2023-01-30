package model

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type Api struct {
	Router *mux.Router
	DB     *sql.DB
}
