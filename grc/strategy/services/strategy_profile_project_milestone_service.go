package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	projectMilestoneServices "jethur/grc/project/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileProjectMilestoneRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileProjectMilestoneServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectMilestoneListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectMilestoneListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectMilestone, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileProjectMilestoneDetailDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileProjectMilestoneFormDTO) ([]uint, error)
	Delete(strategyProfileID uint, kpiID uint, currentUserID uint) error
}

// Service and implementations
type StrategyProfileProjectMilestoneService struct {
	repo                      strategyProfileProjectMilestoneRepositories.StrategyProfileProjectMilestoneRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	userSettingSvc            organizationServices.UserSettingServiceInterface
	strategyProfileSvc        StrategyProfileServiceInterface
	strategyProfileProjectSvc StrategyProfileProjectServiceInterface
	projectMilestoneSvc       projectMilestoneServices.ProjectMilestoneServiceInterface
}

func NewStrategyProfileProjectMilestoneService(
	repo strategyProfileProjectMilestoneRepositories.StrategyProfileProjectMilestoneRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	strategyProfileProjectSvc StrategyProfileProjectServiceInterface,
	projectMilestoneSvc projectMilestoneServices.ProjectMilestoneServiceInterface,
) StrategyProfileProjectMilestoneServiceInterface {
	return &StrategyProfileProjectMilestoneService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		userSvc:                   userSvc,
		userSettingSvc:            userSettingSvc,
		strategyProfileSvc:        strategyProfileSvc,
		strategyProfileProjectSvc: strategyProfileProjectSvc,
		projectMilestoneSvc:       projectMilestoneSvc,
	}
}

func (sppm *StrategyProfileProjectMilestoneService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectMilestoneListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectMilestoneListDTO, error) {
	if _, err := sppm.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}
	return sppm.repo.GetAll(strategyProfileID, params)
}

func (sppm *StrategyProfileProjectMilestoneService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectMilestone, error) {
	strategyProfileProjectMilestone, err := sppm.repo.GetByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileProjectMilestone, nil
}

func (ssp *StrategyProfileProjectMilestoneService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileProjectMilestoneDetailDTO, error) {
	var strategyProfileProjectMilestoneDTO strategyDTOs.StrategyProfileProjectMilestoneDetailDTO

	// Check if the PROJECT register exists
	if _, err := ssp.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	// Get full details of the PROJECT action plan
	strategyProfileProjectMilestone, err := ssp.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the PROJECT action plan into the DTO
	err = mapstructure.Decode(strategyProfileProjectMilestone, &strategyProfileProjectMilestoneDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileProjectMilestoneDTO.CreatedAt = strategyProfileProjectMilestone.CreatedAt
	strategyProfileProjectMilestoneDTO.UpdatedAt = strategyProfileProjectMilestone.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &strategyProfileProjectMilestoneDTO, nil
}

func (sppm *StrategyProfileProjectMilestoneService) Create(strategyProfileID uint, currentUserID uint, strategyProfileProjectMilestoneFormDTO *strategyDTOs.StrategyProfileProjectMilestoneFormDTO) ([]uint, error) {

	if _, err := sppm.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sppm.repo.BeginTransaction()
	var createdIds []uint

	for _, projectMilestoneID := range strategyProfileProjectMilestoneFormDTO.ProjectMilestoneIDs {
		strategyProfileProjectMilestone := &models.StrategyProfileProjectMilestone{
			StrategyProfileID:        strategyProfileID,
			StrategyProfileProjectID: strategyProfileProjectMilestoneFormDTO.StrategyProfileProjectID,
			ProjectMilestoneID:       &projectMilestoneID,
			CreatedBYID:              currentUserID,
		}

		id, err := sppm.repo.Create(strategyProfileProjectMilestone, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		createdIds = append(createdIds, id)
	}

	tx.Commit()

	// Log bulk create activities
	go func() {
		for idx, projectMilestoneID := range strategyProfileProjectMilestoneFormDTO.ProjectMilestoneIDs {
			id := createdIds[idx]
			err := sppm.activityLogSvc.LogActivity(
				models.CreateAction,
				models.StrategyProfileProjectMilestone{}.TableName(),
				id,
				map[string]interface{}{
					"strategy_profile_id":         strategyProfileID,
					"strategy_profile_project_id": strategyProfileProjectMilestoneFormDTO.StrategyProfileProjectID,
					"project_milestone_id":        projectMilestoneID,
				},
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.StrategyProfileProjectMilestone{}.TableName(),
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

func (sppm *StrategyProfileProjectMilestoneService) Delete(strategyProfileID uint, id uint, currentUserID uint) error {
	if _, err := sppm.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sppm.repo.BeginTransaction()

	// Find the strategy profile focus area record
	strategyProfileProjectMilestone, err := sppm.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Verify it belongs to the correct strategy profile
	if strategyProfileProjectMilestone.StrategyProfileID != strategyProfileID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("project milestone not found in strategy profile"), id)
	}

	// Delete the record
	if err := sppm.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the deletion
	go func() {
		err := sppm.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.StrategyProfileProjectMilestone{}.TableName(),
			id,
			map[string]interface{}{
				"strategy_profile_id":         strategyProfileID,
				"strategy_profile_project_id": strategyProfileProjectMilestone.StrategyProfileProjectID,
				"project_milestone_id":        strategyProfileProjectMilestone.ProjectMilestoneID,
			},
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.StrategyProfileProjectMilestone{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
