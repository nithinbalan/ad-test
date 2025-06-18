package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
)

type StrategicObjectiveKPIServiceInterface interface {
	GetAll(strategicObjectiveID uint, params *dtos.StrategicObjectiveKPIListQueryParamsDTO) (*dtos.StrategicObjectiveKPIListResponseDTO, error)
	Create(strategicObjectiveID uint, currentUserID uint, input *dtos.StrategicObjectiveKPIFormDTO) error
	DeleteStrategicObjectiveKPI(strategicObjectiveID uint, kpiRegisterID uint) error
}

// Service and implementations
type StrategicObjectiveKPIService struct {
	repo                  repositories.StrategicObjectiveKPIRepositoryInterface
	activityLogSvc        services.ActivityLogServiceInterface
	strategicObjectiveSvc StrategicObjectiveServiceInterface
	kpiRegisterSvc        KPIRegisterServiceInterface
}

func NewStrategicObjectiveKPIService(repo repositories.StrategicObjectiveKPIRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	strategicObjectiveSvc StrategicObjectiveServiceInterface,
	kpiRegisterSvc KPIRegisterServiceInterface,
) StrategicObjectiveKPIServiceInterface {
	return &StrategicObjectiveKPIService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		strategicObjectiveSvc: strategicObjectiveSvc,
		kpiRegisterSvc:        kpiRegisterSvc,
	}
}

func (soks *StrategicObjectiveKPIService) GetAll(strategicObjectiveID uint, params *dtos.StrategicObjectiveKPIListQueryParamsDTO) (*dtos.StrategicObjectiveKPIListResponseDTO, error) {
	return soks.repo.GetAllStrategicObjectiveKPI(strategicObjectiveID, params)
}

func (soks *StrategicObjectiveKPIService) Create(strategicObjectiveID uint, currentUserID uint, input *dtos.StrategicObjectiveKPIFormDTO) error {
	// Validate StrategicObjective existence via service
	if _, err := soks.strategicObjectiveSvc.GetByID(strategicObjectiveID); err != nil {
		return err
	}

	// Validate KPIRegister IDs
	if err := soks.validateKPIRegisterIDs(input.KPIRegisterIDs); err != nil {
		return err
	}

	// Build models
	var mappings []models.StrategicObjectiveKPI

	for _, kpiID := range input.KPIRegisterIDs {
		mappings = append(mappings, models.StrategicObjectiveKPI{
			StrategicObjectiveID: strategicObjectiveID,
			KPIRegisterID:        kpiID,
			CreatedBYID:          currentUserID,
		})
	}
	tx := soks.repo.BeginTransaction()

	// Bulk insert
	if len(mappings) > 0 {
		if err := soks.repo.CreateStrategicObjectiveKPIs(mappings, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	// Async bulk log
	go func() {
		ids := make([]uint, len(mappings))
		for i, m := range mappings {
			ids[i] = m.ID
		}
		err := soks.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.StrategicObjectiveKPI{}.TableName(),
			ids,
			mappings,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log StrategicObjectiveKPI map activity", map[string]interface{}{"error": err})
		}
	}()

	return nil
}
func (soks *StrategicObjectiveKPIService) DeleteStrategicObjectiveKPI(strategicObjectiveID uint, kpiRegisterID uint) error {
	// Ensure parent StrategicObjective exists
	if _, err := soks.strategicObjectiveSvc.GetByID(strategicObjectiveID); err != nil {
		return err
	}

	// Ensure KPIRegister exists
	if _, err := soks.kpiRegisterSvc.GetByID(kpiRegisterID); err != nil {
		return err
	}

	// Prepare mapping snapshot for activity log
	deletedMapping := map[string]interface{}{
		"strategic_objective_id": strategicObjectiveID,
		"kpi_register_id":        kpiRegisterID,
	}

	tx := soks.repo.BeginTransaction()

	if err := soks.repo.DeleteSingleStrategicObjectiveKPI(strategicObjectiveID, kpiRegisterID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Async log
	go func() {
		err := soks.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.StrategicObjectiveKPI{}.TableName(),
			0, // no single PK, use composite in log
			deletedMapping,
			0, // optional: pass currentUserID if you have it in context
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log StrategicObjectiveKPI unmapping activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (soks *StrategicObjectiveKPIService) validateKPIRegisterIDs(ids []uint) error {
	for _, id := range ids {
		_, err := soks.kpiRegisterSvc.GetByID(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("invalid KPIRegister ID: %d", id), "not_exist", "kpi_register_id")
		}
	}
	return nil
}
