package services

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type AuditFinalReportAppendixServiceInterface interface {
	GetByID(auditID uint, finalReportID uint, sectionID uint, id uint) (*models.AuditFinalReportAppendix, error)
	Create(auditID uint, finalReportID uint, sectionID uint, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixFormDTO) (uint, error)
	Update(auditID uint, finalReportID uint, sectionID uint, id, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixFormDTO) error
	Delete(auditID uint, finalReportID uint, sectionID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AuditFinalReportAppendixService struct {
	repo           repositories.AuditFinalReportAppendixRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	auditSvc       AuditServiceInterface
}

func NewAuditFinalReportAppendixService(
	repo repositories.AuditFinalReportAppendixRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
) AuditFinalReportAppendixServiceInterface {
	return &AuditFinalReportAppendixService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		auditSvc:       auditSvc,
	}
}

func (afras *AuditFinalReportAppendixService) GetByID(auditID uint, finalReportID uint, sectionID uint, id uint) (*models.AuditFinalReportAppendix, error) {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	appendix, err := afras.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return appendix, nil
}

func (afras *AuditFinalReportAppendixService) Create(auditID uint, finalReportID uint, sectionID uint, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixFormDTO) (uint, error) {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}
	tx := afras.repo.BeginTransaction()

	document := auditFinalReportAppendixFormDTO.Document
	dataToEncrypt := []byte(auditFinalReportAppendixFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return 0, err
	}

	url := utility.GenerateDocumentURL("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "appendix")

	appendix := &models.AuditFinalReportAppendix{
		AuditFinalReportID:        finalReportID,
		AuditFinalReportSectionID: sectionID,
		URL:                       url,
		Name:                      document.FileName,
		Size:                      document.Size,
		Token:                     encryptedFilename,
		OriginalName:              document.OriginalName,
		Type:                      document.Type,
		CreatedBYID:               currentUserID,
	}

	id, err := afras.repo.Create(appendix, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "appendix")
	if err := utility.StoreFile(auditFinalReportAppendixFormDTO.Document.FileName, storagePath); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	// Log activity asynchronously
	go func() {
		err = afras.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFinalReportAppendix{}.TableName(),
			id,
			appendix,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFinalReportAppendix{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (afras *AuditFinalReportAppendixService) Update(auditID uint, finalReportID uint, sectionID uint, id uint, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixFormDTO) error {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	tx := afras.repo.BeginTransaction()

	existingAppendix, err := afras.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}
	document := auditFinalReportAppendixFormDTO.Document
	if document.IsNewFile != nil && *document.IsNewFile {
		dataToEncrypt := []byte(document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}
		storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "appendix")
		if existingAppendix.Name != "" {
			if err := utility.RemoveFile(existingAppendix.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}

		url := utility.GenerateDocumentURL("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "appendix")
		if err := utility.StoreFile(document.FileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

		existingAppendix.URL = url
		existingAppendix.Name = document.FileName
		existingAppendix.Size = document.Size
		existingAppendix.Token = encryptedFilename
		existingAppendix.OriginalName = document.OriginalName
	}

	if err := afras.repo.Update(existingAppendix, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (afras *AuditFinalReportAppendixService) Delete(auditID uint, finalReportID uint, sectionID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	tx := afras.repo.BeginTransaction()
	var successfulIDs []uint
	var failedIDs []uint
	var deletedAppendix []map[string]interface{}

	for _, id := range ids {
		appendix, err := afras.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "appendix")
		if err := utility.RemoveFile(appendix.Name, storagePath); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := afras.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
		}
		successfulIDs = append(successfulIDs, id)
		deletedAppendix = append(deletedAppendix, map[string]interface{}{
			"id": id,
		})

	}
	tx.Commit()
	if len(successfulIDs) > 0 {

		go func() {
			for _, id := range successfulIDs {
				err := afras.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditFinalReportAppendix{}.TableName(),
					id,
					deletedAppendix,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditFinalReportAppendix{}.TableName(),
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
