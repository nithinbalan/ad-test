package services

import (
	"errors"
	"fmt"
	"jethur/config"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type AuditFindingActionPlanStatusUpdateServiceInterface interface {
	GetAll(auditFindingID uint, auditFindingActionPlanID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	Create(auditFindingID uint, auditFindingActionPlanID uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) (uint, error)
	Update(auditFindingID uint, auditFindingActionPlanID uint, id uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) error
	GetLatestAuditFindingActionPlanStatusUpdate(auditFindingID uint, actionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(auditFindingID uint, auditFindingActionPlanID uint, id uint, token string) (string, error)
}

type AuditFindingActionPlanStatusUpdateService struct {
	repo                      repositories.AuditFindingActionPlanStatusUpdateRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditFindingActionPlanSvc AuditFindingActionPlanServiceInterface
	actionPlanStatusSvc       masterServices.ActionPlanStatusServiceInterface
}

func NewAuditFindingActionPlanStatusUpdateService(
	repo repositories.AuditFindingActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditFindingActionPlanSvc AuditFindingActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) AuditFindingActionPlanStatusUpdateServiceInterface {
	return &AuditFindingActionPlanStatusUpdateService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditFindingActionPlanSvc: auditFindingActionPlanSvc,
		actionPlanStatusSvc:       actionPlanStatusSvc,
	}
}

func (afs *AuditFindingActionPlanStatusUpdateService) GetAll(auditFindingID uint, auditFindingActionPlanID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID); err != nil {
		return nil, err
	}

	return afs.repo.GetAll(auditFindingActionPlanID, sort, sortBy, keyword)
}

func (afs AuditFindingActionPlanStatusUpdateService) GetByID(auditFindingID uint, auditFindingActionPlanID uint, id uint) (*models.AuditFindingActionPlanStatusUpdate, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID); err != nil {
		return nil, err
	}

	auditFindingActionPlanStatusUpdate, err := afs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFindingActionPlanStatusUpdate, nil
}

func (afs *AuditFindingActionPlanStatusUpdateService) Create(auditFindingID uint, auditFindingActionPlanID uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := afs.repo.BeginTransaction()
	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID); err != nil {
		return 0, err
	}

	if err := afs.validateAuditFindingActionPlanStatusUpdateFormDTO(auditFindingActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	auditFindingActionPlanUpdate := &models.AuditFindingActionPlanStatusUpdate{
		AuditFindingActionPlanID: auditFindingActionPlanID,
		ActionPlanStatusID:       auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                  auditFindingActionPlanStatusUpdateFormDTO.Comment,
		Percentage:               auditFindingActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:              currentUserID,
	}

	id, err := afs.repo.Create(auditFindingActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = afs.auditFindingActionPlanSvc.StatusUpdate(auditFindingID, auditFindingActionPlanID, currentUserID, auditFindingActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.AuditFindingActionPlanStatusUpdateDocument

	for _, doc := range auditFindingActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.AuditFindingActionPlanStatusUpdateDocument{
			AuditFindingActionPlanStatusUpdateID: id,
			URL:                                  url,
			Name:                                 doc.FileName,
			Size:                                 doc.Size,
			Token:                                encryptedFilename,
			OriginalName:                         doc.OriginalName,
			Type:                                 doc.Type,
			CreatedByID:                          currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := afs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	tx.Commit()

	go func() {
		existingAuditFindingActionPlanStatusUpdate, err := afs.GetByID(auditFindingID, auditFindingActionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch audit finding action plan status", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = afs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingActionPlanStatusUpdate{}.TableName(),
			id,
			existingAuditFindingActionPlanStatusUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFindingActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (afs *AuditFindingActionPlanStatusUpdateService) Update(auditFindingID uint, auditFindingActionPlanID uint, id uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) error {
	tx := afs.repo.BeginTransaction()

	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID); err != nil {
		return err
	}

	existingAuditFindingActionPlanStatusUpdate, err := afs.GetByID(auditFindingID, auditFindingActionPlanID, id)

	if err != nil {
		return err
	}

	oldAuditFindingActionPlanStatusUpdate, err := afs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := afs.validateAuditFindingActionPlanStatusUpdateFormDTO(auditFindingActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Audit Finding StatusUpdate-- only comments and percentage
	existingAuditFindingActionPlanStatusUpdate.Comment = auditFindingActionPlanStatusUpdateFormDTO.Comment
	existingAuditFindingActionPlanStatusUpdate.Percentage = auditFindingActionPlanStatusUpdateFormDTO.Percentage
	existingAuditFindingActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = afs.repo.Update(existingAuditFindingActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.AuditFindingActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditFindingActionPlanStatusUpdateDocument

	for _, doc := range auditFindingActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.AuditFindingActionPlanStatusUpdateDocument{
				AuditFindingActionPlanStatusUpdateID: id,
				URL:                                  url,
				Name:                                 doc.FileName,
				Size:                                 doc.Size,
				Token:                                encryptedFileName,
				OriginalName:                         doc.OriginalName,
				Type:                                 doc.Type,
				CreatedByID:                          currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := afs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := afs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := afs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingFindingActionPlanStatusUpdate, err := afs.repo.GetByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"percentage": oldAuditFindingActionPlanStatusUpdate.Percentage,
			"comment":    oldAuditFindingActionPlanStatusUpdate.Comment,
		}

		newData := map[string]interface{}{
			"percentage": existingFindingActionPlanStatusUpdate.Percentage,
			"comment":    existingFindingActionPlanStatusUpdate.Comment,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := afs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditFindingActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditFindingActionPlanStatusUpdate{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (afs AuditFindingActionPlanStatusUpdateService) GetLatestAuditFindingActionPlanStatusUpdate(auditFindingID uint, actionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, actionPlanID); err != nil {
		return nil, err
	}

	return afs.repo.GetLatestAuditFindingActionPlanStatusUpdate(actionPlanID)
}

func (afs *AuditFindingActionPlanStatusUpdateService) validateAuditFindingActionPlanStatusUpdateFormDTO(auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) error {

	if _, err := afs.actionPlanStatusSvc.GetByID(auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (afs *AuditFindingActionPlanStatusUpdateService) GetDocumentPath(auditFindingId uint, actionPlanId uint, id uint, token string) (string, error) {
	existingAuditFindingActionPlanStatusUpdate, err := afs.GetByID(auditFindingId, actionPlanId, id)
	if err != nil {
		return "", err
	}

	if len(existingAuditFindingActionPlanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the audit finding action plan update")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingId), 10), "action-plans", strconv.FormatUint(uint64(actionPlanId), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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
