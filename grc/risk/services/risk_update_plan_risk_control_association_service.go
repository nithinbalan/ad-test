package services

import (
	"errors"
	"fmt"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskControlAssociationServiceInterface interface {
	GetRiskControlAssociationList(planID uint, isCritical bool) ([]*riskDTOs.RiskControlAssociationListDTO, error)
	CreateRiskControlAssociation(planID uint, currentUserID uint, formDTO *riskDTOs.RiskControlAssociationFormDTO) error
	MarkRiskControlAssociationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanRiskControlAssociationService struct {
	repo                         riskRepositories.RiskUpdatePlanRiskControlAssociationRepositoryInterface
	riskUpdatePlanSvc            RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc              RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc              RiskRegisterServiceInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	userSvc                      organizationServices.UserServiceInterface
	organizationUnitSvc          settingsServices.OrganizationUnitServiceInterface
	riskStatusSvc                RiskStatusServiceInterface
	riskMatrixSvc                RiskMatrixServiceInterface
	organizationUnitsSvc         organizationServices.OrganizationUnitsServiceInterface
	notificationSvc              commonServices.NotificationServiceInterface
	riskSettingsSvc              RiskSettingsServiceInterface
	riskCategorySvc              RiskCategoryServiceInterface
	riskMatrixImpactSvc          RiskMatrixImpactServiceInterface
	riskMatrixLikelihoodSvc      RiskMatrixLikelihoodServiceInterface
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface
	riskImpactCategorySvc        RiskImpactCategoryServiceInterface
	riskMatrixImpactGuidelineSvc RiskMatrixImpactGuidelineServiceInterface
	riskInherentEvaluationSvc    RiskInherentEvaluationServiceInterface
	controlSvc                   organizationServices.ControlServiceInterface
	controlTestSvc               organizationServices.ControlTestServiceInterface
	riskTimelineSvc              RiskTimelineServiceInterface
}

func NewRiskUpdatePlanRiskControlAssociationService(
	repo riskRepositories.RiskUpdatePlanRiskControlAssociationRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	organizationUnitsSvc organizationServices.OrganizationUnitsServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	riskSettingsSvc RiskSettingsServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface,
	riskMatrixImpactSvc RiskMatrixImpactServiceInterface,
	riskMatrixLikelihoodSvc RiskMatrixLikelihoodServiceInterface,
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface,
	riskImpactCategorySvc RiskImpactCategoryServiceInterface,
	riskMatrixImpactGuidelineSvc RiskMatrixImpactGuidelineServiceInterface,
	riskInherentEvaluationSvc RiskInherentEvaluationServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	controlTestSvc organizationServices.ControlTestServiceInterface,
	riskTimelineSvc RiskTimelineServiceInterface,
) RiskUpdatePlanRiskControlAssociationServiceInterface {
	return &RiskUpdatePlanRiskControlAssociationService{
		repo:                         repo,
		riskUpdatePlanSvc:            riskUpdatePlanSvc,
		timelineTypeSvc:              timelineTypeSvc,
		riskRegisterSvc:              riskRegisterSvc,
		activityLogSvc:               activityLogSvc,
		userSvc:                      userSvc,
		organizationUnitSvc:          organizationUnitSvc,
		riskStatusSvc:                riskStatusSvc,
		riskMatrixSvc:                riskMatrixSvc,
		organizationUnitsSvc:         organizationUnitsSvc,
		notificationSvc:              notificationSvc,
		riskSettingsSvc:              riskSettingsSvc,
		riskCategorySvc:              riskCategorySvc,
		riskMatrixImpactSvc:          riskMatrixImpactSvc,
		riskMatrixLikelihoodSvc:      riskMatrixLikelihoodSvc,
		riskMatrixRiskRatingValueSvc: riskMatrixRiskRatingValueSvc,
		riskImpactCategorySvc:        riskImpactCategorySvc,
		riskMatrixImpactGuidelineSvc: riskMatrixImpactGuidelineSvc,
		riskInherentEvaluationSvc:    riskInherentEvaluationSvc,
		controlSvc:                   controlSvc,
		controlTestSvc:               controlTestSvc,
		riskTimelineSvc:              riskTimelineSvc,
	}
}

func (rupras *RiskUpdatePlanRiskControlAssociationService) CreateRiskControlAssociation(planID uint, currentUserID uint, formDTO *riskDTOs.RiskControlAssociationFormDTO) error {
	// Verify plan exists
	_, err := rupras.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Workshop Assessment
	timelineType, err := rupras.timelineTypeSvc.GetByType("risk_control_association")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupras.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	// Validate that form has at least one control association
	if len(formDTO.Residual) == 0 {
		return utility.NewValidationError(fmt.Errorf("at least one control association is required"), "required", "control associations")
	}

	// Pre-validate all risk registers and controls outside transaction
	riskRegisterMap := make(map[uint]bool) // Track unique risk register IDs
	riskControlCount := make(map[uint]int) // Track control count per risk

	for _, item := range formDTO.Residual {
		// Track unique risk register IDs
		riskRegisterMap[item.RiskRegisterID] = true

		// Count controls per risk
		riskControlCount[item.RiskRegisterID]++

		// Validate risk register exists
		if _, err := rupras.riskRegisterSvc.GetByID(item.RiskRegisterID); err != nil {
			return err
		}

		// Validate control exists
		if _, err := rupras.controlSvc.GetByID(item.ControlID); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", item.ControlID), "not_exist", "control_id")
		}
	}

	// Build risk register IDs slice and validate that every risk has at least one control
	var riskRegisterIDs []uint
	for riskID := range riskRegisterMap {
		riskRegisterIDs = append(riskRegisterIDs, riskID)
		if riskControlCount[riskID] == 0 {
			return utility.NewValidationError(fmt.Errorf("risk register ID %d must have at least one control", riskID), "required", "control for risk")
		}
	}

	// Pre-fetch control details and latest control tests outside transaction
	controlDetailsMap := make(map[uint]*models.Control)
	controlTestMap := make(map[uint]*models.ControlTest)

	for _, item := range formDTO.Residual {
		if _, exists := controlDetailsMap[item.ControlID]; !exists {
			control, err := rupras.controlSvc.GetByID(item.ControlID)
			if err != nil {
				return err
			}
			controlDetailsMap[item.ControlID] = control

			latestControlTest, err := rupras.controlTestSvc.GetLatestByControlIDAndPlanID(item.ControlID, planID)
			if err != nil {
				if _, notFound := err.(*utility.ResourceNotFoundError); !notFound {
					return err
				}

				latestControlTest, err = rupras.controlTestSvc.GetLatestByControlID(item.ControlID)
				if err != nil {
					return err
				}
			}
			controlTestMap[item.ControlID] = latestControlTest
		}
	}

	// Pre-fetch inherent evaluations outside transaction
	inherentEvaluationMap := make(map[uint]*riskDTOs.RiskInherentEvaluationDetailsDTO)
	for _, riskRegisterID := range riskRegisterIDs {
		if _, exists := inherentEvaluationMap[riskRegisterID]; !exists {
			inherentEvaluation, err := rupras.riskInherentEvaluationSvc.GetLatestPlanDraftFullDetailsByRiskID(riskRegisterID)
			if err != nil {
				return err
			}
			if inherentEvaluation == nil {
				return utility.NewValidationError(fmt.Errorf("inherent evaluation must be completed before creating residual evaluation"), "not_exist", "inherent evaluation")
			}
			inherentEvaluationMap[riskRegisterID] = inherentEvaluation
		}
	}

	tx := rupras.repo.BeginTransaction()

	// Step 1: Delete existing control analyses for these risk registers (must be deleted first due to foreign key constraint)
	if err := rupras.repo.DeleteControlAnalysesByRiskRegisterIDs(riskRegisterIDs, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Step 2: Risk-Control Mapping
	// Delete existing risk controls for the given risk register IDs (after deleting control analyses)
	for _, riskRegisterID := range riskRegisterIDs {
		if err := rupras.repo.DeleteRiskControlsByRiskRegisterID(riskRegisterID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create new risk-control mappings
	var riskControls []*models.RiskControl
	currentTime := time.Now()

	for _, item := range formDTO.Residual {
		riskControl := &models.RiskControl{
			RiskRegisterID: item.RiskRegisterID,
			ControlID:      item.ControlID,
			CreatedByID:    currentUserID,
			CreatedAt:      currentTime,
		}
		riskControls = append(riskControls, riskControl)
	}

	if err := rupras.repo.CreateBulkRiskControls(riskControls, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Step 3: Control Analysis (using latest control test results)

	// Create control analyses using pre-fetched data
	var controlAnalyses []*models.RiskControlAnalysis

	for _, riskControl := range riskControls {
		control := controlDetailsMap[riskControl.ControlID]
		latestControlTest := controlTestMap[riskControl.ControlID]

		// Validate that control test exists
		if latestControlTest == nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("control test not completed for control ID %d", riskControl.ControlID), "not_exist", "control test")
		}

		riskControlAnalysis := &models.RiskControlAnalysis{
			RiskControlID:                 riskControl.ID,
			Type:                          "inherent",
			ControlEfficiencyStatusID:     control.ControlEfficiencyStatusID,
			ControlImplementationStatusID: *control.ControlImplementationStatusID,
			Score:                         control.EfficiencyScore,
			ControlTestID:                 latestControlTest.ID,
			ControlEffectivenessScore:     latestControlTest.ControlEffectivenessScore,
			ControlEffectivenessRatingID:  latestControlTest.ControlEffectivenessRatingID,
			RiskReduction:                 latestControlTest.RiskReduction,
			CreatedByID:                   currentUserID,
			CreatedAt:                     currentTime,
		}

		controlAnalyses = append(controlAnalyses, riskControlAnalysis)
	}

	if err := rupras.repo.CreateBulkControlAnalyses(controlAnalyses, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Build mapping from risk register ID to risk controls
	riskControlMap := make(map[uint][]*models.RiskControl)
	for _, rc := range riskControls {
		riskControlMap[rc.RiskRegisterID] = append(riskControlMap[rc.RiskRegisterID], rc)
	}

	// Build mapping from risk register ID to control analyses
	controlAnalysisMap := make(map[uint][]*models.RiskControlAnalysis)
	for _, ca := range controlAnalyses {
		for riskRegisterID, controls := range riskControlMap {
			for _, rc := range controls {
				if ca.RiskControlID == rc.ID {
					controlAnalysisMap[riskRegisterID] = append(controlAnalysisMap[riskRegisterID], ca)
					break
				}
			}
		}
	}

	// Step 3: Residual Evaluation (using inherent evaluation's likelihood and impact ID)
	// Get risk matrix for calculations
	riskMatrixType, err := rupras.riskMatrixSvc.GetByType("erm")
	if err != nil {
		tx.Rollback()
		return err
	}

	riskMatrix, err := rupras.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Build risk rating value array
	var riskRatingValueArray []*riskDTOs.RiskMatrixRiskRatingDataDTO
	for _, riskRating := range riskMatrix.RiskMatrixRiskRating {
		for _, value := range riskRating.RiskMatrixRiskRatingValue {
			detail := &riskDTOs.RiskMatrixRiskRatingDataDTO{
				RiskMatrixID:                riskRating.RiskMatrixID,
				RiskMatrixRiskRatingID:      value.RiskMatrixRiskRatingID,
				RiskRatingTitleEN:           riskRating.TitleEn,
				RiskRatingTitleAR:           riskRating.TitleAr,
				RiskRatingTitleType:         riskRating.Type,
				ColorCode:                   riskRating.ColorCode,
				RiskMatrixRiskRatingValueID: value.ID,
				Value:                       value.Value,
			}
			riskRatingValueArray = append(riskRatingValueArray, detail)
		}
	}

	// Create residual evaluations using pre-fetched inherent evaluations
	var evaluations []*models.RiskEvaluation

	for _, item := range formDTO.Residual {
		inherentEvaluation := inherentEvaluationMap[item.RiskRegisterID]

		// Get control analyses from in-memory mapping
		controlAnalysesForRisk := controlAnalysisMap[item.RiskRegisterID]

		// Calculate average risk reduction from control analyses
		var totalRiskReduction float64
		var controlCount int
		for _, ca := range controlAnalysesForRisk {
			if ca.RiskReduction != nil {
				totalRiskReduction += *ca.RiskReduction
				controlCount++
			}
		}

		var averageRiskReduction *float64
		if controlCount > 0 {
			avg := totalRiskReduction / float64(controlCount)
			averageRiskReduction = &avg
		}

		// Calculate residual score: Inherent Risk Rating Value × Average Risk Reduction
		var residualScore uint
		if averageRiskReduction != nil {
			residualScore = uint(float64(inherentEvaluation.RiskMatrixRiskRatingValue.Value)*(1-*averageRiskReduction/100) + 0.5)
		}

		riskRating, found := rupras.riskMatrixRiskRatingValueSvc.FindClosestRiskRatingValue(riskRatingValueArray, residualScore)
		if !found || riskRating == nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("no risk rating found for score %d", residualScore), "not_found", "risk rating")
		}

		// Create new draft evaluation for this plan
		isLatest := false
		riskResidualEvaluation := &models.RiskEvaluation{
			RiskRegisterID:                   item.RiskRegisterID,
			Type:                             "residual",
			RiskMatrixLikelihoodID:           inherentEvaluation.RiskMatrixLikelihood.ID,
			RiskMatrixImpactID:               inherentEvaluation.RiskMatrixImpact.ID,
			RiskMatrixRiskRatingID:           riskRating.RiskMatrixRiskRatingID,
			RiskMatrixRiskRatingValueID:      riskRating.RiskMatrixRiskRatingValueID,
			AverageControlEffectivenessScore: averageRiskReduction,
			CreatedByID:                      currentUserID,
			IsLatest:                         &isLatest,
			RiskRegisterUpdatePlanID:         &planID,
		}

		if err := rupras.repo.CreateResidualEvaluation(riskResidualEvaluation, tx); err != nil {
			tx.Rollback()
			return err
		}
		evaluations = append(evaluations, riskResidualEvaluation)

		// Create history record for this evaluation
		history := &models.RiskResidualEvaluationHistory{
			RiskEvaluationID:                 evaluations[len(evaluations)-1].ID,
			RiskRegisterID:                   item.RiskRegisterID,
			RiskMatrixLikelihoodID:           inherentEvaluation.RiskMatrixLikelihood.ID,
			RiskMatrixRiskRatingID:           riskRating.RiskMatrixRiskRatingID,
			RiskMatrixRiskRatingValueID:      riskRating.RiskMatrixRiskRatingValueID,
			AverageControlEffectivenessScore: averageRiskReduction,
			CreatedBYID:                      currentUserID,
			CreatedAt:                        time.Now(),
		}

		if err := rupras.repo.CreateResidualEvaluationHistory(history, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update individual risk timelines for control_analysis
	controlAnalysisTimelineType, err := rupras.riskTimelineSvc.GetByType("control_analysis")
	if err != nil {
		tx.Rollback()
		return err
	}
	completedStatus, err := rupras.riskTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return err
	}

	currentTime = time.Now()
	for _, riskID := range riskRegisterIDs {
		timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
			PerformedByID:    &currentUserID,
			TimelineStatusID: completedStatus.ID,
			PerformedAt:      &currentTime,
		}
		if err := rupras.riskTimelineSvc.Update(riskID, timelineFormDTO, controlAnalysisTimelineType.ID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update individual risk timelines for residual_risk_analysis
	residualAnalysisTimelineType, err := rupras.riskTimelineSvc.GetByType("residual_risk_analysis")
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, riskID := range riskRegisterIDs {
		timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
			PerformedByID:    &currentUserID,
			TimelineStatusID: completedStatus.ID,
			PerformedAt:      &currentTime,
		}
		if err := rupras.riskTimelineSvc.Update(riskID, timelineFormDTO, residualAnalysisTimelineType.ID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Mark timeline as completed
	if err := rupras.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := rupras.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity for risk controls
	go func() {
		for _, item := range formDTO.Residual {
			err := rupras.activityLogSvc.LogActivity(
				models.CreateAction,
				models.RiskControl{}.TableName(),
				item.RiskRegisterID,
				map[string]interface{}{
					"risk_register_id": item.RiskRegisterID,
					"control_id":       item.ControlID,
				},
				currentUserID,
				"risk", "risk_register_update_plan", "risk_control_association", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	// Log activity for residual evaluations
	go func() {
		for _, evaluation := range evaluations {
			err := rupras.activityLogSvc.LogActivity(
				models.CreateAction,
				models.RiskEvaluation{}.TableName(),
				evaluation.ID,
				map[string]interface{}{
					"risk_register_id": evaluation.RiskRegisterID,
					"type":             evaluation.Type,
				},
				currentUserID,
				"risk", "risk_register_update_plan", "risk_control_association", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (rupras *RiskUpdatePlanRiskControlAssociationService) MarkRiskControlAssociationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := rupras.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := rupras.timelineTypeSvc.GetByType("risk_control_association")
	if err != nil {
		return err
	}

	timeline, err := rupras.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupras.repo.BeginTransaction()

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

	if err := rupras.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (rupras *RiskUpdatePlanRiskControlAssociationService) GetRiskControlAssociationList(planID uint, isCritical bool) ([]*riskDTOs.RiskControlAssociationListDTO, error) {
	return rupras.repo.GetRiskControlAssociationList(planID, isCritical)
}
