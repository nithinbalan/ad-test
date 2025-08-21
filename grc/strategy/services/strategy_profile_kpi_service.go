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

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileKpiServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileKpiListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileKpiListDTO, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileFocusAreaDetailDTO, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileKpiFormDTO) ([]uint, error)
	Delete(strategyProfileID uint, kpiID uint, currentUserID uint) error
	GetObjectivesForKPI(kpiID uint, tx utility.DatabaseTransaction) ([]*models.StrategyProfileStrategicObjective, error)
	GetKPIsForObjective(objectiveID uint, tx utility.DatabaseTransaction) ([]*strategyDTOs.KpiRegisterMinimalDetails, error)
	UpdateWeightage(id uint, currentUserID uint, kpi *strategyDTOs.StrategyProfileKpiWeightageUpdateDTO) error
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

func (sps *StrategyProfileKpiService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileFocusAreaDetailDTO, error) {
	var strategyProfileFocusAreaDTO strategyDTOs.StrategyProfileFocusAreaDetailDTO

	// Check if the PROJECT register exists
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	// Get full details of the PROJECT action plan
	strategyProfileFocusArea, err := sps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the PROJECT action plan into the DTO
	err = mapstructure.Decode(strategyProfileFocusArea, &strategyProfileFocusAreaDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileFocusAreaDTO.CreatedAt = strategyProfileFocusArea.CreatedAt
	strategyProfileFocusAreaDTO.UpdatedAt = strategyProfileFocusArea.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &strategyProfileFocusAreaDTO, nil
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

func (sps *StrategyProfileKpiService) UpdateWeightage(id uint, currentUserID uint, kpiWeightageUpdateDTO *strategyDTOs.StrategyProfileKpiWeightageUpdateDTO) error {
	tx := sps.repo.BeginTransaction()

	// Get the KPI to update
	kpi, err := sps.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldData := map[string]interface{}{
		"weightage": kpi.Weightage,
	}

	// Validate the weightage
	if err := sps.validateKPIWeightage(kpi, kpiWeightageUpdateDTO.Weightage, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Update the weightage
	kpi.Weightage = kpiWeightageUpdateDTO.Weightage
	kpi.UpdatedBYID = &currentUserID

	// Save the changes
	if err := sps.repo.UpdateWeightage(kpi, tx); err != nil {
		tx.Rollback()
		return err
	}

	newData := map[string]interface{}{
		"weightage": kpi.Weightage,
	}

	// Commit first before logging
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	// Log the activity asynchronously
	go func() {
		err := sps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.StrategyProfileKpi{}.TableName(),
			kpi.ID,
			map[string]interface{}{
				"old": oldData,
				"new": newData,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.StrategyProfileKpi{}.TableName(),
				"recordID":  kpi.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (sps *StrategyProfileKpiService) validateKPIWeightage(kpi *models.StrategyProfileKpi, newWeightage *float64, tx utility.DatabaseTransaction) error {
	// Validate weightage is not nil
	if newWeightage == nil {
		return utility.NewValidationError(
			fmt.Errorf("weightage is required"),
			"required",
			"weightage",
		)
	}

	// Skip total weightage validation if no objective is set
	if kpi.StrategyProfileStrategicObjectiveID == nil {
		return nil
	}

	// Get all KPIs for the same objective
	kpis, err := sps.GetKPIsForObjective(*kpi.StrategyProfileStrategicObjectiveID, tx)
	if err != nil {
		return err
	}

	// Calculate total weightage of other KPIs (excluding current one being updated)
	totalWeightage := 0.0
	for _, k := range kpis {
		if k.ID != kpi.ID && k.Weightage != nil {
			totalWeightage += *k.Weightage
		}
	}

	// Add the new weightage to the total
	newTotal := totalWeightage + *newWeightage

	// Validate total doesn't exceed 100
	if newTotal > 100 {
		return utility.NewValidationError(
			fmt.Errorf("total weightage of all KPIs for this objective cannot exceed 100%%. Current total would be %.2f%%", newTotal),
			"invalid_total_weightage",
			"total_weightage",
		)
	}

	return nil
}
