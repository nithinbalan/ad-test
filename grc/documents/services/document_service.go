package services

import (
	"errors"
	"fmt"
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	settingDTOs "jethur/grc/settings/dtos"
	settinsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentServiceInterface interface {
	GetAll(params *documentDTOs.DocumentListQueryParamsDTO) (*documentDTOs.DocumentListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentDetailsDTO, error)
	Create(currentUserID uint, documentFormDTO *documentDTOs.DocumentFormDTO) (uint, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	CheckoutDocument(id uint, currentUserID uint, token string) (string, error)
	CheckinDocument(id uint, currentUserID uint, documentCheckinFormDTO *documentDTOs.DocumentCheckinFormDTO) error
	RevertDocumentCheckout(id uint, currentUserID uint) error

	SubmitDocument(id uint, currentUserID uint) error
	PublishDocument(id uint, currentUserID uint) error
	GetComments(id uint) (*documentDTOs.DocumentVersionWorkflowDTO, error)

	CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error
	CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error

	GetDocumentPath(id uint, token string, fileType string, tx ...utility.DatabaseTransaction) (string, error)
	AddComment(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO) error
	GetDocumentWorkflowList(id uint) (*documentDTOs.DocumentVersionWorkflowListDTO, error)

	CreateDocumentAuditLog(documentID uint, documentVersionID uint, userID uint, logType string, tx ...utility.DatabaseTransaction) error
	DirectDocumentAuditLog(documentID uint, userID uint, logType string) error

	GetSharedUsers(documentId uint) ([]*organizationDTOs.UserMinimalDTO, error)
	GetSharedUnits(documentId uint) ([]*settingDTOs.OrganizationUnitDTO, error)
}

type DocumentService struct {
	repo                      repositories.DocumentRepositoryInterface
	documentCategorySvc       DocumentCategoryServiceInterface
	documentTypeSvc           DocumentTypeServiceInterface
	documentTemplateSvc       DocumentTemplateServiceInterface
	documentAccessTypeSvc     DocumentAccessTypeServiceInterface
	documentStatusSvc         DocumentStatusServiceInterface
	userSvc                   organizationSvc.UserServiceInterface
	workflowSvc               settinsSvc.WorkflowServiceInterface
	documentWorkflowSvc       DocumentWorkflowStatusServiceInterface
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface
	activityLogSvc            services.ActivityLogServiceInterface
	documentAuditLogSvc       DocumentAuditLogServiceInterface
}

func NewDocumentService(
	repo repositories.DocumentRepositoryInterface,
	documentCategorySvc DocumentCategoryServiceInterface,
	documentTypeSvc DocumentTypeServiceInterface,
	documentTemplateSvc DocumentTemplateServiceInterface,
	documentAccessTypeSvc DocumentAccessTypeServiceInterface,
	documentStatusSvc DocumentStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	documentWorkflowSvc DocumentWorkflowStatusServiceInterface,
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	documentAuditLogSvc DocumentAuditLogServiceInterface,
) DocumentServiceInterface {
	return &DocumentService{
		repo:                      repo,
		documentCategorySvc:       documentCategorySvc,
		documentTypeSvc:           documentTypeSvc,
		documentTemplateSvc:       documentTemplateSvc,
		documentAccessTypeSvc:     documentAccessTypeSvc,
		documentStatusSvc:         documentStatusSvc,
		userSvc:                   userSvc,
		workflowSvc:               workflowSvc,
		documentWorkflowSvc:       documentWorkflowSvc,
		documentWorkflowStatusSvc: documentWorkflowStatusSvc,
		activityLogSvc:            activityLogSvc,
		documentAuditLogSvc:       documentAuditLogSvc,
	}
}

func (ds *DocumentService) GetAll(params *documentDTOs.DocumentListQueryParamsDTO) (*documentDTOs.DocumentListResponseDTO, error) {
	return ds.repo.GetAll(params)
}

func (crs DocumentService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error) {
	document, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return document, nil
}

func (ds *DocumentService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentDetailsDTO, error) {
	var documentDTO documentDTOs.DocumentDetailsDTO

	document, err := ds.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(document, &documentDTO)
	if err != nil {
		return nil, err
	}

	return &documentDTO, nil
}

func (ds *DocumentService) Create(currentUserID uint, documentFormDTO *documentDTOs.DocumentFormDTO) (uint, error) {
	tx := ds.repo.BeginTransaction()

	if err := ds.validateDocumentFormDTO(documentFormDTO); err != nil {
		return 0, err
	}

	referenceNumber, err := ds.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	responsibleUnitID, err := ds.userSvc.GetUserPrimaryUnitID(uint(documentFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
	}

	documentStatus, err := ds.documentStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	document := &models.Document{
		ReferenceNumber:      referenceNumber,
		DocumentTypeID:       documentFormDTO.DocumentTypeID,
		TitleEn:              documentFormDTO.TitleEn,
		TitleAr:              documentFormDTO.TitleAr,
		DescriptionEn:        documentFormDTO.DescriptionEn,
		DescriptionAr:        documentFormDTO.DescriptionAr,
		DocumentCategoryID:   documentFormDTO.DocumentCategoryID,
		ResponsibleUserID:    documentFormDTO.ResponsibleUserID,
		ResponsibleUnitID:    responsibleUnitID,
		DocumentAccessTypeID: documentFormDTO.DocumentAccessTypeID,
		DocumentStatusID:     documentStatus.ID,
		CreatedBYID:          currentUserID,
	}

	id, err := ds.repo.Create(document, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	documentAccessType, err := ds.documentAccessTypeSvc.GetByID(documentFormDTO.DocumentAccessTypeID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sharedUsers := make([]*models.DocumentSharedUser, 0)
	createdUser := &models.DocumentSharedUser{
		DocumentID:  id,
		UserID:      currentUserID,
		Type:        "created-user",
		CreatedBYID: currentUserID,
	}

	sharedUsers = append(sharedUsers, createdUser)

	responsibleUser := &models.DocumentSharedUser{
		DocumentID:  id,
		UserID:      documentFormDTO.ResponsibleUserID,
		Type:        "responsible-user",
		CreatedBYID: currentUserID,
	}
	sharedUsers = append(sharedUsers, responsibleUser)

	switch documentAccessType.Type {
	case "shared-users":
		for _, userID := range documentFormDTO.SharedUserIDs {
			if userID != nil && *userID != currentUserID && *userID != documentFormDTO.ResponsibleUserID {
				sharedUser := &models.DocumentSharedUser{
					DocumentID:  id,
					UserID:      *userID,
					CreatedBYID: currentUserID,
				}
				sharedUsers = append(sharedUsers, sharedUser)
			}
		}
	case "shared-units":
		sharedUnits := make([]*models.DocumentSharedUnit, 0, len(documentFormDTO.SharedUnitIDs))
		for _, unitID := range documentFormDTO.SharedUnitIDs {
			sharedUnit := &models.DocumentSharedUnit{
				DocumentID:         id,
				OrganizationUnitID: *unitID,
				CreatedBYID:        currentUserID,
			}
			sharedUnits = append(sharedUnits, sharedUnit)
		}
		if err := ds.CreateSharedUnits(sharedUnits, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := ds.CreateSharedUsers(sharedUsers, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	url := utility.GenerateDocumentURL("document", "draft", strconv.FormatUint(uint64(id), 10), "documents")
	destinationDir := filepath.Join("document", "draft", strconv.FormatUint(uint64(id), 10), "documents")
	var docVersionId uint
	if documentFormDTO.DocumentTemplateID != nil {
		var documentTemplate *documentDTOs.DocumentTemplateDetailsDTO
		documentTemplate, err = ds.documentTemplateSvc.GetFullDetailsByID(*documentFormDTO.DocumentTemplateID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		//generate token
		dataToEncrypt := []byte(documentTemplate.Document.Name)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return 0, err
		}

		documentVersion := &models.DocumentVersion{
			DocumentID:      id,
			Version:         documentFormDTO.DocumentVersion,
			IsLatestVersion: true,
			URL:             url,
			Name:            documentTemplate.Document.Name,
			Size:            documentTemplate.Document.Size,
			Token:           encryptedFilename,
			OriginalName:    &documentTemplate.Document.OriginalName,
			Type:            documentTemplate.Document.Type,
			CreatedBYID:     currentUserID,
		}

		docVersionId, err = ds.repo.CreateDocumentVersion(documentVersion, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		sourceDir := filepath.Join("document", "document-template", strconv.FormatUint(uint64(documentTemplate.ID), 10), "documents")
		fileName := documentTemplate.Document.Name

		err = utility.CopyFileFromPaths(sourceDir, destinationDir, fileName)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to copy file: %w", err)
		}

		//add document history
		documentEditHistory := &documentDTOs.DocumentEditHistoryDTO{
			DocumentID:        id,
			DocumentVersionID: docVersionId,
			ActionType:        "create",
			FileName:          &documentTemplate.Document.Name,
			OriginalName:      &documentTemplate.Document.OriginalName,
			Type:              &documentTemplate.Document.Type,
			Size:              &documentTemplate.Document.Size,
		}

		err = ds.CreateDocumentEditHistory(documentEditHistory, &destinationDir, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

	} else {
		dataToEncrypt := []byte(documentFormDTO.Document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return 0, err
		}
		documentVersion := &models.DocumentVersion{
			DocumentID:      id,
			Version:         documentFormDTO.DocumentVersion,
			IsLatestVersion: true,
			URL:             url,
			Name:            documentFormDTO.Document.FileName,
			Size:            documentFormDTO.Document.Size,
			Token:           encryptedFilename,
			OriginalName:    &documentFormDTO.Document.OriginalName,
			Type:            documentFormDTO.Document.Type,
			CreatedBYID:     currentUserID,
		}

		docVersionId, err = ds.repo.CreateDocumentVersion(documentVersion, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := utility.StoreFile(documentFormDTO.Document.FileName, destinationDir); err != nil {
			tx.Rollback()
			return 0, err
		}

		// if err := utility.ConvertDocxToPDF(destinationDir, documentFormDTO.Document.FileName, destinationDir); err != nil {
		// 	tx.Rollback()
		// 	return 0, err
		// }

		//add document history
		documentEditHistory := &documentDTOs.DocumentEditHistoryDTO{
			DocumentID:        id,
			DocumentVersionID: docVersionId,
			ActionType:        "create",
			FileName:          &documentFormDTO.Document.FileName,
			OriginalName:      &documentFormDTO.Document.OriginalName,
			Type:              &documentFormDTO.Document.Type,
			Size:              &documentFormDTO.Document.Size,
		}

		err = ds.CreateDocumentEditHistory(documentEditHistory, &destinationDir, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = ds.CreateDocumentWorkflow(docVersionId, documentFormDTO.DocumentTypeID, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		if err := utility.ConvertDocxToPDF(destinationDir, documentFormDTO.Document.FileName, destinationDir); err != nil {
			log.Printf("Failed to convert document to PDF: %v", err)
		}
	}()

	go func() {
		err := ds.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Document{}.TableName(),
			id,
			documentFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Document{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ds *DocumentService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ds.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*documentDTOs.DocumentVersionDTO)
	documentEditHistoryToDelete := make(map[uint][]*documentDTOs.DocumentEditHistoryDTO)

	for _, id := range ids {
		document, err := ds.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if document.DocumentStatus.Type == "published" {
			failedIds = append(failedIds, id)
			continue
		}

		// Check if document shared user exist
		if len(document.DocumentSharedUsers) > 0 {
			var deleteSharedUserIds []*uint
			for _, user := range document.DocumentSharedUsers {
				deleteSharedUserIds = append(deleteSharedUserIds, &user.ID)
			}
			if err := ds.repo.DeleteSharedUsers(deleteSharedUserIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(document.DocumentSharedUnits) > 0 {
			var deleteSharedUnitIds []*uint
			for _, unit := range document.DocumentSharedUnits {
				deleteSharedUnitIds = append(deleteSharedUnitIds, &unit.ID)
			}
			if err := ds.repo.DeleteSharedUnits(deleteSharedUnitIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(document.DocumentVersions) > 0 {
			var deleteDocumentIds []*uint
			var deleteDocumentEditHistoryIds []*uint
			for _, doc := range document.DocumentVersions {
				documentsToDelete[id] = append(documentsToDelete[id], &doc)
				deleteDocumentIds = append(deleteDocumentIds, &doc.ID)

				if len(doc.DocumentEditHistory) > 0 {
					for _, history := range doc.DocumentEditHistory {
						documentEditHistoryToDelete[id] = append(documentEditHistoryToDelete[id], &history)
						deleteDocumentEditHistoryIds = append(deleteDocumentEditHistoryIds, history.ID)
					}
				}
			}

			if err := ds.repo.DeleteDocumentEditHistory(deleteDocumentEditHistoryIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

			if err := ds.repo.DeleteDocumentVersions(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := ds.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		for documentID, history := range documentEditHistoryToDelete {
			for _, doc := range history {
				if doc.FileName != nil {
					storagePath := filepath.Join("document", "draft", strconv.FormatUint(uint64(documentID), 10), "document-history", strconv.FormatUint(uint64(*doc.ID), 10), "documents")
					if err := utility.RemoveFile(*doc.FileName, storagePath); err != nil {
						tx.Rollback()
						fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
					}
				}

			}
			destinationPath := filepath.Join("compliance", "draft", strconv.FormatUint(uint64(documentID), 10), "document-history")
			if err := utility.RemoveFolder(destinationPath); err != nil {
				tx.Rollback()
			}
		}

		for documentID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("document", "draft", strconv.FormatUint(uint64(documentID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
			destinationPath := filepath.Join("document", "draft", strconv.FormatUint(uint64(documentID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				tx.Rollback()
			}
		}
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := ds.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Document{}.TableName(),
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

func (ds *DocumentService) validateDocumentFormDTO(documentFormDTO *documentDTOs.DocumentFormDTO) error {
	if _, err := ds.documentCategorySvc.GetByID(documentFormDTO.DocumentCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document category does not exist"), "not_exist", "document category")
	}

	if _, err := ds.documentTypeSvc.GetByID(documentFormDTO.DocumentTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document type does not exist"), "not_exist", "document type")
	}

	if _, err := ds.documentAccessTypeSvc.GetByID(documentFormDTO.DocumentAccessTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document access type does not exist"), "not_exist", "document access type")
	}

	if documentFormDTO.DocumentTemplateID != nil {
		if _, err := ds.documentTemplateSvc.GetByID(*documentFormDTO.DocumentTemplateID); err != nil {
			return utility.NewValidationError(fmt.Errorf("document template does not exist"), "not_exist", "document template")
		}
	} else {
		err := ds.ValidateDocument(documentFormDTO.Document)
		if err != nil {
			return err
		}
	}

	accessType, err := ds.documentAccessTypeSvc.GetByID(documentFormDTO.DocumentAccessTypeID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("failed to get document access type"), "not_exist", "document access type")
	}

	switch accessType.Type {
	case "shared-users":
		if documentFormDTO.SharedUserIDs == nil || len(documentFormDTO.SharedUserIDs) == 0 {
			return utility.NewValidationError(fmt.Errorf("shared user IDs are required for shared-user access type"), "required", "shared_users")
		}
	case "shared-units":
		if documentFormDTO.SharedUnitIDs == nil || len(documentFormDTO.SharedUnitIDs) == 0 {
			return utility.NewValidationError(fmt.Errorf("shared unit IDs are required for shared-unit access type"), "required", "shared_units")
		}
	}
	return nil
}

func (ds *DocumentService) ValidateDocument(document *commonDTOs.DocumentFormDTO) error {
	if document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	}
	// Validate file type (DOC)
	if document.Type != ".docx" {
		return utility.NewValidationError(fmt.Errorf("only word files are allowed"), "invalid_word_file_type", "type")
	}

	return nil
}

func (ds *DocumentService) CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error {
	if len(sharedUsers) > 0 {
		if err := ds.repo.CreateSharedUsers(sharedUsers, tx...); err != nil {
			return err
		}
	}
	return nil
}

func (ds *DocumentService) CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error {
	if len(sharedUnits) > 0 {
		if err := ds.repo.CreateSharedUnits(sharedUnits, tx...); err != nil {
			return err
		}
	}
	return nil
}

func (ds *DocumentService) CreateDocumentEditHistory(documentEditHistoryDTO *documentDTOs.DocumentEditHistoryDTO, sourceDir *string, currentUserId uint, tx ...utility.DatabaseTransaction) error {
	if sourceDir != nil {
		dataToEncrypt := []byte(*documentEditHistoryDTO.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return err
		}
		url := utility.GenerateDocumentURL("document", "draft", strconv.FormatUint(uint64(documentEditHistoryDTO.DocumentID), 10), "document-history")
		documentEditHistory := &models.DocumentEditHistory{
			DocumentID:        documentEditHistoryDTO.DocumentID,
			DocumentVersionID: documentEditHistoryDTO.DocumentVersionID,
			ActionType:        documentEditHistoryDTO.ActionType,
			Comment:           documentEditHistoryDTO.Comment,
			URL:               &url,
			Name:              documentEditHistoryDTO.FileName,
			Size:              documentEditHistoryDTO.Size,
			Token:             &encryptedFilename,
			OriginalName:      documentEditHistoryDTO.OriginalName,
			Type:              documentEditHistoryDTO.Type,
			CreatedBYID:       currentUserId,
		}
		id, err := ds.repo.CreateDocumentEditHistory(documentEditHistory, tx...)
		if err != nil {
			return err
		}
		destinationDir := filepath.Join("document", "draft", strconv.FormatUint(uint64(documentEditHistoryDTO.DocumentID), 10), "document-history", strconv.FormatUint(uint64(id), 10), "documents")
		err = utility.CopyFileFromPaths(*sourceDir, destinationDir, *documentEditHistoryDTO.FileName)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	} else {
		documentEditHistory := &models.DocumentEditHistory{
			DocumentID:        documentEditHistoryDTO.DocumentID,
			DocumentVersionID: documentEditHistoryDTO.DocumentVersionID,
			ActionType:        documentEditHistoryDTO.ActionType,
			CreatedBYID:       currentUserId,
		}

		_, err := ds.repo.CreateDocumentEditHistory(documentEditHistory, tx...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *DocumentService) CheckoutDocument(id uint, currentUserID uint, token string) (string, error) {
	tx := ds.repo.BeginTransaction()
	document, err := ds.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if document.DocumentStatus.Type != "draft" && document.DocumentStatus.Type != "in-review" {
		tx.Rollback()
		return "", fmt.Errorf("checkout only allowed for documents in 'draft' or 'in-review' status")
	}

	if document.IsLocked {
		tx.Rollback()
		return "", utility.NewForbiddenError("document already checkouted")
	}

	if document.CreatedByUser.ID != int(currentUserID) && document.ResponsibleUser.ID != int(currentUserID) {
		tx.Rollback()
		return "", fmt.Errorf("access denied: user is neither creator nor responsible user")
	}

	currentDateStr := time.Now().Format("2006-01-02")
	lastCheckoutAt, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("error parsing current date: %w", err)
	}

	existingDocument, err := ds.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	existingDocument.IsLocked = true
	existingDocument.LastCheckoutBYID = &currentUserID
	existingDocument.LastCheckoutAt = lastCheckoutAt
	err = ds.repo.Update(existingDocument, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	//add document history
	documentEditHistory := &documentDTOs.DocumentEditHistoryDTO{
		DocumentID:        id,
		DocumentVersionID: document.LatestDocumentVersion.ID,
		ActionType:        "check-out",
	}

	err = ds.CreateDocumentEditHistory(documentEditHistory, nil, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	var documentPath string
	fileType := "docx"
	documentPath, err = ds.GetDocumentPath(id, token, fileType, tx)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("error retrieving document path: %w", err)
	}

	tx.Commit()
	return documentPath, nil
}

func (ds *DocumentService) CheckinDocument(id uint, currentUserID uint, documentCheckinFormDTO *documentDTOs.DocumentCheckinFormDTO) error {
	tx := ds.repo.BeginTransaction()
	document, err := ds.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = ds.ValidateDocument(&documentCheckinFormDTO.Document)
	if err != nil {
		return err
	}

	if document.DocumentStatus.Type != "draft" && document.DocumentStatus.Type != "in-review" {
		tx.Rollback()
		return fmt.Errorf("checkin only allowed for documents in 'draft' or 'in-review' status")
	}

	if !document.IsLocked {
		tx.Rollback()
		return fmt.Errorf("document already checked in")
	}

	if document.CreatedByUser.ID != int(currentUserID) && document.ResponsibleUser.ID != int(currentUserID) {
		tx.Rollback()
		return fmt.Errorf("access denied: user is neither creator nor responsible user")
	}

	dataToEncrypt := []byte(documentCheckinFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentVersion, err := ds.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentVersion.Name = documentCheckinFormDTO.Document.FileName
	documentVersion.Size = documentCheckinFormDTO.Document.Size
	documentVersion.Token = encryptedFilename
	documentVersion.OriginalName = &documentCheckinFormDTO.Document.OriginalName
	documentVersion.Type = documentCheckinFormDTO.Document.Type
	documentVersion.UpdatedBYID = &currentUserID

	err = ds.repo.UpdateDocumentVersion(documentVersion, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentDateStr := time.Now().Format("2006-01-02")
	lastCheckoutIn, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	existingDocument, err := ds.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingDocument.IsLocked = false
	existingDocument.LastCheckinBYID = &currentUserID
	existingDocument.LastCheckinAt = lastCheckoutIn
	err = ds.repo.Update(existingDocument, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	destinationDir := filepath.Join("document", "draft", strconv.FormatUint(uint64(id), 10), "documents")

	deleteFilePath := filepath.Join("document", "draft", strconv.FormatUint(uint64(id), 10), "documents")
	if err := utility.RemoveFolder(deleteFilePath); err != nil {
		tx.Rollback()
	}

	if err := utility.StoreFile(documentCheckinFormDTO.Document.FileName, destinationDir); err != nil {
		tx.Rollback()
		return err
	}

	//add document history
	documentEditHistory := &documentDTOs.DocumentEditHistoryDTO{
		DocumentID:        id,
		DocumentVersionID: document.LatestDocumentVersion.ID,
		ActionType:        "check-in",
		FileName:          &documentCheckinFormDTO.Document.FileName,
		OriginalName:      &documentCheckinFormDTO.Document.OriginalName,
		Type:              &documentCheckinFormDTO.Document.Type,
		Size:              &documentCheckinFormDTO.Document.Size,
		Comment:           documentCheckinFormDTO.Comment,
	}

	err = ds.CreateDocumentEditHistory(documentEditHistory, &destinationDir, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		if err := utility.ConvertDocxToPDF(destinationDir, documentCheckinFormDTO.Document.FileName, destinationDir); err != nil {
			log.Printf("Failed to convert document to PDF: %v", err)
		}
	}()

	return nil
}

func (ds *DocumentService) RevertDocumentCheckout(id uint, currentUserID uint) error {
	tx := ds.repo.BeginTransaction()
	document, err := ds.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !document.IsLocked {
		tx.Rollback()
		return fmt.Errorf("document already checked in")
	}

	existingDocument, err := ds.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingDocument.IsLocked = false

	err = ds.repo.Update(existingDocument, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//add document history
	documentEditHistory := &documentDTOs.DocumentEditHistoryDTO{
		DocumentID:        id,
		DocumentVersionID: document.LatestDocumentVersion.ID,
		ActionType:        "revert-checkout",
	}

	err = ds.CreateDocumentEditHistory(documentEditHistory, nil, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (ds *DocumentService) GetDocumentPath(id uint, token string, fileType string, tx ...utility.DatabaseTransaction) (string, error) {
	existingDocument, err := ds.GetFullDetailsByID(id, tx...)
	if err != nil {
		return "", err
	}
	fmt.Println("1")

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}
	fmt.Println("2")
	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}
	fmt.Println("3")
	baseFilename := strings.TrimSuffix(decryptedFilename, filepath.Ext(decryptedFilename))

	var storagePath string
	if existingDocument.DocumentStatus.Type == "published" {
		storagePath = filepath.Join("document", "published", strconv.FormatUint(uint64(id), 10), "documents")
	} else {
		storagePath = filepath.Join("document", "draft", strconv.FormatUint(uint64(id), 10), "documents")
	}

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, baseFilename+"."+fileType)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (ds *DocumentService) CreateDocumentWorkflow(documentVersionId uint, documentTypeId uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflow, err := ds.workflowSvc.GetWorkflowByDocumentType(documentTypeId, tx...)
	if err != nil {
		return err
	}

	workflowLevels, err := ds.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	documentWorkflowStatus, err := ds.documentWorkflowSvc.GetByType("pending")

	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		documentWorkflow := &models.DocumentWorkflow{
			DocumentVersionID:        documentVersionId,
			Level:                    level.Level,
			WorkflowLevelTypeID:      level.WorkflowLevelType.ID,
			ApprovalTypes:            level.ApprovalType,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			CreatedBYID:              currentUserID,
		}

		docWorkflowId, err := ds.repo.CreateDocumentWorkflow(documentWorkflow, tx...)
		if err != nil {
			return err
		}

		// Handle participants for each level
		for _, participant := range level.Participants {
			if participant.Role != nil {

				userIds, err := ds.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					err = ds.createWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}

				}
			} else if participant.User != nil {
				err = ds.createWorkflowParticipant(uint(participant.User.ID), docWorkflowId, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (ds *DocumentService) createWorkflowParticipant(userID uint, documentWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	docWorkflowParticipant := &models.DocumentWorkflowParticipant{
		DocumentWorkflowID: documentWorkflowID,
		UserID:             userID,
		CreatedBYID:        currentUserID,
	}

	if participant.Role != nil {
		docWorkflowParticipant.RoleID = &participant.Role.ID
	}

	err := ds.repo.CreateDocumentWorkflowParticipant(docWorkflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DocumentService) SubmitDocument(id uint, currentUserID uint) error {
	tx := ds.repo.BeginTransaction()
	document, err := ds.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentStatus, err := ds.documentStatusSvc.GetByType("in-review")
	if err != nil {
		tx.Rollback()
		return err
	}

	existingDocument, err := ds.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingDocument.DocumentStatusID = documentStatus.ID

	err = ds.repo.Update(existingDocument, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentVersion, err := ds.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentDateStr := time.Now().Format("2006-01-02")
	submittedAt, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	var level int = 1

	documentVersion.SubmittedBYID = &currentUserID
	documentVersion.SubmittedAt = submittedAt
	documentVersion.NextReviewLevel = &level

	err = ds.repo.UpdateDocumentVersion(documentVersion, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	//to do
	//activity log

	return nil
}

func (ds *DocumentService) PublishDocument(id uint, currentUserID uint) error {
	tx := ds.repo.BeginTransaction()
	document, err := ds.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentStatus, err := ds.documentStatusSvc.GetByType("published")
	if err != nil {
		tx.Rollback()
		return err
	}

	existingDocument, err := ds.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingDocument.DocumentStatusID = documentStatus.ID

	err = ds.repo.Update(existingDocument, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentVersion, err := ds.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentDateStr := time.Now().Format("2006-01-02")
	publishedAt, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	documentVersion.PublishedBYID = &currentUserID
	documentVersion.PublishedAt = publishedAt

	err = ds.repo.UpdateDocumentVersion(documentVersion, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Add audit log for the document retrieval
	if err := ds.CreateDocumentAuditLog(uint(id), documentVersion.ID, currentUserID, "published", tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	//to do
	//activity log

	return nil
}

func (ds *DocumentService) GetComments(id uint) (*documentDTOs.DocumentVersionWorkflowDTO, error) {
	var documentVersionWorkflowDTO documentDTOs.DocumentVersionWorkflowDTO

	document, err := ds.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	currentApprovalLevel := document.LatestDocumentVersion.NextReviewLevel

	documentVersionWorkflows, err := ds.repo.GetComments(document.LatestDocumentVersion.ID, *currentApprovalLevel)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(documentVersionWorkflows, &documentVersionWorkflowDTO)
	if err != nil {
		return nil, err
	}

	return &documentVersionWorkflowDTO, err
}

func (ds *DocumentService) AddComment(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO) error {
	document, err := ds.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	currentApprovalLevel := document.LatestDocumentVersion.NextReviewLevel

	currentDocumentWorkflow, err := ds.repo.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
	if err != nil {
		return err
	}

	if documentReviewCommentFormDTO.IsApproved {
		documentWorkflowStatus, err := ds.documentWorkflowStatusSvc.GetByType("approve")
		if err != nil {
			return err
		}
		switch currentDocumentWorkflow.WorkflowLevelType.Type {
		case "individual":
			if currentUserID != currentDocumentWorkflow.DocumentWorkflowParticipants[0].UserID {
				return utility.NewForbiddenError("you have no permission to approve")
			}

			documentWorkflowReview := &models.DocumentWorkflowReview{
				DocumentWorkflowID:       currentDocumentWorkflow.ID,
				DocumentWorkflowStatusID: documentWorkflowStatus.ID,
				Comment:                  documentReviewCommentFormDTO.Comment,
				CreatedBYID:              currentUserID,
			}

			err = ds.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
			if err != nil {
				return err
			}

			documentWorkflow, err := ds.repo.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
			if err != nil {
				return err
			}

			documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
			err = ds.repo.UpdateDocumentWorkflow(documentWorkflow)
			if err != nil {
				return err
			}

			nextDocWorkflowLevel, err := ds.repo.GetNextLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err // Handle other errors that might occur
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				documentStatus, err := ds.documentStatusSvc.GetByType("approved")
				if err != nil {
					return err
				}
				existingDocument, err := ds.repo.GetByID(id)
				if err != nil {
					return err
				}

				existingDocument.DocumentStatusID = documentStatus.ID
				err = ds.repo.Update(existingDocument)
				if err != nil {
					return err
				}

			} else {
				existingdocVersion, err := ds.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
				if err != nil {
					return err
				}

				existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
				err = ds.repo.UpdateDocumentVersion(existingdocVersion)
				if err != nil {
					return err
				}
			}
		case "multiple-users", "role":
			// Check if the current user ID is in the list of participants
			allowed := false
			for _, participant := range currentDocumentWorkflow.DocumentWorkflowParticipants {
				if participant.UserID == currentUserID {
					allowed = true
					break
				}
			}

			if !allowed {
				return utility.NewForbiddenError("you have no permission to approve")
			}

			documentWorkflowReview := &models.DocumentWorkflowReview{
				DocumentWorkflowID:       currentDocumentWorkflow.ID,
				DocumentWorkflowStatusID: documentWorkflowStatus.ID,
				Comment:                  documentReviewCommentFormDTO.Comment,
				CreatedBYID:              currentUserID,
			}

			err = ds.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
			if err != nil {
				return err
			}

			if currentDocumentWorkflow.ApprovalTypes == "all" {
				allApproved, err := ds.CheckAllParticipantsApproved(currentDocumentWorkflow.ID, currentDocumentWorkflow.DocumentWorkflowParticipants)
				if err != nil {
					return err
				}

				if !allApproved {
					return nil // Not all users have approved, so do not proceed to next level
				} else {
					documentWorkflow, err := ds.repo.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
					if err != nil {
						return err
					}

					documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
					err = ds.repo.UpdateDocumentWorkflow(documentWorkflow)
					if err != nil {
						return err
					}

					nextDocWorkflowLevel, err := ds.repo.GetNextLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return err // Handle other errors that might occur
					}

					if errors.Is(err, gorm.ErrRecordNotFound) {
						documentStatus, err := ds.documentStatusSvc.GetByType("approved")
						if err != nil {
							return err
						}
						existingDocument, err := ds.repo.GetByID(id)
						if err != nil {
							return err
						}

						existingDocument.DocumentStatusID = documentStatus.ID
						err = ds.repo.Update(existingDocument)
						if err != nil {
							return err
						}

					} else {
						existingdocVersion, err := ds.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
						if err != nil {
							return err
						}

						existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
						err = ds.repo.UpdateDocumentVersion(existingdocVersion)
						if err != nil {
							return err
						}
					}
				}
			} else {
				documentWorkflow, err := ds.repo.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
				if err != nil {
					return err
				}

				documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
				err = ds.repo.UpdateDocumentWorkflow(documentWorkflow)
				if err != nil {
					return err
				}

				nextDocWorkflowLevel, err := ds.repo.GetNextLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				if errors.Is(err, gorm.ErrRecordNotFound) {
					documentStatus, err := ds.documentStatusSvc.GetByType("approved")
					if err != nil {
						return err
					}
					existingDocument, err := ds.repo.GetByID(id)
					if err != nil {
						return err
					}

					existingDocument.DocumentStatusID = documentStatus.ID
					err = ds.repo.Update(existingDocument)
					if err != nil {
						return err
					}

				} else {
					existingdocVersion, err := ds.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
					if err != nil {
						return err
					}

					existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
					err = ds.repo.UpdateDocumentVersion(existingdocVersion)
					if err != nil {
						return err
					}
				}
			}
		}

	} else if documentReviewCommentFormDTO.IsRejected {
		documentWorkflowStatus, err := ds.documentWorkflowStatusSvc.GetByType("reject")
		if err != nil {
			return err
		}

		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       currentDocumentWorkflow.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  documentReviewCommentFormDTO.Comment,
			CreatedBYID:              currentUserID,
		}

		err = ds.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
		if err != nil {
			return err
		}

		documentWorkflow, err := ds.repo.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
		if err != nil {
			return err
		}

		documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
		err = ds.repo.UpdateDocumentWorkflow(documentWorkflow)
		if err != nil {
			return err
		}

		documentStatus, err := ds.documentStatusSvc.GetByType("rejected")
		if err != nil {
			return err
		}
		existingDocument, err := ds.repo.GetByID(id)
		if err != nil {
			return err
		}

		existingDocument.DocumentStatusID = documentStatus.ID
		err = ds.repo.Update(existingDocument)
		if err != nil {
			return err
		}

	} else if documentReviewCommentFormDTO.IsReverted {
		documentWorkflowStatus, err := ds.documentWorkflowStatusSvc.GetByType("revert")
		if err != nil {
			return err
		}
		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       currentDocumentWorkflow.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  documentReviewCommentFormDTO.Comment,
			CreatedBYID:              currentUserID,
		}

		err = ds.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
		if err != nil {
			return err
		}

	} else {
		documentWorkflowStatus, err := ds.documentWorkflowStatusSvc.GetByType("comment")
		if err != nil {
			return err
		}
		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       currentDocumentWorkflow.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  documentReviewCommentFormDTO.Comment,
			CreatedBYID:              currentUserID,
		}

		err = ds.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *DocumentService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.DocumentWorkflowParticipant) (bool, error) {
	count, err := ds.repo.CountApprovedParticipants(workflowID)
	if err != nil {
		return false, err
	}

	// If the count of unique approvers is less than the number of participants, not all have approved.
	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}

func (ds *DocumentService) GetDocumentWorkflowList(id uint) (*documentDTOs.DocumentVersionWorkflowListDTO, error) {
	document, err := ds.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	documentVersion, err := ds.repo.GetWorkflowList(document.LatestDocumentVersion.ID)

	if err != nil {
		return nil, err
	}

	dto := &documentDTOs.DocumentVersionWorkflowListDTO{
		ID:              documentVersion.ID,
		DocumentID:      documentVersion.DocumentID,
		SubmittedAt:     *documentVersion.SubmittedAt,
		SubmittedByUser: mapUserToDTO(documentVersion.SubmittedByUser),
	}

	for _, workflow := range documentVersion.DocumentWorkflows {
		workflowDTO := &documentDTOs.DocumentWorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: mapStatusToDTO(workflow.DocumentWorkflowStatus),
		}

		//Since all reviews are approved, add directly to ApprovedUsers
		var lastApprovedAt time.Time
		for _, review := range workflow.DocumentWorkflowReviews {
			userDTO := mapUserToDTO(review.CreatedByUser)
			workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, &userDTO)
			// Assume ApprovedAt is the last review's CreatedAt
			if review.CreatedAt.After(lastApprovedAt) {
				lastApprovedAt = review.CreatedAt
			}
		}

		for _, participant := range workflow.DocumentWorkflowParticipants {
			userDTO := mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, &userDTO)
		}
		workflowDTO.ApprovedAt = lastApprovedAt
		dto.DocumentWorkflows = append(dto.DocumentWorkflows, workflowDTO)
	}

	return dto, nil
}

func mapUserToDTO(user *models.User) organizationDTOs.UserMinimalDTO {
	return organizationDTOs.UserMinimalDTO{
		ID:               int(user.ID),
		FirstNameEN:      user.FirstNameEn,
		LastNameEN:       user.LastNameEn,
		FirstNameAR:      user.FirstNameAr,
		LastNameAR:       user.LastNameAr,
		ProfileImagePath: user.ProfileImagePath,
		Email:            user.Email,
		Initial:          user.Initial,
		Color:            user.Color,
	}
}

func mapStatusToDTO(status *models.DocumentWorkflowStatus) documentDTOs.DocumentStatusDTO {
	return documentDTOs.DocumentStatusDTO{
		ID:        int(status.ID),
		Type:      status.Type,
		TitleEn:   status.TitleEn,
		TitleAr:   status.TitleAr,
		ColorCode: status.ColorCode,
		Status:    status.Status,
	}
}

func (ds *DocumentService) CreateDocumentAuditLog(documentID uint, documentVersionID uint, userID uint, logType string, tx ...utility.DatabaseTransaction) error {
	// Get the audit log type by the provided logType
	documentAuditLogType, err := ds.documentAuditLogSvc.GetByType(logType)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	documentAuditLogFormDTO := &documentDTOs.DocumentAuditLogFormDTO{
		DocumentID:             documentID,
		DocumentVersionID:      documentVersionID,
		UserID:                 userID,
		DocumentAuditLogTypeID: documentAuditLogType.ID,
		CreatedAt:              currentTime,
		CreatedByID:            userID,
	}
	if err := ds.documentAuditLogSvc.Create(documentAuditLogFormDTO, tx...); err != nil {
		return err
	}

	return nil
}

func (ds *DocumentService) DirectDocumentAuditLog(documentID uint, userID uint, logType string) error {
	tx := ds.repo.BeginTransaction()
	document, err := ds.GetFullDetailsByID(documentID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	documentVersion, err := ds.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Add audit log for the document retrieval
	if err := ds.CreateDocumentAuditLog(uint(documentID), documentVersion.ID, userID, logType, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil

}

func (ds *DocumentService) GetSharedUsers(documentId uint) ([]*organizationDTOs.UserMinimalDTO, error) {
	users, err := ds.repo.GetSharedUsers(documentId)
	if err != nil {
		return nil, err
	}

	dtos := make([]*organizationDTOs.UserMinimalDTO, 0) // Initialize as empty slice, not nil
	for _, user := range users {
		dto := &organizationDTOs.UserMinimalDTO{
			ID:               int(user.ID),
			FirstNameEN:      user.FirstNameEn,
			LastNameEN:       user.LastNameEn,
			FirstNameAR:      user.FirstNameAr,
			LastNameAR:       user.LastNameAr,
			ProfileImagePath: user.ProfileImagePath,
			Email:            user.Email,
			Initial:          user.Initial,
			Color:            user.Color,
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (svc *DocumentService) GetSharedUnits(documentId uint) ([]*settingDTOs.OrganizationUnitDTO, error) {
	units, err := svc.repo.GetSharedUnits(documentId)
	if err != nil {
		return nil, err
	}

	dtos := make([]*settingDTOs.OrganizationUnitDTO, 0)
	for _, unit := range units {
		dto := &settingDTOs.OrganizationUnitDTO{
			ID:      int(unit.ID),
			TitleEn: unit.TitleEn,
			TitleAr: unit.TitleAr,
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}
