// Package models contains data structures used across the program/application.
package models

type Transaction struct {
	Type   string
	Amount int64
	Note   string
}

type BotConfig struct {
	TelebotToken  string
	DatabaseToken string
	DatabaseURL   string
}
