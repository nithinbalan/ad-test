package utility

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
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

func GenerateExcelFileWithSheets(dataBySheet map[string][][]interface{}) (*os.File, error) {
	f := excelize.NewFile()

	// Ensure Template is created and set as the first sheet
	if data, exists := dataBySheet["Template"]; exists {
		if err := createSheet(f, "Template", data); err != nil {
			return nil, err
		}
		sheetIndex, err := f.GetSheetIndex("Template")
		if err != nil {
			return nil, fmt.Errorf("failed to get Template sheet index: %v", err)
		}
		if sheetIndex != -1 {
			f.SetActiveSheet(sheetIndex)
		}
		delete(dataBySheet, "Template") // Remove Template from the map to avoid duplicate processing
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

// Helper function to create a sheet and populate it with data
func createSheet(f *excelize.File, sheetName string, data [][]interface{}) error {
	if _, err := f.NewSheet(sheetName); err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	headers := data[0]
	rows := data[1:]

	// Set headers as bold
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return err
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, boldStyle)
	}

	// Populate rows
	for rowIndex, row := range rows {
		for colIndex, value := range row {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Adjust column width
	for colIndex, header := range headers {
		maxWidth := len(fmt.Sprint(header))
		for _, row := range rows {
			if colIndex < len(row) {
				cellWidth := len(fmt.Sprint(row[colIndex]))
				if cellWidth > maxWidth {
					maxWidth = cellWidth
				}
			}
		}
		columnName := string(rune('A' + colIndex))
		f.SetColWidth(sheetName, columnName, columnName, float64(maxWidth+2))
	}

	return nil
}

func ConvertToInterfaceSlice(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
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

func GetOptionalStringValue(value interface{}) *string {
	strValue, ok := value.(string)
	if ok && strValue != "" {
		return &strValue
	}
	return nil
}

// GetOptionalUintValue safely retrieves an optional uint value from a map
func GetOptionalUintValue(value interface{}) *uint {
	if val, ok := value.(uint); ok {
		return &val
	}
	return nil
}

// control import
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

func ParseUint(value string) uint {
	val, _ := strconv.Atoi(value)
	return uint(val)
}

func ParseOptionalUint(value string) *uint {
	if value == "" {
		return nil
	}
	parsed := ParseUint(value)
	return &parsed
}

func ParseOptionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
