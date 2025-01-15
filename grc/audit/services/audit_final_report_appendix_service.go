package services

import (
	"errors"
	"fmt"
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
	Create(auditID uint, finalReportID uint, sectionID uint, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixFormDTO) ([]uint, error)
	Update(auditID uint, finalReportID uint, sectionID uint, id, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixUpdateFormDTO) error
	Delete(auditID uint, finalReportID uint, sectionID uint, id uint, currentUserID uint) (uint, error)
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

func (afras *AuditFinalReportAppendixService) Create(auditID uint, finalReportID uint, sectionID uint, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixFormDTO) ([]uint, error) {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	tx := afras.repo.BeginTransaction()
	var appendixIDs []uint
	var addedAppendix []map[string]interface{}
	fmt.Println("Final ID=", finalReportID)
	fmt.Println("Sect ID=", sectionID)
	for _, document := range auditFinalReportAppendixFormDTO.Documents {
		dataToEncrypt := []byte(document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return nil, err
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
			return nil, err
		}
		storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "appendix")
		if err := utility.StoreFile(document.FileName, storagePath); err != nil {
			tx.Rollback()
			return nil, err
		}
		// Append the created workpaper ID to the list
		appendixIDs = append(appendixIDs, id)

		addedAppendix = append(addedAppendix, map[string]interface{}{
			"id":       id,
			"appendix": appendix,
		})
	}
	tx.Commit()
	// Log activity asynchronously
	go func() {
		err := afras.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.AuditFinalReportAppendix{}.TableName(),
			appendixIDs,
			addedAppendix,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFinalReportAppendix{}.TableName(),
				"recordID":  auditID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()
	return appendixIDs, nil
}

func (afras *AuditFinalReportAppendixService) Update(auditID uint, finalReportID uint, sectionID uint, id uint, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixUpdateFormDTO) error {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	tx := afras.repo.BeginTransaction()

	existingAppendix, err := afras.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}
	document := auditFinalReportAppendixFormDTO.Documents
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

func (afras *AuditFinalReportAppendixService) Delete(auditID uint, finalReportID uint, sectionID uint, id uint, currentUserID uint) (uint, error) {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}
	appendix, err := afras.repo.GetByID(id)
	if err != nil {
		return 0, err
	}
	tx := afras.repo.BeginTransaction()
	var deletedAppendix []map[string]interface{}

	storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "appendix")
	if err := utility.RemoveFile(appendix.Name, storagePath); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := afras.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	deletedAppendix = append(deletedAppendix, map[string]interface{}{
		"id": id,
	})

	tx.Commit()

	go func() {
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
	}()
	return id, nil
}
