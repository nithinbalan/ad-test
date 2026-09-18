package utility

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"jethur/assets"
	"jethur/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/xuri/excelize/v2"
)

func GeneratePDFContent(dto interface{}, templatePaths []string, logoImagePath string, coverBGImagePath string) (string, error) {
	// Read and encode fonts
	lexandDecaFont100, err := fontToDataURL("fonts/LexendDeca-Thin.woff")
	if err != nil {
		return "", err
	}

	lexandDecaFont200, err := fontToDataURL("fonts/LexendDeca-ExtraLight.woff")
	if err != nil {
		return "", err
	}

	lexandDecaFont300, err := fontToDataURL("fonts/LexendDeca-Light.woff")
	if err != nil {
		return "", err
	}

	lexandDecaFont400, err := fontToDataURL("fonts/LexendDeca-Regular.woff")
	if err != nil {
		return "", err
	}

	lexandDecaFont500, err := fontToDataURL("fonts/LexendDeca-Medium.woff")
	if err != nil {
		return "", err
	}

	lexandDecaFont600, err := fontToDataURL("fonts/LexendDeca-SemiBold.woff")
	if err != nil {
		return "", err
	}

	// Prepare the data to include image paths
	data := map[string]interface{}{
		"LogoImagePath":     template.HTML(logoImagePath),
		"CoverBGImagePath":  template.HTML(coverBGImagePath),
		"LexandDecaFont100": template.URL(lexandDecaFont100),
		"LexandDecaFont200": template.URL(lexandDecaFont200),
		"LexandDecaFont300": template.URL(lexandDecaFont300),
		"LexandDecaFont400": template.URL(lexandDecaFont400),
		"LexandDecaFont500": template.URL(lexandDecaFont500),
		"LexandDecaFont600": template.URL(lexandDecaFont600),
		"Data":              dto,
		"ApiPath":           config.GetEnv("API_URL", ""),
	}

	templates := template.New("pdf").Funcs(template.FuncMap{
		"formatDate": func(date string) string { return date },
		"safeHTML":   func(content string) template.HTML { return template.HTML(content) },
		"stripHTML": func(content interface{}) string {
			switch v := content.(type) {
			case string:
				return StripHTMLTags(v)
			case *string:
				if v == nil {
					return ""
				}
				return StripHTMLTags(*v)
			case template.HTML:
				return StripHTMLTags(string(v))
			}
			return ""
		},
		"add": func(a, b int) int { return a + b },
		"boolToYesNo": func(b *bool) string {
			if b == nil {
				return ""
			}
			if *b {
				return "Yes"
			}
			return "No"
		},
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

var statusBadgeHexColorPattern = regexp.MustCompile(`^[0-9A-Fa-f]{6}$`)

func ResolveStatusBadgeColors(statusColorCode string, statusType string, statusLabel string) (string, string, string) {
	baseColor := defaultStatusBadgeColor(statusType, statusLabel)
	if normalizedColor, ok := normalizeStatusBadgeHexColor(statusColorCode); ok {
		baseColor = normalizedColor
	}

	textColor := darkenStatusBadgeHexColor(baseColor, 20)
	borderColor := darkenStatusBadgeHexColor(baseColor, 8)
	backgroundColor := blendStatusBadgeHexColorWithWhite(baseColor, 88)

	return textColor, borderColor, backgroundColor
}

func defaultStatusBadgeColor(statusType string, statusLabel string) string {
	switch strings.ToLower(strings.TrimSpace(statusType)) {
	case "new", "reopen":
		return "#2563EB"
	case "in_progress":
		return "#D97706"
	case "completed", "close":
		return "#1F7D58"
	case "on_hold":
		return "#6B7280"
	case "terminated":
		return "#B91C1C"
	}

	if strings.Contains(strings.ToLower(strings.TrimSpace(statusLabel)), "overdue") {
		return "#DF3F47"
	}

	return "#DF3F47"
}

func normalizeStatusBadgeHexColor(colorValue string) (string, bool) {
	trimmedColor := strings.TrimPrefix(strings.TrimSpace(colorValue), "#")
	if len(trimmedColor) == 3 {
		trimmedColor = strings.Repeat(string(trimmedColor[0]), 2) +
			strings.Repeat(string(trimmedColor[1]), 2) +
			strings.Repeat(string(trimmedColor[2]), 2)
	}

	if len(trimmedColor) != 6 || !statusBadgeHexColorPattern.MatchString(trimmedColor) {
		return "", false
	}

	return "#" + strings.ToUpper(trimmedColor), true
}

func parseStatusBadgeHexColor(colorValue string) (int, int, int, bool) {
	normalizedColor, ok := normalizeStatusBadgeHexColor(colorValue)
	if !ok {
		return 0, 0, 0, false
	}

	parsedValue, err := strconv.ParseInt(strings.TrimPrefix(normalizedColor, "#"), 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}

	return int((parsedValue >> 16) & 0xFF), int((parsedValue >> 8) & 0xFF), int(parsedValue & 0xFF), true
}

func darkenStatusBadgeHexColor(colorValue string, percentage int) string {
	red, green, blue, ok := parseStatusBadgeHexColor(colorValue)
	if !ok {
		return "#DF3F47"
	}

	darkenPercent := clampStatusBadgeColorPercentage(percentage)
	red = (red * (100 - darkenPercent)) / 100
	green = (green * (100 - darkenPercent)) / 100
	blue = (blue * (100 - darkenPercent)) / 100

	return fmt.Sprintf("#%02X%02X%02X", red, green, blue)
}

func blendStatusBadgeHexColorWithWhite(colorValue string, whitePercentage int) string {
	red, green, blue, ok := parseStatusBadgeHexColor(colorValue)
	if !ok {
		return "#FDF0F2"
	}

	whitePercent := clampStatusBadgeColorPercentage(whitePercentage)
	red = ((100-whitePercent)*red + whitePercent*255) / 100
	green = ((100-whitePercent)*green + whitePercent*255) / 100
	blue = ((100-whitePercent)*blue + whitePercent*255) / 100

	return fmt.Sprintf("#%02X%02X%02X", red, green, blue)
}

func clampStatusBadgeColorPercentage(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
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

	// tocPath, err := filepath.Abs("assets/document/templates/pdf/toc.xsl")
	// if err != nil {
	// 	return err
	// }
	// tocPath = filepath.ToSlash(tocPath)

	basePath, _ := os.Getwd()
	tocPath := filepath.Join(basePath, "assets", "document", "templates", "pdf", "toc.xsl")
	println(tocPath)

	// Convert to double backslashes
	tocPath = strings.ReplaceAll(tocPath, "\\", "\\\\")

	addTOC := true
	if _, err := os.Stat(tocPath); os.IsNotExist(err) {
		addTOC = false
		fmt.Printf("TOC XSL not found at %s, skipping TOC generation\n", tocPath)
	}

	fmt.Println(basePath)
	fmt.Println(tocPath)

	// Helper to generate PDF
	generatePDF := func(htmlPaths []string, outPDF string, top, bottom, left, right string, withTOC bool) error {
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

		// if addTOC {
		// 	args = append(args,
		// 		"toc",
		// 		"--xsl-style-sheet", tocPath,
		// 		"--toc-header-text", "Table of Contents",
		// 		"--toc-level-indentation", "0",
		// 	)
		// }

		if withTOC && addTOC {
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

func StripHTMLTags(input string) string {
	// Convert common line-breaking tags to new lines before stripping tags
	input = strings.ReplaceAll(input, "<br>", "\n")
	input = strings.ReplaceAll(input, "<br/>", "\n")
	input = strings.ReplaceAll(input, "<br />", "\n")
	input = strings.ReplaceAll(input, "</p>", "\n")
	input = strings.ReplaceAll(input, "</div>", "\n")
	input = strings.ReplaceAll(input, "</li>", "\n")
	input = strings.ReplaceAll(input, "<li>", "• ")

	// Remove HTML entities
	input = strings.ReplaceAll(input, "&nbsp;", " ")
	input = strings.ReplaceAll(input, "&lt;", "<")
	input = strings.ReplaceAll(input, "&gt;", ">")
	input = strings.ReplaceAll(input, "&amp;", "&")

	// Remove all remaining HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	input = re.ReplaceAllString(input, "")

	// Clean up whitespace
	input = strings.Join(strings.Fields(input), " ")

	return input
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

// ZipEntry is a file on disk and its path inside the archive.
type ZipEntry struct {
	SourcePath string
	ZipPath    string
}

// CreateZipFileWithFolders zips each entry at its own ZipPath, keeping folders.
func CreateZipFileWithFolders(zipPath string, entries []ZipEntry) error {
	newZipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	for _, entry := range entries {
		if err := AddFileToZipAtPath(zipWriter, entry.SourcePath, entry.ZipPath); err != nil {
			zipWriter.Close()
			return err
		}
	}

	// Not deferred: Close writes the central directory and can fail
	return zipWriter.Close()
}

// AddFileToZipAtPath is AddFileToZip storing the file under zipName.
func AddFileToZipAtPath(zipWriter *zip.Writer, filename string, zipName string) error {
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

	// Zip entry names always use forward slashes
	header.Name = filepath.ToSlash(zipName)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

func GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent string) error {
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

	// Create temp files
	coverPath, err := createTempHTML(coverContent, "cover_")
	if err != nil {
		return err
	}
	defer os.Remove(coverPath)

	contentPath, err := createTempHTML(htmlContent, "content_")
	if err != nil {
		return err
	}
	defer os.Remove(contentPath)

	headerPath, err := createTempHTML(headerContent, "header_")
	if err != nil {
		return err
	}
	defer os.Remove(headerPath)

	footerPath, err := createTempHTML(footerContent, "footer_")
	if err != nil {
		return err
	}
	defer os.Remove(footerPath)

	// Helper to generate PDF
	generatePDF := func(htmlPaths []string, outPDF string, top, bottom, left, right string, addHeaderFooter bool) error {
		args := []string{
			"--enable-local-file-access",
			"--margin-top", top,
			"--margin-bottom", bottom,
			"--margin-left", left,
			"--margin-right", right,
		}

		if addHeaderFooter {
			args = append(args,
				"--header-html", headerPath,
				"--header-spacing", "5",
				"--footer-html", footerPath,
				"--footer-spacing", "5",
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

	// 1. Cover (no header/footer, no margins)
	if err := generatePDF([]string{coverPath}, coverPDF, "0", "0", "0", "0", false); err != nil {
		return err
	}

	// 2. Body (content with header/footer and margins)
	if err := generatePDF([]string{contentPath}, bodyPDF, "20", "20", "0", "0", true); err != nil {
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

func GenerateExcelFileWithSheetsAndDropDownAnswer(dataBySheet map[string][][]interface{}, dropDownConfig map[string]bool, readOnlyConfig map[string]bool) (*os.File, error) {
	f := excelize.NewFile()

	if intro, exists := dataBySheet["Introduction"]; exists {
		if err := createSheet(f, "Introduction", intro); err != nil {
			return nil, err
		}
		// Set Introduction as the active sheet
		idx, _ := f.GetSheetIndex("Introduction")
		if idx > -1 {
			f.SetActiveSheet(idx)
		}
		// Remove to avoid creating again in loop
		delete(dataBySheet, "Introduction")
	}

	// Ensure Template is created and set as the first sheet
	if data, exists := dataBySheet["Template"]; exists {
		// delete(dataBySheet, "Template") // Remove Template from the map to avoid duplicate processing
		if err := createSheetWithDropDownAnswer(f, "Template", data, dataBySheet, dropDownConfig, readOnlyConfig); err != nil {
			return nil, err
		}
		// sheetIndex, err := f.GetSheetIndex("Template")
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to get Template sheet index: %v", err)
		// }
		// if sheetIndex != -1 {
		// 	f.SetActiveSheet(sheetIndex)
		// }
	}

	// Create the remaining sheets
	for sheetName, data := range dataBySheet {
		if err := createSheet(f, sheetName, data); err != nil {
			return nil, err
		}
	}

	if index, err := f.GetSheetIndex("Sheet1"); err == nil && index != -1 {
		f.DeleteSheet("Sheet1")
	}

	// Save the file
	tempFile, err := os.CreateTemp("", "template_*.xlsx")
	if err != nil {
		return nil, err
	}

	if err := f.SaveAs(tempFile.Name()); err != nil {
		return nil, err
	}

	return tempFile, nil
}

func createSheetWithDropDownAnswer(f *excelize.File, sheetName string, data [][]interface{}, columnToSheet map[string][][]interface{}, dropDownConfig map[string]bool, readOnlyConfig map[string]bool) error {
	if _, err := f.NewSheet(sheetName); err != nil {
		return fmt.Errorf("failed to create sheet %s: %v", sheetName, err)
	}

	if len(data) == 0 {
		return nil
	}

	headers := data[0]
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	lockedStyle, _ := f.NewStyle(&excelize.Style{Protection: &excelize.Protection{Locked: true}})
	unlockedStyle, _ := f.NewStyle(&excelize.Style{Protection: &excelize.Protection{Locked: false}})
	if err != nil {
		return fmt.Errorf("failed to create bold style: %v", err)
	}

	for row := 1; row < len(data); row++ {
		for col := 0; col < len(headers); col++ {

			colName := toAlphaString(col)
			cell := fmt.Sprintf("%s%d", colName, row+1)
			headerText := headers[col].(string)

			f.SetCellValue(sheetName, cell, data[row][col])

			// APPLY LOCK / UNLOCK BASED ON CONFIG
			if readOnlyConfig[strings.TrimSpace(headerText)] {
				f.SetCellStyle(sheetName, cell, cell, lockedStyle)
			} else {
				f.SetCellStyle(sheetName, cell, cell, unlockedStyle)
			}
		}
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", toAlphaString(i))
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, boldStyle)
	}

	for colIndex, header := range headers {
		headerStr, ok := header.(string)
		if !ok {
			continue
		}

		headerKey := strings.TrimSpace(headerStr)

		apply, exists := dropDownConfig[headerKey]
		if !exists {
			log.Printf("Dropdown config missing for header: %q", headerKey)
			continue
		}
		if !apply {
			log.Printf("Dropdown disabled for: %q", headerKey)
			continue
		}

		optionsData, exists := columnToSheet[headerKey]
		if !exists || len(optionsData) <= 1 {
			log.Printf("No data for dropdown source: %s", headerKey)
			continue
		}

		colLetter := toAlphaString(colIndex)

		// Create Data Validation
		dv := excelize.NewDataValidation(true)
		dv.Sqref = fmt.Sprintf("%s2:%s1000", colLetter, colLetter)

		// Range reference (bypasses 255-char limit)
		var sourceSheet, valueCol string
		if headerKey == "Clause" {
			sourceSheet = "Clause"
			valueCol = "B"
		} else {
			sourceSheet = headerKey
			valueCol = "A"
		}
		lastRow := len(optionsData) // includes header → data from row 2
		dv.Formula1 = fmt.Sprintf("=%s!$%s$2:$%s$%d", sourceSheet, valueCol, valueCol, lastRow+1)
		dv.Type = "list"

		errorTitle := "Invalid Selection"
		errorMsg := fmt.Sprintf("Please select a valid %s from the list.", strings.ToLower(headerKey))
		// promptTitle := headerKey
		// promptMsg := fmt.Sprintf("Select from available %ss", strings.ToLower(headerKey))

		dv.ErrorTitle = &errorTitle
		dv.Error = &errorMsg
		// dv.PromptTitle = &promptTitle
		// dv.Prompt = &promptMsg

		dv.ShowErrorMessage = true
		dv.ShowInputMessage = true

		if err := f.AddDataValidation(sheetName, dv); err != nil {
			return fmt.Errorf("failed to add dropdown for %s: %v", headerKey, err)
		}

		log.Printf("Dropdown added: %s → %s!$%s$2:$%s$%d", headerKey, sourceSheet, valueCol, valueCol, lastRow+1)
	}

	f.ProtectSheet(sheetName, &excelize.SheetProtectionOptions{
		SelectLockedCells:   false,
		SelectUnlockedCells: true,
		FormatCells:         false,
		FormatColumns:       false,
		FormatRows:          false,
		InsertColumns:       false,
		InsertRows:          false,
		InsertHyperlinks:    false,
		EditObjects:         false,
		EditScenarios:       false,
	})

	return nil
}

func ReadExcelFileWithIntroduction(file io.Reader) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(1))
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Helper function to convert font file to base64 data URL
func fontToDataURL(fontPath string) (string, error) {
	fontData, err := assets.FS.ReadFile(fontPath)
	if err != nil {
		return "", fmt.Errorf("failed to read font %s: %w", fontPath, err)
	}
	encoded := base64.StdEncoding.EncodeToString(fontData)
	return "data:font/woff;charset=utf-8;base64," + encoded, nil
}

func ExtractBulletPointsFromHTML(htmlContent string) []string {

	// Convert common line-breaking tags to new lines
	replacements := map[string]string{
		"<br>":   "\n",
		"<br/>":  "\n",
		"<br />": "\n",
		"</p>":   "\n",
		"</div>": "\n",
		"</li>":  "\n",
		"<li>":   "• ",
	}

	for old, newVal := range replacements {
		htmlContent = strings.ReplaceAll(htmlContent, old, newVal)
	}

	// Remove HTML entities
	htmlContent = strings.ReplaceAll(htmlContent, "&nbsp;", " ")

	// Remove all remaining HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	htmlContent = re.ReplaceAllString(htmlContent, "")

	// Split into lines
	lines := strings.Split(htmlContent, "\n")

	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// Remove bullet symbol if present
		line = strings.TrimPrefix(line, "•")
		line = strings.TrimSpace(line)

		result = append(result, line)
	}

	return result
}
