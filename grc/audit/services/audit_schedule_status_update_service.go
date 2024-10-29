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

type AuditScheduleStatusUpdateServiceInterface interface {
	GetAll(auditScheduleID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditScheduleStatusUpdateListDTO, error)
	Create(auditScheduleID uint, currentUserID uint, auditScheduleStatusUpdateFormDTO *auditDTOs.AuditScheduleStatusUpdateFormDTO) (uint, error)
	Update(auditScheduleID uint, id uint, currentUserID uint, auditScheduleStatusUpdateFormDTO *auditDTOs.AuditScheduleStatusUpdateFormDTO) error
	GetLatestAuditScheduleStatusUpdate(auditScheduleID uint) (*auditDTOs.AuditScheduleStatusUpdateListDTO, error)
	GetDocumentPath(auditID uint, auditScheduleID uint, id uint, token string) (string, error)
}

type AuditScheduleStatusUpdateService struct {
	repo                   repositories.AuditScheduleStatusUpdateRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	auditScheduleSvc       AuditScheduleServiceInterface
	auditSchedulestatusSvc AuditScheduleStatusServiceInterface
}

func NewAuditScheduleStatusUpdateService(
	repo repositories.AuditScheduleStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditScheduleSvc AuditScheduleServiceInterface,
	auditSchedulestatusSvc AuditScheduleStatusServiceInterface,
) AuditScheduleStatusUpdateServiceInterface {
	return &AuditScheduleStatusUpdateService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		auditScheduleSvc:       auditScheduleSvc,
		auditSchedulestatusSvc: auditSchedulestatusSvc,
	}
}

func (asus *AuditScheduleStatusUpdateService) GetAll(auditScheduleID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditScheduleStatusUpdateListDTO, error) {
	return asus.repo.GetAll(auditScheduleID, sort, sortBy, keyword)
}

func (asus AuditScheduleStatusUpdateService) GetByID(auditID uint, auditScheduleID uint, id uint) (*models.AuditScheduleStatusUpdate, error) {
	if _, err := asus.auditScheduleSvc.GetByID(auditScheduleID); err != nil {
		return nil, err
	}

	auditScheduleStatusUpdate, err := asus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditScheduleStatusUpdate, nil
}

func (asus *AuditScheduleStatusUpdateService) Create(auditScheduleID uint, currentUserID uint, auditScheduleStatusUpdateFormDTO *auditDTOs.AuditScheduleStatusUpdateFormDTO) (uint, error) {
	tx := asus.repo.BeginTransaction()
	auditSchedule, err := asus.auditScheduleSvc.GetByID(auditScheduleID)
	if err != nil {
		return 0, err
	}
	auditScheduleStatusUpdate := &models.AuditScheduleStatusUpdate{
		AuditScheduleID:       auditScheduleID,
		AuditScheduleStatusID: auditScheduleStatusUpdateFormDTO.AuditScheduleStatusID,
		Comment:               auditScheduleStatusUpdateFormDTO.Comment,
		Percentage:            auditScheduleStatusUpdateFormDTO.Percentage,
		CreatedBYID:           currentUserID,
	}

	id, err := asus.repo.Create(auditScheduleStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update the audit schedule status
	err = asus.auditScheduleSvc.StatusUpdate(auditScheduleID, auditScheduleStatusUpdateFormDTO, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// This section handles the processing and storage of documents associated with the status update
	var documents []*models.AuditScheduleStatusUpdateDocument
	for _, doc := range auditScheduleStatusUpdateFormDTO.Documents {
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("audit", "audits", strconv.FormatUint(uint64(auditSchedule.AuditID), 10), "audit-schedules", strconv.FormatUint(uint64(auditScheduleID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents", encryptedFilename)

		document := &models.AuditScheduleStatusUpdateDocument{
			AuditScheduleStatusUpdateID: id,
			URL:                         url,
			Name:                        doc.FileName,
			Size:                        doc.Size,
			Token:                       encryptedFilename,
			OriginalName:                doc.OriginalName,
			Type:                        doc.Type,
			CreatedBYID:                 currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := asus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditSchedule.AuditID), 10), "audit-schedules", strconv.FormatUint(uint64(auditScheduleID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	tx.Commit()

	// Logging the creation asynchronously
	go func() {
		err := asus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditScheduleStatusUpdate{}.TableName(),
			id,
			auditScheduleStatusUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditScheduleStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (asus *AuditScheduleStatusUpdateService) Update(auditScheduleID uint, id uint, currentUserID uint, auditScheduleStatusUpdateFormDTO *auditDTOs.AuditScheduleStatusUpdateFormDTO) error {
	tx := asus.repo.BeginTransaction()
	auditSchedule, err := asus.auditScheduleSvc.GetByID(auditScheduleID)
	if err != nil {
		return err
	}
	existingUpdate, err := asus.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update fields based on the form DTO
	existingUpdate.Comment = auditScheduleStatusUpdateFormDTO.Comment
	existingUpdate.Percentage = auditScheduleStatusUpdateFormDTO.Percentage
	existingUpdate.UpdatedBYID = &currentUserID

	// Update the existing update record
	err = asus.repo.Update(existingUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.AuditScheduleStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditScheduleStatusUpdateDocument

	for _, doc := range auditScheduleStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Generate URL
			url := utility.GenerateDocumentURL("audit", "audits", strconv.FormatUint(uint64(auditSchedule.AuditID), 10), "audit-schedules", strconv.FormatUint(uint64(auditScheduleID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
			document := &models.AuditScheduleStatusUpdateDocument{
				AuditScheduleStatusUpdateID: id,
				URL:                         url,
				Name:                        doc.FileName,
				Size:                        doc.Size,
				Token:                       encryptedFileName,
				OriginalName:                doc.OriginalName,
				Type:                        doc.Type,
				CreatedBYID:                 currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := asus.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, *doc.ID)
					}
					return err
				}
				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := asus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := asus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditSchedule.AuditID), 10), "audit-schedules", strconv.FormatUint(uint64(auditScheduleID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditSchedule.AuditID), 10), "audit-schedules", strconv.FormatUint(uint64(auditScheduleID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (asus *AuditScheduleStatusUpdateService) GetLatestAuditScheduleStatusUpdate(auditScheduleID uint) (*auditDTOs.AuditScheduleStatusUpdateListDTO, error) {
	return asus.repo.GetLatestAuditScheduleStatusUpdate(auditScheduleID)
}

func (asus *AuditScheduleStatusUpdateService) GetDocumentPath(auditID uint, auditScheduleID uint, id uint, token string) (string, error) {

	auditSchedule, err := asus.auditScheduleSvc.GetByID(auditScheduleID)
	if err != nil {
		return "", err
	}

	existingAuditScheduleStatusUpdate, err := asus.GetByID(auditID, auditScheduleID, id)
	if err != nil {
		return "", err
	}

	if len(existingAuditScheduleStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the audit finding")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditSchedule.AuditID), 10), "audit-schedules", strconv.FormatUint(uint64(auditScheduleID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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
