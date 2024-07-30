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

type RiskInherentEvaluationServiceInterface interface {
	GetFullDetailsByID(riskID uint) (*riskDTOs.RiskInherentEvaluationDetailsDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskEvaluation, error)
	Create(riskID uint, currentUserID uint, inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
}

// Service and implementations
type RiskInherentEvaluationService struct {
	repo                         riskRepositories.RiskInherentEvaluationRepositoryInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	riskRegisterSvc              RiskRegisterServiceInterface
	riskMatrixSvc                RiskMatrixServiceInterface
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface
	riskMatrixImpactSvc          RiskMatrixImpactServiceInterface
	riskMatrixLikelihoodSvc      RiskMatrixLikelihoodServiceInterface
	riskImpactCategorySvc        RiskImpactCategoryServiceInterface
	riskMatrixImpactGuidelineSvc RiskMatrixImpactGuidelineServiceInterface
	riskAppetiteStatementSvc     RiskAppetiteStatementServiceInterface
	riskTimelineSvc              RiskTimelineServiceInterface
}

func NewRiskInherentEvaluationService(
	repo riskRepositories.RiskInherentEvaluationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface,
	riskMatrixImpactSvc RiskMatrixImpactServiceInterface,
	riskMatrixLikelihoodSvc RiskMatrixLikelihoodServiceInterface,
	riskImpactCategorySvc RiskImpactCategoryServiceInterface,
	riskMatrixImpactGuidelineSvc RiskMatrixImpactGuidelineServiceInterface,
	riskAppetiteStatementSvc RiskAppetiteStatementServiceInterface,
	riskTimelineSvc RiskTimelineServiceInterface,
) RiskInherentEvaluationServiceInterface {
	return &RiskInherentEvaluationService{
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
		riskTimelineSvc:              riskTimelineSvc,
	}
}

func (ries RiskInherentEvaluationService) GetFullDetailsByID(riskID uint) (*riskDTOs.RiskInherentEvaluationDetailsDTO, error) {
	// Check if risk register exists
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	var riskInherentEvaluationDetailsDTO riskDTOs.RiskInherentEvaluationDetailsDTO

	// Get full details from the repository
	riskInherentEvaluation, err := ries.repo.GetFullDetailsByID(riskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskID)
		}
		return nil, err
	}

	// Map the results to the DTO
	err = mapstructure.Decode(riskInherentEvaluation, &riskInherentEvaluationDetailsDTO)
	if err != nil {
		return nil, err
	}

	return &riskInherentEvaluationDetailsDTO, nil
}

func (ries *RiskInherentEvaluationService) GetByID(riskID uint, id uint) (*models.RiskEvaluation, error) {
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskInherentEvaluation, err := ries.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskInherentEvaluation, nil
}

