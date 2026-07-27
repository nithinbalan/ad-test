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

type AssetRiskResidualEvaluationServiceInterface interface {
	GetFullDetailsByID(riskID uint) (*assetDTOs.AssetRiskResidualEvaluationDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskEvaluation, error)
	Create(riskID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskResidualEvaluationFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskResidualEvaluationFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
}

type AssetRiskResidualEvaluationService struct {
	repo                     repositories.AssetRiskResidualEvaluationRepositoryInterface
	assetRiskRegisterSvc     AssetRiskRegisterServiceInterface
	assetRiskInherentEvalSvc AssetRiskInherentEvaluationServiceInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
}

func NewAssetRiskResidualEvaluationService(
	repo repositories.AssetRiskResidualEvaluationRepositoryInterface,
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface,
	assetRiskInherentEvalSvc AssetRiskInherentEvaluationServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) AssetRiskResidualEvaluationServiceInterface {
	return &AssetRiskResidualEvaluationService{
		repo:                     repo,
		assetRiskRegisterSvc:     assetRiskRegisterSvc,
		assetRiskInherentEvalSvc: assetRiskInherentEvalSvc,
		activityLogSvc:           activityLogSvc,
	}
}

func (ares *AssetRiskResidualEvaluationService) GetFullDetailsByID(riskID uint) (*assetDTOs.AssetRiskResidualEvaluationDTO, error) {
	riskRegister, err := ares.assetRiskRegisterSvc.GetByID(riskID)
	if err != nil {
		return nil, err
	}

	riskEvaluation, err := ares.repo.GetFullDetailsByRiskID(riskID)
	if err != nil {
		return nil, err
	}
	if riskEvaluation == nil {
		return nil, nil
	}

	return toAssetRiskResidualEvaluationDTO(riskEvaluation, riskRegister), nil
}

func (ares *AssetRiskResidualEvaluationService) GetByID(riskID uint, id uint) (*models.RiskEvaluation, error) {
	if _, err := ares.assetRiskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskEvaluation, err := ares.repo.GetByID(id)
	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, id)
	}

	return riskEvaluation, nil
}

