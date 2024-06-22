package utility

import (
	"fmt"
	"os"

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
