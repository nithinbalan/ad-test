package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpBiaMatrixRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessContinuityPlan, error)
	GetByID(bcpProcessId uint) (*models.BcpProcess, error)
	Update(bcpProcess *models.BcpProcess, tx ...utility.DatabaseTransaction) error
}

type BcpBiaMatrixRepository struct {
	db *gorm.DB
}

func NewBcpBiaMatrixRepository(db *gorm.DB) *BcpBiaMatrixRepository {
	return &BcpBiaMatrixRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BcpBiaMatrixRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BcpBiaMatrixRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessContinuityPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var businessContinuityPlan *models.BusinessContinuityPlan
	result := db.
		Preload("RtoUnit").
		Preload("MtpdUnit").
		Preload("RpoUnit").
		Preload("Processes").
		Preload("Processes.Process").
		Preload("Processes.RtoUnit").
		Preload("Processes.MtpdUnit").
		Preload("Processes.RpoUnit").
		Preload("Processes.RiskLevel").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&businessContinuityPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return businessContinuityPlan, nil
}

func (repo *BcpBiaMatrixRepository) GetByID(id uint) (*models.BcpProcess, error) {
	var bcpProcess *models.BcpProcess

	result := repo.db.First(&bcpProcess, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return bcpProcess, nil

}

func (repo *BcpBiaMatrixRepository) Update(bcpProcess *models.BcpProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(bcpProcess).Error; err != nil {
		return err
	}

	return nil
}
