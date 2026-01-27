package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	violationReportRepositories "jethur/grc/violation-report/repositories"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ViolationReportInvestigationServiceInterface interface {
	GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportInvestigationListQueryParamsDTO) (*violationReportDTOs.ViolationReportInvestigationListResponseDTO, error)
	GetFullDetailsByID(violationReportID uint, tx ...utility.DatabaseTransaction) (*violationReportDTOs.ViolationReportInvestigationDetailDTO, error)
	GetByID(violationReportID uint) (*models.ViolationReportInvestigation, error)
	Create(violationReportID uint, currentUserID uint, violationReportInvestigationFormDTO *violationReportDTOs.ViolationReportInvestigationFormDTO) (uint, error)
	Update(violationReportID uint, currentUserID uint, violationReportInvestigationFormDTO *violationReportDTOs.ViolationReportInvestigationFormDTO) error
	Delete(violationReportID uint, currentUserID uint) (uint, error)
	GetDocumentPath(violationReportID uint, token string) (string, error)

	UpdateInvestigationNoteAndObservations(violationReportID uint, currentUserID uint, violationReportInvestigationFormDTO *violationReportDTOs.ViolationReportInvestigationNoteAndObservationDTO) error
}

// Service and implementations
type ViolationReportInvestigationService struct {
	repo               violationReportRepositories.ViolationReportInvestigationRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	violationReportSvc ViolationReportServiceInterface
	userSvc            organizationServices.UserServiceInterface
}

func NewViolationReportInvestigationService(
	repo violationReportRepositories.ViolationReportInvestigationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	violationReportSvc ViolationReportServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) ViolationReportInvestigationServiceInterface {
	return &ViolationReportInvestigationService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		violationReportSvc: violationReportSvc,
		userSvc:            userSvc,
	}
}

func (vris *ViolationReportInvestigationService) GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportInvestigationListQueryParamsDTO) (*violationReportDTOs.ViolationReportInvestigationListResponseDTO, error) {
	if _, err := vris.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}
	return vris.repo.GetAll(violationReportID, params)
}

func (vris *ViolationReportInvestigationService) GetByID(violationReportID uint) (*models.ViolationReportInvestigation, error) {
	violationReportInvestigation, err := vris.repo.GetByID(violationReportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, violationReportID)
		}
		return nil, err
	}
	return violationReportInvestigation, nil
}

func (vris *ViolationReportInvestigationService) GetFullDetailsByID(violationReportID uint, tx ...utility.DatabaseTransaction) (*violationReportDTOs.ViolationReportInvestigationDetailDTO, error) {
	if _, err := vris.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}
	// Initialize your DTO
	var violationReportInvestigationDTO violationReportDTOs.ViolationReportInvestigationDetailDTO

	violationReportInvestigation, err := vris.repo.GetFullDetailsByID(violationReportID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, violationReportID)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(violationReportInvestigation, &violationReportInvestigationDTO)
	if err != nil {
		return nil, err
	}

	// todo :- also add extra fields which required in details screen
	investigationDate := violationReportInvestigation.InvestigationDate.Format("2006-01-02")
	violationReportInvestigationDTO.InvestigationDateString = investigationDate

	violationReportInvestigationDTO.CreatedAt = violationReportInvestigation.CreatedAt
	violationReportInvestigationDTO.UpdatedAt = violationReportInvestigation.UpdatedAt

	return &violationReportInvestigationDTO, nil
}

