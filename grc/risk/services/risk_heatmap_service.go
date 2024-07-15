package services

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
)

type RiskHeatmapServiceInterface interface {
	GetRiskHeatmapByProcessID(processID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapByOwnerUnit(params *riskDTOs.RiskRegisterListQueryParamsDTO) ([][]interface{}, error)
}

type RiskHeatmapService struct {
	repo                         repositories.RiskHeatmapRepositoryInterface
	riskMatrixSvc                RiskMatrixServiceInterface
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface
	riskRegisterSvc              RiskRegisterServiceInterface
}

func NewRiskHeatmapService(repo repositories.RiskHeatmapRepositoryInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	riskMatrixRiskRatingValueSvc RiskMatrixRiskRatingValueServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
) RiskHeatmapServiceInterface {
	return &RiskHeatmapService{
		repo:                         repo,
		riskMatrixSvc:                riskMatrixSvc,
		riskMatrixRiskRatingValueSvc: riskMatrixRiskRatingValueSvc,
		riskRegisterSvc:              riskRegisterSvc,
	}
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByProcessID(processID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error) {
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

	calculationMethod := riskMatrix.CalculationType
	var heatmap [][]*riskDTOs.RiskHeatmapResponseDTO

	for _, likelihood := range riskMatrix.RiskMatrixLikelihood {
		if likelihood.Status == "active" {
			var likelihoodDetails []*riskDTOs.RiskHeatmapResponseDTO
			for _, impact := range riskMatrix.RiskMatrixImpact {
				if impact.Status == "active" {
					score := rhs.riskMatrixRiskRatingValueSvc.EvaluateScore(impact, likelihood, calculationMethod)
					riskRating, _ := rhs.riskMatrixRiskRatingValueSvc.FindRiskRatingValue(riskRatingValueArray, score)

					riskRegisterCount, _ := rhs.repo.CountRiskRegistersWithLikelihoodAndImpact(riskRegisterIDs, riskRating.RiskMatrixRiskRatingValueID, likelihood.ID, impact.ID)

					impactDetail := &riskDTOs.RiskHeatmapResponseDTO{
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

func (rhs *RiskHeatmapService) GetRiskHeatmapByOwnerUnit(params *riskDTOs.RiskRegisterListQueryParamsDTO) ([][]interface{}, error) {
	// Fetch all risk registers using the existing method
	riskRegisterResponse, err := rhs.riskRegisterSvc.GetAll(params)
	if err != nil {
		return nil, err
	}

	riskRegisters := riskRegisterResponse.Items

	// Fetch risk matrix details
	riskMatrixType, err := rhs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return nil, err
	}

	riskMatrix, err := rhs.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
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

	// Create a map to store the risk counts by owner unit and risk rating value
	riskCountMap := make(map[uint]map[uint]int)
	ownerUnitTitleMap := make(map[uint]string)

	// Initialize the map and store owner unit titles
	for _, register := range riskRegisters {
		if _, exists := riskCountMap[register.OwnerUnitID]; !exists {
			riskCountMap[register.OwnerUnitID] = make(map[uint]int)
			ownerUnitTitleMap[register.OwnerUnitID] = register.OwnerUnitTitleEn // Assuming the English title is used
		}
	}

	// Iterate over the risk registers to count risks by owner unit and risk rating value
	for _, register := range riskRegisters {
		if register.CurrentRiskRatingValueID != nil {
			riskCountMap[register.OwnerUnitID][*register.CurrentRiskRatingValueID]++
		}
	}

	// Create the heatmap response
	var heatmap [][]interface{}
	for ownerUnitID, riskRatings := range riskCountMap {
		var row []interface{}
		// Add owner unit details
		row = append(row, map[string]interface{}{
			"owner_unit_id":    ownerUnitID,
			"owner_unit_title": ownerUnitTitleMap[ownerUnitID],
		})
		// Add risk rating details
		var riskDetails []map[string]interface{}
		for _, riskRating := range riskRatingValueArray {
			count := riskRatings[riskRating.RiskMatrixRiskRatingValueID]
			riskDetails = append(riskDetails, map[string]interface{}{
				"risk_matrix_score":       riskRating.Value,
				"risk_matrix_risk_rating": riskRating,
				"risk_register_count":     count,
			})
		}
		row = append(row, riskDetails)
		heatmap = append(heatmap, row)
	}

	return heatmap, nil
}
