package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanServiceInterface interface {
	GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO) (*auditDTOs.AuditPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanDetailDTO, error)
	GetByID(id uint) (*models.AuditPlan, error)
	Create(currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanFormDTO) error
	Publish(id uint, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, auditStatusID uint, tx utility.DatabaseTransaction) error
}

type AuditPlanService struct {
	repo             repositories.AuditPlanRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditCategorySvc AuditCategoryServiceInterface
	auditStatusSvc   AuditStatusServiceInterface
	auditTeamSvc     AuditTeamServiceInterface
	userSvc          organizationServices.UserServiceInterface
	auditableItemSvc AuditableItemServiceInterface
}

func NewAuditPlanService(
	repo repositories.AuditPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditCategorySvc AuditCategoryServiceInterface,
	auditStatusSvc AuditStatusServiceInterface,
	auditTeamSvc AuditTeamServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
) AuditPlanServiceInterface {
	return &AuditPlanService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditCategorySvc: auditCategorySvc,
		auditStatusSvc:   auditStatusSvc,
		auditTeamSvc:     auditTeamSvc,
		userSvc:          userSvc,
		auditableItemSvc: auditableItemSvc,
	}
}

func (aps *AuditPlanService) GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO) (*auditDTOs.AuditPlanListResponseDTO, error) {
	return aps.repo.GetAll(params)
}

func (aps *AuditPlanService) GetByID(id uint) (*models.AuditPlan, error) {
	auditPlan, err := aps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlan, nil
}

func (aps *AuditPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanDetailDTO, error) {
	var auditPlanDTO auditDTOs.AuditPlanDetailDTO

	auditPlan, err := aps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditPlan, &auditPlanDTO)
	if err != nil {
		return nil, err
	}

	if auditPlan.StartDate != nil {
		auditPlanDTO.StartDateString = auditPlan.StartDate.Format("2006-01-02")
	}
	if auditPlan.TargetDate != nil {
		auditPlanDTO.TargetDateString = auditPlan.TargetDate.Format("2006-01-02")
	}

	auditPlanDTO.CreatedAt = auditPlan.CreatedAt
	auditPlanDTO.UpdatedAt = auditPlan.UpdatedAt

	return &auditPlanDTO, nil
}

