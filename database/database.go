// Package database provides database connection management and transaction operations.
package database

import (
	"database/sql"

	"github.com/bbeetlesam/imalrightjack-bot/messages"
	"github.com/bbeetlesam/imalrightjack-bot/models"
	"github.com/bbeetlesam/imalrightjack-bot/utils"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func Open(botCfg *models.BotConfig) (*sql.DB, error) {
	connectionString := botCfg.DatabaseURL + "?authToken=" + botCfg.DatabaseToken
	db, err := sql.Open("libsql", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	utils.LogColor("info", messages.LogDBConnected) // supersister

	return db, nil
}

func InitSchema(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS transactions ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		type TEXT NOT NULL CHECK(type IN ('spend', 'earn')),
		timestamp TEXT NOT NULL,
		amount INTEGER NOT NULL,
		note TEXT
	);`

	_, err := db.Exec(query)
	return err
}
