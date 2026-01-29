package database

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open Database
	log.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test Connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set Connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database Connected!")
	return db, nil
}