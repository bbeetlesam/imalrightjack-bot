package main

import (
	"database/sql"
	"log"

	"github.com/bbeetlesam/imalrightjack-bot/messages"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func openDatabase(botCfg *BotConfig) (*sql.DB, error) {
	connectionString := botCfg.DatabaseURL + "?authToken=" + botCfg.DatabaseToken
	db, err := sql.Open("libsql", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println(messages.LogDBConnected) // supersister

	return db, nil
}

func initSchemaDB(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS transactions ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		type TEXT NOT NULL CHECK(type IN ('spend', 'earn')),
		timestamp INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		note TEXT
	);`

	_, err := db.Exec(query)	
	return err
}
