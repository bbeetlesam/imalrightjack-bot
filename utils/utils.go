// Package utils provides utility functions used across the application.
package utils

import "strings"

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
