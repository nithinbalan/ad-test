package services

import (
	"fmt"
	"time"

	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type AssetRiskInherentEvaluationServiceInterface interface {
	GetFullDetailsByID(riskID uint) (*assetDTOs.AssetRiskInherentEvaluationDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskEvaluation, error)
	Create(riskID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskInherentEvaluationFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskInherentEvaluationFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
}

type AssetRiskInherentEvaluationService struct {
	repo                 repositories.AssetRiskInherentEvaluationRepositoryInterface
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
}

func NewAssetRiskInherentEvaluationService(
	repo repositories.AssetRiskInherentEvaluationRepositoryInterface,
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) AssetRiskInherentEvaluationServiceInterface {
	return &AssetRiskInherentEvaluationService{
		repo:                 repo,
		assetRiskRegisterSvc: assetRiskRegisterSvc,
		activityLogSvc:       activityLogSvc,
	}
}

func (aries *AssetRiskInherentEvaluationService) GetFullDetailsByID(riskID uint) (*assetDTOs.AssetRiskInherentEvaluationDTO, error) {
	riskRegister, err := aries.assetRiskRegisterSvc.GetByID(riskID)
	if err != nil {
		return nil, err
	}

	riskEvaluation, err := aries.repo.GetFullDetailsByRiskID(riskID)
	if err != nil {
		return nil, err
	}
	if riskEvaluation == nil {
		return nil, nil
	}

	return toAssetRiskInherentEvaluationDTO(riskEvaluation, riskRegister), nil
}

func (aries *AssetRiskInherentEvaluationService) GetByID(riskID uint, id uint) (*models.RiskEvaluation, error) {
	if _, err := aries.assetRiskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskEvaluation, err := aries.repo.GetByID(id)
	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, id)
	}

	return riskEvaluation, nil
}

