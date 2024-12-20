package services

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
)

type RiskHeatmapServiceInterface interface {
	GetRiskHeatmapByProcessID(processID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapByControlID(controlID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapByComplianceID(complianceID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapBySLAID(slaID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)

	GetRiskHeatmapByResponsibleUnit() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error)
	GetRiskHeatmapByRiskCategory() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error)
	GetRiskHeatmapByRiskSource() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error)
	GetRiskHeatmapByRiskMatrix() ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetRiskHeatmapByResponseStrategy() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error)
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

func (rhs *RiskHeatmapService) GetRiskHeatmapBySLAID(slaID uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error) {
	riskSLAs, err := rhs.repo.GetRiskSLAsBySLAID(slaID)
	if err != nil {
		return nil, err
	}

	var riskRegisterIDs []uint
	for _, rs := range riskSLAs {
		riskRegisterIDs = append(riskRegisterIDs, rs.RiskRegisterID)
	}

	return rhs.getRiskHeatmap(riskRegisterIDs)
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByResponsibleUnit() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error) {
	riskRegisters, err := rhs.repo.GetAllRiskRegisters()
	if err != nil {
		return nil, err
	}

	riskMatrixType, err := rhs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return nil, err
	}

	riskMatrix, err := rhs.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		return nil, err
	}

	// Initialize map for responsible units
	responsibleUnitMap := make(map[uint]*riskDTOs.RiskHeatmapResponseByItemDTO)
	for _, register := range riskRegisters {
		if _, exists := responsibleUnitMap[register.ResponsibleUnitID]; !exists {
			responsibleUnitMap[register.ResponsibleUnitID] = &riskDTOs.RiskHeatmapResponseByItemDTO{
				ItemID:       register.ResponsibleUnitID,
				ItemTitleEn:  register.ResponsibleUnit.TitleEn,
				ItemTitleAr:  register.ResponsibleUnit.TitleAr,
				HeatmapCells: []*riskDTOs.RiskHeatmapResponseByItemDetailsDTO{},
			}

			// Initialize HeatmapCells with all possible values from the risk matrix
			for _, rating := range riskMatrix.RiskMatrixRiskRating {
				for _, value := range rating.RiskMatrixRiskRatingValue {
					responsibleUnitMap[register.ResponsibleUnitID].HeatmapCells = append(responsibleUnitMap[register.ResponsibleUnitID].HeatmapCells, &riskDTOs.RiskHeatmapResponseByItemDetailsDTO{
						RiskMatrixRiskRatingID:      rating.ID,
						RiskMatrixRiskRatingTitleEN: rating.TitleEn,
						RiskMatrixRiskRatingTitleAR: rating.TitleAr,
						ColorCode:                   rating.ColorCode,
						Value:                       value.Value,
						RiskRegisterCount:           0, // Initialize count to 0
					})
				}
			}
		}
	}

	// Populate each responsible unit's heatmap cells
	for _, register := range riskRegisters {
		if register.CurrentRiskRatingID != nil && register.CurrentRiskRatingValueID != nil {
			unitResponse := responsibleUnitMap[register.ResponsibleUnitID]
			ratingID := *register.CurrentRiskRatingID
			valueID := *register.CurrentRiskRatingValueID

			// Find and increment the corresponding cell
			for _, cell := range unitResponse.HeatmapCells {
				if cell.RiskMatrixRiskRatingID == ratingID && cell.Value == valueID {
					cell.RiskRegisterCount++
					break
				}
			}
		}
	}

	// Prepare the final response
	heatmapResponses := make([]*riskDTOs.RiskHeatmapResponseByItemDTO, 0, len(responsibleUnitMap))
	for _, unitResponse := range responsibleUnitMap {
		heatmapResponses = append(heatmapResponses, unitResponse)
	}

	return heatmapResponses, nil
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByRiskCategory() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error) {
	riskRegisters, err := rhs.repo.GetAllRiskRegisters()
	if err != nil {
		return nil, err
	}

	riskMatrixType, err := rhs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return nil, err
	}

	riskMatrix, err := rhs.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		return nil, err
	}

	// Initialize map for risk categories
	riskCategoryMap := make(map[uint]*riskDTOs.RiskHeatmapResponseByItemDTO)
	for _, register := range riskRegisters {
		if _, exists := riskCategoryMap[register.RiskCategoryID]; !exists {
			riskCategoryMap[register.RiskCategoryID] = &riskDTOs.RiskHeatmapResponseByItemDTO{
				ItemID:       register.RiskCategoryID,
				ItemTitleEn:  register.RiskCategory.TitleEn,
				ItemTitleAr:  register.RiskCategory.TitleAr,
				HeatmapCells: []*riskDTOs.RiskHeatmapResponseByItemDetailsDTO{},
			}

			// Initialize HeatmapCells with all possible values from the risk matrix
			for _, rating := range riskMatrix.RiskMatrixRiskRating {
				for _, value := range rating.RiskMatrixRiskRatingValue {
					riskCategoryMap[register.RiskCategoryID].HeatmapCells = append(riskCategoryMap[register.RiskCategoryID].HeatmapCells, &riskDTOs.RiskHeatmapResponseByItemDetailsDTO{
						RiskMatrixRiskRatingID:      rating.ID,
						RiskMatrixRiskRatingTitleEN: rating.TitleEn,
						RiskMatrixRiskRatingTitleAR: rating.TitleAr,
						ColorCode:                   rating.ColorCode,
						Value:                       value.Value,
						RiskRegisterCount:           0, // Initialize count to 0
					})
				}
			}
		}
	}

	// Populate each risk category's heatmap cells
	for _, register := range riskRegisters {
		if register.CurrentRiskRatingID != nil && register.CurrentRiskRatingValueID != nil {
			categoryResponse := riskCategoryMap[register.RiskCategoryID]
			ratingID := *register.CurrentRiskRatingID
			valueID := *register.CurrentRiskRatingValueID

			// Find and increment the corresponding cell
			for _, cell := range categoryResponse.HeatmapCells {
				if cell.RiskMatrixRiskRatingID == ratingID && cell.Value == valueID {
					cell.RiskRegisterCount++
					break
				}
			}
		}
	}

	// Prepare the final response
	heatmapResponses := make([]*riskDTOs.RiskHeatmapResponseByItemDTO, 0, len(riskCategoryMap))
	for _, categoryResponse := range riskCategoryMap {
		heatmapResponses = append(heatmapResponses, categoryResponse)
	}

	return heatmapResponses, nil
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByRiskSource() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error) {
	riskRegisters, err := rhs.repo.GetAllRiskRegisters()
	if err != nil {
		return nil, err
	}

	riskMatrixType, err := rhs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return nil, err
	}

	riskMatrix, err := rhs.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
	if err != nil {
		return nil, err
	}

	// Initialize map for risk sources
	riskSourceMap := make(map[uint]*riskDTOs.RiskHeatmapResponseByItemDTO)
	for _, register := range riskRegisters {
		if register.RiskSourceID != nil {
			if _, exists := riskSourceMap[*register.RiskSourceID]; !exists {
				riskSourceMap[*register.RiskSourceID] = &riskDTOs.RiskHeatmapResponseByItemDTO{
					ItemID:       *register.RiskSourceID,
					ItemTitleEn:  register.RiskSource.TitleEn,
					ItemTitleAr:  register.RiskSource.TitleAr,
					HeatmapCells: []*riskDTOs.RiskHeatmapResponseByItemDetailsDTO{},
				}

				// Initialize HeatmapCells with all possible values from the risk matrix
				for _, rating := range riskMatrix.RiskMatrixRiskRating {
					for _, value := range rating.RiskMatrixRiskRatingValue {
						riskSourceMap[*register.RiskSourceID].HeatmapCells = append(riskSourceMap[*register.RiskSourceID].HeatmapCells, &riskDTOs.RiskHeatmapResponseByItemDetailsDTO{
							RiskMatrixRiskRatingID:      rating.ID,
							RiskMatrixRiskRatingTitleEN: rating.TitleEn,
							RiskMatrixRiskRatingTitleAR: rating.TitleAr,
							ColorCode:                   rating.ColorCode,
							Value:                       value.Value,
							RiskRegisterCount:           0, // Initialize count to 0
						})
					}
				}
			}
		}
	}

	// Populate each risk source's heatmap cells
	for _, register := range riskRegisters {
		if register.RiskSourceID != nil && register.CurrentRiskRatingID != nil && register.CurrentRiskRatingValueID != nil {
			sourceResponse := riskSourceMap[*register.RiskSourceID]
			ratingID := *register.CurrentRiskRatingID
			valueID := *register.CurrentRiskRatingValueID

			// Find and increment the corresponding cell
			for _, cell := range sourceResponse.HeatmapCells {
				if cell.RiskMatrixRiskRatingID == ratingID && cell.Value == valueID {
					cell.RiskRegisterCount++
					break
				}
			}
		}
	}

	// Prepare the final response
	heatmapResponses := make([]*riskDTOs.RiskHeatmapResponseByItemDTO, 0, len(riskSourceMap))
	for _, sourceResponse := range riskSourceMap {
		heatmapResponses = append(heatmapResponses, sourceResponse)
	}

	return heatmapResponses, nil
}

