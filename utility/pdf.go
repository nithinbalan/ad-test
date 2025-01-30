package utility

import (
	"bytes"
	"fmt"
	"html/template"
	"jethur/assets"
	"os"
	"os/exec"
)

func GeneratePDFContent(dto interface{}, templatePaths []string, logoImagePath string, coverBGImagePath string) (string, error) {
	// Prepare the data to include image paths
	data := map[string]interface{}{
		"LogoImagePath":    template.HTML(logoImagePath),
		"CoverBGImagePath": template.HTML(coverBGImagePath),
		"Data":             dto,
	}

	templates := template.New("pdf").Funcs(template.FuncMap{
		"formatDate": func(date string) string { return date },
		"safeHTML":   func(content string) template.HTML { return template.HTML(content) },
	})

	// Parse all templates at once
	for _, path := range templatePaths {
		content, err := assets.FS.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read template %s: %w", path, err)
		}
		_, err = templates.Parse(string(content))
		if err != nil {
			return "", fmt.Errorf("failed to parse template %s: %w", path, err)
		}
	}

	// Render the template into a buffer
	var renderedHTML bytes.Buffer
	if err := templates.Execute(&renderedHTML, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return renderedHTML.String(), nil
}

func ConvertHTMLContentToPDF(htmlContent, pdfFileName string) error {
	// Get the path to `wkhtmltopdf` from environment variables
	wkhtmltopdfPath := os.Getenv("WKHTMLTOPDF_PATH")
	if wkhtmltopdfPath == "" {
		// Provide a default path if the environment variable is not set
		wkhtmltopdfPath = "C:\\Program Files\\wkhtmltopdf\\bin\\wkhtmltopdf.exe"
	}

	// Check if `wkhtmltopdf` exists
	if _, err := os.Stat(wkhtmltopdfPath); os.IsNotExist(err) {
		return fmt.Errorf("wkhtmltopdf executable not found at %s", wkhtmltopdfPath)
	}

	// Create a temporary HTML file
	tempHTMLFile, err := os.CreateTemp("", "export_pdf_*.html")
	if err != nil {
		return fmt.Errorf("failed to create temporary HTML file: %w", err)
	}
	defer os.Remove(tempHTMLFile.Name())

	// Write the HTML content to the temporary file
	if _, err := tempHTMLFile.Write([]byte(htmlContent)); err != nil {
		return fmt.Errorf("failed to write HTML content to file: %w", err)
	}
	tempHTMLFile.Close()

	// Execute `wkhtmltopdf` to generate the PDF
	cmd := exec.Command(wkhtmltopdfPath, "--enable-local-file-access", tempHTMLFile.Name(), pdfFileName)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert HTML to PDF: %s: %w", stderr.String(), err)
	}

	return nil
}