func (aries *AssetRiskInherentEvaluationService) Create(riskID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskInherentEvaluationFormDTO) (id uint, err error) {
	tx := aries.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	riskRegister, err := aries.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return 0, err
	}

	if riskRegister.InherentRiskRatingID != nil {
		return 0, utility.NewValidationError(fmt.Errorf("inherent evaluation already exists for this risk"), "unique", "inherent evaluation")
	}

	likelihood, err := aries.assetRiskRegisterSvc.GetLikelihoodByID(formDTO.RiskMatrixLikelihoodID, tx)
	if err != nil {
		return 0, err
	}

	if riskRegister.ConfidentialityScore == nil {
		return 0, fmt.Errorf("missing cofidentaility score in risk")
	}
	if riskRegister.IntegrityScore == nil {
		return 0, fmt.Errorf("missing integrity score in risk")
	}
	if riskRegister.AvailabilityScore == nil {
		return 0, fmt.Errorf("missing avaiability score in risk")
	}

	impactScore := max(*riskRegister.ConfidentialityScore, *riskRegister.IntegrityScore, *riskRegister.AvailabilityScore)

	overallImpact, err := aries.assetRiskRegisterSvc.GetImpactByScore(impactScore, tx)
	if err != nil {
		return 0, err
	}

	riskValue := int(overallImpact.Score * likelihood.Score)

	riskRatingValue, err := aries.assetRiskRegisterSvc.GetRiskRatingValueByValue(riskValue, tx)
	if err != nil {
		return 0, err
	}

	riskEvaluation := &models.RiskEvaluation{
		RiskRegisterID:              riskID,
		Type:                        "inherent",
		RiskMatrixLikelihoodID:      likelihood.ID,
		RiskMatrixImpactID:          overallImpact.ID,
		RiskMatrixRiskRatingID:      riskRatingValue.RiskMatrixRiskRatingID,
		RiskMatrixRiskRatingValueID: riskRatingValue.ID,
		CreatedByID:                 currentUserID,
	}

	id, err = aries.repo.Create(riskEvaluation, tx)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	riskRegister.InherentLikelihoodID = &likelihood.ID
	riskRegister.InherentImpactID = &overallImpact.ID
	riskRegister.InherentRiskRatingID = &riskRatingValue.RiskMatrixRiskRatingID
	riskRegister.InherentRiskRatingValueID = &riskRatingValue.ID
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	if err = aries.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return 0, err
	}

	timelineType, err := aries.assetRiskRegisterSvc.GetTimelineTypeByType("risk_analysis", tx)
	if err != nil {
		return 0, err
	}
	timelineStatus, err := aries.assetRiskRegisterSvc.GetTimelineStatusByType("completed", tx)
	if err != nil {
		return 0, err
	}
	riskTimeline, err := aries.assetRiskRegisterSvc.GetRiskTimelineByRiskAndType(riskID, timelineType.ID, tx)
	if err != nil {
		return 0, err
	}
	riskTimeline.TimelineStatusID = timelineStatus.ID
	riskTimeline.PerformedByID = &currentUserID
	riskTimeline.PerformedAt = &now

	if err = aries.assetRiskRegisterSvc.UpdateRiskTimeline(riskTimeline, tx); err != nil {
		return 0, err
	}

	go func() {
		existingEvaluation, logErr := aries.GetFullDetailsByID(riskID)
		if logErr != nil {
			utility.LogError("Failed to fetch risk inherent evaluation", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = aries.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskEvaluation{}.TableName(),
			id,
			existingEvaluation,
			currentUserID,
			"asset",
			"asset_risk",
			"risk",
			nil,
		)

		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.CreateAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return id, nil
}

func (aries *AssetRiskInherentEvaluationService) Update(riskID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskInherentEvaluationFormDTO) (err error) {
	tx := aries.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	riskRegister, err := aries.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return err
	}

	if riskRegister.ResidualRiskRatingID != nil {
		return utility.NewValidationError(fmt.Errorf("cannot update inherent evaluation after residual evaluation has been completed"), "subsequent_stage_completed", "residual evaluation")
	}

	existingEvaluation, err := aries.repo.GetByID(id, tx)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldEvaluation, err := aries.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}

	likelihood, err := aries.assetRiskRegisterSvc.GetLikelihoodByID(formDTO.RiskMatrixLikelihoodID, tx)
	if err != nil {
		return err
	}

	if riskRegister.ConfidentialityScore == nil {
		return fmt.Errorf("missing cofidentaility score in risk")
	}
	if riskRegister.IntegrityScore == nil {
		return fmt.Errorf("missing integrity score in risk")
	}
	if riskRegister.AvailabilityScore == nil {
		return fmt.Errorf("missing avaiability score in risk")
	}

	impactScore := max(*riskRegister.ConfidentialityScore, *riskRegister.IntegrityScore, *riskRegister.AvailabilityScore)

	overallImpact, err := aries.assetRiskRegisterSvc.GetImpactByScore(impactScore, tx)
	if err != nil {
		return err
	}

	riskValue := int(overallImpact.Score * likelihood.Score)

	riskRatingValue, err := aries.assetRiskRegisterSvc.GetRiskRatingValueByValue(riskValue, tx)
	if err != nil {
		return err
	}

	existingEvaluation.RiskMatrixLikelihoodID = likelihood.ID
	existingEvaluation.RiskMatrixImpactID = overallImpact.ID
	existingEvaluation.RiskMatrixRiskRatingID = riskRatingValue.RiskMatrixRiskRatingID
	existingEvaluation.RiskMatrixRiskRatingValueID = riskRatingValue.ID
	existingEvaluation.UpdatedByID = &currentUserID

	if err = aries.repo.Update(existingEvaluation, tx); err != nil {
		return err
	}

	now := time.Now()
	riskRegister.InherentLikelihoodID = &likelihood.ID
	riskRegister.InherentImpactID = &overallImpact.ID
	riskRegister.InherentRiskRatingID = &riskRatingValue.RiskMatrixRiskRatingID
	riskRegister.InherentRiskRatingValueID = &riskRatingValue.ID
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	if err = aries.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return err
	}

	go func() {
		newEvaluation, logErr := aries.GetFullDetailsByID(riskID)
		if logErr != nil {
			utility.LogError("Failed to fetch risk inherent evaluation", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		oldData := map[string]interface{}{
			"risk_matrix_likelihood":        oldEvaluation.RiskMatrixLikelihood.TitleEn,
			"risk_matrix_impact":            oldEvaluation.RiskMatrixImpact.TitleEn,
			"risk_matrix_risk_rating":       oldEvaluation.RiskMatrixRiskRating.TitleEn,
			"risk_matrix_risk_rating_value": oldEvaluation.RiskMatrixRiskRatingValue.Value,
		}

		newData := map[string]interface{}{
			"risk_matrix_likelihood":        newEvaluation.RiskMatrixLikelihood.TitleEn,
			"risk_matrix_impact":            newEvaluation.RiskMatrixImpact.TitleEn,
			"risk_matrix_risk_rating":       newEvaluation.RiskMatrixRiskRating.TitleEn,
			"risk_matrix_risk_rating_value": newEvaluation.RiskMatrixRiskRatingValue.Value,
		}

		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) == 0 {
			return
		}

		logErr = aries.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskEvaluation{}.TableName(),
			id,
			changes,
			currentUserID,
			"asset",
			"asset_risk",
			"risk",
			nil,
		)

		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.UpdateAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return nil
}

