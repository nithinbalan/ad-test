package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
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
	Create(auditID uint, currentUserID uint, auditWorkpaperFormDTO *auditDTOs.AuditWorkpaperFormDTO) (uint, error)
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

func (aws *AuditWorkpaperService) Create(auditID uint, currentUserID uint, formDTO *auditDTOs.AuditWorkpaperFormDTO) (uint, error) {
	tx := aws.repo.BeginTransaction()
	if _, err := aws.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := aws.validateAuditWorkpaperFormDTO(formDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	//generate token
	dataToEncrypt := []byte(formDTO.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return 0, err
	}

	//genereate URL
	url := utility.GenerateDocumentURL("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "audit-workpaper")

	workpaper := &models.AuditWorkpaper{
		AuditID:      auditID,
		SourceID:     formDTO.SourceID,
		Title:        formDTO.Title,
		Description:  formDTO.Description,
		URL:          &url,
		Name:         &formDTO.FileName,
		Size:         &formDTO.Size,
		Token:        &encryptedFilename,
		OriginalName: &formDTO.OriginalName,
		CreatedBYID:  currentUserID,
	}

	id, err := aws.repo.Create(workpaper, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Store Compliance Document File using utility.StoreFile
	storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "audit-workpaper")

	if err := utility.StoreFile(formDTO.FileName, storagePath); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := aws.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditWorkpaper{}.TableName(),
			id,
			formDTO,
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
	}()

	return id, nil
}

func (aws *AuditWorkpaperService) Update(id uint, auditID uint, currentUserID uint, formDTO *auditDTOs.AuditWorkpaperFormDTO) error {
	tx := aws.repo.BeginTransaction()
	if _, err := aws.auditSvc.GetByID(auditID); err != nil {
		return err
	}

	if err := aws.validateAuditWorkpaperFormDTO(formDTO); err != nil {
		tx.Rollback()
		return err
	}

	workpaper, err := aws.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	workpaper.SourceID = formDTO.SourceID
	workpaper.Title = formDTO.Title
	workpaper.Description = formDTO.Description
	workpaper.OriginalName = &formDTO.OriginalName
	workpaper.UpdatedBYID = &currentUserID
	// Check if the file is new or updated
	if formDTO.IsNewFile != nil && *formDTO.IsNewFile {
		// Generate new token
		dataToEncrypt := []byte(formDTO.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Generate new URL
		url := utility.GenerateDocumentURL("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "audit-workpaper")

		// Store new document
		storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "audit-workpaper")
		if err := utility.StoreFile(formDTO.FileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

		// Update workpaper with new file details
		workpaper.URL = &url
		workpaper.Name = &formDTO.FileName
		workpaper.Size = &formDTO.Size
		workpaper.Token = &encryptedFilename
		workpaper.OriginalName = &formDTO.OriginalName
	}

	err = aws.repo.Update(workpaper, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := aws.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditWorkpaper{}.TableName(),
			id,
			formDTO,
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
	}()

	return nil
}

func (aws *AuditWorkpaperService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := aws.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		workpaper, err := aws.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		
		if err := aws.repo.Delete(id, tx); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}
		storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(workpaper.AuditID), 10), "audit-workpaper")
		if err := utility.RemoveFile(*workpaper.Name, storagePath); err != nil {
			tx.Rollback()
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
