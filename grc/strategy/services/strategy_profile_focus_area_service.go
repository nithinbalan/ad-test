package services

import (
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileFocusAreaRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"
)

type StrategyProfileFocusAreaServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileFocusAreaListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileFocusAreaListDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileFocusAreaFormDTO) ([]uint, error)
	Delete(strategyProfileID uint, focusAreaID uint, currentUserID uint) error
}

// Service and implementations
type StrategyProfileFocusAreaService struct {
	repo               strategyProfileFocusAreaRepositories.StrategyProfileFocusAreaRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	strategyProfileSvc StrategyProfileServiceInterface
	focusAreaSvc       FocusAreaServiceInterface
}

func NewStrategyProfileFocusAreaService(
	repo strategyProfileFocusAreaRepositories.StrategyProfileFocusAreaRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	focusAreaSvc FocusAreaServiceInterface,
) StrategyProfileFocusAreaServiceInterface {
	return &StrategyProfileFocusAreaService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		userSettingSvc:     userSettingSvc,
		strategyProfileSvc: strategyProfileSvc,
		focusAreaSvc:       focusAreaSvc,
	}
}

func (sps *StrategyProfileFocusAreaService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileFocusAreaListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileFocusAreaListDTO, error) {
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps *StrategyProfileFocusAreaService) Create(strategyProfileID uint, currentUserID uint, strategyProfileFocusAreaFormDTO *strategyDTOs.StrategyProfileFocusAreaFormDTO) ([]uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var createdIds []uint

	for _, focusAreaID := range strategyProfileFocusAreaFormDTO.FocusAreaIDs {
		strategyProfileFocusArea := &models.StrategyProfileFocusArea{
			StrategyProfileID: strategyProfileID,
			FocusAreaID:       &focusAreaID,
			CreatedBYID:       currentUserID,
		}
		id, err := sps.repo.Create(strategyProfileFocusArea, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		createdIds = append(createdIds, id)
	}

	tx.Commit()

	// Log bulk create activities
	go func() {
		for i, id := range createdIds {
			err := sps.activityLogSvc.LogActivity(
				models.CreateAction,
				models.StrategyProfileFocusArea{}.TableName(),
				id,
				map[string]interface{}{
					"strategy_profile_id": strategyProfileID,
					"focus_area_id":       strategyProfileFocusAreaFormDTO.FocusAreaIDs[i],
				},
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.StrategyProfileFocusArea{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return createdIds, nil
}

func (sps *StrategyProfileFocusAreaService) Delete(strategyProfileID uint, focusAreaID uint, currentUserID uint) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	// Find the strategy profile focus area record
	strategyProfileFocusArea, err := sps.repo.GetByID(focusAreaID)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, focusAreaID)
	}

	// Verify it belongs to the correct strategy profile
	if strategyProfileFocusArea.StrategyProfileID != strategyProfileID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("focus area not found in strategy profile"), focusAreaID)
	}

	// Delete the record
	if err := sps.repo.Delete(focusAreaID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the deletion
	go func() {
		err := sps.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.StrategyProfileFocusArea{}.TableName(),
			focusAreaID,
			map[string]interface{}{
				"strategy_profile_id": strategyProfileID,
				"focus_area_id":       strategyProfileFocusArea.FocusAreaID,
			},
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.StrategyProfileFocusArea{}.TableName(),
				"recordID":  focusAreaID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
