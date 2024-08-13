package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanChecklistServiceInterface interface {
	GetAll(params *auditDTOs.AuditPlanChecklistListQueryParamsDTO) (*auditDTOs.AuditPlanChecklistListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanChecklistDetailDTO, error)
	GetByID(id uint) (*models.AuditPlanChecklist, error)
	Create(currentUserID uint, auditPlanChecklistFormDTO *auditDTOs.AuditPlanChecklistFormDTO) (uint, error)
	Update(id uint, currentUserID uint, auditPlanChecklistUpdateFormDTO *auditDTOs.AuditPlanChecklistUpdateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditPlanChecklistService struct {
	repo             repositories.AuditPlanChecklistRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditPlanSvc     AuditPlanServiceInterface
	auditableItemSvc AuditableItemServiceInterface
}

func NewAuditPlanChecklistService(
	repo repositories.AuditPlanChecklistRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
) AuditPlanChecklistServiceInterface {
	return &AuditPlanChecklistService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditPlanSvc:     auditPlanSvc,
		auditableItemSvc: auditableItemSvc,
	}
}

func (aps *AuditPlanChecklistService) GetAll(params *auditDTOs.AuditPlanChecklistListQueryParamsDTO) (*auditDTOs.AuditPlanChecklistListResponseDTO, error) {
	return aps.repo.GetAll(params)
}

func (aps *AuditPlanChecklistService) GetByID(id uint) (*models.AuditPlanChecklist, error) {
	auditPlanChecklist, err := aps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlanChecklist, nil
}

func (aps *AuditPlanChecklistService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanChecklistDetailDTO, error) {
	var auditPlanChecklistDTO auditDTOs.AuditPlanChecklistDetailDTO

	auditPlanChecklist, err := aps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditPlanChecklist, &auditPlanChecklistDTO)
	if err != nil {
		return nil, err
	}

	auditPlanChecklistDTO.CreatedAt = auditPlanChecklist.CreatedAt
	auditPlanChecklistDTO.UpdatedAt = auditPlanChecklist.UpdatedAt

	return &auditPlanChecklistDTO, nil
}

