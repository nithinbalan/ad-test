package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskMitigationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	UpdateRiskResponseStrategy(riskID uint, responseStrategyID *uint, tx ...utility.DatabaseTransaction) error
	GetRiskByID(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
}

type RiskUpdatePlanRiskMitigationRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanRiskMitigationRepository(db *gorm.DB) *RiskUpdatePlanRiskMitigationRepository {
	return &RiskUpdatePlanRiskMitigationRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanRiskMitigationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanRiskMitigationRepository) UpdateRiskResponseStrategy(riskID uint, responseStrategyID *uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskRegister{}).
		Where("id = ?", riskID).
		Update("risk_response_strategy_id", responseStrategyID).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskUpdatePlanRiskMitigationRepository) GetRiskByID(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var risk models.RiskRegister
	if err := db.
		First(&risk, riskID).Error; err != nil {
		return nil, err
	}

	return &risk, nil
}
