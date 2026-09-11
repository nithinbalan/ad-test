package repositories

import (
	"jethur/grc/models"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThreatIntelligenceRiskLevelValueRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll() ([]*tiDTOs.ThreatIntelligenceRiskLevelValueListResponseDTO, error)
	Create(values []*models.ThreatIntelligenceRiskLevelValue, riskLevelIDs []uint, tx ...utility.DatabaseTransaction) error
}

type ThreatIntelligenceRiskLevelValueRepository struct {
	db *gorm.DB
}

func NewThreatIntelligenceRiskLevelValueRepository(db *gorm.DB) *ThreatIntelligenceRiskLevelValueRepository {
	return &ThreatIntelligenceRiskLevelValueRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ThreatIntelligenceRiskLevelValueRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThreatIntelligenceRiskLevelValueRepository) GetAll() ([]*tiDTOs.ThreatIntelligenceRiskLevelValueListResponseDTO, error) {
	var values []*tiDTOs.ThreatIntelligenceRiskLevelValueListResponseDTO

	result := repo.db.Table("threat_intelligence_risk_level_values").
		Select("threat_intelligence_risk_level_values.*").
		Find(&values)

	return values, result.Error
}

func (repo *ThreatIntelligenceRiskLevelValueRepository) Create(values []*models.ThreatIntelligenceRiskLevelValue, riskLevelIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(riskLevelIDs) > 0 {
		result := db.Where("threat_intelligence_risk_level_values.threat_intelligence_risk_level_id IN (?)", riskLevelIDs).
			Delete(&models.ThreatIntelligenceRiskLevelValue{})
		if result.Error != nil {
			return result.Error
		}
	}

	if len(values) > 0 {
		if err := db.Create(&values).Error; err != nil {
			return err
		}
	}
	return nil
}
