package utils

import "strings"

func CleanGeminiResponse(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.ReplaceAll(text, "\\\"", "\"")
	text = strings.ReplaceAll(text, "\"{", "{")
	text = strings.ReplaceAll(text, "}\"", "}")
	return strings.TrimSpace(text)
}
