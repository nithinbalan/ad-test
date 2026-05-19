package services

import (
	"errors"
	logServices "jethur/grc/logs/services"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type RiskUpdatePlanControlRegisterUpdateServiceInterface interface {
	GetControlCountsByImplementationStatus(planID uint) (*riskDTOs.ControlRegisterUpdatedDTO, error)
	CreateControlTest(planID uint, currentUserID uint, controlTestFormDTO *riskDTOs.UpdatePlanControlTestFormDTO) (uint, error)
	MarkControlRegisterUpdateAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanControlRegisterUpdateService struct {
	repo              riskRepositories.RiskUpdatePlanControlRegisterUpdateRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	controlTestSvc    organizationServices.ControlTestServiceInterface
	controlSvc        organizationServices.ControlServiceInterface
}

func NewRiskUpdatePlanControlRegisterUpdateService(
	repo riskRepositories.RiskUpdatePlanControlRegisterUpdateRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlTestSvc organizationServices.ControlTestServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
) RiskUpdatePlanControlRegisterUpdateServiceInterface {
	return &RiskUpdatePlanControlRegisterUpdateService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
		controlTestSvc:    controlTestSvc,
		controlSvc:        controlSvc,
	}
}

func (rupcrus *RiskUpdatePlanControlRegisterUpdateService) GetControlCountsByImplementationStatus(planID uint) (*riskDTOs.ControlRegisterUpdatedDTO, error) {
	return rupcrus.repo.GetControlCountsByImplementationStatus(planID)
}

func (rupcrus *RiskUpdatePlanControlRegisterUpdateService) CreateControlTest(planID uint, currentUserID uint, controlTestFormDTO *riskDTOs.UpdatePlanControlTestFormDTO) (uint, error) {
	// Verify plan exists
	_, err := rupcrus.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utility.NewResourceNotFoundError(err, planID)
		}
		return 0, err
	}

	// Get timeline type for control_register_updated
	timelineType, err := rupcrus.timelineTypeSvc.GetByType("control_register_updated")
	if err != nil {
		return 0, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupcrus.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return 0, err
	}

	tx := rupcrus.repo.BeginTransaction()
	defer tx.Rollback()

	// Convert UpdatePlanControlTestFormDTO to ControlTestFormDTO
	orgControlTestDTO := &organizationDTOs.ControlTestFormDTO{
		ControlImplementationStatusID: controlTestFormDTO.ControlImplementationStatusID,
		ControlEfficiencyStatusID:     controlTestFormDTO.ControlEfficiencyStatusID,
		Title:                         controlTestFormDTO.Title,
		Description:                   controlTestFormDTO.Description,
		Score:                         controlTestFormDTO.Score,
		Comments:                      controlTestFormDTO.Comments,
		ControlAutomationID:           controlTestFormDTO.ControlAutomationID,
		ControlConsistencyID:          controlTestFormDTO.ControlConsistencyID,
		TestedAt:                      controlTestFormDTO.TestedAt,
		TestedBy:                      controlTestFormDTO.TestedBy,
		Documents:                     controlTestFormDTO.Documents,
	}

	// Create control test using existing service with update_plan source
	controlTestID, err := rupcrus.controlTestSvc.Create(controlTestFormDTO.ControlID, currentUserID, nil, orgControlTestDTO, "update_plan", tx)
	if err != nil {
		return 0, err
	}

	// Update the control test with RiskRegisterUpdatePlanID
	if err := rupcrus.repo.UpdateControlTestPlanID(controlTestID, planID, tx); err != nil {
		return 0, err
	}

	// Update timeline completion
	if err := rupcrus.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		return 0, err
	}

	// Recalculate percentage
	if err := rupcrus.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		return 0, err
	}

	tx.Commit()

	return controlTestID, nil
}

func (rupcrus *RiskUpdatePlanControlRegisterUpdateService) MarkControlRegisterUpdateAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := rupcrus.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := rupcrus.timelineTypeSvc.GetByType("control_register_updated")
	if err != nil {
		return err
	}

	timeline, err := rupcrus.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupcrus.repo.BeginTransaction()

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

	if err := rupcrus.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
