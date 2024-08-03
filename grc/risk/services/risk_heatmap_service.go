package services

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
)

type RiskHeatmapServiceInterface interface {
	GetRiskHeatmapByProcessID(processID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapByControlID(controlID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapByComplianceID(complianceID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapByResponsibleUnit() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error)
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

	return rhs.getRiskHeatmap(riskRegisterIDs)
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByControlID(controlID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error) {
	riskControls, err := rhs.repo.GetRiskControlsByControlID(controlID)
	if err != nil {
		return nil, err
	}

	var riskRegisterIDs []uint
	for _, rc := range riskControls {
		riskRegisterIDs = append(riskRegisterIDs, rc.RiskRegisterID)
	}

	return rhs.getRiskHeatmap(riskRegisterIDs)
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByComplianceID(complianceID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error) {
	riskCompliances, err := rhs.repo.GetRiskCompliancesByComplianceID(complianceID)
	if err != nil {
		return nil, err
	}

	var riskRegisterIDs []uint
	for _, rc := range riskCompliances {
		riskRegisterIDs = append(riskRegisterIDs, rc.RiskRegisterID)
	}

	return rhs.getRiskHeatmap(riskRegisterIDs)
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByResponsibleUnit() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error) {
	riskRegisters, err := rhs.repo.GetAllRiskRegisters()
	if err != nil {
		return nil, err
	}

	riskRegisterIDs := make([]uint, len(riskRegisters))
	for i, register := range riskRegisters {
		riskRegisterIDs[i] = register.ID
	}

	heatmap, err := rhs.getRiskHeatmap(riskRegisterIDs)
	if err != nil {
		return nil, err
	}

	responsibleUnitMap := make(map[uint]*riskDTOs.RiskHeatmapResponseByItemDTO)
	for _, register := range riskRegisters {
		if _, exists := responsibleUnitMap[register.ResponsibleUnitID]; !exists {
			responsibleUnitMap[register.ResponsibleUnitID] = &riskDTOs.RiskHeatmapResponseByItemDTO{
				ItemID:       register.ResponsibleUnitID,
				ItemTitleEn:  register.ResponsibleUnit.TitleEn,
				ItemTitleAr:  register.ResponsibleUnit.TitleAr,
				HeatmapCells: []*riskDTOs.RiskHeatmapResponseDTO{},
			}
		}
	}

	for unitID, unitResponse := range responsibleUnitMap {
		for _, likelihoodDetails := range heatmap {
			for _, detail := range likelihoodDetails {
				// Set the RiskRegisterCount to 0 initially
				detailCopy := *detail
				detailCopy.RiskRegisterCount = 0

				// Find risk registers for this unit with the same likelihood and impact
				for _, register := range riskRegisters {
					if register.ResponsibleUnitID == unitID &&
						register.CurrentLikelihoodID != nil &&
						register.CurrentImpactID != nil &&
						*register.CurrentLikelihoodID == detailCopy.RiskMatrixLikelihoodID &&
						*register.CurrentImpactID == detailCopy.RiskMatrixImpactID {
						detailCopy.RiskRegisterCount++
					}
				}

				unitResponse.HeatmapCells = append(unitResponse.HeatmapCells, &detailCopy)
			}
		}
	}

	heatmapResponses := make([]*riskDTOs.RiskHeatmapResponseByItemDTO, 0, len(responsibleUnitMap))
	for _, unitResponse := range responsibleUnitMap {
		heatmapResponses = append(heatmapResponses, unitResponse)
	}

	return heatmapResponses, nil
}

// function to get risk heatmap matrix
func (rhs *RiskHeatmapService) getRiskHeatmap(riskRegisterIDs []uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error) {
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
