package services

import (
	"errors"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskMatrixRiskRatingValueServiceInterface interface {
	Create(riskMatrixRiskRatingValueFormDTO *riskDTOs.RiskMatrixRiskRatingValueFormDTO) error
	GetRiskMatrix() ([][]*riskDTOs.RiskMatrixRiskRatingResponseDTO, error)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, "erm")
		}
		return err
	}

	var riskMatrixRiskRatings []*riskDTOs.RiskMatrixRiskRatingDTO
	riskMatrixRiskRatings, err = rmrrvs.riskMatrixRiskRatingSvc.GetAllByRiskMatrix(riskMatrix.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, riskMatrix.ID)
		}
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
func (rmrrvs *RiskMatrixRiskRatingValueService) GetRiskMatrix() ([][]*riskDTOs.RiskMatrixRiskRatingResponseDTO, error) {
	riskMatrixType, err := rmrrvs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, "erm")
		}
		return nil, err
	}
	riskMatrix, err := rmrrvs.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskMatrixType.ID)
		}
		return nil, err
	}

	var riskRatingValueArray []*riskDTOs.RiskMatrixRiskRatingDataDTO

	for _, riskRating := range riskMatrix.RiskMatrixRiskRating {
		for _, value := range riskRating.RiskMatrixRiskRatingValue {
			detail := &riskDTOs.RiskMatrixRiskRatingDataDTO{
				RiskMatrixID:           riskRating.RiskMatrixID,
				RiskMatrixRiskRatingID: value.RiskMatrixRiskRatingID,
				RiskRatingTitleEN:      riskRating.TitleEn,
				RiskRatingTitleAR:      riskRating.TitleAr,
				RiskRatingTitleType:    riskRating.Type,
				ColorCode:              riskRating.ColorCode,
				Value:                  value.Value,
			}
			riskRatingValueArray = append(riskRatingValueArray, detail)
		}
	}

	var riskMatrixResponse [][]*riskDTOs.RiskMatrixRiskRatingResponseDTO

	calculationMethod := riskMatrix.CalculationType

	for _, likelihood := range riskMatrix.RiskMatrixLikelihood {

		if likelihood.Status == "active" {

			var likelihoodDetails []*riskDTOs.RiskMatrixRiskRatingResponseDTO

			// Iterate over each impact
			for _, impact := range riskMatrix.RiskMatrixImpact {
				if impact.Status == "active" {

					score := evaluateScore(impact, likelihood, calculationMethod)
					riskRating, _ := findRiskRatingValue(riskRatingValueArray, score)

					impactDetail := &riskDTOs.RiskMatrixRiskRatingResponseDTO{
						RiskMatrixImpactID:          impact.ID,
						RiskMatrixImpactTitleEN:     impact.TitleEn,
						RiskMatrixImpactTitleAR:     impact.TitleAr,
						RiskMatrixLikelihoodID:      likelihood.ID,
						RiskMatrixLikelihoodTitleEN: likelihood.TitleEn,
						RiskMatrixLikelihoodTitleAR: likelihood.TitleAr,
						RiskMatrixScore:             score,
						RiskMatrixRiskRating:        *riskRating,
					}

					likelihoodDetails = append(likelihoodDetails, impactDetail)
				}
			}

			riskMatrixResponse = append(riskMatrixResponse, likelihoodDetails)
		}
	}

	return riskMatrixResponse, nil
}

func evaluateScore(impact *models.RiskMatrixImpact, likelihood *models.RiskMatrixLikelihood, calculationMethod string) uint {
	var score uint
	if calculationMethod == "addition" {
		score = impact.Score + likelihood.Score
	} else if calculationMethod == "multiplication" {
		score = impact.Score * likelihood.Score
	}

	return score
}

func findRiskRatingValue(riskRatingValues []*riskDTOs.RiskMatrixRiskRatingDataDTO, score uint) (*riskDTOs.RiskMatrixRiskRatingDataDTO, bool) {
	for _, riskRatingValues := range riskRatingValues {
		if riskRatingValues.Value == score {
			return riskRatingValues, true
		}
	}
	return nil, false
}
