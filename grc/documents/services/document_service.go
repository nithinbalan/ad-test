package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	settinsSvc "jethur/grc/settings/services"
	taskServices "jethur/grc/task/services"
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
	SetDocumentWorkflowService(svc DocumentWorkflowServiceInterface)
	SetDocumentAccessService(svc DocumentAccessServiceInterface)

	GetAll(params *documentDTOs.DocumentListQueryParamsDTO, currentUserID uint) (*documentDTOs.DocumentListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentDetailsDTO, error)
	Create(currentUserID uint, documentFormDTO *documentDTOs.DocumentFormDTO) (uint, error)
	Update(id uint, currentUserID uint, documentFormDTO *documentDTOs.DocumentFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	UpdateDocument(document *models.Document, tx ...utility.DatabaseTransaction) error

	CheckoutDocument(id uint, currentUserID uint, token string) (string, error)
	CheckinDocument(id uint, currentUserID uint, documentCheckinFormDTO *documentDTOs.DocumentCheckinFormDTO) error
	RevertDocumentCheckout(id uint, currentUserID uint) error

	SubmitDocument(id uint, currentUserID uint) error
	PublishDocument(id uint, currentUserID uint) error

	GetDocumentPath(id uint, token string, fileType string, tx ...utility.DatabaseTransaction) (string, error)

	CreateDocumentAuditLog(documentID uint, documentVersionID uint, userID uint, logType string, tx ...utility.DatabaseTransaction) error
	DirectDocumentAuditLog(documentID uint, userID uint, logType string) error

	GetDocumentLastEditSummary(documentID uint) (*documentDTOs.DocumentLastEditSummaryDTO, error)

	ExportDocuments(data []*documentDTOs.DocumentListDTO, currentUserID uint) (*os.File, error)

	GetDocumentVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersion, error)
	UpdateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) error
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
	documentWorkflowSvc       DocumentWorkflowServiceInterface
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface
	activityLogSvc            services.ActivityLogServiceInterface
	documentAuditLogSvc       DocumentAuditLogServiceInterface
	userSettingSvc            organizationSvc.UserSettingServiceInterface
	documentAccessSvc         DocumentAccessServiceInterface
	taskSvc                   taskServices.TaskServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
	orgLevelSvc               settinsSvc.OrganizationLevelServiceInterface
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
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	documentAuditLogSvc DocumentAuditLogServiceInterface,
	userSettingSvc organizationSvc.UserSettingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	orgLevelSvc settinsSvc.OrganizationLevelServiceInterface,
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
		documentWorkflowStatusSvc: documentWorkflowStatusSvc,
		activityLogSvc:            activityLogSvc,
		documentAuditLogSvc:       documentAuditLogSvc,
		userSettingSvc:            userSettingSvc,
		taskSvc:                   taskSvc,
		notificationSvc:           notificationSvc,
		orgLevelSvc:               orgLevelSvc,
	}
}

func (ds *DocumentService) SetDocumentWorkflowService(svc DocumentWorkflowServiceInterface) {
	ds.documentWorkflowSvc = svc
}

func (ds *DocumentService) SetDocumentAccessService(svc DocumentAccessServiceInterface) {
	ds.documentAccessSvc = svc
}

