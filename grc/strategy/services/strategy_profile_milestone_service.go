package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileMilestoneRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileMilestoneServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileMilestoneListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileMilestoneListDTO, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileMilestoneDetailDTO, error)
	GetByID(id uint) (*models.StrategyProfileMilestone, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileMilestoneFormDTO) (uint, error)
	Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileMilestone *strategyDTOs.StrategyProfileMilestoneFormDTO) error
	Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategyProfileMilestoneService struct {
	repo               strategyProfileMilestoneRepositories.StrategyProfileMilestoneRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	strategyProfileSvc StrategyProfileServiceInterface
}

func NewStrategyProfileMilestoneService(
	repo strategyProfileMilestoneRepositories.StrategyProfileMilestoneRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
) StrategyProfileMilestoneServiceInterface {
	return &StrategyProfileMilestoneService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		userSettingSvc:     userSettingSvc,
		strategyProfileSvc: strategyProfileSvc,
	}
}

func (sps *StrategyProfileMilestoneService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileMilestoneListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileMilestoneListDTO, error) {
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileMilestoneService) GetByID(id uint) (*models.StrategyProfileMilestone, error) {
	strategyProfileMilestone, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileMilestone, nil
}

func (sps StrategyProfileMilestoneService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileMilestoneDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileMilestoneDTO strategyDTOs.StrategyProfileMilestoneDetailDTO

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	strategyProfileMilestone, err := sps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfileMilestone, &strategyProfileMilestoneDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileMilestoneDTO.StrategyProfile.StartDateString = strategyProfileMilestone.StrategyProfile.StartDate.Format("2006-01-02")
	strategyProfileMilestoneDTO.StrategyProfile.EndDateString = strategyProfileMilestone.StrategyProfile.EndDate.Format("2006-01-02")

	strategyProfileMilestoneDTO.CreatedAt = strategyProfileMilestone.CreatedAt
	strategyProfileMilestoneDTO.UpdatedAt = strategyProfileMilestone.UpdatedAt

	return &strategyProfileMilestoneDTO, nil
}

func (sps *StrategyProfileMilestoneService) Create(strategyProfileID uint, currentUserID uint, strategyProfileMilestoneFormDTO *strategyDTOs.StrategyProfileMilestoneFormDTO) (uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return 0, err
	}

	tx := sps.repo.BeginTransaction()

	strategyProfileMilestone := &models.StrategyProfileMilestone{
		StrategyProfileID:           strategyProfileID,
		StrategyProfileFocusAreaID:  strategyProfileMilestoneFormDTO.StrategyProfileFocusAreaID,
		StrategyProfileGoalID:       strategyProfileMilestoneFormDTO.StrategyProfileGoalID,
		StrategyProfileObjectiveID:  strategyProfileMilestoneFormDTO.StrategyProfileObjectiveID,
		StrategyProfileInitiativeID: strategyProfileMilestoneFormDTO.StrategyProfileInitiativeID,
		StrategyProfileProjectID:    strategyProfileMilestoneFormDTO.StrategyProfileProjectID,
		Title:                       strategyProfileMilestoneFormDTO.Title,
		CreatedBYID:                 currentUserID,
	}
	id, err := sps.repo.Create(strategyProfileMilestone, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileMilestone, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			utility.LogError("Failed to fetch strategy profile milestone", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileMilestone{}.TableName(),
			id,
			existingStrategyProfileMilestone,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileMilestone{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (sps *StrategyProfileMilestoneService) Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileMilestoneFormDTO *strategyDTOs.StrategyProfileMilestoneFormDTO) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	existingStrategyProfileMilestone, err := sps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfileMilestone, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStrategyProfileMilestone.Title = strategyProfileMilestoneFormDTO.Title
	existingStrategyProfileMilestone.UpdatedBYID = &currentUserID

	err = sps.repo.Update(existingStrategyProfileMilestone, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileMilestone, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title": oldStrategyProfileMilestone.Title,
		}

		newData := map[string]interface{}{
			"title": existingStrategyProfileMilestone.Title,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := sps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StrategyProfileMilestone{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileMilestone{}.TableName(),
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

func (sps *StrategyProfileMilestoneService) Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategyProfileMilestone []map[string]interface{}

	for _, id := range ids {
		strategyProfileMilestone, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStrategyProfileMilestone = append(deletedStrategyProfileMilestone, map[string]interface{}{
			"id":       id,
			"title_en": strategyProfileMilestone.Title,
		})

		// Delete the business
		if err := sps.repo.Delete(id, tx); err != nil {
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
			err := sps.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StrategyProfileMilestone{}.TableName(),
				successfulIds,
				deletedStrategyProfileMilestone,
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
