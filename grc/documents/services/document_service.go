package services

import (
	"errors"
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentServiceInterface interface {
	GetAll(params *documentDTOs.DocumentListQueryParamsDTO) (*documentDTOs.DocumentListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentDetailsDTO, error)
	Create(currentUserID uint, documentFormDTO *documentDTOs.DocumentFormDTO) (uint, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error
	CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error
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

	url := utility.GenerateDocumentURL("document", strconv.FormatUint(uint64(id), 10), "documents")
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
			IsLatestVersion: 1,
			URL:             url,
			Name:            documentTemplate.Document.Name,
			Size:            documentTemplate.Document.Size,
			Token:           encryptedFilename,
			OriginalName:    &documentTemplate.Document.OriginalName,
			Type:            documentTemplate.Document.Type,
			CreatedBYID:     currentUserID,
		}

		if err := ds.repo.CreateDocument(documentVersion, tx); err != nil {
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

	} else {
		dataToEncrypt := []byte(documentFormDTO.Document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return 0, err
		}
		documentVersion := &models.DocumentVersion{
			DocumentID:      id,
			Version:         documentFormDTO.DocumentVersion,
			IsLatestVersion: 1,
			URL:             url,
			Name:            documentFormDTO.Document.FileName,
			Size:            documentFormDTO.Document.Size,
			Token:           encryptedFilename,
			OriginalName:    &documentFormDTO.Document.OriginalName,
			Type:            documentFormDTO.Document.Type,
			CreatedBYID:     currentUserID,
		}

		if err := ds.repo.CreateDocument(documentVersion, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := utility.StoreFile(documentFormDTO.Document.FileName, destinationDir); err != nil {
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

	for _, id := range ids {
		document, err := ds.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
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
			for _, doc := range document.DocumentVersions {
				documentsToDelete[id] = append(documentsToDelete[id], &doc)
				deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
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
		for documentID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("document", "draft", strconv.FormatUint(uint64(documentID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
			destinationPath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(documentID), 10))
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
		if documentFormDTO.Document == nil {
			return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
		}
		// Validate file type (DOC)
		if documentFormDTO.Document.Type != ".docx" {
			return utility.NewValidationError(fmt.Errorf("only word files are allowed"), "invalid_word_file_type", "type")
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
