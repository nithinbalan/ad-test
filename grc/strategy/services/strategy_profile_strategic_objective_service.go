package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileStrategicObjectiveRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileStrategicObjectiveServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileStrategicObjectiveListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileStrategicObjectiveListDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileStrategicObjectiveFormDTO) ([]uint, error)
	Delete(strategyProfileID uint, objectiveID uint, currentUserID uint) error
}

// Service and implementations
type StrategyProfileStrategicObjectiveService struct {
	repo                        strategyProfileStrategicObjectiveRepositories.StrategyProfileStrategicObjectiveRepositoryInterface
	activityLogSvc              logServices.ActivityLogServiceInterface
	userSvc                     organizationServices.UserServiceInterface
	userSettingSvc              organizationServices.UserSettingServiceInterface
	strategyProfileSvc          StrategyProfileServiceInterface
	strategyProfileFocusAreaSvc StrategyProfileFocusAreaServiceInterface
	strategyProfileGoalSvc      StrategyProfileGoalServiceInterface
	strategicObjectiveSvc       StrategicObjectiveServiceInterface
}

func NewStrategyProfileStrategicObjectiveService(
	repo strategyProfileStrategicObjectiveRepositories.StrategyProfileStrategicObjectiveRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	strategyProfileFocusAreaSvc StrategyProfileFocusAreaServiceInterface,
	strategyProfileGoalSvc StrategyProfileGoalServiceInterface,
	strategicObjectiveSvc StrategicObjectiveServiceInterface,
) StrategyProfileStrategicObjectiveServiceInterface {
	return &StrategyProfileStrategicObjectiveService{
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

func (sps *StrategyProfileStrategicObjectiveService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileStrategicObjectiveListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileStrategicObjectiveListDTO, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileStrategicObjectiveService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileStrategicObjective, error) {
	strategyProfileStrategicObjective, err := sps.repo.GetByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileStrategicObjective, nil
}

func (sps *StrategyProfileStrategicObjectiveService) Create(strategyProfileID uint, currentUserID uint, strategyProfileStrategicObjectiveFormDTO *strategyDTOs.StrategyProfileStrategicObjectiveFormDTO) ([]uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var createdIds []uint

	for _, strategicObjectiveID := range strategyProfileStrategicObjectiveFormDTO.StrategicObjectiveIDs {
		if strategyProfileStrategicObjectiveFormDTO.StrategyProfileGoalIDs != nil && len(*strategyProfileStrategicObjectiveFormDTO.StrategyProfileGoalIDs) > 0 {
			for _, goalID := range *strategyProfileStrategicObjectiveFormDTO.StrategyProfileGoalIDs {
				strategyProfileStrategicObjective := &models.StrategyProfileStrategicObjective{
					StrategyProfileID:          strategyProfileID,
					StrategyProfileFocusAreaID: strategyProfileStrategicObjectiveFormDTO.StrategyProfileFocusAreaID,
					StrategyProfileGoalID:      &goalID,
					StrategicObjectiveID:       &strategicObjectiveID,
					Score:                      strategyProfileStrategicObjectiveFormDTO.Score,
					CreatedBYID:                currentUserID,
				}
				id, err := sps.repo.Create(strategyProfileStrategicObjective, tx)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				createdIds = append(createdIds, id)
			}
		} else {
			strategyProfileStrategicObjective := &models.StrategyProfileStrategicObjective{
				StrategyProfileID:          strategyProfileID,
				StrategyProfileFocusAreaID: strategyProfileStrategicObjectiveFormDTO.StrategyProfileFocusAreaID,
				StrategyProfileGoalID:      nil,
				StrategicObjectiveID:       &strategicObjectiveID,
				Score:                      strategyProfileStrategicObjectiveFormDTO.Score,
				CreatedBYID:                currentUserID,
			}
			id, err := sps.repo.Create(strategyProfileStrategicObjective, tx)
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
		for _, strategicObjectiveID := range strategyProfileStrategicObjectiveFormDTO.StrategicObjectiveIDs {
			if strategyProfileStrategicObjectiveFormDTO.StrategyProfileGoalIDs != nil && len(*strategyProfileStrategicObjectiveFormDTO.StrategyProfileGoalIDs) > 0 {
				for _, goalID := range *strategyProfileStrategicObjectiveFormDTO.StrategyProfileGoalIDs {
					id := createdIds[idx]
					err := sps.activityLogSvc.LogActivity(
						models.CreateAction,
						models.StrategyProfileStrategicObjective{}.TableName(),
						id,
						map[string]interface{}{
							"strategy_profile_id":            strategyProfileID,
							"strategic_objective_id":         strategicObjectiveID,
							"strategy_profile_goal_id":       goalID,
							"strategy_profile_focus_area_id": strategyProfileStrategicObjectiveFormDTO.StrategyProfileFocusAreaID,
							"score":                          strategyProfileStrategicObjectiveFormDTO.Score,
						},
						currentUserID,
						nil,
					)
					if err != nil {
						fields := map[string]interface{}{
							"action":    models.CreateAction,
							"tableName": models.StrategyProfileStrategicObjective{}.TableName(),
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
					models.StrategyProfileStrategicObjective{}.TableName(),
					id,
					map[string]interface{}{
						"strategy_profile_id":            strategyProfileID,
						"strategic_objective_id":         strategicObjectiveID,
						"strategy_profile_goal_id":       nil,
						"strategy_profile_focus_area_id": strategyProfileStrategicObjectiveFormDTO.StrategyProfileFocusAreaID,
						"score":                          strategyProfileStrategicObjectiveFormDTO.Score,
					},
					currentUserID,
					nil,
				)
				if err != nil {
					fields := map[string]interface{}{
						"action":    models.CreateAction,
						"tableName": models.StrategyProfileStrategicObjective{}.TableName(),
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

func (sps *StrategyProfileStrategicObjectiveService) Delete(strategyProfileID uint, id uint, currentUserID uint) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	// Find the strategy profile focus area record
	strategyProfileStrategicObjective, err := sps.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Verify it belongs to the correct strategy profile
	if strategyProfileStrategicObjective.StrategyProfileID != strategyProfileID {
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
			models.StrategyProfileStrategicObjective{}.TableName(),
			id,
			map[string]interface{}{
				"strategy_profile_id":    strategyProfileID,
				"strategic_objective_id": strategyProfileStrategicObjective.StrategicObjectiveID,
				"score":                  strategyProfileStrategicObjective.Score,
			},
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.StrategyProfileStrategicObjective{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
