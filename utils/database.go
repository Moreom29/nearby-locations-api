package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "om"
	dbname := "postgres"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	fmt.Println("Connected to the database!")
}
