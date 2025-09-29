package services

import (
	"errors"
	"fmt"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	projectActionPlanServices "jethur/grc/project/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileProjectActionPlanRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileProjectActionPlanServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectActionPlanListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectActionPlan, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileProjectActionPlanDetailDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileProjectActionPlanFormDTO) ([]uint, error)
	Delete(strategyProfileID uint, kpiID uint, currentUserID uint) error
}

// Service and implementations
type StrategyProfileProjectActionPlanService struct {
	repo                               strategyProfileProjectActionPlanRepositories.StrategyProfileProjectActionPlanRepositoryInterface
	activityLogSvc                     logServices.ActivityLogServiceInterface
	userSvc                            organizationServices.UserServiceInterface
	userSettingSvc                     organizationServices.UserSettingServiceInterface
	strategyProfileSvc                 StrategyProfileServiceInterface
	strategyProfileProjectSvc          StrategyProfileProjectServiceInterface
	strategyProfileProjectMilestoneSvc StrategyProfileProjectMilestoneServiceInterface
	projectActionPlanSvc               projectActionPlanServices.ProjectActionPlanServiceInterface
}

func NewStrategyProfileProjectActionPlanService(
	repo strategyProfileProjectActionPlanRepositories.StrategyProfileProjectActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	strategyProfileProjectSvc StrategyProfileProjectServiceInterface,
	strategyProfileProjectMilestoneSvc StrategyProfileProjectMilestoneServiceInterface,
	projectActionPlanSvc projectActionPlanServices.ProjectActionPlanServiceInterface,
) StrategyProfileProjectActionPlanServiceInterface {
	return &StrategyProfileProjectActionPlanService{
		repo:                               repo,
		activityLogSvc:                     activityLogSvc,
		userSvc:                            userSvc,
		userSettingSvc:                     userSettingSvc,
		strategyProfileSvc:                 strategyProfileSvc,
		strategyProfileProjectSvc:          strategyProfileProjectSvc,
		strategyProfileProjectMilestoneSvc: strategyProfileProjectMilestoneSvc,
		projectActionPlanSvc:               projectActionPlanSvc,
	}
}

func (spap *StrategyProfileProjectActionPlanService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectActionPlanListDTO, error) {
	if _, err := spap.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}
	return spap.repo.GetAll(strategyProfileID, params)
}

func (spap *StrategyProfileProjectActionPlanService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectActionPlan, error) {
	strategyProfileProjectActionPlan, err := spap.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileProjectActionPlan, nil
}

func (ssp *StrategyProfileProjectActionPlanService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileProjectActionPlanDetailDTO, error) {
	var strategyProfileProjectActionPlanDTO strategyDTOs.StrategyProfileProjectActionPlanDetailDTO

	// Check if the PROJECT register exists
	if _, err := ssp.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	// Get full details of the PROJECT action plan
	strategyProfileProjectActionPlan, err := ssp.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the PROJECT action plan into the DTO
	err = mapstructure.Decode(strategyProfileProjectActionPlan, &strategyProfileProjectActionPlanDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileProjectActionPlanDTO.CreatedAt = strategyProfileProjectActionPlan.CreatedAt
	strategyProfileProjectActionPlanDTO.UpdatedAt = strategyProfileProjectActionPlan.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &strategyProfileProjectActionPlanDTO, nil
}

func (spap *StrategyProfileProjectActionPlanService) Create(strategyProfileID uint, currentUserID uint, strategyProfileProjectActionPlanFormDTO *strategyDTOs.StrategyProfileProjectActionPlanFormDTO) ([]uint, error) {
	if _, err := spap.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := spap.repo.BeginTransaction()
	var createdIds []uint

	for _, projectActionPlanID := range strategyProfileProjectActionPlanFormDTO.ProjectActionPlanIDs {
		strategyProfileProjectActionPlan := &models.StrategyProfileProjectActionPlan{
			StrategyProfileID:                 strategyProfileID,
			StrategyProfileProjectID:          strategyProfileProjectActionPlanFormDTO.StrategyProfileProjectID,
			StrategyProfileProjectMilestoneID: strategyProfileProjectActionPlanFormDTO.StrategyProfileProjectMilestoneID,
			ProjectActionPlanID:               &projectActionPlanID,
			CreatedBYID:                       currentUserID,
		}

		id, err := spap.repo.Create(strategyProfileProjectActionPlan, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		createdIds = append(createdIds, id)
	}

	tx.Commit()

	// Log bulk create activities
	go func() {
		for idx, projectActionPlanID := range strategyProfileProjectActionPlanFormDTO.ProjectActionPlanIDs {
			id := createdIds[idx]
			err := spap.activityLogSvc.LogActivity(
				models.CreateAction,
				models.StrategyProfileProjectActionPlan{}.TableName(),
				id,
				map[string]interface{}{
					"strategy_profile_id":                   strategyProfileID,
					"strategy_profile_project_id":           strategyProfileProjectActionPlanFormDTO.StrategyProfileProjectID,
					"strategy_profile_project_milestone_id": strategyProfileProjectActionPlanFormDTO.StrategyProfileProjectMilestoneID,
					"project_action_plan_id":                projectActionPlanID,
				},
				currentUserID,
				"strategy", "strategy_profile", "strategy_profile_action_plan", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.StrategyProfileProjectActionPlan{}.TableName(),
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

func (spap *StrategyProfileProjectActionPlanService) Delete(strategyProfileID uint, id uint, currentUserID uint) error {
	if _, err := spap.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := spap.repo.BeginTransaction()

	// Find the strategy profile focus area record
	strategyProfileProjectActionPlan, err := spap.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Verify it belongs to the correct strategy profile
	if strategyProfileProjectActionPlan.StrategyProfileID != strategyProfileID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("project action plan not found in strategy profile"), id)
	}

	// Delete the record
	if err := spap.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the deletion
	go func() {
		err := spap.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.StrategyProfileProjectActionPlan{}.TableName(),
			id,
			map[string]interface{}{
				"strategy_profile_id":         strategyProfileID,
				"strategy_profile_project_id": strategyProfileProjectActionPlan.StrategyProfileProjectID,
				"project_action_plan_id":      strategyProfileProjectActionPlan.ProjectActionPlanID,
			},
			currentUserID,
			"strategy", "strategy_profile", "strategy_profile_action_plan", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.StrategyProfileProjectActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
