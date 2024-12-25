package utility

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GenerateExcelFile(headers []string, data [][]interface{}) (*os.File, error) {
	f := excelize.NewFile()

	// Create a new sheet
	sheet := "Sheet1"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

	// Set headers as bold
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, boldStyle)
	}

	// Populate data
	for rowIndex, row := range data {
		for colIndex, value := range row {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2) // +2 because rowIndex starts from 0 and first row is header
			f.SetCellValue(sheet, cell, value)
		}
	}

	// Adjust column width according to content
	for colIndex, header := range headers {
		maxWidth := len(header) // Start with header width
		for _, rowData := range data {
			if colIndex < len(rowData) {
				cellWidth := len(fmt.Sprint(rowData[colIndex]))
				if cellWidth > maxWidth {
					maxWidth = cellWidth
				}
			}
		}
		columnName := string(rune('A' + colIndex))
		f.SetColWidth(sheet, columnName, columnName, float64(maxWidth+2)) // Add some padding
	}

	// Save the file
	tempFile, err := os.CreateTemp("", "export_*.xlsx")
	if err != nil {
		return nil, err
	}

	err = f.SaveAs(tempFile.Name())
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func ParseExcelFile(file multipart.File, mandatoryFields []string) ([]map[string]string, error) {
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Excel file: %v", err)
	}

	// Get the first sheet
	sheetName := xlsx.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("the sheet name cannot be blank")
	}

	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from the sheet: %v", err)
	}

	// Validate the template has data
	if len(rows) < 2 {
		return nil, fmt.Errorf("no data in the file")
	}

	// Map the headers to row values
	headers := rows[0]
	var data []map[string]string
	var validationErrors []string
	for rowIndex, row := range rows[1:] {
		rowData := make(map[string]string)
		for i, cell := range row {
			rowData[headers[i]] = cell
		}

		// Validate mandatory fields
		for _, field := range mandatoryFields {
			if value, exists := rowData[field]; !exists || value == "" {
				validationErrors = append(validationErrors, fmt.Sprintf("mandatory field '%s' is missing or empty in row %d", field, rowIndex+2))
			}
		}

		data = append(data, rowData)
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf(strings.Join(validationErrors, "; "))
	}

	return data, nil
}

func ReadExcelFile(file io.Reader) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func ValidateHeaders(expected, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i, header := range expected {
		if header != actual[i] {
			return false
		}
	}
	return true
}

func ParseCommaSeparatedValues(input string) []string {
	return strings.Split(strings.TrimSpace(input), ";")
}

// GetOptionalString retrieves a pointer to a string from an interface.
// If the input is nil or not a string, it returns nil.
func GetOptionalString(value interface{}) *string {
	if value == nil {
		return nil
	}

	// Assert the value as a string
	str, ok := value.(string)
	if !ok || str == "" {
		return nil
	}

	return &str
}

func GetOptionalStringValue(value interface{}) *string {
	strValue, ok := value.(string)
	if ok && strValue != "" {
		return &strValue
	}
	return nil
}
