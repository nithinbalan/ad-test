package utility

import "strings"

func SetDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func RemoveLangCodeSuffix(text string) string {
	text = strings.TrimSuffix(text, "En")
	text = strings.TrimSuffix(text, "Ar")
	return strings.ToLower(text)
}
