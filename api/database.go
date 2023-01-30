package api

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	host = "localhost"
	port = "5432"
)

var (
	user   string
	pw     string
	dbname string
)

func getKeysFromEnvironment() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Configuration file not found")
	}

	user = viper.GetString("APP_DB_USERNAME")
	pw = viper.GetString("APP_DB_PASSWORD")
	dbname = viper.GetString("APP_DB_NAME")
}

func getDatabaseConnection() *sql.DB {
	getKeysFromEnvironment()

	cs := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pw, dbname)
	db, err := sql.Open("postgres", cs)
	checkError(err)

	err = db.Ping()
	checkError(err)

	fmt.Println("Connected to Database!")
	return db
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
