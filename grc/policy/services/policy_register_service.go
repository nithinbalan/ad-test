package services

import (
	"encoding/json"
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	documentServices "jethur/grc/documents/services"
	"jethur/grc/logs/services"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/repositories"
	settinsSvc "jethur/grc/settings/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"log"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PolicyRegisterServiceInterface interface {
	GetAll(params *policyDTOs.PolicyRegisterListQueryParamsDTO, currentUserID uint) (*policyDTOs.PolicyRegisterListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*policyDTOs.PolicyRegisterDetailsDTO, error)
	Create(currentUserID uint, documentFormDTO *policyDTOs.PolicyRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, documentFormDTO *policyDTOs.PolicyRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	SubmitPolicy(id uint, currentUserID uint) error

	ExportPolicies(data []*policyDTOs.PolicyRegisterListDTO, currentUserID uint) (*os.File, error)
}

type PolicyRegisterService struct {
	repo                      repositories.PolicyRegisterRepositoryInterface
	documentCategorySvc       documentServices.DocumentCategoryServiceInterface
	documentTypeSvc           documentServices.DocumentTypeServiceInterface
	documentAccessTypeSvc     documentServices.DocumentAccessTypeServiceInterface
	documentStatusSvc         documentServices.DocumentStatusServiceInterface
	userSvc                   organizationSvc.UserServiceInterface
	activityLogSvc            services.ActivityLogServiceInterface
	userSettingSvc            organizationSvc.UserSettingServiceInterface
	documentAccessSvc         documentServices.DocumentAccessServiceInterface
	taskSvc                   taskServices.TaskServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
	orgLevelSvc               settinsSvc.OrganizationLevelServiceInterface
	reviewFrequencySvc        masterSvc.ReviewFrequencyServiceInterface
	systemPreferenceSvc       settinsSvc.SystemPreferenceServiceInterface
	policyWorkflowSvc         documentServices.DocumentWorkflowServiceInterface
	documentWorkflowStatusSvc documentServices.DocumentWorkflowStatusServiceInterface
	workflowContextSvc        masterSvc.WorkflowContextServiceInterface
	workflowSvc               settinsSvc.WorkflowServiceInterface
}

func NewPolicyRegisterService(
	repo repositories.PolicyRegisterRepositoryInterface,
	documentCategorySvc documentServices.DocumentCategoryServiceInterface,
	documentTypeSvc documentServices.DocumentTypeServiceInterface,
	documentAccessTypeSvc documentServices.DocumentAccessTypeServiceInterface,
	documentStatusSvc documentServices.DocumentStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	userSettingSvc organizationSvc.UserSettingServiceInterface,
	documentAccessSvc documentServices.DocumentAccessServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	orgLevelSvc settinsSvc.OrganizationLevelServiceInterface,
	reviewFrequencySvc masterSvc.ReviewFrequencyServiceInterface,
	systemPreferenceSvc settinsSvc.SystemPreferenceServiceInterface,
	policyWorkflowSvc documentServices.DocumentWorkflowServiceInterface,
	documentWorkflowStatusSvc documentServices.DocumentWorkflowStatusServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
) PolicyRegisterServiceInterface {
	return &PolicyRegisterService{
		repo:                      repo,
		documentCategorySvc:       documentCategorySvc,
		documentTypeSvc:           documentTypeSvc,
		documentAccessTypeSvc:     documentAccessTypeSvc,
		documentStatusSvc:         documentStatusSvc,
		userSvc:                   userSvc,
		activityLogSvc:            activityLogSvc,
		userSettingSvc:            userSettingSvc,
		documentAccessSvc:         documentAccessSvc,
		taskSvc:                   taskSvc,
		notificationSvc:           notificationSvc,
		orgLevelSvc:               orgLevelSvc,
		reviewFrequencySvc:        reviewFrequencySvc,
		systemPreferenceSvc:       systemPreferenceSvc,
		policyWorkflowSvc:         policyWorkflowSvc,
		documentWorkflowStatusSvc: documentWorkflowStatusSvc,
		workflowContextSvc:        workflowContextSvc,
		workflowSvc:               workflowSvc,
	}
}

func (prs *PolicyRegisterService) GetAll(params *policyDTOs.PolicyRegisterListQueryParamsDTO, currentUserID uint) (*policyDTOs.PolicyRegisterListResponseDTO, error) {
	// Get user details to check for super admin role
	user, err := prs.userSvc.GetFullDetailsByID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Check if user has super admin role
	isUserSuperAdmin := false
	if user.Roles != nil {
		for _, role := range *user.Roles {
			if role.IsSuperAdmin {
				isUserSuperAdmin = true
				break
			}
		}
	}

	// Get primary organization unit field
	primaryLevel, err := prs.orgLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	// Determine the field name based on the primary level
	var primaryOrgUnitField string
	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return nil, fmt.Errorf("unsupported primary level")
	}

	return prs.repo.GetAll(params, currentUserID, isUserSuperAdmin, primaryOrgUnitField)
}

