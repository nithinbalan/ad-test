package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileObjectiveRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileObjectiveServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileObjectiveListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileObjectiveListDTO, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileObjectiveDetailDTO, error)
	GetByID(id uint) (*models.StrategyProfileObjective, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileObjectiveFormDTO) (uint, error)
	Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileObjective *strategyDTOs.StrategyProfileObjectiveFormDTO) error
	Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategyProfileObjectiveService struct {
	repo               strategyProfileObjectiveRepositories.StrategyProfileObjectiveRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	strategyProfileSvc StrategyProfileServiceInterface
}

func NewStrategyProfileObjectiveService(
	repo strategyProfileObjectiveRepositories.StrategyProfileObjectiveRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
) StrategyProfileObjectiveServiceInterface {
	return &StrategyProfileObjectiveService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		userSettingSvc:     userSettingSvc,
		strategyProfileSvc: strategyProfileSvc,
	}
}

func (sps *StrategyProfileObjectiveService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileObjectiveListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileObjectiveListDTO, error) {
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileObjectiveService) GetByID(id uint) (*models.StrategyProfileObjective, error) {
	strategyProfileObjective, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileObjective, nil
}

func (sps StrategyProfileObjectiveService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileObjectiveDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileObjectiveDTO strategyDTOs.StrategyProfileObjectiveDetailDTO

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	strategyProfileObjective, err := sps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfileObjective, &strategyProfileObjectiveDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileObjectiveDTO.StrategyProfile.StartDateString = strategyProfileObjective.StrategyProfile.StartDate.Format("2006-01-02")
	strategyProfileObjectiveDTO.StrategyProfile.EndDateString = strategyProfileObjective.StrategyProfile.EndDate.Format("2006-01-02")

	strategyProfileObjectiveDTO.CreatedAt = strategyProfileObjective.CreatedAt
	strategyProfileObjectiveDTO.UpdatedAt = strategyProfileObjective.UpdatedAt

	return &strategyProfileObjectiveDTO, nil
}

func (sps *StrategyProfileObjectiveService) Create(strategyProfileID uint, currentUserID uint, strategyProfileObjectiveFormDTO *strategyDTOs.StrategyProfileObjectiveFormDTO) (uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return 0, err
	}

	tx := sps.repo.BeginTransaction()

	strategyProfileObjective := &models.StrategyProfileObjective{
		StrategyProfileID:          strategyProfileID,
		StrategyProfileFocusAreaID: strategyProfileObjectiveFormDTO.StrategyProfileFocusAreaID,
		StrategyProfileGoalID:      strategyProfileObjectiveFormDTO.StrategyProfileGoalID,
		Title:                      strategyProfileObjectiveFormDTO.Title,
		Description:                strategyProfileObjectiveFormDTO.Description,
		CreatedBYID:                currentUserID,
	}
	id, err := sps.repo.Create(strategyProfileObjective, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileObjective, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			utility.LogError("Failed to fetch strategy profile objective", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileObjective{}.TableName(),
			id,
			existingStrategyProfileObjective,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileObjective{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (sps *StrategyProfileObjectiveService) Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileObjectiveFormDTO *strategyDTOs.StrategyProfileObjectiveFormDTO) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	existingStrategyProfileObjective, err := sps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfileObjective, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStrategyProfileObjective.Title = strategyProfileObjectiveFormDTO.Title
	existingStrategyProfileObjective.Description = strategyProfileObjectiveFormDTO.Description
	existingStrategyProfileObjective.UpdatedBYID = &currentUserID

	err = sps.repo.Update(existingStrategyProfileObjective, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileObjective, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":       oldStrategyProfileObjective.Title,
			"description": oldStrategyProfileObjective.Description,
		}

		newData := map[string]interface{}{
			"title":       existingStrategyProfileObjective.Title,
			"description": existingStrategyProfileObjective.Description,
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
				models.StrategyProfileObjective{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileObjective{}.TableName(),
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

func (sps *StrategyProfileObjectiveService) Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategyProfileObjective []map[string]interface{}

	for _, id := range ids {
		strategyProfileObjective, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStrategyProfileObjective = append(deletedStrategyProfileObjective, map[string]interface{}{
			"id":       id,
			"title_en": strategyProfileObjective.Title,
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
				models.StrategyProfileObjective{}.TableName(),
				successfulIds,
				deletedStrategyProfileObjective,
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
