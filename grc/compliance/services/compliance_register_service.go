package services

import (
	"errors"
	"fmt"
	"jethur/config"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceRegisterServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRegisterDetailDTO, error)
	GetByID(id uint) (*models.ComplianceRegister, error)
	Create(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error
	ReviewStatusUpdate(complianceRegister *models.ComplianceRegister, tx utility.DatabaseTransaction) error

	GetDocumentPath(id uint, token string) (string, error)
}

// Service and implementations
type ComplianceRegisterService struct {
	repo                  complianceRepositories.ComplianceRegisterRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	lifecycleStatusSvc    masterServices.LifecycleStatusServiceInterface
	complianceCategorySvc ComplianceCategoryServiceInterface
	complianceStatusSvc   ComplianceStatusServiceInterface
	reviewFrequencySvc    masterServices.ReviewFrequencyServiceInterface
	severityLevelSvc      masterServices.SeverityLevelServiceInterface
	controlSvc            organizationServices.ControlServiceInterface
}

func NewComplianceRegisterService(
	repo complianceRepositories.ComplianceRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	complianceCategorySvc ComplianceCategoryServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	controlSvc organizationServices.ControlServiceInterface) ComplianceRegisterServiceInterface {
	return &ComplianceRegisterService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		lifecycleStatusSvc:    lifecycleStatusSvc,
		complianceCategorySvc: complianceCategorySvc,
		complianceStatusSvc:   complianceStatusSvc,
		reviewFrequencySvc:    reviewFrequencySvc,
		severityLevelSvc:      severityLevelSvc,
		controlSvc:            controlSvc,
	}
}

func (crs *ComplianceRegisterService) GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error) {
	return crs.repo.GetAll(params)
}

func (crs ComplianceRegisterService) GetByID(id uint) (*models.ComplianceRegister, error) {
	compliance, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return compliance, nil
}

func (crs ComplianceRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRegisterDetailDTO, error) {
	var complianceDTO complianceDTOs.ComplianceRegisterDetailDTO

	compliance, err := crs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(compliance, &complianceDTO)
	if err != nil {
		return nil, err
	}

	// Format time fields without time component
	if compliance.Deadline != nil {
		deadline := compliance.Deadline.Format("2006-01-02")
		complianceDTO.DeadlineString = &deadline
	}
	if compliance.IssueDate != nil {
		issueDate := compliance.IssueDate.Format("2006-01-02")
		complianceDTO.IssueDateString = &issueDate
	}
	if compliance.ExpiryDate != nil {
		expiryDate := compliance.ExpiryDate.Format("2006-01-02")
		complianceDTO.ExpiryDateString = &expiryDate
	}
	if compliance.NextReviewDate != nil {
		nextReviewDate := compliance.NextReviewDate.Format("2006-01-02")
		complianceDTO.NextReviewDateString = &nextReviewDate
	}
	if compliance.LastReviewedDate != nil {
		lastReviewDate := compliance.LastReviewedDate.Format("2006-01-02")
		complianceDTO.LastReviewedDateString = &lastReviewDate
	}

	complianceDTO.CreatedAt = compliance.CreatedAt
	complianceDTO.UpdatedAt = compliance.UpdatedAt
	//todo :- also add extra fields which required in details screen

	return &complianceDTO, nil
}

