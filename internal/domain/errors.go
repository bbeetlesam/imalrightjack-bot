package domain

import "errors"

var (
	// config errors
	ErrMissingTeletoken = errors.New("$TELETOKEN environment variable is missing")
	ErrMissingDBToken   = errors.New("$TURSOTOKEN environment variable is missing")
	ErrMissingDBURL     = errors.New("$TURSOURL environment variable is missing")

	// transaction errors
	ErrInvalidAmount        = errors.New("amount must be greater than zero")
	ErrInvalidType          = errors.New("transaction type must be 'earn' or 'spend'")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrInvalidTransactionID = errors.New("invalid transaction ID")

	// user errors
	ErrInvalidUserID = errors.New("invalid user ID")
)
