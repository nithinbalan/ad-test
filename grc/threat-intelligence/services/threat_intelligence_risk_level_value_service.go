package services

import (
	"jethur/grc/models"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/grc/threat-intelligence/repositories"
)

type ThreatIntelligenceRiskLevelValueServiceInterface interface {
	GetRiskAll() (*tiDTOs.ThreatIntelligenceRiskLevelValueFormDTO, error)
	Create(riskLevelValueFormDTO *tiDTOs.ThreatIntelligenceRiskLevelValueFormDTO) error
	GetRiskMatrix() ([][]*tiDTOs.ThreatIntelligenceRiskLevelMatrixResponseDTO, error)
	EvaluateScore(likelihood *models.ThreatIntelligenceLikelihood, severity *models.ThreatIntelligenceSeverity, calculationMethod string) uint
	FindRiskLevelValue(riskLevelValues []*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, score uint) (*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, bool)
	FindClosestRiskLevelValue(riskLevelValues []*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, score uint) (*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, bool)
}

// Service and implementations
type ThreatIntelligenceRiskLevelValueService struct {
	repo          repositories.ThreatIntelligenceRiskLevelValueRepositoryInterface
	likelihoodSvc ThreatIntelligenceLikelihoodServiceInterface
	severitySvc   ThreatIntelligenceSeverityServiceInterface
	riskLevelSvc  ThreatIntelligenceRiskLevelServiceInterface
}

func NewThreatIntelligenceRiskLevelValueService(
	repo repositories.ThreatIntelligenceRiskLevelValueRepositoryInterface,
	likelihoodSvc ThreatIntelligenceLikelihoodServiceInterface,
	severitySvc ThreatIntelligenceSeverityServiceInterface,
	riskLevelSvc ThreatIntelligenceRiskLevelServiceInterface,
) ThreatIntelligenceRiskLevelValueServiceInterface {
	return &ThreatIntelligenceRiskLevelValueService{
		repo:          repo,
		likelihoodSvc: likelihoodSvc,
		severitySvc:   severitySvc,
		riskLevelSvc:  riskLevelSvc,
	}
}

func (tirlvs *ThreatIntelligenceRiskLevelValueService) GetRiskAll() (*tiDTOs.ThreatIntelligenceRiskLevelValueFormDTO, error) {
	values, err := tirlvs.repo.GetAll()
	if err != nil {
		return nil, err
	}

	groupedValues := make(map[uint][]uint)
	for _, value := range values {
		groupedValues[value.ThreatIntelligenceRiskLevelID] = append(groupedValues[value.ThreatIntelligenceRiskLevelID], value.Value)
	}

	var riskLevelValueDTOs []*tiDTOs.ThreatIntelligenceRiskLevelValueArrayDTO
	for id, vals := range groupedValues {
		riskLevelValueDTOs = append(riskLevelValueDTOs, &tiDTOs.ThreatIntelligenceRiskLevelValueArrayDTO{
			ID:                            id,
			ThreatIntelligenceRiskLevelID: id,
			Values:                        vals,
		})
	}

	result := &tiDTOs.ThreatIntelligenceRiskLevelValueFormDTO{
		RiskLevels: riskLevelValueDTOs,
	}

	return result, nil
}

