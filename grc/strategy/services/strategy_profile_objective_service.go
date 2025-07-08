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

	"gorm.io/gorm"
)

type StrategyProfileObjectiveServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileObjectiveListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileObjectiveListDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileObjectiveFormDTO, focusAreaID *uint) ([]uint, error)
	Delete(strategyProfileID uint, objectiveID uint, currentUserID uint) error
}

// Service and implementations
type StrategyProfileObjectiveService struct {
	repo                        strategyProfileObjectiveRepositories.StrategyProfileObjectiveRepositoryInterface
	activityLogSvc              logServices.ActivityLogServiceInterface
	userSvc                     organizationServices.UserServiceInterface
	userSettingSvc              organizationServices.UserSettingServiceInterface
	strategyProfileSvc          StrategyProfileServiceInterface
	strategyProfileFocusAreaSvc StrategyProfileFocusAreaServiceInterface
	strategyProfileGoalSvc      StrategyProfileGoalServiceInterface
	strategicObjectiveSvc       StrategicObjectiveServiceInterface
}

func NewStrategyProfileObjectiveService(
	repo strategyProfileObjectiveRepositories.StrategyProfileObjectiveRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	strategyProfileFocusAreaSvc StrategyProfileFocusAreaServiceInterface,
	strategyProfileGoalSvc StrategyProfileGoalServiceInterface,
	strategicObjectiveSvc StrategicObjectiveServiceInterface,
) StrategyProfileObjectiveServiceInterface {
	return &StrategyProfileObjectiveService{
		repo:                        repo,
		activityLogSvc:              activityLogSvc,
		userSvc:                     userSvc,
		userSettingSvc:              userSettingSvc,
		strategyProfileSvc:          strategyProfileSvc,
		strategyProfileFocusAreaSvc: strategyProfileFocusAreaSvc,
		strategyProfileGoalSvc:      strategyProfileGoalSvc,
		strategicObjectiveSvc:       strategicObjectiveSvc,
	}
}

func (sps *StrategyProfileObjectiveService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileObjectiveListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileObjectiveListDTO, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}
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

func (sps *StrategyProfileObjectiveService) Create(strategyProfileID uint, currentUserID uint, strategyProfileObjectiveFormDTO *strategyDTOs.StrategyProfileObjectiveFormDTO, focusAreaID *uint) ([]uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var createdIds []uint

	for _, strategicObjectiveID := range strategyProfileObjectiveFormDTO.StrategicObjectiveIDs {
		if strategyProfileObjectiveFormDTO.StrategyProfileGoalIDs != nil && len(*strategyProfileObjectiveFormDTO.StrategyProfileGoalIDs) > 0 {
			for _, goalID := range *strategyProfileObjectiveFormDTO.StrategyProfileGoalIDs {
				strategyProfileObjective := &models.StrategyProfileObjective{
					StrategyProfileID:          strategyProfileID,
					StrategyProfileFocusAreaID: focusAreaID,
					StrategyProfileGoalID:      &goalID,
					StrategicObjectiveID:       &strategicObjectiveID,
					CreatedBYID:                currentUserID,
				}
				id, err := sps.repo.Create(strategyProfileObjective, tx)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				createdIds = append(createdIds, id)
			}
		} else {
			strategyProfileObjective := &models.StrategyProfileObjective{
				StrategyProfileID:          strategyProfileID,
				StrategyProfileFocusAreaID: focusAreaID,
				StrategyProfileGoalID:      nil,
				StrategicObjectiveID:       &strategicObjectiveID,
				CreatedBYID:                currentUserID,
			}
			id, err := sps.repo.Create(strategyProfileObjective, tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			createdIds = append(createdIds, id)
		}
	}

	tx.Commit()

	// Log bulk create activities
	go func() {
		idx := 0
		for _, strategicObjectiveID := range strategyProfileObjectiveFormDTO.StrategicObjectiveIDs {
			if strategyProfileObjectiveFormDTO.StrategyProfileGoalIDs != nil && len(*strategyProfileObjectiveFormDTO.StrategyProfileGoalIDs) > 0 {
				for _, goalID := range *strategyProfileObjectiveFormDTO.StrategyProfileGoalIDs {
					id := createdIds[idx]
					err := sps.activityLogSvc.LogActivity(
						models.CreateAction,
						models.StrategyProfileObjective{}.TableName(),
						id,
						map[string]interface{}{
							"strategy_profile_id":            strategyProfileID,
							"strategic_objective_id":         strategicObjectiveID,
							"strategy_profile_goal_id":       goalID,
							"strategy_profile_focus_area_id": focusAreaID,
						},
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
					idx++
				}
			} else {
				id := createdIds[idx]
				err := sps.activityLogSvc.LogActivity(
					models.CreateAction,
					models.StrategyProfileObjective{}.TableName(),
					id,
					map[string]interface{}{
						"strategy_profile_id":            strategyProfileID,
						"strategic_objective_id":         strategicObjectiveID,
						"strategy_profile_goal_id":       nil,
						"strategy_profile_focus_area_id": focusAreaID,
					},
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
				idx++
			}
		}
	}()

	return createdIds, nil
}

func (sps *StrategyProfileObjectiveService) Delete(strategyProfileID uint, id uint, currentUserID uint) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	// Find the strategy profile focus area record
	strategyProfileObjective, err := sps.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Verify it belongs to the correct strategy profile
	if strategyProfileObjective.StrategyProfileID != strategyProfileID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("objective not found in strategy profile"), id)
	}

	// Delete the record
	if err := sps.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the deletion
	go func() {
		err := sps.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.StrategyProfileObjective{}.TableName(),
			id,
			map[string]interface{}{
				"strategy_profile_id":    strategyProfileID,
				"strategic_objective_id": strategyProfileObjective.StrategicObjectiveID,
			},
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.StrategyProfileObjective{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
