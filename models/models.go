// Package models contains data structures used across the program/application.
package models

type Transaction struct {
	Type   string
	Amount int64
	Note   string
	Time   string
}

type BotConfig struct {
	TelebotToken  string
	DatabaseToken string
	DatabaseURL   string
}

type Command struct {
	Action string
	Bot    string
}