func (ries *RiskInherentEvaluationService) Create(riskID uint, currentUserID uint, inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) (uint, error) {
	tx := ries.repo.BeginTransaction()

	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	if err := ries.validateRiskInherentEvaluationFormDTO(inherentEvaluationFormDTO); err != nil {
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

	impact, err := ries.riskMatrixImpactSvc.GetByID(inherentEvaluationFormDTO.RiskMatrixImpactID)
	if err != nil {
		return 0, err
	}
	likelihood, err := ries.riskMatrixLikelihoodSvc.GetByID(inherentEvaluationFormDTO.RiskMatrixLikelihoodID)
	if err != nil {
		return 0, err
	}
	calculationMethod := riskMatrix.CalculationType

	score := ries.riskMatrixRiskRatingValueSvc.EvaluateScore(impact, likelihood, calculationMethod)
	riskRating, _ := ries.riskMatrixRiskRatingValueSvc.FindRiskRatingValue(riskRatingValueArray, score)

	riskInherentEvaluation := &models.RiskEvaluation{
		RiskRegisterID:              riskID,
		Type:                        "inherent",
		RiskMatrixLikelihoodID:      inherentEvaluationFormDTO.RiskMatrixLikelihoodID,
		RiskMatrixImpactID:          inherentEvaluationFormDTO.RiskMatrixImpactID,
		RiskMatrixRiskRatingID:      riskRating.RiskMatrixRiskRatingID,
		RiskMatrixRiskRatingValueID: riskRating.RiskMatrixRiskRatingValueID,
		CreatedByID:                 currentUserID,
	}

	id, err := ries.repo.Create(riskInherentEvaluation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//update current risk rating of risk register
	err = ries.riskRegisterSvc.UpdateRiskRating(riskID, currentUserID, riskInherentEvaluation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, impactEval := range inherentEvaluationFormDTO.RiskImpactAnalyses {
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
	for _, appetiteAlignment := range inherentEvaluationFormDTO.RiskAppetiteAlignments {
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
	timelineType, err := ries.riskTimelineSvc.GetByType("risk_analysis")
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
			inherentEvaluationFormDTO,
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

func (ries *RiskInherentEvaluationService) Update(riskID uint, id uint, currentUserID uint, inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) error {
	tx := ries.repo.BeginTransaction()

	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	if err := ries.validateRiskInherentEvaluationFormDTO(inherentEvaluationFormDTO); err != nil {
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

	impact, err := ries.riskMatrixImpactSvc.GetByID(inherentEvaluationFormDTO.RiskMatrixImpactID)
	if err != nil {
		return err
	}
	likelihood, err := ries.riskMatrixLikelihoodSvc.GetByID(inherentEvaluationFormDTO.RiskMatrixLikelihoodID)
	if err != nil {
		return err
	}
	calculationMethod := riskMatrix.CalculationType

	score := ries.riskMatrixRiskRatingValueSvc.EvaluateScore(impact, likelihood, calculationMethod)
	riskRating, _ := ries.riskMatrixRiskRatingValueSvc.FindRiskRatingValue(riskRatingValueArray, score)

	existingRiskInherentEvaluation, err := ries.GetByID(riskID, id)
	if err != nil {
		return err
	}

	existingRiskInherentEvaluation.RiskMatrixImpactID = inherentEvaluationFormDTO.RiskMatrixImpactID
	existingRiskInherentEvaluation.RiskMatrixLikelihoodID = inherentEvaluationFormDTO.RiskMatrixLikelihoodID
	existingRiskInherentEvaluation.RiskMatrixRiskRatingID = riskRating.RiskMatrixRiskRatingID
	existingRiskInherentEvaluation.RiskMatrixRiskRatingValueID = riskRating.RiskMatrixRiskRatingValueID
	existingRiskInherentEvaluation.UpdatedByID = &currentUserID

	if err := ries.repo.Update(existingRiskInherentEvaluation, tx); err != nil {
		tx.Rollback()
		return err
	}

	err = ries.riskRegisterSvc.UpdateRiskRating(riskID, currentUserID, existingRiskInherentEvaluation, tx)
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

	for _, impactEval := range inherentEvaluationFormDTO.RiskImpactAnalyses {
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
	for _, appetiteAlignment := range inherentEvaluationFormDTO.RiskAppetiteAlignments {
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

func (ries *RiskInherentEvaluationService) Delete(riskID uint, id uint, currentUserID uint) (uint, error) {
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

	if err := ries.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}

	updateDTO := &riskDTOs.RiskRegisterUpdateDTO{
		CurrentImpactID:          nil,
		CurrentLikelihoodID:      nil,
		CurrentRiskRatingID:      nil,
		CurrentRiskRatingValueID: nil,
	}
	err := ries.riskRegisterSvc.UpdateCurrentRiskRating(riskID, updateDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	//Delete-> update timeline performedAt as null
	//todo
	timelineType, err := ries.riskTimelineSvc.GetByType("risk_analysis")
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

func (ries *RiskInherentEvaluationService) validateRiskInherentEvaluationFormDTO(inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) error {
	if _, err := ries.riskMatrixImpactSvc.GetByID(uint(inherentEvaluationFormDTO.RiskMatrixImpactID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk matrix impact does not exist"), "not_exist", "risk matrix impact")
	}

	if _, err := ries.riskMatrixLikelihoodSvc.GetByID(uint(inherentEvaluationFormDTO.RiskMatrixLikelihoodID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk matrix likelihood does not exist"), "not_exist", "risk matrix likelihood")
	}

	for _, appetiteAlignment := range inherentEvaluationFormDTO.RiskAppetiteAlignments {
		if _, err := ries.riskAppetiteStatementSvc.GetByID(uint(appetiteAlignment.RiskAppetiteStatementID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk appetite statement does not exist"), "not_exist", "risk appetite statement")
		}
	}

	for _, impactEval := range inherentEvaluationFormDTO.RiskImpactAnalyses {
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
