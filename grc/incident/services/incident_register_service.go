package services

import (
	"errors"
	"fmt"
	complianceServices "jethur/grc/compliance/services"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentRegisterServiceInterface interface {
	GetAll(params *incidentDTOs.IncidentRegisterListQueryParamsDTO) (*incidentDTOs.IncidentRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentRegisterDetailDTO, error)
	GetByID(id uint) (*models.IncidentRegister, error)
	Create(currentUserID uint, incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, incident *incidentDTOs.IncidentRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type IncidentRegisterService struct {
	repo                     incidentRepositories.IncidentRegisterRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	incidentCategorySvc      IncidentCategoryServiceInterface
	incidentSubCategorySvc   IncidentSubCategoryServiceInterface
	incidentTypeSvc          IncidentTypeServiceInterface
	incidentStatusSvc        IncidentStatusServiceInterface
	severityLevelSvc         masterServices.SeverityLevelServiceInterface
	controlSvc               organizationServices.ControlServiceInterface
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface
}

func NewIncidentRegisterService(
	repo incidentRepositories.IncidentRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentCategorySvc IncidentCategoryServiceInterface,
	incidentSubCategorySvc IncidentSubCategoryServiceInterface,
	incidentTypeSvc IncidentTypeServiceInterface,
	incidentStatusSvc IncidentStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	controlSvc organizationServices.ControlServiceInterface,
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface,
) IncidentRegisterServiceInterface {
	return &IncidentRegisterService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		incidentCategorySvc:      incidentCategorySvc,
		incidentSubCategorySvc:   incidentSubCategorySvc,
		incidentTypeSvc:          incidentTypeSvc,
		incidentStatusSvc:        incidentStatusSvc,
		severityLevelSvc:         severityLevelSvc,
		controlSvc:               controlSvc,
		complianceRequirementSvc: complianceRequirementSvc,
	}
}

func (irs *IncidentRegisterService) GetAll(params *incidentDTOs.IncidentRegisterListQueryParamsDTO) (*incidentDTOs.IncidentRegisterListResponseDTO, error) {
	return irs.repo.GetAll(params)
}

func (irs IncidentRegisterService) GetByID(id uint) (*models.IncidentRegister, error) {
	incident, err := irs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incident, nil
}

func (irs IncidentRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentRegisterDetailDTO, error) {
	// Initialize your DTO
	var incidentDTO incidentDTOs.IncidentRegisterDetailDTO

	incident, err := irs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(incident, &incidentDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &incidentDTO, nil
}

func (irs *IncidentRegisterService) Create(currentUserID uint, incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) (uint, error) {
	tx := irs.repo.BeginTransaction()
	if err := irs.validateIncidentRegisterFormDTO(incidentRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := irs.parseIncidentRegisterDates(incidentRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := irs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	if err := irs.validateIncidentDocumentFormDTO(incidentRegisterFormDTO.Documents); err != nil {
		return 0, err
	}

	incidentRegister := &models.IncidentRegister{
		ReferenceNumber:       referenceNumber,
		IncidentCategoryID:    incidentRegisterFormDTO.IncidentCategoryID,
		IncidentSubCategoryID: incidentRegisterFormDTO.IncidentSubCategoryID,
		IncidentTypeID:        incidentRegisterFormDTO.IncidentTypeID,
		SeverityLevelID:       incidentRegisterFormDTO.SeverityLevelID,
		IncidentStatusID:      incidentRegisterFormDTO.IncidentStatusID,
		Title:                 incidentRegisterFormDTO.Title,
		Description:           incidentRegisterFormDTO.Description,
		IncidentAt:            dates.IncidentAt,
		Location:              incidentRegisterFormDTO.Location,
		ReportedBy:            incidentRegisterFormDTO.ReportedBy,
		ReportedAt:            dates.ReportedAt,
		QuickAction:           incidentRegisterFormDTO.QuickAction,
		CreatedBy:             currentUserID,
	}

	id, err := irs.repo.Create(incidentRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Save ComplianceRequirements
	var complianceRequirements []*models.IncidentComplianceRequirement
	for _, req := range incidentRegisterFormDTO.ComplianceRequirements {
		requirement := &models.IncidentComplianceRequirement{
			IncidentRegisterID:      id,
			ComplianceRequirementID: *req,
			CreatedBy:               currentUserID,
		}
		complianceRequirements = append(complianceRequirements, requirement)
	}

	if len(complianceRequirements) > 0 {
		if err := irs.repo.SaveIncidentComplianceRequirements(complianceRequirements, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Controls
	var controls []*models.IncidentControl
	for _, ctrl := range incidentRegisterFormDTO.Controls {
		control := &models.IncidentControl{
			IncidentRegisterID: id,
			ControlID:          *ctrl,
			CreatedBy:          currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := irs.repo.SaveIncidentControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	var documents []*models.IncidentDocument

	// Iterate over each document in the array
	for _, doc := range incidentRegisterFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("incident-register", strconv.FormatUint(uint64(incidentRegister.ID), 10), encryptedFilename)

		// Create a new IncidentDocument instance
		document := &models.IncidentDocument{
			IncidentRegisterID: id,
			Language:           doc.Language,
			URL:                url,
			Icon:               doc.Icon,
			Name:               doc.FileName,
			Size:               doc.Size,
			Token:              encryptedFilename,
			Type:               doc.Type,
			CreatedBy:          currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		// Save the IncidentDocument
		if err := irs.repo.SaveIncidentDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentRegister.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := irs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentRegister{}.TableName(),
			id,
			incidentRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (irs *IncidentRegisterService) Update(id uint, currentUserID uint, incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) error {
	tx := irs.repo.BeginTransaction()

	existingIncidentRegister, err := irs.repo.GetByID(id)

	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := irs.validateIncidentRegisterFormDTO(incidentRegisterFormDTO); err != nil {
		return err
	}

	dates, err := irs.parseIncidentRegisterDates(incidentRegisterFormDTO)
	if err != nil {
		return err
	}
	if err := irs.validateIncidentDocumentFormDTO(incidentRegisterFormDTO.Documents); err != nil {
		return err
	}
	existingIncidentRegister.IncidentCategoryID = incidentRegisterFormDTO.IncidentCategoryID
	existingIncidentRegister.IncidentSubCategoryID = incidentRegisterFormDTO.IncidentSubCategoryID
	existingIncidentRegister.IncidentTypeID = incidentRegisterFormDTO.IncidentTypeID
	existingIncidentRegister.SeverityLevelID = incidentRegisterFormDTO.SeverityLevelID
	existingIncidentRegister.IncidentStatusID = incidentRegisterFormDTO.IncidentStatusID
	existingIncidentRegister.Title = incidentRegisterFormDTO.Title
	existingIncidentRegister.Description = incidentRegisterFormDTO.Description
	existingIncidentRegister.IncidentAt = dates.IncidentAt
	existingIncidentRegister.Location = incidentRegisterFormDTO.Location
	existingIncidentRegister.ReportedBy = incidentRegisterFormDTO.ReportedBy
	existingIncidentRegister.ReportedAt = dates.ReportedAt
	existingIncidentRegister.QuickAction = incidentRegisterFormDTO.QuickAction
	existingIncidentRegister.UpdatedBy = &currentUserID

	err = irs.repo.Update(existingIncidentRegister, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	//documents update
	var documents []*models.IncidentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.IncidentDocument

	// Iterate over each document in the array
	for _, doc := range incidentRegisterFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("incident-register", strconv.FormatUint(uint64(id), 10), encryptedFilename)
			// Create a new IncidentDocument instance
			document := &models.IncidentDocument{
				IncidentRegisterID: id,
				Language:           doc.Language,
				URL:                url,
				Icon:               doc.Icon,
				Name:               doc.FileName,
				Size:               doc.Size,
				Token:              encryptedFilename,
				Type:               doc.Type,
				CreatedBy:          currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := irs.repo.GetIncidentDocumentByID(*doc.ID)
				if err != nil {
					return err
				}
				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
				documentsToDelete = append(documentsToDelete, existingDocument)
			}
		}
	}
	if len(documents) > 0 {
		if err := irs.repo.SaveIncidentDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := irs.repo.DeleteIncidentDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(existingIncidentRegister.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(existingIncidentRegister.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()

	go func() {
		err := irs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentRegister{}.TableName(),
			id,
			incidentRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.IncidentRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (irs *IncidentRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := irs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*incidentDTOs.IncidentDocumentDTO)

	for _, id := range ids {
		incident, err := irs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Check if controls exist
		if len(incident.Controls) > 0 {
			var deleteControlIds []*uint
			for _, ctrl := range incident.Controls {
				deleteControlIds = append(deleteControlIds, ctrl.ID)
			}
			if err := irs.repo.DeleteIncidentControls(deleteControlIds, tx); err != nil {

				tx.Rollback()
				return nil, err
			}
		}

		// Delete each compliance requirement if requirements exist
		if len(incident.ComplianceRequirements) > 0 {
			var deleteRequirementIds []*uint
			for _, req := range incident.ComplianceRequirements {
				deleteRequirementIds = append(deleteRequirementIds, req.ID)
			}
			if err := irs.repo.DeleteIncidentComplianceRequirements(deleteRequirementIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		// Delete each incident document if documents exist
		if len(incident.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range incident.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := irs.repo.DeleteIncidentDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		// Delete the incident
		if err := irs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for incidentID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}
			}
		}
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := irs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.IncidentRegister{}.TableName(),
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

func (irs *IncidentRegisterService) validateIncidentRegisterFormDTO(incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) error {
	if _, err := irs.incidentCategorySvc.GetByID(incidentRegisterFormDTO.IncidentCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident category does not exist"), "not_exist", "incident category")
	}

	if _, err := irs.incidentSubCategorySvc.GetByID(incidentRegisterFormDTO.IncidentSubCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident sub category does not exist"), "not_exist", "incident sub category")
	}

	if _, err := irs.incidentTypeSvc.GetByID(incidentRegisterFormDTO.IncidentTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident type does not exist"), "not_exist", "incident type")
	}

	if _, err := irs.incidentStatusSvc.GetByID(incidentRegisterFormDTO.IncidentStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident status does not exist"), "not_exist", "incident status")
	}

	if _, err := irs.severityLevelSvc.GetByID(incidentRegisterFormDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}
	//Controls exist
	controls := incidentRegisterFormDTO.Controls
	// If no controls are provided, simply return nil to exit the function
	if len(controls) == 0 {
		return nil
	}

	// Create a map to store control IDs
	controlIDs := make(map[uint]bool)

	// Check if each control exists
	for _, control := range controls {
		// Check if the control ID already exists in the map
		if controlIDs[*control] {
			// If the control ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", control), "duplicate", "control")
		}

		// Check if the control exists in the database
		if _, err := irs.controlSvc.GetByID(*control); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", control), "not_exist", "control")
		}

		// Add the control ID to the map
		controlIDs[*control] = true
	}
	//Compliance Requirements Exist
	requirements := incidentRegisterFormDTO.ComplianceRequirements

	// If no controls are provided, simply return nil to exit the function
	if len(requirements) == 0 {
		return nil
	}

	// Create a map to store requirement IDs
	requirementIDs := make(map[uint]bool)

	// Check if each requirement exists
	for _, requirement := range requirements {
		// Check if the requirement ID already exists in the map
		if requirementIDs[*requirement] {
			// If the requirement ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate compliance requirement ID: %d", requirement), "duplicate", "compliance requirement")
		}

		// Check if the requirement exists in the database
		if _, err := irs.complianceRequirementSvc.GetByComplianceRequirementID(*requirement); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement does not exist for ID: %d", requirement), "not_exist", "compliance requirement")
		}

		// Add the requirement ID to the map
		requirementIDs[*requirement] = true
	}

	return nil
}

func (irs *IncidentRegisterService) parseIncidentRegisterDates(dto *incidentDTOs.IncidentRegisterFormDTO) (*incidentDTOs.IncidentRegisterParsedDates, error) {
	parsed := &incidentDTOs.IncidentRegisterParsedDates{}

	var err error
	if dto.ReportedAt != "" {
		parsed.ReportedAt, err = utility.ParseDate(dto.ReportedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported date"), "invalid_date", "reported at")
		}
	}

	if dto.IncidentAt != "" {
		parsed.IncidentAt, err = utility.ParseDate(dto.IncidentAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid incident date"), "invalid_date", "incident at")
		}
	}

	return parsed, nil
}

func (irs *IncidentRegisterService) validateIncidentDocumentFormDTO(documents []*incidentDTOs.IncidentDocumentFormDTO) error {
	// Check if documents array is empty
	if len(documents) == 0 {
		return utility.NewValidationError(fmt.Errorf("documents array is required"), "required", "documents")
	}

	for _, document := range documents {
		//@Review
		//All items are mandatory during creation

		// Validate icon
		if document.Icon == "" {
			return utility.NewValidationError(fmt.Errorf("icon is required"), "required", "icon")
		}

		// @todo
		// Validate file size (if needed)
		// For example:
		// if document.Size > MaxFileSize {
		//     return utility.NewValidationError(fmt.Errorf("file size exceeds the limit"), "invalid_file_size", "file_size")
		// }
	}

	return nil
}