func (aps *AuditPlanService) Create(currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanFormDTO) (uint, error) {
	tx := aps.repo.BeginTransaction()
	dates, err := aps.parseAuditPlanDates(auditPlanFormDTO)
	if err != nil {
		return 0, err
	}

	if err := aps.validateAuditPlanFormDTO(auditPlanFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	referenceNumber, err := aps.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}
	auditStatus, err := aps.auditStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	auditPlan := &models.AuditPlan{
		ReferenceNumber: referenceNumber,
		AuditCategoryID: auditPlanFormDTO.AuditCategoryID,
		AuditTeamID:     auditPlanFormDTO.AuditTeamID,
		TeamLeadID:      auditPlanFormDTO.TeamLeadID,
		Title:           auditPlanFormDTO.Title,
		Description:     auditPlanFormDTO.Description,
		StartDate:       dates.StartDate,
		TargetDate:      dates.TargetDate,
		Purpose:         auditPlanFormDTO.Purpose,
		Objective:       auditPlanFormDTO.Objective,
		AuditStatusID:   uint(auditStatus.ID),
		CreatedBYID:     currentUserID,
	}

	id, err := aps.repo.Create(auditPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var auditMembers []*models.AuditPlanAuditMember
	for _, userID := range auditPlanFormDTO.UserIds {
		auditMember := &models.AuditPlanAuditMember{
			AuditPlanID: id,
			UserID:      *userID,
			CreatedBYID: currentUserID,
		}
		auditMembers = append(auditMembers, auditMember)
	}

	if len(auditMembers) > 0 {
		if err := aps.repo.CreateAuditPlanMembers(auditMembers, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var auditableItems []*models.AuditPlanAuditableItem
	for _, itemID := range auditPlanFormDTO.AuditableItemIds {
		auditableItem := &models.AuditPlanAuditableItem{
			AuditPlanID:     id,
			AuditableItemID: *itemID,
			CreatedBYID:     currentUserID,
		}
		auditableItems = append(auditableItems, auditableItem)
	}

	if len(auditableItems) > 0 {
		if err := aps.repo.CreateAuditPlanAuditableItems(auditableItems, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditPlan{}.TableName(),
			id,
			auditPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (aps *AuditPlanService) Update(id uint, currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanFormDTO) error {
	tx := aps.repo.BeginTransaction()
	dates, err := aps.parseAuditPlanDates(auditPlanFormDTO)
	if err != nil {
		return err
	}

	existingAuditPlan, err := aps.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := aps.validateAuditPlanFormDTO(auditPlanFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlan.AuditCategoryID = auditPlanFormDTO.AuditCategoryID
	existingAuditPlan.AuditTeamID = auditPlanFormDTO.AuditTeamID
	existingAuditPlan.TeamLeadID = auditPlanFormDTO.TeamLeadID
	existingAuditPlan.Title = auditPlanFormDTO.Title
	existingAuditPlan.Description = auditPlanFormDTO.Description
	existingAuditPlan.StartDate = dates.StartDate
	existingAuditPlan.TargetDate = dates.TargetDate
	existingAuditPlan.Purpose = auditPlanFormDTO.Purpose
	existingAuditPlan.Objective = auditPlanFormDTO.Objective
	existingAuditPlan.UpdatedBYID = &currentUserID

	err = aps.repo.Update(existingAuditPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlanDetails, err := aps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingUserIDs := make(map[uint]bool)
	for _, member := range existingAuditPlanDetails.Members {
		existingUserIDs[member.UserID] = true
	}

	newUserIDs := make(map[uint]bool)
	for _, userID := range auditPlanFormDTO.UserIds {
		newUserIDs[*userID] = true
	}

	var usersToAdd []*models.AuditPlanAuditMember
	for _, userID := range auditPlanFormDTO.UserIds {
		if _, exists := existingUserIDs[*userID]; !exists {
			user := &models.AuditPlanAuditMember{
				AuditPlanID: id,
				UserID:      *userID,
				CreatedBYID: currentUserID,
			}
			usersToAdd = append(usersToAdd, user)
		}
	}

	var usersToRemove []*models.AuditPlanAuditMember
	for userID := range existingUserIDs {
		if _, exists := newUserIDs[userID]; !exists {
			usersToRemove = append(usersToRemove, &models.AuditPlanAuditMember{
				AuditPlanID: id,
				UserID:      userID,
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
		if err := aps.repo.CreateAuditPlanMembers(usersToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove old members
	if len(userIDsToRemove) > 0 {
		if err := aps.repo.DeleteAuditPlanMembers(userIDsToRemove, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Handle auditable items
	existingAuditableItemIDs := make(map[uint]bool)
	for _, item := range existingAuditPlanDetails.AuditableItems {
		existingAuditableItemIDs[item.AuditableItemID] = true
	}

	newAuditableItemIDs := make(map[uint]bool)
	for _, itemID := range auditPlanFormDTO.AuditableItemIds {
		newAuditableItemIDs[*itemID] = true
	}

	var itemsToAdd []*models.AuditPlanAuditableItem
	for _, itemID := range auditPlanFormDTO.AuditableItemIds {
		if _, exists := existingAuditableItemIDs[*itemID]; !exists {
			item := &models.AuditPlanAuditableItem{
				AuditPlanID:     id,
				AuditableItemID: *itemID,
				CreatedBYID:     currentUserID,
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
		if err := aps.repo.CreateAuditPlanAuditableItems(itemsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(auditableItemIDsToRemove) > 0 {
		if err := aps.repo.DeleteAuditPlanAuditableItems(auditableItemIDsToRemove, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditPlan{}.TableName(),
			id,
			auditPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (aps *AuditPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := aps.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		// Get existing audit plan details
		existingAuditPlan, err := aps.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if len(existingAuditPlan.AuditableItems) > 0 {
			var deleteAuditableItemIDs []*uint
			for _, item := range existingAuditPlan.AuditableItems {
				deleteAuditableItemIDs = append(deleteAuditableItemIDs, item.ID)
			}
			if err := aps.repo.DeleteAuditPlanAuditableItems(deleteAuditableItemIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(existingAuditPlan.Members) > 0 {
			var deleteMemberIDs []*uint
			for _, member := range existingAuditPlan.Members {
				deleteMemberIDs = append(deleteMemberIDs, member.ID)
			}
			if err := aps.repo.DeleteAuditPlanMembers(deleteMemberIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the audit plan itself
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
					models.AuditPlan{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlan{}.TableName(),
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

func (aps *AuditPlanService) StatusUpdate(id uint, currentUserID uint, auditStatusID uint, tx utility.DatabaseTransaction) error {
	existingAuditPlan, err := aps.GetByID(id)
	if err != nil {
		return err
	}

	existingAuditPlan.AuditStatusID =  auditStatusID
	existingAuditPlan.UpdatedBYID = &currentUserID

	err = aps.repo.Update(existingAuditPlan, tx)

	if err != nil {
		return err
	}

	return nil
}

func (aps *AuditPlanService) Publish(id uint, currentUserID uint) error {
	tx := aps.repo.BeginTransaction()
	existingAuditPlan, err := aps.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}
	auditStatus, err := aps.auditStatusSvc.GetByType("published")
	if err != nil {
		return err
	}
	if existingAuditPlan.AuditStatusID == uint(auditStatus.ID) {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report has already been published.")
	}
	existingAuditPlan.AuditStatusID = uint(auditStatus.ID)
	existingAuditPlan.UpdatedBYID = &currentUserID

	err = aps.repo.Update(existingAuditPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditPlan{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (aps *AuditPlanService) validateAuditPlanFormDTO(dto *auditDTOs.AuditPlanFormDTO) error {
	if _, err := aps.auditCategorySvc.GetByID(dto.AuditCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit category does not exist"), "not_exist", "audit category")
	}

	if dto.AuditTeamID != nil {
		if _, err := aps.auditTeamSvc.GetByID(*dto.AuditTeamID); err != nil {
			return utility.NewValidationError(fmt.Errorf("audit team does not exist"), "not_exist", "audit team")
		}
	}

	if dto.TeamLeadID != nil {
		if _, err := aps.userSvc.GetByID(*dto.TeamLeadID); err != nil {
			return utility.NewValidationError(fmt.Errorf("team lead does not exist"), "not_exist", "team lead")
		}
	}

	// Validate UserIds
	for _, userID := range dto.UserIds {
		if userID != nil {
			if _, err := aps.userSvc.GetByID(*userID); err != nil {
				return utility.NewValidationError(fmt.Errorf("user does not exist"), "not_exist", "user")
			}
		}
	}

	// Validate AuditableItemIds
	for _, itemID := range dto.AuditableItemIds {
		if itemID != nil {
			if _, err := aps.auditableItemSvc.GetByID(*itemID); err != nil {
				return utility.NewValidationError(fmt.Errorf("auditable item does not exist"), "not_exist", "auditable item")
			}
		}
	}

	return nil
}

func (aps *AuditPlanService) parseAuditPlanDates(dto *auditDTOs.AuditPlanFormDTO) (*auditDTOs.AuditPlanParsedDates, error) {
	parsed := &auditDTOs.AuditPlanParsedDates{}

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
