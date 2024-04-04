package services

import (
	"fmt"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	complianceModels "jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	activityLogModels "jethur/grc/logs/models"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"
)

type SLARegisterServiceInterface interface {
	GetAll(params *slaDTOs.SLARegisterListQueryParamsDTO) (*slaDTOs.SLARegisterListResponseDTO, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(currentUserID uint, slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, sla *slaDTOs.SLARegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
	Renew(id uint, currentUserID uint, document *slaDTOs.SLADocumentFormDTO) (*models.SLADocument, error)
	// Revert(ids []uint, currentUserID uint) ([]uint, error)
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

	dates, err := srs.parseSLARegisterDates(slaRegisterFormDTO)
	if err != nil {
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
	// Document is provided, validate it
	if err := srs.validateSLADocumentFormDTO(slaRegisterFormDTO.Document); err != nil {
		return 0, err
	}

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

	if err := srs.repo.StoreSLADocumentFile(id, slaRegisterFormDTO.Document.FileName); err != nil {
		return 0, err
	}

	slaDocument := &models.SLADocument{
		SLARegisterID: id,
		Version:       slaRegisterFormDTO.Version,
		IssueDate:     dates.IssueDate,
		ExpiryDate:    dates.ExpiryDate,
		SA1Date:       dates.SA1Date,
		SA2Date:       dates.SA2Date,
		Language:      slaRegisterFormDTO.Document.Language,
		URL:           slaRegisterFormDTO.Document.URL, //URL To be decided
		Icon:          slaRegisterFormDTO.Document.Icon,
		Name:          slaRegisterFormDTO.Document.FileName,
		Size:          slaRegisterFormDTO.Document.Size,
		Token:         slaRegisterFormDTO.Document.FileName, // Assuming token is same as fileName
		Type:          slaRegisterFormDTO.Document.Type,
		IsLatest:      true,
		CreatedBy:     currentUserID,
	}

	// Save SLA document
	if err := srs.repo.SaveSLADocument(slaDocument); err != nil {
		return 0, err
	}

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

	if err := srs.validateSLADocumentFormDTO(slaRegisterFormDTO.Document); err != nil {
		return err
	}

	dates, err := srs.parseSLARegisterDates(slaRegisterFormDTO)
	if err != nil {
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

	if err := srs.repo.Update(existingSLARegister); err != nil {
		return err
	}

	// Check if the file is updated
	if slaRegisterFormDTO.Document != nil && slaRegisterFormDTO.Document.ID != nil {
		existingSLADocument, err := srs.repo.GetSLADocumentByID(*slaRegisterFormDTO.Document.ID)
		if err != nil {
			return err
		}
		// Update existing SLA document fields
		existingSLADocument.SLARegisterID = id
		existingSLADocument.Version = slaRegisterFormDTO.Version
		existingSLADocument.IssueDate = dates.IssueDate
		existingSLADocument.ExpiryDate = dates.ExpiryDate
		existingSLADocument.SA1Date = dates.SA1Date
		existingSLADocument.SA2Date = dates.SA2Date
		existingSLADocument.UpdatedBy = currentUserID

		// Check if is_new_file flag is set to true
		if slaRegisterFormDTO.Document.IsNewFile != nil && *slaRegisterFormDTO.Document.IsNewFile {
			// File is updated, delete existing file
			if err := srs.repo.RemoveSLADocumentFile(existingSLADocument.ID,existingSLADocument.SLARegisterID, existingSLADocument.Name); err != nil {
				return err
			}
			// Move new file from temp to storage
			if err := srs.repo.StoreSLADocumentFile(existingSLADocument.SLARegisterID, slaRegisterFormDTO.Document.FileName); err != nil {
				return err
			}

			// Update existing SLA document fields
			existingSLADocument.Language = slaRegisterFormDTO.Document.Language
			existingSLADocument.URL = slaRegisterFormDTO.Document.URL
			existingSLADocument.Icon = slaRegisterFormDTO.Document.Icon
			existingSLADocument.Name = slaRegisterFormDTO.Document.FileName
			existingSLADocument.Size = slaRegisterFormDTO.Document.Size
			existingSLADocument.Token = slaRegisterFormDTO.Document.FileName // Assuming token is same as fileName
			existingSLADocument.Type = slaRegisterFormDTO.Document.Type
			existingSLADocument.IsLatest = true
		}
		// Save updated SLA document
		if err := srs.repo.SaveSLADocument(existingSLADocument); err != nil {
			return err
		}
	}

	go func() {
		err := srs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			complianceModels.SLARegister{}.TableName(),
			id,
			slaRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.UpdateAction,
				"tableName": complianceModels.SLARegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

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

		// Fetch the SLA document by ID
		slaDocument, err := srs.repo.GetLatestSLADocumentBySLARegisterID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete associated documents
		err = srs.repo.RemoveSLADocumentFile(slaDocument.ID,slaDocument.SLARegisterID, slaDocument.Name)
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

func (srs *SLARegisterService) Renew(id uint, currentUserID uint, document *slaDTOs.SLADocumentFormDTO) (*models.SLADocument, error) {
	// Parse SLA register dates
    dates, err := srs.parseSLARegisterDates(document)
    if err != nil {
        return nil, err
    }
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
    // Save updated SLA document
    if err := srs.repo.Renew(slaDocument); err != nil {
        return nil, err
    }

    // Set is_latest to false for remaining SLA documents associated with the same register
    if err := srs.repo.SetIsLatestToFalseForRemainingDocuments(slaDocument.ID, slaDocument.SLARegisterID); err != nil {
        return nil, err
    }

    return id, nil
}


func (srs *SLARegisterService) validateSLARegisterFormDTO(slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) error {
	// Check if Document is provided
	if slaRegisterFormDTO.Document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	}
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

func (srs *SLARegisterService) parseSLARegisterDates(slaRegisterForm *slaDTOs.SLARegisterFormDTO) (*slaDTOs.ParsedSLADates, error) {
	parsed := &slaDTOs.ParsedSLADates{}
	var err error

	// Validate issue_date
	if slaRegisterForm.IssueDate != "" {
		parsed.IssueDate, err = utility.ParseDate(slaRegisterForm.IssueDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
		}
	}

	// Validate expiry_date
	if slaRegisterForm.ExpiryDate != "" {
		parsed.ExpiryDate, err = utility.ParseDate(slaRegisterForm.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	// Validate SA1Date
	if slaRegisterForm.SA1Date != "" {
		parsed.SA1Date, err = utility.ParseDate(slaRegisterForm.SA1Date)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid sa1 date"), "invalid_date", "sa1 date")
		}
	}

	// Validate SA2Date
	if slaRegisterForm.SA2Date != "" {
		parsed.SA2Date, err = utility.ParseDate(slaRegisterForm.SA2Date)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid sa2 date"), "invalid_date", "sa2 date")
		}
	}

	return parsed, nil
}

func (srs *SLARegisterService) validateSLADocumentFormDTO(document *slaDTOs.SLADocumentFormDTO) error {

	// Validate language
	if document.Language == "" {
		return utility.NewValidationError(fmt.Errorf("language is required"), "required", "language")
	}

	// Validate URL
	if document.URL == "" {
		return utility.NewValidationError(fmt.Errorf("URL is required"), "required", "URL")
	}

	// Validate icon
	if document.Icon == "" {
		return utility.NewValidationError(fmt.Errorf("icon is required"), "required", "icon")
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
