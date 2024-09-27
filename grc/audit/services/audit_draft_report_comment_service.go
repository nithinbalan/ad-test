package services

import (
	"errors"
	"fmt"
	"jethur/config"
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditDraftReportCommentServiceInterface interface {
	GetAll(auditID uint, auditDraftReportID uint, params *dtos.AuditDraftReportCommentListQueryParamsDTO) ([]*dtos.AuditDraftReportCommentListDTO, error)
	GetByID(id uint) (*models.AuditDraftReportComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AuditDraftReportCommentDetailDTO, error)
	Create(auditID uint, auditDraftReportID uint, currentUserID uint, formDTO *dtos.AuditDraftReportCommentFormDTO) (uint, error)
	Update(auditID uint, auditDraftReportID uint, id uint, currentUserID uint, auditDraftReportCommentFormDTO *dtos.AuditDraftReportCommentFormDTO) error
	Delete(auditID uint, auditDraftReportID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(auditID uint, auditDraftReportID uint, id uint, currentUserID uint, auditDraftReportCommentStatusUpdateFormDTO *dtos.AuditDraftReportCommentStatusUpdateFormDTO) error
	GetDocumentPath(auditId uint, auditDraftReportID uint, id uint, token string) (string, error)
}

type AuditDraftReportCommentService struct {
	repo                repositories.AuditDraftReportCommentRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	auditSvc            AuditServiceInterface
	auditDraftReportSvc AuditDraftReportServiceInterface
}

func NewAuditDraftReportCommentService(
	repo repositories.AuditDraftReportCommentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDraftReportSvc AuditDraftReportServiceInterface,
) AuditDraftReportCommentServiceInterface {
	return &AuditDraftReportCommentService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		auditSvc:            auditSvc,
		auditDraftReportSvc: auditDraftReportSvc,
	}
}

func (adrcs *AuditDraftReportCommentService) GetAll(auditID uint, auditDraftReportID uint, params *dtos.AuditDraftReportCommentListQueryParamsDTO) ([]*dtos.AuditDraftReportCommentListDTO, error) {
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	if _, err := adrcs.auditDraftReportSvc.GetByID(auditDraftReportID); err != nil {
		return nil, err
	}

	return adrcs.repo.GetAll(params)
}

func (adrcs AuditDraftReportCommentService) GetByID(id uint) (*models.AuditDraftReportComment, error) {
	auditDraftReportComment, err := adrcs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDraftReportComment, nil
}

func (adrcs AuditDraftReportCommentService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AuditDraftReportCommentDetailDTO, error) {
	var auditDraftReportCommentDTO dtos.AuditDraftReportCommentDetailDTO
	auditDraftReportComment, err := adrcs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditDraftReportComment, &auditDraftReportCommentDTO)
	if err != nil {
		return nil, err
	}

	auditDraftReportCommentDTO.CreatedAt = auditDraftReportComment.CreatedAt
	auditDraftReportCommentDTO.UpdatedAt = auditDraftReportComment.UpdatedAt

	return &auditDraftReportCommentDTO, nil
}
func (adrcs *AuditDraftReportCommentService) Create(auditID uint, auditDraftReportID uint, currentUserID uint, formDTO *dtos.AuditDraftReportCommentFormDTO) (uint, error) {
	tx := adrcs.repo.BeginTransaction()
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}
	_, err := adrcs.auditDraftReportSvc.GetByID(auditDraftReportID)
	if err != nil {
		return 0, err
	}

	auditDraftReportComment := &models.AuditDraftReportComment{
		AuditDraftReportID:  auditDraftReportID,
		AuditDraftFindingID: formDTO.AuditDraftFindingID,
		Comment:             formDTO.Comment,
		Type:                formDTO.Type,
		CreatedBYID:         currentUserID,
	}

	id, err := adrcs.repo.Create(auditDraftReportComment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var documents []*models.AuditDraftReportCommentDocument
	for _, doc := range formDTO.Documents {
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "draft-report", strconv.FormatUint(uint64(auditDraftReportID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents", encryptedFilename)

		document := &models.AuditDraftReportCommentDocument{
			AuditDraftReportCommentID: id,
			URL:                       url,
			Name:                      doc.FileName,
			Size:                      doc.Size,
			Token:                     encryptedFilename,
			OriginalName:              doc.OriginalName,
			Type:                      doc.Type,
			CreatedBYID:               currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := adrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "draft-report", strconv.FormatUint(uint64(auditDraftReportID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := adrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftReportComment{}.TableName(),
			id,
			auditDraftReportComment,
			currentUserID,
			nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftReportComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (adrcs *AuditDraftReportCommentService) Update(auditID uint, auditDraftReportID uint, id uint, currentUserID uint, auditDraftReportCommentFormDTO *dtos.AuditDraftReportCommentFormDTO) error {
	tx := adrcs.repo.BeginTransaction()
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	if _, err := adrcs.auditDraftReportSvc.GetByID(auditDraftReportID); err != nil {
		return err
	}

	existingAuditDraftReportComment, err := adrcs.GetByID(id)

	if err != nil {
		return err
	}

	// Update Audit Draft Report Comment
	existingAuditDraftReportComment.Comment = auditDraftReportCommentFormDTO.Comment
	existingAuditDraftReportComment.Type = auditDraftReportCommentFormDTO.Type
	existingAuditDraftReportComment.AuditDraftFindingID = auditDraftReportCommentFormDTO.AuditDraftFindingID
	existingAuditDraftReportComment.UpdatedBYID = &currentUserID

	err = adrcs.repo.Update(existingAuditDraftReportComment, tx)
	if err != nil {
		return err
	}

	//update document
	var documents []*models.AuditDraftReportCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditDraftReportCommentDocument

	for _, doc := range auditDraftReportCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "draft-report", strconv.FormatUint(uint64(auditDraftReportID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

			// Create a new Control test Document instance
			document := &models.AuditDraftReportCommentDocument{
				AuditDraftReportCommentID: id,
				URL:                       url,
				Name:                      doc.FileName,
				Size:                      doc.Size,
				Token:                     encryptedFileName,
				OriginalName:              doc.OriginalName,
				Type:                      doc.Type,
				CreatedBYID:               currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := adrcs.repo.GetDocumentByID(*doc.ID)
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
		if err := adrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := adrcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "draft-report", strconv.FormatUint(uint64(auditDraftReportID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "draft-report", strconv.FormatUint(uint64(auditDraftReportID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	go func() {
		err := adrcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditDraftReportComment{}.TableName(),
			id,
			auditDraftReportCommentFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditDraftReportComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (adrcs *AuditDraftReportCommentService) Delete(auditID uint, auditDraftReportID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := adrcs.repo.BeginTransaction()
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	if _, err := adrcs.auditDraftReportSvc.GetByID(auditDraftReportID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*dtos.AuditDraftReportCommentDocumentDTO)
	for _, id := range ids {

		auditDraftReportComment, err := adrcs.GetFullDetailsByID(id, tx)

		if err != nil {
			continue
		}
		if len(auditDraftReportComment.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range auditDraftReportComment.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := adrcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		err = adrcs.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}
	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for auditDraftReportCommentID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "draft-report", strconv.FormatUint(uint64(auditDraftReportID), 10), "comments", strconv.FormatUint(uint64(auditDraftReportCommentID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := adrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditFinding{}.TableName(),
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

func (adrcs *AuditDraftReportCommentService) StatusUpdate(auditID uint, auditDraftReportID uint, id uint, currentUserID uint, auditDraftReportCommentStatusUpdateFormDTO *dtos.AuditDraftReportCommentStatusUpdateFormDTO) error {
	tx := adrcs.repo.BeginTransaction()
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	if _, err := adrcs.auditDraftReportSvc.GetByID(auditDraftReportID); err != nil {
		return err
	}

	existingAuditDraftReportComment, err := adrcs.GetByID(id)
	if err != nil {
		return err
	}

	// Update Audit Draft Report Comment
	existingAuditDraftReportComment.IsResolved = *auditDraftReportCommentStatusUpdateFormDTO.IsResolved
	existingAuditDraftReportComment.UpdatedBYID = &currentUserID

	err = adrcs.repo.Update(existingAuditDraftReportComment, tx)
	if err != nil {
		return err
	}
	tx.Commit()

	go func() {
		err := adrcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditDraftReportComment{}.TableName(),
			id,
			auditDraftReportCommentStatusUpdateFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditDraftReportComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (adrcs *AuditDraftReportCommentService) GetDocumentPath(auditID uint, auditDraftReportID uint, id uint, token string) (string, error) {
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return "", err
	}
	if _, err := adrcs.auditDraftReportSvc.GetByID(auditDraftReportID); err != nil {
		return "", err
	}

	_, err := adrcs.GetByID(id)
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

	storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "draft-report", strconv.FormatUint(uint64(auditDraftReportID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

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
