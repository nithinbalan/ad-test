package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditScheduleServiceInterface interface {
	GetAll(auditID uint, params *auditDTOs.AuditScheduleListQueryParamsDTO) (*auditDTOs.AuditScheduleListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditScheduleDetailDTO, error)
	GetByID(id uint) (*models.AuditSchedule, error)
	Create(auditID uint, currentUserID uint, auditScheduleFormDTO *auditDTOs.AuditScheduleFormDTO) (uint, error)
	Update(id uint, auditID uint, currentUserID uint, auditScheduleFormDTO *auditDTOs.AuditScheduleFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetAllSchedulesByAuditID(auditID uint, tx ...utility.DatabaseTransaction) ([]*models.AuditSchedule, error)
	GetScheduleSummary(auditID uint) (*auditDTOs.AuditScheduleSummaryDTO, error)

	StatusUpdate(auditScheduleID uint, auditScheduleStatusUpdateFormDTO *auditDTOs.AuditScheduleStatusUpdateFormDTO, currentUserID uint, tx utility.DatabaseTransaction) error
}

type AuditScheduleService struct {
	repo                   repositories.AuditScheduleRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	auditSvc               AuditServiceInterface
	auditActivityTypeSvc   AuditActivityTypeServiceInterface
	userSvc                organizationServices.UserServiceInterface
	auditScheduleStatusSvc AuditScheduleStatusServiceInterface
	auditableItemSvc       AuditableItemServiceInterface
}

func NewAuditScheduleService(
	repo repositories.AuditScheduleRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditActivityTypeSvc AuditActivityTypeServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditScheduleStatusSvc AuditScheduleStatusServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
) AuditScheduleServiceInterface {
	return &AuditScheduleService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		auditSvc:               auditSvc,
		auditActivityTypeSvc:   auditActivityTypeSvc,
		userSvc:                userSvc,
		auditScheduleStatusSvc: auditScheduleStatusSvc,
		auditableItemSvc:       auditableItemSvc,
	}
}

func (ass *AuditScheduleService) GetAll(auditID uint, params *auditDTOs.AuditScheduleListQueryParamsDTO) (*auditDTOs.AuditScheduleListResponseDTO, error) {
	if _, err := ass.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	return ass.repo.GetAll(auditID, params)
}

func (ass *AuditScheduleService) GetByID(id uint) (*models.AuditSchedule, error) {
	auditSchedule, err := ass.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditSchedule, nil
}

func (ass *AuditScheduleService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditScheduleDetailDTO, error) {
	var auditScheduleDTO auditDTOs.AuditScheduleDetailDTO

	auditSchedule, err := ass.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditSchedule, &auditScheduleDTO)
	if err != nil {
		return nil, err
	}
	if auditSchedule.StartDate != nil {
		auditScheduleDTO.StartDateString = auditSchedule.StartDate.Format("2006-01-02")
	}
	if auditSchedule.TargetDate != nil {
		auditScheduleDTO.TargetDateString = auditSchedule.TargetDate.Format("2006-01-02")
	}
	auditScheduleDTO.CreatedAt = auditSchedule.CreatedAt
	auditScheduleDTO.UpdatedAt = auditSchedule.UpdatedAt

	return &auditScheduleDTO, nil
}