func (aries *AssetRiskInherentEvaluationService) Delete(riskID uint, id uint, currentUserID uint) (failedID uint, err error) {
	tx := aries.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	riskRegister, err := aries.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return 0, err
	}

	if riskRegister.ResidualRiskRatingID != nil {
		return 0, utility.NewValidationError(fmt.Errorf("cannot delete inherent evaluation after residual evaluation has been completed"), "subsequent_stage_completed", "residual evaluation")
	}

	if _, err = aries.repo.GetByID(id, tx); err != nil {
		return 0, utility.NewResourceNotFoundError(err, id)
	}

	if err = aries.repo.Delete(id, tx); err != nil {
		return id, err
	}

	riskRegister.InherentLikelihoodID = nil
	riskRegister.InherentImpactID = nil
	riskRegister.InherentRiskRatingID = nil
	riskRegister.InherentRiskRatingValueID = nil
	riskRegister.UpdatedBYID = &currentUserID
	now := time.Now()
	riskRegister.UpdatedAt = &now

	if err = aries.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return id, err
	}

	timelineType, err := aries.assetRiskRegisterSvc.GetTimelineTypeByType("risk_analysis", tx)
	if err != nil {
		return id, err
	}
	timelineStatus, err := aries.assetRiskRegisterSvc.GetTimelineStatusByType("pending", tx)
	if err != nil {
		return id, err
	}
	riskTimeline, err := aries.assetRiskRegisterSvc.GetRiskTimelineByRiskAndType(riskID, timelineType.ID, tx)
	if err != nil {
		return id, err
	}
	riskTimeline.TimelineStatusID = timelineStatus.ID
	riskTimeline.PerformedByID = nil
	riskTimeline.PerformedAt = nil

	if err = aries.assetRiskRegisterSvc.UpdateRiskTimeline(riskTimeline, tx); err != nil {
		return id, err
	}

	go func() {
		deletedEvaluation := map[string]interface{}{
			"id":               id,
			"risk_register_id": riskID,
		}

		logErr := aries.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskEvaluation{}.TableName(),
			id,
			deletedEvaluation,
			currentUserID,
			"asset",
			"asset_risk",
			"risk",
			nil,
		)

		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.DeleteAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return 0, nil
}

func toAssetRiskInherentEvaluationDTO(riskEvaluation *models.RiskEvaluation, riskRegister *models.RiskRegister) *assetDTOs.AssetRiskInherentEvaluationDTO {
	result := &assetDTOs.AssetRiskInherentEvaluationDTO{
		ID:                   riskEvaluation.ID,
		ConfidentialityScore: riskRegister.ConfidentialityScore,
		IntegrityScore:       riskRegister.IntegrityScore,
		AvailabilityScore:    riskRegister.AvailabilityScore,
	}

	if riskEvaluation.RiskMatrixLikelihood != nil {
		likelihood := riskEvaluation.RiskMatrixLikelihood
		result.RiskMatrixLikelihood = &assetDTOs.AssetRiskMatrixLikelihoodDTO{
			ID:            likelihood.ID,
			RiskMatrixID:  likelihood.RiskMatrixID,
			TitleEn:       likelihood.TitleEn,
			TitleAr:       likelihood.TitleAr,
			Score:         likelihood.Score,
			DescriptionEn: likelihood.DescriptionEn,
			DescriptionAr: likelihood.DescriptionAr,
			Probability:   likelihood.Probability,
			Status:        likelihood.Status,
		}
	}

	if riskEvaluation.RiskMatrixImpact != nil {
		impact := riskEvaluation.RiskMatrixImpact
		result.RiskMatrixImpact = &assetDTOs.AssetRiskMatrixImpactDTO{
			ID:            impact.ID,
			RiskMatrixID:  impact.RiskMatrixID,
			TitleEn:       impact.TitleEn,
			TitleAr:       impact.TitleAr,
			Score:         impact.Score,
			DescriptionEn: impact.DescriptionEn,
			DescriptionAr: impact.DescriptionAr,
			Status:        impact.Status,
		}
	}

	if riskEvaluation.RiskMatrixRiskRating != nil {
		rating := riskEvaluation.RiskMatrixRiskRating
		result.RiskMatrixRiskRating = &assetDTOs.AssetRiskMatrixRiskRatingDTO{
			ID:           rating.ID,
			RiskMatrixID: rating.RiskMatrixID,
			TitleEn:      rating.TitleEn,
			TitleAr:      rating.TitleAr,
			ColorCode:    rating.ColorCode,
			Type:         rating.Type,
			Status:       rating.Status,
		}
	}

	if riskEvaluation.RiskMatrixRiskRatingValue != nil {
		ratingValue := riskEvaluation.RiskMatrixRiskRatingValue
		result.RiskMatrixRiskRatingValue = &assetDTOs.AssetRiskMatrixRiskRatingValueDTO{
			ID:                     ratingValue.ID,
			RiskMatrixRiskRatingID: ratingValue.RiskMatrixRiskRatingID,
			Value:                  ratingValue.Value,
		}
	}

	return result
}
