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

type RiskUpdatePlanRiskAnalysisServiceInterface interface {
	GetRiskAnalysisList(planID uint, params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskAnalysisListResponseDTO, error)
	CreateRiskAnalysis(planID uint, currentUserID uint, formDTO *riskDTOs.RiskAnalysisFormDTO) error
	MarkRiskAnalysisAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanRiskAnalysisService struct {
	repo                         riskRepositories.RiskUpdatePlanRiskAnalysisRepositoryInterface
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
}

func NewRiskUpdatePlanRiskAnalysisService(
	repo riskRepositories.RiskUpdatePlanRiskAnalysisRepositoryInterface,
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
) RiskUpdatePlanRiskAnalysisServiceInterface {
	return &RiskUpdatePlanRiskAnalysisService{
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
	}
}

func (rrups *RiskUpdatePlanRiskAnalysisService) GetRiskAnalysisList(planID uint, params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskAnalysisListResponseDTO, error) {
	// Verify plan exists
	_, err := rrups.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	return rrups.repo.GetRiskAnalysisList(planID, params)
}

// Bulk Inherent Evaluation

func (rrups *RiskUpdatePlanRiskAnalysisService) CreateRiskAnalysis(planID uint, currentUserID uint, formDTO *riskDTOs.RiskAnalysisFormDTO) error {
	// Verify plan exists
	_, err := rrups.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Workshop Assessment
	timelineType, err := rrups.timelineTypeSvc.GetByType("risk_analysis")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rrups.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rrups.repo.BeginTransaction()

	// Extract risk register IDs
	var riskRegisterIDs []uint
	for _, evaluation := range formDTO.Analysis {
		riskRegisterIDs = append(riskRegisterIDs, evaluation.RiskRegisterID)
	}

	// Delete existing inherent evaluations for the given risk register IDs
	if err := rrups.repo.DeleteInherentEvaluationsByRiskRegisterIDs(riskRegisterIDs, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Get risk matrix for calculations
	riskMatrixType, err := rrups.riskMatrixSvc.GetByType("erm")
	if err != nil {
		tx.Rollback()
		return err
	}

	riskMatrix, err := rrups.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
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

	// Create evaluations
	var evaluations []*models.RiskEvaluation
	var impactAnalyses []*models.RiskImpactAnalysis
	var analysisGuidelines []*models.RiskImpactAnalysisGuideline

	for _, evaluationForm := range formDTO.Analysis {
		// Validate risk register exists
		if _, err := rrups.riskRegisterSvc.GetByID(evaluationForm.RiskRegisterID); err != nil {
			tx.Rollback()
			return err
		}

		// Validate form data
		if _, err := rrups.riskMatrixImpactSvc.GetByID(evaluationForm.RiskMatrixImpactID); err != nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("risk matrix impact does not exist"), "not_exist", "risk matrix impact")
		}
		if _, err := rrups.riskMatrixLikelihoodSvc.GetByID(evaluationForm.RiskMatrixLikelihoodID); err != nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("risk matrix likelihood does not exist"), "not_exist", "risk matrix likelihood")
		}

		// Calculate risk rating
		impact, err := rrups.riskMatrixImpactSvc.GetByID(evaluationForm.RiskMatrixImpactID)
		if err != nil {
			tx.Rollback()
			return err
		}
		likelihood, err := rrups.riskMatrixLikelihoodSvc.GetByID(evaluationForm.RiskMatrixLikelihoodID)
		if err != nil {
			tx.Rollback()
			return err
		}
		calculationMethod := riskMatrix.CalculationType
		score := rrups.riskMatrixRiskRatingValueSvc.EvaluateScore(impact, likelihood, calculationMethod)
		riskRating, _ := rrups.riskMatrixRiskRatingValueSvc.FindRiskRatingValue(riskRatingValueArray, score)

		// Create risk evaluation
		riskEvaluation := &models.RiskEvaluation{
			RiskRegisterID:              evaluationForm.RiskRegisterID,
			Type:                        "inherent",
			RiskMatrixLikelihoodID:      evaluationForm.RiskMatrixLikelihoodID,
			RiskMatrixImpactID:          evaluationForm.RiskMatrixImpactID,
			RiskMatrixRiskRatingID:      riskRating.RiskMatrixRiskRatingID,
			RiskMatrixRiskRatingValueID: riskRating.RiskMatrixRiskRatingValueID,
			CreatedByID:                 currentUserID,
		}
		evaluations = append(evaluations, riskEvaluation)

		// Create impact analyses (will need evaluation ID after create)
		for _, impactEval := range evaluationForm.RiskImpactAnalyses {
			// Validate impact category
			if _, err := rrups.riskImpactCategorySvc.GetByID(impactEval.RiskImpactCategoryID); err != nil {
				tx.Rollback()
				return utility.NewValidationError(fmt.Errorf("risk impact category does not exist"), "not_exist", "risk impact category")
			}
			// Validate impact
			if _, err := rrups.riskMatrixImpactSvc.GetByID(impactEval.RiskMatrixImpactID); err != nil {
				tx.Rollback()
				return utility.NewValidationError(fmt.Errorf("risk matrix impact does not exist"), "not_exist", "risk matrix impact")
			}

			riskImpactAnalysis := &models.RiskImpactAnalysis{
				RiskEvaluationID:     0, // Will be set after evaluation creation
				RiskImpactCategoryID: impactEval.RiskImpactCategoryID,
				RiskMatrixImpactID:   impactEval.RiskMatrixImpactID,
				CreatedByID:          currentUserID,
			}
			impactAnalyses = append(impactAnalyses, riskImpactAnalysis)
		}
	}

	// Create evaluations in bulk
	if err := rrups.repo.CreateBulkInherentEvaluations(evaluations, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Update impact analyses with evaluation IDs and create them
	evaluationIndex := 0
	for i := 0; i < len(impactAnalyses); i++ {
		impactAnalyses[i].RiskEvaluationID = evaluations[evaluationIndex].ID
		evaluationIndex++
		if evaluationIndex >= len(evaluations) {
			evaluationIndex = 0
		}
	}

	if err := rrups.repo.CreateBulkRiskImpactAnalyses(impactAnalyses, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Auto-fetch and create guidelines for each impact analysis
	for _, impactAnalysis := range impactAnalyses {
		guidelines, err := rrups.riskMatrixImpactGuidelineSvc.GetByCategoryAndImpact(impactAnalysis.RiskImpactCategoryID, impactAnalysis.RiskMatrixImpactID)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, guideline := range guidelines {
			riskImpactAnalysisGuideline := &models.RiskImpactAnalysisGuideline{
				RiskImpactAnalysisID:        impactAnalysis.ID,
				RiskMatrixImpactGuidelineID: guideline.ID,
			}
			analysisGuidelines = append(analysisGuidelines, riskImpactAnalysisGuideline)
		}
	}

	if err := rrups.repo.CreateBulkRiskImpactAnalysisGuidelines(analysisGuidelines, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Update risk register ratings
	for _, evaluation := range evaluations {
		if err := rrups.riskRegisterSvc.UpdateRiskRating(evaluation.RiskRegisterID, currentUserID, evaluation, tx, "inherent"); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Mark timeline as completed
	if err := rrups.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := rrups.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		for _, evaluation := range evaluations {
			err := rrups.activityLogSvc.LogActivity(
				models.CreateAction,
				models.RiskEvaluation{}.TableName(),
				evaluation.ID,
				map[string]interface{}{
					"risk_register_id": evaluation.RiskRegisterID,
					"type":             evaluation.Type,
				},
				currentUserID,
				"risk", "risk_register_update_plan", "bulk_inherent_evaluation", nil,
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

func (rrups *RiskUpdatePlanRiskAnalysisService) MarkRiskAnalysisAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := rrups.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := rrups.timelineTypeSvc.GetByType("risk_analysis")
	if err != nil {
		return err
	}

	timeline, err := rrups.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rrups.repo.BeginTransaction()

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

	if err := rrups.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
