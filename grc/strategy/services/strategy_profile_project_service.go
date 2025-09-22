package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	projectServices "jethur/grc/project/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileProjectRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileProjectServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProject, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileProjectDetailDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileProjectFormDTO) ([]uint, error)
	Delete(strategyProfileID uint, projectID uint, currentUserID uint) error
}

// Service and implementations
type StrategyProfileProjectService struct {
	repo                                 strategyProfileProjectRepositories.StrategyProfileProjectRepositoryInterface
	activityLogSvc                       logServices.ActivityLogServiceInterface
	userSvc                              organizationServices.UserServiceInterface
	userSettingSvc                       organizationServices.UserSettingServiceInterface
	strategyProfileSvc                   StrategyProfileServiceInterface
	projectSvc                           projectServices.ProjectServiceInterface
	strategyProfileStrategicObjectiveSvc StrategyProfileStrategicObjectiveServiceInterface
	strategyProfileStrategicIntiativeSvc StrategyProfileStrategicIntiativeServiceInterface
}

func NewStrategyProfileProjectService(
	repo strategyProfileProjectRepositories.StrategyProfileProjectRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	projectSvc projectServices.ProjectServiceInterface,
	strategyProfileStrategicObjectiveSvc StrategyProfileStrategicObjectiveServiceInterface,
	strategyProfileStrategicIntiativeSvc StrategyProfileStrategicIntiativeServiceInterface,
) StrategyProfileProjectServiceInterface {
	return &StrategyProfileProjectService{
		repo:                                 repo,
		activityLogSvc:                       activityLogSvc,
		userSvc:                              userSvc,
		userSettingSvc:                       userSettingSvc,
		strategyProfileSvc:                   strategyProfileSvc,
		projectSvc:                           projectSvc,
		strategyProfileStrategicObjectiveSvc: strategyProfileStrategicObjectiveSvc,
		strategyProfileStrategicIntiativeSvc: strategyProfileStrategicIntiativeSvc,
	}
}

func (ssp *StrategyProfileProjectService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectListDTO, error) {
	if _, err := ssp.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}
	return ssp.repo.GetAll(strategyProfileID, params)
}

func (ssp *StrategyProfileProjectService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProject, error) {
	strategyProfileProject, err := ssp.repo.GetByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileProject, nil
}

func (ssp *StrategyProfileProjectService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileProjectDetailDTO, error) {
	var strategyProfileProjectDTO strategyDTOs.StrategyProfileProjectDetailDTO

	// Check if the PROJECT register exists
	if _, err := ssp.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	// Get full details of the PROJECT action plan
	strategyProfileProject, err := ssp.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the PROJECT action plan into the DTO
	err = mapstructure.Decode(strategyProfileProject, &strategyProfileProjectDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileProjectDTO.CreatedAt = strategyProfileProject.CreatedAt
	strategyProfileProjectDTO.UpdatedAt = strategyProfileProject.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &strategyProfileProjectDTO, nil
}

func (ssp *StrategyProfileProjectService) Create(strategyProfileID uint, currentUserID uint, strategyProfileProjectFormDTO *strategyDTOs.StrategyProfileProjectFormDTO) ([]uint, error) {

	if _, err := ssp.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := ssp.repo.BeginTransaction()
	var createdIds []uint

	for _, projectID := range strategyProfileProjectFormDTO.ProjectIDs {
		strategyProfileProject := &models.StrategyProfileProject{
			StrategyProfileID:                   strategyProfileID,
			ProjectID:                           &projectID,
			StrategyProfileStrategicObjectiveID: strategyProfileProjectFormDTO.StrategyProfileStrategicObjectiveID,
			StrategyProfileStrategicIntiativeID: strategyProfileProjectFormDTO.StrategyProfileStrategicIntiativeID,
			CreatedBYID:                         currentUserID,
		}
		id, err := ssp.repo.Create(strategyProfileProject, tx)
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
			err := ssp.activityLogSvc.LogActivity(
				models.CreateAction,
				models.StrategyProfileProject{}.TableName(),
				id,
				map[string]interface{}{
					"strategy_profile_id": strategyProfileID,
					"project_id":          strategyProfileProjectFormDTO.ProjectIDs[i],
					"strategy_profile_strategic_objective_id": strategyProfileProjectFormDTO.StrategyProfileStrategicObjectiveID,
					"strategy_profile_strategic_intiative_id": strategyProfileProjectFormDTO.StrategyProfileStrategicIntiativeID,
				},
				currentUserID,
				"", "", "", nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.StrategyProfileProject{}.TableName(),
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

func (ssp *StrategyProfileProjectService) Delete(strategyProfileID uint, id uint, currentUserID uint) error {
	if _, err := ssp.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := ssp.repo.BeginTransaction()

	// Find the strategy profile focus area record
	strategyProfileProject, err := ssp.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Verify it belongs to the correct strategy profile
	if strategyProfileProject.StrategyProfileID != strategyProfileID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("project not found in strategy profile"), id)
	}

	// Delete the record
	if err := ssp.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the deletion
	go func() {
		err := ssp.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.StrategyProfileProject{}.TableName(),
			id,
			map[string]interface{}{
				"strategy_profile_id": strategyProfileID,
				"project_id":          strategyProfileProject.ProjectID,
				"strategy_profile_strategic_objective_id": strategyProfileProject.StrategyProfileStrategicObjectiveID,
				"strategy_profile_strategic_intiative_id": strategyProfileProject.StrategyProfileStrategicIntiativeID,
			},
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.StrategyProfileProject{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