func (ass *AuditScheduleService) Create(auditID uint, currentUserID uint, auditScheduleFormDTO *auditDTOs.AuditScheduleFormDTO) (uint, error) {
	tx := ass.repo.BeginTransaction()
	if _, err := ass.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}
	dates, err := ass.parseAuditScheduleDates(auditScheduleFormDTO)
	if err != nil {
		return 0, err
	}
	// auditScheduleStatus, err := ass.auditScheduleStatusSvc.GetByType("scheduled")
	// if err != nil {
	// 	return 0, err
	// }

	if err := ass.validateAuditScheduleFormDTO(auditScheduleFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	audit, err := ass.auditSvc.GetByID(auditID)
	if err != nil {
		return 0, err
	}

	auditSchedule := &models.AuditSchedule{
		AuditID:               auditID,
		AuditPlanID:           audit.AuditPlanID,
		AuditActivityTypeID:   auditScheduleFormDTO.AuditActivityTypeID,
		Title:                 auditScheduleFormDTO.Title,
		Description:           auditScheduleFormDTO.Description,
		StartDate:             dates.StartDate,
		TargetDate:            dates.TargetDate,
		ResponsibleUserID:     auditScheduleFormDTO.ResponsibleUserID,
		AuditScheduleStatusID: auditScheduleFormDTO.AuditScheduleStatusID,
		CreatedBYID:           currentUserID,
	}

	id, err := ass.repo.Create(auditSchedule, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	scheduleStatusID := auditScheduleFormDTO.AuditScheduleStatusID
	auditScheduleStatusUpdate := &models.AuditScheduleStatusUpdate{
		AuditScheduleID:       id,
		AuditScheduleStatusID: scheduleStatusID,
		Comment:               "Initial audit schedule status",
		CreatedBYID:           currentUserID,
	}

	_, err = ass.repo.CreateInitialAuditScheduleStatusUpdate(auditScheduleStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var auditableItems []*models.AuditScheduleAuditableItem
	for _, itemID := range auditScheduleFormDTO.AuditableItemIds {
		auditableItem := &models.AuditScheduleAuditableItem{
			AuditScheduleID: id,
			AuditableItemID: *itemID,
			CreatedBYID:     currentUserID,
		}
		auditableItems = append(auditableItems, auditableItem)
	}

	if len(auditableItems) > 0 {
		if err := ass.repo.CreateAuditableItems(auditableItems, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingAuditSchedule, err := ass.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch audit schedule", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ass.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditSchedule{}.TableName(),
			id,
			existingAuditSchedule,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditSchedule{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ass *AuditScheduleService) Update(id uint, auditID uint, currentUserID uint, auditScheduleFormDTO *auditDTOs.AuditScheduleFormDTO) error {
	tx := ass.repo.BeginTransaction()
	if _, err := ass.auditSvc.GetByID(auditID); err != nil {
		return err
	}

	oldAuditSchedule, err := ass.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	dates, err := ass.parseAuditScheduleDates(auditScheduleFormDTO)
	if err != nil {
		return err
	}

	if err := ass.validateAuditScheduleFormDTO(auditScheduleFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditSchedule, err := ass.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingAuditSchedule.AuditActivityTypeID = auditScheduleFormDTO.AuditActivityTypeID
	existingAuditSchedule.Title = auditScheduleFormDTO.Title
	existingAuditSchedule.Description = auditScheduleFormDTO.Description
	existingAuditSchedule.StartDate = dates.StartDate
	existingAuditSchedule.TargetDate = dates.TargetDate
	existingAuditSchedule.ResponsibleUserID = auditScheduleFormDTO.ResponsibleUserID
	existingAuditSchedule.AuditScheduleStatusID = auditScheduleFormDTO.AuditScheduleStatusID
	existingAuditSchedule.UpdatedBYID = &currentUserID

	err = ass.repo.Update(existingAuditSchedule, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingSchedule, err := ass.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create a map of existing auditable item IDs
	existingAuditableItemIDs := make(map[uint]bool)
	for _, item := range existingSchedule.AuditableItems {
		existingAuditableItemIDs[item.AuditableItemID] = true
	}

	// Create a map of new auditable item IDs
	newAuditableItemIDs := make(map[uint]bool)
	for _, itemID := range auditScheduleFormDTO.AuditableItemIds {
		newAuditableItemIDs[*itemID] = true
	}

	// Add new auditable items
	var auditableItemsToAdd []*models.AuditScheduleAuditableItem
	for _, newItemID := range auditScheduleFormDTO.AuditableItemIds {
		if _, exists := existingAuditableItemIDs[*newItemID]; !exists {
			// Create a new auditable item record
			auditableItem := &models.AuditScheduleAuditableItem{
				AuditScheduleID: id,
				AuditableItemID: *newItemID,
				CreatedBYID:     currentUserID,
			}
			auditableItemsToAdd = append(auditableItemsToAdd, auditableItem)
		}
	}

	// Save the new auditable items
	if len(auditableItemsToAdd) > 0 {
		if err := ass.repo.CreateAuditableItems(auditableItemsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove auditable items not in the new list
	var auditableItemsToDelete []*uint
	for _, existingItem := range existingSchedule.AuditableItems {
		if _, exists := newAuditableItemIDs[existingItem.AuditableItemID]; !exists {
			// Append item to delete list
			auditableItemsToDelete = append(auditableItemsToDelete, &existingItem.ID)
		}
	}

	// Delete auditable items not in the new list
	if len(auditableItemsToDelete) > 0 {
		if err := ass.repo.DeleteAuditableItems(auditableItemsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	go func() {
		existingAuditSchedule, err := ass.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		auditableItemsAdded := []uint{}
		auditableItemsRemoved := []uint{}
		for id := range newAuditableItemIDs {
			if !existingAuditableItemIDs[id] {
				auditableItemsAdded = append(auditableItemsAdded, id)
			}
		}
		for id := range existingAuditableItemIDs {
			if !newAuditableItemIDs[id] {
				auditableItemsRemoved = append(auditableItemsRemoved, id)
			}
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(auditableItemsAdded) > 0 || len(auditableItemsRemoved) > 0 {
			associatedItems["auditable_items_ids"] = map[string]interface{}{
				"added":   auditableItemsAdded,
				"removed": auditableItemsRemoved,
			}
		}
		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		//log changes
		oldData := map[string]interface{}{
			"audit_activity_type":   oldAuditSchedule.AuditActivityType.TitleEn,
			"title":                 oldAuditSchedule.Title,
			"description":           oldAuditSchedule.Description,
			"start_date":            oldAuditSchedule.StartDate,
			"target_date":           oldAuditSchedule.TargetDate,
			"responsible_user":      fmt.Sprintf("%s %s", oldAuditSchedule.ResponsibleUser.FirstNameEn, oldAuditSchedule.ResponsibleUser.LastNameEn),
			"audit_schedule_status": oldAuditSchedule.AuditScheduleStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"audit_activity_type":   existingAuditSchedule.AuditActivityType.TitleEn,
			"title":                 existingAuditSchedule.Title,
			"description":           existingAuditSchedule.Description,
			"start_date":            existingAuditSchedule.StartDate,
			"target_date":           existingAuditSchedule.TargetDate,
			"responsible_user":      fmt.Sprintf("%s %s", existingAuditSchedule.ResponsibleUser.FirstNameEn, existingAuditSchedule.ResponsibleUser.LastNameEn),
			"audit_schedule_status": existingAuditSchedule.AuditScheduleStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := ass.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditSchedule{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditSchedule{}.TableName(),
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

func (ass *AuditScheduleService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ass.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	var deletedAuditSchedule []map[string]interface{}
	for _, id := range ids {
		auditScheduleDetail, err := ass.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}
		// Check if auditable items exist
		if len(auditScheduleDetail.AuditableItems) > 0 {
			var deletAuditableItemIds []*uint
			for _, item := range auditScheduleDetail.AuditableItems {
				deletAuditableItemIds = append(deletAuditableItemIds, &item.ID)
			}
			if err := ass.repo.DeleteAuditableItems(deletAuditableItemIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := ass.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedAuditSchedule = append(deletedAuditSchedule, map[string]interface{}{
			"id":       id,
			"audit_id": auditScheduleDetail.AuditID,
			"title":    auditScheduleDetail.Title,
		})
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := ass.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditSchedule{}.TableName(),
					id,
					deletedAuditSchedule,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditSchedule{}.TableName(),
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

func (ass *AuditScheduleService) validateAuditScheduleFormDTO(auditScheduleFormDTO *auditDTOs.AuditScheduleFormDTO) error {
	if _, err := ass.auditActivityTypeSvc.GetByID(auditScheduleFormDTO.AuditActivityTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit activity type does not exist"), "not_exist", "audit activity type")
	}
	if _, err := ass.userSvc.GetByID(*auditScheduleFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}
	// Validate each AuditableItemId
	for _, itemID := range auditScheduleFormDTO.AuditableItemIds {
		if itemID != nil {
			if _, err := ass.auditableItemSvc.GetByID(*itemID); err != nil {
				return utility.NewValidationError(fmt.Errorf("auditable item with ID %d does not exist", *itemID), "not_exist", "auditable item")
			}
		}
	}
	return nil
}

func (ass *AuditScheduleService) parseAuditScheduleDates(dto *auditDTOs.AuditScheduleFormDTO) (*auditDTOs.AuditScheduleParsedDates, error) {
	parsed := &auditDTOs.AuditScheduleParsedDates{}

	var err error
	if dto.StartDate != nil && *dto.StartDate != "" {
		parsed.StartDate, err = utility.ParseDate(*dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.TargetDate != nil && *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (ass *AuditScheduleService) GetAllSchedulesByAuditID(auditID uint, tx ...utility.DatabaseTransaction) ([]*models.AuditSchedule, error) {
	return ass.repo.GetAllSchedulesByAuditID(auditID, tx...)
}

func (ass *AuditScheduleService) GetScheduleSummary(auditID uint) (*auditDTOs.AuditScheduleSummaryDTO, error) {
	pendingCount, completedCount, err := ass.repo.GetScheduleSummary(auditID)
	if err != nil {
		return nil, err
	}

	summary := &auditDTOs.AuditScheduleSummaryDTO{
		PendingCount:   pendingCount,
		CompletedCount: completedCount,
	}

	return summary, nil
}

func (ass *AuditScheduleService) StatusUpdate(auditScheduleID uint, auditScheduleStatusUpdateFormDTO *auditDTOs.AuditScheduleStatusUpdateFormDTO, currentUserID uint, tx utility.DatabaseTransaction) error {
	existingAuditSchedule, err := ass.GetByID(auditScheduleID)
	if err != nil {
		return err
	}

	existingAuditSchedule.AuditScheduleStatusID = auditScheduleStatusUpdateFormDTO.AuditScheduleStatusID
	existingAuditSchedule.Comment = auditScheduleStatusUpdateFormDTO.Comment
	existingAuditSchedule.Percentage = auditScheduleStatusUpdateFormDTO.Percentage
	existingAuditSchedule.UpdatedBYID = &currentUserID

	err = ass.repo.Update(existingAuditSchedule, tx)
	if err != nil {
		return err
	}

	return nil
}
