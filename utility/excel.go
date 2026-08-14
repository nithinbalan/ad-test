package utility

import (
	"fmt"
	"io"
	"log"
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

func GenerateExcelFileWithSheetsAndDropDown(dataBySheet map[string][][]interface{}, dropDownConfig map[string]bool, dateColumnConfigs ...map[string]string) (*os.File, error) {
	f := excelize.NewFile()
	dateColumnConfig := map[string]string{}
	if len(dateColumnConfigs) > 0 && dateColumnConfigs[0] != nil {
		dateColumnConfig = dateColumnConfigs[0]
	}

	// Ensure Template is created and set as the first sheet
	if data, exists := dataBySheet["Template"]; exists {
		delete(dataBySheet, "Template") // Remove Template from the map to avoid duplicate processing
		if err := createSheetWithDropDown(f, "Template", data, dataBySheet, dropDownConfig, dateColumnConfig); err != nil {
			return nil, err
		}
		sheetIndex, err := f.GetSheetIndex("Template")
		if err != nil {
			return nil, fmt.Errorf("failed to get Template sheet index: %v", err)
		}
		if sheetIndex != -1 {
			f.SetActiveSheet(sheetIndex)
		}
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

func toAlphaString(colIndex int) string {
	col := colIndex + 1
	result := ""
	for col > 0 {
		col--
		result = string(rune('A'+(col%26))) + result
		col /= 26
	}
	return result
}

func createSheetWithDropDown(f *excelize.File, sheetName string, data [][]interface{}, columnToSheet map[string][][]interface{}, dropDownConfig map[string]bool, dateColumnConfig map[string]string) error {
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
	if err != nil {
		return fmt.Errorf("failed to create bold style: %v", err)
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

		if apply, exists := dropDownConfig[headerStr]; !exists || !apply {
			log.Printf("skipping the drop down for: %s", headerStr)
			continue
		}

		optionsData, exists := columnToSheet[headerStr]
		if !exists {
			log.Printf("No data for dropdown column: %s", headerStr)
			continue
		}

		// Flatten dropdown values
		dropList := []string{}
		for _, row := range optionsData[1:] { // Skip header row
			if len(row) > 0 {
				if headerStr == "Clause" && len(row) > 1 {
					clauseNo, ok1 := row[0].(string)
					titleEn, ok2 := row[1].(string)
					if ok1 && ok2 && clauseNo != "" && titleEn != "" {
						dropList = append(dropList, titleEn)
					}
				} else if val, ok := row[0].(string); ok && val != "" {
					dropList = append(dropList, val)
				}
			}
		}

		if len(dropList) == 0 {
			log.Printf("No valid dropdown values for column %s", headerStr)
			continue // No values to create dropdown
		}

		// Convert column index to letter
		colLetter := toAlphaString(colIndex)

		// Set data validation for rows 2 to 100
		dv := excelize.NewDataValidation(true)
		dv.Sqref = fmt.Sprintf("%s2:%s100", colLetter, colLetter)
		dv.SetDropList(dropList)
		dv.ShowErrorMessage = true
		errTitle := "Invalid Input"
		errMsg := "Please select a value from the dropdown."
		dv.ErrorTitle = &errTitle
		dv.Error = &errMsg

		if err := f.AddDataValidation(sheetName, dv); err != nil {
			return fmt.Errorf("failed to add dropdown to column %s: %v", headerStr, err)
		}
	}

	for colIndex, header := range headers {
		headerStr, ok := header.(string)
		if !ok {
			continue
		}

		dateFormat, apply := dateColumnConfig[headerStr]
		if !apply {
			continue
		}

		colLetter := toAlphaString(colIndex)
		startCell := fmt.Sprintf("%s2", colLetter)
		endCell := fmt.Sprintf("%s100", colLetter)

		dateStyle, err := f.NewStyle(&excelize.Style{
			CustomNumFmt: &dateFormat,
			Alignment: &excelize.Alignment{
				Horizontal: "center",
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create date style for column %s: %v", headerStr, err)
		}
		if err := f.SetColWidth(sheetName, colLetter, colLetter, 22); err != nil {
			return fmt.Errorf("failed to set width for date column %s: %v", headerStr, err)
		}
		if err := f.SetCellStyle(sheetName, startCell, endCell, dateStyle); err != nil {
			return fmt.Errorf("failed to apply date style to column %s: %v", headerStr, err)
		}

		dv := excelize.NewDataValidation(true)
		dv.Sqref = fmt.Sprintf("%s:%s", startCell, endCell)
		if err := dv.SetRange(1, 2958465, excelize.DataValidationTypeDate, excelize.DataValidationOperatorBetween); err != nil {
			return fmt.Errorf("failed to configure date validation for column %s: %v", headerStr, err)
		}
		dv.SetError(excelize.DataValidationErrorStyleStop, "Invalid Date and Time", "Please enter a valid date and time in YYYY-MM-DD HH:mm format.")
		if err := f.AddDataValidation(sheetName, dv); err != nil {
			return fmt.Errorf("failed to add date validation to column %s: %v", headerStr, err)
		}
	}

	return nil
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

func MapRowsToRecords(headers []string, rows [][]string) []map[string]interface{} {
	var records []map[string]interface{}

	for _, row := range rows {
		record := make(map[string]interface{})
		for i, header := range headers {
			if i < len(row) {
				record[header] = row[i]
			} else {
				record[header] = ""
			}
		}
		records = append(records, record)
	}

	return records
}
