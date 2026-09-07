package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestExecutionRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(penetrationTestID uint) (*models.PenetrationTestExecution, error)
	GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestExecution, error)
	Create(execution *models.PenetrationTestExecution, tx ...utility.DatabaseTransaction) (uint, error)
	Update(execution *models.PenetrationTestExecution, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type PenetrationTestExecutionRepository struct {
	db *gorm.DB
}

func NewPenetrationTestExecutionRepository(db *gorm.DB) *PenetrationTestExecutionRepository {
	return &PenetrationTestExecutionRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *PenetrationTestExecutionRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestExecutionRepository) GetByID(penetrationTestID uint) (*models.PenetrationTestExecution, error) {
	var execution *models.PenetrationTestExecution
	result := repo.db.Where("penetration_test_id = ?", penetrationTestID).First(&execution)

	if result.Error != nil {
		return nil, result.Error
	}

	return execution, nil
}

func (repo *PenetrationTestExecutionRepository) GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestExecution, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var execution *models.PenetrationTestExecution

	result := db.Where("penetration_test_id = ?", penetrationTestID).
		Preload("ExecutionStatus").
		Preload("SecurityResult").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&execution)

	if result.Error != nil {
		return nil, result.Error
	}

	return execution, nil
}

func (repo *PenetrationTestExecutionRepository) Create(execution *models.PenetrationTestExecution, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(execution).Error; err != nil {
		return 0, err
	}
	return execution.ID, nil
}

func (repo *PenetrationTestExecutionRepository) Update(execution *models.PenetrationTestExecution, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(execution).Error; err != nil {
		return err
	}
	return nil
}

func (repo *PenetrationTestExecutionRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var execution *models.PenetrationTestExecution
	if err := db.First(&execution, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.PenetrationTestExecution{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
