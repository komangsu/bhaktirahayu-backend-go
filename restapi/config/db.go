package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	dbURI := os.Getenv("DATABASE_URI")

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("Failed connect to database: ", err)
	}

	return db
}