func (tirlvs *ThreatIntelligenceRiskLevelValueService) Create(riskLevelValueFormDTO *tiDTOs.ThreatIntelligenceRiskLevelValueFormDTO) error {
	tx := tirlvs.repo.BeginTransaction()

	var riskLevelIDs []uint
	var riskLevelValueModels []*models.ThreatIntelligenceRiskLevelValue

	for _, riskLevelValue := range riskLevelValueFormDTO.RiskLevels {
		_, err := tirlvs.riskLevelSvc.GetByID(riskLevelValue.ThreatIntelligenceRiskLevelID)
		if err != nil {
			tx.Rollback()
			return err
		}

		riskLevelIDs = append(riskLevelIDs, riskLevelValue.ThreatIntelligenceRiskLevelID)

		for _, item := range riskLevelValue.Values {
			riskLevelValueModels = append(riskLevelValueModels, &models.ThreatIntelligenceRiskLevelValue{
				ThreatIntelligenceRiskLevelID: riskLevelValue.ThreatIntelligenceRiskLevelID,
				Value:                         item,
			})
		}
	}

	if len(riskLevelValueModels) > 0 {
		if err := tirlvs.repo.Create(riskLevelValueModels, riskLevelIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (tirlvs *ThreatIntelligenceRiskLevelValueService) GetRiskMatrix() ([][]*tiDTOs.ThreatIntelligenceRiskLevelMatrixResponseDTO, error) {
	likelihoods, err := tirlvs.likelihoodSvc.GetAll("asc", "score", "", "active")
	if err != nil {
		return nil, err
	}

	severities, err := tirlvs.severitySvc.GetAll("asc", "score", "", "active")
	if err != nil {
		return nil, err
	}

	riskLevels, err := tirlvs.riskLevelSvc.GetAll(&tiDTOs.ThreatIntelligenceRiskLevelListQueryParamsDTO{
		Status: "active",
	})
	if err != nil {
		return nil, err
	}

	var riskLevelValueArray []*tiDTOs.ThreatIntelligenceRiskLevelDataDTO
	for _, riskLevel := range riskLevels {
		riskLevelDetails, err := tirlvs.riskLevelSvc.GetByID(riskLevel.ID)
		if err != nil {
			continue
		}

		for _, value := range riskLevelDetails.ThreatIntelligenceRiskLevelValue {
			riskLevelValueArray = append(riskLevelValueArray, &tiDTOs.ThreatIntelligenceRiskLevelDataDTO{
				ThreatIntelligenceRiskLevelID:      riskLevel.ID,
				ThreatIntelligenceRiskLevelValueID: value.ID,
				RiskLevelTitleEN:                   riskLevel.TitleEn,
				RiskLevelTitleAR:                   riskLevel.TitleAr,
				RiskLevelType:                      riskLevel.Type,
				ColorCode:                          riskLevel.ColorCode,
				Value:                              value.Value,
			})
		}
	}

	calculationMethod := "multiplication"
	var response [][]*tiDTOs.ThreatIntelligenceRiskLevelMatrixResponseDTO

	for _, likelihood := range likelihoods {
		if likelihood.Status != "active" {
			continue
		}

		var severityDetails []*tiDTOs.ThreatIntelligenceRiskLevelMatrixResponseDTO
		for _, severity := range severities {
			if severity.Status != "active" {
				continue
			}

			score := tirlvs.EvaluateScore(&models.ThreatIntelligenceLikelihood{Score: likelihood.Score}, &models.ThreatIntelligenceSeverity{Score: severity.Score}, calculationMethod)
			riskLevel, _ := tirlvs.FindRiskLevelValue(riskLevelValueArray, score)

			if riskLevel == nil {
				riskLevel, _ = tirlvs.FindClosestRiskLevelValue(riskLevelValueArray, score)
			}

			if riskLevel == nil {
				riskLevel = &tiDTOs.ThreatIntelligenceRiskLevelDataDTO{}
			}

			detail := &tiDTOs.ThreatIntelligenceRiskLevelMatrixResponseDTO{
				ThreatIntelligenceLikelihoodID:      likelihood.ID,
				ThreatIntelligenceLikelihoodTitleEN: likelihood.TitleEn,
				ThreatIntelligenceLikelihoodTitleAR: likelihood.TitleAr,
				ThreatIntelligenceSeverityID:        severity.ID,
				ThreatIntelligenceSeverityTitleEN:   severity.TitleEn,
				ThreatIntelligenceSeverityTitleAR:   severity.TitleAr,
				ThreatIntelligenceScore:             score,
				ThreatIntelligenceRiskLevel:         *riskLevel,
			}

			severityDetails = append(severityDetails, detail)
		}

		response = append(response, severityDetails)
	}

	return response, nil
}

func (tirlvs *ThreatIntelligenceRiskLevelValueService) EvaluateScore(likelihood *models.ThreatIntelligenceLikelihood, severity *models.ThreatIntelligenceSeverity, calculationMethod string) uint {
	var score uint
	if calculationMethod == "addition" {
		score = likelihood.Score + severity.Score
	} else if calculationMethod == "multiplication" {
		score = likelihood.Score * severity.Score
	}

	return score
}

func (tirlvs *ThreatIntelligenceRiskLevelValueService) FindRiskLevelValue(riskLevelValues []*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, score uint) (*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, bool) {
	for _, riskLevelValue := range riskLevelValues {
		if riskLevelValue.Value == score {
			return riskLevelValue, true
		}
	}
	return nil, false
}

// FindClosestRiskLevelValue finds the closest risk level value to the given score
func (tirlvs *ThreatIntelligenceRiskLevelValueService) FindClosestRiskLevelValue(riskLevelValues []*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, score uint) (*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, bool) {
	if len(riskLevelValues) == 0 {
		return nil, false
	}

	for _, riskLevelValue := range riskLevelValues {
		if riskLevelValue.Value == score {
			return riskLevelValue, true
		}
	}

	var closestValue *tiDTOs.ThreatIntelligenceRiskLevelDataDTO
	var minDifference uint
	firstIteration := true

	for _, riskLevelValue := range riskLevelValues {
		var difference uint
		if riskLevelValue.Value > score {
			difference = riskLevelValue.Value - score
		} else {
			difference = score - riskLevelValue.Value
		}

		if firstIteration || difference < minDifference {
			minDifference = difference
			closestValue = riskLevelValue
			firstIteration = false
		}
	}

	return closestValue, closestValue != nil
}