func (vris *ViolationReportInvestigationService) Create(violationReportID uint, currentUserID uint, violationReportInvestigationFormDTO *violationReportDTOs.ViolationReportInvestigationFormDTO) (uint, error) {
	tx := vris.repo.BeginTransaction()
	if _, err := vris.violationReportSvc.GetByID(violationReportID); err != nil {
		return 0, err
	}
	if err := vris.validateViolationReportInvestigationFormDTO(violationReportInvestigationFormDTO); err != nil {
		return 0, err
	}
	dates, err := vris.parseViolationReportInvestigationDates(violationReportInvestigationFormDTO)
	if err != nil {
		return 0, err
	}

	// Create Violation Report Investigation
	violationReportInvestigation := &models.ViolationReportInvestigation{
		ViolationReportID:                  violationReportID,
		InvestigationDate:                  dates.InvestigationDate, // parsed from string
		StepsTaken:                         violationReportInvestigationFormDTO.StepsTaken,
		SuggestedRecommendation:            violationReportInvestigationFormDTO.SuggestedRecommendation,
		FindingSummary:                     violationReportInvestigationFormDTO.FindingSummary,
		DisciplinaryProcedure:              violationReportInvestigationFormDTO.DisciplinaryProcedure,
		IsDisciplinaryRegulatoryCompliance: violationReportInvestigationFormDTO.IsDisciplinaryRegulatoryCompliance,
		ActualRecommendation:               violationReportInvestigationFormDTO.ActualRecommendation,
		CreatedBYID:                        currentUserID,
	}

	id, err := vris.repo.Create(violationReportInvestigation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save Departments
	var departments []*models.ViolationReportInvestigationDepartment
	for _, departmentId := range violationReportInvestigationFormDTO.DepartmentIds {
		department := &models.ViolationReportInvestigationDepartment{
			ViolationReportInvestigationID: id,
			DepartmentID:                   *departmentId,
			CreatedBYID:                    currentUserID,
		}
		departments = append(departments, department)
	}

	if len(departments) > 0 {
		if err := vris.repo.CreateDepartments(departments, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Documents
	var documents []*models.ViolationReportInvestigationDocument
	// Iterate over each document in the array
	for _, doc := range violationReportInvestigationFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		// Construct the URL for the document
		url := utility.GenerateDocumentURL("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "investigations", strconv.FormatUint(uint64(id), 10), "documents")

		// Create a new Violation Report Investigation Document instance
		document := &models.ViolationReportInvestigationDocument{
			ViolationReportInvestigationID: id,
			URL:                            url,
			Name:                           doc.FileName,
			Size:                           doc.Size,
			Token:                          encryptedFilename,
			OriginalName:                   doc.OriginalName,
			Type:                           doc.Type,
			CreatedByID:                    currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		if err := vris.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "investigations", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		existingViolationReportInvestigation, err := vris.GetFullDetailsByID(violationReportID)
		if err != nil {
			utility.LogError("Failed to fetch violation report investigation", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = vris.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportInvestigation{}.TableName(),
			id,
			existingViolationReportInvestigation,
			currentUserID,
			"violation_report", "violation_report_register", "investigation", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ViolationReportInvestigation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (vris *ViolationReportInvestigationService) Update(violationReportID uint, currentUserID uint, violationReportInvestigationFormDTO *violationReportDTOs.ViolationReportInvestigationFormDTO) error {
	tx := vris.repo.BeginTransaction()
	if _, err := vris.violationReportSvc.GetByID(violationReportID); err != nil {
		return err
	}
	if err := vris.validateViolationReportInvestigationFormDTO(violationReportInvestigationFormDTO); err != nil {
		return err
	}
	dates, err := vris.parseViolationReportInvestigationDates(violationReportInvestigationFormDTO)
	if err != nil {
		return err
	}
	oldViolationReportInvestigation, err := vris.GetFullDetailsByID(violationReportID)
	if err != nil {
		return err
	}
	// Update existing Violation Report Investigation
	existingViolationReportInvestigation, err := vris.repo.GetByID(violationReportID)
	if err != nil {
		return err
	}
	if existingViolationReportInvestigation == nil {
		return errors.New("existing violation report investigation not found")
	}

	existingViolationReportInvestigation.ViolationReportID = violationReportID
	existingViolationReportInvestigation.InvestigationDate = dates.InvestigationDate
	existingViolationReportInvestigation.StepsTaken = violationReportInvestigationFormDTO.StepsTaken
	existingViolationReportInvestigation.SuggestedRecommendation = violationReportInvestigationFormDTO.SuggestedRecommendation
	existingViolationReportInvestigation.FindingSummary = violationReportInvestigationFormDTO.FindingSummary
	existingViolationReportInvestigation.DisciplinaryProcedure = violationReportInvestigationFormDTO.DisciplinaryProcedure
	existingViolationReportInvestigation.IsDisciplinaryRegulatoryCompliance = violationReportInvestigationFormDTO.IsDisciplinaryRegulatoryCompliance
	existingViolationReportInvestigation.ActualRecommendation = violationReportInvestigationFormDTO.ActualRecommendation
	existingViolationReportInvestigation.UpdatedBYID = &currentUserID

	err = vris.repo.Update(existingViolationReportInvestigation, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Collect existing department IDs
	existingViolationReportInvestigationDetails, err := vris.repo.GetFullDetailsByID(violationReportID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingDepartmentIDs := make(map[uint]bool)
	for _, dept := range existingViolationReportInvestigationDetails.Departments {
		existingDepartmentIDs[dept.DepartmentID] = true
	}
	// Collect new department IDs
	newDepartmentIDs := make(map[uint]bool)
	for _, dept := range violationReportInvestigationFormDTO.DepartmentIds {
		newDepartmentIDs[*dept] = true
	}
	// Add new departments
	var departments []*models.ViolationReportInvestigationDepartment
	for _, newDept := range violationReportInvestigationFormDTO.DepartmentIds {
		if _, exists := existingDepartmentIDs[*newDept]; !exists {
			// Create a new violation report investigation department record
			department := &models.ViolationReportInvestigationDepartment{
				ViolationReportInvestigationID: existingViolationReportInvestigationDetails.ID,
				DepartmentID:                   *newDept,
				CreatedBYID:                    currentUserID,
			}
			departments = append(departments, department)
		}
	}
	// Save the new departments
	if len(departments) > 0 {
		if err := vris.repo.CreateDepartments(departments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove departments not in the new list
	var deleteDepartmentIDs []*uint
	for _, existingDept := range existingViolationReportInvestigationDetails.Departments {
		if _, exists := newDepartmentIDs[existingDept.DepartmentID]; !exists {
			// Delete the existing department from the repository
			deleteDepartmentIDs = append(deleteDepartmentIDs, &existingDept.ID)
		}
	}
	// Delete the departments
	if len(deleteDepartmentIDs) > 0 {
		if err := vris.repo.DeleteDepartments(deleteDepartmentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update or delete documents
	var documents []*models.ViolationReportInvestigationDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.ViolationReportInvestigationDocument

	for _, doc := range violationReportInvestigationFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "investigations", strconv.FormatUint(uint64(existingViolationReportInvestigation.ID), 10), "documents")

			document := &models.ViolationReportInvestigationDocument{
				ViolationReportInvestigationID: existingViolationReportInvestigation.ID,
				URL:                            url,
				Name:                           doc.FileName,
				Size:                           doc.Size,
				Token:                          encryptedFilename,
				Type:                           doc.Type,
				OriginalName:                   doc.OriginalName,
				CreatedByID:                    currentUserID,
			}
			documents = append(documents, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := vris.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					return err
				}
				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
				documentsToDelete = append(documentsToDelete, existingDocument)
			}
		}
	}

	// Save Document
	if len(documents) > 0 {
		if err := vris.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIDs) > 0 {
		if err := vris.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Store File
	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "investigations", strconv.FormatUint(uint64(existingViolationReportInvestigation.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Remove File
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "investigations", strconv.FormatUint(uint64(existingViolationReportInvestigation.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	go func() {
		existingViolationReportInvestigation, err := vris.GetFullDetailsByID(violationReportID)
		if err != nil {
			return
		}

		departmentsAdded := []uint{}
		departmentsRemoved := []uint{}
		for id := range newDepartmentIDs {
			if !existingDepartmentIDs[id] {
				departmentsAdded = append(departmentsAdded, id)
			}
		}
		for id := range existingDepartmentIDs {
			if !newDepartmentIDs[id] {
				departmentsRemoved = append(departmentsRemoved, id)
			}
		}

		// log changes
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(departmentsAdded) > 0 || len(departmentsRemoved) > 0 {
			associatedItems["department_ids"] = map[string]interface{}{
				"added":   departmentsAdded,
				"removed": departmentsRemoved,
			}
		}

		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		oldData := map[string]interface{}{
			"investigation_date":                    oldViolationReportInvestigation.InvestigationDateString,
			"steps_taken":                           oldViolationReportInvestigation.StepsTaken,
			"suggested_recommendation":              oldViolationReportInvestigation.SuggestedRecommendation,
			"finding_summary":                       oldViolationReportInvestigation.FindingSummary,
			"disciplinary_procedure":                oldViolationReportInvestigation.DisciplinaryProcedure,
			"is_disciplinary_regulatory_compliance": oldViolationReportInvestigation.IsDisciplinaryRegulatoryCompliance,
			"actual_recommendation":                 oldViolationReportInvestigation.ActualRecommendation,
			"note_and_observations":                 oldViolationReportInvestigation.NoteAndObservations,
		}

		newData := map[string]interface{}{
			"investigation_date":                    existingViolationReportInvestigation.InvestigationDateString,
			"steps_taken":                           existingViolationReportInvestigation.StepsTaken,
			"suggested_recommendation":              existingViolationReportInvestigation.SuggestedRecommendation,
			"finding_summary":                       existingViolationReportInvestigation.FindingSummary,
			"disciplinary_procedure":                existingViolationReportInvestigation.DisciplinaryProcedure,
			"is_disciplinary_regulatory_compliance": existingViolationReportInvestigation.IsDisciplinaryRegulatoryCompliance,
			"actual_recommendation":                 existingViolationReportInvestigation.ActualRecommendation,
			"note_and_observations":                 existingViolationReportInvestigation.NoteAndObservations,
		}

		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := vris.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ViolationReportInvestigation{}.TableName(),
				existingViolationReportInvestigation.ID,
				changedData,
				currentUserID,
				"violation_report", "violation_report_register", "investigation", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ViolationReportInvestigation{}.TableName(),
					"recordID":  existingViolationReportInvestigation.ID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (vris *ViolationReportInvestigationService) Delete(violationReportID uint, currentUserID uint) (uint, error) {
	tx := vris.repo.BeginTransaction()
	var failedId uint
	var deletedViolationReportInvestigation []map[string]interface{}
	documentsToDelete := make(map[uint][]*models.ViolationReportInvestigationDocument)

	if _, err := vris.violationReportSvc.GetByID(violationReportID); err != nil {
		return 0, err
	}
	violationReportInvestigation, err := vris.repo.GetFullDetailsByID(violationReportID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(violationReportInvestigation.Departments) > 0 {
		var deleteDepartmentIds []*uint
		for _, dept := range violationReportInvestigation.Departments {
			deleteDepartmentIds = append(deleteDepartmentIds, &dept.ID)
		}
		if err := vris.repo.DeleteDepartments(deleteDepartmentIds, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Delete each violation report investigation document if documents exist
	if len(violationReportInvestigation.Documents) > 0 {
		var deleteDocumentIds []*uint
		for _, doc := range violationReportInvestigation.Documents {
			documentsToDelete[violationReportInvestigation.ID] = append(documentsToDelete[violationReportInvestigation.ID], doc)
			deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
		}
		if err := vris.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err = vris.repo.Delete(uint(violationReportInvestigation.ID), tx); err != nil {
		failedId = uint(violationReportInvestigation.ID)
		return failedId, err
	}
	if len(documentsToDelete) > 0 {
		for violationReportInvestigationID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "investigations", strconv.FormatUint(uint64(violationReportInvestigationID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}
	}
	deletedViolationReportInvestigation = append(deletedViolationReportInvestigation, map[string]interface{}{
		"id":                  violationReportInvestigation.ID,
		"violation_report_id": violationReportID,
	})
	tx.Commit()
	// Log successful delete activity
	go func() {
		err := vris.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.ViolationReportInvestigation{}.TableName(),
			uint(violationReportInvestigation.ID),
			deletedViolationReportInvestigation,
			currentUserID,
			"violation_report", "violation_report_register", "investigation", nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}

func (vris *ViolationReportInvestigationService) validateViolationReportInvestigationFormDTO(violationReportInvestigationFormDTO *violationReportDTOs.ViolationReportInvestigationFormDTO) error {
	// Validate ActualRecommendation when IsDisciplinaryRegulatoryCompliance is false
	if violationReportInvestigationFormDTO.IsDisciplinaryRegulatoryCompliance != nil &&
		!*violationReportInvestigationFormDTO.IsDisciplinaryRegulatoryCompliance {
		if violationReportInvestigationFormDTO.ActualRecommendation == nil ||
			*violationReportInvestigationFormDTO.ActualRecommendation == "" {
			return utility.NewValidationError(
				fmt.Errorf("actual recommendation is required when disciplinary regulatory compliance is false"),
				"required",
				"actual_recommendation")
		}
	}

	// Department exist
	departments := violationReportInvestigationFormDTO.DepartmentIds
	// If no departments are provided, simply return nil to exit the function
	if len(departments) == 0 {
		return nil
	}

	// Create a map to store department IDs
	departmentIDs := make(map[uint]bool)

	// Check if each department exists
	for _, dept := range departments {
		// Check if the department ID already exists in the map
		if departmentIDs[*dept] {
			// If the department ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate department ID: %d", dept), "duplicate", "department")
		}

		// Note: Add department validation here if you have an organization service
		// if _, err := vris.organizationSvc.GetOrganizationUnitByID(*dept); err != nil {
		// 	return utility.NewValidationError(fmt.Errorf("department does not exist for ID: %d", dept), "not_exist", "department")
		// }

		// Add the department ID to the map
		departmentIDs[*dept] = true
	}

	return nil
}

func (vris *ViolationReportInvestigationService) parseViolationReportInvestigationDates(dto *violationReportDTOs.ViolationReportInvestigationFormDTO) (*violationReportDTOs.ViolationReportInvestigationParsedDates, error) {
	if dto.InvestigationDate == "" {
		return nil, utility.NewValidationError(fmt.Errorf("investigation date is required"), "required", "investigation date")
	}

	parsed := &violationReportDTOs.ViolationReportInvestigationParsedDates{}
	parsedDate, err := utility.ParseDate(dto.InvestigationDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid investigation date format"), "invalid_format", "investigation date")
	}
	if parsedDate == nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid investigation date"), "invalid_date", "investigation date")
	}

	parsed.InvestigationDate = *parsedDate
	return parsed, nil
}
func (vris *ViolationReportInvestigationService) GetDocumentPath(violationReportID uint, token string) (string, error) {
	existingViolationReportInvestigation, err := vris.GetFullDetailsByID(violationReportID)
	if err != nil {
		return "", err
	}

	if existingViolationReportInvestigation == nil {
		return "", errors.New("existing violation report investigation not found")
	}

	if len(existingViolationReportInvestigation.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the violation report investigation")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "investigations", strconv.FormatUint(uint64(existingViolationReportInvestigation.ID), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (vris *ViolationReportInvestigationService) UpdateInvestigationNoteAndObservations(id uint, currentUserID uint, dto *violationReportDTOs.ViolationReportInvestigationNoteAndObservationDTO) error {
	tx := vris.repo.BeginTransaction()

	// Get existing investigation for audit log
	existing, err := vris.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the note
	if err := vris.repo.UpdateNoteAndObservations(id, dto.NoteAndObservations, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the activity
	go func() {
		changedData := map[string]interface{}{
			"note_and_observations": map[string]interface{}{
				"old": existing.NoteAndObservations,
				"new": dto.NoteAndObservations,
			},
		}

		_ = vris.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ViolationReportInvestigation{}.TableName(),
			id,
			changedData,
			currentUserID,
			"violation_report", "violation_report_register", "investigation", nil,
		)
	}()

	return nil
}