func (ds *DocumentService) GetAll(params *documentDTOs.DocumentListQueryParamsDTO, currentUserID uint) (*documentDTOs.DocumentListResponseDTO, error) {
	return ds.repo.GetAll(params, currentUserID)
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

func (ds *DocumentService) GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentDetailsDTO, error) {
	var documentDTO documentDTOs.DocumentDetailsDTO

	document, err := ds.repo.GetFullDetailsByID(id, currentUserID, tx...)
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

	responsibleUnitID, err := ds.userSvc.GetUserPrimaryUnitID(uint(documentFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
	}

	documentStatus, err := ds.documentStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	document := &models.Document{
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
	case "shared":
		for _, userID := range documentFormDTO.SharedUserIDs {
			sharedUser := &models.DocumentSharedUser{
				DocumentID:  id,
				UserID:      *userID,
				CreatedBYID: currentUserID,
			}
			sharedUsers = append(sharedUsers, sharedUser)
		}

		sharedUnits := make([]*models.DocumentSharedUnit, 0, len(documentFormDTO.SharedUnitIDs))
		for _, unitID := range documentFormDTO.SharedUnitIDs {
			sharedUnit := &models.DocumentSharedUnit{
				DocumentID:         id,
				OrganizationUnitID: *unitID,
				CreatedBYID:        currentUserID,
			}
			sharedUnits = append(sharedUnits, sharedUnit)
		}
		if err := ds.documentAccessSvc.CreateUnits(sharedUnits, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := ds.documentAccessSvc.CreateUsers(sharedUsers, tx); err != nil {
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

	err = ds.documentWorkflowSvc.CreateDocumentWorkflow(id, docVersionId, documentFormDTO.DocumentTypeID, currentUserID, tx)
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

func (ds *DocumentService) Update(id uint, currentUserID uint, documentFormDTO *documentDTOs.DocumentFormDTO) error {
	tx := ds.repo.BeginTransaction()
	// Fetch existing document details
	existingDocument, err := ds.repo.GetByID(id)
	if err != nil {
		return err
	}
	existingDocumentAccessTypeID := existingDocument.DocumentAccessTypeID
	documentAccessType, err := ds.documentAccessTypeSvc.GetByID(documentFormDTO.DocumentAccessTypeID)
	if err != nil {
		return err
	}
	responsibleUnitID, err := ds.userSvc.GetUserPrimaryUnitID(uint(documentFormDTO.ResponsibleUserID))
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
	}

	// Check if the document type has changed
	documentTypeChanged := existingDocument.DocumentTypeID != documentFormDTO.DocumentTypeID
	if documentTypeChanged {
		documentStatus, err := ds.documentStatusSvc.GetByType("draft")
		if err != nil {
			return err
		}
		if existingDocument.DocumentStatusID != documentStatus.ID {
			return utility.NewForbiddenError("Document type can only be changed in draft status")
		}
	}

	if documentFormDTO.Document.IsNewFile != nil && *documentFormDTO.Document.IsNewFile {
		documentStatus, err := ds.documentStatusSvc.GetByType("draft")
		if err != nil {
			return err
		}
		if existingDocument.DocumentStatusID != documentStatus.ID {
			return utility.NewForbiddenError("Document can only be changed in draft status")
		}
	}

	if err := ds.validateDocumentFormDTO(documentFormDTO); err != nil {
		return err
	}

	existingDocument.DocumentTypeID = documentFormDTO.DocumentTypeID
	existingDocument.TitleEn = documentFormDTO.TitleEn
	existingDocument.TitleAr = documentFormDTO.TitleAr
	existingDocument.DescriptionEn = documentFormDTO.DescriptionEn
	existingDocument.DescriptionAr = documentFormDTO.DescriptionAr
	existingDocument.DocumentCategoryID = documentFormDTO.DocumentCategoryID
	existingDocument.ResponsibleUserID = documentFormDTO.ResponsibleUserID
	existingDocument.ResponsibleUnitID = responsibleUnitID
	existingDocument.DocumentAccessTypeID = documentFormDTO.DocumentAccessTypeID
	existingDocument.UpdatedBYID = &currentUserID

	if err := ds.repo.Update(existingDocument, tx); err != nil {
		tx.Rollback()
		return err
	}

	documentFullDetails, err := ds.repo.GetFullDetailsByID(id, nil, tx)
	if err != nil {
		return err
	}

	var sharedUsers []*models.DocumentSharedUser
	var sharedUnits []*models.DocumentSharedUnit

	if existingDocument.ResponsibleUserID != documentFormDTO.ResponsibleUserID {
		sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
			DocumentID:  id,
			UserID:      documentFormDTO.ResponsibleUserID,
			Type:        "responsible-user",
			CreatedBYID: currentUserID,
		})

		err = ds.documentAccessSvc.DeleteResponsibleUser(id, &existingDocument.ResponsibleUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if existingDocumentAccessTypeID == documentFormDTO.DocumentAccessTypeID {
		fmt.Println("0")
		switch documentAccessType.Type {
		case "shared":
			// Fetch existing shared users
			existingSharedUsers, err := ds.documentAccessSvc.GetSharedUsers(id, 1, 10000, tx)
			if err != nil {
				tx.Rollback()
				log.Printf("Error fetching shared users: %v", err)
				return err
			}

			// Fetch existing shared units
			existingSharedUnits, err := ds.documentAccessSvc.GetSharedUnits(id, tx)
			if err != nil {
				tx.Rollback()
				log.Printf("Error fetching shared units: %v", err)
				return err
			}

			// Map existing shared users
			existingUserIDs := make(map[uint]bool)
			for _, user := range existingSharedUsers.Items {
				existingUserIDs[uint(user.ID)] = true
			}

			// Map existing shared units
			existingUnitIDs := make(map[uint]bool)
			for _, unit := range existingSharedUnits {
				existingUnitIDs[uint(unit.ID)] = true
			}

			// Add new shared users and build list to delete old users
			userIDsToDelete := []*uint{}
			for _, userID := range documentFormDTO.SharedUserIDs {
				if !existingUserIDs[*userID] {
					sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
						DocumentID:  id,
						UserID:      *userID,
						Type:        "shared-user",
						CreatedBYID: currentUserID,
					})
				} else {
					delete(existingUserIDs, *userID)
				}
			}
			for userID := range existingUserIDs {
				userIDsToDelete = append(userIDsToDelete, &userID)
			}

			// Delete old shared users
			if len(userIDsToDelete) > 0 {
				err = ds.documentAccessSvc.DeleteSharedUsers(id, userIDsToDelete, tx)
				if err != nil {
					tx.Rollback()
					log.Printf("Error deleting shared users: %v", err)
					return err
				}
			}

			// Add new shared units and build list to delete old units
			unitIDsToDelete := []*uint{}
			for _, unitID := range documentFormDTO.SharedUnitIDs {
				if !existingUnitIDs[*unitID] {
					sharedUnits = append(sharedUnits, &models.DocumentSharedUnit{
						DocumentID:         id,
						OrganizationUnitID: *unitID,
						CreatedBYID:        currentUserID,
					})
				} else {
					delete(existingUnitIDs, *unitID)
				}
			}
			for unitID := range existingUnitIDs {
				unitIDsToDelete = append(unitIDsToDelete, &unitID)
			}

			// Delete old shared units
			if len(unitIDsToDelete) > 0 {
				err = ds.documentAccessSvc.DeleteSharedUnits(id, unitIDsToDelete, tx)
				if err != nil {
					tx.Rollback()
					log.Printf("Error deleting shared units: %v", err)
					return err
				}
			}
		}
	} else {
		// Handle case when DocumentAccessTypeID changes
		fmt.Println("1")
		switch documentAccessType.Type {
		case "shared":
			fmt.Println("2")
			// Add new shared users and units
			for _, userID := range documentFormDTO.SharedUserIDs {
				sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
					DocumentID:  id,
					UserID:      *userID,
					Type:        "shared-user",
					CreatedBYID: currentUserID,
				})
			}
			for _, unitID := range documentFormDTO.SharedUnitIDs {
				sharedUnits = append(sharedUnits, &models.DocumentSharedUnit{
					DocumentID:         id,
					OrganizationUnitID: *unitID,
					CreatedBYID:        currentUserID,
				})
			}

		case "public":
			fmt.Println("3")
			// Delete all existing shared users
			existingSharedUsers, err := ds.documentAccessSvc.GetSharedUsers(id, 1, 10000, tx)
			if err != nil {
				tx.Rollback()
				log.Printf("Error fetching shared users for deletion: %v", err)
				return err
			}
			userIDsToDelete := []*uint{}
			for _, user := range existingSharedUsers.Items {
				userID := uint(user.ID)
				userIDsToDelete = append(userIDsToDelete, &userID)
			}
			if len(userIDsToDelete) > 0 {
				err = ds.documentAccessSvc.DeleteSharedUsers(id, userIDsToDelete, tx)
				if err != nil {
					tx.Rollback()
					log.Printf("Error deleting shared users: %v", err)
					return err
				}
			}

			// Delete all existing shared units
			existingSharedUnits, err := ds.documentAccessSvc.GetSharedUnits(id, tx)
			if err != nil {
				tx.Rollback()
				log.Printf("Error fetching shared units for deletion: %v", err)
				return err
			}
			unitIDsToDelete := []*uint{}
			for _, unit := range existingSharedUnits {
				unitID := uint(unit.ID)
				unitIDsToDelete = append(unitIDsToDelete, &unitID)
			}
			if len(unitIDsToDelete) > 0 {
				err = ds.documentAccessSvc.DeleteSharedUnits(id, unitIDsToDelete, tx)
				if err != nil {
					tx.Rollback()
					log.Printf("Error deleting shared units: %v", err)
					return err
				}
			}
		}
	}

	// Add new shared users and units
	if err := ds.documentAccessSvc.CreateUsers(sharedUsers, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ds.documentAccessSvc.CreateUnits(sharedUnits, tx); err != nil {
		tx.Rollback()
		return err
	}

	destinationDir := filepath.Join("document", "draft", strconv.FormatUint(uint64(id), 10), "documents")
	if documentFormDTO.Document.IsNewFile != nil && *documentFormDTO.Document.IsNewFile {
		existingDocumentFileName := documentFullDetails.LatestDocumentVersion.Name
		baseFileName := strings.TrimSuffix(existingDocumentFileName, filepath.Ext(existingDocumentFileName))

		url := utility.GenerateDocumentURL("document", "draft", strconv.FormatUint(uint64(id), 10), "documents")

		dataToEncrypt := []byte(documentFormDTO.Document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingDocVersion, err := ds.repo.GetDocumentVersionByID(documentFullDetails.LatestDocumentVersion.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		existingDocVersion.Version = documentFormDTO.DocumentVersion
		existingDocVersion.URL = url
		existingDocVersion.Name = documentFormDTO.Document.FileName
		existingDocVersion.Size = documentFormDTO.Document.Size
		existingDocVersion.Token = encryptedFilename
		existingDocVersion.OriginalName = &documentFormDTO.Document.OriginalName
		existingDocVersion.Type = documentFormDTO.Document.Type
		existingDocVersion.UpdatedBYID = &currentUserID

		err = ds.repo.UpdateDocumentVersion(existingDocVersion, tx)

		if err != nil {
			tx.Rollback()
			return err
		}

		if err := utility.RemoveFile(existingDocumentFileName, destinationDir); err != nil {
			tx.Rollback()
			return err
		}
		pdfFileName := baseFileName + ".pdf"
		if err := utility.RemoveFile(pdfFileName, destinationDir); err != nil {
			tx.Rollback()
			return err
		}

		if err := utility.StoreFile(documentFormDTO.Document.FileName, destinationDir); err != nil {
			tx.Rollback()
			return err
		}

		//add document history
		documentEditHistory := &documentDTOs.DocumentEditHistoryDTO{
			DocumentID:        id,
			DocumentVersionID: existingDocVersion.ID,
			ActionType:        "update",
			FileName:          &documentFormDTO.Document.FileName,
			OriginalName:      &documentFormDTO.Document.OriginalName,
			Type:              &documentFormDTO.Document.Type,
			Size:              &documentFormDTO.Document.Size,
		}

		err = ds.CreateDocumentEditHistory(documentEditHistory, &destinationDir, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	if existingDocument.DocumentTypeID != documentFormDTO.DocumentTypeID {

		var deleteDocumentWorkflowIds []*uint
		var deleteDocumentWorkflowParticipantIds []*uint
		var deleteDocumentWorkflowReviewIds []*uint

		if len(documentFullDetails.LatestDocumentVersion.DocumentWorkflows) > 0 {
			for _, docWorkflow := range documentFullDetails.LatestDocumentVersion.DocumentWorkflows {
				deleteDocumentWorkflowIds = append(deleteDocumentWorkflowIds, &docWorkflow.ID)

				if len(docWorkflow.DocumentWorkflowParticipants) > 0 {
					for _, participant := range docWorkflow.DocumentWorkflowParticipants {
						deleteDocumentWorkflowParticipantIds = append(deleteDocumentWorkflowParticipantIds, &participant.ID)
					}
				}

				if len(docWorkflow.DocumentWorkflowReviews) > 0 {
					for _, reviews := range docWorkflow.DocumentWorkflowReviews {
						deleteDocumentWorkflowReviewIds = append(deleteDocumentWorkflowReviewIds, &reviews.ID)
					}
				}
			}
		}

		if err := ds.documentWorkflowSvc.DeleteDocumentWorkflowParticipants(deleteDocumentWorkflowParticipantIds, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := ds.documentWorkflowSvc.DeleteDocumentWorkflowReviews(deleteDocumentWorkflowReviewIds, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := ds.documentWorkflowSvc.DeleteDocumentWorkflows(deleteDocumentWorkflowIds, tx); err != nil {
			tx.Rollback()
			return err
		}

		existingReviewUsers, err := ds.documentAccessSvc.GetReviewUsers(id, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		userIDsToDelete := make([]*uint, 0)
		for _, user := range existingReviewUsers {
			userID := uint(user.ID)
			userIDsToDelete = append(userIDsToDelete, &userID)
		}

		err = ds.documentAccessSvc.DeleteReviewUsers(id, userIDsToDelete, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = ds.documentWorkflowSvc.CreateDocumentWorkflow(id, documentFullDetails.LatestDocumentVersion.ID, documentFormDTO.DocumentTypeID, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	tx.Commit()

	go func() {
		if documentFormDTO.Document.IsNewFile != nil && *documentFormDTO.Document.IsNewFile {
			if err := utility.ConvertDocxToPDF(destinationDir, documentFormDTO.Document.FileName, destinationDir); err != nil {
				log.Printf("Failed to convert document to PDF: %v", err)
			}
		}
	}()

	return nil
}

func (ds *DocumentService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ds.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*models.DocumentVersion)
	documentEditHistoryToDelete := make(map[uint][]*models.DocumentEditHistory)

	for _, id := range ids {
		document, err := ds.repo.GetFullDetailsByID(id, nil, tx)
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
				deleteSharedUserIds = append(deleteSharedUserIds, &user.UserID)
			}
			if err := ds.documentAccessSvc.DeleteAllSharedUsers(id, deleteSharedUserIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(document.DocumentSharedUnits) > 0 {
			var deleteSharedUnitIds []*uint
			for _, unit := range document.DocumentSharedUnits {
				deleteSharedUnitIds = append(deleteSharedUnitIds, &unit.OrganizationUnitID)
			}
			if err := ds.documentAccessSvc.DeleteSharedUnits(id, deleteSharedUnitIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(document.DocumentAuditLogs) > 0 {
			if err := ds.documentAuditLogSvc.DeleteByDocumentID(id, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(document.DocumentVersions) > 0 {
			var deleteDocumentIds []*uint
			var deleteDocumentEditHistoryIds []*uint
			var deleteDocumentWorkflowIds []*uint
			var deleteDocumentWorkflowParticipantIds []*uint
			var deleteDocumentWorkflowReviewIds []*uint
			for _, doc := range document.DocumentVersions {
				documentsToDelete[id] = append(documentsToDelete[id], doc)
				deleteDocumentIds = append(deleteDocumentIds, &doc.ID)

				if len(doc.DocumentEditHistory) > 0 {
					for _, history := range doc.DocumentEditHistory {
						documentEditHistoryToDelete[id] = append(documentEditHistoryToDelete[id], history)
						deleteDocumentEditHistoryIds = append(deleteDocumentEditHistoryIds, &history.ID)
					}
				}

				if len(doc.DocumentWorkflows) > 0 {
					for _, docWorkflow := range doc.DocumentWorkflows {
						deleteDocumentWorkflowIds = append(deleteDocumentWorkflowIds, &docWorkflow.ID)

						if len(docWorkflow.DocumentWorkflowParticipants) > 0 {
							for _, participant := range docWorkflow.DocumentWorkflowParticipants {
								deleteDocumentWorkflowParticipantIds = append(deleteDocumentWorkflowParticipantIds, &participant.ID)
							}
						}

						if len(docWorkflow.DocumentWorkflowReviews) > 0 {
							for _, reviews := range docWorkflow.DocumentWorkflowReviews {
								deleteDocumentWorkflowReviewIds = append(deleteDocumentWorkflowReviewIds, &reviews.ID)
							}
						}
					}
				}
			}

			if err := ds.documentWorkflowSvc.DeleteDocumentWorkflowParticipants(deleteDocumentWorkflowParticipantIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

			if err := ds.documentWorkflowSvc.DeleteDocumentWorkflowReviews(deleteDocumentWorkflowReviewIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

			if err := ds.documentWorkflowSvc.DeleteDocumentWorkflows(deleteDocumentWorkflowIds, tx); err != nil {
				tx.Rollback()
				return nil, err
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
				if doc.Name != nil {
					storagePath := filepath.Join("document", "draft", strconv.FormatUint(uint64(documentID), 10), "document-history", strconv.FormatUint(uint64(doc.ID), 10), "documents")
					if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
						tx.Rollback()
						fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
					}
				}

			}
			destinationPath := filepath.Join("document", "draft", strconv.FormatUint(uint64(documentID), 10), "document-history")
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

func (ds *DocumentService) UpdateDocument(document *models.Document, tx ...utility.DatabaseTransaction) error {
	err := ds.repo.Update(document, tx...)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DocumentService) validateDocumentFormDTO(documentFormDTO *documentDTOs.DocumentFormDTO) error {
	if _, err := ds.documentCategorySvc.GetByID(documentFormDTO.DocumentCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document category does not exist"), "not_exist", "document category")
	}

	if _, err := ds.documentTypeSvc.GetByID(documentFormDTO.DocumentTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document type does not exist"), "not_exist", "document type")
	}

	if _, err := ds.workflowSvc.GetWorkflowByDocumentType(documentFormDTO.DocumentTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document type workflow does not exist"), "not_exist", "document type workflow")
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

	// accessType, err := ds.documentAccessTypeSvc.GetByID(documentFormDTO.DocumentAccessTypeID)
	// if err != nil {
	// 	return utility.NewValidationError(fmt.Errorf("failed to get document access type"), "not_exist", "document access type")
	// }

	// switch accessType.Type {
	// case "shared-users":
	// 	if documentFormDTO.SharedUserIDs == nil || len(documentFormDTO.SharedUserIDs) == 0 {
	// 		return utility.NewValidationError(fmt.Errorf("shared user IDs are required for shared-user access type"), "required", "shared_users")
	// 	}
	// case "shared-units":
	// 	if documentFormDTO.SharedUnitIDs == nil || len(documentFormDTO.SharedUnitIDs) == 0 {
	// 		return utility.NewValidationError(fmt.Errorf("shared unit IDs are required for shared-unit access type"), "required", "shared_units")
	// 	}
	// }
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
	document, err := ds.GetFullDetailsByID(id, nil, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if document.DocumentStatus.Type != "draft" && document.DocumentStatus.Type != "in-review" {
		tx.Rollback()
		return "", utility.NewForbiddenError("Checkout only allowed for documents in 'draft' or 'in-review' status.")
	}

	if document.IsLocked {
		tx.Rollback()
		return "", utility.NewForbiddenError("Document already checkouted.")
	}

	if document.CreatedByUser.ID != int(currentUserID) && document.ResponsibleUser.ID != int(currentUserID) {
		tx.Rollback()
		return "", utility.NewForbiddenError("Access denied: user is neither creator nor responsible user.")
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
	document, err := ds.GetFullDetailsByID(id, nil, tx)
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
	document, err := ds.GetFullDetailsByID(id, nil, tx)
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
	existingDocument, err := ds.GetFullDetailsByID(id, nil, tx...)
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

	baseFilename := strings.TrimSuffix(decryptedFilename, filepath.Ext(decryptedFilename))

	var storagePath string

	if existingDocument.DocumentStatus.Type == "published" {
		storagePath = filepath.Join("document", "published", strconv.FormatUint(uint64(id), 10), "documents")
		fileType = "pdf"
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

func (ds *DocumentService) SubmitDocument(id uint, currentUserID uint) error {
	tx := ds.repo.BeginTransaction()
	document, err := ds.GetFullDetailsByID(id, nil, tx)
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

	// Fetch the workflow level details, including participants (review users)
	docWorkflowLevel, err := ds.documentWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	urlParam := fmt.Sprintf(`{"id":%d}`, id)
	var taskIDs []uint

	for _, participant := range docWorkflowLevel.DocumentWorkflowParticipants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:   "document",
			DocumentID:         &id,
			DocumentWorkflowID: &docWorkflowLevel.ID,
			Title:              existingDocument.TitleEn,
			Description:        existingDocument.DescriptionEn,
			URLParam:           &urlParam,
			QueryParam:         nil,
			ResponsibleUserID:  participant.User.ID,
		}

		taskID, err := ds.taskSvc.Create(taskDTO, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		taskIDs = append(taskIDs, taskID)

	}

	for _, taskID := range taskIDs {
		relatedTaskIDs := []uint{}
		for _, otherTaskID := range taskIDs {
			if otherTaskID != taskID {
				relatedTaskIDs = append(relatedTaskIDs, otherTaskID)
			}
		}

		relatedTaskIDsJSON, err := json.Marshal(relatedTaskIDs)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = ds.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON), tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	//to do
	//activity log

	// Prepare notification
	go func() {
		urlParam = fmt.Sprintf(`{"id":%d}`, id)
		entity := "document"
		action := "review-user-assign"

		for _, participant := range docWorkflowLevel.DocumentWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               participant.User.ID,
				ModuleIdentifier:     "document",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      existingDocument.ReferenceNumber,
				HeadingEn:            "Document Review",
				HeadingAr:            "مراجعة الوثيقة",
				TitleEn:              existingDocument.TitleEn,
				TitleAr:              existingDocument.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := ds.userSvc.GetByID(participant.User.ID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := ds.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Document/%s: Review Assigned", existingDocument.TitleEn)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"DocumentTitle":         existingDocument.TitleEn,
				"DocumentType":          document.DocumentType.TitleEn,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/document/register/%d", appURL, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "document/templates/email/document_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = ds.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (ds *DocumentService) PublishDocument(id uint, currentUserID uint) error {
	tx := ds.repo.BeginTransaction()
	referenceNumber, err := ds.repo.GenerateReferenceNumber()
	if err != nil {
		return err
	}
	document, err := ds.GetFullDetailsByID(id, nil, tx)
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

	currentDateStr := time.Now().Format("2006-01-02")
	publishedAt, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	existingDocument.DocumentStatusID = documentStatus.ID
	existingDocument.ReferenceNumber = referenceNumber
	existingDocument.PublishedBYID = &currentUserID
	existingDocument.PublishedAt = publishedAt

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

	fileName := document.LatestDocumentVersion.Name
	baseFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) // Remove the extension to get the base name

	sourceDir := filepath.Join("document", "draft", strconv.FormatUint(uint64(document.ID), 10), "documents")
	destinationDir := filepath.Join("document", "published", strconv.FormatUint(uint64(document.ID), 10), "documents")

	//copy docx file
	err = utility.CopyFileFromPaths(sourceDir, destinationDir, fileName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error copying .docx file: %w", err)
	}

	//copy pdf file
	pdfFileName := baseFileName + ".pdf"
	err = utility.CopyFileFromPaths(sourceDir, destinationDir, pdfFileName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error copying .pdf file: %w", err)
	}

	tx.Commit()

	ds.notifyUsersOnPublish(document, id, currentUserID)

	//to do
	//activity log

	return nil
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
	document, err := ds.GetFullDetailsByID(documentID, nil, tx)
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

func (ds *DocumentService) GetDocumentLastEditSummary(documentID uint) (*documentDTOs.DocumentLastEditSummaryDTO, error) {
	document, err := ds.GetFullDetailsByID(documentID, nil)
	if err != nil {
		return nil, err
	}

	activity, err := ds.repo.GetDocumentLatestEditHistory(document.LatestDocumentVersion.ID)
	if err != nil {
		return nil, err
	}

	userDTO := mapUserToDTO(activity.CreatedByUser)

	latestActivityDTO := documentDTOs.ActivityDTO{
		Activity:  activity.ActionType,
		By:        userDTO,
		Comment:   activity.Comment,
		Timestamp: activity.CreatedAt,
	}
	// Initialize variables for the current document workflow
	var currentDocumentWorkflow *models.DocumentWorkflow
	var participants []*organizationDTOs.UserMinimalDTO
	totalComments := 0

	// Check if there is a next review level and fetch the corresponding workflow
	if document.LatestDocumentVersion.NextReviewLevel != nil {
		currentDocumentWorkflow, err = ds.documentWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, *document.LatestDocumentVersion.NextReviewLevel)
		if err != nil {
			return nil, err
		}
	}

	// Process workflow participants and count comments
	if currentDocumentWorkflow != nil {
		for _, participant := range currentDocumentWorkflow.DocumentWorkflowParticipants {
			userDTO := mapUserToDTO(participant.User)
			if userDTO != nil {
				participants = append(participants, userDTO)
			}
		}

		for _, review := range currentDocumentWorkflow.DocumentWorkflowReviews {
			if review.Comment != "" {
				totalComments++
			}
		}
	}

	// Prepare the response DTO
	responseDTO := &documentDTOs.DocumentLastEditSummaryDTO{
		LatestActivity: latestActivityDTO,
		CurrentLevel:   document.LatestDocumentVersion.NextReviewLevel,
		Participants:   participants,
		TotalComments:  &totalComments,
	}

	return responseDTO, nil
}

func (ds *DocumentService) ExportDocuments(data []*documentDTOs.DocumentListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := ds.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Type",
			"Title",
			"Responsible User",
			"Category",
			"Status",
		}
	} else {
		headers = []string{
			"SL No.",
			"Type",
			"Title",
			"Responsible User",
			"Category",
			"Status",
		}
	}
	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, document := range data {
		var category, responsibleUser, status, title, documentType string

		// Assign values based on language preference and check for nil pointers
		if userSetting.Language == "ar" {
			documentType = document.DocumentTypeTitleAr
			title = document.TitleAr
			responsibleUser = document.ResponsibleUserFirstNameAr + " " + document.ResponsibleUserLastNameAr
			category = document.DocumentCategoryTitleAr
			status = document.DocumentStatusTitleAr

		} else {
			documentType = document.DocumentTypeTitleEn
			title = document.TitleEn
			responsibleUser = document.ResponsibleUserFirstNameEn + " " + document.ResponsibleUserLastNameEn
			category = document.DocumentCategoryTitleEn
			status = document.DocumentStatusTitleEn
		}

		rows[i] = []interface{}{
			i + 1,
			documentType,
			title,
			responsibleUser,
			category,
			status,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (ds *DocumentService) GetDocumentVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersion, error) {
	documentVersion, err := ds.repo.GetDocumentVersionByID(id, tx...)
	if err != nil {
		return nil, err
	}

	return documentVersion, nil
}

func (ds *DocumentService) UpdateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) error {
	err := ds.repo.UpdateDocumentVersion(documentVersion, tx...)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DocumentService) notifyUsersOnPublish(document *documentDTOs.DocumentDetailsDTO, documentID uint, publishedBy uint) {
	go func() {
		entity := "document"
		action := "published"
		appURL := utility.GetAppURL()
		subject := fmt.Sprintf("Document/%s: Published Notification", document.TitleEn)
		redirectURL := fmt.Sprintf("%s/document/register/%d", appURL, documentID)

		publishedUser, _ := ds.userSvc.GetByID(publishedBy)

		// Collect users to notify based on access type
		var usersToNotify []*organizationDTOs.UserMinimalDTO

		if document.DocumentAccessType.Type == "public" {
			params := &organizationDTOs.UserListQueryParamsDTO{
				Sort:        "asc",
				SortBy:      "first_name_en",
				Keyword:     "",
				Status:      "",
				Page:        1,
				PageSize:    0,
				UserIDs:     nil,
				Filters:     nil,
				QuickFilter: nil,
				Export:      true,
			}

			userResponse, err := ds.userSvc.GetAll(params)
			if err != nil {
				utility.LogError("Failed to fetch all users for public document notification", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			for _, user := range userResponse.Items {
				minimalUser := &organizationDTOs.UserMinimalDTO{
					ID:          int(user.ID),
					FirstNameEN: user.FirstNameEn,
					LastNameEN:  user.LastNameEn,
					Email:       user.Email,
				}
				usersToNotify = append(usersToNotify, minimalUser)
			}
		} else if document.DocumentAccessType.Type == "shared" {
			const page = 1
			const pageSize = 10000
			sharedUsers, _ := ds.documentAccessSvc.GetSharedUsers(documentID, page, pageSize)
			usersToNotify = append(usersToNotify, sharedUsers.Items...)

			reviewUsers, _ := ds.documentAccessSvc.GetReviewUsers(documentID)
			usersToNotify = append(usersToNotify, reviewUsers...)

			primaryLevel, _ := ds.orgLevelSvc.GetPrimaryLevel()
			columnName := fmt.Sprintf("organization_unit_%d_id", primaryLevel.Level)

			sharedUnits, _ := ds.documentAccessSvc.GetSharedUnits(documentID)
			if len(sharedUnits) > 0 {

				unitIDs := make([]uint, len(sharedUnits))
				for i, unit := range sharedUnits {
					unitIDs[i] = uint(unit.ID)
				}

				// Convert unit IDs to string slice
				unitIDsStr := make([]string, len(unitIDs))
				for i, id := range unitIDs {
					unitIDsStr[i] = strconv.Itoa(int(id))
				}

				params := &organizationDTOs.UserListQueryParamsDTO{
					Sort:        "asc",
					SortBy:      "first_name_en",
					Keyword:     "",
					Status:      "",
					Page:        1,
					PageSize:    0,
					UserIDs:     nil,
					Filters:     map[string][]string{columnName: unitIDsStr},
					QuickFilter: nil,
					Export:      true,
				}

				unitUserResponse, err := ds.userSvc.GetAll(params)
				if err != nil {
					utility.LogError("Failed to fetch all users for public document notification", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				for _, user := range unitUserResponse.Items {
					minimalUser := &organizationDTOs.UserMinimalDTO{
						ID:          int(user.ID),
						FirstNameEN: user.FirstNameEn,
						LastNameEN:  user.LastNameEn,
						Email:       user.Email,
					}
					usersToNotify = append(usersToNotify, minimalUser)
				}
			}

		}

		// Send notifications to all users in the usersToNotify list
		uniqueUsers := make(map[uint]bool)
		for _, user := range usersToNotify {
			if _, exists := uniqueUsers[uint(user.ID)]; exists {
				continue
			}

			uniqueUsers[uint(user.ID)] = true

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               uint(user.ID),
				ModuleIdentifier:     "document",
				NotificationType:     string(commonDTOs.PublishType),
				HeadingEn:            "Document Published",
				HeadingAr:            "تم نشر الوثيقة",
				ReferenceNumber:      document.ReferenceNumber,
				TitleEn:              document.TitleEn,
				TitleAr:              document.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             fmt.Sprintf(`{"id":%d}`, documentID),
				IsSystemNotification: false,
				CreatedByID:          &publishedBy,
			}

			// Prepare email data
			emailTo := []string{user.Email}
			emailData := map[string]interface{}{
				"FirstNameEn":            user.FirstNameEN,
				"LastNameEn":             user.LastNameEN,
				"DocumentTitle":          document.TitleEn,
				"DocumentType":           document.DocumentType.TitleEn,
				"PublishedByFirstNameEn": publishedUser.FirstNameEn,
				"PublishedByFirstLastEn": publishedUser.LastNameEn,
				"PublishedAt":            time.Now().Format("2006-01-02"),
				"RedirectURL":            redirectURL,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "document/templates/email/document_publish_user_mail.html",
				EmailData:            emailData,
			}

			// Send both database and email notifications
			err := ds.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
			if err != nil {
				utility.LogError("Failed to send document publish notification", map[string]interface{}{
					"error": err.Error(),
				})
			}

		}
	}()
}
