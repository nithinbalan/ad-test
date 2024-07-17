package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type RiskHeatmapRepositoryInterface interface {
	GetRiskProcessesByProcessID(processID uint) ([]models.RiskProcess, error)

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

func (repo *RiskHeatmapRepository) CountRiskRegistersWithLikelihoodAndImpact(riskRegisterIDs []uint, riskRatingValueID, likelihoodID, impactID uint) (int, error) {
	var count int64
	err := repo.db.Model(&models.RiskRegister{}).
		Where("id IN (?) AND current_risk_rating_value_id = ? AND risk_matrix_likelihood_id = ? AND risk_matrix_impact_id = ?", riskRegisterIDs, riskRatingValueID, likelihoodID, impactID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
