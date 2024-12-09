package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	name     = os.Getenv("DB_NAME")
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	stage    = os.Getenv("STAGE_ENV")
)

func NewPostgresDb() (*sql.DB, error) {
	var connStr string

	if stage == "development" {
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	} else {
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, user, password, name)
	}

	db, err := sql.Open("postgres", connStr)

	return db, err
}

func InitializeDb(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Couldn not connect to database: %s", err)
	}

	log.Println("Database connected")
}
