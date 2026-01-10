package main

import (
	"errors"
	"os"

	"github.com/bbeetlesam/imalrightjack-bot/messages"
)

type BotConfig struct {
	TelebotToken  string
	DatabaseToken string
	DatabaseURL   string
}

func loadBotConfig() (*BotConfig, error) {
	teleToken := os.Getenv("TELETOKEN")
	if teleToken == "" {
		return nil, errors.New(messages.ErrTeletokenMissing)
	}

	tursoToken := os.Getenv("TURSOTOKEN")
	if tursoToken == "" {
		return nil, errors.New(messages.ErrDBTokenMissing)
	}

	tursoURL := os.Getenv("TURSOURL")
	if tursoURL == "" {
		return nil, errors.New(messages.ErrDBUrlMissing)
	}

	return &BotConfig{
		TelebotToken:  teleToken,
		DatabaseToken: tursoToken,
		DatabaseURL:   tursoURL,
	}, nil
}
