package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	mastersServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditWorkpaperServiceInterface interface {
	GetAll(auditID uint, params *auditDTOs.AuditWorkpaperListQueryParamsDTO) (*auditDTOs.AuditWorkpaperListResponseDTO, error)
	GetByID(id uint) (*models.AuditWorkpaper, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditWorkpaperDetailDTO, error)
	Create(auditID uint, currentUserID uint, auditWorkpaperFormDTO *auditDTOs.AuditWorkpaperFormDTO) ([]uint, error)
	Update(id uint, auditID uint, currentUserID uint, auditWorkpaperFormDTO *auditDTOs.AuditWorkpaperFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditWorkpaperService struct {
	repo           repositories.AuditWorkpaperRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	auditSvc       AuditServiceInterface
	sourceSvc      mastersServices.SourceServiceInterface
}

func NewAuditWorkpaperService(
	repo repositories.AuditWorkpaperRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	sourceSvc mastersServices.SourceServiceInterface,
) AuditWorkpaperServiceInterface {
	return &AuditWorkpaperService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		auditSvc:       auditSvc,
		sourceSvc:      sourceSvc,
	}
}

func (aws *AuditWorkpaperService) GetAll(auditID uint, params *auditDTOs.AuditWorkpaperListQueryParamsDTO) (*auditDTOs.AuditWorkpaperListResponseDTO, error) {
	if _, err := aws.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	return aws.repo.GetAll(auditID, params)
}

func (aws *AuditWorkpaperService) GetByID(id uint) (*models.AuditWorkpaper, error) {
	workpaper, err := aws.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return workpaper, nil
}

func (aws *AuditWorkpaperService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditWorkpaperDetailDTO, error) {
	var workpaperDTO auditDTOs.AuditWorkpaperDetailDTO

	workpaper, err := aws.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(workpaper, &workpaperDTO)
	if err != nil {
		return nil, err
	}
	workpaperDTO.CreatedAt = workpaper.CreatedAt
	if workpaper.UpdatedAt != nil {
		workpaperDTO.UpdatedAt = workpaper.UpdatedAt
	}

	return &workpaperDTO, nil
}

func (aws *AuditWorkpaperService) Create(auditID uint, currentUserID uint, formDTO *auditDTOs.AuditWorkpaperFormDTO) ([]uint, error) {
	tx := aws.repo.BeginTransaction()

	// Validate audit ID
	if _, err := aws.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Validate the form DTO
	if err := aws.validateAuditWorkpaperFormDTO(formDTO); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Prepare to collect IDs of created workpapers
	var workpaperIDs []uint

	// Loop through each document and create a workpaper
	for _, document := range formDTO.Documents {
		// Generate token for the document file name
		dataToEncrypt := []byte(document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Generate URL for the document
		url := utility.GenerateDocumentURL("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "audit-workpaper")

		// Create workpaper model
		workpaper := &models.AuditWorkpaper{
			AuditID:      auditID,
			SourceID:     formDTO.SourceID,
			Title:        formDTO.Title,
			Description:  formDTO.Description,
			URL:          &url,
			Name:         &document.FileName,
			Size:         &document.Size,
			Token:        &encryptedFilename,
			OriginalName: &document.OriginalName,
			CreatedBYID:  currentUserID,
		}

		// Save the workpaper in the database
		id, err := aws.repo.Create(workpaper, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Store the document file
		storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "audit-workpaper")
		if err := utility.StoreFile(document.FileName, storagePath); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Append the created workpaper ID to the list
		workpaperIDs = append(workpaperIDs, id)

		// Log activity asynchronously
		go func(id uint, document *commonDTOs.DocumentFormDTO) {
			err := aws.activityLogSvc.LogActivity(
				models.CreateAction,
				models.AuditWorkpaper{}.TableName(),
				id,
				document,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.AuditWorkpaper{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}
				utility.LogError("Failed to log activity", fields)
			}
		}(id, document)
	}

	tx.Commit()
	return workpaperIDs, nil
}

func (aws *AuditWorkpaperService) Update(id uint, auditID uint, currentUserID uint, formDTO *auditDTOs.AuditWorkpaperFormDTO) error {
	tx := aws.repo.BeginTransaction()

	// Validate audit ID
	if _, err := aws.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return err
	}

	// Validate the form DTO
	if err := aws.validateAuditWorkpaperFormDTO(formDTO); err != nil {
		tx.Rollback()
		return err
	}

	// Loop through each document and update the corresponding workpaper
	for _, document := range formDTO.Documents {
		// Retrieve the existing workpaper by ID
		workpaper, err := aws.repo.GetByID(id)
		if err != nil {
			tx.Rollback()
			return utility.NewResourceNotFoundError(err, id)
		}

		// Update workpaper fields
		workpaper.SourceID = formDTO.SourceID
		workpaper.Title = formDTO.Title
		workpaper.Description = formDTO.Description
		workpaper.UpdatedBYID = &currentUserID

		// Check if the file is new or updated
		if document.IsNewFile != nil && *document.IsNewFile {
			// Generate new token
			dataToEncrypt := []byte(document.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Generate new URL
			url := utility.GenerateDocumentURL("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "audit-workpaper")

			// Store new document
			storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "audit-workpaper")
			if err := utility.StoreFile(document.FileName, storagePath); err != nil {
				tx.Rollback()
				return err
			}

			// Update workpaper with new file details
			workpaper.URL = &url
			workpaper.Name = &document.FileName
			workpaper.Size = &document.Size
			workpaper.Token = &encryptedFilename
			workpaper.OriginalName = &document.OriginalName
		}

		// Update the workpaper in the database
		err = aws.repo.Update(workpaper, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Log activity asynchronously
		go func(id uint, document *commonDTOs.DocumentFormDTO) {
			err := aws.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditWorkpaper{}.TableName(),
				id,
				document,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditWorkpaper{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}
				utility.LogError("Failed to log activity", fields)
			}
		}(id, document)
	}

	tx.Commit()
	return nil
}

func (aws *AuditWorkpaperService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := aws.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		workpaper, err := aws.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			return nil, err
		}

		if err := aws.repo.Delete(id, tx); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			return nil, err
		}
		storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(workpaper.AuditID), 10), "audit-workpaper")
		if err := utility.RemoveFile(*workpaper.Name, storagePath); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := aws.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditWorkpaper{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditWorkpaper{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}
					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIDs, nil
}

func (aws *AuditWorkpaperService) validateAuditWorkpaperFormDTO(formDTO *auditDTOs.AuditWorkpaperFormDTO) error {
	if _, err := aws.sourceSvc.GetByID(formDTO.SourceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("source does not exist"), "not_exist", "source")
	}
	return nil
}
