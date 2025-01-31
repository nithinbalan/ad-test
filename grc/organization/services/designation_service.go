package services

import (
	"errors"
	"fmt"
	"io"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

type DesignationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Designation, error)
	GetByID(id uint) (*models.Designation, error)
	Create(currentUserID uint, designation *models.Designation) (*models.Designation, error)
	Update(id uint, currentUserID uint, designation *models.Designation) (*models.Designation, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error

	GenerateDesignationTemplate() (*os.File, error)
	ValidateDesignationTemplate(file io.Reader) ([]map[string]interface{}, error)
	ValidateImportDesignation(input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportDesignations(input []map[string]interface{}, currentUserID uint) ([]uint, error)
}

// Service and implementations
type DesignationService struct {
	repo           repositories.DesignationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewDesignationService(repo repositories.DesignationRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) DesignationServiceInterface {
	return &DesignationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ds *DesignationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Designation, error) {
	return ds.repo.GetAll(sort, sortBy, keyword, status)
}

func (ds DesignationService) GetByID(id uint) (*models.Designation, error) {
	designation, err := ds.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return designation, nil
}

func (ds *DesignationService) Create(currentUserID uint, designation *models.Designation) (*models.Designation, error) {
	if exists, err := ds.repo.ExistsByTitle(designation.TitleEn, designation.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdDesignation, err := ds.repo.Create(designation)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ds.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Designation{}.TableName(),
			createdDesignation.ID,
			createdDesignation,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdDesignation, nil
}

func (ds *DesignationService) Update(id uint, currentUserID uint, updatedDesignation *models.Designation) (*models.Designation, error) {
	existingDesignation, err := ds.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedDesignation.TitleEn != existingDesignation.TitleEn || updatedDesignation.TitleAr != existingDesignation.TitleAr {
		if exists, err := ds.repo.ExistsByTitleExcludingID(id, updatedDesignation.TitleEn, updatedDesignation.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := ds.repo.Update(id, updatedDesignation)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingDesignation, updatedDesignation, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ds.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Designation{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (ds *DesignationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ds.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = ds.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := ds.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.Designation{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (ds *DesignationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ds.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ds.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.Designation{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (ds *DesignationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ds.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ds.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.Designation{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (ds *DesignationService) GenerateDesignationTemplate() (*os.File, error) {
	// Headers for the first sheet (Template)
	templateHeaders := []string{
		"TitleEn", "TitleAr",
	}
	templateData := [][]interface{}{} // Empty template data

	file, err := utility.GenerateExcelFileWithSheets(map[string][][]interface{}{
		"Template": append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
	})
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (ds *DesignationService) ValidateDesignationTemplate(file io.Reader) ([]map[string]interface{}, error) {
	// Read Excel file
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	// Define headers in the input format and map them to snake case
	inputHeaders := []string{"TitleEn", "TitleAr"}
	snakeCaseHeaders := []string{"title_en", "title_ar"}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}

	// Validate headers
	if !utility.ValidateHeaders(inputHeaders, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers in the Excel file"), "invalid_header", "Designation Template")
	}

	// Convert rows to input format
	var input []map[string]interface{}
	for _, row := range rows[1:] { // Skip header row
		record := make(map[string]interface{})
		for i, header := range snakeCaseHeaders {
			if i < len(row) {
				record[header] = row[i]
			} else {
				record[header] = "" // Pad missing columns
			}
		}
		input = append(input, record)
	}

	validatedData, err := ds.validateDesignations(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (ds *DesignationService) ValidateImportDesignation(input []map[string]interface{}) ([]map[string]interface{}, error) {

	validatedData, err := ds.validateDesignations(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (ds *DesignationService) ImportDesignations(input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validatedData, err := ds.validateDesignations(input)
	if err != nil {
		return nil, err
	}

	for _, record := range validatedData {
		if errors, ok := record["errors"].(map[string]interface{}); ok && len(errors) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data in the import file"), "invalid_data", "Invalid Data")
		}
	}

	var createdDesignationIDs []uint
	for _, record := range validatedData {
		designationFormDTO := &models.Designation{
			TitleEn: record["title_en"].(string),
			TitleAr: record["title_ar"].(string),
			Status:  "active",
		}

		designation, err := ds.Create(currentUserID, designationFormDTO)
		if err != nil {
			return nil, fmt.Errorf("failed to create designation for record: %+v, error: %v", record, err)
		}

		createdDesignationIDs = append(createdDesignationIDs, designation.ID)
	}

	return createdDesignationIDs, nil
}

func (ds *DesignationService) validateDesignations(input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{"title_en", "title_ar"}

	// Fetch existing designations for validation
	designations, err := ds.repo.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching designations: %v", err)
	}

	// Create maps of existing designation titles for validation
	designationTitlesEn := make(map[string]bool)
	designationTitlesAr := make(map[string]bool)
	for _, designation := range designations {
		designationTitlesEn[strings.TrimSpace(designation.TitleEn)] = true
		designationTitlesAr[strings.TrimSpace(designation.TitleAr)] = true
	}

	// Create maps to track duplicates within the input sheet
	inputTitlesEn := make(map[string]bool)
	inputTitlesAr := make(map[string]bool)
	var validatedData []map[string]interface{}

	for _, record := range input {
		rowErrors := make(map[string]interface{})
		processedRecord := make(map[string]interface{})

		// Copy original data to processed record
		for key, value := range record {
			processedRecord[key] = value
		}

		// Validate fields
		for _, header := range headers {
			value, ok := record[header].(string)
			trimmedValue := strings.TrimSpace(value)

			// Required field check
			if !ok || trimmedValue == "" {
				rowErrors[header] = "Field is required"
				continue
			}

			switch header {
			case "title_en":
				// Check for duplicates within the input sheet
				if inputTitlesEn[trimmedValue] {
					rowErrors["title_en"] = "Duplicate designation title in sheet"
					continue
				}

				// Check for duplicates in the database
				if designationTitlesEn[trimmedValue] {
					rowErrors["title_en"] = "Duplicate designation title in database"
					continue
				}

				// Mark the title as seen in the input sheet
				inputTitlesEn[trimmedValue] = true

			case "title_ar":
				// Check for duplicates within the input sheet
				if inputTitlesAr[trimmedValue] {
					rowErrors["title_ar"] = "Duplicate designation title in sheet"
					continue
				}

				// Check for duplicates in the database
				if designationTitlesAr[trimmedValue] {
					rowErrors["title_ar"] = "Duplicate designation title in database"
					continue
				}

				// Mark the title as seen in the input sheet
				inputTitlesAr[trimmedValue] = true
			}
		}

		// Add errors or success record
		if len(rowErrors) == 0 {
			processedRecord["errors"] = nil
		} else {
			for _, header := range headers {
				if _, exists := rowErrors[header]; !exists {
					rowErrors[header] = "" // Empty string for other fields
				}
			}
			processedRecord["errors"] = rowErrors
		}
		validatedData = append(validatedData, processedRecord)
	}

	return validatedData, nil
}
