package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskResponseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	GetByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	GetFullDetailsByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	Create(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AssetRiskResponseRepository struct {
	db *gorm.DB
}

func NewAssetRiskResponseRepository(db *gorm.DB) AssetRiskResponseRepositoryInterface {
	return &AssetRiskResponseRepository{
		db: db,
	}
}

func (repo *AssetRiskResponseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskResponseRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	if err := db.Where("id = ?", id).First(&riskResponse).Error; err != nil {
		return nil, err
	}
	return &riskResponse, nil
}

func (repo *AssetRiskResponseRepository) GetByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	if err := db.Where("risk_register_id = ?", riskRegisterID).First(&riskResponse).Error; err != nil {
		return nil, err
	}
	return &riskResponse, nil
}

func (repo *AssetRiskResponseRepository) GetFullDetailsByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	err := db.Where("risk_register_id = ?", riskRegisterID).
		Preload("RiskRegister").
		Preload("RiskRegister.RiskResponseStrategy").
		Preload("CreatedBy").
		Preload("UpdatedBy").
		First(&riskResponse).Error
	if err != nil {
		return nil, err
	}
	return &riskResponse, nil
}

func (repo *AssetRiskResponseRepository) Create(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(riskResponse).Error; err != nil {
		return 0, err
	}
	return riskResponse.ID, nil
}

func (repo *AssetRiskResponseRepository) Update(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(riskResponse).Error
}

func (repo *AssetRiskResponseRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.RiskResponse{}, id).Error
}
