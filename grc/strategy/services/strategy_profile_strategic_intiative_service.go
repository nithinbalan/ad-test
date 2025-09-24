package services

import (
	"errors"
	"fmt"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileStrategicIntiativeRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileStrategicIntiativeServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileStrategicIntiativeListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileStrategicIntiativeListDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileStrategicIntiativeFormDTO) ([]uint, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileStrategicIntiativeDetailDTO, error)
	Delete(strategyProfileID uint, kpiID uint, currentUserID uint) error
}

// Service and implementations
type StrategyProfileStrategicIntiativeService struct {
	repo                                 strategyProfileStrategicIntiativeRepositories.StrategyProfileStrategicIntiativeRepositoryInterface
	activityLogSvc                       logServices.ActivityLogServiceInterface
	userSvc                              organizationServices.UserServiceInterface
	userSettingSvc                       organizationServices.UserSettingServiceInterface
	strategyProfileSvc                   StrategyProfileServiceInterface
	strategyProfileStrategicObjectiveSvc StrategyProfileStrategicObjectiveServiceInterface
	strategicIntiativeSvc                StrategicIntiativeServiceInterface
}

func NewStrategyProfileStrategicIntiativeService(
	repo strategyProfileStrategicIntiativeRepositories.StrategyProfileStrategicIntiativeRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	strategyProfileStrategicObjectiveSvc StrategyProfileStrategicObjectiveServiceInterface,
	strategicIntiativeSvc StrategicIntiativeServiceInterface,
) StrategyProfileStrategicIntiativeServiceInterface {
	return &StrategyProfileStrategicIntiativeService{
		repo:                                 repo,
		activityLogSvc:                       activityLogSvc,
		userSvc:                              userSvc,
		userSettingSvc:                       userSettingSvc,
		strategyProfileSvc:                   strategyProfileSvc,
		strategyProfileStrategicObjectiveSvc: strategyProfileStrategicObjectiveSvc,
		strategicIntiativeSvc:                strategicIntiativeSvc,
	}
}

func (sps *StrategyProfileStrategicIntiativeService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileStrategicIntiativeListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileStrategicIntiativeListDTO, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileStrategicIntiativeService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileStrategicIntiative, error) {
	strategyProfileKpi, err := sps.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileKpi, nil
}

func (spsi *StrategyProfileStrategicIntiativeService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileStrategicIntiativeDetailDTO, error) {
	var strategyProfileStrategicIntiativeDTO strategyDTOs.StrategyProfileStrategicIntiativeDetailDTO

	// Check if the PROJECT register exists
	if _, err := spsi.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	// Get full details of the PROJECT action plan
	strategyProfileStrategicIntiative, err := spsi.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the PROJECT action plan into the DTO
	err = mapstructure.Decode(strategyProfileStrategicIntiative, &strategyProfileStrategicIntiativeDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileStrategicIntiativeDTO.CreatedAt = strategyProfileStrategicIntiative.CreatedAt
	strategyProfileStrategicIntiativeDTO.UpdatedAt = strategyProfileStrategicIntiative.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &strategyProfileStrategicIntiativeDTO, nil
}

func (spsi *StrategyProfileStrategicIntiativeService) Create(strategyProfileID uint, currentUserID uint, strategyProfileStrategicIntiativeFormDTO *strategyDTOs.StrategyProfileStrategicIntiativeFormDTO) ([]uint, error) {
	if _, err := spsi.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := spsi.repo.BeginTransaction()
	var createdIds []uint

	for _, strategicIntiativeID := range strategyProfileStrategicIntiativeFormDTO.StrategicIntiativeIDs {
		strategyProfileStrategicIntiative := &models.StrategyProfileStrategicIntiative{
			StrategyProfileID:                   strategyProfileID,
			StrategyProfileStrategicObjectiveID: strategyProfileStrategicIntiativeFormDTO.StrategyProfileStrategicObjectiveID,
			StrategicIntiativeID:                &strategicIntiativeID,
			CreatedBYID:                         currentUserID,
		}

		id, err := spsi.repo.Create(strategyProfileStrategicIntiative, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		createdIds = append(createdIds, id)
	}

	tx.Commit()

	// Log bulk create activities
	go func() {
		for idx, strategicIntiativeID := range strategyProfileStrategicIntiativeFormDTO.StrategicIntiativeIDs {
			id := createdIds[idx]
			err := spsi.activityLogSvc.LogActivity(
				models.CreateAction,
				models.StrategyProfileStrategicIntiative{}.TableName(),
				id,
				map[string]interface{}{
					"strategy_profile_id":                     strategyProfileID,
					"strategy_profile_strategic_objective_id": strategyProfileStrategicIntiativeFormDTO.StrategyProfileStrategicObjectiveID,
					"strategic_intiative_id":                  strategicIntiativeID,
				},
				currentUserID,
				"strategy", "strategy_profile", "strategic_profile_strategic_objective", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.StrategyProfileStrategicIntiative{}.TableName(),
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

func (spsi *StrategyProfileStrategicIntiativeService) Delete(strategyProfileID uint, id uint, currentUserID uint) error {
	if _, err := spsi.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := spsi.repo.BeginTransaction()

	// Find the strategy profile focus area record
	strategyProfileStrategicIntiative, err := spsi.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Verify it belongs to the correct strategy profile
	if strategyProfileStrategicIntiative.StrategyProfileID != strategyProfileID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("strategic intiative not found in strategy profile"), id)
	}

	// Delete the record
	if err := spsi.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the deletion
	go func() {
		err := spsi.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.StrategyProfileStrategicIntiative{}.TableName(),
			id,
			map[string]interface{}{
				"strategy_profile_id":                     strategyProfileID,
				"strategy_profile_strategic_objective_id": strategyProfileStrategicIntiative.StrategyProfileStrategicObjectiveID,
				"strategic_intiative_id":                  strategyProfileStrategicIntiative.StrategicIntiativeID,
			},
			currentUserID,
			"strategy", "strategy_profile", "strategic_profile_strategic_objective", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.StrategyProfileStrategicIntiative{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
