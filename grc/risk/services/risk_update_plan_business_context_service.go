package services

import (
	"errors"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskUpdatePlanBusinessContextServiceInterface interface {
	CreateBusinessContext(planID uint, currentUserID uint) error
	GetBusinessContextByPlanID(planID uint) ([]*riskDTOs.BusinessContextDetailDTO, error)
	UpdateBusinessContext(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextFormDTO) error
	UpdateBusinessContextStatus(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextStatusFormDTO) error
	MarkBusinessContextAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanBusinessContextService struct {
	repo              riskRepositories.RiskUpdatePlanBusinessContextRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
}

func NewRiskUpdatePlanBusinessContextService(
	repo riskRepositories.RiskUpdatePlanBusinessContextRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) RiskUpdatePlanBusinessContextServiceInterface {
	return &RiskUpdatePlanBusinessContextService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
	}
}

func (rupbcs *RiskUpdatePlanBusinessContextService) CreateBusinessContext(planID uint, currentUserID uint) error {
	// Verify plan exists
	_, err := rupbcs.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Business Context
	timelineType, err := rupbcs.timelineTypeSvc.GetByType("business_context")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupbcs.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	// Get business unit ID from plan
	businessUnitID, err := rupbcs.repo.GetBusinessUnitIDByPlanID(planID)
	if err != nil {
		return err
	}

	// Get OrganizationUnitBusinessContext data
	orgBusinessContext, err := rupbcs.repo.GetOrganizationUnitBusinessContextByBusinessUnit(businessUnitID)
	if err != nil {
		return err
	}

	tx := rupbcs.repo.BeginTransaction()

	// Create a single business context record with data from OrganizationUnitBusinessContext
	businessContext := &models.RiskUpdatePlanBusinessContext{
		RiskUpdatePlanTimelineID:   timeline.ID,
		InternalAuditNote:          orgBusinessContext.InternalAuditNote,
		ExternalContext:            orgBusinessContext.ExternalContext,
		AnnualObjectiveDescription: orgBusinessContext.AnnualObjectiveDescription,
		CreatedBYID:                &currentUserID,
		CreatedAt:                  time.Now(),
	}

	if err := rupbcs.repo.CreateBusinessContext(businessContext, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Mark timeline as completed
	if err := rupbcs.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := rupbcs.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := rupbcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskUpdatePlanBusinessContext{}.TableName(),
			planID,
			map[string]interface{}{
				"plan_id":          planID,
				"business_unit_id": businessUnitID,
				"timeline_id":      timeline.ID,
			},
			currentUserID,
			"risk", "risk_update_plan_business_context", "create", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.CreateAction,
				"planID": planID,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return nil
}

func (rupbcs *RiskUpdatePlanBusinessContextService) MarkBusinessContextAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	// Get timeline type for Business Context
	timelineType, err := rupbcs.timelineTypeSvc.GetByType("business_context")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupbcs.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupbcs.repo.BeginTransaction()

	// Update timeline - mark as not applicable and completed
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
	if err := rupbcs.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rupbcs *RiskUpdatePlanBusinessContextService) GetBusinessContextByPlanID(planID uint) ([]*riskDTOs.BusinessContextDetailDTO, error) {
	// Verify plan exists
	_, err := rupbcs.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get timeline type for Business Context
	timelineType, err := rupbcs.timelineTypeSvc.GetByType("business_context")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupbcs.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get business contexts by timeline ID
	businessContexts, err := rupbcs.repo.GetBusinessContextByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	var dtos []*riskDTOs.BusinessContextDetailDTO

	err = mapstructure.Decode(businessContexts, &dtos)
	if err != nil {
		return nil, err
	}

	return dtos, nil
}

func (rupbcs *RiskUpdatePlanBusinessContextService) UpdateBusinessContext(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextFormDTO) error {
	// Get existing business context
	businessContext, err := rupbcs.repo.GetBusinessContextByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	// Update fields from DTO
	if formDTO.InternalAuditNote != nil {
		businessContext.InternalAuditNote = formDTO.InternalAuditNote
	}
	if formDTO.ExternalContext != nil {
		businessContext.ExternalContext = formDTO.ExternalContext
	}
	if formDTO.AnnualObjectiveDescription != nil {
		businessContext.AnnualObjectiveDescription = formDTO.AnnualObjectiveDescription
	}

	tx := rupbcs.repo.BeginTransaction()

	// Update business context
	if err := rupbcs.repo.UpdateBusinessContext(businessContext, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := rupbcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanBusinessContext{}.TableName(),
			id,
			map[string]interface{}{
				"internal_audit_note":          formDTO.InternalAuditNote,
				"external_context":             formDTO.ExternalContext,
				"annual_objective_description": formDTO.AnnualObjectiveDescription,
			},
			currentUserID,
			"risk", "risk_update_plan_business_context", "update", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.UpdateAction,
				"id":     id,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return nil
}

func (rupbcs *RiskUpdatePlanBusinessContextService) UpdateBusinessContextStatus(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextStatusFormDTO) error {
	// Verify plan exists
	_, err := rupbcs.riskUpdatePlanSvc.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	// Get timeline type for Business Context
	timelineType, err := rupbcs.timelineTypeSvc.GetByType("business_context")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupbcs.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(id, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupbcs.repo.BeginTransaction()

	// Update status
	if err := rupbcs.repo.UpdateBusinessContextStatus(timeline.ID, *formDTO.IsBusinessContextUpdated, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := rupbcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanBusinessContext{}.TableName(),
			id,
			map[string]interface{}{
				"is_business_context_updated": *formDTO.IsBusinessContextUpdated,
			},
			currentUserID,
			"risk", "risk_update_plan_business_context", "update_status", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.UpdateAction,
				"id":     id,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return nil
}
