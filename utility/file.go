package utility

import (
	"path/filepath"
	"regexp"
	"strings"
)

func SanitizeFileName(fileName string) string {
	extension := filepath.Ext(fileName)

	fileNameWithoutExtension := strings.ToLower(fileName[:len(fileName)-len(extension)])

	// Remove special characters
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return fileName
	}
	fileNameWithoutExtension = reg.ReplaceAllString(fileNameWithoutExtension, " ")

	fileNameWithoutExtension = strings.ReplaceAll(fileNameWithoutExtension, " ", "_")

	return fileNameWithoutExtension + extension
}
