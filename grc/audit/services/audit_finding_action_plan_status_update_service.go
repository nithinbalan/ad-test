package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type AuditFindingActionPlanStatusUpdateServiceInterface interface {
	GetAll(auditFindingID uint, auditFindingActionPlanID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	Create(auditFindingID uint, auditFindingActionPlanID uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) (uint, error)
	Update(auditFindingID uint, auditFindingActionPlanID uint, id uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) error
	GetLatestAuditFindingActionPlanStatusUpdate(auditFindingID uint, actionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
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
	if _, err := afs.auditFindingActionPlanSvc.GetByFindingID(auditFindingID, auditFindingActionPlanID); err != nil {
		return nil, err
	}

	return afs.repo.GetAll(auditFindingActionPlanID, sort, sortBy, keyword)
}

func (afs AuditFindingActionPlanStatusUpdateService) GetByID(auditFindingID uint, auditFindingActionPlanID uint, id uint) (*models.AuditFindingActionPlanStatusUpdate, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByFindingID(auditFindingID, auditFindingActionPlanID); err != nil {
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

	if _, err := afs.auditFindingActionPlanSvc.GetByFindingID(auditFindingID, auditFindingActionPlanID); err != nil {
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

	err = afs.auditFindingActionPlanSvc.StatusUpdate(auditFindingID, auditFindingActionPlanID, currentUserID, auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID, tx)
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

		url := utility.GenerateDocumentURL("audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "status-updates", encryptedFilename)

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
		err := afs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingActionPlanStatusUpdate{}.TableName(),
			id,
			auditFindingActionPlanUpdate,
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

	if _, err := afs.auditFindingActionPlanSvc.GetByFindingID(auditFindingID, auditFindingActionPlanID); err != nil {
		return err
	}

	existingAuditFindingActionPlanStatusUpdate, err := afs.GetByID(auditFindingID, auditFindingActionPlanID, id)

	if err != nil {
		return err
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
			url := utility.GenerateDocumentURL("audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "status-updates", encryptedFileName)
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

	return nil
}

func (afs AuditFindingActionPlanStatusUpdateService) GetLatestAuditFindingActionPlanStatusUpdate(auditFindingID uint, actionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByFindingID(auditFindingID, actionPlanID); err != nil {
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
