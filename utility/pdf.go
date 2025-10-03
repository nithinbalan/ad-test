package utility

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"jethur/assets"
	"jethur/config"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func GeneratePDFContent(dto interface{}, templatePaths []string, logoImagePath string, coverBGImagePath string) (string, error) {
	// Prepare the data to include image paths
	data := map[string]interface{}{
		"LogoImagePath":    template.HTML(logoImagePath),
		"CoverBGImagePath": template.HTML(coverBGImagePath),
		"Data":             dto,
		"ApiPath":          config.GetEnv("API_URL", ""),
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

func ConvertPolicyHTMLContentToPDF(headerHTML, footerHTML, coverHTML, disclaimerHTML, contentHTML, pdfFileName string) error {
	wkhtmltopdfPath := os.Getenv("WKHTMLTOPDF_PATH")
	if wkhtmltopdfPath == "" {
		wkhtmltopdfPath = "C:\\Program Files\\wkhtmltopdf\\bin\\wkhtmltopdf.exe"
	}
	if _, err := os.Stat(wkhtmltopdfPath); os.IsNotExist(err) {
		return fmt.Errorf("wkhtmltopdf executable not found at %s", wkhtmltopdfPath)
	}

	// Helper to create temp HTML files
	createTempHTML := func(content, prefix string) (string, error) {
		f, err := os.CreateTemp("", prefix+"*.html")
		if err != nil {
			return "", err
		}
		if _, err := f.Write([]byte(content)); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
		return f.Name(), nil
	}

	headerPath, err := createTempHTML(headerHTML, "header_")
	if err != nil {
		return err
	}
	defer os.Remove(headerPath)

	tempFooterFile, err := os.CreateTemp("", "export_pdf_footer_*.html")
	if err != nil {
		return fmt.Errorf("failed to create temporary footer file: %w", err)
	}
	defer os.Remove(tempFooterFile.Name())

	if _, err := tempFooterFile.Write([]byte(footerHTML)); err != nil {
		return fmt.Errorf("failed to write footer content to file: %w", err)
	}
	tempFooterFile.Close()

	coverPath, err := createTempHTML(coverHTML, "cover_")
	if err != nil {
		return err
	}
	defer os.Remove(coverPath)

	contentPath, err := createTempHTML(contentHTML, "content_")
	if err != nil {
		return err
	}
	defer os.Remove(contentPath)

	disclaimerPath, err := createTempHTML(disclaimerHTML, "disclaimer_")
	if err != nil {
		return err
	}
	defer os.Remove(disclaimerPath)

	tocPath, err := filepath.Abs("assets/document/templates/pdf/toc.xsl")
	if err != nil {
		return err
	}
	tocPath = filepath.ToSlash(tocPath)

	// Helper to generate PDF
	generatePDF := func(htmlPaths []string, outPDF string, top, bottom, left, right string, addTOC bool) error {
		args := []string{
			"--enable-local-file-access",
			"--margin-top", top,
			"--margin-bottom", bottom,
			"--margin-left", left,
			"--margin-right", right,
			"--header-html", headerPath,
			"--header-spacing", "5",
			"--footer-html", tempFooterFile.Name(),
			// "--footer-left", appURL,
			// "--footer-right", "[page]",
			// "--footer-font-size", "8",
			// "--footer-spacing", "5",
			// "--footer-line",
		}

		if addTOC {
			args = append(args,
				"toc",
				"--xsl-style-sheet", tocPath,
				"--toc-header-text", "Table of Contents",
				"--toc-level-indentation", "0",
			)
		}

		args = append(args, htmlPaths...)
		args = append(args, outPDF)

		cmd := exec.Command(wkhtmltopdfPath, args...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("wkhtmltopdf error: %s: %w", stderr.String(), err)
		}
		return nil
	}

	// Temporary PDFs
	coverPDF := "cover.pdf"
	bodyPDF := "body.pdf"

	// 1. Cover without header/footer
	if err := generatePDF([]string{coverPath}, coverPDF, "0", "0", "0", "0", false); err != nil {
		return err
	}

	// 2. TOC + Content + Disclaimer as a single PDF with header/footer
	if err := generatePDF([]string{contentPath, disclaimerPath}, bodyPDF, "25", "18", "10", "10", true); err != nil {
		return err
	}

	// 3. Merge cover + body
	if err := api.MergeCreateFile([]string{coverPDF, bodyPDF}, pdfFileName, false, nil); err != nil {
		return err
	}

	// Cleanup
	os.Remove(coverPDF)
	os.Remove(bodyPDF)

	return nil
}

func StripPTags(input string) string {
	re := regexp.MustCompile(`(?i)</?p[^>]*>`)
	return re.ReplaceAllString(input, "")
}

func CreateZipFile(zipPath string, files []string) error {
	newZipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		err = AddFileToZip(zipWriter, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}
