package services

import (
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
)

type RiskHeatmapServiceInterface interface {
	GetRiskHeatmapByProcessID(processID uint) ([][]*dtos.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapByOwnerUnit() ([][]*dtos.RiskHeatmapResponseDTO, error)
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

func (rhs *RiskHeatmapService) GetRiskHeatmapByOwnerUnit() ([][]*dtos.RiskHeatmapResponseDTO, error) {
	// Fetch all risk registers
	riskRegisters, err := rhs.repo.GetAllRiskRegisters()
	if err != nil {
		return nil, err
	}

	// Fetch risk matrix details
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

	// Create a map to store the risk counts by owner unit and risk rating value
	riskCountMap := make(map[uint]map[uint]int)

	// Initialize the map
	for _, register := range riskRegisters {
		if _, exists := riskCountMap[register.OwnerUnitID]; !exists {
			riskCountMap[register.OwnerUnitID] = make(map[uint]int)
		}
	}

	// Iterate over the risk registers to count risks by owner unit and risk rating value
	for _, register := range riskRegisters {
		if register.CurrentRiskRatingValueID != nil {
			riskCountMap[register.OwnerUnitID][*register.CurrentRiskRatingValueID]++
		}
	}

	// Create the heatmap response
	var heatmap [][]*dtos.RiskHeatmapResponseDTO
	for ownerUnitID, riskRatings := range riskCountMap {
		var row []*dtos.RiskHeatmapResponseDTO
		for _, riskRating := range riskRatingValueArray {
			count := riskRatings[riskRating.RiskMatrixRiskRatingValueID]
			row = append(row, &dtos.RiskHeatmapResponseDTO{
				OwnerUnitID:                 ownerUnitID,
				RiskMatrixRiskRatingValueID: riskRating.RiskMatrixRiskRatingValueID,
				RiskRegisterCount:           count,
				ColorCode:                   riskRating.ColorCode,
			})
		}
		heatmap = append(heatmap, row)
	}

	return heatmap, nil
}
