package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var (
	TURSO_TOKEN string = os.Getenv("TURSOTOKEN")
	TURSO_URL string = os.Getenv("TURSOURL")
	// currently unused
	query string = `CREATE TABLE IF NOT EXISTS transactions ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp INTEGER NOT NULL,
		amount REAL NOT NULL,
		note TEXT
	);`
)

func openDatabase() (*sql.DB, error) {
	connectionString := TURSO_URL + "?authToken=" + TURSO_TOKEN
	db, err := sql.Open("libsql", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfully connected to the database. A present from Nancy!") // supersister

	return db, nil
}
