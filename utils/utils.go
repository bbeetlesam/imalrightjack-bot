// Package utils provides utility functions used across the application.
package utils

import (
	"strconv"
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
// the command type (e.g. "earn") and the bot name.
func ParseCommand(command string) models.Command {
	cmd := strings.TrimPrefix(command, "/")
	cmdType, cmdBot, _ := strings.Cut(cmd, "@")

	return models.Command{
		Action: cmdType,
		Bot:    cmdBot,
	}
}

// Itoa64 serves the same purpose as strconv.Itoa(), but for in64 type.
func Itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}

// ParseCommandMsg parses the whole command message and returns a string list.
// Not to be confused with ParseCommand() which only parses the prefix command from the message.
func ParseCommandMsg(msg string) []string {
	args := strings.Fields(msg)

	return args
}

// StringFieldsN serves as a combination between strings.Fields() (smart whitespaces parsing)
// and strings.SplitN (explicit separation amount).
func StringsFieldsN(str string, n int) []string {
	if n <= 0 {
		return nil
	}

	fields := strings.Fields(str)

	if len(fields) <= n {
		return fields
	}

	head := fields[:n-1]
	tail := strings.Join(fields[n-1:], " ")

	return append(head, tail)
}
