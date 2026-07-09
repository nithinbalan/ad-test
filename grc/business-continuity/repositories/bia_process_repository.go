package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BiaProcessRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error)
	GetFullDetailsByProcessID(processID uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error)
	GetByProcessID(processID uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error)
	Create(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error
	Update(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error
	DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
}

type BiaProcessRepository struct {
	db *gorm.DB
}

func NewBiaProcessRepository(db *gorm.DB) *BiaProcessRepository {
	return &BiaProcessRepository{
		db: db,
	}
}

func (repo *BiaProcessRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BiaProcessRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var biaProcess models.BiaProcess
	err := db.Where("id = ?", id).First(&biaProcess).Error
	if err != nil {
		return nil, err
	}
	return &biaProcess, nil
}

func (repo *BiaProcessRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var biaProcess models.BiaProcess
	err := db.Where("id = ?", id).
		Preload("Process").
		Preload("WorkflowStatus").
		Preload("SubmittedBy").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Preload("UpdatedBy").
		First(&biaProcess).Error
	if err != nil {
		return nil, err
	}
	return &biaProcess, nil
}

func (repo *BiaProcessRepository) GetFullDetailsByProcessID(processID uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var biaProcess models.BiaProcess
	err := db.Where("process_id = ?", processID).
		Preload("Process").
		Preload("WorkflowStatus").
		Preload("SubmittedBy").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Preload("UpdatedBy").
		First(&biaProcess).Error
	if err != nil {
		return nil, err
	}
	return &biaProcess, nil
}

func (repo *BiaProcessRepository) GetByProcessID(processID uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var biaProcess models.BiaProcess
	err := db.Where("process_id = ?", processID).First(&biaProcess).Error
	if err != nil {
		return nil, err
	}
	return &biaProcess, nil
}

func (repo *BiaProcessRepository) Create(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(biaProcess).Error
}

func (repo *BiaProcessRepository) Update(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(biaProcess).Error
}

func (repo *BiaProcessRepository) DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("process_id = ?", processID).Delete(&models.BiaProcess{}).Error
}
