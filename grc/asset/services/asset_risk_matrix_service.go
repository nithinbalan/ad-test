package services

import (
	"errors"

	"jethur/grc/asset/repositories"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskMatrixServiceInterface interface {
	GetAllImpacts() ([]*riskDTOs.RiskMatrixImpactDTO, error)
	GetAllLikelihoods() ([]*riskDTOs.RiskMatrixLikelihoodDTO, error)
	GetAllRiskRatings() ([]*riskDTOs.RiskMatrixRiskRatingDTO, error)
	GetAllRiskRatingValues() (*riskDTOs.RiskMatrixRiskRatingValueFormDTO, error)
	GetRiskMatrix() ([][]*riskDTOs.RiskMatrixRiskRatingResponseDTO, error)
}

type AssetRiskMatrixService struct {
	repo repositories.AssetRiskMatrixRepositoryInterface
}

func NewAssetRiskMatrixService(repo repositories.AssetRiskMatrixRepositoryInterface) AssetRiskMatrixServiceInterface {
	return &AssetRiskMatrixService{
		repo: repo,
	}
}

func (arms *AssetRiskMatrixService) GetAllImpacts() ([]*riskDTOs.RiskMatrixImpactDTO, error) {
	return arms.repo.GetAllImpacts()
}

func (arms *AssetRiskMatrixService) GetAllLikelihoods() ([]*riskDTOs.RiskMatrixLikelihoodDTO, error) {
	return arms.repo.GetAllLikelihoods()
}

func (arms *AssetRiskMatrixService) GetAllRiskRatings() ([]*riskDTOs.RiskMatrixRiskRatingDTO, error) {
	return arms.repo.GetAllRiskRatings()
}

func (arms *AssetRiskMatrixService) GetAllRiskRatingValues() (*riskDTOs.RiskMatrixRiskRatingValueFormDTO, error) {
	riskMatrix, err := arms.repo.GetMatrix()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, "isms")
		}
		return nil, err
	}

	riskRatingValues, err := arms.repo.GetAllRiskRatingValues()
	if err != nil {
		return nil, err
	}

	groupedValues := make(map[uint][]uint)
	for _, value := range riskRatingValues {
		groupedValues[value.RiskMatrixRiskRatingID] = append(groupedValues[value.RiskMatrixRiskRatingID], value.Value)
	}

	riskMatrixRiskRatingValueDTO := []*riskDTOs.RiskMatrixRiskRatingValueArrayDTO{}
	for id, values := range groupedValues {
		riskMatrixRiskRatingValueDTO = append(riskMatrixRiskRatingValueDTO, &riskDTOs.RiskMatrixRiskRatingValueArrayDTO{
			ID:                     id,
			RiskMatrixRiskRatingID: id,
			Values:                 values,
		})
	}

	return &riskDTOs.RiskMatrixRiskRatingValueFormDTO{
		CalculationType:      riskMatrix.CalculationType,
		RiskMatrixRiskRating: riskMatrixRiskRatingValueDTO,
	}, nil
}

func (arms *AssetRiskMatrixService) GetRiskMatrix() ([][]*riskDTOs.RiskMatrixRiskRatingResponseDTO, error) {
	riskMatrix, err := arms.repo.GetMatrixDetails()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, "isms")
		}
		return nil, err
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

	var riskMatrixResponse [][]*riskDTOs.RiskMatrixRiskRatingResponseDTO

	calculationMethod := riskMatrix.CalculationType

	for _, likelihood := range riskMatrix.RiskMatrixLikelihood {
		if likelihood.Status == "active" {
			var likelihoodDetails []*riskDTOs.RiskMatrixRiskRatingResponseDTO

			for _, impact := range riskMatrix.RiskMatrixImpact {
				if impact.Status == "active" {
					score := evaluateAssetRiskMatrixScore(impact, likelihood, calculationMethod)
					riskRating, _ := findAssetRiskMatrixRatingValue(riskRatingValueArray, score)

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

func evaluateAssetRiskMatrixScore(impact *models.RiskMatrixImpact, likelihood *models.RiskMatrixLikelihood, calculationMethod string) uint {
	var score uint
	if calculationMethod == "addition" {
		score = impact.Score + likelihood.Score
	} else if calculationMethod == "multiplication" {
		score = impact.Score * likelihood.Score
	}

	return score
}

func findAssetRiskMatrixRatingValue(riskRatingValues []*riskDTOs.RiskMatrixRiskRatingDataDTO, score uint) (*riskDTOs.RiskMatrixRiskRatingDataDTO, bool) {
	for _, riskRatingValue := range riskRatingValues {
		if riskRatingValue.Value == score {
			return riskRatingValue, true
		}
	}
	return nil, false
}
