package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type StrategicObjectiveKPIRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateStrategicObjectiveKPIs(mappings []models.StrategicObjectiveKPI, tx utility.DatabaseTransaction) error
	DeleteSingleStrategicObjectiveKPI(strategicObjectiveID uint, kpiRegisterID uint, tx utility.DatabaseTransaction) error
}

// implementation
type StrategicObjectiveKPIRepository struct {
	db *gorm.DB
}

func NewStrategicObjectiveKPIRepository(db *gorm.DB) *StrategicObjectiveKPIRepository {
	return &StrategicObjectiveKPIRepository{
		db: db,
	}
}
func (repo *StrategicObjectiveKPIRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategicObjectiveKPIRepository) CreateStrategicObjectiveKPIs(mappings []models.StrategicObjectiveKPI, tx utility.DatabaseTransaction) error {
	db := tx.DB()
	return db.Create(&mappings).Error
}
func (repo *StrategicObjectiveKPIRepository) DeleteSingleStrategicObjectiveKPI(strategicObjectiveID uint, kpiRegisterID uint, tx utility.DatabaseTransaction) error {
	db := tx.DB()
	return db.
		Where("strategic_objective_id = ? AND kpi_register_id = ?", strategicObjectiveID, kpiRegisterID).
		Delete(&models.StrategicObjectiveKPI{}).Error
}
