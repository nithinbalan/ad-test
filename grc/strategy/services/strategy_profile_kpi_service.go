package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileKpiRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileKpiServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileKpiListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileKpiListDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileKpiFormDTO) ([]uint, error)
	Delete(strategyProfileID uint, kpiID uint, currentUserID uint) error
	GetObjectivesForKPI(kpiID uint, tx utility.DatabaseTransaction) ([]*models.StrategyProfileStrategicObjective, error)
	GetKPIsForObjective(objectiveID uint, tx utility.DatabaseTransaction) ([]*strategyDTOs.KpiRegisterMinimalDetails, error)
}

// Service and implementations
type StrategyProfileKpiService struct {
	repo                                 strategyProfileKpiRepositories.StrategyProfileKpiRepositoryInterface
	activityLogSvc                       logServices.ActivityLogServiceInterface
	userSvc                              organizationServices.UserServiceInterface
	userSettingSvc                       organizationServices.UserSettingServiceInterface
	strategyProfileSvc                   StrategyProfileServiceInterface
	strategyProfileStrategicObjectiveSvc StrategyProfileStrategicObjectiveServiceInterface
	kpiRegisterSvc                       KPIRegisterServiceInterface
}

func NewStrategyProfileKpiService(
	repo strategyProfileKpiRepositories.StrategyProfileKpiRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	strategyProfileStrategicObjectiveSvc StrategyProfileStrategicObjectiveServiceInterface,
	kpiRegisterSvc KPIRegisterServiceInterface,
) StrategyProfileKpiServiceInterface {
	return &StrategyProfileKpiService{
		repo:                                 repo,
		activityLogSvc:                       activityLogSvc,
		userSvc:                              userSvc,
		userSettingSvc:                       userSettingSvc,
		strategyProfileSvc:                   strategyProfileSvc,
		strategyProfileStrategicObjectiveSvc: strategyProfileStrategicObjectiveSvc,
		kpiRegisterSvc:                       kpiRegisterSvc,
	}
}

func (sps *StrategyProfileKpiService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileKpiListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileKpiListDTO, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileKpiService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileKpi, error) {
	strategyProfileKpi, err := sps.repo.GetByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileKpi, nil
}

func (sps *StrategyProfileKpiService) Create(strategyProfileID uint, currentUserID uint, strategyProfileKpiFormDTO *strategyDTOs.StrategyProfileKpiFormDTO) ([]uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var createdIds []uint

	for _, kpiID := range strategyProfileKpiFormDTO.KpiIDs {
		strategyProfileKpi := &models.StrategyProfileKpi{
			StrategyProfileID:                   strategyProfileID,
			StrategyProfileStrategicObjectiveID: strategyProfileKpiFormDTO.StrategyProfileStrategicObjectiveID,
			KpiID:                               &kpiID,
			Weightage:                           strategyProfileKpiFormDTO.Weightage,
			CreatedBYID:                         currentUserID,
		}
		id, err := sps.repo.Create(strategyProfileKpi, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		createdIds = append(createdIds, id)
	}

	tx.Commit()

	// Log bulk create activities
	go func() {
		for idx, kpiID := range strategyProfileKpiFormDTO.KpiIDs {
			id := createdIds[idx]
			err := sps.activityLogSvc.LogActivity(
				models.CreateAction,
				models.StrategyProfileKpi{}.TableName(),
				id,
				map[string]interface{}{
					"strategy_profile_id":    strategyProfileID,
					"strategic_objective_id": strategyProfileKpiFormDTO.StrategyProfileStrategicObjectiveID,
					"kpi_id":                 kpiID,
					"weightage":              strategyProfileKpiFormDTO.Weightage,
				},
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.StrategyProfileKpi{}.TableName(),
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

func (sps *StrategyProfileKpiService) Delete(strategyProfileID uint, id uint, currentUserID uint) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	// Find the strategy profile focus area record
	strategyProfileKpi, err := sps.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Verify it belongs to the correct strategy profile
	if strategyProfileKpi.StrategyProfileID != strategyProfileID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("kpi not found in strategy profile"), id)
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
			models.StrategyProfileKpi{}.TableName(),
			id,
			map[string]interface{}{
				"strategy_profile_id":    strategyProfileID,
				"strategic_objective_id": strategyProfileKpi.KpiID,
				"weightage":              strategyProfileKpi.Weightage,
			},
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.StrategyProfileKpi{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (sps *StrategyProfileKpiService) GetObjectivesForKPI(kpiID uint, tx utility.DatabaseTransaction) ([]*models.StrategyProfileStrategicObjective, error) {
	return sps.repo.GetObjectivesForKPI(kpiID, tx)
}

func (sps *StrategyProfileKpiService) GetKPIsForObjective(objectiveID uint, tx utility.DatabaseTransaction) ([]*strategyDTOs.KpiRegisterMinimalDetails, error) {
	return sps.repo.GetKPIsForObjective(objectiveID, tx)
}
