package utility

import (
	"jethur/config"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

// StoreFile moves a file from a temporary location to a permanent location.
func StoreFile(fileName, destinationPath string) error {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	temporaryStoragePath := filepath.Join(storagePath, "tmp", fileName)
	destinationStoragePath := filepath.Join(storagePath, destinationPath, fileName)

	if err := os.MkdirAll(filepath.Dir(destinationStoragePath), os.ModePerm); err != nil {
		return err
	}
	if err := os.Rename(temporaryStoragePath, destinationStoragePath); err != nil {
		return err
	}
	return nil
}

func RemoveFile(fileName, destinationPath string) error {
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	filePath := filepath.Join(baseStoragePath, destinationPath, fileName)
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func RemoveFolder(destinationPath string) error {
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	directoryPath := filepath.Join(baseStoragePath, destinationPath)
	if err := os.RemoveAll(directoryPath); err != nil {
		return err
	}
	return nil
}

func GenerateDocumentURL(endpoint, resourceType string, ids ...interface{}) string {
	// baseURL := config.GetEnv("API_URL", "http://localhost:8080")
	var url strings.Builder
	// url.WriteString(baseURL)
	// url.WriteString("/")
	url.WriteString(endpoint)
	url.WriteString("/")
	url.WriteString(resourceType)

	for _, id := range ids {
		url.WriteString("/")
		switch v := id.(type) {
		case uint64:
			url.WriteString(strconv.FormatUint(v, 10))
		case string:
			url.WriteString(v)
		default:
			// Handle other types or ignore if not needed
		}
	}

	return url.String()
}
