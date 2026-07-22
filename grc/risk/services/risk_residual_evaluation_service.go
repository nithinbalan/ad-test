package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskResidualEvaluationServiceInterface interface {
	GetFullDetailsByID(riskID uint) (*riskDTOs.RiskResidualEvaluationDetailsDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskEvaluation, error)
	Create(riskID uint, currentUserID uint, residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
	GetHistoryByRiskID(riskID uint) ([]*riskDTOs.RiskResidualEvaluationHistoryDTO, error)
}

// Service and implementations
type RiskResidualEvaluationService struct {
	repo                         riskRepositories.RiskResidualEvaluationRepositoryInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	riskRegisterSvc              RiskRegisterServiceInterface
	riskMatrixSvc                RiskMatrixServiceInterface
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface
	riskMatrixImpactSvc          RiskMatrixImpactServiceInterface
	riskMatrixLikelihoodSvc      RiskMatrixLikelihoodServiceInterface
	riskImpactCategorySvc        RiskImpactCategoryServiceInterface
	riskMatrixImpactGuidelineSvc RiskMatrixImpactGuidelineServiceInterface
	riskAppetiteStatementSvc     RiskAppetiteStatementServiceInterface
	riskInherentEvaluationSvc    RiskInherentEvaluationServiceInterface
	riskTimelineSvc              RiskTimelineServiceInterface
	controlStatusSvc             organizationServices.ControlStatusServiceInterface
	controlSvc                   organizationServices.ControlServiceInterface
}

func NewRiskResidualEvaluationService(
	repo riskRepositories.RiskResidualEvaluationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface,
	riskMatrixImpactSvc RiskMatrixImpactServiceInterface,
	riskMatrixLikelihoodSvc RiskMatrixLikelihoodServiceInterface,
	riskImpactCategorySvc RiskImpactCategoryServiceInterface,
	riskMatrixImpactGuidelineSvc RiskMatrixImpactGuidelineServiceInterface,
	riskAppetiteStatementSvc RiskAppetiteStatementServiceInterface,
	riskInherentEvaluationSvc RiskInherentEvaluationServiceInterface,
	riskTimelineSvc RiskTimelineServiceInterface,
	controlStatusSvc organizationServices.ControlStatusServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
) RiskResidualEvaluationServiceInterface {
	return &RiskResidualEvaluationService{
		repo:                         repo,
		activityLogSvc:               activityLogSvc,
		riskRegisterSvc:              riskRegisterSvc,
		riskMatrixSvc:                riskMatrixSvc,
		riskMatrixRiskRatingValueSvc: riskMatrixRiskRatingValueSvc,
		riskMatrixImpactSvc:          riskMatrixImpactSvc,
		riskMatrixLikelihoodSvc:      riskMatrixLikelihoodSvc,
		riskImpactCategorySvc:        riskImpactCategorySvc,
		riskMatrixImpactGuidelineSvc: riskMatrixImpactGuidelineSvc,
		riskAppetiteStatementSvc:     riskAppetiteStatementSvc,
		riskInherentEvaluationSvc:    riskInherentEvaluationSvc,
		riskTimelineSvc:              riskTimelineSvc,
		controlStatusSvc:             controlStatusSvc,
		controlSvc:                   controlSvc,
	}
}

func (ries RiskResidualEvaluationService) GetFullDetailsByID(riskID uint) (*riskDTOs.RiskResidualEvaluationDetailsDTO, error) {
	// Check if risk register exists
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	// Get full details from the repository
	riskResidualEvaluation, err := ries.repo.GetFullDetailsByID(riskID)
	if err != nil {
		return nil, err
	}

	// If no residual evaluation exists, return nil
	if riskResidualEvaluation == nil {
		return nil, nil
	}

	var riskResidualEvaluationDetailsDTO riskDTOs.RiskResidualEvaluationDetailsDTO

	// Map the results to the DTO
	err = mapstructure.Decode(riskResidualEvaluation, &riskResidualEvaluationDetailsDTO)
	if err != nil {
		return nil, err
	}

	return &riskResidualEvaluationDetailsDTO, nil
}

func (ries *RiskResidualEvaluationService) GetByID(riskID uint, id uint) (*models.RiskEvaluation, error) {
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskResidualEvaluation, err := ries.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskResidualEvaluation, nil
}

func (ries *RiskResidualEvaluationService) Create(riskID uint, currentUserID uint, residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO, tx ...utility.DatabaseTransaction) (uint, error) {
	var transaction utility.DatabaseTransaction

	if len(tx) > 0 && tx[0] != nil {
		transaction = tx[0]
	} else {
		transaction = ries.repo.BeginTransaction()
		defer func() {
			if r := recover(); r != nil {
				transaction.Rollback()
			}
		}()
	}
	if _, err := ries.riskRegisterSvc.GetByID(riskID, transaction); err != nil {
		transaction.Rollback()
		return 0, err
	}

	if err := ries.validateRiskResidualEvaluationFormDTO(residualEvaluationFormDTO); err != nil {
		transaction.Rollback()
		return 0, err
	}

	// Get inherent evaluation to use as base for residual calculation
	inherentEvaluation, err := ries.riskInherentEvaluationSvc.GetFullDetailsByID(riskID, transaction)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	if inherentEvaluation == nil {
		transaction.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("inherent evaluation must be completed before creating residual evaluation"), "not_exist", "inherent evaluation")
	}

	// Get control analyses to calculate average risk reduction
	controlAnalyses, err := ries.repo.GetControlAnalysesByRiskID(riskID, transaction)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	// Calculate average risk reduction from control analyses
	var totalRiskReduction float64
	var controlCount int
	for _, ca := range controlAnalyses {
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

	riskMatrixType, err := ries.riskMatrixSvc.GetByType("erm")
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	riskMatrix, err := ries.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

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

	// Use inherent evaluation's likelihood and impact as base (for guideline purposes)
	// Calculate residual score: Inherent Risk Rating Value × Average Risk Reduction
	var residualScore uint
	if averageRiskReduction != nil {
		// Residual Risk Score = Inherent Risk Rating Value × (1-Average Risk Reduction / 100)
		// Average Risk Reduction is stored as percentage (0-100), so divide by 100
		// Add 0.5 for proper rounding (e.g., 2.75 becomes 3 instead of 2)
		residualScore = uint(float64(inherentEvaluation.RiskMatrixRiskRatingValue.Value)*(1-*averageRiskReduction/100) + 0.5)
	} else {
		// If no control analysis exists, residual = inherent
		residualScore = inherentEvaluation.RiskMatrixRiskRatingValue.Value
	}

	riskRating, found := ries.riskMatrixRiskRatingValueSvc.FindClosestRiskRatingValue(riskRatingValueArray, residualScore)
	if !found || riskRating == nil {
		transaction.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("no risk rating found for score %d", residualScore), "not_found", "risk rating")
	}

	// Use form values if provided, otherwise inherit from inherent evaluation
	likelihoodID := inherentEvaluation.RiskMatrixLikelihood.ID
	impactID := inherentEvaluation.RiskMatrixImpact.ID

	if residualEvaluationFormDTO.RiskMatrixLikelihoodID > 0 {
		likelihoodID = residualEvaluationFormDTO.RiskMatrixLikelihoodID
	}

	if residualEvaluationFormDTO.RiskMatrixImpactID > 0 {
		impactID = residualEvaluationFormDTO.RiskMatrixImpactID
	}

	isLatest := true
	riskResidualEvaluation := &models.RiskEvaluation{
		RiskRegisterID:                   riskID,
		Type:                             "residual",
		RiskMatrixLikelihoodID:           likelihoodID,
		RiskMatrixImpactID:               impactID,
		RiskMatrixRiskRatingID:           riskRating.RiskMatrixRiskRatingID,
		RiskMatrixRiskRatingValueID:      riskRating.RiskMatrixRiskRatingValueID,
		AverageControlEffectivenessScore: averageRiskReduction,
		CreatedByID:                      currentUserID,
		IsLatest:                         &isLatest,
	}

	id, err := ries.repo.Create(riskResidualEvaluation, transaction)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	// Create history record after creating residual evaluation
	history := &models.RiskResidualEvaluationHistory{
		RiskEvaluationID:                 id,
		RiskRegisterID:                   riskID,
		RiskMatrixLikelihoodID:           likelihoodID,
		RiskMatrixRiskRatingID:           riskRating.RiskMatrixRiskRatingID,
		RiskMatrixRiskRatingValueID:      riskRating.RiskMatrixRiskRatingValueID,
		AverageControlEffectivenessScore: averageRiskReduction,
		CreatedBYID:                      currentUserID,
		CreatedAt:                        time.Now(),
	}

	if err := ries.repo.CreateHistory(history, transaction); err != nil {
		transaction.Rollback()
		return 0, err
	}

	err = ries.riskRegisterSvc.UpdateRiskRating(riskID, currentUserID, riskResidualEvaluation, transaction, "residual")
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	for _, impactEval := range residualEvaluationFormDTO.RiskImpactAnalyses {
		riskImpactEvaluation := &models.RiskImpactAnalysis{
			RiskEvaluationID:     id,
			RiskImpactCategoryID: impactEval.RiskImpactCategoryID,
			RiskMatrixImpactID:   impactEval.RiskMatrixImpactID,
			CreatedByID:          currentUserID,
		}

		err := ries.repo.CreateRiskImpactAnalysis(riskImpactEvaluation, transaction)
		if err != nil {
			transaction.Rollback()
			return 0, err
		}

		// Auto-fetch guidelines based on category and impact combination
		guidelines, err := ries.riskMatrixImpactGuidelineSvc.GetByCategoryAndImpact(impactEval.RiskImpactCategoryID, impactEval.RiskMatrixImpactID)
		if err != nil {
			transaction.Rollback()
			return 0, err
		}

		// Store all fetched guidelines
		for _, guideline := range guidelines {
			riskImpactAnalysisGuideline := &models.RiskImpactAnalysisGuideline{
				RiskImpactAnalysisID:        riskImpactEvaluation.ID,
				RiskMatrixImpactGuidelineID: guideline.ID,
			}
			err := ries.repo.CreateRiskImpactAnalysisGuideline(riskImpactAnalysisGuideline, transaction)
			if err != nil {
				transaction.Rollback()
				return 0, err
			}
		}

	}

	var riskAppetiteAlignments []*models.RiskAppetiteAlignment
	for _, appetiteAlignment := range residualEvaluationFormDTO.RiskAppetiteAlignments {
		riskAppetiteAlignment := &models.RiskAppetiteAlignment{
			RiskEvaluationID:             id,
			RiskAppetiteStatementID:      appetiteAlignment.RiskAppetiteStatementID,
			RiskAppetiteToleranceLevelID: appetiteAlignment.RiskAppetiteToleranceLevelID,
			CreatedByID:                  currentUserID,
		}
		riskAppetiteAlignments = append(riskAppetiteAlignments, riskAppetiteAlignment)
	}

	if len(riskAppetiteAlignments) > 0 {
		if err := ries.repo.CreateRiskAppetiteAlignments(riskAppetiteAlignments, transaction); err != nil {
			transaction.Rollback()
			return 0, err
		}
	}
	// Update timeline
	// todo
	timelineType, err := ries.riskTimelineSvc.GetByType("residual_risk_analysis")
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	timelineStatus, err := ries.riskTimelineSvc.GetByStatusType("completed")
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	currentTime := time.Now()
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    &currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      &currentTime,
	}
	if err := ries.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, transaction); err != nil {
		transaction.Rollback()
		return 0, err
	}

	// Update control status based on residual risk rating
	// If residual risk is critical or high = needs improvement, otherwise = doesn't need improvement
	var controlStatusType string
	if riskRating.RiskRatingTitleType == "critical" || riskRating.RiskRatingTitleType == "high" {
		controlStatusType = "needs_improvement"
	} else {
		controlStatusType = "doesnt_need_improvement"
	}

	controlStatus, err := ries.controlStatusSvc.GetByType(controlStatusType)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	// Update all controls associated with this risk
	for i, ca := range controlAnalyses {
		if ca.RiskControl != nil && ca.RiskControl.Control != nil {
			controlID := ca.RiskControl.Control.ID
			if err := ries.controlSvc.UpdateControlStatus(controlID, controlStatus.ID, transaction); err != nil {
				transaction.Rollback()
				return 0, err
			}
		} else {
			fmt.Printf("Skipping Control Analysis %d: RiskControl or Control is nil\n", i)
		}
	}

	if len(tx) == 0 {
		transaction.Commit()
	}
	// Perform activity logging
	go func() {
		existingRiskResidualEvaluation, err := ries.GetFullDetailsByID(riskID)
		if err != nil {
			utility.LogError("Failed to fetch risk residual evaluation", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ries.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskEvaluation{}.TableName(),
			id,
			existingRiskResidualEvaluation,
			currentUserID,
			"risk", "risk_register", "residual_evaluation", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskEvaluation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ries *RiskResidualEvaluationService) Update(riskID uint, id uint, currentUserID uint, residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO) error {
	tx := ries.repo.BeginTransaction()

	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	if err := ries.validateRiskResidualEvaluationFormDTO(residualEvaluationFormDTO); err != nil {
		return err
	}

	// Get inherent evaluation to use as base for residual calculation
	inherentEvaluation, err := ries.riskInherentEvaluationSvc.GetFullDetailsByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if inherentEvaluation == nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("inherent evaluation must be completed before updating residual evaluation"), "not_exist", "inherent evaluation")
	}

	// Get control analyses to calculate average risk reduction
	controlAnalyses, err := ries.repo.GetControlAnalysesByRiskID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Calculate average risk reduction from control analyses
	var totalRiskReduction float64
	var controlCount int
	for _, ca := range controlAnalyses {
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

	riskMatrixType, err := ries.riskMatrixSvc.GetByType("erm")
	if err != nil {
		tx.Rollback()
		return err
	}

	riskMatrix, err := ries.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldRiskEvaluation, err := ries.GetFullDetailsByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	var existingImpactAnalysis []map[string]interface{}
	for _, impactAnalysis := range oldRiskEvaluation.RiskImpactAnalyses {
		existingImpactAnalysis = append(existingImpactAnalysis, map[string]interface{}{
			"id":                      impactAnalysis.ID,
			"risk_evaluation_id":      impactAnalysis.RiskEvaluationID,
			"risk_impact_category_id": impactAnalysis.RiskImpactCategory.ID,
			"risk_matrix_impact_id":   impactAnalysis.RiskMatrixImpact.ID,
		})
	}

	var existingAppetiteAlignment []map[string]interface{}
	for _, appetiteAlignment := range oldRiskEvaluation.RiskAppetiteAlignments {
		existingAppetiteAlignment = append(existingAppetiteAlignment, map[string]interface{}{
			"id":                               appetiteAlignment.ID,
			"risk_evaluation_id":               appetiteAlignment.RiskEvaluationID,
			"risk_appetite_statement_id":       appetiteAlignment.RiskAppetiteStatementID,
			"risk_appetite_tolerance_level_id": appetiteAlignment.RiskAppetiteToleranceLevelID,
		})
	}

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

	// Use inherent evaluation's likelihood and impact as base (for guideline purposes)
	// Calculate residual score: Inherent Risk Rating Value × Average Risk Reduction
	var residualScore uint
	if averageRiskReduction != nil {
		// Residual Risk Score = Inherent Risk Rating Value × (1 - Average Risk Reduction / 100)
		// Average Risk Reduction is stored as percentage (0-100), so divide by 100
		// For example: 20 (critical) × (1 - 75/100) = 20 × 0.25 = 5 (medium)
		// Add 0.5 for proper rounding (e.g., 2.75 becomes 3 instead of 2)
		residualScore = uint(float64(inherentEvaluation.RiskMatrixRiskRatingValue.Value)*(1-*averageRiskReduction/100) + 0.5)
	} else {
		// If no control analysis exists, residual = inherent
		residualScore = inherentEvaluation.RiskMatrixRiskRatingValue.Value
	}

	riskRating, found := ries.riskMatrixRiskRatingValueSvc.FindClosestRiskRatingValue(riskRatingValueArray, residualScore)
	if !found || riskRating == nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("no risk rating found for score %d", residualScore), "not_found", "risk rating")
	}

	existingRiskResidualEvaluation, err := ries.GetByID(riskID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingRiskResidualEvaluation.RiskMatrixImpactID = inherentEvaluation.RiskMatrixImpact.ID
	existingRiskResidualEvaluation.RiskMatrixLikelihoodID = inherentEvaluation.RiskMatrixLikelihood.ID
	existingRiskResidualEvaluation.RiskMatrixRiskRatingID = riskRating.RiskMatrixRiskRatingID
	existingRiskResidualEvaluation.RiskMatrixRiskRatingValueID = riskRating.RiskMatrixRiskRatingValueID
	existingRiskResidualEvaluation.AverageControlEffectivenessScore = averageRiskReduction
	existingRiskResidualEvaluation.UpdatedByID = &currentUserID

	if err := ries.repo.Update(existingRiskResidualEvaluation, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create another history record after update to capture new values
	historyAfterUpdate := &models.RiskResidualEvaluationHistory{
		RiskEvaluationID:                 existingRiskResidualEvaluation.ID,
		RiskRegisterID:                   existingRiskResidualEvaluation.RiskRegisterID,
		RiskMatrixLikelihoodID:           existingRiskResidualEvaluation.RiskMatrixLikelihoodID,
		RiskMatrixRiskRatingID:           existingRiskResidualEvaluation.RiskMatrixRiskRatingID,
		RiskMatrixRiskRatingValueID:      existingRiskResidualEvaluation.RiskMatrixRiskRatingValueID,
		AverageControlEffectivenessScore: existingRiskResidualEvaluation.AverageControlEffectivenessScore,
		CreatedBYID:                      currentUserID,
		CreatedAt:                        time.Now(),
	}

	if err := ries.repo.CreateHistory(historyAfterUpdate, tx); err != nil {
		tx.Rollback()
		return err
	}

	err = ries.riskRegisterSvc.UpdateRiskRating(riskID, currentUserID, existingRiskResidualEvaluation, tx, "residual")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing RiskImpactAnalyses and related RiskMatrixImpactGuidelines
	if err := ries.repo.DeleteRiskImpactAnalysesByRiskEvaluationID(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing RiskAppetiteAlignments
	if err := ries.repo.DeleteRiskAppetiteAlignmentsByRiskEvaluationID(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	var riskImpactAnalysis []*models.RiskImpactAnalysis
	for _, impactEval := range residualEvaluationFormDTO.RiskImpactAnalyses {
		riskImpactEvaluation := &models.RiskImpactAnalysis{
			RiskEvaluationID:     id,
			RiskImpactCategoryID: impactEval.RiskImpactCategoryID,
			RiskMatrixImpactID:   impactEval.RiskMatrixImpactID,
			CreatedByID:          currentUserID,
		}
		riskImpactAnalysis = append(riskImpactAnalysis, riskImpactEvaluation)
		err := ries.repo.CreateRiskImpactAnalysis(riskImpactEvaluation, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Auto-fetch guidelines based on category and impact combination
		guidelines, err := ries.riskMatrixImpactGuidelineSvc.GetByCategoryAndImpact(impactEval.RiskImpactCategoryID, impactEval.RiskMatrixImpactID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Store all fetched guidelines
		for _, guideline := range guidelines {
			riskImpactAnalysisGuideline := &models.RiskImpactAnalysisGuideline{
				RiskImpactAnalysisID:        riskImpactEvaluation.ID,
				RiskMatrixImpactGuidelineID: guideline.ID,
			}
			err := ries.repo.CreateRiskImpactAnalysisGuideline(riskImpactAnalysisGuideline, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	var riskAppetiteAlignments []*models.RiskAppetiteAlignment
	for _, appetiteAlignment := range residualEvaluationFormDTO.RiskAppetiteAlignments {
		riskAppetiteAlignment := &models.RiskAppetiteAlignment{
			RiskEvaluationID:             id,
			RiskAppetiteStatementID:      appetiteAlignment.RiskAppetiteStatementID,
			RiskAppetiteToleranceLevelID: appetiteAlignment.RiskAppetiteToleranceLevelID,
			CreatedByID:                  currentUserID,
		}
		riskAppetiteAlignments = append(riskAppetiteAlignments, riskAppetiteAlignment)
	}

	if len(riskAppetiteAlignments) > 0 {
		if err := ries.repo.CreateRiskAppetiteAlignments(riskAppetiteAlignments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update control status based on residual risk rating
	// If residual risk is critical or high = needs improvement, otherwise = doesn't need improvement
	var controlStatusType string
	if riskRating.RiskRatingTitleType == "critical" || riskRating.RiskRatingTitleType == "high" {
		controlStatusType = "needs_improvement"
	} else {
		controlStatusType = "doesnt_need_improvement"
	}

	controlStatus, err := ries.controlStatusSvc.GetByType(controlStatusType)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update all controls associated with this risk
	for _, ca := range controlAnalyses {
		if ca.RiskControl != nil && ca.RiskControl.Control != nil {
			controlID := ca.RiskControl.Control.ID
			if err := ries.controlSvc.UpdateControlStatus(controlID, controlStatus.ID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	go func() {
		existingRiskEvaluation, err := ries.GetFullDetailsByID(riskID)
		if err != nil {
			return
		}
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"risk_matrix_impact":            oldRiskEvaluation.RiskMatrixImpact.TitleEn,
			"risk_matrix_likelihood":        oldRiskEvaluation.RiskMatrixLikelihood.TitleEn,
			"risk_matrix_risk_rating":       oldRiskEvaluation.RiskMatrixRiskRating.TitleEn,
			"risk_matirx_risk_rating_value": oldRiskEvaluation.RiskMatrixRiskRatingValue.Value,
		}

		newData := map[string]interface{}{
			"risk_matrix_impact":            existingRiskEvaluation.RiskMatrixImpact.TitleEn,
			"risk_matrix_likelihood":        existingRiskEvaluation.RiskMatrixLikelihood.TitleEn,
			"risk_matrix_risk_rating":       existingRiskEvaluation.RiskMatrixRiskRating.TitleEn,
			"risk_matirx_risk_rating_value": existingRiskEvaluation.RiskMatrixRiskRatingValue.Value,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(riskImpactAnalysis) > 0 || len(existingImpactAnalysis) > 0 {
			changedData["impact_analysis"] = map[string]interface{}{
				"added":   riskImpactAnalysis,
				"removed": existingImpactAnalysis,
			}
		}

		if len(riskAppetiteAlignments) > 0 || len(existingAppetiteAlignment) > 0 {
			changedData["appetite_analysis"] = map[string]interface{}{
				"added":   riskAppetiteAlignments,
				"removed": existingAppetiteAlignment,
			}
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err = ries.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskEvaluation{}.TableName(),
				id,
				changedData,
				currentUserID,
				"risk", "risk_register", "residual_evaluation", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.RiskEvaluation{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (ries *RiskResidualEvaluationService) Delete(riskID uint, id uint, currentUserID uint) (uint, error) {
	tx := ries.repo.BeginTransaction()

	var deletedResidualEvaluation []map[string]interface{}

	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}
	var failedId uint

	if err := ries.repo.DeleteRiskImpactAnalysesByRiskEvaluationID(id, tx); err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}
	if err := ries.repo.DeleteRiskAppetiteAlignmentsByRiskEvaluationID(id, tx); err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}
	// Delete-> update timeline performedAt as null
	// todo
	timelineType, err := ries.riskTimelineSvc.GetByType("residual_risk_analysis")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	timelineStatus, err := ries.riskTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    nil,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      nil,
	}
	if err := ries.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := ries.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}

	inherentEvaluationDetailsDTO, err := ries.riskInherentEvaluationSvc.GetFullDetailsByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	updateDTO := &riskDTOs.RiskRegisterRiskRatingUpdateDTO{
		InherentImpactID:          &inherentEvaluationDetailsDTO.RiskMatrixImpact.ID,
		InherentLikelihoodID:      &inherentEvaluationDetailsDTO.RiskMatrixLikelihood.ID,
		InherentRiskRatingID:      &inherentEvaluationDetailsDTO.RiskMatrixRiskRating.ID,
		InherentRiskRatingValueID: &inherentEvaluationDetailsDTO.RiskMatrixRiskRatingValue.ID,
		ResidualImpactID:          nil,
		ResidualLikelihoodID:      nil,
		ResidualRiskRatingID:      nil,
		ResidualRiskRatingValueID: nil,
		CurrentImpactID:           &inherentEvaluationDetailsDTO.RiskMatrixImpact.ID,
		CurrentLikelihoodID:       &inherentEvaluationDetailsDTO.RiskMatrixLikelihood.ID,
		CurrentRiskRatingID:       &inherentEvaluationDetailsDTO.RiskMatrixRiskRating.ID,
		CurrentRiskRatingValueID:  &inherentEvaluationDetailsDTO.RiskMatrixRiskRatingValue.ID,
	}

	err = ries.riskRegisterSvc.ClearRiskRating(riskID, updateDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	// Log activity
	go func() {
		deletedResidualEvaluation = append(deletedResidualEvaluation, map[string]interface{}{
			"id":               id,
			"risk_register_id": riskID,
		})

		err := ries.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskEvaluation{}.TableName(),
			id,
			deletedResidualEvaluation,
			currentUserID,
			"risk", "risk_register", "residual_evaluation", nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}

func (ries RiskResidualEvaluationService) GetHistoryByRiskID(riskID uint) ([]*riskDTOs.RiskResidualEvaluationHistoryDTO, error) {
	// Check if risk register exists
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	// Get residual evaluation for this risk
	residualEvaluation, err := ries.repo.GetFullDetailsByID(riskID)
	if err != nil {
		return nil, err
	}

	if residualEvaluation == nil {
		return []*riskDTOs.RiskResidualEvaluationHistoryDTO{}, nil
	}

	// Get history records
	histories, err := ries.repo.GetHistoryByRiskID(riskID)
	if err != nil {
		return nil, err
	}

	// Map histories to DTOs
	var historyDTOs []*riskDTOs.RiskResidualEvaluationHistoryDTO

	err = mapstructure.Decode(histories, &historyDTOs)
	if err != nil {
		return nil, err
	}

	// Format created_at (because mapstructure won't format time)
	for _, dto := range historyDTOs {
		for _, history := range histories {
			if dto.ID == history.ID {
				dto.CreatedAtString = history.CreatedAt.Format("2006-01-02")
				break
			}
		}
	}

	return historyDTOs, nil
}

func (ries *RiskResidualEvaluationService) validateRiskResidualEvaluationFormDTO(residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO) error {
	// Only validate RiskMatrixLikelihoodID and RiskMatrixImpactID if they are provided
	// Otherwise they will be inherited from inherent evaluation
	if residualEvaluationFormDTO.RiskMatrixImpactID > 0 {
		if _, err := ries.riskMatrixImpactSvc.GetByID(uint(residualEvaluationFormDTO.RiskMatrixImpactID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk matrix impact does not exist"), "not_exist", "risk matrix impact")
		}
	}

	if residualEvaluationFormDTO.RiskMatrixLikelihoodID > 0 {
		if _, err := ries.riskMatrixLikelihoodSvc.GetByID(uint(residualEvaluationFormDTO.RiskMatrixLikelihoodID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk matrix likelihood does not exist"), "not_exist", "risk matrix likelihood")
		}
	}

	for _, appetiteAlignment := range residualEvaluationFormDTO.RiskAppetiteAlignments {
		if _, err := ries.riskAppetiteStatementSvc.GetByID(uint(appetiteAlignment.RiskAppetiteStatementID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk appetite statement does not exist"), "not_exist", "risk appetite statement")
		}
	}

	for _, impactEval := range residualEvaluationFormDTO.RiskImpactAnalyses {

		if _, err := ries.riskImpactCategorySvc.GetByID(uint(impactEval.RiskImpactCategoryID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk impact category does not exist"), "not_exist", "risk impact category")
		}
		if _, err := ries.riskMatrixImpactSvc.GetByID(uint(impactEval.RiskMatrixImpactID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk matrix impact does not exist"), "not_exist", "risk matrix impact")
		}
	}

	return nil
}
