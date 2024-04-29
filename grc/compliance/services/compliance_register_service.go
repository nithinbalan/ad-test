package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceModels "jethur/grc/compliance/models"
	complianceRepositories "jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceRegisterServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint) (*complianceDTOs.ComplianceRegisterDetailDTO, error)
	GetByID(id uint) (*complianceModels.ComplianceRegister, error)
	Create(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, complianceStatusID uint) error
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
}

func NewComplianceRegisterService(
	repo complianceRepositories.ComplianceRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	complianceCategorySvc ComplianceCategoryServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface) ComplianceRegisterServiceInterface {
	return &ComplianceRegisterService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		lifecycleStatusSvc:    lifecycleStatusSvc,
		complianceCategorySvc: complianceCategorySvc,
		complianceStatusSvc:   complianceStatusSvc,
		reviewFrequencySvc:    reviewFrequencySvc,
		severityLevelSvc:      severityLevelSvc,
	}
}

func (crs *ComplianceRegisterService) GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error) {
	return crs.repo.GetAll(params)
}

func (crs ComplianceRegisterService) GetByID(id uint) (*complianceModels.ComplianceRegister, error) {
	compliance, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return compliance, nil
}

func (crs ComplianceRegisterService) GetFullDetailsByID(id uint) (*complianceDTOs.ComplianceRegisterDetailDTO, error) {

	// Initialize your DTO
	var complianceDTO complianceDTOs.ComplianceRegisterDetailDTO

	compliance, err := crs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(compliance, &complianceDTO)
	if err != nil {
		return nil, err
	}

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

	complianceRegister := &complianceModels.ComplianceRegister{
		ComplianceCategoryID:    complianceRegisterFormDTO.ComplianceCategoryID,
		ComplianceStatusID:      complianceRegisterFormDTO.ComplianceStatusID,
		ResponsibleUnitID:       complianceRegisterFormDTO.ResponsibleUnitID,
		ResponsibleUserID:       complianceRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID:       complianceRegisterFormDTO.ReviewFrequencyID,
		ReviewResponsibleUserID: complianceRegisterFormDTO.ReviewResponsibleUserID,
		SeverityLevelID:         complianceRegisterFormDTO.SeverityLevelID,
		ReferenceNumber:         referenceNumber,
		StandardCode:            complianceRegisterFormDTO.StandardCode,
		TitleEn:                 complianceRegisterFormDTO.TitleEn,
		TitleAr:                 complianceRegisterFormDTO.TitleAr,
		DescriptionEn:           complianceRegisterFormDTO.DescriptionEn,
		DescriptionAr:           complianceRegisterFormDTO.DescriptionAr,
		Version:                 complianceRegisterFormDTO.Version,
		Deadline:                dates.Deadline,
		IssueDate:               dates.IssueDate,
		ExpiryDate:              dates.ExpiryDate,
		NextReviewDate:          dates.NextReviewDate,
		CreatedBy:               currentUserID,
		LifecycleStatusID:       lifecycleStatus.ID,
	}

	id, err := crs.repo.Create(complianceRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Store Compliance Document File using utility.StoreFile
	storagePath := filepath.Join("compliance", "compliance", strconv.FormatUint(uint64(id), 10), "documents")

	if err := utility.StoreFile(complianceRegisterFormDTO.Document.FileName, storagePath); err != nil {
		tx.Rollback()
		return 0, err
	}

	//generate token
	dataToEncrypt := []byte(complianceRegisterFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return 0, err
	}

	//genereate URL
	url := utility.GenerateDocumentURL("compliance-register", strconv.FormatUint(uint64(complianceRegister.ID), 10), encryptedFilename)
	// url := fmt.Sprintf("%s/compliance-register/%d/%s", config.GetEnv("API_URL", "http://localhost:8080"), complianceRegister.ID, encryptedFilename)

	complianceDocument := &complianceModels.ComplianceDocument{
		ComplianceRegisterID: id,
		Language:             complianceRegisterFormDTO.Document.Language,
		URL:                  url,
		Icon:                 complianceRegisterFormDTO.Document.Icon,
		Name:                 complianceRegisterFormDTO.Document.FileName,
		Size:                 complianceRegisterFormDTO.Document.Size,
		Token:                encryptedFilename,
		Type:                 complianceRegisterFormDTO.Document.Type,
		CreatedBy:            currentUserID,
	}

	// Save Compliance document
	if err := crs.repo.SaveComplianceDocument(complianceDocument, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := crs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			complianceModels.ComplianceRegister{}.TableName(),
			id,
			complianceRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": complianceModels.ComplianceRegister{}.TableName(),
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

	existingComplianceRegister, err := crs.repo.GetByID(id)

	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
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
	existingComplianceRegister.UpdatedBy = &currentUserID

	err = crs.repo.Update(existingComplianceRegister, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	if complianceRegisterFormDTO.Document.ID != nil {
		existingComplianceDocument, err := crs.repo.GetComplianceDocumentByID(*complianceRegisterFormDTO.Document.ID)
		if err != nil {
			return err
		}
		existingComplianceDocument.ComplianceRegisterID = id
		existingComplianceDocument.UpdatedBy = currentUserID

		if complianceRegisterFormDTO.Document.IsNewFile != nil && *complianceRegisterFormDTO.Document.IsNewFile {
			storagePath := filepath.Join("compliance", "compliance", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.RemoveFile(existingComplianceDocument.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
			// Store Compliance Document File using utility.StoreFile
			if err := utility.StoreFile(complianceRegisterFormDTO.Document.FileName, storagePath); err != nil {
				tx.Rollback()
				return err
			}
			//generate token
			dataToEncrypt := []byte(complianceRegisterFormDTO.Document.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("compliance-register", strconv.FormatUint(uint64(id), 10), encryptedFilename)

			existingComplianceDocument.Language = complianceRegisterFormDTO.Document.Language
			existingComplianceDocument.URL = url
			existingComplianceDocument.Icon = complianceRegisterFormDTO.Document.Icon
			existingComplianceDocument.Name = complianceRegisterFormDTO.Document.FileName
			existingComplianceDocument.Size = complianceRegisterFormDTO.Document.Size
			existingComplianceDocument.Token = encryptedFilename
			existingComplianceDocument.Type = complianceRegisterFormDTO.Document.Type

		}
		if err := crs.repo.SaveComplianceDocument(existingComplianceDocument, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (crs *ComplianceRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := crs.repo.GetByID(id)

		if err != nil {
			continue
		}

		// Fetch the Compliance document by ID
		complianceDocument, err := crs.repo.GetComplianceDocumentByComplianceID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		storagePath := filepath.Join("compliance", "compliance", strconv.FormatUint(uint64(complianceDocument.ComplianceRegisterID), 10), "documents")

		if err := utility.RemoveFile(complianceDocument.Name, storagePath); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		err = crs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				complianceModels.ComplianceRegister{}.TableName(),
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
				logModels.ArchiveAction,
				complianceModels.ComplianceRegister{}.TableName(),
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
				logModels.RestoreAction,
				complianceModels.ComplianceRegister{}.TableName(),
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

func (crs *ComplianceRegisterService) StatusUpdate(id uint, currentUserID uint, complianceStatusID uint) error {
	existingCompliance, err := crs.GetByID(id)

	if err != nil {
		return err
	}

	existingCompliance.ComplianceStatusID = complianceStatusID
	existingCompliance.UpdatedBy = &currentUserID

	err = crs.repo.Update(existingCompliance)

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

	return nil
}

func (crs *ComplianceRegisterService) parseComplianceRegisterDates(dto *complianceDTOs.ComplianceRegisterFormDTO) (*complianceDTOs.ComplianceRegisterParsedDates, error) {
	parsed := &complianceDTOs.ComplianceRegisterParsedDates{}

	var err error
	if dto.Deadline != "" {
		parsed.Deadline, err = utility.ParseDate(dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	if dto.IssueDate != "" {
		parsed.IssueDate, err = utility.ParseDate(dto.IssueDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
		}
	}

	if dto.ExpiryDate != "" {
		parsed.ExpiryDate, err = utility.ParseDate(dto.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	if dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	return parsed, nil
}

func (crs *ComplianceRegisterService) validateComplianceDocumentFormDTO(document *complianceDTOs.ComplianceDocumentFormDTO) error {
	// Validate document
	if document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
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
