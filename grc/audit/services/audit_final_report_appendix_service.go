package services

import (
	"errors"
	"fmt"
	"jethur/config"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type AuditFinalReportAppendixServiceInterface interface {
	GetByID(auditID uint, finalReportID uint, id uint) (*models.AuditFinalReportAppendix, error)
	Create(auditID uint, finalReportID uint, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixFormDTO) ([]uint, error)
	Delete(auditID uint, finalReportID uint, id uint, currentUserID uint) (uint, error)
	GetDocumentPath(auditId uint, finalReportID uint, id uint, token string) (string, error)
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

func (afras *AuditFinalReportAppendixService) GetByID(auditID uint, finalReportID uint, id uint) (*models.AuditFinalReportAppendix, error) {
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

func (afras *AuditFinalReportAppendixService) Create(auditID uint, finalReportID uint, currentUserID uint, auditFinalReportAppendixFormDTO *auditDTOs.AuditFinalReportAppendixFormDTO) ([]uint, error) {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	tx := afras.repo.BeginTransaction()
	var appendixIDs []uint
	var addedAppendix []map[string]interface{}
	for _, document := range auditFinalReportAppendixFormDTO.Documents {
		dataToEncrypt := []byte(document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return nil, err
		}

		url := utility.GenerateDocumentURL("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "final-report", "appendix")

		appendix := &models.AuditFinalReportAppendix{
			AuditFinalReportID: finalReportID,
			URL:                url,
			Name:               document.FileName,
			Size:               document.Size,
			Token:              encryptedFilename,
			OriginalName:       document.OriginalName,
			Type:               document.Type,
			CreatedBYID:        currentUserID,
		}
		id, err := afras.repo.Create(appendix, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "final-report", "appendix")
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

func (afras *AuditFinalReportAppendixService) Delete(auditID uint, finalReportID uint, id uint, currentUserID uint) (uint, error) {
	if _, err := afras.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}
	appendix, err := afras.repo.GetByID(id)
	if err != nil {
		return 0, err
	}
	tx := afras.repo.BeginTransaction()
	var deletedAppendix []map[string]interface{}

	storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "final-report", "appendix")
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

func (afras *AuditFinalReportAppendixService) GetDocumentPath(auditId uint, finalReportID uint, id uint, token string) (string, error) {
	_, err := afras.auditSvc.GetByID(auditId)
	if err != nil {
		return "", err
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditId), 10), "final-report", "appendix")

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
