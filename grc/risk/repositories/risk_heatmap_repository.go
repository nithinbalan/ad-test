package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"

	"gorm.io/gorm"
)

type RiskHeatmapRepositoryInterface interface {
	GetRiskProcessesByProcessID(processID uint) ([]models.RiskProcess, error)
	GetRiskControlsByControlID(controlID uint) ([]models.RiskControl, error)
	GetRiskCompliancesByComplianceID(complianceID uint) ([]models.RiskComplianceRequirement, error)
	GetRiskSLAsBySLAID(slaID uint) ([]models.RiskSLA, error)
	GetRiskAssetByAssetID(assetID uint) ([]models.RiskAsset, error)
	GetAllRiskRegisters(params *riskDTOs.RiskHeatmapQueryParamsDTO) ([]models.RiskRegister, error)

	CountRiskRegistersWithLikelihoodAndImpact(riskRegisterIDs []uint, riskRatingValueID, likelihoodID, impactID uint) (int, error)

	GetVulnerabilityIDsByAssetID(assetID uint) ([]uint, error)
	GetThreatIDsByVulnerabilityIDs(vulnIDs []uint) ([]uint, error)
	GetRiskRegisterIDsByThreatIDs(threatIDs []uint) ([]uint, error)
}

type RiskHeatmapRepository struct {
	db *gorm.DB
}

func NewRiskHeatmapRepository(db *gorm.DB) *RiskHeatmapRepository {
	return &RiskHeatmapRepository{
		db: db,
	}
}

func (repo *RiskHeatmapRepository) GetRiskProcessesByProcessID(processID uint) ([]models.RiskProcess, error) {
	var riskProcesses []models.RiskProcess
	if err := repo.db.Where("process_id = ?", processID).Find(&riskProcesses).Error; err != nil {
		return nil, err
	}
	return riskProcesses, nil
}

func (repo *RiskHeatmapRepository) GetRiskControlsByControlID(controlID uint) ([]models.RiskControl, error) {
	var riskControls []models.RiskControl
	if err := repo.db.Where("control_id = ?", controlID).Find(&riskControls).Error; err != nil {
		return nil, err
	}
	return riskControls, nil
}

func (repo *RiskHeatmapRepository) GetRiskCompliancesByComplianceID(complianceID uint) ([]models.RiskComplianceRequirement, error) {
	var riskCompliances []models.RiskComplianceRequirement
	if err := repo.db.Where("compliance_register_id = ?", complianceID).Find(&riskCompliances).Error; err != nil {
		return nil, err
	}
	return riskCompliances, nil
}

func (repo *RiskHeatmapRepository) GetRiskSLAsBySLAID(slaID uint) ([]models.RiskSLA, error) {
	var riskSLAs []models.RiskSLA
	if err := repo.db.Where("sla_register_id = ?", slaID).Find(&riskSLAs).Error; err != nil {
		return nil, err
	}
	return riskSLAs, nil
}

func (repo *RiskHeatmapRepository) GetRiskAssetByAssetID(assetID uint) ([]models.RiskAsset, error) {
	var riskAssets []models.RiskAsset
	if err := repo.db.Where("asset_register_id = ?", assetID).Find(&riskAssets).Error; err != nil {
		return nil, err
	}
	return riskAssets, nil
}

func (repo *RiskHeatmapRepository) CountRiskRegistersWithLikelihoodAndImpact(riskRegisterIDs []uint, riskRatingValueID, likelihoodID, impactID uint) (int, error) {
	var count int64
	err := repo.db.Model(&models.RiskRegister{}).
		Where("id IN (?) AND current_risk_rating_value_id = ? AND current_likelihood_id = ? AND current_impact_id = ?", riskRegisterIDs, riskRatingValueID, likelihoodID, impactID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *RiskHeatmapRepository) GetAllRiskRegisters(params *riskDTOs.RiskHeatmapQueryParamsDTO) ([]models.RiskRegister, error) {
	var risks []models.RiskRegister

	query := repo.db.Model(&models.RiskRegister{}).
		Preload("ResponsibleUnit").
		Preload("RiskCategory").
		Preload("RiskSource").
		Preload("CurrentRiskRating").
		Preload("CurrentRiskRatingValue").
		Preload("RiskResponseStrategy").
		Preload("RiskStatus").
		Preload("ReviewFrequency")

	// Apply date filter if both DateFrom and DateTo are provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("reported_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err := query.Find(&risks).Error
	if err != nil {
		return nil, err
	}

	return risks, nil
}

func (repo *RiskHeatmapRepository) GetVulnerabilityIDsByAssetID(assetID uint) ([]uint, error) {
	var ids []uint
	err := repo.db.
		Model(&models.VulnerabilityAsset{}).
		Where("asset_register_id = ?", assetID).
		Pluck("vulnerability_id", &ids).Error
	return ids, err
}

func (repo *RiskHeatmapRepository) GetThreatIDsByVulnerabilityIDs(vulnIDs []uint) ([]uint, error) {
	var ids []uint
	err := repo.db.
		Model(&models.VulnerabilityThreat{}).
		Where("vulnerability_id IN ?", vulnIDs).
		Pluck("threat_id", &ids).Error
	return ids, err
}

func (repo *RiskHeatmapRepository) GetRiskRegisterIDsByThreatIDs(threatIDs []uint) ([]uint, error) {
	var ids []uint
	err := repo.db.
		Model(&models.RiskThreat{}).
		Where("threat_id IN ?", threatIDs).
		Pluck("risk_register_id", &ids).Error
	return ids, err
}
