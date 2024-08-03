package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type RiskHeatmapRepositoryInterface interface {
	GetRiskProcessesByProcessID(processID uint) ([]models.RiskProcess, error)
	GetRiskControlsByControlID(controlID uint) ([]models.RiskControl, error)
	GetRiskCompliancesByComplianceID(complianceID uint) ([]models.RiskComplianceRequirement, error)
	GetAllRiskRegisters() ([]models.RiskRegister, error)

	CountRiskRegistersWithLikelihoodAndImpact(riskRegisterIDs []uint, riskRatingValueID, likelihoodID, impactID uint) (int, error)
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

func (repo *RiskHeatmapRepository) GetAllRiskRegisters() ([]models.RiskRegister, error) {
	var risks []models.RiskRegister
	err := repo.db.Preload("ResponsibleUnit").Find(&risks).Error
	if err != nil {
		return nil, err
	}
	return risks, nil
}
