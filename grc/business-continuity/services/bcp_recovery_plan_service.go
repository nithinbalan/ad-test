package services

import (
	"errors"
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	BcpRecoveryPlanRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BcpRecoveryPlanServiceInterface interface {
	GetAll(bcpId uint, params *businessContinuityDTOs.BcpRecoveryPlanListQueryParamsDTO) (*businessContinuityDTOs.BcpRecoveryPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpRecoveryPlanDetailDTO, error)
	GetByID(id uint) (*models.BcpRecoveryPlan, error)
	Create(bcpId uint, currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BcpRecoveryPlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, BcpRecoveryPlan *businessContinuityDTOs.BcpRecoveryPlanFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type BcpRecoveryPlanService struct {
	repo           BcpRecoveryPlanRepositories.BcpRecoveryPlanRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewBcpRecoveryPlanService(
	repo BcpRecoveryPlanRepositories.BcpRecoveryPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) BcpRecoveryPlanServiceInterface {
	return &BcpRecoveryPlanService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (brps *BcpRecoveryPlanService) GetAll(bcpId uint, params *businessContinuityDTOs.BcpRecoveryPlanListQueryParamsDTO) (*businessContinuityDTOs.BcpRecoveryPlanListResponseDTO, error) {
	return brps.repo.GetAll(bcpId, params)
}

func (brps BcpRecoveryPlanService) GetByID(id uint) (*models.BcpRecoveryPlan, error) {
	BcpRecoveryPlan, err := brps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return BcpRecoveryPlan, nil
}

func (brps BcpRecoveryPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpRecoveryPlanDetailDTO, error) {
	// Initialize your DTO
	var BcpRecoveryPlanDTO businessContinuityDTOs.BcpRecoveryPlanDetailDTO

	BcpRecoveryPlan, err := brps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(BcpRecoveryPlan, &BcpRecoveryPlanDTO)
	if err != nil {
		return nil, err
	}

	BcpRecoveryPlanDTO.CreatedAt = BcpRecoveryPlan.CreatedAt
	BcpRecoveryPlanDTO.UpdatedAt = BcpRecoveryPlan.UpdatedAt

	return &BcpRecoveryPlanDTO, nil
}

func (brps *BcpRecoveryPlanService) Create(bcpId uint, currentUserID uint, BcpRecoveryPlanFormDTO *businessContinuityDTOs.BcpRecoveryPlanFormDTO) (uint, error) {
	tx := brps.repo.BeginTransaction()

	responsibleUnitID, err := brps.userSvc.GetUserPrimaryUnitID(uint(BcpRecoveryPlanFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp recovery plan responsible unit")
	}

	BcpRecoveryPlan := &models.BcpRecoveryPlan{
		BcpPhaseID:               *BcpRecoveryPlanFormDTO.BcpPhaseID,
		BusinessContinuityPlanID: bcpId,
		BcpProcessID:             *BcpRecoveryPlanFormDTO.BcpProcessID,
		SeverityLevelID:          BcpRecoveryPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:        BcpRecoveryPlanFormDTO.ResponsibleUserID,
		ResponsibleUnitID:        responsibleUnitID,
		Title:                    BcpRecoveryPlanFormDTO.Title,
		TargetCompletion:         BcpRecoveryPlanFormDTO.TargetCompletion,
		TargetCompletionUnitId:   BcpRecoveryPlanFormDTO.TargetCompletionUnitID,
		CreatedBYID:              currentUserID,
	}

	id, err := brps.repo.Create(BcpRecoveryPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingBcpRecoveryPlan, err := brps.GetFullDetailsByID(BcpRecoveryPlan.ID)
		if err != nil {
			utility.LogError("Failed to fetch bcp recovery plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = brps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BcpRecoveryPlan{}.TableName(),
			id,
			existingBcpRecoveryPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BcpRecoveryPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (brps *BcpRecoveryPlanService) Update(id uint, currentUserID uint, BcpRecoveryPlanFormDTO *businessContinuityDTOs.BcpRecoveryPlanFormDTO) error {
	tx := brps.repo.BeginTransaction()

	existingBcpRecoveryPlan, err := brps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldBcpRecoveryPlan, err := brps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingBcpRecoveryPlan.ResponsibleUserID != BcpRecoveryPlanFormDTO.ResponsibleUserID {
		responsibleUnitID, err := brps.userSvc.GetUserPrimaryUnitID(existingBcpRecoveryPlan.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingBcpRecoveryPlan.ResponsibleUnitID = responsibleUnitID
	}

	existingBcpRecoveryPlan.BcpPhaseID = *BcpRecoveryPlanFormDTO.BcpPhaseID
	existingBcpRecoveryPlan.BcpProcessID = *BcpRecoveryPlanFormDTO.BcpProcessID
	existingBcpRecoveryPlan.SeverityLevelID = BcpRecoveryPlanFormDTO.SeverityLevelID
	existingBcpRecoveryPlan.ResponsibleUserID = BcpRecoveryPlanFormDTO.ResponsibleUserID
	existingBcpRecoveryPlan.Title = BcpRecoveryPlanFormDTO.Title
	existingBcpRecoveryPlan.TargetCompletion = BcpRecoveryPlanFormDTO.TargetCompletion
	existingBcpRecoveryPlan.TargetCompletionUnitId = BcpRecoveryPlanFormDTO.TargetCompletionUnitID
	existingBcpRecoveryPlan.UpdatedBYID = &currentUserID

	err = brps.repo.Update(existingBcpRecoveryPlan, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingBcpRecoveryPlan, err := brps.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"severity_level":         oldBcpRecoveryPlan.SeverityLevel.TitleEn,
			"responsible_user":       fmt.Sprintf("%s %s", oldBcpRecoveryPlan.ResponsibleUser.FirstNameEN, oldBcpRecoveryPlan.ResponsibleUser.LastNameEN),
			"title":                  oldBcpRecoveryPlan.Title,
			"bcp_phase":              oldBcpRecoveryPlan.BcpPhase.TitleEn,
			"bcp_process":            oldBcpRecoveryPlan.BcpProcess.Process.TitleEn,
			"target_completion":      oldBcpRecoveryPlan.TargetCompletion,
			"target_completion_unit": oldBcpRecoveryPlan.TargetCompletionUnit.TitleEn,
		}

		newData := map[string]interface{}{
			"severity_level":         existingBcpRecoveryPlan.SeverityLevel.TitleEn,
			"responsible_user":       fmt.Sprintf("%s %s", existingBcpRecoveryPlan.ResponsibleUser.FirstNameEN, existingBcpRecoveryPlan.ResponsibleUser.LastNameEN),
			"title":                  existingBcpRecoveryPlan.Title,
			"bcp_phase":              existingBcpRecoveryPlan.BcpPhase.TitleEn,
			"bcp_process":            existingBcpRecoveryPlan.BcpProcess.Process.TitleEn,
			"target_completion":      existingBcpRecoveryPlan.TargetCompletion,
			"target_completion_unit": existingBcpRecoveryPlan.TargetCompletionUnit.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := brps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BcpRecoveryPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BcpRecoveryPlan{}.TableName(),
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

func (brps *BcpRecoveryPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := brps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedBcpRecoveryPlan []map[string]interface{}

	for _, id := range ids {
		BcpRecoveryPlan, err := brps.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedBcpRecoveryPlan = append(deletedBcpRecoveryPlan, map[string]interface{}{
			"id":    id,
			"title": BcpRecoveryPlan.Title,
		})

		// Delete the business
		if err := brps.repo.Delete(id, tx); err != nil {
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

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := brps.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.BcpRecoveryPlan{}.TableName(),
				successfulIds,
				deletedBcpRecoveryPlan,
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
