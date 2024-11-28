package services

import (
	"errors"
	"fmt"
	"time"

	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditDocumentRequestServiceInterface interface {
	GetAll(auditID uint, params *auditDTOs.AuditDocumentRequestListQueryParamsDTO) (*auditDTOs.AuditDocumentRequestListResponseDTO, error)
	GetByID(id uint) (*models.AuditDocumentRequest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDocumentRequestDetailDTO, error)
	Create(auditID uint, currentUserID uint, auditDocumentRequestFormDTO *auditDTOs.AuditDocumentRequestFormDTO) ([]uint, error)
	Update(auditID uint, id uint, currentUserID uint, auditDocumentRequestFormDTO *auditDTOs.AuditDocumentRequestFormDTO) error
	Submit(auditID uint, id uint, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditDocumentRequestService struct {
	repo                          repositories.AuditDocumentRequestRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	auditSvc                      AuditServiceInterface
	auditDocumentRequestStatusSvc AuditDocumentRequestStatusServiceInterface
	auditableItemSvc              AuditableItemServiceInterface
	userSvc                       organizationServices.UserServiceInterface
	notificationSvc               commonServices.NotificationServiceInterface
}

func NewAuditDocumentRequestService(
	repo repositories.AuditDocumentRequestRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDocumentRequestStatusSvc AuditDocumentRequestStatusServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AuditDocumentRequestServiceInterface {
	return &AuditDocumentRequestService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		auditSvc:                      auditSvc,
		auditDocumentRequestStatusSvc: auditDocumentRequestStatusSvc,
		auditableItemSvc:              auditableItemSvc,
		userSvc:                       userSvc,
		notificationSvc:               notificationSvc,
	}
}

func (adrs *AuditDocumentRequestService) GetAll(auditID uint, params *auditDTOs.AuditDocumentRequestListQueryParamsDTO) (*auditDTOs.AuditDocumentRequestListResponseDTO, error) {
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	return adrs.repo.GetAll(auditID, params)
}

func (adrs *AuditDocumentRequestService) GetByID(id uint) (*models.AuditDocumentRequest, error) {
	documentRequest, err := adrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return documentRequest, nil
}

func (adrs *AuditDocumentRequestService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDocumentRequestDetailDTO, error) {
	var documentRequestDTO auditDTOs.AuditDocumentRequestDetailDTO

	documentRequest, err := adrs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(documentRequest, &documentRequestDTO)
	if err != nil {
		return nil, err
	}
	// Call the repository function to get the conversation count
	conversationCount, err := adrs.repo.CountConversationsByRequestID(id, tx...)
	if err != nil {
		return nil, err
	}
	documentRequestDTO.ConversationCount = int(conversationCount)

	if documentRequest.RequestedAt != nil {
		documentRequestDTO.RequestedAtString = documentRequest.RequestedAt.Format("2006-01-02")
	}
	documentRequestDTO.CreatedAt = documentRequest.CreatedAt
	if documentRequest.UpdatedAt != nil {
		documentRequestDTO.UpdatedAt = documentRequest.UpdatedAt
	}

	return &documentRequestDTO, nil
}

func (adrs *AuditDocumentRequestService) Create(auditID uint, currentUserID uint, formDTO *auditDTOs.AuditDocumentRequestFormDTO) ([]uint, error) {
	tx := adrs.repo.BeginTransaction()

	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := adrs.validateAuditDocumentRequestFormDTO(formDTO); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Group auditable items by ResponsibleUserID
	var responsibleUserGroups []*auditDTOs.ResponsibleUserGroupDTO
	for _, itemID := range formDTO.AuditableItemIds {
		auditableItem, err := adrs.auditableItemSvc.GetByID(*itemID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		responsibleUnitID := auditableItem.ResponsibleUnitID
		responsibleUserID := auditableItem.ResponsibleUserID
		// Find or create the group for this ResponsibleUnitID
		var group *auditDTOs.ResponsibleUserGroupDTO
		for _, existingGroup := range responsibleUserGroups {
			if existingGroup.ResponsibleUserID == responsibleUserID {
				group = existingGroup
				break
			}
		}
		if group == nil {
			group = &auditDTOs.ResponsibleUserGroupDTO{
				ResponsibleUnitID: responsibleUnitID,
				ResponsibleUserID: responsibleUserID,
				AuditableItems:    []*auditDTOs.AuditDocumentRequestAuditableItemDTO{},
			}
			responsibleUserGroups = append(responsibleUserGroups, group)
		}

		// Add the auditable item to the group
		group.AuditableItems = append(group.AuditableItems, &auditDTOs.AuditDocumentRequestAuditableItemDTO{
			AuditableItemID: *itemID,
			CreatedBYID:     currentUserID,
		})
	}
	auditDocumentRequestStatus, err := adrs.auditDocumentRequestStatusSvc.GetByType("draft")
	if err != nil {
		return nil, err
	}
	// Loop through the grouped data and create requests
	var requestIDs []uint
	for _, group := range responsibleUserGroups {
		// Create a new AuditDocumentRequest for each group
		documentRequest := &models.AuditDocumentRequest{
			AuditID:                      auditID,
			AuditDocumentRequestStatusID: uint(auditDocumentRequestStatus.ID),
			Title:                        formDTO.Title,
			Description:                  formDTO.Description,
			ResponsibleUnitID:            group.ResponsibleUnitID,
			ResponsibleUserID:            group.ResponsibleUserID,
			CreatedBYID:                  currentUserID,
		}

		// Save the AuditDocumentRequest
		id, err := adrs.repo.Create(documentRequest, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Convert DTOs to models for saving
		var auditableItems []*models.AuditDocumentRequestAuditableItem
		for _, itemDTO := range group.AuditableItems {
			auditableItems = append(auditableItems, &models.AuditDocumentRequestAuditableItem{
				AuditDocumentRequestID: id,
				AuditableItemID:        itemDTO.AuditableItemID,
				CreatedBYID:            itemDTO.CreatedBYID,
			})
		}

		// Save the AuditableItems
		if len(auditableItems) > 0 {
			if err := adrs.repo.CreateAuditableItems(auditableItems, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		requestIDs = append(requestIDs, id)

		// Log activity asynchronously
		go func(requestID uint, formDTO *auditDTOs.AuditDocumentRequestFormDTO) {
			err := adrs.activityLogSvc.LogActivity(
				models.CreateAction,
				models.AuditDocumentRequest{}.TableName(),
				requestID,
				formDTO,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.AuditDocumentRequest{}.TableName(),
					"recordID":  requestID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}
				utility.LogError("Failed to log activity", fields)
			}
		}(id, formDTO)
	}

	tx.Commit()
	return requestIDs, nil
}

func (adrs *AuditDocumentRequestService) Update(auditID uint, id uint, currentUserID uint, formDTO *auditDTOs.AuditDocumentRequestFormDTO) error {
	tx := adrs.repo.BeginTransaction()
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return err
	}
	documentRequest, err := adrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	documentRequest.Title = formDTO.Title
	documentRequest.Description = formDTO.Description
	documentRequest.UpdatedBYID = &currentUserID

	err = adrs.repo.Update(documentRequest, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := adrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditDocumentRequest{}.TableName(),
			id,
			formDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditDocumentRequest{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (adrs *AuditDocumentRequestService) Submit(auditID uint, id uint, currentUserID uint) error {
	tx := adrs.repo.BeginTransaction()
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return err
	}

	documentRequest, err := adrs.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}
	documentRequestFullDetails, err := adrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}
	if documentRequestFullDetails.AuditDocumentRequestStatus.Type != "draft" {
		tx.Rollback()
		return utility.NewForbiddenError("The document request has already been submitted.")
	}

	auditDocumentRequestStatus, err := adrs.auditDocumentRequestStatusSvc.GetByType("open")
	if err != nil {
		return err
	}
	now := time.Now()
	documentRequest.AuditDocumentRequestStatusID = uint(auditDocumentRequestStatus.ID)
	documentRequest.RequestedAt = &now
	documentRequest.RequestedBYID = &currentUserID

	err = adrs.repo.Update(documentRequest, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// todo -> send email notifications to responsible users

	tx.Commit()

	go func() {
		err := adrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditDocumentRequest{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditDocumentRequest{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}

		if documentRequest.ResponsibleUserID == 0 {
			utility.LogError("Responsible user is missing", map[string]interface{}{
				"documentRequestID": id,
			})
			return
		}

		responsibleUserID := documentRequest.ResponsibleUserID
		responsibleUser, err := adrs.userSvc.GetByID(responsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"userID": responsibleUserID,
				"error":  err.Error(),
			})
			return
		}

		requestedBy, err := adrs.userSvc.GetByID(*documentRequest.RequestedBYID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		urlParam := fmt.Sprintf(`{"audit_id":%d,"request_id":%d}`, auditID, id)
		entity := "audit-document-submitted"
		action := "review-document"
		appURL := utility.GetAppURL()

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserID:               responsibleUserID,
			ModuleIdentifier:     "audit",
			NotificationType:     string(commonDTOs.InformationType),
			ReferenceNumber:      documentRequest.Audit.AuditPlan.ReferenceNumber,
			HeadingEn:            "New Document Submitted",
			HeadingAr:            "تم تقديم مستند جديد",
			TitleEn:              documentRequest.Title,
			TitleAr:              documentRequest.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Document Submitted: %s", documentRequest.Title),
			TemplateRelativePath: "audit/templates/email/audit_document_submitted.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":     responsibleUser.FirstNameEn,
				"LastNameEn":      responsibleUser.LastNameEn,
				"AuditTitle":      documentRequest.Audit.AuditPlan.Title,
				"Title":           documentRequest.Title,
				"Description":     documentRequest.Description,
				"ResponsibleUnit": documentRequest.ResponsibleUnit.TitleEn,
				"RequestedBy":     fmt.Sprintf("%s %s", requestedBy.FirstNameEn, requestedBy.LastNameEn),
				"RequestedAt":     time.Now().Format("2006-01-02"),
				"RedirectURL":     fmt.Sprintf("%s/audit/audits", appURL),
			},
		}

		// Send Notifications
		if err := adrs.notificationSvc.SendNotification(dbNotification, emailNotification, "audit"); err != nil {
			utility.LogError("Failed to send document request notifications", map[string]interface{}{
				"documentRequestID": id,
				"userID":            responsibleUserID,
				"error":             err.Error(),
			})
		}
	}()

	return nil
}

func (adrs *AuditDocumentRequestService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := adrs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	for _, id := range ids {
		documentRequest, err := adrs.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}
		// Only Delete If the Audit Document Request Status is Draft
		if documentRequest.AuditDocumentRequestStatus.Type != "draft" {
			failedIDs = append(failedIDs, id)
			continue
		}
		if len(documentRequest.AuditableItems) > 0 {
			var deleteAuditableItemIDs []*uint
			for _, item := range documentRequest.AuditableItems {
				deleteAuditableItemIDs = append(deleteAuditableItemIDs, item.ID)
			}
			if err := adrs.repo.DeleteAuditableItems(deleteAuditableItemIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if err := adrs.repo.Delete(id, tx); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}
		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := adrs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditDocumentRequest{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)
				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditDocumentRequest{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}
					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIDs, nil
}

func (adrs *AuditDocumentRequestService) validateAuditDocumentRequestFormDTO(formDTO *auditDTOs.AuditDocumentRequestFormDTO) error {
	// Validate AuditableItemIds
	if formDTO.AuditableItemIds == nil {
		return utility.NewValidationError(fmt.Errorf("auditable item IDs are required"), "required", "auditable item IDs")
	}

	if len(formDTO.AuditableItemIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("at least one auditable item ID is required"), "required", "auditable item ID")
	}

	// Map to track seen AuditableItemIds
	seenItems := make(map[uint]bool)
	var uniqueAuditableItemIds []*uint

	for _, itemID := range formDTO.AuditableItemIds {
		if itemID != nil {
			// Check if the item is already seen
			if seenItems[*itemID] {
				continue // Skip duplicate
			}

			// Check if the auditable item exists
			if _, err := adrs.auditableItemSvc.GetByID(*itemID); err != nil {
				return utility.NewValidationError(fmt.Errorf("auditable item with ID %d does not exist", *itemID), "not_exist", "auditable item")
			}

			// Mark the item as seen and add it to the unique list
			seenItems[*itemID] = true
			uniqueAuditableItemIds = append(uniqueAuditableItemIds, itemID)
		}
	}

	// Update the formDTO with the unique list of AuditableItemIds
	formDTO.AuditableItemIds = uniqueAuditableItemIds

	return nil
}
