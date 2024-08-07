package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskResidualEvaluationServiceInterface interface {
	GetFullDetailsByID(riskID uint) (*riskDTOs.RiskResidualEvaluationDetailsDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskEvaluation, error)
	Create(riskID uint, currentUserID uint, residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
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
	}
}

func (ries RiskResidualEvaluationService) GetFullDetailsByID(riskID uint) (*riskDTOs.RiskResidualEvaluationDetailsDTO, error) {
	// Check if risk register exists
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	var riskResidualEvaluationDetailsDTO riskDTOs.RiskResidualEvaluationDetailsDTO

	// Get full details from the repository
	riskResidualEvaluation, err := ries.repo.GetFullDetailsByID(riskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskID)
		}
		return nil, err
	}

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

func (ries *RiskResidualEvaluationService) Create(riskID uint, currentUserID uint, residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO) (uint, error) {
	tx := ries.repo.BeginTransaction()

	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	if err := ries.validateRiskResidualEvaluationFormDTO(residualEvaluationFormDTO); err != nil {
		return 0, err
	}

	riskMatrixType, err := ries.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	riskMatrix, err := ries.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
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

	impact, err := ries.riskMatrixImpactSvc.GetByID(residualEvaluationFormDTO.RiskMatrixImpactID)
	if err != nil {
		return 0, err
	}
	likelihood, err := ries.riskMatrixLikelihoodSvc.GetByID(residualEvaluationFormDTO.RiskMatrixLikelihoodID)
	if err != nil {
		return 0, err
	}
	calculationMethod := riskMatrix.CalculationType

	score := ries.riskMatrixRiskRatingValueSvc.EvaluateScore(impact, likelihood, calculationMethod)
	riskRating, _ := ries.riskMatrixRiskRatingValueSvc.FindRiskRatingValue(riskRatingValueArray, score)

	riskResidualEvaluation := &models.RiskEvaluation{
		RiskRegisterID:              riskID,
		Type:                        "residual",
		RiskMatrixLikelihoodID:      residualEvaluationFormDTO.RiskMatrixLikelihoodID,
		RiskMatrixImpactID:          residualEvaluationFormDTO.RiskMatrixImpactID,
		RiskMatrixRiskRatingID:      riskRating.RiskMatrixRiskRatingID,
		RiskMatrixRiskRatingValueID: riskRating.RiskMatrixRiskRatingValueID,
		CreatedByID:                 currentUserID,
	}

	id, err := ries.repo.Create(riskResidualEvaluation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = ries.riskRegisterSvc.UpdateRiskRating(riskID, currentUserID, riskResidualEvaluation, tx, "residual")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, impactEval := range residualEvaluationFormDTO.RiskImpactAnalyses {
		riskImpactEvaluation := &models.RiskImpactAnalysis{
			RiskEvaluationID:     id,
			RiskImpactCategoryID: impactEval.RiskImpactCategoryID,
			RiskMatrixImpactID:   impactEval.RiskMatrixImpactID,
			CreatedByID:          currentUserID,
		}

		err := ries.repo.CreateRiskImpactAnalysis(riskImpactEvaluation, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, guidelineID := range impactEval.RiskMatrixImpactGuidelineIDs {
			riskImpactAnalysisGuideline := &models.RiskImpactAnalysisGuideline{
				RiskImpactAnalysisID:        riskImpactEvaluation.ID,
				RiskMatrixImpactGuidelineID: guidelineID,
			}
			err := ries.repo.CreateRiskImpactAnalysisGuideline(riskImpactAnalysisGuideline, tx)
			if err != nil {
				tx.Rollback()
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
		if err := ries.repo.CreateRiskAppetiteAlignments(riskAppetiteAlignments, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	//Update timeline
	//todo
	timelineType, err := ries.riskTimelineSvc.GetByType("residual_risk_analysis")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	timelineStatus, err := ries.riskTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	currentTime := time.Now()
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    &currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      &currentTime,
	}
	if err := ries.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	// Perform activity logging
	go func() {
		err := ries.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskEvaluation{}.TableName(),
			id,
			residualEvaluationFormDTO,
			currentUserID,
			nil,
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

	riskMatrixType, err := ries.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return err
	}

	riskMatrix, err := ries.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		return err
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

	impact, err := ries.riskMatrixImpactSvc.GetByID(residualEvaluationFormDTO.RiskMatrixImpactID)
	if err != nil {
		return err
	}
	likelihood, err := ries.riskMatrixLikelihoodSvc.GetByID(residualEvaluationFormDTO.RiskMatrixLikelihoodID)
	if err != nil {
		return err
	}
	calculationMethod := riskMatrix.CalculationType

	score := ries.riskMatrixRiskRatingValueSvc.EvaluateScore(impact, likelihood, calculationMethod)
	riskRating, _ := ries.riskMatrixRiskRatingValueSvc.FindRiskRatingValue(riskRatingValueArray, score)

	existingRiskResidualEvaluation, err := ries.GetByID(riskID, id)
	if err != nil {
		return err
	}

	existingRiskResidualEvaluation.RiskMatrixImpactID = residualEvaluationFormDTO.RiskMatrixImpactID
	existingRiskResidualEvaluation.RiskMatrixLikelihoodID = residualEvaluationFormDTO.RiskMatrixLikelihoodID
	existingRiskResidualEvaluation.RiskMatrixRiskRatingID = riskRating.RiskMatrixRiskRatingID
	existingRiskResidualEvaluation.RiskMatrixRiskRatingValueID = riskRating.RiskMatrixRiskRatingValueID
	existingRiskResidualEvaluation.UpdatedByID = &currentUserID

	if err := ries.repo.Update(existingRiskResidualEvaluation, tx); err != nil {
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

	for _, impactEval := range residualEvaluationFormDTO.RiskImpactAnalyses {
		riskImpactEvaluation := &models.RiskImpactAnalysis{
			RiskEvaluationID:     id,
			RiskImpactCategoryID: impactEval.RiskImpactCategoryID,
			RiskMatrixImpactID:   impactEval.RiskMatrixImpactID,
			CreatedByID:          currentUserID,
		}

		err := ries.repo.CreateRiskImpactAnalysis(riskImpactEvaluation, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, guidelineID := range impactEval.RiskMatrixImpactGuidelineIDs {
			riskImpactAnalysisGuideline := &models.RiskImpactAnalysisGuideline{
				RiskImpactAnalysisID:        riskImpactEvaluation.ID,
				RiskMatrixImpactGuidelineID: guidelineID,
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
	tx.Commit()
	return nil
}

func (ries *RiskResidualEvaluationService) Delete(riskID uint, id uint, currentUserID uint) (uint, error) {
	tx := ries.repo.BeginTransaction()
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
	//Delete-> update timeline performedAt as null
	//todo
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

	inherentEvaluationDetailsDTO, err := ries.riskInherentEvaluationSvc.GetFullDetailsByID(riskID,tx)
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
		err := ries.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskEvaluation{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}

func (ries *RiskResidualEvaluationService) validateRiskResidualEvaluationFormDTO(residualEvaluationFormDTO *riskDTOs.RiskResidualEvaluationFormDTO) error {
	if _, err := ries.riskMatrixImpactSvc.GetByID(uint(residualEvaluationFormDTO.RiskMatrixImpactID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk matrix impact does not exist"), "not_exist", "risk matrix impact")
	}

	if _, err := ries.riskMatrixLikelihoodSvc.GetByID(uint(residualEvaluationFormDTO.RiskMatrixLikelihoodID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk matrix likelihood does not exist"), "not_exist", "risk matrix likelihood")
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
		for _, impactGuidelineID := range impactEval.RiskMatrixImpactGuidelineIDs {
			if _, err := ries.riskMatrixImpactGuidelineSvc.GetByID(uint(impactGuidelineID)); err != nil {
				return utility.NewValidationError(fmt.Errorf("risk matrix impact guideline does not exist"), "not_exist", "risk matrix impact guideline")
			}
		}
	}

	return nil
}
