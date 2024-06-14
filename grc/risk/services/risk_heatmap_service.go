package services

import (
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
)

type RiskHeatmapServiceInterface interface {
	GetRiskHeatmapByProcessID(processID uint) ([][]*dtos.RiskHeatmapResponseDTO, error)
}

type RiskHeatmapService struct {
	repo                         repositories.RiskHeatmapRepositoryInterface
	riskMatrixSvc                RiskMatrixServiceInterface
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface
}

func NewRiskHeatmapService(repo repositories.RiskHeatmapRepositoryInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface,
) RiskHeatmapServiceInterface {
	return &RiskHeatmapService{
		repo:                         repo,
		riskMatrixSvc:                riskMatrixSvc,
		riskMatrixRiskRatingValueSvc: riskMatrixRiskRatingValueSvc,
	}
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByProcessID(processID uint) ([][]*dtos.RiskHeatmapResponseDTO, error) {
	riskProcesses, err := rhs.repo.GetRiskProcessesByProcessID(processID)
	if err != nil {
		return nil, err
	}

	var riskRegisterIDs []uint
	for _, rp := range riskProcesses {
		riskRegisterIDs = append(riskRegisterIDs, rp.RiskRegisterID)
	}

	riskMatrixType, err := rhs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return nil, err
	}

	riskMatrix, err := rhs.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		return nil, err
	}

	var riskRatingValueArray []*dtos.RiskMatrixRiskRatingDataDTO
	for _, riskRating := range riskMatrix.RiskMatrixRiskRating {
		for _, value := range riskRating.RiskMatrixRiskRatingValue {
			detail := &dtos.RiskMatrixRiskRatingDataDTO{
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

	calculationMethod := riskMatrix.CalculationType
	var heatmap [][]*dtos.RiskHeatmapResponseDTO

	for _, likelihood := range riskMatrix.RiskMatrixLikelihood {
		if likelihood.Status == "active" {
			var likelihoodDetails []*dtos.RiskHeatmapResponseDTO
			for _, impact := range riskMatrix.RiskMatrixImpact {
				if impact.Status == "active" {
					score := rhs.riskMatrixRiskRatingValueSvc.EvaluateScore(impact, likelihood, calculationMethod)
					riskRating, _ := rhs.riskMatrixRiskRatingValueSvc.FindRiskRatingValue(riskRatingValueArray, score)

					riskRegisterCount, _ := rhs.repo.CountRiskRegistersWithLikelihoodAndImpact(riskRegisterIDs, riskRating.RiskMatrixRiskRatingValueID, likelihood.ID, impact.ID)

					impactDetail := &dtos.RiskHeatmapResponseDTO{
						RiskMatrixImpactID:          impact.ID,
						RiskMatrixImpactTitleEN:     impact.TitleEn,
						RiskMatrixImpactTitleAR:     impact.TitleAr,
						RiskMatrixLikelihoodID:      likelihood.ID,
						RiskMatrixLikelihoodTitleEN: likelihood.TitleEn,
						RiskMatrixLikelihoodTitleAR: likelihood.TitleAr,
						RiskMatrixScore:             score,
						RiskMatrixRiskRating:        *riskRating,
						RiskRegisterCount:           riskRegisterCount,
					}

					likelihoodDetails = append(likelihoodDetails, impactDetail)
				}
			}
			heatmap = append(heatmap, likelihoodDetails)
		}
	}

	return heatmap, nil
}