func (aps *AuditPlanChecklistService) Create(currentUserID uint, auditPlanChecklistFormDTO *auditDTOs.AuditPlanChecklistFormDTO) (uint, error) {
	tx := aps.repo.BeginTransaction()

	if err := aps.validateAuditPlanChecklistFormDTO(auditPlanChecklistFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditPlanChecklist := &models.AuditPlanChecklist{
		AuditPlanID:     referenceNumber,
		AuditableItemID: auditPlanChecklistFormDTO.AuditableItemId,
		Question:        auditPlanChecklistFormDTO.Question,
		CreatedBYID:     currentUserID,
	}

	id, err := aps.repo.Create(auditPlanChecklist, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditPlanChecklist{}.TableName(),
			id,
			auditPlanChecklistFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlanChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (aps *AuditPlanChecklistService) Update(id uint, currentUserID uint, auditPlanChecklistFormDTO *auditDTOs.AuditPlanChecklistFormDTO) error {
	tx := aps.repo.BeginTransaction()

	existingAuditPlanChecklist, err := aps.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := aps.validateAuditPlanChecklistFormDTO(auditPlanChecklistFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlanChecklist.AuditCategoryID = auditPlanChecklistFormDTO.AuditCategoryID
	existingAuditPlanChecklist.AuditTeamID = auditPlanChecklistFormDTO.AuditTeamID
	existingAuditPlanChecklist.TeamLeadID = auditPlanChecklistFormDTO.TeamLeadID
	existingAuditPlanChecklist.Title = auditPlanChecklistFormDTO.Title
	existingAuditPlanChecklist.Description = auditPlanChecklistFormDTO.Description
	existingAuditPlanChecklist.StartDate = dates.StartDate
	existingAuditPlanChecklist.TargetDate = dates.TargetDate
	existingAuditPlanChecklist.Purpose = auditPlanChecklistFormDTO.Purpose
	existingAuditPlanChecklist.Objective = auditPlanChecklistFormDTO.Objective
	existingAuditPlanChecklist.UpdatedBYID = &currentUserID

	err = aps.repo.Update(existingAuditPlanChecklist, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlanChecklistDetails, err := aps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingUserIDs := make(map[uint]bool)
	for _, member := range existingAuditPlanChecklistDetails.Members {
		existingUserIDs[member.UserID] = true
	}

	newUserIDs := make(map[uint]bool)
	for _, userID := range auditPlanChecklistFormDTO.UserIds {
		newUserIDs[*userID] = true
	}

	var usersToAdd []*models.AuditPlanChecklistAuditMember
	for _, userID := range auditPlanChecklistFormDTO.UserIds {
		if _, exists := existingUserIDs[*userID]; !exists {
			user := &models.AuditPlanChecklistAuditMember{
				AuditPlanChecklistID: id,
				UserID:               *userID,
				CreatedBYID:          currentUserID,
			}
			usersToAdd = append(usersToAdd, user)
		}
	}

	var usersToRemove []*models.AuditPlanChecklistAuditMember
	for userID := range existingUserIDs {
		if _, exists := newUserIDs[userID]; !exists {
			usersToRemove = append(usersToRemove, &models.AuditPlanChecklistAuditMember{
				AuditPlanChecklistID: id,
				UserID:               userID,
			})
		}
	}

	// Extract user IDs to remove
	var userIDsToRemove []*uint
	for _, user := range usersToRemove {
		userID := user.UserID
		userIDsToRemove = append(userIDsToRemove, &userID)
	}

	// Add new members
	if len(usersToAdd) > 0 {
		if err := aps.repo.CreateAuditPlanChecklistMembers(usersToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove old members
	if len(userIDsToRemove) > 0 {
		if err := aps.repo.DeleteAuditPlanChecklistMembers(userIDsToRemove, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Handle auditable items
	existingAuditableItemIDs := make(map[uint]bool)
	for _, item := range existingAuditPlanChecklistDetails.AuditableItems {
		existingAuditableItemIDs[item.AuditableItemID] = true
	}

	newAuditableItemIDs := make(map[uint]bool)
	for _, itemID := range auditPlanChecklistFormDTO.AuditableItemIds {
		newAuditableItemIDs[*itemID] = true
	}

	var itemsToAdd []*models.AuditPlanChecklistAuditableItem
	for _, itemID := range auditPlanChecklistFormDTO.AuditableItemIds {
		if _, exists := existingAuditableItemIDs[*itemID]; !exists {
			item := &models.AuditPlanChecklistAuditableItem{
				AuditPlanChecklistID: id,
				AuditableItemID:      *itemID,
				CreatedBYID:          currentUserID,
			}
			itemsToAdd = append(itemsToAdd, item)
		}
	}

	var auditableItemIDsToRemove []*uint
	for itemID := range existingAuditableItemIDs {
		if _, exists := newAuditableItemIDs[itemID]; !exists {
			auditableItemIDsToRemove = append(auditableItemIDsToRemove, &itemID)
		}
	}

	if len(itemsToAdd) > 0 {
		if err := aps.repo.CreateAuditPlanChecklistAuditableItems(itemsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(auditableItemIDsToRemove) > 0 {
		if err := aps.repo.DeleteAuditPlanChecklistAuditableItems(auditableItemIDsToRemove, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditPlanChecklist{}.TableName(),
			id,
			auditPlanChecklistFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditPlanChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (aps *AuditPlanChecklistService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := aps.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		// Get existing audit plan checklist details
		existingAuditPlanChecklist, err := aps.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if len(existingAuditPlanChecklist.AuditableItems) > 0 {
			var deleteAuditableItemIDs []*uint
			for _, item := range existingAuditPlanChecklist.AuditableItems {
				deleteAuditableItemIDs = append(deleteAuditableItemIDs, item.ID)
			}
			if err := aps.repo.DeleteAuditPlanChecklistAuditableItems(deleteAuditableItemIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(existingAuditPlanChecklist.Members) > 0 {
			var deleteMemberIDs []*uint
			for _, member := range existingAuditPlanChecklist.Members {
				deleteMemberIDs = append(deleteMemberIDs, member.ID)
			}
			if err := aps.repo.DeleteAuditPlanChecklistMembers(deleteMemberIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the audit plan checklist itself
		if err := aps.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := aps.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditPlanChecklist{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlanChecklist{}.TableName(),
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

func (aps *AuditPlanChecklistService) validateAuditPlanChecklistFormDTO(dto *auditDTOs.AuditPlanChecklistFormDTO) error {
	if _, err := aps.auditableItemSvc.GetByID(dto.AuditableItemId); err != nil {
		return utility.NewValidationError(fmt.Errorf("auditable item does not exist"), "not_exist", "auditable item")
	}

	return nil
}

func (aps *AuditPlanChecklistService) parseAuditPlanChecklistDates(dto *auditDTOs.AuditPlanChecklistFormDTO) (*auditDTOs.AuditPlanChecklistParsedDates, error) {
	parsed := &auditDTOs.AuditPlanChecklistParsedDates{}

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
