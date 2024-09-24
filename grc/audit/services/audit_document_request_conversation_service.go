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
	"time"
)

type AuditDocumentRequestConversationServiceInterface interface {
	GetAll(auditID uint, auditDocumentRequestID uint, queryParams *dtos.AuditDocumentRequestConversationListQueryParamsDTO) ([]*dtos.AuditDocumentRequestConversationListDTO, error)
	Create(auditID uint, auditDocumentRequestID uint, currentUserID uint, formDTO *dtos.AuditDocumentRequestConversationFormDTO) (uint, error)
	GetDocumentPath(auditId uint, auditDocumentRequestID uint, conversationID uint, token string) (string, error)
}

type AuditDocumentRequestConversationService struct {
	repo                          repositories.AuditDocumentRequestConversationRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	auditSvc                      AuditServiceInterface
	auditDocumentRequestSvc       AuditDocumentRequestServiceInterface
	auditDocumentRequestStatusSvc AuditDocumentRequestStatusServiceInterface
}

func NewAuditDocumentRequestConversationService(
	repo repositories.AuditDocumentRequestConversationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDocumentRequestSvc AuditDocumentRequestServiceInterface,
	auditDocumentRequestStatusSvc AuditDocumentRequestStatusServiceInterface,
) AuditDocumentRequestConversationServiceInterface {
	return &AuditDocumentRequestConversationService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		auditSvc:                      auditSvc,
		auditDocumentRequestSvc:       auditDocumentRequestSvc,
		auditDocumentRequestStatusSvc: auditDocumentRequestStatusSvc,
	}
}

func (adrcs *AuditDocumentRequestConversationService) GetAll(auditID uint, auditDocumentRequestID uint, queryParams *dtos.AuditDocumentRequestConversationListQueryParamsDTO) ([]*dtos.AuditDocumentRequestConversationListDTO, error) {
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	if _, err := adrcs.auditDocumentRequestSvc.GetByID(auditDocumentRequestID); err != nil {
		return nil, err
	}

	return adrcs.repo.GetAll(auditDocumentRequestID, queryParams)
}

func (adrcs *AuditDocumentRequestConversationService) Create(auditID uint, auditDocumentRequestID uint, currentUserID uint, formDTO *dtos.AuditDocumentRequestConversationFormDTO) (uint, error) {
	tx := adrcs.repo.BeginTransaction()

	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}
	documentRequest, err := adrcs.auditDocumentRequestSvc.GetByID(auditDocumentRequestID)
	if err != nil {
		return 0, err
	}

	if err := adrcs.validateAuditDocumentRequestConversationFormDTO(formDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	now := time.Now()

	conversation := &models.AuditDocumentRequestConversation{
		AuditDocumentRequestID: auditDocumentRequestID,
		Comment:                formDTO.Comment,
		CommentedBYID:          &currentUserID,
		CommentedAt:            now,
	}

	id, err := adrcs.repo.Create(conversation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.AuditDocumentRequestConversationDocument
	for _, doc := range formDTO.Documents {
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "document-requests", strconv.FormatUint(uint64(auditDocumentRequestID), 10), "conversations", encryptedFilename)

		document := &models.AuditDocumentRequestConversationDocument{
			AuditDocumentRequestConversationID: id,
			URL:                                url,
			Name:                               doc.FileName,
			Size:                               doc.Size,
			Token:                              encryptedFilename,
			OriginalName:                       doc.OriginalName,
			Type:                               doc.Type,
			CreatedByID:                        currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := adrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditID), 10), "document-requests", strconv.FormatUint(uint64(auditDocumentRequestID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	auditDocumentRequestClosedStatus, err := adrcs.auditDocumentRequestStatusSvc.GetByType("closed")
	if err != nil {
		return 0, err
	}
	// Check if IsClose is not nil and its value is true
	if formDTO.IsClose != nil && *formDTO.IsClose {
		if documentRequest.AuditDocumentRequestStatusID != uint(auditDocumentRequestClosedStatus.ID) {
			documentRequest.AuditDocumentRequestStatusID = uint(auditDocumentRequestClosedStatus.ID)
			documentRequest.UpdatedAt = &now
			documentRequest.UpdatedBYID = &currentUserID

			err = adrcs.repo.AuditDocumentRequestStatusUpdate(documentRequest, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	} else {
		//Reopen the status when a new comment is added, if the status is closed
		if documentRequest.AuditDocumentRequestStatusID == uint(auditDocumentRequestClosedStatus.ID) {
			auditDocumentRequestOpenStatus, err := adrcs.auditDocumentRequestStatusSvc.GetByType("open")
			if err != nil {
				return 0, err
			}
			if documentRequest.AuditDocumentRequestStatusID != uint(auditDocumentRequestOpenStatus.ID) {
				documentRequest.AuditDocumentRequestStatusID = uint(auditDocumentRequestOpenStatus.ID)
				documentRequest.UpdatedAt = &now
				documentRequest.UpdatedBYID = &currentUserID

				err = adrcs.repo.AuditDocumentRequestStatusUpdate(documentRequest, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
			}

		}
	}

	//todo -send notifications to users
	tx.Commit()

	go func() {
		err := adrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDocumentRequestConversation{}.TableName(),
			id,
			conversation,
			currentUserID,
			nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDocumentRequestConversation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (adrcs *AuditDocumentRequestConversationService) validateAuditDocumentRequestConversationFormDTO(formDTO *dtos.AuditDocumentRequestConversationFormDTO) error {
	if len(formDTO.Comment) == 0 {
		return errors.New("comment is required")
	}
	if len(formDTO.Comment) > 500 {
		return errors.New("comment exceeds maximum length of 500 characters")
	}
	return nil
}

func (adrcs *AuditDocumentRequestConversationService) GetDocumentPath(auditId uint, auditDocumentRequestID uint, conversationID uint, token string) (string, error) {
	_, err := adrcs.auditSvc.GetByID(auditId)
	if err != nil {
		return "", err
	}
	_, err = adrcs.auditDocumentRequestSvc.GetByID(auditDocumentRequestID)
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

	storagePath := filepath.Join("audits", "audit", strconv.FormatUint(uint64(auditId), 10), "document-requests", strconv.FormatUint(uint64(auditDocumentRequestID), 10), "conversations", strconv.FormatUint(uint64(conversationID), 10), "documents")

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
