package services

import (
	"errors"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskMitigationServiceInterface interface {
	UpdateRiskMitigation(planID uint, currentUserID uint, formDTO *riskDTOs.RiskMitigationUpdatePlanFormDTO) error
	DeleteRiskMitigation(planID uint, currentUserID uint, riskID uint) error
	GetActionPlanDetails(planID uint, riskID uint, actionPlanID uint) (*riskDTOs.RiskActionPlanDetailDTO, error)
	CreateActionPlan(planID uint, riskID uint, currentUserID uint, formDTO *riskDTOs.RiskActionPlanFormDTO) (uint, error)
	UpdateActionPlan(planID uint, riskID uint, actionPlanID uint, currentUserID uint, formDTO *riskDTOs.RiskActionPlanFormDTO) error
	DeleteActionPlan(planID uint, riskID uint, actionPlanIDs []uint, currentUserID uint) ([]uint, error)
	MarkRiskMitigationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanRiskMitigationService struct {
	repo              riskRepositories.RiskUpdatePlanRiskMitigationRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	riskResponseSvc   RiskResponseServiceInterface
	actionPlanSvc     RiskActionPlanServiceInterface
}

func NewRiskUpdatePlanRiskMitigationService(
	repo riskRepositories.RiskUpdatePlanRiskMitigationRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskResponseSvc RiskResponseServiceInterface,
	actionPlanSvc RiskActionPlanServiceInterface,
) RiskUpdatePlanRiskMitigationServiceInterface {
	return &RiskUpdatePlanRiskMitigationService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
		riskResponseSvc:   riskResponseSvc,
		actionPlanSvc:     actionPlanSvc,
	}
}

func (ruprms *RiskUpdatePlanRiskMitigationService) UpdateRiskMitigation(planID uint, currentUserID uint, formDTO *riskDTOs.RiskMitigationUpdatePlanFormDTO) error {
	// Verify plan exists
	_, err := ruprms.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Mitigation
	timelineType, err := ruprms.timelineTypeSvc.GetByType("risk_mitigation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruprms.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Begin transaction
	tx := ruprms.repo.BeginTransaction()

	// Verify risk exists
	_, err = ruprms.repo.GetRiskByID(formDTO.RiskRegisterID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, formDTO.RiskRegisterID)
		}
		return err
	}

	// Update risk response strategy
	if err := ruprms.repo.UpdateRiskResponseStrategy(formDTO.RiskRegisterID, formDTO.RiskResponseStrategyID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Count risk mitigations to check if at least one exists
	riskMitigationCount, err := ruprms.repo.CountRiskMitigationsByPlanID(planID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update timeline completion only if this is the first risk mitigation (count >= 1)
	if riskMitigationCount >= 1 && !timeline.IsCompleted {
		// Update timeline completion
		if err := ruprms.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Calculate and update percentage
		if err := ruprms.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Log activity
	go func() {
		err := ruprms.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			formDTO.RiskRegisterID,
			map[string]interface{}{
				"risk_response_strategy_id": formDTO.RiskResponseStrategyID,
			},
			currentUserID,
			"risk", "risk_register_update_plan", "risk_mitigation", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (ruprms *RiskUpdatePlanRiskMitigationService) DeleteRiskMitigation(planID uint, currentUserID uint, riskID uint) error {
	// Verify plan exists
	_, err := ruprms.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Mitigation
	timelineType, err := ruprms.timelineTypeSvc.GetByType("risk_mitigation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruprms.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Begin transaction
	tx := ruprms.repo.BeginTransaction()

	// Verify risk exists
	_, err = ruprms.repo.GetRiskByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, riskID)
		}
		return err
	}

	// Delete risk response strategy (set to null)
	if err := ruprms.repo.UpdateRiskResponseStrategy(riskID, nil, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Mark timeline as incomplete
	if err := ruprms.riskUpdatePlanSvc.MarkTimelineAsIncomplete(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Calculate and update percentage
	if err := ruprms.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Log activity
	go func() {
		err := ruprms.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"risk_response_strategy_id": nil,
			},
			currentUserID,
			"risk", "risk_register_update_plan", "risk_mitigation", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (ruprms *RiskUpdatePlanRiskMitigationService) GetActionPlanDetails(planID uint, riskID uint, actionPlanID uint) (*riskDTOs.RiskActionPlanDetailDTO, error) {
	// Verify plan exists
	_, err := ruprms.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get action plan full details using the action plan service
	actionPlanDetails, err := ruprms.actionPlanSvc.GetFullDetailsByID(riskID, actionPlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, actionPlanID)
		}
		return nil, err
	}

	return actionPlanDetails, nil
}

func (ruprms *RiskUpdatePlanRiskMitigationService) CreateActionPlan(planID uint, riskID uint, currentUserID uint, formDTO *riskDTOs.RiskActionPlanFormDTO) (uint, error) {
	// Verify plan exists
	_, err := ruprms.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utility.NewResourceNotFoundError(err, planID)
		}
		return 0, err
	}

	// Get timeline type for Risk Mitigation
	timelineType, err := ruprms.timelineTypeSvc.GetByType("risk_mitigation")
	if err != nil {
		return 0, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruprms.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utility.NewResourceNotFoundError(err, planID)
		}
		return 0, err
	}

	// Begin transaction
	tx := ruprms.repo.BeginTransaction()

	// Create action plan using the action plan service
	actionPlanID, err := ruprms.actionPlanSvc.Create(riskID, currentUserID, formDTO)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update timeline completion only if not already completed
	if !timeline.IsCompleted {
		// Update timeline completion
		if err := ruprms.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		// Calculate and update percentage
		if err := ruprms.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	// Log activity
	go func() {
		err := ruprms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlan{}.TableName(),
			actionPlanID,
			map[string]interface{}{
				"risk_register_id": riskID,
				"title":            formDTO.Title,
			},
			currentUserID,
			"risk", "risk_register_update_plan", "action_plan", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return actionPlanID, nil
}

func (ruprms *RiskUpdatePlanRiskMitigationService) UpdateActionPlan(planID uint, riskID uint, actionPlanID uint, currentUserID uint, formDTO *riskDTOs.RiskActionPlanFormDTO) error {
	// Verify plan exists
	_, err := ruprms.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Update action plan using the action plan service
	if err := ruprms.actionPlanSvc.Update(riskID, actionPlanID, currentUserID, formDTO); err != nil {
		return err
	}

	// Log activity
	go func() {
		err := ruprms.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskActionPlan{}.TableName(),
			actionPlanID,
			map[string]interface{}{
				"risk_register_id": riskID,
				"title":            formDTO.Title,
			},
			currentUserID,
			"risk", "risk_register_update_plan", "action_plan", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (ruprms *RiskUpdatePlanRiskMitigationService) DeleteActionPlan(planID uint, riskID uint, actionPlanIDs []uint, currentUserID uint) ([]uint, error) {
	// Verify plan exists
	_, err := ruprms.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Delete action plans using the action plan service
	failedIDs, err := ruprms.actionPlanSvc.Delete(riskID, actionPlanIDs, currentUserID)
	if err != nil {
		return failedIDs, err
	}

	// Log activity for each successfully deleted action plan
	go func() {
		for _, id := range actionPlanIDs {
			isFailedID := false
			for _, failedID := range failedIDs {
				if id == failedID {
					isFailedID = true
					break
				}
			}
			if !isFailedID {
				err := ruprms.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.RiskActionPlan{}.TableName(),
					id,
					map[string]interface{}{
						"risk_register_id": riskID,
					},
					currentUserID,
					"risk", "risk_register_update_plan", "action_plan", nil,
				)
				if err != nil {
					utility.LogError("Failed to log activity", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
	}()

	return failedIDs, nil
}

func (ruprms *RiskUpdatePlanRiskMitigationService) MarkRiskMitigationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	// Verify plan exists
	_, err := ruprms.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Mitigation
	timelineType, err := ruprms.timelineTypeSvc.GetByType("risk_mitigation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruprms.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Begin transaction
	tx := ruprms.repo.BeginTransaction()

	now := time.Now()
	timeline.IsNotApplicable = true
	timeline.IsCompleted = true
	timeline.Justification = &formDTO.Justification
	timeline.UpdatedBYID = &currentUserID
	timeline.UpdatedAt = &now

	if err := tx.DB().Save(&timeline).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := ruprms.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