func (crs *ComplianceRegisterService) Create(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	if err := crs.validateComplianceRegisterFormDTO(complianceRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRegisterDates(complianceRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := crs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	lifecycleStatus, err := crs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}

	if err := crs.validateComplianceDocumentFormDTO(&complianceRegisterFormDTO.Document); err != nil {
		return 0, err
	}

	complianceRegister := &models.ComplianceRegister{
		ComplianceCategoryID:     complianceRegisterFormDTO.ComplianceCategoryID,
		ComplianceStatusID:       complianceRegisterFormDTO.ComplianceStatusID,
		ResponsibleUnitID:        complianceRegisterFormDTO.ResponsibleUnitID,
		ResponsibleUserID:        complianceRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID:        complianceRegisterFormDTO.ReviewFrequencyID,
		LastReviewNote:           complianceRegisterFormDTO.LastReviewNote,
		ReviewResponsibleUserID:  complianceRegisterFormDTO.ReviewResponsibleUserID,
		SeverityLevelID:          complianceRegisterFormDTO.SeverityLevelID,
		ReferenceNumber:          referenceNumber,
		StandardCode:             complianceRegisterFormDTO.StandardCode,
		TitleEn:                  complianceRegisterFormDTO.TitleEn,
		TitleAr:                  complianceRegisterFormDTO.TitleAr,
		DescriptionEn:            complianceRegisterFormDTO.DescriptionEn,
		DescriptionAr:            complianceRegisterFormDTO.DescriptionAr,
		Version:                  complianceRegisterFormDTO.Version,
		Deadline:                 dates.Deadline,
		IssueDate:                dates.IssueDate,
		ExpiryDate:               dates.ExpiryDate,
		NextReviewDate:           dates.NextReviewDate,
		LastReviewedDate:         dates.LastReviewedDate,
		LastReviewedUserID:       complianceRegisterFormDTO.LastReviewedUserID,
		ComplianceReviewStatusID: complianceRegisterFormDTO.ComplianceReviewStatusID,
		CreatedByID:              currentUserID,
		LifecycleStatusID:        lifecycleStatus.ID,
	}

	id, err := crs.repo.Create(complianceRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save Controls
	var controls []*models.ComplianceControl
	for _, ctrl := range complianceRegisterFormDTO.ControlIDs {
		control := &models.ComplianceControl{
			ComplianceRegisterID:    id,
			ControlID:               *ctrl,
			ComplianceRequirementID: nil,
			CreatedByID:             currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := crs.repo.CreateComplianceControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	//generate token
	dataToEncrypt := []byte(complianceRegisterFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return 0, err
	}

	//genereate URL
	url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")

	complianceDocument := &models.ComplianceDocument{
		ComplianceRegisterID: id,
		Language:             complianceRegisterFormDTO.Document.Language,
		URL:                  url,
		Name:                 complianceRegisterFormDTO.Document.FileName,
		Size:                 complianceRegisterFormDTO.Document.Size,
		Token:                encryptedFilename,
		OriginalName:         complianceRegisterFormDTO.Document.OriginalName,
		Type:                 complianceRegisterFormDTO.Document.Type,
		CreatedByID:          currentUserID,
	}

	// Save Compliance document
	if err := crs.repo.CreateDocument(complianceDocument, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	// Store Compliance Document File using utility.StoreFile
	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")

	if err := utility.StoreFile(complianceRegisterFormDTO.Document.FileName, storagePath); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceRegister{}.TableName(),
			id,
			complianceRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crs *ComplianceRegisterService) Update(id uint, currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {
	tx := crs.repo.BeginTransaction()

	existingComplianceRegister, err := crs.GetByID(id)

	if err != nil {
		return err
	}

	if err := crs.validateComplianceRegisterFormDTO(complianceRegisterFormDTO); err != nil {
		return err
	}

	dates, err := crs.parseComplianceRegisterDates(complianceRegisterFormDTO)
	if err != nil {
		return err
	}

	if err := crs.validateComplianceDocumentFormDTO(&complianceRegisterFormDTO.Document); err != nil {
		return err
	}

	existingComplianceRegister.ComplianceCategoryID = complianceRegisterFormDTO.ComplianceCategoryID
	existingComplianceRegister.ComplianceStatusID = complianceRegisterFormDTO.ComplianceStatusID
	existingComplianceRegister.ResponsibleUnitID = complianceRegisterFormDTO.ResponsibleUnitID
	existingComplianceRegister.ResponsibleUserID = complianceRegisterFormDTO.ResponsibleUserID
	existingComplianceRegister.ReviewFrequencyID = complianceRegisterFormDTO.ReviewFrequencyID
	existingComplianceRegister.ReviewResponsibleUserID = complianceRegisterFormDTO.ReviewResponsibleUserID
	existingComplianceRegister.SeverityLevelID = complianceRegisterFormDTO.SeverityLevelID
	existingComplianceRegister.StandardCode = complianceRegisterFormDTO.StandardCode
	existingComplianceRegister.TitleEn = complianceRegisterFormDTO.TitleEn
	existingComplianceRegister.TitleAr = complianceRegisterFormDTO.TitleAr
	existingComplianceRegister.DescriptionEn = complianceRegisterFormDTO.DescriptionEn
	existingComplianceRegister.DescriptionAr = complianceRegisterFormDTO.DescriptionAr
	existingComplianceRegister.Version = complianceRegisterFormDTO.Version
	existingComplianceRegister.Deadline = dates.Deadline
	existingComplianceRegister.IssueDate = dates.IssueDate
	existingComplianceRegister.ExpiryDate = dates.ExpiryDate
	existingComplianceRegister.NextReviewDate = dates.NextReviewDate
	existingComplianceRegister.LastReviewedDate = dates.LastReviewedDate
	existingComplianceRegister.LastReviewedUserID = complianceRegisterFormDTO.LastReviewedUserID
	existingComplianceRegister.ComplianceReviewStatusID = complianceRegisterFormDTO.ComplianceReviewStatusID
	existingComplianceRegister.LastReviewNote = complianceRegisterFormDTO.LastReviewNote
	existingComplianceRegister.UpdatedByID = &currentUserID

	err = crs.repo.Update(existingComplianceRegister, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Collect existing control IDs
	existingCompliance, err := crs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range existingCompliance.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	// Collect new control IDs
	newControlIDs := make(map[uint]bool)
	for _, ctrl := range complianceRegisterFormDTO.ControlIDs {
		newControlIDs[*ctrl] = true
	}

	// Add new controls
	var controlsToAdd []*models.ComplianceControl
	for _, newControlID := range complianceRegisterFormDTO.ControlIDs {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			// Create a new compliance control record
			control := &models.ComplianceControl{
				ComplianceRegisterID:    id,
				ComplianceRequirementID: nil,
				ControlID:               *newControlID,
				CreatedByID:             currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	// Save the new controls
	if len(controlsToAdd) > 0 {
		if err := crs.repo.CreateComplianceControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove controls not in the new list
	var deleteControlIDs []*uint

	for _, existingControl := range existingCompliance.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			// Append control to delete list
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	// Delete controls not in the new list
	if len(deleteControlIDs) > 0 {
		if err := crs.repo.DeleteComplianceControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingComplianceDocument, err := crs.repo.GetLatestDocumentByComplianceID(id)
	if err != nil {
		return err
	}

	if complianceRegisterFormDTO.Document.IsNewFile != nil && *complianceRegisterFormDTO.Document.IsNewFile {
		existingComplianceDocumentFileName := existingComplianceDocument.Name

		//generate token
		dataToEncrypt := []byte(complianceRegisterFormDTO.Document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}

		//genereate URL
		url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")

		existingComplianceDocument.Language = complianceRegisterFormDTO.Document.Language
		existingComplianceDocument.URL = url
		existingComplianceDocument.Name = complianceRegisterFormDTO.Document.FileName
		existingComplianceDocument.Size = complianceRegisterFormDTO.Document.Size
		existingComplianceDocument.Token = encryptedFilename
		existingComplianceDocument.OriginalName = complianceRegisterFormDTO.Document.OriginalName
		existingComplianceDocument.Type = complianceRegisterFormDTO.Document.Type
		existingComplianceDocument.UpdatedByID = currentUserID

		if err := crs.repo.CreateDocument(existingComplianceDocument, tx); err != nil {
			tx.Rollback()
			return err
		}

		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(existingComplianceDocumentFileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

		if err := utility.StoreFile(complianceRegisterFormDTO.Document.FileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

	}

	tx.Commit()
	return nil
}

func (crs *ComplianceRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := crs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*complianceDTOs.ComplianceDocumentDTO)

	for _, id := range ids {
		complianceRegister, err := crs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Check if controls exist
		if len(complianceRegister.Controls) > 0 {
			var deleteControlIds []*uint
			for _, ctrl := range complianceRegister.Controls {
				deleteControlIds = append(deleteControlIds, ctrl.ID)
			}
			if err := crs.repo.DeleteComplianceControls(deleteControlIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete each Compliance document if documents exist
		if len(complianceRegister.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range complianceRegister.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], &doc)
				deleteDocumentIds = append(deleteDocumentIds, doc.ID)
			}
			if err := crs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the Compliance Register
		if err := crs.repo.Delete(id, tx); err != nil {
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
		for complianceID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
			destinationPath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				tx.Rollback()
			}
		}
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRegister{}.TableName(),
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

func (crs *ComplianceRegisterService) Archive(ids []uint, currentUserID uint) ([]uint, error) {

	lifecycleStatus, err := crs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return nil, err
	}
	idsToUpdate, err := crs.repo.GetIDsByLifecycleStatus(ids, "active")
	if err != nil {
		return nil, err
	}

	successfulIds, err := crs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.ArchiveAction,
				models.ComplianceRegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (crs *ComplianceRegisterService) Restore(ids []uint, currentUserID uint) ([]uint, error) {

	lifecycleStatus, err := crs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return nil, err
	}
	idsToUpdate, err := crs.repo.GetIDsByLifecycleStatus(ids, "archived")
	if err != nil {
		return nil, err
	}

	successfulIds, err := crs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.RestoreAction,
				models.ComplianceRegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (crs *ComplianceRegisterService) StatusUpdate(id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error {
	existingCompliance, err := crs.GetByID(id)

	if err != nil {
		return err
	}

	existingCompliance.ComplianceStatusID = complianceStatusID
	existingCompliance.UpdatedByID = &currentUserID

	err = crs.repo.Update(existingCompliance, tx)

	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRegisterService) ReviewStatusUpdate(complianceRegister *models.ComplianceRegister, tx utility.DatabaseTransaction) error {

	err := crs.repo.Update(complianceRegister, tx)
	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRegisterService) validateComplianceRegisterFormDTO(complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {
	if _, err := crs.complianceCategorySvc.GetByID(complianceRegisterFormDTO.ComplianceCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
	}

	if _, err := crs.complianceStatusSvc.GetByID(complianceRegisterFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	if complianceRegisterFormDTO.ReviewFrequencyID != nil {
		if _, err := crs.reviewFrequencySvc.GetByID(*complianceRegisterFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}
	if complianceRegisterFormDTO.SeverityLevelID != nil {
		if _, err := crs.severityLevelSvc.GetByID(*complianceRegisterFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	//Controls exist
	controls := complianceRegisterFormDTO.ControlIDs
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
		if _, err := crs.controlSvc.GetByID(*control); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", control), "not_exist", "control")
		}

		// Add the control ID to the map
		controlIDs[*control] = true
	}

	return nil
}

func (crs *ComplianceRegisterService) parseComplianceRegisterDates(dto *complianceDTOs.ComplianceRegisterFormDTO) (*complianceDTOs.ComplianceRegisterParsedDates, error) {
	parsed := &complianceDTOs.ComplianceRegisterParsedDates{}

	var err error

	if dto.Deadline != nil {
		parsed.Deadline, err = utility.ParseDate(*dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	if dto.IssueDate != nil {
		parsed.IssueDate, err = utility.ParseDate(*dto.IssueDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
		}
	}

	if dto.ExpiryDate != nil {
		parsed.ExpiryDate, err = utility.ParseDate(*dto.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	if dto.NextReviewDate != nil {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	if dto.LastReviewedDate != nil {
		parsed.LastReviewedDate, err = utility.ParseDate(*dto.LastReviewedDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last reviewed date"), "invalid_date", "last reviewed date")
		}
	}

	return parsed, nil
}

func (crs *ComplianceRegisterService) validateComplianceDocumentFormDTO(document *complianceDTOs.ComplianceDocumentFormDTO) error {
	// Validate document
	if document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	}
	// Validate file type (PDF)
	if document.Type != ".pdf" {
		return utility.NewValidationError(fmt.Errorf("only PDF files are allowed"), "invalid_file_type", "type")
	}

	// @todo
	// Validate file size (if needed)
	// For example:
	// if document.Size > MaxFileSize {
	//     return utility.NewValidationError(fmt.Errorf("file size exceeds the limit"), "invalid_file_size", "file_size")
	// }

	return nil
}

func (crs *ComplianceRegisterService) GetDocumentPath(id uint, token string) (string, error) {
	existingComplianceRegister, err := crs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	if len(existingComplianceRegister.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the compliance register")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found!")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found!")
	}

	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")

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
