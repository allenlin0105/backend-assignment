package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func OpenMySQLDatabase() (*sql.DB, error) {
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseMySQLDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database closed")
}
