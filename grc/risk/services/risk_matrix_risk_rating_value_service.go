package services

import (
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
)

type RiskMatrixRiskRatingValueServiceInterface interface {
	Create(riskMatrixRiskRatingValueFormDTO *riskDTOs.RiskMatrixRiskRatingValueFormDTO) error
}

// Service and implementations
type RiskMatrixRiskRatingValueService struct {
	repo                    riskRepositories.RiskMatrixRiskRatingValueRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	riskMatrixSvc           RiskMatrixServiceInterface
	riskMatrixRiskRatingSvc RiskMatrixRiskRatingServiceInterface
}

func NewRiskMatrixRiskRatingValueService(
	repo riskRepositories.RiskMatrixRiskRatingValueRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	riskMatrixRiskRatingSvc RiskMatrixRiskRatingServiceInterface,
) RiskMatrixRiskRatingValueServiceInterface {
	return &RiskMatrixRiskRatingValueService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		riskMatrixSvc:           riskMatrixSvc,
		riskMatrixRiskRatingSvc: riskMatrixRiskRatingSvc,
	}
}

func (rmrrvs *RiskMatrixRiskRatingValueService) Create(riskMatrixRiskRatingValueFormDTO *riskDTOs.RiskMatrixRiskRatingValueFormDTO) error {
	tx := rmrrvs.repo.BeginTransaction()
	riskMatrix, err := rmrrvs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return err
	}

	var riskMatrixRiskRatings []*riskDTOs.RiskMatrixRiskRatingDTO
	riskMatrixRiskRatings, err = rmrrvs.riskMatrixRiskRatingSvc.GetAllByRiskMatrix(riskMatrix.ID)
	if err != nil {
		return err
	}

	var riskMatrixRiskRatingIds []uint
	for _, riskMatrixRiskRating := range riskMatrixRiskRatings {
		riskMatrixRiskRatingIds = append(riskMatrixRiskRatingIds, riskMatrixRiskRating.ID)
	}

	var riskValues []*models.RiskMatrixRiskRatingValue
	for _, riskRatingValues := range riskMatrixRiskRatingValueFormDTO.RiskMatrixRiskRating {

		_, err = rmrrvs.riskMatrixRiskRatingSvc.GetByID(riskRatingValues.RiskMatrixRiskRatingID)
		if err != nil {
			return err
		}

		for _, item := range riskRatingValues.Values {
			// Create a new RiskMatrixRiskRatingValue instance
			riskValue := &models.RiskMatrixRiskRatingValue{
				RiskMatrixRiskRatingID: riskRatingValues.RiskMatrixRiskRatingID,
				Value:                  item,
			}
			// Append the document to the slice
			riskValues = append(riskValues, riskValue)
		}
	}

	riskMatrix.CalculationType = riskMatrixRiskRatingValueFormDTO.CalculationType
	err = rmrrvs.riskMatrixSvc.Update(riskMatrix.ID, riskMatrix, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(riskValues) > 0 {
		if err := rmrrvs.repo.Create(riskValues, riskMatrixRiskRatingIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