func (rhs *RiskHeatmapService) GetRiskHeatmapByRiskMatrix() ([][]*riskDTOs.RiskHeatmapResponseDTO, error) {
	risks, err := rhs.repo.GetAllRiskRegisters()
	if err != nil {
		return nil, err
	}

	var riskRegisterIDs []uint
	for _, rp := range risks {
		riskRegisterIDs = append(riskRegisterIDs, rp.ID)
	}

	return rhs.getRiskHeatmap(riskRegisterIDs)
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

func (rhs *RiskHeatmapService) GetRiskHeatmapByResponseStrategy() ([]*riskDTOs.RiskHeatmapResponseByItemDTO, error) {
    riskRegisters, err := rhs.repo.GetAllRiskRegisters()
    if err != nil {
        return nil, err
    }

    riskMatrixType, err := rhs.riskMatrixSvc.GetByType("erm")
    if err != nil {
        return nil, err
    }

    riskMatrix, err := rhs.riskMatrixSvc.GetDetailsByID(riskMatrixType.ID)
    if err != nil {
        return nil, err
    }

    // Initialize map for Risk Response Strategies
    responseStrategyMap := make(map[uint]*riskDTOs.RiskHeatmapResponseByItemDTO)
    for _, register := range riskRegisters {
        if register.RiskResponseStrategyID != nil {
            if _, exists := responseStrategyMap[*register.RiskResponseStrategyID]; !exists {
                responseStrategyMap[*register.RiskResponseStrategyID] = &riskDTOs.RiskHeatmapResponseByItemDTO{
                    ItemID:       *register.RiskResponseStrategyID,
                    ItemTitleEn:  register.RiskResponseStrategy.TitleEn,
                    ItemTitleAr:  register.RiskResponseStrategy.TitleAr,
                    HeatmapCells: []*riskDTOs.RiskHeatmapResponseByItemDetailsDTO{},
                }

                // Initialize HeatmapCells with all possible values from the risk matrix
                for _, rating := range riskMatrix.RiskMatrixRiskRating {
                    for _, value := range rating.RiskMatrixRiskRatingValue {
                        responseStrategyMap[*register.RiskResponseStrategyID].HeatmapCells = append(responseStrategyMap[*register.RiskResponseStrategyID].HeatmapCells, &riskDTOs.RiskHeatmapResponseByItemDetailsDTO{
                            RiskMatrixRiskRatingID:      rating.ID,
                            RiskMatrixRiskRatingTitleEN: rating.TitleEn,
                            RiskMatrixRiskRatingTitleAR: rating.TitleAr,
                            ColorCode:                   rating.ColorCode,
                            Value:                       value.Value,
                            RiskRegisterCount:           0, // Initialize count to 0
                        })
                    }
                }
            }
        }
    }

    // Populate each strategy's heatmap cells
    for _, register := range riskRegisters {
        if register.RiskResponseStrategyID != nil && register.CurrentRiskRatingID != nil && register.CurrentRiskRatingValueID != nil {
            strategyResponse := responseStrategyMap[*register.RiskResponseStrategyID]
            ratingID := *register.CurrentRiskRatingID
            valueID := *register.CurrentRiskRatingValueID

            // Increment the corresponding cell count
            for _, cell := range strategyResponse.HeatmapCells {
                if cell.RiskMatrixRiskRatingID == ratingID && cell.Value == valueID {
                    cell.RiskRegisterCount++
                    break
                }
            }
        }
    }

    // Prepare the final response
    heatmapResponses := make([]*riskDTOs.RiskHeatmapResponseByItemDTO, 0, len(responseStrategyMap))
    for _, strategyResponse := range responseStrategyMap {
        heatmapResponses = append(heatmapResponses, strategyResponse)
    }

    return heatmapResponses, nil
}