func (ares *AssetRiskResidualEvaluationService) Create(riskID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskResidualEvaluationFormDTO) (id uint, err error) {
	tx := ares.repo.BeginTransaction()

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

	riskRegister, err := ares.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return 0, err
	}

	inherentEvaluation, err := ares.assetRiskInherentEvalSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return 0, err
	}
	if inherentEvaluation == nil {
		return 0, utility.NewValidationError(fmt.Errorf("inherent evaluation must be completed before creating residual evaluation"), "not_exist", "inherent evaluation")
	}

	likelihood, err := ares.assetRiskRegisterSvc.GetLikelihoodByID(formDTO.RiskMatrixLikelihoodID, tx)
	if err != nil {
		return 0, err
	}

	impactScore := max(formDTO.ConfidentialityScore, formDTO.IntegrityScore, formDTO.AvailabilityScore)

	overallImpact, err := ares.assetRiskRegisterSvc.GetImpactByScore(impactScore, tx)
	if err != nil {
		return 0, err
	}

	riskValue := int(overallImpact.Score * likelihood.Score)

	riskRatingValue, err := ares.assetRiskRegisterSvc.GetRiskRatingValueByValue(riskValue, tx)
	if err != nil {
		return 0, err
	}

	riskEvaluation := &models.RiskEvaluation{
		RiskRegisterID:              riskID,
		Type:                        "residual",
		RiskMatrixLikelihoodID:      likelihood.ID,
		RiskMatrixImpactID:          overallImpact.ID,
		RiskMatrixRiskRatingID:      riskRatingValue.RiskMatrixRiskRatingID,
		RiskMatrixRiskRatingValueID: riskRatingValue.ID,
		CreatedByID:                 currentUserID,
	}

	id, err = ares.repo.Create(riskEvaluation, tx)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	riskRegister.ResidualLikelihoodID = &likelihood.ID
	riskRegister.ResidualImpactID = &overallImpact.ID
	riskRegister.ResidualRiskRatingID = &riskRatingValue.RiskMatrixRiskRatingID
	riskRegister.ResidualRiskRatingValueID = &riskRatingValue.ID
	riskRegister.ConfidentialityScore = &formDTO.ConfidentialityScore
	riskRegister.IntegrityScore = &formDTO.IntegrityScore
	riskRegister.AvailabilityScore = &formDTO.AvailabilityScore
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	if err = ares.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return 0, err
	}

	timelineType, err := ares.assetRiskRegisterSvc.GetTimelineTypeByType("residual_risk_analysis", tx)
	if err != nil {
		return 0, err
	}
	timelineStatus, err := ares.assetRiskRegisterSvc.GetTimelineStatusByType("completed", tx)
	if err != nil {
		return 0, err
	}
	riskTimeline, err := ares.assetRiskRegisterSvc.GetRiskTimelineByRiskAndType(riskID, timelineType.ID, tx)
	if err != nil {
		return 0, err
	}
	riskTimeline.TimelineStatusID = timelineStatus.ID
	riskTimeline.PerformedByID = &currentUserID
	riskTimeline.PerformedAt = &now

	if err = ares.assetRiskRegisterSvc.UpdateRiskTimeline(riskTimeline, tx); err != nil {
		return 0, err
	}

	go func() {
		existingEvaluation, logErr := ares.GetFullDetailsByID(riskID)
		if logErr != nil {
			utility.LogError("Failed to fetch risk residual evaluation", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = ares.activityLogSvc.LogActivity(
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

func (ares *AssetRiskResidualEvaluationService) Update(riskID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskResidualEvaluationFormDTO) (err error) {
	tx := ares.repo.BeginTransaction()

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

	riskRegister, err := ares.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return err
	}

	existingEvaluation, err := ares.repo.GetByID(id, tx)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	inherentEvaluation, err := ares.assetRiskInherentEvalSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}
	if inherentEvaluation == nil {
		return utility.NewValidationError(fmt.Errorf("inherent evaluation must be completed before updating residual evaluation"), "not_exist", "inherent evaluation")
	}

	oldEvaluation, err := ares.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}

	likelihood, err := ares.assetRiskRegisterSvc.GetLikelihoodByID(formDTO.RiskMatrixLikelihoodID, tx)
	if err != nil {
		return err
	}

	impactScore := max(formDTO.ConfidentialityScore, formDTO.IntegrityScore, formDTO.AvailabilityScore)

	overallImpact, err := ares.assetRiskRegisterSvc.GetImpactByScore(impactScore, tx)
	if err != nil {
		return err
	}

	riskValue := int(overallImpact.Score * likelihood.Score)

	riskRatingValue, err := ares.assetRiskRegisterSvc.GetRiskRatingValueByValue(riskValue, tx)
	if err != nil {
		return err
	}

	existingEvaluation.RiskMatrixLikelihoodID = likelihood.ID
	existingEvaluation.RiskMatrixImpactID = overallImpact.ID
	existingEvaluation.RiskMatrixRiskRatingID = riskRatingValue.RiskMatrixRiskRatingID
	existingEvaluation.RiskMatrixRiskRatingValueID = riskRatingValue.ID
	existingEvaluation.UpdatedByID = &currentUserID

	if err = ares.repo.Update(existingEvaluation, tx); err != nil {
		return err
	}

	now := time.Now()
	riskRegister.ResidualLikelihoodID = &likelihood.ID
	riskRegister.ResidualImpactID = &overallImpact.ID
	riskRegister.ResidualRiskRatingID = &riskRatingValue.RiskMatrixRiskRatingID
	riskRegister.ResidualRiskRatingValueID = &riskRatingValue.ID
	riskRegister.ConfidentialityScore = &formDTO.ConfidentialityScore
	riskRegister.IntegrityScore = &formDTO.IntegrityScore
	riskRegister.AvailabilityScore = &formDTO.AvailabilityScore
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	if err = ares.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return err
	}

	go func() {
		newEvaluation, logErr := ares.GetFullDetailsByID(riskID)
		if logErr != nil {
			utility.LogError("Failed to fetch risk residual evaluation", map[string]interface{}{
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

		logErr = ares.activityLogSvc.LogActivity(
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

func (ares *AssetRiskResidualEvaluationService) Delete(riskID uint, id uint, currentUserID uint) (failedID uint, err error) {
	tx := ares.repo.BeginTransaction()

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

	riskRegister, err := ares.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return 0, err
	}

	if _, err = ares.repo.GetByID(id, tx); err != nil {
		return 0, utility.NewResourceNotFoundError(err, id)
	}

	if err = ares.repo.Delete(id, tx); err != nil {
		return id, err
	}

	riskRegister.ResidualLikelihoodID = nil
	riskRegister.ResidualImpactID = nil
	riskRegister.ResidualRiskRatingID = nil
	riskRegister.ResidualRiskRatingValueID = nil
	riskRegister.UpdatedBYID = &currentUserID
	now := time.Now()
	riskRegister.UpdatedAt = &now

	if err = ares.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return id, err
	}

	timelineType, err := ares.assetRiskRegisterSvc.GetTimelineTypeByType("residual_risk_analysis", tx)
	if err != nil {
		return id, err
	}
	timelineStatus, err := ares.assetRiskRegisterSvc.GetTimelineStatusByType("pending", tx)
	if err != nil {
		return id, err
	}
	riskTimeline, err := ares.assetRiskRegisterSvc.GetRiskTimelineByRiskAndType(riskID, timelineType.ID, tx)
	if err != nil {
		return id, err
	}
	riskTimeline.TimelineStatusID = timelineStatus.ID
	riskTimeline.PerformedByID = nil
	riskTimeline.PerformedAt = nil

	if err = ares.assetRiskRegisterSvc.UpdateRiskTimeline(riskTimeline, tx); err != nil {
		return id, err
	}

	go func() {
		deletedEvaluation := map[string]interface{}{
			"id":               id,
			"risk_register_id": riskID,
		}

		logErr := ares.activityLogSvc.LogActivity(
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

func toAssetRiskResidualEvaluationDTO(riskEvaluation *models.RiskEvaluation, riskRegister *models.RiskRegister) *assetDTOs.AssetRiskResidualEvaluationDTO {
	result := &assetDTOs.AssetRiskResidualEvaluationDTO{
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
