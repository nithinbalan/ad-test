package utility

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// RunPDFConversionCheck checks if PDF conversion using LibreOffice is working.
func RunPDFConversionCheck() {
	fmt.Println("🔍 Checking PDF conversion using LibreOffice...")

	// Step 1: Check if LibreOffice is installed
	if !isLibreOfficeInstalled() {
		fmt.Println("❌ LibreOffice is not installed or not found in PATH.")
		os.Exit(1)
	}
	fmt.Println("✅ LibreOffice is installed.")

	// Define paths for input and output files
	rootPath := os.Getenv("STORAGE_PATH") // Ensure STORAGE_PATH is defined in your .env file
	if rootPath == "" {
		fmt.Println("❌ STORAGE_PATH is not set in the environment variables.")
		os.Exit(1)
	}

	inputFile := filepath.Join(rootPath, "test.docx")
	outputFile := filepath.Join(rootPath, "test.pdf")

	// Step 2: Check if the input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("❌ Input file does not exist at %s\n", inputFile)
		os.Exit(1)
	}
	fmt.Println("✅ Input DOCX file found.")

	// Step 3: Convert DOCX to PDF
	if err := convertToPDF(inputFile, filepath.Dir(outputFile)); err != nil {
		fmt.Printf("❌ PDF conversion failed: %v\n", err)
		os.Exit(1)
	}

	// Step 4: Verify the PDF Exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		fmt.Printf("❌ PDF file not created at %s. Conversion failed.\n", outputFile)
		os.Exit(1)
	}
	fmt.Println("✅ PDF conversion test passed successfully!")
}

// Check if LibreOffice is installed
func isLibreOfficeInstalled() bool {
	//var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		_, err := exec.LookPath("soffice.exe")
		return err == nil
	}
	_, err := exec.LookPath("libreoffice")
	return err == nil
}

// Convert DOCX to PDF
func convertToPDF(inputFile, outputDir string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("soffice.exe", "--headless", "--convert-to", "pdf", "--outdir", outputDir, inputFile)
	} else {
		cmd = exec.Command("libreoffice", "--headless", "--convert-to", "pdf", "--outdir", outputDir, inputFile)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("LibreOffice conversion failed: %s\n%s", err, output)
	}
	return nil
}
