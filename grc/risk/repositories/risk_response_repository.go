package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskResponseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	GetFullDetailsByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	GetByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	Create(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error
	Update(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error
	DeleteByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) error
}

type RiskResponseRepository struct {
	db *gorm.DB
}

func NewRiskResponseRepository(db *gorm.DB) *RiskResponseRepository {
	return &RiskResponseRepository{
		db: db,
	}
}

func (repo *RiskResponseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskResponseRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	err := db.Where("id = ?", id).First(&riskResponse).Error
	if err != nil {
		return nil, err
	}
	return &riskResponse, nil
}

func (repo *RiskResponseRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	err := db.Where("id = ?", id).
		Preload("RiskRegister").
		Preload("RiskRegister.ResidualRiskRating").
		Preload("RiskRegister.ResidualRiskRatingValue").
		Preload("RiskRegister.ActionPlans").
		Preload("WorkflowStatus").
		Preload("SubmittedBy").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Preload("UpdatedBy").
		First(&riskResponse).Error
	if err != nil {
		return nil, err
	}
	return &riskResponse, nil
}

func (repo *RiskResponseRepository) GetFullDetailsByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	err := db.Where("risk_register_id = ?", riskRegisterID).
		Preload("RiskRegister").
		Preload("RiskRegister.ResidualRiskRating").
		Preload("RiskRegister.ResidualRiskRatingValue").
		Preload("WorkflowStatus").
		Preload("SubmittedBy").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Preload("UpdatedBy").
		First(&riskResponse).Error
	if err != nil {
		return nil, err
	}
	return &riskResponse, nil
}

func (repo *RiskResponseRepository) GetByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	err := db.Where("risk_register_id = ?", riskRegisterID).First(&riskResponse).Error
	if err != nil {
		return nil, err
	}
	return &riskResponse, nil
}

func (repo *RiskResponseRepository) Create(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(riskResponse).Error
}

func (repo *RiskResponseRepository) Update(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(riskResponse).Error
}

func (repo *RiskResponseRepository) DeleteByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_register_id = ?", riskRegisterID).Delete(&models.RiskResponse{}).Error
}
