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

type AuditPlanScheduleServiceInterface interface {
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanScheduleListQueryParamsDTO) (*auditDTOs.AuditPlanScheduleListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanScheduleDetailDTO, error)
	GetByID(id uint) (*models.AuditPlanSchedule, error)
	Create(auditPlanID uint, currentUserID uint, auditPlanScheduleFormDTO *auditDTOs.AuditPlanScheduleFormDTO) (uint, error)
	Update(id uint, auditPlanID uint, currentUserID uint, auditPlanScheduleFormDTO *auditDTOs.AuditPlanScheduleFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetAllSchedulesByAuditPlanID(auditPlanID uint, tx ...utility.DatabaseTransaction) ([]*models.AuditPlanSchedule, error)
}

type AuditPlanScheduleService struct {
	repo                   repositories.AuditPlanScheduleRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	auditPlanSvc           AuditPlanServiceInterface
	auditActivityTypeSvc   AuditActivityTypeServiceInterface
	userSvc                organizationServices.UserServiceInterface
	auditScheduleStatusSvc AuditScheduleStatusServiceInterface
	auditableItemSvc       AuditableItemServiceInterface
}

func NewAuditPlanScheduleService(
	repo repositories.AuditPlanScheduleRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	auditActivityTypeSvc AuditActivityTypeServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditScheduleStatusSvc AuditScheduleStatusServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
) AuditPlanScheduleServiceInterface {
	return &AuditPlanScheduleService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		auditPlanSvc:           auditPlanSvc,
		auditActivityTypeSvc:   auditActivityTypeSvc,
		userSvc:                userSvc,
		auditScheduleStatusSvc: auditScheduleStatusSvc,
		auditableItemSvc:       auditableItemSvc,
	}
}

func (apss *AuditPlanScheduleService) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanScheduleListQueryParamsDTO) (*auditDTOs.AuditPlanScheduleListResponseDTO, error) {
	if _, err := apss.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return nil, err
	}
	return apss.repo.GetAll(auditPlanID, params)
}

func (apss *AuditPlanScheduleService) GetByID(id uint) (*models.AuditPlanSchedule, error) {
	auditPlanSchedule, err := apss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlanSchedule, nil
}

func (apss *AuditPlanScheduleService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanScheduleDetailDTO, error) {
	var auditPlanScheduleDTO auditDTOs.AuditPlanScheduleDetailDTO

	auditPlanSchedule, err := apss.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditPlanSchedule, &auditPlanScheduleDTO)
	if err != nil {
		return nil, err
	}
	if auditPlanSchedule.StartDate != nil {
		auditPlanScheduleDTO.StartDateString = auditPlanSchedule.StartDate.Format("2006-01-02")
	}
	if auditPlanSchedule.TargetDate != nil {
		auditPlanScheduleDTO.TargetDateString = auditPlanSchedule.TargetDate.Format("2006-01-02")
	}
	auditPlanScheduleDTO.CreatedAt = auditPlanSchedule.CreatedAt
	auditPlanScheduleDTO.UpdatedAt = auditPlanSchedule.UpdatedAt

	return &auditPlanScheduleDTO, nil
}

