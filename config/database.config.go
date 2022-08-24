package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {

	//loading parameters from env file
	// username := os.Getenv("DATABASE_USERNAME")
	// password := os.Getenv("DATABASE_PASSWORD")
	// databaseHost := os.Getenv("DATABASE_HOST")
	databaseName := os.Getenv("DATABASE_NAME")
	//formatting
	// dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)
	dbUri := os.Getenv("DB_SOURCE")
	//Opens database
	db, err := sql.Open("postgres", dbUri)

	if err != nil {
		log.Fatal(err)
	}

	// verifies connection to the database is still alive
	err = db.Ping()
	if err != nil {
		fmt.Println("error in pinging")
		log.Fatal(err)

	}

	log.Println("\nConnected to database:", databaseName)

	return db

}
