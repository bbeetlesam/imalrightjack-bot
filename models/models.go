// Package models contains data structures used across the program/application.
package models

type TransactionType string

const (
	TransactionTypeEarn  TransactionType = "earn"
	TransactionTypeSpend TransactionType = "spend"
)

type Transaction struct {
	ID     int64
	Type   TransactionType
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
	Action TransactionType
	Bot    string
}