func (prs PolicyRegisterService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error) {
	document, err := prs.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return document, nil
}

func (prs *PolicyRegisterService) GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*policyDTOs.PolicyRegisterDetailsDTO, error) {
	var documentDTO policyDTOs.PolicyRegisterDetailsDTO

	document, err := prs.repo.GetFullDetailsByID(id, currentUserID, tx...)
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

	documentDTO.CreatedAt = document.CreatedAt
	documentDTO.UpdatedAt = document.UpdatedAt

	return &documentDTO, nil
}

func (prs *PolicyRegisterService) Create(currentUserID uint, documentFormDTO *policyDTOs.PolicyRegisterFormDTO) (uint, error) {
	tx := prs.repo.BeginTransaction()

	if err := prs.validateDocumentFormDTO(documentFormDTO); err != nil {
		return 0, err
	}

	responsibleUnitID, err := prs.userSvc.GetUserPrimaryUnitID(uint(documentFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
	}

	documentStatus, err := prs.documentStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	document := &models.Document{
		TitleEn:              documentFormDTO.TitleEn,
		TitleAr:              documentFormDTO.TitleAr,
		DescriptionEn:        documentFormDTO.DescriptionEn,
		DescriptionAr:        documentFormDTO.DescriptionAr,
		DocumentTypeID:       &documentFormDTO.DocumentTypeID,
		DocumentCategoryID:   &documentFormDTO.DocumentCategoryID,
		ResponsibleUserID:    &documentFormDTO.ResponsibleUserID,
		ResponsibleUnitID:    &responsibleUnitID,
		ReviewFrequencyID:    documentFormDTO.ReviewFrequencyID,
		DocumentAccessTypeID: documentFormDTO.DocumentAccessTypeID,
		DocumentStatusID:     &documentStatus.ID,
		PolicyFolderID:       documentFormDTO.PolicyFolderID,
		ParentDocumentID:     documentFormDTO.ParentDocumentID,
		CreatedBYID:          currentUserID,
	}

	id, err := prs.repo.Create(document, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	documentAccessType, err := prs.documentAccessTypeSvc.GetByID(documentFormDTO.DocumentAccessTypeID)
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
				DocumentID:  id,
				OrgUnitID:   *unitID,
				CreatedBYID: currentUserID,
			}
			sharedUnits = append(sharedUnits, sharedUnit)
		}
		if err := prs.documentAccessSvc.CreateUnits(sharedUnits, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err := prs.documentAccessSvc.CreateUsers(sharedUsers, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	var docVersionId uint

	documentVersion := &models.DocumentVersion{
		DocumentID:      id,
		Version:         documentFormDTO.DocumentVersion,
		IsLatestVersion: true,
		CreatedBYID:     currentUserID,
	}

	docVersionId, err = prs.repo.CreateDocumentVersion(documentVersion, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = prs.CreatePolicyWorkflow(id, docVersionId, documentFormDTO.DocumentTypeID, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingDocument, err := prs.GetFullDetailsByID(id, &currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch document", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = prs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Document{}.TableName(),
			id,
			existingDocument,
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

func (prs *PolicyRegisterService) Update(id uint, currentUserID uint, documentFormDTO *policyDTOs.PolicyRegisterFormDTO) error {
	// Fetch existing document details
	existingDocument, err := prs.repo.GetByID(id)
	if err != nil {
		return err
	}
	existingDocumentAccessTypeID := existingDocument.DocumentAccessTypeID
	documentAccessType, err := prs.documentAccessTypeSvc.GetByID(documentFormDTO.DocumentAccessTypeID)
	if err != nil {
		return err
	}
	responsibleUnitID, err := prs.userSvc.GetUserPrimaryUnitID(uint(documentFormDTO.ResponsibleUserID))
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
	}

	tx := prs.repo.BeginTransaction()
	oldDocument, err := prs.repo.GetFullDetailsByID(id, nil, tx)
	if err != nil {
		return err
	}
	oldDocVersion, err := prs.repo.GetDocumentVersionByID(oldDocument.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.validateDocumentFormDTO(documentFormDTO); err != nil {
		return err
	}

	existingDocument.TitleEn = documentFormDTO.TitleEn
	existingDocument.TitleAr = documentFormDTO.TitleAr
	existingDocument.DescriptionEn = documentFormDTO.DescriptionEn
	existingDocument.DescriptionAr = documentFormDTO.DescriptionAr
	existingDocument.DocumentTypeID = &documentFormDTO.DocumentTypeID
	existingDocument.DocumentCategoryID = &documentFormDTO.DocumentCategoryID
	existingDocument.ResponsibleUserID = &documentFormDTO.ResponsibleUserID
	existingDocument.ResponsibleUnitID = &responsibleUnitID
	existingDocument.ReviewFrequencyID = documentFormDTO.ReviewFrequencyID
	existingDocument.DocumentAccessTypeID = documentFormDTO.DocumentAccessTypeID
	existingDocument.UpdatedBYID = &currentUserID

	if err := prs.repo.Update(existingDocument, tx); err != nil {
		tx.Rollback()
		return err
	}

	documentFullDetails, err := prs.repo.GetFullDetailsByID(id, nil, tx)
	if err != nil {
		return err
	}

	existingDocVersion, err := prs.repo.GetDocumentVersionByID(documentFullDetails.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingDocVersion.Version = documentFormDTO.DocumentVersion

	err = prs.repo.UpdateDocumentVersion(existingDocVersion, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var sharedUsers []*models.DocumentSharedUser
	var sharedUnits []*models.DocumentSharedUnit

	if existingDocument.ResponsibleUserID != &documentFormDTO.ResponsibleUserID {
		sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
			DocumentID:  id,
			UserID:      documentFormDTO.ResponsibleUserID,
			Type:        "responsible-user",
			CreatedBYID: currentUserID,
		})

		err = prs.documentAccessSvc.DeleteResponsibleUser(id, *existingDocument.ResponsibleUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if existingDocumentAccessTypeID == documentFormDTO.DocumentAccessTypeID {
		switch documentAccessType.Type {
		case "shared":
			// Fetch existing shared users
			existingSharedUsers, err := prs.documentAccessSvc.GetSharedUsers(id, 1, 10000, tx)
			if err != nil {
				tx.Rollback()
				log.Printf("Error fetching shared users: %v", err)
				return err
			}

			// Fetch existing shared units
			existingSharedUnits, err := prs.documentAccessSvc.GetSharedUnits(id, tx)
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
			userIDsToDelete := []uint{}
			// Create a map of new shared user IDs for quick lookup
			newUserIDs := make(map[uint]bool)
			for _, userID := range documentFormDTO.SharedUserIDs {
				newUserIDs[*userID] = true
			}

			// Check each existing user to determine if they should be deleted
			for userID := range existingUserIDs {
				if !newUserIDs[userID] {
					// User is in existing list but not in the new list, so mark for deletion
					// userIDCopy := userID // Create a copy of userID to avoid pointer reuse issues
					userIDsToDelete = append(userIDsToDelete, userID)
				} else {
					// User is still valid, remove from newUserIDs to avoid re-adding
					delete(newUserIDs, userID)
				}
			}

			// Add new users (remaining in newUserIDs) to the sharedUsers slice
			for userID := range newUserIDs {
				sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
					DocumentID:  id,
					UserID:      userID,
					Type:        "shared-user",
					CreatedBYID: currentUserID,
				})
			}

			// Delete old shared users
			if len(userIDsToDelete) > 0 {
				err = prs.documentAccessSvc.DeleteSharedUsers(id, userIDsToDelete, tx)
				if err != nil {
					tx.Rollback()
					log.Printf("Error deleting shared users: %v", err)
					return err
				}
			}

			// Add new shared units and build list to delete old units
			unitIDsToDelete := []uint{}
			for _, unitID := range documentFormDTO.SharedUnitIDs {
				if !existingUnitIDs[*unitID] {
					sharedUnits = append(sharedUnits, &models.DocumentSharedUnit{
						DocumentID:  id,
						OrgUnitID:   *unitID,
						CreatedBYID: currentUserID,
					})
				} else {
					delete(existingUnitIDs, *unitID)
				}
			}
			for unitID := range existingUnitIDs {
				unitIDsToDelete = append(unitIDsToDelete, unitID)
			}

			// Delete old shared units
			if len(unitIDsToDelete) > 0 {
				err = prs.documentAccessSvc.DeleteSharedUnits(id, unitIDsToDelete, tx)
				if err != nil {
					tx.Rollback()
					log.Printf("Error deleting shared units: %v", err)
					return err
				}
			}
		}
	} else {
		// Handle case when DocumentAccessTypeID changes
		switch documentAccessType.Type {
		case "shared":
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
					DocumentID:  id,
					OrgUnitID:   *unitID,
					CreatedBYID: currentUserID,
				})
			}

		case "public":
			// Delete all existing shared users
			existingSharedUsers, err := prs.documentAccessSvc.GetSharedUsers(id, 1, 10000, tx)
			if err != nil {
				tx.Rollback()
				log.Printf("Error fetching shared users for deletion: %v", err)
				return err
			}
			userIDsToDelete := []uint{}
			for _, user := range existingSharedUsers.Items {
				userID := uint(user.ID)
				userIDsToDelete = append(userIDsToDelete, userID)
			}
			if len(userIDsToDelete) > 0 {
				err = prs.documentAccessSvc.DeleteSharedUsers(id, userIDsToDelete, tx)
				if err != nil {
					tx.Rollback()
					log.Printf("Error deleting shared users: %v", err)
					return err
				}
			}

			// Delete all existing shared units
			existingSharedUnits, err := prs.documentAccessSvc.GetSharedUnits(id, tx)
			if err != nil {
				tx.Rollback()
				log.Printf("Error fetching shared units for deletion: %v", err)
				return err
			}
			unitIDsToDelete := []uint{}
			for _, unit := range existingSharedUnits {
				unitID := uint(unit.ID)
				unitIDsToDelete = append(unitIDsToDelete, unitID)
			}
			if len(unitIDsToDelete) > 0 {
				err = prs.documentAccessSvc.DeleteSharedUnits(id, unitIDsToDelete, tx)
				if err != nil {
					tx.Rollback()
					log.Printf("Error deleting shared units: %v", err)
					return err
				}
			}
		}
	}

	// Add new shared users and units
	if err := prs.documentAccessSvc.CreateUsers(sharedUsers, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.documentAccessSvc.CreateUnits(sharedUnits, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingDoc, err := prs.repo.GetFullDetailsByID(id, nil)
		if err != nil {
			return
		}

		existingDocumentVersion, err := prs.repo.GetDocumentVersionByID(existingDoc.LatestDocumentVersion.ID)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"document_type":        oldDocument.DocumentType.TitleEn,
			"title_en":             oldDocument.TitleEn,
			"title_ar":             oldDocument.TitleAr,
			"description_en":       oldDocument.DescriptionEn,
			"description_ar":       oldDocument.DescriptionAr,
			"document_category":    oldDocument.DocumentCategory.TitleEn,
			"responsible_user":     fmt.Sprintf("%s %s", oldDocument.ResponsibleUser.FirstNameEn, oldDocument.ResponsibleUser.LastNameEn),
			"responsible_unit":     oldDocument.ResponsibleUnit.TitleEn,
			"document_access_type": oldDocument.DocumentAccessType.TitleEn,
			"version":              oldDocVersion,
		}

		newData := map[string]interface{}{
			"document_type":        existingDoc.DocumentType.TitleEn,
			"title_en":             existingDoc.TitleEn,
			"title_ar":             existingDoc.TitleAr,
			"description_en":       existingDoc.DescriptionEn,
			"description_ar":       existingDoc.DescriptionAr,
			"document_category":    existingDoc.DocumentCategory.TitleEn,
			"responsible_user":     fmt.Sprintf("%s %s", existingDoc.ResponsibleUser.FirstNameEn, existingDoc.ResponsibleUser.LastNameEn),
			"responsible_unit":     existingDoc.ResponsibleUnit.TitleEn,
			"document_access_type": existingDoc.DocumentAccessType.TitleEn,
			"version":              existingDocumentVersion,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := prs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.Document{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.Document{}.TableName(),
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

func (prs *PolicyRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := prs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedDocument []map[string]interface{}
	documentsToDelete := make(map[uint][]*models.DocumentVersion)

	for _, id := range ids {
		document, err := prs.repo.GetFullDetailsByID(id, nil, tx)
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
			var deleteSharedUserIds []uint
			for _, user := range document.DocumentSharedUsers {
				deleteSharedUserIds = append(deleteSharedUserIds, user.UserID)
			}
			if err := prs.documentAccessSvc.DeleteAllSharedUsers(id, deleteSharedUserIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if len(document.DocumentSharedUnits) > 0 {
			var deleteSharedUnitIds []uint
			for _, unit := range document.DocumentSharedUnits {
				deleteSharedUnitIds = append(deleteSharedUnitIds, unit.OrgUnitID)
			}
			if err := prs.documentAccessSvc.DeleteSharedUnits(id, deleteSharedUnitIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if len(document.DocumentVersions) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range document.DocumentVersions {
				documentsToDelete[id] = append(documentsToDelete[id], doc)
				deleteDocumentIds = append(deleteDocumentIds, &doc.ID)

			}

			if err := prs.repo.DeleteDocumentVersions(deleteDocumentIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}

		}

		if err := prs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedDocument = append(deletedDocument, map[string]interface{}{
			"id":       id,
			"title_en": document.TitleEn,
			"title_ar": document.TitleAr,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		tx.Commit()
	}

	if len(successfulIds) > 0 {
		go func() {
			err := prs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Document{}.TableName(),
				successfulIds,
				deletedDocument,
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

func (prs *PolicyRegisterService) validateDocumentFormDTO(documentFormDTO *policyDTOs.PolicyRegisterFormDTO) error {
	if _, err := prs.documentCategorySvc.GetByID(documentFormDTO.DocumentCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document category does not exist"), "not_exist", "document category")
	}

	if _, err := prs.documentAccessTypeSvc.GetByID(documentFormDTO.DocumentAccessTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document access type does not exist"), "not_exist", "document access type")
	}
	return nil
}

func (prs *PolicyRegisterService) ExportPolicies(data []*policyDTOs.PolicyRegisterListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := prs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
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

func (prs *PolicyRegisterService) SubmitPolicy(id uint, currentUserID uint) error {
	tx := prs.repo.BeginTransaction()
	document, err := prs.GetFullDetailsByID(id, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentStatus, err := prs.documentStatusSvc.GetByType("in-review")
	if err != nil {
		tx.Rollback()
		return err
	}

	existingDocument, err := prs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingDocument.DocumentStatusID = &documentStatus.ID

	err = prs.repo.Update(existingDocument, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentVersion, err := prs.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var level int
	if document.DocumentStatus.Type == "reverted" {
		// If status is "reverted", use the NextReviewLevel for resubmission
		if document.LatestDocumentVersion.NextReviewLevel != nil {
			level = *document.LatestDocumentVersion.NextReviewLevel
		} else {
			tx.Rollback()
			return fmt.Errorf("next review level is not set for reverted document")
		}
	} else {
		// For non-reverted documents, start from level 1
		level = 1
	}

	currentTime := time.Now()

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	if document.DocumentStatus.Type != "reverted" {
		documentVersion.SubmittedBYID = &currentUserID
		documentVersion.SubmittedAt = &currentTime
		documentVersion.NextReviewLevel = &level

		err = prs.repo.UpdateDocumentVersion(documentVersion, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Fetch the workflow level details, including participants (review users)
	docWorkflowLevel, err := prs.policyWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if document.DocumentStatus.Type == "reverted" {
		documentWorkflowStatus, err := prs.documentWorkflowStatusSvc.GetByType("resubmit")
		if err != nil {
			return err
		}
		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       docWorkflowLevel.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  "Resubmitted",
			CreatedBYID:              currentUserID,
		}

		err = prs.repo.CreatePolicyWorkflowReview(documentWorkflowReview, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	urlParam := fmt.Sprintf(`{"id":%d}`, id)
	var taskIDs []uint

	for _, participant := range docWorkflowLevel.DocumentWorkflowParticipants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:   "document_register",
			DocumentID:         &id,
			DocumentWorkflowID: &docWorkflowLevel.ID,
			Title:              existingDocument.TitleEn,
			Description:        existingDocument.DescriptionEn,
			URLParam:           &urlParam,
			QueryParam:         nil,
			ResponsibleUserID:  participant.User.ID,
		}

		taskID, err := prs.taskSvc.Create(taskDTO, currentUserID)
		if err != nil {
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
			return err
		}

		err = prs.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON))
		if err != nil {
			return err
		}
	}

	// to do
	// activity log

	// Prepare notification
	go func() {
		urlParam = fmt.Sprintf(`{"id":%d}`, id)
		entity := "policy"
		action := "review-user-assign"

		for _, participant := range docWorkflowLevel.DocumentWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.User.ID},
				ModuleIdentifier:     "policy_management",
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

			responsibleUser, err := prs.userSvc.GetByID(participant.User.ID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := prs.userSvc.GetByID(currentUserID)
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
				"DocumentCategory":      document.DocumentCategory.TitleEn,
				"DocumentVersion":       document.LatestDocumentVersion.Version,
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
			err = prs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (prs *PolicyRegisterService) CreatePolicyWorkflow(documentId uint, documentVersionId uint, documentTypeId uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := prs.workflowContextSvc.GetByType("policy_register")
	if err != nil {
		return err
	}

	workflow, err := prs.workflowSvc.GetWorkflowByDocumentType(documentTypeId, workflowContext.ID, tx...)
	if err != nil {
		return err
	}

	workflowLevels, err := prs.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	documentWorkflowStatus, err := prs.documentWorkflowStatusSvc.GetByType("pending")

	if err != nil {
		return err
	}

	var sharedUsers []*models.DocumentSharedUser
	uniqueUsers := make(map[uint]bool) // Map to track added users

	for _, level := range workflowLevels.WorkflowLevels {
		documentWorkflow := &models.DocumentWorkflow{
			DocumentVersionID:        documentVersionId,
			Level:                    level.Level,
			WorkflowLevelTypeID:      level.WorkflowLevelType.ID,
			ApprovalTypes:            level.ApprovalType,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			CreatedBYID:              currentUserID,
		}

		docWorkflowId, err := prs.repo.CreateDocumentWorkflow(documentWorkflow, tx...)
		if err != nil {
			return err
		}

		// Handle participants for each level
		for _, participant := range level.Participants {
			if participant.Role != nil {

				userIds, err := prs.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					err = prs.policyWorkflowSvc.CreateWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}

					if _, exists := uniqueUsers[userId]; exists {
						continue
					}

					uniqueUsers[userId] = true

					// Append to sharedUsers slice
					sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
						DocumentID:  documentId,
						UserID:      userId,
						Type:        "review-user",
						CreatedBYID: currentUserID,
					})

				}
			} else if participant.User != nil {
				userId := uint(participant.User.ID)
				err = prs.policyWorkflowSvc.CreateWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}

				if _, exists := uniqueUsers[userId]; exists {
					continue
				}

				uniqueUsers[userId] = true

				// Append to sharedUsers slice
				sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
					DocumentID:  documentId,
					UserID:      userId,
					Type:        "review-user",
					CreatedBYID: currentUserID,
				})
			}
		}
	}

	if err := prs.documentAccessSvc.CreateUsers(sharedUsers, tx...); err != nil {
		return err
	}

	return nil
}
