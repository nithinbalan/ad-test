package services

import (
	"fmt"
	"jethur/config"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	complianceModels "jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	activityLogModels "jethur/grc/logs/models"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
)

type SLARegisterServiceInterface interface {
	GetAll(params *slaDTOs.SLARegisterListQueryParamsDTO) (*slaDTOs.SLARegisterListResponseDTO, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(currentUserID uint, slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, sla *slaDTOs.SLARegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type SLARegisterService struct {
	repo               repositories.SLARegisterRepositoryInterface
	slaCategorySvc     SLACategoryServiceInterface
	severityLevelSvc   masterServices.SeverityLevelServiceInterface
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
}

func NewSLARegisterService(
	repo repositories.SLARegisterRepositoryInterface,
	slaCategorySvc SLACategoryServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface) SLARegisterServiceInterface {
	return &SLARegisterService{
		repo:               repo,
		slaCategorySvc:     slaCategorySvc,
		severityLevelSvc:   severityLevelSvc,
		lifecycleStatusSvc: lifecycleStatusSvc,
		activityLogSvc:     activityLogSvc,
	}
}

func (srs *SLARegisterService) GetAll(params *slaDTOs.SLARegisterListQueryParamsDTO) (*slaDTOs.SLARegisterListResponseDTO, error) {
	return srs.repo.GetAll(params)
}

func (srs SLARegisterService) GetByID(id uint) (*models.SLARegister, error) {
	return srs.repo.GetByID(id)
}

func (srs *SLARegisterService) Create(currentUserID uint, slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error) {
	if err := srs.validateSLARegisterFormDTO(slaRegisterFormDTO); err != nil {
		return 0, err
	}

	referenceNumber, err := srs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByStatus("active")
	if err != nil {
		return 0, err
	}

	// Create SLA register
	slaRegister := &models.SLARegister{
		SLACategoryID:     slaRegisterFormDTO.SLACategoryID,
		SeverityLevelID:   slaRegisterFormDTO.SeverityLevelID,
		ResponsibleUnitID: slaRegisterFormDTO.ResponsibleUnitID,
		ResponsibleUserID: slaRegisterFormDTO.ResponsibleUserID,
		LifecycleStatusID: lifecycleStatus.ID,
		ReferenceNumber:   referenceNumber,
		TitleEn:           slaRegisterFormDTO.TitleEn,
		TitleAr:           slaRegisterFormDTO.TitleAr,
		Description:       slaRegisterFormDTO.Description,
		Client:            slaRegisterFormDTO.Client,
		CreatedBy:         currentUserID,
	}

	id, err := srs.repo.Create(slaRegister)
	if err != nil {
		return 0, err
	}

	//file move to storage folder.

	// srs.repo.SaveSLADocumentFile(id,slaRegisterFormDTO.Document)

	uploadedDocument, err := srs.UploadDocument(currentUserID, id, slaRegisterFormDTO.Document)
	if err != nil {
		return 0, err
	}

	if err := srs.repo.SaveSLADocument(uploadedDocument); err != nil {
		return 0, err
	}
	// Log activity
	go func() {
		err := srs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			complianceModels.SLARegister{}.TableName(),
			id,
			slaRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": complianceModels.SLARegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (srs *SLARegisterService) Update(id uint, currentUserID uint, slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) error {
	existingSLARegister, err := srs.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := srs.validateSLARegisterFormDTO(slaRegisterFormDTO); err != nil {
		return err
	}

	existingSLARegister.SLACategoryID = slaRegisterFormDTO.SLACategoryID
	existingSLARegister.ResponsibleUnitID = slaRegisterFormDTO.ResponsibleUnitID
	existingSLARegister.ResponsibleUserID = slaRegisterFormDTO.ResponsibleUserID
	existingSLARegister.SeverityLevelID = slaRegisterFormDTO.SeverityLevelID
	existingSLARegister.LifecycleStatusID = slaRegisterFormDTO.LifecycleStatusID
	existingSLARegister.TitleEn = slaRegisterFormDTO.TitleEn
	existingSLARegister.TitleAr = slaRegisterFormDTO.TitleAr
	existingSLARegister.Description = slaRegisterFormDTO.Description
	existingSLARegister.Client = slaRegisterFormDTO.Client
	existingSLARegister.UpdatedBy = currentUserID

	err = srs.repo.Update(existingSLARegister)
	if err != nil {
		return err
	}

	if slaRegisterFormDTO.Document != nil {
		// Upload or update SLA document
		_, err := srs.UploadDocument(currentUserID, id, slaRegisterFormDTO.Document)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srs *SLARegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := srs.repo.GetByID(id)

		if err != nil {
			continue
		}
		// Delete associated documents
		err = srs.DeleteDocument(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		err = srs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeleteAction,
				models.SLARegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (srs *SLARegisterService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByStatus("archived")
	if err != nil {
		return nil, err
	}

	idsToUpdate, err := srs.repo.GetIDsByLifecycleStatus(ids, "active")
	if err != nil {
		return nil, err
	}

	successfulIds, err := srs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				logModels.ArchiveAction,
				complianceModels.SLARegister{}.TableName(),
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

func (srs *SLARegisterService) Restore(ids []uint, currentUserID uint) ([]uint, error) {

	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByStatus("active")
	if err != nil {
		return nil, err
	}
	idsToUpdate, err := srs.repo.GetIDsByLifecycleStatus(ids, "archived")
	if err != nil {
		return nil, err
	}

	successfulIds, err := srs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				logModels.RestoreAction,
				complianceModels.SLARegister{}.TableName(),
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

func (srs *SLARegisterService) validateSLARegisterFormDTO(slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) error {
	if _, err := srs.slaCategorySvc.GetByID(slaRegisterFormDTO.SLACategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("sla category does not exist"), "not_exist", "sla category")
	}

	if slaRegisterFormDTO.LifecycleStatusID != 0 {
		if _, err := srs.lifecycleStatusSvc.GetByID(slaRegisterFormDTO.LifecycleStatusID); err != nil {
			return utility.NewValidationError(fmt.Errorf("status does not exist"), "not_exist", "lifecycle_status")
		}
	}
	if slaRegisterFormDTO.SeverityLevelID != 0 {
		if _, err := srs.severityLevelSvc.GetByID(slaRegisterFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	return nil
}

func (srs *SLARegisterService) parseSLARegisterDates(document *slaDTOs.SLADocumentFormDTO) (*slaDTOs.ParsedSLADates, error) {
	parsed := &slaDTOs.ParsedSLADates{}
	var err error

	// Validate issue_date
	if document.IssueDate != "" {
		parsed.IssueDate, err = utility.ParseDate(document.IssueDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
		}
	}

	// Validate expiry_date
	if document.ExpiryDate != "" {
		parsed.ExpiryDate, err = utility.ParseDate(document.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	// Validate SA1Date
	if document.SA1Date != "" {
		parsed.SA1Date, err = utility.ParseDate(document.SA1Date)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid sa1 date"), "invalid_date", "sa1 date")
		}
	}

	// Validate SA2Date
	if document.SA2Date != "" {
		parsed.SA2Date, err = utility.ParseDate(document.SA2Date)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid sa2 date"), "invalid_date", "sa2 date")
		}
	}

	return parsed, nil
}

func (srs *SLARegisterService) UploadDocument(currentUserID uint, id uint, document *slaDTOs.SLADocumentFormDTO) (*models.SLADocument, error) {
	storagePath := config.GetEnv("STORAGE_PATH", "")

	// Construct the directory path for the SLA register
	directoryPath := filepath.Join(storagePath, "compliance", "sla", strconv.FormatUint(uint64(id), 10), "document")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		return nil, err
	}

	// Read the file from temporary storage
	fileBytes, err := os.ReadFile(filepath.Join(storagePath, "tmp", document.Name))
	if err != nil {
		return nil, err
	}

	// Determine the destination path in permanent storage
	permanentStoragePath := filepath.Join(directoryPath, document.Name)

	// Write the file to permanent storage
	if err := os.WriteFile(permanentStoragePath, fileBytes, 0644); err != nil {
		return nil, err
	}

	// Parse SLA register dates
	dates, err := srs.parseSLARegisterDates(document)
	if err != nil {
		return nil, err
	}
	// Create or update SLA document based on whether it already exists
	var slaDocument *models.SLADocument
	if document.ID != nil && *document.ID != 0 {
		// Fetch existing SLA document
		existingSLADocument, err := srs.repo.GetSLADocumentByID(*document.ID)
		if err != nil {
			return nil, err
		}

		// Update SLA document details
		existingSLADocument.Version = document.Version
		existingSLADocument.IssueDate = dates.IssueDate
		existingSLADocument.ExpiryDate = dates.ExpiryDate
		existingSLADocument.SA1Date = dates.SA1Date
		existingSLADocument.SA2Date = dates.SA2Date
		existingSLADocument.Comment = document.Comment
		existingSLADocument.Language = document.Language
		existingSLADocument.URL = document.URL
		existingSLADocument.Name = document.Name
		existingSLADocument.Size = document.Size
		existingSLADocument.Token = document.Token
		existingSLADocument.Type = document.Type
		existingSLADocument.Icon = document.Icon
		existingSLADocument.UpdatedBy = currentUserID

		// Update SLA document in the repository
		err = srs.repo.SaveSLADocument(existingSLADocument)
		if err != nil {
			return nil, err
		}

		slaDocument = existingSLADocument
	} else {
		// Create new SLA document
		slaDocument = &models.SLADocument{
			SLARegisterID: id,
			Version:       document.Version,
			IssueDate:     dates.IssueDate,
			ExpiryDate:    dates.ExpiryDate,
			SA1Date:       dates.SA1Date,
			SA2Date:       dates.SA2Date,
			Comment:       document.Comment,
			Language:      document.Language,
			URL:           document.URL,
			Name:          document.Name,
			Size:          document.Size,
			Token:         document.Token,
			Type:          document.Type,
			Icon:          document.Icon,
			IsLatest:      true,
			CreatedBy:     currentUserID,
		}
	}

	return slaDocument, nil
}

func (srs *SLARegisterService) DeleteDocument(id uint) error {
	// Fetch the SLA document by ID
	slaDocument, err := srs.repo.GetLatestSLADocumentBySLARegisterID(id)
	if err != nil {
		return err
	}

	// Determine the path of the file to delete
	storagePath := config.GetEnv("STORAGE_PATH", "")
	filePath := filepath.Join(storagePath, "compliance", "sla", strconv.FormatUint(uint64(slaDocument.SLARegisterID), 10), "document", slaDocument.Name)

	// Delete the file from storage
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
