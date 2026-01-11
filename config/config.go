// Package config handles loading and validating bot configuration from environment variables.
package config

import (
	"errors"
	"os"

	"github.com/bbeetlesam/imalrightjack-bot/messages"
	"github.com/bbeetlesam/imalrightjack-bot/models"
)

func LoadBotConfig() (*models.BotConfig, error) {
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

	return &models.BotConfig{
		TelebotToken:  teleToken,
		DatabaseToken: tursoToken,
		DatabaseURL:   tursoURL,
	}, nil
}