func (apss *AuditPlanScheduleService) Create(auditPlanID uint, currentUserID uint, auditPlanScheduleFormDTO *auditDTOs.AuditPlanScheduleFormDTO) (uint, error) {
	tx := apss.repo.BeginTransaction()
	if _, err := apss.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return 0, err
	}
	dates, err := apss.parseAuditPlanScheduleDates(auditPlanScheduleFormDTO)
	if err != nil {
		return 0, err
	}
	auditScheduleStatus, err := apss.auditScheduleStatusSvc.GetByType("scheduled")
	if err != nil {
		return 0, err
	}

	if err := apss.validateAuditPlanScheduleFormDTO(auditPlanScheduleFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditPlanSchedule := &models.AuditPlanSchedule{
		AuditPlanID:           auditPlanID,
		AuditActivityTypeID:   auditPlanScheduleFormDTO.AuditActivityTypeID,
		Title:                 auditPlanScheduleFormDTO.Title,
		Description:           auditPlanScheduleFormDTO.Description,
		StartDate:             dates.StartDate,
		TargetDate:            dates.TargetDate,
		ResponsibleUserID:     auditPlanScheduleFormDTO.ResponsibleUserID,
		AuditScheduleStatusID: uint(auditScheduleStatus.ID),
		CreatedByID:           currentUserID,
	}

	id, err := apss.repo.Create(auditPlanSchedule, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var auditableItems []*models.AuditPlanScheduleAuditableItem
	for _, itemID := range auditPlanScheduleFormDTO.AuditableItemIds {
		auditableItem := &models.AuditPlanScheduleAuditableItem{
			AuditPlanScheduleID: id,
			AuditableItemID:     *itemID,
			CreatedBYID:         currentUserID,
		}
		auditableItems = append(auditableItems, auditableItem)
	}

	if len(auditableItems) > 0 {
		if err := apss.repo.CreateAuditableItems(auditableItems, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := apss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditPlanSchedule{}.TableName(),
			id,
			auditPlanScheduleFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlanSchedule{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (apss *AuditPlanScheduleService) Update(id uint, auditPlanID uint, currentUserID uint, auditPlanScheduleFormDTO *auditDTOs.AuditPlanScheduleFormDTO) error {
	tx := apss.repo.BeginTransaction()
	if _, err := apss.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return err
	}
	dates, err := apss.parseAuditPlanScheduleDates(auditPlanScheduleFormDTO)
	if err != nil {
		return err
	}

	if err := apss.validateAuditPlanScheduleFormDTO(auditPlanScheduleFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlanSchedule, err := apss.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingAuditPlanSchedule.AuditActivityTypeID = auditPlanScheduleFormDTO.AuditActivityTypeID
	existingAuditPlanSchedule.Title = auditPlanScheduleFormDTO.Title
	existingAuditPlanSchedule.Description = auditPlanScheduleFormDTO.Description
	existingAuditPlanSchedule.StartDate = dates.StartDate
	existingAuditPlanSchedule.TargetDate = dates.TargetDate
	existingAuditPlanSchedule.ResponsibleUserID = auditPlanScheduleFormDTO.ResponsibleUserID
	existingAuditPlanSchedule.UpdatedByID = &currentUserID

	err = apss.repo.Update(existingAuditPlanSchedule, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingSchedule, err := apss.GetFullDetailsByID(id, tx)
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
	for _, itemID := range auditPlanScheduleFormDTO.AuditableItemIds {
		newAuditableItemIDs[*itemID] = true
	}

	// Add new auditable items
	var auditableItemsToAdd []*models.AuditPlanScheduleAuditableItem
	for _, newItemID := range auditPlanScheduleFormDTO.AuditableItemIds {
		if _, exists := existingAuditableItemIDs[*newItemID]; !exists {
			// Create a new auditable item record
			auditableItem := &models.AuditPlanScheduleAuditableItem{
				AuditPlanScheduleID: id,
				AuditableItemID:     *newItemID,
				CreatedBYID:         currentUserID,
			}
			auditableItemsToAdd = append(auditableItemsToAdd, auditableItem)
		}
	}

	// Save the new auditable items
	if len(auditableItemsToAdd) > 0 {
		if err := apss.repo.CreateAuditableItems(auditableItemsToAdd, tx); err != nil {
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
		if err := apss.repo.DeleteAuditableItems(auditableItemsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	go func() {
		err := apss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditPlanSchedule{}.TableName(),
			id,
			auditPlanScheduleFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditPlanSchedule{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (apss *AuditPlanScheduleService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := apss.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		auditPlanScheduleDetail, err := apss.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}
		// Check if controls exist
		if len(auditPlanScheduleDetail.AuditableItems) > 0 {
			var deletAuditableItemIds []*uint
			for _, item := range auditPlanScheduleDetail.AuditableItems {
				deletAuditableItemIds = append(deletAuditableItemIds, &item.ID)
			}
			if err := apss.repo.DeleteAuditableItems(deletAuditableItemIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := apss.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := apss.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditPlanSchedule{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlanSchedule{}.TableName(),
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

func (apss *AuditPlanScheduleService) validateAuditPlanScheduleFormDTO(auditPlanScheduleFormDTO *auditDTOs.AuditPlanScheduleFormDTO) error {
	if _, err := apss.auditActivityTypeSvc.GetByID(auditPlanScheduleFormDTO.AuditActivityTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit activity type does not exist"), "not_exist", "audit activity type")
	}
	if _, err := apss.userSvc.GetByID(*auditPlanScheduleFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}
	// Validate each AuditableItemId
	for _, itemID := range auditPlanScheduleFormDTO.AuditableItemIds {
		if itemID != nil {
			if _, err := apss.auditableItemSvc.GetByID(*itemID); err != nil {
				return utility.NewValidationError(fmt.Errorf("auditable item with ID %d does not exist", *itemID), "not_exist", "auditable item")
			}
		}
	}
	return nil
}

func (apss *AuditPlanScheduleService) parseAuditPlanScheduleDates(dto *auditDTOs.AuditPlanScheduleFormDTO) (*auditDTOs.AuditPlanScheduleParsedDates, error) {
	parsed := &auditDTOs.AuditPlanScheduleParsedDates{}

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

func (apss *AuditPlanScheduleService) GetAllSchedulesByAuditPlanID(auditPlanID uint, tx ...utility.DatabaseTransaction) ([]*models.AuditPlanSchedule, error) {
	return apss.repo.GetAllSchedulesByAuditPlanID(auditPlanID, tx...)
}
