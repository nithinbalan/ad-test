package utility

import (
	"fmt"
	"io"
	"jethur/config"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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

func CopyFileFromPaths(srcDir, dstDir, fileName string) error {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	srcPath := filepath.Join(storagePath, srcDir, fileName)
	dstPath := filepath.Join(storagePath, dstDir, fileName)

	// Open the source file for reading
	sourceFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer sourceFile.Close()

	// Ensure that the destination directory exists
	destinationDir := filepath.Dir(dstPath) // Get the directory part of dstPath
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create the destination file for writing
	destinationFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Commit the file contents to stable storage
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file contents: %w", err)
	}

	return nil
}

// ConvertDocxToPDF converts a DOCX file to a PDF file. by  using LibreOffice:
func ConvertDocxToPDF(storagePath, fileName string, destinationPath string) error {
	basePath := config.GetEnv("STORAGE_PATH", "")
	inputFilePath := filepath.Join(basePath, storagePath, fileName)
	outputFilePath := filepath.Join(basePath, destinationPath, fileName[:len(fileName)-5]+".pdf")

	// Log paths for debugging
	//fmt.Printf("Input file path: %s\n", inputFilePath)
	//fmt.Printf("Output file path: %s\n", outputFilePath)

	// Check if the input file exists
	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %w", err)
	}

	// Read LibreOffice path from environment variable
	libreOfficePath := os.Getenv("LIBREOFFICE_PATH")
	if libreOfficePath == "" {
		return fmt.Errorf("LIBREOFFICE_PATH environment variable is not set")
	}

	// Determine the command based on the operating system
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(filepath.Join(libreOfficePath, "soffice.exe"), "--headless", "--convert-to", "pdf", "--outdir", filepath.Dir(outputFilePath), inputFilePath)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("libreoffice", "--headless", "--convert-to", "pdf", "--outdir", filepath.Dir(outputFilePath), inputFilePath)
	} else {
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to convert docx to pdf: %s\n%s", err, output)
	}

	return nil
}
