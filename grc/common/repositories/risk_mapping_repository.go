package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type RiskMappingRepositoryInterface interface {
	IsRiskMapped(tableName string, masterColumn string, riskID, masterID uint) (bool, error)
	RiskExists(id uint) (*models.RiskRegister, error)
}

type RiskMappingRepository struct {
	db *gorm.DB
}

func NewRiskMappingRepository(db *gorm.DB) RiskMappingRepositoryInterface {
	return &RiskMappingRepository{
		db: db,
	}
}

// Generic function to check if a risk is already mapped to a master entity
func (r *RiskMappingRepository) IsRiskMapped(masterName string, columnName string, riskID, masterID uint) (bool, error) {
	var count int64
	err := r.db.
		Table(masterName).
		Where("risk_register_id = ? AND "+columnName+" = ?", riskID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *RiskMappingRepository) RiskExists(id uint) (*models.RiskRegister, error) {
	db := repo.db

	var risk *models.RiskRegister

	result := db.
		First(&risk, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return risk, nil

}
