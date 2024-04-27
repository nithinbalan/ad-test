package services

import (
	logServices "jethur/grc/logs/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
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
	riskMatrix, err := rmrrvs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return err
	}

	riskMatrix.CalculationType = riskMatrixRiskRatingValueFormDTO.CalculationType
	err = rmrrvs.riskMatrixSvc.Update(riskMatrix.ID, riskMatrix)
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
	if len(riskValues) > 0 {
		if err := rmrrvs.repo.Create(riskValues, riskMatrixRiskRatingIds); err != nil {
			return err
		}
	}

	return nil
}
