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
	organizationSvc "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
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

	DocumentCheckout(id uint, currentUserID uint, token string) (string, error)
	DocumentCheckin(id uint, currentUserID uint, documentCheckinFormDTO *documentDTOs.DocumentCheckinFormDTO) error
	DocumentCheckoutRevert(id uint, currentUserID uint) error

	CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error
	CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error

	GetDocumentPath(id uint, token string, tx ...utility.DatabaseTransaction) (string, error)
}

type DocumentService struct {
	repo                  repositories.DocumentRepositoryInterface
	documentCategorySvc   DocumentCategoryServiceInterface
	documentTypeSvc       DocumentTypeServiceInterface
	documentTemplateSvc   DocumentTemplateServiceInterface
	documentAccessTypeSvc DocumentAccessTypeServiceInterface
	documentStatusSvc     DocumentStatusServiceInterface
	userSvc               organizationSvc.UserServiceInterface
	activityLogSvc        services.ActivityLogServiceInterface
}

func NewDocumentService(
	repo repositories.DocumentRepositoryInterface,
	documentCategorySvc DocumentCategoryServiceInterface,
	documentTypeSvc DocumentTypeServiceInterface,
	documentTemplateSvc DocumentTemplateServiceInterface,
	documentAccessTypeSvc DocumentAccessTypeServiceInterface,
	documentStatusSvc DocumentStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	activityLogSvc services.ActivityLogServiceInterface) DocumentServiceInterface {
	return &DocumentService{
		repo:                  repo,
		documentCategorySvc:   documentCategorySvc,
		documentTypeSvc:       documentTypeSvc,
		documentTemplateSvc:   documentTemplateSvc,
		documentAccessTypeSvc: documentAccessTypeSvc,
		documentStatusSvc:     documentStatusSvc,
		userSvc:               userSvc,
		activityLogSvc:        activityLogSvc,
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
		fmt.Println("asdfg")
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
		DocumentStatusID:     uint(documentStatus.ID),
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

		docVersionId, err := ds.repo.CreateDocument(documentVersion, tx)
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
		documentHistory := &documentDTOs.DcumentHistoryDTO{
			DocumentID:        id,
			DocumentVersionID: docVersionId,
			ActionType:        "create",
			FileName:          &documentTemplate.Document.Name,
			OriginalName:      &documentTemplate.Document.OriginalName,
			Type:              &documentTemplate.Document.Type,
			Size:              &documentTemplate.Document.Size,
		}

		err = ds.CreateDocumentHistory(documentHistory, &destinationDir, currentUserID, tx)
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

		docVersionId, err := ds.repo.CreateDocument(documentVersion, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := utility.StoreFile(documentFormDTO.Document.FileName, destinationDir); err != nil {
			tx.Rollback()
			return 0, err
		}

		//add document history
		documentHistory := &documentDTOs.DcumentHistoryDTO{
			DocumentID:        id,
			DocumentVersionID: docVersionId,
			ActionType:        "create",
			FileName:          &documentFormDTO.Document.FileName,
			OriginalName:      &documentFormDTO.Document.OriginalName,
			Type:              &documentFormDTO.Document.Type,
			Size:              &documentFormDTO.Document.Size,
		}

		err = ds.CreateDocumentHistory(documentHistory, &destinationDir, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

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
	documentEditHistoryToDelete := make(map[uint][]*documentDTOs.DcumentHistoryDTO)

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

			if err := ds.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
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
		// if documentFormDTO.Document == nil {
		// 	return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
		// }
		// // Validate file type (DOC)
		// if documentFormDTO.Document.Type != ".docx" {
		// 	return utility.NewValidationError(fmt.Errorf("only word files are allowed"), "invalid_word_file_type", "type")
		// }
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

func (ds *DocumentService) CreateDocumentHistory(documentHistoryDTO *documentDTOs.DcumentHistoryDTO, sourceDir *string, currentUserId uint, tx ...utility.DatabaseTransaction) error {
	if sourceDir != nil {
		dataToEncrypt := []byte(*documentHistoryDTO.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return err
		}
		url := utility.GenerateDocumentURL("document", "draft", strconv.FormatUint(uint64(documentHistoryDTO.DocumentID), 10), "document-history")
		documentEditHistory := &models.DocumentEditHistory{
			DocumentID:        documentHistoryDTO.DocumentID,
			DocumentVersionID: documentHistoryDTO.DocumentVersionID,
			ActionType:        documentHistoryDTO.ActionType,
			Comment:           documentHistoryDTO.Comment,
			URL:               &url,
			Name:              documentHistoryDTO.FileName,
			Size:              documentHistoryDTO.Size,
			Token:             &encryptedFilename,
			OriginalName:      documentHistoryDTO.OriginalName,
			Type:              documentHistoryDTO.Type,
			CreatedBYID:       currentUserId,
		}
		id, err := ds.repo.CreateDocumentHistory(documentEditHistory, tx...)
		if err != nil {
			return err
		}
		destinationDir := filepath.Join("document", "draft", strconv.FormatUint(uint64(documentHistoryDTO.DocumentID), 10), "document-history", strconv.FormatUint(uint64(id), 10), "documents")
		err = utility.CopyFileFromPaths(*sourceDir, destinationDir, *documentHistoryDTO.FileName)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	} else {
		documentEditHistory := &models.DocumentEditHistory{
			DocumentID:        documentHistoryDTO.DocumentID,
			DocumentVersionID: documentHistoryDTO.DocumentVersionID,
			ActionType:        documentHistoryDTO.ActionType,
			CreatedBYID:       currentUserId,
		}

		_, err := ds.repo.CreateDocumentHistory(documentEditHistory, tx...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *DocumentService) DocumentCheckout(id uint, currentUserID uint, token string) (string, error) {
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
		return "", fmt.Errorf("document already checked out")
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
	documentHistory := &documentDTOs.DcumentHistoryDTO{
		DocumentID:        id,
		DocumentVersionID: document.DocumentLetastVersion.ID,
		ActionType:        "check-out",
	}

	err = ds.CreateDocumentHistory(documentHistory, nil, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	var documentPath string
	documentPath, err = ds.GetDocumentPath(id, token, tx)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("error retrieving document path: %w", err)
	}

	tx.Commit()
	return documentPath, nil
}

func (ds *DocumentService) DocumentCheckin(id uint, currentUserID uint, documentCheckinFormDTO *documentDTOs.DocumentCheckinFormDTO) error {
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

	documentVersion, err := ds.repo.GetDocumentByID(document.DocumentLetastVersion.ID, tx)
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

	err = ds.repo.UpdateDocument(documentVersion, tx)
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
	documentHistory := &documentDTOs.DcumentHistoryDTO{
		DocumentID:        id,
		DocumentVersionID: document.DocumentLetastVersion.ID,
		ActionType:        "check-in",
		FileName:          &documentCheckinFormDTO.Document.FileName,
		OriginalName:      &documentCheckinFormDTO.Document.OriginalName,
		Type:              &documentCheckinFormDTO.Document.Type,
		Size:              &documentCheckinFormDTO.Document.Size,
		Comment:           documentCheckinFormDTO.Comment,
	}

	err = ds.CreateDocumentHistory(documentHistory, &destinationDir, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (ds *DocumentService) DocumentCheckoutRevert(id uint, currentUserID uint) error {
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
	documentHistory := &documentDTOs.DcumentHistoryDTO{
		DocumentID:        id,
		DocumentVersionID: document.DocumentLetastVersion.ID,
		ActionType:        "revert",
	}

	err = ds.CreateDocumentHistory(documentHistory, nil, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (ds *DocumentService) GetDocumentPath(id uint, token string, tx ...utility.DatabaseTransaction) (string, error) {
	existingDocument, err := ds.GetFullDetailsByID(id, tx...)
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

	var storagePath string
	if existingDocument.DocumentStatus.Type == "published" {
		storagePath = filepath.Join("document", "published", strconv.FormatUint(uint64(id), 10), "documents")
	} else {
		storagePath = filepath.Join("document", "draft", strconv.FormatUint(uint64(id), 10), "documents")
	}

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
