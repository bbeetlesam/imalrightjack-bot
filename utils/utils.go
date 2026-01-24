// Package utils provides utility functions used across the application.
package utils

import (
	"strings"

	"github.com/bbeetlesam/imalrightjack-bot/models"
)

// EscapeMarkdownV2 escapes special characters for Telegram MarkdownV2 parse mode.
// This prevents parse errors when sending messages with special characters.
// Blame this mess on MarkdownV2's clunky escapism mind.
func EscapeMarkdownV2(text string) string {
	specialChars := []string{
		"_", "*", "[", "]", "(", ")", "~", "`", ">",
		"#", "+", "-", "=", "|", "{", "}", ".", "!",
	}

	result := text
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}

	return result
}

// ParseCommand parses a command string, e.g. "/earn" and "/earn@bot" and returns
// the command type ("spend" or "earn") and the bot name.
func ParseCommand(command string) models.Command {
	cmd := strings.TrimPrefix(command, "/")
	cmdType, cmdBot, _ := strings.Cut(cmd, "@")

	return models.Command{
		Action: cmdType,
		Bot:    cmdBot,
	}
}
